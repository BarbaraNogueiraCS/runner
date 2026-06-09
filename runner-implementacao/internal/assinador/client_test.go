package assinador

import "testing"

func TestNewUsesDefaultPort(t *testing.T) {
	c := New(0, "")
	if c.Port != DefaultPort {
		t.Fatalf("porta esperada %d, obtida %d", DefaultPort, c.Port)
	}
}
