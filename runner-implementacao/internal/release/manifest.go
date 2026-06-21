package release

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/BarbaraNogueiraCS/runner/internal/paths"
)

const DefaultManifestURL = "https://raw.githubusercontent.com/BarbaraNogueiraCS/runner/main/runner-implementacao/release.json"

type Artifact struct {
	URL     string `json:"url"`
	Version string `json:"version"`
	Tag     string `json:"tag,omitempty"`
	SHA256  string `json:"sha256,omitempty"`
}

type Manifest struct {
	Jar       Artifact          `json:"jar"`
	Validador Artifact          `json:"validador"`
	Simulador Artifact          `json:"simulador"`
	JRE       map[string]string `json:"jre"`
	JDK       map[string]string `json:"jdk"`
}

func FetchManifest(ctx context.Context, manifestURL string) (Manifest, error) {
	if manifestURL == "" {
		manifestURL = DefaultManifestURL
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, manifestURL, nil)
	if err != nil {
		return Manifest{}, err
	}
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return Manifest{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return Manifest{}, fmt.Errorf("release.json retornou HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}
	var m Manifest
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return Manifest{}, fmt.Errorf("release.json inválido: %w", err)
	}
	return m, nil
}

func (m Manifest) Artifact(name string) (Artifact, error) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "", "simulador":
		if m.Simulador.URL != "" {
			return m.Simulador, nil
		}
		if m.Validador.URL != "" {
			return m.Validador, nil
		}
	case "validador", "hubsaude-validador-api":
		if m.Validador.URL != "" {
			return m.Validador, nil
		}
	case "assinador", "jar":
		if m.Jar.URL != "" {
			return m.Jar, nil
		}
	}
	return Artifact{}, fmt.Errorf("artefato %q não encontrado no release.json", name)
}

func EnsureArtifact(ctx context.Context, manifestURL, artifactName string) (string, Artifact, bool, error) {
	m, err := FetchManifest(ctx, manifestURL)
	if err != nil {
		return "", Artifact{}, false, err
	}
	a, err := m.Artifact(artifactName)
	if err != nil {
		return "", Artifact{}, false, err
	}
	if a.URL == "" || a.Version == "" {
		return "", Artifact{}, false, errors.New("release.json sem url ou version para o artefato solicitado")
	}
	dir, err := managedArtifactDir(artifactName)
	if err != nil {
		return "", a, false, err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", a, false, err
	}
	jarPath := filepath.Join(dir, artifactFileName(artifactName, a))
	metaPath := jarPath + ".version"
	if fileMatches(jarPath, a.SHA256) && versionMatches(metaPath, a.Version) {
		return jarPath, a, true, nil
	}
	tmp := jarPath + ".download"
	if err := downloadToFile(ctx, a.URL, tmp); err != nil {
		_ = os.Remove(tmp)
		return "", a, false, err
	}
	if a.SHA256 != "" {
		ok, err := sha256Matches(tmp, a.SHA256)
		if err != nil {
			_ = os.Remove(tmp)
			return "", a, false, err
		}
		if !ok {
			_ = os.Remove(tmp)
			return "", a, false, fmt.Errorf("checksum SHA256 inválido para %s", a.URL)
		}
	}
	if err := os.Rename(tmp, jarPath); err != nil {
		_ = os.Remove(tmp)
		return "", a, false, err
	}
	if err := os.WriteFile(metaPath, []byte(a.Version+"\n"), 0o644); err != nil {
		return "", a, false, err
	}
	return jarPath, a, false, nil
}

func JREURL(m Manifest) (string, error) {
	key := PlatformKey()
	if m.JDK != nil && m.JDK[key] != "" {
		return m.JDK[key], nil
	}
	if m.JRE != nil && m.JRE[key] != "" {
		return m.JRE[key], nil
	}
	return "", fmt.Errorf("release.json não possui JRE/JDK para %s", key)
}

func PlatformKey() string {
	osName := runtime.GOOS
	arch := runtime.GOARCH
	switch osName {
	case "darwin":
		osName = "mac"
	}
	switch arch {
	case "amd64":
		arch = "x64"
	case "arm64":
		arch = "arm64"
	}
	return osName + "_" + arch
}

