package release

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureArtifactDownloadsAndReusesVersion(t *testing.T) {
	t.Setenv("RUNNER_HOME", t.TempDir())
	jarBody := []byte("jar simulado")
	sum := sha256.Sum256(jarBody)
	downloads := 0
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/release.json":
			fmt.Fprintf(w, `{"simulador":{"url":"%s/artifact.jar","version":"0.1.0","sha256":"%x"},"jre":{"linux_x64":"%s/jre.tar.gz"}}`, srv.URL, sum, srv.URL)
		case "/artifact.jar":
			downloads++
			_, _ = w.Write(jarBody)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	first, artifact, reused, err := EnsureArtifact(context.Background(), srv.URL+"/release.json", "simulador")
	if err != nil {
		t.Fatalf("EnsureArtifact retornou erro: %v", err)
	}
	if reused {
		t.Fatalf("primeiro download não deve ser marcado como reuso")
	}
	if artifact.Version != "0.1.0" {
		t.Fatalf("versão esperada 0.1.0, obtida %s", artifact.Version)
	}
	if _, err := os.Stat(first); err != nil {
		t.Fatalf("arquivo baixado não existe: %v", err)
	}

	second, _, reused, err := EnsureArtifact(context.Background(), srv.URL+"/release.json", "simulador")
	if err != nil {
		t.Fatalf("EnsureArtifact no reuso retornou erro: %v", err)
	}
	if !reused {
		t.Fatalf("segundo uso deveria reaproveitar artefato local")
	}
	if first != second {
		t.Fatalf("caminho reaproveitado deveria ser igual")
	}
	if downloads != 1 {
		t.Fatalf("esperava 1 download, obtido %d", downloads)
	}
}

func TestPlatformKeyReturnsNonEmptyValue(t *testing.T) {
	if PlatformKey() == "" {
		t.Fatalf("PlatformKey não deveria ser vazio")
	}
}

func TestArtifactFallbackToValidador(t *testing.T) {
	m := Manifest{Validador: Artifact{URL: "https://example.org/validador.jar", Version: "1.0.0"}}
	a, err := m.Artifact("simulador")
	if err != nil {
		t.Fatalf("fallback simulador->validador deveria funcionar: %v", err)
	}
	if filepath.Base(a.URL) != "validador.jar" {
		t.Fatalf("artefato inesperado: %s", a.URL)
	}
}

func TestEnsureArtifactFromURLDownloadsAndReusesSource(t *testing.T) {
	t.Setenv("RUNNER_HOME", t.TempDir())
	body := []byte("simulador jar via source")
	sum := sha256.Sum256(body)
	downloads := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/simulador.jar" {
			http.NotFound(w, r)
			return
		}
		downloads++
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	first, artifact, reused, err := EnsureArtifactFromURL(context.Background(), srv.URL+"/simulador.jar", fmt.Sprintf("%x", sum), "simulador")
	if err != nil {
		t.Fatalf("EnsureArtifactFromURL retornou erro: %v", err)
	}
	if reused {
		t.Fatalf("primeiro uso não deve ser reuso")
	}
	if artifact.URL == "" || artifact.SHA256 == "" {
		t.Fatalf("artefato direto deveria preservar URL e SHA: %#v", artifact)
	}
	second, _, reused, err := EnsureArtifactFromURL(context.Background(), srv.URL+"/simulador.jar", fmt.Sprintf("%x", sum), "simulador")
	if err != nil {
		t.Fatalf("EnsureArtifactFromURL no reuso retornou erro: %v", err)
	}
	if !reused || first != second || downloads != 1 {
		t.Fatalf("cache inesperado: first=%s second=%s reused=%v downloads=%d", first, second, reused, downloads)
	}
}
