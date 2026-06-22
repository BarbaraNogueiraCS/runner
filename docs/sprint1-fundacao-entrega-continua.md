# Sprint 1 — Fundação e Entrega Contínua

## US-01.1 — Estrutura base do CLI em Go

| Critério | Status | Evidência |
|---|---|---|
| Projeto Go inicializado | Implementado | `go.mod` na raiz |
| CLI com Cobra | Implementado | `cmd/assinatura/main.go`, `cmd/simulador/main.go`, dependência `github.com/spf13/cobra` |
| Comando `assinatura version` | Implementado | `./dist/assinatura version` |
| Estrutura documentada | Implementado | `README.md`, `docs/criterios-aceitacao.md` |
| Build Windows/Linux/macOS | Implementado | `.github/workflows/build.yml`, `.github/workflows/release.yml` |

O usuário final não precisa instalar `cobra-cli`; ele é útil apenas para scaffolding em desenvolvimento. O projeto usa a biblioteca Cobra diretamente nos comandos.

## US-05.1 — Pipeline CI/CD multiplataforma

| Critério | Status | Evidência |
|---|---|---|
| GitHub Actions configurado | Implementado | `.github/workflows/build.yml` |
| Cross-compilation | Implementado | `GOOS=linux`, `GOOS=windows`, `GOOS=darwin` |
| Build a cada push | Implementado | branches `main`, `develop`, `feature/**` |
| Artifacts do workflow | Implementado | `actions/upload-artifact@v4`, `runner-dev-binaries`, artifacts do workflow |

## US-05.2 — Publicação com versionamento semântico

| Critério | Status | Evidência |
|---|---|---|
| Tags SemVer | Implementado | trigger `v[0-9]+.[0-9]+.[0-9]+` e validação `^v[0-9]+\.[0-9]+\.[0-9]+$` |
| Binários por plataforma | Implementado | release Linux, Windows e macOS |
| GitHub Releases automático | Implementado | `softprops/action-gh-release` |
| Convenção de nomes | Implementado | `assinatura-<versão>-<os>-<arch>` e `simulador-<versão>-<os>-<arch>` |

## US-05.3 — Checksums SHA256 e Cosign

| Critério | Status | Evidência |
|---|---|---|
| `checksums.txt` | Implementado | `sha256sum * > checksums.txt` |
| Cosign keyless/OIDC | Implementado | `id-token: write`, `cosign sign-blob` |
| Transparency log | Implementado | `--tlog-upload=true` |
| `.sig`, `.pem`, `.bundle` | Implementado | workflow gera os três arquivos por artefato |
| Documentação de verificação | Implementado | `docs/integridade-assinatura-artefatos.md` com `cosign verify-blob` |

## Automação local

O `Makefile` na raiz automatiza:

```text
make deps
make tidy
make test
make vet
make cover
make check
make build
make java-test
make samples
make clean
make all
```

Observação de rastreabilidade: a Sprint 1 exige checksums SHA256 e Cosign para verificação com cosign verify-blob.


Observação de CI: o workflow de build também roda em branches `release/**` e `refactor/**`, além de `main`, `develop` e `feature/**`, para apoiar o fluxo local de release/refatoração usado no projeto.
