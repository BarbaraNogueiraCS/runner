# Runner

Este repositório usa a estrutura:

```text
runner/
├── .github/workflows/        # GitHub Actions na raiz do repositório
├── .gitignore                # ignorados do repositório
└── runner-implementacao/     # código-fonte, go.mod, docs, scripts e assinador Java
```

A pasta `.github` precisa ficar nesta raiz para que o GitHub Actions reconheça os workflows `build.yml` e `release.yml`. Como o código está em `runner-implementacao`, os workflows usam `working-directory: runner-implementacao`.

Para executar localmente:

```bash
cd runner-implementacao
go test ./...
go vet ./...
./scripts/check-release-artifacts.sh
cd assinador && make clean all test && cd ..
```

Para gerar release no GitHub:

```bash
git checkout main
git pull origin main
git tag v1.0.1
git push origin v1.0.1
```

O workflow de release gera binários, checksums e assinaturas Cosign (`.sig`, `.pem` e `.bundle`) e publica tudo em GitHub Releases.
