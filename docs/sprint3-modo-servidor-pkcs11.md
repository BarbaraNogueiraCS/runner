# Sprint 3 — Modo Servidor e Material Criptográfico

Esta sprint adiciona ao Sistema Runner o modo servidor do `assinador.jar`, a invocação HTTP pelo CLI e o suporte operacional a material criptográfico via PKCS#11.

## US-02.4 — Endpoints HTTP do assinador.jar

Implementado em:

- `assinador/src/br/ufg/hubsaude/assinador/SignatureController.java`
- `assinador/src/br/ufg/hubsaude/assinador/Main.java`
- `assinador/src/br/ufg/hubsaude/assinador/FakeSignatureService.java`

Endpoints disponíveis:

```text
GET  /health
POST /sign
POST /validate
POST /shutdown
```

O `SignatureController` recebe as requisições HTTP, converte o JSON em `SignRequest` ou `ValidateRequest` e delega para `SignatureService`. Assim, o modo servidor reutiliza a mesma lógica do modo CLI local.

## US-02.5 — Integração com dispositivo criptográfico via PKCS#11

Implementado em:

- `assinador/src/br/ufg/hubsaude/assinador/Pkcs11Support.java`
- `assinador/src/br/ufg/hubsaude/assinador/GuideSignatureService.java`

O assinador aceita material criptográfico com `type` igual a `TOKEN` ou `SMARTCARD`. Quando `simulation=true`, o fluxo simula o token. Quando a simulação não está habilitada, o assinador exige:

- `pin`
- `pkcs11Library` ou `library`
- opcionalmente `identifier`, `alias`, `slot` ou `tokenLabel`

Exemplo simulado:

```json
{"type":"TOKEN","pin":"123456","identifier":"runner","simulation":"true"}
```

Exemplo real com PKCS#11:

```json
{"type":"TOKEN","pin":"123456","identifier":"meu-certificado","pkcs11Library":"/usr/lib/softhsm/libsofthsm2.so"}
```

Se a biblioteca não existir, a mensagem é clara:

```text
PKCS11.LIBRARY-NOT-FOUND
```

Se o PIN estiver ausente:

```text
CRYPTO.PIN-REQUIRED
```

## US-01.5 — Iniciar assinador.jar no modo servidor

Implementado em:

- `cmd/assinatura/main.go`
- `internal/assinador/client.go`
- `internal/process/state.go`

Comandos:

```bash
assinatura start --jar assinador/target/assinador.jar --port 8080 --timeout 10
assinatura server start --jar assinador/target/assinador.jar --port 8080 --timeout 10
```

O estado é salvo em:

```text
~/.hubsaude/processos/assinador.json
```

Esse arquivo registra nome, PID, porta, URL, caminho do JAR e data de início.

## US-01.6 — Invocar assinador.jar via HTTP

Por padrão, `assinatura sign` e `assinatura validate` tentam usar o servidor HTTP em execução. Se não houver servidor, o CLI tenta iniciar o servidor quando `--jar` ou `RUNNER_ASSINADOR_JAR` estiver disponível. Se o servidor não puder ser iniciado, o CLI tenta fallback local com `java -jar`.

Para forçar modo local:

```bash
assinatura sign --local --jar assinador/target/assinador.jar ...
assinatura validate --local --jar assinador/target/assinador.jar ...
```

## US-01.7 — Detectar instância em execução

O CLI usa:

- health check HTTP em `/health`
- estado em `~/.hubsaude/processos/assinador.json`

Se o health check responder, a instância é reutilizada. Se o processo registrado não responder, ele é tratado como inativo.

## US-01.8 — Interromper execução

Comandos:

```bash
assinatura stop --port 8080
assinatura server stop --port 8080
```

O CLI chama `POST /shutdown` e remove o registro local em `~/.hubsaude/processos/assinador.json`.

## US-01.9 — Timeout por inatividade

O parâmetro documentado é:

```bash
--timeout <minutos>
```

Também existe o alias de compatibilidade:

```bash
--idle-timeout-minutes <minutos>
```

Exemplo:

```bash
assinatura start --jar assinador/target/assinador.jar --port 8080 --timeout 10
```

Após o período configurado sem requisições, o servidor encerra automaticamente.

## Testes de aceitação

O alvo abaixo cobre CLI local, HTTP, PKCS#11 simulado, erro de biblioteca PKCS#11 ausente, criação, validação e endpoints HTTP:

```bash
make java-test
```

Os testes Go cobrem parsing, aliases de servidor, timeout, cliente HTTP, reuso de instância saudável e mensagens de erro.

```bash
go test ./...
```
