package e2e_test

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	implementskill "github.com/socrates-tools/socrates/implement/skill"
	skillkit "github.com/socrates-tools/socrates/internal/skill"
	"github.com/socrates-tools/socrates/internal/state"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

const (
	installerFile        = "install.sh"
	releasedVersion      = "9.9.9"
	releasedTag          = "v" + releasedVersion
	pinnedVersion        = "8.8.8"
	pinnedTag            = "v" + pinnedVersion
	checksumsFile        = "checksums.txt"
	checksumsBundle      = "checksums.txt.bundle"
	decoyExitCode        = 97
	defaultInstallSuffix = ".local/bin"
)

var releasedPairs = []string{"linux_amd64", "linux_arm64", "darwin_amd64", "darwin_arm64"}

var toolsTheInstallerRuns = []string{"sh", "tar", "gzip", "mkdir", "mktemp", "rm", "mv", "chmod", "grep", "sed", "tr", "cut", "cat"}

type machine struct {
	kernel       string
	hardware     string
	pair         string
	fetcher      string
	checksumTool string
}

type release struct {
	directory string
	server    *httptest.Server
	real      []byte
}

type installation struct {
	home       string
	installDir string
	tools      string
	printed    string
	failures   string
	code       int
}

func archiveHolding(t *testing.T, binary []byte) []byte {
	t.Helper()

	var packed bytes.Buffer
	compressed := gzip.NewWriter(&packed)
	archive := tar.NewWriter(compressed)
	for _, member := range []struct {
		name    string
		mode    int64
		content []byte
	}{
		{"LICENSE", 0o644, []byte("MIT\n")},
		{"NOTICE", 0o644, []byte("notice\n")},
		{"socrates", 0o755, binary},
	} {
		if err := archive.WriteHeader(&tar.Header{Name: member.name, Mode: member.mode, Size: int64(len(member.content))}); err != nil {
			t.Fatal(err)
		}
		if _, err := archive.Write(member.content); err != nil {
			t.Fatal(err)
		}
	}
	if err := archive.Close(); err != nil {
		t.Fatal(err)
	}
	if err := compressed.Close(); err != nil {
		t.Fatal(err)
	}
	return packed.Bytes()
}

