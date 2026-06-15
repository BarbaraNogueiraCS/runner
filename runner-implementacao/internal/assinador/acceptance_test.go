package assinador

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAcceptanceDefaultServerModeUsesHTTPPort(t *testing.T) {
	c := New(0, "")
	if c.Port != DefaultPort {
		t.Fatalf("porta padrão esperada %d, obtida %d", DefaultPort, c.Port)
	}
	if got := c.baseURL(); got != "http://127.0.0.1:8080" {
		t.Fatalf("baseURL padrão inesperada: %s", got)
	}
}

func TestAcceptanceHTTPValidateInvokesValidateEndpoint(t *testing.T) {
	var gotPath, gotMethod string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"resourceType":"OperationOutcome","valid":true}`))
	}))
	defer srv.Close()

	c := New(portFromURL(t, srv.URL), "")
	out, err := c.ValidateHTTP(ValidateRequest{Signature: "assinatura.json"})
	if err != nil {
		t.Fatalf("ValidateHTTP retornou erro: %v", err)
	}
	if gotPath != "/validate" || gotMethod != http.MethodPost {
		t.Fatalf("esperado POST /validate, obtido %s %s", gotMethod, gotPath)
	}
	if !strings.Contains(string(out), `"valid":true`) {
		t.Fatalf("resposta inesperada: %s", string(out))
	}
}

func TestAcceptanceLocalInvocationFailsClearlyWithoutJar(t *testing.T) {
	c := New(0, "")
	_, stderr, code, err := c.SignLocal(SignRequest{Input: "documento.txt"})
	if err == nil {
		t.Fatalf("SignLocal sem jar deveria falhar")
	}
	if code != 2 {
		t.Fatalf("código esperado 2 para dependência ausente, obtido %d", code)
	}
	if !strings.Contains(string(stderr), "assinador.jar não informado") {
		t.Fatalf("erro deveria orientar sobre jar ausente, obtido: %s", string(stderr))
	}
}
