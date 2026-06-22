package httpx

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func New(timeout time.Duration, insecure bool) *http.Client {
	tr := http.DefaultTransport.(*http.Transport).Clone()
	if insecure {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec // usado somente para simuladores locais com certificado autoassinado
	}
	return &http.Client{Timeout: timeout, Transport: tr}
}

func Get(c *http.Client, url string) ([]byte, int, error) {
	resp, err := c.Get(url)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return b, resp.StatusCode, fmt.Errorf("resposta HTTP inesperada: %d", resp.StatusCode)
	}
	return b, resp.StatusCode, nil
}

func PostJSON(c *http.Client, url string, payload any) ([]byte, int, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, 0, err
	}
	resp, err := c.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return out, resp.StatusCode, fmt.Errorf("resposta HTTP inesperada: %d", resp.StatusCode)
	}
	return out, resp.StatusCode, nil
}

func PostEmpty(c *http.Client, url string) ([]byte, int, error) {
	resp, err := c.Post(url, "application/json", bytes.NewReader([]byte("{}")))
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return out, resp.StatusCode, fmt.Errorf("resposta HTTP inesperada: %d", resp.StatusCode)
	}
	return out, resp.StatusCode, nil
}
