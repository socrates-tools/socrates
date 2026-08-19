package fidelity_test

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"slices"
	"testing"

	"github.com/rui-ktei/socrates/internal/testing/source"
)

var behindATag = []string{"record_test.go", "drift_test.go", "livesmoke_test.go"}

func TestTheOrdinarySuiteReachesNoNetwork(t *testing.T) {
	listed := listPackage(t)

	compiled := slices.Concat(listed.GoFiles, listed.TestGoFiles, listed.XTestGoFiles)
	for _, guarded := range behindATag {
		if slices.Contains(compiled, guarded) {
			t.Errorf("the default build compiles %s, which reaches the real CLI", guarded)
		}
		if !slices.Contains(listed.IgnoredGoFiles, guarded) {
			t.Errorf("%s is not excluded from the default build: give it a build tag of its own", guarded)
		}
	}
}

func TestEveryFileBehindATagIsThere(t *testing.T) {
	listed := listPackage(t)

	for _, guarded := range behindATag {
		if _, err := os.Stat(guarded); err != nil {
			t.Errorf("%s is named as being behind a build tag and is not there: %v", guarded, err)
		}
	}
	if len(listed.IgnoredGoFiles) != len(behindATag) {
		t.Errorf("the default build ignores %v, and the files behind a tag are %v", listed.IgnoredGoFiles, behindATag)
	}
}

type packageFiles struct {
	GoFiles        []string
	TestGoFiles    []string
	XTestGoFiles   []string
	IgnoredGoFiles []string
}

func listPackage(t *testing.T) packageFiles {
	t.Helper()

	command := exec.Command("go", "list", "-json", "./internal/fidelity")
	command.Dir = source.Root(t)
	command.Env = append(os.Environ(), source.Toolchain()...)

	var out, failures bytes.Buffer
	command.Stdout = &out
	command.Stderr = &failures
	if err := command.Run(); err != nil {
		t.Fatalf("listing the package: %v: %s", err, failures.String())
	}

	var listed packageFiles
	if err := json.Unmarshal(out.Bytes(), &listed); err != nil {
		t.Fatalf("reading what go list said: %v", err)
	}
	return listed
}
