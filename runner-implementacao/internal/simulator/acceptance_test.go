package simulator

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAcceptanceStartRejectsBusyPortThatIsNotSimulator(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	defer srv.Close()

	c := New(srv.URL, false, "")
	_, reused, err := c.Start()
	if err == nil {
		t.Fatalf("Start deveria falhar quando porta está ocupada e /api/info não responde")
	}
	if reused {
		t.Fatalf("serviço não identificado como simulador não deve ser reutilizado")
	}
	if !strings.Contains(err.Error(), "porta") || !strings.Contains(err.Error(), "/api/info") {
		t.Fatalf("erro deveria explicar porta ocupada e /api/info, obtido: %v", err)
	}
}

func TestAcceptanceArtifactCanBeSelectedByEnvironment(t *testing.T) {
	t.Setenv("RUNNER_SIMULADOR_ARTIFACT", "validador")
	c := New("http://127.0.0.1:8443", false, "")
	if c.Artifact != "validador" {
		t.Fatalf("artefato esperado validador, obtido %s", c.Artifact)
	}
}
