#!/usr/bin/env bash
# Verifica se a estrutura raiz do projeto, os workflows de CI/release,
# checksums, Cosign e rastreabilidade das Sprints 1 e 2 estão presentes.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

WORKFLOW="$REPO_ROOT/.github/workflows/release.yml"
BUILD_WORKFLOW="$REPO_ROOT/.github/workflows/build.yml"
DOC="$REPO_ROOT/docs/integridade-assinatura-artefatos.md"
SPRINT1_DOC="$REPO_ROOT/docs/sprint1-fundacao-entrega-continua.md"
SPRINT2_DOC="$REPO_ROOT/docs/sprint2-assinatura-digital-local.md"
SPRINT2_JAVA_DIR="$REPO_ROOT/projetos/assinador-java"
GITIGNORE="$REPO_ROOT/.gitignore"
GITATTRIBUTES="$REPO_ROOT/.gitattributes"
GENERATED_CHECK="$REPO_ROOT/scripts/check-generated-files.sh"
TEXT_FORMAT_CHECK="$REPO_ROOT/scripts/check-text-format.sh"
INTERNAL_PACKAGES_CHECK="$REPO_ROOT/scripts/check-internal-packages.sh"
MAKEFILE="$REPO_ROOT/Makefile"

for required in \
  "$WORKFLOW" \
  "$BUILD_WORKFLOW" \
  "$DOC" \
  "$SPRINT1_DOC" \
  "$SPRINT2_DOC" \
  "$GITIGNORE" \
  "$GITATTRIBUTES" \
  "$GENERATED_CHECK" \
  "$TEXT_FORMAT_CHECK" \
  "$INTERNAL_PACKAGES_CHECK" \
  "$MAKEFILE" \
  "$REPO_ROOT/go.mod"; do
  if [[ ! -e "$required" ]]; then
    echo "ERRO: item obrigatório não encontrado: $required" >&2
    exit 1
  fi
done

if [[ ! -d "$SPRINT2_JAVA_DIR" ]]; then
  echo "ERRO: diretório projetos/assinador-java não encontrado: $SPRINT2_JAVA_DIR" >&2
  exit 1
fi

if [[ -d "$REPO_ROOT/runner-implementacao" ]]; then
  echo "ERRO: a pasta runner-implementacao não deve existir nesta versão. O módulo Go deve ficar na raiz." >&2
  exit 1
fi

if grep -Fxq "release/" "$GITIGNORE"; then
  echo "ERRO: .gitignore contém release/ sem barra inicial; isso ignora internal/release/. Use /release/." >&2
  exit 1
fi
if ! grep -Fxq "/release/" "$GITIGNORE"; then
  echo "ERRO: .gitignore deve conter /release/ para ignorar apenas releases locais na raiz." >&2
  exit 1
fi


root_items=(
  "$REPO_ROOT/go.mod"
  "$REPO_ROOT/Makefile"
  "$REPO_ROOT/cmd/assinatura"
  "$REPO_ROOT/cmd/simulador"
  "$REPO_ROOT/internal"
  "$REPO_ROOT/internal/release"
  "$REPO_ROOT/internal/release/manifest.go"
  "$REPO_ROOT/assinador"
  "$REPO_ROOT/scripts"
  "$REPO_ROOT/examples"
  "$REPO_ROOT/docs"
  "$REPO_ROOT/projetos"
)
for item in "${root_items[@]}"; do
  if [[ ! -e "$item" ]]; then
    echo "ERRO: estrutura raiz esperada ausente: $item" >&2
    exit 1
  fi
done

for file in "$WORKFLOW" "$BUILD_WORKFLOW" "$GITIGNORE" "$GITATTRIBUTES" "$SPRINT1_DOC" "$SPRINT2_DOC" "$DOC" "$REPO_ROOT/README.md"; do
  if grep -Fq -- "runner-implementacao" "$file"; then
    echo "ERRO: referência antiga a runner-implementacao encontrada em $file" >&2
    exit 1
  fi
done

