package promotion_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates-review/internal/promotion"
)

const baseline = `# run-store Specification

## Purpose
What a run is.

## Requirements

### Requirement: A run holds one of six states

A run SHALL be in exactly one of six states.

#### Scenario: A terminal run does not move

- **WHEN** a run has reached a terminal state
- **THEN** no later write changes it

#### Scenario: A failed run reports no findings

- **WHEN** a run ends as failed
- **THEN** its verdict is the failure

### Requirement: A run is a directory

Every run SHALL own one directory.

#### Scenario: Runs sort by when they started

- **WHEN** three runs start
- **THEN** their names sort
`

func TestWhatArchivingWouldDrop(t *testing.T) {
	t.Run("A delta that carries every scenario drops nothing", func(t *testing.T) {
		root := repositoryWith(t, "carries-them", `## MODIFIED Requirements

### Requirement: A run holds one of six states

A run SHALL be in exactly one of six states, and one of them is new.

#### Scenario: A terminal run does not move

- **WHEN** a run has reached a terminal state
- **THEN** no later write changes it

#### Scenario: A failed run reports no findings

- **WHEN** a run ends as failed
- **THEN** its verdict is the failure

#### Scenario: A parked run is not finished

- **WHEN** a run is parked
- **THEN** it is not finished
`)

		if drops := found(t, root); len(drops) != 0 {
			t.Errorf("a delta carrying every scenario reports %+v", drops)
		}
	})

	t.Run("A scenario the delta does not carry is named", func(t *testing.T) {
		root := repositoryWith(t, "drops-one", `## MODIFIED Requirements

### Requirement: A run holds one of six states

A run SHALL be in exactly one of six states.

#### Scenario: A terminal run does not move

- **WHEN** a run has reached a terminal state
- **THEN** no later write changes it
`)

		want := []promotion.Drop{{
			Change:      "drops-one",
			Spec:        filepath.Join("openspec", "changes", "drops-one", "specs", "run-store", "spec.md"),
			Requirement: "A run holds one of six states",
			Condition:   promotion.ScenariosDropped,
			Scenarios:   []string{"A failed run reports no findings"},
		}}
		if diff := cmp.Diff(want, found(t, root)); diff != "" {
			t.Errorf("what archiving would drop differs (-want +got):\n%s", diff)
		}
	})

	t.Run("A modified requirement that substitutes nothing is named", func(t *testing.T) {
		root := repositoryWith(t, "stale-heading", `## MODIFIED Requirements

### Requirement: A run holds one of seven states

A run SHALL be in exactly one of seven states.

#### Scenario: A terminal run does not move

- **WHEN** a run has reached a terminal state
- **THEN** no later write changes it
`)

		drops := found(t, root)
		if len(drops) != 1 || drops[0].Condition != promotion.SubstitutesNothing {
			t.Fatalf("a modified requirement matching no baseline reports %+v", drops)
		}
		if drops[0].Requirement != "A run holds one of seven states" {
			t.Errorf("the report names %q", drops[0].Requirement)
		}
	})

	t.Run("A rename is followed to the requirement it substitutes", func(t *testing.T) {
		root := repositoryWith(t, "renames-it", `## RENAMED Requirements

- FROM: `+"`### Requirement: A run holds one of six states`"+`
- TO: `+"`### Requirement: A run holds one of seven states`"+`

## MODIFIED Requirements

### Requirement: A run holds one of seven states

A run SHALL be in exactly one of seven states.

#### Scenario: A terminal run does not move

- **WHEN** a run has reached a terminal state
- **THEN** no later write changes it

#### Scenario: A failed run reports no findings

- **WHEN** a run ends as failed
- **THEN** its verdict is the failure
`)

		if drops := found(t, root); len(drops) != 0 {
			t.Errorf("a renamed requirement carrying every scenario reports %+v", drops)
		}
	})

	t.Run("A renamed scenario is followed rather than reported as dropped", func(t *testing.T) {
		root := repositoryWith(t, "renames-a-scenario", `## RENAMED Requirements

- FROM: `+"`#### Scenario: A terminal run does not move`"+`
- TO: `+"`#### Scenario: A terminal run stays where it ended`"+`

## MODIFIED Requirements

### Requirement: A run holds one of six states

A run SHALL be in exactly one of six states.

#### Scenario: A terminal run stays where it ended

- **WHEN** a run has reached a terminal state
- **THEN** no later write changes it

#### Scenario: A failed run reports no findings

- **WHEN** a run ends as failed
- **THEN** its verdict is the failure
`)

		if drops := found(t, root); len(drops) != 0 {
			t.Errorf("a scenario the delta renames deliberately reports %+v", drops)
		}
	})

	t.Run("A scenario dropped without a rename is still named", func(t *testing.T) {
		root := repositoryWith(t, "drops-a-scenario", `## MODIFIED Requirements

### Requirement: A run holds one of six states

A run SHALL be in exactly one of six states.

#### Scenario: A terminal run stays where it ended

- **WHEN** a run has reached a terminal state
- **THEN** no later write changes it

#### Scenario: A failed run reports no findings

- **WHEN** a run ends as failed
- **THEN** its verdict is the failure
`)

		drops := found(t, root)
		if len(drops) != 1 || drops[0].Condition != promotion.ScenariosDropped {
			t.Fatalf("renaming a scenario without declaring it reports %+v", drops)
		}
		if len(drops[0].Scenarios) != 1 || drops[0].Scenarios[0] != "A terminal run does not move" {
			t.Errorf("the report names %v", drops[0].Scenarios)
		}
	})

	t.Run("A requirement added beside one that exists is named", func(t *testing.T) {
		root := repositoryWith(t, "adds-what-exists", `## ADDED Requirements

### Requirement: A run is a directory

Every run SHALL own one directory, again.

#### Scenario: Runs sort by when they started

- **WHEN** three runs start
- **THEN** their names sort
`)

		drops := found(t, root)
		if len(drops) != 1 || drops[0].Condition != promotion.AddsWhatExists {
			t.Fatalf("adding a requirement the specification already holds reports %+v", drops)
		}
	})

	t.Run("A removed requirement is a decision rather than a loss", func(t *testing.T) {
		root := repositoryWith(t, "removes-it", `## REMOVED Requirements

### Requirement: A run is a directory

**Reason**: runs stopped being directories.
`)

		if drops := found(t, root); len(drops) != 0 {
			t.Errorf("a deliberate removal reports %+v", drops)
		}
	})

	t.Run("An archived change is not checked", func(t *testing.T) {
		root := repositoryWith(t, filepath.Join("archive", "2026-01-01-old"), `## MODIFIED Requirements

### Requirement: A run holds one of six states

A run SHALL be in exactly one of six states.
`)

		if drops := found(t, root); len(drops) != 0 {
			t.Errorf("an archived change reports %+v", drops)
		}
	})
}

