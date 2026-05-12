package process

import "testing"

func TestIsPortAvailableWithInvalidPort(t *testing.T) {
	if IsPortAvailable(-1) {
		t.Fatal("porta inválida não deve estar disponível")
	}
}