required_workflow_patterns=(
  "cache-dependency-path: go.sum"
  "go-version-file: go.mod"
  "path: dist/*"
  "actions/download-artifact@v4"
  "merge-multiple: true"
  "path: dist"
  "files: dist/*"
  "fail_on_unmatched_files: true"
  "assinatura-\${VERSION}-windows-amd64.exe"
  "assinatura-\${VERSION}-linux-amd64.AppImage"
  "assinatura-\${VERSION}-macos-amd64.dmg"
  "simulador-\${VERSION}-windows-amd64.exe"
  "simulador-\${VERSION}-linux-amd64.AppImage"
  "simulador-\${VERSION}-macos-amd64.dmg"
  "sha256sum"
  "checksums.txt"
  "sigstore/cosign-installer"
  "cosign sign-blob"
  "--tlog-upload=true"
  "--output-signature \"\$f.sig\""
  "--output-certificate \"\$f.pem\""
  "--bundle \"\$f.bundle\""
  "id-token: write"
  "contents: write"
  "softprops/action-gh-release"
)
for pattern in "${required_workflow_patterns[@]}"; do
  if ! grep -Fq -- "$pattern" "$WORKFLOW"; then
    echo "ERRO: padrão ausente no workflow de release: $pattern" >&2
    exit 1
  fi
done

for forbidden in "working-directory: runner-implementacao" "runner-implementacao/dist" "runner-implementacao/go.mod"; do
  if grep -Fq -- "$forbidden" "$WORKFLOW" "$BUILD_WORKFLOW"; then
    echo "ERRO: workflow ainda contém caminho antigo: $forbidden" >&2
    exit 1
  fi
done

required_build_patterns=(
  "branches:"
  "develop"
  "feature/**"
  "release/**"
  "refactor/**"
  "go test ./..."
  "go vet ./..."
  "Verificar pacotes internos obrigatórios"
  "./scripts/check-internal-packages.sh"
  "GOOS=linux GOARCH=amd64"
  "GOOS=windows GOARCH=amd64"
  "GOOS=darwin GOARCH=amd64"
  "actions/upload-artifact@v4"
  "runner-dev-binaries"
)
for pattern in "${required_build_patterns[@]}"; do
  if ! grep -Fq -- "$pattern" "$BUILD_WORKFLOW"; then
    echo "ERRO: padrão ausente no build.yml: $pattern" >&2
    exit 1
  fi
done

required_gitignore_patterns=(
  "**/target/"
  "**/out/"
  "**/dist/"
  "assinador/target/"
  "assinador/out/"
  "dist/"
  "examples/*.json"
)
for pattern in "${required_gitignore_patterns[@]}"; do
  if ! grep -Fq -- "$pattern" "$GITIGNORE"; then
    echo "ERRO: padrão ausente no .gitignore: $pattern" >&2
    exit 1
  fi
done

required_makefile_patterns=(
  ".PHONY: deps tidy test vet cover check text-check internal-check build java-test samples clean all help"
  "deps:"
  "go mod download"
  "tidy:"
  "go mod tidy"
  "test: deps"
  "go test ./..."
  "vet: deps"
  "go vet ./..."
  "cover: deps"
  "go test -cover ./..."
  "text-check:"
  "./scripts/check-text-format.sh"
  "internal-check:"
  "./scripts/check-internal-packages.sh"
  "check: deps text-check internal-check"
  "./scripts/check-generated-files.sh"
  "./scripts/check-release-artifacts.sh"
  "./scripts/check-text-format.sh"
  "build: deps"
  "go build -ldflags"
  "java-test:"
  "\$(MAKE) -C assinador clean all test"
  "all: deps test vet cover check java-test build"
)
for pattern in "${required_makefile_patterns[@]}"; do
  if ! grep -Fq -- "$pattern" "$MAKEFILE"; then
    echo "ERRO: padrão ausente no Makefile: $pattern" >&2
    exit 1
  fi
done

required_doc_patterns=(
  "<artefato>.sig"
  "<artefato>.pem"
  "cosign verify-blob"
  "--certificate assinatura-1.0.0-linux-amd64.AppImage.pem"
  "--signature assinatura-1.0.0-linux-amd64.AppImage.sig"
  "assinatura-1.0.0-linux-amd64.AppImage"
  "OIDC"
  "transparency log"
)
for pattern in "${required_doc_patterns[@]}"; do
  if ! grep -Fq -- "$pattern" "$DOC"; then
    echo "ERRO: padrão ausente na documentação de integridade: $pattern" >&2
    exit 1
  fi
