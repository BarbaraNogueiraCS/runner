#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/../projetos/assinador-java"
rm -rf build
mkdir -p build/classes
javac -encoding UTF-8 -d build/classes $(find src/main/java -name '*.java')
jar --create --file build/assinador.jar --main-class br.ufg.hubsaude.assinador.Main -C build/classes .
echo "Gerado: projetos/assinador-java/build/assinador.jar"
