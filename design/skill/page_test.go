package skill_test

import (
	skillkit "github.com/socrates-tools/socrates/internal/skill"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/design/skill"
)

const (
	roundsHeading = "## Rounds"
	pageHeading   = "## Putting a round on a page"
	factsHeading  = "## Facts"
)

func TestTheSkillStatesTheFrontierAndTheRound(t *testing.T) {
	document := skill.Document(skillkit.ClaudeCode)

	for _, said := range []string{
		"The frontier is every question whose prerequisites are settled.",
		"A round is the whole frontier",
		"A question whose answer depends on one still open belongs to a later round",
	} {
		if !strings.Contains(document, said) {
			t.Errorf("the skill does not say %q, so a session has no rule for which questions belong in front of a person at once", said)
		}
	}
}

func TestTheAnswerArrivesChosenAndAtMostOneIsMarked(t *testing.T) {
	page := sectionOf(t, skill.Document(skillkit.ClaudeCode), pageHeading)

	for _, said := range []string{
		`"recommended": true`,
		"at most one is ever marked",
		"A document marking two is refused before anything is served",
		"it arrives chosen when the page renders",
	} {
		if !strings.Contains(page, said) {
			t.Errorf("the page section does not say %q, so a session cannot tell whether agreeing with the round costs one press:\n%s", said, page)
		}
	}
}

func TestTheSessionIsMintedOnceAndKeptForTheWholeDesign(t *testing.T) {
	page := sectionOf(t, skill.Document(skillkit.ClaudeCode), pageHeading)

	for _, said := range []string{
		"Mint the session once, with the first",
		"socrates agora serve",
		"so the person is given one address",
		"Keep that session for the whole design",
		"socrates agora serve --session",
	} {
		if !strings.Contains(page, said) {
			t.Errorf("the page section does not say %q, so a person could be handed a second address for a later round:\n%s", said, page)
		}
	}
}

func TestTheSkillSaysWhenAnIllustrationEarnsItsPlace(t *testing.T) {
	page := sectionOf(t, skill.Document(skillkit.ClaudeCode), pageHeading)

	for _, said := range []string{
		"A decision may also carry one illustration, an HTML file you write and name by path",
		"Draw one where the choice is a shape",
		"Skip it where the choice is a name, a threshold or a policy",
		"An illustration that only repeats what the words already said costs the person attention and costs the skill its ceiling for nothing",
	} {
		if !strings.Contains(page, said) {
			t.Errorf("the page section does not say %q, so a session has no rule for when an illustration earns its place:\n%s", said, page)
		}
	}
}

func TestTheSkillSaysCodeOnACardGoesInEvidenceAndNeverInAnIllustrationAlone(t *testing.T) {
	// scenario: A session shows code on a design card
	for _, host := range []skillkit.Host{skillkit.ClaudeCode, skillkit.Codex} {
		page := sectionOf(t, skill.Document(host), pageHeading)

		for _, said := range []string{
			"Code goes in the decision's `evidence`",
			"An illustration never draws code on its own",
			"it may still show code inside a drawn screen or layout",
			"a card that only shows code uses evidence",
			"A before and after is a `patch` block",
		} {
			if !strings.Contains(page, said) {
				t.Errorf("the %s page section does not say %q, so a session hand-draws code in an illustration:\n%s", host.Flag, said, page)
			}
		}
	}
}

