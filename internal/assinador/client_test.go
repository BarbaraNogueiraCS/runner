package assinador

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestNewUsesDefaultPort(t *testing.T) {
	c := New(0, "")
	if c.Port != DefaultPort {
		t.Fatalf("porta esperada %d, obtida %d", DefaultPort, c.Port)
	}
}

func TestStartServerReusesHealthyInstance(t *testing.T) {
	port, closeFn := startAssinadorHealthServer(t)
	defer closeFn()

	c := New(port, "")
	state, reused, err := c.StartServer(0)
	if err != nil {
		t.Fatalf("StartServer deveria reutilizar instância saudável: %v", err)
	}
	if !reused {
		t.Fatalf("StartServer deveria retornar reused=true")
	}
	if state.Port != port || !strings.Contains(state.URL, strconv.Itoa(port)) {
		t.Fatalf("estado inesperado: %#v", state)
	}
}

func TestStartServerRejectsOccupiedPortWithoutHealth(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	defer srv.Close()
	port := portFromURL(t, srv.URL)

	c := New(port, "")
	_, reused, err := c.StartServer(0)
	if err == nil {
		t.Fatalf("StartServer deveria falhar quando a porta está ocupada sem /health válido")
	}
	if reused {
		t.Fatalf("porta ocupada por serviço não saudável não deve ser reuso")
	}
	if !strings.Contains(err.Error(), "porta") {
		t.Fatalf("erro deveria explicar porta ocupada, obtido: %v", err)
	}
}

func TestSignHTTPPostsToServer(t *testing.T) {
	var gotMethod, gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"resourceType":"Signature","status":"SUCCESS"}`))
	}))
	defer srv.Close()

	c := New(portFromURL(t, srv.URL), "")
	out, err := c.SignHTTP(SignRequest{Input: "doc.txt", Signer: "Maria"})
	if err != nil {
		t.Fatalf("SignHTTP retornou erro: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/sign" {
		t.Fatalf("requisição inesperada: %s %s", gotMethod, gotPath)
	}
	if !strings.Contains(string(out), "Signature") {
		t.Fatalf("resposta inesperada: %s", string(out))
	}
}

func startAssinadorHealthServer(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"status":"UP"}`))
	})
	srv := &http.Server{Handler: mux}
	go func() { _ = srv.Serve(ln) }()
	return ln.Addr().(*net.TCPAddr).Port, func() { _ = srv.Close() }
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

func TestWriteOutputWritesFile(t *testing.T) {
	path := t.TempDir() + "/out.json"
	var stdout strings.Builder
	if err := WriteOutput([]byte(`{"ok":true}`), path, &stdout); err != nil {
		t.Fatalf("WriteOutput retornou erro: %v", err)
	}
	if !strings.Contains(stdout.String(), fmt.Sprintf("Resultado gravado em %s", path)) {
		t.Fatalf("mensagem legível esperada, obtida: %s", stdout.String())
	}
}

func TestSprint2RunLocalRejectsMissingJarBeforeJavaProvisioning(t *testing.T) {
	client := New(DefaultPort, filepath.Join(t.TempDir(), "assinador-inexistente.jar"))
	_, stderr, code, err := client.SignLocal(SignRequest{Input: "documento.txt", Signer: "Teste"})
	if err == nil {
		t.Fatalf("SignLocal deveria falhar quando o jar não existe")
	}
	if code != 2 {
		t.Fatalf("código esperado 2, obtido %d", code)
	}
	if !strings.Contains(string(stderr), "assinador.jar não encontrado") {
		t.Fatalf("mensagem deveria orientar sobre jar ausente, obtida: %s", string(stderr))
	}
}
