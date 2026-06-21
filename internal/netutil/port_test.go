package netutil

import (
	"net"
	"testing"
	"time"
)

func TestIsTCPPortFreeDetectsBusyPort(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("não foi possível abrir porta temporária: %v", err)
	}
	defer ln.Close()
	port := ln.Addr().(*net.TCPAddr).Port
	if IsTCPPortFree(port) {
		t.Fatalf("porta %d deveria estar ocupada", port)
	}
	if !IsTCPPortReachable(port, time.Second) {
		t.Fatalf("porta %d deveria estar alcançável", port)
	}
}
