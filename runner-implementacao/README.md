# Sistema Runner

Implementação de referência do Sistema Runner para a disciplina de Implementação e Integração de Software. O projeto fornece dois CLIs multiplataforma em Go e uma aplicação Java executável (`assinador.jar`).

## 1. Visão geral

O Runner facilita a execução e integração de aplicações Java por linha de comando. Ele contém:

- `assinatura`: CLI para criar e validar assinaturas digitais simuladas por meio do `assinador.jar`.
- `simulador`: CLI para consultar, iniciar e parar o Simulador do HubSaúde.
- `assinador.jar`: aplicação Java 21 que valida parâmetros, simula assinatura, simula validação e expõe endpoints HTTP.

O Simulador do HubSaúde não é reimplementado neste repositório. Ele é tratado como sistema externo, normalmente disponível em `https://localhost:8443/`.

## 2. Conformidade com os guias SES-GO

A criação e a validação de assinatura foram reimplementadas para seguir a estrutura dos guias **Criar assinatura digital** e **Validar assinatura digital** da SES-GO, respeitando o escopo acadêmico de simulação definido para o trabalho.

A implementação atual:

- recebe `Bundle` FHIR R4 em JSON;
- recebe `Provenance` FHIR R4 em JSON;
- valida `Provenance.target` com referências `urn:uuid:`;
- valida timestamp Unix UTC no intervalo definido pelo guia;
- valida URI versionada da política de assinatura;
- recebe cadeia de certificados em array JSON com certificados DER codificados em base64;
- gera estrutura JWS JSON Serialization;
- usa `protected header` com `alg`, `x5c`, `sigPId` e `iat`;
- usa `payload` como hash SHA-256 codificado em base64Url;
- encapsula o JWS em `Signature.data` usando base64 padrão FHIR;
- valida `Signature.data`, JWS, `protected header`, política, timestamp, assinatura simulada e integridade opcional do conteúdo.

Limitação assumida: a operação criptográfica real, OCSP, CRL, TSA real e PKCS#11 real continuam simulados, porque a própria especificação acadêmica do Runner define que assinatura e validação criptográfica reais estão fora do escopo.

## 3. Tecnologias

- Go 1.23+ para os CLIs.
- Java 21 para o `assinador.jar`.
- HTTP local para o modo servidor.
- Diretório de estado: `~/.hubsaude/`.

## 4. Estrutura de pastas

```text
runner/                         # raiz do repositório Git/GitHub
├── .github/workflows/           # GitHub Actions: build.yml e release.yml
├── .gitignore                   # regras de arquivos gerados/temporários
└── runner-implementacao/        # raiz do módulo Go e do código-fonte
    ├── cmd/
    │   ├── assinatura/
    │   └── simulador/
    ├── internal/
    │   ├── apperrors/
    │   ├── assinador/
    │   ├── httpx/
    │   ├── jdk/
    │   ├── release/
    │   ├── paths/
    │   ├── process/
    │   └── simulador/
    ├── assinador/
    │   ├── Makefile
    │   └── src/
    ├── docs/adr/
    └── examples/
```

A pasta `.github` e o arquivo `.gitignore` ficam na raiz do repositório porque o GitHub Actions só reconhece workflows em `.github/workflows` a partir da raiz do repositório Git. Como o código está em `runner-implementacao`, os workflows usam `working-directory: runner-implementacao`.

## 5. Pré-requisitos

No Ubuntu 20.04 ou superior:

```bash
sudo apt update
sudo apt install -y git make
```

Instale Go e Java 21. Verifique:

```bash
go version
java -version
javac -version
keytool -help >/dev/null
```

O `keytool` vem com o JDK e é usado pelo `make test` para gerar um certificado de teste.

## 6. Build

Compile os CLIs:

```bash
go build ./...
go build -o dist/assinatura ./cmd/assinatura
go build -o dist/simulador ./cmd/simulador
```

Compile o `assinador.jar`:

```bash
cd assinador
make all
cd ..
```

O artefato será gerado em:

```text
assinador/target/assinador.jar
```

## 7. Testes

Execute os testes Go:

```bash
go test ./...
go vet ./...
```

Execute o teste local do Java:

