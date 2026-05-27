# Sistema Runner

## 1. Visão Geral

O **Sistema Runner** é uma implementação acadêmica de referência criada para facilitar a execução e o gerenciamento de aplicações Java por meio de comandos de terminal. O sistema foi desenvolvido no contexto da disciplina de **Implementação e Integração de Software**, com aplicação relacionada à Plataforma HubSaúde, voltada à interoperabilidade de dados em saúde.

A proposta central do Sistema Runner é reduzir a complexidade de uso de aplicações Java. Em vez de exigir que o usuário conheça comandos como `java -jar`, portas de execução, controle de processos, localização de arquivos `.jar`, endpoints HTTP e dependências de ambiente, o sistema oferece comandos CLI mais simples e padronizados.

A solução é composta principalmente por três partes:

| Parte | Tecnologia | Descrição |
|---|---|---|
| **`assinatura`** | Go + Cobra | CLI responsável por criar e validar assinaturas simuladas, invocando o `assinador.jar` localmente ou via HTTP. |
| **`assinador.jar`** | Java 21 | Aplicação Java responsável por validar parâmetros e simular operações de criação e validação de assinatura digital. |
| **`simulador`** | Go + Cobra | CLI responsável por iniciar, parar e consultar o status do Simulador do HubSaúde. |

O sistema utiliza também um diretório local gerenciado, normalmente em `~/.hubsaude/`, para armazenar arquivos `.jar`, cache, metadados, registros de processos, logs e artefatos auxiliares.

> **Importante:** esta implementação não realiza assinatura digital criptográfica real. A criação e a validação de assinaturas são simuladas, com finalidade acadêmica e de integração.

---

## 2. Objetivo do Projeto

O objetivo do projeto é desenvolver uma ferramenta de linha de comandos capaz de facilitar a execução, integração e gerenciamento de aplicações Java relacionadas ao HubSaúde.

De forma mais específica, o projeto tem os seguintes objetivos:

- fornecer comandos simples para criação e validação simulada de assinatura digital;
- permitir a execução do `assinador.jar` em modo local, por meio de `java -jar`;
- permitir a execução do `assinador.jar` em modo servidor HTTP;
- expor endpoints HTTP para assinatura e validação simuladas;
- controlar o ciclo de vida do `assinador.jar`, incluindo início, verificação de saúde e parada;
- fornecer um CLI separado para gerenciamento do Simulador do HubSaúde;
- registrar processos em execução no diretório local `~/.hubsaude/`;
- verificar disponibilidade de portas antes de iniciar servidores;
- estruturar download e verificação de integridade de artefatos;
- apoiar testes, documentação, CI/CD e distribuição multiplataforma;
- aplicar boas práticas de Engenharia de Software, como modularidade, rastreabilidade, tratamento de erros, testes e documentação.

---

## 3. Problema Resolvido

Sem o Sistema Runner, o usuário precisaria conhecer detalhes técnicos para executar aplicações Java, como:

- instalar e configurar Java;
- localizar arquivos `.jar`;
- executar comandos `java -jar`;
- informar parâmetros técnicos corretamente;
- saber qual porta utilizar;
- verificar se uma aplicação já está em execução;
- encerrar processos manualmente;
- chamar endpoints HTTP diretamente;
- interpretar mensagens técnicas de erro.

O Sistema Runner resolve esse problema ao oferecer uma interface de terminal mais simples. Assim, o usuário pode executar comandos como:

```bash
./bin/assinatura sign
./bin/assinatura validate
./bin/assinatura start
./bin/assinatura stop
./bin/simulador start
./bin/simulador status
./bin/simulador stop
```

Com isso, o sistema oculta parte da complexidade técnica e permite que o usuário execute os fluxos principais de forma mais padronizada.

---

## 4. Funcionalidades

### 4.1 Funcionalidades do CLI `assinatura`

| Funcionalidade | Comando | Descrição |
|---|---|---|
| Exibir versão | `assinatura version` | Mostra a versão atual do CLI. |
| Criar assinatura simulada | `assinatura sign` | Solicita ao `assinador.jar` a criação de uma assinatura simulada. |
| Validar assinatura simulada | `assinatura validate` | Solicita ao `assinador.jar` a validação simulada de uma assinatura. |
| Iniciar assinador em modo servidor | `assinatura start` | Inicia o `assinador.jar` em modo servidor HTTP. |
| Parar assinador | `assinatura stop` | Encerra o `assinador.jar` em execução na porta informada. |
| Forçar modo local | `--local` | Executa a operação diretamente por `java -jar`. |
| Definir JAR | `--jar` | Informa o caminho do `assinador.jar`. |
| Definir porta | `--port` | Informa a porta usada pelo servidor do assinador. |

