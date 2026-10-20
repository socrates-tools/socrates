package e2e_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	designskill "github.com/socrates-tools/socrates/design/skill"
	"github.com/socrates-tools/socrates/internal/exit"
	skillkit "github.com/socrates-tools/socrates/internal/skill"
)

func evidenceDocument(blocks string) string {
	return `{
  "title": "two shapes for the component",
  "decisions": [
    {
      "id": "q-1",
      "question": "which shape should the component take",
      "mode": "single-choice",
      "options": [{"label": "as a prop", "meaning": "options are data"}, {"label": "as children", "meaning": "options are elements"}],
      "evidence": ` + blocks + `
    }
  ]
}`
}

const twoCaptionedBlocks = `[
  {"caption": "as a prop", "language": "go", "code": "func Select(options []Option) {}\n"},
  {"caption": "as children", "language": "python", "code": "def select(children): pass\n"}
]`

const aBlockOfEachKind = `[
  {"caption": "as a prop", "language": "go", "code": "func Select(options []Option) {}\n"},
  {"language": "python", "code": "def select(children): pass\n"},
  {"caption": "the lock", "patch": "--- a/internal/cache/write.go\n+++ b/internal/cache/write.go\n@@ -1,1 +1,1 @@\n-a\n+b\n"},
  {"patch": "--- a/internal/cache/read.go\n+++ b/internal/cache/read.go\n@@ -1,1 +1,1 @@\n-a\n+b\n"}
]`

const blockHeadings = `Array.from(document.querySelectorAll('.evidence:not(.illustrated)')).map((held) => held.querySelector('.evidence-label').innerText.replace(/\s+/g, ' ').trim())`

func TestADecisionListingTwoCodeBlocksShowsBothInTheOrderWritten(t *testing.T) {
	t.Parallel()
	_, page := read(t, evidenceDocument(twoCaptionedBlocks))

	var shown []string
	page.evaluate(`Array.from(document.querySelectorAll('.decision[data-decision="q-1"] .evidence.code pre')).map((held) => held.innerText.trim())`, &shown)

	want := []string{"func Select(options []Option) {}", "def select(children): pass"}
	if diff := cmp.Diff(want, shown); diff != "" {
		t.Errorf("the card does not show both blocks in the order written (-want +got):\n%s", diff)
	}
}

func TestEachBlockHeadingShowsItsCaptionBesideItsLanguageOrItsFile(t *testing.T) {
	t.Parallel()
	_, page := read(t, evidenceDocument(aBlockOfEachKind))

	var headings []string
	page.evaluate(blockHeadings, &headings)

	want := []string{
		"as a prop go",
		"python",
		"the lock internal/cache/write.go",
		"internal/cache/read.go",
	}
	if diff := cmp.Diff(want, headings); diff != "" {
		t.Errorf("the headings do not show the caption where there is one with the language or the file beside it, and the language or file alone where there is none (-want +got):\n%s", diff)
	}
}

func TestServeRefusesEvidenceWrittenAsAnObject(t *testing.T) {
	// scenario: A skill installed from an earlier release still writes the object
	t.Parallel()
	held := standing(t)

	printed, failures, code := held.reading(evidenceDocument(`{"language": "go", "code": "x := 1\n"}`), "serve", "--document", "-")
	if code == exit.NothingAsked {
		t.Fatalf("a document whose evidence is an object was served at %q", printed)
	}
	for _, said := range []string{"q-1", "list of blocks", "install the skill again"} {
		if !strings.Contains(failures, said) {
			t.Errorf("refusing evidence written as an object does not say %q: %s", said, failures)
		}
	}
	if strings.TrimSpace(printed) != "" {
		t.Errorf("refusing evidence written as an object printed %q, and an address would say it was served", printed)
	}
	if _, err := os.Stat(filepath.Join(held.surfaceHome(), "surface.json")); err == nil {
		t.Error("a refused document left a session behind, and it starts no server")
	}
}

func TestServeRefusesABlockHoldingBothCodeAndAPatch(t *testing.T) {
	t.Parallel()
	held := standing(t)

	blocks := `[{"language": "go", "code": "x := 1\n", "patch": "--- a/x\n+++ b/x\n@@ -1,1 +1,1 @@\n-a\n+b\n"}]`
	printed, failures, code := held.reading(evidenceDocument(blocks), "serve", "--document", "-")
	if code == exit.NothingAsked {
		t.Fatalf("a block holding both code and a patch was served at %q", printed)
	}
	if !strings.Contains(failures, "q-1") {
		t.Errorf("refusing a block holding both code and a patch does not name the decision: %s", failures)
	}
}

