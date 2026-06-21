# Projeto Java do Assinador

Este diretório existe para atender à organização pedida na Sprint 2: `projetos/assinador-java`.

A implementação Java canônica do `assinador.jar` fica em:

```text
assinador/
```

O `Makefile` desta pasta é um wrapper que encaminha os comandos para `../../assinador`.

## Comandos

```bash
make -C projetos/assinador-java clean
make -C projetos/assinador-java all
make -C projetos/assinador-java test
make -C projetos/assinador-java samples
```
