package simulator

import "testing"

func TestNewDefaultURL(t *testing.T) {
	c := New("", false, "")
	if c.BaseURL != DefaultURL {
		t.Fatalf("URL esperada %s, obtida %s", DefaultURL, c.BaseURL)
	}
}
