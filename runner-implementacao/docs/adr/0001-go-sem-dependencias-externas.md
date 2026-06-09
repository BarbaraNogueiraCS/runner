# ADR 0001 - CLIs em Go sem dependências externas

## Status

Aceito.

## Contexto

O projeto exige CLIs multiplataforma para Windows, Linux e macOS. O plano menciona o uso de Cobra, mas a implementação de referência foi preparada em ambiente sem download de dependências externas.

## Decisão

Os comandos foram implementados com a biblioteca padrão de Go. Isso reduz dependências, facilita build reprodutível e mantém a portabilidade.

## Consequências

A interface de linha de comando é mais simples do que uma implementação com Cobra, mas atende aos comandos principais e pode ser evoluída futuramente.