func TestServeRefusesABlockHoldingNeitherCodeNorAPatch(t *testing.T) {
	t.Parallel()
	held := standing(t)

	printed, failures, code := held.reading(evidenceDocument(`[{"caption": "empty", "language": "go"}]`), "serve", "--document", "-")
	if code == exit.NothingAsked {
		t.Fatalf("a block holding neither code nor a patch was served at %q", printed)
	}
	if !strings.Contains(failures, "q-1") {
		t.Errorf("refusing a block holding neither code nor a patch does not name the decision: %s", failures)
	}
}

var cardsTheTicketsSkillShows = regexp.MustCompile(`(?s)\n    (\{"title": .*?\]\}\]\})\n`)

func TestAnExampleTheTicketsSkillShowsIsServedAndShowsTheTicketsFile(t *testing.T) {
	// scenario: A ticket round served from what socrates-tickets teaches
	t.Parallel()
	shown := cardsTheTicketsSkillShows.FindStringSubmatch(designskill.TicketsDocument(skillkit.ClaudeCode))
	if shown == nil {
		t.Fatal("the tickets skill shows no document for a breakdown")
	}

	_, page := read(t, shown[1])

	var ticket string
	page.evaluate(`document.querySelector('.decision[data-decision="2.1-parse"] .evidence.code pre').innerText`, &ticket)
	if !strings.Contains(ticket, "# 2.1:") || !strings.Contains(ticket, "## What to build") {
		t.Errorf("the card the tickets skill shows does not show its ticket's file: %q", ticket)
	}
}

var evidenceTheDesignSkillShows = regexp.MustCompile(`(?s)\n    \{"evidence": (\[.*?\])\}\n`)

func TestTheCardTheDesignSkillShowsForTwoShapesIsServedAsTwoCaptionedHighlightedTsxBlocks(t *testing.T) {
	// scenario: A card comparing two shapes served from what socrates-design teaches
	t.Parallel()
	for _, host := range []skillkit.Host{skillkit.ClaudeCode, skillkit.Codex} {
		t.Run(host.Flag, func(t *testing.T) {
			t.Parallel()
			shown := evidenceTheDesignSkillShows.FindStringSubmatch(designskill.Document(host))
			if shown == nil {
				t.Fatal("the design skill shows no evidence for a card comparing two shapes")
			}

			held := standing(t)
			page := opening(t, held.serve(evidenceDocument(shown[1])))
			page.waitFor(`document.querySelectorAll('.evidence.code code.hljs').length === 2`)

			var read map[string]bool
			page.evaluate(highlightedBlocks, &read)
			if diff := cmp.Diff(map[string]bool{"language-tsx": true}, read); diff != "" {
				t.Errorf("the blocks the design skill shows are not highlighted as tsx (-want +got):\n%s", diff)
			}

			var captioned []string
			page.evaluate(`Array.from(document.querySelectorAll('.evidence.code .evidence-caption')).map((held) => held.innerText.trim())`, &captioned)
			if len(captioned) != 2 || captioned[0] == "" || captioned[1] == "" || captioned[0] == captioned[1] {
				t.Errorf("the blocks the design skill shows are not two captioned blocks: %q", captioned)
			}
		})
	}
}

func TestEveryBlockOfEvidencePaintsTheSameBackgroundAsTheCardItSitsOn(t *testing.T) {
	// scenario: A code block and a patch block on one card are painted the same white as the card
	t.Parallel()
	_, page := read(t, evidenceDocument(aBlockOfEachKind))

	var painted struct {
		Card   string   `json:"card"`
		Blocks []string `json:"blocks"`
	}
	page.evaluate(`({
  card: getComputedStyle(document.querySelector('.decision')).backgroundColor,
  blocks: Array.from(document.querySelectorAll('.evidence:not(.illustrated)')).map((held) => getComputedStyle(held.querySelector('pre, .hunk')).backgroundColor === 'rgba(0, 0, 0, 0)' ? getComputedStyle(held).backgroundColor : getComputedStyle(held.querySelector('pre, .hunk')).backgroundColor),
})`, &painted)

	for position, block := range painted.Blocks {
		if block != painted.Card {
			t.Errorf("block %d of the evidence is painted %s on a card painted %s", position+1, block, painted.Card)
		}
	}
}
