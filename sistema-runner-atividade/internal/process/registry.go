package process

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kyriosdata/assinatura/internal/config"
)

type ProcessRegistry struct {
	Dir string
}

func NewProcessRegistry() ProcessRegistry {
	return ProcessRegistry{Dir: config.ProcessDir()}
}

func (r ProcessRegistry) path(application string, port int) string {
	return filepath.Join(r.Dir, fmt.Sprintf("%s-%d.json", application, port))
}

func (r ProcessRegistry) Save(metadata ProcessMetadata) error {
	if err := os.MkdirAll(r.Dir, 0o755); err != nil {
		return err
	}
	content, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.path(metadata.Application, metadata.Port), content, 0o644)
}

func (r ProcessRegistry) Find(application string, port int) (ProcessMetadata, error) {
	content, err := os.ReadFile(r.path(application, port))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ProcessMetadata{}, os.ErrNotExist
		}
		return ProcessMetadata{}, err
	}
	var metadata ProcessMetadata
	if err := json.Unmarshal(content, &metadata); err != nil {
		return ProcessMetadata{}, err
	}
	return metadata, nil
}

func (r ProcessRegistry) Remove(application string, port int) error {
	err := os.Remove(r.path(application, port))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
