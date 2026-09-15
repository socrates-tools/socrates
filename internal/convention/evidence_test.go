package convention

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

func TestSocratesHoldsItselfToTheTestQualityRule(t *testing.T) {
	// scenario: socrates holds itself to it
	reading, err := TestsReadingGoFiles(source.Root(t), ContractsTestsRead)
	if err != nil {
		t.Fatalf("reading the test suite: %v", err)
	}

	for _, unnamed := range reading.Unnamed {
		t.Errorf("%s reads a Go file without naming the contract it reads: assert behaviour through an interface instead, or name the contract in ContractsTestsRead", unnamed)
	}
	for _, stale := range reading.Stale {
		t.Errorf("ContractsTestsRead names %s, which no longer reads a Go file: remove the entry", stale)
	}
	for path, contract := range ContractsTestsRead {
		if strings.TrimSpace(contract) == "" {
			t.Errorf("%s reads a Go file with no contract named", path)
		}
	}
}

func TestASourceMatchingTestIsReported(t *testing.T) {
	root := t.TempDir()
	write(t, root, filepath.Join("example", "example_test.go"), `package example

import (
	"go/parser"
	"go/token"
	"testing"
)

func TestTheFeatureIsImplemented(t *testing.T) {
	set := token.NewFileSet()
	if _, err := parser.ParseFile(set, "feature.go", nil, 0); err != nil {
		t.Fatal(err)
	}
}
`)

	reading, err := TestsReadingGoFiles(root, nil)
	if err != nil {
		t.Fatalf("reading the fabricated suite: %v", err)
	}

	want := []string{filepath.Join("example", "example_test.go")}
	if diff := cmp.Diff(want, reading.Unnamed); diff != "" {
		t.Errorf("the reported readers differ (-want +got):\n%s", diff)
	}
}

func TestAHelperThatListsTheSourceIsReported(t *testing.T) {
	root := t.TempDir()
	write(t, root, filepath.Join("example", "listing_test.go"), `package example

import (
	"os"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/testing/source"
)

func TestTheRuleExists(t *testing.T) {
	for _, path := range source.Files(t, source.Root(t)) {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(content), "the rule") {
			return
		}
	}
	t.Error("the rule is not in the source")
}
`)

	reading, err := TestsReadingGoFiles(root, nil)
	if err != nil {
		t.Fatalf("reading the fabricated suite: %v", err)
	}

	want := []string{filepath.Join("example", "listing_test.go")}
	if diff := cmp.Diff(want, reading.Unnamed); diff != "" {
		t.Errorf("the reported readers differ (-want +got):\n%s", diff)
	}
}

func TestARegisteredReaderIsNotReported(t *testing.T) {
	root := t.TempDir()
	relative := filepath.Join("example", "example_test.go")
	write(t, root, relative, `package example

import (
	"go/parser"
	"go/token"
	"testing"
)

func TestTheGeneratedFileParses(t *testing.T) {
	set := token.NewFileSet()
	if _, err := parser.ParseFile(set, "generated.go", nil, 0); err != nil {
		t.Fatal(err)
	}
}
`)

	reading, err := TestsReadingGoFiles(root, map[string]string{relative: "the generated file, which is a public output of this module"})
	if err != nil {
		t.Fatalf("reading the fabricated suite: %v", err)
	}

	if len(reading.Unnamed) != 0 {
		t.Errorf("a test that names its contract was reported: %v", reading.Unnamed)
	}
	if diff := cmp.Diff([]string{relative}, reading.Named); diff != "" {
		t.Errorf("the named readers differ (-want +got):\n%s", diff)
	}
	if len(reading.Stale) != 0 {
		t.Errorf("a live entry was reported as stale: %v", reading.Stale)
	}
}

func TestATestThatExercisesBehaviourIsNotReported(t *testing.T) {
	root := t.TempDir()
	write(t, root, filepath.Join("example", "behaviour_test.go"), `package example

import "testing"

func TestAnEmptyDocumentIsRejected(t *testing.T) {
	if err := Parse(""); err == nil {
		t.Error("an empty document was accepted")
	}
}
`)

	reading, err := TestsReadingGoFiles(root, nil)
	if err != nil {
		t.Fatalf("reading the fabricated suite: %v", err)
	}
	if len(reading.Unnamed) != 0 {
		t.Errorf("a behavioural test was reported as reading source: %v", reading.Unnamed)
	}
}

func TestAnEntryForATestThatNoLongerReadsSourceIsStale(t *testing.T) {
	root := t.TempDir()
	relative := filepath.Join("example", "behaviour_test.go")
	write(t, root, relative, `package example

import "testing"

func TestAnEmptyDocumentIsRejected(t *testing.T) {
	if err := Parse(""); err == nil {
		t.Error("an empty document was accepted")
	}
}
`)

	reading, err := TestsReadingGoFiles(root, map[string]string{relative: "a contract that stopped being read"})
	if err != nil {
		t.Fatalf("reading the fabricated suite: %v", err)
	}

	if diff := cmp.Diff([]string{relative}, reading.Stale); diff != "" {
		t.Errorf("the stale entries differ (-want +got):\n%s", diff)
	}
}

func write(t *testing.T, root, relative, content string) {
	t.Helper()

	path := filepath.Join(root, relative)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("creating the directory for %s: %v", relative, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing %s: %v", relative, err)
	}
}
