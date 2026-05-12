#!/usr/bin/env bash
set -euo pipefail
if [ "$#" -ne 3 ]; then
  echo "Uso: $0 <artefato> <assinatura.sig> <certificado.pem>"
  exit 1
fi
cosign verify-blob --certificate "$3" --signature "$2" "$1"
