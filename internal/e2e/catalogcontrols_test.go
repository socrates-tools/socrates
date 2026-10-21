package e2e_test

import (
	"bytes"
	"fmt"
	"image/png"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/google/go-cmp/cmp"
)

const (
	catalogCopiesThis  = "Copy this exactly\n\tas it was written, with a trailing space "
	theOwnerOffersThis = "what an owner offered instead"
	somethingElseHeld  = "something else entirely"
	anIdlePressWindow  = 400 * time.Millisecond
)

func stateOf(primitive, key string) string {
	return fmt.Sprintf(`[data-primitive="%s"] [data-state="%s"]`, primitive, key)
}

func countPresses(page *reader, within string) {
	page.evaluate(`(() => {
  const held = document.querySelector(`+jsonOf(page.t, within)+`);
  held.__presses = 0;
  held.addEventListener('click', () => { held.__presses += 1; });
  return true;
})()`, nil)
}

func pressesIn(page *reader, within string) int {
	var counted int
	page.evaluate(`document.querySelector(`+jsonOf(page.t, within)+`).__presses`, &counted)
	return counted
}

func TestTheCatalogPressesEveryEnabledButtonStateWithThePointer(t *testing.T) {
	// scenario: Pressing an enabled button state in the catalog is one press
	t.Parallel()
	page := catalogued(t)

	for _, key := range []string{"primary", "raised", "secondary", "quiet", "link"} {
		t.Run(key, func(t *testing.T) {
			state := stateOf("ui-button", key)
			countPresses(page, state)

			page.click(state + " ui-button")

			if got := pressesIn(page, state); got != 1 {
				t.Errorf("pressing the %s button counted %d presses, want 1", key, got)
			}
		})
	}
}

func TestTheCatalogPressesAButtonWithTheKeyboard(t *testing.T) {
	// scenario: Enter and Space on a focused button state in the catalog each press it once
	t.Parallel()
	page := catalogued(t)

	state := stateOf("ui-button", "primary")
	countPresses(page, state)
	page.evaluate(`document.querySelector('`+state+` ui-button').focus()`, nil)

	page.press(kb.Enter)
	page.press(" ")

	if got := pressesIn(page, state); got != 2 {
		t.Errorf("Enter and Space on a focused button counted %d presses, want 2", got)
	}
}

func TestTheCatalogDoesNotPressADisabledButton(t *testing.T) {
	// scenario: A disabled button state in the catalog takes no press
	t.Parallel()
	page := catalogued(t)

	state := stateOf("ui-button", "disabled")
	countPresses(page, state)
	var disabled bool
	page.evaluate(`document.querySelector('`+state+` ui-button').disabled`, &disabled)
	page.evaluate(`document.querySelector('`+state+` ui-button').click()`, nil)

	if !disabled {
		t.Error("the disabled state draws a button that is not disabled")
	}
	if got := pressesIn(page, state); got != 0 {
		t.Errorf("a disabled button counted %d presses, want none", got)
	}
}

func TestTheCatalogTakesTypingInAnEnabledField(t *testing.T) {
	// scenario: Typing into an enabled field state in the catalog fills it and tells the page on every key
	t.Parallel()
	page := catalogued(t)

	state := stateOf("ui-field", "empty")
	page.evaluate(`(() => {
  const control = document.querySelector('`+state+` ui-field');
  window.__typed = 0;
  control.addEventListener('input', () => { window.__typed += 1; });
  return true;
})()`, nil)

	page.write(state+" ui-field", "in my own words")

	var held struct {
		Value string `json:"value"`
		Typed int    `json:"typed"`
	}
	page.evaluate(`({ value: document.querySelector('`+state+` ui-field').value, typed: window.__typed })`, &held)
	if held.Value != "in my own words" {
		t.Errorf("the field holds %q after typing, want what was typed", held.Value)
	}
	if held.Typed != len("in my own words") {
		t.Errorf("the field announced %d inputs for %d keys", held.Typed, len("in my own words"))
	}
}

