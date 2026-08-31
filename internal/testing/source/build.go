package source

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var toolchain = capture()

var toolchainVariables = []string{"GOPATH", "GOMODCACHE", "GOCACHE"}

func capture() []string {
	command := exec.Command("go", "env", strings.Join(toolchainVariables, " "))
	command.Args = append([]string{"go", "env"}, toolchainVariables...)

	var out bytes.Buffer
	command.Stdout = &out
	if err := command.Run(); err != nil {
		return nil
	}

	values := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(values) != len(toolchainVariables) {
		return nil
	}

	captured := make([]string, 0, len(values))
	for position, name := range toolchainVariables {
		captured = append(captured, name+"="+strings.TrimSpace(values[position]))
	}
	return captured
}

func Toolchain() []string {
	return append([]string(nil), toolchain...)
}

func Build(t *testing.T, target, output string) {
	t.Helper()
	BuildWithTags(t, target, output)
}

func BuildWithTags(t *testing.T, target, output string, tags ...string) {
	t.Helper()

	arguments := []string{"build"}
	if len(tags) > 0 {
		arguments = append(arguments, "-tags", strings.Join(tags, ","))
	}
	arguments = append(arguments, "-o", output, target)

	build := exec.Command("go", arguments...)
	build.Dir = Root(t)
	build.Env = append(os.Environ(), toolchain...)

	var failures bytes.Buffer
	build.Stderr = &failures
	if err := build.Run(); err != nil {
		t.Fatalf("building %s: %v: %s", target, err, failures.String())
	}
}
