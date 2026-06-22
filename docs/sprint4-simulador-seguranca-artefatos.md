# Sprint 4 — Simulador do HubSaúde e Segurança de Artefatos

## Objetivo

Entregar a gestão completa do Simulador do HubSaúde pelo CLI `simulador` e manter a distribuição dos artefatos com verificação de integridade por checksums SHA-256 e assinatura criptográfica com Cosign/Sigstore.

## US-03.1 — Iniciar o Simulador via CLI

Status: implementado.

Artefatos principais:

- `cmd/simulador/main.go`
- `internal/simulador/client.go`
- `internal/netutil/port.go`
- `internal/process/state.go`
- `internal/release/manifest.go`

O comando `simulador start` inicia o `simulador.jar` sem exigir que o usuário saiba o comando Java. Antes de iniciar, o CLI consulta `/api/info`; se a instância já estiver ativa, ela é reaproveitada. Se a porta estiver ocupada por outro serviço que não responde como Simulador HubSaúde, o comando falha com mensagem clara.

Exemplos:

```bash
simulador start --port 8443 --insecure
simulador start --url https://localhost:9443 --insecure
simulador start --jar /caminho/simulador.jar --port 8443 --insecure
```

## US-03.2 — Parar e monitorar o Simulador

Status: implementado.

Comandos disponíveis:

```bash
simulador status --port 8443 --insecure
simulador info --port 8443 --insecure
simulador stop --port 8443 --insecure
```

O processo é registrado em `~/.hubsaude/processos/simulador.json`, com dados como PID, porta, URL e caminho do JAR. O encerramento usa o endpoint `/shutdown` e remove o registro local após a solicitação.

## US-03.3 — Estrutura base do CLI `simulador` em Go

Status: implementado.

O CLI `simulador` segue a mesma organização do CLI `assinatura`:

- pacote principal em `cmd/simulador`;
- lógica de integração em `internal/simulador`;
- registro de processo em `internal/process`;
- verificação de porta em `internal/netutil`;
- provisionamento de JRE/JDK e artefatos em `internal/jdk` e `internal/release`.

Os workflows `.github/workflows/build.yml` e `.github/workflows/release.yml` geram binários multiplataforma do CLI `simulador` junto com o CLI `assinatura`:

```text
simulador-<versão>-linux-amd64.AppImage
simulador-<versão>-windows-amd64.exe
simulador-<versão>-macos-amd64.dmg
```

## US-03.4 — Obter `simulador.jar` dinamicamente

Status: implementado.

Quando o usuário não informa `--jar`, o CLI usa `internal/release/manifest.go` para consultar o `release.json`, selecionar o artefato `simulador` e baixar o JAR do GitHub Releases. O manifesto contém URL, versão, tag e checksum SHA-256 do artefato.

O cache local é mantido em `~/.hubsaude/simulador/`. Se o JAR já tiver sido baixado e o checksum continuar válido, o download não é repetido.

Exemplo usando o manifesto padrão:

```bash
simulador start --port 8443 --insecure
```

Exemplo usando um `release.json` alternativo:

```bash
simulador start --release-json http://127.0.0.1:8099/release.json --port 8443 --insecure
```

Exemplo usando uma URL direta alternativa:

```bash
simulador start \
  --source https://example.org/simulador.jar \
  --sha256 <sha256-esperado> \
  --port 8443 \
  --insecure
```

A opção `--checksum` também é aceita como alias de `--sha256`.

## Segurança dos artefatos

Status: implementado.

A Sprint 4 reaproveita o fluxo de segurança da US-05.3:

- `release.yml` gera `checksums.txt` com SHA-256 dos artefatos distribuídos;
- todos os artefatos principais são assinados com `cosign sign-blob`;
- a assinatura usa identidade OIDC e transparency log;
- cada artefato publicado possui `.sig`, `.pem` e `.bundle`;
- a documentação de verificação fica em `docs/integridade-assinatura-artefatos.md` e no script `scripts/verify-release-artifact.sh`.

Exemplo de verificação local de um artefato baixado da release:

```bash
sha256sum -c checksums.txt

cosign verify-blob \
  --certificate simulador-1.0.6-linux-amd64.AppImage.pem \
  --signature simulador-1.0.6-linux-amd64.AppImage.sig \
  simulador-1.0.6-linux-amd64.AppImage
```

## Evidências de teste

- `cmd/simulador/main_test.go`: version, parsing de flags, `--port`, `--source`, `--sha256`.
- `internal/simulador/client_test.go`: `/api/info`, `/shutdown`, reuso de instância ativa.
- `internal/simulador/acceptance_test.go`: porta ocupada por serviço inválido, seleção de artefato e source por variável de ambiente.
- `internal/release/manifest_test.go`: download, cache, checksum e URL direta alternativa.
- `internal/release/acceptance_test.go`: rejeição de checksum inválido.
- `scripts/check-release-artifacts.sh`: valida presença da Sprint 4, workflow de release, checksums, Cosign, binários `simulador` e ausência da pasta intermediária antiga.
