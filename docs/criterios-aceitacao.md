# Critérios de aceitação rastreados

| Critério | Evidência no projeto |
|---|---|
| CLI Go estruturado | `go.mod`, `cmd/assinatura`, `cmd/simulador`, `internal/*` |
| Automação local | `Makefile` na raiz |
| Assinador Java | `assinador/src`, `assinador/Makefile` |
| Diretório solicitado na Sprint 2 | `projetos/assinador-java` |
| GitHub Actions | `.github/workflows/build.yml`, `.github/workflows/release.yml` |
| Binário Windows amd64 | release gera `*.exe` |
| Binário Linux amd64 | release gera `*.AppImage` |
| Binário macOS amd64 | release gera `*.dmg` |
| Checksums SHA256 | `checksums.txt` no workflow de release |
| Assinatura Cosign | `.sig`, `.pem`, `.bundle` no workflow de release |
| Higiene de arquivos gerados | `.gitignore`, `.gitattributes`, `scripts/check-generated-files.sh` |