### 4.2 Funcionalidades do `assinador.jar`

| Funcionalidade | Descrição |
|---|---|
| Modo CLI | Permite executar `sign` e `validate` diretamente por `java -jar`. |
| Modo servidor HTTP | Mantém o assinador em execução para receber requisições HTTP. |
| Validação de parâmetros | Verifica se os campos obrigatórios foram informados. |
| Simulação de assinatura | Retorna uma assinatura fictícia para entradas válidas. |
| Simulação de validação | Retorna resultado simulado de validade da assinatura. |
| Endpoint `/health` | Indica se o servidor está em execução. |
| Endpoint `/sign` | Recebe requisição HTTP para criação de assinatura simulada. |
| Endpoint `/validate` | Recebe requisição HTTP para validação simulada. |
| Endpoint `/shutdown` | Solicita encerramento do servidor. |
| Adaptador PKCS#11 | Estrutura a integração futura com token ou smart card, sem assinatura real nesta versão. |

### 4.3 Funcionalidades do CLI `simulador`

| Funcionalidade | Comando | Descrição |
|---|---|---|
| Iniciar simulador | `simulador start` | Inicia o `simulador.jar`. |
| Consultar status | `simulador status` | Consulta o status do processo do simulador. |
| Parar simulador | `simulador stop` | Encerra o simulador. |
| Definir JAR local | `--jar` | Informa o caminho local do `simulador.jar`. |
| Definir origem remota | `--source` | Informa uma URL para download do `simulador.jar`. |
| Definir porta | `--port` | Informa a porta do simulador. |
| Verificar artefato | `--sha256` | Permite conferir a integridade do arquivo baixado. |

---

## 5. Tecnologias Utilizadas

| Tecnologia | Uso no projeto |
|---|---|
| **Go** | Desenvolvimento dos CLIs `assinatura` e `simulador`. |
| **Cobra** | Estruturação de comandos, subcomandos e flags dos CLIs. |
| **Java 21** | Desenvolvimento e execução do `assinador.jar`. |
| **Maven** | Build e execução dos testes do projeto Java. |
| **Javac/Jar** | Alternativa de build Java por script local. |
| **HTTP** | Comunicação entre o CLI `assinatura` e o `assinador.jar` em modo servidor. |
| **Sistema de arquivos local** | Persistência de cache, metadados e registros em `~/.hubsaude/`. |
| **SHA256** | Verificação de integridade de artefatos. |
| **Cosign / Sigstore** | Assinatura e verificação de artefatos de release. |
| **GitHub Actions** | Automação de build, testes e releases. |
| **GitHub Releases** | Publicação de binários, arquivos `.jar`, checksums e assinaturas. |

---

## 6. Pré-requisitos

Para executar o projeto em uma máquina local, recomenda-se instalar:

| Dependência | Versão recomendada | Finalidade |
|---|---|---|
| **Git** | Atual | Clonar e versionar o projeto. |
| **Go** | Go 1.23+ para execução local; Go 1.25 como versão-alvo documentada | Compilar os CLIs. |
| **Java JDK** | Java 21 | Compilar e executar o `assinador.jar`. |
| **Maven** | 3.9+ recomendado | Compilar e testar o projeto Java. |
| **Bash** | Compatível | Executar scripts auxiliares em Linux/macOS. |
| **Curl ou wget** | Atual | Apoiar downloads e testes HTTP. |
| **Cosign** | Opcional | Verificar ou assinar artefatos. |

No Ubuntu, as dependências básicas podem ser instaladas com:

```bash
sudo apt update
sudo apt install -y git curl wget unzip zip build-essential maven
```

O Java 21 pode ser instalado via Temurin/Adoptium ou outro provedor compatível.

---

## 7. Instalação

### 7.1 Obter o projeto

Clone o repositório ou extraia o arquivo `.zip` do projeto.

```bash
git clone <url-do-repositorio>
cd sistema-runner-atividade
```

Caso esteja usando o pacote `.zip`:

```bash
unzip sistema_runner_implementacao.zip
cd sistema-runner-atividade
```

### 7.2 Preparar o diretório local do Runner

```bash
mkdir -p ~/.hubsaude/assinador
mkdir -p ~/.hubsaude/simulador
mkdir -p ~/.hubsaude/processos
mkdir -p ~/.hubsaude/cache
mkdir -p ~/.hubsaude/logs
```

### 7.3 Compilar o `assinador.jar` com Maven

```bash
cd projetos/assinador-java
mvn clean package
cd ../..
```

