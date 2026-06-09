package simulator

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/kyriosdata/runner/internal/httpx"
	"github.com/kyriosdata/runner/internal/jdk"
	"github.com/kyriosdata/runner/internal/paths"
	"github.com/kyriosdata/runner/internal/process"
)

const DefaultURL = "https://localhost:8443"

type Client struct {
	BaseURL  string
	Insecure bool
	Timeout  time.Duration
	JarPath  string
}

func New(baseURL string, insecure bool, jar string) Client {
	if baseURL == "" {
		baseURL = os.Getenv("HUBSAUDE_SIMULATOR_URL")
	}
	if baseURL == "" {
		baseURL = DefaultURL
	}
	return Client{BaseURL: strings.TrimRight(baseURL, "/"), Insecure: insecure, Timeout: 5 * time.Second, JarPath: jar}
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
	if state.Name == "" {
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
	if c.JarPath == "" {
		return process.State{}, false, errors.New("simulador não está acessível e nenhum simulador.jar foi informado. Use --jar para iniciar um arquivo local")
	}
	if _, err := os.Stat(c.JarPath); err != nil {
		return process.State{}, false, fmt.Errorf("simulador.jar não encontrado em %s", c.JarPath)
	}
	java, err := jdk.FindJava()
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
	args := []string{"-jar", c.JarPath}
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
	state := process.State{Name: "simulador", PID: cmd.Process.Pid, Port: c.port(), URL: c.BaseURL, JarPath: c.JarPath, StartedAt: time.Now()}
	if err := process.Save(state); err != nil {
		return state, false, err
	}
	_ = cmd.Process.Release()
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		if _, _, err := c.Status(); err == nil {
			return state, false, nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return state, false, errors.New("simulador foi iniciado, mas /api/info não respondeu no tempo limite")
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