```bash
cd assinador
make clean all test
cd ..
```

O `make test` gera arquivos de amostra em `assinador/target/`, cria uma assinatura FHIR `Signature` e valida o resultado.

## 8. Uso do CLI assinatura conforme guia SES-GO

Primeiro gere os arquivos de amostra:

```bash
cd assinador
make samples
cd ..
```

Guarde o timestamp gerado:

```bash
TS=$(cat assinador/target/timestamp.txt)
POLICY='https://fhir.saude.go.gov.br/r4/seguranca/ImplementationGuide/br.go.ses.seguranca|0.1.2'
```

### 8.1 Modo local

Criar assinatura:

```bash
go run ./cmd/assinatura sign \
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
  --output examples/assinatura.json
```

Validar assinatura:

```bash
go run ./cmd/assinatura validate \
  --local \
  --jar assinador/target/assinador.jar \
  --signature examples/assinatura.json \
  --timestamp "$TS" \
  --policy "$POLICY" \
  --bundle assinador/target/bundle.json \
  --provenance assinador/target/provenance.json \
  --config assinador/target/config.json
```

### 8.2 Modo servidor

Inicie o servidor do assinador:

```bash
go run ./cmd/assinatura server start \
  --jar assinador/target/assinador.jar \
  --port 8080 \
  --idle-timeout-minutes 10
```

Crie assinatura via HTTP:

```bash
go run ./cmd/assinatura sign \
  --bundle assinador/target/bundle.json \
  --provenance assinador/target/provenance.json \
  --crypto-material assinador/target/crypto.json \
  --cert-chain assinador/target/certs.json \
  --timestamp "$TS" \
  --policy "$POLICY" \
  --config assinador/target/config.json \
  --signer "Maria Runner" \
  --output examples/assinatura-servidor.json \
  --port 8080
```

Valide via HTTP:

```bash
go run ./cmd/assinatura validate \
  --signature examples/assinatura-servidor.json \
  --timestamp "$TS" \
  --policy "$POLICY" \
  --bundle assinador/target/bundle.json \
  --provenance assinador/target/provenance.json \
  --config assinador/target/config.json \
  --port 8080
```

Encerre o servidor:

```bash
go run ./cmd/assinatura server stop --port 8080
```

## 9. Compatibilidade com comando antigo

O comando antigo continua disponível para testes simples:

```bash
echo "documento de teste" > examples/documento.txt

go run ./cmd/assinatura sign \
  --local \
  --jar assinador/target/assinador.jar \
  --input examples/documento.txt \
  --signer "Maria Runner" \
  --output examples/assinatura-legada.json
```

Esse modo gera uma assinatura simulada simples encapsulada no novo formato `Signature`, mas o fluxo recomendado é o fluxo com `Bundle`, `Provenance`, `cert-chain`, `timestamp` e `policy`.

## 10. Uso do CLI simulador

Como o Simulador HubSaúde informado no trabalho executa em HTTPS local com certificado autoassinado, use `--insecure` em ambiente de desenvolvimento.

Consultar status:

```bash
go run ./cmd/simulador status --url https://localhost:8443 --insecure
```

Consultar `/api/info`:

```bash
go run ./cmd/simulador info --url https://localhost:8443 --insecure
```

Iniciar o simulador. Se ele já estiver ativo em `https://localhost:8443`, o comando reaproveita a instância existente. Se não estiver ativo, verifica a porta 8443 antes de iniciar. Quando `--jar` não é informado, baixa dinamicamente o artefato indicado pelo `release.json`, sem repetir o download quando a versão local já está atualizada:

```bash
go run ./cmd/simulador start --url https://localhost:8443 --insecure
```

Por padrão, o arquivo de manifesto usado é:

```text
https://raw.githubusercontent.com/kyriosdata/runner/main/release.json
```

Também é possível indicar explicitamente o artefato desejado do manifesto:

```bash
go run ./cmd/simulador start \
  --url https://localhost:8443 \
  --insecure \
  --artifact validador
```

Parar o simulador por `/shutdown`:

```bash
go run ./cmd/simulador stop --url https://localhost:8443 --insecure
```

## 11. Variáveis de ambiente

