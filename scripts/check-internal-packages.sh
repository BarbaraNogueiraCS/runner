#!/usr/bin/env bash
# Verifica se os pacotes internos obrigatórios existem e são visíveis para o Go.
# Esse check evita o erro de CI:
#   no required module provides package github.com/BarbaraNogueiraCS/runner/internal/release
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"

required_files=(
  internal/release/manifest.go
  internal/release/manifest_test.go
  internal/release/acceptance_test.go
  internal/release/artifacts_acceptance_test.go
  internal/jdk/jdk.go
)

for f in "${required_files[@]}"; do
  if [[ ! -f "$f" ]]; then
    echo "ERRO: arquivo interno obrigatório ausente: $f" >&2
    echo "Esse arquivo deve ser commitado. Rode: git add internal/release internal/jdk" >&2
    exit 1
  fi
done

if ! grep -Fq "package release" internal/release/manifest.go; then
  echo "ERRO: internal/release/manifest.go não declara package release." >&2
  exit 1
fi

if ! grep -Fq 'github.com/BarbaraNogueiraCS/runner/internal/release' internal/jdk/jdk.go; then
  echo "ERRO: internal/jdk/jdk.go não importa o pacote interno release esperado." >&2
  exit 1
fi

# go list nesses pacotes internos não depende do Cobra e detecta cedo se a pasta
# internal/release ficou fora do commit enviado ao GitHub.
go list ./internal/release >/dev/null
go list ./internal/jdk >/dev/null

echo "OK: pacotes internos obrigatórios presentes e visíveis para o Go."
