#!/usr/bin/env bash
# Verifica, localmente, se o workflow de release contém os nomes, checksums,
# publicação em GitHub Releases e assinatura Cosign keyless/OIDC esperados.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WORKFLOW="$ROOT/.github/workflows/release.yml"
DOC="$ROOT/docs/integridade-assinatura-artefatos.md"

if [[ ! -f "$WORKFLOW" ]]; then
  echo "ERRO: workflow de release não encontrado: $WORKFLOW" >&2
  exit 1
fi

if [[ ! -f "$DOC" ]]; then
  echo "ERRO: documento de integridade de artefatos não encontrado: $DOC" >&2
  exit 1
fi

required_workflow_patterns=(
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
    echo "ERRO: padrão ausente no workflow: $pattern" >&2
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

echo "Workflow e documentação contêm artefatos, checksums, Cosign, OIDC, transparency log, .sig e .pem esperados."
