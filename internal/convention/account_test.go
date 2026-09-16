package convention_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/internal/convention"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

func TestNothingNamesAnAccountOtherThanTheOrganisation(t *testing.T) {
	named, err := convention.AccountsOtherThanTheOrganisation(source.Root(t))
	if err != nil {
		t.Fatalf("reading what the project names itself under: %v", err)
	}

	for _, account := range named {
		t.Errorf("%s:%d names the account %s, and the project answers to %s", account.File, account.Line, account.Account, convention.Organisation)
	}
}

func TestASourceFileNamingAnotherAccountFailsTheCheck(t *testing.T) {
	root := fabricatedTree(t)
	writeNested(t, filepath.Join(root, "elsewhere.go"), "package elsewhere\n\nimport _ \""+hostedBy("github.com", "a-stranger", convention.Suite)+"/internal/report\"\n")

	named, err := convention.AccountsOtherThanTheOrganisation(root)
	if err != nil {
		t.Fatalf("reading a fabricated tree: %v", err)
	}

	want := []convention.NamedAccount{{File: "elsewhere.go", Line: 3, Account: "a-stranger"}}
	if diff := cmp.Diff(want, named); diff != "" {
		t.Errorf("what a file naming another account reports differs (-want +got):\n%s", diff)
	}
}

func TestAnAccountOnAnotherHostIsNamed(t *testing.T) {
	root := fabricatedTree(t)
	writeNested(t, filepath.Join(root, "README.md"), "Clone it from https://"+hostedBy("gitlab.com", "a-mirror", convention.Repository)+".\n")

	named, err := convention.AccountsOtherThanTheOrganisation(root)
	if err != nil {
		t.Fatalf("reading a fabricated tree: %v", err)
	}

	want := []convention.NamedAccount{{File: "README.md", Line: 1, Account: "a-mirror"}}
	if diff := cmp.Diff(want, named); diff != "" {
		t.Errorf("what a mirror on another host reports differs (-want +got):\n%s", diff)
	}
}

func fabricatedTree(t *testing.T) string {
	t.Helper()

	root := t.TempDir()
	writeNested(t, filepath.Join(root, "go.mod"), "module "+hostedBy("github.com", convention.Organisation, convention.Suite)+"\n")
	return root
}

func hostedBy(host, account, project string) string {
	return host + "/" + account + "/" + project
}

func writeNested(t *testing.T, path, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("making %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
}

func TestTheModulePathIsNamedForTheSuite(t *testing.T) {
	root := source.Root(t)

	path, err := convention.ModulePath(root)
	if err != nil {
		t.Fatalf("reading the module path: %v", err)
	}
	segments := strings.Split(path, "/")
	if len(segments) < 3 {
		t.Fatalf("the module path %q names no organisation and suite", path)
	}

	named, organisation := segments[len(segments)-1], segments[len(segments)-2]
	if named != convention.Suite {
		t.Errorf("the module path ends in %q and the suite is %q, so the module is not named for the suite", named, convention.Suite)
	}
	if organisation != convention.Organisation {
		t.Errorf("the module path names %q before the suite, and the organisation is %q", organisation, convention.Organisation)
	}
	if organisation == named {
		t.Errorf("the organisation and the suite are both %q, so a tool the organisation publishes has no name to take", named)
	}
	if named == convention.Repository {
		t.Errorf("the module path ends in %q and the tool installs the same name, so the suite has no room for a second tool", named)
	}

	entry, err := os.Stat(filepath.Join(root, "cmd", convention.Suite))
	if err != nil || !entry.IsDir() {
		t.Errorf("the tool is built from no cmd/%s, so what the suite installs and what the build produces are different names", convention.Suite)
	}
}
