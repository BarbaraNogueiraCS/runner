package process

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/kyriosdata/assinatura/internal/config"
)

type Manager struct {
	Registry ProcessRegistry
}

func NewManager() Manager {
	return Manager{Registry: NewProcessRegistry()}
}

func (m Manager) StartJavaServer(application, javaPath, jarPath string, port int, extraArgs ...string) (ProcessMetadata, error) {
	if !IsPortAvailable(port) {
		return ProcessMetadata{}, fmt.Errorf("porta %d indisponível", port)
	}
	if err := config.EnsureBaseDirs(); err != nil {
		return ProcessMetadata{}, err
	}
	logPath := filepath.Join(config.LogsDir(), fmt.Sprintf("%s-%d.log", application, port))
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return ProcessMetadata{}, err
	}

	args := []string{"-jar", jarPath, "server", "--port", fmt.Sprintf("%d", port)}
	args = append(args, extraArgs...)
	cmd := exec.Command(javaPath, args...)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	if err := cmd.Start(); err != nil {
		_ = logFile.Close()
		return ProcessMetadata{}, err
	}
	metadata := ProcessMetadata{
		Application:      application,
		PID:              cmd.Process.Pid,
		Port:             port,
		Mode:             "server",
		StartedAt:        time.Now(),
		HealthEndpoint:   fmt.Sprintf("http://localhost:%d/health", port),
		ShutdownEndpoint: fmt.Sprintf("http://localhost:%d/shutdown", port),
		Status:           "running",
	}
	if application == "simulador" {
		metadata.InfoEndpoint = fmt.Sprintf("http://localhost:%d/api/info", port)
	}
	if err := m.Registry.Save(metadata); err != nil {
		return ProcessMetadata{}, err
	}
	return metadata, nil
}

func (m Manager) Stop(application string, port int) error {
	metadata, err := m.Registry.Find(application, port)
	if err != nil {
		return err
	}
	if metadata.ShutdownEndpoint != "" {
		_ = callShutdown(metadata.ShutdownEndpoint)
		time.Sleep(300 * time.Millisecond)
	}
	proc, err := os.FindProcess(metadata.PID)
	if err == nil && proc != nil {
		if runtime.GOOS == "windows" {
			_ = proc.Kill()
		} else {
			_ = proc.Signal(os.Interrupt)
			time.Sleep(300 * time.Millisecond)
			_ = proc.Kill()
		}
	}
	return m.Registry.Remove(application, port)
}

func callShutdown(url string) error {
	client := httpClient()
	resp, err := client.Post(url, "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func httpClient() *http.Client {
	return &http.Client{Timeout: 2 * time.Second}
}
