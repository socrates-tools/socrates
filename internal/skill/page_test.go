package skill_test

import (
	"strings"
	"testing"

	"github.com/rui-ktei/socrates/internal/skill"
)

const (
	pageHeading      = "## The page a person reads"
	escalatedHeading = "### escalated"
	parkedHeading    = "### awaiting-response"
	waitingStep      = "poll scrat status"
)

func TestTheBarIsStatedOnceAndRenderedFromThere(t *testing.T) {
	document := skill.Document()

	if !strings.Contains(document, skill.PageBar) {
		t.Fatalf("the skill does not render the bar its source states:\n%s", skill.PageBar)
	}
	for line := range strings.SplitSeq(skill.PageBar, "\n") {
		stated := strings.TrimSpace(line)
		if stated == "" {
			continue
		}
		if counted := strings.Count(document, stated); counted != 1 {
			t.Errorf("the skill states %q %d times, and two copies of a bar drift", stated, counted)
		}
	}
}

func TestTheBarNamesEachWayThePageFailsAPerson(t *testing.T) {
	bar := strings.ToLower(skill.PageBar)

	for _, required := range []struct {
		named  string
		stated []string
	}{
		{named: "getting to the point", stated: []string{"get to the point", "length is not thoroughness"}},
		{named: "plain english about impact", stated: []string{"plain english", "what the finding means for the change"}},
		{named: "an example where a finding is purely technical", stated: []string{"purely technical", "example"}},
		{named: "a diagram only where a workflow needs one", stated: []string{"diagram only where", "not otherwise"}},
		{named: "restraint about tokens", stated: []string{"spend tokens"}},
	} {
		t.Run(required.named, func(t *testing.T) {
			for _, stated := range required.stated {
				if !strings.Contains(bar, stated) {
					t.Errorf("the bar does not say %q, and that is how a page fails a person by %s", stated, required.named)
				}
			}
		})
	}
}

func TestTheBarIsStatedBeforeTheBranchesThatReferToIt(t *testing.T) {
	document := skill.Document()

	bar := strings.Index(document, skill.PageBar)
	if bar < 0 {
		t.Fatal("the skill does not state the bar at all")
	}
	for _, branch := range []string{escalatedHeading, parkedHeading} {
		at := strings.Index(document, branch)
		if at < 0 {
			t.Fatalf("the skill has no %s branch", branch)
		}
		if bar > at {
			t.Errorf("the bar is stated after %s, so the branch sends a session to make a page it has not been told how to make", branch)
		}
	}

	for _, branch := range []string{escalatedHeading, parkedHeading} {
		if strings.Contains(sectionOf(t, document, branch), "Length is not thoroughness") {
			t.Errorf("%s restates the bar, and the copy in front of the session is the one that gets softened", branch)
		}
	}
}

func TestBothBranchesSayToRenderThePage(t *testing.T) {
	document := skill.Document()

	for _, branch := range []string{escalatedHeading, parkedHeading} {
		section := strings.ToLower(sectionOf(t, document, branch))
		if !strings.Contains(section, "page") {
			t.Errorf("%s says nothing about a page, and it is a branch where a person is the reader:\n%s", branch, section)
		}
	}

	t.Run("The escalated branch drafts inline", func(t *testing.T) {
		escalated := strings.ToLower(sectionOf(t, document, escalatedHeading))
		for _, said := range []string{"draft it here", "the run is over"} {
			if !strings.Contains(escalated, said) {
				t.Errorf("the escalated branch does not say %q, so nothing says why the drafting is not delegated:\n%s", said, escalated)
			}
		}
		if strings.Contains(escalated, "delegate") {
			t.Errorf("the escalated branch delegates, and a fresh subagent costs its own start-up where nothing follows:\n%s", escalated)
		}
	})

	t.Run("The parked branch delegates", func(t *testing.T) {
		parked := strings.ToLower(sectionOf(t, document, parkedHeading))
		for _, said := range []string{"delegate", "subagent", "read again on every one of those turns", "scrat findings --json"} {
			if !strings.Contains(parked, said) {
				t.Errorf("the parked branch does not say %q, so the drafting lands in the context that is about to patch:\n%s", said, parked)
			}
		}
	})
}

func TestTheSessionDoesNotDecideTheConditionItself(t *testing.T) {
	page := sectionOf(t, skill.Document(), pageHeading)

	if !strings.Contains(page, "read by a person") {
		t.Error("the skill does not name the line socrates prints, so a session has to decide for itself which rounds a person reads")
	}
	for _, invented := range []string{"unanswered", "escalated", "awaiting-response", "count of"} {
		if strings.Contains(strings.ToLower(page), invented) {
			t.Errorf("the page section names %q, which is a condition of its own over what socrates already decided:\n%s", invented, page)
		}
	}
}

