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

const rowsOfEachSourceLine = `(() => {
  const code = document.querySelector('.evidence.code code');
  const source = code.textContent;
  const sourceLines = source.split('\n');
  const starts = [];
  let offset = 0;
  for (const line of sourceLines) {
    starts.push(offset);
    offset += line.length + 1;
  }
  const marks = [];
  const walker = document.createTreeWalker(code, NodeFilter.SHOW_TEXT);
  let position = 0;
  for (let node = walker.nextNode(); node; node = walker.nextNode()) {
    for (let at = 0; at < node.data.length; at++, position++) {
      if (/\s/.test(node.data[at])) {
        continue;
      }
      const range = document.createRange();
      range.setStart(node, at);
      range.setEnd(node, at + 1);
      const box = range.getBoundingClientRect();
      marks.push({ position, top: Math.round(box.top), left: box.left });
    }
  }
  return sourceLines.map((line, index) => {
    const begin = starts[index];
    const rows = new Map();
    for (const mark of marks) {
      if (mark.position >= begin && mark.position < begin + line.length && !rows.has(mark.top)) {
        rows.set(mark.top, mark.left);
      }
    }
    return Array.from(rows.values());
  });
})()`

const sidewaysScrolling = `Array.from(document.querySelectorAll('.evidence.code pre')).filter((held) => held.scrollWidth > held.clientWidth + 1).length`

func TestALineOfGoWiderThanItsCardWrapsAndItsContinuationStartsFurtherIn(t *testing.T) {
	t.Parallel()
	page := openingAt(t, longGoDocument(t), 700, 900)
	page.waitFor(`document.querySelectorAll('.evidence.code code.hljs').length === 1`)

	var rows [][]float64
	page.evaluate(rowsOfEachSourceLine, &rows)

	for _, index := range []int{0, 2, 4} {
		if len(rows[index]) < 2 {
			t.Errorf("source line %d of the Go block is drawn on %d row, so a line wider than its card is not wrapped", index, len(rows[index]))
			continue
		}
		if rows[index][1] <= rows[index][0]+1 {
			t.Errorf("source line %d continues at %vpx after starting at %vpx, so its continuation does not start further in than the line it continues", index, rows[index][1], rows[index][0])
		}
	}
}

func TestNoCodeBlockScrollsSidewaysAtANarrowWindowOrAWideOne(t *testing.T) {
	t.Parallel()
	for name, width := range map[string]int{"narrow": 360, "wide": 1600} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			page := openingAt(t, longGoDocument(t), width, 900)
			page.waitFor(`document.querySelectorAll('.evidence.code code.hljs').length === 1`)

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
	page.waitFor(`document.querySelectorAll('.evidence.code code.hljs').length === 1`)

	var scrolling int
	page.evaluate(sidewaysScrolling, &scrolling)
	if scrolling != 0 {
		t.Error("markdown evidence scrolls sideways rather than wrapping at the width of its card")
	}

	var rows [][]float64
	page.evaluate(rowsOfEachSourceLine, &rows)
	if len(rows) < 3 || len(rows[2]) < 2 {
		t.Errorf("a long markdown line is not drawn on more than one row: %v", rows)
	}
}

const columnsOfAnIndentedLine = `(() => {
  const code = document.querySelector('.evidence.code code');
  const lefts = [];
  const walker = document.createTreeWalker(code, NodeFilter.SHOW_TEXT);
  const wanted = [{ line: 0, at: 0 }, { line: 0, at: 1 }, { line: 1, at: 0 }];
  const texts = [];
  for (let node = walker.nextNode(); node; node = walker.nextNode()) {
    texts.push(node);
  }
  const whole = texts.map((node) => node.data).join('');
  const lines = whole.split('\n');
  const place = (offset) => {
    let seen = 0;
    for (const node of texts) {
      if (offset < seen + node.data.length) {
        const range = document.createRange();
        range.setStart(node, offset - seen);
        range.setEnd(node, offset - seen + 1);
        return range.getBoundingClientRect().left;
      }
      seen += node.data.length;
    }
    return NaN;
  };
  const secondLineStart = lines[0].length + 1;
  return {
    text: whole,
    column: place(1) - place(0),
    indentedBy: place(secondLineStart + lines[1].indexOf('if')) - place(0),
  };
})()`

func TestATabIndentedLineOfGoKeepsItsFourColumnsAndItsTextIsIntact(t *testing.T) {
	t.Parallel()
	page := openingAt(t, longGoDocument(t), 900, 900)
	page.waitFor(`document.querySelectorAll('.evidence.code code.hljs').length === 1`)

	var measured struct {
		Text       string  `json:"text"`
		Column     float64 `json:"column"`
		IndentedBy float64 `json:"indentedBy"`
	}
	page.evaluate(columnsOfAnIndentedLine, &measured)

	if measured.Text != longGoLines {
		t.Errorf("the code on the page is %q rather than what the document stated", measured.Text)
	}
	if measured.Column <= 0 || measured.IndentedBy < measured.Column*3.9 || measured.IndentedBy > measured.Column*4.1 {
		t.Errorf("a line starting with one tab is indented %vpx where one column is %vpx, not the four columns it states", measured.IndentedBy, measured.Column)
	}
}

func TestEverySourceLineOfACodeBlockIsNumberedOnceInAGutterBesideIt(t *testing.T) {
	// scenario: A code block numbers each source line from 1 in a gutter, a wrapped line once, and the numbers are not part of its text
	t.Parallel()
	page := openingAt(t, longGoDocument(t), 700, 900)

	var numbered struct {
		Numbers   []string `json:"numbers"`
		Sources   int      `json:"sources"`
		Drawn     []bool   `json:"drawn"`
		TextAfter []bool   `json:"textAfter"`
		Text      string   `json:"text"`
	}
	page.evaluate(`(() => {
  const code = document.querySelector('.evidence.code code');
  const lines = Array.from(code.querySelectorAll('.code-line'));
  return {
    numbers: lines.map((line) => line.dataset.number),
    sources: code.textContent.replace(/\n$/, '').split('\n').length,
    drawn: lines.map((line) => parseFloat(getComputedStyle(line, '::before').width) > 0),
    textAfter: lines.map((line) => {
      const range = document.createRange();
      range.selectNodeContents(line);
      const gutter = parseFloat(getComputedStyle(line, '::before').width);
      return range.getClientRects()[0].left >= line.getBoundingClientRect().left + gutter - 0.5;
    }),
    text: code.textContent,
  };
})()`, &numbered)

	if len(numbered.Numbers) != numbered.Sources {
		t.Fatalf("the block draws %d numbered lines for %d source lines", len(numbered.Numbers), numbered.Sources)
	}
	for position, number := range numbered.Numbers {
		if number != fmt.Sprint(position+1) {
			t.Errorf("source line %d is numbered %q", position+1, number)
		}
		if !numbered.Drawn[position] {
			t.Errorf("source line %d draws no number in the gutter", position+1)
		}
		if !numbered.TextAfter[position] {
			t.Errorf("the code of source line %d starts inside the gutter", position+1)
		}
	}
	if numbered.Text != longGoLines {
		t.Errorf("numbering the lines changed the block's text to %q", numbered.Text)
	}
}
