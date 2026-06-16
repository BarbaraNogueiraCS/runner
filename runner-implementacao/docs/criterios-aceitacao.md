# Verificação dos Critérios de Aceitação

Este documento registra a rastreabilidade entre os critérios solicitados e os pontos de implementação no código.

## 1. CLI `assinatura` e invocação do `assinador.jar`

| Critério | Implementação |
|---|---|
| Aceitar comandos para criação e validação | `cmd/assinatura/main.go`, comandos `sign` e `validate` |
| Invocar `assinador.jar` com parâmetros fornecidos | `internal/assinador/client.go`, métodos `SignLocal`, `ValidateLocal`, `SignHTTP` e `ValidateHTTP` |
| Permitir invocação direta/local | `assinatura sign --local --jar ...` e `assinatura validate --local --jar ...` |
| Permitir invocação via HTTP | modo padrão sem `--local`, com `POST /sign` e `POST /validate` |
| Exibir resultado legível | `compactJSON` e `WriteOutput` formatam JSON ou gravam em arquivo |
| Iniciar servidor na porta padrão | `internal/assinador.DefaultPort = 8080` e `StartServer` |
| Detectar instância em execução | `Client.Health()` antes de iniciar; se saudável, `StartServer` retorna `reused=true` |
| Usar modo servidor por padrão | `cmd/assinatura/main.go`: quando `--local` não é informado, usa HTTP e inicia servidor se necessário |
| Interromper execução | `assinatura server stop [--port]` chama `POST /shutdown` |
| Interrupção programada por inatividade | `assinatura server start --idle-timeout-minutes N`; `Main.java` encerra após janela sem requisições |

## 2. `assinador.jar`

| Critério | Implementação |
|---|---|
| Validar parâmetros conforme especificação | `GuideSignatureService`: valida Bundle, Provenance, timestamp, política, estratégia, material criptográfico, cadeia e Signature.data |
| Simular criação de assinatura | `GuideSignatureService.sign` monta `Signature` com JWS JSON Serialization simulado |
| Simular validação de assinatura | `GuideSignatureService.validate` retorna `OperationOutcome` com `valid=true` ou erro claro |
| Suportar PKCS#11 | `Pkcs11Support`: valida SMARTCARD/TOKEN, PIN, biblioteca PKCS#11 e abre `KeyStore` PKCS11 quando não está em simulação |
| Mensagens claras de erro | `UserInputException` + `Json.error`, com códigos como `INPUT.MISSING-PARAMETER`, `CRYPTO.PIN-REQUIRED`, `PKCS11.LIBRARY-REQUIRED` |

## 3. CLI `simulador`

| Critério | Implementação |
|---|---|
| Verificar porta 8443 antes de iniciar | `internal/simulador/client.go`, método `Start`, usa `Status` e `netutil.IsTCPPortFree` |
| Permitir iniciar simulador | `simulador start` |
| Permitir parar simulador | `simulador stop`, endpoint `/shutdown` |
| Exibir status | `simulador status`, endpoint `/api/info` |
| Obter simulador dinamicamente | `internal/release.EnsureArtifact` com `release.json` |
| Baixar JRE/JDK Temurin quando necessário | `internal/jdk.EnsureJava21` e `internal/release.DownloadAndInstallJRE` |
| Não baixar novamente se versão já existe | `EnsureArtifact` compara `.version` e SHA256 antes de baixar |

## 4. Provisionamento Java

| Critério | Implementação |
|---|---|
| Detectar JDK/Java presente | `internal/jdk.FindJava`, `MajorVersion` e `EnsureJava21` |
| Baixar JDK/JRE compatível quando ausente | `release.JREURL` prefere `jdk` no manifesto e usa `jre` como fallback |
| Disponibilizar runtime para Assinador e Simulador | runtime gerenciado em `~/.hubsaude/jdk/bin/java`, usado por ambos os CLIs |
| Download em três plataformas | `PlatformKey`: `windows_x64`, `linux_x64`, `mac_x64`, com suporte adicional a ARM64 quando presente no manifesto |

## 5. Distribuição

| Critério | Implementação |
|---|---|
| Binário Windows amd64 | `.github/workflows/release.yml`, artefatos `*.exe` |
| Binário Linux amd64 | `.github/workflows/release.yml`, artefatos `*.AppImage` |
| Binário macOS amd64 | `.github/workflows/release.yml`, artefatos `*.dmg` |
| GitHub Releases | `softprops/action-gh-release@v2` |
| Checksums SHA256 | etapa `sha256sum * > checksums.txt` |
| SemVer | workflow aceita tags `v*` e valida padrão `vMAJOR.MINOR.PATCH` |
| Assinatura Cosign | workflow usa `sigstore/cosign-installer` e `cosign sign-blob` com OIDC/keyless |
