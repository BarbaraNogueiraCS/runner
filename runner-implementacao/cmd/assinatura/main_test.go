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
