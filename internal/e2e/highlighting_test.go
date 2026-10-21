package e2e_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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

	page.waitFor(`!!document.querySelector('ui-diff-block td:nth-child(even) .hljs-keyword')`)

	var coloured struct {
		Removed bool `json:"removed"`
		Added   bool `json:"added"`
	}
	page.evaluate(`({
  removed: !!document.querySelector('ui-diff-block td:nth-child(even)[data-kind="removed"] .hljs-comment'),
  added: !!document.querySelector('ui-diff-block td:nth-child(even)[data-kind="added"] .hljs-comment'),
})`, &coloured)
	if !coloured.Removed || !coloured.Added {
		t.Errorf("a comment on a changed line of a Go diff is not coloured as one: %+v", coloured)
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

	page.waitFor(`!!document.querySelector('ui-diff-block td:nth-child(even) .hljs-keyword')`)

	var coloured bool
	page.evaluate(`!!document.querySelector('ui-diff-block td:nth-child(even)[data-kind="added"] .hljs-built_in, ui-diff-block td:nth-child(even)[data-kind="added"] .hljs-type')`, &coloured)
	if !coloured {
		t.Error("a type on a changed line of a patch whose language is ts is not coloured as one")
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
