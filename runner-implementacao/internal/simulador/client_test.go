package simulador

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestNewDefaultURL(t *testing.T) {
	c := New("", false, "")
	if c.BaseURL != DefaultURL {
		t.Fatalf("URL esperada %s, obtida %s", DefaultURL, c.BaseURL)
	}
}

func TestStatusUsesApiInfo(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/info" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":"0.1.7"}`))
	}))
	defer srv.Close()

	c := New(srv.URL, false, "")
	state, info, err := c.Status()
	if err != nil {
		t.Fatalf("Status retornou erro: %v", err)
	}
	if state.URL != srv.URL || state.Port == 0 {
		t.Fatalf("estado inesperado: %#v", state)
	}
	if !strings.Contains(string(info), "0.1.7") {
		t.Fatalf("info inesperado: %s", string(info))
	}
}

func TestStartReusesRunningSimulador(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/info" {
			_, _ = w.Write([]byte(`{"version":"0.1.7"}`))
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	c := New(srv.URL, false, "")
	state, reused, err := c.Start()
	if err != nil {
		t.Fatalf("Start deveria reutilizar simulador vivo: %v", err)
	}
	if !reused {
		t.Fatalf("Start deveria retornar reused=true")
	}
	if state.Port != portFromURL(t, srv.URL) {
		t.Fatalf("porta inesperada: %#v", state)
	}
}

func TestStopCallsShutdownEndpoint(t *testing.T) {
	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/shutdown" && r.Method == http.MethodPost {
			called = true
			_, _ = w.Write([]byte(`{"status":"SHUTTING_DOWN"}`))
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	c := New(srv.URL, false, "")
	out, err := c.Stop()
	if err != nil {
		t.Fatalf("Stop retornou erro: %v", err)
	}
	if !called || !strings.Contains(string(out), "SHUTTING_DOWN") {
		t.Fatalf("/shutdown não foi chamado corretamente")
	}
}

func portFromURL(t *testing.T, raw string) int {
	t.Helper()
	idx := strings.LastIndex(raw, ":")
	if idx < 0 {
		t.Fatalf("URL sem porta: %s", raw)
	}
	port, err := strconv.Atoi(raw[idx+1:])
	if err != nil {
		t.Fatalf("porta inválida em %s: %v", raw, err)
	}
	return port
}
