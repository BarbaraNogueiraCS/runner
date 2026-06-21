# Sprint 1 — Fundação e Entrega Contínua

Este documento registra a rastreabilidade dos artefatos da Sprint 1 no código gerado do Sistema Runner.

## Decisão de identidade do repositório

O enunciado original menciona `go mod init github.com/kyriosdata/assinatura`, mas esta implementação foi ajustada para a identidade do repositório de entrega da aluna. Por isso, o módulo Go oficial é:

```go
module github.com/BarbaraNogueiraCS/runner
```

Essa decisão mantém a autoria e a rastreabilidade do projeto no repositório `BarbaraNogueiraCS/runner`. Os pacotes internos passam a aparecer nos testes como `github.com/BarbaraNogueiraCS/runner/...`.

## US-01.1 — Estrutura base do CLI em Go

| Critério | Situação | Evidência no código |
|---|---:|---|
| Projeto Go inicializado | Implementado | `runner-implementacao/go.mod` |
| Identidade do módulo | Implementado com identidade própria | `module github.com/BarbaraNogueiraCS/runner` |
| CLI usando Cobra | Implementado | `cmd/assinatura/main.go` e `cmd/simulador/main.go` importam `github.com/spf13/cobra` |
| Instalação da ferramenta Cobra | Documentado | `go install github.com/spf13/cobra-cli@latest` |
| Comando `assinatura version` | Implementado | comando `version` do Cobra em `cmd/assinatura/main.go` |
| Estrutura de pacotes definida | Implementado | `cmd/`, `internal/`, `assinador/`, `scripts/`, `docs/` |
| Estrutura documentada | Implementado | `README.md`, `runner-implementacao/README.md` e `docs/criterios-aceitacao.md` |
| Compilação para Windows, Linux e macOS | Implementado | `.github/workflows/build.yml` e `.github/workflows/release.yml` |

Comando de instalação da ferramenta usada para criação/manutenção de CLIs com Cobra:

```bash
go install github.com/spf13/cobra-cli@latest
```

Observação: o projeto já contém a implementação com Cobra no código-fonte, então o usuário final não precisa instalar `cobra-cli` para executar os binários. A ferramenta é útil para desenvolvimento.

## US-05.1 — Pipeline CI/CD multiplataforma

| Critério | Situação | Evidência no código |
|---|---:|---|
| GitHub Actions com workflow de build | Implementado | `.github/workflows/build.yml` |
| Build a cada push na branch principal | Implementado | `on: push`, branch `main` |
| Build também em `develop` e `feature/**` | Implementado | `.github/workflows/build.yml` |
| Cross-compilation Windows amd64 | Implementado | `GOOS=windows GOARCH=amd64` |
| Cross-compilation Linux amd64 | Implementado | `GOOS=linux GOARCH=amd64` |
| Cross-compilation macOS amd64 | Implementado | `GOOS=darwin GOARCH=amd64` |
| Artifacts do workflow | Implementado | `actions/upload-artifact@v4`, artifact `runner-dev-binaries` |

O build de desenvolvimento gera artefatos de validação no workflow. A release definitiva é gerada quando uma tag SemVer é publicada.

## US-05.2 — Publicação de releases com versionamento semântico

| Critério | Situação | Evidência no código |
|---|---:|---|
| Tags seguem SemVer | Implementado | `v[0-9]+.[0-9]+.[0-9]+` e validação `^v[0-9]+\.[0-9]+\.[0-9]+$` no release workflow |
| Workflow de release por tag | Implementado | `.github/workflows/release.yml` |
| Binários nomeados por plataforma | Implementado | `assinatura-${VERSION}-linux-amd64.AppImage`, `.exe`, `.dmg` |
| Publicação automática no GitHub Releases | Implementado | `softprops/action-gh-release@v2` |
| Convenção `assinatura-<versão>-<os>-<arch>` | Implementado | nomes gerados no workflow de release |

Exemplos esperados de artefatos:

```text
assinatura-1.0.5-linux-amd64.AppImage
assinatura-1.0.5-windows-amd64.exe
assinatura-1.0.5-macos-amd64.dmg
```

Além do `assinatura`, a implementação também publica o CLI `simulador` e o `assinador.jar`, porque eles fazem parte da execução completa do Sistema Runner.

## US-05.3 — Checksums SHA256 e assinatura de artefatos com Cosign

Esta seção cobre checksums SHA256, assinatura de artefatos com Cosign, identidade OIDC, transparency log, arquivos `.sig`, `.pem` e `.bundle`, e documentação com `cosign verify-blob`.

| Critério | Situação | Evidência no código |
|---|---:|---|
| `checksums.txt` por release | Implementado | `sha256sum * > checksums.txt` |
| Assinatura com Cosign | Implementado | `sigstore/cosign-installer@v3` e `cosign sign-blob` |
| Identidade OIDC | Implementado | `permissions: id-token: write` |
| Transparency log | Implementado | `--tlog-upload=true` |
| `.sig` por artefato | Implementado | `--output-signature "$f.sig"` |
| `.pem` por artefato | Implementado | `--output-certificate "$f.pem"` |
| `.bundle` por artefato | Implementado adicionalmente | `--bundle "$f.bundle"` |
| Documentação de verificação | Implementado | `docs/integridade-assinatura-artefatos.md` |

Exemplo de verificação com Cosign:

```bash
cosign verify-blob \
  --certificate assinatura-1.0.5-linux-amd64.AppImage.pem \
  --signature assinatura-1.0.5-linux-amd64.AppImage.sig \
  assinatura-1.0.5-linux-amd64.AppImage
```

Resultado esperado:

```text
Verified OK
```

## Resumo de entrega da Sprint 1

A Sprint 1 está coberta por código, workflow e documentação:

```text
US-01.1: CLI Go estruturado, com Cobra e comando assinatura version.
US-05.1: CI/CD com build multiplataforma e artifacts do workflow.
US-05.2: Releases SemVer com binários por plataforma no GitHub Releases.
US-05.3: Checksums SHA256, Cosign OIDC, transparency log, .sig, .pem e documentação de verificação.
```
