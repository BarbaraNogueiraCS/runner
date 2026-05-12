package runtime

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/kyriosdata/assinatura/internal/config"
)

type JavaRuntime struct {
	Version  string
	JavaPath string
	Source   string
}

type Resolver struct{}

func NewResolver() Resolver { return Resolver{} }

func (r Resolver) Resolve() (JavaRuntime, error) {
	if err := config.EnsureBaseDirs(); err != nil {
		return JavaRuntime{}, err
	}
	candidates := []struct {
		path   string
		source string
	}{}
	if env := os.Getenv("JAVA_HOME"); env != "" {
		candidates = append(candidates, struct{ path, source string }{filepath.Join(env, "bin", javaExecutable()), "JAVA_HOME"})
	}
	candidates = append(candidates, struct{ path, source string }{filepath.Join(config.HubSaudeHome(), "jdk", "java-21", "bin", javaExecutable()), "gerenciado"})
	candidates = append(candidates, struct{ path, source string }{filepath.Join(config.HubSaudeHome(), "jre", "java-21", "bin", javaExecutable()), "gerenciado"})
	if path, err := exec.LookPath("java"); err == nil {
		candidates = append(candidates, struct{ path, source string }{path, "PATH"})
	}
	for _, candidate := range candidates {
		if candidate.path == "" {
			continue
		}
		if _, err := os.Stat(candidate.path); err != nil {
			continue
		}
		version, _ := javaVersion(candidate.path)
		if version == "" || strings.Contains(version, "21") || strings.Contains(version, "22") || strings.Contains(version, "23") || strings.Contains(version, "24") || strings.Contains(version, "25") {
			return JavaRuntime{Version: version, JavaPath: candidate.path, Source: candidate.source}, nil
		}
	}
	return JavaRuntime{}, fmt.Errorf("Java 21 compatível não encontrado. Instale Java 21 ou configure JAVA_HOME")
}

func javaVersion(javaPath string) (string, error) {
	cmd := exec.Command(javaPath, "-version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func javaExecutable() string {
	if runtime.GOOS == "windows" {
		return "java.exe"
	}
	return "java"
}
