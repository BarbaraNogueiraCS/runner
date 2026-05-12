# Sistema Runner — Implementação de Referência Acadêmica

Este repositório contém uma implementação de referência do **Sistema Runner**, desenvolvida a partir dos documentos de requisitos, arquitetura, projeto detalhado e Modelo C4 elaborados para a disciplina de Implementação e Integração de Software.

O objetivo do sistema é facilitar a execução e o gerenciamento de aplicações Java por meio de comandos CLI, especialmente para os fluxos de assinatura simulada (`assinatura` + `assinador.jar`) e gerenciamento do Simulador do HubSaúde (`simulador` + `simulador.jar`).

## Principais partes

- `assinatura`: CLI em Go para criar e validar assinaturas simuladas.
- `simulador`: CLI em Go para iniciar, parar e consultar status do Simulador do HubSaúde.
- `assinador.jar`: aplicação Java 21 que valida parâmetros e simula assinatura/validação.
- `~/.hubsaude/`: diretório local gerenciado para cache, metadados, processos e artefatos.
- GitHub Actions: workflows de build e release multiplataforma.

## Escopo desta implementação

Implementado:

- Estrutura modular dos CLIs em Go.
- Comandos `assinatura version`, `assinatura sign`, `assinatura validate`, `assinatura start`, `assinatura stop`.
- Comandos `simulador start`, `simulador stop`, `simulador status`.
- Invocação local do `assinador.jar` via `java -jar`.
- Invocação HTTP do `assinador.jar` via `/sign`, `/validate` e `/health`.
- Aplicação Java `assinador.jar` com modo CLI e modo servidor.
- Validação de parâmetros no Java e validação preliminar no Go.
- Registro local de processos em `~/.hubsaude/processos/`.
- Verificação de porta disponível.
- Download genérico de artefatos por URL.
- Verificação SHA256.
- Estrutura de CI/CD para builds e releases.

Limitações conhecidas desta versão acadêmica:

- A assinatura e a validação são **simuladas**, não criptográficas reais.
- O provisionamento automático completo de JDK/JRE foi estruturado, mas esta implementação prioriza o uso de Java 21 instalado no sistema ou já presente no diretório gerenciado.
- A integração PKCS#11 está isolada em adaptador no Java, mas não executa assinatura real neste escopo.
- O `simulador.jar` real não acompanha este pacote; o CLI aceita `--jar` ou `--source` para indicar uma origem.

## Pré-requisitos para desenvolvimento

- Go compatível com o módulo do projeto. A documentação do projeto define Go 1.25 como versão-alvo. Para permitir execução local/offline neste pacote, o `go.mod` foi deixado compatível com Go 1.23+ e contém um `replace` para uma implementação mínima local de Cobra. Em um repositório conectado à internet, recomenda-se remover o `replace` e usar o módulo oficial `github.com/spf13/cobra`.
- Java 21.
- Maven 3.9+ para o fluxo Maven, ou `javac`/`jar` para o script local `scripts/build-assinador.sh`.

## Como compilar o `assinador.jar`

```bash
cd projetos/assinador-java
mvn clean package
```

O arquivo será gerado em:

```text
projetos/assinador-java/target/assinador.jar
```

Para instalar no diretório gerenciado:

```bash
mkdir -p ~/.hubsaude/assinador
cp projetos/assinador-java/target/assinador.jar ~/.hubsaude/assinador/assinador.jar
```

## Como compilar os CLIs

```bash
go mod tidy
go build -o bin/assinatura ./cmd/assinatura
go build -o bin/simulador ./cmd/simulador
```

## Exemplos de uso

### Exibir versão

```bash
./bin/assinatura version
```

### Assinar em modo local

```bash
./bin/assinatura sign --local --jar ~/.hubsaude/assinador/assinador.jar --documento documento.json --certificado certificado.pem
```

### Validar em modo local

```bash
./bin/assinatura validate --local --jar ~/.hubsaude/assinador/assinador.jar --documento documento.json --assinatura assinatura-simulada
```

### Iniciar o assinador em modo servidor

```bash
./bin/assinatura start --jar ~/.hubsaude/assinador/assinador.jar --port 8080
```

### Assinar via HTTP

```bash
./bin/assinatura sign --documento documento.json --certificado certificado.pem --port 8080
```

### Parar o assinador

```bash
./bin/assinatura stop --port 8080
```

### Iniciar simulador com jar local

```bash
./bin/simulador start --jar ~/.hubsaude/simulador/simulador.jar --port 8443
```

## Testes

```bash
go test ./...
cd projetos/assinador-java && mvn test
```

## Observação importante

Esta implementação tem finalidade acadêmica e de integração. Ela não deve ser usada como solução de assinatura digital juridicamente válida, pois as operações de assinatura e validação são simuladas.
