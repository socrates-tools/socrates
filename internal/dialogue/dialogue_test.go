package dialogue_test

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/rui-ktei/socrates/internal/dialogue"
	"github.com/rui-ktei/socrates/internal/findings"
	"github.com/rui-ktei/socrates/internal/runstate"
	"github.com/rui-ktei/socrates/internal/testing/source"
)

var everySeverity = []findings.Severity{findings.Error, findings.Warning, findings.Info}

var everyAction = []findings.Action{findings.NoOp, findings.AutoFix, findings.AskUser}

func TestWhatARoundLeadsTo(t *testing.T) {
	t.Run("A clean round finishes the run", func(t *testing.T) {
		classified := classify(finding("f1", findings.Info, findings.NoOp))
		if got := dialogue.After(classified, 1); got != runstate.Clean {
			t.Errorf("a round with nothing blocking and nothing patchable is %q, want %q", got, runstate.Clean)
		}
		if got := dialogue.After(classify(), 1); got != runstate.Clean {
			t.Errorf("a round with no findings at all is %q, want %q", got, runstate.Clean)
		}
	})

	t.Run("A question about intent stops the loop at once", func(t *testing.T) {
		classified := classify(
			finding("f1", findings.Error, findings.AskUser),
			finding("f2", findings.Error, findings.AutoFix),
			finding("f3", findings.Warning, findings.AutoFix),
			finding("f4", findings.Warning, findings.AutoFix),
			finding("f5", findings.Error, findings.AutoFix),
			finding("f6", findings.Warning, findings.AutoFix),
		)
		if got := dialogue.After(classified, 1); got != runstate.Escalated {
			t.Errorf("one question about intent among five mechanical defects is %q at the first round, want %q", got, runstate.Escalated)
		}
	})

	t.Run("Mechanical findings park for the author", func(t *testing.T) {
		classified := classify(finding("f1", findings.Error, findings.AutoFix), finding("f2", findings.Warning, findings.AutoFix))
		if got := dialogue.After(classified, 1); got != runstate.AwaitingResponse {
			t.Errorf("blocking findings that are all patchable are %q, want %q", got, runstate.AwaitingResponse)
		}
	})

	t.Run("The cap escalates", func(t *testing.T) {
		classified := classify(finding("f1", findings.Error, findings.AutoFix))
		if got := dialogue.After(classified, dialogue.Cap); got != runstate.Escalated {
			t.Errorf("blocking findings at the final permitted round are %q, want %q", got, runstate.Escalated)
		}
	})

	t.Run("the whole input space", func(t *testing.T) {
		kinds := everyKind()

		for selection := range 1 << len(kinds) {
			held := selected(kinds, selection)
			classified := classify(held...)

			for round := range dialogue.Cap + 2 {
				asked := round + 1
				got := dialogue.After(classified, asked)

				if !slices.Contains([]runstate.State{runstate.Clean, runstate.Escalated, runstate.AwaitingResponse}, got) {
					t.Fatalf("round %d of %v is %q, which is not an outcome of a round", asked, describe(held), got)
				}
				switch {
				case asksAPerson(held):
					if got != runstate.Escalated {
						t.Fatalf("round %d of %v holds a question about intent and is %q, want %q", asked, describe(held), got, runstate.Escalated)
					}
				case !blocks(held) && !patchable(held):
					if got != runstate.Clean {
						t.Fatalf("round %d of %v blocks nothing and patches nothing and is %q, want %q", asked, describe(held), got, runstate.Clean)
					}
				case asked >= dialogue.Cap:
					if got != runstate.Escalated {
						t.Fatalf("round %d of %v is at the cap and is %q, want %q", asked, describe(held), got, runstate.Escalated)
					}
				default:
					if got != runstate.AwaitingResponse {
						t.Fatalf("round %d of %v is %q, want %q", asked, describe(held), got, runstate.AwaitingResponse)
					}
				}
			}
		}
	})
}

func TestTheRoundCapIsTwo(t *testing.T) {
	if dialogue.Cap != 2 {
		t.Fatalf("the round cap is %d, want 2", dialogue.Cap)
	}
	if !strings.Contains(dialogue.CapIsForContextRatherThanCost, "compacting") {
		t.Errorf("the cap does not record why it is what it is: %q", dialogue.CapIsForContextRatherThanCost)
	}

	blocking := classify(finding("f1", findings.Error, findings.AutoFix))

	t.Run("A second round is permitted", func(t *testing.T) {
		if got := dialogue.After(blocking, 1); got != runstate.AwaitingResponse {
			t.Errorf("the first round is %q, so no second round starts", got)
		}
		if dialogue.CapReached(1) {
			t.Error("the first round already reads as the cap")
		}
	})

	t.Run("A third round is not", func(t *testing.T) {
		if got := dialogue.After(blocking, 2); got != runstate.Escalated {
			t.Errorf("the second round with blocking findings is %q, want %q", got, runstate.Escalated)
		}
		if !dialogue.CapReached(2) {
			t.Error("the second round does not read as the cap")
		}
	})
}

func TestTheRuleReachesNothingOutward(t *testing.T) {
	forbidden := []string{
		"embed",
		"io/fs",
		"io/ioutil",
		"log",
		"math/rand",
		"os",
		"os/exec",
		"path/filepath",
		"syscall",
		"time",
	}

	for _, imported := range source.Imports(t, source.Root(t), filepath.Join("internal", "dialogue")) {
		reaching := slices.Contains(forbidden, imported) ||
			strings.HasPrefix(imported, "os/") ||
			strings.HasPrefix(imported, "net") ||
			strings.HasPrefix(imported, "golang.org/x/sys")
		if reaching {
			t.Errorf("the package holding the round rule imports %s", imported)
		}
	}
}

type kind struct {
	severity findings.Severity
	action   findings.Action
}

func everyKind() []kind {
	var kinds []kind
	for _, severity := range everySeverity {
		for _, action := range everyAction {
			kinds = append(kinds, kind{severity: severity, action: action})
		}
	}
	return kinds
}

func selected(kinds []kind, selection int) []findings.Finding {
	var held []findings.Finding
	for position, one := range kinds {
		if selection&(1<<position) == 0 {
			continue
		}
		held = append(held, finding("f"+string(rune('a'+position)), one.severity, one.action))
	}
	return held
}

func describe(held []findings.Finding) []string {
	described := make([]string, 0, len(held))
	for _, one := range held {
		described = append(described, string(one.Severity)+"/"+string(one.Action))
	}
	return described
}

func asksAPerson(held []findings.Finding) bool {
	for _, one := range held {
		if one.EffectiveAction() == findings.AskUser {
			return true
		}
	}
	return false
}

func blocks(held []findings.Finding) bool {
	for _, one := range held {
		if one.Severity.Blocks() {
			return true
		}
	}
	return false
}

func patchable(held []findings.Finding) bool {
	for _, one := range held {
		if one.EffectiveAction() == findings.AutoFix {
			return true
		}
	}
	return false
}

func finding(identifier string, severity findings.Severity, action findings.Action) findings.Finding {
	return findings.Finding{
		ID:          identifier,
		Severity:    severity,
		File:        "internal/added.go",
		Line:        1,
		Description: "what is wrong",
		Action:      action,
	}
}

func classify(held ...findings.Finding) findings.Classification {
	return findings.Classify(findings.Document{
		Findings:      held,
		RiskLevel:     findings.LowRisk,
		RiskRationale: "nothing much",
	})
}