func found(t *testing.T, root string) []promotion.Drop {
	t.Helper()

	report, err := promotion.Check(root)
	if err != nil {
		t.Fatalf("checking what archiving would drop: %v", err)
	}
	return report.Drops
}

func TestWhatIsUnderWay(t *testing.T) {
	t.Run("A change under way is counted", func(t *testing.T) {
		checked, err := promotion.Check(repositoryWith(t, "under-way", `## ADDED Requirements

### Requirement: Something is proposed

Something SHALL be proposed.

#### Scenario: It is proposed

- **WHEN** it is read
- **THEN** it is there
`))
		if err != nil {
			t.Fatalf("checking what archiving would drop: %v", err)
		}
		if checked.Changes != 1 {
			t.Errorf("one change is under way and the check counts %d, so the live gate would skip itself", checked.Changes)
		}
	})

	t.Run("An archived change is not under way", func(t *testing.T) {
		root := t.TempDir()
		write(t, filepath.Join(root, "openspec", "specs", "run-store", "spec.md"), baseline)
		write(t, filepath.Join(root, "openspec", "changes", "archive", "2026-01-01-done", "tasks.md"), "- [x] 1.1 done\n")

		checked, err := promotion.Check(root)
		if err != nil {
			t.Fatalf("checking what archiving would drop: %v", err)
		}
		if checked.Changes != 0 {
			t.Errorf("nothing is under way and the check counts %d", checked.Changes)
		}
	})
}

func repositoryWith(t *testing.T, change, delta string) string {
	t.Helper()

	return repositoryPromoting(t, change, baseline, delta)
}

func repositoryPromoting(t *testing.T, change, promoted, delta string) string {
	t.Helper()

	root := t.TempDir()
	write(t, filepath.Join(root, "openspec", "specs", "run-store", "spec.md"), promoted)
	write(t, filepath.Join(root, "openspec", "changes", change, "specs", "run-store", "spec.md"), delta)
	write(t, filepath.Join(root, "openspec", "changes", change, "tasks.md"), "- [ ] 1.1 do it\n")
	return root
}

