package e2e_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func threeOptions(mode string) sketch {
	decision := aCacheDecision(mode)
	decision["options"] = append(decision["options"].([]sketch), sketch{"label": "shard the map", "meaning": "contention matters", "command": "patch c-1 --shard"})
	return decision
}

func TestADocumentRendersWithoutTheCallerWritingMarkup(t *testing.T) {
	t.Parallel()
	_, page := read(t, written(t, oneCacheDocument("single-choice")))

	var card struct {
		Cards    int    `json:"cards"`
		ID       string `json:"id"`
		Text     string `json:"text"`
		Groups   int    `json:"groups"`
		Multiple bool   `json:"multiple"`
		Marks    map[string]*struct {
			Data  string `json:"data"`
			Label string `json:"label"`
		} `json:"marks"`
	}
	page.evaluate(`(() => {
  const cards = Array.from(document.querySelectorAll('.decision'));
  const held = cards[0];
  const marks = {};
  for (const kind of ['severity', 'audience']) {
    const mark = held.querySelector('.mark.' + kind);
    marks[kind] = mark ? { data: mark.dataset[kind], label: mark.getAttribute('aria-label') || '' } : null;
  }
  const groups = held.querySelectorAll('ui-choice-group');
  return { cards: cards.length, id: held.dataset.decision, text: held.textContent, groups: groups.length, multiple: groups[0] ? groups[0].multiple : false, marks };
})()`, &card)

	if card.Cards != 1 || card.ID != "c-1" {
		t.Fatalf("one decision drew %d cards, the first naming %q, want one naming c-1", card.Cards, card.ID)
	}
	for _, expected := range []string{
		"should the cache be read before the lock is taken",
		"two callers arriving together read the same stale entry",
		"take the lock first",
		"leave it",
		"internal/cache/read.go:41",
	} {
		if !strings.Contains(card.Text, expected) {
			t.Errorf("the card does not carry %q: %q", expected, card.Text)
		}
	}
	for kind, value := range map[string]string{"severity": "error", "audience": "a person"} {
		mark := card.Marks[kind]
		if mark == nil {
			t.Errorf("the card carries no %s mark", kind)
			continue
		}
		if mark.Data != value {
			t.Errorf("the card's %s mark carries %q, and the document stated %q", kind, mark.Data, value)
		}
		if !strings.Contains(strings.ToLower(mark.Label), value) {
			t.Errorf("the card's %s mark is announced as %q, and a person who cannot see the icon hears %q in it", kind, mark.Label, value)
		}
	}
	if card.Groups != 1 {
		t.Fatalf("a single-choice decision draws %d groups of choices, want one", card.Groups)
	}
	if card.Multiple {
		t.Error("a single-choice decision draws a group that lets several be chosen at once")
	}
}

func TestAMultipleChoiceDecisionRendersOneGroupLettingSeveralBeChosenAtOnce(t *testing.T) {
	t.Parallel()
	_, page := read(t, written(t, oneCacheDocument("multiple-choice")))

	var groups []bool
	page.evaluate(`Array.from(document.querySelectorAll('ui-choice-group'), (held) => held.multiple)`, &groups)
	if len(groups) != 1 {
		t.Fatalf("a multiple-choice decision draws %d groups of choices, want one", len(groups))
	}
	if !groups[0] {
		t.Error("a multiple-choice decision draws a group that keeps one choice chosen")
	}
}

func TestEveryOptionPostsBackThePositionTheDocumentGaveIt(t *testing.T) {
	t.Parallel()
	decision := threeOptions("multiple-choice")
	decision["options"].([]sketch)[2]["recommended"] = true
	_, page := read(t, written(t, aDocument("one thing to settle", nil, decision)))

	var positions map[string]string
	page.evaluate(`Object.fromEntries(Array.from(document.querySelectorAll('.decision .option'), (held) => [held.querySelector('.option-label').textContent.trim(), held.value]))`, &positions)

	want := map[string]string{"take the lock first": "0", "leave it": "1", "shard the map": "2"}
	if diff := cmp.Diff(want, positions); diff != "" {
		t.Errorf("what each option posts back differs from the position the document gave it (-want +got):\n%s", diff)
	}
}

