package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/kyriosdata/runner/internal/apperrors"
	"github.com/kyriosdata/runner/internal/simulator"
)

var version = "dev"

type flags map[string]string

func main() { os.Exit(run(os.Args[1:])) }

func run(args []string) int {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		usage()
		return apperrors.OK
	}
	switch args[0] {
	case "version", "--version", "-v":
		fmt.Println(version)
		return apperrors.OK
	case "status":
		return status(args[1:])
	case "info":
		return info(args[1:])
	case "stop":
		return stop(args[1:])
	case "start":
		return start(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "Comando desconhecido: %s\n\n", args[0])
		usage()
		return apperrors.UsageError
	}
}

func status(args []string) int {
	f, err := parseFlags(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return apperrors.UsageError
	}
	c := simulator.New(f["url"], boolFlag(f, "insecure"), f["jar"])
	c.Artifact = defaultText(f["artifact"], c.Artifact)
	c.ManifestURL = defaultText(f["release-json"], c.ManifestURL)
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
	c := simulator.New(f["url"], boolFlag(f, "insecure"), f["jar"])
	c.Artifact = defaultText(f["artifact"], c.Artifact)
	c.ManifestURL = defaultText(f["release-json"], c.ManifestURL)
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
	c := simulator.New(f["url"], boolFlag(f, "insecure"), f["jar"])
	c.Artifact = defaultText(f["artifact"], c.Artifact)
	c.ManifestURL = defaultText(f["release-json"], c.ManifestURL)
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
	c := simulator.New(f["url"], boolFlag(f, "insecure"), f["jar"])
	c.Artifact = defaultText(f["artifact"], c.Artifact)
	c.ManifestURL = defaultText(f["release-json"], c.ManifestURL)
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
  simulador status [--url https://localhost:8443] [--insecure]
  simulador info   [--url https://localhost:8443] [--insecure]
  simulador stop   [--url https://localhost:8443] [--insecure]
  simulador start  [--url https://localhost:8443] [--insecure] [--jar <simulador.jar>] [--artifact simulador|validador] [--release-json <url>]

Observação:
  O Simulador HubSaúde informado no trabalho usa HTTPS local com certificado autoassinado.
  Por isso, em ambiente de desenvolvimento, use --insecure ao consultar https://localhost:8443.

Variáveis de ambiente:
  HUBSAUDE_SIMULATOR_URL  URL padrão do simulador. Padrão: https://localhost:8443
  RUNNER_JAVA             Caminho explícito do executável java
  RUNNER_HOME             Diretório de estado, logs e cache. Padrão: ~/.hubsaude
  RUNNER_RELEASE_JSON     URL opcional para release.json
  RUNNER_SIMULADOR_ARTIFACT Artefato padrão do release.json: simulador ou validador`)
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
