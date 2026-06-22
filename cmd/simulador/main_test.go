package main

import "testing"

func TestVersionCommand(t *testing.T) {
	if code := run([]string{"version"}); code != 0 {
		t.Fatalf("exit code esperado 0, obtido %d", code)
	}
}

func TestUnknownCommand(t *testing.T) {
	if code := run([]string{"desconhecido"}); code == 0 {
		t.Fatal("comando desconhecido deveria retornar erro")
	}
}

func TestParseFlagsAcceptsEqualsSyntax(t *testing.T) {
	f, err := parseFlags([]string{"--url=https://localhost:8443", "--insecure"})
	if err != nil {
		t.Fatalf("parseFlags retornou erro: %v", err)
	}
	if f["url"] != "https://localhost:8443" || f["insecure"] != "true" {
		t.Fatalf("flags parseadas incorretamente: %#v", f)
	}
}

func TestURLWithPortBuildsDefaultHTTPSURL(t *testing.T) {
	u, err := urlWithPort("", "9443")
	if err != nil {
		t.Fatalf("urlWithPort retornou erro: %v", err)
	}
	if u != "https://localhost:9443" {
		t.Fatalf("URL inesperada: %s", u)
	}
}

func TestURLWithPortOverridesURLPort(t *testing.T) {
	u, err := urlWithPort("http://127.0.0.1:8080", "9443")
	if err != nil {
		t.Fatalf("urlWithPort retornou erro: %v", err)
	}
	if u != "http://127.0.0.1:9443" {
		t.Fatalf("URL inesperada: %s", u)
	}
}

func TestNewClientFromFlagsAcceptsSourceAndChecksum(t *testing.T) {
	c, err := newClientFromFlags(flags{"port": "9443", "source": "https://example.org/simulador.jar", "sha256": "abc"})
	if err != nil {
		t.Fatalf("newClientFromFlags retornou erro: %v", err)
	}
	if c.BaseURL != "https://localhost:9443" || c.SourceURL != "https://example.org/simulador.jar" || c.SourceSHA256 != "abc" {
		t.Fatalf("cliente inesperado: %#v", c)
	}
}
