package e2e_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const patchedWrite = "--- a/internal/cache/write.go\n+++ b/internal/cache/write.go\n@@ -10,5 +10,7 @@ func Write(key string, entry Entry) {\n func Write(key string, entry Entry) {\n-\tstore[key] = entry\n+\tmu.Lock()\n+\tdefer mu.Unlock()\n+\tstore[key] = entry\n }\n"

func withEvidence(blocks ...sketch) sketch {
	decision := aCacheDecision("single-choice")
	decision["evidence"] = blocks
	return aDocument("one thing to settle", nil, decision)
}

func TestADiffRendersFromAPatch(t *testing.T) {
	t.Parallel()
	_, page := read(t, written(t, withEvidence(sketch{"patch": patchedWrite})))

	var drawn struct {
		Found  bool   `json:"found"`
		Text   string `json:"text"`
		Blocks int    `json:"blocks"`
		Held   string `json:"held"`
		Tables int    `json:"tables"`
		Drawn  int    `json:"drawn"`
	}
	page.evaluate(`(() => {
  const diff = document.querySelector('.diff');
  const blocks = diff ? diff.querySelectorAll('ui-diff-block') : [];
  return {
    found: !!diff,
    text: diff ? diff.textContent : '',
    blocks: blocks.length,
    held: blocks[0] ? blocks[0].textContent : '',
    tables: blocks[0] ? blocks[0].querySelectorAll('table').length : -1,
    drawn: blocks[0] ? blocks[0].shadowRoot.querySelectorAll('table').length : 0,
  };
})()`, &drawn)

	if !drawn.Found {
		t.Fatal("a decision carrying a patch draws no diff")
	}
	for _, named := range []string{"internal/cache/write.go", "func Write(key string, entry Entry) {"} {
		if !strings.Contains(drawn.Text, named) {
			t.Errorf("the diff does not carry %q: %q", named, drawn.Text)
		}
	}
	if drawn.Blocks != 1 {
		t.Fatalf("the diff holds %d diff blocks, want one", drawn.Blocks)
	}
	if drawn.Held != patchedWrite {
		t.Errorf("the diff block holds %q, want the patch as the document stated it, %q", drawn.Held, patchedWrite)
	}
	if drawn.Tables != 0 {
		t.Error("the page writes the tables of the diff, and the diff block draws them itself")
	}
	if drawn.Drawn == 0 {
		t.Error("the diff block draws no table of its own")
	}
}

type patchedSide struct {
	Kind   string `json:"kind"`
	Number string `json:"number"`
	Text   string `json:"text"`
}

type patchedRow struct {
	Left  patchedSide `json:"left"`
	Right patchedSide `json:"right"`
}

func TestAPatchIsReadIntoBothSidesWithTheirNumbers(t *testing.T) {
	t.Parallel()
	_, page := read(t, written(t, withEvidence(sketch{"patch": patchedWrite})))

	var rows []patchedRow
	page.evaluate(`(async () => {
  const block = document.querySelector('.diff ui-diff-block');
  await block.updateComplete;
  return block.drawnHunks.flatMap((hunk) => hunk.rows);
})()`, &rows, awaitingThePromise)

	want := []patchedRow{
		{Left: patchedSide{"context", "10", "func Write(key string, entry Entry) {"}, Right: patchedSide{"context", "10", "func Write(key string, entry Entry) {"}},
		{Left: patchedSide{"removed", "11", "\tstore[key] = entry"}, Right: patchedSide{"added", "11", "\tmu.Lock()"}},
		{Left: patchedSide{"absent", "", ""}, Right: patchedSide{"added", "12", "\tdefer mu.Unlock()"}},
		{Left: patchedSide{"absent", "", ""}, Right: patchedSide{"added", "13", "\tstore[key] = entry"}},
		{Left: patchedSide{"context", "12", "}"}, Right: patchedSide{"context", "14", "}"}},
	}
	if diff := cmp.Diff(want, rows); diff != "" {
		t.Errorf("what the patch was read into differs (-want +got):\n%s", diff)
	}

	var heading string
	page.evaluate(`document.querySelector('.diff .diff-file').textContent.trim()`, &heading)
	if heading != "internal/cache/write.go" {
		t.Errorf("the diff names %q, and a person reads the path they know rather than the one the patch format prefixes", heading)
	}
}

