package httpx

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetReturnsBodyAndStatusOnSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("método esperado GET, obtido %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	body, status, err := Get(New(time.Second, false), server.URL)
	if err != nil {
		t.Fatalf("Get retornou erro: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("status esperado 200, obtido %d", status)
	}
	if string(body) != `{"status":"ok"}` {
		t.Fatalf("corpo inesperado: %s", string(body))
	}
}

func TestGetReturnsErrorOnNon2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "falha", http.StatusBadGateway)
	}))
	defer server.Close()

	_, status, err := Get(New(time.Second, false), server.URL)
	if err == nil {
		t.Fatalf("Get deveria retornar erro para resposta não 2xx")
	}
	if status != http.StatusBadGateway {
		t.Fatalf("status esperado 502, obtido %d", status)
	}
}

func TestPostJSONSendsJSONPayload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("método esperado POST, obtido %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); !strings.Contains(ct, "application/json") {
			t.Fatalf("Content-Type esperado application/json, obtido %s", ct)
		}
		var payload map[string]string
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("payload JSON inválido: %v", err)
		}
		if payload["name"] != "runner" {
			t.Fatalf("payload inesperado: %#v", payload)
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"created":true}`))
	}))
	defer server.Close()

	body, status, err := PostJSON(New(time.Second, false), server.URL, map[string]string{"name": "runner"})
	if err != nil {
		t.Fatalf("PostJSON retornou erro: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("status esperado 201, obtido %d", status)
	}
	if string(body) != `{"created":true}` {
		t.Fatalf("corpo inesperado: %s", string(body))
	}
}

func TestPostEmptySendsPostRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("método esperado POST, obtido %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"stopped":true}`))
	}))
	defer server.Close()

	body, status, err := PostEmpty(New(time.Second, false), server.URL)
	if err != nil {
		t.Fatalf("PostEmpty retornou erro: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("status esperado 200, obtido %d", status)
	}
	if string(body) != `{"stopped":true}` {
		t.Fatalf("corpo inesperado: %s", string(body))
	}
}
