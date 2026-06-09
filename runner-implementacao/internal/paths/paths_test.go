package paths

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHomeUsesRunnerHomeWhenConfigured(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("RUNNER_HOME", tempDir)

	home, err := Home()
	if err != nil {
		t.Fatalf("Home retornou erro: %v", err)
	}
	if home != tempDir {
		t.Fatalf("Home deveria usar RUNNER_HOME=%s, obtido %s", tempDir, home)
	}
}

func TestEnsureRuntimeDirsCreatesManagedDirectories(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("RUNNER_HOME", tempDir)

	if err := EnsureRuntimeDirs(); err != nil {
		t.Fatalf("EnsureRuntimeDirs retornou erro: %v", err)
	}

	expectedDirs := []string{
		tempDir,
		filepath.Join(tempDir, "processos"),
		filepath.Join(tempDir, "logs"),
		filepath.Join(tempDir, "cache"),
	}

	for _, dir := range expectedDirs {
		info, err := os.Stat(dir)
		if err != nil {
			t.Fatalf("diretório esperado não foi criado %s: %v", dir, err)
		}
		if !info.IsDir() {
			t.Fatalf("caminho %s deveria ser diretório", dir)
		}
	}
}
