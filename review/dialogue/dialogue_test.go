package dialogue_test

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/internal/testing/source"
	"github.com/socrates-tools/socrates/review/dialogue"
	"github.com/socrates-tools/socrates/review/findings"
	"github.com/socrates-tools/socrates/review/runstate"
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
		held := []findings.Finding{
			finding("f1", findings.Error, findings.AskUser),
			finding("f2", findings.Error, findings.AutoFix),
			finding("f3", findings.Warning, findings.AutoFix),
			finding("f4", findings.Warning, findings.AutoFix),
			finding("f5", findings.Error, findings.AutoFix),
			finding("f6", findings.Warning, findings.AutoFix),
		}

		if got := dialogue.After(classify(held...), 1); got != runstate.AwaitingResponse {
			t.Errorf("one question about intent among five mechanical defects is %q at the first round, want %q", got, runstate.AwaitingResponse)
		}
		if diff := cmp.Diff([]string{"f1", "f2", "f3", "f4", "f5", "f6"}, (dialogue.Responses{}).Unanswered(document(held...), dialogue.Mine)); diff != "" {
			t.Errorf("a set holding one question and five mechanical defects waits on (-want +got):\n%s", diff)
		}
	})

	t.Run("An informational round finishes the run", func(t *testing.T) {
		classified := classify(finding("f1", findings.Info, findings.AutoFix), finding("f2", findings.Info, findings.NoOp))
		if got := dialogue.After(classified, 1); got != runstate.Clean {
			t.Errorf("a round whose patchable findings all block nothing is %q, want %q: there is nothing the author is being asked to answer", got, runstate.Clean)
		}
	})

	t.Run("A blocking finding nobody can patch goes to the person", func(t *testing.T) {
		alone := finding("f1", findings.Warning, findings.NoOp)
		if got := dialogue.After(classify(alone), 1); got != runstate.Escalated {
			t.Errorf("a finding that holds the change while asking for nothing is %q, want %q: no response addresses it", got, runstate.Escalated)
		}

		for _, beside := range []findings.Finding{
			finding("f2", findings.Error, findings.AutoFix),
			finding("f2", findings.Error, findings.AskUser),
		} {
			alongside := classify(alone, beside)
			if got := dialogue.After(alongside, 1); got != runstate.Escalated {
				t.Errorf("a finding no response addresses is %q when it arrives beside %s, want %q: parking would drop it from the run", got, beside.Action, runstate.Escalated)
			}
			if dialogue.EscalatedByTheCap(alongside, 1) {
				t.Errorf("a first round escalating for a finding no response addresses reads as the cap being why, beside %s", beside.Action)
			}
			if diff := cmp.Diff([]string{"f2"}, (dialogue.Responses{}).Unanswered(document(alone, beside), dialogue.Mine)); diff != "" {
				t.Errorf("what the round waits on beside %s differs (-want +got):\n%s", beside.Action, diff)
			}
		}

		if waiting := (dialogue.Responses{}).Unanswered(document(alone), dialogue.Mine); len(waiting) != 0 {
			t.Errorf("%v is waiting on an answer to a finding presented as asking for nothing", waiting)
		}
	})

	t.Run("A clean second round does not name the cap", func(t *testing.T) {
		classified := classify(finding("f1", findings.Info, findings.NoOp))
		if got := dialogue.After(classified, dialogue.Cap); got != runstate.Clean {
			t.Errorf("the final permitted round with nothing blocking is %q, want %q", got, runstate.Clean)
		}
		if dialogue.EscalatedByTheCap(classified, dialogue.Cap) {
			t.Error("a clean final round reads as escalated by the cap, so its verdict would tell the author something is left for them")
		}
	})

	t.Run("A question alone parks the run", func(t *testing.T) {
		question := finding("f1", findings.Error, findings.AskUser)
		classified := classify(question, finding("f2", findings.Info, findings.NoOp))

		if got := dialogue.After(classified, 1); got != runstate.AwaitingResponse {
			t.Errorf("a round holding one question about intent and nothing else blocking is %q, want %q: the person's answer has nowhere to go otherwise", got, runstate.AwaitingResponse)
		}
		if diff := cmp.Diff([]string{"f1"}, (dialogue.Responses{}).Unanswered(document(question, finding("f2", findings.Info, findings.NoOp)), dialogue.Mine)); diff != "" {
			t.Errorf("what the round waits on differs (-want +got):\n%s", diff)
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

	t.Run("A question left unanswered at the cap escalates", func(t *testing.T) {
		for _, held := range [][]findings.Finding{
			{finding("f1", findings.Error, findings.AskUser)},
			{finding("f1", findings.Error, findings.AskUser), finding("f2", findings.Error, findings.AutoFix)},
		} {
			classified := classify(held...)
			if got := dialogue.After(classified, dialogue.Cap); got != runstate.Escalated {
				t.Errorf("%v at the final permitted round is %q, want %q: no round remains in which an answer could be judged", describe(held), got, runstate.Escalated)
			}
			if !dialogue.EscalatedByTheCap(classified, dialogue.Cap) {
				t.Errorf("%v at the final permitted round escalates without the cap being why, so the verdict would not say what is left for the person", describe(held))
			}
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
				case !blocks(held):
					if got != runstate.Clean {
						t.Fatalf("round %d of %v blocks nothing and is %q, want %q", asked, describe(held), got, runstate.Clean)
					}
				case blocksBeyondAnswering(held):
					if got != runstate.Escalated {
						t.Fatalf("round %d of %v holds something no response addresses and is %q, want %q", asked, describe(held), got, runstate.Escalated)
					}
				case asked >= dialogue.Cap:
					if got != runstate.Escalated {
						t.Fatalf("round %d of %v is at the cap and is %q, want %q", asked, describe(held), got, runstate.Escalated)
					}
				default:
					if got != runstate.AwaitingResponse {
						t.Fatalf("round %d of %v is %q, want %q", asked, describe(held), got, runstate.AwaitingResponse)
					}
					waiting := (dialogue.Responses{}).Unanswered(document(held...), dialogue.Mine)
					if diff := cmp.Diff(blockingIdentifiers(held), waiting); diff != "" {
						t.Fatalf("round %d of %v parks and waits on (-want +got):\n%s", asked, describe(held), diff)
					}
				}
			}
		}
	})
}

