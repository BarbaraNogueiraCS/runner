#!/usr/bin/env bash
# Garante que saídas geradas por build/teste não sejam versionadas.
# Deve ser executado a partir de runner-implementacao ou de qualquer subpasta do repositório.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MODULE_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

find_repo_root() {
  local dir="$MODULE_ROOT"
  while [[ "$dir" != "/" ]]; do
    if [[ -d "$dir/.git" || -d "$dir/.github" ]]; then
      printf '%s\n' "$dir"
      return 0
    fi
    dir="$(dirname "$dir")"
  done
  return 1
}

REPO_ROOT="$(find_repo_root || true)"
if [[ -z "${REPO_ROOT:-}" ]]; then
  echo "ERRO: raiz do repositório não encontrada." >&2
  exit 1
fi

GITIGNORE="$REPO_ROOT/.gitignore"
if [[ ! -f "$GITIGNORE" ]]; then
  echo "ERRO: .gitignore não encontrado na raiz do repositório: $GITIGNORE" >&2
  exit 1
fi

required_ignore_patterns=(
  "**/target/"
  "**/out/"
  "**/dist/"
  "runner-implementacao/assinador/target/"
  "runner-implementacao/assinador/out/"
  "runner-implementacao/dist/"
  "runner-implementacao/examples/*.json"
)

for pattern in "${required_ignore_patterns[@]}"; do
  if ! grep -Fq -- "$pattern" "$GITIGNORE"; then
    echo "ERRO: padrão ausente no .gitignore: $pattern" >&2
    exit 1
  fi
done

if command -v git >/dev/null 2>&1 && git -C "$REPO_ROOT" rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  tracked_generated="$(git -C "$REPO_ROOT" ls-files | grep -E '(^|/)(dist|target|out)/|^runner-implementacao/examples/.*\.json$' || true)"
  if [[ -n "$tracked_generated" ]]; then
    echo "ERRO: arquivos gerados estão rastreados pelo Git. Remova-os do índice com git rm --cached." >&2
    echo "$tracked_generated" >&2
    exit 1
  fi
fi

echo "OK: target/, out/, dist/ e JSONs gerados estão protegidos pelo .gitignore e não aparecem como arquivos rastreados."