func TestARenamedFileIsNamedByBothItsPathsInTheDiffHeading(t *testing.T) {
	t.Parallel()
	_, page := read(t, written(t, withEvidence(sketch{"patch": "--- a/internal/cache/old.go\n+++ b/internal/cache/new.go\n@@ -1,1 +1,1 @@\n-a\n+b\n"})))

	var named struct {
		From string `json:"from"`
		To   string `json:"to"`
		One  int    `json:"one"`
	}
	page.evaluate(`({
  from: document.querySelector('.diff .diff-from').textContent.trim(),
  to: document.querySelector('.diff .diff-to').textContent.trim(),
  one: document.querySelectorAll('.diff .diff-file').length,
})`, &named)

	if named.From != "internal/cache/old.go" || named.To != "internal/cache/new.go" || named.One != 0 {
		t.Errorf("the heading of a renamed file names %+v, want both of its paths and no single file", named)
	}
}

func TestAnEmptyPatchRendersNoDiff(t *testing.T) {
	t.Parallel()
	_, page := read(t, written(t, withEvidence(
		sketch{"language": "go", "code": "first()\n"},
		sketch{"caption": "nothing to show", "patch": "not a patch at all"},
	)))

	var drawn struct{ Diffs, Blocks int }
	page.evaluate(`({ Diffs: document.querySelectorAll('.diff').length, Blocks: document.querySelectorAll('.evidence:not(.illustrated)').length })`, &drawn)
	if drawn.Diffs != 0 {
		t.Errorf("a patch holding no hunk draws %d diffs", drawn.Diffs)
	}
	if drawn.Blocks != 1 {
		t.Errorf("a code block and a patch holding no hunk draw %d blocks, want only the code block", drawn.Blocks)
	}
}

func TestADiffStatesWhatItAddsAndRemoves(t *testing.T) {
	t.Parallel()
	_, page := read(t, written(t, withEvidence(sketch{"patch": "--- a/internal/cache/write.go\n+++ b/internal/cache/write.go\n@@ -10,3 +10,4 @@\n func Write(key string, entry Entry) {\n-\tstore[key] = entry\n+\tmu.Lock()\n+\tstore[key] = entry\n }\n"})))

	var said string
	page.evaluate(`document.querySelector('.evidence-aside').textContent.replace(/\s+/g, ' ').trim()`, &said)
	if said != "+2 −1" {
		t.Errorf("the diff says %q of itself, and it adds two lines and removes one", said)
	}
}

func TestADiffIsWrittenInTheLanguageStatedOrElseItsFilesOwn(t *testing.T) {
	t.Parallel()
	const patch = "--- a/internal/cache/write.go\n+++ b/internal/cache/write.go\n@@ -1,1 +1,1 @@\n-a\n+b\n"
	shown := showing(t, withEvidence(sketch{"patch": patch}))

	for _, stated := range []struct {
		language string
		patch    string
		written  string
	}{
		{"", patch, "go"},
		{"python", patch, "python"},
		{"", "--- a/Makefile\n+++ b/Makefile\n@@ -1,1 +1,1 @@\n-a\n+b\n", ""},
		{"", "--- a/web.v2/build.mjs\n+++ /dev/null\n@@ -1,1 +0,0 @@\n-a\n", "mjs"},
	} {
		block := sketch{"patch": stated.patch}
		if stated.language != "" {
			block["language"] = stated.language
		}
		shown.show(withEvidence(block))

		var languages []string
		shown.page.evaluate(`Array.from(document.querySelectorAll('ui-diff-block'), (held) => held.language)`, &languages)
		if len(languages) != 1 {
			t.Fatalf("a diff of %q is drawn as %d diff blocks, want one", stated.patch, len(languages))
		}
		if languages[0] != stated.written {
			t.Errorf("a diff of %q stating %q is written in %q, want %q", stated.patch, stated.language, languages[0], stated.written)
		}
	}
}

