# Documento de Gerenciamento de Projeto — Sistema Runner

## 1. Identificação do documento

| Campo | Informação |
|---|---|
| **Nome do sistema** | Sistema Runner |
| **Nome do documento** | Documento de Gerenciamento de Projet |
| **Versão do documento** | 1.0 |
| **Data de elaboração** | 08/05/2026 |
| **Responsável pela elaboração** | Equipe do projeto / Disciplina de Implementação e Integração de Software |
| **Instituição / contexto acadêmico** | Bacharelado em Engenharia de Software — Universidade Federal de Goiás (UFG) |
| **Contexto de aplicação** | Plataforma HubSaúde — interoperabilidade de dados em saúde |
| **Documentos relacionados** | Especificação de Requisitos de Software; Documento de Arquitetura de Software; Documento de Projeto Detalhado de Software; Documento de Modelo C4; Documento de Implementação e Integração; Documento de Teste de Software; Documento de Implantação de Software; Plano Revisado #2; README da implementação de referência |
| **Tipo de sistema** | Ferramenta de linha de comandos, integração com aplicações Java, gerenciamento de processos, provisionamento de ambiente Java e simulação de assinatura digital |

---

## 2. Histórico de versões

| Versão | Data | Autor / Responsável | Descrição da alteração |
|---|---|---|---|
| 1.0 | 08/05/2026 | Equipe do projeto | Elaboração inicial do Documento de Gerenciamento de Projeto com Cronograma do Sistema Runner, organizado em 8 entregas principais, com tarefas, artefatos, critérios de aceite, riscos, qualidade e acompanhamento do desenvolvimento iterativo e incremental. |

---

## 3. Sumário

