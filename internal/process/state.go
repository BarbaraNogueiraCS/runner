package process

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/BarbaraNogueiraCS/runner/internal/paths"
)

type State struct {
	Name      string    `json:"name"`
	PID       int       `json:"pid"`
	Port      int       `json:"port"`
	URL       string    `json:"url"`
	JarPath   string    `json:"jar_path,omitempty"`
	StartedAt time.Time `json:"started_at"`
}

func Save(s State) error {
	if err := paths.EnsureRuntimeDirs(); err != nil {
		return err
	}
	d, err := paths.ProcessDir()
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(d, s.Name+".json"), b, 0o644)
}

func Load(name string) (State, error) {
	d, err := paths.ProcessDir()
	if err != nil {
		return State{}, err
	}
	b, err := os.ReadFile(filepath.Join(d, name+".json"))
	if err != nil {
		return State{}, err
	}
	var s State
	if err := json.Unmarshal(b, &s); err != nil {
		return State{}, err
	}
	return s, nil
}

func Delete(name string) error {
	d, err := paths.ProcessDir()
	if err != nil {
		return err
	}
	p := filepath.Join(d, name+".json")
	if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("não foi possível remover registro de processo %s: %w", p, err)
	}
	return nil
}
