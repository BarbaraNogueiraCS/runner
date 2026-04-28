# Entrega da Sprint 1

## Itens implementados

- Estrutura base do projeto em Go.
- CLI `assinatura` com comando `version`.
- Stub da CLI `simulador`.
- Estrutura de diretórios prevista para evolução do sistema.
- Teste de aceitação do comando `assinatura version`.
- Workflow de build multiplataforma.
- Workflow de release com SemVer, checksums e Cosign.

## Definição de pronto atendida

- O projeto contém código-fonte compilável para as CLIs da Sprint 1.
- O comando `assinatura version` exibe `dev` localmente.
- O workflow de build está configurado para branch `main`.
- O workflow de release está configurado para tags `v*`.
- A release gera binários para Linux, Windows e macOS.
- A release gera `checksums.txt`.
- A release assina os binários com Cosign.
