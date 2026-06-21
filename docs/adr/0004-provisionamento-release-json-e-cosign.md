# ADR 0004 — Provisionamento por release.json e integridade de artefatos

## Decisão

O Runner usa `https://raw.githubusercontent.com/BarbaraNogueiraCS/runner/main/release.json` como manifesto remoto padrão.

O manifesto aponta para artefatos reutilizáveis, como JDK/JRE gerenciado e componentes externos necessários ao simulador. A integridade dos downloads é validada com SHA256 quando disponível.

As releases oficiais do Runner publicam binários multiplataforma, `checksums.txt` e assinaturas Cosign com OIDC e transparency log.
