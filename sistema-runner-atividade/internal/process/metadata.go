package process

import "time"

type ProcessMetadata struct {
	Application      string    `json:"application"`
	PID              int       `json:"pid"`
	Port             int       `json:"port"`
	Mode             string    `json:"mode"`
	StartedAt        time.Time `json:"startedAt"`
	HealthEndpoint   string    `json:"healthEndpoint,omitempty"`
	InfoEndpoint     string    `json:"infoEndpoint,omitempty"`
	ShutdownEndpoint string    `json:"shutdownEndpoint,omitempty"`
	Status           string    `json:"status"`
}
