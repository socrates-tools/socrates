package e2e_test

import (
	"bytes"
	"fmt"
	"image/png"
	"strings"
	"testing"

	"github.com/chromedp/cdproto/input"
	"github.com/google/go-cmp/cmp"
)

const (
	catalogShortLine         = "short-line"
	catalogWrappedLine       = "wrapped-line"
	catalogTabIndentedLine   = "tab-indented-line"
	catalogUnknownLanguage   = "unknown-language"
	catalogPatch             = "patch"
	catalogPatchInUnknownTag = "patch-in-an-unknown-language"
)

var sourceLinesOfEachCodeBlockState = map[string]int{
	catalogShortLine:       1,
	catalogWrappedLine:     3,
	catalogTabIndentedLine: 4,
	catalogUnknownLanguage: 2,
}

func codeBlockOf(key string) string {
	return stateOf("ui-code-block", key) + " ui-code-block"
}

func diffBlockOf(key string) string {
	return stateOf("ui-diff-block", key) + " ui-diff-block"
}

func numberedInTheCatalog(t *testing.T) *reader {
	t.Helper()

	page := catalogued(t)
	for key := range sourceLinesOfEachCodeBlockState {
		page.waitFor(`document.querySelector('` + codeBlockOf(key) + `').drawnLines.length > 0`)
	}
	return page
}

type drawnLine struct {
	Number string    `json:"number"`
	Gutter float64   `json:"gutter"`
	Rows   []float64 `json:"rows"`
}

func (r *reader) drawnLinesOf(selector string) []drawnLine {
	r.t.Helper()

	var lines []drawnLine
	r.evaluate(`document.querySelector(`+jsonOf(r.t, selector)+`).drawnLines.map((line) => ({
  number: line.number,
  gutter: line.gutter,
  rows: line.rows.map((columns) => Math.round(columns * 10) / 10),
}))`, &lines)
	return lines
}

func (r *reader) colouredPixelsOf(selector, token string) int {
	r.t.Helper()

	r.evaluate(`document.querySelector(`+jsonOf(r.t, selector)+`).scrollIntoView({ block: 'center' }); new Promise((settled) => requestAnimationFrame(() => requestAnimationFrame(settled)))`, nil, awaitingThePromise)
	var box struct{ Left, Top, Width, Height float64 }
	r.evaluate(`(() => { const box = document.querySelector(`+jsonOf(r.t, selector)+`).getBoundingClientRect(); return { Left: box.left, Top: box.top, Width: box.width, Height: box.height }; })()`, &box)
	var want [3]int
	if _, err := fmt.Sscanf(r.pageColourOf(token), "rgb(%d, %d, %d)", &want[0], &want[1], &want[2]); err != nil {
		r.t.Fatalf("reading the colour of %s: %v", token, err)
	}

	shot, err := png.Decode(bytes.NewReader(r.screenshot()))
	if err != nil {
		r.t.Fatalf("reading what the page painted: %v", err)
	}
	matching := 0
	for x := int(box.Left); x < int(box.Left+box.Width); x++ {
		for y := int(box.Top); y < int(box.Top+box.Height); y++ {
			red, green, blue, _ := shot.At(x, y).RGBA()
			if abs(int(red>>8)-want[0])+abs(int(green>>8)-want[1])+abs(int(blue>>8)-want[2]) <= 12 {
				matching++
			}
		}
	}
	return matching
}

const (
	insideTheFirstLine = 18
	besideTheGutter    = 49
)

