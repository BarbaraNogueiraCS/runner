package usecase

import (
	"fmt"
	"os"
	"time"

	"github.com/kyriosdata/assinatura/internal/artifacts"
	"github.com/kyriosdata/assinatura/internal/config"
	apperrors "github.com/kyriosdata/assinatura/internal/errors"
	javaruntime "github.com/kyriosdata/assinatura/internal/java/runtime"
	"github.com/kyriosdata/assinatura/internal/process"
	"github.com/kyriosdata/assinatura/internal/simulador/dto"
)

type StartSimulatorUseCase struct{}
type StopSimulatorUseCase struct{}
type StatusSimulatorUseCase struct{}

func NewStartSimulatorUseCase() StartSimulatorUseCase   { return StartSimulatorUseCase{} }
func NewStopSimulatorUseCase() StopSimulatorUseCase     { return StopSimulatorUseCase{} }
func NewStatusSimulatorUseCase() StatusSimulatorUseCase { return StatusSimulatorUseCase{} }

func (uc StartSimulatorUseCase) Execute(cfg dto.SimulatorConfig) (dto.SimulatorStatus, error) {
	if cfg.Port == 0 {
		cfg.Port = config.DefaultSimulatorPort
	}
	jarPath := cfg.JarPath
	if jarPath == "" {
		jarPath = config.SimulatorJarPath()
	}
	if _, err := os.Stat(jarPath); err != nil {
		if cfg.Source == "" {
			return dto.SimulatorStatus{}, apperrors.New(apperrors.JarNotFound, "simulador.jar não encontrado", "informe --jar ou --source para download")
		}
		if err := artifacts.DownloadFile(cfg.Source, jarPath); err != nil {
			return dto.SimulatorStatus{}, apperrors.Wrap(apperrors.DownloadFailed, "falha ao baixar simulador.jar", err)
		}
		if err := artifacts.VerifySHA256(jarPath, cfg.SHA256); err != nil {
			return dto.SimulatorStatus{}, apperrors.Wrap(apperrors.ChecksumMismatch, "falha na verificação de integridade", err)
		}
	}
	if !process.IsPortAvailable(cfg.Port) {
		status, _ := uc.statusFromRegistry(cfg.Port)
		if status.Running {
			return status, nil
		}
		return dto.SimulatorStatus{}, apperrors.New(apperrors.PortUnavailable, "porta indisponível", fmt.Sprintf("a porta %d já está em uso", cfg.Port))
	}
	javaRt, err := javaruntime.NewResolver().Resolve()
	if err != nil {
		return dto.SimulatorStatus{}, err
	}
	metadata, err := process.NewManager().StartJavaServer("simulador", javaRt.JavaPath, jarPath, cfg.Port)
	if err != nil {
		return dto.SimulatorStatus{}, err
	}
	return dto.SimulatorStatus{Running: true, Port: metadata.Port, PID: metadata.PID, Message: "Simulador iniciado"}, nil
}

func (uc StartSimulatorUseCase) statusFromRegistry(port int) (dto.SimulatorStatus, error) {
	metadata, err := process.NewProcessRegistry().Find("simulador", port)
	if err != nil {
		return dto.SimulatorStatus{}, err
	}
	running := false
	if metadata.InfoEndpoint != "" {
		running = process.IsHTTPHealthy(metadata.InfoEndpoint, 700*time.Millisecond)
	}
	return dto.SimulatorStatus{Running: running, Port: metadata.Port, PID: metadata.PID, Message: "Status obtido do registro local"}, nil
}

func (uc StopSimulatorUseCase) Execute(port int) error {
	if port == 0 {
		port = config.DefaultSimulatorPort
	}
	return process.NewManager().Stop("simulador", port)
}

func (uc StatusSimulatorUseCase) Execute(port int) (dto.SimulatorStatus, error) {
	if port == 0 {
		port = config.DefaultSimulatorPort
	}
	metadata, err := process.NewProcessRegistry().Find("simulador", port)
	if err != nil {
		return dto.SimulatorStatus{Running: false, Port: port, Message: "Simulador não está registrado"}, nil
	}
	running := false
	if metadata.InfoEndpoint != "" {
		running = process.IsHTTPHealthy(metadata.InfoEndpoint, 700*time.Millisecond)
	}
	msg := "Simulador não respondeu ao health check"
	if running {
		msg = "Simulador em execução"
	}
	return dto.SimulatorStatus{Running: running, Port: metadata.Port, PID: metadata.PID, Message: msg}, nil
}