func TestEvidenceIsDrawnBlockAfterBlockInTheOrderWrittenEachUnderItsOwnHeading(t *testing.T) {
	t.Parallel()
	_, page := read(t, written(t, withEvidence(
		sketch{"caption": "as a prop", "language": "go", "code": "first()\n"},
		sketch{"language": "python", "code": "second()\n"},
		sketch{"caption": "the lock", "patch": "--- a/internal/cache/write.go\n+++ b/internal/cache/write.go\n@@ -1,1 +1,1 @@\n-a\n+b\n"},
		sketch{"caption": "nothing to show", "patch": "not a patch at all"},
	)))

	var headings []string
	page.evaluate(blockHeadings, &headings)

	want := []string{"as a prop go", "python", "the lock internal/cache/write.go"}
	if diff := cmp.Diff(want, headings); diff != "" {
		t.Errorf("the headings of the blocks differ from the order and captions written, and a patch with no hunk is left out (-want +got):\n%s", diff)
	}
}

type patchedHunk struct {
	Heading string       `json:"heading"`
	Rows    []patchedRow `json:"rows"`
}

func readHunksOf(page *reader) []patchedHunk {
	page.t.Helper()
	var hunks []patchedHunk
	page.evaluate(`(async () => {
  const block = document.querySelector('.diff ui-diff-block');
  await block.updateComplete;
  return block.drawnHunks;
})()`, &hunks, awaitingThePromise)
	return hunks
}

func TestALineOfADiffThatLooksLikeAFileHeaderInsideAHunkIsStillALine(t *testing.T) {
	// scenario: Removing the SQL comment -- note gives the hunk line --- note, and the diff keeps it as a removed line
	t.Parallel()
	_, page := read(t, written(t, withEvidence(sketch{"patch": "--- a/db/notes.sql\n+++ b/db/notes.sql\n@@ -1,3 +1,3 @@\n select 1;\n--- note\n+++ note\n select 2;\n"})))

	want := []patchedHunk{{Rows: []patchedRow{
		{Left: patchedSide{"context", "1", "select 1;"}, Right: patchedSide{"context", "1", "select 1;"}},
		{Left: patchedSide{"removed", "2", "-- note"}, Right: patchedSide{"added", "2", "++ note"}},
		{Left: patchedSide{"context", "3", "select 2;"}, Right: patchedSide{"context", "3", "select 2;"}},
	}}}
	if diff := cmp.Diff(want, readHunksOf(page)); diff != "" {
		t.Errorf("what the patch was read into differs (-want +got):\n%s", diff)
	}

	var named struct {
		File  string `json:"file"`
		Aside string `json:"aside"`
	}
	page.evaluate(`({
  file: document.querySelector('.diff .diff-file').textContent.trim(),
  aside: document.querySelector('.evidence-aside').textContent.replace(/\s+/g, ' ').trim(),
})`, &named)
	if named.File != "db/notes.sql" {
		t.Errorf("the diff names %q, and the lines inside the hunk are not file headers", named.File)
	}
	if named.Aside != "+1 −1" {
		t.Errorf("the diff says %q of itself, and it adds one line and removes one", named.Aside)
	}
}

func TestAHunkHeaderThatOmitsItsCountsStatesOneLineOnEachSide(t *testing.T) {
	// scenario: A count omitted in a hunk header means one, so the next file's headers still start the next file
	t.Parallel()
	_, page := read(t, written(t, withEvidence(sketch{"patch": "--- a/one.sql\n+++ b/one.sql\n@@ -5 +5 @@\n--- note\n+++ note\n--- a/two.sql\n+++ b/two.sql\n@@ -1 +1 @@\n-c\n+d\n"})))

	want := []patchedHunk{
		{Rows: []patchedRow{{Left: patchedSide{"removed", "5", "-- note"}, Right: patchedSide{"added", "5", "++ note"}}}},
		{Rows: []patchedRow{{Left: patchedSide{"removed", "1", "c"}, Right: patchedSide{"added", "1", "d"}}}},
	}
	if diff := cmp.Diff(want, readHunksOf(page)); diff != "" {
		t.Errorf("what the patch was read into differs (-want +got):\n%s", diff)
	}
}
