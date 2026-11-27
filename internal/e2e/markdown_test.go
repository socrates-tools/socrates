package e2e_test

import (
	"encoding/json"
	"fmt"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type drawnMarkdownBlock struct {
	Kind string `json:"kind"`
	Text string `json:"text"`
}

func markdownBlocksDrawn(count int) string {
	return fmt.Sprintf(`Array.from(document.querySelectorAll('.evidence.markdown ui-markdown-block')).filter((block) => block.drawnBlocks.length > 0).length === %d`, count)
}

func markdownDocument(t *testing.T, markdown string) string {
	t.Helper()
	evidence, err := json.Marshal([]map[string]string{{"language": "markdown", "code": markdown}})
	if err != nil {
		t.Fatal(err)
	}
	return `{"title": "one thing", "decisions": [{"id": "c-1", "question": "keep it", "mode": "single-choice",
 "options": [{"label": "keep it", "meaning": "as cut"}], "evidence": ` + string(evidence) + `}]}`
}

const aTicket = "# 1.2: Classify a document\n\n" +
	"## What to build\n\n" +
	"A **sure** pick wins.\nA second sentence.\n\n" +
	"## Acceptance criteria\n\n" +
	"- [ ] the type is shown\n" +
	"- [x] the calls are named\n\n" +
	"| Ticket | Why |\n| --- | --- |\n| 1.1 | the menu |\n\n" +
	"```go\nfunc Classify() {}\n```\n"

const markdownBlock = `document.querySelector('.evidence.markdown ui-markdown-block')`

func TestMarkdownEvidenceIsDrawnAsFormattedTextRatherThanAsItsSource(t *testing.T) {
	// scenario: A ticket's file shown as markdown evidence reads as headings, a task list, a table and a fence
	t.Parallel()
	_, page := read(t, markdownDocument(t, aTicket))
	page.waitFor(markdownBlocksDrawn(1))

	var drawn []drawnMarkdownBlock
	page.evaluate(markdownBlock+`.drawnBlocks`, &drawn)
	want := []drawnMarkdownBlock{
		{Kind: "h1", Text: "1.2: Classify a document"},
		{Kind: "h2", Text: "What to build"},
		{Kind: "p", Text: "A sure pick wins. A second sentence."},
		{Kind: "h2", Text: "Acceptance criteria"},
		{Kind: "ul", Text: "the type is shown\nthe calls are named"},
		{Kind: "table", Text: "Ticket\tWhy\n1.1\tthe menu"},
		{Kind: "pre", Text: "func Classify() {}"},
	}
	if diff := cmp.Diff(want, drawn); diff != "" {
		t.Errorf("the markdown is not drawn as formatted text (-want +got):\n%s", diff)
	}

	var tasks []bool
	page.evaluate(markdownBlock+`.drawnTasks`, &tasks)
	if diff := cmp.Diff([]bool{false, true}, tasks); diff != "" {
		t.Errorf("the task list does not draw a box for each criterion, ticked where the file ticks it (-want +got):\n%s", diff)
	}

	var elements []string
	page.evaluate(markdownBlock+`.drawnElements`, &elements)
	if !slices.Contains(elements, "strong") {
		t.Errorf("the bold word is not drawn bold: %v", elements)
	}
}

func TestMarkupWrittenInMarkdownEvidenceIsShownAsTheWordsItIsAndNeverBecomesMarkup(t *testing.T) {
	// scenario: Raw HTML and a script link in a markdown block are shown as text and no element or script runs from them
	t.Parallel()
	_, page := read(t, markdownDocument(t, "A line holding <img src=x onerror=\"window.ran=1\"> and <b>bold</b>.\n\n<script>window.ran=1</script>\n\n[a script](javascript:window.ran=1) and [a page](https://example.com) and [a file](../spec.md)\n"))
	page.waitFor(markdownBlocksDrawn(1))

	var elements []string
	page.evaluate(markdownBlock+`.drawnElements`, &elements)
	for _, forbidden := range []string{"img", "b", "script"} {
		if slices.Contains(elements, forbidden) {
			t.Errorf("markup written in the markdown became a %s element: %v", forbidden, elements)
		}
	}

	var drawn []drawnMarkdownBlock
	page.evaluate(markdownBlock+`.drawnBlocks`, &drawn)
	if len(drawn) < 2 || drawn[0].Text != `A line holding <img src=x onerror="window.ran=1"> and <b>bold</b>.` || drawn[1].Text != "<script>window.ran=1</script>" {
		t.Errorf("the markup is not shown as the words it is: %+v", drawn)
	}

	var links []string
	page.evaluate(markdownBlock+`.drawnLinks`, &links)
	if diff := cmp.Diff([]string{"https://example.com"}, links); diff != "" {
		t.Errorf("only a link to a web page is drawn as a link (-want +got):\n%s", diff)
	}

	var ran bool
	page.evaluate(`window.ran === 1`, &ran)
	if ran {
		t.Error("a script written in the markdown ran on the page")
	}
}
