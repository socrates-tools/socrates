package outcomes_test

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/socrates-tools/socrates/review/corpus/outcomes"
	"github.com/socrates-tools/socrates/review/runstate"
)

func TestNamingCountsTheTestsThatNameAnOutcome(t *testing.T) {
	tree := t.TempDir()
	write(t, filepath.Join(tree, "selector_test.go"), `package fabricated

func TestOneNamesTheStateBySelector(t *testing.T) {
	if run.State != runstate.AwaitingResponse {
		t.Error("no")
	}
}
`)
	write(t, filepath.Join(tree, "literal_test.go"), `package fabricated

func TestTwoNamesTheStateAsWritten(t *testing.T) {
	if printed != "escalated" {
		t.Error("no")
	}
}

func helperNamingAStateOutsideATest() string {
	return string(runstate.Clean)
}
`)
	write(t, filepath.Join(tree, "nothing_test.go"), `package fabricated

func TestThreeNamesNoState(t *testing.T) {
	if 1 != 1 {
		t.Error("no")
	}
}
`)
	write(t, filepath.Join(tree, "elsewhere.go"), `package fabricated

func NotATest() string {
	return string(runstate.Escalated)
}
`)
	write(t, filepath.Join(tree, "testdata", "ignored_test.go"), `package fabricated

func TestFourIsTestData(t *testing.T) {
	if run.State != runstate.Clean {
		t.Error("no")
	}
}
`)

	naming, err := outcomes.Naming(tree)
	if err != nil {
		t.Fatalf("counting the tests that name an outcome: %v", err)
	}

	for state, want := range map[runstate.State][]string{
		runstate.AwaitingResponse: {"TestOneNamesTheStateBySelector"},
		runstate.Escalated:        {"TestTwoNamesTheStateAsWritten"},
		runstate.Clean:            nil,
	} {
		if !slices.Equal(naming[state], want) {
			t.Errorf("%s is named by %v, want %v", state, naming[state], want)
		}
	}
}

func write(t *testing.T, path, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("creating the directory for %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
}
