package e2e_test

import (
	"fmt"
	"testing"

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
		page.waitFor(`document.querySelectorAll('` + codeBlockOf(key) + ` .code-line').length > 0`)
	}
	return page
}

const columnOfEachRowOfEachLine = `(block) => {
  const code = document.querySelector(block + ' [data-part="code"]');
  const probe = document.createElement('span');
  probe.style.cssText = 'display:inline-block;width:1ch;position:absolute;visibility:hidden';
  code.appendChild(probe);
  const column = probe.getBoundingClientRect().width;
  probe.remove();
  return Array.from(code.querySelectorAll('.code-line')).map((line) => {
    const left = line.getBoundingClientRect().left + parseFloat(getComputedStyle(line, '::before').width);
    const rows = new Map();
    const walker = document.createTreeWalker(line, NodeFilter.SHOW_TEXT);
    for (let node = walker.nextNode(); node; node = walker.nextNode()) {
      for (let at = 0; at < node.data.length; at++) {
        if (/\s/.test(node.data[at])) {
          continue;
        }
        const range = document.createRange();
        range.setStart(node, at);
        range.setEnd(node, at + 1);
        const box = range.getBoundingClientRect();
        const top = Math.round(box.top);
        if (!rows.has(top)) {
          rows.set(top, (box.left - left) / column);
        }
      }
    }
    return Array.from(rows.values()).map((columns) => Math.round(columns * 10) / 10);
  });
}`

func TestTheCatalogNumbersEverySourceLineOfACodeBlockOnceFromOne(t *testing.T) {
	// scenario: Every code block state in the catalog numbers its source lines from 1, a wrapped line once
	t.Parallel()
	page := numberedInTheCatalog(t)

	for key, sources := range sourceLinesOfEachCodeBlockState {
		t.Run(key, func(t *testing.T) {
			var numbers []string
			page.evaluate(`Array.from(document.querySelectorAll('`+codeBlockOf(key)+` .code-line')).map((line) => line.dataset.number)`, &numbers)

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
			var drawn struct {
				Gutters   []float64 `json:"gutters"`
				TextAfter []bool    `json:"textAfter"`
			}
			page.evaluate(`(() => {
  const code = document.querySelector('`+codeBlockOf(key)+` [data-part="code"]');
  const lines = Array.from(code.querySelectorAll('.code-line'));
  return {
    gutters: lines.map((line) => parseFloat(getComputedStyle(line, '::before').width)),
    textAfter: lines.map((line) => {
      const range = document.createRange();
      range.selectNodeContents(line);
      const gutter = parseFloat(getComputedStyle(line, '::before').width);
      return range.getClientRects()[0].left >= line.getBoundingClientRect().left + gutter - 0.5;
    }),
  };
})()`, &drawn)

			for position, gutter := range drawn.Gutters {
				if gutter <= 0 {
					t.Errorf("source line %d draws no gutter", position+1)
				}
				if !drawn.TextAfter[position] {
					t.Errorf("the code of source line %d starts inside the gutter", position+1)
				}
			}
		})
	}
}

func TestTheCatalogKeepsATabIndentedCodeBlocksTextExactlyAsStated(t *testing.T) {
	// scenario: A code block that numbers and highlights its lines still holds the code exactly as it was written
	t.Parallel()
	page := numberedInTheCatalog(t)

	var held struct {
		Drawn  string `json:"drawn"`
		Stated string `json:"stated"`
	}
	page.evaluate(`({
  drawn: document.querySelector('`+codeBlockOf(catalogTabIndentedLine)+` [data-part="code"]').textContent,
  stated: document.querySelector('`+codeBlockOf(catalogTabIndentedLine)+`').statedText(),
})`, &held)

	const written = "func Write() {\n\tmu.Lock()\n\treturn\n}\n"
	if held.Drawn != written {
		t.Errorf("the block draws %q, want what the catalog wrote, %q", held.Drawn, written)
	}
	if held.Stated != written {
		t.Errorf("the block states %q as its code, want %q", held.Stated, written)
	}
}

