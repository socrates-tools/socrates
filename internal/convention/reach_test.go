package convention_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rui-ktei/socrates/internal/convention"
	"github.com/rui-ktei/socrates/internal/testing/source"
)

func TestSocratesReachesNothingNew(t *testing.T) {
	root := source.Root(t)

	reach, err := convention.NetworkReach(root, convention.BuildConstrainedFiles)
	if err != nil {
		t.Fatalf("reading what socrates reaches: %v", err)
	}

	for _, reaching := range reach.Reaching {
		t.Errorf("%s, and socrates reaches no network outside a build tag", reaching)
	}
	for _, undeclared := range reach.Undeclared {
		t.Errorf("%s is behind a build tag and nothing says why, so a file compiled out of the ordinary suite is invisible", undeclared)
	}
	for _, stale := range reach.Stale {
		t.Errorf("%s is declared as behind a build tag and is not", stale)
	}

	if len(reach.Tagged) == 0 {
		t.Fatal("no file is behind a build tag, so this check has nothing to compare and would pass forever")
	}
}

func TestAFileReachingTheNetworkIsReported(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "reaching.go"), "package reaching\n\nimport \"net/http\"\n\nvar client = http.DefaultClient\n")

	reach, err := convention.NetworkReach(root, nil)
	if err != nil {
		t.Fatalf("reading what a fabricated tree reaches: %v", err)
	}
	if diff := cmp.Diff([]string{"reaching.go imports net/http"}, reach.Reaching); diff != "" {
		t.Errorf("what a file importing net/http reports differs (-want +got):\n%s", diff)
	}
}

func TestAFileReachingTheNetworkBehindABuildTagIsDeclaredRatherThanReported(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "live.go"), "//go:build live\n\npackage live\n\nimport \"net/http\"\n\nvar client = http.DefaultClient\n")

	reach, err := convention.NetworkReach(root, map[string]string{"live.go": "one real call against the real thing"})
	if err != nil {
		t.Fatalf("reading what a fabricated tree reaches: %v", err)
	}
	if len(reach.Reaching) != 0 {
		t.Errorf("a file behind a build tag is reported as reaching the network: %v", reach.Reaching)
	}
	if diff := cmp.Diff([]string{"live.go"}, reach.Tagged); diff != "" {
		t.Errorf("what a build-tagged file reports differs (-want +got):\n%s", diff)
	}
	if len(reach.Undeclared) != 0 {
		t.Errorf("a declared file reads as undeclared: %v", reach.Undeclared)
	}
}

func TestABuildTagNothingDeclaresIsReported(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "quiet.go"), "//go:build quiet\n\npackage quiet\n")

	reach, err := convention.NetworkReach(root, map[string]string{"gone.go": "a file that is not there"})
	if err != nil {
		t.Fatalf("reading what a fabricated tree reaches: %v", err)
	}
	if diff := cmp.Diff([]string{"quiet.go"}, reach.Undeclared); diff != "" {
		t.Errorf("an undeclared build tag differs (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]string{"gone.go"}, reach.Stale); diff != "" {
		t.Errorf("a declaration for a file that is not there differs (-want +got):\n%s", diff)
	}
}

func write(t *testing.T, path, content string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
}
