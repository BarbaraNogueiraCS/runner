# Sistema Runner — Sprint 1

Este repositório contém a entrega da Sprint 1 do Sistema Runner.

## O que existe nesta sprint

- Projeto Go inicializado como `github.com/kyriosdata/runner`.
- CLI `assinatura` com comando `version`.
- Stub da CLI `simulador`.
- Estrutura de pacotes em `internal/` para evolução futura.
- Diretório `assinador/` reservado para o projeto Java/Maven da Sprint 2.
- Workflow de build multiplataforma.
- Workflow de release com versionamento por tag, checksums SHA256 e assinatura Cosign.

## Estrutura

```text
runner/
├── cmd/
│   ├── assinatura/
│   │   ├── main.go
│   │   └── version_test.go
│   └── simulador/
│       └── main.go
├── internal/
│   ├── cli/
│   ├── invoker/
│   ├── jdk/
│   └── release/
├── assinador/
├── .github/workflows/
│   ├── build.yml
│   └── release.yml
├── go.mod
├── go.sum
└── README.md
```

## Requisitos

- Go 1.25.
- Git.
- Acesso à internet para baixar dependências Go na primeira execução.

## Como executar localmente

Instale as dependências:

```bash
go mod tidy
```

Execute a CLI de assinatura:

```bash
go run ./cmd/assinatura version
```

Saída esperada em ambiente local:

```text
dev
```

Execute o stub do simulador:

```bash
go run ./cmd/simulador
```

Saída esperada:

```text
simulador vdev — em construção
```

## Como testar

```bash
go vet ./...
go test ./...
```

## Como gerar binários localmente

Linux:

```bash
GOOS=linux GOARCH=amd64 go build -o dist/assinatura-linux-amd64 ./cmd/assinatura
```

Windows:

```bash
GOOS=windows GOARCH=amd64 go build -o dist/assinatura-windows-amd64.exe ./cmd/assinatura
```

macOS:

```bash
GOOS=darwin GOARCH=amd64 go build -o dist/assinatura-darwin-amd64 ./cmd/assinatura
```

## Como criar uma release

Crie e envie uma tag semântica:

```bash
git tag -a v0.1.0 -m "Primeira entrega da Sprint 1"
git push origin v0.1.0
```

O workflow `.github/workflows/release.yml` irá:

1. Executar `go vet ./...` e `go test ./...` em Linux, Windows e macOS.
2. Gerar binários versionados.
3. Gerar `checksums.txt` com SHA256.
4. Assinar os binários com Cosign usando identidade OIDC.
5. Publicar tudo no GitHub Releases.

## Como verificar um artefato assinado

Exemplo para Linux:

```bash
cosign verify-blob \
  --certificate assinatura-v0.1.0-linux-amd64.pem \
  --signature assinatura-v0.1.0-linux-amd64.sig \
  assinatura-v0.1.0-linux-amd64
```
