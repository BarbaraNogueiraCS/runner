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
