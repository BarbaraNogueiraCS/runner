package release

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func projectRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("não foi possível obter diretório atual: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("go.mod não encontrado a partir de %s", dir)
		}
		dir = parent
	}
}

func TestReleaseWorkflowDeclaresExecutableArtifacts(t *testing.T) {
	root := projectRoot(t)
	workflowPath := filepath.Join(root, ".github", "workflows", "release.yml")
	content, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatalf("workflow de release não encontrado: %v", err)
	}
	workflow := string(content)

	required := []string{
		"assinatura-${VERSION}-windows-amd64.exe",
		"assinatura-${VERSION}-linux-amd64.AppImage",
		"assinatura-${VERSION}-macos-amd64.dmg",
		"simulador-${VERSION}-windows-amd64.exe",
		"simulador-${VERSION}-linux-amd64.AppImage",
		"simulador-${VERSION}-macos-amd64.dmg",
		"softprops/action-gh-release",
		"sha256sum",
		"checksums.txt",
		"cosign sign-blob",
		"id-token: write",
		"contents: write",
	}
	for _, item := range required {
		if !strings.Contains(workflow, item) {
			t.Fatalf("workflow de release não contém item obrigatório: %s", item)
		}
	}
}

func TestReleaseWorkflowRequiresSemVerTags(t *testing.T) {
	root := projectRoot(t)
	workflowPath := filepath.Join(root, ".github", "workflows", "release.yml")
	content, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatalf("workflow de release não encontrado: %v", err)
	}
	workflow := string(content)

	if !strings.Contains(workflow, "^v[0-9]+\\.[0-9]+\\.[0-9]+$") {
		t.Fatalf("workflow deve validar tags SemVer no formato vMAJOR.MINOR.PATCH")
	}
}
