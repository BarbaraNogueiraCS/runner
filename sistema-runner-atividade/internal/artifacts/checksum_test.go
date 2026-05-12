package artifacts

import (
	"os"
	"testing"
)

func TestSHA256File(t *testing.T) {
	file, err := os.CreateTemp("", "checksum-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	if _, err := file.WriteString("runner"); err != nil {
		t.Fatal(err)
	}
	_ = file.Close()
	got, err := SHA256File(file.Name())
	if err != nil {
		t.Fatal(err)
	}
	want := "527aa9f431539da8e151d5434d1d5e611d973f601d8e970790882624554146b0"
	if got != want {
		t.Fatalf("checksum esperado %s, obtido %s", want, got)
	}
}
