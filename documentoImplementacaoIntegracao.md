# Documento de Implementação e Integração de Software — Sistema Runner

## 1. Identificação do documento

| Campo | Informação |
|---|---|
| **Nome do sistema** | Sistema Runner |
| **Nome do documento** | Documento de Implementação e Integração de Software |
| **Versão do documento** | 1.0 |
| **Data de elaboração** | 07/05/2026 |
| **Responsável pela elaboração** | Equipe do projeto / Disciplina de Implementação e Integração de Software |
| **Instituição / contexto acadêmico** | Bacharelado em Engenharia de Software — Universidade Federal de Goiás (UFG) |
| **Contexto de aplicação** | Plataforma HubSaúde — interoperabilidade de dados em saúde |
| **Documentos relacionados** | Especificação de Requisitos de Software; Documento de Arquitetura de Software; Documento de Projeto Detalhado de Software; Documento de Modelo C4; Plano Revisado #2; README da implementação |
| **Tipo de sistema** | Ferramenta de linha de comandos, integração com aplicações Java, gerenciamento de processos, provisionamento de ambiente Java e simulação de assinatura digital |

---

## 2. Histórico de versões

| Versão | Data | Autor / Responsável | Descrição da alteração |
|---|---|---|---|
| 1.0 | 07/05/2026 | Equipe do projeto | Elaboração inicial do Documento de Implementação e Integração de Software do Sistema Runner, com base nos arquivos fornecidos, documentos de requisitos, arquitetura, projeto detalhado, Modelo C4 e implementação de referência acadêmica. |

---

## 3. Sumário