func (r *reader) selectByDraggingAcross(selector string) string {
	r.t.Helper()

	r.evaluate(`document.querySelector(`+jsonOf(r.t, selector)+`).scrollIntoView({ block: 'center' }); new Promise((settled) => requestAnimationFrame(() => requestAnimationFrame(settled)))`, nil, awaitingThePromise)
	var box struct{ Left, Top, Right, Bottom float64 }
	r.evaluate(`(() => { const box = document.querySelector(`+jsonOf(r.t, selector)+`).getBoundingClientRect(); return { Left: box.left, Top: box.top, Right: box.right, Bottom: box.bottom }; })()`, &box)
	r.do(
		sending(input.DispatchMouseEvent, input.DispatchMouseEventParams{Type: input.DispatchMouseEventTypeMouseMoved, X: box.Left + besideTheGutter, Y: box.Top + insideTheFirstLine}),
		sending(input.DispatchMouseEvent, input.DispatchMouseEventParams{Type: input.DispatchMouseEventTypeMousePressed, X: box.Left + besideTheGutter, Y: box.Top + insideTheFirstLine, Button: input.MouseButtonLeft, Buttons: 1, ClickCount: 1}),
		sending(input.DispatchMouseEvent, input.DispatchMouseEventParams{Type: input.DispatchMouseEventTypeMouseMoved, X: box.Right - 1, Y: box.Bottom - insideTheFirstLine, Button: input.MouseButtonLeft, Buttons: 1}),
		sending(input.DispatchMouseEvent, input.DispatchMouseEventParams{Type: input.DispatchMouseEventTypeMouseReleased, X: box.Right - 1, Y: box.Bottom - insideTheFirstLine, Button: input.MouseButtonLeft, Buttons: 1, ClickCount: 1}),
	)
	var selected string
	r.evaluate(`window.getSelection().toString()`, &selected)
	return selected
}

const everyBlockOfEvidenceAdoptsTheSharedStylesheet = `(() => {
  const blocks = Array.from(document.querySelectorAll('ui-code-block, ui-diff-block'));
  return blocks.length > 0 && blocks.every((block) => block.shadowRoot !== null && block.shadowRoot.adoptedStyleSheets.length > 0);
})()`

func TestAShadowRootHoldsEveryCodeBlockAndDiffBlockOfTheCatalogAndThePage(t *testing.T) {
	// scenario: Every code block and diff block renders into a shadow root that adopts the shared stylesheet, on the catalog and on the page, and the policy blocks nothing
	t.Parallel()
	catalog := numberedInTheCatalog(t)
	_, page := read(t, evidenceDocument(aBlockOfEachKind))

	for place, reader := range map[string]*reader{"catalog": catalog, "page": page} {
		reader.waitFor(everyBlockOfEvidenceAdoptsTheSharedStylesheet)
		reader.evaluate(afterTwoFrames, nil, awaitingThePromise)
		if violations := reader.violations(); len(violations) != 0 {
			t.Errorf("the %s's policy blocked %v while the blocks of evidence adopted the shared stylesheet", place, violations)
		}
	}
}

func TestTheCatalogNumbersEverySourceLineOfACodeBlockOnceFromOne(t *testing.T) {
	// scenario: Every code block state in the catalog numbers its source lines from 1, a wrapped line once
	t.Parallel()
	page := numberedInTheCatalog(t)

	for key, sources := range sourceLinesOfEachCodeBlockState {
		t.Run(key, func(t *testing.T) {
			var numbers []string
			for _, line := range page.drawnLinesOf(codeBlockOf(key)) {
				numbers = append(numbers, line.Number)
			}

			var want []string
			for number := 1; number <= sources; number++ {
				want = append(want, fmt.Sprint(number))
			}
			if diff := cmp.Diff(want, numbers); diff != "" {
				t.Errorf("the numbers of the %s block differ (-want +got):\n%s", key, diff)
			}
		})
	}
}

func TestTheCatalogDrawsEachCodeBlockNumberInAGutterTheCodeStartsAfter(t *testing.T) {
	// scenario: The numbers of a code block are drawn in a gutter the code of its line starts after
	t.Parallel()
	page := numberedInTheCatalog(t)

	for key := range sourceLinesOfEachCodeBlockState {
		t.Run(key, func(t *testing.T) {
			for position, line := range page.drawnLinesOf(codeBlockOf(key)) {
				if line.Gutter <= 0 {
					t.Errorf("source line %d draws no gutter", position+1)
				}
				if len(line.Rows) == 0 || line.Rows[0] < 0 {
					t.Errorf("the code of source line %d starts inside the gutter: %v", position+1, line.Rows)
				}
			}
		})
	}
}

