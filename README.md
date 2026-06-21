# Runner

Este repositório usa a estrutura:

```text
runner/
├── .github/workflows/        # GitHub Actions na raiz do repositório
├── .gitignore                # regras para não versionar saídas geradas
├── .gitattributes            # normalização de finais de linha
├── docs/                     # documentação do projeto na raiz
└── runner-implementacao/     # código-fonte, go.mod, scripts e assinador Java
```

A pasta `.github` precisa ficar na raiz para que o GitHub Actions reconheça os workflows `build.yml` e `release.yml`. Como o código está em `runner-implementacao`, os workflows usam `working-directory: runner-implementacao`.

## Arquivos que não devem ser commitados

As pastas abaixo são geradas por build, teste ou execução local e não devem ser versionadas:

```text
runner-implementacao/assinador/target/
runner-implementacao/assinador/out/
runner-implementacao/dist/
runner-implementacao/examples/*.json
```

O `.gitignore` da raiz protege esses caminhos. Antes de commitar, execute:

```bash
cd runner-implementacao
./scripts/check-generated-files.sh
```

Se algum arquivo gerado já tiver sido adicionado ao Git por engano, remova apenas do índice, mantendo-o no computador:

```bash
git rm -r --cached runner-implementacao/assinador/target || true
git rm -r --cached runner-implementacao/assinador/out || true
git rm -r --cached runner-implementacao/dist || true
git rm -r --cached runner-implementacao/examples/*.json || true
```

## Execução local rápida

A automação local fica no `runner-implementacao/Makefile`. Ele evita repetir manualmente comandos como `go mod download`, `go test`, `go vet` e `go build`.

```bash
cd runner-implementacao
make deps
make test
make vet
make cover
make check
make java-test
make build
```

Para executar o fluxo local completo de uma vez:

```bash
cd runner-implementacao
make all
```

Principais alvos:

```text
make deps       baixa dependências Go com go mod download
make tidy       organiza go.mod e go.sum com go mod tidy
make test       roda go test ./...
make vet        roda go vet ./...
make cover      roda go test -cover ./...
make check      valida higiene do repositório e artefatos de release
make build      compila dist/assinatura e dist/simulador
make java-test  compila e testa o assinador.jar
make clean      remove dist/, assinador/out/, assinador/target/ e examples/*.json
make all        executa deps, test, vet, cover, check, java-test e build
```

## Release

Para gerar release no GitHub, faça a tag a partir da `main`:

```bash
git checkout main
git pull origin main
git tag v1.0.5
git push origin v1.0.5
```

O workflow de release gera binários, checksums e assinaturas Cosign (`.sig`, `.pem` e `.bundle`) e publica tudo em GitHub Releases. A documentação fica em `docs/`, também na raiz do repositório.


## Sprint 1 — Fundação e Entrega Contínua

A Sprint 1 está rastreada em `docs/sprint1-fundacao-entrega-continua.md`. O CLI usa Cobra (`github.com/spf13/cobra`) e mantém a identidade do repositório no módulo Go:

```go
module github.com/BarbaraNogueiraCS/runner
```

Para instalar a ferramenta de scaffolding do Cobra em ambiente de desenvolvimento:

```bash
go install github.com/spf13/cobra-cli@latest
```

O usuário final não precisa instalar `cobra-cli`; os binários são gerados e publicados automaticamente pelo GitHub Actions em GitHub Releases.
