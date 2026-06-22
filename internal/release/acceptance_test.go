package release

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestErrorScenarioEnsureArtifactRejectsInvalidChecksum(t *testing.T) {
	t.Setenv("RUNNER_HOME", t.TempDir())
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/release.json":
			fmt.Fprintf(w, `{"simulador":{"url":"%s/artifact.jar","version":"1.0.0","sha256":"0000000000000000000000000000000000000000000000000000000000000000"}}`, srv.URL)
		case "/artifact.jar":
			_, _ = w.Write([]byte("conteudo diferente"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	_, _, _, err := EnsureArtifact(context.Background(), srv.URL+"/release.json", "simulador")
	if err == nil {
		t.Fatalf("EnsureArtifact deveria rejeitar checksum inválido")
	}
	if !strings.Contains(err.Error(), "checksum SHA256 inválido") {
		t.Fatalf("erro inesperado: %v", err)
	}
}

func TestErrorScenarioManifestWithoutRequestedArtifactFails(t *testing.T) {
	m := Manifest{}
	_, err := m.Artifact("simulador")
	if err == nil {
		t.Fatalf("Artifact deveria falhar quando release.json não define simulador/validador")
	}
}

func TestAcceptanceJREURLUsesPlatformKey(t *testing.T) {
	key := PlatformKey()
	m := Manifest{JRE: map[string]string{key: "https://example.org/jre"}}
	url, err := JREURL(m)
	if err != nil {
		t.Fatalf("JREURL retornou erro: %v", err)
	}
	if url != "https://example.org/jre" {
		t.Fatalf("URL inesperada: %s", url)
	}
}
