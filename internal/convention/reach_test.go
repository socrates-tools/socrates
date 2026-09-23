package convention_test

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/internal/convention"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

var socratesOwn = []string{
	filepath.Join("review", "command"),
	filepath.Join("review", "agent"),
	filepath.Join("review", "config"),
	filepath.Join("review", "dialogue"),
	filepath.Join("review", "findings"),
	filepath.Join("review", "prompt"),
	filepath.Join("review", "report"),
	filepath.Join("review", "reviewer"),
	filepath.Join("review", "runstore"),
	filepath.Join("review", "skill"),
	filepath.Join("review", "snapshot"),
}

func TestSocratesReachesNothingNew(t *testing.T) {
	root := source.Root(t)

	reach, err := convention.NetworkReach(root, convention.Socrates, convention.BuildConstrainedFiles)
	if err != nil {
		t.Fatalf("reading what socrates reaches: %v", err)
	}

	for _, own := range socratesOwn {
		if !slices.Contains(reach.Measured, own) {
			t.Errorf("%s is socrates' own and is not measured, so what socrates reaches is going unread", own)
		}
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

func TestTheDispatcherAndItsBinaryReachNoNetworkDirectly(t *testing.T) {
	root := source.Root(t)

	reach, err := convention.NetworkReach(root, convention.Dispatcher, convention.BuildConstrainedFiles)
	if err != nil {
		t.Fatalf("reading what the dispatcher and its binary reach: %v", err)
	}

	want := []string{filepath.Join("cmd", "socrates"), filepath.Join("internal", "dispatcher")}
	if diff := cmp.Diff(want, reach.Measured); diff != "" {
		t.Errorf("what is measured is not the dispatcher's and the binary's own files (-want +got):\n%s", diff)
	}

	for _, reaching := range reach.Reaching {
		t.Errorf("%s, and neither the dispatcher nor the binary's own main reaches the network directly", reaching)
	}
}

func TestAToolTheReviewTreeDoesNotReachImportingTheNetworkIsNotReported(t *testing.T) {
	root := t.TempDir()
	scope := convention.Scope{Module: "example.com/two", From: []string{filepath.Join("review", "command")}}

	write := func(relative, content string) {
		t.Helper()
		writeNested(t, filepath.Join(root, relative), content)
	}
	write(filepath.Join("review", "command", "root.go"), "package command\n\nimport \"example.com/two/review/reviewer\"\n\nvar _ = reviewer.Begin\n")
	write(filepath.Join("review", "reviewer", "reviewer.go"), "package reviewer\n\nfunc Begin() {}\n")

	before, err := convention.NetworkReach(root, scope, nil)
	if err != nil {
		t.Fatalf("reading what a fabricated tree reaches: %v", err)
	}

	write(filepath.Join("agora", "surface", "server.go"), "package surface\n\nimport \"net/http\"\n\nvar server = &http.Server{}\n")
	write(filepath.Join("internal", "dispatcher", "root.go"), "package dispatcher\n\nimport (\n\t_ \"example.com/two/agora/surface\"\n\t_ \"example.com/two/review/command\"\n)\n")

	after, err := convention.NetworkReach(root, scope, nil)
	if err != nil {
		t.Fatalf("reading what a fabricated tree reaches once the page server and the dispatcher are there: %v", err)
	}

	if len(after.Reaching) != 0 {
		t.Errorf("the page server importing the network is reported as the review tree reaching it: %v", after.Reaching)
	}
	if diff := cmp.Diff(before.Measured, after.Measured); diff != "" {
		t.Errorf("what the review tree reaches changed when the page server and the dispatcher appeared (-before +after):\n%s", diff)
	}
}

func TestSomethingSocratesReachesImportingTheNetworkIsReported(t *testing.T) {
	root := t.TempDir()
	scope := convention.Scope{Module: "example.com/two", From: []string{filepath.Join("cmd", "socrates")}}

	mkdir(t, filepath.Join(root, "cmd", "socrates"))
	mkdir(t, filepath.Join(root, "internal", "reviewer"))

	write(t, filepath.Join(root, "cmd", "socrates", "main.go"), "package main\n\nimport \"example.com/two/internal/reviewer\"\n\nvar _ = reviewer.Begin\n")
	write(t, filepath.Join(root, "internal", "reviewer", "reviewer.go"), "package reviewer\n\nimport \"net/http\"\n\nvar client = http.DefaultClient\n\nfunc Begin() {}\n")

	reach, err := convention.NetworkReach(root, scope, nil)
	if err != nil {
		t.Fatalf("reading what a fabricated tree reaches: %v", err)
	}

	want := []string{filepath.Join("internal", "reviewer", "reviewer.go") + " imports net/http"}
	if diff := cmp.Diff(want, reach.Reaching); diff != "" {
		t.Errorf("what a package socrates reaches reports differs (-want +got):\n%s", diff)
	}
}

func TestAFileReachingTheNetworkIsReported(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "reaching.go"), "package reaching\n\nimport \"net/http\"\n\nvar client = http.DefaultClient\n")

	reach, err := convention.NetworkReach(root, convention.Scope{From: []string{"."}}, nil)
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

	reach, err := convention.NetworkReach(root, convention.Scope{From: []string{"."}}, map[string]string{"live.go": "one real call against the real thing"})
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

	reach, err := convention.NetworkReach(root, convention.Scope{From: []string{"."}}, map[string]string{"gone.go": "a file that is not there"})
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

func mkdir(t *testing.T, path string) {
	t.Helper()

	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("making %s: %v", path, err)
	}
}

func write(t *testing.T, path, content string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
}
