package source

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
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

	if err := build(Root(t), target, output, tags); err != nil {
		t.Fatalf("building %s: %v", target, err)
	}
}

func build(root, target, output string, tags []string) error {
	arguments := []string{"build"}
	if len(tags) > 0 {
		arguments = append(arguments, "-tags", strings.Join(tags, ","))
	}
	arguments = append(arguments, "-o", output, target)

	command := exec.Command("go", arguments...)
	command.Dir = root
	command.Env = append(os.Environ(), toolchain...)

	var failures bytes.Buffer
	command.Stderr = &failures
	if err := command.Run(); err != nil {
		return fmt.Errorf("%w: %s", err, failures.String())
	}
	return nil
}

type builtKey struct {
	target string
	tags   string
}

type builtBinary struct {
	once sync.Once
	path string
	err  error
}

var (
	sharedMutex     sync.Mutex
	sharedDirectory string
	sharedBinaries  = map[builtKey]*builtBinary{}
)

func Shared(t *testing.T, target string, tags ...string) string {
	t.Helper()

	key := builtKey{target: target, tags: strings.Join(tags, ",")}

	sharedMutex.Lock()
	built, found := sharedBinaries[key]
	if !found {
		built = &builtBinary{}
		sharedBinaries[key] = built
	}
	sharedMutex.Unlock()

	root := Root(t)
	built.once.Do(func() {
		directory, err := sharedBuildDirectory()
		if err != nil {
			built.err = err
			return
		}
		built.path = filepath.Join(directory, strings.NewReplacer("/", "_", ".", "_").Replace(target)+"-"+strings.ReplaceAll(key.tags, ",", "_"))
		built.err = build(root, target, built.path, tags)
	})

	if built.err != nil {
		t.Fatalf("building %s: %v", target, built.err)
	}
	return built.path
}

func sharedBuildDirectory() (string, error) {
	sharedMutex.Lock()
	defer sharedMutex.Unlock()

	if sharedDirectory != "" {
		return sharedDirectory, nil
	}
	directory, err := os.MkdirTemp("", "socrates-shared-builds-")
	if err != nil {
		return "", err
	}
	sharedDirectory = directory
	return directory, nil
}

func RunSharingBuilds(m *testing.M) int {
	code := m.Run()

	sharedMutex.Lock()
	defer sharedMutex.Unlock()
	if sharedDirectory != "" {
		_ = os.RemoveAll(sharedDirectory)
	}
	return code
}
