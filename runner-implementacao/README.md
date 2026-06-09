# Sistema Runner

Implementação de referência do Sistema Runner para a disciplina de Implementação e Integração de Software. O projeto fornece dois CLIs multiplataforma em Go e uma aplicação Java executável (`assinador.jar`).

## 1. Visão geral

O Runner facilita a execução e integração de aplicações Java por linha de comando. Ele contém:

- `assinatura`: CLI para criar e validar assinaturas digitais simuladas por meio do `assinador.jar`.
- `simulador`: CLI para consultar, iniciar e parar o Simulador do HubSaúde.
- `assinador.jar`: aplicação Java 21 que valida parâmetros, simula assinatura, simula validação e expõe endpoints HTTP.

O Simulador do HubSaúde não é reimplementado neste repositório. Ele é tratado como sistema externo, normalmente disponível em `https://localhost:8443/`.

## 2. Tecnologias

- Go 1.23+ para os CLIs.
- Java 21 para o `assinador.jar`.
- HTTP local para o modo servidor.
- Diretório de estado: `~/.hubsaude/`.

## 3. Estrutura de pastas

```text
runner/
├── cmd/
│   ├── assinatura/
│   └── simulador/
├── internal/
│   ├── apperrors/
│   ├── assinador/
│   ├── httpx/
│   ├── jdk/
│   ├── paths/
│   ├── process/
│   └── simulator/
├── assinador/
│   ├── Makefile
│   └── src/
├── docs/adr/
├── examples/
└── .github/workflows/
```

## 4. Pré-requisitos

No Ubuntu 20.04 ou superior:

```bash
sudo apt update
sudo apt install -y git make
```

Instale Go e Java 21. Verifique:

```bash
go version
java -version
```

## 5. Build

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

## 6. Testes

Execute os testes Go:

```bash
go test ./...
```

Execute o teste local do Java:

```bash
cd assinador
make test
cd ..
```

## 7. Uso do CLI assinatura

### 7.1 Modo local

```bash
echo "documento de teste" > examples/documento.txt

go run ./cmd/assinatura sign \
  --local \
  --jar assinador/target/assinador.jar \
  --input examples/documento.txt \
  --signer "Maria Runner" \
  --output examples/assinatura.json

go run ./cmd/assinatura validate \
  --local \
  --jar assinador/target/assinador.jar \
  --signature examples/assinatura.json
```

### 7.2 Modo servidor

Inicie o servidor do assinador:

```bash
go run ./cmd/assinatura server start \
  --jar assinador/target/assinador.jar \
  --port 8080 \
  --idle-timeout-minutes 10
```

Consulte o status:

```bash
go run ./cmd/assinatura server status --port 8080
```

Assine usando HTTP:

```bash
go run ./cmd/assinatura sign \
  --input examples/documento.txt \
  --signer "Maria Runner" \
  --output examples/assinatura.json
```

Valide usando HTTP:

```bash
go run ./cmd/assinatura validate \
  --signature examples/assinatura.json
```

Encerre o servidor:

```bash
go run ./cmd/assinatura server stop --port 8080
```

## 8. Uso do CLI simulador

Como o Simulador HubSaúde informado no trabalho executa em HTTPS local com certificado autoassinado, use `--insecure` em ambiente de desenvolvimento.

Consultar status:

```bash
go run ./cmd/simulador status --url https://localhost:8443 --insecure
```

Consultar `/api/info`:

```bash
go run ./cmd/simulador info --url https://localhost:8443 --insecure
```

Parar o simulador por `/shutdown`:

```bash
go run ./cmd/simulador stop --url https://localhost:8443 --insecure
```

Iniciar um `simulador.jar` local, caso ele não esteja em execução:

```bash
go run ./cmd/simulador start \
  --url https://localhost:8443 \
  --insecure \
  --jar /caminho/para/simulador.jar
```

## 9. Variáveis de ambiente

| Variável | Finalidade |
|---|---|
| `RUNNER_HOME` | Altera o diretório gerenciado. Padrão: `~/.hubsaude`. |
| `RUNNER_JAVA` | Informa o caminho exato do executável `java`. |
| `RUNNER_ASSINADOR_JAR` | Caminho padrão do `assinador.jar`. |
| `HUBSAUDE_SIMULATOR_URL` | URL padrão do Simulador do HubSaúde. |

## 10. Códigos de saída

| Código | Significado |
|---:|---|
| 0 | Sucesso |
| 1 | Erro de uso do usuário |
| 2 | Dependência ausente |
| 3 | Falha de integração HTTP ou subprocesso |
| 4 | Erro interno inesperado |

## 11. Observações de qualidade

- A validação de negócio da assinatura fica no `assinador.jar`.
- O CLI apenas coleta parâmetros, invoca o modo local ou HTTP e formata a saída.
- `stdout` é reservado para resultados; `stderr` é usado para diagnóstico.
- O estado de processos é salvo em `~/.hubsaude/processos/`.
- Logs de processos iniciados pelo Runner ficam em `~/.hubsaude/logs/`.

## 12. Limitações conhecidas

- O provisionamento automático real do JDK/JRE está preparado por estrutura, mas não baixa artefatos automaticamente nesta versão entregue.
- A integração PKCS#11 real não foi implementada; o `assinador.jar` contém simulação de assinatura e validação, conforme escopo acadêmico.
- O CLI `simulador` depende do Simulador HubSaúde real para os comandos `status`, `info` e `stop`.
