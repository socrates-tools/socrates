package e2e_test

import (
	"bytes"
	"fmt"
	"image/png"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/chromedp/kb"
	"github.com/google/go-cmp/cmp"
)

const (
	singleChoiceKey          = "single-choice"
	singleChoiceChosenKey    = "single-choice-with-one-chosen"
	singleChoiceNoneKey      = "single-choice-with-none-of-these"
	singleChoiceRecommendKey = "single-choice-with-a-recommendation"
	multipleChoiceKey        = "multiple-choice"
	disabledChoicesKey       = "disabled"
)

func groupOf(key string) string {
	return stateOf("ui-choice-group", key) + " ui-choice-group"
}

func choiceOf(group string, at int) string {
	return fmt.Sprintf("%s ui-choice:nth-of-type(%d)", group, at+1)
}

func (r *reader) checkedIn(group string) []bool {
	r.t.Helper()

	var checked []bool
	r.evaluate(`Array.from(document.querySelectorAll(`+jsonOf(r.t, group+" ui-choice")+`)).map((held) => held.checked)`, &checked)
	return checked
}

func (r *reader) focusedChoiceIn(group string) int {
	r.t.Helper()

	var at int
	r.evaluate(`Array.from(document.querySelectorAll(`+jsonOf(r.t, group+" ui-choice")+`)).indexOf(document.activeElement)`, &at)
	return at
}

func (r *reader) focus(selector string) {
	r.t.Helper()
	r.evaluate(`document.querySelector(`+jsonOf(r.t, selector)+`).focus(); true`, nil)
}

func TestTheCatalogUnchoosesTheChosenChoiceWhenAnotherOfASingleChoiceGroupIsChosen(t *testing.T) {
	// scenario: Choosing one choice of a single-choice group unchecks the one checked before
	t.Parallel()
	page := catalogued(t)
	group := groupOf(singleChoiceChosenKey)

	if diff := cmp.Diff([]bool{false, true, false}, page.checkedIn(group)); diff != "" {
		t.Fatalf("the catalog draws the group differently (-want +got):\n%s", diff)
	}

	page.click(choiceOf(group, 0))
	if diff := cmp.Diff([]bool{true, false, false}, page.checkedIn(group)); diff != "" {
		t.Errorf("after choosing the first the group differs (-want +got):\n%s", diff)
	}

	page.click(choiceOf(group, 2))
	if diff := cmp.Diff([]bool{false, false, true}, page.checkedIn(group)); diff != "" {
		t.Errorf("after choosing the third the group differs (-want +got):\n%s", diff)
	}
}

func TestTheCatalogUnchoosesAnOptionWhenNoneOfTheseIsChosenAndTheOtherWayRound(t *testing.T) {
	// scenario: Choosing none of these unchecks the option checked before, and choosing an option unchecks none of these
	t.Parallel()
	page := catalogued(t)
	group := groupOf(singleChoiceNoneKey)

	page.click(choiceOf(group, 0))
	page.click(choiceOf(group, 3))
	if diff := cmp.Diff([]bool{false, false, false, true}, page.checkedIn(group)); diff != "" {
		t.Errorf("after choosing none of these the group differs (-want +got):\n%s", diff)
	}

	page.click(choiceOf(group, 1))
	if diff := cmp.Diff([]bool{false, true, false, false}, page.checkedIn(group)); diff != "" {
		t.Errorf("after choosing an option again the group differs (-want +got):\n%s", diff)
	}
}

func TestTheCatalogKeepsOneChosenWhenAChoiceIsCheckedThroughItsProperty(t *testing.T) {
	// scenario: Setting checked on a choice of a single-choice group unchecks the one checked before
	t.Parallel()
	page := catalogued(t)
	group := groupOf(singleChoiceChosenKey)

	page.evaluate(`document.querySelector(`+jsonOf(t, choiceOf(group, 2))+`).checked = true; true`, nil)
	page.waitFor(`document.querySelectorAll(` + jsonOf(t, group+" ui-choice[checked]") + `).length === 1`)

	if diff := cmp.Diff([]bool{false, false, true}, page.checkedIn(group)); diff != "" {
		t.Errorf("after setting the third's checked the group differs (-want +got):\n%s", diff)
	}
}

