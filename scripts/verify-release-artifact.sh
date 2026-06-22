#!/usr/bin/env bash
# Verifica um artefato de release usando Cosign, o certificado .pem e a assinatura .sig.
set -euo pipefail

if [[ $# -ne 1 ]]; then
  echo "Uso: $0 <artefato>" >&2
  echo "Exemplo: $0 assinatura-1.0.0-linux-amd64.AppImage" >&2
  exit 2
fi

ARTIFACT="$1"
CERTIFICATE="${ARTIFACT}.pem"
SIGNATURE="${ARTIFACT}.sig"

if [[ ! -f "$ARTIFACT" ]]; then
  echo "ERRO: artefato não encontrado: $ARTIFACT" >&2
  exit 1
fi
if [[ ! -f "$CERTIFICATE" ]]; then
  echo "ERRO: certificado não encontrado: $CERTIFICATE" >&2
  exit 1
fi
if [[ ! -f "$SIGNATURE" ]]; then
  echo "ERRO: assinatura não encontrada: $SIGNATURE" >&2
  exit 1
fi

cosign verify-blob \
  --certificate "$CERTIFICATE" \
  --signature "$SIGNATURE" \
  "$ARTIFACT"