func TestTheCatalogInksEachCodeBlockNumberInTheMutedForeground(t *testing.T) {
	// scenario: The numbers of a code block are inked in the muted foreground, so a number is read as the margin and not as the code
	t.Parallel()
	page := numberedInTheCatalog(t)

	if inked := page.colouredPixelsOf(codeBlockOf(catalogUnknownLanguage), "--fg-3"); inked == 0 {
		t.Error("the block draws no pixel in the muted foreground, so its numbers are not drawn in the gutter's ink")
	}
}

func TestTheCatalogKeepsATabIndentedCodeBlocksTextExactlyAsStated(t *testing.T) {
	// scenario: A code block that numbers and highlights its lines still holds the code exactly as it was written
	t.Parallel()
	page := numberedInTheCatalog(t)

	var held struct {
		Lines  int    `json:"lines"`
		Stated string `json:"stated"`
	}
	page.evaluate(`({
  lines: document.querySelector('`+codeBlockOf(catalogTabIndentedLine)+`').drawnLines.length,
  stated: document.querySelector('`+codeBlockOf(catalogTabIndentedLine)+`').statedText(),
})`, &held)

	const written = "func Write() {\n\tmu.Lock()\n\treturn\n}\n"
	if held.Stated != written {
		t.Errorf("the block states %q as its code, want %q", held.Stated, written)
	}
	if held.Lines != 4 {
		t.Errorf("the block draws %d lines, want the 4 it was written with", held.Lines)
	}
}

func TestSelectingACodeBlocksTextByHandSelectsNoLineNumber(t *testing.T) {
	// scenario: Dragging a selection across a numbered code block selects its code and none of its line numbers
	t.Parallel()
	page := numberedInTheCatalog(t)

	selected := page.selectByDraggingAcross(codeBlockOf(catalogTabIndentedLine))

	const written = "func Write() {\n\tmu.Lock()\n\treturn\n}"
	if got := strings.TrimRight(strings.ReplaceAll(selected, "\r", ""), "\n"); got != written {
		t.Errorf("selecting the block by hand selects %q, want its code %q and no number from its gutter", selected, written)
	}
}

func TestTheCatalogIndentsATabIndentedLineByFourColumns(t *testing.T) {
	// scenario: A line starting with one tab is drawn four columns in
	t.Parallel()
	page := numberedInTheCatalog(t)

	var starts [][]float64
	for _, line := range page.drawnLinesOf(codeBlockOf(catalogTabIndentedLine)) {
		starts = append(starts, line.Rows)
	}

	if len(starts) != 4 || len(starts[1]) != 1 || len(starts[2]) != 1 {
		t.Fatalf("the tab-indented block is drawn as rows %v, want four lines of one row each", starts)
	}
	if starts[0][0] != 0 || starts[1][0] != 4 || starts[2][0] != 4 {
		t.Errorf("the lines are indented by %v, %v and %v columns, want 0, 4 and 4", starts[0][0], starts[1][0], starts[2][0])
	}
}

func TestTheCatalogWrapsALineWiderThanItsBlockAndNeverScrollsItSideways(t *testing.T) {
	// scenario: A line wider than its code block is drawn on more than one row and the block does not scroll
	t.Parallel()
	page := numberedInTheCatalog(t)

	lines := page.drawnLinesOf(codeBlockOf(catalogWrappedLine))
	var scrolling []string
	page.evaluate(`Array.from(document.querySelectorAll('ui-code-block')).filter((held) => held.scrollWidth > held.clientWidth + 1).map((held) => held.closest('[data-state]').dataset.state)`, &scrolling)

	if len(lines) != 3 {
		t.Fatalf("the wrapped block is drawn as %d lines, want 3", len(lines))
	}
	for position, line := range lines {
		if len(line.Rows) < 2 {
			t.Errorf("source line %d is drawn on %d row, so a line wider than its block is not wrapped", position+1, len(line.Rows))
		}
	}
	if len(scrolling) != 0 {
		t.Errorf("the code blocks %v scroll sideways", scrolling)
	}
}