func TestTheCatalogMovesTheChosenChoiceAndTheCursorOnTheArrowKeys(t *testing.T) {
	// scenario: The arrow keys move the checked choice and the focus within a single-choice group
	t.Parallel()
	page := catalogued(t)
	group := groupOf(singleChoiceChosenKey)
	page.focus(choiceOf(group, 1))

	for _, step := range []struct {
		key  string
		name string
		at   int
	}{
		{kb.ArrowDown, "down from the second", 2},
		{kb.ArrowDown, "down from the last, round to the first", 0},
		{kb.ArrowUp, "up from the first, round to the last", 2},
		{kb.ArrowLeft, "left", 1},
		{kb.ArrowRight, "right", 2},
	} {
		page.press(step.key)

		want := []bool{false, false, false}
		want[step.at] = true
		if diff := cmp.Diff(want, page.checkedIn(group)); diff != "" {
			t.Errorf("after %s the chosen choice differs (-want +got):\n%s", step.name, diff)
		}
		if got := page.focusedChoiceIn(group); got != step.at {
			t.Errorf("after %s the cursor is on choice %d, want %d", step.name, got, step.at)
		}
	}
}

func TestTheCatalogMovesOntoNoneOfTheseOnTheArrowKeys(t *testing.T) {
	// scenario: The arrow keys move between the options and none of these
	t.Parallel()
	page := catalogued(t)
	group := groupOf(singleChoiceNoneKey)
	page.click(choiceOf(group, 2))

	page.press(kb.ArrowDown)

	if diff := cmp.Diff([]bool{false, false, false, true}, page.checkedIn(group)); diff != "" {
		t.Errorf("after the arrow key the group differs (-want +got):\n%s", diff)
	}
	if got := page.focusedChoiceIn(group); got != 3 {
		t.Errorf("the cursor is on choice %d, want none of these", got)
	}
}

func (r *reader) aroundWithButtons(group string) {
	r.t.Helper()
	r.evaluate(`(() => {
  const held = document.querySelector(`+jsonOf(r.t, group)+`);
  const before = document.createElement('button');
  before.id = 'before-the-group';
  const after = document.createElement('button');
  after.id = 'after-the-group';
  held.before(before);
  held.after(after);
  return true;
})()`, nil)
}

func (r *reader) onTheButton(id string) bool {
	r.t.Helper()

	var held bool
	r.evaluate(`document.activeElement === document.getElementById(`+jsonOf(r.t, id)+`)`, &held)
	return held
}

func TestTheCatalogEntersASingleChoiceGroupAtItsChosenChoiceAndLeavesItInOnePress(t *testing.T) {
	// scenario: Tab enters a single-choice group at its checked choice and leaves it in one press
	t.Parallel()
	page := catalogued(t)
	group := groupOf(singleChoiceChosenKey)
	page.aroundWithButtons(group)
	page.focus("#before-the-group")

	page.press(kb.Tab)
	if got := page.focusedChoiceIn(group); got != 1 {
		t.Fatalf("Tab entered the group at choice %d, want its checked second", got)
	}

	page.press(kb.Tab)
	if !page.onTheButton("after-the-group") {
		t.Error("one more Tab did not leave the group")
	}
}

func TestTheCatalogEntersASingleChoiceGroupHoldingNothingChosenAtItsFirstChoice(t *testing.T) {
	// scenario: Tab enters a single-choice group with nothing checked at its first choice
	t.Parallel()
	page := catalogued(t)
	group := groupOf(singleChoiceKey)
	page.aroundWithButtons(group)
	page.focus("#before-the-group")

	page.press(kb.Tab)

	if got := page.focusedChoiceIn(group); got != 0 {
		t.Errorf("Tab entered the group at choice %d, want the first", got)
	}
}

func TestTheCatalogReachesEveryChoiceOfAMultipleChoiceGroupWithTab(t *testing.T) {
	// scenario: Tab visits each choice of a multiple-choice group
	t.Parallel()
	page := catalogued(t)
	group := groupOf(multipleChoiceKey)
	page.aroundWithButtons(group)
	page.focus("#before-the-group")

	var visited []int
	for range 3 {
		page.press(kb.Tab)
		visited = append(visited, page.focusedChoiceIn(group))
	}

	if diff := cmp.Diff([]int{0, 1, 2}, visited); diff != "" {
		t.Errorf("Tab visited the choices in an order that differs (-want +got):\n%s", diff)
	}
}

func TestTheCatalogChecksSeveralChoicesOfAMultipleChoiceGroupAtOnce(t *testing.T) {
	// scenario: A multiple-choice group checks several choices at once and unchecks one pressed again
	t.Parallel()
	page := catalogued(t)
	group := groupOf(multipleChoiceKey)

	page.click(choiceOf(group, 1))
	page.click(choiceOf(group, 2))
	if diff := cmp.Diff([]bool{true, true, true}, page.checkedIn(group)); diff != "" {
		t.Errorf("after choosing the second and the third the group differs (-want +got):\n%s", diff)
	}

	page.click(choiceOf(group, 0))
	if diff := cmp.Diff([]bool{false, true, true}, page.checkedIn(group)); diff != "" {
		t.Errorf("after pressing the first again the group differs (-want +got):\n%s", diff)
	}
}

