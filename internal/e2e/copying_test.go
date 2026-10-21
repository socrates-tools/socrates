package e2e_test

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp/kb"
	"github.com/google/go-cmp/cmp"
)

const indentedGo = "func Select(options []Option) error {\n" +
	"\tfor _, option := range options {\n" +
	"\t\tif option.Disabled {\n" +
	"\n" +
	"\t\t\tcontinue\n" +
	"\t\t}\n" +
	"\t\tlog.Printf(\"selecting %q among %d options, none of which a caller could have asked for without a long line\", option.Name, len(options))\n" +
	"\t}\n" +
	"    return nil   \n" +
	"}\n"

const aMarkdownNote = "# Why\n\n- first\n- second\n\n```go\nx := 1\n```\n"

func copyingDocument(t *testing.T) string {
	t.Helper()
	blocks, err := json.Marshal([]map[string]string{
		{"caption": "as a prop", "language": "go", "code": indentedGo},
		{"language": "markdown", "code": aMarkdownNote},
		{"caption": "the lock", "patch": "--- a/internal/cache/write.go\n+++ b/internal/cache/write.go\n@@ -1,1 +1,1 @@\n-a\n+b\n"},
	})
	if err != nil {
		t.Fatal(err)
	}
	return `{
  "title": "copying",
  "decisions": [
    {"id": "c-1", "question": "should the lookup retry", "mode": "single-choice",
     "options": [{"label": "retry", "meaning": "tolerate a blip"}, {"label": "fail", "meaning": "report at once"}],
     "evidence": ` + string(blocks) + `}
  ]
}`
}

func (r *reader) mayCopy() {
	r.t.Helper()
	r.do(sending(emulation.SetFocusEmulationEnabled, emulation.SetFocusEmulationEnabledParams{Enabled: true}), sending(browser.SetPermission, browser.SetPermissionParams{
		Permission: &browser.PermissionDescriptor{Name: "clipboard-write"},
		Setting:    browser.PermissionSettingGranted,
	}), sending(browser.SetPermission, browser.SetPermissionParams{
		Permission: &browser.PermissionDescriptor{Name: "clipboard-read"},
		Setting:    browser.PermissionSettingGranted,
	}))
}

func (r *reader) onTheClipboard() string {
	r.t.Helper()
	var held string
	r.evaluate(`navigator.clipboard.readText()`, &held, awaitingThePromise)
	return held
}

func (r *reader) clipboardHolds(text string) {
	r.t.Helper()
	r.evaluate(`navigator.clipboard.writeText(`+jsonOf(r.t, text)+`)`, nil, awaitingThePromise)
}

func jsonOf(t *testing.T, value string) string {
	t.Helper()
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return string(encoded)
}

const copyControls = `.evidence.code ui-copy.copy`

func TestEveryCodeBlockOffersACopyControlAndASplitDiffOffersNone(t *testing.T) {
	t.Parallel()
	_, page := read(t, copyingDocument(t))

	var counts struct {
		Code     int `json:"code"`
		Controls int `json:"controls"`
		Diffs    int `json:"diffs"`
		InDiffs  int `json:"inDiffs"`
	}
	page.evaluate(`({
  code: document.querySelectorAll('.evidence.code').length,
  controls: document.querySelectorAll('.evidence.code > .evidence-heading ui-copy.copy').length,
  diffs: document.querySelectorAll('.evidence.diff').length,
  inDiffs: document.querySelectorAll('.evidence.diff ui-copy.copy').length,
})`, &counts)

	if counts.Code != 2 || counts.Diffs != 1 {
		t.Fatalf("the page shows %d code blocks and %d split diffs, want 2 and 1", counts.Code, counts.Diffs)
	}
	if counts.Controls != counts.Code {
		t.Errorf("%d of %d code blocks show a copy control in their heading", counts.Controls, counts.Code)
	}
	if counts.InDiffs != 0 {
		t.Errorf("a split diff shows %d copy controls, want none", counts.InDiffs)
	}
}

func TestPressingACopyControlPutsTheBlocksCodeOnTheClipboardExactlyAsTheDocumentStatedIt(t *testing.T) {
	t.Parallel()
	_, page := read(t, copyingDocument(t))
	page.mayCopy()
	page.waitFor(`document.querySelectorAll('.evidence.code code.hljs').length === 2`)

	for position, want := range []string{indentedGo, aMarkdownNote} {
		page.clipboardHolds("something else entirely")
		page.evaluate(`document.querySelectorAll('`+copyControls+`')[`+strconv.Itoa(position)+`].click()`, nil)
		page.waitFor(`navigator.clipboard.readText().then((held) => held !== "something else entirely")`)

		if diff := cmp.Diff(want, page.onTheClipboard()); diff != "" {
			t.Errorf("the clipboard does not hold the block's code as the document stated it (-want +got):\n%s", diff)
		}
	}
}

