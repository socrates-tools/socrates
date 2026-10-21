package e2e_test

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/agora/document"
)

const inTwoLanguagesDocument = `{
  "title": "two things to settle",
  "decisions": [
    {
      "id": "c-1",
      "question": "should the settlement be read before the lock is taken",
      "mode": "single-choice",
      "options": [{"label": "take the lock first", "meaning": "correctness first", "command": "patch c-1"}],
      "evidence": [{
        "language": "csharp",
        "code": "public sealed class Settlement\n{\n    public string Reference { get; init; } = \"unset\";\n}\n"
      }]
    },
    {
      "id": "c-2",
      "question": "is the change to the reader the one that was wanted",
      "mode": "single-choice",
      "options": [{"label": "yes", "meaning": "that is the change", "command": "agree c-2"}],
      "evidence": [{
        "language": "diff",
        "code": "--- a/internal/cache/read.go\n+++ b/internal/cache/read.go\n@@ -1,2 +1,2 @@\n-\tdelete(locks, key)\n+\tdefer unlock(key)\n"
      }]
    }
  ]
}`

const highlightedBlocks = `Object.fromEntries(Array.from(document.querySelectorAll('.evidence.code ui-code-block')).map((held) => [held.language, held.highlighted]))`

func TestTheLanguagesACallerActuallyPassesAreHighlighted(t *testing.T) {
	t.Parallel()
	held := standing(t)
	page := opening(t, held.serve(inTwoLanguagesDocument))

	page.waitFor(highlightedCodeBlocks(2))

	var read map[string]bool
	page.evaluate(highlightedBlocks, &read)

	want := map[string]bool{"csharp": true, "diff": true}
	if diff := cmp.Diff(want, read); diff != "" {
		t.Errorf("what the page highlighted differs (-want +got):\n%s\na language the caller named and the page ignored renders as a grey block with nothing saying so", diff)
	}
}

const colouredDiffDocument = `{
  "title": "one thing to settle",
  "decisions": [
    {"id": "c-1", "question": "should the write take the lock", "mode": "single-choice",
     "options": [{"label": "take it", "meaning": "correctness first"}, {"label": "leave it", "meaning": "acceptable here"}],
     "evidence": [{"patch": "--- a/internal/cache/write.go\n+++ b/internal/cache/write.go\n@@ -10,3 +10,4 @@\n func Write(key string, entry Entry) {\n-\t// store it, unlocked (fine for now)\n+\t// store it, under the lock\n+\tmu.Lock()\n \treturn\n }\n"}]}
  ]
}`

func TestTheLinesOfADiffAreColouredLikeCodeAndStillMeetTheFloor(t *testing.T) {
	t.Parallel()
	_, page := read(t, colouredDiffDocument)

	page.waitFor(`document.querySelector('ui-diff-block').highlighted`)

	for _, kind := range []string{"removed", "added"} {
		var comments []string
		page.evaluate(`document.querySelector('ui-diff-block').drawnHunks.flatMap((hunk) => hunk.rows).map((row) => row.`+map[string]string{"removed": "left", "added": "right"}[kind]+`).filter((side) => side.kind === '`+kind+`').flatMap((side) => side.tokens).filter((token) => token.text.startsWith('//')).map((token) => token.ink)`, &comments)
		if len(comments) == 0 {
			t.Errorf("a comment on a %s line of a Go diff is not coloured as one", kind)
		}
		for _, ink := range comments {
			if ink != "quiet-on-change" {
				t.Errorf("a comment on a %s line of a Go diff is drawn in the %s ink, want the quiet ink of a change", kind, ink)
			}
		}
	}
	if inked := page.colouredPixelsOf("ui-diff-block", "--token-keyword"); inked == 0 {
		t.Error("a keyword in a Go diff is painted in no pixel of the keyword token")
	}

	var weighedDiff weighed
	page.evaluate(measuringTheContrast, &weighedDiff)
	for _, failing := range weighedDiff.Failing {
		t.Errorf("%s, and the floor the page is held to is %v", failing, contrastFloor)
	}
}

