package assinador

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/kyriosdata/runner/internal/httpx"
	"github.com/kyriosdata/runner/internal/jdk"
	"github.com/kyriosdata/runner/internal/paths"
	"github.com/kyriosdata/runner/internal/process"
)

const DefaultPort = 8080

type Client struct {
	Port    int
	JarPath string
	Timeout time.Duration
}

type SignRequest struct {
	Input  string `json:"input"`
	Signer string `json:"signer"`
}

type ValidateRequest struct {
	Signature string `json:"signature"`
	Input     string `json:"input,omitempty"`
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

func (c Client) baseURL() string { return fmt.Sprintf("http://localhost:%d", c.Port) }

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
	if state.Name == "" {
		state = process.State{Name: "assinador", Port: c.Port, URL: c.baseURL()}
	}
	return state, nil
}

func (c Client) StartServer(idleMinutes int) (process.State, bool, error) {
	if err := c.Health(); err == nil {
		state, _ := process.Load("assinador")
		if state.Name == "" {
			state = process.State{Name: "assinador", Port: c.Port, URL: c.baseURL(), StartedAt: time.Now()}
		}
		return state, true, nil
	}
	if c.JarPath == "" {
		return process.State{}, false, errors.New("assinador.jar não informado. Use --jar ou RUNNER_ASSINADOR_JAR")
	}
	if _, err := os.Stat(c.JarPath); err != nil {
		return process.State{}, false, fmt.Errorf("assinador.jar não encontrado em %s", c.JarPath)
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
	stdout, err := os.OpenFile(filepath.Join(logDir, "assinador.out.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return process.State{}, false, err
	}
	stderr, err := os.OpenFile(filepath.Join(logDir, "assinador.err.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
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
	state := process.State{Name: "assinador", PID: cmd.Process.Pid, Port: c.Port, URL: c.baseURL(), JarPath: c.JarPath, StartedAt: time.Now()}
	if err := process.Save(state); err != nil {
		return state, false, err
	}
	_ = cmd.Process.Release()
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		if err := c.Health(); err == nil {
			return state, false, nil
		}
		time.Sleep(300 * time.Millisecond)
	}
	return state, false, errors.New("servidor do assinador foi iniciado, mas não ficou pronto no tempo limite")
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
	args := []string{"sign", "--input", req.Input, "--signer", req.Signer}
	return c.runLocal(args)
}

func (c Client) ValidateLocal(req ValidateRequest) ([]byte, []byte, int, error) {
	args := []string{"validate", "--signature", req.Signature}
	if req.Input != "" {
		args = append(args, "--input", req.Input)
	}
	return c.runLocal(args)
}

func (c Client) runLocal(args []string) ([]byte, []byte, int, error) {
	if c.JarPath == "" {
		return nil, []byte("assinador.jar não informado. Use --jar ou RUNNER_ASSINADOR_JAR\n"), 2, errors.New("jar ausente")
	}
	java, err := jdk.FindJava()
	if err != nil {
		return nil, []byte(err.Error() + "\n"), 2, err
	}
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