func DownloadAndInstallJRE(ctx context.Context, manifestURL string) (string, error) {
	m, err := FetchManifest(ctx, manifestURL)
	if err != nil {
		return "", err
	}
	url, err := JREURL(m)
	if err != nil {
		return "", err
	}
	home, err := paths.Home()
	if err != nil {
		return "", err
	}
	cache, err := paths.CacheDir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(cache, 0o755); err != nil {
		return "", err
	}
	ext := ".tar.gz"
	if runtime.GOOS == "windows" {
		ext = ".zip"
	}
	archivePath := filepath.Join(cache, "temurin-21-runtime"+ext)
	if err := downloadToFile(ctx, url, archivePath); err != nil {
		return "", err
	}
	target := filepath.Join(home, "jdk")
	tmp := target + ".tmp"
	_ = os.RemoveAll(tmp)
	if err := os.MkdirAll(tmp, 0o755); err != nil {
		return "", err
	}
	if runtime.GOOS == "windows" {
		err = extractZipStripRoot(archivePath, tmp)
	} else {
		err = extractTarGzStripRoot(archivePath, tmp)
	}
	if err != nil {
		_ = os.RemoveAll(tmp)
		return "", err
	}
	_ = os.RemoveAll(target)
	if err := os.Rename(tmp, target); err != nil {
		_ = os.RemoveAll(tmp)
		return "", err
	}
	java := filepath.Join(target, "bin", "java")
	if runtime.GOOS == "windows" {
		java += ".exe"
	}
	return java, nil
}

func managedArtifactDir(name string) (string, error) {
	home, err := paths.Home()
	if err != nil {
		return "", err
	}
	safe := strings.ToLower(strings.TrimSpace(name))
	if safe == "" {
		safe = "simulador"
	}
	safe = strings.NewReplacer("/", "-", "\\", "-", " ", "-").Replace(safe)
	return filepath.Join(home, safe), nil
}

func artifactFileName(name string, a Artifact) string {
	base := filepath.Base(a.URL)
	if strings.Contains(base, ".") && strings.HasSuffix(base, ".jar") {
		return base
	}
	safe := strings.ToLower(strings.TrimSpace(name))
	if safe == "" {
		safe = "simulador"
	}
	return safe + "-" + a.Version + ".jar"
}

func versionMatches(path, expected string) bool {
	b, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(b)) == expected
}

func fileMatches(path, expectedSHA string) bool {
	st, err := os.Stat(path)
	if err != nil || st.IsDir() || st.Size() == 0 {
		return false
	}
	if expectedSHA == "" {
		return true
	}
	ok, err := sha256Matches(path, expectedSHA)
	return err == nil && ok
}

func sha256Matches(path, expected string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return false, err
	}
	got := hex.EncodeToString(h.Sum(nil))
	return strings.EqualFold(got, strings.TrimSpace(expected)), nil
}

func downloadToFile(ctx context.Context, url, dst string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: 10 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("download retornou HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := io.Copy(f, resp.Body); err != nil {
		return err
	}
	return nil
}

func extractTarGzStripRoot(archivePath, dst string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		rel := stripFirstPathComponent(hdr.Name)
		if rel == "" {
			continue
		}
		target := filepath.Join(dst, rel)
		if !strings.HasPrefix(filepath.Clean(target), filepath.Clean(dst)+string(os.PathSeparator)) && filepath.Clean(target) != filepath.Clean(dst) {
			return fmt.Errorf("entrada insegura no arquivo: %s", hdr.Name)
		}
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(hdr.Mode)&0o777)
			if err != nil {
				return err
			}
			_, copyErr := io.Copy(out, tr)
			closeErr := out.Close()
			if copyErr != nil {
				return copyErr
			}
			if closeErr != nil {
				return closeErr
			}
		}
	}
}

func extractZipStripRoot(archivePath, dst string) error {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, zf := range r.File {
		rel := stripFirstPathComponent(zf.Name)
		if rel == "" {
			continue
		}
		target := filepath.Join(dst, rel)
		if !strings.HasPrefix(filepath.Clean(target), filepath.Clean(dst)+string(os.PathSeparator)) && filepath.Clean(target) != filepath.Clean(dst) {
			return fmt.Errorf("entrada insegura no arquivo: %s", zf.Name)
		}
		if zf.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		in, err := zf.Open()
		if err != nil {
			return err
		}
		out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, zf.Mode())
		if err != nil {
			_ = in.Close()
			return err
		}
		_, copyErr := io.Copy(out, in)
		closeInErr := in.Close()
		closeOutErr := out.Close()
		if copyErr != nil {
			return copyErr
		}
		if closeInErr != nil {
			return closeInErr
		}
		if closeOutErr != nil {
			return closeOutErr
		}
	}
	return nil
}

func stripFirstPathComponent(p string) string {
	p = filepath.ToSlash(filepath.Clean(p))
	p = strings.TrimPrefix(p, "./")
	if p == "." || p == "/" || p == "" {
		return ""
	}
	parts := strings.Split(p, "/")
	if len(parts) <= 1 {
		return ""
	}
	return filepath.Join(parts[1:]...)
}
