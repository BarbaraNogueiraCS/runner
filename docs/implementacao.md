# Implementação

A implementação principal do Runner fica na raiz do repositório.

```text
runner/
├── go.mod
├── go.sum
├── Makefile
├── cmd/
├── internal/
├── assinador/
├── scripts/
├── examples/
├── docs/
├── projetos/
└── .github/
```

## Componentes

- `cmd/assinatura`: CLI para assinatura e validação simuladas.
- `cmd/simulador`: CLI para gerenciamento/consulta do Simulador do HubSaúde.
- `internal/assinador`: cliente Go responsável por invocar o `assinador.jar` localmente ou via servidor.
- `internal/jdk`: detecção e provisionamento de Java/JDK 21.
- `assinador`: projeto Java que gera o `assinador.jar`.
- `projetos/assinador-java`: wrapper de compatibilidade exigido pela Sprint 2.
- `scripts`: verificações de higiene, release e integridade.

## Build e testes

```bash
make deps
make test
make vet
make cover
make check
make java-test
make samples
make build
```

## Release

A geração dos binários, checksums SHA256, assinatura Cosign e publicação no GitHub Releases é feita por `.github/workflows/release.yml`.
