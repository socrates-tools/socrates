package coverage_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-devtool/socrates/internal/coverage"
)

const specification = `## ADDED Requirements

### Requirement: A thing happens

#### Scenario: A named test claims this one

- **WHEN** something
- **THEN** something else

#### Scenario: A subtest claims this one

#### Scenario: A marker claims this one

#### Scenario: Nothing claims this one
`

func TestEveryWayATestCanNameAScenario(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, filepath.Join("openspec", "specs", "thing", "spec.md"), specification)
	writeFile(t, root, filepath.Join("internal", "thing", "thing_test.go"), `package thing

import "testing"

func TestANamedTestClaimsThisOne(t *testing.T) {
	t.Run("A subtest claims this one", func(t *testing.T) {})
}

func TestSomethingElseEntirely(t *testing.T) {
	// scenario: A marker claims this one
}
`)

	scenarios, err := coverage.Scenarios(root)
	if err != nil {
		t.Fatalf("reading the specifications: %v", err)
	}
	if len(scenarios) != 4 {
		t.Fatalf("%d scenarios were found, want 4", len(scenarios))
	}

	claims, err := coverage.Claims(root)
	if err != nil {
		t.Fatalf("reading the tests: %v", err)
	}

	var missing []string
	for _, gap := range coverage.Gaps(scenarios, claims) {
		missing = append(missing, gap.Text)
	}
	if diff := cmp.Diff([]string{"Nothing claims this one"}, missing); diff != "" {
		t.Errorf("the uncovered scenarios differ (-want +got):\n%s", diff)
	}
}

func TestAGapNamesWhereTheScenarioIs(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, filepath.Join("openspec", "specs", "thing", "spec.md"), "#### Scenario: Nobody tests this\n")

	scenarios, err := coverage.Scenarios(root)
	if err != nil {
		t.Fatalf("reading the specifications: %v", err)
	}
	claims, err := coverage.Claims(root)
	if err != nil {
		t.Fatalf("reading the tests: %v", err)
	}

	gaps := coverage.Gaps(scenarios, claims)
	if len(gaps) != 1 {
		t.Fatalf("%d gaps were reported, want one", len(gaps))
	}
	if gaps[0].Spec != filepath.Join("openspec", "specs", "thing", "spec.md") || gaps[0].Line != 1 {
		t.Errorf("the gap is reported at %s:%d", gaps[0].Spec, gaps[0].Line)
	}
	if !gaps[0].Underway {
		t.Error("a scenario in openspec/specs is not under way, and a deployed specification always is")
	}
}

func TestAPromotedScenarioAChangeRenamesIsReadUnderItsNewName(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, filepath.Join("openspec", "specs", "thing", "spec.md"), "#### Scenario: The old name nobody will keep\n")
	writeFile(t, root, filepath.Join("openspec", "changes", "now", "tasks.md"), "- [x] 1.1 Did the thing\n")
	writeFile(t, root, filepath.Join("openspec", "changes", "now", "specs", "thing", "spec.md"),
		"## RENAMED Requirements\n\n- FROM: `#### Scenario: The old name nobody will keep`\n- TO: `#### Scenario: The name a test actually claims`\n\n## MODIFIED Requirements\n\n#### Scenario: The name a test actually claims\n")

	scenarios, err := coverage.Scenarios(root)
	if err != nil {
		t.Fatalf("reading the specifications: %v", err)
	}

	for _, scenario := range scenarios {
		if scenario.Text == "The old name nobody will keep" {
			t.Errorf("%s:%d is still read under the name the change renames away, so a test naming the new one reads as a gap", scenario.Spec, scenario.Line)
		}
	}

	named := 0
	for _, scenario := range scenarios {
		if scenario.Text == "The name a test actually claims" {
			named++
		}
	}
	if named != 2 {
		t.Errorf("%d scenarios read under the new name, want the promoted one and the delta's own", named)
	}
}

func TestAChangeNobodyHasStartedRenamesNothing(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, filepath.Join("openspec", "specs", "thing", "spec.md"), "#### Scenario: The old name nobody will keep\n")
	writeFile(t, root, filepath.Join("openspec", "changes", "later", "tasks.md"), "- [ ] 1.1 Not started\n")
	writeFile(t, root, filepath.Join("openspec", "changes", "later", "specs", "thing", "spec.md"),
		"## RENAMED Requirements\n\n- FROM: `#### Scenario: The old name nobody will keep`\n- TO: `#### Scenario: A name no test claims`\n")

	scenarios, err := coverage.Scenarios(root)
	if err != nil {
		t.Fatalf("reading the specifications: %v", err)
	}

	for _, scenario := range scenarios {
		if scenario.Spec == filepath.Join("openspec", "specs", "thing", "spec.md") && scenario.Text != "The old name nobody will keep" {
			t.Errorf("a change nobody has started renamed a promoted scenario to %q, which hides whether a test names the promoted one", scenario.Text)
		}
	}
}

