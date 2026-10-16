package skill_test

import (
	skillkit "github.com/socrates-tools/socrates/internal/skill"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/review/skill"
)

const (
	showing = "what a person sees"
	drawing = "diagram"
)

func barLine(t *testing.T, naming string) string {
	t.Helper()

	var found []string
	for line := range strings.SplitSeq(skill.PageBar, "\n") {
		if strings.HasPrefix(line, "- ") && strings.Contains(strings.ToLower(line), naming) {
			found = append(found, strings.TrimSpace(line))
		}
	}

	if len(found) != 1 {
		t.Fatalf("the bar states %d lines about %q, want one: %v", len(found), naming, found)
	}
	return found[0]
}

func TestShowingAndDrawingAreSeparateLines(t *testing.T) {
	document := skill.Document(skillkit.ClaudeCode)

	shown := barLine(t, showing)
	drawn := barLine(t, drawing)

	if shown == drawn {
		t.Fatalf("one line carries both showing a person their own screen and drawing a diagram: %q", shown)
	}
	if strings.Contains(strings.ToLower(shown), drawing) {
		t.Errorf("the line about what a person sees restates the one about diagrams: %q", shown)
	}
	if strings.Contains(strings.ToLower(drawn), showing) {
		t.Errorf("the line about diagrams restates the one about what a person sees: %q", drawn)
	}
	for _, line := range []string{shown, drawn} {
		if counted := strings.Count(document, line); counted != 1 {
			t.Errorf("the skill states %q %d times, and the copy nobody reads while working is the one that stays", line, counted)
		}
	}
}

func TestTheSkillNamesTheIllustrationAsAnHTMLFileDrawnInASandboxedFrame(t *testing.T) {
	// scenario: Both skills stop teaching the fragment
	document := skill.Document(skillkit.ClaudeCode)
	page := strings.ToLower(sectionOf(t, document, pageHeading))

	for _, said := range []string{
		"a decision may also carry one illustration",
		"beside the evidence rather than inside it",
		`"illustration": {"path": "q3.html"}`,
		"an illustration is an html file you write and name by path",
		"drawn in a sandboxed frame",
	} {
		if !strings.Contains(page, strings.ToLower(said)) {
			t.Errorf("the page section does not say %q, so nothing says how an illustration reaches the page:\n%s", said, page)
		}
	}

	for _, retired := range []string{"fragment", "markup", "allowlist", `"stylesheet"`} {
		if strings.Contains(strings.ToLower(document), retired) {
			t.Errorf("the skill still mentions %q, and an illustration is an HTML file named by path", retired)
		}
	}
}

func TestHowToIllustrateIsStatedBesideWhenTo(t *testing.T) {
	document := skill.Document(skillkit.ClaudeCode)

	bar := strings.Index(document, skill.PageBar)
	mechanism := strings.Index(document, "A decision may also carry one illustration")
	rest := strings.Index(document, "Give every finding in the round a card of its own.")
	if bar < 0 || mechanism < 0 || rest < 0 {
		t.Fatal("the skill states either the bar or the mechanism nowhere at all")
	}
	if mechanism < bar+len(skill.PageBar) || mechanism > rest {
		t.Error("the mechanism is not stated beside the bar, so a session reads when and how in two places")
	}

	beside := document[mechanism:rest]
	for _, judged := range []string{barLine(t, showing), barLine(t, drawing), "not otherwise", "only where"} {
		if strings.Contains(beside, judged) {
			t.Errorf("the mechanism restates %q, and the copy beside the mechanism is the one that gets softened:\n%s", judged, beside)
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

func TestAFindingDecidedByLookingCarriesABeforeAndAfterDrawingOrWordsWhereThereIsNoBrowser(t *testing.T) {
	// scenario: A UI finding is drawn as a before and after pair
	document := skill.Document(skillkit.ClaudeCode)
	mechanism := strings.Index(document, "A decision may also carry one illustration")
	rest := strings.Index(document, "Give every finding in the round a card of its own.")
	if mechanism < 0 || rest < mechanism {
		t.Fatal("the skill states the mechanism nowhere")
	}
	beside := strings.ToLower(document[mechanism:rest])

	for _, said := range []string{
		"a finding decided by looking carries a before and after drawing",
		"capture the screen before the change and after it with whatever browser you have",
		"save the screenshots beside the html file",
		"side by side",
		"where you have no browser, say what changed in words",
		"socrates review captures nothing itself",
	} {
		if !strings.Contains(beside, said) {
			t.Errorf("the mechanism does not say %q, so a finding decided by looking is described instead of shown:\n%s", said, beside)
		}
	}
}

func TestThePageSendsNothingBeyondThisMachineAndNoLongerSaysItReachesNowhere(t *testing.T) {
	// scenario: The network promise is restated everywhere it was stated
	for name, document := range map[string]string{
		"review":   skill.Document(skillkit.ClaudeCode),
		"handover": skill.HandoverDocument(skillkit.ClaudeCode),
	} {
		if !strings.Contains(document, "The page is served on this machine and sends nothing beyond it.") {
			t.Errorf("the %s skill does not promise that the page sends nothing beyond this machine", name)
		}
		if strings.Contains(document, "reaches nowhere") {
			t.Errorf("the %s skill still says the page reaches nowhere beyond this machine, and a drawing may load from three CDNs", name)
		}
	}
}