func TestTheCatalogHangsAWrappedLineOneTabStopFurtherInThanThatLinesOwnIndentation(t *testing.T) {
	// scenario: A continuation starts at the next tab stop beyond the indentation of the line it continues
	t.Parallel()
	page := numberedInTheCatalog(t)

	lines := page.drawnLinesOf(codeBlockOf(catalogWrappedLine))

	startsAndContinuations := [][2]float64{{0, 4}, {4, 8}, {6, 8}}
	if len(lines) != len(startsAndContinuations) {
		t.Fatalf("the wrapped block is drawn as %d lines, want %d", len(lines), len(startsAndContinuations))
	}
	for position, want := range startsAndContinuations {
		rows := lines[position].Rows
		if len(rows) < 2 {
			t.Errorf("source line %d is not wrapped", position+1)
			continue
		}
		if rows[0] != want[0] || rows[1] != want[1] {
			t.Errorf("source line %d starts at column %v and continues at column %v, want %v and %v", position+1, rows[0], rows[1], want[0], want[1])
		}
	}
}

func TestTheCatalogPaintsACodeBlockInALanguageThePageKnowsByItsTokens(t *testing.T) {
	// scenario: A code block naming a language the page highlights is painted by its tokens
	t.Parallel()
	page := numberedInTheCatalog(t)

	var highlighted bool
	page.evaluate(`document.querySelector('`+codeBlockOf(catalogShortLine)+`').highlighted`, &highlighted)

	if !highlighted {
		t.Error("the short Go line is not highlighted")
	}
	if inked := page.colouredPixelsOf(codeBlockOf(catalogShortLine), "--token-keyword"); inked == 0 {
		t.Error("the keyword of the short Go line is painted in no pixel of the keyword token")
	}
}

func TestTheCatalogDrawsACodeBlockInAnUnknownLanguageAsPlainTextThatIsStillNumbered(t *testing.T) {
	// scenario: A code block naming a language the page does not highlight is plain text, numbered like any other
	t.Parallel()
	page := numberedInTheCatalog(t)

	var held struct {
		Highlighted bool   `json:"highlighted"`
		Lines       int    `json:"lines"`
		Text        string `json:"text"`
	}
	page.evaluate(`(() => {
  const block = document.querySelector('`+codeBlockOf(catalogUnknownLanguage)+`');
  return { highlighted: block.highlighted, lines: block.drawnLines.length, text: block.statedText() };
})()`, &held)

	if held.Highlighted {
		t.Error("the block in an unknown language is highlighted")
	}
	for _, token := range []string{"--token-keyword", "--token-string", "--token-number", "--token-call", "--token-type"} {
		if inked := page.colouredPixelsOf(codeBlockOf(catalogUnknownLanguage), token); inked != 0 {
			t.Errorf("the block in an unknown language is painted in %d pixels of %s, want plain text", inked, token)
		}
	}
	if held.Lines != 2 || held.Text != "qapla'\nbatlh\n" {
		t.Errorf("the block in an unknown language is %+v, want its two lines numbered and its text as written", held)
	}
}

func TestTheCatalogLeavesACodeBlockAsItWasWhenItIsMovedWithinThePage(t *testing.T) {
	// scenario: A code block attached to the page a second time is not numbered or highlighted again
	t.Parallel()
	page := numberedInTheCatalog(t)

	var held struct {
		Before []string `json:"before"`
		After  []string `json:"after"`
		Text   string   `json:"text"`
	}
	page.evaluate(`(() => {
  const block = document.querySelector('`+codeBlockOf(catalogTabIndentedLine)+`');
  const numbers = () => block.drawnLines.map((line) => line.number);
  const before = numbers();
  block.parentNode.appendChild(block);
  return { before, after: numbers(), text: block.statedText() };
})()`, &held)

	if diff := cmp.Diff(held.Before, held.After); diff != "" {
		t.Errorf("moving the block changed its lines (-before +after):\n%s", diff)
	}
	if held.Text != "func Write() {\n\tmu.Lock()\n\treturn\n}\n" {
		t.Errorf("moving the block changed the code it states to %q", held.Text)
	}
}