func TestADocumentThatMarksNothingIsUnaffected(t *testing.T) {
	t.Parallel()
	_, page := read(t, written(t, aDocument("one thing to settle", nil, threeOptions("multiple-choice"))))

	var drawn struct {
		Labels      []string `json:"labels"`
		Recommended int      `json:"recommended"`
		Marks       int      `json:"marks"`
		Chosen      []string `json:"chosen"`
	}
	page.evaluate(`({
  labels: Array.from(document.querySelectorAll('.option .option-label'), (held) => held.textContent.trim()),
  recommended: document.querySelectorAll('.option[data-recommended]').length,
  marks: document.querySelectorAll('.option .option-mark').length,
  chosen: Array.from(document.querySelectorAll('ui-choice'), (held) => held.checked ? held.value : '').filter(Boolean),
})`, &drawn)

	if want := []string{"take the lock first", "leave it", "shard the map"}; !slices.Equal(want, drawn.Labels) {
		t.Errorf("a document that marks nothing draws its options as %v, want the order it listed them %v", drawn.Labels, want)
	}
	if drawn.Recommended != 0 {
		t.Errorf("%d options are marked, and the document marked nothing", drawn.Recommended)
	}
	if drawn.Marks != 0 {
		t.Errorf("the page carries %d marks, and the document marked nothing", drawn.Marks)
	}
	if len(drawn.Chosen) != 0 {
		t.Errorf("the options %v arrive chosen, and the document marked nothing", drawn.Chosen)
	}
}

func TestASingleChoiceDecisionOffersTakingNoneOfThem(t *testing.T) {
	t.Parallel()
	_, page := read(t, written(t, oneCacheDocument("single-choice")))

	var none *struct {
		Option      bool   `json:"option"`
		Choice      bool   `json:"choice"`
		Variant     string `json:"variant"`
		Interaction string `json:"interaction"`
		Options     int    `json:"options"`
		SameGroup   bool   `json:"sameGroup"`
	}
	page.evaluate(`(() => {
  const held = document.querySelector('.decision[data-decision="c-1"] ui-choice.none');
  if (!held) return null;
  const options = Array.from(document.querySelectorAll('.decision[data-decision="c-1"] .option'));
  return {
    option: held.classList.contains('option'),
    choice: held.classList.contains('choice'),
    variant: held.variant,
    interaction: held.dataset.interaction || '',
    options: options.length,
    sameGroup: options.every((option) => option.parentElement === held.parentElement),
  };
})()`, &none)

	if none == nil {
		t.Fatal("a single-choice decision draws no control for taking none of the options")
	}
	if none.Option {
		t.Error("the control for taking none of the options is counted as an option")
	}
	if none.Choice {
		t.Error("the control for taking none of the options is counted as a choice")
	}
	if none.Variant != "dashed" {
		t.Errorf("the control for taking none of the options is drawn %q, want dashed", none.Variant)
	}
	if none.Interaction != "answering" {
		t.Errorf("the control for taking none of the options carries data-interaction %q, want answering", none.Interaction)
	}
	if none.Options != 2 {
		t.Errorf("a single-choice decision offering two options and a way to take none of them counts %d options", none.Options)
	}
	if !none.SameGroup {
		t.Error("the control for taking none of the options is not in the group its options are in")
	}
}

func TestADecisionWithNoOptionsRendersAsACardThatAsksNothing(t *testing.T) {
	t.Parallel()
	_, page := read(t, written(t, oneCacheDocumentStatingNothingToChoose()))

	var card struct {
		Cards    int    `json:"cards"`
		Text     string `json:"text"`
		Options  int    `json:"options"`
		Controls int    `json:"controls"`
	}
	page.evaluate(`(() => {
  const cards = Array.from(document.querySelectorAll('.decision'));
  return {
    cards: cards.length,
    text: cards[0].textContent,
    options: document.querySelectorAll('.option').length,
    controls: cards[0].querySelectorAll('input, select, fieldset, legend').length,
  };
})()`, &card)

	if card.Cards != 1 {
		t.Fatalf("one decision that asks nothing drew %d cards", card.Cards)
	}
	for _, stated := range []string{"should the cache be read before the lock is taken", "two callers arriving together read the same stale entry"} {
		if !strings.Contains(card.Text, stated) {
			t.Errorf("the card does not carry %q: %q", stated, card.Text)
		}
	}
	if card.Options != 0 {
		t.Errorf("a decision that asks nothing drew %d options", card.Options)
	}
	if card.Controls != 0 {
		t.Errorf("a decision that asks nothing draws %d controls, and nothing on it is a choice", card.Controls)
	}
	for _, said := range []string{"Choose one", "Choose as many as apply"} {
		if strings.Contains(card.Text, said) {
			t.Errorf("a decision that asks nothing says %q, and there is nothing on it to choose among", said)
		}
	}
}

