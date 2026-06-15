#!/usr/bin/env bash
# Verifica, localmente, se o workflow de release contém os nomes, checksums,
# publicação em GitHub Releases e assinatura Cosign esperados.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WORKFLOW="$ROOT/.github/workflows/release.yml"

if [[ ! -f "$WORKFLOW" ]]; then
  echo "ERRO: workflow de release não encontrado: $WORKFLOW" >&2
  exit 1
fi

required_patterns=(
  "assinatura-\${VERSION}-windows-amd64.exe"
  "assinatura-\${VERSION}-linux-amd64.AppImage"
  "assinatura-\${VERSION}-macos-amd64.dmg"
  "simulador-\${VERSION}-windows-amd64.exe"
  "simulador-\${VERSION}-linux-amd64.AppImage"
  "simulador-\${VERSION}-macos-amd64.dmg"
  "sha256sum"
  "checksums.txt"
  "cosign sign-blob"
  "id-token: write"
  "softprops/action-gh-release"
)

for pattern in "${required_patterns[@]}"; do
  if ! grep -Fq "$pattern" "$WORKFLOW"; then
    echo "ERRO: padrão ausente no workflow: $pattern" >&2
    exit 1
  fi
done

echo "Workflow de release contém todos os artefatos, checksums, Cosign e GitHub Releases esperados."
