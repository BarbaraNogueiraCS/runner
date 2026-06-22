# Runner

Sistema Runner com estrutura de projeto Go na raiz do repositório. Esta organização evita a pasta intermediária a raiz do repositório e deixa o projeto no formato mais comum para aplicações Go com CLI.

```text
runner/
├── go.mod
├── go.sum
├── Makefile
├── cmd/
│   ├── assinatura/
│   └── simulador/
├── internal/
├── assinador/
├── scripts/
├── examples/
├── docs/
├── projetos/
├── .github/
├── .gitignore
├── .gitattributes
└── README.md
```

## Visão geral

O projeto entrega:

- `assinatura`: CLI em Go para criar e validar assinaturas digitais simuladas.
- `simulador`: CLI em Go para consultar, iniciar e parar o Simulador do HubSaúde.
- `assinador.jar`: componente Java 21 que valida parâmetros, simula assinatura, simula validação, expõe endpoints HTTP `/sign` e `/validate` e pode ser invocado localmente ou em modo servidor pelo CLI.
- GitHub Actions para CI, build multiplataforma, GitHub Releases, checksums SHA256 e assinatura Cosign.

O Simulador do HubSaúde é tratado como sistema externo. O comando `simulador` gerencia/consulta esse serviço, normalmente em `https://localhost:8443/`.

## Pré-requisitos locais

No Ubuntu 20.04 ou superior:

```bash
sudo apt update
sudo apt install -y git make
```

Verifique Go e Java:

```bash
go version
java -version
javac -version
keytool -help >/dev/null
```

O projeto usa Go 1.23.2 e Java 21.

## Automação local

A automação fica no `Makefile` da raiz:

```bash
make help
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
make samples    gera arquivos de exemplo em assinador/target/
make clean      remove dist/, assinador/out/, assinador/target/ e examples/*.json
make all        executa deps, test, vet, cover, check, java-test e build
```

Fluxo local recomendado:

```bash
make deps
make test
make vet
make cover
make check
make java-test
make samples
make build
```

Ou tudo de uma vez:

```bash
make all
```

## Arquivos que não devem ser commitados

As pastas abaixo são geradas por build, teste ou execução local e não devem ser versionadas:

```text
assinador/target/
assinador/out/
dist/
examples/*.json
```

Antes de commitar, rode:

```bash
make check
```

Se algum arquivo gerado já tiver sido adicionado ao Git por engano, remova do índice, sem apagar necessariamente o arquivo local:

```bash
git rm -r --cached --ignore-unmatch assinador/target
git rm -r --cached --ignore-unmatch assinador/out
git rm -r --cached --ignore-unmatch dist
git rm --cached --ignore-unmatch examples/*.json
```

## Assinatura local

Gere o `assinador.jar` e os arquivos de exemplo:

```bash
make java-test
make samples
make build
```

Assine localmente:

```bash
TS=$(date -u +%s)
POLICY='https://fhir.saude.go.gov.br/r4/seguranca/ImplementationGuide/br.go.ses.seguranca|0.1.2'
mkdir -p examples

./dist/assinatura sign \
  --local \
  --jar assinador/target/assinador.jar \
  --bundle assinador/target/bundle.json \
  --provenance assinador/target/provenance.json \
  --crypto-material assinador/target/crypto.json \
  --cert-chain assinador/target/certs.json \
  --timestamp "$TS" \
  --policy "$POLICY" \
  --config assinador/target/config.json \
  --signer "Maria Runner" \
  --output examples/assinatura-local.json
```

Valide localmente:

```bash
./dist/assinatura validate \
  --local \
  --jar assinador/target/assinador.jar \
  --signature examples/assinatura-local.json \
  --timestamp "$TS" \
  --policy "$POLICY" \
  --bundle assinador/target/bundle.json \
  --provenance assinador/target/provenance.json \
  --config assinador/target/config.json
```

## Assinador em modo servidor

Depois de `make java-test` e `make build`, inicie o assinador HTTP:

```bash
./dist/assinatura start   --jar assinador/target/assinador.jar   --port 8080   --timeout 10
```

Consulte status:

```bash
./dist/assinatura status --port 8080
```

Pare o servidor:

```bash
./dist/assinatura stop --port 8080
```

Os comandos `sign` e `validate` usam o servidor por padrão quando ele está ativo. Use `--local` para forçar a execução direta com `java -jar`.

## PKCS#11 / token ou smart card

O material criptográfico aceita `type` igual a `TOKEN` ou `SMARTCARD`. Para testes sem dispositivo real, use `simulation=true`. Para dispositivo real, informe `pin` e `pkcs11Library`. Consulte `docs/sprint3-modo-servidor-pkcs11.md`.

## GitHub Actions e release

Os workflows ficam em `.github/workflows/` na raiz:

- `build.yml`: roda em `main`, `develop`, `feature/**`, `release/**`, `refactor/**` e pull requests.
- `release.yml`: roda quando uma tag SemVer, por exemplo `v1.0.5`, é enviada.

Como sua release atual é `v1.0.4`, a próxima versão desta refatoração deve ser:

```bash
git tag v1.0.5
git push origin v1.0.5
```

A release publica binários multiplataforma, `checksums.txt` e assinaturas Cosign (`.sig`, `.pem` e `.bundle`).

## Documentação

- `docs/sprint1-fundacao-entrega-continua.md`
- `docs/sprint2-assinatura-digital-local.md`
- `docs/sprint3-modo-servidor-pkcs11.md`
- `docs/integridade-assinatura-artefatos.md`
- `docs/higiene-repositorio.md`
- `docs/implementacao.md`


## Observação sobre `internal/release`

A pasta `internal/release/` é código-fonte interno do projeto Go e deve ser versionada. O `.gitignore` usa `/release/` com barra inicial para ignorar apenas uma pasta local de release na raiz, sem ignorar `internal/release/`. Antes de enviar para o GitHub, rode:

```bash
./scripts/check-internal-packages.sh
git ls-files internal/release/manifest.go
```

O segundo comando deve imprimir `internal/release/manifest.go`.
