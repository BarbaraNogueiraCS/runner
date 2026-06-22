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

func TestSprint2HelpCommandsReturnOK(t *testing.T) {
	for _, args := range [][]string{
		{"sign", "--help"},
		{"validate", "--help"},
		{"server", "--help"},
	} {
		if code := run(args); code != 0 {
			t.Fatalf("%v deveria retornar 0, obtido %d", args, code)
		}
	}
}

func TestSprint2UnknownFlagIsRejected(t *testing.T) {
	if code := run([]string{"sign", "--flag-inexistente", "x"}); code == 0 {
		t.Fatalf("flag desconhecida deveria retornar erro")
	}
}

func TestSprint3RootServerAliasesHelpReturnOK(t *testing.T) {
	for _, args := range [][]string{
		{"start", "--help"},
		{"status", "--help"},
		{"stop", "--help"},
	} {
		if code := run(args); code != 0 {
			t.Fatalf("%v deveria retornar 0, obtido %d", args, code)
		}
	}
}

func TestSprint3TimeoutAliasIsAccepted(t *testing.T) {
	f, err := parseFlagsKnown([]string{"--timeout", "7"}, serverAllowedFlags)
	if err != nil {
		t.Fatalf("--timeout deveria ser aceito: %v", err)
	}
	if got := idleMinutesFlag(f); got != 7 {
		t.Fatalf("timeout esperado 7, obtido %d", got)
	}
}
