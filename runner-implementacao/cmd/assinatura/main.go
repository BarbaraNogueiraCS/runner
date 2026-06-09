package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kyriosdata/runner/internal/apperrors"
	"github.com/kyriosdata/runner/internal/assinador"
)

var version = "dev"

type flags map[string]string

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		usage()
		return apperrors.OK
	}
	switch args[0] {
	case "version", "--version", "-v":
		fmt.Println(version)
		return apperrors.OK
	case "sign":
		return sign(args[1:])
	case "validate":
		return validate(args[1:])
	case "server":
		return server(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "Comando desconhecido: %s\n\n", args[0])
		usage()
		return apperrors.UsageError
	}
}

func sign(args []string) int {
	f, err := parseFlags(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return apperrors.UsageError
	}
	input := f["input"]
	signer := f["signer"]
	if input == "" || signer == "" {
		fmt.Fprintln(os.Stderr, "Uso inválido: informe --input e --signer. A validação detalhada é realizada pelo assinador.jar.")
		return apperrors.UsageError
	}
	client := assinador.New(intFlag(f, "port", assinador.DefaultPort), f["jar"])
	req := assinador.SignRequest{Input: input, Signer: signer}
	var out []byte
	if boolFlag(f, "local") {
		stdout, stderr, code, err := client.SignLocal(req)
		if len(stderr) > 0 {
			_, _ = os.Stderr.Write(stderr)
		}
		if err != nil {
			return code
		}
		out = stdout
	} else {
		if err := client.Health(); err != nil {
			if _, reused, startErr := client.StartServer(intFlag(f, "idle-timeout-minutes", 0)); startErr == nil {
				if !reused {
					fmt.Fprintln(os.Stderr, "assinador.jar iniciado em modo servidor.")
				}
			} else {
				fmt.Fprintf(os.Stderr, "Servidor do assinador indisponível: %v\n", startErr)
				fmt.Fprintln(os.Stderr, "Como resolver: informe --jar para iniciar o servidor ou use --local --jar para invocação direta.")
				return apperrors.IntegrationError
			}
		}
		var err error
		out, err = client.SignHTTP(req)
		if err != nil {
			if len(out) > 0 {
				_, _ = os.Stderr.Write(append(out, '\n'))
			}
			fmt.Fprintf(os.Stderr, "Falha ao invocar /sign: %v\n", err)
			return apperrors.IntegrationError
		}
	}
	if err := assinador.WriteOutput(compactJSON(out), f["output"], os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Não foi possível gravar resultado: %v\n", err)
		return apperrors.InternalError
	}
	return apperrors.OK
}

func validate(args []string) int {
	f, err := parseFlags(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return apperrors.UsageError
	}
	signature := f["signature"]
	if signature == "" {
		fmt.Fprintln(os.Stderr, "Uso inválido: informe --signature. A validação detalhada é realizada pelo assinador.jar.")
		return apperrors.UsageError
	}
	client := assinador.New(intFlag(f, "port", assinador.DefaultPort), f["jar"])
	req := assinador.ValidateRequest{Signature: signature, Input: f["input"]}
	var out []byte
	if boolFlag(f, "local") {
		stdout, stderr, code, err := client.ValidateLocal(req)
		if len(stderr) > 0 {
			_, _ = os.Stderr.Write(stderr)
		}
		if err != nil {
			return code
		}
		out = stdout
	} else {
		if err := client.Health(); err != nil {
			if _, reused, startErr := client.StartServer(intFlag(f, "idle-timeout-minutes", 0)); startErr == nil {
				if !reused {
					fmt.Fprintln(os.Stderr, "assinador.jar iniciado em modo servidor.")
				}
			} else {
				fmt.Fprintf(os.Stderr, "Servidor do assinador indisponível: %v\n", startErr)
				fmt.Fprintln(os.Stderr, "Como resolver: informe --jar para iniciar o servidor ou use --local --jar para invocação direta.")
				return apperrors.IntegrationError
			}
		}
		var err error
		out, err = client.ValidateHTTP(req)
		if err != nil {
			if len(out) > 0 {
				_, _ = os.Stderr.Write(append(out, '\n'))
			}
			fmt.Fprintf(os.Stderr, "Falha ao invocar /validate: %v\n", err)
			return apperrors.IntegrationError
		}
	}
	if err := assinador.WriteOutput(compactJSON(out), f["output"], os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Não foi possível gravar resultado: %v\n", err)
		return apperrors.InternalError
	}
	return apperrors.OK
}

