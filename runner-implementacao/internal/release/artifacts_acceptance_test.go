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

func readWorkflow(t *testing.T) string {
	t.Helper()
	root := projectRoot(t)
	workflowPath := filepath.Join(root, ".github", "workflows", "release.yml")
	content, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatalf("workflow de release não encontrado: %v", err)
	}
	return string(content)
}

func TestReleaseWorkflowDeclaresExecutableArtifacts(t *testing.T) {
	workflow := readWorkflow(t)

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
		"contents: write",
	}
	for _, item := range required {
		if !strings.Contains(workflow, item) {
			t.Fatalf("workflow de release não contém item obrigatório: %s", item)
		}
	}
}

func TestReleaseWorkflowRequiresSemVerTags(t *testing.T) {
	workflow := readWorkflow(t)

	if !strings.Contains(workflow, "^v[0-9]+\\.[0-9]+\\.[0-9]+$") {
		t.Fatalf("workflow deve validar tags SemVer no formato vMAJOR.MINOR.PATCH")
	}
}

func TestReleaseWorkflowSignsArtifactsWithCosignOIDCAndTransparencyLog(t *testing.T) {
	workflow := readWorkflow(t)

	required := []string{
		"id-token: write",
		"sigstore/cosign-installer",
		"cosign sign-blob",
		"--tlog-upload=true",
		"--output-signature \"$f.sig\"",
		"--output-certificate \"$f.pem\"",
		"--bundle \"$f.bundle\"",
		"for f in assinatura-* simulador-* assinador-*.jar checksums.txt; do",
	}
	for _, item := range required {
		if !strings.Contains(workflow, item) {
			t.Fatalf("workflow de release não contém requisito de assinatura Cosign/OIDC: %s", item)
		}
	}
}

func TestReleaseDocumentationExplainsCosignVerification(t *testing.T) {
	root := projectRoot(t)
	docPath := filepath.Join(root, "docs", "integridade-assinatura-artefatos.md")
	content, err := os.ReadFile(docPath)
	if err != nil {
		t.Fatalf("documento de integridade de artefatos não encontrado: %v", err)
	}
	doc := string(content)

	required := []string{
		"cosign verify-blob",
		"--certificate assinatura-1.0.0-linux-amd64.AppImage.pem",
		"--signature assinatura-1.0.0-linux-amd64.AppImage.sig",
		"assinatura-1.0.0-linux-amd64.AppImage",
		"OIDC",
		"transparency log",
		"<artefato>.sig",
		"<artefato>.pem",
	}
	for _, item := range required {
		if !strings.Contains(doc, item) {
			t.Fatalf("documentação de integridade não contém item obrigatório: %s", item)
		}
	}
}