O arquivo esperado será gerado em:

```text
projetos/assinador-java/target/assinador.jar
```

Copie o JAR para o diretório local gerenciado:

```bash
cp projetos/assinador-java/target/assinador.jar ~/.hubsaude/assinador/assinador.jar
```

### 7.4 Alternativa: compilar o `assinador.jar` por script

Se Maven não estiver disponível:

```bash
chmod +x scripts/build-assinador.sh
./scripts/build-assinador.sh
```

Depois, localize o arquivo gerado:

```bash
find . -name "assinador.jar"
```

### 7.5 Compilar os CLIs Go

```bash
mkdir -p bin
go mod tidy
go build -o bin/assinatura ./cmd/assinatura
go build -o bin/simulador ./cmd/simulador
```

Verifique se os binários foram gerados:

```bash
ls -l bin
```

---

## 8. Configuração

### 8.1 Diretório local gerenciado

Por padrão, o Sistema Runner usa:

```text
~/.hubsaude/
```

Estrutura esperada:

```text
~/.hubsaude/
├── assinador/
│   └── assinador.jar
├── simulador/
│   └── simulador.jar
├── processos/
├── cache/
└── logs/
```

### 8.2 Variáveis de ambiente opcionais

| Variável | Finalidade | Exemplo |
|---|---|---|
| `JAVA_HOME` | Indicar instalação Java 21. | `/usr/lib/jvm/java-21` |
| `HUBSAUDE_HOME` | Sobrescrever o diretório padrão `~/.hubsaude/`. | `/home/usuario/.hubsaude` |
| `PATH` | Permitir localizar `go`, `java`, `javac` e `mvn`. | Inclui caminhos do Go e Java. |

### 8.3 Configuração de portas

| Serviço | Porta padrão |
|---|---:|
| `assinador.jar` em modo servidor | `8080` |
| `simulador.jar` | `8443` |

As portas podem ser alteradas com a flag `--port`.

### 8.4 Configuração dos arquivos `.jar`

O `assinador.jar` pode ser informado com:

```bash
--jar ~/.hubsaude/assinador/assinador.jar
```

O `simulador.jar` pode ser informado com:

```bash
--jar ~/.hubsaude/simulador/simulador.jar
```

Também é possível informar uma origem remota para o simulador:

```bash
--source https://exemplo.com/simulador.jar
```

---

## 9. Execução

### 9.1 Verificar versão

```bash
./bin/assinatura version
```

### 9.2 Criar assinatura simulada em modo local

```bash
./bin/assinatura sign \
  --local \
  --jar ~/.hubsaude/assinador/assinador.jar \
  --documento documento.json \
  --certificado certificado.pem
```

### 9.3 Validar assinatura simulada em modo local

```bash
./bin/assinatura validate \
  --local \
  --jar ~/.hubsaude/assinador/assinador.jar \
  --documento documento.json \
  --assinatura assinatura-simulada-123
```

### 9.4 Iniciar o assinador em modo servidor

```bash
./bin/assinatura start \
  --jar ~/.hubsaude/assinador/assinador.jar \
  --port 8080
```

### 9.5 Verificar servidor do assinador

```bash
curl http://localhost:8080/health
```

Resposta esperada:

```json
{"status":"UP"}
```

### 9.6 Criar assinatura via HTTP

```bash
./bin/assinatura sign \
  --documento documento.json \
  --certificado certificado.pem \
  --port 8080
```

### 9.7 Validar assinatura via HTTP

```bash
./bin/assinatura validate \
  --documento documento.json \
  --assinatura assinatura-simulada-123 \
  --port 8080
```

### 9.8 Parar o servidor do assinador

```bash
./bin/assinatura stop --port 8080
```

### 9.9 Iniciar o simulador

```bash
./bin/simulador start \
  --jar ~/.hubsaude/simulador/simulador.jar \
  --port 8443
```

### 9.10 Consultar status do simulador

```bash
./bin/simulador status --port 8443
```

### 9.11 Parar o simulador

```bash
./bin/simulador stop --port 8443
```

---

## 10. Exemplos de Uso

### 10.1 Criar arquivos de teste

Crie um documento fictício:

```bash
cat > documento.json <<'EOF_DOC'
{
  "id": "doc-001",
  "tipo": "teste",
  "conteudo": "Documento de teste para assinatura simulada."
}
EOF_DOC
```

Crie um certificado fictício:

```bash
cat > certificado.pem <<'EOF_CERT'
-----BEGIN CERTIFICATE-----
CERTIFICADO-SIMULADO-PARA-TESTE
-----END CERTIFICATE-----
EOF_CERT
```