| Variável | Finalidade |
|---|---|
| `RUNNER_HOME` | Altera o diretório gerenciado. Padrão: `~/.hubsaude`. |
| `RUNNER_JAVA` | Informa o caminho exato do executável `java`. |
| `RUNNER_ASSINADOR_JAR` | Caminho padrão do `assinador.jar`. |
| `HUBSAUDE_SIMULADOR_URL` | URL padrão do Simulador do HubSaúde. |
| `RUNNER_RELEASE_JSON` | URL do `release.json` usado para baixar simulador/validador e runtime Java. |
| `RUNNER_SIMULADOR_ARTIFACT` | Artefato a baixar para o CLI `simulador`: `simulador` ou `validador`. |

## 12. Códigos de saída

| Código | Significado |
|---:|---|
| 0 | Sucesso |
| 1 | Erro de uso do usuário |
| 2 | Dependência ausente |
| 3 | Falha de integração HTTP ou subprocesso |
| 4 | Erro interno inesperado |

## 13. Observações de qualidade

- A validação de negócio da assinatura fica no `assinador.jar`.
- O CLI apenas coleta parâmetros, invoca o modo local ou HTTP e formata a saída.
- `stdout` é reservado para resultados; `stderr` é usado para diagnóstico.
- O estado de processos é salvo em `~/.hubsaude/processos/`.
- Logs de processos iniciados pelo Runner ficam em `~/.hubsaude/logs/`.

## 14. Provisionamento automático

A implementação possui provisionamento automático de runtime Java e artefatos do HubSaúde:

- `internal/jdk` detecta Java 21+ em `RUNNER_JAVA`, em `~/.hubsaude/jdk/bin/java` ou no `PATH`.
- Se Java 21+ não for encontrado, o Runner baixa um JRE/JDK Temurin 21 a partir da URL da plataforma definida no `release.json`.
- `internal/release` busca o `release.json`, compara versão instalada localmente e baixa apenas quando ausente, desatualizado ou com checksum incompatível.
- O artefato baixado é armazenado em `~/.hubsaude/<artefato>/` com arquivo `.version` para rastreabilidade.
- O checksum SHA256 do manifesto é verificado quando informado.

## 15. Releases e assinatura de artefatos

O workflow `.github/workflows/release.yml`, localizado na raiz do repositório, publica binários pré-compilados para as três plataformas no GitHub Releases quando uma tag `v*` é criada. Para a tag `v1.0.0`, os nomes esperados são:

```text
assinatura-1.0.0-windows-amd64.exe
assinatura-1.0.0-linux-amd64.AppImage
assinatura-1.0.0-macos-amd64.dmg
simulador-1.0.0-windows-amd64.exe
simulador-1.0.0-linux-amd64.AppImage
simulador-1.0.0-macos-amd64.dmg
```

O workflow também gera `checksums.txt` e assina todos os artefatos com Cosign em modo keyless/OIDC, com envio explícito ao transparency log do Sigstore (`--tlog-upload=true`). Para cada arquivo assinado são publicados `.sig`, `.pem` e `.bundle`.

## 16. Limitações conhecidas

- A integração PKCS#11 real não foi implementada; o `assinador.jar` contém simulação de assinatura e validação, conforme escopo acadêmico.
- A validação de cadeia ICP-Brasil, OCSP, CRL e TSA real foi substituída por validação estrutural e simulação controlada.
- O CLI `simulador` depende do Simulador HubSaúde real para os comandos `status`, `info` e `stop`.

## Artefatos executáveis de release

A geração dos binários de Windows, Linux e macOS, os checksums SHA256, a assinatura Cosign e a publicação no GitHub Releases estão documentados em [`docs/artefatos-executaveis.md`](docs/artefatos-executaveis.md). A política de integridade, os arquivos obrigatórios `<artefato>.sig` e `<artefato>.pem`, o uso de OIDC/transparency log e os comandos de verificação estão em [`docs/integridade-assinatura-artefatos.md`](docs/integridade-assinatura-artefatos.md). O workflow responsável é `.github/workflows/release.yml` na raiz do repositório Git. Como o módulo Go está em `runner-implementacao`, o workflow usa `working-directory: runner-implementacao`.