func TestTheCatalogChecksTheFocusedChoiceWithTheSpaceBar(t *testing.T) {
	// scenario: Space on a focused choice checks it
	t.Parallel()
	page := catalogued(t)
	group := groupOf(singleChoiceKey)
	page.focus(choiceOf(group, 2))

	page.press(" ")

	if diff := cmp.Diff([]bool{false, false, true}, page.checkedIn(group)); diff != "" {
		t.Errorf("after Space the group differs (-want +got):\n%s", diff)
	}
}

func TestTheCatalogAnnouncesAGroupAndItsChoices(t *testing.T) {
	// scenario: A single-choice group is a radio group of radios saying whether each is checked, and a multiple-choice group's choices are checkboxes
	t.Parallel()
	page := catalogued(t)

	single := groupOf(singleChoiceChosenKey)
	if got := page.attributeOf(single, "role"); got != "radiogroup" {
		t.Errorf("a single-choice group is announced as %q, want a radio group", got)
	}
	for at, want := range []string{"false", "true", "false"} {
		choice := choiceOf(single, at)
		if got := page.attributeOf(choice, "role"); got != "radio" {
			t.Errorf("choice %d of a single-choice group is announced as %q, want a radio", at, got)
		}
		if got := page.attributeOf(choice, "aria-checked"); got != want {
			t.Errorf("choice %d of a single-choice group says aria-checked %q, want %q", at, got, want)
		}
	}

	multiple := groupOf(multipleChoiceKey)
	if got := page.attributeOf(multiple, "role"); got != "group" {
		t.Errorf("a multiple-choice group is announced as %q, want a group", got)
	}
	for at, want := range []string{"true", "false", "false"} {
		choice := choiceOf(multiple, at)
		if got := page.attributeOf(choice, "role"); got != "checkbox" {
			t.Errorf("choice %d of a multiple-choice group is announced as %q, want a checkbox", at, got)
		}
		if got := page.attributeOf(choice, "aria-checked"); got != want {
			t.Errorf("choice %d of a multiple-choice group says aria-checked %q, want %q", at, got, want)
		}
	}

	page.click(choiceOf(single, 2))
	if got := page.attributeOf(choiceOf(single, 2), "aria-checked"); got != "true" {
		t.Errorf("a choice just chosen says aria-checked %q, want true", got)
	}
	if got := page.attributeOf(choiceOf(single, 1), "aria-checked"); got != "false" {
		t.Errorf("the choice it replaced says aria-checked %q, want false", got)
	}
}

func TestTheCatalogDrawsARecommendedChoiceChosenFromTheStartWithItsTagAndMeaning(t *testing.T) {
	// scenario: A recommended choice arrives chosen and carries its tag and its meaning only where it has them
	t.Parallel()
	page := catalogued(t)
	group := groupOf(singleChoiceRecommendKey)

	if diff := cmp.Diff([]bool{true, false, false}, page.checkedIn(group)); diff != "" {
		t.Errorf("the recommendation does not arrive chosen (-want +got):\n%s", diff)
	}
	if got := page.countOf(choiceOf(group, 0) + ` [slot="tag"]`); got != 1 {
		t.Errorf("the recommended choice carries %d tags, want 1", got)
	}
	if got := page.countOf(choiceOf(group, 0) + ` [slot="meaning"]`); got != 1 {
		t.Errorf("the recommended choice carries %d meanings, want 1", got)
	}
	if got := page.countOf(choiceOf(group, 1) + ` [slot="tag"], ` + choiceOf(group, 1) + ` [slot="meaning"]`); got != 0 {
		t.Errorf("a choice nobody recommended carries %d tags and meanings, want none", got)
	}
}

func TestTheCatalogTakesNoPressOnADisabledGroup(t *testing.T) {
	// scenario: A group drawn disabled takes no press and no cursor and says so to a screen reader
	t.Parallel()
	page := catalogued(t)
	group := groupOf(disabledChoicesKey)

	page.evaluate(`Array.from(document.querySelectorAll(`+jsonOf(t, group+" ui-choice")+`)).forEach((held) => held.click()); true`, nil)
	page.focus(choiceOf(group, 1))

	if diff := cmp.Diff([]bool{true, false, false}, page.checkedIn(group)); diff != "" {
		t.Errorf("a disabled group changed (-want +got):\n%s", diff)
	}
	for at := range 3 {
		choice := choiceOf(group, at)
		if got := page.attributeOf(choice, "aria-disabled"); got != "true" {
			t.Errorf("disabled choice %d says aria-disabled %q, want true", at, got)
		}
		if page.countOf(choice+"[tabindex]") != 0 {
			t.Errorf("disabled choice %d can take the cursor", at)
		}
	}
	if got := page.focusedChoiceIn(group); got != -1 {
		t.Errorf("the cursor reached choice %d of a disabled group", got)
	}
}