done

required_sprint1_doc_patterns=(
  "US-01.1"
  "Cobra"
  "assinatura version"
  "US-05.1"
  "GitHub Actions"
  "artifacts do workflow"
  "US-05.2"
  "SemVer"
  "assinatura-<versão>-<os>-<arch>"
  "US-05.3"
  "checksums SHA256"
  "Cosign"
  "cosign verify-blob"
)
for pattern in "${required_sprint1_doc_patterns[@]}"; do
  if ! grep -Fq -- "$pattern" "$SPRINT1_DOC"; then
    echo "ERRO: padrão ausente na documentação da Sprint 1: $pattern" >&2
    exit 1
  fi
done

required_sprint2_doc_patterns=(
  "US-02.1"
  "projetos/assinador-java"
  "SignatureService"
  "FakeSignatureService"
  "US-02.2"
  "INPUT.MISSING-PARAMETER"
  "US-02.3"
  "valid: true"
  "US-01.2"
  "assinatura sign --help"
  "US-01.3"
  "java -jar assinador.jar"
  "US-01.4"
  "compactJSON"
  "US-04.1"
  "~/.hubsaude/jdk"
)
for pattern in "${required_sprint2_doc_patterns[@]}"; do
  if ! grep -Fq -- "$pattern" "$SPRINT2_DOC"; then
    echo "ERRO: padrão ausente na documentação da Sprint 2: $pattern" >&2
    exit 1
  fi
done

required_sprint2_code_patterns=(
  "$REPO_ROOT/assinador/src/br/ufg/hubsaude/assinador/SignatureService.java:interface SignatureService"
  "$REPO_ROOT/assinador/src/br/ufg/hubsaude/assinador/FakeSignatureService.java:final class FakeSignatureService"
  "$REPO_ROOT/assinador/src/br/ufg/hubsaude/assinador/Main.java:new FakeSignatureService"
  "$REPO_ROOT/cmd/assinatura/main.go:signUsage"
  "$REPO_ROOT/cmd/assinatura/main.go:validateUsage"
  "$REPO_ROOT/cmd/assinatura/main.go:parseFlagsKnown"
  "$REPO_ROOT/internal/assinador/client.go:java -jar"
  "$REPO_ROOT/internal/assinador/client.go:assinador.jar não encontrado"
  "$REPO_ROOT/internal/jdk/jdk.go:EnsureJava21"
  "$REPO_ROOT/internal/jdk/jdk.go:ManagedJavaPathFor"
  "$REPO_ROOT/internal/release/manifest.go:package release"
  "$REPO_ROOT/internal/release/manifest.go:DownloadAndInstallJRE"
)
for item in "${required_sprint2_code_patterns[@]}"; do
  file="${item%%:*}"
  pattern="${item#*:}"
  if ! grep -Fq -- "$pattern" "$file"; then
    echo "ERRO: padrão da Sprint 2 ausente em $file: $pattern" >&2
    exit 1
  fi
done

if ! grep -Fq -- "github.com/spf13/cobra" "$REPO_ROOT/go.mod"; then
  echo "ERRO: go.mod não declara dependência github.com/spf13/cobra." >&2
  exit 1
fi

if ! grep -Fq -- "module github.com/BarbaraNogueiraCS/runner" "$REPO_ROOT/go.mod"; then
  echo "ERRO: go.mod não está com a identidade do repositório BarbaraNogueiraCS/runner." >&2
  exit 1
fi

"$GENERATED_CHECK" >/dev/null

echo "OK: Sprint 1 e Sprint 2 cobertas; projeto Go está na raiz; não existe runner-implementacao; workflows usam go.mod e dist/ na raiz; release contém artefatos, checksums, Cosign, OIDC, transparency log, .sig e .pem; Makefile automatiza deps/test/vet/cover/check/build/all e check-text-format protege arquivos críticos contra perda de quebras de linha; check-internal-packages garante que internal/release foi commitado."
