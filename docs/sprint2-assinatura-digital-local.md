# Sprint 2 — Assinatura Digital Simulada em modo local

## US-02.1 — Simulação de criação de assinatura digital

| Critério | Status | Evidência |
|---|---|---|
| Projeto Java base | Implementado | `projetos/assinador-java` e implementação canônica em `assinador/` |
| Interface `SignatureService` | Implementado | `assinador/src/br/ufg/hubsaude/assinador/SignatureService.java` |
| `FakeSignatureService` | Implementado | `assinador/src/br/ufg/hubsaude/assinador/FakeSignatureService.java` |
| Resposta simulada com campos esperados | Implementado | `Signature`, `JWS`, `policy`, `timestamp` |
| Testes de sucesso | Implementado | `make java-test` |

## US-02.2 — Validação de parâmetros de criação

O `assinador.jar` verifica presença e formato dos parâmetros obrigatórios. Erros como `INPUT.MISSING-PARAMETER` são retornados antes do processamento.

## US-02.3 — Validação de assinatura

O modo `validate` verifica parâmetros de validação e retorna resultado pré-determinado baseado em critérios simples. O cenário válido retorna `valid: true`.

## US-01.2 — Parsing de comandos e parâmetros no CLI

O CLI aceita:

```bash
./dist/assinatura sign --help
./dist/assinatura validate --help
```

A ajuda documenta flags como `--jar`, `--bundle`, `--provenance`, `--crypto-material`, `--cert-chain`, `--timestamp`, `--policy`, `--signature`, `--local` e `--output`.

## US-01.3 — Invocação local do assinador.jar

O CLI localiza Java 21 e executa `java -jar assinador.jar` com os parâmetros mapeados. A saída do JAR é capturada e repassada ao usuário. Erros como JDK ausente ou JAR não encontrado são tratados com mensagens claras.

## US-01.4 — Exibição legível

A saída é JSON compacto/legível para terminal, usando `compactJSON` quando adequado. A validação indica claramente `valid: true` ou `valid: false`.

## US-04.1 — Detecção e provisionamento automático do JDK

O pacote `internal/jdk` contém `EnsureJava21` e `ManagedJavaPathFor`. O sistema procura Java 21 no PATH e em `~/.hubsaude/jdk`; se necessário, usa o manifesto `release.json` para provisionamento. O JDK baixado é reutilizado para não repetir download.

## Execução local

```bash
make deps
make java-test
make samples
make build
```

Assinatura local:

```bash
TS=$(date -u +%s)
POLICY='https://fhir.saude.go.gov.br/r4/seguranca/ImplementationGuide/br.go.ses.seguranca|0.1.2'
./dist/assinatura sign --local --jar assinador/target/assinador.jar \
  --bundle assinador/target/bundle.json \
  --provenance assinador/target/provenance.json \
  --crypto-material assinador/target/crypto.json \
  --cert-chain assinador/target/certs.json \
  --timestamp "$TS" \
  --policy "$POLICY" \
  --config assinador/target/config.json \
  --signer "Maria Runner" \
  --output examples/assinatura-local.json
```