func write(t *testing.T, path, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("creating %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
}

func TestARenameTheGateCannotCheck(t *testing.T) {
	renamed := strings.Replace(baseline, "#### Scenario: A terminal run does not move", "#### Scenario: A terminal run stays where it ended", 1)
	delta := `## RENAMED Requirements

- FROM: ` + "`#### Scenario: A terminal run does not move`" + `
- TO: ` + "`#### Scenario: A terminal run stays where it ended`" + `

## MODIFIED Requirements

### Requirement: A run holds one of six states

A run SHALL be in exactly one of six states.

#### Scenario: A terminal run stays where it ended

- **WHEN** a run has reached a terminal state
- **THEN** no later write changes it

#### Scenario: A failed run reports no findings

- **WHEN** a run ends as failed
- **THEN** its verdict is the failure
`

	t.Run("A rename the promoted specification already holds is reported", func(t *testing.T) {
		drops := found(t, repositoryPromoting(t, "promoted-ahead", renamed, delta))
		if len(drops) != 1 || drops[0].Condition != promotion.PromotedAhead {
			t.Fatalf("a rename against a specification that already holds the new name reports %+v", drops)
		}
		if len(drops[0].Scenarios) != 1 || drops[0].Scenarios[0] != "A terminal run stays where it ended" {
			t.Errorf("the report names %v, and the new name is what the promoted specification already holds", drops[0].Scenarios)
		}
	})

	t.Run("A rename with a before-state to read is still followed", func(t *testing.T) {
		if drops := found(t, repositoryPromoting(t, "renames-it", baseline, delta)); len(drops) != 0 {
			t.Errorf("a rename the promoted specification can be compared against reports %+v", drops)
		}
	})

	t.Run("A rename that declares the same name on both sides is left alone", func(t *testing.T) {
		same := `## RENAMED Requirements

- FROM: ` + "`#### Scenario: A terminal run does not move`" + `
- TO: ` + "`#### Scenario: A terminal run does not move`" + `

## MODIFIED Requirements

### Requirement: A run holds one of six states

A run SHALL be in exactly one of six states.

#### Scenario: A terminal run does not move

- **WHEN** a run has reached a terminal state
- **THEN** no later write changes it

#### Scenario: A failed run reports no findings

- **WHEN** a run ends as failed
- **THEN** its verdict is the failure
`
		if drops := found(t, repositoryPromoting(t, "renames-nothing", baseline, same)); len(drops) != 0 {
			t.Errorf("a rename declaring one name on both sides reports %+v", drops)
		}
	})
}

func TestARenameBelongingToAnotherRequirementIsNotReported(t *testing.T) {
	promoted := `# run-store Specification

## Purpose
What a run is.

## Requirements

### Requirement: A response names its resolution

Every response SHALL name one resolution.

#### Scenario: It is rejected

- **WHEN** a response rejects a finding
- **THEN** the reason is required

### Requirement: A refusal names the condition

Every refusal SHALL name the condition.

#### Scenario: It is refused

- **WHEN** a run is already under way
- **THEN** the second one is refused
`
	delta := `## RENAMED Requirements

- FROM: ` + "`#### Scenario: It is rejected`" + `
- TO: ` + "`#### Scenario: It is refused`" + `

## MODIFIED Requirements

### Requirement: A response names its resolution

Every response SHALL name one resolution, again.

#### Scenario: It is refused

- **WHEN** a response rejects a finding
- **THEN** the reason is required

### Requirement: A refusal names the condition

Every refusal SHALL name the condition, again.

#### Scenario: It is refused

- **WHEN** a run is already under way
- **THEN** the second one is refused
`

	if drops := found(t, repositoryPromoting(t, "renames-one-of-two", promoted, delta)); len(drops) != 0 {
		t.Errorf("a rename the promoted specification can still be compared against reports %+v, and the name it moves to belongs to another requirement", drops)
	}
}

func TestARequirementBothAheadAndShortReportsBoth(t *testing.T) {
	promoted := strings.Replace(baseline,
		`#### Scenario: A terminal run does not move`,
		`#### Scenario: A terminal run stays where it ended`, 1)
	promoted = strings.Replace(promoted,
		`### Requirement: A run is a directory`,
		`#### Scenario: A cancelled run keeps its worktree

- **WHEN** a run is cancelled
- **THEN** the worktree is kept

### Requirement: A run is a directory`, 1)

	delta := `## RENAMED Requirements

- FROM: ` + "`#### Scenario: A terminal run does not move`" + `
- TO: ` + "`#### Scenario: A terminal run stays where it ended`" + `

## MODIFIED Requirements

### Requirement: A run holds one of six states

A run SHALL be in exactly one of six states.

#### Scenario: A terminal run stays where it ended

- **WHEN** a run has reached a terminal state
- **THEN** no later write changes it

#### Scenario: A failed run reports no findings

- **WHEN** a run ends as failed
- **THEN** its verdict is the failure
`

	drops := found(t, repositoryPromoting(t, "ahead-and-short", promoted, delta))
	if len(drops) != 2 {
		t.Fatalf("a requirement whose rename cannot be checked and which also drops a scenario reports %+v, want both", drops)
	}

	reported := map[promotion.Condition][]string{}
	for _, drop := range drops {
		reported[drop.Condition] = drop.Scenarios
	}
	if named := reported[promotion.PromotedAhead]; len(named) != 1 || named[0] != "A terminal run stays where it ended" {
		t.Errorf("the rename it cannot check is reported as %v", named)
	}
	if named := reported[promotion.ScenariosDropped]; len(named) != 1 || named[0] != "A cancelled run keeps its worktree" {
		t.Errorf("the scenario the promotion would lose is reported as %v", named)
	}
}
