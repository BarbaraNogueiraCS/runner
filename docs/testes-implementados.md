# Testes implementados no Sistema Runner

Este documento registra a rastreabilidade entre os tipos de teste exigidos e os arquivos de teste presentes no código.

## 1. Testes unitários

Implementados nos pacotes internos do Go:

- `internal/apperrors/exitcode_test.go`: códigos de saída estáveis.
- `internal/httpx/client_test.go`: cliente HTTP, GET, POST JSON e POST vazio.
- `internal/jdk/jdk_test.go` e `internal/jdk/acceptance_test.go`: localização e validação de versão do Java.
- `internal/paths/paths_test.go`: criação dos diretórios gerenciados em `~/.hubsaude` ou `RUNNER_HOME`.
- `internal/process/state_test.go`: persistência, leitura e remoção de estado de processo.
- `internal/netutil/port_test.go`: detecção de porta ocupada/alcançável.
- `internal/release/manifest_test.go`: leitura de manifest, download e reuso de artefato.

## 2. Testes de integração

Implementados com servidores HTTP reais em memória, via `httptest`:

- `internal/assinador/client_test.go`: integração HTTP com `/sign` e reutilização de instância saudável via `/health`.
- `internal/assinador/acceptance_test.go`: integração HTTP com `/validate`.
- `internal/simulador/client_test.go`: integração com `/api/info`, `/shutdown` e reuso de simulador ativo.
- `internal/simulador/acceptance_test.go`: comportamento quando a porta está ocupada por serviço que não é o simulador.
- `assinador/Makefile`: fluxo real `java -jar assinador.jar sign` e `java -jar assinador.jar validate`.

## 3. Casos de teste para cenários de erro

Implementados em Go e Java:

- Comando desconhecido nos CLIs: `cmd/assinatura/main_test.go` e `cmd/simulador/main_test.go`.
- Parâmetros obrigatórios ausentes: `cmd/assinatura/acceptance_test.go`.
- Porta ocupada sem health check válido: `internal/assinador/client_test.go`.
- Porta ocupada por serviço que não é o Simulador HubSaúde: `internal/simulador/acceptance_test.go`.
- Checksum inválido no download de artefato: `internal/release/acceptance_test.go`.
- Java com saída de versão inválida: `internal/jdk/acceptance_test.go`.
- `assinador.jar` com PIN ausente para TOKEN/SMARTCARD, política inválida, Bundle ausente, Provenance inconsistente e arquivo que não contém `Signature.data`: `assinador/Makefile`.

## 4. Testes de aceitação baseados nos critérios definidos

- CLI `assinatura` aceita `sign`, `validate`, `server start`, `server status`, `server stop`: `cmd/assinatura/main_test.go` e `cmd/assinatura/acceptance_test.go`.
- Modo servidor como padrão, porta padrão 8080 e integração HTTP: `internal/assinador/acceptance_test.go`.
- Invocação local falha bem quando o JAR não é informado: `internal/assinador/acceptance_test.go`.
- Gerenciador do Simulador consulta `/api/info`, para via `/shutdown`, reutiliza instância ativa e rejeita porta ocupada não identificada como simulador: `internal/simulador/client_test.go` e `internal/simulador/acceptance_test.go`.
- Download dinâmico e reuso de versão local do artefato: `internal/release/manifest_test.go`.
- Garantia de que `target/`, `out/`, `dist/` e JSONs gerados não sejam versionados: `scripts/check-generated-files.sh` e `internal/release/artifacts_acceptance_test.go`.

## Como executar

```bash
go test ./...
go vet ./...
./scripts/check-generated-files.sh
./scripts/check-release-artifacts.sh

cd assinador
make clean all test
cd ..
```

## Sprint 4 — Simulador do HubSaúde e Segurança de Artefatos

- `cmd/simulador/main_test.go`: valida `simulador version`, parsing de flags, montagem de URL com `--port`, override de porta e configuração de `--source`/`--sha256`.
- `internal/simulador/client_test.go`: valida consulta `/api/info`, parada por `/shutdown` e reuso de instância ativa.
- `internal/simulador/acceptance_test.go`: valida erro quando a porta está ocupada por serviço que não é o Simulador, seleção de artefato via ambiente e configuração de URL direta por `RUNNER_SIMULADOR_SOURCE`.
- `internal/release/manifest_test.go`: valida download, cache, seleção por `release.json`, fallback `simulador -> validador` e URL direta alternativa por `EnsureArtifactFromURL` com SHA-256.
- `internal/release/acceptance_test.go`: valida rejeição de checksum SHA-256 inválido.
- `scripts/check-release-artifacts.sh`: valida rastreabilidade da Sprint 4, incluindo binários `simulador-*`, `checksums.txt`, Cosign, `--source`, `--sha256` e ausência da pasta intermediária antiga.