func (r *reader) inkInsideTheWritingOf(selector string) int {
	r.t.Helper()

	r.evaluate(`document.querySelector(`+jsonOf(r.t, selector)+`).scrollIntoView({ block: 'center' }); new Promise((settled) => requestAnimationFrame(() => requestAnimationFrame(settled)))`, nil, awaitingThePromise)
	var box struct{ Left, Top, Width, Height float64 }
	r.evaluate(`(() => { const box = document.querySelector(`+jsonOf(r.t, selector)+`).getBoundingClientRect(); return { Left: box.left, Top: box.top, Width: box.width, Height: box.height }; })()`, &box)

	shot, err := png.Decode(bytes.NewReader(r.screenshot()))
	if err != nil {
		r.t.Fatalf("reading what the page painted: %v", err)
	}
	left, top, right, bottom := int(box.Left)+8, int(box.Top+box.Height)-56, int(box.Left+box.Width)-28, int(box.Top+box.Height)-14
	red, green, blue, _ := shot.At(left, top).RGBA()
	inked := 0
	for x := left; x < right; x++ {
		for y := top; y < bottom; y++ {
			r2, g2, b2, _ := shot.At(x, y).RGBA()
			if abs(int(r2>>8)-int(red>>8))+abs(int(g2>>8)-int(green>>8))+abs(int(b2>>8)-int(blue>>8)) > 60 {
				inked++
			}
		}
	}
	return inked
}

func TestTheCatalogShowsAFieldsPlaceholderUntilSomethingIsTyped(t *testing.T) {
	// scenario: A field state that carries a placeholder draws words in its writing area while empty and draws none once a space is typed
	t.Parallel()
	page := catalogued(t)

	state := stateOf("ui-field", "placeholder") + " ui-field"
	before := page.inkInsideTheWritingOf(state)

	page.write(state, " ")
	page.evaluate(`document.activeElement.blur(); true`, nil)

	after := page.inkInsideTheWritingOf(state)
	if before == 0 || after != 0 {
		t.Errorf("the writing area held %d inked pixels before typing and %d after a space, want some and then none", before, after)
	}
}

func TestTheCatalogTakesNoTypingInADisabledField(t *testing.T) {
	// scenario: A disabled field state in the catalog cannot be focused and holds nothing
	t.Parallel()
	page := catalogued(t)

	state := stateOf("ui-field", "disabled")
	var held struct {
		Disabled bool   `json:"disabled"`
		Focused  bool   `json:"focused"`
		Value    string `json:"value"`
	}
	page.evaluate(`(() => {
  const control = document.querySelector('`+state+` ui-field');
  control.focus();
  return { disabled: control.disabled, focused: document.activeElement === control, value: control.value };
})()`, &held)

	if !held.Disabled || held.Focused || held.Value != "" {
		t.Errorf("the disabled field is %+v, want disabled, unfocusable and empty", held)
	}
}

func TestTheCatalogChoosesAnUnchosenRadioAndUnchoosesItsGroupMate(t *testing.T) {
	// scenario: Choosing one radio of a group in the catalog unchooses the one that was chosen
	t.Parallel()
	page := catalogued(t)

	unchosen := stateOf("ui-choice", "radio") + " input"
	chosen := stateOf("ui-choice", "radio-chosen") + " input"
	var before struct{ Unchosen, Chosen bool }
	page.evaluate(`({ Unchosen: document.querySelector('`+unchosen+`').checked, Chosen: document.querySelector('`+chosen+`').checked })`, &before)

	page.click(unchosen)

	var after struct{ Unchosen, Chosen bool }
	page.evaluate(`({ Unchosen: document.querySelector('`+unchosen+`').checked, Chosen: document.querySelector('`+chosen+`').checked })`, &after)
	if before.Unchosen || !before.Chosen {
		t.Fatalf("the catalog draws the group as %+v, want the first unchosen and the second chosen", before)
	}
	if !after.Unchosen || after.Chosen {
		t.Errorf("after choosing the first the group is %+v, want only the first chosen", after)
	}
}

func TestTheCatalogChoosesAndUnchoosesACheckbox(t *testing.T) {
	// scenario: Pressing a checkbox state in the catalog chooses it and pressing it again unchooses it
	t.Parallel()
	page := catalogued(t)

	control := stateOf("ui-choice", "checkbox") + " input"
	checked := `document.querySelector('` + control + `').checked`

	page.click(control)
	var chosen bool
	page.evaluate(checked, &chosen)
	page.click(control)
	var unchosen bool
	page.evaluate(checked, &unchosen)

	if !chosen || unchosen {
		t.Errorf("the checkbox was chosen after the first press: %v, and after the second: %v; want true then false", chosen, unchosen)
	}
}

func TestTheCatalogChoosesWhatThePressLandsOnInsideTheOptionsFrame(t *testing.T) {
	// scenario: Pressing the words of a choice in the catalog chooses it
	t.Parallel()
	page := catalogued(t)

	state := stateOf("ui-choice", "checkbox")
	page.click(state + " [data-part='label']")

	var chosen bool
	page.evaluate(`document.querySelector('`+state+` input').checked`, &chosen)
	if !chosen {
		t.Error("pressing the label of a choice did not choose it")
	}
}

