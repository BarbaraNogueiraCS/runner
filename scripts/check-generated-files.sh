#!/usr/bin/env bash
# Garante que saídas geradas por build/teste não sejam versionadas.
# Pode ser executado a partir da raiz ou de qualquer subpasta do repositório.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

GITIGNORE="$REPO_ROOT/.gitignore"
if [[ ! -f "$GITIGNORE" ]]; then
  echo "ERRO: .gitignore não encontrado na raiz do repositório: $GITIGNORE" >&2
  exit 1
fi

required_ignore_patterns=(
  "**/target/"
  "**/out/"
  "**/dist/"
  "assinador/target/"
  "assinador/out/"
  "dist/"
  "examples/*.json"
)

for pattern in "${required_ignore_patterns[@]}"; do
  if ! grep -Fq -- "$pattern" "$GITIGNORE"; then
    echo "ERRO: padrão ausente no .gitignore: $pattern" >&2
    exit 1
  fi
done

if command -v git >/dev/null 2>&1 && git -C "$REPO_ROOT" rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  tracked_generated="$(git -C "$REPO_ROOT" ls-files | grep -E '(^|/)(dist|target|out)/|^examples/.*\.json$|^release/' || true)"
  if [[ -n "$tracked_generated" ]]; then
    echo "ERRO: arquivos gerados estão rastreados pelo Git. Remova-os do índice com git rm --cached." >&2
    echo "$tracked_generated" >&2
    exit 1
  fi
fi

echo "OK: target/, out/, dist/ e JSONs gerados estão protegidos pelo .gitignore e não aparecem como arquivos rastreados."
