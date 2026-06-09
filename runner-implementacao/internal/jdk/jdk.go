package jdk

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