func TestTheCatalogIndentsATabIndentedLineByFourColumns(t *testing.T) {
	// scenario: A line starting with one tab is drawn four columns in
	t.Parallel()
	page := numberedInTheCatalog(t)

	var rows [][]float64
	page.evaluate(`(`+columnOfEachRowOfEachLine+`)('`+codeBlockOf(catalogTabIndentedLine)+`')`, &rows)

	if len(rows) != 4 || len(rows[1]) != 1 || len(rows[2]) != 1 {
		t.Fatalf("the tab-indented block is drawn as rows %v, want four lines of one row each", rows)
	}
	if rows[0][0] != 0 || rows[1][0] != 4 || rows[2][0] != 4 {
		t.Errorf("the lines are indented by %v, %v and %v columns, want 0, 4 and 4", rows[0][0], rows[1][0], rows[2][0])
	}
}

func TestTheCatalogWrapsALineWiderThanItsBlockAndNeverScrollsItSideways(t *testing.T) {
	// scenario: A line wider than its code block is drawn on more than one row and the block does not scroll
	t.Parallel()
	page := numberedInTheCatalog(t)

	var rows [][]float64
	page.evaluate(`(`+columnOfEachRowOfEachLine+`)('`+codeBlockOf(catalogWrappedLine)+`')`, &rows)
	var scrolling []string
	page.evaluate(`Array.from(document.querySelectorAll('[data-primitive="ui-code-block"] pre')).filter((held) => held.scrollWidth > held.clientWidth + 1).map((held) => held.closest('[data-state]').dataset.state)`, &scrolling)

	if len(rows) != 3 {
		t.Fatalf("the wrapped block is drawn as %d lines, want 3", len(rows))
	}
	for position, columns := range rows {
		if len(columns) < 2 {
			t.Errorf("source line %d is drawn on %d row, so a line wider than its block is not wrapped", position+1, len(columns))
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

	var rows [][]float64
	page.evaluate(`(`+columnOfEachRowOfEachLine+`)('`+codeBlockOf(catalogWrappedLine)+`')`, &rows)

	startsAndContinuations := [][2]float64{{0, 4}, {4, 8}, {6, 8}}
	if len(rows) != len(startsAndContinuations) {
		t.Fatalf("the wrapped block is drawn as %d lines, want %d", len(rows), len(startsAndContinuations))
	}
	for position, want := range startsAndContinuations {
		if len(rows[position]) < 2 {
			t.Errorf("source line %d is not wrapped", position+1)
			continue
		}
		if rows[position][0] != want[0] || rows[position][1] != want[1] {
			t.Errorf("source line %d starts at column %v and continues at column %v, want %v and %v", position+1, rows[position][0], rows[position][1], want[0], want[1])
		}
	}
}

func TestTheCatalogHighlightsACodeBlockInALanguageThePageKnows(t *testing.T) {
	// scenario: A code block naming a language the page highlights is painted by its tokens
	t.Parallel()
	page := numberedInTheCatalog(t)

	var held struct {
		Highlighted bool   `json:"highlighted"`
		Keywords    int    `json:"keywords"`
		Painted     string `json:"painted"`
		Token       string `json:"token"`
		Plain       string `json:"plain"`
	}
	page.evaluate(`(() => {
  const code = document.querySelector('`+codeBlockOf(catalogShortLine)+` [data-part="code"]');
  const probe = document.createElement('span');
  probe.style.color = 'var(--token-keyword)';
  code.appendChild(probe);
  const token = getComputedStyle(probe).color;
  probe.remove();
  return {
    highlighted: code.classList.contains('hljs'),
    keywords: code.querySelectorAll('.hljs-keyword').length,
    painted: getComputedStyle(code.querySelector('.hljs-keyword')).color,
    token,
    plain: getComputedStyle(code).color,
  };
})()`, &held)

	if !held.Highlighted || held.Keywords != 1 {
		t.Errorf("the short Go line is %+v, want highlighted with its one keyword painted", held)
	}
	if held.Painted != held.Token || held.Painted == held.Plain {
		t.Errorf("the keyword is painted %s, want the keyword token %s and not the plain %s", held.Painted, held.Token, held.Plain)
	}
}

func TestTheCatalogDrawsACodeBlockInAnUnknownLanguageAsPlainTextThatIsStillNumbered(t *testing.T) {
	// scenario: A code block naming a language the page does not highlight is plain text, numbered like any other
	t.Parallel()
	page := numberedInTheCatalog(t)

	var held struct {
		Highlighted bool   `json:"highlighted"`
		Tokens      int    `json:"tokens"`
		Lines       int    `json:"lines"`
		Text        string `json:"text"`
	}
	page.evaluate(`(() => {
  const code = document.querySelector('`+codeBlockOf(catalogUnknownLanguage)+` [data-part="code"]');
  return {
    highlighted: code.classList.contains('hljs'),
    tokens: code.querySelectorAll('span[class^="hljs-"]').length,
    lines: code.querySelectorAll('.code-line').length,
    text: code.textContent,
  };
})()`, &held)

	if held.Highlighted || held.Tokens != 0 {
		t.Errorf("the block in an unknown language is %+v, want no highlighting", held)
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
  const numbers = () => Array.from(block.querySelectorAll('.code-line')).map((line) => line.dataset.number);
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

func TestTheCatalogDrawsAPatchAsAFixedTableOfNumberedSides(t *testing.T) {
	// scenario: A patch in the catalog is a table of rows, each side numbered and kinded as the patch states
	t.Parallel()
	page := catalogued(t)

	var rows []struct {
		LeftKind    string `json:"leftKind"`
		LeftNumber  string `json:"leftNumber"`
		RightKind   string `json:"rightKind"`
		RightNumber string `json:"rightNumber"`
		LeftText    string `json:"leftText"`
		RightText   string `json:"rightText"`
	}
	page.evaluate(`Array.from(document.querySelectorAll('`+diffBlockOf(catalogPatch)+` table.hunk tr')).map((row) => {
  const [leftNumber, left, rightNumber, right] = Array.from(row.children);
  return {
    leftKind: left.dataset.kind,
    leftNumber: leftNumber.textContent,
    rightKind: right.dataset.kind,
    rightNumber: rightNumber.textContent,
    leftText: left.querySelector('code').textContent,
    rightText: right.querySelector('code').textContent,
  };
})`, &rows)

	type row = struct {
		LeftKind    string `json:"leftKind"`
		LeftNumber  string `json:"leftNumber"`
		RightKind   string `json:"rightKind"`
		RightNumber string `json:"rightNumber"`
		LeftText    string `json:"leftText"`
		RightText   string `json:"rightText"`
	}
	want := []row{
		{"context", "1", "context", "1", "func Limit() int {", "func Limit() int {"},
		{"removed", "2", "added", "2", "\treturn 2", "\t// the limit, raised"},
		{"absent", "", "added", "3", "", "\treturn 3"},
		{"context", "3", "context", "4", "}", "}"},
	}
	if diff := cmp.Diff(want, rows); diff != "" {
		t.Errorf("the rows of the patch differ (-want +got):\n%s", diff)
	}
}

func TestTheCatalogPaintsAPatchInALanguageThePageKnowsAndQuietsWhatChanged(t *testing.T) {
	// scenario: A patch naming a language the page highlights paints each side, and quiets a comment on a changed line
	t.Parallel()
	page := catalogued(t)
	page.waitFor(`!!document.querySelector('` + diffBlockOf(catalogPatch) + ` .hljs-keyword')`)

	var held struct {
		Keywords      int    `json:"keywords"`
		ChangedQuiet  string `json:"changedQuiet"`
		ChangedLoud   string `json:"changedLoud"`
		ContextLoud   string `json:"contextLoud"`
		QuietOnChange string `json:"quietOnChange"`
	}
	page.evaluate(`(() => {
  const block = document.querySelector('`+diffBlockOf(catalogPatch)+`');
  const colour = (selector) => getComputedStyle(block.querySelector(selector)).color;
  const probe = document.createElement('span');
  probe.style.color = 'var(--token-quiet-on-change)';
  block.appendChild(probe);
  const quiet = getComputedStyle(probe).color;
  probe.remove();
  return {
    keywords: block.querySelectorAll('.line .hljs-keyword').length,
    changedQuiet: colour('.line[data-kind="added"] .hljs-comment'),
    changedLoud: colour('.line[data-kind="added"] .hljs-keyword'),
    contextLoud: colour('.line[data-kind="context"] .hljs-keyword'),
    quietOnChange: quiet,
  };
})()`, &held)

	if held.Keywords < 3 {
		t.Errorf("the patch paints %d keywords, want every side painted", held.Keywords)
	}
	if held.ChangedQuiet != held.QuietOnChange {
		t.Errorf("a comment on an added line is %s, want the quiet colour %s", held.ChangedQuiet, held.QuietOnChange)
	}
	if held.ChangedLoud != held.ContextLoud {
		t.Errorf("a keyword on an added line is %s and on a context line %s, want the same", held.ChangedLoud, held.ContextLoud)
	}
}

func TestTheCatalogDrawsAPatchInAnUnknownLanguageAsPlainTextRows(t *testing.T) {
	// scenario: A patch naming a language the page does not highlight keeps its rows and paints none of them
	t.Parallel()
	page := catalogued(t)

	var held struct {
		Tokens int `json:"tokens"`
		Rows   int `json:"rows"`
	}
	page.evaluate(`({
  tokens: document.querySelectorAll('`+diffBlockOf(catalogPatchInUnknownTag)+` span[class^="hljs-"]').length,
  rows: document.querySelectorAll('`+diffBlockOf(catalogPatchInUnknownTag)+` table.hunk tr').length,
})`, &held)

	if held.Tokens != 0 || held.Rows == 0 {
		t.Errorf("the patch in an unknown language is %+v, want its rows and no painted token", held)
	}
}

func TestTheCatalogGivesAPatchTooWideForItsBlockAScrollOfItsOwn(t *testing.T) {
	// scenario: At a window narrower than the patch's table the table scrolls inside its block and not the page
	t.Parallel()
	page := catalogued(t)
	page.narrow(360, 780)

	var held struct {
		Overflow     string  `json:"overflow"`
		Scrolls      bool    `json:"scrolls"`
		BlockWithin  bool    `json:"blockWithin"`
		PageSideways float64 `json:"pageSideways"`
	}
	page.evaluate(`(() => {
  const block = document.querySelector('`+diffBlockOf(catalogPatch)+`');
  const hunk = block.querySelector('table.hunk');
  const state = block.closest('[data-state]');
  const page = document.documentElement;
  return {
    overflow: getComputedStyle(hunk).overflowX,
    scrolls: hunk.scrollWidth > hunk.clientWidth + 1,
    blockWithin: block.getBoundingClientRect().right <= state.getBoundingClientRect().right + 1,
    pageSideways: page.scrollWidth - page.clientWidth,
  };
})()`, &held)

	if held.Overflow != "auto" && held.Overflow != "scroll" {
		t.Errorf("the table overflows as %q, want a scroll of its own", held.Overflow)
	}
	if !held.Scrolls {
		t.Error("the patch does not scroll inside its table at 360px")
	}
	if !held.BlockWithin || held.PageSideways > 1 {
		t.Errorf("the patch spills out of its block %v or scrolls the page by %vpx", !held.BlockWithin, held.PageSideways)
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