1. [Identificação do documento](#1-identificação-do-documento)  
2. [Histórico de versões](#2-histórico-de-versões)  
3. [Sumário](#3-sumário)  
4. [Objetivo do documento](#4-objetivo-do-documento)  
5. [Visão geral do sistema](#5-visão-geral-do-sistema)  
6. [Escopo da implementação](#6-escopo-da-implementação)  
7. [Tecnologias utilizadas](#7-tecnologias-utilizadas)  
8. [Organização de pastas e arquivos](#8-organização-de-pastas-e-arquivos)  
9. [Arquitetura da implementação](#9-arquitetura-da-implementação)  
10. [Módulos implementados](#10-módulos-implementados)  
11. [Funcionalidades implementadas](#11-funcionalidades-implementadas)  
12. [Modelo de dados](#12-modelo-de-dados)  
13. [Integração entre as partes do sistema](#13-integração-entre-as-partes-do-sistema)  
14. [APIs e endpoints implementados](#14-apis-e-endpoints-implementados)  
15. [Integração com banco de dados](#15-integração-com-banco-de-dados)  
16. [Configuração do ambiente](#16-configuração-do-ambiente)  
17. [Variáveis de ambiente](#17-variáveis-de-ambiente)  
18. [Procedimento de instalação](#18-procedimento-de-instalação)  
19. [Procedimento de execução](#19-procedimento-de-execução)  
20. [Testes realizados](#20-testes-realizados)  
21. [Tratamento de erros](#21-tratamento-de-erros)  
22. [Segurança implementada](#22-segurança-implementada)  
23. [Controle de versão](#23-controle-de-versão)  
24. [Evidências da implementação e integração](#24-evidências-da-implementação-e-integração)  
25. [Limitações e melhorias futuras](#25-limitações-e-melhorias-futuras)  
26. [Referências](#26-referências)  

---

## 4. Objetivo do documento

O objetivo deste documento é descrever a **implementação e integração de software do Sistema Runner**, registrando como a solução foi construída, organizada, integrada, testada e preparada para execução.

Este documento tem como finalidade:

- apresentar a estrutura implementada do Sistema Runner;
- documentar as tecnologias, módulos, pastas e arquivos utilizados;
- explicar como ocorre a integração entre os CLIs em Go e as aplicações Java;
- registrar as funcionalidades implementadas nos componentes `assinatura`, `simulador` e `assinador.jar`;
- demonstrar como o sistema trata erros, persistência local, segurança, endpoints e execução;
- apoiar a avaliação acadêmica da atividade de Implementação e Integração de Software;
- servir como material de consulta para instalação, execução, testes e manutenção futura.

O documento foi elaborado considerando os requisitos, a arquitetura, o projeto detalhado, o Modelo C4 e a implementação de referência acadêmica do Sistema Runner.

---

## 5. Visão geral do sistema

O **Sistema Runner** é uma solução de software criada para facilitar a execução e o gerenciamento de aplicações Java por meio de comandos de terminal. A proposta central é permitir que usuários, estudantes, integradores e operadores técnicos executem aplicações associadas à Plataforma HubSaúde sem precisar conhecer detalhes internos de instalação do Java, execução de arquivos `.jar`, configuração de portas, controle de processos ou comunicação HTTP.

A solução é composta por três partes principais:

| Parte | Tecnologia | Finalidade |
|---|---|---|
| **`assinatura`** | Go + Cobra | CLI responsável por criar e validar assinaturas simuladas, invocando o `assinador.jar` localmente ou via HTTP. |
| **`simulador`** | Go + Cobra | CLI responsável por iniciar, parar e consultar status do Simulador do HubSaúde. |
| **`assinador.jar`** | Java 21 | Aplicação Java responsável por validar parâmetros e simular criação/validação de assinatura digital. |

Além desses elementos, o sistema utiliza o diretório local `~/.hubsaude/` para armazenar artefatos, cache, registros de processos, metadados e dependências. A implementação também possui estrutura de CI/CD para builds e releases multiplataforma, incluindo verificação de integridade por SHA256 e assinatura de artefatos com Cosign.

O sistema não realiza assinatura digital criptográfica real. A assinatura e a validação são **simuladas**, conforme o escopo definido nos documentos do projeto. O foco da implementação é a integração, a execução controlada das aplicações Java, a validação de parâmetros, a qualidade do código, os testes e a distribuição dos artefatos.

---

## 6. Escopo da implementação

### 6.1 Itens implementados

A implementação de referência contempla os seguintes itens:

| Item | Situação | Descrição |
|---|---|---|
| CLI `assinatura` | Implementado | CLI em Go para operações de assinatura e validação simuladas. |
| CLI `simulador` | Implementado | CLI em Go para gerenciamento do Simulador do HubSaúde. |
| `assinador.jar` | Implementado | Aplicação Java 21 com modo CLI e modo servidor HTTP. |
| Comando `assinatura version` | Implementado | Exibe a versão atual do CLI. |
| Comando `assinatura sign` | Implementado | Cria assinatura simulada. |
| Comando `assinatura validate` | Implementado | Valida assinatura simulada. |
| Comando `assinatura start` | Implementado | Inicia o `assinador.jar` em modo servidor. |
| Comando `assinatura stop` | Implementado | Interrompe o `assinador.jar` em modo servidor. |
| Comando `simulador start` | Implementado | Inicia o Simulador do HubSaúde a partir de `.jar` local ou obtido por URL. |
| Comando `simulador stop` | Implementado | Solicita parada do simulador. |
| Comando `simulador status` | Implementado | Consulta o status do simulador. |
| Invocação local via `java -jar` | Implementado | Permite chamar o `assinador.jar` diretamente. |
| Invocação HTTP do `assinador.jar` | Implementado | Usa endpoints `/sign`, `/validate` e `/health`. |
| Validação preliminar no Go | Implementado | Verifica parâmetros essenciais antes da invocação. |
| Validação no Java | Implementado | Valida parâmetros recebidos pelo `assinador.jar`. |
| Registro local de processos | Implementado | Armazena metadados em `~/.hubsaude/processos/`. |
| Verificação de porta | Implementado | Verifica disponibilidade de porta antes de iniciar processos. |
| Download genérico de artefatos | Implementado | Permite obter artefatos remotos por URL. |
| Verificação SHA256 | Implementado | Calcula e valida integridade de arquivos. |
| Estrutura CI/CD | Implementado | Workflows de build e release em GitHub Actions. |
| Testes unitários iniciais | Implementado | Testes Go e Java para componentes centrais. |

### 6.2 Itens parcialmente implementados ou estruturados

| Item | Situação | Observação |
|---|---|---|
| Provisionamento automático completo de JDK/JRE | Estruturado | A implementação prioriza o uso de Java 21 já instalado ou disponível no diretório gerenciado. O fluxo completo de download automático pode ser expandido. |
| Integração PKCS#11 | Estruturada | Existe adaptador Java isolado, mas não há assinatura criptográfica real nesta versão. |
| Assinatura real de documentos | Fora do escopo | As operações são simuladas, conforme requisitos. |
| `simulador.jar` real | Não incluso | O CLI permite informar `--jar` ou `--source` para indicar origem do simulador. |

### 6.3 Itens fora do escopo

Não fazem parte desta implementação:

- assinatura digital criptográfica real;
- validação criptográfica real de assinatura;
- integração real com autoridades certificadoras;
- geração de certificados digitais;
- autenticação de usuários;
- armazenamento persistente de assinaturas;
- interface gráfica;
- banco de dados relacional ou não relacional;
- implantação distribuída em produção.

---

## 7. Tecnologias utilizadas

| Tecnologia | Uso na implementação | Justificativa |
|---|---|---|
| **Go** | Implementação dos CLIs `assinatura` e `simulador`. | Facilita a criação de executáveis multiplataforma. |
| **Cobra** | Estruturação dos comandos CLI. | Permite comandos, subcomandos, flags e ajuda integrada. |
| **Java 21** | Implementação do `assinador.jar`. | Atende à restrição técnica do projeto. |
| **Maven** | Build e testes do projeto Java. | Padroniza empacotamento e execução de testes. |
| **Javac/Jar** | Alternativa de build local do `assinador.jar`. | Permite compilar o Java sem Maven, por script auxiliar. |
| **HTTP** | Integração em modo servidor entre CLI e `assinador.jar`. | Comunicação simples e adequada para execução local. |
| **Sistema de arquivos local** | Persistência em `~/.hubsaude/`. | Armazena cache, processos e metadados sem banco de dados. |
| **SHA256** | Verificação de integridade de artefatos. | Garante que arquivos baixados não foram alterados ou corrompidos. |
| **Cosign / Sigstore** | Assinatura de artefatos de release. | Melhora segurança da cadeia de suprimentos. |
| **GitHub Actions** | CI/CD. | Automatiza testes, build e publicação de releases. |
| **GitHub Releases** | Distribuição dos artefatos. | Disponibiliza binários versionados aos usuários. |
| **PKCS#11** | Adaptador de integração criptográfica. | Padrão para comunicação com token/smart card, estruturado nesta versão. |

---

## 8. Organização de pastas e arquivos

A implementação foi organizada para separar pontos de entrada, lógica interna, integração com Java, gerenciamento de processos, artefatos, testes, documentação e automações.

```text
sistema-runner-atividade/
  README.md
  go.mod

  docs/
    relatorio_implementacao.md

  cmd/
    assinatura/
      main.go
    simulador/
      main.go

  internal/
    appinfo/
      version.go

    config/
      paths.go

    errors/
      app_error.go

    cli/
      formatter/
        formatter.go
      assinatura/
        commands/
          root.go
          version.go
          sign.go
          validate.go
          start_stop.go
      simulador/
        commands/
          root.go
          start_stop_status.go

    assinatura/
      dto/
        dto.go
      usecase/
        usecase.go

    simulador/
      dto/
        dto.go
      usecase/
        usecase.go

    java/
      runtime/
        resolver.go
      invoker/
        local.go
        http.go

    process/
      metadata.go
      registry.go
      port.go
      health.go
      manager.go
      port_test.go

    artifacts/
      checksum.go
      downloader.go
      release.go
      checksum_test.go

  projetos/
    assinador-java/
      pom.xml
      src/
        main/java/br/ufg/hubsaude/assinador/
          Main.java
          cli/
            CliAdapter.java
          controller/
            SignatureHttpServer.java
          dto/
            SignRequest.java
            SignResponse.java
            ValidateRequest.java
            ValidateResponse.java
            ErrorResponse.java
          validation/
            ValidationResult.java
            ParameterValidator.java
          service/
            SignatureService.java
            FakeSignatureService.java
          pkcs11/
            PKCS11Adapter.java
          error/
            InvalidParameterException.java
          util/
            JsonUtil.java
        test/java/br/ufg/hubsaude/assinador/
          validation/
            ParameterValidatorTest.java
          service/
            FakeSignatureServiceTest.java

  third_party/
    cobra/
      go.mod
      command.go

  scripts/
    build-assinador.sh
    install-local.sh
    verify-artifact.sh

  .github/
    workflows/
      build.yml
      release.yml

  release.json
```

### 8.1 Justificativa da organização

| Diretório | Justificativa |
|---|---|
| `cmd/` | Contém os pontos de entrada dos executáveis Go. |
| `internal/` | Contém módulos internos do sistema, evitando acoplamento externo indevido. |
| `internal/cli/` | Agrupa comandos e formatação da interface de linha de comandos. |
| `internal/java/` | Centraliza resolução de Java e invocação de aplicações `.jar`. |
| `internal/process/` | Centraliza controle de portas, processos, health check e registros locais. |
| `internal/artifacts/` | Contém download, metadados de release e verificação SHA256. |
| `projetos/assinador-java/` | Contém a aplicação Java `assinador.jar`. |
| `third_party/cobra/` | Implementação mínima local de Cobra para permitir execução acadêmica/offline. |
| `scripts/` | Scripts auxiliares de build, instalação e verificação. |
| `.github/workflows/` | Automação de build e release. |

---

## 9. Arquitetura da implementação

A implementação segue uma arquitetura modular orientada a comandos, com separação entre camada CLI, casos de uso, infraestrutura, persistência local e aplicação Java.

### 9.1 Visão em camadas

| Camada | Responsabilidade | Exemplos na implementação |
|---|---|---|
| **Interface CLI** | Receber comandos e parâmetros do usuário. | `internal/cli/assinatura/commands`, `internal/cli/simulador/commands`. |
| **Casos de uso** | Orquestrar os fluxos de assinatura, validação, start, stop e status. | `internal/assinatura/usecase`, `internal/simulador/usecase`. |
| **Infraestrutura Java** | Resolver Java e invocar `.jar` local ou via HTTP. | `internal/java/runtime`, `internal/java/invoker`. |
| **Infraestrutura de processo** | Controlar portas, PID, health check e encerramento. | `internal/process`. |
| **Infraestrutura de artefatos** | Baixar artefatos e verificar integridade. | `internal/artifacts`. |
| **Persistência local** | Armazenar metadados em arquivos locais. | `~/.hubsaude/`, `ProcessRegistry`. |
| **Aplicação Java** | Validar parâmetros e simular assinatura/validação. | `projetos/assinador-java`. |
| **Automação** | Compilar, testar e publicar artefatos. | GitHub Actions. |

### 9.2 Estilo arquitetural aplicado

A solução aplica os seguintes estilos:

- **arquitetura modular orientada a comandos**, adequada a aplicações CLI;
- **arquitetura em camadas**, separando interface, caso de uso e infraestrutura;
- **cliente-servidor local**, quando o CLI `assinatura` se comunica com `assinador.jar` via HTTP;
- **portas e adaptadores**, pois a comunicação com Java, HTTP, sistema de arquivos, PKCS#11 e downloads externos é isolada em módulos próprios.

---

## 10. Módulos implementados

### 10.1 Módulos Go

| Módulo | Arquivos principais | Responsabilidade |
|---|---|---|
| `cmd/assinatura` | `main.go` | Inicializa o CLI `assinatura`. |
| `cmd/simulador` | `main.go` | Inicializa o CLI `simulador`. |
| `internal/appinfo` | `version.go` | Mantém informações de versão. |
| `internal/config` | `paths.go` | Define caminhos locais, especialmente `~/.hubsaude/`. |
| `internal/errors` | `app_error.go` | Define erros padronizados da aplicação. |
| `internal/cli/formatter` | `formatter.go` | Formata respostas e erros para o terminal. |
| `internal/cli/assinatura/commands` | `root.go`, `sign.go`, `validate.go`, `start_stop.go`, `version.go` | Implementa comandos do CLI `assinatura`. |
| `internal/cli/simulador/commands` | `root.go`, `start_stop_status.go` | Implementa comandos do CLI `simulador`. |
| `internal/assinatura/dto` | `dto.go` | Define estruturas de entrada e saída das operações de assinatura. |
| `internal/assinatura/usecase` | `usecase.go` | Implementa fluxos de assinatura, validação e controle do assinador. |
| `internal/simulador/dto` | `dto.go` | Define estruturas de status/configuração do simulador. |
| `internal/simulador/usecase` | `usecase.go` | Implementa início, parada e status do simulador. |
| `internal/java/runtime` | `resolver.go` | Resolve o executável Java a ser usado. |
| `internal/java/invoker` | `local.go`, `http.go` | Invoca o `assinador.jar` via processo local ou HTTP. |
| `internal/process` | `metadata.go`, `registry.go`, `port.go`, `health.go`, `manager.go` | Gerencia processos, portas, registros e health check. |
| `internal/artifacts` | `checksum.go`, `downloader.go`, `release.go` | Gerencia download, release e verificação SHA256. |

### 10.2 Módulos Java

| Módulo / pacote | Arquivos principais | Responsabilidade |
|---|---|---|
| `br.ufg.hubsaude.assinador` | `Main.java` | Ponto de entrada do `assinador.jar`. |
| `cli` | `CliAdapter.java` | Recebe comandos quando o `.jar` é executado localmente. |
| `controller` | `SignatureHttpServer.java` | Expõe endpoints HTTP do assinador. |
| `dto` | `SignRequest`, `SignResponse`, `ValidateRequest`, `ValidateResponse`, `ErrorResponse` | Define objetos de entrada e saída. |
| `validation` | `ParameterValidator`, `ValidationResult` | Valida parâmetros de assinatura e validação. |
| `service` | `SignatureService`, `FakeSignatureService` | Define e implementa a lógica simulada. |
| `pkcs11` | `PKCS11Adapter` | Isola integração futura com PKCS#11. |
| `error` | `InvalidParameterException` | Representa erros de parâmetro. |
| `util` | `JsonUtil` | Apoia serialização e respostas JSON. |

---

## 11. Funcionalidades implementadas

### 11.1 Funcionalidades do CLI `assinatura`

| Funcionalidade | Comando | Descrição |
|---|---|---|
| Exibir versão | `assinatura version` | Mostra a versão atual do CLI. |
| Criar assinatura simulada | `assinatura sign` | Envia dados ao `assinador.jar` para gerar assinatura simulada. |
| Validar assinatura simulada | `assinatura validate` | Envia dados ao `assinador.jar` para validação simulada. |
| Iniciar assinador em servidor | `assinatura start` | Executa o `assinador.jar` em modo servidor HTTP. |
| Parar assinador | `assinatura stop` | Encerra o processo do `assinador.jar` registrado. |
| Invocação local | `--local` | Força execução via `java -jar`. |
| Invocação HTTP | padrão quando servidor ativo | Usa `/sign`, `/validate` e `/health`. |
| Indicação de JAR | `--jar` | Permite informar caminho do `assinador.jar`. |
| Porta customizada | `--port` | Permite definir porta do servidor. |

### 11.2 Funcionalidades do CLI `simulador`

| Funcionalidade | Comando | Descrição |
|---|---|---|
| Iniciar simulador | `simulador start` | Inicia o `simulador.jar` local ou obtido por URL. |
| Parar simulador | `simulador stop` | Solicita encerramento do processo/serviço. |
| Consultar status | `simulador status` | Consulta status via registro local e endpoint, quando disponível. |
| Usar JAR local | `--jar` | Permite informar caminho do `simulador.jar`. |
| Baixar por URL | `--source` | Permite indicar URL alternativa para download do simulador. |
| Porta customizada | `--port` | Permite usar porta diferente da padrão 8443. |

### 11.3 Funcionalidades do `assinador.jar`

| Funcionalidade | Descrição |
|---|---|
| Modo CLI | Permite executar operações quando chamado por `java -jar`. |
| Modo servidor | Permite iniciar servidor HTTP local. |
| Endpoint `/sign` | Recebe requisições de criação de assinatura simulada. |
| Endpoint `/validate` | Recebe requisições de validação simulada. |
| Endpoint `/health` | Permite verificar se o servidor está ativo. |
| Endpoint `/shutdown` | Permite solicitar encerramento controlado do servidor. |
| Validação de parâmetros | Rejeita requisições incompletas ou inválidas. |
| Simulação de assinatura | Retorna assinatura simulada para entradas válidas. |
| Simulação de validação | Retorna resultado simulado de validação. |
| Adaptador PKCS#11 | Mantém estrutura para integração futura com dispositivo criptográfico. |

---

## 12. Modelo de dados

O Sistema Runner não possui modelo de dados de domínio complexo nem banco de dados tradicional. Os dados são representados por estruturas em memória, objetos de requisição/resposta e arquivos locais de metadados.

### 12.1 Dados principais manipulados

| Dado | Local / estrutura | Descrição |
|---|---|---|
| Comando de assinatura | `SignCommand` / `SignRequest` | Dados informados pelo usuário para criar assinatura simulada. |
| Comando de validação | `ValidateCommand` / `ValidateRequest` | Dados informados para validar assinatura simulada. |
| Resultado de assinatura | `OperationResult` / `SignResponse` | Resultado contendo assinatura simulada e mensagem. |
| Resultado de validação | `OperationResult` / `ValidateResponse` | Resultado contendo status de validade simulada. |
| Processo em execução | `ProcessMetadata` | PID, porta, aplicação, modo e data de início. |
| Artefato baixado | `ReleaseMetadata` | URL, versão, checksum e caminho local. |
| Erro da aplicação | `AppError` / `ErrorResponse` | Código, mensagem e detalhes do erro. |

### 12.2 Estrutura de metadados de processo

Exemplo conceitual de arquivo em `~/.hubsaude/processos/`:

```json
{
  "application": "assinador",
  "pid": 15342,
  "port": 8080,
  "mode": "server",
  "startedAt": "2026-05-07T20:30:00Z",
  "healthEndpoint": "http://localhost:8080/health",
  "status": "running"
}
```

### 12.3 Estrutura de metadados de artefato

```json
{
  "artifact": "simulador.jar",
  "version": "1.2.0",
  "url": "https://exemplo/releases/latest/download/simulador.jar",
  "checksumSha256": "valor-sha256-esperado",
  "downloadedAt": "2026-05-07T20:30:00Z"
}
```

### 12.4 Restrições do modelo de dados

- assinaturas simuladas não devem ser armazenadas permanentemente;
- registros de processo devem ser atualizados após encerramento;
- downloads devem ser reaproveitados quando possível;
- arquivos com checksum divergente devem ser rejeitados;
- metadados locais devem permanecer em diretório gerenciado.

---

## 13. Integração entre as partes do sistema

### 13.1 Integração CLI `assinatura` → `assinador.jar` em modo local

No modo local, o CLI `assinatura` executa o `assinador.jar` por meio de processo do sistema operacional.

```text
Usuário → assinatura CLI → java -jar assinador.jar → resposta → assinatura CLI → Usuário
```

Exemplo:

```bash
./bin/assinatura sign --local --jar ~/.hubsaude/assinador/assinador.jar --documento documento.json --certificado certificado.pem
```

Fluxo:

1. o usuário executa o comando `assinatura sign`;
2. o CLI valida os argumentos básicos;
3. o módulo de runtime localiza o Java;
4. o invocador local monta o comando `java -jar`;
5. o `assinador.jar` recebe os parâmetros;
6. o Java valida os parâmetros e gera resposta simulada;
7. o CLI captura e formata o resultado.

### 13.2 Integração CLI `assinatura` → `assinador.jar` em modo HTTP

No modo servidor, o `assinador.jar` é iniciado e passa a receber requisições HTTP.

```text
Usuário → assinatura CLI → HTTP POST /sign ou /validate → assinador.jar → resposta JSON → Usuário
```

Exemplo de inicialização:

```bash
./bin/assinatura start --jar ~/.hubsaude/assinador/assinador.jar --port 8080
```

Exemplo de uso:

```bash
./bin/assinatura sign --documento documento.json --certificado certificado.pem --port 8080
```

### 13.3 Integração CLI `simulador` → `simulador.jar`

O CLI `simulador` gerencia o ciclo de vida do Simulador do HubSaúde.

```text
Usuário → simulador CLI → processo Java simulador.jar → endpoints /api/info e /shutdown
```

Exemplo:

```bash
./bin/simulador start --jar ~/.hubsaude/simulador/simulador.jar --port 8443
```

### 13.4 Integração com diretório local `~/.hubsaude/`

O diretório local gerenciado é usado para armazenar:

- JDK/JRE ou referência ao runtime Java;
- `assinador.jar`;
- `simulador.jar`;
- metadados de artefatos;
- registros de processos;
- cache de downloads;
- logs, quando habilitados.

### 13.5 Integração com CI/CD

A integração com GitHub Actions ocorre por meio dos arquivos:

```text
.github/workflows/build.yml
.github/workflows/release.yml
```

Esses workflows são responsáveis por:

- executar testes;
- compilar CLIs Go;
- preparar o build Java;
- gerar binários multiplataforma;
- gerar checksums SHA256;
- assinar artefatos com Cosign;
- publicar releases.

---

## 14. APIs e endpoints implementados

### 14.1 Endpoints do `assinador.jar`

| Método | Endpoint | Finalidade | Entrada | Saída esperada |
|---|---|---|---|---|
| `GET` | `/health` | Verificar se o servidor está ativo. | Nenhuma. | Status de saúde do servidor. |
| `POST` | `/sign` | Criar assinatura simulada. | JSON de assinatura. | JSON com assinatura simulada. |
| `POST` | `/validate` | Validar assinatura simulada. | JSON com documento e assinatura. | JSON com resultado simulado. |
| `POST` ou `GET` | `/shutdown` | Solicitar encerramento do servidor. | Nenhuma ou simples requisição. | Confirmação de parada. |

### 14.2 Exemplo conceitual de requisição `/sign`

```json
{
  "document": "documento.json",
  "certificate": "certificado.pem"
}
```

### 14.3 Exemplo conceitual de resposta `/sign`

```json
{
  "success": true,
  "operation": "sign",
  "signature": "assinatura-simulada",
  "message": "Assinatura simulada criada com sucesso."
}
```

### 14.4 Exemplo conceitual de requisição `/validate`

```json
{
  "document": "documento.json",
  "signature": "assinatura-simulada"
}
```

### 14.5 Exemplo conceitual de resposta `/validate`

```json
{
  "success": true,
  "operation": "validate",
  "valid": true,
  "message": "Assinatura simulada considerada válida."
}
```

### 14.6 Endpoints esperados do Simulador do HubSaúde

| Método | Endpoint | Uso pelo CLI `simulador` |
|---|---|---|
| `GET` | `/api/info` | Consultar status e informações do simulador. |
| `POST` ou `GET` | `/shutdown` | Solicitar encerramento do simulador. |

Observação: o `simulador.jar` real não acompanha esta implementação de referência. O CLI foi preparado para executar um JAR informado por `--jar` ou obtido por `--source`.

---

## 15. Integração com banco de dados

O Sistema Runner **não utiliza banco de dados** nesta versão.

A persistência é realizada por meio de arquivos locais no diretório gerenciado `~/.hubsaude/`. Essa decisão é adequada ao contexto da aplicação, pois o Runner é uma ferramenta local de linha de comandos e não precisa armazenar dados de negócio complexos.

### 15.1 Justificativa para não usar banco de dados

| Justificativa | Explicação |
|---|---|
| Sistema local | O Runner executa na máquina do usuário, sem necessidade de servidor central. |
| Dados simples | Os dados persistidos são apenas metadados, cache e registros de processo. |
| Menor complexidade | Evita dependência de instalação e configuração de banco. |
| Escopo acadêmico | A proposta é integração e execução de aplicações Java, não gestão persistente de dados. |
| Segurança | O sistema não deve armazenar assinaturas ou dados sensíveis permanentemente. |

### 15.2 Persistência alternativa utilizada

| Tipo de dado | Forma de persistência |
|---|---|
| Registros de processo | Arquivos JSON em `~/.hubsaude/processos/`. |
| Artefatos baixados | Arquivos em subpastas de `~/.hubsaude/`. |
| Metadados de artefatos | Arquivos JSON ou metadados locais. |
| Checksums | Arquivos de texto ou metadados associados. |
| Logs | Arquivos locais, quando habilitados. |

---

## 16. Configuração do ambiente

### 16.1 Pré-requisitos para desenvolvimento

| Ferramenta | Versão recomendada | Finalidade |
|---|---|---|
| Go | Go 1.25 como versão-alvo do projeto | Compilar os CLIs `assinatura` e `simulador`. |
| Java | Java 21 | Executar e compilar o `assinador.jar`. |
| Maven | 3.9+ | Build e testes do projeto Java. |
| Git | Versão atual | Controle de versão. |
| Bash | Compatível com scripts do projeto | Execução de scripts auxiliares em Linux/macOS. |
| GitHub Actions | Ambiente remoto | CI/CD e releases. |
| Cosign | Versão compatível | Verificação/assinatura de artefatos, quando aplicável. |

### 16.2 Observação sobre Cobra

A documentação do projeto define o uso da biblioteca `github.com/spf13/cobra`. Para permitir execução local/offline no pacote acadêmico, o `go.mod` foi configurado com um `replace` apontando para uma implementação mínima local em `third_party/cobra/`.

Em um repositório com acesso à internet, recomenda-se remover o `replace` e utilizar o módulo oficial:

```bash
go get github.com/spf13/cobra
```

### 16.3 Diretório local gerenciado

O sistema utiliza o diretório:

```text
~/.hubsaude/
```

Estrutura sugerida:

```text
~/.hubsaude/
  assinador/
    assinador.jar
  simulador/
    simulador.jar
  processos/
    assinador-8080.json
    simulador-8443.json
  cache/
  logs/
```

---

## 17. Variáveis de ambiente

A implementação pode funcionar com parâmetros via CLI, mas algumas variáveis de ambiente podem ser usadas ou adotadas para facilitar configuração em ambientes diferentes.

| Variável | Obrigatória | Finalidade | Valor exemplo |
|---|---|---|---|
| `JAVA_HOME` | Não, mas recomendada | Indicar instalação Java 21. | `/usr/lib/jvm/java-21` |
| `PATH` | Sim, indiretamente | Permitir localizar `java`, `go`, `mvn` e scripts. | Inclui caminhos de Go, Java e Maven. |
| `RUNNER_HOME` | Não | Poderia sobrescrever o diretório padrão `~/.hubsaude/`. | `/home/usuario/.hubsaude` |
| `ASSINADOR_JAR` | Não | Caminho padrão para o `assinador.jar`. | `~/.hubsaude/assinador/assinador.jar` |
| `SIMULADOR_JAR` | Não | Caminho padrão para o `simulador.jar`. | `~/.hubsaude/simulador/simulador.jar` |
| `ASSINADOR_PORT` | Não | Porta padrão alternativa para o assinador. | `8080` |
| `SIMULADOR_PORT` | Não | Porta padrão alternativa para o simulador. | `8443` |
| `COSIGN_EXPERIMENTAL` | Condicional | Pode ser necessário em alguns fluxos de assinatura keyless com Cosign. | `1` |

Observação: as variáveis opcionais acima servem como diretriz de configuração. A implementação prioriza flags explícitas, como `--jar`, `--port`, `--source` e `--local`.

---

## 18. Procedimento de instalação

### 18.1 Obter o projeto

```bash
git clone <url-do-repositorio>
cd sistema-runner-atividade
```

Quando utilizado a partir do pacote `.zip`, basta extrair o arquivo e entrar na pasta do projeto.

### 18.2 Compilar os CLIs Go

```bash
go mod tidy
go build -o bin/assinatura ./cmd/assinatura
go build -o bin/simulador ./cmd/simulador
```

### 18.3 Compilar o `assinador.jar` com Maven

```bash
cd projetos/assinador-java
mvn clean package
```

O artefato esperado é:

```text
projetos/assinador-java/target/assinador.jar
```

### 18.4 Compilar o `assinador.jar` com script alternativo

Quando Maven não estiver disponível, pode ser usado o script auxiliar:

```bash
./scripts/build-assinador.sh
```

### 18.5 Instalar o `assinador.jar` no diretório gerenciado

```bash
mkdir -p ~/.hubsaude/assinador
cp projetos/assinador-java/target/assinador.jar ~/.hubsaude/assinador/assinador.jar
```

### 18.6 Instalação local por script

```bash
./scripts/install-local.sh
```

Esse script deve preparar diretórios e copiar artefatos principais para locais adequados, conforme a implementação disponível.

---

## 19. Procedimento de execução

### 19.1 Exibir versão do CLI `assinatura`

```bash
./bin/assinatura version
```

Resultado esperado:

```text
assinatura version <versão-atual>
```

### 19.2 Criar assinatura simulada em modo local

```bash
./bin/assinatura sign \
  --local \
  --jar ~/.hubsaude/assinador/assinador.jar \
  --documento documento.json \
  --certificado certificado.pem
```

Resultado esperado:

```text
Operação: sign
Status: sucesso
Mensagem: Assinatura simulada criada com sucesso.
Assinatura: assinatura-simulada
```

### 19.3 Validar assinatura simulada em modo local

```bash
./bin/assinatura validate \
  --local \
  --jar ~/.hubsaude/assinador/assinador.jar \
  --documento documento.json \
  --assinatura assinatura-simulada
```

Resultado esperado:

```text
Operação: validate
Status: sucesso
Resultado: assinatura válida
```

### 19.4 Iniciar o `assinador.jar` em modo servidor

```bash
./bin/assinatura start \
  --jar ~/.hubsaude/assinador/assinador.jar \
  --port 8080
```

Resultado esperado:

```text
Assinador iniciado em modo servidor.
Porta: 8080
PID: <pid>
```

### 19.5 Criar assinatura via HTTP

```bash
./bin/assinatura sign \
  --documento documento.json \
  --certificado certificado.pem \
  --port 8080
```

### 19.6 Parar o assinador

```bash
./bin/assinatura stop --port 8080
```

### 19.7 Iniciar o simulador com JAR local

```bash
./bin/simulador start \
  --jar ~/.hubsaude/simulador/simulador.jar \
  --port 8443
```

### 19.8 Consultar status do simulador

```bash
./bin/simulador status --port 8443
```

### 19.9 Parar o simulador

```bash
./bin/simulador stop --port 8443
```

---

## 20. Testes realizados

A implementação inclui testes iniciais em Go e Java, além de suporte para testes de integração manual e automatizável.

### 20.1 Testes Go

Comando:

```bash
go test ./...
```

Testes previstos/implementados:

| Teste | Arquivo | Finalidade |
|---|---|---|
| Teste de porta | `internal/process/port_test.go` | Verificar comportamento do verificador de portas. |
| Teste de checksum | `internal/artifacts/checksum_test.go` | Validar cálculo e comparação de SHA256. |
| Testes dos comandos CLI | Estrutura preparada | Validar parsing de comandos e flags. |
| Testes dos casos de uso | Estrutura preparada | Validar fluxos de assinatura, validação e simulador. |

### 20.2 Testes Java

Comando:

```bash
cd projetos/assinador-java
mvn test
```

Testes implementados:

| Teste | Arquivo | Finalidade |
|---|---|---|
| Validação de parâmetros | `ParameterValidatorTest.java` | Verificar entradas válidas e inválidas. |
| Serviço de assinatura simulada | `FakeSignatureServiceTest.java` | Verificar criação e validação simuladas. |

### 20.3 Testes de integração recomendados

| ID | Teste | Resultado esperado |
|---|---|---|
| TI-001 | Executar `assinatura version`. | CLI exibe versão atual. |
| TI-002 | Executar `assinatura sign --local`. | CLI invoca `assinador.jar` via `java -jar` e exibe assinatura simulada. |
| TI-003 | Executar `assinatura validate --local`. | CLI invoca `assinador.jar` e exibe resultado simulado. |
| TI-004 | Executar `assinatura start`. | Servidor Java é iniciado e registrado em `~/.hubsaude/processos/`. |
| TI-005 | Executar `assinatura sign` via HTTP. | CLI envia requisição a `/sign` e exibe resposta. |
| TI-006 | Executar `assinatura stop`. | Processo do assinador é encerrado e registro atualizado. |
| TI-007 | Executar `simulador start` com JAR local. | Processo é iniciado na porta informada. |
| TI-008 | Executar `simulador status`. | Status do simulador é exibido. |
| TI-009 | Executar `simulador stop`. | Simulador é encerrado. |
| TI-010 | Verificar checksum de artefato. | Sistema aceita arquivo válido e rejeita divergente. |

### 20.4 Observação sobre execução dos testes

Em ambiente de geração deste documento, a estrutura de testes e build está documentada para execução em ambiente com Go, Java 21 e Maven. Caso Maven não esteja instalado, o build Java pode ser feito com o script alternativo `scripts/build-assinador.sh`.

---

## 21. Tratamento de erros

A implementação adota tratamento de erros com mensagens compreensíveis, evitando expor detalhes técnicos desnecessários ao usuário final.

### 21.1 Tipos de erro tratados

| Erro | Tratamento esperado |
|---|---|
| Parâmetro obrigatório ausente | Exibir nome do parâmetro e orientação de uso. |
| Caminho do `.jar` inválido | Informar que o arquivo não foi encontrado. |
| Java não encontrado | Orientar instalação/configuração do Java 21 ou uso do diretório gerenciado. |
| Porta ocupada | Informar que a porta está indisponível e sugerir outra porta. |
| Processo não encontrado | Informar que não há instância ativa registrada. |
| Servidor HTTP indisponível | Informar falha de conexão e permitir uso local ou reinício. |
| Resposta HTTP inválida | Indicar resposta inesperada do servidor. |
| Falha de download | Informar URL, causa provável e tentativa de correção. |
| Checksum divergente | Rejeitar arquivo e orientar novo download. |
| Dispositivo PKCS#11 ausente | Informar que token/smart card não foi encontrado. |

### 21.2 Estrutura de erro no CLI

```text
Erro: <descrição resumida>
Motivo: <causa provável>
Orientação: <ação recomendada>
```

### 21.3 Estrutura de erro HTTP

```json
{
  "success": false,
  "error": "Parâmetro obrigatório ausente: documento"
}
```

### 21.4 Boas práticas adotadas

- validação antes do processamento;
- captura de erros de processo;
- uso de mensagens de erro orientativas;
- separação entre erro técnico e erro apresentado ao usuário;
- retorno HTTP adequado para falhas de requisição;
- não exposição de stack trace ao usuário final em uso normal.

---

## 22. Segurança implementada

A segurança implementada está concentrada na integridade de artefatos, assinatura de releases, isolamento de componentes e limitação da persistência local.

### 22.1 Segurança de artefatos

A estrutura de CI/CD prevê:

- geração de checksums SHA256;
- assinatura de artefatos com Cosign;
- publicação de arquivos `.sig` e `.pem`;
- verificação de artefatos por script auxiliar.

Exemplo de verificação:

```bash
./scripts/verify-artifact.sh <artefato> <checksum-esperado>
```

### 22.2 Verificação SHA256

O módulo `internal/artifacts/checksum.go` calcula e compara o hash SHA256 de artefatos. Essa verificação reduz o risco de uso de arquivos corrompidos ou alterados.

### 22.3 Segurança no gerenciamento de processos

O sistema registra processos iniciados em `~/.hubsaude/processos/`, permitindo controle posterior por comandos `stop` e `status`. Antes de reutilizar uma instância, o sistema pode confirmar sua disponibilidade por health check.

### 22.4 Segurança no uso de PKCS#11

A implementação inclui `PKCS11Adapter.java`, mantendo a integração criptográfica isolada. Nesta versão, o adaptador não executa assinatura real, mas prepara a estrutura para evolução futura sem misturar lógica criptográfica com os demais serviços.

### 22.5 Segurança de dados locais

O sistema não deve armazenar assinaturas digitais nem dados sensíveis. A persistência local é limitada a:

- metadados;
- registros de processo;
- arquivos `.jar`;
- cache;
- checksums;
- logs operacionais, quando utilizados.

### 22.6 Limite de segurança

O Sistema Runner, nesta versão, é uma implementação acadêmica de integração e simulação. Ele não deve ser usado como solução de assinatura digital juridicamente válida.

---

## 23. Controle de versão

### 23.1 Estratégia de versionamento

O projeto deve usar **versionamento semântico**, no formato:

```text
MAJOR.MINOR.PATCH
```

Exemplo:

```text
v0.1.0
v1.0.0
```

### 23.2 Controle por Git

O código-fonte deve ser mantido em repositório Git, com organização mínima recomendada:

```text
main
feature/<nome-da-funcionalidade>
fix/<nome-da-correcao>
release/<versao>
```

### 23.3 Tags de release

As releases devem ser geradas a partir de tags:

```bash
git tag v0.1.0
git push origin v0.1.0
```

### 23.4 CI/CD

O controle de versão é integrado aos workflows:

| Workflow | Finalidade |
|---|---|
| `build.yml` | Compilar e testar o projeto a cada mudança relevante. |
| `release.yml` | Gerar binários, checksums, assinaturas e publicar release. |

### 23.5 Artefatos de release esperados

```text
assinatura-v0.1.0-linux-amd64
assinatura-v0.1.0-windows-amd64.exe
assinatura-v0.1.0-darwin-amd64
simulador-v0.1.0-linux-amd64
simulador-v0.1.0-windows-amd64.exe
simulador-v0.1.0-darwin-amd64
checksums.txt
<artefato>.sig
<artefato>.pem
```

---

## 24. Evidências da implementação e integração

### 24.1 Evidência de estrutura do projeto

A presença dos seguintes diretórios evidencia a separação modular da implementação:

```text
cmd/assinatura
cmd/simulador
internal/assinatura
internal/simulador
internal/java
internal/process
internal/artifacts
projetos/assinador-java
.github/workflows
scripts
```

### 24.2 Evidência dos CLIs implementados

Comandos previstos:

```bash
./bin/assinatura version
./bin/assinatura sign --local --jar ~/.hubsaude/assinador/assinador.jar --documento documento.json --certificado certificado.pem
./bin/assinatura validate --local --jar ~/.hubsaude/assinador/assinador.jar --documento documento.json --assinatura assinatura-simulada
./bin/assinatura start --jar ~/.hubsaude/assinador/assinador.jar --port 8080
./bin/assinatura stop --port 8080

./bin/simulador start --jar ~/.hubsaude/simulador/simulador.jar --port 8443
./bin/simulador status --port 8443
./bin/simulador stop --port 8443
```

### 24.3 Evidência da integração local

A integração local é evidenciada pelo módulo:

```text
internal/java/invoker/local.go
```

Esse módulo executa o `assinador.jar` com o mecanismo `java -jar`.

### 24.4 Evidência da integração HTTP

A integração HTTP é evidenciada pelos módulos:

```text
internal/java/invoker/http.go
projetos/assinador-java/src/main/java/br/ufg/hubsaude/assinador/controller/SignatureHttpServer.java
```

Esses arquivos implementam o cliente HTTP em Go e o servidor HTTP no Java.

### 24.5 Evidência da validação de parâmetros

A validação de parâmetros é evidenciada por:

```text
projetos/assinador-java/src/main/java/br/ufg/hubsaude/assinador/validation/ParameterValidator.java
projetos/assinador-java/src/test/java/br/ufg/hubsaude/assinador/validation/ParameterValidatorTest.java
```

### 24.6 Evidência da simulação de assinatura

A simulação é evidenciada por:

```text
projetos/assinador-java/src/main/java/br/ufg/hubsaude/assinador/service/SignatureService.java
projetos/assinador-java/src/main/java/br/ufg/hubsaude/assinador/service/FakeSignatureService.java
projetos/assinador-java/src/test/java/br/ufg/hubsaude/assinador/service/FakeSignatureServiceTest.java
```

### 24.7 Evidência do gerenciamento de processos

O gerenciamento de processos é evidenciado por:

```text
internal/process/manager.go
internal/process/registry.go
internal/process/metadata.go
internal/process/health.go
internal/process/port.go
```

### 24.8 Evidência de segurança de artefatos

A verificação de integridade e release seguro são evidenciadas por:

```text
internal/artifacts/checksum.go
scripts/verify-artifact.sh
.github/workflows/release.yml
```

### 24.9 Evidência de CI/CD

A automação é evidenciada pelos workflows:

```text
.github/workflows/build.yml
.github/workflows/release.yml
```

### 24.10 Evidência de documentação

A documentação da implementação é evidenciada por:

```text
README.md
docs/relatorio_implementacao.md
```

---

## 25. Limitações e melhorias futuras

### 25.1 Limitações conhecidas

| Limitação | Descrição |
|---|---|
| Assinatura simulada | A assinatura não é criptográfica real. |
| Validação simulada | A validação segue critério predeterminado e não verifica assinatura real. |
| PKCS#11 estruturado | O adaptador existe, mas não executa assinatura real no escopo atual. |
| Provisionamento de JDK/JRE | Estrutura existe, mas a implementação prioriza Java 21 instalado ou já disponível localmente. |
| `simulador.jar` real ausente | O JAR real do simulador deve ser informado por `--jar` ou obtido por `--source`. |
| Sem banco de dados | Persistência é feita apenas por arquivos locais. |
| Sem interface gráfica | O sistema é exclusivamente CLI. |
| Sem autenticação | Não há controle de usuários, pois o uso é local. |
| Sem validade jurídica | O sistema não deve ser usado como solução real de assinatura digital. |

### 25.2 Melhorias futuras

| Melhoria | Descrição |
|---|---|
| Completar provisionamento automático de JDK/JRE | Implementar download automático completo por plataforma, com extração e cache. |
| Expandir suporte PKCS#11 | Integrar SoftHSM2 e preparar cenários com token/smart card real. |
| Adicionar comando `doctor` | Verificar ambiente, Java, portas, arquivos e permissões. |
| Melhorar logs estruturados | Registrar eventos operacionais em formato JSON ou texto estruturado. |
| Ampliar testes end-to-end | Testar fluxos completos local e HTTP. |
| Ampliar testes multiplataforma | Validar execução em Windows, Linux e macOS. |
| Suporte ARM64 | Gerar binários para arquiteturas adicionais. |
| Configuração por arquivo | Permitir `config.yaml` ou `config.json`. |
| Melhorar gerenciamento de releases | Integrar totalmente `release.json`, SemVer e cache de versões. |
| Documentação de troubleshooting | Criar guia para erros comuns e soluções. |

---

## 26. Referências

- Especificação original do Sistema Runner — Trabalho Prático.
- Plano Revisado #2 do Sistema Runner.
- Documento de Design do Sistema Runner baseado no Modelo C4.
- Especificação de Requisitos de Software — Sistema Runner.
- Documento de Arquitetura de Software — Sistema Runner.
- Documento de Projeto Detalhado de Software — Sistema Runner.
- Documento de Modelo C4 de Software — Sistema Runner.
- README da implementação de referência acadêmica do Sistema Runner.
- Boas práticas de Engenharia de Software para implementação, integração, testes, documentação, segurança de artefatos e controle de versão.
