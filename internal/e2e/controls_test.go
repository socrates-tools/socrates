package e2e_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/cdproto/accessibility"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/google/go-cmp/cmp"
)

const everyControlHoldsAShadowRootThatAdoptsTheSharedStylesheet = `['ui-button', 'ui-field', 'ui-copy', 'ui-expand'].map((name) => {
  const drawn = Array.from(document.querySelectorAll(name));
  return { Name: name, Drawn: drawn.length, Adopting: drawn.filter((held) => held.shadowRoot !== null && held.shadowRoot.adoptedStyleSheets.length > 0).length };
})`

type shadowed struct {
	Name     string
	Drawn    int
	Adopting int
}

func (r *reader) controlsInShadowRoots() []shadowed {
	r.t.Helper()
	var held []shadowed
	r.evaluate(everyControlHoldsAShadowRootThatAdoptsTheSharedStylesheet, &held)
	return held
}

func TestAShadowRootHoldsTheButtonTheFieldTheCopyControlAndTheExpandControlOfTheCatalog(t *testing.T) {
	// scenario: Every ui-button, ui-field, ui-copy and ui-expand the catalog draws renders into a shadow root that adopts the shared stylesheet
	t.Parallel()
	page := catalogued(t)

	for _, each := range page.controlsInShadowRoots() {
		if each.Drawn == 0 {
			t.Errorf("the catalog draws no %s, so this proves nothing of it", each.Name)
		}
		if each.Adopting != each.Drawn {
			t.Errorf("of %d %s the catalog draws, %d have a shadow root that adopts the shared stylesheet, want all", each.Drawn, each.Name, each.Adopting)
		}
	}
	if violations := page.violations(); len(violations) != 0 {
		t.Errorf("the catalog's policy blocked %v while the controls adopted the shared stylesheet", violations)
	}
}

func TestAShadowRootHoldsTheButtonTheFieldTheCopyControlAndTheExpandControlOfThePage(t *testing.T) {
	// scenario: Every ui-button, ui-field, ui-copy and ui-expand the page draws renders into a shadow root that adopts the shared stylesheet
	t.Parallel()
	held := standing(t)
	aDrawing(t, held, "counting.html", drawingWithACounter)
	page := opening(t, held.serve(illustrating(t, copyingDocument(t), "counting.html")))
	page.waitFor(highlightedCodeBlocks(2))

	for _, each := range page.controlsInShadowRoots() {
		if each.Drawn == 0 {
			t.Errorf("the page draws no %s, so this proves nothing of it", each.Name)
		}
		if each.Adopting != each.Drawn {
			t.Errorf("of %d %s the page draws, %d have a shadow root that adopts the shared stylesheet, want all", each.Drawn, each.Name, each.Adopting)
		}
	}
	if violations := page.violations(); len(violations) != 0 {
		t.Errorf("the page's policy blocked %v while the controls adopted the shared stylesheet", violations)
	}
}

func TestAPressOnAFieldsElementPutsTheCursorInTheTextareaItHoldsAndTypingFillsIt(t *testing.T) {
	// scenario: Clicking a field's element, and calling focus() on it, each let a person type into the control inside it and read it back through the field's value
	t.Parallel()
	_, page := read(t, oneDecisionDocument)
	field := within("c-1", ".words-field")

	page.click(field)
	page.press("by click ")
	page.evaluate(`document.querySelector('.message-field').focus(); true`, nil)
	page.press("by focus")

	var held struct{ Words, Message string }
	page.evaluate(`({ Words: document.querySelector(`+jsonOf(t, field)+`).value, Message: document.querySelector('.message-field').value })`, &held)
	if diff := cmp.Diff(struct{ Words, Message string }{"by click ", "by focus"}, held); diff != "" {
		t.Errorf("what a person typed differs from what the fields hold (-want +got):\n%s", diff)
	}
}

func TestACallToClickOnAButtonsElementPressesItAndASubmitButtonSubmitsWithItselfAsTheSubmitter(t *testing.T) {
	// scenario: click() on the finish button's element sends and ends, and on the send button's element sends and waits
	t.Parallel()
	for name, ends := range map[string]bool{".send": false, ".finish": true} {
		t.Run(name, func(t *testing.T) {
			held, page := read(t, oneDecisionDocument)
			page.evaluate(chose("c-1", 1), nil)

			page.evaluate(`document.querySelector(`+jsonOf(t, name)+`).click(); true`, nil)

			delivered, arrived := held.awaiting(5 * time.Second)
			if !arrived {
				t.Fatalf("click() on %s delivered nothing", name)
			}
			if delivered.Ended != ends {
				t.Errorf("click() on %s delivered ended=%v, want %v", name, delivered.Ended, ends)
			}
			if violations := page.violations(); len(violations) != 0 {
				t.Errorf("a real form navigation was attempted and blocked by the policy: %v", violations)
			}
		})
	}
}

func TestEnterInTheMessageFieldBreaksTheLineAndSendsNothing(t *testing.T) {
	// scenario: Enter in a field's textarea writes a new line as it does in any textarea, and delivers nothing
	t.Parallel()
	held, page := read(t, oneDecisionDocument)

	page.evaluate(`document.querySelector('.message-field').focus(); true`, nil)
	page.press("one")
	page.press(kb.Enter)
	page.press("two")

	var message string
	page.evaluate(`document.querySelector('.message-field').value`, &message)
	if message != "one\ntwo" {
		t.Errorf("the message holds %q, want the two lines typed", message)
	}
	if delivered, arrived := held.awaiting(shortWindow); arrived {
		t.Errorf("Enter in the message delivered an answer: %+v", delivered)
	}
}

