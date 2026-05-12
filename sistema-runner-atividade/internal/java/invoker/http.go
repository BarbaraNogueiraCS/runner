package invoker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	assinaturadto "github.com/kyriosdata/assinatura/internal/assinatura/dto"
)

type AssinadorHTTPClient struct {
	BaseURL string
	Client  http.Client
}

func NewAssinadorHTTPClient(port int) AssinadorHTTPClient {
	return AssinadorHTTPClient{
		BaseURL: fmt.Sprintf("http://localhost:%d", port),
		Client:  http.Client{Timeout: 10 * time.Second},
	}
}

func (c AssinadorHTTPClient) Sign(command assinaturadto.SignCommand) (assinaturadto.OperationResult, error) {
	payload := map[string]any{"document": command.Documento, "certificate": command.Certificado, "parameters": map[string]string{"algorithm": command.Algoritmo}}
	return c.post("/sign", "sign", payload)
}

func (c AssinadorHTTPClient) Validate(command assinaturadto.ValidateCommand) (assinaturadto.OperationResult, error) {
	payload := map[string]any{"document": command.Documento, "signature": command.Assinatura, "certificate": command.Certificado}
	return c.post("/validate", "validate", payload)
}

func (c AssinadorHTTPClient) post(path string, operation string, payload any) (assinaturadto.OperationResult, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return assinaturadto.OperationResult{}, err
	}
	resp, err := c.Client.Post(c.BaseURL+path, "application/json", bytes.NewReader(body))
	if err != nil {
		return assinaturadto.OperationResult{}, err
	}
	defer resp.Body.Close()
	content, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return assinaturadto.OperationResult{}, fmt.Errorf("assinador retornou status %d: %s", resp.StatusCode, string(content))
	}
	return assinaturadto.OperationResult{Success: true, Operation: operation, Message: "Operação executada via HTTP", ExecutionMode: "servidor HTTP", Raw: string(content)}, nil
}