1. [Identificação do documento](#1-identificação-do-documento)  
2. [Histórico de versões](#2-histórico-de-versões)  
3. [Sumário](#3-sumário)  
4. [Introdução](#4-introdução)  
5. [Objetivo do documento](#5-objetivo-do-documento)  
6. [Visão geral do projeto](#6-visão-geral-do-projeto)  
7. [Escopo do gerenciamento do projeto](#7-escopo-do-gerenciamento-do-projeto)  
8. [Estratégia de desenvolvimento](#8-estratégia-de-desenvolvimento)  
9. [Premissas, restrições e dependências](#9-premissas-restrições-e-dependências)  
10. [Papéis e responsabilidades](#10-papéis-e-responsabilidades)  
11. [Estrutura analítica do projeto](#11-estrutura-analítica-do-projeto)  
12. [Cronograma com 8 entregas principais](#12-cronograma-com-8-entregas-principais)  
13. [Detalhamento das entregas](#13-detalhamento-das-entregas)  
14. [Matriz de artefatos por entrega](#14-matriz-de-artefatos-por-entrega)  
15. [Rastreabilidade entre entregas, requisitos e histórias de usuário](#15-rastreabilidade-entre-entregas-requisitos-e-histórias-de-usuário)  
16. [Plano de qualidade](#16-plano-de-qualidade)  
17. [Plano de testes por entrega](#17-plano-de-testes-por-entrega)  
18. [Plano de comunicação](#18-plano-de-comunicação)  
19. [Gestão de riscos](#19-gestão-de-riscos)  
20. [Gestão de mudanças](#20-gestão-de-mudanças)  
21. [Gestão de configuração e versionamento](#21-gestão-de-configuração-e-versionamento)  
22. [Indicadores de acompanhamento](#22-indicadores-de-acompanhamento)  
23. [Critérios de aceite do projeto](#23-critérios-de-aceite-do-projeto)  
24. [Plano de encerramento](#24-plano-de-encerramento)  
25. [Referências](#25-referências)  

---

## 4. Introdução

O Sistema Runner é uma solução de software desenvolvida no contexto da disciplina de Implementação e Integração de Software, relacionada à Plataforma HubSaúde. Seu propósito é facilitar a execução e o gerenciamento de aplicações Java por meio de interfaces de linha de comandos, reduzindo a necessidade de o usuário conhecer detalhes técnicos de instalação do Java, execução de arquivos `.jar`, configuração de portas, controle de processos, comunicação HTTP e validação de parâmetros.

O sistema é composto principalmente pelos CLIs `assinatura` e `simulador`, desenvolvidos em Go, e pela aplicação Java `assinador.jar`, desenvolvida em Java 21. O CLI `assinatura` permite criar e validar assinaturas simuladas, invocando o `assinador.jar` em modo local ou HTTP. O CLI `simulador` permite iniciar, parar e consultar o status do Simulador do HubSaúde. O projeto também utiliza o diretório local `~/.hubsaude/` para armazenar artefatos, metadados, cache, registros de processos e dependências.

Este documento apresenta o gerenciamento do projeto, com foco na definição de cronograma, distribuição de tarefas e artefatos em 8 entregas principais. A organização foi feita de forma iterativa e incremental, permitindo que o sistema evolua por partes, com validação contínua da qualidade, integração progressiva e geração de evidências ao longo do desenvolvimento.

---

## 5. Objetivo do documento

O objetivo deste documento é planejar e organizar o desenvolvimento do Sistema Runner, estabelecendo:

- as entregas principais do projeto;
- o cronograma de desenvolvimento em 8 entregas;
- as tarefas previstas em cada entrega;
- os artefatos esperados em cada etapa;
- os critérios de aceite por entrega;
- a rastreabilidade com requisitos, histórias de usuário e documentos técnicos;
- as responsabilidades da equipe;
- a estratégia de qualidade, testes, integração e implantação;
- os riscos, indicadores e mecanismos de acompanhamento.

Este documento deve apoiar o planejamento, a execução, o monitoramento e o encerramento do projeto, servindo como referência para estudantes, desenvolvedores, avaliadores e demais envolvidos.

---

## 6. Visão geral do projeto

### 6.1 Descrição resumida

O projeto consiste no desenvolvimento de uma solução local, operada por linha de comandos, composta por:

| Parte | Tecnologia | Finalidade |
|---|---|---|
| **CLI `assinatura`** | Go + Cobra | Receber comandos de criação e validação de assinatura simulada e invocar o `assinador.jar`. |
| **`assinador.jar`** | Java 21 | Validar parâmetros, simular criação/validação de assinatura e expor endpoints HTTP. |
| **CLI `simulador`** | Go + Cobra | Iniciar, parar e consultar status do Simulador do HubSaúde. |
| **Diretório `~/.hubsaude/`** | Sistema de arquivos | Armazenar JDK/JRE, arquivos `.jar`, cache, metadados e registros de processo. |
| **CI/CD** | GitHub Actions | Compilar, testar, gerar binários multiplataforma, checksums e assinaturas de artefatos. |

### 6.2 Objetivo geral do projeto

Desenvolver e integrar o Sistema Runner com qualidade, permitindo que usuários executem aplicações Java relacionadas ao HubSaúde por comandos simples, com validação de parâmetros, tratamento de erros, testes, documentação, distribuição multiplataforma e segurança básica de artefatos.

### 6.3 Objetivos específicos

- Construir os CLIs `assinatura` e `simulador` em Go.
- Construir o `assinador.jar` em Java 21.
- Implementar simulação de criação e validação de assinatura digital.
- Implementar validação rigorosa de parâmetros.
- Integrar CLI e Java por modo local via `java -jar`.
- Integrar CLI e Java por modo servidor HTTP.
- Gerenciar processos locais com PID, porta, status e health check.
- Controlar o Simulador do HubSaúde pelo CLI `simulador`.
- Estruturar download, cache e verificação de artefatos.
- Preparar CI/CD, releases, checksums SHA256 e assinatura com Cosign.
- Elaborar documentação técnica, testes e plano de implantação.

---

## 7. Escopo do gerenciamento do projeto

### 7.1 Está no escopo

- Planejamento do projeto em 8 entregas principais.
- Organização iterativa e incremental das tarefas.
- Priorização do núcleo funcional antes das funcionalidades complementares.
- Desenvolvimento dos CLIs `assinatura` e `simulador`.
- Desenvolvimento do `assinador.jar`.
- Integração local e HTTP.
- Gerenciamento de processos e diretório `~/.hubsaude/`.
- Estrutura de testes unitários, integração, aceitação e implantação.
- Documentação de requisitos, arquitetura, projeto detalhado, C4, implementação, testes e implantação.
- Preparação de CI/CD e artefatos de distribuição.

### 7.2 Não está no escopo

- Implementação real de assinatura digital criptográfica.
- Validação criptográfica real de assinatura.
- Integração real com autoridades certificadoras.
- Geração de certificados digitais.
- Autenticação de usuários.
- Interface gráfica.
- Banco de dados relacional ou não relacional.
- Implantação distribuída em ambiente de produção.
- Garantia jurídica de validade de assinatura digital.

### 7.3 Limitações consideradas no planejamento

- A assinatura e a validação são simuladas.
- A integração PKCS#11 é estruturada, mas não executa assinatura real nesta versão.
- O `simulador.jar` real pode não estar incluso no pacote e pode precisar ser informado por `--jar` ou `--source`.
- O provisionamento automático completo de JDK/JRE pode ser evoluído incrementalmente.
- A execução em Windows, Linux e macOS deve ser validada progressivamente por CI/CD.

---

## 8. Estratégia de desenvolvimento

O desenvolvimento será conduzido de forma **iterativa e incremental**.

### 8.1 Desenvolvimento iterativo

A cada entrega, o sistema será revisado, testado e ajustado. Isso permite corrigir problemas cedo, melhorar a documentação e reduzir riscos técnicos.

### 8.2 Desenvolvimento incremental

Cada entrega adiciona funcionalidades novas ao sistema. O projeto começa com a base documental e estrutural, avança para a implementação dos componentes principais, depois integra os fluxos, adiciona automação, realiza testes e prepara a implantação.

### 8.3 Estratégia de qualidade contínua

A qualidade será acompanhada durante todo o projeto, por meio de:

- revisão de requisitos e escopo;
- padronização de estrutura de pastas;
- separação de responsabilidades;
- validação de parâmetros;
- tratamento de erros;
- testes automatizados;
- testes de integração;
- documentação contínua;
- rastreabilidade entre requisitos, código, testes e entregas;
- versionamento semântico;
- geração de evidências.

---

## 9. Premissas, restrições e dependências

### 9.1 Premissas

| ID | Premissa |
|---|---|
| PR-001 | A equipe terá acesso a ambiente de desenvolvimento com terminal, Go, Java 21 e Git. |
| PR-002 | Os CLIs serão desenvolvidos em Go. |
| PR-003 | O `assinador.jar` será desenvolvido em Java 21. |
| PR-004 | O sistema será usado localmente, por linha de comandos. |
| PR-005 | O diretório `~/.hubsaude/` será usado para cache, registros e artefatos locais. |
| PR-006 | As operações de assinatura e validação serão simuladas. |
| PR-007 | O projeto será controlado por Git e poderá usar GitHub Actions para CI/CD. |
| PR-008 | Os documentos já elaborados serão usados como base de rastreabilidade. |

### 9.2 Restrições

| ID | Restrição |
|---|---|
| RE-001 | O sistema não deve implementar assinatura digital real. |
| RE-002 | O sistema não deve possuir interface gráfica. |
| RE-003 | O sistema deve funcionar por linha de comandos. |
| RE-004 | O projeto deve manter coerência com os documentos de requisitos, arquitetura, projeto detalhado, Modelo C4, implementação, testes e implantação. |
| RE-005 | Os artefatos distribuíveis devem seguir versionamento semântico. |
| RE-006 | As releases devem prever checksums SHA256 e assinatura com Cosign. |
| RE-007 | O Simulador do HubSaúde deve usar a porta padrão 8443, salvo configuração alternativa. |
| RE-008 | O `assinador.jar` em modo servidor deve expor endpoints compatíveis com os fluxos `/sign`, `/validate`, `/health` e `/shutdown`. |

### 9.3 Dependências

| ID | Dependência | Impacto |
|---|---|---|
| DP-001 | Go | Necessário para compilar os CLIs. |
| DP-002 | Java 21 | Necessário para compilar e executar o `assinador.jar`. |
| DP-003 | Maven ou script alternativo | Necessário para empacotar e testar o projeto Java. |
| DP-004 | GitHub Actions | Necessário para automação de build e release. |
| DP-005 | GitHub Releases | Necessário para publicação de artefatos. |
| DP-006 | Cosign | Necessário para assinatura/verificação de artefatos. |
| DP-007 | `simulador.jar` real ou URL configurada | Necessário para teste completo do CLI `simulador`. |
| DP-008 | Acesso à internet | Necessário para downloads automáticos e dependências externas. |

---

## 10. Papéis e responsabilidades

| Papel | Responsabilidades principais |
|---|---|
| **Gerente do projeto / líder da equipe** | Planejar entregas, acompanhar cronograma, organizar prioridades, consolidar evidências e controlar riscos. |
| **Analista de requisitos** | Garantir coerência entre requisitos, histórias de usuário, critérios de aceitação e entregas. |
| **Arquiteto de software** | Manter coerência arquitetural, modularização, integração local/HTTP e organização C4. |
| **Desenvolvedor Go** | Implementar CLIs `assinatura` e `simulador`, comandos, casos de uso, processos e artefatos. |
| **Desenvolvedor Java** | Implementar `assinador.jar`, validação, simulação, endpoints HTTP e adaptador PKCS#11. |
| **Responsável por testes** | Planejar, executar e registrar testes unitários, integração, aceitação e implantação. |
| **Responsável por DevOps/CI-CD** | Configurar workflows, builds multiplataforma, releases, checksums e assinatura de artefatos. |
| **Responsável por documentação** | Atualizar documentos técnicos, guias de uso, evidências e relatório final. |
| **Avaliador / professor** | Validar entregas, critérios de aceite, qualidade e coerência do projeto. |

---

## 11. Estrutura analítica do projeto

A Estrutura Analítica do Projeto, EAP, organiza o trabalho em pacotes de entrega.

```text
Sistema Runner
├── 1. Planejamento e documentação inicial
│   ├── requisitos
│   ├── arquitetura
│   ├── projeto detalhado
│   └── Modelo C4
├── 2. Base técnica e ambiente
│   ├── repositório
│   ├── estrutura Go
│   ├── estrutura Java
│   └── CI inicial
├── 3. CLI assinatura
│   ├── comandos
│   ├── validação preliminar
│   ├── formatação
│   └── invocação local
├── 4. assinador.jar
│   ├── validação Java
│   ├── simulação de assinatura
│   ├── simulação de validação
│   └── testes Java
├── 5. Integração HTTP e processos
│   ├── endpoints
│   ├── start/stop
│   ├── health check
│   └── registro em ~/.hubsaude
├── 6. CLI simulador
│   ├── start
│   ├── status
│   ├── stop
│   └── download/cached artifacts
├── 7. Qualidade, segurança e distribuição
│   ├── testes
│   ├── checksums
│   ├── Cosign
│   ├── GitHub Releases
│   └── documentação técnica
└── 8. Implantação, validação final e entrega
    ├── testes finais
    ├── implantação local
    ├── evidências
    └── encerramento
```

---

## 12. Cronograma com 8 entregas principais

O cronograma abaixo está organizado em 8 entregas principais, considerando um ciclo acadêmico de desenvolvimento iterativo e incremental. Cada entrega pode corresponder a uma semana, sprint ou marco de acompanhamento, conforme a organização da disciplina.

| Entrega | Período sugerido | Nome da entrega | Objetivo principal | Resultado esperado |
|---|---|---|---|---|
| **E1** | Semana 1 | Planejamento, requisitos e arquitetura inicial | Consolidar entendimento do problema, escopo, requisitos e arquitetura de alto nível. | Documentos iniciais revisados e backlog priorizado. |
| **E2** | Semana 2 | Estrutura base do projeto e ambiente | Criar estrutura do repositório, CLIs, projeto Java e CI inicial. | Projeto compilável em estrutura mínima. |
| **E3** | Semana 3 | Núcleo do `assinador.jar` | Implementar validação e simulação de assinatura/validação no Java. | `assinador.jar` com lógica central e testes unitários. |
| **E4** | Semana 4 | CLI `assinatura` e integração local | Implementar comandos principais e invocação local via `java -jar`. | Fluxo local de assinatura e validação funcionando. |
| **E5** | Semana 5 | Modo servidor HTTP e gerenciamento de processos | Implementar endpoints, start/stop, health check e registro de processo. | Integração HTTP funcional entre CLI e `assinador.jar`. |
| **E6** | Semana 6 | CLI `simulador` e gerenciamento de artefatos | Implementar start, status, stop, download/cache e verificação de porta. | CLI `simulador` funcional com estrutura de download e controle. |
| **E7** | Semana 7 | Qualidade, segurança e distribuição | Consolidar testes, CI/CD, checksums, Cosign, releases e documentação de implementação/testes. | Artefatos testados, verificáveis e prontos para distribuição. |
| **E8** | Semana 8 | Implantação, validação final e encerramento | Executar testes finais, implantação local, rollback, evidências e entrega final. | Sistema validado, documentado e entregue com evidências. |

---

## 13. Detalhamento das entregas

## 13.1 Entrega E1 — Planejamento, requisitos e arquitetura inicial

### Objetivo

Consolidar a visão do projeto, o escopo, as funcionalidades esperadas, os requisitos, os riscos iniciais e a arquitetura de alto nível do Sistema Runner.

### Tarefas principais

| ID | Tarefa |
|---|---|
| E1-T01 | Analisar a especificação original do Sistema Runner. |
| E1-T02 | Identificar o escopo: CLI `assinatura`, `assinador.jar`, CLI `simulador`, integração e distribuição. |
| E1-T03 | Consolidar requisitos funcionais e não funcionais. |
| E1-T04 | Elaborar ou revisar histórias de usuário e critérios de aceitação. |
| E1-T05 | Elaborar arquitetura inicial com visão C4 de contexto e contêineres. |
| E1-T06 | Definir riscos iniciais e premissas técnicas. |
| E1-T07 | Priorizar backlog inicial do projeto. |

### Artefatos esperados

- Especificação de Requisitos de Software.
- Documento de Arquitetura de Software.
- Documento de Modelo C4.
- Backlog inicial.
- Matriz de rastreabilidade inicial.
- Lista inicial de riscos.

### Critérios de aceite

- Escopo do sistema claramente definido.
- Itens fora do escopo explicitamente registrados.
- Requisitos principais vinculados às histórias de usuário.
- Arquitetura inicial coerente com os elementos do sistema.
- Backlog organizado por prioridade.

### Atividades de qualidade

- Revisão dos requisitos por pelo menos um membro da equipe.
- Validação de coerência entre escopo, requisitos e arquitetura.
- Registro de dúvidas e decisões técnicas.

---

## 13.2 Entrega E2 — Estrutura base do projeto e ambiente

### Objetivo

Preparar a base técnica do projeto, com estrutura de repositório, módulos iniciais, pontos de entrada dos CLIs, projeto Java e automação inicial.

### Tarefas principais

| ID | Tarefa |
|---|---|
| E2-T01 | Criar repositório Git do projeto. |
| E2-T02 | Definir organização de pastas e arquivos. |
| E2-T03 | Inicializar módulo Go do projeto. |
| E2-T04 | Criar ponto de entrada do CLI `assinatura`. |
| E2-T05 | Criar ponto de entrada do CLI `simulador`. |
| E2-T06 | Criar estrutura base do projeto Java `assinador.jar`. |
| E2-T07 | Configurar comando `version` no CLI `assinatura`. |
| E2-T08 | Configurar workflows iniciais de build. |
| E2-T09 | Criar diretório de documentação e scripts auxiliares. |

### Artefatos esperados

- Repositório Git estruturado.
- Estrutura base em Go.
- Estrutura base Java.
- Comando `assinatura version`.
- Estrutura inicial de CI/CD.
- Documento de Projeto Detalhado atualizado.
- README inicial com instruções de build.

### Critérios de aceite

- Projeto Go compila sem erros.
- Projeto Java possui estrutura de build definida.
- Comando `assinatura version` executa corretamente.
- Estrutura de pastas corresponde ao projeto detalhado.
- Pipeline inicial de build está definido.

### Atividades de qualidade

- Padronização de nomes de pastas e pacotes.
- Revisão da organização do código.
- Verificação de compilação local.
- Registro de decisões técnicas.

---

## 13.3 Entrega E3 — Núcleo do `assinador.jar`

### Objetivo

Implementar a lógica central do `assinador.jar`, incluindo validação de parâmetros, simulação de criação de assinatura e simulação de validação de assinatura.

### Tarefas principais

| ID | Tarefa |
|---|---|
| E3-T01 | Criar DTOs `SignRequest`, `SignResponse`, `ValidateRequest`, `ValidateResponse` e `ErrorResponse`. |
| E3-T02 | Implementar interface `SignatureService`. |
| E3-T03 | Implementar `FakeSignatureService`. |
| E3-T04 | Implementar `ParameterValidator`. |
| E3-T05 | Implementar tratamento de erro de parâmetros inválidos. |
| E3-T06 | Implementar modo CLI básico do `assinador.jar`. |
| E3-T07 | Criar testes unitários para validação de parâmetros. |
| E3-T08 | Criar testes unitários para serviço de assinatura simulada. |
| E3-T09 | Estruturar adaptador PKCS#11 sem assinatura real. |

### Artefatos esperados

- Código Java do `assinador.jar`.
- Classes DTO.
- Serviço de assinatura simulada.
- Validador de parâmetros.
- Testes unitários Java.
- Build do `assinador.jar`.
- Documentação técnica atualizada.

### Critérios de aceite

- `assinador.jar` compila.
- Parâmetros obrigatórios são validados.
- Entradas inválidas retornam erro claro.
- Entradas válidas retornam assinatura simulada.
- Validação simulada retorna resultado predeterminado.
- Testes unitários Java passam.

### Atividades de qualidade

- Cobertura de cenários de sucesso e erro.
- Revisão da clareza das mensagens de erro.
- Separação entre validação, serviço e DTOs.
- Garantia explícita de que a assinatura é simulada.

---

## 13.4 Entrega E4 — CLI `assinatura` e integração local

### Objetivo

Implementar os comandos principais do CLI `assinatura` e integrá-los ao `assinador.jar` em modo local por meio de `java -jar`.

### Tarefas principais

| ID | Tarefa |
|---|---|
| E4-T01 | Implementar comando `assinatura sign`. |
| E4-T02 | Implementar comando `assinatura validate`. |
| E4-T03 | Implementar validação preliminar de parâmetros no Go. |
| E4-T04 | Implementar resolução do executável Java. |
| E4-T05 | Implementar invocação local por `java -jar`. |
| E4-T06 | Capturar stdout, stderr e código de saída do processo Java. |
| E4-T07 | Implementar formatação de resultados no terminal. |
| E4-T08 | Implementar tratamento de erro para JAR inexistente, Java ausente e parâmetros inválidos. |
| E4-T09 | Criar testes Go para validação de comandos e invocação local. |

### Artefatos esperados

- CLI `assinatura` com comandos `sign` e `validate`.
- Módulo de invocação local.
- Módulo de resolução de Java.
- Formatador de saída.
- Testes Go iniciais.
- Evidências de execução local.
- README atualizado com exemplos.

### Critérios de aceite

- `assinatura sign --local` invoca o `assinador.jar`.
- `assinatura validate --local` invoca o `assinador.jar`.
- Saídas de sucesso são legíveis.
- Erros são claros e orientativos.
- O usuário não precisa digitar `java -jar` manualmente.
- Fluxos locais possuem evidências de teste.

### Atividades de qualidade

- Testes de parâmetros ausentes.
- Testes de caminho inválido do JAR.
- Testes de Java indisponível, quando possível.
- Revisão de usabilidade das mensagens CLI.

---

## 13.5 Entrega E5 — Modo servidor HTTP e gerenciamento de processos

### Objetivo

Implementar o modo servidor do `assinador.jar`, os endpoints HTTP e o gerenciamento de processos pelo CLI `assinatura`.

### Tarefas principais

| ID | Tarefa |
|---|---|
| E5-T01 | Implementar servidor HTTP no `assinador.jar`. |
| E5-T02 | Implementar endpoint `GET /health`. |
| E5-T03 | Implementar endpoint `POST /sign`. |
| E5-T04 | Implementar endpoint `POST /validate`. |
| E5-T05 | Implementar endpoint `/shutdown`. |
| E5-T06 | Implementar comando `assinatura start`. |
| E5-T07 | Implementar comando `assinatura stop`. |
| E5-T08 | Implementar verificação de porta disponível. |
| E5-T09 | Implementar registro de processo em `~/.hubsaude/processos/`. |
| E5-T10 | Implementar health check antes de reutilizar instância. |
| E5-T11 | Implementar cliente HTTP no Go. |
| E5-T12 | Criar testes de integração HTTP. |

### Artefatos esperados

- `assinador.jar` com modo servidor.
- Endpoints `/health`, `/sign`, `/validate` e `/shutdown`.
- CLI `assinatura start`.
- CLI `assinatura stop`.
- Módulo de processo e registro local.
- Módulo de health check.
- Testes de integração HTTP.
- Evidências de execução HTTP.

### Critérios de aceite

- Servidor Java inicia em porta configurável.
- Processo é registrado em `~/.hubsaude/processos/`.
- CLI consegue enviar requisições HTTP ao `assinador.jar`.
- Endpoints retornam respostas coerentes.
- `assinatura stop` encerra o processo.
- Registros desatualizados não são reutilizados sem health check.

### Atividades de qualidade

- Teste de porta ocupada.
- Teste de servidor indisponível.
- Teste de JSON inválido.
- Teste de reutilização de instância ativa.
- Verificação de atualização do registro de processo.

---

## 13.6 Entrega E6 — CLI `simulador` e gerenciamento de artefatos

### Objetivo

Implementar o CLI `simulador` para iniciar, parar e consultar status do Simulador do HubSaúde, incluindo suporte a JAR local, download por URL, cache e verificação de porta.

### Tarefas principais

| ID | Tarefa |
|---|---|
| E6-T01 | Implementar comando `simulador start`. |
| E6-T02 | Implementar comando `simulador status`. |
| E6-T03 | Implementar comando `simulador stop`. |
| E6-T04 | Definir porta padrão 8443. |
| E6-T05 | Permitir parâmetro `--jar` para JAR local. |
| E6-T06 | Permitir parâmetro `--source` para URL alternativa. |
| E6-T07 | Implementar download genérico de artefatos por URL. |
| E6-T08 | Implementar cache local em `~/.hubsaude/`. |
| E6-T09 | Implementar verificação SHA256 para artefatos, quando disponível. |
| E6-T10 | Integrar status via endpoint `/api/info`, quando disponível. |
| E6-T11 | Integrar parada via endpoint `/shutdown`, quando disponível. |
| E6-T12 | Criar testes do CLI `simulador`. |

### Artefatos esperados

- CLI `simulador`.
- Comandos `start`, `status` e `stop`.
- Módulo de download de artefatos.
- Módulo de verificação SHA256.
- Registro de processo do simulador.
- Testes e evidências do fluxo do simulador.
- Documentação de uso do CLI `simulador`.

### Critérios de aceite

- `simulador start` inicia o JAR informado.
- Porta 8443 é usada como padrão.
- Porta ocupada é detectada.
- `simulador status` informa situação do processo.
- `simulador stop` encerra o processo quando possível.
- Download por `--source` é tratado com mensagem clara.
- Artefatos podem ser verificados por checksum.

### Atividades de qualidade

- Teste de JAR inexistente.
- Teste de URL inválida.
- Teste de porta ocupada.
- Teste de checksum divergente.
- Teste de registro de processo.

---

## 13.7 Entrega E7 — Qualidade, segurança e distribuição

### Objetivo

Consolidar a qualidade técnica do sistema, ampliar testes, preparar CI/CD, configurar distribuição multiplataforma e garantir segurança dos artefatos.

### Tarefas principais

| ID | Tarefa |
|---|---|
| E7-T01 | Revisar requisitos e rastreabilidade com implementação. |
| E7-T02 | Executar testes unitários Go com `go test ./...`. |
| E7-T03 | Executar testes Java com Maven ou script alternativo. |
| E7-T04 | Executar testes de integração local. |
| E7-T05 | Executar testes de integração HTTP. |
| E7-T06 | Executar testes dos comandos `simulador`. |
| E7-T07 | Configurar build multiplataforma para Windows, Linux e macOS. |
| E7-T08 | Configurar geração de checksums SHA256. |
| E7-T09 | Configurar assinatura de artefatos com Cosign. |
| E7-T10 | Configurar publicação em GitHub Releases. |
| E7-T11 | Atualizar Documento de Implementação e Integração. |
| E7-T12 | Atualizar Documento de Teste de Software. |

### Artefatos esperados

- Relatório de testes atualizado.
- Logs de testes Go e Java.
- Workflows `build.yml` e `release.yml`.
- Binários multiplataforma.
- Arquivo `checksums.txt`.
- Arquivos `.sig` e `.pem`, quando release for executada.
- Documento de Implementação e Integração atualizado.
- Documento de Teste atualizado.
- Evidências de execução.

### Critérios de aceite

- Testes automatizados principais executam sem falhas críticas.
- Builds são gerados para as plataformas previstas.
- Artefatos possuem checksums.
- Processo de release está documentado.
- Segurança de artefatos está prevista por Cosign.
- Não existem defeitos críticos abertos nos fluxos principais.

### Atividades de qualidade

- Revisão de código.
- Reexecução de testes após correções.
- Registro de defeitos.
- Validação de mensagens de erro.
- Verificação de integridade dos artefatos.

---

## 13.8 Entrega E8 — Implantação, validação final e encerramento

### Objetivo

Executar a validação final, preparar o pacote de implantação, registrar evidências, revisar documentação e formalizar o encerramento do projeto.

### Tarefas principais

| ID | Tarefa |
|---|---|
| E8-T01 | Preparar pacote de implantação. |
| E8-T02 | Executar procedimento de instalação local. |
| E8-T03 | Executar procedimento de execução dos CLIs. |
| E8-T04 | Validar fluxo local de assinatura e validação. |
| E8-T05 | Validar fluxo HTTP de assinatura e validação. |
| E8-T06 | Validar comandos do simulador, quando `simulador.jar` estiver disponível. |
| E8-T07 | Executar plano de testes de implantação. |
| E8-T08 | Registrar prints, logs e arquivos gerados. |
| E8-T09 | Validar plano de backup e rollback. |
| E8-T10 | Atualizar Documento de Implantação. |
| E8-T11 | Consolidar relatório final do projeto. |
| E8-T12 | Apresentar entrega final e registrar aceite. |

### Artefatos esperados

- Pacote final do projeto.
- Documento de Implantação de Software.
- Evidências de implantação.
- Evidências de teste.
- Relatório final.
- Lista final de limitações e melhorias futuras.
- Registro de aceite.
- Repositório organizado e versionado.

### Critérios de aceite

- Sistema pode ser instalado e executado no ambiente definido.
- Comandos principais funcionam conforme documentação.
- Testes finais foram executados ou justificados.
- Evidências foram registradas.
- Documentação está coerente com a implementação.
- Riscos e limitações finais foram declarados.
- Entrega final foi revisada e aprovada.

### Atividades de qualidade

- Checklist final de entrega.
- Revisão de documentação.
- Validação de rastreabilidade.
- Revisão dos critérios de aceite.
- Registro de lições aprendidas.

---

## 14. Matriz de artefatos por entrega

| Artefato | E1 | E2 | E3 | E4 | E5 | E6 | E7 | E8 |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| Especificação de Requisitos | X |  |  |  |  |  | Revisão | Revisão final |
| Documento de Arquitetura | X | Revisão |  |  | Revisão |  | Revisão | Revisão final |
| Documento de Projeto Detalhado |  | X | Revisão | Revisão | Revisão | Revisão | Revisão | Revisão final |
| Modelo C4 | X | Revisão |  |  | Revisão |  | Revisão | Revisão final |
| Código CLI `assinatura` |  | Base |  | X | X |  | Revisão | Final |
| Código `assinador.jar` |  | Base | X | Ajuste | X |  | Revisão | Final |
| Código CLI `simulador` |  | Base |  |  |  | X | Revisão | Final |
| Testes Go |  | Base |  | X | X | X | X | Final |
| Testes Java |  | Base | X | Ajuste | X |  | X | Final |
| Documento de Implementação e Integração |  |  |  |  |  |  | X | Revisão final |
| Documento de Teste |  |  |  |  |  |  | X | Revisão final |
| Documento de Implantação |  |  |  |  |  |  |  | X |
| README |  | X | Revisão | Revisão | Revisão | Revisão | Revisão | Final |
| Workflows CI/CD |  | Base |  |  |  |  | X | Final |
| Evidências |  |  | X | X | X | X | X | X |
| Pacote final |  |  |  |  |  |  |  | X |

---

## 15. Rastreabilidade entre entregas, requisitos e histórias de usuário

| Entrega | Requisitos / funcionalidades relacionadas | Histórias de usuário relacionadas |
|---|---|---|
| **E1** | Definição do escopo, requisitos, arquitetura e documentação inicial | Todas as histórias em nível de planejamento |
| **E2** | Estrutura base do CLI, versionamento, CI inicial, projeto Java | US-01.1, US-05.1 |
| **E3** | Simulação de assinatura, validação e lógica Java | US-02.1, US-02.2, US-02.3, US-02.5 |
| **E4** | Comandos `sign` e `validate`, invocação local e saída legível | US-01.2, US-01.3, US-01.4 |
| **E5** | Modo servidor, HTTP, start/stop, health check e processos | US-01.5, US-01.6, US-01.7, US-01.8, US-01.9, US-02.4 |
| **E6** | CLI `simulador`, start/status/stop, download e cache | US-03.1, US-03.2, US-03.3, US-03.4 |
| **E7** | Provisionamento, releases, checksums, Cosign, testes e CI/CD | US-04.1, US-05.1, US-05.2, US-05.3 |
| **E8** | Implantação, validação final, evidências e aceite | Todas as histórias com foco em validação final |

---

## 16. Plano de qualidade

### 16.1 Objetivos de qualidade

| Objetivo | Descrição |
|---|---|
| **Correção funcional** | O sistema deve executar os comandos principais de acordo com os requisitos. |
| **Integração confiável** | O CLI deve se comunicar corretamente com o `assinador.jar` em modo local e HTTP. |
| **Usabilidade por terminal** | Mensagens e comandos devem ser claros para o usuário. |
| **Manutenibilidade** | O código deve ser modular, organizado e documentado. |
| **Testabilidade** | Componentes devem ser testáveis isoladamente e em integração. |
| **Portabilidade** | Binários devem ser geráveis para Windows, Linux e macOS. |
| **Segurança de artefatos** | Releases devem prever SHA256 e Cosign. |
| **Rastreabilidade** | Requisitos, entregas, testes e artefatos devem estar relacionados. |

### 16.2 Práticas de qualidade adotadas

- Revisão de requisitos antes da implementação.
- Modularização por responsabilidade.
- Separação entre CLI, casos de uso, infraestrutura e aplicação Java.
- Validação de parâmetros no CLI e no Java.
- Tratamento padronizado de erros.
- Testes unitários e de integração.
- Registro de evidências.
- Controle de versão por Git.
- Uso de CI/CD.
- Documentação incremental.

### 16.3 Critérios de qualidade por entrega

| Entrega | Critérios mínimos de qualidade |
|---|---|
| E1 | Documentos coerentes e revisados. |
| E2 | Estrutura compila e comando de versão funciona. |
| E3 | Testes Java passam e validação rejeita dados inválidos. |
| E4 | Fluxo local funciona e mensagens são claras. |
| E5 | Endpoints HTTP funcionam e processos são registrados. |
| E6 | CLI `simulador` controla processo e trata erros básicos. |
| E7 | Testes e CI/CD executam sem falhas críticas. |
| E8 | Implantação e validação final possuem evidências. |

---

## 17. Plano de testes por entrega

| Entrega | Testes previstos |
|---|---|
| E1 | Revisão documental, checklist de escopo, validação de rastreabilidade inicial. |
| E2 | Teste de compilação Go, teste de estrutura Java, teste do comando `version`. |
| E3 | Testes unitários Java de `ParameterValidator` e `FakeSignatureService`. |
| E4 | Testes funcionais de `assinatura sign --local` e `assinatura validate --local`; testes de erro de parâmetro. |
| E5 | Testes HTTP de `/health`, `/sign`, `/validate` e `/shutdown`; testes de start/stop; teste de porta ocupada. |
| E6 | Testes de `simulador start`, `simulador status`, `simulador stop`; testes de JAR ausente, URL inválida e checksum. |
| E7 | Testes automatizados Go e Java; testes de integração; testes de build; testes de release e checksum. |
| E8 | Testes de implantação; testes de aceitação; testes de rollback; validação de evidências. |

---

## 18. Plano de comunicação

### 18.1 Reuniões e acompanhamento

| Comunicação | Frequência | Participantes | Objetivo |
|---|---|---|---|
| Reunião de planejamento | Início do projeto | Equipe completa | Alinhar escopo, entregas e responsabilidades. |
| Reunião de acompanhamento | Semanal ou por entrega | Equipe completa | Revisar progresso, impedimentos e riscos. |
| Revisão técnica | Ao final de cada entrega | Desenvolvedores e avaliador interno | Verificar qualidade técnica e aderência aos requisitos. |
| Revisão de documentação | Ao final de entregas documentais | Responsável por documentação e equipe | Garantir coerência entre documentos. |
| Demonstração | Ao final de entregas E4, E5, E6 e E8 | Equipe e avaliador | Demonstrar fluxos implementados. |

### 18.2 Canais de comunicação

- Repositório Git e issues.
- Quadro de tarefas, como Kanban ou lista de backlog.
- Grupo de mensagens da equipe.
- Reuniões presenciais ou online.
- Documentos Markdown no repositório.

### 18.3 Registro de decisões

As decisões importantes devem ser registradas em um arquivo de decisões técnicas ou na documentação do projeto, incluindo:

- data da decisão;
- problema analisado;
- alternativas consideradas;
- decisão tomada;
- justificativa;
- impacto no projeto.

---

## 19. Gestão de riscos

| ID | Risco | Probabilidade | Impacto | Plano de mitigação | Plano de contingência |
|---|---|---|---|---|---|
| R-001 | Falta de Java 21 no ambiente | Média | Alto | Documentar pré-requisito e preparar resolução de Java. | Usar Java instalado manualmente ou diretório `~/.hubsaude/`. |
| R-002 | Maven indisponível | Média | Médio | Disponibilizar script alternativo com `javac`/`jar`. | Compilar Java por script local. |
| R-003 | Falha na integração Go → Java | Média | Alto | Testar invocação local cedo na E4. | Simplificar argumentos e revisar captura de saída. |
| R-004 | Endpoints HTTP inconsistentes | Média | Alto | Definir contrato HTTP e testar na E5. | Ajustar cliente e servidor com DTOs estáveis. |
| R-005 | Porta 8080 ou 8443 ocupada | Alta | Médio | Implementar verificação de porta e permitir `--port`. | Usar porta alternativa. |
| R-006 | `simulador.jar` real indisponível | Alta | Médio | Permitir `--jar` e `--source`; documentar limitação. | Usar JAR de teste ou simulação controlada. |
| R-007 | Baixa cobertura de testes | Média | Alto | Planejar testes por entrega. | Priorizar testes dos fluxos críticos. |
| R-008 | Inconsistência entre documentos e código | Média | Médio | Atualizar documentação incrementalmente. | Fazer revisão final de rastreabilidade. |
| R-009 | Erros em CI/CD | Média | Médio | Validar workflows antes da entrega final. | Executar builds localmente e registrar evidências. |
| R-010 | Confusão entre assinatura simulada e real | Média | Alto | Documentar claramente o caráter simulado. | Adicionar aviso em README e mensagens. |
| R-011 | Falha na verificação SHA256/Cosign | Baixa | Médio | Automatizar verificação no release. | Publicar checksums e documentar verificação manual. |
| R-012 | Atraso no desenvolvimento | Média | Alto | Priorizar núcleo funcional antes de melhorias. | Reduzir escopo complementar mantendo entregas essenciais. |

---

## 20. Gestão de mudanças

### 20.1 Processo de solicitação de mudança

1. Registrar a mudança solicitada.
2. Identificar motivo e impacto.
3. Avaliar impacto em requisitos, arquitetura, código, testes e cronograma.
4. Classificar prioridade.
5. Aprovar, rejeitar ou adiar a mudança.
6. Atualizar backlog e documentação.
7. Implementar e testar.
8. Registrar evidências.

### 20.2 Classificação de mudanças

| Tipo de mudança | Exemplo | Tratamento |
|---|---|---|
| Correção | Ajuste de erro no comando `sign`. | Pode ser priorizada na entrega atual. |
| Melhoria | Melhorar mensagem de erro. | Entra no backlog se não bloquear aceite. |
| Nova funcionalidade | Adicionar comando `doctor`. | Avaliar para evolução futura. |
| Mudança de escopo | Implementar assinatura real. | Requer replanejamento, pois está fora do escopo atual. |
| Mudança técnica | Trocar biblioteca CLI. | Requer análise de impacto arquitetural. |

---

## 21. Gestão de configuração e versionamento

### 21.1 Repositório

O projeto deve ser controlado por Git, com a seguinte organização recomendada:

```text
main
feature/<nome-da-funcionalidade>
fix/<nome-da-correcao>
release/<versao>
```

### 21.2 Versionamento semântico

As versões devem seguir SemVer:

```text
MAJOR.MINOR.PATCH
```

Exemplos:

```text
v0.1.0
v0.2.0
v1.0.0
```

### 21.3 Convenção de artefatos

```text
assinatura-<versao>-<os>-<arch>
simulador-<versao>-<os>-<arch>
checksums.txt
<artefato>.sig
<artefato>.pem
```

### 21.4 Itens sob controle de configuração

- Código-fonte Go.
- Código-fonte Java.
- Scripts.
- Workflows CI/CD.
- Documentação Markdown.
- Testes.
- Arquivos de configuração.
- Artefatos de release.
- Evidências finais.

---

## 22. Indicadores de acompanhamento

| Indicador | Descrição | Forma de cálculo / verificação |
|---|---|---|
| Percentual de entregas concluídas | Mede avanço do cronograma. | Entregas concluídas / 8 × 100. |
| Percentual de tarefas concluídas | Mede avanço dentro das entregas. | Tarefas concluídas / tarefas planejadas × 100. |
| Testes aprovados | Mede qualidade funcional. | Testes aprovados / testes executados × 100. |
| Defeitos abertos | Mede pendências técnicas. | Quantidade de defeitos ainda não corrigidos. |
| Defeitos críticos | Mede risco para entrega. | Quantidade de defeitos críticos abertos. |
| Cobertura de requisitos | Mede rastreabilidade. | Requisitos com teste / requisitos testáveis × 100. |
| Builds bem-sucedidos | Mede estabilidade técnica. | Builds aprovados / builds executados × 100. |
| Artefatos entregues | Mede completude documental e técnica. | Artefatos entregues / artefatos planejados × 100. |

---

## 23. Critérios de aceite do projeto

O projeto será considerado aceito quando:

- o CLI `assinatura` estiver funcional para comandos principais;
- o `assinador.jar` validar parâmetros e simular assinatura/validação;
- a integração local via `java -jar` estiver demonstrada;
- a integração HTTP estiver demonstrada;
- os endpoints principais estiverem documentados;
- o CLI `simulador` estiver implementado conforme escopo possível;
- o sistema registrar processos em `~/.hubsaude/`;
- os testes principais estiverem documentados e, quando possível, executados;
- os documentos técnicos estiverem coerentes entre si;
- a implantação local estiver documentada;
- as limitações conhecidas estiverem explicitadas;
- existirem evidências de implementação, integração, teste e implantação;
- não existirem defeitos críticos abertos nos fluxos principais;
- o pacote final estiver organizado e versionado.

---

## 24. Plano de encerramento

### 24.1 Atividades de encerramento

| ID | Atividade |
|---|---|
| ENC-001 | Conferir se as 8 entregas foram concluídas ou justificadas. |
| ENC-002 | Revisar documentação final. |
| ENC-003 | Executar checklist de testes finais. |
| ENC-004 | Registrar evidências de execução. |
| ENC-005 | Consolidar lista de limitações e melhorias futuras. |
| ENC-006 | Validar repositório e artefatos finais. |
| ENC-007 | Preparar apresentação ou demonstração final. |
| ENC-008 | Registrar aceite ou pendências finais. |

### 24.2 Checklist final

| Item | Status |
|---|---|
| Requisitos documentados | A verificar |
| Arquitetura documentada | A verificar |
| Projeto detalhado documentado | A verificar |
| Modelo C4 documentado | A verificar |
| Implementação disponível | A verificar |
| Testes documentados | A verificar |
| Implantação documentada | A verificar |
| Cronograma e gerenciamento documentados | A verificar |
| Evidências anexadas | A verificar |
| Limitações declaradas | A verificar |
| Aceite registrado | A verificar |

---

## 25. Referências

- Especificação original do Sistema Runner — Trabalho Prático.
- Plano Revisado #2 do Sistema Runner.
- Documento de Design do Sistema Runner baseado no Modelo C4.
- Especificação de Requisitos de Software — Sistema Runner.
- Documento de Arquitetura de Software — Sistema Runner.
- Documento de Projeto Detalhado de Software — Sistema Runner.
- Documento de Modelo C4 de Software — Sistema Runner.
- Documento de Implementação e Integração de Software — Sistema Runner.
- Documento de Teste de Software — Sistema Runner.
- Documento de Implantação de Software — Sistema Runner.
- README da implementação de referência acadêmica do Sistema Runner.
- Boas práticas de Engenharia de Software para gerenciamento de projetos, desenvolvimento iterativo e incremental, rastreabilidade, testes, integração, implantação e controle de qualidade.
