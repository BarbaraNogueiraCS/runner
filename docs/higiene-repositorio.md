# Higiene do repositório

O repositório mantém apenas código-fonte, documentação e configuração de build. Saídas geradas localmente não devem ser versionadas.

Não commitar:

```text
assinador/target/
assinador/out/
dist/
examples/*.json
```

Esses caminhos são protegidos pelo `.gitignore` da raiz.

## Verificação

Execute na raiz do repositório:

```bash
make check
```

ou diretamente:

```bash
./scripts/check-generated-files.sh
```

Resultado esperado:

```text
OK: target/, out/, dist/ e JSONs gerados estão protegidos pelo .gitignore e não aparecem como arquivos rastreados.
```

Se algum arquivo gerado já estiver rastreado, remova apenas do índice:

```bash
git rm -r --cached --ignore-unmatch assinador/target
git rm -r --cached --ignore-unmatch assinador/out
git rm -r --cached --ignore-unmatch dist
git rm --cached --ignore-unmatch examples/*.json
```
