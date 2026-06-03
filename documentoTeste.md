# Documento de Teste de Software — Sistema Runner

## 1. Identificação do documento

### 1.1 Nome do projeto

**Sistema Runner**.

### 1.2 Nome do documento

**Documento de Teste de Software**.

### 1.3 Versão

**1.0**.

### 1.4 Data

**07/05/2026**.

### 1.5 Responsáveis

| Papel | Responsável |
|---|---|
| Elaboração | Equipe do projeto / Disciplina de Implementação e Integração de Software |
| Instituição | Bacharelado em Engenharia de Software — Universidade Federal de Goiás (UFG) |
| Contexto de aplicação | Plataforma HubSaúde — interoperabilidade de dados em saúde |
| Documentos relacionados | Especificação de Requisitos de Software; Documento de Arquitetura de Software; Documento de Projeto Detalhado de Software; Documento de Modelo C4; Documento de Implementação e Integração; README da implementação |

---

## 2. Histórico de versões

| Versão | Data | Autor / Responsável | Descrição da alteração |
|---|---|---|---|
| 1.0 | 07/05/2026 | Equipe do projeto | Elaboração inicial do Documento de Teste de Software do Sistema Runner, com base nos arquivos enviados, nos documentos já elaborados e na implementação de referência acadêmica. |

---

## 3. Sumário