### 10.2 Fluxo completo em modo local

```bash
./bin/assinatura sign \
  --local \
  --jar ~/.hubsaude/assinador/assinador.jar \
  --documento documento.json \
  --certificado certificado.pem

./bin/assinatura validate \
  --local \
  --jar ~/.hubsaude/assinador/assinador.jar \
  --documento documento.json \
  --assinatura assinatura-simulada-123
```

### 10.3 Fluxo completo em modo servidor

```bash
./bin/assinatura start \
  --jar ~/.hubsaude/assinador/assinador.jar \
  --port 8080

./bin/assinatura sign \
  --documento documento.json \
  --certificado certificado.pem \
  --port 8080

./bin/assinatura validate \
  --documento documento.json \
  --assinatura assinatura-simulada-123 \
  --port 8080

./bin/assinatura stop --port 8080
```

### 10.4 Testar endpoints diretamente

```bash
curl http://localhost:8080/health
```

```bash
curl -X POST http://localhost:8080/sign \
  -H "Content-Type: application/json" \
  -d '{"document":"documento.json","certificate":"certificado.pem"}'
```

```bash
curl -X POST http://localhost:8080/validate \
  -H "Content-Type: application/json" \
  -d '{"document":"documento.json","signature":"assinatura-simulada-123"}'
```

---

## 11. Estrutura de Pastas

```text
sistema-runner-atividade/
├── README.md
├── go.mod
├── release.json
├── cmd/
│   ├── assinatura/
│   │   └── main.go
│   └── simulador/
│       └── main.go
├── internal/
│   ├── appinfo/
│   │   └── version.go
│   ├── config/
│   │   └── paths.go
│   ├── errors/
│   │   └── app_error.go
│   ├── cli/
│   │   ├── formatter/
│   │   ├── assinatura/
│   │   └── simulador/
│   ├── assinatura/
│   │   ├── dto/
│   │   └── usecase/
│   ├── simulador/
│   │   ├── dto/
│   │   └── usecase/
│   ├── java/
│   │   ├── runtime/
│   │   └── invoker/
│   ├── process/
│   └── artifacts/
├── projetos/
│   └── assinador-java/
│       ├── pom.xml
│       └── src/
├── third_party/
│   └── cobra/
├── scripts/
│   ├── build-assinador.sh
│   ├── install-local.sh
│   └── verify-artifact.sh
└── .github/
    └── workflows/
        ├── build.yml
        └── release.yml
```

### 11.1 Descrição das principais pastas

| Pasta | Descrição |
|---|---|
| `cmd/` | Contém os pontos de entrada dos executáveis `assinatura` e `simulador`. |
| `internal/` | Contém a lógica interna dos CLIs em Go. |
| `internal/cli/` | Define comandos, subcomandos, flags e formatação de saída. |
| `internal/assinatura/` | Contém DTOs e casos de uso relacionados à assinatura simulada. |
| `internal/simulador/` | Contém DTOs e casos de uso relacionados ao simulador. |
| `internal/java/` | Resolve o Java instalado e invoca aplicações `.jar` localmente ou via HTTP. |
| `internal/process/` | Gerencia portas, processos, health check e registros locais. |
| `internal/artifacts/` | Realiza download, leitura de metadados e verificação SHA256 de artefatos. |
| `projetos/assinador-java/` | Contém a aplicação Java `assinador.jar`. |
| `third_party/cobra/` | Contém uma implementação mínima local do Cobra para uso acadêmico/offline. |
| `scripts/` | Contém scripts auxiliares de build, instalação e verificação. |
| `.github/workflows/` | Contém workflows de CI/CD para build e release. |

---

## 12. Testes

### 12.1 Executar testes Go

Na raiz do projeto:

```bash
go test ./...
```

Os testes Go validam partes como:

- verificação de porta;
- cálculo e conferência de checksum SHA256;
- comportamento de módulos internos;
- estrutura dos comandos e casos de uso.

### 12.2 Executar testes Java

```bash
cd projetos/assinador-java
mvn test
cd ../..
```

Os testes Java validam:

- validação de parâmetros;
- serviço de assinatura simulada;
- serviço de validação simulada.

### 12.3 Testes de integração recomendados