func oneCacheDocumentStatingNothingToChoose() sketch {
	decision := aCacheDecision("")
	delete(decision, "options")
	delete(decision, "mode")
	return aDocument("one thing to settle", nil, decision)
}

const whatEveryCardCarries = `(() => {
  const cards = Array.from(document.querySelectorAll('.decision'));
  return cards.map((card) => {
    const fields = Array.from(card.querySelectorAll('ui-field.words-field'));
    const caption = fields[0] ? fields[0].querySelector('[slot="caption"]') : null;
    return { fields: fields.length, rows: fields[0] ? fields[0].rows : 0, caption: caption ? caption.textContent.trim() : '' };
  });
})()`

type cardCarrying struct {
	Fields  int    `json:"fields"`
	Rows    int    `json:"rows"`
	Caption string `json:"caption"`
}

func (r *revisions) cardsCarrying() []cardCarrying {
	r.t.Helper()

	var carried []cardCarrying
	r.page.evaluate(whatEveryCardCarries, &carried)
	return carried
}

func withoutField(decision sketch, emptied string) sketch {
	held := sketch{}
	for key, value := range decision {
		if key != emptied {
			held[key] = value
		}
	}
	if emptied == "options" {
		delete(held, "mode")
	}
	return held
}

func TestADecisionWithNoOptionsStillCarriesTheField(t *testing.T) {
	t.Parallel()
	asking := oneCacheDocumentStatingNothingToChoose()["decisions"].([]sketch)[0]
	shown := showing(t, aDocument("one thing to settle", nil, asking))

	for _, emptied := range []string{"", "severity", "audience", "location", "meaning"} {
		shown.show(aDocument("one thing to settle", nil, withoutField(asking, emptied)))

		carried := shown.cardsCarrying()
		if len(carried) != 1 || carried[0].Fields != 1 {
			t.Fatalf("emptying %q removed the field a person writes in from a decision that asks nothing: %+v", emptied, carried)
		}
		if carried[0].Rows != 2 {
			t.Errorf("the field a person writes in is %d rows, and its size is the page's", carried[0].Rows)
		}
		if carried[0].Caption != "In your own words" {
			t.Errorf("the field a person writes in is labelled %q, and its label is the page's", carried[0].Caption)
		}
	}
}

func TestTheFreeTextFieldCannotBeSuppressed(t *testing.T) {
	t.Parallel()
	shown := showing(t, oneCacheDocument("single-choice"))

	for _, mode := range []string{"single-choice", "multiple-choice"} {
		for _, emptied := range []string{"", "severity", "audience", "location", "meaning", "options"} {
			shown.show(aDocument("one thing to settle", nil, withoutField(aCacheDecision(mode), emptied)))

			carried := shown.cardsCarrying()
			if len(carried) != 1 || carried[0].Fields != 1 {
				t.Errorf("emptying %q of a %s decision removed the field a person writes in: %+v", emptied, mode, carried)
			}
		}
	}
}

var namesOfClasses = []string{
	"none", "decline", "option", "choice", "words-field", "amend", "confirmation", "exchange", "shell", "finish", "go-back",
	"message-field", "receipt", "send", "send-anyway", "drawing-close", "expand",
}

