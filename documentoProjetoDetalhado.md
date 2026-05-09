# Documento de Projeto Detalhado de Software — Sistema Runner

## 1. Identificação do documento

| Campo | Informação |
|---|---|
| **Nome do sistema** | Sistema Runner |
| **Nome do documento** | Documento de Projeto Detalhado de Software |
| **Versão do documento** | 1.0 |
| **Data de elaboração** | 07/05/2026 |
| **Responsável pela elaboração** | Equipe do projeto / Disciplina de Implementação e Integração de Software |
| **Instituição / contexto acadêmico** | Bacharelado em Engenharia de Software — Universidade Federal de Goiás (UFG) |
| **Contexto de aplicação** | Plataforma HubSaúde — interoperabilidade de dados em saúde |
| **Documentos relacionados** | Especificação de Requisitos de Software — Sistema Runner; Documento de Arquitetura de Software — Sistema Runner; Plano Revisado #2; Documento de Design C4 |
| **Tipo de sistema** | Ferramenta de linha de comandos, integração com aplicações Java, gerenciamento de processos, provisionamento de JDK/JRE e simulação de assinatura digital |

---

## 2. Histórico de versões

| Versão | Data | Autor / Responsável | Descrição da alteração |
|---|---|---|---|
| 1.0 | 07/05/2026 | Equipe do projeto | Elaboração inicial do Documento de Projeto Detalhado de Software do Sistema Runner, com base nos arquivos de especificação, plano revisitado, design arquitetural C4, documento de requisitos e documento de arquitetura previamente elaborados. |

---

## 3. Sumário

