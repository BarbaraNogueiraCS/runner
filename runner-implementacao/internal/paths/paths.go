package paths

import (
	"os"
	"path/filepath"
)

// Home returns the managed Runner directory. It can be overridden in tests by
// RUNNER_HOME. By default it uses ~/.hubsaude.
func Home() (string, error) {
	if v := os.Getenv("RUNNER_HOME"); v != "" {
		return v, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".hubsaude"), nil
}

func ProcessDir() (string, error) {
	h, err := Home()
	if err != nil {
		return "", err
	}
	return filepath.Join(h, "processos"), nil
}

func LogDir() (string, error) {
	h, err := Home()
	if err != nil {
		return "", err
	}
	return filepath.Join(h, "logs"), nil
}

func CacheDir() (string, error) {
	h, err := Home()
	if err != nil {
		return "", err
	}
	return filepath.Join(h, "cache"), nil
}

func EnsureRuntimeDirs() error {
	for _, fn := range []func() (string, error){Home, ProcessDir, LogDir, CacheDir} {
		d, err := fn()
		if err != nil {
			return err
		}
		if err := os.MkdirAll(d, 0o755); err != nil {
			return err
		}
	}
	return nil
}