func TestWhichRoundsWaitOnAPerson(t *testing.T) {
	t.Run("An escalated round waits on a person", func(t *testing.T) {
		beyondAnswering := classify(finding("f1", findings.Error, findings.NoOp))
		if !dialogue.AwaitsAPerson(beyondAnswering, dialogue.After(beyondAnswering, 1)) {
			t.Error("a round escalating for a finding that holds the change while asking for nothing waits on nobody, and nobody else can act on it")
		}

		atTheCap := classify(finding("f1", findings.Error, findings.AutoFix))
		if !dialogue.AwaitsAPerson(atTheCap, dialogue.After(atTheCap, dialogue.Cap)) {
			t.Error("a round escalating at the cap waits on nobody, and no round remains in which anything else could act on it")
		}
	})

	t.Run("A parked round holding a question waits on a person", func(t *testing.T) {
		alone := classify(finding("f1", findings.Error, findings.AskUser))
		if !dialogue.AwaitsAPerson(alone, dialogue.After(alone, 1)) {
			t.Error("a parked round holding one question about intent waits on nobody, and only a person closes it")
		}

		alongside := classify(
			finding("f1", findings.Error, findings.AskUser),
			finding("f2", findings.Error, findings.AutoFix),
			finding("f3", findings.Warning, findings.AutoFix),
		)
		if !dialogue.AwaitsAPerson(alongside, dialogue.After(alongside, 1)) {
			t.Error("a parked round holding a question beside mechanical defects waits on nobody, and the question is still the person's")
		}
	})

	t.Run("A parked round the author answers alone waits on nobody", func(t *testing.T) {
		classified := classify(finding("f1", findings.Error, findings.AutoFix), finding("f2", findings.Warning, findings.AutoFix))
		if dialogue.AwaitsAPerson(classified, dialogue.After(classified, 1)) {
			t.Error("a parked round whose findings the author patches waits on a person, so the session stops for an answer nobody was asked for")
		}
	})

	t.Run("A clean round waits on nobody", func(t *testing.T) {
		for _, classified := range []findings.Classification{
			classify(),
			classify(finding("f1", findings.Info, findings.NoOp)),
			classify(finding("f1", findings.Info, findings.AskUser)),
			classify(finding("f1", findings.Info, findings.AutoFix)),
		} {
			for _, round := range []int{1, dialogue.Cap} {
				if dialogue.AwaitsAPerson(classified, dialogue.After(classified, round)) {
					t.Errorf("round %d blocks nothing and waits on a person, and nothing is being asked of anybody", round)
				}
			}
		}
	})

	t.Run("The rule is decided over its whole input space", func(t *testing.T) {
		kinds := everyKind()

		for selection := range 1 << len(kinds) {
			held := selected(kinds, selection)
			classified := classify(held...)

			for _, state := range runstate.States() {
				got := dialogue.AwaitsAPerson(classified, state)

				want := false
				switch state {
				case runstate.Escalated:
					want = true
				case runstate.AwaitingResponse:
					want = blocksAskingAPerson(held)
				}
				if got != want {
					t.Fatalf("%v in state %q waits on a person: %t, want %t", describe(held), state, got, want)
				}
				if !state.Terminal() && state != runstate.AwaitingResponse && got {
					t.Fatalf("%v in state %q waits on a person, and no round has led anywhere yet", describe(held), state)
				}
			}

			for round := range dialogue.Cap + 2 {
				asked := round + 1
				reached := dialogue.After(classified, asked)
				got := dialogue.AwaitsAPerson(classified, reached)

				want := false
				switch {
				case !blocks(held):
					want = false
				case blocksBeyondAnswering(held), asked >= dialogue.Cap:
					want = true
				default:
					want = blocksAskingAPerson(held)
				}
				if got != want {
					t.Fatalf("round %d of %v led to %q and waits on a person: %t, want %t", asked, describe(held), reached, got, want)
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

	for _, imported := range source.Imports(t, source.Root(t), filepath.Join("review", "dialogue")) {
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

func blockingIdentifiers(held []findings.Finding) []string {
	var named []string
	for _, one := range held {
		if slices.Contains([]string{"error", "warning"}, string(one.Severity)) {
			named = append(named, one.ID)
		}
	}
	return named
}

func blocks(held []findings.Finding) bool {
	for _, one := range held {
		if slices.Contains([]string{"error", "warning"}, string(one.Severity)) {
			return true
		}
	}
	return false
}

func blocksBeyondAnswering(held []findings.Finding) bool {
	for _, one := range held {
		if slices.Contains([]string{"error", "warning"}, string(one.Severity)) && string(one.Action) == "no-op" {
			return true
		}
	}
	return false
}

func blocksAskingAPerson(held []findings.Finding) bool {
	for _, one := range held {
		if slices.Contains([]string{"error", "warning"}, string(one.Severity)) && string(one.Action) == "ask-user" {
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