func TestAPromotedScenarioNoChangeRenamesKeepsItsName(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, filepath.Join("openspec", "specs", "thing", "spec.md"), "#### Scenario: The old name nobody will keep\n")
	writeFile(t, root, filepath.Join("openspec", "changes", "now", "tasks.md"), "- [x] 1.1 Did the thing\n")
	writeFile(t, root, filepath.Join("openspec", "changes", "now", "specs", "thing", "spec.md"),
		"## MODIFIED Requirements\n\n#### Scenario: The name a test actually claims\n")

	scenarios, err := coverage.Scenarios(root)
	if err != nil {
		t.Fatalf("reading the specifications: %v", err)
	}

	for _, scenario := range scenarios {
		if scenario.Spec == filepath.Join("openspec", "specs", "thing", "spec.md") && scenario.Text != "The old name nobody will keep" {
			t.Errorf("a promoted scenario no change renames reads as %q", scenario.Text)
		}
	}
}

func TestAChangeNobodyHasStartedIsNotUnderWay(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, filepath.Join("openspec", "changes", "later", "tasks.md"), "- [ ] 1.1 Do the thing\n")
	writeFile(t, root, filepath.Join("openspec", "changes", "later", "specs", "thing", "spec.md"), "#### Scenario: Nobody tests this yet\n")
	writeFile(t, root, filepath.Join("openspec", "changes", "now", "tasks.md"), "- [x] 1.1 Did the thing\n")
	writeFile(t, root, filepath.Join("openspec", "changes", "now", "specs", "thing", "spec.md"), "#### Scenario: Nobody tests this either\n")

	scenarios, err := coverage.Scenarios(root)
	if err != nil {
		t.Fatalf("reading the specifications: %v", err)
	}

	underway := map[string]bool{}
	for _, scenario := range scenarios {
		underway[scenario.Text] = scenario.Underway
	}
	if underway["Nobody tests this yet"] {
		t.Error("a change with no task ticked is treated as under way")
	}
	if !underway["Nobody tests this either"] {
		t.Error("a change with a task ticked is not treated as under way")
	}
}

func TestAScenarioWithNoTextIsAGap(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, filepath.Join("openspec", "specs", "thing", "spec.md"), "#### Scenario:\n\n#### Scenario: ***\n")
	writeFile(t, root, filepath.Join("internal", "thing", "thing_test.go"), `package thing

import "testing"

func TestSomething(t *testing.T) {}
`)

	scenarios, err := coverage.Scenarios(root)
	if err != nil {
		t.Fatalf("reading the specifications: %v", err)
	}
	claims, err := coverage.Claims(root)
	if err != nil {
		t.Fatalf("reading the tests: %v", err)
	}

	if gaps := coverage.Gaps(scenarios, claims); len(gaps) != 2 {
		t.Errorf("%d of 2 headings with no readable text were reported as a gap, and a heading nothing can name is never covered", len(gaps))
	}
}

func TestAnArchivedChangeIsNotCountedTwice(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, filepath.Join("openspec", "specs", "thing", "spec.md"), "#### Scenario: The thing happens\n")
	writeFile(t, root, filepath.Join("openspec", "changes", "archive", "2026-01-01-thing", "tasks.md"), "- [x] 1.1 Did the thing\n")
	writeFile(t, root, filepath.Join("openspec", "changes", "archive", "2026-01-01-thing", "specs", "thing", "spec.md"), "#### Scenario: The thing happens\n")

	scenarios, err := coverage.Scenarios(root)
	if err != nil {
		t.Fatalf("reading the specifications: %v", err)
	}
	if len(scenarios) != 1 {
		t.Errorf("%d scenarios were found, and an archived change's promoted specification is the same promise once", len(scenarios))
	}
	if len(scenarios) > 0 && scenarios[0].Spec != filepath.Join("openspec", "specs", "thing", "spec.md") {
		t.Errorf("the scenario is reported at %s, want the promoted specification", scenarios[0].Spec)
	}
}

func TestNormalisingIgnoresEverythingButLettersAndDigits(t *testing.T) {
	if coverage.Normalise("A run's state: `clean`, not failed!") != coverage.Normalise("A runs state clean not failed") {
		t.Error("two spellings of one scenario normalise differently")
	}
}

func writeFile(t *testing.T, root, relative, content string) {
	t.Helper()

	path := filepath.Join(root, relative)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("creating the directory for %s: %v", relative, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing %s: %v", relative, err)
	}
}