func TestARoundTheAuthorAnswersAloneGetsNoPage(t *testing.T) {
	parked := sectionOf(t, skill.Document(), parkedHeading)

	patching, _, found := strings.Cut(parked, "A finding it reports as needing a person")
	if !found {
		t.Fatalf("the parked branch does not split the work in two:\n%s", parked)
	}
	if strings.Contains(strings.ToLower(patching), "page") {
		t.Errorf("the half of the parked branch a session answers alone names a page:\n%s", patching)
	}
	for _, said := range []string{"scrat respond", "scrat continue"} {
		if !strings.Contains(parked, said) {
			t.Errorf("the parked branch no longer says %q", said)
		}
	}
}

func TestThePageIsNotInThePollingLoop(t *testing.T) {
	document := skill.Document()

	step := stepOf(t, document, waitingStep)
	if strings.Contains(strings.ToLower(step), "page") {
		t.Errorf("the step a session runs in a loop names a page, and a page is made once after the run has stopped moving:\n%s", step)
	}
	if counted := strings.Count(document, "scrat watch"); counted != 1 {
		t.Errorf("the skill names scrat watch %d times, want once", counted)
	}
}

func TestTheSkillSaysThePageIsReadOnly(t *testing.T) {
	// scenario: The page carries no control that implies an action
	page := strings.ToLower(sectionOf(t, skill.Document(), pageHeading))

	for _, said := range []string{"collects nothing", "declare no capabilities", "reaches no run", "tells you here"} {
		if !strings.Contains(page, said) {
			t.Errorf("the page section does not say %q, and a control on the page takes a decision that goes nowhere", said)
		}
	}
	for _, control := range []string{"a form", "a button", "a checkbox", "a stored answer"} {
		if !strings.Contains(page, control) {
			t.Errorf("the page section does not name %q as something the page does not carry", control)
		}
	}
}

func TestTheSkillSaysThePageAndTheTerminalUseTheSameIdentifiers(t *testing.T) {
	// scenario: Every finding keeps the identifier socrates gave it
	page := strings.ToLower(sectionOf(t, skill.Document(), pageHeading))

	for _, said := range []string{"identifier", "severity", "audience", "f3"} {
		if !strings.Contains(page, said) {
			t.Errorf("the page section does not say %q, so a person naming a finding leaves the session guessing which one", said)
		}
	}
}

func TestThePageSaysWhatThePersonDoesBackInTheTerminal(t *testing.T) {
	page := strings.ToLower(sectionOf(t, skill.Document(), pageHeading))

	for _, said := range []string{"names its finding", "what answering it looks like here"} {
		if !strings.Contains(page, said) {
			t.Errorf("the page section does not say %q, so a person who has decided is left working out how to say so", said)
		}
	}
}

func TestTheDecisionsAreInFullAndTheRestIsAList(t *testing.T) {
	page := strings.ToLower(sectionOf(t, skill.Document(), pageHeading))

	for _, said := range []string{"in full", "short list"} {
		if !strings.Contains(page, said) {
			t.Errorf("the page section does not say %q, so the only thing the person has to act on is buried among the rest", said)
		}
	}
}

func sectionOf(t *testing.T, document, heading string) string {
	t.Helper()

	_, after, found := strings.Cut(document, heading)
	if !found {
		t.Fatalf("the skill has no %s section", heading)
	}

	depth := strings.Count(strings.Fields(heading)[0], "#")
	for line := range strings.SplitSeq(after, "\n") {
		if !strings.HasPrefix(line, "#") {
			continue
		}
		if strings.Count(strings.Fields(line)[0], "#") > depth {
			continue
		}
		if next := strings.Index(after, "\n"+line); next >= 0 {
			return after[:next]
		}
	}
	return after
}

func stepOf(t *testing.T, document, named string) string {
	t.Helper()

	at := strings.Index(document, named)
	if at < 0 {
		t.Fatalf("the skill has no step naming %q", named)
	}

	start := strings.LastIndex(document[:at], "\n3.")
	if start < 0 {
		t.Fatalf("the step naming %q is not numbered", named)
	}
	end := strings.Index(document[start+1:], "\n4.")
	if end < 0 {
		t.Fatalf("the step naming %q does not end", named)
	}
	return document[start : start+1+end]
}
