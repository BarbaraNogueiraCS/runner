# Higiene do repositório

Este projeto mantém os arquivos de configuração do repositório na raiz:

- `.github/workflows/`: workflows do GitHub Actions.
- `.gitignore`: regras para ignorar saídas geradas.
- `.gitattributes`: normalização de finais de linha e tratamento de binários.
- `docs/`: documentação do projeto.

## Arquivos gerados que não devem ser versionados

As seguintes pastas são criadas por testes, builds ou execução local:

```text
runner-implementacao/assinador/target/
runner-implementacao/assinador/out/
runner-implementacao/dist/
runner-implementacao/examples/*.json
```

Esses caminhos estão no `.gitignore`. Eles podem existir localmente depois de `make clean all test`, `go build` ou execução de assinatura, mas não devem aparecer em commits.

## Como verificar antes do commit

```bash
cd runner-implementacao
./scripts/check-generated-files.sh
```

Resultado esperado:

```text
OK: target/, out/, dist/ e JSONs gerados estão protegidos pelo .gitignore e não aparecem como arquivos rastreados.
```

## Como remover do índice se foram adicionados por engano

Execute na raiz do repositório:

```bash
git rm -r --cached runner-implementacao/assinador/target || true
git rm -r --cached runner-implementacao/assinador/out || true
git rm -r --cached runner-implementacao/dist || true
git rm -r --cached runner-implementacao/examples/*.json || true
```

Depois faça um novo `git status` e commite apenas código-fonte, documentação, scripts e workflows.
