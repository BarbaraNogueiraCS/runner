package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/BarbaraNogueiraCS/runner/internal/apperrors"
	"github.com/BarbaraNogueiraCS/runner/internal/assinador"
	"github.com/spf13/cobra"
)

var version = "dev"

type flags map[string]string

var signAllowedFlags = map[string]bool{
	"bundle": true, "provenance": true, "crypto-material": true, "cert-chain": true,
	"timestamp": true, "strategy": true, "policy": true, "config": true,
	"signer": true, "input": true, "output": true, "local": true,
	"jar": true, "port": true, "idle-timeout-minutes": true, "timeout": true, "release-json": true,
}

var validateAllowedFlags = map[string]bool{
	"signature": true, "timestamp": true, "policy": true, "config": true,
	"bundle": true, "provenance": true, "input": true, "output": true, "local": true,
	"jar": true, "port": true, "idle-timeout-minutes": true, "timeout": true, "release-json": true,
}

var serverAllowedFlags = map[string]bool{
	"jar": true, "port": true, "idle-timeout-minutes": true, "timeout": true, "release-json": true,
}

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
		Use:                "start",
		Short:              "Inicia o assinador.jar em modo servidor",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			*exitCode = server(append([]string{"start"}, args...))
		},
	})

	root.AddCommand(&cobra.Command{
		Use:                "status",
		Short:              "Mostra o status do assinador.jar em modo servidor",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			*exitCode = server(append([]string{"status"}, args...))
		},
	})

	root.AddCommand(&cobra.Command{
		Use:                "stop",
		Short:              "Interrompe o assinador.jar em modo servidor",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			*exitCode = server(append([]string{"stop"}, args...))
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
	if containsHelp(args) {
		signUsage(os.Stdout)
		return apperrors.OK
	}
	f, err := parseFlagsKnown(args, signAllowedFlags)
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
		var err error
		out, err = signViaServerOrLocal(client, req, idleMinutesFlag(f))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
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
	if containsHelp(args) {
		validateUsage(os.Stdout)
		return apperrors.OK
	}
	f, err := parseFlagsKnown(args, validateAllowedFlags)
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
		var err error
		out, err = validateViaServerOrLocal(client, req, idleMinutesFlag(f))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
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
		serverUsage(os.Stderr)
		return apperrors.UsageError
	}
	if containsHelp(args) {
		serverUsage(os.Stdout)
		return apperrors.OK
	}
	action := args[0]
	f, err := parseFlagsKnown(args[1:], serverAllowedFlags)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return apperrors.UsageError
	}
	client := assinador.New(intFlag(f, "port", assinador.DefaultPort), f["jar"])
	client.ManifestURL = f["release-json"]
	switch action {
	case "start":
		state, reused, err := client.StartServer(idleMinutesFlag(f))
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
  assinatura start --jar <assinador.jar> [--port 8080] [--timeout 10]
  assinatura status [--port 8080]
  assinatura stop [--port 8080]
  assinatura server start --jar <assinador.jar> [--port 8080] [--timeout 10] [--release-json <url>]
  assinatura server status [--port 8080]
  assinatura server stop [--port 8080]

Variáveis de ambiente:
  RUNNER_ASSINADOR_JAR  Caminho padrão do assinador.jar
  RUNNER_JAVA           Caminho explícito do executável java
  RUNNER_HOME           Diretório de estado, logs e cache. Padrão: ~/.hubsaude
  release-json          URL opcional para release.json com URLs de runtime Java`)
}

func signUsage(w io.Writer) {
	fmt.Fprintln(w, `Uso:
  assinatura sign --bundle <bundle.json> --provenance <provenance.json> --crypto-material <crypto.json> --cert-chain <certs.json> --timestamp <unix> --policy <uri> [opções]

Parâmetros obrigatórios no fluxo principal:
  --bundle              Caminho do Bundle FHIR R4 em JSON
  --provenance          Caminho do Provenance FHIR R4 em JSON
  --crypto-material     Caminho do JSON de material criptográfico simulado
  --cert-chain          Caminho do JSON com cadeia de certificados em base64
  --timestamp           Timestamp Unix de referência
  --policy              URI da política de assinatura

Opções:
  --strategy iat        Estratégia simulada. Padrão: iat
  --config <arquivo>    Configuração adicional do assinador
  --signer <nome>       Nome exibido no campo who.display da assinatura
  --output <arquivo>    Grava o resultado no arquivo informado
  --local               Invoca diretamente java -jar assinador.jar
  --jar <arquivo>       Caminho do assinador.jar
  --port <porta>        Porta do servidor do assinador. Padrão: 8080
  --timeout <minutos>   Timeout de inatividade se o servidor for iniciado automaticamente
  --release-json <url>  Manifesto usado para provisionar Java/JDK
  --help                Exibe esta ajuda

Compatibilidade:
  assinatura sign --input <arquivo> --signer <nome> --local --jar <assinador.jar>`)
}

func validateUsage(w io.Writer) {
	fmt.Fprintln(w, `Uso:
  assinatura validate --signature <signature.json> --timestamp <unix> --policy <uri> [opções]

Parâmetros principais:
  --signature           Caminho da assinatura JSON a validar
  --timestamp           Timestamp Unix de referência
  --policy              URI da política esperada

Opções:
  --bundle <arquivo>    Bundle FHIR R4 usado para conferir integridade
  --provenance <arquivo> Provenance FHIR R4 usado para conferir integridade
  --config <arquivo>    Configuração adicional do assinador
  --output <arquivo>    Grava o resultado no arquivo informado
  --local               Invoca diretamente java -jar assinador.jar
  --jar <arquivo>       Caminho do assinador.jar
  --port <porta>        Porta do servidor do assinador. Padrão: 8080
  --timeout <minutos>   Timeout de inatividade se o servidor for iniciado automaticamente
  --release-json <url>  Manifesto usado para provisionar Java/JDK
  --help                Exibe esta ajuda`)
}

func serverUsage(w io.Writer) {
	fmt.Fprintln(w, `Uso:
  assinatura start --jar <assinador.jar> [--port 8080] [--timeout 10]
  assinatura status [--port 8080]
  assinatura stop [--port 8080]

Compatibilidade:
  assinatura server start --jar <assinador.jar> [--port 8080] [--timeout 10]
  assinatura server status [--port 8080]
  assinatura server stop [--port 8080]

Opções:
  --jar <arquivo>       Caminho do assinador.jar
  --port <porta>        Porta do servidor do assinador. Padrão: 8080
  --timeout <minutos> Encerra o servidor após período ocioso
  --idle-timeout-minutes Alias compatível de --timeout
  --release-json <url>  Manifesto usado para provisionar Java/JDK
  --help                Exibe esta ajuda`)
}

func containsHelp(args []string) bool {
	for _, arg := range args {
		if arg == "--help" || arg == "-h" || arg == "help" {
			return true
		}
	}
	return false
}

func parseFlagsKnown(args []string, allowed map[string]bool) (flags, error) {
	parsed, err := parseFlags(args)
	if err != nil {
		return parsed, err
	}
	for key := range parsed {
		if !allowed[key] {
			return parsed, fmt.Errorf("flag desconhecida para este comando: --%s. Use --help para ver as opções disponíveis", key)
		}
	}
	return parsed, nil
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

func idleMinutesFlag(f flags) int {
	if f["timeout"] != "" {
		return intFlag(f, "timeout", 0)
	}
	return intFlag(f, "idle-timeout-minutes", 0)
}

func signViaServerOrLocal(client assinador.Client, req assinador.SignRequest, idleMinutes int) ([]byte, error) {
	if err := client.Health(); err != nil {
		if _, reused, startErr := client.StartServer(idleMinutes); startErr == nil {
			if !reused {
				fmt.Fprintln(os.Stderr, "assinador.jar iniciado em modo servidor.")
			}
		} else {
			stdout, stderr, _, localErr := client.SignLocal(req)
			if localErr == nil {
				fmt.Fprintf(os.Stderr, "Servidor indisponível; usando fallback local: %v\n", startErr)
				if len(stderr) > 0 {
					_, _ = os.Stderr.Write(stderr)
				}
				return stdout, nil
			}
			return nil, fmt.Errorf("servidor do assinador indisponível: %v. Fallback local também falhou: %s", startErr, strings.TrimSpace(string(stderr)))
		}
	}
	out, err := client.SignHTTP(req)
	if err != nil {
		if len(out) > 0 {
			_, _ = os.Stderr.Write(append(out, '\n'))
		}
		stdout, stderr, _, localErr := client.SignLocal(req)
		if localErr == nil {
			fmt.Fprintf(os.Stderr, "Falha ao invocar /sign; usando fallback local: %v\n", err)
			if len(stderr) > 0 {
				_, _ = os.Stderr.Write(stderr)
			}
			return stdout, nil
		}
		return nil, fmt.Errorf("falha ao invocar /sign: %v. Fallback local também falhou: %s", err, strings.TrimSpace(string(stderr)))
	}
	return out, nil
}

func validateViaServerOrLocal(client assinador.Client, req assinador.ValidateRequest, idleMinutes int) ([]byte, error) {
	if err := client.Health(); err != nil {
		if _, reused, startErr := client.StartServer(idleMinutes); startErr == nil {
			if !reused {
				fmt.Fprintln(os.Stderr, "assinador.jar iniciado em modo servidor.")
			}
		} else {
			stdout, stderr, _, localErr := client.ValidateLocal(req)
			if localErr == nil {
				fmt.Fprintf(os.Stderr, "Servidor indisponível; usando fallback local: %v\n", startErr)
				if len(stderr) > 0 {
					_, _ = os.Stderr.Write(stderr)
				}
				return stdout, nil
			}
			return nil, fmt.Errorf("servidor do assinador indisponível: %v. Fallback local também falhou: %s", startErr, strings.TrimSpace(string(stderr)))
		}
	}
	out, err := client.ValidateHTTP(req)
	if err != nil {
		if len(out) > 0 {
			_, _ = os.Stderr.Write(append(out, '\n'))
		}
		stdout, stderr, _, localErr := client.ValidateLocal(req)
		if localErr == nil {
			fmt.Fprintf(os.Stderr, "Falha ao invocar /validate; usando fallback local: %v\n", err)
			if len(stderr) > 0 {
				_, _ = os.Stderr.Write(stderr)
			}
			return stdout, nil
		}
		return nil, fmt.Errorf("falha ao invocar /validate: %v. Fallback local também falhou: %s", err, strings.TrimSpace(string(stderr)))
	}
	return out, nil
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
