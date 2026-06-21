package assinador

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/BarbaraNogueiraCS/runner/internal/httpx"
	"github.com/BarbaraNogueiraCS/runner/internal/jdk"
	"github.com/BarbaraNogueiraCS/runner/internal/netutil"
	"github.com/BarbaraNogueiraCS/runner/internal/paths"
	"github.com/BarbaraNogueiraCS/runner/internal/process"
)

const DefaultPort = 8080

type Client struct {
	Port        int
	JarPath     string
	Timeout     time.Duration
	ManifestURL string
}

type SignRequest struct {
	Bundle           string `json:"bundle,omitempty"`
	Provenance       string `json:"provenance,omitempty"`
	CryptoMaterial   string `json:"cryptoMaterial,omitempty"`
	CertificateChain string `json:"certificateChain,omitempty"`
	Timestamp        string `json:"timestamp,omitempty"`
	Strategy         string `json:"strategy,omitempty"`
	Policy           string `json:"policy,omitempty"`
	Config           string `json:"config,omitempty"`
	Signer           string `json:"signer,omitempty"`
	Input            string `json:"input,omitempty"` // compatibilidade com modo legado
}

type ValidateRequest struct {
	Signature  string `json:"signature"`
	Timestamp  string `json:"timestamp,omitempty"`
	Policy     string `json:"policy,omitempty"`
	Config     string `json:"config,omitempty"`
	Bundle     string `json:"bundle,omitempty"`
	Provenance string `json:"provenance,omitempty"`
	Input      string `json:"input,omitempty"` // compatibilidade com modo legado
}

func New(port int, jar string) Client {
	if port == 0 {
		port = DefaultPort
	}
	if jar == "" {
		jar = os.Getenv("RUNNER_ASSINADOR_JAR")
	}
	return Client{Port: port, JarPath: jar, Timeout: 10 * time.Second}
}

func (c Client) baseURL() string { return fmt.Sprintf("http://127.0.0.1:%d", c.Port) }

func (c Client) matchesRequestedInstance(state process.State) bool {
	return state.Name == "assinador" && state.Port == c.Port && strings.TrimRight(state.URL, "/") == c.baseURL()
}

func (c Client) Health() error {
	client := httpx.New(2*time.Second, false)
	_, _, err := httpx.Get(client, c.baseURL()+"/health")
	return err
}

func (c Client) Status() (process.State, error) {
	if err := c.Health(); err != nil {
		return process.State{}, err
	}
	state, _ := process.Load("assinador")
	if !c.matchesRequestedInstance(state) {
		state = process.State{Name: "assinador", Port: c.Port, URL: c.baseURL()}
	}
	return state, nil
}

