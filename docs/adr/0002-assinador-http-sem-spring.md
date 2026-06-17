# ADR 0002 - assinador.jar com HTTP nativo do Java

## Status

Aceito.

## Contexto

O `assinador.jar` precisa expor endpoints HTTP, porém o ambiente de entrega não possui Maven nem acesso garantido a repositórios externos.

## Decisão

Foi usado `com.sun.net.httpserver.HttpServer`, disponível no JDK, evitando dependências externas.

## Consequências

A aplicação não depende de Spring Boot nesta versão. A estrutura é suficiente para simulação, testes e contrato CLI/JAR, mas pode ser migrada posteriormente para Spring Boot.
