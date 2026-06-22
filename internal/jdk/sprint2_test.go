package jdk

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestSprint2ManagedJavaPathForSupportedPlatforms(t *testing.T) {
	home := filepath.Join("tmp", "runner-home")
	cases := []struct {
		goos   string
		suffix string
	}{
		{goos: "linux", suffix: filepath.Join(".hubsaude", "jdk", "bin", "java")},
		{goos: "darwin", suffix: filepath.Join(".hubsaude", "jdk", "bin", "java")},
		{goos: "windows", suffix: filepath.Join(".hubsaude", "jdk", "bin", "java.exe")},
	}
	for _, tc := range cases {
		got := ManagedJavaPathFor(home, tc.goos)
		if !strings.HasSuffix(got, tc.suffix) {
			t.Fatalf("caminho gerenciado para %s deveria terminar com %s, obtido %s", tc.goos, tc.suffix, got)
		}
	}
}
