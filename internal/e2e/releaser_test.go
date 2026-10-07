package e2e_test

import (
	"bytes"
	"encoding/json"
	"go/version"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/testing/source"
)

type moduleFile struct {
	Go        string
	Toolchain string
}

func TestTheReleaserBuildsUnderTheToolchainAVerificationRunInstallsFromGoMod(t *testing.T) {
	root := source.Root(t)

	pinned := releaserPinnedBy(t, root)
	releaser := moduleFileAt(t, downloadedModuleFile(t, pinned))
	module := moduleFileAt(t, filepath.Join(root, "go.mod"))

	installed := lowestToolchainSetupGoInstallsFor(module)
	asked := "go" + releaser.Go
	if version.Compare(asked, installed) > 0 {
		t.Errorf("the Makefile pins %s, whose go.mod asks for %s, and a verification run under GOTOOLCHAIN=local may be running %s, the lowest toolchain go.mod lets setup-go install, so go run cannot build the releaser there", pinned, asked, installed)
	}
}

func lowestToolchainSetupGoInstallsFor(module moduleFile) string {
	if module.Toolchain != "" {
		return module.Toolchain
	}
	stated := "go" + module.Go
	if version.Lang(stated) == stated {
		return stated + ".0"
	}
	return stated
}

func releaserPinnedBy(t *testing.T, root string) string {
	t.Helper()

	planned, code := runMake(t, root, "-n", "release-snapshot")
	if code != 0 {
		t.Fatalf("make -n release-snapshot exited %d: %s", code, planned)
	}
	arguments := strings.Fields(planned)
	for position := range len(arguments) - 2 {
		if arguments[position] == "go" && arguments[position+1] == "run" && strings.Contains(arguments[position+2], "@") {
			return arguments[position+2]
		}
	}
	t.Fatalf("make release-snapshot runs no pinned releaser through go run, so there is no pin to read:\n%s", planned)
	return ""
}

func downloadedModuleFile(t *testing.T, pinned string) string {
	t.Helper()

	command := exec.Command("go", "mod", "download", "-json", pinned)
	command.Dir = t.TempDir()
	command.Env = append(os.Environ(), source.Toolchain()...)

	var out, failures bytes.Buffer
	command.Stdout = &out
	command.Stderr = &failures
	if err := command.Run(); err != nil {
		t.Fatalf("fetching the go.mod of %s, which make release-snapshot fetches too: %v: %s%s", pinned, err, out.String(), failures.String())
	}

	var downloaded struct{ GoMod string }
	if err := json.Unmarshal(out.Bytes(), &downloaded); err != nil || downloaded.GoMod == "" {
		t.Fatalf("go mod download named no go.mod for %s (%v): %s", pinned, err, out.String())
	}
	return downloaded.GoMod
}

func moduleFileAt(t *testing.T, path string) moduleFile {
	t.Helper()

	command := exec.Command("go", "mod", "edit", "-json", path)
	command.Env = append(os.Environ(), source.Toolchain()...)

	var out, failures bytes.Buffer
	command.Stdout = &out
	command.Stderr = &failures
	if err := command.Run(); err != nil {
		t.Fatalf("reading %s: %v: %s", path, err, failures.String())
	}

	var read moduleFile
	if err := json.Unmarshal(out.Bytes(), &read); err != nil {
		t.Fatalf("reading what %s states: %v", path, err)
	}
	if read.Go == "" {
		t.Fatalf("%s states no go directive, so there is no toolchain to compare", path)
	}
	return read
}
