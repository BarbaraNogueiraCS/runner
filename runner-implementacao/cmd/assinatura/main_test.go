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
	f, err := parseFlags([]string{"--input=doc.txt", "--signer", "Maria", "--local"})
	if err != nil {
		t.Fatalf("parseFlags retornou erro: %v", err)
	}
	if f["input"] != "doc.txt" || f["signer"] != "Maria" || f["local"] != "true" {
		t.Fatalf("flags parseadas incorretamente: %#v", f)
	}
}
