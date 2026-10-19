package e2e_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
)

const (
	aLaptopWindow       = 1200
	aWideWindow         = 1440
	anUltrawideWindow   = 1920
	exactlyTwelveEighty = 1280
	aTallWindow         = 1000
)

const pairingDocument = `{
  "title": "options in pairs",
  "preamble": "what follows is what was seen",
  "decisions": [
    {
      "id": "p-1",
      "question": "which of three",
      "meaning": "three options and the way out",
      "mode": "single-choice",
      "options": [
        {"label": "first", "meaning": "the first", "command": "agree p-1"},
        {"label": "second", "meaning": "the second", "command": "patch p-1"},
        {"label": "third", "meaning": "the third", "command": "defer p-1"}
      ]
    },
    {
      "id": "p-2",
      "question": "which of these to keep",
      "meaning": "four to tick",
      "mode": "multiple-choice",
      "options": [
        {"label": "alpha", "meaning": "the alpha", "command": "keep p-2 alpha"},
        {"label": "beta", "meaning": "the beta", "command": "keep p-2 beta"},
        {"label": "gamma", "meaning": "the gamma", "command": "keep p-2 gamma"},
        {"label": "delta", "meaning": "the delta", "command": "keep p-2 delta"}
      ]
    },
    {
      "id": "p-3",
      "question": "which of two",
      "meaning": "two options and the way out make an odd count",
      "mode": "single-choice",
      "options": [
        {"label": "one", "meaning": "the one", "command": "agree p-3"},
        {"label": "two", "meaning": "the two", "command": "patch p-3"}
      ]
    }
  ]
}`

type cell struct {
	Label string `json:"label"`
	Left  int    `json:"left"`
	Top   int    `json:"top"`
	Width int    `json:"width"`
}

func theCellsOf(decision string) string {
	return fmt.Sprintf(`(() => {
  const group = document.querySelector('.decision[data-decision="%s"] .options');
  return {
    innerWidth: window.innerWidth,
    groupWidth: Math.round(group.getBoundingClientRect().width),
    cells: [...group.querySelectorAll('.option, .none')].map((held) => {
      const box = held.getBoundingClientRect();
      return {
        label: held.querySelector('.option-label').textContent.trim(),
        left: Math.round(box.left),
        top: Math.round(box.top),
        width: Math.round(box.width),
      };
    }),
  };
})()`, decision)
}

type laidOut struct {
	InnerWidth int    `json:"innerWidth"`
	GroupWidth int    `json:"groupWidth"`
	Cells      []cell `json:"cells"`
}

func laidOutIn(page *reader, decision string) laidOut {
	page.t.Helper()
	var read laidOut
	page.evaluate(theCellsOf(decision), &read)
	return read
}

func (l laidOut) rows() [][]cell {
	var rows [][]cell
	for _, held := range l.Cells {
		if len(rows) > 0 && rows[len(rows)-1][0].Top == held.Top {
			rows[len(rows)-1] = append(rows[len(rows)-1], held)
			continue
		}
		rows = append(rows, []cell{held})
	}
	return rows
}

func (l laidOut) labelsByRow() string {
	var rows []string
	for _, row := range l.rows() {
		var labels []string
		for _, held := range row {
			labels = append(labels, held.Label)
		}
		rows = append(rows, strings.Join(labels, " | "))
	}
	return strings.Join(rows, " / ")
}

func TestAWindow1440PxWidePairsTheOptionsOfASingleChoiceDecision(t *testing.T) {
	t.Parallel()
	page := openingAt(t, pairingDocument, aWideWindow, aTallWindow)

	read := laidOutIn(page, "p-1")
	want := "first | second / third | None of these"
	if got := read.labelsByRow(); got != want {
		t.Fatalf("in a window %dpx wide the options of a single-choice decision sit as %q, and they sit as %q", aWideWindow, got, want)
	}
	first, second, third := read.Cells[0], read.Cells[1], read.Cells[2]
	if first.Width != second.Width || first.Width != third.Width || first.Width != read.Cells[3].Width {
		t.Errorf("the columns are %d, %d, %d and %dpx wide, and they are all the same width", first.Width, second.Width, third.Width, read.Cells[3].Width)
	}
	if first.Left != third.Left || second.Left != read.Cells[3].Left {
		t.Errorf("the second row starts at %d and %d and the first at %d and %d, and the rows share two columns", third.Left, read.Cells[3].Left, first.Left, second.Left)
	}
}

func TestAWindow1440PxWidePairsTheCheckboxesOfAMultiSelectDecision(t *testing.T) {
	t.Parallel()
	page := openingAt(t, pairingDocument, aWideWindow, aTallWindow)

	read := laidOutIn(page, "p-2")
	want := "alpha | beta / gamma | delta"
	if got := read.labelsByRow(); got != want {
		t.Errorf("in a window %dpx wide the checkboxes of a multi-select decision sit as %q, and they sit as %q", aWideWindow, got, want)
	}
}

