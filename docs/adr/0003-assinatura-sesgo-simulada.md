# ADR 0003 — Estrutura de assinatura conforme guia SES-GO em modo simulado

## Status

Aceita.

## Contexto

Os guias da SES-GO para criação e validação de assinatura digital definem um fluxo baseado em Bundle FHIR R4, Provenance FHIR R4, cadeia de certificados, política versionada, timestamp de referência, JWS JSON Serialization e encapsulamento em `Signature.data`.

A especificação acadêmica do Sistema Runner, por outro lado, delimita que a implementação deve simular a criação e a validação de assinatura digital, sem exigir assinatura criptográfica real, validação real com ACs, OCSP, CRL, TSA ou PKCS#11 real.

## Decisão

A implementação do `assinador.jar` passa a seguir a estrutura de entrada e saída dos guias SES-GO, mas mantém a operação criptográfica em modo simulado.

O serviço Java agora:

- recebe Bundle, Provenance, material criptográfico, cadeia de certificados, timestamp, estratégia e política;
- valida parâmetros estruturais;
- gera JWS JSON Serialization;
- usa `protected header` com `alg`, `x5c`, `sigPId` e `iat`;
- usa payload como SHA-256 em base64Url;
- encapsula JWS em FHIR `Signature.data` com base64 padrão;
- valida a estrutura de `Signature.data` e JWS;
- verifica política, timestamp, assinatura simulada e integridade opcional do conteúdo.

## Consequências

A solução fica mais aderente aos guias de interoperabilidade da SES-GO sem ultrapassar o escopo acadêmico do Runner.

As validações criptográficas reais continuam explicitamente fora do escopo desta versão.
