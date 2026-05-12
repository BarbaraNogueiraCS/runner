package dto

type SignCommand struct {
	Documento      string
	Certificado    string
	Algoritmo      string
	Local          bool
	Port           int
	TimeoutMinutes int
	JarPath        string
}

type ValidateCommand struct {
	Documento   string
	Assinatura  string
	Certificado string
	Local       bool
	Port        int
	JarPath     string
}

type OperationResult struct {
	Success       bool
	Operation     string
	Message       string
	Signature     string
	Valid         *bool
	ExecutionMode string
	Raw           string
}