1. [Identificação do documento](#1-identificação-do-documento)  
2. [Histórico de versões](#2-histórico-de-versões)  
3. [Sumário](#3-sumário)  
4. [Introdução](#4-introdução)  
   4.1 [Objetivo do documento](#41-objetivo-do-documento)  
   4.2 [Contexto do sistema](#42-contexto-do-sistema)  
5. [Escopo do projeto detalhado](#5-escopo-do-projeto-detalhado)  
6. [Visão geral da solução](#6-visão-geral-da-solução)  
7. [Organização de módulos](#7-organização-de-módulos)  
8. [Organização de pastas e arquivos](#8-organização-de-pastas-e-arquivos)  
9. [Projeto das classes principais](#9-projeto-das-classes-principais)  
   9.1 [Classes de entidade](#91-classes-de-entidade)  
   9.2 [Classes de serviço](#92-classes-de-serviço)  
   9.3 [Classes de repositório](#93-classes-de-repositório)  
   9.4 [Classes de controle](#94-classes-de-controle)  
10. [Projeto dos métodos principais](#10-projeto-dos-métodos-principais)  
11. [Projeto dos dados](#11-projeto-dos-dados)  
   11.1 [Entidades](#111-entidades)  
   11.2 [Campos](#112-campos)  
   11.3 [Relacionamentos](#113-relacionamentos)  
   11.4 [Restrições](#114-restrições)  
12. [Projeto das interfaces ou endpoints](#12-projeto-das-interfaces-ou-endpoints)  
13. [Fluxos principais do sistema](#13-fluxos-principais-do-sistema)  
14. [Regras de negócio detalhadas](#14-regras-de-negócio-detalhadas)  
15. [Tratamento de erros](#15-tratamento-de-erros)  
16. [Projeto de segurança](#16-projeto-de-segurança)  
17. [Projeto de persistência](#17-projeto-de-persistência)  
18. [Projeto de testes](#18-projeto-de-testes)  
19. [Rastreabilidade](#19-rastreabilidade)  
20. [Decisões técnicas principais](#20-decisões-técnicas-principais)  
21. [Apêndice A — Exemplos conceituais de comandos](#apêndice-a--exemplos-conceituais-de-comandos)  
22. [Apêndice B — Exemplos conceituais de estruturas JSON](#apêndice-b--exemplos-conceituais-de-estruturas-json)  
23. [Referências documentais utilizadas](#referências-documentais-utilizadas)  

---

## 4. Introdução

O Sistema Runner é uma solução de software voltada para facilitar a execução e o gerenciamento de aplicações Java relacionadas ao ecossistema HubSaúde, por meio de interfaces de linha de comandos. O sistema é composto, principalmente, pelos CLIs `assinatura` e `simulador`, desenvolvidos em Go, e pela aplicação Java `assinador.jar`, responsável por validar parâmetros e simular operações de criação e validação de assinatura digital.

Este documento apresenta o projeto detalhado de software do Sistema Runner. Enquanto a especificação de requisitos define **o que** o sistema deve fazer e o documento de arquitetura define **como a solução é organizada em alto nível**, este documento detalha **como os módulos, classes, métodos, dados, interfaces, fluxos e testes devem ser projetados para orientar a implementação**.

### 4.1 Objetivo do documento

O objetivo deste documento é detalhar a estrutura interna do Sistema Runner, especificando:

- os módulos que compõem a solução;
- a organização sugerida de pastas e arquivos;
- as principais classes, estruturas, serviços, controladores e repositórios;
- os principais métodos e suas responsabilidades;
- os dados manipulados e persistidos localmente;
- as interfaces CLI e endpoints HTTP;
- os fluxos principais de execução;
- as regras de negócio em nível detalhado;
- o tratamento de erros;
- o projeto de segurança;
- o projeto de persistência;
- o projeto de testes;
- a rastreabilidade entre requisitos, módulos e testes;
- as principais decisões técnicas que orientam a implementação.

Este documento deve servir como referência para desenvolvedores, estudantes, avaliadores, integradores e mantenedores do Sistema Runner.

### 4.2 Contexto do sistema

O Sistema Runner foi definido no contexto de uma disciplina de Implementação e Integração de Software do Bacharelado em Engenharia de Software, com aplicação prática relacionada à Plataforma HubSaúde. O objetivo do sistema é facilitar o acesso a aplicações Java por linha de comandos, reduzindo a necessidade de o usuário conhecer detalhes de instalação do Java, execução de arquivos `.jar`, parâmetros técnicos, portas, processos em segundo plano e comunicação com dispositivos criptográficos.

No contexto do projeto:

- o usuário interage com o sistema por comandos CLI;
- o CLI `assinatura` invoca o `assinador.jar` para operações simuladas de assinatura e validação;
- o `assinador.jar` valida parâmetros e retorna respostas simuladas;
- o CLI `simulador` gerencia o ciclo de vida do Simulador do HubSaúde;
- o sistema pode provisionar automaticamente JDK/JRE compatível;
- o sistema pode armazenar arquivos, cache e registros locais no diretório `~/.hubsaude/`;
- os artefatos distribuídos devem possuir checksums SHA256 e assinatura Cosign;
- o sistema deve ser multiplataforma, com suporte a Windows, Linux e macOS em arquitetura amd64.

---

## 5. Escopo do projeto detalhado

Este documento detalha o projeto interno das partes do Sistema Runner que estão no escopo da implementação.

### 5.1 Está no escopo deste projeto detalhado

- projeto interno do CLI `assinatura`;
- projeto interno do CLI `simulador`;
- projeto interno do `assinador.jar`;
- definição de módulos, responsabilidades e dependências;
- definição de entidades, serviços, repositórios e controladores;
- projeto dos comandos CLI;
- projeto dos endpoints HTTP `/sign` e `/validate`;
- projeto da invocação local via `java -jar`;
- projeto da invocação HTTP em modo servidor;
- projeto do gerenciamento de processos;
- projeto do provisionamento automático de JDK/JRE;
- projeto do download e cache do `simulador.jar`;
- projeto da validação de parâmetros;
- projeto da simulação de criação e validação de assinatura;
- projeto de persistência local em `~/.hubsaude/`;
- projeto de tratamento de erros;
- projeto de testes unitários, integração, aceitação e multiplataforma;
- rastreabilidade entre requisitos, histórias, módulos e testes.

### 5.2 Não está no escopo deste projeto detalhado

- implementação real de assinatura digital criptográfica;
- implementação real de validação criptográfica de assinatura digital;
- integração real com autoridades certificadoras;
- geração de certificados digitais;
- autenticação de usuários;
- armazenamento persistente de assinaturas digitais;
- interface gráfica;
- implantação em ambiente distribuído de produção;
- substituição do Simulador do HubSaúde por uma aplicação própria de negócio;
- definição final dos parâmetros oficiais dos casos de uso FHIR além dos exemplos conceituais apresentados neste documento.

---

## 6. Visão geral da solução

A solução é composta por três aplicações principais e por módulos auxiliares compartilhados.

| Aplicação / parte | Tecnologia | Responsabilidade principal |
|---|---|---|
| **`assinatura`** | Go 1.25 + Cobra | CLI responsável por receber comandos de criação e validação de assinatura, invocar o `assinador.jar`, formatar resultados e gerenciar o modo servidor do assinador. |
| **`simulador`** | Go 1.25 + Cobra | CLI responsável por iniciar, parar, consultar status e baixar dinamicamente o Simulador do HubSaúde. |
| **`assinador.jar`** | Java 21 | Aplicação Java responsável por validar parâmetros, simular criação de assinatura, simular validação de assinatura e expor endpoints HTTP. |
| **Diretório local gerenciado** | Sistema de arquivos | Armazena JDK/JRE, arquivos `.jar`, metadados, cache, registros de processos e logs. |
| **Pipeline CI/CD** | GitHub Actions | Gera binários multiplataforma, publica releases, checksums SHA256 e assinaturas Cosign. |

### 6.1 Funcionamento conceitual

O funcionamento geral do Sistema Runner pode ser resumido da seguinte forma:

```text
Usuário
  ↓
CLI assinatura ou CLI simulador
  ↓
Serviços internos do Runner
  ↓
JDK/JRE local ou provisionado automaticamente
  ↓
assinador.jar ou simulador.jar
  ↓
Resposta formatada ao usuário
```

### 6.2 Responsabilidades principais

| Responsabilidade | Onde será implementada |
|---|---|
| Receber comandos do usuário | Camada CLI em Go, usando Cobra |
| Validar argumentos básicos do CLI | Módulos de validação nos CLIs |
| Validar parâmetros de assinatura | `assinador.jar`, em `ParameterValidator` |
| Simular assinatura | `FakeSignatureService` |
| Simular validação | `FakeSignatureService` |
| Invocar Java localmente | `LocalJarInvoker` |
| Invocar assinador via HTTP | `AssinadorHttpClient` |
| Iniciar/parar processos Java | `ProcessManager` |
| Registrar PID e porta | `ProcessRegistry` |
| Baixar JDK/JRE | `JavaRuntimeDownloader` |
| Baixar `simulador.jar` | `ArtifactDownloader` |
| Verificar checksum | `ChecksumVerifier` |
| Exibir mensagens ao usuário | `OutputFormatter` |
| Padronizar erros | `ErrorHandler` |

---

## 7. Organização de módulos

A organização de módulos proposta segue separação de responsabilidades, baixo acoplamento e alta coesão.

### 7.1 Módulos do CLI `assinatura`

| Módulo | Responsabilidade |
|---|---|
| `cmd/assinatura` | Ponto de entrada do binário `assinatura`. |
| `internal/cli/assinatura/commands` | Definição dos comandos `sign`, `validate`, `stop`, `version` e `help`. |
| `internal/cli/assinatura/formatter` | Formatação das respostas de sucesso e erro. |
| `internal/assinatura/usecase` | Orquestra os casos de uso de assinatura e validação. |
| `internal/assinatura/validation` | Validação preliminar dos argumentos recebidos pelo CLI. |
| `internal/java/runtime` | Detecção e provisionamento de JDK/JRE. |
| `internal/java/invoker` | Invocação local por `java -jar` e invocação HTTP. |
| `internal/process` | Gerenciamento de processos, portas, health check e registros. |
| `internal/config` | Resolução de caminhos locais e configurações. |
| `internal/errors` | Tipos e mensagens padronizadas de erro. |

### 7.2 Módulos do CLI `simulador`

| Módulo | Responsabilidade |
|---|---|
| `cmd/simulador` | Ponto de entrada do binário `simulador`. |
| `internal/cli/simulador/commands` | Definição dos comandos `start`, `stop`, `status` e `help`. |
| `internal/cli/simulador/formatter` | Formatação das respostas do simulador. |
| `internal/simulador/usecase` | Orquestra início, parada e consulta de status do simulador. |
| `internal/artifacts` | Consulta de release, download, checksum e cache do `simulador.jar`. |
| `internal/java/runtime` | Detecção e provisionamento de JDK/JRE. |
| `internal/process` | Gerenciamento do processo do simulador. |
| `internal/config` | Caminhos locais em `~/.hubsaude/`. |
| `internal/errors` | Mensagens padronizadas de erro. |

### 7.3 Módulos do `assinador.jar`

| Módulo / pacote | Responsabilidade |
|---|---|
| `br.ufg.hubsaude.assinador` | Ponto de entrada da aplicação Java. |
| `cli` | Recebimento de comandos quando o `.jar` é invocado localmente. |
| `controller` | Endpoints HTTP `/sign` e `/validate`. |
| `dto` | Objetos de entrada e saída de assinatura e validação. |
| `validation` | Validação dos parâmetros obrigatórios e formatos. |
| `service` | Interface e implementação da lógica de assinatura simulada. |
| `pkcs11` | Adaptador para interação com dispositivo criptográfico. |
| `error` | Tratamento global de exceções e respostas de erro. |
| `config` | Configurações da aplicação, porta e modo de execução. |

### 7.4 Módulos de CI/CD

| Módulo | Responsabilidade |
|---|---|
| `.github/workflows/build.yml` | Executar testes e gerar builds para cada push ou pull request. |
| `.github/workflows/release.yml` | Gerar artefatos de release, checksums SHA256 e assinaturas Cosign. |
| `scripts/` | Scripts auxiliares para geração de diagramas, empacotamento ou automações locais. |

---

## 8. Organização de pastas e arquivos

A organização abaixo é uma proposta detalhada e coerente com a arquitetura definida. Os nomes podem ser ajustados durante a implementação, desde que as responsabilidades sejam preservadas.

```text
sistema-runner/
  README.md
  LICENSE
  go.mod
  go.sum

  docs/
    requisitos.md
    arquitetura.md
    projeto-detalhado.md
    manual-usuario.md
    implantacao.md
    testes.md
    diagramas/
      contexto.puml
      conteineres.puml
      componentes.puml
      codigo.puml
      implantacao.puml
      imagens/

  cmd/
    assinatura/
      main.go
    simulador/
      main.go

  internal/
    cli/
      assinatura/
        commands/
          root.go
          sign.go
          validate.go
          stop.go
          version.go
        formatter/
          output_formatter.go
      simulador/
        commands/
          root.go
          start.go
          stop.go
          status.go
        formatter/
          output_formatter.go

    assinatura/
      usecase/
        sign_usecase.go
        validate_usecase.go
        stop_assinador_usecase.go
      validation/
        cli_argument_validator.go
      dto/
        sign_command.go
        validate_command.go
        operation_result.go

    simulador/
      usecase/
        start_simulator_usecase.go
        stop_simulator_usecase.go
        status_simulator_usecase.go
      dto/
        simulator_status.go
        simulator_config.go

    java/
      runtime/
        java_runtime.go
        java_detector.go
        java_downloader.go
        java_resolver.go
      invoker/
        local_jar_invoker.go
        assinador_http_client.go

    process/
      process_manager.go
      process_registry.go
      health_checker.go
      port_checker.go
      process_metadata.go

    artifacts/
      release_resolver.go
      artifact_downloader.go
      checksum_verifier.go
      cache_manager.go
      release_metadata.go

    config/
      paths.go
      app_config.go
      constants.go

    errors/
      app_error.go
      error_handler.go
      messages.go

  projetos/
    assinador-java/
      pom.xml
      README.md
      src/
        main/
          java/
            br/
              ufg/
                hubsaude/
                  assinador/
                    Main.java
                    cli/
                      CliAdapter.java
                      CliCommandParser.java
                    config/
                      ServerConfig.java
                      Pkcs11Config.java
                    controller/
                      SignatureController.java
                      HealthController.java
                    dto/
                      SignRequest.java
                      SignResponse.java
                      ValidateRequest.java
                      ValidateResponse.java
                      ErrorResponse.java
                    validation/
                      ParameterValidator.java
                      ValidationResult.java
                    service/
                      SignatureService.java
                      FakeSignatureService.java
                    pkcs11/
                      PKCS11Adapter.java
                      PKCS11Exception.java
                    error/
                      GlobalExceptionHandler.java
                      InvalidParameterException.java
        test/
          java/
            br/
              ufg/
                hubsaude/
                  assinador/
                    validation/
                    service/
                    controller/
                    cli/

  scripts/
    geraimagens.sh
    geraimagens.bat
    verify-artifacts.sh

  .github/
    workflows/
      build.yml
      release.yml
```

### 8.1 Justificativa da organização

| Diretório | Justificativa |
|---|---|
| `cmd/` | Convenção comum em projetos Go para pontos de entrada de executáveis. |
| `internal/` | Impede uso externo dos pacotes internos e organiza a lógica da aplicação. |
| `internal/cli/` | Separa comandos e formatação da lógica de negócio. |
| `internal/java/` | Centraliza detecção de Java, provisionamento e invocação de `.jar`. |
| `internal/process/` | Reúne controle de processos, portas, health checks e PID. |
| `internal/artifacts/` | Centraliza download, releases, cache e checksum. |
| `projetos/assinador-java/` | Separa claramente a aplicação Java do restante do código Go. |
| `.github/workflows/` | Centraliza automações de build, testes e releases. |
| `docs/` | Mantém documentação técnica e de uso no próprio repositório. |

---

## 9. Projeto das classes principais

Observação: em Go, o conceito de “classe” é normalmente representado por `structs`, `interfaces`, funções e pacotes. Neste documento, o termo “classe” é usado em sentido amplo, contemplando tanto classes Java quanto structs/interfaces Go.

### 9.1 Classes de entidade

As entidades representam dados estruturados manipulados pelo sistema.

#### 9.1.1 Entidades em Go

| Entidade / Struct | Pacote sugerido | Responsabilidade |
|---|---|---|
| `SignCommand` | `internal/assinatura/dto` | Representa os parâmetros recebidos no comando `assinatura sign`. |
| `ValidateCommand` | `internal/assinatura/dto` | Representa os parâmetros recebidos no comando `assinatura validate`. |
| `OperationResult` | `internal/assinatura/dto` | Representa resultado de sucesso ou erro de uma operação de assinatura. |
| `JavaRuntime` | `internal/java/runtime` | Representa uma instalação de Java detectada ou provisionada. |
| `ProcessMetadata` | `internal/process` | Representa PID, porta, aplicação, modo e data de início de um processo. |
| `SimulatorStatus` | `internal/simulador/dto` | Representa o status do Simulador do HubSaúde. |
| `SimulatorConfig` | `internal/simulador/dto` | Representa porta, origem de download e configurações do simulador. |
| `ReleaseMetadata` | `internal/artifacts` | Representa informações de versão, URL e checksum de artefato. |
| `AppError` | `internal/errors` | Representa erro padronizado da aplicação. |

#### 9.1.2 Entidades em Java

| Classe | Pacote sugerido | Responsabilidade |
|---|---|---|
| `SignRequest` | `dto` | Representa a requisição de criação de assinatura. |
| `SignResponse` | `dto` | Representa a resposta de criação de assinatura simulada. |
| `ValidateRequest` | `dto` | Representa a requisição de validação de assinatura. |
| `ValidateResponse` | `dto` | Representa a resposta de validação simulada. |
| `ErrorResponse` | `dto` | Representa resposta padronizada de erro. |
| `ValidationResult` | `validation` | Representa resultado interno da validação de parâmetros. |
| `Pkcs11Config` | `config` | Representa configurações de acesso ao dispositivo PKCS#11. |
| `ServerConfig` | `config` | Representa configurações de porta e modo servidor. |

### 9.2 Classes de serviço

As classes de serviço executam regras de aplicação e coordenam operações.

#### 9.2.1 Serviços em Go

| Serviço | Pacote sugerido | Responsabilidade |
|---|---|---|
| `SignUseCase` | `internal/assinatura/usecase` | Orquestra o fluxo de criação de assinatura simulada. |
| `ValidateUseCase` | `internal/assinatura/usecase` | Orquestra o fluxo de validação de assinatura simulada. |
| `StopAssinadorUseCase` | `internal/assinatura/usecase` | Orquestra a interrupção do `assinador.jar`. |
| `StartSimulatorUseCase` | `internal/simulador/usecase` | Orquestra o início do Simulador do HubSaúde. |
| `StopSimulatorUseCase` | `internal/simulador/usecase` | Orquestra a parada do Simulador do HubSaúde. |
| `StatusSimulatorUseCase` | `internal/simulador/usecase` | Orquestra a consulta de status do Simulador. |
| `JavaRuntimeResolver` | `internal/java/runtime` | Decide se usa Java local ou se baixa uma versão compatível. |
| `JavaDetector` | `internal/java/runtime` | Detecta Java no `PATH` ou em `~/.hubsaude/`. |
| `JavaDownloader` | `internal/java/runtime` | Baixa e prepara JDK/JRE compatível. |
| `LocalJarInvoker` | `internal/java/invoker` | Executa aplicações Java via `java -jar`. |
| `AssinadorHttpClient` | `internal/java/invoker` | Envia requisições HTTP ao `assinador.jar`. |
| `ProcessManager` | `internal/process` | Inicia, encerra e monitora processos. |
| `HealthChecker` | `internal/process` | Verifica se um serviço HTTP está respondendo. |
| `PortChecker` | `internal/process` | Verifica disponibilidade de portas. |
| `ReleaseResolver` | `internal/artifacts` | Obtém metadados de release e identifica versão mais recente. |
| `ArtifactDownloader` | `internal/artifacts` | Baixa artefatos remotos. |
| `ChecksumVerifier` | `internal/artifacts` | Verifica integridade de arquivos baixados. |
| `CacheManager` | `internal/artifacts` | Controla cache local de artefatos. |
| `OutputFormatter` | `internal/cli/*/formatter` | Converte resultados internos em mensagens legíveis no terminal. |
| `ErrorHandler` | `internal/errors` | Converte erros técnicos em mensagens orientativas. |

#### 9.2.2 Serviços em Java

| Serviço | Pacote sugerido | Responsabilidade |
|---|---|---|
| `SignatureService` | `service` | Interface para operações `sign` e `validate`. |
| `FakeSignatureService` | `service` | Implementa assinatura e validação simuladas. |
| `ParameterValidator` | `validation` | Valida presença, formato e consistência dos parâmetros. |
| `PKCS11Adapter` | `pkcs11` | Encapsula a comunicação com dispositivo criptográfico. |
| `CliCommandParser` | `cli` | Interpreta argumentos recebidos via modo local. |
| `GlobalExceptionHandler` | `error` | Converte exceções em respostas padronizadas. |

### 9.3 Classes de repositório

As classes de repositório controlam acesso a dados persistidos localmente. O sistema não possui banco de dados relacional; a persistência é feita em arquivos locais, preferencialmente no diretório `~/.hubsaude/`.

#### 9.3.1 Repositórios em Go

| Repositório | Pacote sugerido | Responsabilidade |
|---|---|---|
| `ProcessRegistry` | `internal/process` | Salvar, ler, atualizar e remover registros de processos. |
| `MetadataRepository` | `internal/config` ou `internal/artifacts` | Salvar e ler metadados de versões, artefatos e downloads. |
| `CacheRepository` | `internal/artifacts` | Gerenciar arquivos em cache, como `release.json` e checksums. |
| `LocalFileRepository` | `internal/config` | Centralizar operações de leitura e escrita no diretório gerenciado. |
| `LogRepository` | `internal/config` ou `internal/logging` | Registrar eventos operacionais, quando implementado. |

#### 9.3.2 Repositórios em Java

O `assinador.jar` não deve persistir assinaturas nem resultados de assinatura. Assim, os repositórios Java são opcionais e limitados a configuração local.

| Repositório | Responsabilidade |
|---|---|
| `Pkcs11ConfigRepository` | Ler configuração local de PKCS#11, se adotado arquivo de configuração. |
| `ServerConfigRepository` | Ler parâmetros de execução do servidor, quando necessário. |

### 9.4 Classes de controle

As classes de controle recebem comandos ou requisições externas e chamam serviços internos.

#### 9.4.1 Controladores em Go

| Controlador / comando | Pacote sugerido | Responsabilidade |
|---|---|---|
| `RootCommand` | `internal/cli/assinatura/commands` | Define comando raiz `assinatura`. |
| `SignCommandHandler` | `internal/cli/assinatura/commands` | Recebe comando `assinatura sign`. |
| `ValidateCommandHandler` | `internal/cli/assinatura/commands` | Recebe comando `assinatura validate`. |
| `StopCommandHandler` | `internal/cli/assinatura/commands` | Recebe comando `assinatura stop`. |
| `VersionCommandHandler` | `internal/cli/assinatura/commands` | Recebe comando `assinatura version`. |
| `SimulatorRootCommand` | `internal/cli/simulador/commands` | Define comando raiz `simulador`. |
| `StartCommandHandler` | `internal/cli/simulador/commands` | Recebe comando `simulador start`. |
| `SimulatorStopCommandHandler` | `internal/cli/simulador/commands` | Recebe comando `simulador stop`. |
| `StatusCommandHandler` | `internal/cli/simulador/commands` | Recebe comando `simulador status`. |

#### 9.4.2 Controladores em Java

| Controlador | Endpoint / entrada | Responsabilidade |
|---|---|---|
| `CliAdapter` | Argumentos de linha de comando | Receber chamadas locais via `java -jar`. |
| `SignatureController` | `POST /sign`, `POST /validate` | Receber requisições HTTP de assinatura e validação. |
| `HealthController` | `GET /health` | Informar se o `assinador.jar` está ativo. |
| `GlobalExceptionHandler` | Exceções internas | Padronizar respostas de erro. |

---

## 10. Projeto dos métodos principais

Esta seção descreve os principais métodos previstos. As assinaturas são conceituais e podem ser adaptadas conforme a linguagem e bibliotecas utilizadas.

### 10.1 Métodos principais do CLI `assinatura`

#### `SignUseCase.Execute(command SignCommand) (OperationResult, error)`

| Aspecto | Descrição |
|---|---|
| Responsabilidade | Executar o fluxo de criação de assinatura simulada. |
| Entrada | `SignCommand` com parâmetros informados pelo usuário. |
| Saída | `OperationResult` com assinatura simulada ou erro. |
| Passos | Validar comando; resolver Java; verificar modo servidor/local; invocar `assinador.jar`; formatar resultado. |
| Erros tratados | Parâmetro ausente, Java ausente, jar não encontrado, servidor indisponível, erro HTTP, erro de execução local. |

#### `ValidateUseCase.Execute(command ValidateCommand) (OperationResult, error)`

| Aspecto | Descrição |
|---|---|
| Responsabilidade | Executar o fluxo de validação simulada de assinatura. |
| Entrada | `ValidateCommand`. |
| Saída | `OperationResult` indicando assinatura válida ou inválida. |
| Passos | Validar comando; resolver Java; escolher invocação; chamar `assinador.jar`; formatar retorno. |
| Erros tratados | Parâmetros inválidos, servidor indisponível, erro de validação, falha de comunicação. |

#### `StopAssinadorUseCase.Execute(port int) error`

| Aspecto | Descrição |
|---|---|
| Responsabilidade | Encerrar o `assinador.jar` em modo servidor. |
| Entrada | Porta informada ou porta padrão. |
| Saída | Erro ou confirmação de encerramento. |
| Passos | Consultar registro; verificar processo; enviar shutdown ou encerrar processo; atualizar registro. |
| Erros tratados | Processo inexistente, porta inválida, falha ao encerrar, registro desatualizado. |

#### `JavaRuntimeResolver.Resolve() (JavaRuntime, error)`

| Aspecto | Descrição |
|---|---|
| Responsabilidade | Obter Java compatível para execução das aplicações Java. |
| Entrada | Configuração de versão esperada. |
| Saída | Caminho do executável Java. |
| Passos | Verificar Java em `~/.hubsaude/`; verificar Java no PATH; validar versão; baixar se ausente; retornar caminho. |
| Erros tratados | Download indisponível, versão incompatível, permissões de arquivo, checksum inválido. |

#### `LocalJarInvoker.Invoke(jarPath string, args []string) (OperationResult, error)`

| Aspecto | Descrição |
|---|---|
| Responsabilidade | Executar `java -jar` com argumentos. |
| Entrada | Caminho do `.jar` e argumentos. |
| Saída | Resultado capturado de stdout/stderr. |
| Passos | Montar comando; executar processo; capturar saída; mapear retorno. |
| Erros tratados | Java não encontrado, jar ausente, código de saída diferente de zero, timeout de execução. |

#### `AssinadorHttpClient.Sign(request SignCommand) (OperationResult, error)`

| Aspecto | Descrição |
|---|---|
| Responsabilidade | Enviar requisição HTTP para criação de assinatura. |
| Entrada | Dados de assinatura. |
| Saída | Resultado da assinatura simulada. |
| Endpoint | `POST /sign`. |
| Erros tratados | Conexão recusada, timeout, resposta inválida, status HTTP de erro. |

#### `AssinadorHttpClient.Validate(request ValidateCommand) (OperationResult, error)`

| Aspecto | Descrição |
|---|---|
| Responsabilidade | Enviar requisição HTTP para validação de assinatura. |
| Entrada | Dados de validação. |
| Saída | Resultado da validação simulada. |
| Endpoint | `POST /validate`. |
| Erros tratados | Conexão recusada, timeout, resposta inválida, status HTTP de erro. |

### 10.2 Métodos principais do CLI `simulador`

#### `StartSimulatorUseCase.Execute(config SimulatorConfig) (SimulatorStatus, error)`

| Aspecto | Descrição |
|---|---|
| Responsabilidade | Iniciar o Simulador do HubSaúde. |
| Entrada | Porta, URL de origem e configurações. |
| Saída | Status do simulador iniciado. |
| Passos | Verificar porta; resolver release; baixar jar se necessário; verificar checksum; resolver Java; iniciar processo; registrar PID; consultar status. |
| Erros tratados | Porta ocupada, download falho, checksum inválido, Java ausente, processo não iniciado. |

#### `StopSimulatorUseCase.Execute(port int) error`

| Aspecto | Descrição |
|---|---|
| Responsabilidade | Encerrar o Simulador do HubSaúde. |
| Entrada | Porta padrão ou informada. |
| Saída | Confirmação ou erro. |
| Passos | Consultar registro; chamar `/shutdown`; aguardar encerramento; atualizar registro. |
| Erros tratados | Simulador não encontrado, endpoint indisponível, falha no encerramento. |

#### `StatusSimulatorUseCase.Execute(port int) (SimulatorStatus, error)`

| Aspecto | Descrição |
|---|---|
| Responsabilidade | Consultar o status do Simulador. |
| Entrada | Porta padrão ou informada. |
| Saída | Status, PID e porta. |
| Passos | Consultar registro; chamar `/api/info`; validar resposta; retornar status. |
| Erros tratados | Registro inexistente, processo inativo, endpoint indisponível. |

#### `ReleaseResolver.Resolve(source string) (ReleaseMetadata, error)`

| Aspecto | Descrição |
|---|---|
| Responsabilidade | Identificar versão e URL do artefato mais recente. |
| Entrada | URL padrão ou alternativa. |
| Saída | `ReleaseMetadata`. |
| Passos | Baixar `release.json`; interpretar conteúdo; comparar versão local; retornar metadados. |
| Erros tratados | URL inválida, JSON inválido, versão ausente, rede indisponível. |

#### `ChecksumVerifier.Verify(filePath string, expected string) error`

| Aspecto | Descrição |
|---|---|
| Responsabilidade | Verificar a integridade de arquivo baixado. |
| Entrada | Caminho do arquivo e hash esperado. |
| Saída | Sucesso ou erro. |
| Passos | Calcular SHA256; comparar com valor esperado; rejeitar divergência. |
| Erros tratados | Arquivo inexistente, leitura falha, hash divergente. |

### 10.3 Métodos principais do `assinador.jar`

#### `SignatureController.sign(SignRequest request): SignResponse`

| Aspecto | Descrição |
|---|---|
| Responsabilidade | Receber requisição HTTP de criação de assinatura. |
| Entrada | `SignRequest`. |
| Saída | `SignResponse`. |
| Passos | Validar requisição; chamar `SignatureService.sign`; retornar resposta padronizada. |
| Erros tratados | Parâmetros inválidos, falha no serviço, dispositivo indisponível. |

#### `SignatureController.validate(ValidateRequest request): ValidateResponse`

| Aspecto | Descrição |
|---|---|
| Responsabilidade | Receber requisição HTTP de validação de assinatura. |
| Entrada | `ValidateRequest`. |
| Saída | `ValidateResponse`. |
| Passos | Validar requisição; chamar `SignatureService.validate`; retornar resultado. |
| Erros tratados | Parâmetros inválidos, assinatura ausente, falha de validação simulada. |

#### `ParameterValidator.validateSignRequest(SignRequest request): ValidationResult`

| Aspecto | Descrição |
|---|---|
| Responsabilidade | Validar parâmetros de criação de assinatura. |
| Entrada | `SignRequest`. |
| Saída | `ValidationResult`. |
| Validações | Presença de campos obrigatórios, formato, consistência e valores permitidos. |

#### `ParameterValidator.validateValidateRequest(ValidateRequest request): ValidationResult`

| Aspecto | Descrição |
|---|---|
| Responsabilidade | Validar parâmetros de validação de assinatura. |
| Entrada | `ValidateRequest`. |
| Saída | `ValidationResult`. |
| Validações | Presença de documento, assinatura e parâmetros obrigatórios; formato dos campos. |

#### `FakeSignatureService.sign(SignRequest request): SignResponse`

| Aspecto | Descrição |
|---|---|
| Responsabilidade | Gerar resposta simulada de assinatura. |
| Entrada | `SignRequest` válido. |
| Saída | `SignResponse` com assinatura simulada. |
| Observação | Não deve executar assinatura criptográfica real. |

#### `FakeSignatureService.validate(ValidateRequest request): ValidateResponse`

| Aspecto | Descrição |
|---|---|
| Responsabilidade | Gerar resultado simulado de validação. |
| Entrada | `ValidateRequest` válido. |
| Saída | `ValidateResponse` indicando válido ou inválido por critério predeterminado. |
| Observação | Não deve executar validação criptográfica real. |

#### `PKCS11Adapter.checkAvailability(): boolean`

| Aspecto | Descrição |
|---|---|
| Responsabilidade | Verificar disponibilidade de dispositivo criptográfico ou simulador. |
| Entrada | Configuração PKCS#11. |
| Saída | `true` se disponível; `false` ou exceção controlada se ausente. |
| Erros tratados | Biblioteca ausente, token não conectado, configuração inválida. |

---

## 11. Projeto dos dados

### 11.1 Entidades

As entidades principais do sistema são apresentadas a seguir.

| Entidade | Descrição | Origem |
|---|---|---|
| `SignCommand` | Dados recebidos pelo CLI para criação de assinatura simulada. | CLI `assinatura`. |
| `ValidateCommand` | Dados recebidos pelo CLI para validação de assinatura simulada. | CLI `assinatura`. |
| `SignRequest` | Requisição interna ou HTTP para criação de assinatura no `assinador.jar`. | `assinador.jar`. |
| `ValidateRequest` | Requisição interna ou HTTP para validação no `assinador.jar`. | `assinador.jar`. |
| `SignResponse` | Resultado simulado de criação de assinatura. | `assinador.jar`. |
| `ValidateResponse` | Resultado simulado de validação. | `assinador.jar`. |
| `ProcessMetadata` | Registro de processo em execução. | CLIs. |
| `JavaRuntime` | Informações sobre JDK/JRE disponível. | Módulo Java Runtime. |
| `ReleaseMetadata` | Informações sobre artefato remoto. | Módulo de artefatos. |
| `SimulatorStatus` | Estado atual do Simulador do HubSaúde. | CLI `simulador`. |
| `ErrorResponse` | Representação padronizada de erro. | CLIs e `assinador.jar`. |

### 11.2 Campos

#### 11.2.1 `SignCommand`

| Campo | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `documento` | string | Sim | Caminho ou identificador do documento a ser assinado. |
| `certificado` | string | Condicional | Caminho ou identificador do certificado, quando exigido. |
| `algoritmo` | string | Não | Algoritmo conceitual informado para simulação. |
| `modo` | string | Não | Modo de execução: `local` ou `server`. |
| `porta` | int | Não | Porta do `assinador.jar` em modo servidor. |
| `timeout` | int | Não | Tempo de inatividade para encerramento automático. |

#### 11.2.2 `ValidateCommand`

| Campo | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `documento` | string | Sim | Caminho ou identificador do documento associado à assinatura. |
| `assinatura` | string | Sim | Caminho ou valor da assinatura simulada. |
| `certificado` | string | Condicional | Certificado usado na validação, quando exigido. |
| `modo` | string | Não | Modo de execução: `local` ou `server`. |
| `porta` | int | Não | Porta do `assinador.jar` em modo servidor. |

#### 11.2.3 `SignRequest`

| Campo | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `document` | string | Sim | Conteúdo, caminho ou identificador do documento. |
| `certificate` | string | Condicional | Certificado ou identificador associado. |
| `parameters` | map/string | Não | Parâmetros adicionais da operação. |

#### 11.2.4 `SignResponse`

| Campo | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `success` | boolean | Sim | Indica se a operação simulada foi bem-sucedida. |
| `operation` | string | Sim | Nome da operação, por exemplo `sign`. |
| `signature` | string | Sim em sucesso | Assinatura simulada retornada. |
| `message` | string | Sim | Mensagem explicativa. |
| `errors` | list | Não | Lista de erros, quando houver. |

#### 11.2.5 `ValidateRequest`

| Campo | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `document` | string | Sim | Documento associado à assinatura. |
| `signature` | string | Sim | Assinatura a validar. |
| `certificate` | string | Condicional | Certificado ou identificador usado na validação. |
| `parameters` | map/string | Não | Parâmetros adicionais. |

#### 11.2.6 `ValidateResponse`

| Campo | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `success` | boolean | Sim | Indica se a requisição foi processada. |
| `operation` | string | Sim | Nome da operação, por exemplo `validate`. |
| `valid` | boolean | Sim em sucesso | Resultado simulado de validade. |
| `message` | string | Sim | Mensagem explicativa. |
| `errors` | list | Não | Lista de erros, quando houver. |

#### 11.2.7 `ProcessMetadata`

| Campo | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `application` | string | Sim | Nome da aplicação: `assinador` ou `simulador`. |
| `pid` | int | Sim | Identificador do processo. |
| `port` | int | Sim | Porta usada pelo processo. |
| `mode` | string | Sim | Modo de execução, por exemplo `server`. |
| `startedAt` | datetime | Sim | Data/hora de início. |
| `healthEndpoint` | string | Não | Endpoint usado para health check. |
| `status` | string | Não | Estado conhecido do processo. |

#### 11.2.8 `ReleaseMetadata`

| Campo | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `artifact` | string | Sim | Nome do artefato. |
| `version` | string | Sim | Versão do artefato. |
| `url` | string | Sim | URL de download. |
| `checksumSha256` | string | Condicional | Hash esperado do arquivo. |
| `downloadedAt` | datetime | Não | Data/hora do download local. |

### 11.3 Relacionamentos

| Relacionamento | Descrição |
|---|---|
| `SignCommand` → `SignRequest` | O CLI converte parâmetros do terminal em requisição para o `assinador.jar`. |
| `ValidateCommand` → `ValidateRequest` | O CLI converte parâmetros de validação em requisição para o `assinador.jar`. |
| `SignRequest` → `ParameterValidator` | Toda requisição de assinatura deve ser validada antes da simulação. |
| `ValidateRequest` → `ParameterValidator` | Toda requisição de validação deve ser validada antes da resposta simulada. |
| `ParameterValidator` → `FakeSignatureService` | Apenas requisições válidas seguem para o serviço de simulação. |
| `FakeSignatureService` → `SignResponse` / `ValidateResponse` | O serviço gera respostas simuladas padronizadas. |
| `ProcessMetadata` → `ProcessRegistry` | O registro de processo é salvo e consultado por repositório local. |
| `ReleaseMetadata` → `ArtifactDownloader` | O downloader usa metadados para obter artefatos. |
| `ReleaseMetadata` → `ChecksumVerifier` | O checksum esperado é usado para verificar arquivo baixado. |
| `JavaRuntime` → `LocalJarInvoker` | O invocador usa o executável Java resolvido pelo runtime. |

### 11.4 Restrições

| Restrição | Descrição |
|---|---|
| Parâmetros obrigatórios | Requisições sem campos obrigatórios devem ser rejeitadas. |
| Formato de porta | Porta deve ser número inteiro válido e disponível. |
| Diretório local | Dados operacionais devem ser armazenados em `~/.hubsaude/`. |
| Sem persistência de assinatura | Assinaturas simuladas não devem ser armazenadas permanentemente. |
| Versão Java | O `assinador.jar` deve ser executado com Java 21 compatível. |
| Arquitetura suportada | Binários devem ser gerados para Windows, Linux e macOS em amd64. |
| Checksum | Artefatos baixados devem ser aceitos apenas se checksum for válido, quando disponível. |
| Processo ativo | Registro de PID não basta; deve ser confirmado por health check quando possível. |
| Modo servidor | Deve usar HTTP para `/sign` e `/validate`. |
| Simulador | A porta padrão do Simulador do HubSaúde é 8443. |

---

## 12. Projeto das interfaces ou endpoints

### 12.1 Interface CLI `assinatura`

#### 12.1.1 Comando `assinatura version`

```bash
assinatura version
```

| Item | Descrição |
|---|---|
| Finalidade | Exibir a versão atual do CLI. |
| Entrada | Nenhuma obrigatória. |
| Saída esperada | Versão do binário, por exemplo `assinatura version v0.1.0`. |
| Requisitos relacionados | RF-002, US-01.1 |

#### 12.1.2 Comando `assinatura sign`

```bash
assinatura sign --documento <arquivo> --certificado <arquivo> [--local] [--port <porta>] [--timeout <minutos>]
```

| Item | Descrição |
|---|---|
| Finalidade | Criar assinatura digital simulada. |
| Entrada obrigatória | Documento e demais parâmetros definidos para o caso de uso. |
| Saída esperada | Resultado formatado com assinatura simulada. |
| Modos | Local via `java -jar` ou servidor via HTTP. |
| Requisitos relacionados | RF-003, RF-005, RF-006, RF-008, RF-011, RF-019 |

#### 12.1.3 Comando `assinatura validate`

```bash
assinatura validate --documento <arquivo> --assinatura <arquivo> [--local] [--port <porta>]
```

| Item | Descrição |
|---|---|
| Finalidade | Validar assinatura simulada. |
| Entrada obrigatória | Documento e assinatura. |
| Saída esperada | Resultado simulado indicando assinatura válida ou inválida. |
| Modos | Local via `java -jar` ou servidor via HTTP. |
| Requisitos relacionados | RF-004, RF-005, RF-006, RF-008, RF-011, RF-022 |

#### 12.1.4 Comando `assinatura stop`

```bash
assinatura stop [--port <porta>]
```

| Item | Descrição |
|---|---|
| Finalidade | Encerrar `assinador.jar` em modo servidor. |
| Entrada | Porta opcional. |
| Saída esperada | Confirmação de encerramento ou mensagem de processo não encontrado. |
| Requisitos relacionados | RF-015 |

### 12.2 Interface CLI `simulador`

#### 12.2.1 Comando `simulador start`

```bash
simulador start [--port <porta>] [--source <url>]
```

| Item | Descrição |
|---|---|
| Finalidade | Iniciar o Simulador do HubSaúde. |
| Entrada | Porta opcional e URL alternativa opcional. |
| Saída esperada | Status, porta e PID do simulador. |
| Requisitos relacionados | RF-030, RF-031, RF-035, RF-036 |

#### 12.2.2 Comando `simulador stop`

```bash
simulador stop [--port <porta>]
```

| Item | Descrição |
|---|---|
| Finalidade | Parar o Simulador do HubSaúde. |
| Entrada | Porta opcional. |
| Saída esperada | Confirmação de encerramento. |
| Requisitos relacionados | RF-032 |

#### 12.2.3 Comando `simulador status`

```bash
simulador status [--port <porta>]
```

| Item | Descrição |
|---|---|
| Finalidade | Consultar status do Simulador do HubSaúde. |
| Entrada | Porta opcional. |
| Saída esperada | Indicação de execução, porta e PID. |
| Requisitos relacionados | RF-033 |

### 12.3 Endpoints HTTP do `assinador.jar`

#### 12.3.1 `POST /sign`

| Item | Descrição |
|---|---|
| Finalidade | Criar assinatura simulada. |
| Método | POST |
| Entrada | JSON compatível com `SignRequest`. |
| Saída | JSON compatível com `SignResponse`. |
| Sucesso | HTTP 200. |
| Erro de parâmetro | HTTP 400. |
| Erro interno | HTTP 500. |

Exemplo conceitual de requisição:

```json
{
  "document": "documento.json",
  "certificate": "certificado.pem",
  "parameters": {
    "algorithm": "simulado"
  }
}
```

Exemplo conceitual de resposta:

```json
{
  "success": true,
  "operation": "sign",
  "signature": "assinatura-simulada",
  "message": "Assinatura simulada criada com sucesso."
}
```

#### 12.3.2 `POST /validate`

| Item | Descrição |
|---|---|
| Finalidade | Validar assinatura simulada. |
| Método | POST |
| Entrada | JSON compatível com `ValidateRequest`. |
| Saída | JSON compatível com `ValidateResponse`. |
| Sucesso | HTTP 200. |
| Erro de parâmetro | HTTP 400. |
| Erro interno | HTTP 500. |

Exemplo conceitual de requisição:

```json
{
  "document": "documento.json",
  "signature": "assinatura-simulada"
}
```

Exemplo conceitual de resposta:

```json
{
  "success": true,
  "operation": "validate",
  "valid": true,
  "message": "Assinatura simulada considerada válida."
}
```

#### 12.3.3 `GET /health`

| Item | Descrição |
|---|---|
| Finalidade | Verificar se o `assinador.jar` está ativo em modo servidor. |
| Método | GET |
| Entrada | Nenhuma. |
| Saída | Status simples de saúde da aplicação. |
| Uso | Health check pelo CLI antes de reutilizar instância ativa. |

### 12.4 Endpoints do Simulador do HubSaúde

| Método | Endpoint | Finalidade | Uso pelo CLI |
|---|---|---|---|
| GET | `/api/info` | Consultar informações e status do simulador. | Usado por `simulador status`. |
| POST ou GET | `/shutdown` | Solicitar encerramento do simulador, conforme implementação disponível. | Usado por `simulador stop`. |

### 12.5 Interface com PKCS#11

| Item | Descrição |
|---|---|
| Finalidade | Permitir interação com token, smart card ou simulador criptográfico. |
| Implementação | `PKCS11Adapter`. |
| Tecnologia | Java `SunPKCS11` provider ou equivalente. |
| Cenários tratados | Dispositivo disponível, dispositivo ausente, biblioteca ausente, configuração inválida. |
| Observação | A assinatura continua simulada no escopo atual, salvo evolução futura. |

---

## 13. Fluxos principais do sistema

### 13.1 Fluxo de criação de assinatura em modo local

```text
1. Usuário executa `assinatura sign --local ...`.
2. `SignCommandHandler` recebe os parâmetros.
3. `CliArgumentValidator` valida parâmetros básicos.
4. `SignUseCase` solicita Java ao `JavaRuntimeResolver`.
5. `JavaRuntimeResolver` localiza ou provisiona JDK/JRE.
6. `LocalJarInvoker` monta comando `java -jar assinador.jar sign ...`.
7. `assinador.jar` recebe os argumentos pelo `CliAdapter`.
8. `ParameterValidator` valida parâmetros.
9. `FakeSignatureService` gera assinatura simulada.
10. `assinador.jar` retorna resposta.
11. `OutputFormatter` apresenta resultado ao usuário.
```

### 13.2 Fluxo de criação de assinatura em modo servidor

```text
1. Usuário executa `assinatura sign ...`.
2. CLI verifica se existe servidor `assinador.jar` registrado.
3. `HealthChecker` confirma se o servidor responde.
4. Se servidor estiver ativo, `AssinadorHttpClient` envia POST `/sign`.
5. Se servidor não estiver ativo, o CLI pode iniciar o servidor ou usar fallback local, conforme configuração.
6. `SignatureController` recebe a requisição.
7. `ParameterValidator` valida parâmetros.
8. `FakeSignatureService` gera assinatura simulada.
9. `SignatureController` retorna JSON de sucesso.
10. CLI formata e apresenta a resposta.
```

### 13.3 Fluxo de validação de assinatura

```text
1. Usuário executa `assinatura validate ...`.
2. CLI valida argumentos básicos.
3. CLI escolhe modo local ou HTTP.
4. `assinador.jar` recebe a requisição.
5. `ParameterValidator` valida documento e assinatura.
6. `FakeSignatureService` aplica critério predeterminado de validação simulada.
7. Resposta indica assinatura válida ou inválida.
8. CLI apresenta resultado legível ao usuário.
```

### 13.4 Fluxo de inicialização do `assinador.jar` em modo servidor

```text
1. Usuário executa operação que requer servidor ou comando específico de inicialização.
2. CLI verifica porta padrão ou porta informada.
3. CLI verifica se já existe instância registrada.
4. Se houver registro, executa health check.
5. Se não houver instância ativa, resolve Java compatível.
6. CLI inicia `assinador.jar` em segundo plano.
7. `ProcessRegistry` salva PID, porta e endpoint de saúde.
8. CLI informa ao usuário que o servidor foi iniciado.
```

### 13.5 Fluxo de parada do `assinador.jar`

```text
1. Usuário executa `assinatura stop`.
2. CLI consulta registro local pela porta.
3. CLI verifica se o processo ainda está ativo.
4. CLI tenta encerrar o processo de forma controlada.
5. Registro local é atualizado ou removido.
6. Usuário recebe confirmação.
```

### 13.6 Fluxo de início do Simulador do HubSaúde

```text
1. Usuário executa `simulador start`.
2. CLI verifica disponibilidade da porta 8443 ou porta informada.
3. CLI verifica se `simulador.jar` está em cache local.
4. Se não estiver, `ReleaseResolver` consulta release ou `release.json`.
5. `ArtifactDownloader` baixa o `simulador.jar`.
6. `ChecksumVerifier` valida integridade do arquivo.
7. `JavaRuntimeResolver` localiza ou baixa JDK/JRE compatível.
8. `ProcessManager` inicia o simulador.
9. `ProcessRegistry` registra PID e porta.
10. `SimulatorHttpClient` consulta `/api/info`.
11. CLI exibe status ao usuário.
```

### 13.7 Fluxo de parada do Simulador

```text
1. Usuário executa `simulador stop`.
2. CLI consulta registro local.
3. CLI envia requisição ao endpoint `/shutdown`, conforme disponível.
4. CLI aguarda encerramento do processo.
5. Registro local é atualizado.
6. Usuário recebe confirmação.
```

### 13.8 Fluxo de provisionamento de JDK/JRE

```text
1. Alguma operação exige execução Java.
2. `JavaRuntimeResolver` procura Java em `~/.hubsaude/`.
3. Se não encontrar, procura Java no PATH.
4. Se encontrar, valida compatibilidade com Java 21.
5. Se não encontrar versão compatível, baixa JDK/JRE adequado.
6. Arquivo baixado é armazenado em diretório gerenciado.
7. Caminho do Java é retornado ao invocador.
```

### 13.9 Fluxo de publicação de release

```text
1. Desenvolvedor cria tag SemVer, como `v0.1.0`.
2. GitHub Actions executa workflow de release.
3. Pipeline compila binários para Windows, Linux e macOS.
4. Pipeline gera checksums SHA256.
5. Pipeline assina artefatos com Cosign.
6. Pipeline publica binários, checksums, `.sig` e `.pem` no GitHub Releases.
```

---

## 14. Regras de negócio detalhadas

| ID | Regra | Detalhamento |
|---|---|---|
| RND-001 | O sistema deve simular assinatura digital | A criação de assinatura não deve executar operação criptográfica real. O retorno deve ser uma assinatura simulada e identificável como tal. |
| RND-002 | O sistema deve simular validação de assinatura | A validação deve retornar resultado predeterminado baseado em critérios simples definidos no serviço de simulação. |
| RND-003 | Parâmetros obrigatórios devem ser validados | Operações sem parâmetros obrigatórios devem ser rejeitadas antes do processamento. |
| RND-004 | Erros devem indicar causa e correção | Mensagens devem indicar parâmetro inválido, motivo e, quando possível, sugestão de uso. |
| RND-005 | O modo servidor deve ser preferido quando disponível | Quando houver instância ativa do `assinador.jar`, o CLI deve reutilizá-la, salvo se o usuário solicitar `--local`. |
| RND-006 | O modo local deve estar disponível | O usuário deve poder forçar execução local via `--local` ou fallback equivalente. |
| RND-007 | O sistema não deve criar processos duplicados desnecessariamente | Antes de iniciar novo servidor, o CLI deve verificar registro local e health check. |
| RND-008 | Porta ocupada deve bloquear inicialização | O sistema não deve iniciar assinador ou simulador em porta já ocupada sem avisar o usuário. |
| RND-009 | A porta padrão do Simulador é 8443 | O comando `simulador start` deve usar 8443 quando nenhuma porta for informada. |
| RND-010 | O JDK/JRE deve ser reutilizado | Após download, o Java deve ser armazenado para evitar novo download em execuções futuras. |
| RND-011 | O `simulador.jar` não deve ser baixado repetidamente | Se a versão local for a mais recente, o sistema deve reutilizá-la. |
| RND-012 | Downloads devem ser verificados | Artefatos com checksum disponível devem ser verificados antes do uso. |
| RND-013 | Registros de processo devem ser atualizados | Ao iniciar ou parar processos, o sistema deve atualizar `~/.hubsaude/processos/`. |
| RND-014 | Artefatos de release devem ser verificáveis | Binários publicados devem possuir SHA256 e assinatura Cosign. |
| RND-015 | O usuário não deve precisar conhecer `java -jar` | O CLI deve encapsular os comandos Java. |
| RND-016 | Ausência de PKCS#11 deve ser tratada claramente | Falhas envolvendo token, smart card ou biblioteca PKCS#11 devem gerar erro compreensível. |
| RND-017 | Assinaturas simuladas não devem ser persistidas | Resultados de assinatura devem ser exibidos, mas não armazenados permanentemente pelo Runner. |
| RND-018 | A documentação deve deixar claro o caráter simulado | O sistema não deve ser apresentado como solução de assinatura digital juridicamente válida. |

---

## 15. Tratamento de erros

O tratamento de erros deve ser padronizado, com mensagens claras, orientativas e adequadas para terminal. A implementação deve evitar stack traces para o usuário final, exceto em modo de depuração.

### 15.1 Estrutura padronizada de erro

#### Em CLI

```text
Erro: <descrição resumida>
Motivo: <causa identificada>
Orientação: <como o usuário pode corrigir>
Código: <código opcional do erro>
```

#### Em HTTP

```json
{
  "success": false,
  "error": {
    "code": "INVALID_PARAMETER",
    "message": "Parâmetro obrigatório ausente: documento",
    "details": "Informe o parâmetro documento ou consulte a ajuda do comando."
  }
}
```

### 15.2 Códigos de erro sugeridos

| Código | Situação |
|---|---|
| `INVALID_PARAMETER` | Parâmetro ausente, vazio ou em formato inválido. |
| `JAVA_NOT_FOUND` | Java compatível não encontrado. |
| `JAVA_DOWNLOAD_FAILED` | Falha ao baixar JDK/JRE. |
| `JAR_NOT_FOUND` | Arquivo `.jar` não encontrado. |
| `JAR_EXECUTION_FAILED` | Falha ao executar `java -jar`. |
| `HTTP_CONNECTION_FAILED` | Falha de conexão com servidor local. |
| `HTTP_INVALID_RESPONSE` | Resposta HTTP inválida ou inesperada. |
| `PORT_UNAVAILABLE` | Porta ocupada ou inválida. |
| `PROCESS_NOT_FOUND` | Processo não encontrado ou já encerrado. |
| `PROCESS_REGISTRY_STALE` | Registro de processo desatualizado. |
| `DOWNLOAD_FAILED` | Falha ao baixar artefato remoto. |
| `CHECKSUM_MISMATCH` | Checksum do arquivo baixado não confere. |
| `PKCS11_DEVICE_NOT_FOUND` | Token ou smart card não encontrado. |
| `PKCS11_CONFIGURATION_ERROR` | Configuração PKCS#11 inválida. |
| `INTERNAL_ERROR` | Erro inesperado interno. |

### 15.3 Tratamento por camada

| Camada | Responsabilidade no tratamento de erro |
|---|---|
| CLI | Capturar erro, formatar mensagem e retornar código de saída apropriado. |
| Use case | Identificar falhas de fluxo e propagar erros de domínio. |
| Serviços de infraestrutura | Encapsular erros técnicos, como rede, processo, arquivo e download. |
| `assinador.jar` | Validar parâmetros e retornar erro estruturado. |
| Controladores HTTP | Converter exceções em HTTP status adequado. |
| Repositórios locais | Tratar falhas de leitura/escrita e permissões. |

### 15.4 Exemplos de mensagens

#### Parâmetro ausente

```text
Erro: parâmetro obrigatório ausente.
Motivo: o parâmetro 'documento' não foi informado.
Orientação: execute 'assinatura sign --help' para ver os parâmetros necessários.
```

#### Porta ocupada

```text
Erro: não foi possível iniciar o simulador.
Motivo: a porta 8443 já está em uso.
Orientação: finalize o processo atual ou informe outra porta com '--port'.
```

#### Java ausente

```text
Erro: Java 21 não encontrado.
Motivo: não foi localizado JDK/JRE compatível no sistema.
Orientação: conecte-se à internet para permitir o provisionamento automático ou instale Java 21 manualmente.
```

#### Checksum inválido

```text
Erro: o arquivo baixado não passou na verificação de integridade.
Motivo: o checksum SHA256 calculado é diferente do esperado.
Orientação: remova o arquivo em cache e tente novamente.
```

---

## 16. Projeto de segurança

O projeto de segurança deve considerar a distribuição de artefatos, downloads, execução local, dispositivos criptográficos e armazenamento de dados operacionais.

### 16.1 Segurança dos artefatos de release

Cada release deve conter:

```text
<artefato>
<artefato>.sig
<artefato>.pem
checksums.txt
```

Exemplo:

```text
assinatura-1.0.0-linux-amd64
assinatura-1.0.0-linux-amd64.sig
assinatura-1.0.0-linux-amd64.pem
checksums.txt
```

A assinatura deve ser feita com Cosign, usando identidade OIDC e registro em transparency log, conforme definido nos documentos do projeto.

### 16.2 Verificação de integridade

O sistema deve utilizar SHA256 para verificar artefatos baixados quando houver checksum disponível.

Fluxo:

```text
1. Baixar artefato.
2. Calcular SHA256 local.
3. Comparar com SHA256 esperado.
4. Aceitar arquivo somente se os valores coincidirem.
5. Em caso de divergência, remover ou isolar arquivo e informar erro.
```

### 16.3 Segurança no provisionamento de JDK/JRE

- Preferir download por HTTPS.
- Usar fontes confiáveis, como Eclipse Temurin / Adoptium.
- Armazenar o runtime em diretório gerenciado.
- Evitar baixar novamente versões já disponíveis.
- Informar ao usuário quando um download for realizado.
- Tratar falhas de rede e permissões.

### 16.4 Segurança na execução de processos

- Executar apenas caminhos de `.jar` conhecidos e gerenciados.
- Não montar comandos por concatenação insegura de strings.
- Usar APIs de execução de processo com lista de argumentos.
- Validar caminhos antes da execução.
- Registrar PID e porta para controle posterior.
- Encerrar processos de forma controlada quando possível.

### 16.5 Segurança no uso de PKCS#11

- Isolar lógica PKCS#11 em `PKCS11Adapter`.
- Não registrar PINs, senhas ou dados sensíveis em logs.
- Tratar ausência de dispositivo com erro claro.
- Permitir testes com SoftHSM2 ou equivalente.
- Manter documentação clara sobre configuração necessária.
- Reforçar que a operação de assinatura é simulada no escopo atual.

### 16.6 Segurança de dados locais

O Sistema Runner não deve persistir assinaturas digitais nem resultados sensíveis. A persistência local deve se limitar a:

- JDK/JRE;
- arquivos `.jar`;
- metadados de versão;
- registros de processo;
- cache de releases;
- logs operacionais, quando necessários.

### 16.7 Limites de segurança

O sistema não deve ser tratado como uma solução completa de assinatura digital real. Ele não substitui infraestrutura criptográfica de produção, autoridade certificadora ou solução de assinatura juridicamente válida. Essa limitação deve estar explícita na documentação de uso.

---

## 17. Projeto de persistência

O Sistema Runner não utiliza banco de dados tradicional. A persistência é local, baseada em arquivos no diretório `~/.hubsaude/`.

### 17.1 Diretório local gerenciado

Estrutura sugerida:

```text
~/.hubsaude/
  jdk/
    java-21/
  jre/
    java-21/
  assinador/
    assinador.jar
    metadata.json
  simulador/
    simulador.jar
    metadata.json
  processos/
    assinador-8080.json
    simulador-8443.json
  cache/
    release.json
    checksums.txt
  logs/
    runner.log
```

### 17.2 Arquivos persistidos

| Arquivo | Finalidade |
|---|---|
| `metadata.json` | Guardar versão, URL de origem e checksum de artefato. |
| `assinador-<porta>.json` | Registrar processo do `assinador.jar`. |
| `simulador-<porta>.json` | Registrar processo do `simulador.jar`. |
| `release.json` | Cache de metadados remotos. |
| `checksums.txt` | Hashes para verificação de integridade. |
| `runner.log` | Log operacional, se implementado. |

### 17.3 Exemplo de registro do assinador

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

### 17.4 Exemplo de registro do simulador

```json
{
  "application": "simulador",
  "pid": 15890,
  "port": 8443,
  "mode": "server",
  "startedAt": "2026-05-07T20:40:00Z",
  "infoEndpoint": "http://localhost:8443/api/info",
  "shutdownEndpoint": "http://localhost:8443/shutdown",
  "status": "running"
}
```

### 17.5 Políticas de persistência

| Política | Descrição |
|---|---|
| Reutilização | JDK/JRE e `.jar` baixados devem ser reutilizados. |
| Atualização | Metadados devem ser atualizados após download ou mudança de versão. |
| Limpeza | Registros de processo devem ser removidos ou marcados como inativos após parada. |
| Validação | Registros de processo devem ser confirmados por health check. |
| Segurança | Dados sensíveis não devem ser gravados. |
| Portabilidade | Caminhos devem considerar diferenças entre Windows, Linux e macOS. |

---

## 18. Projeto de testes

O projeto de testes deve cobrir as funcionalidades principais, cenários de erro, integração entre CLIs e aplicações Java, endpoints HTTP, gerenciamento de processos e distribuição multiplataforma.

### 18.1 Tipos de teste

| Tipo de teste | Objetivo |
|---|---|
| Teste unitário | Verificar funções e classes isoladas. |
| Teste de integração | Verificar comunicação entre CLI, Java, HTTP, arquivos e processos. |
| Teste de aceitação | Validar critérios de aceitação das histórias de usuário. |
| Teste de contrato HTTP | Validar formato de requisições e respostas dos endpoints. |
| Teste multiplataforma | Verificar builds e execução em Windows, Linux e macOS. |
| Teste de segurança de artefatos | Verificar checksums e assinatura Cosign. |
| Teste de erro | Confirmar mensagens claras em cenários de falha. |

### 18.2 Testes unitários sugeridos — Go

| ID | Módulo | Caso de teste |
|---|---|---|
| TU-GO-001 | `CliArgumentValidator` | Deve rejeitar `sign` sem documento obrigatório. |
| TU-GO-002 | `CliArgumentValidator` | Deve rejeitar `validate` sem assinatura. |
| TU-GO-003 | `PortChecker` | Deve identificar porta disponível. |
| TU-GO-004 | `PortChecker` | Deve identificar porta ocupada. |
| TU-GO-005 | `ProcessRegistry` | Deve salvar e ler registro de processo. |
| TU-GO-006 | `ProcessRegistry` | Deve remover registro após parada. |
| TU-GO-007 | `ChecksumVerifier` | Deve aceitar arquivo com SHA256 correto. |
| TU-GO-008 | `ChecksumVerifier` | Deve rejeitar arquivo com SHA256 divergente. |
| TU-GO-009 | `ReleaseResolver` | Deve interpretar `release.json` válido. |
| TU-GO-010 | `ErrorHandler` | Deve converter erro técnico em mensagem orientativa. |

### 18.3 Testes unitários sugeridos — Java

| ID | Módulo | Caso de teste |
|---|---|---|
| TU-JAVA-001 | `ParameterValidator` | Deve aceitar `SignRequest` válido. |
| TU-JAVA-002 | `ParameterValidator` | Deve rejeitar `SignRequest` sem documento. |
| TU-JAVA-003 | `ParameterValidator` | Deve aceitar `ValidateRequest` válido. |
| TU-JAVA-004 | `ParameterValidator` | Deve rejeitar `ValidateRequest` sem assinatura. |
| TU-JAVA-005 | `FakeSignatureService` | Deve retornar assinatura simulada para entrada válida. |
| TU-JAVA-006 | `FakeSignatureService` | Deve retornar resultado simulado de validação. |
| TU-JAVA-007 | `GlobalExceptionHandler` | Deve converter erro de parâmetro em resposta padronizada. |
| TU-JAVA-008 | `PKCS11Adapter` | Deve retornar erro claro quando dispositivo não estiver disponível. |

### 18.4 Testes de integração

| ID | Caso de teste | Resultado esperado |
|---|---|---|
| TI-001 | CLI `assinatura` invoca `assinador.jar` via `java -jar`. | Resposta simulada é capturada e exibida. |
| TI-002 | CLI `assinatura` envia `POST /sign`. | Endpoint retorna assinatura simulada. |
| TI-003 | CLI `assinatura` envia `POST /validate`. | Endpoint retorna validação simulada. |
| TI-004 | CLI detecta servidor registrado e ativo. | Reutiliza instância existente. |
| TI-005 | CLI detecta registro desatualizado. | Não reutiliza processo inativo e informa adequadamente. |
| TI-006 | CLI `simulador` inicia `simulador.jar`. | Processo é registrado com PID e porta. |
| TI-007 | CLI `simulador` consulta `/api/info`. | Status é exibido ao usuário. |
| TI-008 | CLI `simulador` chama `/shutdown`. | Processo é encerrado e registro atualizado. |
| TI-009 | Provisionamento de JDK/JRE ausente. | Sistema baixa ou informa falha de forma controlada. |
| TI-010 | Download de `simulador.jar` ausente. | Sistema baixa, verifica e armazena em cache. |

### 18.5 Testes de aceitação por história

| História | Testes associados |
|---|---|
| US-01 | TI-001, TI-002, TI-003, TI-004, TI-005 |
| US-02 | TU-JAVA-001 a TU-JAVA-008, TI-002, TI-003 |
| US-03 | TI-006, TI-007, TI-008, TI-010 |
| US-04 | TI-009 |
| US-05 | Testes de build, release, checksum e Cosign |

### 18.6 Testes de CI/CD

| ID | Caso de teste | Resultado esperado |
|---|---|---|
| TCICD-001 | Build em push na branch principal. | Workflow executa com sucesso. |
| TCICD-002 | Cross-compilation para Windows. | Binário Windows é gerado. |
| TCICD-003 | Cross-compilation para Linux. | Binário Linux é gerado. |
| TCICD-004 | Cross-compilation para macOS. | Binário macOS é gerado. |
| TCICD-005 | Criação de tag SemVer. | Release é criada automaticamente. |
| TCICD-006 | Geração de checksums. | Arquivo de checksums é publicado. |
| TCICD-007 | Assinatura Cosign. | `.sig` e `.pem` são publicados para os artefatos. |

### 18.7 Critérios mínimos de qualidade dos testes

- Toda validação de parâmetro deve possuir teste de sucesso e teste de falha.
- Todo comando CLI principal deve possuir teste de comportamento.
- Todo endpoint HTTP deve possuir teste de contrato.
- Todo fluxo que manipula processo deve testar cenário ativo e inativo.
- Todo download com checksum deve testar sucesso e falha.
- Os testes devem ser executáveis em pipeline de CI/CD.
- As mensagens de erro devem ser verificadas nos cenários de falha mais importantes.

---

## 19. Rastreabilidade

### 19.1 Rastreabilidade entre histórias, requisitos e módulos

| História | Requisitos principais | Módulos principais |
|---|---|---|
| US-01.1 | RF-001, RF-002 | `cmd/assinatura`, `commands`, `version` |
| US-01.2 | RF-003, RF-004, RF-005 | `sign.go`, `validate.go`, `CliArgumentValidator` |
| US-01.3 | RF-006, RF-007 | `LocalJarInvoker`, `JavaRuntimeResolver` |
| US-01.4 | RF-008 | `OutputFormatter`, `ErrorHandler` |
| US-01.5 | RF-009, RF-010 | `ProcessManager`, `ProcessRegistry` |
| US-01.6 | RF-011, RF-012, RF-013 | `AssinadorHttpClient`, `HealthChecker` |
| US-01.7 | RF-014 | `ProcessRegistry`, `HealthChecker` |
| US-01.8 | RF-015 | `StopAssinadorUseCase`, `ProcessManager` |
| US-01.9 | RF-016 | `ProcessManager`, mecanismo de timeout |
| US-02.1 | RF-017, RF-018, RF-019 | `SignatureService`, `FakeSignatureService` |
| US-02.2 | RF-020, RF-021 | `ParameterValidator`, `InvalidParameterException` |
| US-02.3 | RF-022, RF-023 | `FakeSignatureService`, `ParameterValidator` |
| US-02.4 | RF-024, RF-025, RF-026 | `SignatureController`, DTOs HTTP |
| US-02.5 | RF-027, RF-028 | `PKCS11Adapter`, `Pkcs11Config` |
| US-03.1 | RF-030, RF-031 | `StartSimulatorUseCase`, `PortChecker` |
| US-03.2 | RF-032, RF-033, RF-034 | `StopSimulatorUseCase`, `StatusSimulatorUseCase`, `ProcessRegistry` |
| US-03.3 | RF-029 | `cmd/simulador`, `simulador/commands` |
| US-03.4 | RF-035, RF-036, RF-037, RF-038 | `ReleaseResolver`, `ArtifactDownloader`, `ChecksumVerifier`, `CacheManager` |
| US-04.1 | RF-039, RF-040, RF-041 | `JavaDetector`, `JavaDownloader`, `JavaRuntimeResolver` |
| US-05.1 | RF-042 | `.github/workflows/build.yml` |
| US-05.2 | RF-043 | `.github/workflows/release.yml` |
| US-05.3 | RF-044, RF-045, RF-046, RF-047 | `release.yml`, Cosign, checksum scripts |

### 19.2 Rastreabilidade entre módulos e testes

| Módulo | Testes recomendados |
|---|---|
| `CliArgumentValidator` | TU-GO-001, TU-GO-002 |
| `LocalJarInvoker` | TI-001 |
| `AssinadorHttpClient` | TI-002, TI-003 |
| `ProcessRegistry` | TU-GO-005, TU-GO-006, TI-004, TI-005 |
| `PortChecker` | TU-GO-003, TU-GO-004 |
| `JavaRuntimeResolver` | TI-009 |
| `ReleaseResolver` | TU-GO-009, TI-010 |
| `ArtifactDownloader` | TI-010 |
| `ChecksumVerifier` | TU-GO-007, TU-GO-008 |
| `StartSimulatorUseCase` | TI-006 |
| `StatusSimulatorUseCase` | TI-007 |
| `StopSimulatorUseCase` | TI-008 |
| `ParameterValidator` | TU-JAVA-001 a TU-JAVA-004 |
| `FakeSignatureService` | TU-JAVA-005, TU-JAVA-006 |
| `SignatureController` | TI-002, TI-003 |
| `PKCS11Adapter` | TU-JAVA-008 |
| CI/CD | TCICD-001 a TCICD-007 |

### 19.3 Rastreabilidade entre atributos de qualidade e decisões

| Atributo | Decisões relacionadas |
|---|---|
| Segurança | SHA256, Cosign, OIDC, PKCS#11 isolado, não persistência de assinatura. |
| Desempenho | Modo servidor HTTP, reutilização de instância ativa, cache de JDK/JRE e `.jar`. |
| Manutenibilidade | Módulos separados, serviços reutilizáveis, DTOs, controladores e repositórios. |
| Usabilidade | CLI com Cobra, `--help`, mensagens claras, saída formatada. |
| Confiabilidade | Validação de parâmetros, health check, fallback local, tratamento de erros. |
| Escalabilidade | Estrutura modular para novos comandos e novas aplicações Java. |
| Portabilidade | Go, cross-compilation, GitHub Actions e suporte a Windows/Linux/macOS. |

---

## 20. Decisões técnicas principais

| ID | Decisão técnica | Justificativa | Impacto |
|---|---|---|---|
| DT-001 | Desenvolver os CLIs em Go | Go facilita geração de binários multiplataforma e reduz dependências no ambiente do usuário. | Melhora portabilidade e distribuição. |
| DT-002 | Usar Cobra para comandos CLI | Cobra oferece suporte a comandos, subcomandos, flags e help integrado. | Melhora organização e usabilidade. |
| DT-003 | Desenvolver `assinador.jar` em Java 21 | Java 21 é restrição definida pelo projeto. | Padroniza execução Java. |
| DT-004 | Separar CLI `assinatura` e CLI `simulador` | Cada CLI possui responsabilidade clara. | Reduz acoplamento e melhora clareza de uso. |
| DT-005 | Suportar modo local e modo servidor | Modo local é simples; modo servidor reduz cold start em múltiplas operações. | Melhora flexibilidade e desempenho. |
| DT-006 | Expor endpoints `/sign` e `/validate` | Padroniza integração HTTP com o `assinador.jar`. | Facilita testes e integração. |
| DT-007 | Isolar validação de parâmetros no `assinador.jar` | A validação é essencial e deve ser reutilizada no modo local e HTTP. | Evita duplicação de regras. |
| DT-008 | Usar `~/.hubsaude/` para armazenamento local | Centraliza cache, JDK/JRE, metadados e processos. | Melhora previsibilidade operacional. |
| DT-009 | Registrar PID e porta dos processos | Permite status, reutilização e encerramento. | Melhora controle operacional. |
| DT-010 | Confirmar processos por health check | Evita confiar em registros desatualizados. | Aumenta confiabilidade. |
| DT-011 | Provisionar JDK/JRE automaticamente | Atende objetivo de ocultar instalação Java do usuário. | Melhora usabilidade. |
| DT-012 | Baixar `simulador.jar` dinamicamente | Permite uso da versão mais recente sem download manual. | Melhora atualização e integração. |
| DT-013 | Verificar checksum de downloads | Evita uso de arquivos corrompidos ou adulterados. | Aumenta segurança e confiabilidade. |
| DT-014 | Assinar releases com Cosign | Garante autenticidade dos artefatos distribuídos. | Melhora segurança da cadeia de suprimentos. |
| DT-015 | Não persistir assinaturas simuladas | O sistema não é repositório de assinaturas. | Reduz risco e mantém coerência com escopo. |
| DT-016 | Isolar PKCS#11 em adaptador | Facilita testes e tratamento de variações de dispositivos. | Reduz acoplamento com tecnologia criptográfica. |
| DT-017 | Automatizar builds com GitHub Actions | Garante geração consistente de binários. | Melhora qualidade de entrega. |
| DT-018 | Usar SemVer em releases | Facilita controle de versões e comunicação com usuários. | Melhora rastreabilidade de entregas. |

---

## Apêndice A — Exemplos conceituais de comandos

### A.1 Exibir versão

```bash
assinatura version
```

Saída esperada:

```text
assinatura version v0.1.0
```

### A.2 Criar assinatura simulada

```bash
assinatura sign --documento documento.json --certificado certificado.pem
```

Saída esperada:

```text
Operação: criação de assinatura
Status: sucesso
Mensagem: Assinatura simulada criada com sucesso.
Assinatura: assinatura-simulada
```

### A.3 Validar assinatura simulada

```bash
assinatura validate --documento documento.json --assinatura assinatura.txt
```

Saída esperada:

```text
Operação: validação de assinatura
Status: sucesso
Resultado: assinatura válida
```

### A.4 Iniciar simulador

```bash
simulador start
```

Saída esperada:

```text
Verificando porta 8443...
Porta disponível.
Iniciando Simulador do HubSaúde...

Status: em execução
Porta: 8443
PID: 15342
```

### A.5 Consultar status do simulador

```bash
simulador status
```

Saída esperada:

```text
Simulador do HubSaúde: em execução
Porta: 8443
PID: 15342
```

### A.6 Parar simulador

```bash
simulador stop
```

Saída esperada:

```text
Simulador do HubSaúde encerrado com sucesso.
```

---

## Apêndice B — Exemplos conceituais de estruturas JSON

### B.1 `release.json`

```json
{
  "jar": {
    "url": "https://github.com/kyriosdata/assinador/releases/latest/download/assinador.jar",
    "version": "1.2.0"
  },
  "jre": {
    "windows_x64": "https://api.adoptium.net/v3/binary/latest/21/ga/windows/x64/jre/hotspot/normal/eclipse",
    "linux_x64": "https://api.adoptium.net/v3/binary/latest/21/ga/linux/x64/jre/hotspot/normal/eclipse",
    "mac_x64": "https://api.adoptium.net/v3/binary/latest/21/ga/mac/x64/jre/hotspot/normal/eclipse"
  }
}
```

### B.2 Resposta de erro

```json
{
  "success": false,
  "error": {
    "code": "INVALID_PARAMETER",
    "message": "Parâmetro obrigatório ausente: documento",
    "details": "Informe o documento ou consulte a ajuda do comando."
  }
}
```

### B.3 Metadados de artefato

```json
{
  "artifact": "simulador.jar",
  "version": "1.2.0",
  "url": "https://github.com/kyriosdata/assinador/releases/latest/download/simulador.jar",
  "checksumSha256": "valor-sha256-esperado",
  "downloadedAt": "2026-05-07T20:30:00Z"
}
```

---

## Referências documentais utilizadas

- Especificação original do Sistema Runner — Trabalho Prático.
- Plano Revisado #2 do Sistema Runner.
- Documento de Design do Sistema Runner com Modelo C4.
- Especificação de Requisitos de Software — Sistema Runner.
- Documento de Arquitetura de Software — Sistema Runner.
- Diagramas de Contexto e Contêineres fornecidos nos arquivos do projeto.
- Boas práticas de Engenharia de Software para projeto detalhado, modularidade, separação de responsabilidades, tratamento de erros, rastreabilidade, testes, segurança de artefatos e integração entre sistemas.