type drawnSide struct {
	Kind   string `json:"kind"`
	Number string `json:"number"`
	Text   string `json:"text"`
}

type drawnToken struct {
	Text string `json:"text"`
	Ink  string `json:"ink"`
}

func (r *reader) tokensOfTheSidesOf(selector, kind string) []drawnToken {
	r.t.Helper()

	var tokens []drawnToken
	r.evaluate(`document.querySelector(`+jsonOf(r.t, selector)+`).drawnHunks.flatMap((hunk) => hunk.rows).flatMap((row) => [row.left, row.right]).filter((side) => side.kind === `+jsonOf(r.t, kind)+`).flatMap((side) => side.tokens)`, &tokens)
	return tokens
}

type drawnRow struct {
	Left  drawnSide `json:"left"`
	Right drawnSide `json:"right"`
}

type drawnHunk struct {
	Heading string     `json:"heading"`
	Rows    []drawnRow `json:"rows"`
}

func TestTheCatalogDrawsAPatchAsAFixedTableOfNumberedSides(t *testing.T) {
	// scenario: A patch in the catalog is a table of rows, each side numbered and kinded as the patch states
	t.Parallel()
	page := catalogued(t)
	page.waitFor(`document.querySelector('` + diffBlockOf(catalogPatch) + `').drawnHunks.length > 0`)

	var hunks []drawnHunk
	page.evaluate(`document.querySelector('`+diffBlockOf(catalogPatch)+`').drawnHunks.map((hunk) => ({ heading: hunk.heading, rows: hunk.rows }))`, &hunks)

	want := []drawnHunk{{Heading: "func Limit", Rows: []drawnRow{
		{drawnSide{"context", "1", "func Limit() int {"}, drawnSide{"context", "1", "func Limit() int {"}},
		{drawnSide{"removed", "2", "\treturn 2"}, drawnSide{"added", "2", "\t// the limit, raised"}},
		{drawnSide{"absent", "", ""}, drawnSide{"added", "3", "\treturn 3"}},
		{drawnSide{"context", "3", "}"}, drawnSide{"context", "4", "}"}},
	}}}
	if diff := cmp.Diff(want, hunks); diff != "" {
		t.Errorf("the hunks of the patch differ (-want +got):\n%s", diff)
	}
}

func (r *reader) hunksReadFrom(patch string) []drawnHunk {
	r.t.Helper()

	var hunks []drawnHunk
	r.evaluate(`(async (patch) => {
  const block = document.createElement('ui-diff-block');
  block.textContent = patch;
  document.body.append(block);
  await block.updateComplete;
  const hunks = block.drawnHunks.map((hunk) => ({ heading: hunk.heading, rows: hunk.rows }));
  block.remove();
  return hunks;
})(`+jsonOf(r.t, patch)+`)`, &hunks, awaitingThePromise)
	return hunks
}

