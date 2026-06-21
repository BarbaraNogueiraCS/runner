package process

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSaveLoadAndDeleteState(t *testing.T) {
	runnerHome := t.TempDir()
	t.Setenv("RUNNER_HOME", runnerHome)

	state := State{
		Name:      "assinador",
		PID:       12345,
		Port:      8080,
		URL:       "http://127.0.0.1:8080",
		JarPath:   "assinador/target/assinador.jar",
		StartedAt: time.Date(2026, 6, 8, 10, 0, 0, 0, time.UTC),
	}

	if err := Save(state); err != nil {
		t.Fatalf("Save retornou erro: %v", err)
	}

	loaded, err := Load("assinador")
	if err != nil {
		t.Fatalf("Load retornou erro: %v", err)
	}

	if loaded.Name != state.Name || loaded.PID != state.PID || loaded.Port != state.Port || loaded.URL != state.URL || loaded.JarPath != state.JarPath {
		t.Fatalf("estado carregado diferente do salvo: %#v", loaded)
	}
	if !loaded.StartedAt.Equal(state.StartedAt) {
		t.Fatalf("StartedAt esperado %v, obtido %v", state.StartedAt, loaded.StartedAt)
	}

	if err := Delete("assinador"); err != nil {
		t.Fatalf("Delete retornou erro: %v", err)
	}

	stateFile := filepath.Join(runnerHome, "processos", "assinador.json")
	if _, err := os.Stat(stateFile); !os.IsNotExist(err) {
		t.Fatalf("arquivo de estado deveria ter sido removido; erro obtido: %v", err)
	}
}
