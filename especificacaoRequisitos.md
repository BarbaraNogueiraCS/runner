# Especificação de Requisitos de Software — Sistema Runner

## 1. Capa

| Campo | Informação |
|---|---|
| **Nome do sistema** | Sistema Runner |
| **Nome do documento** | Especificação de Requisitos de Software |
| **Versão do documento** | 1.0 |
| **Data de elaboração** | 07/05/2026 |
| **Responsável pela elaboração** | Bárbara Nogueira |
| **Instituição / contexto** | Bacharelado em Engenharia de Software — Universidade Federal de Goiás (UFG) |
| **Contexto de aplicação** | Plataforma HubSaúde — interoperabilidade de dados em saúde |
| **Tipo de sistema** | Ferramenta de linha de comandos, integração com aplicações Java e simulação de assinatura digital |

---

## 2. Histórico de versões

| Versão | Data | Autor / Responsável | Descrição da alteração |
|---|---|---|---|
| 1.0 | 07/05/2026 | Equipe do projeto | Elaboração inicial da Especificação de Requisitos de Software do Sistema Runner, com base nos arquivos de especificação, design arquitetural e plano revisitado. |

---

## Sumário

1. [Capa](#1-capa)  
2. [Histórico de versões](#2-histórico-de-versões)  
3. [Introdução](#3-introdução)  
   3.1 [Objetivo do documento](#31-objetivo-do-documento)  
   3.2 [Escopo do sistema](#32-escopo-do-sistema)  
   3.3 [Definições, siglas e abreviações](#33-definições-siglas-e-abreviações)  
4. [Descrição geral do sistema](#4-descrição-geral-do-sistema)  
   4.1 [Contexto do sistema](#41-contexto-do-sistema)  
   4.2 [Objetivos do sistema](#42-objetivos-do-sistema)  
   4.3 [Principais funcionalidades](#43-principais-funcionalidades)  
   4.4 [Classes de usuários](#44-classes-de-usuários)  
   4.5 [Restrições gerais](#45-restrições-gerais)  
   4.6 [Premissas e dependências](#46-premissas-e-dependências)  
5. [Requisitos funcionais](#5-requisitos-funcionais)  
6. [Requisitos não funcionais](#6-requisitos-não-funcionais)  
7. [Regras de negócio](#7-regras-de-negócio)  
8. [Requisitos de interface externa](#8-requisitos-de-interface-externa)  
9. [Requisitos de dados](#9-requisitos-de-dados)  
10. [Histórias de usuário](#10-histórias-de-usuário)  
11. [Critérios de aceitação](#11-critérios-de-aceitação)  
12. [Priorização dos requisitos](#12-priorização-dos-requisitos)  
13. [Restrições técnicas](#13-restrições-técnicas)  
14. [Protótipos de tela](#14-protótipos-de-tela)  
15. [Casos de teste iniciais](#15-casos-de-teste-iniciais)  
16. [Glossário](#16-glossário)  
17. [Apêndices e anexos](#17-apêndices-e-anexos)  

---

## 3. Introdução

O Sistema Runner é uma solução de software voltada para facilitar a execução e o gerenciamento de aplicações Java por meio de interfaces de linha de comandos. O sistema está relacionado ao contexto da Plataforma HubSaúde, iniciativa associada à interoperabilidade de dados em saúde, e tem como propósito reduzir a complexidade técnica enfrentada por usuários e integradores ao executar aplicações Java, especialmente o `assinador.jar` e o `simulador.jar`.

Em vez de exigir que o usuário conheça detalhes de instalação do Java, comandos `java -jar`, portas de execução, parâmetros técnicos, processos em segundo plano e integração com dispositivos criptográficos, o Sistema Runner fornece comandos CLI mais simples e padronizados.

Este documento apresenta a especificação de requisitos do Sistema Runner, descrevendo o escopo, funcionalidades, restrições, regras de negócio, interfaces externas, dados manipulados, histórias de usuário, critérios de aceitação e casos de teste iniciais.

### 3.1 Objetivo do documento

O objetivo deste documento é especificar, de forma clara e estruturada, os requisitos de software do Sistema Runner, servindo como referência para:

- orientar a implementação do sistema;
- apoiar o planejamento das atividades de desenvolvimento;
- permitir rastreabilidade entre requisitos, histórias de usuário e critérios de aceitação;
- apoiar a validação do sistema por meio de testes;
- documentar o comportamento esperado das aplicações `assinatura`, `assinador.jar` e `simulador`;
- registrar as restrições técnicas, funcionais e de escopo do projeto.

Este documento deve ser utilizado por estudantes, desenvolvedores, avaliadores, professores, integradores e demais partes interessadas no desenvolvimento e validação do Sistema Runner.

### 3.2 Escopo do sistema

O Sistema Runner tem como escopo principal fornecer ferramentas de linha de comandos para facilitar a execução de aplicações Java relacionadas ao HubSaúde, especialmente no fluxo de assinatura digital simulada e no gerenciamento do Simulador do HubSaúde.

#### Está no escopo

- desenvolvimento da aplicação CLI `assinatura`;
- desenvolvimento da aplicação Java `assinador.jar`;
- integração entre o CLI `assinatura` e o `assinador.jar`;
- validação rigorosa de parâmetros pelo `assinador.jar`;
- simulação de criação de assinatura digital;
- simulação de validação de assinatura digital;
- tratamento de erros de parâmetros e exceções;
- suporte a invocação local do `assinador.jar` via `java -jar`;
- suporte a invocação do `assinador.jar` via HTTP, em modo servidor;
- gerenciamento do ciclo de vida do `assinador.jar` em modo servidor;
- desenvolvimento da aplicação CLI `simulador`;
- início, parada e consulta de status do Simulador do HubSaúde;
- obtenção dinâmica do `simulador.jar`;
- provisionamento automático do JDK/JRE necessário;
- geração de binários multiplataforma para Windows, Linux e macOS;
- publicação de releases com versionamento semântico;
- geração de checksums SHA256;
- assinatura dos artefatos com Cosign;
- testes unitários, testes de integração e testes de aceitação;
- documentação de uso e documentação técnica da integração.

#### Não está no escopo

- implementação real de assinatura digital criptográfica;
- implementação real de validação criptográfica de assinatura digital;
- integração real com autoridades certificadoras;
- armazenamento persistente de assinaturas digitais;
- criação de interface gráfica;
- autenticação de usuários;
- geração de certificados digitais;
- implementação de um sistema completo de gestão de certificados;
- substituição do Simulador do HubSaúde por uma nova aplicação de negócio.

### 3.3 Definições, siglas e abreviações

| Termo / Sigla | Definição |
|---|---|
| **Runner** | Nome do sistema responsável por facilitar a execução e o gerenciamento de aplicações Java por meio de comandos CLI. |
| **CLI** | Command Line Interface, ou interface de linha de comandos. É uma forma de interagir com o sistema digitando comandos no terminal. |
| **assinatura** | Aplicação CLI desenvolvida em Go, responsável por receber comandos de assinatura e validação e invocar o `assinador.jar`. |
| **assinador.jar** | Aplicação Java responsável por validar parâmetros e simular operações de criação e validação de assinatura digital. |
| **simulador** | Aplicação CLI desenvolvida em Go, responsável por iniciar, parar e consultar o status do Simulador do HubSaúde. |
| **simulador.jar** | Aplicação Java do Simulador do HubSaúde, gerenciada pelo CLI `simulador`. |
| **JDK** | Java Development Kit. Conjunto de ferramentas para executar e desenvolver aplicações Java. |
| **JRE** | Java Runtime Environment. Ambiente necessário para executar aplicações Java. |
| **HTTP** | Protocolo usado para comunicação entre aplicações em rede. |
| **Endpoint** | Ponto de acesso de uma aplicação HTTP, como `/sign`, `/validate`, `/api/info` ou `/shutdown`. |
| **PKCS#11** | Padrão de interface para comunicação com dispositivos criptográficos, como tokens USB e smart cards. |
| **Token / Smart card** | Dispositivo físico usado para armazenar certificados e realizar operações criptográficas. |
| **SemVer** | Versionamento semântico, padrão de versionamento no formato `MAJOR.MINOR.PATCH`, por exemplo, `1.0.0`. |
| **GitHub Releases** | Recurso do GitHub usado para publicar versões executáveis do sistema. |
| **Checksum SHA256** | Código gerado a partir de um arquivo para permitir verificação de integridade. |
| **Cosign** | Ferramenta do ecossistema Sigstore usada para assinar e verificar artefatos de software. |
| **OIDC** | OpenID Connect, mecanismo de identidade usado no processo de assinatura dos artefatos. |
| **Transparency log** | Registro público de assinaturas, usado para rastreabilidade e verificação de autenticidade. |
| **Cold start** | Execução em que a aplicação Java é iniciada do zero a cada comando. |
| **Warm start** | Execução em que a aplicação Java já está em execução, reduzindo o tempo de resposta. |

---

## 4. Descrição geral do sistema

### 4.1 Contexto do sistema

O Sistema Runner está inserido no contexto de uma disciplina de Implementação e Integração de Software do Bacharelado em Engenharia de Software, com aplicação prática relacionada à Plataforma HubSaúde.

A Plataforma HubSaúde está relacionada à interoperabilidade de dados em saúde. Nesse contexto, existem aplicações Java que precisam ser executadas por usuários ou integradores, como o `assinador.jar` e o `simulador.jar`. O Sistema Runner atua como uma camada de simplificação, permitindo que essas aplicações sejam executadas por comandos mais simples, sem exigir conhecimento detalhado sobre Java, instalação de JDK/JRE, portas, processos ou parâmetros internos.

No contexto arquitetural, o sistema envolve os seguintes elementos principais:

| Elemento | Tipo | Descrição |
|---|---|---|
| Usuário | Ator | Pessoa que interage com o sistema por linha de comandos. |
| Sistema Runner | Sistema principal | Conjunto de CLIs e integrações responsáveis por executar e gerenciar aplicações Java. |
| Dispositivo de Assinatura Digital | Sistema externo | Token ou smart card compatível com PKCS#11. |
| Simulador do HubSaúde | Sistema externo | Aplicação Java/Web gerenciada pelo CLI `simulador`. |

A comunicação geral ocorre da seguinte forma:

- o usuário interage com os CLIs `assinatura` e `simulador`;
- o CLI `assinatura` invoca o `assinador.jar` localmente ou via HTTP;
- o `assinador.jar` valida parâmetros e simula operações de assinatura;
- o `assinador.jar` pode interagir com dispositivo criptográfico via PKCS#11;
- o CLI `simulador` inicia, para e monitora o `simulador.jar`;
- o sistema pode baixar automaticamente JDK/JRE e arquivos `.jar` necessários.

### 4.2 Objetivos do sistema

#### Objetivo geral

Facilitar a execução, integração e gerenciamento de aplicações Java relacionadas ao HubSaúde, por meio de interfaces de linha de comandos simples, multiplataforma e adequadas ao uso por integradores e usuários técnicos.

#### Objetivos específicos

- permitir a execução de aplicações Java sem exigir instalação ou configuração manual do Java pelo usuário;
- fornecer comandos CLI simples para criação e validação simulada de assinaturas digitais;
- permitir que o `assinador.jar` seja invocado localmente ou via HTTP;
- permitir que o usuário inicie, pare e monitore o Simulador do HubSaúde;
- validar rigorosamente os parâmetros recebidos pelo `assinador.jar`;
- retornar mensagens de erro claras e orientativas;
- disponibilizar binários multiplataforma para Windows, Linux e macOS;
- automatizar a geração, assinatura e publicação de artefatos executáveis;
- apoiar a qualidade do software por meio de testes, documentação e organização modular.

### 4.3 Principais funcionalidades

| Código | Funcionalidade | Descrição resumida |
|---|---|---|
| F01 | CLI de assinatura | Permite executar comandos para criação e validação simulada de assinaturas. |
| F02 | Simulação de assinatura | Retorna uma assinatura simulada quando os parâmetros são válidos. |
| F03 | Simulação de validação | Retorna resultado simulado de validade da assinatura. |
| F04 | Validação de parâmetros | Verifica presença, formato e consistência dos parâmetros recebidos. |
| F05 | Invocação local do `assinador.jar` | Executa o `assinador.jar` via `java -jar`. |
| F06 | Invocação HTTP do `assinador.jar` | Envia requisições HTTP para `/sign` e `/validate`. |
| F07 | Modo servidor do assinador | Inicia, reutiliza e encerra o `assinador.jar` em segundo plano. |
| F08 | Provisionamento de JDK/JRE | Detecta ou baixa automaticamente o ambiente Java necessário. |
| F09 | CLI do simulador | Permite iniciar, parar e consultar o status do Simulador do HubSaúde. |
| F10 | Download do `simulador.jar` | Obtém dinamicamente a versão mais recente do simulador. |
| F11 | Distribuição multiplataforma | Gera binários para Windows, Linux e macOS. |
| F12 | Segurança de artefatos | Gera checksums SHA256 e assinatura com Cosign. |
| F13 | Testes e documentação | Disponibiliza testes unitários, integração, aceitação e documentação de uso. |

### 4.4 Classes de usuários

| Classe de usuário | Descrição | Necessidades principais |
|---|---|---|
| Usuário do Sistema Runner | Pessoa que usa os CLIs para criar/validar assinaturas ou controlar o simulador. | Comandos simples, mensagens claras e execução sem configuração complexa. |
| Integrador da Plataforma HubSaúde | Profissional que integra sistemas com aplicações do HubSaúde. | Automatizar execução, validar fluxos e controlar aplicações Java. |
| Desenvolvedor do Sistema Runner | Pessoa responsável pela implementação e manutenção do código. | Estrutura organizada, testes, CI/CD e documentação técnica. |
| Avaliador / Professor | Responsável por avaliar o trabalho no contexto acadêmico. | Requisitos claros, rastreabilidade, código funcional e evidências de qualidade. |
| Operador técnico | Pessoa que executa o sistema em ambiente local ou de teste. | Binários prontos, instalação simples, logs e diagnóstico de erros. |

### 4.5 Restrições gerais

- O sistema não deve implementar assinatura digital criptográfica real.
- O sistema não deve implementar validação criptográfica real de assinaturas digitais.
- O sistema não deve gerar certificados digitais.
- O sistema não deve oferecer interface gráfica.
- O sistema deve ser operado por linha de comandos.
- O sistema deve ser compatível com Windows, Linux e macOS em arquitetura amd64.
- Os CLIs devem ser desenvolvidos em Go.
- O `assinador.jar` deve ser desenvolvido em Java 21.
- O sistema deve tratar erros de forma clara e compreensível.
- Os artefatos distribuídos devem usar versionamento semântico.
- As releases devem incluir checksums SHA256 e assinatura com Cosign.
- O sistema deve armazenar arquivos e metadados próprios em diretório gerenciado, preferencialmente `~/.hubsaude/`.

### 4.6 Premissas e dependências

#### Premissas

- O usuário terá acesso a um terminal ou prompt de comando.
- O usuário poderá executar binários locais em seu sistema operacional.
- O sistema poderá acessar a internet quando precisar baixar JDK/JRE, `simulador.jar` ou artefatos de release.
- O `assinador.jar` terá comportamento de simulação, e não de assinatura real.
- O `simulador.jar` será disponibilizado em local acessível, como GitHub Releases ou URL configurável.
- O sistema poderá armazenar arquivos locais em `~/.hubsaude/`.

#### Dependências

- Go 1.25 para desenvolvimento dos CLIs.
- Java 21 para desenvolvimento e execução do `assinador.jar`.
- Biblioteca Cobra para estruturação dos comandos CLI.
- GitHub Actions para CI/CD.
- GitHub Releases para distribuição de binários.
- Eclipse Temurin / Adoptium para obtenção do JDK/JRE.
- Cosign / Sigstore para assinatura e verificação dos artefatos.
- Interface PKCS#11 para comunicação com dispositivo criptográfico.
- Possível uso de SoftHSM2 ou simulador equivalente para testes de integração com PKCS#11.

---

## 5. Requisitos funcionais

| ID | Requisito funcional | Descrição | Origem / Rastreamento |
|---|---|---|---|
| RF-001 | Estruturar CLI `assinatura` | O sistema deve possuir um projeto Go organizado para o CLI `assinatura`, com pacotes, módulo Go e build funcional. | US-01.1 |
| RF-002 | Exibir versão do CLI | O CLI `assinatura` deve possuir comando `version` para exibir a versão atual da aplicação. | US-01.1 |
| RF-003 | Aceitar comando de criação de assinatura | O CLI `assinatura` deve aceitar comando para solicitar criação de assinatura simulada. | US-01 / US-01.2 |
| RF-004 | Aceitar comando de validação de assinatura | O CLI `assinatura` deve aceitar comando para solicitar validação de assinatura simulada. | US-01 / US-01.2 |
| RF-005 | Validar parâmetros no CLI | O CLI deve identificar parâmetros ausentes ou inválidos antes de invocar o `assinador.jar`, quando aplicável. | US-01.2 |
| RF-006 | Invocar `assinador.jar` localmente | O CLI `assinatura` deve executar o `assinador.jar` via `java -jar` quando em modo local. | US-01.3 |
| RF-007 | Capturar saída do `assinador.jar` | O CLI deve capturar a saída de sucesso ou erro do `assinador.jar` e repassá-la ao usuário. | US-01.3 |
| RF-008 | Exibir resultados legíveis | O CLI deve apresentar resultados de criação e validação em formato claro e estruturado. | US-01.4 |
| RF-009 | Iniciar `assinador.jar` em modo servidor | O CLI deve iniciar o `assinador.jar` em modo servidor, usando porta padrão ou porta informada pelo usuário. | US-01.5 |
| RF-010 | Registrar processo do assinador | O CLI deve registrar PID, porta e informações do processo em diretório gerenciado. | US-01.5 |
| RF-011 | Invocar assinador via HTTP | O CLI deve enviar requisições HTTP para os endpoints `/sign` e `/validate` quando o servidor estiver em execução. | US-01.6 |
| RF-012 | Preferir modo servidor quando disponível | O CLI deve reutilizar o `assinador.jar` em modo servidor quando houver instância ativa, salvo orientação contrária. | US-01.6 / US-01.7 |
| RF-013 | Fazer fallback para modo local | O CLI deve permitir fallback para modo local quando o servidor não estiver disponível ou quando a flag correspondente for usada. | US-01.6 |
| RF-014 | Detectar instância ativa do assinador | O CLI deve verificar registro local e health check HTTP para confirmar se o `assinador.jar` está em execução. | US-01.7 |
| RF-015 | Interromper `assinador.jar` | O CLI deve oferecer comando para interromper o `assinador.jar` na porta padrão ou em porta especificada. | US-01.8 |
| RF-016 | Agendar interrupção por inatividade | O sistema deve permitir definir timeout de inatividade para encerramento automático do `assinador.jar`. | US-01.9 |
| RF-017 | Inicializar projeto Java `assinador.jar` | O sistema deve conter projeto Java para o `assinador.jar`, com estrutura adequada. | US-02.1 |
| RF-018 | Definir interface de assinatura | O `assinador.jar` deve possuir interface de serviço com métodos de criação e validação de assinatura. | US-02.1 |
| RF-019 | Simular criação de assinatura | O `assinador.jar` deve retornar assinatura simulada quando receber parâmetros válidos. | US-02.1 |
| RF-020 | Validar parâmetros de criação | O `assinador.jar` deve validar presença e formato dos parâmetros obrigatórios de criação de assinatura. | US-02.2 |
| RF-021 | Retornar erros claros de criação | O `assinador.jar` deve retornar mensagens claras quando parâmetros de criação forem inválidos. | US-02.2 |
| RF-022 | Simular validação de assinatura | O `assinador.jar` deve retornar resultado simulado de validação com base em critérios simples e predeterminados. | US-02.3 |
| RF-023 | Validar parâmetros de validação | O `assinador.jar` deve verificar presença e formato dos parâmetros usados para validação de assinatura. | US-02.3 |
| RF-024 | Expor endpoint `/sign` | O `assinador.jar` deve expor endpoint HTTP `POST /sign` para criação simulada de assinatura. | US-02.4 |
| RF-025 | Expor endpoint `/validate` | O `assinador.jar` deve expor endpoint HTTP `POST /validate` para validação simulada de assinatura. | US-02.4 |
| RF-026 | Padronizar respostas HTTP | O `assinador.jar` deve retornar respostas HTTP consistentes para sucesso e erro. | US-02.4 |
| RF-027 | Suportar PKCS#11 | O `assinador.jar` deve suportar interação com dispositivo criptográfico via PKCS#11. | US-02.5 |
| RF-028 | Tratar ausência de dispositivo criptográfico | O sistema deve informar claramente quando o dispositivo criptográfico não estiver disponível. | US-02.5 |
| RF-029 | Criar CLI `simulador` | O sistema deve possuir CLI dedicado para gerenciamento do Simulador do HubSaúde. | US-03.3 |
| RF-030 | Iniciar Simulador do HubSaúde | O CLI `simulador` deve iniciar o `simulador.jar`. | US-03.1 |
| RF-031 | Verificar porta do simulador | O CLI deve verificar se a porta padrão 8443 está disponível antes de iniciar o simulador. | US-03.1 |
| RF-032 | Parar Simulador do HubSaúde | O CLI `simulador` deve interromper o Simulador do HubSaúde, preferencialmente via endpoint `/shutdown`. | US-03.2 |
| RF-033 | Consultar status do simulador | O CLI `simulador` deve consultar e exibir o status do Simulador, preferencialmente via endpoint `/api/info`. | US-03.2 |
| RF-034 | Registrar processo do simulador | O CLI deve registrar PID, porta e informações do processo do simulador em diretório gerenciado. | US-03.2 |
| RF-035 | Baixar `simulador.jar` dinamicamente | O CLI deve baixar automaticamente a versão mais recente do `simulador.jar` quando não estiver disponível localmente. | US-03.4 |
| RF-036 | Permitir URL alternativa do simulador | O CLI deve aceitar parâmetro `--source <url>` para indicar origem alternativa de download do `simulador.jar`. | US-03.4 |
| RF-037 | Evitar download repetido | O CLI não deve baixar novamente o `simulador.jar` quando a versão mais recente já estiver em cache local. | US-03.4 |
| RF-038 | Verificar integridade do download | O sistema deve verificar checksum do `simulador.jar` baixado. | US-03.4 / US-05.3 |
| RF-039 | Detectar JDK/JRE disponível | O sistema deve verificar se há JDK/JRE compatível no PATH ou em diretório gerenciado. | US-04.1 |
| RF-040 | Baixar JDK/JRE automaticamente | O sistema deve baixar JDK/JRE compatível quando ausente. | US-04.1 |
| RF-041 | Armazenar JDK/JRE localmente | O sistema deve armazenar o JDK/JRE baixado em diretório gerenciado para reuso. | US-04.1 |
| RF-042 | Gerar binários multiplataforma | O pipeline deve gerar binários para Windows, Linux e macOS. | US-05.1 |
| RF-043 | Publicar releases versionadas | O sistema deve publicar binários em GitHub Releases com versionamento semântico. | US-05.2 |
| RF-044 | Gerar checksums SHA256 | Cada release deve incluir checksums SHA256 para os artefatos. | US-05.3 |
| RF-045 | Assinar artefatos com Cosign | Cada artefato distribuído deve ser assinado com Cosign. | US-05.3 |
| RF-046 | Disponibilizar arquivos `.sig` e `.pem` | Cada artefato assinado deve possuir arquivos de assinatura e certificado correspondentes. | US-05.3 |
| RF-047 | Documentar verificação de artefatos | A documentação deve explicar como verificar artefatos com `cosign verify-blob`. | US-05.3 |

---

## 6. Requisitos não funcionais

| ID | Categoria | Requisito não funcional | Critério de qualidade esperado |
|---|---|---|---|
| RNF-001 | Usabilidade | Os comandos CLI devem ser simples, consistentes e acompanhados de ajuda por `--help`. | Usuário deve compreender comandos e parâmetros sem consultar código-fonte. |
| RNF-002 | Usabilidade | Mensagens de erro devem ser claras, orientativas e indicar o parâmetro ou causa do problema. | Erros devem sugerir correção quando possível. |
| RNF-003 | Portabilidade | Os CLIs devem funcionar em Windows, Linux e macOS. | Binários devem ser gerados para `windows/amd64`, `linux/amd64` e `darwin/amd64`. |
| RNF-004 | Manutenibilidade | O código deve ser organizado em pacotes, módulos e serviços com responsabilidades claras. | Estrutura deve facilitar evolução e testes. |
| RNF-005 | Testabilidade | O sistema deve possuir testes unitários, de integração e de aceitação. | Funcionalidades principais e cenários de erro devem estar cobertos. |
| RNF-006 | Confiabilidade | O sistema deve tratar falhas de execução, portas ocupadas, arquivos ausentes e dependências indisponíveis. | Falhas não devem causar encerramentos sem mensagem compreensível. |
| RNF-007 | Segurança | Artefatos distribuídos devem permitir verificação de integridade e autenticidade. | Releases devem incluir checksums SHA256 e assinatura Cosign. |
| RNF-008 | Desempenho | O modo servidor deve reduzir o tempo de resposta em múltiplas operações sucessivas. | Quando servidor estiver ativo, CLI deve preferir invocação HTTP. |
| RNF-009 | Interoperabilidade | O `assinador.jar` deve expor endpoints HTTP e suporte a PKCS#11 conforme especificado. | Integração deve ocorrer por interfaces padronizadas. |
| RNF-010 | Reusabilidade | O JDK/JRE e arquivos `.jar` baixados devem ser armazenados para reuso. | Downloads não devem ser repetidos sem necessidade. |
| RNF-011 | Rastreabilidade | Requisitos devem estar associados a histórias de usuário e critérios de aceitação. | Cada requisito deve poder ser vinculado a uma necessidade do projeto. |
| RNF-012 | Observabilidade operacional | O sistema deve registrar informações mínimas de processo, como PID, porta e status. | Usuário deve conseguir consultar ou encerrar processos iniciados pelo Runner. |
| RNF-013 | Compatibilidade | O `assinador.jar` deve ser compatível com Java 21. | Execução deve ocorrer com JDK/JRE 21. |
| RNF-014 | Evolutibilidade | O sistema deve permitir inclusão futura de novas aplicações Java gerenciadas por CLI. | Estrutura deve separar responsabilidades e evitar acoplamento desnecessário. |
| RNF-015 | Documentação | O sistema deve possuir documentação técnica e manual de uso. | Deve incluir instalação, comandos, exemplos, erros comuns e verificação de artefatos. |

---

## 7. Regras de negócio

| ID | Regra de negócio | Descrição |
|---|---|---|
| RN-001 | Operações de assinatura são simuladas | O sistema deve simular criação e validação de assinatura, sem realizar assinatura criptográfica real. |
| RN-002 | Parâmetros devem ser validados antes do processamento | Nenhuma operação de assinatura ou validação deve ser processada sem validação dos parâmetros obrigatórios. |
| RN-003 | Erros devem ser explicativos | Quando um parâmetro estiver inválido, a mensagem deve indicar qual parâmetro falhou e o motivo. |
| RN-004 | Modo servidor deve ser reutilizado | Quando houver instância ativa do `assinador.jar`, o CLI deve reutilizá-la, salvo se o usuário solicitar modo local. |
| RN-005 | Porta ocupada impede inicialização | O sistema não deve iniciar o simulador ou assinador em porta já ocupada sem informar o problema ao usuário. |
| RN-006 | Porta padrão do Simulador é 8443 | O CLI `simulador` deve considerar a porta 8443 como padrão para o Simulador do HubSaúde. |
| RN-007 | Download não deve ser repetido sem necessidade | JDK/JRE e `simulador.jar` não devem ser baixados novamente se versão compatível já estiver disponível localmente. |
| RN-008 | Diretório local gerenciado | Arquivos baixados, metadados, versões e registros de processo devem ser mantidos em diretório gerenciado, como `~/.hubsaude/`. |
| RN-009 | Releases devem usar SemVer | As versões publicadas devem seguir versionamento semântico, como `v0.1.0` ou `v1.0.0`. |
| RN-010 | Artefatos devem ser verificáveis | Cada artefato de release deve possuir checksum SHA256 e assinatura Cosign. |
| RN-011 | CLI deve ocultar complexidade Java | O usuário não deve precisar conhecer comandos `java -jar` para utilizar o sistema. |
| RN-012 | Simulador deve ser controlado pelo CLI | O usuário deve iniciar, parar e consultar o status do Simulador por meio do CLI `simulador`. |
| RN-013 | Fallback deve preservar operação | Quando o modo servidor não estiver disponível, o sistema deve permitir execução local, quando aplicável. |
| RN-014 | Ausência de dispositivo criptográfico deve ser tratada | O sistema deve informar claramente quando token, smart card ou interface PKCS#11 não estiverem disponíveis. |

---

## 8. Requisitos de interface externa

### 8.1 Interface com o usuário — CLI

O sistema deve fornecer comandos de linha de comando para interação com o usuário.

#### CLI `assinatura`

Comandos previstos:

```bash
assinatura version
assinatura sign [parâmetros]
assinatura validate [parâmetros]
assinatura stop [--port <porta>]
assinatura --help
```

Parâmetros específicos devem ser definidos conforme o formato esperado pelos casos de uso de criação e validação de assinatura do HubSaúde. O CLI deve apresentar ajuda, mensagens de sucesso, mensagens de erro e resultados estruturados.

#### CLI `simulador`

Comandos previstos:

```bash
simulador start [--port <porta>] [--source <url>]
simulador stop [--port <porta>]
simulador status [--port <porta>]
simulador --help
```

O CLI deve informar claramente se o simulador foi iniciado, se já estava em execução, se foi encerrado ou se ocorreu erro.

### 8.2 Interface com `assinador.jar` em modo local

O CLI `assinatura` deve invocar o `assinador.jar` por meio de chamada local ao Java:

```bash
java -jar assinador.jar [comando] [parâmetros]
```

O CLI deve:

- localizar o executável Java adequado;
- montar corretamente os argumentos;
- executar o processo;
- capturar saída padrão e saída de erro;
- apresentar a resposta ao usuário.

### 8.3 Interface HTTP com `assinador.jar`

Quando em modo servidor, o `assinador.jar` deve expor endpoints HTTP:

| Método | Endpoint | Finalidade |
|---|---|---|
| POST | `/sign` | Criar assinatura simulada. |
| POST | `/validate` | Validar assinatura simulada. |

As respostas devem possuir estrutura consistente para sucesso e erro.

Exemplo conceitual de resposta de sucesso:

```json
{
  "success": true,
  "operation": "sign",
  "signature": "assinatura-simulada",
  "message": "Assinatura simulada criada com sucesso."
}
```

Exemplo conceitual de resposta de erro:

```json
{
  "success": false,
  "error": "Parâmetro obrigatório ausente: documento"
}
```

### 8.4 Interface HTTP com o Simulador do HubSaúde

O CLI `simulador` deve interagir com o Simulador do HubSaúde por meio de endpoints HTTP.

| Método | Endpoint | Finalidade |
|---|---|---|
| GET | `/api/info` | Consultar informações e status do simulador. |
| POST ou GET, conforme implementação do simulador | `/shutdown` | Solicitar encerramento do simulador. |

### 8.5 Interface com dispositivo criptográfico

O `assinador.jar` deve possuir suporte à interface PKCS#11 para comunicação com dispositivos criptográficos, como token USB ou smart card.

O sistema deve tratar, no mínimo, os seguintes cenários:

- dispositivo disponível;
- dispositivo ausente;
- biblioteca PKCS#11 não configurada;
- falha na comunicação com o dispositivo;
- uso de simulador equivalente para testes, como SoftHSM2.

### 8.6 Interface com serviços externos de download

O sistema deve acessar serviços externos para baixar dependências e artefatos.

| Serviço externo | Uso |
|---|---|
| GitHub Releases | Baixar binários, `assinador.jar` ou `simulador.jar`. |
| URL `release.json` | Identificar versão e URL mais recentes do `simulador.jar`. |
| Eclipse Temurin / Adoptium | Baixar JDK/JRE compatível com a plataforma. |

### 8.7 Interface com sistema de arquivos

O sistema deve utilizar diretório local gerenciado para armazenar arquivos e metadados.

Estrutura conceitual:

```text
~/.hubsaude/
  jdk/
  jre/
  assinador/
  simulador/
  processos/
  cache/
  logs/
```

---

## 9. Requisitos de dados

| ID | Dado | Descrição | Persistência esperada |
|---|---|---|---|
| RD-001 | Versão do CLI | Número da versão atual do CLI `assinatura` ou `simulador`. | Embutido no binário e exibido por comando `version`. |
| RD-002 | Parâmetros de assinatura | Dados informados pelo usuário para criação de assinatura simulada. | Não persistente, usado durante a execução. |
| RD-003 | Parâmetros de validação | Dados informados pelo usuário para validação de assinatura simulada. | Não persistente, usado durante a execução. |
| RD-004 | Resultado de assinatura simulada | Resposta gerada pelo `assinador.jar` para uma criação simulada. | Não persistente, exibido ao usuário. |
| RD-005 | Resultado de validação simulada | Indicação simulada de assinatura válida ou inválida. | Não persistente, exibido ao usuário. |
| RD-006 | Registro de processo | PID, porta, aplicação e status de processos iniciados. | Persistente em `~/.hubsaude/processos/`. |
| RD-007 | Metadados do JDK/JRE | Versão, caminho local e plataforma do JDK/JRE baixado. | Persistente em `~/.hubsaude/`. |
| RD-008 | Metadados do `simulador.jar` | Versão, URL de origem, caminho local e checksum. | Persistente em cache local. |
| RD-009 | Arquivo `release.json` | Informações sobre versão e URL do `simulador.jar` e JDK/JRE. | Consultado remotamente e possivelmente armazenado em cache. |
| RD-010 | Checksums SHA256 | Hash dos artefatos distribuídos. | Publicado em release e usado para verificação. |
| RD-011 | Assinaturas Cosign | Arquivos `.sig` e `.pem` associados aos artefatos. | Publicados em release. |
| RD-012 | Logs operacionais | Informações de erro e execução, quando implementado. | Opcional, em diretório de logs local. |

### 9.1 Qualidade dos dados

- Parâmetros obrigatórios não podem estar vazios.
- Parâmetros devem respeitar formatos definidos pelas operações de assinatura e validação.
- Registros de processo devem ser atualizados após parada de aplicações.
- Checksums devem corresponder aos arquivos baixados.
- Versões em cache devem ser comparadas com versões remotas antes de novos downloads.

---

## 10. Histórias de usuário

### US-01 — Invocar `assinador.jar` via CLI

**Como** usuário do Sistema Runner,  
**quero** executar comandos de assinatura digital através da linha de comandos,  
**para que** eu possa invocar o `assinador.jar` sem conhecer detalhes técnicos de configuração Java, tanto para assinar quanto para validar assinaturas digitais.

#### Histórias derivadas

| ID | História |
|---|---|
| US-01.1 | Como usuário, quero que o projeto CLI esteja estruturado com organização de pacotes e build funcional, para que o desenvolvimento progrida de forma organizada. |
| US-01.2 | Como usuário, quero executar comandos `sign` e `validate` com parâmetros via linha de comandos, para solicitar operações de assinatura de forma intuitiva. |
| US-01.3 | Como usuário, quero que o CLI invoque o `assinador.jar` via `java -jar`, para não executar comandos Java manualmente. |
| US-01.4 | Como usuário, quero que o CLI apresente resultados legíveis, para compreender facilmente o retorno da operação. |
| US-01.5 | Como usuário, quero que o CLI inicie o `assinador.jar` em modo servidor, para reduzir latência em múltiplas requisições. |
| US-01.6 | Como usuário, quero que o CLI envie requisições HTTP ao `assinador.jar`, para usar o modo servidor por padrão quando disponível. |
| US-01.7 | Como usuário, quero que o CLI detecte instância ativa do `assinador.jar`, para evitar processos duplicados. |
| US-01.8 | Como usuário, quero interromper o `assinador.jar`, para controlar processos em execução. |
| US-01.9 | Como usuário, quero agendar interrupção por inatividade, para liberar recursos automaticamente. |

### US-02 — Simular assinatura digital com validação de parâmetros

**Como** usuário do Sistema Runner,  
**quero** que o `assinador.jar` valide rigorosamente os parâmetros de entrada antes de simular uma operação de assinatura digital,  
**para que** eu receba feedback imediato sobre erros e apenas requisições bem formadas sejam aceitas.

#### Histórias derivadas

| ID | História |
|---|---|
| US-02.1 | Como usuário, quero que o `assinador.jar` retorne uma assinatura simulada quando receber parâmetros válidos, para testar o fluxo sem infraestrutura criptográfica real. |
| US-02.2 | Como usuário, quero que o `assinador.jar` valide os parâmetros de criação de assinatura, para receber feedback claro sobre erros. |
| US-02.3 | Como usuário, quero que o `assinador.jar` valide parâmetros de validação e retorne resultado predeterminado, para testar o fluxo de validação. |
| US-02.4 | Como usuário, quero que o `assinador.jar` exponha endpoints HTTP `/sign` e `/validate`, para permitir invocação via servidor. |
| US-02.5 | Como usuário, quero que o `assinador.jar` suporte interação com dispositivo criptográfico via PKCS#11, para permitir uso real ou simulado de material criptográfico. |

### US-03 — Gerenciar ciclo de vida do Simulador do HubSaúde

**Como** usuário do Sistema Runner,  
**quero** iniciar, parar e monitorar o Simulador do HubSaúde através do CLI,  
**para que** eu possa gerenciar o ciclo de vida do Simulador sem conhecer comandos Java subjacentes.

#### Histórias derivadas

| ID | História |
|---|---|
| US-03.1 | Como usuário, quero iniciar o Simulador via CLI, para gerenciá-lo sem conhecer comandos Java. |
| US-03.2 | Como usuário, quero parar e monitorar o Simulador, para ter visibilidade e controle do processo. |
| US-03.3 | Como usuário, quero um CLI dedicado para o Simulador, para ter interface independente e clara. |
| US-03.4 | Como usuário, quero que o CLI baixe automaticamente o `simulador.jar`, para usar a versão atualizada sem download manual. |

### US-04 — Provisionar JDK automaticamente

**Como** usuário do Sistema Runner,  
**quero** que o sistema baixe e configure automaticamente o JDK necessário quando este não estiver disponível,  
**para que** eu possa utilizar o Assinador e o Simulador sem instalar ou configurar Java manualmente.

#### História derivada

| ID | História |
|---|---|
| US-04.1 | Como usuário, quero que o sistema detecte JDK compatível e, caso esteja ausente, baixe e configure automaticamente. |

### US-05 — Disponibilizar binários multiplataforma

**Como** usuário do Sistema Runner,  
**quero** baixar uma versão pré-compilada do CLI para minha plataforma,  
**para que** eu possa utilizar o sistema imediatamente sem necessidade de compilação.

#### Histórias derivadas

| ID | História |
|---|---|
| US-05.1 | Como desenvolvedor, quero pipeline CI/CD multiplataforma, para gerar binários automaticamente. |
| US-05.2 | Como usuário, quero releases com versionamento semântico, para baixar binários com versões claras. |
| US-05.3 | Como usuário, quero checksums SHA256 e assinatura Cosign, para verificar integridade e autenticidade dos artefatos. |

---

## 11. Critérios de aceitação

### CA-US-01 — Invocar `assinador.jar` via CLI

| ID | Critério de aceitação |
|---|---|
| CA-01.1 | O CLI deve aceitar comandos para criação e validação de assinatura. |
| CA-01.2 | O CLI deve invocar o `assinador.jar` com os parâmetros fornecidos. |
| CA-01.3 | O CLI deve permitir invocação direta/local do `assinador.jar`. |
| CA-01.4 | O CLI deve permitir invocação via HTTP do `assinador.jar`. |
| CA-01.5 | O CLI deve exibir o resultado da operação em formato legível. |
| CA-01.6 | O CLI deve iniciar o `assinador.jar` no modo servidor usando a porta padrão quando não orientado de forma diferente. |
| CA-01.7 | O CLI deve detectar instância em execução e reutilizá-la quando possível. |
| CA-01.8 | O CLI deve permitir interromper o `assinador.jar` na porta padrão ou em porta informada. |
| CA-01.9 | O CLI deve permitir interrupção programada após tempo de inatividade. |

### CA-US-02 — Simulação e validação de assinatura

| ID | Critério de aceitação |
|---|---|
| CA-02.1 | O `assinador.jar` deve validar todos os parâmetros obrigatórios. |
| CA-02.2 | O `assinador.jar` deve simular criação de assinatura com resposta pré-construída quando os parâmetros forem válidos. |
| CA-02.3 | O `assinador.jar` deve simular validação de assinatura com resultado predeterminado. |
| CA-02.4 | O `assinador.jar` deve suportar interação com dispositivo criptográfico via PKCS#11. |
| CA-02.5 | O `assinador.jar` deve retornar mensagens claras para parâmetros inválidos. |
| CA-02.6 | Os endpoints `/sign` e `/validate` devem reutilizar a mesma lógica de validação e simulação do modo local. |

### CA-US-03 — Gerenciamento do Simulador

| ID | Critério de aceitação |
|---|---|
| CA-03.1 | O CLI deve verificar se a porta padrão 8443 está disponível antes de iniciar o Simulador. |
| CA-03.2 | O CLI deve permitir iniciar o Simulador. |
| CA-03.3 | O CLI deve permitir parar o Simulador. |
| CA-03.4 | O CLI deve exibir o status atual do Simulador. |
| CA-03.5 | O CLI deve obter dinamicamente a versão mais recente do `simulador.jar`. |
| CA-03.6 | O CLI não deve baixar o `simulador.jar` se a versão mais recente já estiver disponível localmente. |
| CA-03.7 | O CLI deve registrar PID e porta do processo em diretório gerenciado. |

### CA-US-04 — Provisionamento de JDK/JRE

| ID | Critério de aceitação |
|---|---|
| CA-04.1 | O sistema deve detectar JDK/JRE compatível no PATH ou em diretório gerenciado. |
| CA-04.2 | O sistema deve baixar JDK/JRE compatível quando ausente. |
| CA-04.3 | O JDK/JRE baixado deve ser armazenado localmente para reuso. |
| CA-04.4 | O download deve funcionar em Windows, Linux e macOS. |

### CA-US-05 — Distribuição multiplataforma

| ID | Critério de aceitação |
|---|---|
| CA-05.1 | O pipeline deve gerar binários para Windows, Linux e macOS. |
| CA-05.2 | As releases devem ser publicadas com versionamento semântico. |
| CA-05.3 | Os artefatos devem ser publicados em GitHub Releases. |
| CA-05.4 | Cada release deve incluir checksums SHA256. |
| CA-05.5 | Cada artefato deve possuir assinatura Cosign. |
| CA-05.6 | Cada artefato assinado deve possuir arquivos `.sig` e `.pem`. |
| CA-05.7 | A documentação deve explicar como verificar a integridade e autenticidade dos artefatos. |

---

## 12. Priorização dos requisitos

A priorização utiliza uma classificação inspirada em MoSCoW:

- **Obrigatório**: necessário para que o sistema cumpra seu objetivo principal.
- **Importante**: agrega valor significativo, mas pode ser implementado após o núcleo funcional.
- **Desejável**: melhora a qualidade ou a experiência, mas não bloqueia a versão inicial.

| ID | Requisito / Grupo | Prioridade | Justificativa |
|---|---|---|---|
| RF-001 | Estrutura do CLI `assinatura` | Obrigatório | Base para todos os comandos de assinatura. |
| RF-002 | Comando `version` | Importante | Apoia versionamento e suporte ao usuário. |
| RF-003 a RF-005 | Parsing de `sign` e `validate` | Obrigatório | Essencial para interação do usuário com o sistema. |
| RF-006 a RF-008 | Invocação local e exibição de resultados | Obrigatório | Permite uso básico do `assinador.jar`. |
| RF-017 a RF-023 | Simulação e validação no `assinador.jar` | Obrigatório | Núcleo funcional da simulação de assinatura e validação. |
| RF-039 a RF-041 | Detecção e provisionamento de JDK/JRE | Obrigatório | Atende ao objetivo de ocultar configuração Java. |
| RF-009 a RF-016 | Modo servidor do assinador | Importante | Melhora desempenho e experiência em múltiplas requisições. |
| RF-024 a RF-026 | Endpoints HTTP `/sign` e `/validate` | Importante | Necessário para o modo servidor. |
| RF-027 a RF-028 | Suporte PKCS#11 | Importante | Alinha o sistema ao contexto de dispositivos criptográficos. |
| RF-029 a RF-037 | CLI `simulador` e gestão do simulador | Obrigatório | Faz parte do escopo principal do Sistema Runner. |
| RF-038 | Verificação de integridade do download | Importante | Reduz risco de uso de artefatos corrompidos. |
| RF-042 a RF-047 | CI/CD, releases, checksums e Cosign | Importante | Garante distribuição segura e multiplataforma. |
| RNF-001 a RNF-015 | Requisitos não funcionais | Variável | Devem ser atendidos progressivamente durante o desenvolvimento. |

### 12.1 Sugestão de implementação por etapas

#### Etapa 1 — Base funcional mínima

- Estrutura do CLI `assinatura`.
- Comando `version`.
- Comandos `sign` e `validate`.
- Projeto Java `assinador.jar`.
- Simulação básica de assinatura e validação.
- Validação inicial de parâmetros.
- Invocação local via `java -jar`.

#### Etapa 2 — Robustez e automação

- Tratamento completo de erros.
- Provisionamento automático de JDK/JRE.
- Formatação legível dos resultados.
- Testes unitários e de integração.

#### Etapa 3 — Modo servidor e integração HTTP

- Endpoints `/sign` e `/validate`.
- Inicialização do `assinador.jar` em segundo plano.
- Reutilização de instância ativa.
- Stop e timeout por inatividade.

#### Etapa 4 — Simulador e distribuição

- CLI `simulador`.
- Start, stop e status do Simulador.
- Download dinâmico do `simulador.jar`.
- CI/CD multiplataforma.
- Releases com checksums e Cosign.

---

## 13. Restrições técnicas

| ID | Restrição técnica | Descrição |
|---|---|---|
| RT-001 | Linguagem dos CLIs | Os CLIs `assinatura` e `simulador` devem ser desenvolvidos em Go. |
| RT-002 | Versão do Go | O desenvolvimento dos CLIs deve considerar Go 1.25, conforme premissa do projeto. |
| RT-003 | Framework CLI | O CLI deve utilizar Cobra ou estrutura equivalente prevista no projeto. |
| RT-004 | Linguagem do `assinador.jar` | O `assinador.jar` deve ser desenvolvido em Java. |
| RT-005 | Versão Java | O `assinador.jar` deve ser compatível com Java 21. |
| RT-006 | Diretório local | O sistema deve utilizar diretório gerenciado, preferencialmente `~/.hubsaude/`. |
| RT-007 | Plataformas suportadas | O sistema deve suportar Windows, Linux e macOS em arquitetura amd64. |
| RT-008 | Invocação local | A execução local deve ocorrer por `java -jar assinador.jar`. |
| RT-009 | Invocação HTTP | O modo servidor deve utilizar HTTP para comunicação com o `assinador.jar`. |
| RT-010 | Porta do Simulador | A porta padrão do Simulador do HubSaúde deve ser 8443. |
| RT-011 | Endpoints do Simulador | O status deve ser consultado em `/api/info` e a parada deve ocorrer por `/shutdown`, conforme disponibilidade do simulador. |
| RT-012 | PKCS#11 | A comunicação com dispositivo criptográfico deve usar interface PKCS#11. |
| RT-013 | CI/CD | A geração de binários deve ser automatizada por GitHub Actions. |
| RT-014 | Releases | A distribuição deve ocorrer por GitHub Releases. |
| RT-015 | Segurança de artefatos | Os artefatos devem ser acompanhados de SHA256, `.sig` e `.pem`. |
| RT-016 | Sem interface gráfica | O sistema deve operar sem GUI, exclusivamente por linha de comandos. |
| RT-017 | Sem persistência de assinaturas | O sistema não deve armazenar assinaturas digitais de forma persistente. |
| RT-018 | Sem autoridade certificadora | O sistema não deve integrar com autoridades certificadoras reais. |

---

## 14. Protótipos de tela

Como o Sistema Runner não possui interface gráfica, os protótipos representam telas de terminal. Os exemplos abaixo são conceituais e servem para orientar a experiência esperada do usuário.

### 14.1 Protótipo — ajuda do CLI `assinatura`

```text
$ assinatura --help

Sistema Runner - CLI de Assinatura

Uso:
  assinatura [comando]

Comandos disponíveis:
  sign        Cria uma assinatura digital simulada
  validate    Valida uma assinatura digital simulada
  stop        Interrompe o assinador em modo servidor
  version     Exibe a versão atual do CLI
  help        Exibe ajuda sobre os comandos

Opções globais:
  --local            Força execução local via java -jar
  --port <porta>     Define a porta do assinador em modo servidor
  --timeout <min>    Define tempo de inatividade para encerramento automático
```

### 14.2 Protótipo — criação de assinatura simulada com sucesso

```text
$ assinatura sign --documento documento.json --certificado certificado.pem

Operação: criação de assinatura
Status: sucesso
Mensagem: Assinatura simulada criada com sucesso.
Assinatura: assinatura-simulada-2026
Modo de execução: servidor HTTP
```

### 14.3 Protótipo — erro por parâmetro ausente

```text
$ assinatura sign --documento documento.json

Erro: parâmetro obrigatório ausente.
Parâmetro: certificado
Orientação: informe o certificado ou consulte 'assinatura sign --help'.
```

### 14.4 Protótipo — validação de assinatura simulada

```text
$ assinatura validate --documento documento.json --assinatura assinatura.txt

Operação: validação de assinatura
Status: sucesso
Resultado: assinatura válida
Mensagem: A assinatura simulada foi validada conforme os critérios configurados.
```

### 14.5 Protótipo — ajuda do CLI `simulador`

```text
$ simulador --help

Sistema Runner - CLI do Simulador HubSaúde

Uso:
  simulador [comando]

Comandos disponíveis:
  start      Inicia o Simulador do HubSaúde
  stop       Para o Simulador do HubSaúde
  status     Consulta o status atual do Simulador
  help       Exibe ajuda sobre os comandos

Opções:
  --port <porta>      Define a porta do simulador. Padrão: 8443
  --source <url>      Define URL alternativa para download do simulador.jar
```

### 14.6 Protótipo — iniciar Simulador

```text
$ simulador start

Verificando porta 8443...
Porta disponível.
Verificando simulador.jar...
Versão local encontrada: 1.2.0
Iniciando Simulador do HubSaúde...

Status: em execução
Porta: 8443
PID: 15342
```

### 14.7 Protótipo — status do Simulador

```text
$ simulador status

Simulador do HubSaúde: em execução
Porta: 8443
PID: 15342
Endpoint de status: http://localhost:8443/api/info
```

### 14.8 Protótipo — porta ocupada

```text
$ simulador start

Erro: não foi possível iniciar o Simulador.
Motivo: a porta 8443 já está em uso.
Orientação: finalize o processo atual ou informe outra porta com --port.
```

---

## 15. Casos de teste iniciais

| ID | Caso de teste | Pré-condição | Passos | Resultado esperado | Requisito relacionado |
|---|---|---|---|---|---|
| CT-001 | Exibir versão do CLI `assinatura` | CLI instalado | Executar `assinatura version` | Sistema exibe versão atual do CLI. | RF-002 |
| CT-002 | Exibir ajuda do CLI `assinatura` | CLI instalado | Executar `assinatura --help` | Sistema exibe comandos e opções disponíveis. | RNF-001 |
| CT-003 | Criar assinatura simulada com parâmetros válidos | `assinador.jar` disponível | Executar comando `sign` com parâmetros válidos | Sistema retorna assinatura simulada com status de sucesso. | RF-003, RF-019 |
| CT-004 | Criar assinatura com parâmetro ausente | CLI instalado | Executar `sign` sem parâmetro obrigatório | Sistema informa parâmetro ausente e orientação de correção. | RF-005, RF-020, RF-021 |
| CT-005 | Validar assinatura com parâmetros válidos | `assinador.jar` disponível | Executar comando `validate` com parâmetros válidos | Sistema retorna resultado simulado de validação. | RF-004, RF-022 |
| CT-006 | Validar assinatura com parâmetro inválido | CLI instalado | Executar `validate` com parâmetro em formato inválido | Sistema rejeita requisição e informa o motivo. | RF-023 |
| CT-007 | Invocar `assinador.jar` localmente | Java disponível | Executar `assinatura sign --local ...` | CLI executa `java -jar assinador.jar` e exibe resultado. | RF-006, RF-007 |
| CT-008 | Executar sem Java instalado | Java ausente no PATH e em `~/.hubsaude/` | Executar comando que exige Java | Sistema inicia provisionamento automático do JDK/JRE. | RF-039, RF-040 |
| CT-009 | Reutilizar JDK/JRE já baixado | JDK/JRE presente em `~/.hubsaude/` | Executar comando que exige Java | Sistema usa JDK/JRE local sem novo download. | RF-041, RN-007 |
| CT-010 | Iniciar `assinador.jar` em modo servidor | Porta disponível | Executar comando para iniciar servidor | Sistema inicia processo e registra PID/porta. | RF-009, RF-010 |
| CT-011 | Reutilizar instância ativa do assinador | Servidor já em execução | Executar `assinatura sign ...` | CLI usa servidor ativo via HTTP. | RF-011, RF-012, RF-014 |
| CT-012 | Parar `assinador.jar` | Servidor em execução | Executar `assinatura stop` | Sistema encerra processo e atualiza registro local. | RF-015 |
| CT-013 | Encerrar por timeout | Servidor iniciado com timeout | Aguardar período de inatividade | Sistema encerra servidor automaticamente. | RF-016 |
| CT-014 | Iniciar Simulador com porta livre | Porta 8443 livre | Executar `simulador start` | Sistema inicia simulador e registra PID/porta. | RF-030, RF-031 |
| CT-015 | Iniciar Simulador com porta ocupada | Porta 8443 ocupada | Executar `simulador start` | Sistema informa erro de porta ocupada. | RF-031, RN-005 |
| CT-016 | Consultar status do Simulador | Simulador em execução | Executar `simulador status` | Sistema exibe status, porta e PID. | RF-033 |
| CT-017 | Parar Simulador | Simulador em execução | Executar `simulador stop` | Sistema encerra simulador e atualiza registro. | RF-032, RF-034 |
| CT-018 | Baixar `simulador.jar` ausente | `simulador.jar` não existe localmente | Executar `simulador start` | Sistema baixa o arquivo antes de iniciar. | RF-035 |
| CT-019 | Evitar download repetido do simulador | Versão mais recente já em cache | Executar `simulador start` | Sistema não baixa novamente o arquivo. | RF-037 |
| CT-020 | Verificar checksum de artefato | Arquivo baixado e checksum disponível | Executar fluxo de verificação | Sistema confirma integridade ou rejeita arquivo divergente. | RF-038, RF-044 |
| CT-021 | Verificar artefato com Cosign | Artefato, `.sig` e `.pem` disponíveis | Executar comando de verificação documentado | Assinatura é validada ou erro é reportado. | RF-045, RF-046, RF-047 |
| CT-022 | Dispositivo PKCS#11 ausente | Nenhum token/smart card disponível | Executar fluxo que exige dispositivo | Sistema informa ausência do dispositivo de forma clara. | RF-027, RF-028 |

---

## 16. Glossário

| Termo | Significado |
|---|---|
| Artefato | Arquivo produzido pelo processo de desenvolvimento, como binário, checksum, documento ou pacote. |
| Binário | Arquivo executável pronto para uso pelo usuário. |
| Build | Processo de compilação e empacotamento do software. |
| Cache local | Armazenamento temporário ou reutilizável de arquivos baixados, evitando novos downloads. |
| Checksum | Valor calculado a partir de um arquivo para verificar se ele foi alterado ou corrompido. |
| CLI | Interface de linha de comandos usada por meio do terminal. |
| Endpoint | Rota HTTP exposta por uma aplicação. |
| Fallback | Estratégia alternativa usada quando a principal não está disponível. |
| GitHub Actions | Serviço de automação usado para executar pipelines de build, teste e release. |
| GitHub Releases | Área do GitHub usada para publicar versões distribuíveis de um projeto. |
| Health check | Verificação usada para saber se um serviço está ativo e respondendo. |
| Integração | Comunicação entre partes diferentes de um sistema ou entre sistemas externos. |
| JAR | Formato de arquivo usado para empacotar aplicações Java. |
| JDK | Kit de desenvolvimento Java, necessário para desenvolver e executar aplicações Java. |
| JRE | Ambiente de execução Java, necessário para executar aplicações Java. |
| Multiplataforma | Capacidade de funcionar em diferentes sistemas operacionais. |
| PID | Identificador de processo no sistema operacional. |
| Porta | Número usado para identificar um serviço de rede em execução. |
| PKCS#11 | Padrão de comunicação com dispositivos criptográficos. |
| Release | Versão publicada do software. |
| SemVer | Padrão de versionamento semântico, geralmente no formato `MAJOR.MINOR.PATCH`. |
| Simulação | Reprodução controlada de um comportamento sem executar a operação real. |
| Smart card | Cartão físico usado para armazenar certificados ou chaves criptográficas. |
| Token | Dispositivo físico usado para armazenar certificados ou chaves criptográficas. |

---

## 17. Apêndices e anexos

### Apêndice A — Rastreabilidade entre épicos e histórias

| Épico | Descrição | Histórias derivadas |
|---|---|---|
| US-01 | Invocar `assinador.jar` via CLI | US-01.1, US-01.2, US-01.3, US-01.4, US-01.5, US-01.6, US-01.7, US-01.8, US-01.9 |
| US-02 | Simular assinatura digital com validação | US-02.1, US-02.2, US-02.3, US-02.4, US-02.5 |
| US-03 | Gerenciar ciclo de vida do Simulador | US-03.1, US-03.2, US-03.3, US-03.4 |
| US-04 | Provisionar JDK automaticamente | US-04.1 |
| US-05 | Disponibilizar binários multiplataforma | US-05.1, US-05.2, US-05.3 |

### Apêndice B — Fluxo lógico de criação de assinatura

```text
Usuário → assinatura CLI → assinador.jar → assinatura CLI → Usuário

1. Usuário executa comando para criar assinatura.
2. CLI recebe e interpreta os parâmetros.
3. CLI valida parâmetros básicos.
4. CLI localiza JDK/JRE compatível.
5. CLI invoca o assinador.jar localmente ou via HTTP.
6. assinador.jar valida os parâmetros recebidos.
7. assinador.jar retorna assinatura simulada.
8. CLI formata o resultado.
9. Usuário visualiza o resultado no terminal.
```

### Apêndice C — Fluxo lógico de validação de assinatura

```text
Usuário → assinatura CLI → assinador.jar → assinatura CLI → Usuário

1. Usuário executa comando para validar assinatura.
2. CLI recebe e interpreta os parâmetros.
3. CLI valida parâmetros básicos.
4. CLI localiza JDK/JRE compatível.
5. CLI invoca o assinador.jar localmente ou via HTTP.
6. assinador.jar valida os parâmetros recebidos.
7. assinador.jar retorna resultado simulado de validação.
8. CLI formata o resultado.
9. Usuário visualiza se a assinatura é válida ou inválida.
```

### Apêndice D — Exemplo conceitual de `release.json`

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

### Apêndice E — Exemplo de artefatos esperados em release

```text
assinatura-1.0.0-windows-amd64.exe
assinatura-1.0.0-linux-amd64
assinatura-1.0.0-darwin-amd64
simulador-1.0.0-windows-amd64.exe
simulador-1.0.0-linux-amd64
simulador-1.0.0-darwin-amd64
checksums.txt
assinatura-1.0.0-linux-amd64.sig
assinatura-1.0.0-linux-amd64.pem
simulador-1.0.0-linux-amd64.sig
simulador-1.0.0-linux-amd64.pem
```

### Apêndice F — Referências documentais utilizadas

- Documento de especificação do Sistema Runner.
- Documento de design arquitetural do Sistema Runner baseado no Modelo C4.
- Plano revisitado do Sistema Runner com histórias de usuário, critérios de aceitação e rastreabilidade.
- Referências técnicas associadas aos casos de uso de criação e validação de assinatura no contexto FHIR/HubSaúde.
- Boas práticas de engenharia de software para requisitos, rastreabilidade, teste, documentação e integração.
