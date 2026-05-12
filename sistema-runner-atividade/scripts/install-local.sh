#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
mkdir -p "$HOME/.hubsaude/assinador" "$HOME/.hubsaude/simulador" "$ROOT/bin"
cd "$ROOT"
GOTOOLCHAIN=local go build -o bin/assinatura ./cmd/assinatura
GOTOOLCHAIN=local go build -o bin/simulador ./cmd/simulador
bash scripts/build-assinador.sh
cp projetos/assinador-java/build/assinador.jar "$HOME/.hubsaude/assinador/assinador.jar"
echo "Instalação local concluída. Binários em ./bin e assinador.jar em ~/.hubsaude/assinador/"
