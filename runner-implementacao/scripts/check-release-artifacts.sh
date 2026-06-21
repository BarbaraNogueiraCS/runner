#!/usr/bin/env bash
# Verifica, localmente, se o workflow de release está na raiz do repositório
# e contém nomes, checksums, publicação em GitHub Releases e assinatura Cosign OIDC.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MODULE_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

find_repo_root() {
  local dir="$MODULE_ROOT"
  while [[ "$dir" != "/" ]]; do
    if [[ -f "$dir/.github/workflows/release.yml" ]]; then
      printf '%s\n' "$dir"
      return 0
    fi
    dir="$(dirname "$dir")"
  done
  return 1
}

REPO_ROOT="$(find_repo_root || true)"
if [[ -z "${REPO_ROOT:-}" ]]; then
  echo "ERRO: workflow de release não encontrado na raiz do repositório ou em diretórios superiores." >&2
  echo "Esperado: <raiz-do-repositorio>/.github/workflows/release.yml" >&2
  exit 1
fi

WORKFLOW="$REPO_ROOT/.github/workflows/release.yml"
BUILD_WORKFLOW="$REPO_ROOT/.github/workflows/build.yml"
DOC="$REPO_ROOT/docs/integridade-assinatura-artefatos.md"
SPRINT1_DOC="$REPO_ROOT/docs/sprint1-fundacao-entrega-continua.md"
GITIGNORE="$REPO_ROOT/.gitignore"
GITATTRIBUTES="$REPO_ROOT/.gitattributes"
GENERATED_CHECK="$MODULE_ROOT/scripts/check-generated-files.sh"
MAKEFILE="$MODULE_ROOT/Makefile"

if [[ ! -f "$DOC" ]]; then
  echo "ERRO: documento de integridade de artefatos não encontrado: $DOC" >&2
  exit 1
fi

if [[ ! -f "$SPRINT1_DOC" ]]; then
  echo "ERRO: documento de rastreabilidade da Sprint 1 não encontrado: $SPRINT1_DOC" >&2
  exit 1
fi

if [[ ! -f "$GITIGNORE" ]]; then
  echo "ERRO: .gitignore deve ficar na raiz do repositório: $GITIGNORE" >&2
  exit 1
fi

if [[ ! -f "$GITATTRIBUTES" ]]; then
  echo "ERRO: .gitattributes deve ficar na raiz do repositório: $GITATTRIBUTES" >&2
  exit 1
fi

if [[ ! -x "$GENERATED_CHECK" ]]; then
  echo "ERRO: script de verificação de arquivos gerados ausente ou sem permissão de execução: $GENERATED_CHECK" >&2
  exit 1
fi

if [[ ! -f "$MAKEFILE" ]]; then
  echo "ERRO: Makefile de automação local não encontrado: $MAKEFILE" >&2
  exit 1
fi

if [[ -d "$MODULE_ROOT/docs" ]]; then
  echo "ERRO: docs não deve ficar dentro de runner-implementacao. Use <raiz-do-repositorio>/docs." >&2
  exit 1
fi

if [[ -f "$MODULE_ROOT/.github/workflows/release.yml" ]]; then
  echo "ERRO: workflow duplicado encontrado dentro de runner-implementacao/.github." >&2
  echo "Mantenha workflows somente em <raiz-do-repositorio>/.github/workflows/." >&2
  exit 1
fi

required_workflow_patterns=(
  "working-directory: runner-implementacao"
  "cache-dependency-path: runner-implementacao/go.mod"
  "path: runner-implementacao/dist/*"
  "actions/download-artifact@v4"
  "merge-multiple: true"
  "path: runner-implementacao/dist"
  "files: runner-implementacao/dist/*"
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

if [[ -f "$BUILD_WORKFLOW" ]]; then
  required_build_patterns=(
    "working-directory: runner-implementacao"
    "branches:"
    "develop"
    "feature/**"
    "go test ./..."
    "go vet ./..."
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
fi


required_gitignore_patterns=(
  "**/target/"
  "**/out/"
  "**/dist/"
  "runner-implementacao/assinador/target/"
  "runner-implementacao/assinador/out/"
  "runner-implementacao/dist/"
  "runner-implementacao/examples/*.json"
)

for pattern in "${required_gitignore_patterns[@]}"; do
  if ! grep -Fq -- "$pattern" "$GITIGNORE"; then
    echo "ERRO: padrão ausente no .gitignore: $pattern" >&2
    exit 1
  fi
done

required_makefile_patterns=(
  ".PHONY: deps tidy test vet cover check build java-test samples clean all help"
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
  "check: deps"
  "./scripts/check-generated-files.sh"
  "./scripts/check-release-artifacts.sh"
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

"$GENERATED_CHECK" >/dev/null

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

required_sprint_doc_patterns=(
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

for pattern in "${required_sprint_doc_patterns[@]}"; do
  if ! grep -Fq -- "$pattern" "$SPRINT1_DOC"; then
    echo "ERRO: padrão ausente na documentação da Sprint 1: $pattern" >&2
    exit 1
  fi
done

if ! grep -Fq -- "github.com/spf13/cobra" "$MODULE_ROOT/go.mod"; then
  echo "ERRO: go.mod não declara dependência github.com/spf13/cobra." >&2
  exit 1
fi

if ! grep -Fq -- "module github.com/BarbaraNogueiraCS/runner" "$MODULE_ROOT/go.mod"; then
  echo "ERRO: go.mod não está com a identidade do repositório BarbaraNogueiraCS/runner." >&2
  exit 1
fi

echo "OK: Sprint 1 coberta; .github, docs, .gitignore e .gitattributes estão na raiz; arquivos gerados não são versionados; workflows usam runner-implementacao; release contém artefatos, checksums, Cosign, OIDC, transparency log, .sig e .pem; Makefile automatiza deps/test/vet/cover/check/build/all."