func TestTheCatalogChoosesAChoiceWithTheKeyboard(t *testing.T) {
	// scenario: Space on a focused choice state in the catalog chooses it
	t.Parallel()
	page := catalogued(t)

	control := stateOf("ui-choice", "checkbox") + " input"
	page.evaluate(`document.querySelector('`+control+`').focus()`, nil)

	page.press(" ")

	var chosen bool
	page.evaluate(`document.querySelector('`+control+`').checked`, &chosen)
	if !chosen {
		t.Error("Space on a focused choice did not choose it")
	}
}

func TestTheCatalogDrawsAChoiceThatArrivesChosenAndOneThatCannotBeChosen(t *testing.T) {
	// scenario: A disabled choice state in the catalog takes no press and a chosen state arrives chosen
	t.Parallel()
	page := catalogued(t)

	disabled := stateOf("ui-choice", "disabled") + " input"
	page.evaluate(`document.querySelector('`+disabled+`').click()`, nil)

	var held struct {
		Disabled bool `json:"disabled"`
		Checked  bool `json:"checked"`
		Arrives  bool `json:"arrives"`
	}
	page.evaluate(`({
  disabled: document.querySelector('`+disabled+`').disabled,
  checked: document.querySelector('`+disabled+`').checked,
  arrives: document.querySelector('`+stateOf("ui-choice", "radio-chosen")+` input').defaultChecked,
})`, &held)
	if !held.Disabled || held.Checked || !held.Arrives {
		t.Errorf("the choices are %+v, want the disabled one disabled and unchosen and the chosen one chosen from the start", held)
	}
}

func TestTheCatalogDrawsAChoiceWithATagAndAMeaningAndOneWithNeither(t *testing.T) {
	// scenario: A choice state in the catalog draws its tag and its meaning only where it has them
	t.Parallel()
	page := catalogued(t)

	var held struct {
		FullTag, FullMeaning, BareTag, BareMeaning int
	}
	page.evaluate(`({
  FullTag: document.querySelectorAll('`+stateOf("ui-choice", "with-everything")+` [data-part="tag"]').length,
  FullMeaning: document.querySelectorAll('`+stateOf("ui-choice", "with-everything")+` [data-part="meaning"]').length,
  BareTag: document.querySelectorAll('`+stateOf("ui-choice", "checkbox")+` [data-part="tag"]').length,
  BareMeaning: document.querySelectorAll('`+stateOf("ui-choice", "checkbox")+` [data-part="meaning"]').length,
})`, &held)
	want := struct{ FullTag, FullMeaning, BareTag, BareMeaning int }{1, 1, 0, 0}
	if diff := cmp.Diff(want, held); diff != "" {
		t.Errorf("the tags and meanings drawn differ (-want +got):\n%s", diff)
	}
}

func TestTheCatalogCopiesWhatACopyControlHoldsExactlyAsItWasWritten(t *testing.T) {
	// scenario: Pressing a copy control state that holds text puts that text on the clipboard exactly
	t.Parallel()
	page := catalogued(t)
	page.mayCopy()
	page.clipboardHolds(somethingElseHeld)

	page.click(stateOf("ui-copy", "resting") + " ui-copy")
	page.waitFor(`navigator.clipboard.readText().then((held) => held !== ` + jsonOf(t, somethingElseHeld) + `)`)

	if diff := cmp.Diff(catalogCopiesThis, page.onTheClipboard()); diff != "" {
		t.Errorf("the clipboard does not hold the text the control holds (-want +got):\n%s", diff)
	}
}

func TestTheCatalogCopiesWhatAnOwnerOffersWhenTheControlHoldsNothing(t *testing.T) {
	// scenario: A copy control that holds no text copies what a listener above it offers
	t.Parallel()
	page := catalogued(t)
	page.mayCopy()
	page.clipboardHolds(somethingElseHeld)
	page.evaluate(`document.addEventListener('copyrequested', (event) => { event.detail.text = `+jsonOf(t, theOwnerOffersThis)+`; }); true`, nil)

	page.click(stateOf("ui-copy", "offering-nothing") + " ui-copy")
	page.waitFor(`navigator.clipboard.readText().then((held) => held !== ` + jsonOf(t, somethingElseHeld) + `)`)

	if diff := cmp.Diff(theOwnerOffersThis, page.onTheClipboard()); diff != "" {
		t.Errorf("the clipboard does not hold what the owner offered (-want +got):\n%s", diff)
	}
}

