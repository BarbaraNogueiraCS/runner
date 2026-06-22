package release

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func moduleRoot(t *testing.T) string {
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

func repositoryRoot(t *testing.T) string {
	t.Helper()
	dir := moduleRoot(t)
	for {
		workflow := filepath.Join(dir, ".github", "workflows", "release.yml")
		if _, err := os.Stat(workflow); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf(".github/workflows/release.yml não encontrado a partir de %s", moduleRoot(t))
		}
		dir = parent
	}
}

func readWorkflow(t *testing.T) string {
	t.Helper()
	workflowPath := filepath.Join(repositoryRoot(t), ".github", "workflows", "release.yml")
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

func TestReleaseWorkflowIsAtRepositoryRootAndUsesRootModule(t *testing.T) {
	repo := repositoryRoot(t)
	module := moduleRoot(t)
	if repo != module {
		t.Fatalf("o módulo Go deve estar na raiz do repositório; repo=%s module=%s", repo, module)
	}
	workflow := readWorkflow(t)

	required := []string{
		"path: dist/*",
		"path: dist",
		"files: dist/*",
		"fail_on_unmatched_files: true",
		"merge-multiple: true",
		"cache-dependency-path: go.sum",
	}
	for _, item := range required {
		if !strings.Contains(workflow, item) {
			t.Fatalf("workflow não está ajustado para o projeto na raiz: %s", item)
		}
	}
	if strings.Contains(workflow, "runner-implementacao") || strings.Contains(workflow, "working-directory: runner-implementacao") {
		t.Fatalf("workflow ainda contém caminhos antigos com runner-implementacao")
	}
}

func TestReleaseWorkflowRequiresSemVerTags(t *testing.T) {
	workflow := readWorkflow(t)

	if !strings.Contains(workflow, "^v[0-9]+\\.[0-9]+\\.[0-9]+$") {
		t.Fatalf("workflow deve validar tags SemVer no formato vMAJOR.MINOR.PATCH")
	}
	if !strings.Contains(workflow, "tags:") || !strings.Contains(workflow, "v[0-9]+.[0-9]+.[0-9]+") {
		t.Fatalf("workflow deve ser acionado por tags SemVer")
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
	docPath := filepath.Join(repositoryRoot(t), "docs", "integridade-assinatura-artefatos.md")
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

func TestGeneratedBuildOutputsAreIgnoredAndChecked(t *testing.T) {
	repo := repositoryRoot(t)

	gitignorePath := filepath.Join(repo, ".gitignore")
	gitignoreBytes, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf(".gitignore não encontrado na raiz do repositório: %v", err)
	}
	gitignore := string(gitignoreBytes)

	requiredIgnore := []string{
		"**/target/",
		"**/out/",
		"**/dist/",
		"assinador/target/",
		"assinador/out/",
		"dist/",
		"examples/*.json",
	}
	for _, item := range requiredIgnore {
		if !strings.Contains(gitignore, item) {
			t.Fatalf(".gitignore não protege saída gerada: %s", item)
		}
	}

	checkScript := filepath.Join(moduleRoot(t), "scripts", "check-generated-files.sh")
	content, err := os.ReadFile(checkScript)
	if err != nil {
		t.Fatalf("script check-generated-files.sh não encontrado: %v", err)
	}
	for _, item := range []string{"ls-files", "target", "out", "dist"} {
		if !strings.Contains(string(content), item) {
			t.Fatalf("script de verificação de gerados não contém item obrigatório: %s", item)
		}
	}

	for _, workflowName := range []string{"build.yml", "release.yml"} {
		workflow, err := os.ReadFile(filepath.Join(repo, ".github", "workflows", workflowName))
		if err != nil {
			t.Fatalf("workflow %s não encontrado: %v", workflowName, err)
		}
		if !strings.Contains(string(workflow), "./scripts/check-generated-files.sh") {
			t.Fatalf("workflow %s não executa check-generated-files.sh", workflowName)
		}
	}
}

func TestSprint1CLIUsesRepositoryIdentityAndCobra(t *testing.T) {
	moduleBytes, err := os.ReadFile(filepath.Join(moduleRoot(t), "go.mod"))
	if err != nil {
		t.Fatalf("go.mod não encontrado: %v", err)
	}
	module := string(moduleBytes)
	for _, item := range []string{
		"module github.com/BarbaraNogueiraCS/runner",
		"github.com/spf13/cobra",
	} {
		if !strings.Contains(module, item) {
			t.Fatalf("go.mod não contém requisito da Sprint 1: %s", item)
		}
	}

	for _, path := range []string{
		filepath.Join(moduleRoot(t), "cmd", "assinatura", "main.go"),
		filepath.Join(moduleRoot(t), "cmd", "simulador", "main.go"),
	} {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("arquivo CLI não encontrado %s: %v", path, err)
		}
		text := string(content)
		for _, item := range []string{"github.com/spf13/cobra", "Use:", "version", "SetVersionTemplate"} {
			if !strings.Contains(text, item) {
				t.Fatalf("%s não evidencia uso do Cobra/version: %s", path, item)
			}
		}
	}
}

func TestSprint1BuildWorkflowProducesMultiplatformArtifactsOnPush(t *testing.T) {
	repo := repositoryRoot(t)
	workflowPath := filepath.Join(repo, ".github", "workflows", "build.yml")
	content, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatalf("build.yml não encontrado: %v", err)
	}
	workflow := string(content)
	for _, item := range []string{
		"push:",
		"main",
		"GOOS=linux GOARCH=amd64",
		"GOOS=windows GOARCH=amd64",
		"GOOS=darwin GOARCH=amd64",
		"actions/upload-artifact@v4",
		"runner-dev-binaries",
		"assinatura-linux-amd64.AppImage",
		"assinatura-windows-amd64.exe",
		"assinatura-macos-amd64",
	} {
		if !strings.Contains(workflow, item) {
			t.Fatalf("build.yml não contém requisito da Sprint 1 CI/CD: %s", item)
		}
	}
}

func TestSprint1DocumentationCoversAcceptanceCriteria(t *testing.T) {
	docPath := filepath.Join(repositoryRoot(t), "docs", "sprint1-fundacao-entrega-continua.md")
	content, err := os.ReadFile(docPath)
	if err != nil {
		t.Fatalf("documento de rastreabilidade da Sprint 1 não encontrado: %v", err)
	}
	doc := string(content)
	for _, item := range []string{
		"US-01.1",
		"Cobra",
		"assinatura version",
		"Windows",
		"Linux",
		"macOS",
		"US-05.1",
		"GitHub Actions",
		"artifacts do workflow",
		"US-05.2",
		"SemVer",
		"GitHub Releases",
		"assinatura-<versão>-<os>-<arch>",
		"US-05.3",
		"checksums SHA256",
		"Cosign",
		"cosign verify-blob",
	} {
		if !strings.Contains(doc, item) {
			t.Fatalf("documentação da Sprint 1 não contém requisito: %s", item)
		}
	}
}