func server(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Uso inválido: use assinatura server start|status|stop")
		return apperrors.UsageError
	}
	action := args[0]
	f, err := parseFlags(args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return apperrors.UsageError
	}
	client := assinador.New(intFlag(f, "port", assinador.DefaultPort), f["jar"])
	switch action {
	case "start":
		state, reused, err := client.StartServer(intFlag(f, "idle-timeout-minutes", 0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Não foi possível iniciar assinador.jar: %v\n", err)
			return apperrors.IntegrationError
		}
		if reused {
			fmt.Printf("assinador.jar já está em execução em %s\n", state.URL)
		} else {
			fmt.Printf("assinador.jar iniciado em %s (PID %d)\n", state.URL, state.PID)
		}
		return apperrors.OK
	case "status":
		state, err := client.Status()
		if err != nil {
			fmt.Fprintf(os.Stderr, "assinador.jar não está disponível na porta %d: %v\n", client.Port, err)
			return apperrors.IntegrationError
		}
		fmt.Printf("assinador.jar disponível em %s", state.URL)
		if state.PID != 0 {
			fmt.Printf(" (PID %d)", state.PID)
		}
		fmt.Println()
		return apperrors.OK
	case "stop":
		if err := client.StopServer(); err != nil {
			fmt.Fprintf(os.Stderr, "Não foi possível encerrar assinador.jar: %v\n", err)
			return apperrors.IntegrationError
		}
		fmt.Println("Solicitação de encerramento enviada ao assinador.jar.")
		return apperrors.OK
	default:
		fmt.Fprintf(os.Stderr, "Ação desconhecida: %s\n", action)
		return apperrors.UsageError
	}
}

func usage() {
	fmt.Println(`Sistema Runner - CLI assinatura

Uso:
  assinatura version
  assinatura sign --input <arquivo> --signer <nome> [--output <arquivo>] [--local] [--jar <assinador.jar>] [--port 8080]
  assinatura validate --signature <arquivo> [--input <arquivo>] [--output <arquivo>] [--local] [--jar <assinador.jar>] [--port 8080]
  assinatura server start --jar <assinador.jar> [--port 8080] [--idle-timeout-minutes 10]
  assinatura server status [--port 8080]
  assinatura server stop [--port 8080]

Variáveis de ambiente:
  RUNNER_ASSINADOR_JAR  Caminho padrão do assinador.jar
  RUNNER_JAVA           Caminho explícito do executável java
  RUNNER_HOME           Diretório de estado, logs e cache. Padrão: ~/.hubsaude`)
}

func parseFlags(args []string) (flags, error) {
	out := flags{}
	for i := 0; i < len(args); i++ {
		a := args[i]
		if !strings.HasPrefix(a, "--") {
			return out, fmt.Errorf("argumento inesperado: %s", a)
		}
		key := strings.TrimPrefix(a, "--")
		if key == "local" || key == "help" {
			out[key] = "true"
			continue
		}
		if i+1 >= len(args) {
			return out, fmt.Errorf("flag sem valor: --%s", key)
		}
		out[key] = args[i+1]
		i++
	}
	return out, nil
}

func boolFlag(f flags, name string) bool { return f[name] == "true" }

func intFlag(f flags, name string, def int) int {
	v := f[name]
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func compactJSON(b []byte) []byte {
	var v any
	if json.Unmarshal(b, &v) != nil {
		return b
	}
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return b
	}
	return out
}
