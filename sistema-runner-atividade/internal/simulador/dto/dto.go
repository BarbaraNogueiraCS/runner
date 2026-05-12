package dto

type SimulatorConfig struct {
	Port    int
	JarPath string
	Source  string
	SHA256  string
}

type SimulatorStatus struct {
	Running bool
	Port    int
	PID     int
	Message string
	Raw     string
}
