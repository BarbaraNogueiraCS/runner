# 1. Capa

**Documento:** Especificação de Requisitos de Software  
**Sistema:** Sistema Runner  
**Disciplina:** Implementação e Integração de Software  
**Curso:** Engenharia de Software  
**Contexto institucional:** Trabalho prático relacionado à Plataforma HubSaúde, iniciativa de interesse da Secretaria de Estado da Saúde de Goiás (SES-GO) e da Universidade Federal de Goiás (UFG).  
**Tipo de documento:** Especificação de Requisitos de Software (ERS)  
**Versão:** 1.0  
**Data:** 07/06/2026  
**Status:** Versão inicial consolidada a partir dos artefatos de especificação, design, plano revisado, critérios de qualidade e tarefas operacionais do projeto.  
**Autores:** Equipe do Sistema Runner  

---

# 2. Histórico de versões

| Versão | Data | Autor(es) | Descrição da alteração |
|---|---:|---|---|
| 1.0 | 07/06/2026 | Equipe do Sistema Runner | Criação da versão inicial da especificação de requisitos de software, consolidando escopo, requisitos funcionais, requisitos não funcionais, histórias de usuário, critérios de aceitação, restrições técnicas e casos de teste iniciais. |

---

# Sumário

- [1. Capa](#1-capa)
- [2. Histórico de versões](#2-histórico-de-versões)
- [3. Introdução](#3-introdução)
  - [3.1 Objetivo do documento](#31-objetivo-do-documento)
  - [3.2 Escopo do sistema](#32-escopo-do-sistema)
  - [3.3 Definições, siglas e abreviações](#33-definições-siglas-e-abreviações)
- [4. Descrição geral do sistema](#4-descrição-geral-do-sistema)
  - [4.1 Contexto do sistema](#41-contexto-do-sistema)
  - [4.2 Objetivos do sistema](#42-objetivos-do-sistema)
  - [4.3 Principais funcionalidades](#43-principais-funcionalidades)
  - [4.4 Classes de usuários](#44-classes-de-usuários)
  - [4.5 Restrições gerais](#45-restrições-gerais)
  - [4.6 Premissas e dependências](#46-premissas-e-dependências)
- [5. Requisitos funcionais](#5-requisitos-funcionais)
- [6. Requisitos não funcionais](#6-requisitos-não-funcionais)
- [7. Regras de negócio](#7-regras-de-negócio)
- [8. Requisitos de interface externa](#8-requisitos-de-interface-externa)
- [9. Requisitos de dados](#9-requisitos-de-dados)
- [10. Histórias de usuário](#10-histórias-de-usuário)
- [11. Critérios de aceitação](#11-critérios-de-aceitação)
- [12. Priorização dos requisitos](#12-priorização-dos-requisitos)
- [13. Restrições técnicas](#13-restrições-técnicas)
- [14. Protótipos de tela](#14-protótipos-de-tela)
- [15. Casos de teste iniciais](#15-casos-de-teste-iniciais)
- [16. Glossário](#16-glossário)
- [17. Apêndices e anexos](#17-apêndices-e-anexos)

---

# 3. Introdução

## 3.1 Objetivo do documento

Este documento tem como objetivo especificar, de forma clara e rastreável, os requisitos de software do **Sistema Runner**. A especificação consolida as informações apresentadas nos documentos de trabalho, design, planejamento, critérios de qualidade e tarefas operacionais do projeto.

A finalidade deste documento é servir como referência para:

- implementação dos componentes do sistema;
- validação do escopo do trabalho prático;
- planejamento das sprints e entregas;
- elaboração de testes unitários, testes de integração e testes de aceitação;
- documentação técnica e manual de uso;
- avaliação da conformidade entre requisitos, código, testes e artefatos entregues.

Este documento deve ser mantido atualizado durante o desenvolvimento, preservando a rastreabilidade entre requisitos, histórias de usuário, critérios de aceitação, issues, pull requests, commits, código e testes.

## 3.2 Escopo do sistema

O **Sistema Runner** é um sistema de linha de comando cujo propósito é facilitar a execução de aplicações Java associadas à Plataforma HubSaúde, ocultando do usuário detalhes técnicos relacionados à instalação do Java, execução de arquivos `.jar`, gerenciamento de processos, comunicação HTTP, portas, validação de parâmetros e obtenção de artefatos.

O sistema contempla o desenvolvimento de:

- uma aplicação de linha de comando chamada `assinatura`, responsável por invocar o `assinador.jar` para criar e validar assinaturas digitais simuladas;
- uma aplicação Java chamada `assinador.jar`, responsável por validar parâmetros e simular operações de criação e validação de assinatura digital;
- uma aplicação de linha de comando chamada `simulador`, responsável por iniciar, parar e consultar o status do Simulador do HubSaúde;
- mecanismos de integração local e via HTTP entre CLI e aplicações Java;
- mecanismo de provisionamento automático do JDK/JRE necessário para execução dos componentes Java;
- geração de binários multiplataforma para Windows, Linux e macOS;
- publicação de releases com versionamento semântico, checksums SHA256 e assinatura dos artefatos com Cosign/Sigstore;
- testes e documentação de uso.

Não fazem parte do escopo:

- implementação real de assinatura digital criptográfica;
- implementação real de validação criptográfica;
- integração com autoridades certificadoras;
- geração de certificados digitais;
- autenticação de usuários;
- armazenamento persistente de assinaturas;
- interface gráfica de usuário.

## 3.3 Definições, siglas e abreviações

| Termo | Definição |
|---|---|
| ADR | *Architecture Decision Record*. Registro curto de decisão arquitetural relevante. |
| CLI | *Command Line Interface*. Interface de linha de comando usada por terminal. |
| Cosign | Ferramenta do ecossistema Sigstore usada para assinar e verificar artefatos de software. |
| FHIR | *Fast Healthcare Interoperability Resources*. Padrão de interoperabilidade em saúde. |
| GitHub Actions | Plataforma de automação de CI/CD integrada ao GitHub. |
| GitHub Releases | Área do GitHub usada para publicar versões e artefatos executáveis de um projeto. |
| Health check | Verificação que confirma se um serviço está em execução e respondendo corretamente. |
| HTTP | Protocolo de comunicação usado para troca de requisições e respostas entre cliente e servidor. |
| JAR | Arquivo Java empacotado, normalmente executado com `java -jar`. |
| JDK | *Java Development Kit*. Kit de desenvolvimento Java. |
| JRE | *Java Runtime Environment*. Ambiente de execução Java. |
| JVM | *Java Virtual Machine*. Máquina virtual usada para executar aplicações Java. |
| OIDC | *OpenID Connect*. Protocolo de identidade usado no fluxo de assinatura de artefatos. |
| PKCS#11 | Padrão de interface para comunicação com dispositivos criptográficos, como tokens e smart cards. |
| SemVer | *Semantic Versioning*. Esquema de versionamento semântico no formato `MAJOR.MINOR.PATCH`. |
| SHA256 | Algoritmo de hash usado para verificação de integridade de arquivos. |
| Sigstore | Ecossistema para assinatura, verificação e transparência de artefatos de software. |
| SoftHSM2 | Software que simula um dispositivo criptográfico HSM, útil para testes com PKCS#11. |
| stdout | Saída padrão de um processo, usada para resultados esperados. |
| stderr | Saída de erro de um processo, usada para diagnósticos e mensagens de erro. |

---

# 4. Descrição geral do sistema

## 4.1 Contexto do sistema

O Sistema Runner está inserido no contexto da Plataforma HubSaúde, uma plataforma de interoperabilidade de dados em saúde. Seu papel é facilitar a execução de aplicações Java relacionadas a esse ambiente, especialmente por integradores que precisam utilizar ferramentas de assinatura e simulação sem lidar diretamente com comandos Java e configurações de ambiente.

No nível de contexto, o sistema interage com três elementos principais:

- **Usuário:** pessoa que utiliza o Runner por meio de comandos no terminal.
- **Dispositivo de Assinatura Digital:** hardware criptográfico, como token USB ou smart card, compatível com PKCS#11.
- **Simulador do HubSaúde:** aplicação Java/Web que responde a requisições e pode ser controlada por meio do CLI `simulador`.

No nível de contêineres, o sistema é composto por:

- **CLI `assinatura`:** interface de linha de comando para criação e validação de assinaturas simuladas.
- **`assinador.jar`:** aplicação Java responsável por validar parâmetros, simular assinatura, simular validação e, quando aplicável, interagir com dispositivo criptográfico via PKCS#11.
- **CLI `simulador`:** interface de linha de comando para gerenciar o ciclo de vida do Simulador do HubSaúde.
- **Simulador do HubSaúde:** aplicação Java externa gerida pelo CLI.
- **Diretório local gerenciado:** diretório, preferencialmente `~/.hubsaude/`, usado para armazenar JDK/JRE provisionado, JARs baixados, metadados de versão, PID, porta e informações de execução.

## 4.2 Objetivos do sistema

O objetivo geral do Sistema Runner é **facilitar a execução de aplicações Java por linha de comando**, permitindo que usuários executem operações relacionadas ao HubSaúde sem conhecimento aprofundado sobre Java, configuração de ambiente, execução de `.jar`, gerenciamento de processos e chamadas HTTP.

Os objetivos específicos são:

1. Permitir a execução de comandos de assinatura digital simulada por meio de um CLI simples e intuitivo.
2. Permitir a validação simulada de assinaturas digitais por meio do mesmo CLI.
3. Invocar o `assinador.jar` em modo local ou via HTTP.
4. Usar o modo servidor como padrão quando não houver orientação explícita para modo local.
5. Gerenciar o ciclo de vida do `assinador.jar` quando executado como servidor.
6. Gerenciar o ciclo de vida do Simulador do HubSaúde por CLI.
7. Baixar e configurar automaticamente o JDK/JRE necessário quando ausente.
8. Disponibilizar binários multiplataforma para Windows, Linux e macOS.
9. Garantir integridade e autenticidade dos artefatos publicados em releases.
10. Fornecer documentação, testes e rastreabilidade suficientes para comprovar o funcionamento do sistema.

## 4.3 Principais funcionalidades

As principais funcionalidades previstas são:

- criação simulada de assinatura digital;
- validação simulada de assinatura digital;
- validação rigorosa de parâmetros no `assinador.jar`;
- invocação local do `assinador.jar` via subprocesso;
- invocação do `assinador.jar` via HTTP;
- inicialização do `assinador.jar` como servidor;
- detecção de instância ativa do `assinador.jar` por health check;
- interrupção manual do `assinador.jar`;
- interrupção automática do `assinador.jar` por inatividade;
- integração com dispositivo criptográfico ou simulador PKCS#11;
- início, parada e consulta de status do Simulador do HubSaúde;
- download automático do `simulador.jar`;
- download automático do JDK/JRE quando necessário;
- geração de binários multiplataforma;
- publicação de releases com checksums SHA256 e assinatura Cosign;
- comandos de ajuda e versão;
- tratamento estruturado de erros.

## 4.4 Classes de usuários

| Classe de usuário | Descrição | Necessidades principais |
|---|---|---|
| Usuário final/integrador | Pessoa que utiliza os CLIs para assinar, validar ou gerenciar o simulador. | Comandos simples, mensagens claras, instalação mínima, resultados legíveis. |
| Desenvolvedor do projeto | Pessoa que implementa funcionalidades, corrige defeitos e mantém o código. | Requisitos rastreáveis, arquitetura clara, testes automatizados, CI/CD. |
| Avaliador ou professor | Pessoa responsável por avaliar o trabalho prático. | Evidências de funcionamento, rastreabilidade, documentação e conformidade com critérios. |
| Mantenedor de release | Pessoa responsável por publicar versões e artefatos. | Pipeline automatizado, SemVer, checksums, assinatura de artefatos e changelog. |

## 4.5 Restrições gerais

- O sistema deve ser operado por linha de comando, sem interface gráfica.
- Os CLIs devem ser compatíveis com Windows, Linux e macOS em arquitetura `amd64`.
- O CLI deve ocultar do usuário os detalhes de execução de aplicações Java.
- O `assinador.jar` deve ser implementado em Java 21.
- Os CLIs devem ser implementados em Go 1.25.
- O sistema deve validar parâmetros de forma rigorosa no `assinador.jar`.
- O sistema não deve implementar assinatura digital criptográfica real.
- O sistema não deve armazenar assinaturas de forma persistente.
- Os binários devem ser publicados em GitHub Releases.
- Os artefatos publicados devem ter checksums SHA256 e assinatura Cosign/Sigstore.
- O projeto deve adotar CI/CD para build, testes e publicação.
- O modo servidor deve ser o modo preferencial/padrão para invocação do `assinador.jar`, exceto quando o usuário solicitar explicitamente o modo local.

## 4.6 Premissas e dependências

| Identificador | Premissa ou dependência |
|---|---|
| PD-01 | O usuário terá acesso a um terminal compatível com o sistema operacional utilizado. |
| PD-02 | O ambiente terá acesso à internet quando for necessário baixar JDK/JRE, `simulador.jar` ou artefatos de release. |
| PD-03 | O projeto utilizará GitHub como repositório, CI/CD e publicação de releases. |
| PD-04 | Os CLIs serão desenvolvidos em Go 1.25, com suporte à cross-compilation. |
| PD-05 | O `assinador.jar` será desenvolvido em Java 21. |
| PD-06 | O JDK/JRE será obtido preferencialmente a partir do Eclipse Temurin/Adoptium quando não estiver disponível localmente. |
| PD-07 | O Simulador do HubSaúde poderá ser obtido dinamicamente a partir de releases do repositório definido. |
| PD-08 | O dispositivo criptográfico real poderá ser substituído por simulador compatível, como SoftHSM2, durante testes de integração. |
| PD-09 | As especificações externas de criação e validação de assinatura serão usadas como referência para parâmetros e formato de resposta. |
| PD-10 | O diretório local gerenciado `~/.hubsaude/` será usado para cache, metadados e arquivos provisionados. |

---

# 5. Requisitos funcionais

## 5.1 Requisitos do CLI `assinatura`

| ID | Requisito funcional | Prioridade | Origem |
|---|---|---|---|
| RF-01 | O sistema deve disponibilizar um CLI chamado `assinatura` para operações de assinatura e validação. | Alta | US-01 |
| RF-02 | O CLI `assinatura` deve aceitar comando para criação simulada de assinatura digital. | Alta | US-01, US-02 |
| RF-03 | O CLI `assinatura` deve aceitar comando para validação simulada de assinatura digital. | Alta | US-01, US-02 |
| RF-04 | O CLI `assinatura` deve interpretar comandos, parâmetros e flags informados pelo usuário. | Alta | US-01.2 |
| RF-05 | O CLI `assinatura` deve exibir ajuda de uso por meio de `--help` ou comando equivalente. | Alta | Critérios de operabilidade |
| RF-06 | O CLI `assinatura` deve exibir a versão atual por meio de comando `version` ou equivalente. | Média | US-01.1 |
| RF-07 | O CLI `assinatura` deve exibir resultados de assinatura e validação de forma legível e estruturada. | Alta | US-01.4 |
| RF-08 | O CLI `assinatura` deve diferenciar mensagens de resultado em `stdout` e mensagens de diagnóstico em `stderr`. | Alta | Critérios de qualidade |
| RF-09 | O CLI `assinatura` deve propagar códigos de saída coerentes conforme sucesso, erro de usuário ou erro de sistema. | Alta | Critérios de qualidade |

## 5.2 Requisitos do `assinador.jar`

| ID | Requisito funcional | Prioridade | Origem |
|---|---|---|---|
| RF-10 | O sistema deve disponibilizar uma aplicação Java chamada `assinador.jar`. | Alta | US-01, US-02 |
| RF-11 | O `assinador.jar` deve validar rigorosamente os parâmetros de criação de assinatura. | Alta | US-02.2 |
| RF-12 | O `assinador.jar` deve validar rigorosamente os parâmetros de validação de assinatura. | Alta | US-02.3 |
| RF-13 | O `assinador.jar` deve simular a criação de assinatura digital quando os parâmetros forem válidos. | Alta | US-02.1 |
| RF-14 | O `assinador.jar` deve retornar resposta pré-construída contendo os campos esperados para uma assinatura simulada. | Alta | US-02.1 |
| RF-15 | O `assinador.jar` deve simular a validação de assinatura digital com resultado pré-determinado ou baseado em critérios simples. | Alta | US-02.3 |
| RF-16 | O `assinador.jar` deve retornar mensagens de erro claras quando parâmetros forem inválidos. | Alta | US-02 |
| RF-17 | O `assinador.jar` deve distinguir erros do usuário de erros do sistema. | Alta | Critérios de qualidade |
| RF-18 | O `assinador.jar` deve suportar integração com dispositivo criptográfico via PKCS#11. | Média | US-02.5 |
| RF-19 | O `assinador.jar` deve tratar adequadamente a ausência de dispositivo criptográfico ou falha de acesso ao PKCS#11. | Média | US-02.5 |

## 5.3 Requisitos de integração local entre CLI e JAR

| ID | Requisito funcional | Prioridade | Origem |
|---|---|---|---|
| RF-20 | O CLI `assinatura` deve invocar o `assinador.jar` em modo local por meio de subprocesso. | Alta | US-01.3 |
| RF-21 | A invocação local deve funcionar independentemente do diretório atual em que o CLI for executado. | Alta | Critérios E1 |
| RF-22 | A invocação local deve preservar espaços, acentos e aspas em argumentos. | Alta | Critérios E1 |
| RF-23 | O CLI deve capturar `stdout`, `stderr` e código de saída do `assinador.jar`. | Alta | Critérios E1 |
| RF-24 | O CLI deve apresentar mensagem clara quando o `assinador.jar` não for encontrado. | Alta | Critérios de erro |
| RF-25 | O CLI deve apresentar mensagem clara quando Java/JDK/JRE compatível não estiver disponível. | Alta | US-04.1 |

## 5.4 Requisitos de integração HTTP e modo servidor

| ID | Requisito funcional | Prioridade | Origem |
|---|---|---|---|
| RF-26 | O `assinador.jar` deve poder ser iniciado em modo servidor HTTP. | Alta | US-02.4, US-01.5 |
| RF-27 | O `assinador.jar` em modo servidor deve expor endpoint HTTP para criação de assinatura. | Alta | US-02.4 |
| RF-28 | O `assinador.jar` em modo servidor deve expor endpoint HTTP para validação de assinatura. | Alta | US-02.4 |
| RF-29 | O `assinador.jar` em modo servidor deve expor endpoint de health check. | Alta | US-01.7 |
| RF-30 | O CLI `assinatura` deve iniciar o `assinador.jar` como processo em segundo plano no modo servidor. | Alta | US-01.5 |
| RF-31 | O CLI `assinatura` deve registrar PID, porta e metadados do processo em `~/.hubsaude/` ou diretório equivalente. | Alta | US-01.5 |
| RF-32 | O CLI `assinatura` deve verificar, por health check real, se já existe instância ativa do `assinador.jar`. | Alta | US-01.7 |
| RF-33 | O CLI `assinatura` deve reutilizar instância ativa do `assinador.jar` quando disponível. | Alta | US-01.7 |
| RF-34 | O CLI `assinatura` deve iniciar nova instância quando não houver instância ativa. | Alta | US-01.5 |
| RF-35 | O modo servidor deve ser usado como padrão, exceto quando o usuário solicitar explicitamente modo local. | Alta | US-01.6 |
| RF-36 | O CLI `assinatura` deve permitir especificar porta customizada para o servidor. | Média | US-01.5, US-01.8 |
| RF-37 | O CLI `assinatura` deve encerrar o `assinador.jar` na porta padrão ou em porta indicada. | Alta | US-01.8 |
| RF-38 | O CLI `assinatura` deve permitir configurar encerramento automático por inatividade. | Média | US-01.9 |
| RF-39 | O mecanismo de inatividade deve reiniciar a contagem a cada requisição válida recebida. | Média | Critérios E2 |
| RF-40 | O CLI deve tratar explicitamente timeout, conexão recusada e resposta HTTP malformada. | Alta | Critérios E2 |

## 5.5 Requisitos do CLI `simulador`

| ID | Requisito funcional | Prioridade | Origem |
|---|---|---|---|
| RF-41 | O sistema deve disponibilizar um CLI chamado `simulador`. | Alta | US-03.3 |
| RF-42 | O CLI `simulador` deve disponibilizar comando `start`. | Alta | US-03.1 |
| RF-43 | O CLI `simulador` deve disponibilizar comando `stop`. | Alta | US-03.2 |
| RF-44 | O CLI `simulador` deve disponibilizar comando `status`. | Alta | US-03.2 |
| RF-45 | O CLI `simulador` deve verificar se a porta padrão do simulador, `8443`, está disponível antes de iniciar o processo. | Alta | US-03 |
| RF-46 | O CLI `simulador` deve iniciar o `simulador.jar` sem exigir que o usuário execute comandos Java diretamente. | Alta | US-03.1 |
| RF-47 | O CLI `simulador` deve encerrar o Simulador do HubSaúde por meio do endpoint `/shutdown` quando disponível. | Alta | US-03 |
| RF-48 | O CLI `simulador` deve consultar o status do simulador por meio do endpoint `/api/info` quando disponível. | Alta | US-03 |
| RF-49 | O CLI `simulador` deve registrar PID, porta e metadados em `~/.hubsaude/` ou diretório equivalente. | Média | US-03.2 |
| RF-50 | O CLI `simulador` deve diferenciar processo iniciado de serviço realmente pronto para receber requisições. | Alta | Critérios E4 |

## 5.6 Requisitos de download e provisionamento

| ID | Requisito funcional | Prioridade | Origem |
|---|---|---|---|
| RF-51 | O sistema deve detectar se há JDK/JRE compatível disponível no `PATH` ou em diretório gerenciado. | Alta | US-04.1 |
| RF-52 | O sistema deve baixar JDK/JRE compatível quando ausente. | Alta | US-04.1 |
| RF-53 | O sistema deve armazenar JDK/JRE baixado em diretório local gerenciado para reuso. | Alta | US-04.1 |
| RF-54 | O sistema não deve repetir download de JDK/JRE quando versão compatível já estiver disponível localmente. | Média | US-04.1 |
| RF-55 | O CLI `simulador` deve obter dinamicamente o `simulador.jar` quando ele não estiver disponível localmente. | Alta | US-03.4 |
| RF-56 | O CLI `simulador` deve comparar versões disponíveis com versões locais antes de baixar novo artefato. | Média | US-03.4 |
| RF-57 | O CLI `simulador` deve permitir URL alternativa de download por parâmetro `--source <url>`. | Baixa | US-03.4 |
| RF-58 | O sistema deve verificar a integridade dos artefatos baixados por checksum. | Alta | US-03.4, US-05.3 |

## 5.7 Requisitos de build, release e distribuição

| ID | Requisito funcional | Prioridade | Origem |
|---|---|---|---|
| RF-59 | O sistema deve gerar binários do CLI `assinatura` para Windows, Linux e macOS em arquitetura `amd64`. | Alta | US-05 |
| RF-60 | O sistema deve gerar binários do CLI `simulador` para Windows, Linux e macOS em arquitetura `amd64`. | Alta | US-03.3, US-05 |
| RF-61 | Os binários devem ser publicados em GitHub Releases. | Alta | US-05 |
| RF-62 | Os artefatos devem seguir convenção de nome definida pelo projeto. | Média | Tarefas operacionais |
| RF-63 | As versões publicadas devem seguir SemVer. | Alta | US-05 |
| RF-64 | Cada release deve conter checksums SHA256 dos artefatos. | Alta | US-05.3 |
| RF-65 | Cada artefato publicado deve ser assinado com Cosign. | Alta | Integridade e assinatura |
| RF-66 | Cada artefato assinado deve ser acompanhado por arquivos `.sig` e `.pem`. | Alta | Integridade e assinatura |
| RF-67 | A assinatura dos artefatos deve ser automatizada no pipeline de CI/CD. | Alta | Integridade e assinatura |

---

# 6. Requisitos não funcionais

| ID | Categoria | Requisito não funcional | Prioridade |
|---|---|---|---|
| RNF-01 | Portabilidade | O sistema deve funcionar em Windows, Linux e macOS em arquitetura `amd64`. | Alta |
| RNF-02 | Usabilidade | Os comandos CLI devem ser simples, consistentes e acompanhados de ajuda com exemplos. | Alta |
| RNF-03 | Usabilidade | As mensagens de erro devem explicar o que ocorreu, por que ocorreu e como o usuário pode corrigir. | Alta |
| RNF-04 | Confiabilidade | O sistema deve falhar de forma controlada, com códigos de saída coerentes. | Alta |
| RNF-05 | Confiabilidade | O sistema deve tratar cenários de erro como porta ocupada, JAR ausente, JVM ausente, timeout e resposta malformada. | Alta |
| RNF-06 | Desempenho | O modo servidor deve reduzir o custo de inicialização da JVM em chamadas subsequentes. | Média |
| RNF-07 | Segurança | Os artefatos publicados devem permitir verificação de integridade e autenticidade. | Alta |
| RNF-08 | Segurança | O projeto não deve armazenar segredos, tokens ou credenciais no código-fonte. | Alta |
| RNF-09 | Segurança | Caminhos absolutos, IPs e portas não devem ser fixados no código quando deveriam ser configuráveis. | Alta |
| RNF-10 | Manutenibilidade | O código deve ter responsabilidades separadas entre CLI, transporte, domínio, download, provisionamento e integração. | Alta |
| RNF-11 | Manutenibilidade | Decisões arquiteturais não óbvias devem ser registradas em ADRs curtos. | Média |
| RNF-12 | Testabilidade | O sistema deve possuir testes unitários, testes de integração e testes de aceitação. | Alta |
| RNF-13 | Testabilidade | Cenários negativos devem ser tratados como casos de teste de primeira classe. | Alta |
| RNF-14 | Reprodutibilidade | Qualquer pessoa deve conseguir clonar o projeto, compilar e executar testes seguindo a documentação. | Alta |
| RNF-15 | CI/CD | O projeto deve usar GitHub Actions para lint, build, testes e release. | Alta |
| RNF-16 | Rastreabilidade | Requisitos, histórias, issues, PRs, commits, código e testes devem ser rastreáveis. | Alta |
| RNF-17 | Documentação | O README deve funcionar como contrato do projeto, descrevendo uso, build, testes e status. | Alta |
| RNF-18 | Qualidade de código | O código deve seguir estilo idiomático da linguagem e ser validado por ferramentas automatizadas. | Alta |
| RNF-19 | Logs | O sistema deve usar logs adequados e configuráveis, evitando uso indiscriminado de `print` ou `System.out` para diagnósticos. | Média |
| RNF-20 | Operabilidade | Os CLIs devem oferecer opções previsíveis de `--verbose` e/ou `--quiet` quando aplicável. | Baixa |
| RNF-21 | Compatibilidade | O `assinador.jar` deve ser compatível com Java 21. | Alta |
| RNF-22 | Distribuição | Os artefatos de release devem ter nomes padronizados e versão rastreável. | Alta |
| RNF-23 | Organização | O repositório não deve versionar artefatos gerados, caches ou arquivos temporários. | Alta |
| RNF-24 | Encoding | O projeto deve adotar UTF-8 e tratar finais de linha por meio de `.gitattributes`. | Média |

---

# 7. Regras de negócio

| ID | Regra de negócio |
|---|---|
| RB-01 | O Sistema Runner deve facilitar a execução de aplicações Java relacionadas ao HubSaúde, ocultando detalhes técnicos do usuário. |
| RB-02 | A criação de assinatura digital no projeto é simulada, não criptograficamente real. |
| RB-03 | A validação de assinatura digital no projeto é simulada, não criptograficamente real. |
| RB-04 | O `assinador.jar` é a autoridade principal de validação dos parâmetros de assinatura e validação. |
| RB-05 | O CLI não deve duplicar regras complexas de validação que pertencem ao `assinador.jar`. |
| RB-06 | O modo servidor deve ser o modo padrão de uso do `assinador.jar`, salvo quando o usuário solicitar explicitamente modo local. |
| RB-07 | O sistema deve reutilizar uma instância viva do `assinador.jar` em modo servidor em vez de iniciar instâncias duplicadas. |
| RB-08 | Uma porta ocupada não deve ser considerada automaticamente uma instância válida; é obrigatório verificar readiness/health check. |
| RB-09 | O Simulador do HubSaúde deve ser iniciado, parado e monitorado pelo CLI `simulador`. |
| RB-10 | O Simulador do HubSaúde deve usar a porta padrão `8443`, salvo configuração diferente prevista pela implementação. |
| RB-11 | O sistema não deve baixar novamente JDK/JRE ou JARs quando uma versão compatível e íntegra já estiver disponível localmente. |
| RB-12 | Todos os artefatos publicados em release devem possuir checksum SHA256. |
| RB-13 | Todos os artefatos distribuídos em release devem ser assinados com Cosign/Sigstore. |
| RB-14 | O versionamento dos artefatos deve seguir SemVer. |
| RB-15 | O projeto deve priorizar rastreabilidade, reprodutibilidade e qualidade de integração. |

---

# 8. Requisitos de interface externa

## 8.1 Interface de usuário por linha de comando

O sistema deve fornecer interfaces de linha de comando para os binários `assinatura` e `simulador`.

### 8.1.1 CLI `assinatura`

Comandos esperados, em forma conceitual:

```bash
assinatura version
assinatura --help
assinatura sign [parâmetros de assinatura]
assinatura validate [parâmetros de validação]
assinatura stop [--port <porta>]
```

O CLI `assinatura` deve:

- aceitar comandos de criação e validação;
- permitir modo servidor por padrão;
- permitir modo local por flag explícita, por exemplo `--local`;
- permitir porta customizada quando aplicável;
- permitir configuração de timeout de inatividade quando aplicável;
- exibir resultados em formato legível;
- separar resultado (`stdout`) de diagnóstico (`stderr`);
- retornar códigos de saída coerentes.

### 8.1.2 CLI `simulador`

Comandos esperados, em forma conceitual:

```bash
simulador --help
simulador start [--port 8443]
simulador stop [--port 8443]
simulador status [--port 8443]
```

O CLI `simulador` deve:

- iniciar o `simulador.jar`;
- parar o `simulador.jar`;
- consultar status do simulador;
- verificar porta antes de iniciar;
- verificar readiness antes de informar sucesso;
- baixar o `simulador.jar` quando ausente;
- reutilizar arquivos já baixados quando compatíveis.

## 8.2 Interface entre CLI `assinatura` e `assinador.jar` em modo local

A integração local deve ocorrer por subprocesso, com execução conceitual semelhante a:

```bash
java -jar assinador.jar [comando] [parâmetros]
```

Essa interface deve preservar corretamente:

- argumentos com espaços;
- argumentos com acentos;
- argumentos com aspas;
- códigos de saída;
- saída padrão;
- saída de erro.

## 8.3 Interface HTTP do `assinador.jar`

Quando executado em modo servidor, o `assinador.jar` deve expor endpoints HTTP. A nomenclatura exata poderá ser refinada durante o design detalhado, mas os endpoints mínimos esperados são:

| Método | Endpoint | Finalidade |
|---|---|---|
| `GET` | `/health` ou equivalente | Verificar se o servidor está vivo. |
| `POST` | `/sign` | Solicitar criação simulada de assinatura. |
| `POST` | `/validate` | Solicitar validação simulada de assinatura. |
| `POST` ou `GET` | `/shutdown` ou equivalente | Solicitar encerramento controlado, se definido para o componente. |

As respostas devem ser estruturadas e consistentes, distinguindo sucesso e erro.

## 8.4 Interface HTTP do Simulador do HubSaúde

O CLI `simulador` deve interagir com o Simulador do HubSaúde por HTTP. Os endpoints previstos são:

| Método | Endpoint | Finalidade |
|---|---|---|
| `GET` | `/api/info` | Consultar informações/status do simulador. |
| `POST` ou `GET` | `/shutdown` | Solicitar encerramento do simulador. |

A implementação deve confirmar que o serviço está pronto para receber requisições antes de informar ao usuário que o simulador foi iniciado com sucesso.

## 8.5 Interface com dispositivo criptográfico

O `assinador.jar` deve suportar integração com dispositivo criptográfico por meio da interface PKCS#11. A implementação deve permitir:

- uso de token ou smart card compatível;
- uso de simulador equivalente em testes, como SoftHSM2;
- tratamento de erro quando o dispositivo não estiver disponível;
- documentação do setup necessário.

## 8.6 Interface com sistema de arquivos

O sistema deve utilizar um diretório local gerenciado, preferencialmente `~/.hubsaude/`, para armazenar:

- JDK/JRE baixado;
- `simulador.jar` baixado;
- `assinador.jar`, se aplicável ao empacotamento definido;
- metadados de versão;
- registros de PID e porta;
- informações de cache;
- logs, quando aplicável.

## 8.7 Interface com serviços externos de download

O sistema poderá acessar:

- GitHub Releases para baixar artefatos do projeto;
- arquivo `release.json` para identificar versões e URLs;
- Eclipse Temurin/Adoptium para obter JDK/JRE compatível;
- endpoints de verificação de integridade, quando definidos.

---

# 9. Requisitos de dados

## 9.1 Dados de entrada

| Tipo de dado | Origem | Uso |
|---|---|---|
| Parâmetros de assinatura | Usuário via CLI ou requisição HTTP | Criar assinatura simulada. |
| Parâmetros de validação | Usuário via CLI ou requisição HTTP | Validar assinatura simulada. |
| Porta | Usuário ou configuração padrão | Iniciar, localizar ou parar servidores. |
| Timeout de inatividade | Usuário ou configuração padrão | Encerrar servidor após período sem uso. |
| URL de download | Configuração, release ou parâmetro `--source` | Baixar JARs ou dependências. |
| Versão esperada | Release, tag ou arquivo de metadados | Comparar e decidir download. |

## 9.2 Dados de saída

| Tipo de dado | Destino | Observação |
|---|---|---|
| Assinatura simulada | Usuário/CLI ou resposta HTTP | Gerada apenas como simulação. |
| Resultado de validação | Usuário/CLI ou resposta HTTP | Indica validade simulada. |
| Mensagens de erro | `stderr` ou resposta HTTP de erro | Devem ser claras e orientativas. |
| Status do simulador | Usuário via CLI | Deve indicar se está em execução/pronto. |
| Status do servidor do assinador | Usuário via CLI | Deve indicar se há instância ativa. |
| Checksums | Release | Usados para verificação de integridade. |
| Assinaturas Cosign | Release | Usadas para verificação de autenticidade. |

## 9.3 Dados persistidos localmente

O sistema não deve persistir assinaturas digitais como regra de negócio. Entretanto, poderá persistir dados operacionais necessários ao funcionamento:

| Dado | Local sugerido | Finalidade |
|---|---|---|
| JDK/JRE provisionado | `~/.hubsaude/jdk/` ou equivalente | Reutilizar Java sem novo download. |
| JARs baixados | `~/.hubsaude/` ou equivalente | Reutilizar artefatos locais. |
| Versão dos artefatos | Arquivo de metadados local | Evitar downloads desnecessários. |
| PID do processo | Arquivo de controle local | Permitir stop/status. |
| Porta do processo | Arquivo de controle local | Permitir localização e health check. |
| Logs | Diretório de logs, se definido | Diagnóstico e operabilidade. |

## 9.4 Dados não persistidos

| Dado | Justificativa |
|---|---|
| Assinaturas geradas | O projeto não prevê armazenamento persistente de assinaturas. |
| Credenciais ou segredos | O projeto não deve armazenar segredos no código ou em arquivos locais sem necessidade. |
| Dados sensíveis de usuários | O sistema não prevê autenticação nem cadastro de usuários. |

---

# 10. Histórias de usuário

## 10.1 Épicos principais

| ID | História de usuário | Prioridade |
|---|---|---|
| US-01 | Como usuário do Sistema Runner, quero executar comandos de assinatura digital pela linha de comando, para invocar o `assinador.jar` sem conhecer detalhes técnicos de configuração Java. | Alta |
| US-02 | Como usuário do Sistema Runner, quero que o `assinador.jar` valide rigorosamente os parâmetros de entrada antes de simular assinatura digital, para receber feedback imediato sobre erros. | Alta |
| US-03 | Como usuário do Sistema Runner, quero iniciar, parar e monitorar o Simulador do HubSaúde por CLI, para gerenciar seu ciclo de vida sem conhecer comandos Java. | Alta |
| US-04 | Como usuário do Sistema Runner, quero que o sistema baixe e configure automaticamente o JDK necessário quando ausente, para usar o Assinador e o Simulador sem configurar Java manualmente. | Alta |
| US-05 | Como usuário do Sistema Runner, quero baixar binários pré-compilados para minha plataforma, para utilizar o sistema sem necessidade de compilação. | Alta |

## 10.2 Histórias derivadas

| ID | História derivada | Épico de origem | Prioridade |
|---|---|---|---|
| US-01.1 | Como usuário, quero que o CLI tenha estrutura base e comando de versão, para identificar a versão em uso. | US-01 | Média |
| US-01.2 | Como usuário, quero executar comandos `sign` e `validate` com parâmetros, para solicitar operações de assinatura de forma intuitiva. | US-01 | Alta |
| US-01.3 | Como usuário, quero que o CLI invoque o `assinador.jar` diretamente via `java -jar`, para usar o modo local. | US-01 | Alta |
| US-01.4 | Como usuário, quero visualizar resultados de forma legível, para compreender o resultado da operação. | US-01 | Alta |
| US-01.5 | Como usuário, quero iniciar o `assinador.jar` em modo servidor, para reduzir latência em chamadas subsequentes. | US-01 | Alta |
| US-01.6 | Como usuário, quero que o CLI envie requisições HTTP ao `assinador.jar` por padrão, para evitar cold start em múltiplas chamadas. | US-01 | Alta |
| US-01.7 | Como usuário, quero que o CLI detecte instância ativa do `assinador.jar`, para evitar processos duplicados. | US-01 | Alta |
| US-01.8 | Como usuário, quero interromper o `assinador.jar` em uma porta específica ou padrão, para controlar processos em execução. | US-01 | Alta |
| US-01.9 | Como usuário, quero agendar interrupção por inatividade, para liberar recursos automaticamente. | US-01 | Média |
| US-02.1 | Como usuário, quero que o `assinador.jar` retorne assinatura simulada para parâmetros válidos, para testar o fluxo sem criptografia real. | US-02 | Alta |
| US-02.2 | Como usuário, quero que o `assinador.jar` valide parâmetros de criação de assinatura, para receber mensagens claras de erro. | US-02 | Alta |
| US-02.3 | Como usuário, quero que o `assinador.jar` valide parâmetros de validação e retorne resultado simulado, para testar o fluxo de validação. | US-02 | Alta |
| US-02.4 | Como usuário, quero que o `assinador.jar` exponha endpoints HTTP, para permitir invocação via servidor. | US-02 | Alta |
| US-02.5 | Como usuário, quero suporte a PKCS#11, para permitir uso de material criptográfico real ou simulado. | US-02 | Média |
| US-03.1 | Como usuário, quero iniciar o Simulador via CLI, para gerenciá-lo sem comandos Java manuais. | US-03 | Alta |
| US-03.2 | Como usuário, quero parar e consultar o status do Simulador, para ter controle sobre seu ciclo de vida. | US-03 | Alta |
| US-03.3 | Como usuário, quero um CLI dedicado para o Simulador, para ter uma interface independente e clara. | US-03 | Alta |
| US-03.4 | Como usuário, quero que o CLI baixe automaticamente o `simulador.jar`, para usar a versão disponível sem download manual. | US-03 | Média |
| US-04.1 | Como usuário, quero que o sistema detecte e provisione JDK/JRE compatível, para não configurar Java manualmente. | US-04 | Alta |
| US-05.1 | Como desenvolvedor, quero pipeline de CI/CD multiplataforma, para validar builds e testes nas plataformas-alvo. | US-05 | Alta |
| US-05.2 | Como usuário, quero releases com versionamento semântico, para baixar versões claras e rastreáveis. | US-05 | Alta |
| US-05.3 | Como usuário, quero checksums e assinatura Cosign, para verificar integridade e autenticidade dos artefatos. | US-05 | Alta |

---

# 11. Critérios de aceitação

## 11.1 Critérios de aceitação por épico

| Épico | Critérios de aceitação principais |
|---|---|
| US-01 | O CLI aceita comandos de criação e validação; invoca o `assinador.jar`; suporta modo local e HTTP; exibe resultado legível; inicia servidor em porta padrão; detecta e reutiliza instância ativa; interrompe servidor; permite interrupção por inatividade. |
| US-02 | O `assinador.jar` valida parâmetros; simula criação de assinatura; simula validação; suporta PKCS#11; retorna erros claros. |
| US-03 | O CLI `simulador` verifica porta `8443`; inicia o Simulador; para por `/shutdown`; consulta status por `/api/info`; baixa `simulador.jar` quando necessário; evita download redundante. |
| US-04 | O sistema detecta JDK/JRE compatível; baixa quando ausente; disponibiliza para uso próprio; funciona nas plataformas-alvo. |
| US-05 | O sistema gera binários para Windows, Linux e macOS; publica em GitHub Releases; inclui SHA256; usa SemVer; assina artefatos com Cosign. |

## 11.2 Critérios gerais de aceite do produto

| ID | Critério |
|---|---|
| CA-01 | O projeto compila sem erros localmente. |
| CA-02 | O projeto executa testes automatizados com sucesso. |
| CA-03 | O CI executa lint, build e testes em Windows, Linux e macOS. |
| CA-04 | O CLI `assinatura version` exibe a versão correta. |
| CA-05 | O comando de criação de assinatura simulada funciona em modo local. |
| CA-06 | O comando de validação de assinatura simulada funciona em modo local. |
| CA-07 | O `assinador.jar` expõe endpoints HTTP de assinatura e validação em modo servidor. |
| CA-08 | O CLI `assinatura` consegue usar o `assinador.jar` via HTTP. |
| CA-09 | O CLI detecta instância ativa por health check, não apenas por porta ocupada. |
| CA-10 | O CLI `assinatura stop` encerra o servidor corretamente. |
| CA-11 | O CLI `simulador start` inicia o Simulador e confirma readiness. |
| CA-12 | O CLI `simulador status` informa corretamente se o Simulador está ativo. |
| CA-13 | O CLI `simulador stop` encerra o Simulador corretamente. |
| CA-14 | O sistema baixa JDK/JRE quando ausente e reutiliza quando presente. |
| CA-15 | O sistema baixa `simulador.jar` quando ausente e evita download quando já disponível. |
| CA-16 | Os binários são publicados em release com nomes padronizados. |
| CA-17 | A release contém checksums SHA256. |
| CA-18 | A release contém arquivos `.sig` e `.pem` para verificação Cosign. |
| CA-19 | As mensagens de erro são claras e orientam correção. |
| CA-20 | A documentação explica instalação, execução, testes, build e verificação dos artefatos. |

---

# 12. Priorização dos requisitos

A priorização utiliza a classificação **Alta**, **Média** e **Baixa**, considerando valor ao usuário, risco técnico, dependência entre componentes e relevância para avaliação acadêmica.

## 12.1 Requisitos de prioridade alta

São essenciais para o funcionamento mínimo do sistema:

- RF-01 a RF-05: existência e uso básico do CLI `assinatura`;
- RF-10 a RF-17: existência do `assinador.jar`, validação e simulação;
- RF-20 a RF-25: invocação local e tratamento de erros essenciais;
- RF-26 a RF-35: modo servidor, endpoints, health check e reutilização;
- RF-37 e RF-40: encerramento e tratamento de falhas HTTP;
- RF-41 a RF-50: CLI `simulador` e ciclo de vida do Simulador;
- RF-51 a RF-55 e RF-58: provisionamento e download com integridade;
- RF-59 a RF-61 e RF-63 a RF-67: distribuição, SemVer, checksums e assinatura.

## 12.2 Requisitos de prioridade média

São importantes para robustez, operabilidade e completude:

- RF-06: comando de versão;
- RF-18 e RF-19: suporte e tratamento de PKCS#11;
- RF-36: porta customizada;
- RF-38 e RF-39: encerramento automático por inatividade;
- RF-49: registro operacional do Simulador;
- RF-56: comparação de versões antes do download;
- RF-62: convenção de nome de artefatos;
- requisitos não funcionais de ADRs, logs e padronização de encoding.

## 12.3 Requisitos de prioridade baixa

São úteis, mas podem ser implementados após o núcleo funcional:

- RF-57: URL alternativa de download por `--source`;
- opções avançadas de `--verbose` e `--quiet`;
- refinamentos de UX do help;
- melhorias complementares de empacotamento, desde que não comprometam as entregas obrigatórias.

---

# 13. Restrições técnicas

| ID | Restrição técnica |
|---|---|
| RT-01 | Os CLIs devem ser desenvolvidos em Go 1.25. |
| RT-02 | O `assinador.jar` deve ser desenvolvido em Java 21. |
| RT-03 | O projeto deve usar GitHub como repositório e GitHub Actions como plataforma de CI/CD. |
| RT-04 | A branch principal do repositório deve ser `main`. |
| RT-05 | As plataformas-alvo são `windows/amd64`, `linux/amd64` e `darwin/amd64`. |
| RT-06 | O módulo Go deve seguir a definição técnica do projeto, por exemplo `github.com/kyriosdata/runner`. |
| RT-07 | O CLI deve preferencialmente usar a biblioteca Cobra para estruturação de comandos e flags. |
| RT-08 | A estrutura do repositório deve separar binários em `cmd/`, pacotes internos em `internal/` e projeto Java em diretório próprio. |
| RT-09 | A versão do CLI deve ser injetável em tempo de build por `-ldflags`, portanto deve ser variável, não constante. |
| RT-10 | O build deve ser reproduzível e automatizado. |
| RT-11 | O workflow de build deve executar testes antes de gerar artefatos distribuíveis. |
| RT-12 | O workflow de release deve ser disparado por tags no padrão `v*`. |
| RT-13 | Os binários devem ser nomeados com versão, sistema operacional e arquitetura. |
| RT-14 | Artefatos gerados, caches, diretórios de build e binários não devem ser versionados no repositório. |
| RT-15 | A comunicação entre CLI e `assinador.jar` deve ter contrato documentado e testado. |
| RT-16 | O sistema deve tratar corretamente `stdout`, `stderr` e exit code em subprocessos. |
| RT-17 | O sistema deve evitar dependências desnecessárias, abandonadas ou com vulnerabilidades conhecidas. |
| RT-18 | O projeto deve declarar `.gitignore` e `.gitattributes` adequados. |
| RT-19 | Documentos duplicados da especificação upstream devem ser evitados; quando necessário, referências devem usar commit ou tag fixo. |
| RT-20 | Os artefatos de release devem ser assinados automaticamente por Cosign usando identidade OIDC e transparency log. |

---

# 14. Protótipos de tela

Como o Sistema Runner é um sistema de linha de comando, não há telas gráficas. Os protótipos abaixo representam a experiência esperada no terminal.

## 14.1 Protótipo: ajuda do CLI `assinatura`

```text
$ assinatura --help
Sistema Runner - assinatura

Uso:
  assinatura [comando]

Comandos disponíveis:
  sign        Cria uma assinatura digital simulada
  validate    Valida uma assinatura digital simulada
  stop        Interrompe o assinador.jar em modo servidor
  version     Exibe a versão atual do CLI
  help        Exibe ajuda sobre comandos

Opções gerais:
  --local              Executa o assinador.jar diretamente, sem modo servidor
  --port <porta>       Define a porta do assinador.jar em modo servidor
  --timeout <minutos>  Encerra o servidor após período de inatividade
  --verbose            Exibe informações adicionais de diagnóstico
  --quiet              Reduz mensagens não essenciais
```

## 14.2 Protótipo: criação de assinatura simulada

```text
$ assinatura sign --arquivo documento.xml --perfil padrao
Assinatura simulada criada com sucesso.

Arquivo: documento.xml
Status: sucesso
Modo de execução: servidor
Identificador da assinatura: assinatura-simulada-001
```

## 14.3 Protótipo: validação de assinatura simulada

```text
$ assinatura validate --assinatura assinatura.json
Validação concluída.

Assinatura: assinatura.json
Resultado: válida
Observação: validação simulada para fins de integração.
```

## 14.4 Protótipo: erro de parâmetro

```text
$ assinatura sign --arquivo ""
Erro: parâmetro inválido.

Parâmetro: arquivo
Motivo: o caminho do arquivo não pode ser vazio.
Como resolver: informe um arquivo válido e tente novamente.
```

## 14.5 Protótipo: status do Simulador

```text
$ simulador status
Simulador do HubSaúde: em execução
Porta: 8443
Status: pronto para receber requisições
Endpoint de status: /api/info
```

## 14.6 Protótipo: início do Simulador

```text
$ simulador start
Verificando porta 8443...
Verificando simulador.jar local...
Iniciando Simulador do HubSaúde...
Aguardando readiness...
Simulador iniciado com sucesso.
```

## 14.7 Protótipo: porta ocupada

```text
$ simulador start
Erro: não foi possível iniciar o Simulador.

Motivo: a porta 8443 está ocupada por outro processo.
Como resolver: encerre o processo que está usando a porta ou informe outra porta, se suportado.
```

---

# 15. Casos de teste iniciais

| ID | Caso de teste | Pré-condição | Entrada/Ação | Resultado esperado | Requisito relacionado |
|---|---|---|---|---|---|
| CT-01 | Exibir versão do CLI `assinatura` | CLI compilado | Executar `assinatura version` | A versão atual é exibida. | RF-06 |
| CT-02 | Exibir ajuda do CLI `assinatura` | CLI compilado | Executar `assinatura --help` | Ajuda apresenta comandos e opções. | RF-05 |
| CT-03 | Criar assinatura simulada com parâmetros válidos em modo local | `assinador.jar` disponível | Executar `assinatura sign --local ...` | Assinatura simulada é retornada com sucesso. | RF-13, RF-20 |
| CT-04 | Criar assinatura com parâmetro obrigatório ausente | CLI e JAR disponíveis | Executar comando sem parâmetro obrigatório | Erro claro indica parâmetro ausente. | RF-11, RF-16 |
| CT-05 | Validar assinatura simulada com parâmetros válidos | CLI e JAR disponíveis | Executar `assinatura validate ...` | Resultado simulado de validação é exibido. | RF-15 |
| CT-06 | Invocar `assinador.jar` com argumento contendo espaço e acento | CLI e JAR disponíveis | Usar argumento como `documento teste ação.xml` | Argumento é preservado corretamente. | RF-22 |
| CT-07 | Tratar ausência de `assinador.jar` | Remover ou renomear JAR | Executar `assinatura sign ...` | Mensagem informa que o JAR não foi encontrado. | RF-24 |
| CT-08 | Tratar ausência de Java compatível | Simular Java ausente | Executar comando que exige Java | Sistema baixa Java ou informa erro claro. | RF-25, RF-51, RF-52 |
| CT-09 | Iniciar `assinador.jar` em modo servidor | Porta disponível | Executar comando que dispara servidor | Servidor inicia e health check confirma disponibilidade. | RF-30, RF-32 |
| CT-10 | Reutilizar instância ativa do `assinador.jar` | Servidor já ativo | Executar novo comando `sign` | CLI reutiliza instância existente. | RF-33 |
| CT-11 | Porta ocupada por processo que não é o `assinador.jar` | Porta ocupada por outro processo | Iniciar servidor | Sistema falha com mensagem clara, sem assumir instância válida. | RF-32, RF-40 |
| CT-12 | Encerrar `assinador.jar` | Servidor ativo | Executar `assinatura stop` | Servidor é encerrado e metadados são atualizados. | RF-37 |
| CT-13 | Encerrar por inatividade | Servidor ativo com timeout configurado | Aguardar período sem requisições | Servidor encerra automaticamente. | RF-38 |
| CT-14 | Reiniciar timer de inatividade | Servidor ativo com timeout | Enviar requisição antes do timeout | Timer é reiniciado e servidor permanece ativo. | RF-39 |
| CT-15 | Iniciar Simulador com porta livre | Porta 8443 livre | Executar `simulador start` | Simulador inicia e readiness é confirmado. | RF-42, RF-45, RF-50 |
| CT-16 | Consultar status do Simulador ativo | Simulador ativo | Executar `simulador status` | Status indica execução e readiness. | RF-44, RF-48 |
| CT-17 | Parar Simulador ativo | Simulador ativo | Executar `simulador stop` | Simulador é encerrado com sucesso. | RF-43, RF-47 |
| CT-18 | Baixar `simulador.jar` ausente | JAR ausente e internet disponível | Executar `simulador start` | JAR é baixado, verificado e executado. | RF-55, RF-58 |
| CT-19 | Evitar download redundante do Simulador | JAR local atualizado | Executar `simulador start` | Download não é repetido. | RF-56 |
| CT-20 | Gerar binários multiplataforma | CI configurado | Executar workflow de build | Binários são gerados para Windows, Linux e macOS. | RF-59, RF-60 |
| CT-21 | Publicar release com checksums | Tag `v*` criada | Executar workflow de release | Release contém binários e checksums. | RF-61, RF-64 |
| CT-22 | Assinar artefatos com Cosign | Workflow de release ativo | Publicar release | Arquivos `.sig` e `.pem` são publicados. | RF-65, RF-66, RF-67 |
| CT-23 | Testar integração PKCS#11 com simulador | SoftHSM2 configurado | Executar teste de integração | Chamadas PKCS#11 reais ou simuladas são comprovadas. | RF-18 |
| CT-24 | Tratar resposta HTTP malformada | Servidor retorna resposta inválida | Executar comando HTTP | CLI exibe erro claro e código de saída adequado. | RF-40 |

---

# 16. Glossário

| Termo | Significado no projeto |
|---|---|
| Assinador | Aplicação Java `assinador.jar`, responsável por validar parâmetros e simular assinatura/validação. |
| Assinatura simulada | Resultado fictício usado para testar integração sem realizar criptografia real. |
| CLI `assinatura` | Programa de terminal usado para acionar operações do `assinador.jar`. |
| CLI `simulador` | Programa de terminal usado para gerenciar o Simulador do HubSaúde. |
| Cold start | Inicialização completa da JVM e da aplicação Java a cada chamada. |
| Warm start | Uso de aplicação Java já inicializada, geralmente em modo servidor. |
| Health check | Requisição usada para verificar se um serviço está vivo. |
| Readiness | Confirmação de que o serviço não apenas iniciou, mas está pronto para atender requisições. |
| Modo local | Forma de execução em que o CLI chama diretamente `java -jar assinador.jar`. |
| Modo servidor | Forma de execução em que o `assinador.jar` permanece ativo e recebe requisições HTTP. |
| Runner | Sistema que facilita a execução de aplicações Java por CLI. |
| Simulador do HubSaúde | Aplicação Java/Web usada para simular comportamento relacionado à Plataforma HubSaúde. |
| Artefato | Arquivo gerado e distribuído, como binário, checksum, certificado ou assinatura. |
| Checksum | Valor calculado para verificar se um arquivo foi alterado. |
| Release | Versão publicada do sistema contendo binários e arquivos auxiliares. |
| Rastreabilidade | Capacidade de relacionar requisito, história, issue, PR, commit, código e teste. |

---

# 17. Apêndices e anexos

## 17.1 Apêndice A — Relação com os documentos de origem

Este documento foi consolidado a partir dos seguintes artefatos do projeto:

- documento de especificação do trabalho prático do Sistema Runner;
- documento de design com diagramas C4 de contexto e contêineres;
- critérios de qualidade e avaliação do projeto;
- plano revisado com histórias derivadas e organização por entregas;
- tarefas operacionais com decisões técnicas e estrutura inicial do repositório;
- imagens dos diagramas de contexto e de contêineres.

## 17.2 Apêndice B — Visão arquitetural resumida

### Diagrama de contexto

No nível de contexto, o Sistema Runner se posiciona entre o usuário e sistemas externos. O usuário interage por CLI; o Runner invoca o Simulador do HubSaúde e o Dispositivo de Assinatura Digital quando aplicável.

### Diagrama de contêineres

No nível de contêineres, o sistema possui dois CLIs principais:

- `assinatura`, responsável por interagir com `assinador.jar`;
- `simulador`, responsável por controlar o Simulador do HubSaúde.

O `assinador.jar` pode ser invocado localmente ou via HTTP. O Simulador do HubSaúde é gerenciado via HTTP. O dispositivo de assinatura digital é acessado pelo `assinador.jar` por meio de PKCS#11.

## 17.3 Apêndice C — Estrutura de repositório sugerida

```text
runner/
├── cmd/
│   ├── assinatura/
│   │   └── main.go
│   └── simulador/
│       └── main.go
├── internal/
│   ├── cli/
│   ├── invoker/
│   ├── jdk/
│   └── release/
├── assinador/
│   ├── pom.xml
│   └── src/
├── go.mod
└── go.sum
```

## 17.4 Apêndice D — Convenção de artefatos

Exemplos de nomes de artefatos:

```text
assinatura-v0.1.0-linux-amd64
assinatura-v0.1.0-windows-amd64.exe
assinatura-v0.1.0-darwin-amd64
simulador-v0.1.0-linux-amd64
simulador-v0.1.0-windows-amd64.exe
simulador-v0.1.0-darwin-amd64
checksums.txt
```

Para cada artefato assinado, devem ser publicados também:

```text
<artefato>.sig
<artefato>.pem
```

## 17.5 Apêndice E — Matriz resumida de rastreabilidade

| História | Requisitos principais | Casos de teste iniciais |
|---|---|---|
| US-01 | RF-01 a RF-09, RF-20 a RF-40 | CT-01 a CT-14, CT-24 |
| US-02 | RF-10 a RF-19, RF-26 a RF-29 | CT-03 a CT-06, CT-23 |
| US-03 | RF-41 a RF-50, RF-55 a RF-58 | CT-15 a CT-19 |
| US-04 | RF-51 a RF-54 | CT-08 |
| US-05 | RF-59 a RF-67 | CT-20 a CT-22 |

## 17.6 Apêndice F — Observações de escopo

O Sistema Runner deve ser entendido como um projeto de implementação e integração de software. O valor principal do trabalho está na construção de um sistema utilizável, testável, automatizado, documentado e rastreável. O foco não é a criptografia real, mas sim a integração correta entre CLIs, aplicações Java, modo local, modo servidor HTTP, provisionamento de dependências, releases verificáveis e testes automatizados.