func TestTheCatalogCopiesNothingWhenNobodyOffersAnyText(t *testing.T) {
	// scenario: A copy control that holds no text and meets no owner leaves the clipboard and itself alone
	t.Parallel()
	page := catalogued(t)
	page.mayCopy()
	page.clipboardHolds(somethingElseHeld)

	page.click(stateOf("ui-copy", "offering-nothing") + " ui-copy")
	page.do(chromedp.Sleep(anIdlePressWindow))

	var confirmed bool
	page.evaluate(`document.querySelector('`+stateOf("ui-copy", "offering-nothing")+` ui-copy').dataset.copied !== undefined`, &confirmed)
	if confirmed {
		t.Error("the control confirmed a copy it never made")
	}
	if diff := cmp.Diff(somethingElseHeld, page.onTheClipboard()); diff != "" {
		t.Errorf("the clipboard changed though nothing was offered (-want +got):\n%s", diff)
	}
}

func TestTheCatalogConfirmsACopyWhereItWasPressedAndThenRests(t *testing.T) {
	// scenario: A copy control confirms the copy in its title and its live status, then returns to rest
	t.Parallel()
	page := catalogued(t)
	page.mayCopy()
	control := stateOf("ui-copy", "resting") + " ui-copy"

	page.click(control)
	page.waitFor(`document.querySelector('` + control + `').dataset.copied === 'yes'`)

	var confirmed string
	page.evaluate(`document.querySelector('`+control+`').dataset.title`, &confirmed)
	if confirmed != "Copied" {
		t.Errorf("after the copy the title says %q, want Copied", confirmed)
	}

	page.waitFor(`document.querySelector('` + control + `').dataset.copied === undefined && document.querySelector('` + control + `').dataset.title === 'Copy code'`)
}

func TestTheCatalogSaysACopyFailedWhenTheClipboardRefuses(t *testing.T) {
	// scenario: A copy control whose clipboard refuses says it could not copy
	t.Parallel()
	page := catalogued(t)
	page.evaluate(`navigator.clipboard.writeText = () => Promise.reject(new Error('refused')); true`, nil)
	control := stateOf("ui-copy", "resting") + " ui-copy"

	page.click(control)
	page.waitFor(`document.querySelector('` + control + `').dataset.failed === 'yes'`)

	var held struct {
		Title  string `json:"title"`
		Copied bool   `json:"copied"`
	}
	page.evaluate(`({ title: document.querySelector('`+control+`').dataset.title, copied: document.querySelector('`+control+`').dataset.copied !== undefined })`, &held)
	if held.Title != "Could not copy" || held.Copied {
		t.Errorf("after the refusal the control is %+v, want it to say Could not copy and not claim a copy", held)
	}
}

func TestTheCatalogShowsTheTipOfACopyControlAgainWhenItsTitleChanges(t *testing.T) {
	// scenario: A copy control tells the page when its title changes, once for the confirmation and once for the rest
	t.Parallel()
	page := catalogued(t)
	page.mayCopy()
	control := stateOf("ui-copy", "resting") + " ui-copy"
	page.evaluate(`window.__retitled = []; document.addEventListener('retitled', (event) => { window.__retitled.push(event.target.dataset.title); }); true`, nil)

	page.click(control)
	page.waitFor(`window.__retitled.length === 2`)

	var titles []string
	page.evaluate(`window.__retitled`, &titles)
	if diff := cmp.Diff([]string{"Copied", "Copy code"}, titles); diff != "" {
		t.Errorf("the titles announced differ (-want +got):\n%s", diff)
	}
}

func TestTheCatalogExpandsAnExpandControlAndAnnouncesIt(t *testing.T) {
	// scenario: Pressing an expand control state at rest sets it expanded and announces expanded once
	t.Parallel()
	page := catalogued(t)

	state := stateOf("ui-expand", "resting")
	page.evaluate(`window.__expanding = []; document.addEventListener('expanded', () => { window.__expanding.push('expanded'); }); document.addEventListener('collapsed', () => { window.__expanding.push('collapsed'); }); true`, nil)

	page.click(state + " ui-expand")

	var held struct {
		Expanded string   `json:"expanded"`
		Heard    []string `json:"heard"`
	}
	page.evaluate(`({ expanded: String(document.querySelector('`+state+` ui-expand').expanded), heard: window.__expanding })`, &held)
	if held.Expanded != "true" || cmp.Diff([]string{"expanded"}, held.Heard) != "" {
		t.Errorf("after the press the control is expanded=%q and the page heard %v, want true and one expanded", held.Expanded, held.Heard)
	}
}

