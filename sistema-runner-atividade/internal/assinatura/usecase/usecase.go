package usecase

import (
	"fmt"
	"time"

	"github.com/kyriosdata/assinatura/internal/assinatura/dto"
	"github.com/kyriosdata/assinatura/internal/config"
	apperrors "github.com/kyriosdata/assinatura/internal/errors"
	"github.com/kyriosdata/assinatura/internal/java/invoker"
	javaruntime "github.com/kyriosdata/assinatura/internal/java/runtime"
	"github.com/kyriosdata/assinatura/internal/process"
)

type SignUseCase struct{}
type ValidateUseCase struct{}
type StartAssinadorUseCase struct{}
type StopAssinadorUseCase struct{}

func NewSignUseCase() SignUseCase                     { return SignUseCase{} }
func NewValidateUseCase() ValidateUseCase             { return ValidateUseCase{} }
func NewStartAssinadorUseCase() StartAssinadorUseCase { return StartAssinadorUseCase{} }
func NewStopAssinadorUseCase() StopAssinadorUseCase   { return StopAssinadorUseCase{} }

func (uc SignUseCase) Execute(command dto.SignCommand) (dto.OperationResult, error) {
	if command.Documento == "" {
		return dto.OperationResult{}, apperrors.New(apperrors.InvalidParameter, "parâmetro obrigatório ausente", "informe --documento")
	}
	if command.Certificado == "" {
		return dto.OperationResult{}, apperrors.New(apperrors.InvalidParameter, "parâmetro obrigatório ausente", "informe --certificado")
	}
	jarPath := command.JarPath
	if jarPath == "" {
		jarPath = config.AssinadorJarPath()
	}
	if !command.Local && isAssinadorServerHealthy(command.Port) {
		return invoker.NewAssinadorHTTPClient(command.Port).Sign(command)
	}
	return invoker.NewLocalJarInvoker(jarPath).Sign(command)
}

func (uc ValidateUseCase) Execute(command dto.ValidateCommand) (dto.OperationResult, error) {
	if command.Documento == "" {
		return dto.OperationResult{}, apperrors.New(apperrors.InvalidParameter, "parâmetro obrigatório ausente", "informe --documento")
	}
	if command.Assinatura == "" {
		return dto.OperationResult{}, apperrors.New(apperrors.InvalidParameter, "parâmetro obrigatório ausente", "informe --assinatura")
	}
	jarPath := command.JarPath
	if jarPath == "" {
		jarPath = config.AssinadorJarPath()
	}
	if !command.Local && isAssinadorServerHealthy(command.Port) {
		return invoker.NewAssinadorHTTPClient(command.Port).Validate(command)
	}
	return invoker.NewLocalJarInvoker(jarPath).Validate(command)
}

func (uc StartAssinadorUseCase) Execute(jarPath string, port int) (process.ProcessMetadata, error) {
	if jarPath == "" {
		jarPath = config.AssinadorJarPath()
	}
	if !process.IsPortAvailable(port) {
		if isAssinadorServerHealthy(port) {
			metadata, _ := process.NewProcessRegistry().Find("assinador", port)
			return metadata, nil
		}
		return process.ProcessMetadata{}, apperrors.New(apperrors.PortUnavailable, "porta indisponível", fmt.Sprintf("a porta %d já está em uso", port))
	}
	javaRt, err := javaruntime.NewResolver().Resolve()
	if err != nil {
		return process.ProcessMetadata{}, err
	}
	return process.NewManager().StartJavaServer("assinador", javaRt.JavaPath, jarPath, port)
}

func (uc StopAssinadorUseCase) Execute(port int) error {
	return process.NewManager().Stop("assinador", port)
}

func isAssinadorServerHealthy(port int) bool {
	return process.IsHTTPHealthy(fmt.Sprintf("http://localhost:%d/health", port), 700*time.Millisecond)
}