func TestTheSurfaceOffersNoThirdInteraction(t *testing.T) {
	t.Parallel()
	held := standing(t)
	aDrawing(t, held, "smuggled.html", "<!doctype html><p>a drawing</p>")

	first := aCacheDecision("single-choice")
	first["evidence"] = []sketch{{"patch": "--- a/x\n+++ b/x\n@@ -1,2 +1,2 @@\n-old\n+new\n"}}
	first["illustration"] = sketch{"path": "smuggled.html"}
	second := aCacheDecision("multiple-choice")
	second["id"] = "c-4"
	page := opening(t, held.serve(written(t, aDocument("one thing to settle", []sketch{
		anEntry("c-1", "the read is unlocked", "error", "a person", "internal/cache/read.go:41"),
		anEntry("c-2", "the write counts twice", "warning", "the author patches", "internal/cache/write.go:12"),
		anEntry("c-4", "the read is unlocked again", "info", "the author patches", "internal/cache/read.go:80"),
	}, first, second))))

	var drawn struct {
		Links    []string `json:"links"`
		Controls []struct {
			Tag         string   `json:"tag"`
			Interaction string   `json:"interaction"`
			Classes     []string `json:"classes"`
		} `json:"controls"`
	}
	page.evaluate(`({
  links: Array.from(document.querySelectorAll('a'), (held) => held.getAttribute('href') || ''),
  controls: Array.from(document.querySelectorAll('ui-choice, input, textarea, button, select, ui-field, ui-button, ui-copy, ui-expand, form, output, ui-dialog'), (held) => ({
    tag: held.localName,
    interaction: held.dataset.interaction || '',
    classes: Array.from(held.classList),
  })),
})`, &drawn)

	for _, led := range drawn.Links {
		if !strings.HasPrefix(led, "#") {
			t.Errorf("the page offers a link to %q, and reading the source moves a person within the page they are on", led)
		}
	}

	var parts []string
	for _, control := range drawn.Controls {
		if !slices.Contains([]string{"answering", "sending", "viewing"}, control.Interaction) {
			t.Errorf("the page offers a <%s class=%q> that is part of none of answering, sending or viewing", control.Tag, strings.Join(control.Classes, " "))
			continue
		}
		var named []string
		for _, class := range control.Classes {
			if slices.Contains(namesOfClasses, class) {
				named = append(named, class)
			}
		}
		parts = append(parts, control.Interaction+" "+strings.Join(named, " "))
	}

	slices.Sort(parts)
	want := []string{
		"answering none decline",
		"answering option choice",
		"answering option choice",
		"answering option choice",
		"answering option choice",
		"answering words-field",
		"answering words-field",
		"sending amend",
		"sending confirmation",
		"sending exchange",
		"sending finish",
		"sending go-back",
		"sending message-field",
		"sending receipt",
		"sending send",
		"sending send-anyway",
		"viewing drawing-close",
		"viewing expand",
	}
	if diff := cmp.Diff(want, parts); diff != "" {
		t.Errorf("the controls the page offers differ (-want +got):\n%s", diff)
	}
}

func TestTheConfirmationStandsOutsideTheFormAndSubmitsNothing(t *testing.T) {
	t.Parallel()
	_, page := read(t, twoDecisionDocument)

	page.click(".send")
	page.waitFor(confirmationIsOpen)

	var drawn struct {
		Dialogs  int      `json:"dialogs"`
		Inside   int      `json:"inside"`
		Presses  []string `json:"presses"`
		FirstKey string   `json:"first"`
		Named    int      `json:"named"`
		Sending  bool     `json:"sending"`
	}
	page.evaluate(`(() => {
  const dialogs = Array.from(document.querySelectorAll('ui-dialog'));
  const submit = Array.from(document.querySelectorAll('ui-button')).find((held) => held.type === 'submit');
  const named = Array.from(document.querySelectorAll('ui-dialog .confirmation-decision'));
  return {
    dialogs: dialogs.length,
    inside: Array.from(document.querySelectorAll('form')).filter((form) => form.querySelector('ui-dialog')).length,
    presses: Array.from(dialogs[0].querySelectorAll('ui-button'), (held) => held.type),
    first: submit ? submit.classList[0] : '',
    named: named.length,
    sending: named.every((held) => held.dataset.interaction === 'sending'),
  };
})()`, &drawn)

	if drawn.Dialogs != 1 {
		t.Fatalf("the page carries %d dialogs, want the one a blank decision is confirmed in", drawn.Dialogs)
	}
	if drawn.Inside != 0 {
		t.Error("the confirmation stands inside the form, so its presses would submit it")
	}
	if len(drawn.Presses) < 3 {
		t.Errorf("the confirmation offers %v, want its two presses and a name for each blank decision", drawn.Presses)
	}
	for _, press := range drawn.Presses {
		if press != "button" {
			t.Errorf("the confirmation offers a <ui-button type=%q>, and a press inside it submits nothing", press)
		}
	}
	if drawn.Named == 0 || !drawn.Sending {
		t.Errorf("the confirmation names %d blank decisions, and each is part of sending: %v", drawn.Named, drawn.Sending)
	}
	if drawn.FirstKey != "send" {
		t.Errorf("the first submit in document order is %q, and enter reaches the press that waits for a reply", drawn.FirstKey)
	}
}

