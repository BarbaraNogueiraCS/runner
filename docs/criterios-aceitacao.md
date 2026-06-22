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
| Sprint 3 — modo servidor do assinador | `assinador/src/.../SignatureController.java`, `cmd/assinatura`, `internal/assinador`, `internal/process` |
| Sprint 4 — CLI do Simulador | `cmd/simulador`, `internal/simulador`, `internal/release`, `docs/sprint4-simulador-seguranca-artefatos.md` |
| Download dinâmico do simulador.jar | `release.json`, `internal/release/manifest.go`, `simulador start --source`, cache em `~/.hubsaude/simulador/` |
| Verificação de integridade do download | SHA-256 em `release.json` e `--sha256`/`--checksum` para `--source` |