func TestAnOddNumberOfCellsAt1440PxLeavesTheLastOneInTheFirstColumnAtFullColumnWidth(t *testing.T) {
	t.Parallel()
	page := openingAt(t, pairingDocument, aWideWindow, aTallWindow)

	read := laidOutIn(page, "p-3")
	want := "one | two / None of these"
	if got := read.labelsByRow(); got != want {
		t.Fatalf("an odd number of cells sit as %q, and they sit as %q", got, want)
	}
	first, last := read.Cells[0], read.Cells[2]
	if last.Left != first.Left {
		t.Errorf("the last cell starts at %dpx and the first column at %dpx, and the last one sits in the first column", last.Left, first.Left)
	}
	if last.Width != first.Width {
		t.Errorf("the last cell is %dpx wide and the others %dpx, and it is as wide as they are", last.Width, first.Width)
	}
}

func TestAWindow1920PxWideStillHoldsTwoToARowAndNeverThree(t *testing.T) {
	t.Parallel()
	page := openingAt(t, pairingDocument, anUltrawideWindow, aTallWindow)

	for _, decision := range []string{"p-1", "p-2"} {
		read := laidOutIn(page, decision)
		for _, row := range read.rows() {
			if len(row) > 2 {
				t.Errorf("in a window %dpx wide %s has a row of %d options, and a row holds two at most", anUltrawideWindow, decision, len(row))
			}
		}
		if len(read.rows()[0]) != 2 {
			t.Errorf("in a window %dpx wide %s has %q, and its options still sit two to a row", anUltrawideWindow, decision, read.labelsByRow())
		}
	}
}

func TestAWindow1200PxWideStacksEveryOptionThoughTheGroupIsWiderThan600Px(t *testing.T) {
	t.Parallel()
	page := openingAt(t, pairingDocument, aLaptopWindow, aTallWindow)

	read := laidOutIn(page, "p-1")
	if read.GroupWidth < 600 {
		t.Fatalf("the group is %dpx wide in a window %dpx wide, and the scenario needs it wider than 600px", read.GroupWidth, aLaptopWindow)
	}
	if got, want := read.labelsByRow(), "first / second / third / None of these"; got != want {
		t.Errorf("in a window %dpx wide the options sit as %q, and they stack as %q", aLaptopWindow, got, want)
	}
}

func TestAWindowExactly1280PxWideStacksEveryOptionBecauseTheGroupIsBelow600Px(t *testing.T) {
	t.Parallel()
	page := openingAt(t, pairingDocument, aWideWindow, aTallWindow)
	page.narrow(exactlyTwelveEighty, aTallWindow)

	read := laidOutIn(page, "p-1")
	if read.InnerWidth != exactlyTwelveEighty {
		t.Fatalf("the window is %dpx wide, and the scenario needs exactly %d", read.InnerWidth, exactlyTwelveEighty)
	}
	if read.GroupWidth >= 600 {
		t.Fatalf("the group is %dpx wide at %dpx, and the scenario needs it below 600px", read.GroupWidth, exactlyTwelveEighty)
	}
	if got, want := read.labelsByRow(), "first / second / third / None of these"; got != want {
		t.Errorf("at exactly %dpx the options sit as %q, and they stack as %q", exactlyTwelveEighty, got, want)
	}
}

func TestANarrowWindowStacksEveryOptionAsItAlwaysHas(t *testing.T) {
	t.Parallel()
	page := openingAt(t, pairingDocument, 800, aTallWindow)

	read := laidOutIn(page, "p-2")
	if got, want := read.labelsByRow(), "alpha / beta / gamma / delta"; got != want {
		t.Errorf("in a window 800px wide the options sit as %q, and they stack as %q", got, want)
	}
}

func TestTabbingThroughAPairedDecisionMovesFocusInTheListedOrder(t *testing.T) {
	t.Parallel()
	page := openingAt(t, pairingDocument, aWideWindow, aTallWindow)

	read := laidOutIn(page, "p-2")
	if len(read.rows()[0]) != 2 {
		t.Fatalf("the decision sits as %q, and the scenario needs it paired", read.labelsByRow())
	}

	page.do(chromedp.Focus(chromedp.CSS(within("p-2", ".choice"))))
	var visited []string
	for range 4 {
		var held reached
		page.evaluate(focused, &held)
		visited = append(visited, strings.Fields(held.Name)[0])
		page.do(chromedp.KeyEvent(kb.Tab))
	}

	if got, want := strings.Join(visited, " "), "alpha beta gamma delta"; got != want {
		t.Errorf("tabbing through a paired decision reaches %q, and it reaches %q", got, want)
	}
}
