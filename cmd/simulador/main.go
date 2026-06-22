package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/BarbaraNogueiraCS/runner/internal/apperrors"
	"github.com/BarbaraNogueiraCS/runner/internal/simulador"
	"github.com/spf13/cobra"
)

var version = "dev"

type flags map[string]string

func main() { os.Exit(run(os.Args[1:])) }

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
		Use:           "simulador",
		Short:         "Sistema Runner - CLI do Simulador HubSaúde",
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
		Use:                "status",
		Short:              "Consulta o status do Simulador HubSaúde",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			*exitCode = status(args)
		},
	})

	root.AddCommand(&cobra.Command{
		Use:                "info",
		Short:              "Consulta informações do Simulador HubSaúde",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			*exitCode = info(args)
		},
	})

	root.AddCommand(&cobra.Command{
		Use:                "stop",
		Short:              "Solicita encerramento do Simulador HubSaúde",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			*exitCode = stop(args)
		},
	})

	root.AddCommand(&cobra.Command{
		Use:                "start",
		Short:              "Inicia ou reaproveita o Simulador HubSaúde",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			*exitCode = start(args)
		},
	})

	return root
}

func status(args []string) int {
	f, err := parseFlags(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return apperrors.UsageError
	}
	c, err := newClientFromFlags(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return apperrors.UsageError
	}
	state, info, err := c.Status()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Simulador do HubSaúde indisponível em %s. Motivo: %v\n", c.BaseURL, err)
		fmt.Fprintln(os.Stderr, "Como resolver: verifique se o simulador está em execução ou use 'simulador start --jar <simulador.jar> --insecure'.")
		return apperrors.IntegrationError
	}
	fmt.Printf("Simulador do HubSaúde disponível em %s", c.BaseURL)
	if state.PID != 0 {
		fmt.Printf(" (PID %d)", state.PID)
	}
	fmt.Println()
	if version := jsonField(info, "version"); version != "" {
		fmt.Println("Versão:", version)
	}
	return apperrors.OK
}

func info(args []string) int {
	f, err := parseFlags(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return apperrors.UsageError
	}
	c, err := newClientFromFlags(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return apperrors.UsageError
	}
	out, err := c.Info()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Não foi possível consultar /api/info em %s. Motivo: %v\n", c.BaseURL, err)
		return apperrors.IntegrationError
	}
	fmt.Println(string(compactJSON(out)))
	return apperrors.OK
}

func stop(args []string) int {
	f, err := parseFlags(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return apperrors.UsageError
	}
	c, err := newClientFromFlags(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return apperrors.UsageError
	}
	out, err := c.Stop()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Não foi possível encerrar o Simulador do HubSaúde em %s. Motivo: %v\n", c.BaseURL, err)
		return apperrors.IntegrationError
	}
	if len(out) > 0 && strings.TrimSpace(string(out)) != "" {
		fmt.Println(string(out))
	} else {
		fmt.Println("Solicitação de encerramento enviada ao Simulador do HubSaúde.")
	}
	return apperrors.OK
}

func start(args []string) int {
	f, err := parseFlags(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return apperrors.UsageError
	}
	c, err := newClientFromFlags(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return apperrors.UsageError
	}
	state, reused, err := c.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Não foi possível iniciar o Simulador do HubSaúde: %v\n", err)
		return apperrors.IntegrationError
	}
	if reused {
		fmt.Printf("Simulador do HubSaúde já está em execução em %s\n", state.URL)
	} else {
		fmt.Printf("Simulador do HubSaúde iniciado em %s (PID %d)\n", state.URL, state.PID)
	}
	return apperrors.OK
}

func usage() {
	fmt.Println(`Sistema Runner - CLI simulador

Uso:
  simulador version
  simulador status [--url https://localhost:8443] [--port 8443] [--insecure]
  simulador info   [--url https://localhost:8443] [--port 8443] [--insecure]
  simulador stop   [--url https://localhost:8443] [--port 8443] [--insecure]
  simulador start  [--url https://localhost:8443] [--port 8443] [--insecure] [--jar <simulador.jar>]
                   [--artifact simulador|validador] [--release-json <url>] [--source <url>] [--sha256 <hash>]

Observação:
  O Simulador HubSaúde informado no trabalho usa HTTPS local com certificado autoassinado.
  Por isso, em ambiente de desenvolvimento, use --insecure ao consultar https://localhost:8443.

Variáveis de ambiente:
  HUBSAUDE_SIMULADOR_URL  URL padrão do simulador. Padrão: https://localhost:8443
  RUNNER_JAVA             Caminho explícito do executável java
  RUNNER_HOME             Diretório de estado, logs e cache. Padrão: ~/.hubsaude
  RUNNER_RELEASE_JSON        URL opcional para release.json
  RUNNER_SIMULADOR_ARTIFACT  Artefato padrão do release.json: simulador ou validador
  RUNNER_SIMULADOR_SOURCE    URL direta alternativa para baixar simulador.jar
  RUNNER_SIMULADOR_SHA256    SHA-256 esperado quando RUNNER_SIMULADOR_SOURCE for usado`)
}

func newClientFromFlags(f flags) (simulador.Client, error) {
	baseURL, err := urlWithPort(f["url"], f["port"])
	if err != nil {
		return simulador.Client{}, err
	}
	c := simulador.New(baseURL, boolFlag(f, "insecure"), f["jar"])
	c.Artifact = defaultText(f["artifact"], c.Artifact)
	c.ManifestURL = defaultText(f["release-json"], c.ManifestURL)
	c.SourceURL = defaultText(f["source"], c.SourceURL)
	c.SourceSHA256 = defaultText(defaultText(f["sha256"], f["checksum"]), c.SourceSHA256)
	return c, nil
}

func urlWithPort(rawURL, rawPort string) (string, error) {
	if strings.TrimSpace(rawPort) == "" {
		return rawURL, nil
	}
	port, err := strconv.Atoi(rawPort)
	if err != nil || port <= 0 || port > 65535 {
		return "", fmt.Errorf("porta inválida em --port: %s", rawPort)
	}
	if strings.TrimSpace(rawURL) == "" {
		return fmt.Sprintf("https://localhost:%d", port), nil
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("URL inválida em --url: %w", err)
	}
	if u.Scheme == "" {
		u.Scheme = "https"
	}
	if u.Host == "" {
		u.Host = "localhost"
	}
	host := u.Hostname()
	if host == "" {
		host = "localhost"
	}
	u.Host = fmt.Sprintf("%s:%d", host, port)
	return strings.TrimRight(u.String(), "/"), nil
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
		if key == "insecure" || key == "help" {
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

func defaultText(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
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

func jsonField(b []byte, field string) string {
	var m map[string]any
	if json.Unmarshal(b, &m) != nil {
		return ""
	}
	if v, ok := m[field].(string); ok {
		return v
	}
	return ""
}
