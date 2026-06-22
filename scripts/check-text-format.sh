#!/usr/bin/env bash
# Verifica se arquivos críticos foram preservados com quebras de linha LF.
set -euo pipefail

critical_files=(
  go.mod
  go.sum
  .github/workflows/build.yml
  .github/workflows/release.yml
  Makefile
)

for f in "${critical_files[@]}"; do
  if [[ ! -f "$f" ]]; then
    echo "ERRO: arquivo crítico não encontrado: $f" >&2
    exit 1
  fi
  lines=$(wc -l < "$f" | tr -d ' ')
  if [[ "$lines" -lt 2 ]]; then
    echo "ERRO: $f parece ter perdido quebras de linha. Reextraia o .tar.gz/.zip e faça git add -A." >&2
    exit 1
  fi
  if grep -q $'\r' "$f"; then
    echo "ERRO: $f contém CRLF/CR. Use LF." >&2
    exit 1
  fi
done

echo "OK: arquivos críticos preservam quebras de linha LF."
