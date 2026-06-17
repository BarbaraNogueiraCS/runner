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

```bash
cd runner-implementacao
go test ./...
go vet ./...
./scripts/check-release-artifacts.sh
cd assinador && make clean all test && cd ..
```

## Release

Para gerar release no GitHub, faça a tag a partir da `main`:

```bash
git checkout main
git pull origin main
git tag v1.0.2
git push origin v1.0.2
```

O workflow de release gera binários, checksums e assinaturas Cosign (`.sig`, `.pem` e `.bundle`) e publica tudo em GitHub Releases. A documentação fica em `docs/`, também na raiz do repositório.
