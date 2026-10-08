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
      "evidence": {
        "language": "csharp",
        "code": "public sealed class Settlement\n{\n    public string Reference { get; init; } = \"unset\";\n}\n"
      }
    },
    {
      "id": "c-2",
      "question": "is the change to the reader the one that was wanted",
      "mode": "single-choice",
      "options": [{"label": "yes", "meaning": "that is the change", "command": "agree c-2"}],
      "evidence": {
        "language": "diff",
        "code": "--- a/internal/cache/read.go\n+++ b/internal/cache/read.go\n@@ -1,2 +1,2 @@\n-\tdelete(locks, key)\n+\tdefer unlock(key)\n"
      }
    }
  ]
}`

const highlightedBlocks = `Object.fromEntries(Array.from(document.querySelectorAll('.evidence.code code')).map((held) => [
  Array.from(held.classList).find((named) => named.startsWith('language-')),
  held.classList.contains('hljs') && held.querySelectorAll('span[class^="hljs-"]').length > 0,
]))`

func TestTheLanguagesACallerActuallyPassesAreHighlighted(t *testing.T) {
	t.Parallel()
	held := standing(t)
	page := opening(t, held.serve(inTwoLanguagesDocument))

	page.waitFor(`document.querySelectorAll('.evidence.code code.hljs').length === 2`)

	var read map[string]bool
	page.evaluate(highlightedBlocks, &read)

	want := map[string]bool{"language-csharp": true, "language-diff": true}
	if diff := cmp.Diff(want, read); diff != "" {
		t.Errorf("what the page highlighted differs (-want +got):\n%s\na language the caller named and the page ignored renders as a grey block with nothing saying so", diff)
	}
}

const colouredDiffDocument = `{
  "title": "one thing to settle",
  "decisions": [
    {"id": "c-1", "question": "should the write take the lock", "mode": "single-choice",
     "options": [{"label": "take it", "meaning": "correctness first"}, {"label": "leave it", "meaning": "acceptable here"}],
     "evidence": {"patch": "--- a/internal/cache/write.go\n+++ b/internal/cache/write.go\n@@ -10,3 +10,4 @@\n func Write(key string, entry Entry) {\n-\t// store it, unlocked (fine for now)\n+\t// store it, under the lock\n+\tmu.Lock()\n \treturn\n }\n"}}
  ]
}`

func TestTheLinesOfADiffAreColouredLikeCodeAndStillMeetTheFloor(t *testing.T) {
	t.Parallel()
	_, page := read(t, colouredDiffDocument)

	page.waitFor(`!!document.querySelector('.hunk .line .hljs-keyword')`)

	var coloured struct {
		Removed bool `json:"removed"`
		Added   bool `json:"added"`
	}
	page.evaluate(`({
  removed: !!document.querySelector('.hunk .line[data-kind="removed"] .hljs-comment'),
  added: !!document.querySelector('.hunk .line[data-kind="added"] .hljs-comment'),
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
