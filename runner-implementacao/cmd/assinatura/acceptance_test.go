package main

import (
	"testing"

	"github.com/kyriosdata/runner/internal/apperrors"
)

func TestAcceptanceSignCommandRequiresGuideParametersOrLegacyInput(t *testing.T) {
	code := run([]string{"sign", "--bundle", "bundle.json"})
	if code != apperrors.UsageError {
		t.Fatalf("sign sem parâmetros obrigatórios deveria retornar UsageError=%d, obtido %d", apperrors.UsageError, code)
	}
}

func TestAcceptanceValidateCommandRequiresSignature(t *testing.T) {
	code := run([]string{"validate", "--timestamp", "1751328000"})
	if code != apperrors.UsageError {
		t.Fatalf("validate sem --signature deveria retornar UsageError=%d, obtido %d", apperrors.UsageError, code)
	}
}

func TestAcceptanceServerCommandRejectsMissingOrInvalidAction(t *testing.T) {
	if code := run([]string{"server"}); code != apperrors.UsageError {
		t.Fatalf("server sem ação deveria retornar UsageError, obtido %d", code)
	}
	if code := run([]string{"server", "restart"}); code != apperrors.UsageError {
		t.Fatalf("server com ação inválida deveria retornar UsageError, obtido %d", code)
	}
}

func TestAcceptanceDefaultPortIsUsedWhenNotInformed(t *testing.T) {
	f, err := parseFlags([]string{"--signature", "assinatura.json"})
	if err != nil {
		t.Fatalf("parseFlags retornou erro: %v", err)
	}
	if got := intFlag(f, "port", 8080); got != 8080 {
		t.Fatalf("porta padrão esperada 8080, obtida %d", got)
	}
}
