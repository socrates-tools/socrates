package skill_test

import (
	"strings"
	"testing"

	"github.com/rui-ktei/socrates/internal/skill"
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
	document := skill.Document()

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

func TestTheSkillNamesTheIllustrationAsSomethingADecisionMayCarry(t *testing.T) {
	page := strings.ToLower(sectionOf(t, skill.Document(), pageHeading))

	for _, said := range []string{
		"a decision may also carry one illustration",
		"beside the evidence rather than inside it",
		`"fragment": {"markup"`,
	} {
		if !strings.Contains(page, strings.ToLower(said)) {
			t.Errorf("the page section does not say %q, so nothing says how an illustration reaches the page:\n%s", said, page)
		}
	}
}

func TestTheSkillSaysWhatFormsAnIllustrationTakes(t *testing.T) {
	page := strings.ToLower(sectionOf(t, skill.Document(), pageHeading))

	for _, said := range []string{
		"one of two forms",
		"inline svg",
		"is for a workflow",
		"the document's stylesheet paints",
		"a screen the person already knows",
	} {
		if !strings.Contains(page, said) {
			t.Errorf("the page section does not say %q, so a session cannot tell which form is for what:\n%s", said, page)
		}
	}

	document := strings.ToLower(skill.Document())
	for _, engine := range []string{
		"mermaid", "graphviz", "plantuml", "kroki", "excalidraw",
		"diagram engine", "lays it out for you", "renders it for you",
	} {
		if strings.Contains(document, engine) {
			t.Errorf("the skill names %q, and the surface embeds no engine that would lay a drawing out", engine)
		}
	}
}

func TestTheSkillSaysTheLookIsStatedOncePerDocument(t *testing.T) {
	page := strings.ToLower(sectionOf(t, skill.Document(), pageHeading))

	for _, said := range []string{
		`"stylesheet"`,
		"once for the whole document rather than once per decision",
		"one stylesheet paints every illustration in it",
	} {
		if !strings.Contains(page, said) {
			t.Errorf("the page section does not say %q, so a session repeats one look per decision:\n%s", said, page)
		}
	}
}

func TestTheSkillSaysWhatIsStrippedBeforeEitherIsEmbedded(t *testing.T) {
	page := strings.ToLower(sectionOf(t, skill.Document(), pageHeading))

	for _, said := range []string{
		"reduced to an allowlist before they are embedded",
		"out of the markup go script, event handlers, remote references and form controls",
		"out of the stylesheet go imports, font faces",
		"every property the allowlist does not name",
		"every value naming a resource",
	} {
		if !strings.Contains(page, said) {
			t.Errorf("the page section does not say %q, so a session writes what will be dropped:\n%s", said, page)
		}
	}
}

func TestTheSkillSaysTheSurfaceDoesTheScoping(t *testing.T) {
	page := strings.ToLower(sectionOf(t, skill.Document(), pageHeading))

	for _, said := range []string{
		"the surface scopes the stylesheet itself",
		"confined to the illustration",
		"leave the scoping to the surface",
		"do not expect a rule you aimed at the page to reach the page",
	} {
		if !strings.Contains(page, said) {
			t.Errorf("the page section does not say %q, so a session scopes its own selectors or expects one to reach the page:\n%s", said, page)
		}
	}
}

func TestHowToIllustrateIsStatedBesideWhenTo(t *testing.T) {
	document := skill.Document()

	bar := strings.Index(document, skill.PageBar)
	mechanism := strings.Index(document, "A decision may also carry one illustration")
	rest := strings.Index(document, "The findings a person answers go on the page in full.")
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
