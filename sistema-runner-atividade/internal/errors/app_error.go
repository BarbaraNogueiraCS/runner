package errors

import "fmt"

type AppError struct {
	Code    string
	Message string
	Details string
}

func (e AppError) Error() string {
	if e.Details == "" {
		return fmt.Sprintf("%s: %s", e.Code, e.Message)
	}
	return fmt.Sprintf("%s: %s — %s", e.Code, e.Message, e.Details)
}

func New(code, message, details string) AppError {
	return AppError{Code: code, Message: message, Details: details}
}

func Wrap(code, message string, err error) AppError {
	if err == nil {
		return New(code, message, "")
	}
	return New(code, message, err.Error())
}

const (
	InvalidParameter     = "INVALID_PARAMETER"
	JavaNotFound         = "JAVA_NOT_FOUND"
	JarNotFound          = "JAR_NOT_FOUND"
	JarExecutionFailed   = "JAR_EXECUTION_FAILED"
	PortUnavailable      = "PORT_UNAVAILABLE"
	ProcessNotFound      = "PROCESS_NOT_FOUND"
	ProcessRegistryStale = "PROCESS_REGISTRY_STALE"
	HTTPConnectionFailed = "HTTP_CONNECTION_FAILED"
	HTTPInvalidResponse  = "HTTP_INVALID_RESPONSE"
	DownloadFailed       = "DOWNLOAD_FAILED"
	ChecksumMismatch     = "CHECKSUM_MISMATCH"
	PKCS11DeviceNotFound = "PKCS11_DEVICE_NOT_FOUND"
	InternalError        = "INTERNAL_ERROR"
)
