package netutil

import (
	"fmt"
	"net"
	"time"
)

// IsTCPPortFree verifies whether a local TCP port can be bound by a new
// process. It is intentionally conservative: if the check cannot bind the
// port, the caller should assume another process is already using it.
func IsTCPPortFree(port int) bool {
	if port <= 0 || port > 65535 {
		return false
	}
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return false
	}
	_ = ln.Close()
	return true
}

// IsTCPPortReachable verifies whether some process is accepting TCP
// connections on localhost. It does not verify the application protocol.
func IsTCPPortReachable(port int, timeout time.Duration) bool {
	if port <= 0 || port > 65535 {
		return false
	}
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), timeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