func TestFactsAreFoundRatherThanAsked(t *testing.T) {
	facts := sectionOf(t, skill.Document(skillkit.ClaudeCode), factsHeading)

	for _, said := range []string{
		"Facts are found rather than asked for",
		"No question reaches the page that reading the repository would have answered",
	} {
		if !strings.Contains(facts, said) {
			t.Errorf("the facts section does not say %q, so a session could put a question on a page that reading the repository would have answered:\n%s", said, facts)
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

func TestTheRoundsSectionComesBeforeThePageSection(t *testing.T) {
	document := skill.Document(skillkit.ClaudeCode)

	rounds := strings.Index(document, roundsHeading)
	page := strings.Index(document, pageHeading)
	if rounds < 0 || page < 0 {
		t.Fatalf("the skill is missing one of %q or %q", roundsHeading, pageHeading)
	}
	if rounds > page {
		t.Error("the page section is stated before the rounds section explains what a round is")
	}
}

func TestTheSkillNumbersQuestionsFromOneAndNamesEachQuestionN(t *testing.T) {
	page := sectionOf(t, skill.Document(skillkit.ClaudeCode), pageHeading)

	for _, said := range []string{
		`{"id": "q2", "name": "Question 2",`,
		"Number questions from 1, never 0, as they are first asked",
		"identify each |q<N>| and name it |Question <N>|",
	} {
		said = strings.ReplaceAll(said, "|", "`")
		if !strings.Contains(page, said) {
			t.Errorf("the page section does not say %q, so a session titles its cards q0, q1 rather than Question 1, Question 2:\n%s", said, page)
		}
	}
}

func TestTheSkillNamesTheIllustrationAsAnHTMLFileDrawnInASandboxedFrame(t *testing.T) {
	// scenario: Both skills stop teaching the fragment
	page := sectionOf(t, skill.Document(skillkit.ClaudeCode), pageHeading)

	for _, said := range []string{
		`{"illustration": {"path": "q3.html"}}`,
		"An HTML file you write and name by path",
		"drawn in a sandboxed frame",
	} {
		if !strings.Contains(strings.ToLower(page), strings.ToLower(said)) {
			t.Errorf("the page section does not say %q:\n%s", said, page)
		}
	}

	document := strings.ToLower(skill.Document(skillkit.ClaudeCode))
	for _, retired := range []string{"fragment", "markup", "allowlist", `"stylesheet"`} {
		if strings.Contains(document, retired) {
			t.Errorf("the skill still mentions %q, and an illustration is an HTML file named by path", retired)
		}
	}
}

func TestTheSkillTeachesWhatAnIllustrationBringsAndWhatIsRefused(t *testing.T) {
	// scenario: Both skills teach the drawn illustration
	page := strings.ToLower(sectionOf(t, skill.Document(skillkit.ClaudeCode), pageHeading))

	for _, said := range []string{
		"write the files it references beside it",
		"16 mib at most in all",
		"html, css, javascript, json, svg, png, jpeg, webp, gif, woff or woff2",
		"a reference outside the file's folder, a missing file, a type not served or a total over the cap refuses the document, naming the decision and the file",
		"cdnjs.cloudflare.com, cdn.jsdelivr.net and unpkg.com",
		"sends nothing anywhere",
		"/assets/tokens.css",
		"the frame grows to its drawing and can be expanded to fill the window",
	} {
		if !strings.Contains(page, said) {
			t.Errorf("the page section does not say %q, so a drawing could be written that agora refuses or cannot show:\n%s", said, page)
		}
	}
}

func TestAPrototypesResultCanBeThePrototypeItself(t *testing.T) {
	// scenario: Both skills teach the drawn illustration
	document := strings.ToLower(skill.Document(skillkit.ClaudeCode))

	if !strings.Contains(document, "a prototype's result can be the prototype itself") {
		t.Errorf("the skill does not say a prototype can be the illustration, so a prototype's shape is described rather than shown:\n%s", document)
	}
}

func TestThePageSendsNothingBeyondThisMachineAndNoLongerSaysItReachesNowhere(t *testing.T) {
	// scenario: The network promise is restated everywhere it was stated
	document := skill.Document(skillkit.ClaudeCode)

	if !strings.Contains(document, "The page is served on this machine and sends nothing beyond it.") {
		t.Errorf("the skill does not promise that the page sends nothing beyond this machine:\n%s", document)
	}
	if strings.Contains(document, "reaches nowhere") {
		t.Errorf("the skill still says the page reaches nowhere beyond this machine, and a drawing may load from three CDNs")
	}
}
