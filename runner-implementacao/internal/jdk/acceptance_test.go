package jdk

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestAcceptanceMajorVersionParsesJava21Output(t *testing.T) {
	java := fakeJavaExecutable(t, `openjdk version "21.0.5" 2024-10-15`)
	major, text, err := MajorVersion(java)
	if err != nil {
		t.Fatalf("MajorVersion retornou erro: %v", err)
	}
	if major != 21 {
		t.Fatalf("major esperado 21, obtido %d; texto=%s", major, text)
	}
}

func TestErrorScenarioMajorVersionRejectsUnrecognizedOutput(t *testing.T) {
	java := fakeJavaExecutable(t, `versao invalida`)
	_, _, err := MajorVersion(java)
	if err == nil {
		t.Fatalf("MajorVersion deveria falhar para saída não reconhecida")
	}
}

func fakeJavaExecutable(t *testing.T, versionLine string) string {
	t.Helper()
	dir := t.TempDir()
	if runtime.GOOS == "windows" {
		path := filepath.Join(dir, "java.bat")
		content := "@echo off\r\necho " + versionLine + " 1>&2\r\nexit /b 0\r\n"
		if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
			t.Fatalf("não foi possível criar java falso: %v", err)
		}
		return path
	}
	path := filepath.Join(dir, "java")
	content := "#!/bin/sh\necho '" + versionLine + "' 1>&2\nexit 0\n"
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatalf("não foi possível criar java falso: %v", err)
	}
	return path
}