type toldToAScreenReader struct {
	Role     string
	Name     string
	Expanded string
	Said     string
}

func stringOf(value *accessibility.Value) string {
	if value == nil {
		return ""
	}
	var held string
	if err := json.Unmarshal(value.Value, &held); err != nil {
		return strings.Trim(string(value.Value), `"`)
	}
	return held
}

func (r *reader) whatAScreenReaderIsTold() []toldToAScreenReader {
	r.t.Helper()

	var nodes []*accessibility.Node
	r.do(chromedp.Func(func(ctx context.Context, target *chromedp.Target) error {
		if _, err := cdp.Call(ctx, target, accessibility.Enable, cdp.Empty{}); err != nil {
			return err
		}
		result, err := cdp.Call(ctx, target, accessibility.GetFullAXTree, accessibility.GetFullAXTreeParams{})
		if err != nil {
			return err
		}
		nodes = result.Nodes
		return nil
	}))

	byID := map[accessibility.NodeID]*accessibility.Node{}
	for _, node := range nodes {
		byID[node.NodeID] = node
	}
	var said func(node *accessibility.Node) string
	said = func(node *accessibility.Node) string {
		if stringOf(node.Role) == "StaticText" {
			return stringOf(node.Name)
		}
		var words []string
		for _, child := range node.ChildIDs {
			if found, known := byID[child]; known {
				words = append(words, said(found))
			}
		}
		return strings.Join(words, "")
	}

	var told []toldToAScreenReader
	for _, node := range nodes {
		if node.Ignored {
			continue
		}
		each := toldToAScreenReader{Role: stringOf(node.Role), Name: stringOf(node.Name), Said: said(node)}
		for _, property := range node.Properties {
			if property.Name == accessibility.PropertyNameExpanded {
				each.Expanded = stringOf(property.Value)
			}
		}
		told = append(told, each)
	}
	return told
}

func (r *reader) buttonsNamed(name string) []toldToAScreenReader {
	r.t.Helper()
	var found []toldToAScreenReader
	for _, each := range r.whatAScreenReaderIsTold() {
		if each.Role == "button" && each.Name == name {
			found = append(found, each)
		}
	}
	return found
}

func TestTheCatalogsCopyAndExpandControlsAreButtonsNamedToAScreenReaderAndTheExpandStatesItsState(t *testing.T) {
	// scenario: The catalog's copy controls are buttons named Copy code, and its expand controls are buttons named Expand that say whether they are expanded
	t.Parallel()
	page := catalogued(t)

	if got := len(page.buttonsNamed("Copy code")); got != 2 {
		t.Errorf("a screen reader finds %d buttons named Copy code, want the two copy states", got)
	}
	var states []string
	for _, each := range page.buttonsNamed("Expand") {
		states = append(states, each.Expanded)
	}
	if diff := cmp.Diff([]string{"false", "true"}, states); diff != "" {
		t.Errorf("the expand buttons a screen reader finds say they are expanded as (-want +got):\n%s", diff)
	}
}

func TestACopyControlTellsAScreenReaderWhereItWasPressedThatTheCopyWasMade(t *testing.T) {
	// scenario: After a copy the one live status of the page that holds words says Copied, and it is empty again once the control rests
	t.Parallel()
	_, page := read(t, copyingDocument(t))
	page.mayCopy()
	page.waitFor(highlightedCodeBlocks(2))

	statuses := func() []string {
		var said []string
		for _, each := range page.whatAScreenReaderIsTold() {
			if each.Role == "status" && each.Said != "" {
				said = append(said, each.Said)
			}
		}
		return said
	}
	if before := statuses(); len(before) != 0 {
		t.Fatalf("before any copy a live status already says %v", before)
	}

	page.click(copyControls)
	page.waitFor(`document.querySelector('` + copyControls + `').dataset.copied === 'yes'`)
	if diff := cmp.Diff([]string{"Copied"}, statuses()); diff != "" {
		t.Errorf("after the copy the live statuses say (-want +got):\n%s", diff)
	}

	page.waitFor(`document.querySelector('` + copyControls + `').dataset.copied === undefined`)
	if after := statuses(); len(after) != 0 {
		t.Errorf("once the control rested a live status still says %v", after)
	}
}

func TestOneEnterOnAChoiceSubmitsTheFormOnceAndSendsWithoutEnding(t *testing.T) {
	// scenario: A single Enter on a focused choice raises one submit of the exchange form and delivers one answer that does not end
	t.Parallel()
	held, page := read(t, oneDecisionDocument)
	page.evaluate(`window.__submits = []; document.addEventListener('submit', (event) => { window.__submits.push('submit'); }, true); true`, nil)
	page.evaluate(chose("c-1", 1), nil)
	page.focus(within("c-1", ".choice"))

	page.press(kb.Enter)

	delivered, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("Enter on a choice delivered nothing")
	}
	if delivered.Ended {
		t.Error("Enter on a choice finished the exchange, and it presses the send button")
	}
	if again, arrivedAgain := held.awaiting(shortWindow); arrivedAgain {
		t.Errorf("one Enter on a choice delivered a second answer: %+v", again)
	}
	var submits []string
	page.evaluate(`window.__submits`, &submits)
	if len(submits) != 1 {
		t.Errorf("one Enter on a choice raised %d submits of the form, want 1", len(submits))
	}
}
