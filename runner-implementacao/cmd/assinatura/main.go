package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/BarbaraNogueiraCS/runner/internal/apperrors"
	"github.com/BarbaraNogueiraCS/runner/internal/assinador"
	"github.com/spf13/cobra"
)

var version = "dev"

type flags map[string]string

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	exitCode := apperrors.OK
	root := newRootCommand(&exitCode)
	root.SetArgs(args)
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return apperrors.UsageError
	}
	return exitCode
}

func newRootCommand(exitCode *int) *cobra.Command {
	root := &cobra.Command{
		Use:           "assinatura",
		Short:         "Sistema Runner - CLI de assinatura digital",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(cmd *cobra.Command, args []string) {
			usage()
			*exitCode = apperrors.OK
		},
	}
	root.SetVersionTemplate("{{.Version}}\n")

	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Exibe a versão atual do CLI",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), version)
			*exitCode = apperrors.OK
		},
	})

	root.AddCommand(&cobra.Command{
		Use:                "sign",
		Short:              "Cria uma assinatura digital",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			*exitCode = sign(args)
		},
	})

	root.AddCommand(&cobra.Command{
		Use:                "validate",
		Short:              "Valida uma assinatura digital",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			*exitCode = validate(args)
		},
	})

	root.AddCommand(&cobra.Command{
		Use:                "server",
		Short:              "Gerencia o assinador.jar em modo servidor",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			*exitCode = server(args)
		},
	})

	return root
}

func sign(args []string) int {
	f, err := parseFlags(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return apperrors.UsageError
	}
	if f["input"] == "" && (f["bundle"] == "" || f["provenance"] == "" || f["crypto-material"] == "" || f["cert-chain"] == "" || f["timestamp"] == "" || f["policy"] == "") {
		fmt.Fprintln(os.Stderr, "Uso inválido: informe --bundle, --provenance, --crypto-material, --cert-chain, --timestamp e --policy. Para compatibilidade, --input ainda gera uma assinatura simulada simples.")
		return apperrors.UsageError
	}
	client := assinador.New(intFlag(f, "port", assinador.DefaultPort), f["jar"])
	client.ManifestURL = f["release-json"]
	req := assinador.SignRequest{Bundle: f["bundle"], Provenance: f["provenance"], CryptoMaterial: f["crypto-material"], CertificateChain: f["cert-chain"], Timestamp: f["timestamp"], Strategy: f["strategy"], Policy: f["policy"], Config: f["config"], Signer: f["signer"], Input: f["input"]}
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
	client.ManifestURL = f["release-json"]
	req := assinador.ValidateRequest{Signature: signature, Timestamp: f["timestamp"], Policy: f["policy"], Config: f["config"], Bundle: f["bundle"], Provenance: f["provenance"], Input: f["input"]}
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
	client.ManifestURL = f["release-json"]
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
  assinatura sign --bundle <bundle.json> --provenance <provenance.json> --crypto-material <crypto.json> --cert-chain <certs.json> --timestamp <unix> --policy <uri> [--strategy iat] [--output <arquivo>] [--local] [--jar <assinador.jar>] [--port 8080]
  assinatura validate --signature <signature.json> --timestamp <unix> --policy <uri> [--bundle <bundle.json> --provenance <provenance.json>] [--output <arquivo>] [--local] [--jar <assinador.jar>] [--port 8080]

Compatibilidade:
  assinatura sign --input <arquivo> --signer <nome> [--local] [--jar <assinador.jar>]
  assinatura server start --jar <assinador.jar> [--port 8080] [--idle-timeout-minutes 10] [--release-json <url>]
  assinatura server status [--port 8080]
  assinatura server stop [--port 8080]

Variáveis de ambiente:
  RUNNER_ASSINADOR_JAR  Caminho padrão do assinador.jar
  RUNNER_JAVA           Caminho explícito do executável java
  RUNNER_HOME           Diretório de estado, logs e cache. Padrão: ~/.hubsaude
  release-json          URL opcional para release.json com URLs de runtime Java`)
}

func parseFlags(args []string) (flags, error) {
	out := flags{}
	for i := 0; i < len(args); i++ {
		a := args[i]
		if !strings.HasPrefix(a, "--") {
			return out, fmt.Errorf("argumento inesperado: %s", a)
		}
		keyValue := strings.TrimPrefix(a, "--")
		key := keyValue
		value := ""
		if idx := strings.Index(keyValue, "="); idx >= 0 {
			key = keyValue[:idx]
			value = keyValue[idx+1:]
		}
		if key == "" {
			return out, fmt.Errorf("flag inválida: %s", a)
		}
		if key == "local" || key == "help" {
			if value != "" {
				return out, fmt.Errorf("flag booleana não deve receber valor: --%s", key)
			}
			out[key] = "true"
			continue
		}
		if value == "" {
			if i+1 >= len(args) {
				return out, fmt.Errorf("flag sem valor: --%s", key)
			}
			value = args[i+1]
			i++
		}
		out[key] = value
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