func TestTheExchangesFormWrapsTheShellSoEveryControlInItsSlotsIsPartOfTheForm(t *testing.T) {
	t.Parallel()
	_, page := read(t, twoDecisionDocument)

	var drawn struct {
		Forms    int      `json:"forms"`
		Class    string   `json:"class"`
		Shells   int      `json:"shells"`
		Cards    []string `json:"cards"`
		Controls int      `json:"controls"`
	}
	page.evaluate(`(() => {
  const forms = Array.from(document.querySelectorAll('form'));
  const form = forms[0];
  const shells = form ? form.querySelectorAll('ui-app-shell') : [];
  return {
    forms: forms.length,
    class: form ? form.className : '',
    shells: shells.length,
    cards: shells[0] ? Array.from(shells[0].querySelectorAll('ui-card'), (held) => held.className) : [],
    controls: form ? form.querySelectorAll('ui-button, ui-field, ui-choice').length : 0,
  };
})()`, &drawn)

	if drawn.Forms != 1 || !slices.Contains(strings.Fields(drawn.Class), "exchange") {
		t.Fatalf("the page carries %d forms of class %q, want the one exchange", drawn.Forms, drawn.Class)
	}
	if drawn.Shells != 1 {
		t.Fatalf("the exchange form holds %d shells, want the one the page draws", drawn.Shells)
	}
	for _, class := range []string{"decision", "composer"} {
		if !slices.ContainsFunc(drawn.Cards, func(held string) bool { return slices.Contains(strings.Fields(held), class) }) {
			t.Errorf("the shell inside the exchange form holds no card of class %q", class)
		}
	}
	if drawn.Controls == 0 {
		t.Error("the exchange form holds no control")
	}
}

func TestAPageWithNothingToDecideDrawsTheShellWithNoFormAndNoDialog(t *testing.T) {
	t.Parallel()
	held := standing(t)
	page := openingNothingToAnswer(t, held.serve(nothingToAnswerDocument))

	var drawn struct{ Forms, Dialogs, Shells int }
	page.evaluate(`({
  Forms: document.querySelectorAll('form').length,
  Dialogs: document.querySelectorAll('ui-dialog').length,
  Shells: document.querySelectorAll('ui-app-shell').length,
})`, &drawn)

	if drawn.Forms != 0 {
		t.Errorf("the page carries %d forms, and nothing is sent from it", drawn.Forms)
	}
	if drawn.Dialogs != 0 {
		t.Errorf("the page carries %d dialogs, and nothing is confirmed in it", drawn.Dialogs)
	}
	if drawn.Shells != 1 {
		t.Errorf("the page carries %d shells, want one", drawn.Shells)
	}
}

func TestAPageWithNothingToAnswerOffersNeitherControl(t *testing.T) {
	t.Parallel()
	held := standing(t)
	page := openingNothingToAnswer(t, held.serve(nothingToAnswerDocument))

	var offered []string
	page.evaluate(`Array.from(document.querySelectorAll('ui-choice, ui-choice-group, input, textarea, button, select, form, output, dialog, ui-dialog'), (held) => held.localName + '.' + held.className)`, &offered)
	if len(offered) != 0 {
		t.Errorf("a page holding no decision offers %v, and there is nothing on it to answer and nothing to send", offered)
	}

	var said string
	page.evaluate(`document.querySelector('.nothing-said').textContent.trim()`, &said)
	for _, wanted := range []string{"nothing", "waiting"} {
		if !strings.Contains(strings.ToLower(said), wanted) {
			t.Errorf("a page holding no decision says %q, and it says there is nothing to answer and that nobody is waiting", said)
		}
	}
}