| Teste | Comando / ação | Resultado esperado |
|---|---|---|
| Versão do CLI | `./bin/assinatura version` | Exibir versão atual. |
| Assinatura local | `assinatura sign --local ...` | Retornar assinatura simulada. |
| Validação local | `assinatura validate --local ...` | Retornar validação simulada. |
| Start do servidor | `assinatura start --port 8080` | Iniciar `assinador.jar` e registrar processo. |
| Health check | `curl http://localhost:8080/health` | Retornar status do servidor. |
| Assinatura HTTP | `assinatura sign --port 8080 ...` | Usar endpoint `/sign`. |
| Validação HTTP | `assinatura validate --port 8080 ...` | Usar endpoint `/validate`. |
| Stop do servidor | `assinatura stop --port 8080` | Encerrar processo e atualizar registro. |
| Status do simulador | `simulador status --port 8443` | Informar situação do simulador. |

---

## 13. Documentação do Projeto

Este README está relacionado aos documentos técnicos elaborados para o Sistema Runner:

| Documento | Finalidade |
|---|---|
| **Especificação de Requisitos de Software** | Define o que o sistema deve fazer, seus requisitos, regras, histórias de usuário e critérios de aceitação. |
| **Documento de Arquitetura de Software** | Descreve a organização arquitetural, decisões técnicas, atributos de qualidade e restrições. |
| **Documento de Projeto Detalhado de Software** | Detalha módulos, classes, métodos, dados, interfaces, fluxos e testes. |
| **Documento de Modelo C4** | Representa a arquitetura em níveis de contexto, contêineres, componentes e código. |
| **Documento de Implementação e Integração** | Explica como o sistema foi implementado, integrado, configurado e executado. |
| **Documento de Teste de Software** | Planeja e registra os testes, casos de teste, métricas, defeitos e evidências. |
| **Documento de Implantação de Software** | Define ambiente, instalação, execução, rollback, segurança e aceite da implantação. |
| **Documento de Gerenciamento de Projeto com Cronograma** | Organiza o desenvolvimento em entregas, tarefas, artefatos, riscos e indicadores. |

---

## 14. Limitações Conhecidas

| Limitação | Descrição |
|---|---|
| Assinatura simulada | O sistema não realiza assinatura digital criptográfica real. |
| Validação simulada | O sistema não valida criptograficamente assinaturas reais. |
| Sem validade jurídica | O sistema não deve ser usado como solução oficial de assinatura digital. |
| PKCS#11 parcial | O adaptador está estruturado, mas não executa assinatura real com token ou smart card. |
| Provisionamento Java parcial | A estrutura existe, mas a execução prioriza Java 21 instalado ou já disponível localmente. |
| `simulador.jar` externo | O `simulador.jar` real pode não acompanhar o pacote e deve ser informado por `--jar` ou `--source`. |
| Sem banco de dados | A persistência é feita por arquivos locais em `~/.hubsaude/`. |
| Sem interface gráfica | O sistema é utilizado exclusivamente por terminal. |
| Sem autenticação | Não há controle de usuários, pois o uso previsto é local. |
| JSON simplificado | A implementação acadêmica usa tratamento simples de JSON no Java. |
| Suporte inicial amd64 | A distribuição multiplataforma considera inicialmente Windows, Linux e macOS em amd64. |

---

## 15. Equipe

| Papel | Responsabilidade |
|---|---|
| **Equipe do projeto** | Desenvolvimento, integração, documentação, testes e implantação do Sistema Runner. |
| **Disciplina** | Implementação e Integração de Software. |
| **Instituição** | Bacharelado em Engenharia de Software — Universidade Federal de Goiás (UFG). |
| **Contexto de aplicação** | Plataforma HubSaúde — interoperabilidade de dados em saúde. |
| **Responsáveis técnicos** | Bárbara Nogueira Carvalho da Silveira - Estudantes responsáveis pela implementação, integração, testes, documentação e entrega acadêmica. |
| **Avaliador / professor** | Fábio Lucena - Responsável por avaliar a aderência do projeto aos requisitos da disciplina e às boas práticas de engenharia de software. |

---

## 16. Licença ou Finalidade Acadêmica

Este projeto possui **finalidade acadêmica** e foi desenvolvido como uma implementação de referência para a disciplina de **Implementação e Integração de Software**.

A solução demonstra práticas de:

- especificação de requisitos;
- arquitetura de software;
- Modelo C4;
- projeto detalhado;
- implementação modular;
- integração entre Go e Java;
- integração local e HTTP;
- gerenciamento de processos;
- tratamento de erros;
- testes de software;
- CI/CD;
- verificação de artefatos;
- documentação técnica;
- implantação local.

O Sistema Runner **não deve ser utilizado como solução de assinatura digital real ou juridicamente válida**, pois as operações de assinatura e validação são simuladas.

Caso seja necessário definir uma licença formal, recomenda-se que a equipe inclua um arquivo `LICENSE` no repositório, conforme a política da disciplina, da instituição ou do projeto.
