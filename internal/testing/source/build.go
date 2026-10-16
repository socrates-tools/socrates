package source

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/rand/v2"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/socrates-tools/socrates/agora/surface"
)

const (
	suitesSurfaceRangeStart = 30000
	suitesSurfaceRangeSize  = 2000
	suitesSurfaceTries      = 200
)

var (
	surfaceDefaultSymbol = reflect.TypeFor[surface.Record]().PkgPath() + ".defaultAddress"
	suitesSurfaceAddress = sync.OnceValues(chooseTheSuitesSurfaceAddress)
)

func SurfaceDefaultAddress(t *testing.T) string {
	t.Helper()

	address, err := suitesSurfaceAddress()
	if err != nil {
		t.Fatal(err)
	}
	return address
}

func chooseTheSuitesSurfaceAddress() (string, error) {
	for range suitesSurfaceTries {
		address := net.JoinHostPort("127.0.0.1", strconv.Itoa(suitesSurfaceRangeStart+rand.IntN(suitesSurfaceRangeSize)))
		listener, err := net.Listen("tcp", address)
		if err != nil {
			continue
		}
		if err := listener.Close(); err != nil {
			return "", fmt.Errorf("releasing the port chosen for the suite's surface: %w", err)
		}
		return address, nil
	}
	return "", fmt.Errorf("no loopback port in %d..%d is free for the surface the suite's binaries bind where nothing states one", suitesSurfaceRangeStart, suitesSurfaceRangeStart+suitesSurfaceRangeSize)
}

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

	if err := build(Root(t), target, output, tags, ""); err != nil {
		t.Fatalf("building %s: %v", target, err)
	}
}

func build(root, target, output string, tags []string, linkerFlags string) error {
	address, err := suitesSurfaceAddress()
	if err != nil {
		return err
	}
	linked := strings.TrimSpace(linkerFlags + " -X " + surfaceDefaultSymbol + "=" + address)

	arguments := []string{"build", "-ldflags", linked}
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
	target      string
	tags        string
	linkerFlags string
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
	return shared(t, builtKey{target: target, tags: strings.Join(tags, ",")})
}

func SharedLinkedWith(t *testing.T, target, linkerFlags string) string {
	t.Helper()
	return shared(t, builtKey{target: target, linkerFlags: linkerFlags})
}

func shared(t *testing.T, key builtKey) string {
	t.Helper()

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
		built.path = filepath.Join(directory, strings.NewReplacer("/", "_", ".", "_").Replace(key.target)+"-"+strings.ReplaceAll(key.tags, ",", "_")+fmt.Sprintf("-%x", sha256.Sum256([]byte(key.linkerFlags))))
		built.err = build(root, key.target, built.path, strings.FieldsFunc(key.tags, func(r rune) bool { return r == ',' }), key.linkerFlags)
	})

	if built.err != nil {
		t.Fatalf("building %s: %v", key.target, built.err)
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

type LinkedModule struct {
	Path      string
	Version   string
	Directory string
}

func LinkedModules(t *testing.T, target string) []LinkedModule {
	t.Helper()

	command := exec.Command("go", "list", "-deps", "-f", `{{with .Module}}{{if not .Main}}{{.Path}}{{"\t"}}{{.Version}}{{"\t"}}{{.Dir}}{{end}}{{end}}`, target)
	command.Dir = Root(t)
	command.Env = append(os.Environ(), toolchain...)

	var out, failures bytes.Buffer
	command.Stdout = &out
	command.Stderr = &failures
	if err := command.Run(); err != nil {
		t.Fatalf("listing the modules %s links: %v: %s", target, err, failures.String())
	}

	var linked []LinkedModule
	for _, line := range strings.Split(out.String(), "\n") {
		fields := strings.Split(strings.TrimSpace(line), "\t")
		if len(fields) != 3 {
			continue
		}
		module := LinkedModule{Path: fields[0], Version: fields[1], Directory: fields[2]}
		if !slices.ContainsFunc(linked, func(held LinkedModule) bool { return held.Path == module.Path }) {
			linked = append(linked, module)
		}
	}
	slices.SortFunc(linked, func(a, b LinkedModule) int { return strings.Compare(a.Path, b.Path) })
	return linked
}
