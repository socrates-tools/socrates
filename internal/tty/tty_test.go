package tty

import (
	"bytes"
	"os"
	"testing"
)

func TestAPipeAndABufferAreNotTerminals(t *testing.T) {
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("making a pipe: %v", err)
	}
	t.Cleanup(func() { _ = reader.Close() })
	t.Cleanup(func() { _ = writer.Close() })

	if Attached(writer) {
		t.Error("a pipe is reported as a terminal")
	}
	if Attached(&bytes.Buffer{}) {
		t.Error("a buffer is reported as a terminal")
	}
}