func (c Client) StartServer(idleMinutes int) (process.State, bool, error) {
	if err := c.Health(); err == nil {
		state, _ := process.Load("assinador")
		if !c.matchesRequestedInstance(state) {
			state = process.State{Name: "assinador", Port: c.Port, URL: c.baseURL(), StartedAt: time.Now()}
		}
		return state, true, nil
	}
	if !netutil.IsTCPPortFree(c.Port) {
		return process.State{}, false, fmt.Errorf("porta %d já está em uso por outro processo. Use outra porta, por exemplo --port 8080 ou --port 18080", c.Port)
	}
	if c.JarPath == "" {
		return process.State{}, false, errors.New("assinador.jar não informado. Use --jar ou RUNNER_ASSINADOR_JAR")
	}
	if _, err := os.Stat(c.JarPath); err != nil {
		return process.State{}, false, fmt.Errorf("assinador.jar não encontrado em %s", c.JarPath)
	}
	java, err := jdk.EnsureJava21(context.Background(), c.ManifestURL)
	if err != nil {
		return process.State{}, false, err
	}
	if err := paths.EnsureRuntimeDirs(); err != nil {
		return process.State{}, false, err
	}
	logDir, err := paths.LogDir()
	if err != nil {
		return process.State{}, false, err
	}
	stdoutPath := filepath.Join(logDir, "assinador.out.log")
	stderrPath := filepath.Join(logDir, "assinador.err.log")
	stdout, err := os.OpenFile(stdoutPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return process.State{}, false, err
	}
	stderr, err := os.OpenFile(stderrPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		_ = stdout.Close()
		return process.State{}, false, err
	}
	args := []string{"-jar", c.JarPath, "server", "--port", fmt.Sprint(c.Port)}
	if idleMinutes > 0 {
		args = append(args, "--idle-timeout-minutes", fmt.Sprint(idleMinutes))
	}
	cmd := exec.Command(java, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Start(); err != nil {
		_ = stdout.Close()
		_ = stderr.Close()
		return process.State{}, false, err
	}
	_ = stdout.Close()
	_ = stderr.Close()

	exitCh := make(chan error, 1)
	go func() { exitCh <- cmd.Wait() }()

	state := process.State{Name: "assinador", PID: cmd.Process.Pid, Port: c.Port, URL: c.baseURL(), JarPath: c.JarPath, StartedAt: time.Now()}
	if err := process.Save(state); err != nil {
		_ = cmd.Process.Kill()
		return state, false, err
	}
	deadline := time.Now().Add(15 * time.Second)
	var lastHealthErr error
	for time.Now().Before(deadline) {
		select {
		case err := <-exitCh:
			_ = process.Delete("assinador")
			return state, false, fmt.Errorf("processo do assinador.jar encerrou antes de ficar pronto: %v%s", err, logHint(stderrPath))
		default:
		}
		if err := c.Health(); err == nil {
			return state, false, nil
		} else {
			lastHealthErr = err
		}
		time.Sleep(300 * time.Millisecond)
	}
	return state, false, fmt.Errorf("servidor do assinador foi iniciado, mas não ficou pronto no tempo limite; última verificação de saúde: %v%s", lastHealthErr, logHint(stderrPath))
}

func (c Client) StopServer() error {
	client := httpx.New(5*time.Second, false)
	_, _, err := httpx.PostEmpty(client, c.baseURL()+"/shutdown")
	_ = process.Delete("assinador")
	return err
}

func (c Client) SignHTTP(req SignRequest) ([]byte, error) {
	client := httpx.New(c.Timeout, false)
	out, _, err := httpx.PostJSON(client, c.baseURL()+"/sign", req)
	if err != nil {
		return out, err
	}
	return out, nil
}

func (c Client) ValidateHTTP(req ValidateRequest) ([]byte, error) {
	client := httpx.New(c.Timeout, false)
	out, _, err := httpx.PostJSON(client, c.baseURL()+"/validate", req)
	if err != nil {
		return out, err
	}
	return out, nil
}

func (c Client) SignLocal(req SignRequest) ([]byte, []byte, int, error) {
	args := []string{"sign"}
	appendFlag := func(name, value string) {
		if value != "" {
			args = append(args, "--"+name, value)
		}
	}
	appendFlag("bundle", req.Bundle)
	appendFlag("provenance", req.Provenance)
	appendFlag("crypto-material", req.CryptoMaterial)
	appendFlag("cert-chain", req.CertificateChain)
	appendFlag("timestamp", req.Timestamp)
	appendFlag("strategy", req.Strategy)
	appendFlag("policy", req.Policy)
	appendFlag("config", req.Config)
	appendFlag("signer", req.Signer)
	appendFlag("input", req.Input)
	return c.runLocal(args)
}

func (c Client) ValidateLocal(req ValidateRequest) ([]byte, []byte, int, error) {
	args := []string{"validate"}
	appendFlag := func(name, value string) {
		if value != "" {
			args = append(args, "--"+name, value)
		}
	}
	appendFlag("signature", req.Signature)
	appendFlag("timestamp", req.Timestamp)
	appendFlag("policy", req.Policy)
	appendFlag("config", req.Config)
	appendFlag("bundle", req.Bundle)
	appendFlag("provenance", req.Provenance)
	appendFlag("input", req.Input)
	return c.runLocal(args)
}

func (c Client) runLocal(args []string) ([]byte, []byte, int, error) {
	if c.JarPath == "" {
		return nil, []byte("assinador.jar não informado. Use --jar ou RUNNER_ASSINADOR_JAR\n"), 2, errors.New("jar ausente")
	}
	if _, err := os.Stat(c.JarPath); err != nil {
		msg := fmt.Sprintf("assinador.jar não encontrado em %s. Gere com 'make java-test' ou informe --jar com o caminho correto.\n", c.JarPath)
		return nil, []byte(msg), 2, fmt.Errorf("assinador.jar não encontrado: %w", err)
	}
	java, err := jdk.EnsureJava21(context.Background(), c.ManifestURL)
	if err != nil {
		return nil, []byte(err.Error() + "\n"), 2, err
	}
	// Invocação local equivalente a: java -jar assinador.jar <comando> <flags>
	allArgs := append([]string{"-jar", c.JarPath}, args...)
	cmd := exec.Command(java, allArgs...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	code := 0
	if err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			code = ee.ExitCode()
		} else {
			code = 3
		}
	}
	return stdout.Bytes(), stderr.Bytes(), code, err
}

func WriteOutput(out []byte, outputPath string, stdout io.Writer) error {
	if outputPath == "" {
		_, err := stdout.Write(append(out, '\n'))
		return err
	}
	if err := os.WriteFile(outputPath, out, 0o644); err != nil {
		return err
	}
	_, err := fmt.Fprintf(stdout, "Resultado gravado em %s\n", outputPath)
	return err
}

func logHint(path string) string {
	b, err := os.ReadFile(path)
	if err != nil || len(b) == 0 {
		return ""
	}
	text := strings.TrimSpace(string(b))
	if text == "" {
		return ""
	}
	if len(text) > 800 {
		text = text[len(text)-800:]
		if idx := strings.IndexByte(text, '\n'); idx >= 0 && idx+1 < len(text) {
			text = text[idx+1:]
		}
	}
	return "\nÚltimas linhas do log de erro: " + text
}