func TestTheCatalogCollapsesAnExpandedExpandControlWhenItIsPressedAgain(t *testing.T) {
	// scenario: Pressing an expand control state that arrives expanded collapses it and announces collapsed once
	t.Parallel()
	page := catalogued(t)

	state := stateOf("ui-expand", "expanded")
	var arrives string
	page.evaluate(`String(document.querySelector('`+state+` ui-expand').expanded)`, &arrives)
	page.evaluate(`window.__expanding = []; document.addEventListener('collapsed', () => { window.__expanding.push('collapsed'); }); true`, nil)

	page.click(state + " ui-expand")

	var held struct {
		Expanded string   `json:"expanded"`
		Heard    []string `json:"heard"`
	}
	page.evaluate(`({ expanded: String(document.querySelector('`+state+` ui-expand').expanded), heard: window.__expanding })`, &held)
	if arrives != "true" {
		t.Errorf("the expanded state arrives expanded=%q, want true", arrives)
	}
	if held.Expanded != "false" || cmp.Diff([]string{"collapsed"}, held.Heard) != "" {
		t.Errorf("after the press the control is expanded=%q and the page heard %v, want false and one collapsed", held.Expanded, held.Heard)
	}
}

func TestTheCatalogExpandsAnExpandControlWithTheKeyboard(t *testing.T) {
	// scenario: Enter on a focused expand control state at rest expands it
	t.Parallel()
	page := catalogued(t)

	state := stateOf("ui-expand", "resting")
	page.evaluate(`document.querySelector('`+state+` ui-expand').focus()`, nil)

	page.press(kb.Enter)

	var expanded string
	page.evaluate(`String(document.querySelector('`+state+` ui-expand').expanded)`, &expanded)
	if expanded != "true" {
		t.Errorf("Enter on the focused control left it expanded=%q, want true", expanded)
	}
}

func TestTheCatalogReturnsTheCursorToAnExpandControlItsOwnerCollapses(t *testing.T) {
	// scenario: Collapsing an expanded expand control by its method returns the cursor to it
	t.Parallel()
	page := catalogued(t)

	state := stateOf("ui-expand", "expanded")
	page.evaluate(`document.querySelector('`+state+` ui-expand').collapse(); true`, nil)

	var held struct {
		Expanded string `json:"expanded"`
		Focused  bool   `json:"focused"`
	}
	page.evaluate(`({ expanded: String(document.querySelector('`+state+` ui-expand').expanded), focused: document.activeElement === document.querySelector('`+state+` ui-expand') })`, &held)
	if held.Expanded != "false" || !held.Focused {
		t.Errorf("after collapse() the control is %+v, want aria-expanded=false and the cursor on it", held)
	}
}

func TestTheCatalogChoosesAChoiceThroughItsOwnMethodAndLeavesAChosenOneChosen(t *testing.T) {
	// scenario: Calling choose() twice on a checkbox state leaves it chosen, and on a chosen radio changes nothing
	t.Parallel()
	page := catalogued(t)

	unchosen := stateOf("ui-choice", "checkbox")
	chosen := stateOf("ui-choice", "radio-chosen")
	page.evaluate(`window.__changes = 0; document.addEventListener('change', () => { window.__changes += 1; }); true`, nil)

	page.evaluate(`document.querySelector('`+unchosen+` ui-choice').choose(); document.querySelector('`+unchosen+` ui-choice').choose(); true`, nil)
	page.evaluate(`document.querySelector('`+chosen+` ui-choice').choose(); true`, nil)

	var held struct {
		Unchosen bool `json:"unchosen"`
		Chosen   bool `json:"chosen"`
		Changes  int  `json:"changes"`
		Reads    bool `json:"reads"`
	}
	page.evaluate(`({
  unchosen: document.querySelector('`+unchosen+` input').checked,
  chosen: document.querySelector('`+chosen+` input').checked,
  changes: window.__changes,
  reads: document.querySelector('`+unchosen+` ui-choice').chosen && document.querySelector('`+chosen+` ui-choice').chosen,
})`, &held)
	if !held.Unchosen || !held.Chosen || held.Changes != 1 || !held.Reads {
		t.Errorf("after choose() on both the choices are %+v, want both chosen, one change announced and chosen to read true", held)
	}
}
