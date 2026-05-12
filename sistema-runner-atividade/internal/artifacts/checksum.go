package artifacts

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func SHA256File(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func VerifySHA256(path string, expected string) error {
	if expected == "" {
		return nil
	}
	actual, err := SHA256File(path)
	if err != nil {
		return err
	}
	if actual != expected {
		return fmt.Errorf("checksum divergente: esperado %s, obtido %s", expected, actual)
	}
	return nil
}
