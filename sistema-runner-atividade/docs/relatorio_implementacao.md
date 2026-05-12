# Relatório de Implementação — Sistema Runner

## 1. Objetivo

Este relatório descreve a implementação de referência do Sistema Runner, alinhada aos documentos de requisitos, arquitetura, projeto detalhado e Modelo C4.

## 2. O que foi desenvolvido

Foram desenvolvidos:

- CLI `assinatura`, em Go, para operações de assinatura simulada;
- CLI `simulador`, em Go, para gerenciamento do Simulador do HubSaúde;
- aplicação Java `assinador.jar`, compatível com Java 21;
- endpoints HTTP `/sign`, `/validate`, `/health` e `/shutdown` no assinador;
- persistência local de metadados de processo em `~/.hubsaude/processos/`;
- estrutura de download e verificação SHA256 de artefatos;
- workflows de CI/CD para build, teste e release;
- testes unitários iniciais.

## 3. Requisitos atendidos

| Grupo | Atendimento nesta versão |
|---|---|
| US-01 | Implementação do CLI `assinatura`, comandos principais, invocação local e HTTP, start/stop e version. |
| US-02 | Implementação do `assinador.jar`, validação de parâmetros, simulação de assinatura e validação, endpoints HTTP. |
| US-03 | Implementação do CLI `simulador` com start, stop e status, controle de porta e registro de processo. |
| US-04 | Estrutura de resolução de Java 21, priorizando Java instalado ou diretório gerenciado. |
| US-05 | Estrutura de workflows para geração de binários, checksums e assinatura de artefatos. |

## 4. Organização técnica

A solução foi dividida em camadas e módulos:

- `cmd/`: pontos de entrada dos binários;
- `internal/cli/`: comandos e interface com o usuário;
- `internal/assinatura/`: casos de uso e DTOs do fluxo de assinatura;
- `internal/simulador/`: casos de uso e DTOs do fluxo do simulador;
- `internal/java/`: resolução de Java e invocação local/HTTP;
- `internal/process/`: controle de porta, health check e registro de processos;
- `internal/artifacts/`: download e checksum;
- `projetos/assinador-java/`: aplicação Java do assinador.

## 5. Decisões principais

- Uso de Go para CLIs, favorecendo portabilidade.
- Uso de Java 21 para o `assinador.jar`.
- Uso de HTTP para o modo servidor.
- Uso de arquivos JSON locais para registro de processos.
- Uso de respostas JSON no `assinador.jar` para facilitar integração.
- Uso de SHA256 e Cosign nos workflows de release.

## 6. Limitações

- O provisionamento automático completo de JDK/JRE está estruturado, mas não realiza extração automática de distribuições baixadas nesta versão.
- A assinatura digital é simulada.
- A integração PKCS#11 está preparada como adaptador, mas não executa operação criptográfica real.
- O simulador real precisa ser fornecido via `--jar` ou `--source`.

## 7. Como evoluir

- Implementar extração automática de JDK/JRE por plataforma.
- Ampliar os testes end-to-end.
- Integrar SoftHSM2 para testes PKCS#11.
- Implementar logs estruturados.
- Criar comando `doctor` para diagnóstico do ambiente.
