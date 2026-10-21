package e2e_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

const longGoLines = "func Lookup(ctx context.Context, identifier string, options ...LookupOption) (Settlement, error) {\n" +
	"\tif identifier == \"\" {\n" +
	"\t\treturn Settlement{}, fmt.Errorf(\"lookup of %q for the settlement owner failed after %d retries with backoff %s\", identifier, retries, backoff)\n" +
	"\t}\n" +
	"\treturn service.Lookup(ctx, identifier, withTimeout(30*time.Second), withRetries(5), withBackoff(time.Second))\n" +
	"}\n"

func longGoDocument(t *testing.T) string {
	t.Helper()
	evidence, err := json.Marshal([]map[string]string{{"language": "go", "code": longGoLines}})
	if err != nil {
		t.Fatal(err)
	}
	return `{
  "title": "one thing to settle",
  "decisions": [
    {"id": "c-1", "question": "should the lookup retry", "mode": "single-choice",
     "options": [{"label": "retry", "meaning": "tolerate a blip"}, {"label": "fail", "meaning": "report at once"}],
     "evidence": ` + string(evidence) + `}
  ]
}`
}

const columnsOfEachRowOfEachSourceLine = `document.querySelector('.evidence.code ui-code-block').drawnLines.map((line) => line.rows)`

func highlightedCodeBlocks(count int) string {
	return fmt.Sprintf(`Array.from(document.querySelectorAll('.evidence.code ui-code-block')).filter((block) => block.highlighted).length === %d`, count)
}

const sidewaysScrolling = `Array.from(document.querySelectorAll('.evidence.code ui-code-block')).filter((held) => held.scrollWidth > held.clientWidth + 1).length`

func TestALineOfGoWiderThanItsCardWrapsAndItsContinuationStartsFurtherIn(t *testing.T) {
	t.Parallel()
	page := openingAt(t, longGoDocument(t), 700, 900)
	page.waitFor(highlightedCodeBlocks(1))

	var rows [][]float64
	page.evaluate(columnsOfEachRowOfEachSourceLine, &rows)

	for _, index := range []int{0, 2, 4} {
		if len(rows[index]) < 2 {
			t.Errorf("source line %d of the Go block is drawn on %d row, so a line wider than its card is not wrapped", index, len(rows[index]))
			continue
		}
		if rows[index][1] <= rows[index][0]+1 {
			t.Errorf("source line %d continues at column %v after starting at column %v, so its continuation does not start further in than the line it continues", index, rows[index][1], rows[index][0])
		}
	}
}

func TestNoCodeBlockScrollsSidewaysAtANarrowWindowOrAWideOne(t *testing.T) {
	t.Parallel()
	for name, width := range map[string]int{"narrow": 360, "wide": 1600} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			page := openingAt(t, longGoDocument(t), width, 900)
			page.waitFor(highlightedCodeBlocks(1))

			var scrolling int
			page.evaluate(sidewaysScrolling, &scrolling)
			if scrolling != 0 {
				t.Errorf("%d code blocks scroll sideways at %dpx, so a long line is cut off until the reader scrolls", scrolling, width)
			}
		})
	}
}

func TestMarkdownEvidenceStillWrapsAtTheWidthOfItsCard(t *testing.T) {
	t.Parallel()
	prose := strings.Repeat("a ticket sentence that is long enough to need wrapping ", 6)
	evidence, err := json.Marshal([]map[string]string{{"language": "markdown", "code": "# a ticket\n\n- [ ] " + prose + "\n"}})
	if err != nil {
		t.Fatal(err)
	}
	document := `{"title": "one thing", "decisions": [{"id": "c-1", "question": "keep it", "mode": "single-choice",
 "options": [{"label": "keep it", "meaning": "as cut"}], "evidence": ` + string(evidence) + `}]}`
	page := openingAt(t, document, 500, 900)
	page.waitFor(highlightedCodeBlocks(1))

	var scrolling int
	page.evaluate(sidewaysScrolling, &scrolling)
	if scrolling != 0 {
		t.Error("markdown evidence scrolls sideways rather than wrapping at the width of its card")
	}

	var rows [][]float64
	page.evaluate(columnsOfEachRowOfEachSourceLine, &rows)
	if len(rows) < 3 || len(rows[2]) < 2 {
		t.Errorf("a long markdown line is not drawn on more than one row: %v", rows)
	}
}

func TestATabIndentedLineOfGoKeepsItsFourColumnsAndItsTextIsIntact(t *testing.T) {
	t.Parallel()
	page := openingAt(t, longGoDocument(t), 900, 900)
	page.waitFor(highlightedCodeBlocks(1))

	var measured struct {
		Text    string      `json:"text"`
		Columns [][]float64 `json:"columns"`
	}
	page.evaluate(`({
  text: document.querySelector('.evidence.code ui-code-block').statedText(),
  columns: document.querySelector('.evidence.code ui-code-block').drawnLines.map((line) => line.rows.map((columns) => Math.round(columns * 10) / 10)),
})`, &measured)

	if measured.Text != longGoLines {
		t.Errorf("the code on the page is %q rather than what the document stated", measured.Text)
	}
	if len(measured.Columns) < 2 || len(measured.Columns[0]) == 0 || len(measured.Columns[1]) == 0 || measured.Columns[0][0] != 0 || measured.Columns[1][0] != 4 {
		t.Errorf("the first two lines start at columns %v, want a line starting with one tab at column 4 beside a line with none at column 0", measured.Columns)
	}
}

func TestEverySourceLineOfACodeBlockIsNumberedOnceInAGutterBesideIt(t *testing.T) {
	// scenario: A code block numbers each source line from 1 in a gutter, a wrapped line once, and the numbers are not part of its text
	t.Parallel()
	page := openingAt(t, longGoDocument(t), 700, 900)
	page.waitFor(highlightedCodeBlocks(1))

	var numbered struct {
		Lines   []drawnLine `json:"lines"`
		Sources int         `json:"sources"`
		Text    string      `json:"text"`
	}
	page.evaluate(`(() => {
  const block = document.querySelector('.evidence.code ui-code-block');
  return {
    lines: block.drawnLines.map((line) => ({ number: line.number, gutter: line.gutter, rows: line.rows })),
    sources: block.statedText().replace(/\n$/, '').split('\n').length,
    text: block.statedText(),
  };
})()`, &numbered)

	if len(numbered.Lines) != numbered.Sources {
		t.Fatalf("the block draws %d numbered lines for %d source lines", len(numbered.Lines), numbered.Sources)
	}
	for position, line := range numbered.Lines {
		if line.Number != fmt.Sprint(position+1) {
			t.Errorf("source line %d is numbered %q", position+1, line.Number)
		}
		if line.Gutter <= 0 {
			t.Errorf("source line %d draws no number in the gutter", position+1)
		}
		if len(line.Rows) == 0 || line.Rows[0] < 0 {
			t.Errorf("the code of source line %d starts inside the gutter", position+1)
		}
	}
	if numbered.Text != longGoLines {
		t.Errorf("numbering the lines changed the block's text to %q", numbered.Text)
	}
}

const aCodeBlockIsHighlighted = `Array.from(document.querySelectorAll('ui-code-block')).some((block) => block.highlighted)`
