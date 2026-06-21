package jdk

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"time"

	"github.com/BarbaraNogueiraCS/runner/internal/release"
)

// FindJava locates a Java executable. It first honors RUNNER_JAVA, then tries
// the managed ~/.hubsaude/jdk directory, and finally java from PATH.
func FindJava() (string, error) {
	if v := os.Getenv("RUNNER_JAVA"); v != "" {
		if _, err := os.Stat(v); err == nil {
			return v, nil
		}
		return "", fmt.Errorf("RUNNER_JAVA aponta para um arquivo inexistente: %s", v)
	}
	if managed := managedJavaPath(); managed != "" {
		if _, err := os.Stat(managed); err == nil {
			return managed, nil
		}
	}
	if p, err := exec.LookPath("java"); err == nil {
		return p, nil
	}
	return "", errors.New("Java não encontrado. Instale Java 21+ ou configure RUNNER_JAVA")
}

func EnsureJava21(ctx context.Context, manifestURL string) (string, error) {
	if java, err := FindJava(); err == nil {
		major, versionText, verr := MajorVersion(java)
		if verr == nil && major >= 21 {
			return java, nil
		}
		if os.Getenv("RUNNER_JAVA") != "" {
			if verr != nil {
				return "", fmt.Errorf("não foi possível verificar versão do Java definido em RUNNER_JAVA: %w", verr)
			}
			return "", fmt.Errorf("RUNNER_JAVA aponta para Java %s; é necessário Java 21+", versionText)
		}
	}
	java, err := release.DownloadAndInstallJRE(ctx, manifestURL)
	if err != nil {
		return "", fmt.Errorf("Java 21+ não encontrado e não foi possível baixar runtime Temurin: %w", err)
	}
	major, versionText, verr := MajorVersion(java)
	if verr != nil {
		return "", fmt.Errorf("runtime Temurin baixado, mas a versão não pôde ser verificada: %w", verr)
	}
	if major < 21 {
		return "", fmt.Errorf("runtime Temurin baixado é Java %s; é necessário Java 21+", versionText)
	}
	return java, nil
}

func MajorVersion(javaPath string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, javaPath, "-version")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return 0, "", err
	}
	text := stdout.String() + stderr.String()
	re := regexp.MustCompile(`(?m)(?:java|openjdk) version "([0-9]+)(?:\.([0-9]+))?`)
	m := re.FindStringSubmatch(text)
	if len(m) < 2 {
		return 0, text, fmt.Errorf("saída de versão Java não reconhecida: %s", text)
	}
	major, err := strconv.Atoi(m[1])
	if err != nil {
		return 0, text, err
	}
	// Compatibilidade com versões antigas no formato 1.8.x.
	if major == 1 && len(m) >= 3 && m[2] != "" {
		major, _ = strconv.Atoi(m[2])
	}
	return major, m[1], nil
}

func managedJavaPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	name := "java"
	if runtime.GOOS == "windows" {
		name = "java.exe"
	}
	return filepath.Join(home, ".hubsaude", "jdk", "bin", name)
}
