package config

import (
	"os"
	"path/filepath"
)

const (
	DefaultAssinadorPort = 8080
	DefaultSimulatorPort = 8443
)

func HubSaudeHome() string {
	if custom := os.Getenv("HUBSAUDE_HOME"); custom != "" {
		return custom
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ".hubsaude"
	}
	return filepath.Join(home, ".hubsaude")
}

func EnsureBaseDirs() error {
	dirs := []string{
		HubSaudeHome(),
		filepath.Join(HubSaudeHome(), "assinador"),
		filepath.Join(HubSaudeHome(), "simulador"),
		filepath.Join(HubSaudeHome(), "processos"),
		filepath.Join(HubSaudeHome(), "cache"),
		filepath.Join(HubSaudeHome(), "logs"),
		filepath.Join(HubSaudeHome(), "jre"),
		filepath.Join(HubSaudeHome(), "jdk"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return nil
}

func AssinadorJarPath() string {
	return filepath.Join(HubSaudeHome(), "assinador", "assinador.jar")
}

func SimulatorJarPath() string {
	return filepath.Join(HubSaudeHome(), "simulador", "simulador.jar")
}

func ProcessDir() string {
	return filepath.Join(HubSaudeHome(), "processos")
}

func CacheDir() string {
	return filepath.Join(HubSaudeHome(), "cache")
}

func LogsDir() string {
	return filepath.Join(HubSaudeHome(), "logs")
}