func TestACopyControlIsReachedAndPressedByKeyboardAndIsNamedToAScreenReader(t *testing.T) {
	t.Parallel()
	_, page := read(t, copyingDocument(t))
	page.mayCopy()

	var control struct {
		Title    string `json:"title"`
		Hint     string `json:"hint"`
		Disabled bool   `json:"disabled"`
		Part     string `json:"part"`
	}
	page.evaluate(`(() => {
  const held = document.querySelector('`+copyControls+`');
  return { title: held.dataset.title || '', hint: held.dataset.hint || '', disabled: held.hasAttribute('disabled'), part: held.dataset.interaction };
})()`, &control)

	if named := page.buttonsNamed("Copy code"); len(named) != 2 {
		t.Errorf("a screen reader finds %d buttons named Copy code, want one for each of the two code blocks", len(named))
	}
	if control.Title == "" || control.Hint == "" {
		t.Errorf("the copy control is titled %q and hinted %q, and it must say what it does on hover and focus", control.Title, control.Hint)
	}
	if control.Disabled || control.Part != "viewing" {
		t.Errorf("the copy control is %+v, and it must be enabled and take no part in the answer", control)
	}

	reached := false
	for range 40 {
		page.press(kb.Tab)
		page.evaluate(`document.activeElement === document.querySelector('`+copyControls+`')`, &reached)
		if reached {
			break
		}
	}
	if !reached {
		t.Fatal("the copy control cannot be reached by keyboard alone")
	}
	page.clipboardHolds("something else entirely")
	page.waitFor(`!document.querySelector('.tooltip').hidden && document.querySelector('.tooltip').heading !== ''`)
	page.press(kb.Enter)
	page.waitFor(`navigator.clipboard.readText().then((held) => held !== "something else entirely")`)

	if diff := cmp.Diff(indentedGo, page.onTheClipboard()); diff != "" {
		t.Errorf("pressing Enter on the focused control did not copy the block (-want +got):\n%s", diff)
	}
}

func TestACopyControlConfirmsTheCopyWhereItWasPressed(t *testing.T) {
	t.Parallel()
	_, page := read(t, copyingDocument(t))
	page.mayCopy()

	page.click(copyControls)
	page.waitFor(`document.querySelector('` + copyControls + `').dataset.copied === 'yes'`)

	var confirmed struct {
		Title  string `json:"title"`
		Others int    `json:"others"`
	}
	page.evaluate(`({
  title: document.querySelector('.tooltip').heading,
  others: document.querySelectorAll('.copy[data-copied]').length,
})`, &confirmed)

	if confirmed.Title != "Copied" {
		t.Errorf("after the copy the tooltip says %q, want Copied", confirmed.Title)
	}
	if confirmed.Others != 1 {
		t.Errorf("%d controls show a copy, want only the one pressed", confirmed.Others)
	}

	page.waitFor(`!document.querySelector('` + copyControls + `').dataset.copied`)
}

func TestPressingACopyControlChangesNothingAboutTheAnswerAndSendsNothing(t *testing.T) {
	t.Parallel()
	held, page := read(t, copyingDocument(t))
	page.mayCopy()

	page.evaluate(chose("c-1", 1), nil)
	page.write(within("c-1", ".words-field"), "a note")
	before := len(page.requests.all())

	page.click(copyControls)
	page.waitFor(`document.querySelector('` + copyControls + `').dataset.copied === 'yes'`)

	if after := len(page.requests.all()); after != before {
		t.Errorf("pressing copy made the page issue %d more requests", after-before)
	}
	if delivered, arrived := held.awaiting(shortWindow); arrived {
		t.Fatalf("pressing copy delivered an answer: %+v", delivered)
	}

	page.click(".send")
	sent, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("sending delivered nothing at all")
	}
	if len(sent.Decisions) != 1 || len(sent.Decisions[0].Chosen) != 1 || sent.Decisions[0].Chosen[0] != "fail" || sent.Decisions[0].Words != "a note" {
		t.Errorf("after pressing copy the answer is %+v, want the one option and the words that were left in the controls", sent.Decisions)
	}
}
