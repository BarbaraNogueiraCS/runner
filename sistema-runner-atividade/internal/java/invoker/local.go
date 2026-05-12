package invoker

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	assinaturadto "github.com/kyriosdata/assinatura/internal/assinatura/dto"
	javaruntime "github.com/kyriosdata/assinatura/internal/java/runtime"
)

type LocalJarInvoker struct {
	RuntimeResolver javaruntime.Resolver
	JarPath         string
}

func NewLocalJarInvoker(jarPath string) LocalJarInvoker {
	return LocalJarInvoker{RuntimeResolver: javaruntime.NewResolver(), JarPath: jarPath}
}

func (i LocalJarInvoker) Sign(command assinaturadto.SignCommand) (assinaturadto.OperationResult, error) {
	return i.invoke("sign", []string{"--documento", command.Documento, "--certificado", command.Certificado})
}

func (i LocalJarInvoker) Validate(command assinaturadto.ValidateCommand) (assinaturadto.OperationResult, error) {
	return i.invoke("validate", []string{"--documento", command.Documento, "--assinatura", command.Assinatura})
}

func (i LocalJarInvoker) invoke(operation string, args []string) (assinaturadto.OperationResult, error) {
	if _, err := os.Stat(i.JarPath); err != nil {
		return assinaturadto.OperationResult{}, fmt.Errorf("arquivo jar não encontrado em %s", i.JarPath)
	}
	javaRt, err := i.RuntimeResolver.Resolve()
	if err != nil {
		return assinaturadto.OperationResult{}, err
	}
	fullArgs := append([]string{"-jar", i.JarPath, operation}, args...)
	cmd := exec.Command(javaRt.JavaPath, fullArgs...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return assinaturadto.OperationResult{}, err
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		if err != nil {
			return assinaturadto.OperationResult{}, fmt.Errorf("falha ao executar jar: %s %s", err.Error(), stderr.String())
		}
	case <-time.After(30 * time.Second):
		_ = cmd.Process.Kill()
		return assinaturadto.OperationResult{}, fmt.Errorf("tempo limite excedido ao executar jar")
	}
	raw := strings.TrimSpace(stdout.String())
	return assinaturadto.OperationResult{Success: true, Operation: operation, Message: "Operação executada pelo assinador.jar", ExecutionMode: "local", Raw: raw}, nil
}