1. [Identificação do documento](#1-identificação-do-documento)  
2. [Histórico de versões](#2-histórico-de-versões)  
3. [Sumário](#3-sumário)  
4. [Introdução](#4-introdução)  
   4.1 [Contexto do sistema](#41-contexto-do-sistema)  
   4.2 [Finalidade do documento](#42-finalidade-do-documento)  
5. [Objetivo dos testes](#5-objetivo-dos-testes)  
   5.1 [Objetivo geral](#51-objetivo-geral)  
   5.2 [Objetivos específicos](#52-objetivos-específicos)  
6. [Escopo dos testes](#6-escopo-dos-testes)  
   6.1 [O que será testado](#61-o-que-será-testado)  
   6.2 [O que não será testado](#62-o-que-não-será-testado)  
7. [Requisitos ou funcionalidades testadas](#7-requisitos-ou-funcionalidades-testadas)  
8. [Estratégia de testes](#8-estratégia-de-testes)  
9. [Tipos de teste](#9-tipos-de-teste)  
10. [Ambiente de testes](#10-ambiente-de-testes)  
11. [Dados de teste](#11-dados-de-teste)  
12. [Casos de teste](#12-casos-de-teste)  
13. [Execução dos testes](#13-execução-dos-testes)  
14. [Registro de defeitos](#14-registro-de-defeitos)  
15. [Evidências de teste](#15-evidências-de-teste)  
16. [Critérios de aprovação](#16-critérios-de-aprovação)  
17. [Métricas de teste](#17-métricas-de-teste)  
18. [Relatório final de testes](#18-relatório-final-de-testes)  
19. [Conclusão e recomendações](#19-conclusão-e-recomendações)  
20. [Referências](#20-referências)  

---

## 4. Introdução

### 4.1 Contexto do sistema

O **Sistema Runner** é uma solução de software desenvolvida para facilitar a execução e o gerenciamento de aplicações Java por meio de comandos de terminal. O sistema está relacionado ao contexto da Plataforma HubSaúde e foi definido em uma disciplina de Implementação e Integração de Software.

A solução possui três partes principais: o CLI **`assinatura`**, desenvolvido em Go, responsável por criar e validar assinaturas simuladas; o **`assinador.jar`**, desenvolvido em Java 21, responsável por validar parâmetros e simular operações de assinatura e validação; e o CLI **`simulador`**, também desenvolvido em Go, responsável por iniciar, parar e consultar o status do Simulador do HubSaúde.

O sistema também utiliza o diretório local **`~/.hubsaude/`** para armazenar arquivos `.jar`, cache, metadados, registros de processos, artefatos e logs. A implementação prevê integração local via `java -jar`, integração HTTP, verificação de portas, registro de processos, cálculo de checksum SHA256, estrutura de CI/CD e assinatura de artefatos com Cosign.

É importante destacar que o Sistema Runner **não realiza assinatura digital criptográfica real**. A criação e a validação de assinaturas são simuladas, conforme o escopo definido nos documentos do projeto.

### 4.2 Finalidade do documento

Este documento tem a finalidade de planejar, organizar e registrar os testes do Sistema Runner. Ele define o escopo dos testes, as funcionalidades a serem testadas, a estratégia adotada, os tipos de teste, o ambiente necessário, os dados de teste, os casos de teste, os critérios de aprovação, as métricas e a estrutura de registro de defeitos e evidências.

O documento também serve como apoio para a avaliação acadêmica da implementação, pois demonstra a relação entre requisitos, funcionalidades implementadas, integração entre componentes e validação da qualidade do software.

---

## 5. Objetivo dos testes

### 5.1 Objetivo geral

Validar se o Sistema Runner atende aos requisitos definidos e se as partes implementadas funcionam de forma integrada, confiável e compreensível para o usuário, especialmente nos fluxos de criação de assinatura simulada, validação de assinatura simulada, execução local do `assinador.jar`, execução HTTP do `assinador.jar`, gerenciamento do Simulador do HubSaúde e tratamento de erros.

### 5.2 Objetivos específicos

- Verificar se os CLIs `assinatura` e `simulador` compilam corretamente.
- Verificar se o `assinador.jar` compila e executa corretamente.
- Validar os comandos `assinatura version`, `assinatura sign`, `assinatura validate`, `assinatura start` e `assinatura stop`.
- Validar os comandos `simulador start`, `simulador status` e `simulador stop`.
- Validar a invocação local do `assinador.jar` via `java -jar`.
- Validar a invocação HTTP do `assinador.jar` por meio dos endpoints `/sign`, `/validate`, `/health` e `/shutdown`.
- Validar a rejeição de parâmetros ausentes ou inválidos.
- Validar o registro local de processos em `~/.hubsaude/processos/`.
- Validar a verificação de portas disponíveis e ocupadas.
- Validar o cálculo e a conferência de checksum SHA256.
- Verificar se as mensagens de erro são claras e orientativas.
- Validar a estrutura de testes automatizados em Go e Java.
- Registrar evidências, defeitos, limitações e recomendações para melhoria.

---

## 6. Escopo dos testes

### 6.1 O que será testado

| Item | Descrição |
|---|---|
| CLI `assinatura` | Comandos `version`, `sign`, `validate`, `start` e `stop`. |
| CLI `simulador` | Comandos `start`, `status` e `stop`. |
| `assinador.jar` | Execução em modo CLI e modo servidor HTTP. |
| Invocação local | Execução do `assinador.jar` por `java -jar`. |
| Invocação HTTP | Comunicação com os endpoints `/health`, `/sign`, `/validate` e `/shutdown`. |
| Validação de parâmetros | Rejeição de parâmetros ausentes, vazios ou inválidos. |
| Simulação de assinatura | Retorno de assinatura simulada para entradas válidas. |
| Simulação de validação | Retorno de resultado simulado para entradas válidas. |
| Gerenciamento de processos | Registro, consulta e encerramento de processos iniciados pelo Runner. |
| Verificação de portas | Identificação de porta livre e porta ocupada. |
| Persistência local | Uso de arquivos em `~/.hubsaude/`. |
| Download de artefatos | Download por URL quando indicado. |
| Checksum SHA256 | Cálculo e validação da integridade de artefatos. |
| Tratamento de erros | Mensagens claras para falhas comuns. |
| Testes automatizados | Execução de `go test ./...` e `mvn test`. |
| CI/CD | Workflows de build e release. |

### 6.2 O que não será testado

| Item fora do escopo | Justificativa |
|---|---|
| Assinatura digital criptográfica real | O sistema implementa apenas simulação. |
| Validação criptográfica real de assinatura | A validação é simulada. |
| Geração de certificados digitais | Não faz parte do escopo do sistema. |
| Integração real com autoridades certificadoras | Explicitamente fora do escopo. |
| Autenticação de usuários | O sistema é local e não possui controle de usuários. |
| Interface gráfica | O sistema é exclusivamente CLI. |
| Banco de dados relacional ou não relacional | A persistência é feita por arquivos locais. |
| Garantia jurídica da assinatura | O sistema não é solução de assinatura digital juridicamente válida. |
| Implantação distribuída em produção | A proposta é execução local e integração acadêmica. |
| Teste completo com token físico real | O adaptador PKCS#11 está estruturado, mas a assinatura real está fora do escopo atual. |

---

## 7. Requisitos ou funcionalidades testadas

| ID | Requisito / funcionalidade | Descrição | Casos de teste relacionados |
|---|---|---|---|
| FT-001 | CLI `assinatura` | Executar comandos de assinatura e validação. | CT-001 a CT-013 |
| FT-002 | Comando `version` | Exibir a versão atual do CLI. | CT-001 |
| FT-003 | Comando `sign` | Criar assinatura simulada. | CT-003, CT-004, CT-009, CT-014 |
| FT-004 | Comando `validate` | Validar assinatura simulada. | CT-006, CT-007, CT-010, CT-015 |
| FT-005 | Invocação local | Executar `assinador.jar` via `java -jar`. | CT-003, CT-006 |
| FT-006 | Invocação HTTP | Usar `/sign`, `/validate` e `/health`. | CT-008 a CT-016 |
| FT-007 | Modo servidor | Iniciar e parar o `assinador.jar`. | CT-008, CT-012, CT-013 |
| FT-008 | CLI `simulador` | Iniciar, consultar e parar o simulador. | CT-017 a CT-020 |
| FT-009 | Registro de processo | Persistir PID, porta, aplicação e status. | CT-021, CT-022 |
| FT-010 | Verificação de porta | Validar porta livre e porta ocupada. | CT-020 |
| FT-011 | Checksum SHA256 | Validar integridade de artefatos. | CT-023, CT-024 |
| FT-012 | Download por URL | Baixar artefato por `--source`. | CT-025, CT-026 |
| FT-013 | Tratamento de erros | Apresentar mensagens claras. | CT-004, CT-005, CT-007, CT-011, CT-026, CT-027 |
| FT-014 | Testes automatizados | Executar testes Go e Java. | CT-028, CT-029 |
| FT-015 | Build | Compilar CLIs e JAR. | CT-030, CT-031 |
| FT-016 | CI/CD | Validar workflows de build e release. | CT-032, CT-033, CT-034 |

---

## 8. Estratégia de testes

A estratégia de testes adotada combina testes unitários, testes funcionais, testes de integração, testes de tratamento de erro, testes de instalação, testes de segurança de artefatos e testes de aceitação.

A execução deve ocorrer de forma incremental. Primeiro, devem ser validados os componentes isolados, como verificação de portas, checksum, validação de parâmetros e serviço de assinatura simulada. Em seguida, devem ser testados os comandos CLI e a comunicação entre Go e Java. Por fim, devem ser testados os fluxos completos, incluindo modo local, modo HTTP, gerenciamento de processos, simulação de assinatura, simulação de validação e comandos do simulador.

### 8.1 Ordem recomendada de execução

1. Preparar ambiente com Go, Java 21 e Maven.
2. Compilar os CLIs `assinatura` e `simulador`.
3. Compilar o `assinador.jar`.
4. Executar testes unitários Go.
5. Executar testes unitários Java.
6. Executar testes funcionais do CLI `assinatura`.
7. Executar testes de integração local com `java -jar`.
8. Executar testes de integração HTTP.
9. Executar testes do CLI `simulador`.
10. Executar testes de erro.
11. Executar testes de checksum e artefatos.
12. Registrar evidências e defeitos.

### 8.2 Critérios de rastreabilidade

Cada caso de teste deve estar associado a uma funcionalidade, requisito ou critério de aceitação. Caso um defeito seja identificado, ele deve indicar o caso de teste que falhou, o requisito afetado, a severidade, os passos para reprodução e a evidência correspondente.

---

## 9. Tipos de teste

| Tipo de teste | Objetivo | Exemplos no Sistema Runner |
|---|---|---|
| Teste unitário | Validar unidades isoladas de código. | `PortChecker`, `ChecksumVerifier`, `ParameterValidator`, `FakeSignatureService`. |
| Teste funcional | Validar funcionalidades do ponto de vista do usuário. | `assinatura sign`, `assinatura validate`, `simulador status`. |
| Teste de integração | Validar comunicação entre partes do sistema. | CLI Go invocando `assinador.jar`; CLI usando HTTP. |
| Teste de contrato HTTP | Validar formato de requisições e respostas. | `/sign`, `/validate`, `/health`, `/shutdown`. |
| Teste de erro | Validar respostas a situações inválidas. | Parâmetro ausente, JAR inexistente, porta ocupada, Java ausente. |
| Teste de instalação | Validar build e preparação do ambiente. | `go build`, `mvn clean package`, script de instalação local. |
| Teste de regressão | Garantir que correções não quebrem fluxos existentes. | Reexecutar testes automatizados e fluxos principais. |
| Teste de segurança de artefatos | Verificar integridade e autenticidade. | Checksum SHA256, arquivos `.sig` e `.pem`. |
| Teste de compatibilidade | Validar plataformas suportadas. | Windows, Linux e macOS em amd64. |
| Teste de aceitação | Validar atendimento dos critérios de aceitação. | Fluxos completos de assinatura, validação e simulador. |

---

## 10. Ambiente de testes

### 10.1 Ambiente local recomendado

| Item | Configuração recomendada |
|---|---|
| Sistema operacional | Linux, Windows ou macOS |
| Arquitetura | amd64 |
| Go | Versão compatível com o projeto; alvo documentado: Go 1.25 |
| Java | Java 21 |
| Maven | 3.9+ |
| Git | Versão atual |
| Terminal | Bash, PowerShell ou terminal equivalente |
| Diretório gerenciado | `~/.hubsaude/` |
| Porta do assinador | 8080 |
| Porta do simulador | 8443 |
| Ferramentas opcionais | Cosign, curl, jq, SoftHSM2 |

### 10.2 Ambiente CI/CD recomendado

| Item | Configuração |
|---|---|
| Plataforma | GitHub Actions |
| Sistemas de build | Ubuntu, Windows e macOS |
| Testes Go | `go test ./...` |
| Testes Java | `mvn test` |
| Build Go | `go build` para os CLIs |
| Build Java | `mvn clean package` |
| Release | GitHub Releases |
| Segurança | SHA256 e Cosign |

### 10.3 Pré-condições gerais

- Repositório do Sistema Runner disponível localmente.
- Permissão de leitura e escrita no diretório do projeto.
- Permissão de leitura e escrita em `~/.hubsaude/`.
- Java 21 instalado ou disponível em diretório gerenciado.
- `assinador.jar` compilado e disponível para os testes de integração.
- Portas 8080 e 8443 disponíveis, quando utilizadas.
- Maven instalado para execução de testes Java por `mvn test`, salvo uso de script alternativo.

---

## 11. Dados de teste

### 11.1 Arquivos de entrada

| Dado | Nome sugerido | Conteúdo / finalidade |
|---|---|---|
| Documento válido | `documento-valido.json` | Arquivo JSON simples usado em assinatura e validação simuladas. |
| Documento inválido | `documento-invalido.json` | Arquivo vazio, inexistente ou em formato incorreto. |
| Certificado simulado | `certificado-simulado.pem` | Arquivo textual de certificado fictício. |
| Assinatura válida simulada | `assinatura-valida.txt` | Contém o texto `assinatura-simulada`. |
| Assinatura inválida simulada | `assinatura-invalida.txt` | Contém valor diferente do esperado. |
| Artefato de teste | `arquivo-teste.bin` | Arquivo usado para teste de checksum. |

### 11.2 Exemplo de `documento-valido.json`

```json
{
  "id": "doc-001",
  "tipo": "teste",
  "conteudo": "Documento de teste para assinatura simulada."
}
```

### 11.3 Exemplo de `certificado-simulado.pem`

```text
-----BEGIN CERTIFICATE-----
CERTIFICADO-SIMULADO-PARA-TESTE
-----END CERTIFICATE-----
```

### 11.4 Parâmetros utilizados nos testes

| Parâmetro | Valor válido | Valor inválido |
|---|---|---|
| `--documento` | `documento-valido.json` | vazio ou arquivo inexistente |
| `--certificado` | `certificado-simulado.pem` | vazio ou arquivo inexistente |
| `--assinatura` | `assinatura-simulada` | vazio |
| `--jar` | `~/.hubsaude/assinador/assinador.jar` | caminho inexistente |
| `--port` | `8080` ou `8443` | texto, número negativo ou porta ocupada |
| `--source` | URL HTTP/HTTPS válida | URL inválida |

---

## 12. Casos de teste

> A coluna **Resultado obtido** deve ser preenchida após a execução real no ambiente da equipe. Nesta versão, os casos estão documentados como plano executável e rastreável.

| ID | Objetivo | Pré-condições | Dados de entrada | Passos de execução | Resultado esperado | Resultado obtido | Status |
|---|---|---|---|---|---|---|---|
| CT-001 | Verificar versão do CLI `assinatura`. | CLI compilado. | Nenhum. | Executar `./bin/assinatura version`. | Sistema exibe a versão atual. | A registrar. | Planejado |
| CT-002 | Verificar ajuda do CLI `assinatura`. | CLI compilado. | Nenhum. | Executar `./bin/assinatura --help`. | Sistema exibe comandos e flags. | A registrar. | Planejado |
| CT-003 | Criar assinatura simulada em modo local. | Java 21 e `assinador.jar` disponíveis. | Documento e certificado válidos. | Executar `assinatura sign --local --jar ... --documento ... --certificado ...`. | Sistema retorna assinatura simulada com sucesso. | A registrar. | Planejado |
| CT-004 | Rejeitar assinatura sem documento. | CLI compilado. | Certificado válido; documento ausente. | Executar `assinatura sign --certificado certificado.pem`. | Sistema informa parâmetro obrigatório ausente. | A registrar. | Planejado |
| CT-005 | Rejeitar assinatura sem certificado, quando exigido. | CLI compilado. | Documento válido; certificado ausente. | Executar `assinatura sign --documento documento.json`. | Sistema informa ausência do certificado ou parâmetro exigido. | A registrar. | Planejado |
| CT-006 | Validar assinatura simulada em modo local. | Java 21 e `assinador.jar` disponíveis. | Documento válido e assinatura simulada. | Executar `assinatura validate --local --jar ... --documento ... --assinatura ...`. | Sistema retorna resultado simulado de validação. | A registrar. | Planejado |
| CT-007 | Rejeitar validação sem assinatura. | CLI compilado. | Documento válido; assinatura ausente. | Executar `assinatura validate --documento documento.json`. | Sistema informa parâmetro `assinatura` ausente. | A registrar. | Planejado |
| CT-008 | Iniciar `assinador.jar` em modo servidor. | Porta 8080 livre; JAR disponível. | Porta 8080. | Executar `assinatura start --jar ... --port 8080`. | Servidor inicia, PID e porta são registrados. | A registrar. | Planejado |
| CT-009 | Criar assinatura via HTTP. | Servidor `assinador.jar` ativo. | Documento e certificado válidos. | Executar `assinatura sign --documento ... --certificado ... --port 8080`. | CLI usa `/sign` e exibe assinatura simulada. | A registrar. | Planejado |
| CT-010 | Validar assinatura via HTTP. | Servidor `assinador.jar` ativo. | Documento e assinatura válidos. | Executar `assinatura validate --documento ... --assinatura ... --port 8080`. | CLI usa `/validate` e exibe resultado. | A registrar. | Planejado |
| CT-011 | Verificar erro para JAR inexistente. | CLI compilado. | Caminho inválido de JAR. | Executar `assinatura sign --local --jar /inexistente/assinador.jar ...`. | Sistema informa que o JAR não foi encontrado. | A registrar. | Planejado |
| CT-012 | Parar `assinador.jar`. | Servidor ativo e registrado. | Porta 8080. | Executar `assinatura stop --port 8080`. | Processo é encerrado e registro atualizado. | A registrar. | Planejado |
| CT-013 | Verificar endpoint `/health`. | Servidor ativo. | Porta 8080. | Acessar `http://localhost:8080/health`. | Servidor responde indicando que está ativo. | A registrar. | Planejado |
| CT-014 | Verificar endpoint `/sign` diretamente. | Servidor ativo. | JSON de assinatura válido. | Enviar POST `/sign`. | Retorna JSON com assinatura simulada. | A registrar. | Planejado |
| CT-015 | Verificar endpoint `/validate` diretamente. | Servidor ativo. | JSON de validação válido. | Enviar POST `/validate`. | Retorna JSON com resultado simulado. | A registrar. | Planejado |
| CT-016 | Verificar erro HTTP por JSON incompleto. | Servidor ativo. | JSON sem documento. | Enviar POST `/sign` com JSON incompleto. | Retorna erro de parâmetro inválido. | A registrar. | Planejado |
| CT-017 | Iniciar simulador com JAR local. | `simulador.jar` disponível; porta 8443 livre. | Caminho do JAR e porta. | Executar `simulador start --jar ... --port 8443`. | Simulador inicia e processo é registrado. | A registrar. | Planejado |
| CT-018 | Consultar status do simulador. | Simulador ativo. | Porta 8443. | Executar `simulador status --port 8443`. | Sistema exibe status, porta e PID. | A registrar. | Planejado |
| CT-019 | Parar simulador. | Simulador ativo. | Porta 8443. | Executar `simulador stop --port 8443`. | Simulador é encerrado e registro atualizado. | A registrar. | Planejado |
| CT-020 | Bloquear início do simulador em porta ocupada. | Porta 8443 ocupada. | Porta 8443. | Executar `simulador start --port 8443`. | Sistema informa porta ocupada. | A registrar. | Planejado |
| CT-021 | Registrar processo iniciado. | Processo iniciado pelo Runner. | PID e porta. | Iniciar assinador ou simulador. | Arquivo JSON é criado em `~/.hubsaude/processos/`. | A registrar. | Planejado |
| CT-022 | Detectar registro de processo desatualizado. | Registro existe, mas processo não responde. | Registro antigo. | Executar operação que consulta processo. | Sistema não reutiliza processo inativo. | A registrar. | Planejado |
| CT-023 | Verificar checksum correto. | Arquivo de teste disponível. | Hash SHA256 correto. | Executar verificação de checksum. | Sistema aceita arquivo válido. | A registrar. | Planejado |
| CT-024 | Rejeitar checksum divergente. | Arquivo de teste disponível. | Hash SHA256 incorreto. | Executar verificação de checksum. | Sistema rejeita arquivo e informa divergência. | A registrar. | Planejado |
| CT-025 | Baixar artefato por URL. | Internet ou servidor local disponível. | URL válida. | Executar fluxo com `--source`. | Artefato é baixado e salvo em cache. | A registrar. | Planejado |
| CT-026 | Tratar falha de download. | URL inválida. | URL inválida. | Executar fluxo com `--source` inválido. | Sistema informa falha de download. | A registrar. | Planejado |
| CT-027 | Verificar mensagem para Java ausente. | Java indisponível no ambiente. | Comando que exige Java. | Executar `assinatura sign --local ...`. | Sistema informa Java ausente ou orienta configuração. | A registrar. | Planejado |
| CT-028 | Executar testes unitários Go. | Go instalado. | Nenhum. | Executar `go test ./...`. | Testes Go executam sem falhas. | A registrar. | Planejado |
| CT-029 | Executar testes unitários Java. | Java 21 e Maven instalados. | Nenhum. | Executar `cd projetos/assinador-java && mvn test`. | Testes Java executam sem falhas. | A registrar. | Planejado |
| CT-030 | Compilar CLIs Go. | Go instalado. | Nenhum. | Executar builds dos CLIs. | Binários são gerados. | A registrar. | Planejado |
| CT-031 | Compilar `assinador.jar`. | Java 21 e Maven instalados. | Nenhum. | Executar `mvn clean package`. | JAR é gerado em `target/assinador.jar`. | A registrar. | Planejado |
| CT-032 | Verificar workflow de build. | GitHub Actions configurado. | Push ou PR. | Executar workflow `build.yml`. | Workflow executa build e testes. | A registrar. | Planejado |
| CT-033 | Verificar workflow de release. | Tag SemVer criada. | Tag `v0.1.0`. | Executar workflow `release.yml`. | Release gera binários e checksums. | A registrar. | Planejado |
| CT-034 | Verificar artefatos de segurança da release. | Release gerada. | Artefatos publicados. | Conferir `checksums.txt`, `.sig` e `.pem`. | Artefatos de verificação estão presentes. | A registrar. | Planejado |

---

## 13. Execução dos testes

### 13.1 Testes executados

Nesta versão, o documento apresenta o planejamento e a matriz de execução dos testes. A execução deve ser registrada pela equipe após rodar os comandos no ambiente configurado.

| ID | Caso de teste | Data de execução | Executor | Resultado | Status |
|---|---|---|---|---|---|
| CT-001 | Verificar versão do CLI `assinatura` | A preencher | A preencher | A registrar | Planejado |
| CT-003 | Criar assinatura simulada em modo local | A preencher | A preencher | A registrar | Planejado |
| CT-008 | Iniciar `assinador.jar` em modo servidor | A preencher | A preencher | A registrar | Planejado |
| CT-009 | Criar assinatura via HTTP | A preencher | A preencher | A registrar | Planejado |
| CT-017 | Iniciar simulador com JAR local | A preencher | A preencher | A registrar | Planejado |
| CT-028 | Executar testes unitários Go | A preencher | A preencher | A registrar | Planejado |
| CT-029 | Executar testes unitários Java | A preencher | A preencher | A registrar | Planejado |

### 13.2 Testes aprovados

| ID | Caso de teste | Evidência | Observação |
|---|---|---|---|
| A preencher | A preencher | A preencher | Preencher após execução real dos testes. |

### 13.3 Testes reprovados

| ID | Caso de teste | Problema encontrado | Defeito relacionado |
|---|---|---|---|
| A preencher | A preencher | A preencher | A preencher |

### 13.4 Testes bloqueados

| ID | Caso de teste | Motivo do bloqueio | Ação necessária |
|---|---|---|---|
| CT-017 | Iniciar simulador com JAR local | Depende da disponibilidade do `simulador.jar` real. | Fornecer JAR local ou URL válida. |
| CT-029 | Executar testes unitários Java | Depende de Maven instalado. | Instalar Maven 3.9+ ou usar script alternativo. |
| CT-032 | Verificar workflow de build | Depende de repositório GitHub configurado. | Configurar repositório e habilitar Actions. |
| CT-033 | Verificar workflow de release | Depende de tag SemVer e permissões de release. | Criar tag e configurar permissões. |
| CT-034 | Verificar artefatos de segurança da release | Depende da execução do workflow de release. | Executar release e validar artefatos. |

---

## 14. Registro de defeitos

### 14.1 Identificação do defeito

Cada defeito deve receber um identificador único no formato `DF-XXX`, por exemplo: `DF-001`, `DF-002` e `DF-003`.

### 14.2 Descrição do problema

A descrição deve explicar claramente o comportamento incorreto observado, incluindo a funcionalidade afetada.

### 14.3 Passos para reproduzir

Os passos devem ser objetivos e permitir que outro membro da equipe reproduza o erro no mesmo ambiente.

### 14.4 Resultado esperado

Deve indicar o comportamento correto esperado de acordo com os requisitos e casos de teste.

### 14.5 Resultado obtido

Deve indicar o comportamento real observado durante a execução.

### 14.6 Severidade

| Severidade | Significado |
|---|---|
| Crítica | Impede o uso do sistema ou de fluxo principal sem alternativa. |
| Alta | Afeta funcionalidade importante, mas há contorno parcial. |
| Média | Afeta funcionalidade secundária ou depende de configuração. |
| Baixa | Problema textual, visual, documental ou melhoria de usabilidade. |

### 14.7 Status

| Status | Significado |
|---|---|
| Aberto | Defeito registrado e ainda não corrigido. |
| Em análise | Defeito está sendo investigado. |
| Corrigido | Correção foi implementada. |
| Reaberto | Defeito voltou a ocorrer após correção. |
| Fechado | Defeito foi corrigido e validado. |
| Rejeitado | Registro não representa defeito real. |

### 14.8 Tabela de defeitos

| ID | Descrição do problema | Passos para reproduzir | Resultado esperado | Resultado obtido | Severidade | Status |
|---|---|---|---|---|---|---|
| DF-001 | Exemplo: Maven não disponível no ambiente de teste. | Executar `cd projetos/assinador-java && mvn test`. | Testes Java executados. | Comando não executa por ausência do Maven. | Média | Exemplo / A confirmar |
| DF-002 | Exemplo: `simulador.jar` real não disponível. | Executar `simulador start --jar ~/.hubsaude/simulador/simulador.jar`. | Simulador inicia. | JAR não encontrado. | Média | Exemplo / A confirmar |

> Os defeitos acima são exemplos coerentes com as limitações conhecidas. Após a execução real, devem ser mantidos apenas os defeitos efetivamente observados.

---

## 15. Evidências de teste

### 15.1 Prints

| Evidência | Descrição | Local sugerido |
|---|---|---|
| Print 01 | Saída de `assinatura version`. | `evidencias/prints/ct-001-version.png` |
| Print 02 | Saída de `assinatura sign --local`. | `evidencias/prints/ct-003-sign-local.png` |
| Print 03 | Saída de `assinatura validate --local`. | `evidencias/prints/ct-006-validate-local.png` |
| Print 04 | Saída de `assinatura start`. | `evidencias/prints/ct-008-start-server.png` |
| Print 05 | Saída de `simulador status`. | `evidencias/prints/ct-018-simulador-status.png` |

### 15.2 Logs

| Log | Descrição | Local sugerido |
|---|---|---|
| `go-test.log` | Saída de `go test ./...`. | `evidencias/logs/go-test.log` |
| `mvn-test.log` | Saída de `mvn test`. | `evidencias/logs/mvn-test.log` |
| `assinador-server.log` | Saída do servidor `assinador.jar`. | `evidencias/logs/assinador-server.log` |
| `simulador.log` | Saída do Simulador do HubSaúde. | `evidencias/logs/simulador.log` |
| `release-workflow.log` | Saída do workflow de release. | GitHub Actions / exportação em texto |

### 15.3 Vídeos

| Vídeo | Conteúdo sugerido |
|---|---|
| `video-fluxo-assinatura-local.mp4` | Execução do fluxo `assinatura sign --local`. |
| `video-fluxo-assinatura-http.mp4` | Inicialização do servidor e execução de `sign` via HTTP. |
| `video-fluxo-simulador.mp4` | Execução de `simulador start`, `status` e `stop`. |

### 15.4 Arquivos gerados

| Arquivo | Descrição |
|---|---|
| `bin/assinatura` | Binário gerado do CLI `assinatura`. |
| `bin/simulador` | Binário gerado do CLI `simulador`. |
| `projetos/assinador-java/target/assinador.jar` | Artefato Java gerado. |
| `~/.hubsaude/processos/assinador-8080.json` | Registro do processo do assinador. |
| `~/.hubsaude/processos/simulador-8443.json` | Registro do processo do simulador. |
| `checksums.txt` | Arquivo de hashes SHA256. |
| `<artefato>.sig` | Assinatura Cosign do artefato. |
| `<artefato>.pem` | Certificado associado à assinatura Cosign. |

---

## 16. Critérios de aprovação

### 16.1 Condições para aprovação

O Sistema Runner será considerado aprovado na rodada de testes quando:

- os comandos principais do CLI `assinatura` forem executados com sucesso;
- o `assinador.jar` for invocado em modo local via `java -jar`;
- o modo servidor do `assinador.jar` responder aos endpoints principais;
- a criação e a validação simuladas retornarem resultados esperados;
- parâmetros ausentes ou inválidos forem rejeitados com mensagens claras;
- o comando `assinatura stop` encerrar o servidor adequadamente;
- o CLI `simulador` executar os fluxos previstos quando houver `simulador.jar` disponível;
- os registros de processo forem criados e atualizados corretamente;
- os testes unitários Go e Java forem executados sem falhas críticas;
- os artefatos gerados puderem ser verificados por SHA256;
- não houver defeitos críticos ou altos abertos nos fluxos principais.

### 16.2 Condições para reprovação

O sistema será reprovado se:

- o CLI `assinatura` não compilar ou não executar;
- o `assinador.jar` não compilar ou não puder ser executado;
- a criação de assinatura simulada não funcionar;
- a validação de assinatura simulada não funcionar;
- o sistema aceitar parâmetros obrigatórios ausentes sem erro;
- o modo HTTP não responder aos endpoints principais;
- o sistema não tratar erros básicos de JAR inexistente, Java ausente ou porta ocupada;
- houver defeitos críticos sem correção;
- os testes automatizados principais falharem sem justificativa.

### 16.3 Condições para nova rodada de testes

Uma nova rodada de testes deve ser realizada quando:

- defeitos críticos, altos ou médios forem corrigidos;
- houver alteração nos comandos CLI;
- houver alteração nos endpoints HTTP;
- houver alteração na validação de parâmetros;
- houver alteração no gerenciamento de processos;
- houver alteração no empacotamento ou CI/CD;
- for integrado um `simulador.jar` real;
- for expandido o provisionamento automático de JDK/JRE;
- for implementada integração PKCS#11 mais completa.

---

## 17. Métricas de teste

### 17.1 Total de testes planejados

| Métrica | Valor |
|---|---:|
| Total de casos de teste planejados | 34 |

### 17.2 Total de testes executados

| Métrica | Valor |
|---|---:|
| Total de testes executados | A preencher após execução |

### 17.3 Total de testes aprovados

| Métrica | Valor |
|---|---:|
| Total de testes aprovados | A preencher após execução |

### 17.4 Total de testes reprovados

| Métrica | Valor |
|---|---:|
| Total de testes reprovados | A preencher após execução |

### 17.5 Total de defeitos encontrados

| Métrica | Valor |
|---|---:|
| Total de defeitos encontrados | A preencher após execução |

### 17.6 Métricas complementares recomendadas

| Métrica | Fórmula |
|---|---|
| Taxa de aprovação | `(testes aprovados / testes executados) × 100` |
| Taxa de reprovação | `(testes reprovados / testes executados) × 100` |
| Taxa de bloqueio | `(testes bloqueados / testes planejados) × 100` |
| Cobertura de requisitos | `(requisitos com pelo menos um teste / total de requisitos testáveis) × 100` |
| Cobertura automatizada | `(testes automatizados / testes planejados) × 100` |

---

## 18. Relatório final de testes

### 18.1 Resumo da execução

| Item | Resultado |
|---|---|
| Período de execução | A preencher |
| Responsável pela execução | A preencher |
| Ambiente utilizado | A preencher |
| Total de testes planejados | 34 |
| Total de testes executados | A preencher |
| Total de testes aprovados | A preencher |
| Total de testes reprovados | A preencher |
| Total de testes bloqueados | A preencher |
| Total de defeitos encontrados | A preencher |
| Decisão final | A preencher |

### 18.2 Principais resultados

Resultados esperados ao final de uma execução bem-sucedida:

- CLIs Go compilados corretamente.
- `assinador.jar` compilado corretamente.
- Comando `assinatura version` funcionando.
- Assinatura simulada funcionando em modo local.
- Validação simulada funcionando em modo local.
- Servidor HTTP do `assinador.jar` iniciando corretamente.
- Endpoints `/health`, `/sign`, `/validate` e `/shutdown` respondendo corretamente.
- Registro de processos sendo criado em `~/.hubsaude/processos/`.
- Erros de parâmetros sendo tratados de forma clara.
- Testes unitários Go e Java passando.
- Checksums SHA256 funcionando.

### 18.3 Principais problemas encontrados

| Problema potencial | Impacto | Encaminhamento |
|---|---|---|
| Maven ausente | Bloqueia testes Java por Maven. | Instalar Maven ou usar script alternativo. |
| Java 21 ausente | Bloqueia execução do `assinador.jar`. | Instalar Java 21 ou completar provisionamento automático. |
| `simulador.jar` real ausente | Bloqueia testes reais do CLI `simulador`. | Fornecer JAR local ou URL válida. |
| Porta 8080 ou 8443 ocupada | Bloqueia start de servidor. | Liberar porta ou usar `--port`. |
| CI/CD não configurado no GitHub | Bloqueia validação de releases. | Configurar repositório e permissões de Actions. |
| Cosign não instalado/configurado | Bloqueia verificação de assinatura local. | Instalar Cosign ou validar em pipeline. |

---

## 19. Conclusão e recomendações

### 19.1 Avaliação final

Este documento define uma estratégia de teste coerente com o escopo do Sistema Runner, cobrindo CLIs em Go, aplicação Java `assinador.jar`, integração local via `java -jar`, integração HTTP, gerenciamento de processos, verificação de checksum, comandos do simulador, tratamento de erros e estrutura de CI/CD.

A abordagem contempla testes unitários, funcionais, integração, erro, instalação, compatibilidade, segurança de artefatos e aceitação. Como o sistema é uma implementação acadêmica e de simulação, o documento também registra explicitamente que não serão testadas assinatura digital criptográfica real nem validade jurídica de assinatura.

### 19.2 Recomendações para correção

Após a primeira execução real dos testes, recomenda-se que a equipe:

- preencha a coluna “Resultado obtido” dos casos de teste;
- registre defeitos reais encontrados;
- anexe prints, logs e arquivos gerados;
- corrija defeitos críticos e altos antes da entrega final;
- reexecute os testes de regressão após cada correção;
- valide o fluxo completo local e HTTP;
- valide a geração de artefatos por CI/CD;
- documente limitações não resolvidas.

### 19.3 Indicação de nova execução de testes, se necessário

Uma nova rodada de testes será necessária se:

- houver defeito crítico ou alto;
- os endpoints HTTP forem alterados;
- os comandos CLI forem alterados;
- o tratamento de parâmetros for modificado;
- o provisionamento automático de JDK/JRE for concluído;
- o `simulador.jar` real for integrado;
- o adaptador PKCS#11 evoluir para teste com SoftHSM2 ou token real;
- o pipeline de release com Cosign for validado em ambiente GitHub.

---

## 20. Referências

- Especificação original do Sistema Runner — Trabalho Prático.
- Plano Revisado #2 do Sistema Runner.
- Documento de Design do Sistema Runner baseado no Modelo C4.
- Especificação de Requisitos de Software — Sistema Runner.
- Documento de Arquitetura de Software — Sistema Runner.
- Documento de Projeto Detalhado de Software — Sistema Runner.
- Documento de Modelo C4 de Software — Sistema Runner.
- Documento de Implementação e Integração de Software — Sistema Runner.
- README da implementação de referência acadêmica do Sistema Runner.
- Boas práticas de Engenharia de Software para planejamento, execução, rastreabilidade e relatório de testes.
