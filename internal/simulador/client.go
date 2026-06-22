package simulador

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/BarbaraNogueiraCS/runner/internal/httpx"
	"github.com/BarbaraNogueiraCS/runner/internal/jdk"
	"github.com/BarbaraNogueiraCS/runner/internal/netutil"
	"github.com/BarbaraNogueiraCS/runner/internal/paths"
	"github.com/BarbaraNogueiraCS/runner/internal/process"
	"github.com/BarbaraNogueiraCS/runner/internal/release"
)

const DefaultURL = "https://localhost:8443"
const DefaultArtifact = "simulador"

type Client struct {
	BaseURL      string
	Insecure     bool
	Timeout      time.Duration
	JarPath      string
	Artifact     string
	ManifestURL  string
	SourceURL    string
	SourceSHA256 string
}

func New(baseURL string, insecure bool, jar string) Client {
	if baseURL == "" {
		baseURL = os.Getenv("HUBSAUDE_SIMULADOR_URL")
	}
	if baseURL == "" {
		baseURL = DefaultURL
	}
	artifact := os.Getenv("RUNNER_SIMULADOR_ARTIFACT")
	if artifact == "" {
		artifact = DefaultArtifact
	}
	return Client{
		BaseURL:      strings.TrimRight(baseURL, "/"),
		Insecure:     insecure,
		Timeout:      5 * time.Second,
		JarPath:      jar,
		Artifact:     artifact,
		ManifestURL:  os.Getenv("RUNNER_RELEASE_JSON"),
		SourceURL:    os.Getenv("RUNNER_SIMULADOR_SOURCE"),
		SourceSHA256: os.Getenv("RUNNER_SIMULADOR_SHA256"),
	}
}

func (c Client) Info() ([]byte, error) {
	client := httpx.New(c.Timeout, c.Insecure)
	b, _, err := httpx.Get(client, c.BaseURL+"/api/info")
	return b, err
}

func (c Client) Status() (process.State, []byte, error) {
	info, err := c.Info()
	if err != nil {
		return process.State{}, nil, err
	}
	state, _ := process.Load("simulador")
	if state.Name == "" || state.Port != c.port() {
		state = process.State{Name: "simulador", URL: c.BaseURL, Port: c.port()}
	}
	return state, info, nil
}

func (c Client) Stop() ([]byte, error) {
	client := httpx.New(c.Timeout, c.Insecure)
	b, _, err := httpx.PostEmpty(client, c.BaseURL+"/shutdown")
	_ = process.Delete("simulador")
	return b, err
}

func (c Client) Start() (process.State, bool, error) {
	if state, _, err := c.Status(); err == nil {
		return state, true, nil
	}
	port := c.port()
	if port > 0 && !netutil.IsTCPPortFree(port) {
		return process.State{}, false, fmt.Errorf("porta %d já está em uso por outro processo, mas %s/api/info não respondeu como Simulador HubSaúde", port, c.BaseURL)
	}

	jarPath := c.JarPath
	downloaded := false
	if jarPath == "" {
		var err error
		if strings.TrimSpace(c.SourceURL) != "" {
			jarPath, _, downloaded, err = release.EnsureArtifactFromURL(context.Background(), c.SourceURL, c.SourceSHA256, c.Artifact)
			if err != nil {
				return process.State{}, false, fmt.Errorf("não foi possível obter o simulador.jar a partir de --source: %w", err)
			}
		} else {
			jarPath, _, downloaded, err = release.EnsureArtifact(context.Background(), c.ManifestURL, c.Artifact)
			if err != nil {
				return process.State{}, false, fmt.Errorf("não foi possível obter dinamicamente o simulador.jar a partir do release.json/GitHub Releases: %w", err)
			}
		}
	}
	if jarPath == "" {
		return process.State{}, false, errors.New("simulador.jar não definido")
	}
	if _, err := os.Stat(jarPath); err != nil {
		return process.State{}, false, fmt.Errorf("simulador.jar não encontrado em %s", jarPath)
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
	stdout, err := os.OpenFile(filepath.Join(logDir, "simulador.out.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return process.State{}, false, err
	}
	stderr, err := os.OpenFile(filepath.Join(logDir, "simulador.err.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		_ = stdout.Close()
		return process.State{}, false, err
	}
	args := []string{"-jar", jarPath}
	if p := c.port(); p > 0 {
		args = append(args, "--server.port="+strconv.Itoa(p))
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
	state := process.State{Name: "simulador", PID: cmd.Process.Pid, Port: c.port(), URL: c.BaseURL, JarPath: jarPath, StartedAt: time.Now()}
	if err := process.Save(state); err != nil {
		_ = cmd.Process.Kill()
		return state, false, err
	}
	deadline := time.Now().Add(30 * time.Second)
	var lastErr error
	for time.Now().Before(deadline) {
		select {
		case err := <-exitCh:
			_ = process.Delete("simulador")
			return state, false, fmt.Errorf("processo do simulador.jar encerrou antes de ficar pronto: %v%s", err, logHint(filepath.Join(logDir, "simulador.err.log")))
		default:
		}
		if _, _, err := c.Status(); err == nil {
			if downloaded {
				// Mantém o dado no estado para rastreabilidade do artefato obtido por release.json.
				state.JarPath = jarPath
				_ = process.Save(state)
			}
			return state, false, nil
		} else {
			lastErr = err
		}
		time.Sleep(500 * time.Millisecond)
	}
	return state, false, fmt.Errorf("simulador foi iniciado, mas /api/info não respondeu no tempo limite; última verificação: %v%s", lastErr, logHint(filepath.Join(logDir, "simulador.err.log")))
}

func (c Client) port() int {
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return 0
	}
	p := u.Port()
	if p == "" {
		if u.Scheme == "https" {
			return 443
		}
		if u.Scheme == "http" {
			return 80
		}
		return 0
	}
	v, _ := strconv.Atoi(p)
	return v
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