func TestTheCatalogPaintsEachChoiceStateInItsOwnTokens(t *testing.T) {
	// scenario: A choice state paints its border, its wash and its ring in the tokens that state is drawn with
	t.Parallel()
	page := catalogued(t)

	for _, state := range []struct {
		name, key, border, wash string
		style                   string
		opacity                 string
	}{
		{"unchosen", "radio", "--rule", "--surface-card", "solid", "1"},
		{"chosen", "radio-chosen", "--accent", "--accent-selected", "solid", "1"},
		{"dashed", "dashed", "--rule-control", "", "dashed", "1"},
		{"dashed and chosen", "dashed-chosen", "--fg-2", "--surface-muted", "solid", "1"},
		{"disabled", "disabled", "--rule", "--surface-card", "solid", "0.55"},
	} {
		t.Run(state.name, func(t *testing.T) {
			choice := stateOf("ui-choice", state.key) + " ui-choice"

			painted := page.paintedOf(choice, "border-top-color", "border-top-style", "background-color", "opacity", "border-top-width", "border-top-left-radius")

			if want := page.pageColourOf(state.border); painted["border-top-color"] != want {
				t.Errorf("the %s choice's border is %s, want %s", state.name, painted["border-top-color"], want)
			}
			if painted["border-top-style"] != state.style {
				t.Errorf("the %s choice's border is %s, want %s", state.name, painted["border-top-style"], state.style)
			}
			if state.wash != "" {
				if want := page.pageColourOf(state.wash); painted["background-color"] != want {
					t.Errorf("the %s choice is washed %s, want %s", state.name, painted["background-color"], want)
				}
			}
			if painted["opacity"] != state.opacity {
				t.Errorf("the %s choice is drawn at opacity %s, want %s", state.name, painted["opacity"], state.opacity)
			}
			if painted["border-top-width"] != "1px" {
				t.Errorf("the %s choice's border is %s wide, want 1px", state.name, painted["border-top-width"])
			}
			if got, want := page.pixelsOf(choice, "border-top-left-radius"), page.pixelsOfToken("--radius-md"); got != want {
				t.Errorf("the %s choice's corners are %vpx, want %vpx", state.name, got, want)
			}
		})
	}
}

func (r *reader) shadowOfToken(token string) string {
	r.t.Helper()

	var shadow string
	r.evaluate(`(() => {
  const probe = document.createElement('i');
  probe.style.boxShadow = 'var(`+token+`)';
  document.body.append(probe);
  const shadow = getComputedStyle(probe).boxShadow;
  probe.remove();
  return shadow;
})()`, &shadow)
	return shadow
}

func (r *reader) pixelBesideTheIndicatorOf(choice string) [3]int {
	r.t.Helper()

	r.evaluate(`document.querySelector(`+jsonOf(r.t, choice)+`).scrollIntoView({ block: 'center' }); new Promise((settled) => requestAnimationFrame(() => requestAnimationFrame(settled)))`, nil, awaitingThePromise)
	time.Sleep(500 * time.Millisecond)
	var at struct{ X, Y float64 }
	r.evaluate(`(() => {
  const held = document.querySelector(`+jsonOf(r.t, choice)+`);
  const box = held.getBoundingClientRect();
  const style = getComputedStyle(held);
  const edge = box.left + parseFloat(style.borderLeftWidth) + parseFloat(style.paddingLeft);
  const top = box.top + parseFloat(style.borderTopWidth) + parseFloat(style.paddingTop);
  return { X: edge - 1.5, Y: top + 2 + 9 };
})()`, &at)
	shot, err := png.Decode(bytes.NewReader(r.screenshot()))
	if err != nil {
		r.t.Fatalf("reading what the page painted: %v", err)
	}
	red, green, blue, _ := shot.At(int(at.X), int(at.Y)).RGBA()
	return [3]int{int(red >> 8), int(green >> 8), int(blue >> 8)}
}

func (r *reader) focusRingColour() [3]int {
	r.t.Helper()

	var ring [3]int
	colour := r.pageColourOf("--accent-ring")
	if _, err := fmt.Sscanf(colour, "rgb(%d, %d, %d)", &ring[0], &ring[1], &ring[2]); err != nil {
		r.t.Fatalf("reading the focus ring colour %q: %v", colour, err)
	}
	return ring
}

