package convention_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ktei/socrates/internal/convention"
	"github.com/ktei/socrates/internal/testing/source"
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
	writeNested(t, filepath.Join(root, "elsewhere.go"), "package elsewhere\n\nimport _ \""+hostedBy("github.com", "a-stranger")+"/internal/report\"\n")

	named, err := convention.AccountsOtherThanTheOrganisation(root)
	if err != nil {
		t.Fatalf("reading a fabricated tree: %v", err)
	}

	want := []convention.NamedAccount{{File: "elsewhere.go", Line: 3, Account: "a-stranger"}}
	if diff := cmp.Diff(want, named); diff != "" {
		t.Errorf("what a file naming another account reports differs (-want +got):\n%s", diff)
	}
}

func TestAnArchivedChangeNamingTheFormerAccountPasses(t *testing.T) {
	root := fabricatedTree(t)
	writeNested(t, filepath.Join(root, convention.ArchivedChanges, "2026-01-01-a-decision", "design.md"), "The module was `"+hostedBy("github.com", "a-former-account")+"`.\n")

	named, err := convention.AccountsOtherThanTheOrganisation(root)
	if err != nil {
		t.Fatalf("reading a fabricated tree: %v", err)
	}

	if len(named) != 0 {
		t.Errorf("the archive reported %v, and it is the record of what was decided when it was decided", named)
	}
}

func TestAnAccountOnAnotherHostIsNamed(t *testing.T) {
	root := fabricatedTree(t)
	writeNested(t, filepath.Join(root, "README.md"), "Clone it from https://"+hostedBy("gitlab.com", "a-mirror")+".\n")

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
	writeNested(t, filepath.Join(root, "go.mod"), "module "+hostedBy("github.com", convention.Organisation)+"\n")
	return root
}

func hostedBy(host, account string) string {
	return host + "/" + account + "/" + convention.Repository
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