func publishRelease(t *testing.T, directory, version, realPair string, real []byte) {
	t.Helper()

	if err := os.MkdirAll(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	decoy := fmt.Appendf(nil, "#!/bin/sh\necho decoy >&2\nexit %d\n", decoyExitCode)

	var checksums strings.Builder
	for _, pair := range releasedPairs {
		held := decoy
		if pair == realPair {
			held = real
		}
		name := "socrates_" + version + "_" + pair + ".tar.gz"
		packed := archiveHolding(t, held)
		if err := os.WriteFile(filepath.Join(directory, name), packed, 0o644); err != nil {
			t.Fatal(err)
		}
		sum := sha256.Sum256(packed)
		fmt.Fprintf(&checksums, "%s  %s\n", hex.EncodeToString(sum[:]), name)
	}
	for _, name := range []string{checksumsFile, checksumsBundle} {
		content := checksums.String()
		if name == checksumsBundle {
			content = "{}"
		}
		if err := os.WriteFile(filepath.Join(directory, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func servedRelease(t *testing.T, realPair string) *release {
	t.Helper()

	real, err := os.ReadFile(source.Shared(t, "./cmd/socrates"))
	if err != nil {
		t.Fatal(err)
	}

	root := t.TempDir()
	publishRelease(t, filepath.Join(root, "latest", "download"), releasedVersion, realPair, real)
	publishRelease(t, filepath.Join(root, "download", releasedTag), releasedVersion, realPair, real)
	publishRelease(t, filepath.Join(root, "download", pinnedTag), pinnedVersion, realPair, real)

	server := httptest.NewServer(http.FileServer(http.Dir(root)))
	t.Cleanup(server.Close)
	return &release{directory: root, server: server, real: real}
}

func (r *release) latest() string {
	return filepath.Join(r.directory, "latest", "download")
}

func linkTool(t *testing.T, directory, name string) {
	t.Helper()

	found, err := exec.LookPath(name)
	if err != nil {
		t.Fatalf("the installer runs %s, and none is on PATH: %v", name, err)
	}
	if err := os.Symlink(found, filepath.Join(directory, name)); err != nil {
		t.Fatal(err)
	}
}

func writeTool(t *testing.T, directory, name, script string) {
	t.Helper()

	if err := os.WriteFile(filepath.Join(directory, name), []byte("#!/bin/sh\n"+script), 0o755); err != nil {
		t.Fatal(err)
	}
}

func toolsFor(t *testing.T, on machine) string {
	t.Helper()

	directory := t.TempDir()
	for _, name := range toolsTheInstallerRuns {
		linkTool(t, directory, name)
	}
	linkTool(t, directory, on.fetcher)
	linkTool(t, directory, on.checksumTool)
	writeTool(t, directory, "uname", fmt.Sprintf("case \"$1\" in -s) echo %s ;; -m) echo %s ;; *) echo %s ;; esac\n", on.kernel, on.hardware, on.kernel))
	return directory
}

func runInstaller(t *testing.T, tools string, served *release, extra ...string) installation {
	t.Helper()
	return runInstallerIn(t, t.TempDir(), tools, served, extra...)
}

func runInstallerIn(t *testing.T, home, tools string, served *release, extra ...string) installation {
	t.Helper()

	done := installation{home: home, installDir: filepath.Join(home, filepath.FromSlash(defaultInstallSuffix)), tools: tools}

	command := exec.Command("sh", filepath.Join(source.Root(t), installerFile))
	command.Dir = home
	command.Env = append([]string{
		"HOME=" + home,
		"PATH=" + tools,
		state.HomeVariable + "=" + filepath.Join(home, ".socrates"),
		"SOCRATES_RELEASES_URL=" + served.server.URL,
	}, extra...)
	for _, assignment := range extra {
		if named, ok := strings.CutPrefix(assignment, "SOCRATES_INSTALL_DIR="); ok {
			done.installDir = named
		}
	}

	var out, failures bytes.Buffer
	command.Stdout = &out
	command.Stderr = &failures
	if err := command.Run(); err != nil {
		var exited *exec.ExitError
		if !errors.As(err, &exited) {
			t.Fatalf("running the installer: %v", err)
		}
		done.code = exited.ExitCode()
	}
	done.printed, done.failures = out.String(), failures.String()
	return done
}

func (i installation) binary() string {
	return filepath.Join(i.installDir, "socrates")
}

func assertInstalled(t *testing.T, done installation, served *release) {
	t.Helper()

	if done.code != 0 {
		t.Fatalf("the installer exited %d:\n%s\n%s", done.code, done.failures, done.printed)
	}
	placed, err := os.ReadFile(done.binary())
	if err != nil {
		t.Fatalf("the installer placed no binary at %s: %v", done.binary(), err)
	}
	if !bytes.Equal(placed, served.real) {
		t.Fatalf("the binary at %s is not the one from this machine's archive", done.binary())
	}

	written := writtenUnder(t, done.home)
	for _, one := range allSkills() {
		if written["skills/"+one.Name+"/"+skillkit.FileName] != one.Document() {
			t.Errorf("the installer left %s out, or installed something other than what its source generates", one.Name)
		}
	}
	if _, held := written["agents/"+implementskill.AgentFileName()]; !held {
		t.Error("the installer left the implementer's agent definition out")
	}
	if !strings.Contains(done.printed, done.binary()) {
		t.Errorf("the installer does not say where the binary went:\n%s", done.printed)
	}
}

func assertNothingInstalled(t *testing.T, done installation) {
	t.Helper()

	if done.code == 0 {
		t.Fatalf("the installer exited 0:\n%s", done.printed)
	}
	if _, err := os.Stat(done.binary()); !os.IsNotExist(err) {
		t.Errorf("the installer refused and still placed %s", done.binary())
	}
	if written := writtenUnder(t, done.home); len(written) != 0 {
		t.Errorf("the installer refused and still wrote skills: %v", written)
	}
}

func onLinux() machine {
	return machine{kernel: "Linux", hardware: "x86_64", pair: "linux_amd64", fetcher: "curl", checksumTool: "sha256sum"}
}

func TestTheInstallerPlacesTheArchiveForTheMachineAndEverySkill(t *testing.T) {
	t.Parallel()

	for _, on := range []machine{
		onLinux(),
		{kernel: "Linux", hardware: "aarch64", pair: "linux_arm64", fetcher: "curl", checksumTool: "sha256sum"},
		{kernel: "Darwin", hardware: "arm64", pair: "darwin_arm64", fetcher: "curl", checksumTool: "sha256sum"},
		{kernel: "Darwin", hardware: "x86_64", pair: "darwin_amd64", fetcher: "curl", checksumTool: "sha256sum"},
	} {
		t.Run(on.kernel+"_"+on.hardware, func(t *testing.T) {
			t.Parallel()
			served := servedRelease(t, on.pair)
			assertInstalled(t, runInstaller(t, toolsFor(t, on), served), served)
		})
	}
}

func TestTheInstallerRefusesAMachineTheReleaseCarriesNoArchiveFor(t *testing.T) {
	t.Parallel()
	served := servedRelease(t, "linux_amd64")
	on := onLinux()
	on.kernel, on.hardware = "FreeBSD", "amd64"

	done := runInstaller(t, toolsFor(t, on), served)

	assertNothingInstalled(t, done)
	for _, named := range []string{"freebsd/amd64", "linux/amd64", "darwin/arm64"} {
		if !strings.Contains(done.failures, named) {
			t.Errorf("the refusal does not name %s:\n%s", named, done.failures)
		}
	}
}

func TestTheInstallerRefusesAnArchiveWhoseChecksumDoesNotMatch(t *testing.T) {
	t.Parallel()
	served := servedRelease(t, "linux_amd64")
	path := filepath.Join(served.latest(), "socrates_"+releasedVersion+"_linux_amd64.tar.gz")
	if err := os.WriteFile(path, archiveHolding(t, []byte("#!/bin/sh\nexit 0\n")), 0o644); err != nil {
		t.Fatal(err)
	}

	done := runInstaller(t, toolsFor(t, onLinux()), served)

	assertNothingInstalled(t, done)
	if !strings.Contains(done.failures, "checksum") {
		t.Errorf("the refusal does not say the checksum did not match:\n%s", done.failures)
	}
}

func TestTheInstallerChecksWithShasumWhereSha256sumIsAbsent(t *testing.T) {
	t.Parallel()
	served := servedRelease(t, "darwin_arm64")
	on := machine{kernel: "Darwin", hardware: "arm64", pair: "darwin_arm64", fetcher: "curl", checksumTool: "shasum"}

	assertInstalled(t, runInstaller(t, toolsFor(t, on), served), served)
}

func TestTheInstallerDownloadsWithWgetWhereCurlIsAbsent(t *testing.T) {
	t.Parallel()
	served := servedRelease(t, "linux_amd64")
	on := onLinux()
	on.fetcher = "wget"

	assertInstalled(t, runInstaller(t, toolsFor(t, on), served), served)
}

func TestTheInstallerCreatesTheDirectoryAndSaysWhenItIsNotOnPath(t *testing.T) {
	t.Parallel()
	served := servedRelease(t, "linux_amd64")

	done := runInstaller(t, toolsFor(t, onLinux()), served)

	assertInstalled(t, done, served)
	if !strings.Contains(done.printed, "not on your PATH") || !strings.Contains(done.printed, `export PATH="`+done.installDir+`:$PATH"`) {
		t.Errorf("the installer does not say %s is not on PATH, and how to add it:\n%s", done.installDir, done.printed)
	}
}

func TestTheInstallerPlacesTheBinaryWhereSOCRATES_INSTALL_DIRSaysAndSaysNothingOfAPathThatHoldsIt(t *testing.T) {
	t.Parallel()
	served := servedRelease(t, "linux_amd64")
	tools := toolsFor(t, onLinux())
	named := filepath.Join(t.TempDir(), "elsewhere", "bin")

	done := runInstaller(t, tools, served, "PATH="+named+":"+tools, "SOCRATES_INSTALL_DIR="+named)
	if done.code != 0 {
		t.Fatalf("the installer exited %d:\n%s", done.code, done.failures)
	}
	printed := done.printed

	if placed, err := os.ReadFile(filepath.Join(named, "socrates")); err != nil || !bytes.Equal(placed, served.real) {
		t.Fatalf("the binary is not at %s: %v", named, err)
	}
	if _, err := os.Stat(filepath.Join(done.home, filepath.FromSlash(defaultInstallSuffix), "socrates")); !os.IsNotExist(err) {
		t.Error("the installer also placed the binary in the default directory")
	}
	if strings.Contains(printed, "not on your PATH") {
		t.Errorf("the installer says a directory on PATH is not on it:\n%s", printed)
	}
}

func TestTheInstallerInstallsTheReleaseSOCRATES_VERSIONNames(t *testing.T) {
	t.Parallel()
	served := servedRelease(t, "linux_amd64")
	if err := os.RemoveAll(served.latest()); err != nil {
		t.Fatal(err)
	}

	done := runInstaller(t, toolsFor(t, onLinux()), served, "SOCRATES_VERSION="+pinnedTag)

	assertInstalled(t, done, served)
}

func TestRunningTheInstallerAgainReplacesTheBinaryAndTheSkills(t *testing.T) {
	t.Parallel()
	served := servedRelease(t, "linux_amd64")
	tools := toolsFor(t, onLinux())

	first := runInstaller(t, tools, served)
	assertInstalled(t, first, served)
	stale := filepath.Join(first.home, ".claude", "skills", allSkills()[0].Name, skillkit.FileName)
	if err := os.WriteFile(stale, []byte("an older skill\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	older := first.binary() + ".older"
	if err := os.WriteFile(older, []byte("#!/bin/sh\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(older, first.binary()); err != nil {
		t.Fatal(err)
	}

	assertInstalled(t, runInstallerIn(t, first.home, tools, served), served)
}

func TestTheInstallerVerifiesTheSignatureWhereCosignIsOnPath(t *testing.T) {
	t.Parallel()
	served := servedRelease(t, "linux_amd64")
	tools := toolsFor(t, onLinux())
	called := filepath.Join(t.TempDir(), "cosign-arguments")
	writeTool(t, tools, "cosign", "printf '%s\\n' \"$@\" > "+called+"\n")

	done := runInstaller(t, tools, served)

	assertInstalled(t, done, served)
	arguments, err := os.ReadFile(called)
	if err != nil {
		t.Fatalf("cosign is on PATH and the installer never ran it: %v", err)
	}
	for _, wanted := range []string{"verify-blob", "--bundle", "--certificate-identity-regexp", "https://github.com/socrates-tools/socrates/", "--certificate-oidc-issuer", "https://token.actions.githubusercontent.com"} {
		if !strings.Contains(string(arguments), wanted) {
			t.Errorf("the installer ran cosign without %s:\n%s", wanted, arguments)
		}
	}
}

func TestTheInstallerRefusesASignatureThatDoesNotVerify(t *testing.T) {
	t.Parallel()
	served := servedRelease(t, "linux_amd64")
	tools := toolsFor(t, onLinux())
	writeTool(t, tools, "cosign", "exit 1\n")

	done := runInstaller(t, tools, served)

	assertNothingInstalled(t, done)
	if !strings.Contains(done.failures, "signature") {
		t.Errorf("the refusal does not say the signature did not verify:\n%s", done.failures)
	}
}
