package formatter

import (
	"fmt"
	"strings"

	apperrors "github.com/kyriosdata/assinatura/internal/errors"
)

type Result struct {
	Success       bool
	Operation     string
	Message       string
	Signature     string
	Valid         *bool
	ExecutionMode string
	Raw           string
}

func FormatResult(r Result) string {
	var b strings.Builder
	if r.Operation != "" {
		fmt.Fprintf(&b, "Operação: %s\n", r.Operation)
	}
	if r.Success {
		b.WriteString("Status: sucesso\n")
	} else {
		b.WriteString("Status: falha\n")
	}
	if r.Message != "" {
		fmt.Fprintf(&b, "Mensagem: %s\n", r.Message)
	}
	if r.Signature != "" {
		fmt.Fprintf(&b, "Assinatura: %s\n", r.Signature)
	}
	if r.Valid != nil {
		if *r.Valid {
			b.WriteString("Resultado: assinatura válida\n")
		} else {
			b.WriteString("Resultado: assinatura inválida\n")
		}
	}
	if r.ExecutionMode != "" {
		fmt.Fprintf(&b, "Modo de execução: %s\n", r.ExecutionMode)
	}
	if r.Raw != "" {
		fmt.Fprintf(&b, "Resposta bruta: %s\n", r.Raw)
	}
	return strings.TrimRight(b.String(), "\n")
}

func FormatError(err error) string {
	if err == nil {
		return ""
	}
	if appErr, ok := err.(apperrors.AppError); ok {
		var b strings.Builder
		fmt.Fprintf(&b, "Erro: %s\n", appErr.Message)
		if appErr.Details != "" {
			fmt.Fprintf(&b, "Motivo: %s\n", appErr.Details)
		}
		fmt.Fprintf(&b, "Código: %s", appErr.Code)
		return b.String()
	}
	return "Erro: " + err.Error()
}
