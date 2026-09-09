package e2e_test

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates-review/internal/testing/source"
)

const theEntryPoint = "verify"

var everyGate = []string{
	"gofmt -l",
	"go vet",
	"golangci-lint",
	"go test -race",
	"scripts/corpus-outcomes",
	"-tags drift",
}

func TestOneCommandRunsEveryGate(t *testing.T) {
	planned, code := runMake(t, source.Root(t), "-n", theEntryPoint)
	if code != 0 {
		t.Fatalf("make -n %s exited %d: %s", theEntryPoint, code, planned)
	}

	at := 0
	for _, gate := range everyGate {
		found := strings.Index(planned[at:], gate)
		if found < 0 {
			if strings.Contains(planned, gate) {
				t.Errorf("make %s runs %q before a gate it is meant to follow, so the cheapest gate is not the first to fail:\n%s", theEntryPoint, gate, planned)
				continue
			}
			t.Errorf("make %s runs no %q, so a gate a contributor runs is not the gate a run runs:\n%s", theEntryPoint, gate, planned)
			continue
		}
		at += found + len(gate)
	}
}

func TestFormattingFailsBeforeTheSuiteDoes(t *testing.T) {
	unformatted := t.TempDir()

	copyInto(t, filepath.Join(source.Root(t), "Makefile"), filepath.Join(unformatted, "Makefile"))
	write(t, filepath.Join(unformatted, "go.mod"), "module unformatted\n\ngo 1.26\n")
	write(t, filepath.Join(unformatted, "unformatted.go"), "package unformatted\n\nfunc Answer() int {\n\t\t\treturn  1\n}\n")

	out, code := runMake(t, unformatted, theEntryPoint)
	if code == 0 {
		t.Fatalf("make %s passed against a tree gofmt would rewrite:\n%s", theEntryPoint, out)
	}
	if !strings.Contains(out, "unformatted.go") {
		t.Errorf("what a contributor reads does not name the file gofmt would rewrite:\n%s", out)
	}
	if strings.Contains(out, "go test") {
		t.Errorf("the suite was run before formatting failed, so the cheapest gate is not the first to fail:\n%s", out)
	}
}

func TestAnUncheckedErrorReturnFailsTheLinter(t *testing.T) {
	unchecked := t.TempDir()

	copyInto(t, filepath.Join(source.Root(t), "Makefile"), filepath.Join(unchecked, "Makefile"))
	copyInto(t, filepath.Join(source.Root(t), ".golangci.yaml"), filepath.Join(unchecked, ".golangci.yaml"))
	write(t, filepath.Join(unchecked, "go.mod"), "module unchecked\n\ngo 1.26\n")
	write(t, filepath.Join(unchecked, "unchecked.go"), "package unchecked\n\nimport \"os\"\n\nfunc Forget(path string) error {\n\tfile, err := os.Create(path)\n\tif err != nil {\n\t\treturn err\n\t}\n\tfile.Close()\n\treturn nil\n}\n")

	out, code := runMake(t, unchecked, "lint")
	if code == 0 {
		t.Fatalf("make lint passed against a tree that drops the error a close returns:\n%s", out)
	}
	if !strings.Contains(out, "errcheck") {
		t.Errorf("what a contributor reads does not name the analyser that would have caught the dropped error:\n%s", out)
	}
	if !strings.Contains(out, "unchecked.go") {
		t.Errorf("what a contributor reads does not name the file that drops the error:\n%s", out)
	}
}

func runMake(t *testing.T, directory string, arguments ...string) (string, int) {
	t.Helper()

	command := exec.Command("make", arguments...)
	command.Dir = directory
	command.Env = append(os.Environ(), source.Toolchain()...)

	var out bytes.Buffer
	command.Stdout = &out
	command.Stderr = &out

	code := 0
	if err := command.Run(); err != nil {
		var exited *exec.ExitError
		if !errors.As(err, &exited) {
			t.Fatalf("running make %s in %s: %v: %s", strings.Join(arguments, " "), directory, err, out.String())
		}
		code = exited.ExitCode()
	}
	return out.String(), code
}

func copyInto(t *testing.T, from, to string) {
	t.Helper()

	held, err := os.ReadFile(from)
	if err != nil {
		t.Fatalf("reading %s: %v", from, err)
	}
	write(t, to, string(held))
}

func write(t *testing.T, path, content string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
}
