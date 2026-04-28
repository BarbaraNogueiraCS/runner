package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestVersionCommandPrintsDefaultDevVersion(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "version")
	cmd.Dir = "."

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go run . version falhou: %v\nSaída:\n%s", err, string(output))
	}

	if !strings.Contains(string(output), "dev") {
		t.Fatalf("saída esperada contendo %q, obtida: %q", "dev", string(output))
	}
}