func TestADiffBlockReadsAPatchIntoHunksOfPairedNumberedSides(t *testing.T) {
	// scenario: A diff block given a patch pairs its removals with its additions, numbers each side from the hunk's own start and leaves what is not a line of the patch out
	t.Parallel()
	page := catalogued(t)

	side := func(kind, number, text string) drawnSide { return drawnSide{kind, number, text} }
	absent := side("absent", "", "")
	for _, each := range []struct {
		scenario string
		patch    string
		want     []drawnHunk
	}{
		{
			"removals paired with more additions leave the left side absent",
			"--- a/internal/cache/write.go\n+++ b/internal/cache/write.go\n@@ -10,5 +10,7 @@ func Write(key string, entry Entry) {\n func Write(key string, entry Entry) {\n-\tstore[key] = entry\n+\tmu.Lock()\n+\tdefer mu.Unlock()\n+\tstore[key] = entry\n }\n",
			[]drawnHunk{{"func Write(key string, entry Entry) {", []drawnRow{
				{side("context", "10", "func Write(key string, entry Entry) {"), side("context", "10", "func Write(key string, entry Entry) {")},
				{side("removed", "11", "\tstore[key] = entry"), side("added", "11", "\tmu.Lock()")},
				{absent, side("added", "12", "\tdefer mu.Unlock()")},
				{absent, side("added", "13", "\tstore[key] = entry")},
				{side("context", "12", "}"), side("context", "14", "}")},
			}}},
		},
		{
			"a removal alone leaves the right side absent",
			"--- a/x.go\n+++ b/x.go\n@@ -4,2 +4,1 @@\n-gone\n kept\n",
			[]drawnHunk{{"", []drawnRow{
				{side("removed", "4", "gone"), absent},
				{side("context", "5", "kept"), side("context", "4", "kept")},
			}}},
		},
		{
			"two hunks are two tables and a line saying there is no newline is not a row",
			"--- a/x.go\n+++ b/x.go\n@@ -1,1 +1,1 @@\n-a\n+b\n\\ No newline at end of file\n@@ -20,1 +20,1 @@ later\n-c\n+d\n",
			[]drawnHunk{
				{"", []drawnRow{{side("removed", "1", "a"), side("added", "1", "b")}}},
				{"later", []drawnRow{{side("removed", "20", "c"), side("added", "20", "d")}}},
			},
		},
		{
			"a hunk heading naming no range starts both sides at 1",
			"@@ nothing to read\n-a\n+b\n",
			[]drawnHunk{{"", []drawnRow{{side("removed", "1", "a"), side("added", "1", "b")}}}},
		},
		{
			"text before any hunk is not drawn",
			"just some words\nand more\n",
			[]drawnHunk{},
		},
	} {
		t.Run(each.scenario, func(t *testing.T) {
			got := page.hunksReadFrom(each.patch)
			if got == nil {
				got = []drawnHunk{}
			}
			if diff := cmp.Diff(each.want, got); diff != "" {
				t.Errorf("the hunks read from the patch differ (-want +got):\n%s", diff)
			}
		})
	}
}

