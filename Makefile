# Makefile de automação local do Sistema Runner.
#
# Objetivo: evitar repetir manualmente comandos como go mod download,
# go test, go vet e go build. O GitHub Actions continua executando
# go mod download em ambiente limpo; este Makefile automatiza o fluxo local.

.PHONY: deps tidy test vet cover check build java-test samples clean all help

VERSION ?= dev

help:
	@echo "Alvos disponíveis:"
	@echo "  make deps       - baixa dependências Go declaradas no go.mod"
	@echo "  make tidy       - organiza go.mod e go.sum"
	@echo "  make test       - baixa dependências e roda go test ./..."
	@echo "  make vet        - baixa dependências e roda go vet ./..."
	@echo "  make cover      - baixa dependências e roda go test -cover ./..."
	@echo "  make check      - valida higiene do repositório e artefatos de release"
	@echo "  make build      - compila os CLIs locais em dist/"
	@echo "  make java-test  - compila e testa o assinador.jar"
	@echo "  make samples    - gera arquivos de exemplo do assinador em assinador/target/"
	@echo "  make clean      - remove saídas locais geradas por build/teste"
	@echo "  make all        - roda o fluxo local completo"

# Baixa dependências sem modificar go.mod/go.sum.
# Ideal para CI e para preparar o ambiente local.
deps:
	go mod download

# Organiza go.mod e go.sum. Use quando adicionar/remover imports ou dependências.
tidy:
	go mod tidy

# Testes Go.
test: deps
	go test ./...

# Verificação estática Go.
vet: deps
	go vet ./...

# Cobertura de testes Go.
cover: deps
	go test -cover ./...

# Verificações de higiene do repositório e rastreabilidade da release.
check: deps
	./scripts/check-generated-files.sh
	./scripts/check-release-artifacts.sh

# Build local dos dois CLIs. Em release, o GitHub Actions sobrescreve VERSION pela tag.
build: deps
	mkdir -p dist
	go build -ldflags "-X main.version=$(VERSION)" -o dist/assinatura ./cmd/assinatura
	go build -ldflags "-X main.version=$(VERSION)" -o dist/simulador ./cmd/simulador

# Compila e testa o componente Java.
java-test:
	$(MAKE) -C assinador clean all test

# Gera os JSONs de exemplo usados nos fluxos locais de assinatura.
samples:
	$(MAKE) -C assinador samples

# Remove arquivos gerados localmente. Esses diretórios não devem ser versionados.
clean:
	rm -rf dist
	rm -rf assinador/out
	rm -rf assinador/target
	rm -f examples/*.json

# Fluxo completo local para validação antes de commit/tag.
all: deps test vet cover check java-test build
