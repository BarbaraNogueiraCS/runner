package jdk

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindJavaUsesRunnerJavaWhenValid(t *testing.T) {
	tempDir := t.TempDir()
	javaPath := filepath.Join(tempDir, "java")
	if err := os.WriteFile(javaPath, []byte("fake java"), 0o755); err != nil {
		t.Fatalf("não foi possível criar java falso: %v", err)
	}
	t.Setenv("RUNNER_JAVA", javaPath)

	got, err := FindJava()
	if err != nil {
		t.Fatalf("FindJava retornou erro: %v", err)
	}
	if got != javaPath {
		t.Fatalf("FindJava deveria retornar %s, obtido %s", javaPath, got)
	}
}

func TestFindJavaFailsWhenRunnerJavaDoesNotExist(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "java-inexistente")
	t.Setenv("RUNNER_JAVA", missing)

	_, err := FindJava()
	if err == nil {
		t.Fatalf("FindJava deveria falhar quando RUNNER_JAVA aponta para arquivo inexistente")
	}
}