func TestTheCatalogPaintsAPatchInALanguageThePageKnowsAndQuietsWhatChanged(t *testing.T) {
	// scenario: A patch naming a language the page highlights paints each side, and quiets a comment on a changed line
	t.Parallel()
	page := catalogued(t)
	page.waitFor(`document.querySelector('` + diffBlockOf(catalogPatch) + `').highlighted`)

	block := diffBlockOf(catalogPatch)
	if diff := cmp.Diff([]drawnToken{{"func", "keyword"}}, firstOf(page.tokensOfTheSidesOf(block, "context"), "keyword")); diff != "" {
		t.Errorf("the keyword of a context line differs (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]drawnToken{{"return", "keyword"}}, firstOf(page.tokensOfTheSidesOf(block, "added"), "keyword")); diff != "" {
		t.Errorf("the keyword of an added line differs (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]drawnToken{{"// the limit, raised", "quiet-on-change"}}, firstOf(page.tokensOfTheSidesOf(block, "added"), "quiet-on-change")); diff != "" {
		t.Errorf("the comment on the added line differs (-want +got):\n%s", diff)
	}
	for _, token := range page.tokensOfTheSidesOf(block, "added") {
		if token.Ink == "comment" {
			t.Errorf("the comment %q on an added line is drawn in the comment ink, want the quiet ink of a change", token.Text)
		}
	}
	if inked := page.colouredPixelsOf(block, "--token-keyword"); inked == 0 {
		t.Error("the keywords of the patch are painted in no pixel of the keyword token")
	}
}

func firstOf(tokens []drawnToken, ink string) []drawnToken {
	for _, token := range tokens {
		if token.Ink == ink {
			return []drawnToken{token}
		}
	}
	return nil
}

func TestTheCatalogDrawsAPatchInAnUnknownLanguageAsPlainTextRows(t *testing.T) {
	// scenario: A patch naming a language the page does not highlight keeps its rows and paints none of them
	t.Parallel()
	page := catalogued(t)
	page.waitFor(`document.querySelector('` + diffBlockOf(catalogPatchInUnknownTag) + `').drawnHunks.length > 0`)

	var held struct {
		Highlighted bool `json:"highlighted"`
		Rows        int  `json:"rows"`
		Tokens      int  `json:"tokens"`
	}
	page.evaluate(`(() => {
  const block = document.querySelector('`+diffBlockOf(catalogPatchInUnknownTag)+`');
  const sides = block.drawnHunks.flatMap((hunk) => hunk.rows).flatMap((row) => [row.left, row.right]);
  return { highlighted: block.highlighted, rows: block.drawnHunks.flatMap((hunk) => hunk.rows).length, tokens: sides.flatMap((side) => side.tokens).length };
})()`, &held)

	if held.Highlighted || held.Rows == 0 || held.Tokens != 0 {
		t.Errorf("the patch in an unknown language is %+v, want its rows and no highlighting", held)
	}
	if inked := page.colouredPixelsOf(diffBlockOf(catalogPatchInUnknownTag), "--token-keyword"); inked != 0 {
		t.Errorf("the patch in an unknown language is painted in %d pixels of the keyword token, want plain text", inked)
	}
}

func TestTheCatalogGivesAPatchTooWideForItsBlockAScrollOfItsOwn(t *testing.T) {
	// scenario: At a window narrower than the patch's table the table scrolls inside its block and not the page
	t.Parallel()
	page := catalogued(t)
	page.waitFor(`document.querySelector('` + diffBlockOf(catalogPatch) + `').drawnHunks.length > 0`)
	page.narrow(360, 780)

	var held struct {
		Overflow     string  `json:"overflow"`
		Scrolls      bool    `json:"scrolls"`
		BlockWithin  bool    `json:"blockWithin"`
		BlockScrolls bool    `json:"blockScrolls"`
		PageSideways float64 `json:"pageSideways"`
	}
	page.evaluate(`(() => {
  const block = document.querySelector('`+diffBlockOf(catalogPatch)+`');
  const hunk = block.drawnHunks[0];
  const state = block.closest('[data-state]');
  const page = document.documentElement;
  return {
    overflow: hunk.overflowX,
    scrolls: hunk.contentWidth > hunk.width + 1,
    blockWithin: block.getBoundingClientRect().right <= state.getBoundingClientRect().right + 1,
    blockScrolls: block.scrollWidth > block.clientWidth + 1,
    pageSideways: page.scrollWidth - page.clientWidth,
  };
})()`, &held)

	if held.Overflow != "auto" && held.Overflow != "scroll" {
		t.Errorf("the table overflows as %q, want a scroll of its own", held.Overflow)
	}
	if !held.Scrolls {
		t.Error("the patch does not scroll inside its table at 360px")
	}
	if !held.BlockWithin || held.BlockScrolls || held.PageSideways > 1 {
		t.Errorf("the patch spills out of its block %v, scrolls its block %v or scrolls the page by %vpx", !held.BlockWithin, held.BlockScrolls, held.PageSideways)
	}
}

func TestTheCatalogRaisesNoErrorWhenAPatchInAnUnknownLanguageIsAttached(t *testing.T) {
	// scenario: Attaching a patch that names a language the page does not highlight raises nothing
	t.Parallel()
	page := catalogued(t)

	var raised []string
	page.evaluate(`(() => {
  const raised = [];
  window.addEventListener('error', (event) => raised.push(event.message));
  document.body.appendChild(document.querySelector('`+diffBlockOf(catalogPatchInUnknownTag)+`').cloneNode(true));
  return raised;
})()`, &raised)

	if len(raised) != 0 {
		t.Errorf("attaching the patch raised %v", raised)
	}
}