func TestTheCatalogDrawsTheFocusRingAroundTheIndicatorOfAChosenChoiceReachedByKeyboardAndNoOutline(t *testing.T) {
	// scenario: A chosen choice reached by Tab paints the page's focus ring around its indicator beside the chosen ring, and no outline
	t.Parallel()
	page := catalogued(t)
	ring := page.focusRingColour()

	chosen := groupOf(singleChoiceChosenKey)
	page.aroundWithButtons(chosen)
	page.focus("#before-the-group")
	idle := page.pixelBesideTheIndicatorOf(choiceOf(chosen, 1))
	page.press(kb.Tab)

	if got := page.pixelBesideTheIndicatorOf(choiceOf(chosen, 1)); got != ring {
		t.Errorf("beside the indicator of the focused chosen choice the page paints %v, want the focus ring %v", got, ring)
	}
	if idle == ring {
		t.Errorf("beside the indicator of a choice nobody focused the page paints the focus ring %v", ring)
	}
	painted := page.paintedOf(choiceOf(chosen, 1), "outline-style", "box-shadow")
	if painted["outline-style"] != "none" {
		t.Errorf("the focused chosen choice is drawn with the outline %q, want none", painted["outline-style"])
	}
	if want := page.shadowOfToken("--ring-chosen"); !strings.HasSuffix(painted["box-shadow"], want) {
		t.Errorf("the focused chosen choice keeps the shadow %q, want the chosen ring %q", painted["box-shadow"], want)
	}
}

func TestTheCatalogDrawsTheFocusRingAroundTheIndicatorOfAnUnchosenChoiceReachedByKeyboardAndNothingAroundIt(t *testing.T) {
	// scenario: An unchosen choice reached by Tab paints the page's focus ring around its indicator and no outline or shadow around itself
	t.Parallel()
	page := catalogued(t)
	ring := page.focusRingColour()

	multiple := groupOf(multipleChoiceKey)
	page.aroundWithButtons(multiple)
	page.focus("#before-the-group")
	page.press(kb.Tab)
	page.press(kb.Tab)

	if got := page.focusedChoiceIn(multiple); got != 1 {
		t.Fatalf("Tab reached choice %d, want the second", got)
	}
	if got := page.pixelBesideTheIndicatorOf(choiceOf(multiple, 1)); got != ring {
		t.Errorf("beside the indicator of the focused unchosen choice the page paints %v, want the focus ring %v", got, ring)
	}
	painted := page.paintedOf(choiceOf(multiple, 1), "outline-style", "box-shadow")
	nothingPainted := strings.NewReplacer("rgba(0, 0, 0, 0) 0px 0px 0px 0px", "", ",", "", " ", "", "none", "")
	if painted["outline-style"] != "none" || nothingPainted.Replace(painted["box-shadow"]) != "" {
		t.Errorf("the focused unchosen choice is drawn with the outline %q and the shadow %q, want neither", painted["outline-style"], painted["box-shadow"])
	}
}

func TestTheCatalogAnnouncesChoicechangedOnlyWhenAChoiceIsCheckedOrUnchecked(t *testing.T) {
	// scenario: A choice announces choicechanged when its checked changes and not when it is disabled or enabled
	t.Parallel()
	page := catalogued(t)
	group := groupOf(singleChoiceChosenKey)
	page.evaluate(`window.__changed = 0; document.querySelector(`+jsonOf(t, group)+`).addEventListener('choicechanged', () => { window.__changed += 1; }); true`, nil)
	choice := `document.querySelector(` + jsonOf(t, choiceOf(group, 0)) + `)`
	settled := `new Promise((done) => requestAnimationFrame(() => requestAnimationFrame(done)))`

	page.evaluate(choice+`.disabled = true; `+settled, nil, awaitingThePromise)
	page.evaluate(choice+`.disabled = false; `+settled, nil, awaitingThePromise)
	var whileOnlyDisabledChanged int
	page.evaluate(`window.__changed`, &whileOnlyDisabledChanged)

	page.evaluate(choice+`.checked = true; `+settled, nil, awaitingThePromise)
	var afterChecking int
	page.evaluate(`window.__changed`, &afterChecking)

	if whileOnlyDisabledChanged != 0 {
		t.Errorf("disabling and enabling a choice announced choicechanged %d times, want none", whileOnlyDisabledChanged)
	}
	if afterChecking == 0 {
		t.Error("checking a choice announced no choicechanged")
	}
}