const namedByAliasDocument = `{
  "title": "shapes named by alias",
  "decisions": [
    {
      "id": "c-1",
      "question": "which shape of the component",
      "mode": "single-choice",
      "options": [{"label": "the first", "meaning": "the first shape", "command": "take c-1"}],
      "evidence": [
        {"language": "tsx", "code": "export const Select = () => <select className=\"s\">{items}</select>\n"},
        {"language": "jsx", "code": "const Select = () => <select className=\"s\" />\n"},
        {"language": "ts", "code": "export const limit: number = 3\n"},
        {"language": "sh", "code": "if [ -f \"$file\" ]; then echo \"$file\"; fi\n"},
        {"language": "yml", "code": "name: build\non: push\n"},
        {"patch": "--- a/select.tsx\n+++ b/select.tsx\n@@ -1,1 +1,1 @@\n-export const limit: number = 2\n+export const limit: number = 3\n"}
      ]
    }
  ]
}`

func TestCodeNamedByAnAliasIsHighlightedAsItsLanguageIs(t *testing.T) {
	t.Parallel()
	held := standing(t)
	page := opening(t, held.serve(namedByAliasDocument))

	page.waitFor(highlightedCodeBlocks(5))

	var read map[string]bool
	page.evaluate(highlightedBlocks, &read)

	want := map[string]bool{
		"tsx": true,
		"jsx": true,
		"ts":  true,
		"sh":  true,
		"yml": true,
	}
	if diff := cmp.Diff(want, read); diff != "" {
		t.Errorf("what the page highlighted differs (-want +got):\n%s", diff)
	}
}

func TestAPatchWhoseLanguageIsNamedByAnAliasHasItsLinesHighlighted(t *testing.T) {
	t.Parallel()
	held := standing(t)
	page := opening(t, held.serve(patchInAnAliasDocument))

	page.waitFor(`document.querySelector('ui-diff-block').highlighted`)

	var inks []string
	page.evaluate(`document.querySelector('ui-diff-block').drawnHunks.flatMap((hunk) => hunk.rows).map((row) => row.right).filter((side) => side.kind === 'added').flatMap((side) => side.tokens).map((token) => token.ink)`, &inks)
	if !slices.Contains(inks, "keyword") {
		t.Errorf("an added line of a patch whose language is ts is drawn in the inks %v, want its keyword coloured as one", inks)
	}
}

const patchInAnAliasDocument = `{
  "title": "one thing to settle",
  "decisions": [
    {"id": "c-1", "question": "should the limit be raised", "mode": "single-choice",
     "options": [{"label": "raise it", "meaning": "more room"}],
     "evidence": [{"language": "ts", "patch": "--- a/limit.txt\n+++ b/limit.txt\n@@ -1,1 +1,1 @@\n-export const limit: number = 2\n+export const limit: number = 3\n"}]}
  ]
}`

func TestEveryLanguageThePageHighlightsIsHighlightedByACodeBlockNamingIt(t *testing.T) {
	// scenario: A code block naming any language or alias the page highlights is highlighted, none is left plain
	t.Parallel()
	var blocks []map[string]string
	for _, named := range document.Languages() {
		blocks = append(blocks, map[string]string{"language": named, "code": "x = 1\n"})
	}
	evidence, err := json.Marshal(blocks)
	if err != nil {
		t.Fatal(err)
	}
	held := standing(t)
	page := opening(t, held.serve(evidenceDocument(string(evidence))))
	page.waitFor(highlightedCodeBlocks(len(blocks)))

	var plain []string
	page.evaluate(`Array.from(document.querySelectorAll('.evidence.code ui-code-block')).filter((block) => !block.highlighted).map((block) => block.language)`, &plain)
	if len(plain) != 0 {
		t.Errorf("the code blocks naming %v are left plain, though the page highlights those languages", plain)
	}
}
