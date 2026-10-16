package render_test

import (
	"crypto/sha256"
	"encoding/base64"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/agora/document"
	"github.com/socrates-tools/socrates/agora/render"
	"golang.org/x/net/html"
)

const answering, sending = "answering", "sending"

func oneDecision(mode document.Mode) document.Document {
	return document.Document{
		Title:    "one thing to settle",
		Preamble: "what follows is what was seen",
		Decisions: []document.Decision{{
			ID:       "c-1",
			Severity: "error",
			Audience: "a person",
			Location: "internal/cache/read.go:41",
			Question: "should the cache be read before the lock is taken",
			Meaning:  "two callers arriving together read the same stale entry",
			Mode:     mode,
			Options: []document.Option{
				{Label: "take the lock first", Meaning: "correctness first", Command: "patch c-1"},
				{Label: "leave it", Meaning: "the staleness is acceptable", Command: `agree c-1 --reason "acceptable"`},
			},
		}},
	}
}

func rendered(t *testing.T, read document.Document) *html.Node {
	t.Helper()

	written, err := render.Page(read, render.Frame{Base: "/sessions/a1b2c3d4", Token: "a-token", Revision: 1})
	if err != nil {
		t.Fatalf("rendering the page: %v", err)
	}
	parsed, err := html.Parse(strings.NewReader(written))
	if err != nil {
		t.Fatalf("reading back the rendered page: %v", err)
	}
	return parsed
}

func elements(node *html.Node, named ...string) []*html.Node {
	var found []*html.Node
	if node.Type == html.ElementNode && slices.Contains(named, node.Data) {
		found = append(found, node)
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		found = append(found, elements(child, named...)...)
	}
	return found
}

func attribute(node *html.Node, named string) string {
	for _, held := range node.Attr {
		if held.Key == named {
			return held.Val
		}
	}
	return ""
}

func hasAttribute(node *html.Node, named string) bool {
	for _, held := range node.Attr {
		if held.Key == named {
			return true
		}
	}
	return false
}

func hasClass(node *html.Node, named string) bool {
	return slices.Contains(strings.Fields(attribute(node, "class")), named)
}

func text(node *html.Node) string {
	var written strings.Builder
	if node.Type == html.TextNode {
		written.WriteString(node.Data)
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		written.WriteString(text(child))
	}
	return written.String()
}

func within(node *html.Node, named string) *html.Node {
	for _, held := range elements(node, "div") {
		if held != node && hasClass(held, named) {
			return held
		}
	}
	return nil
}

func policyOf(t *testing.T, page *html.Node) string {
	t.Helper()

	for _, stated := range elements(page, "meta") {
		if strings.EqualFold(attribute(stated, "http-equiv"), "content-security-policy") {
			return attribute(stated, "content")
		}
	}
	t.Fatal("the page states no policy at all")
	return ""
}

func embeddedStylesheet(t *testing.T, page *html.Node) string {
	t.Helper()

	embedded := elements(page, "style")
	if len(embedded) != 1 {
		t.Fatalf("the page carries %d embedded stylesheets, want the one the document stated", len(embedded))
	}
	return text(embedded[0])
}

func painted(read document.Document) document.Document {
	read.Stylesheet = `:root{--panel-face:#fbf9f4}.panel{background:var(--panel-face);border:1px solid #d8d2c4}.chart rect{fill:url(#band)}`
	for position := range read.Decisions {
		read.Decisions[position].Fragment = document.Fragment{
			Markup: `<div class="panel"><span class="value" id="total">$1,240.00</span></div>`,
		}
	}
	return read
}

func TestThePageCarriesTheBrandRatherThanTheNameOfTheServer(t *testing.T) {
	page := rendered(t, oneDecision(document.One))

	var wordmark *html.Node
	for _, held := range elements(page, "span") {
		if hasClass(held, "wordmark") {
			wordmark = held
		}
	}
	if wordmark == nil {
		t.Fatal("the page carries no wordmark")
	}
	if said := strings.TrimSpace(text(wordmark)); said != "socrates" {
		t.Errorf("the wordmark says %q, want socrates", said)
	}

	for _, held := range elements(page, "body") {
		if strings.Contains(strings.ToLower(text(held)), "agora") {
			t.Error("the page shows a person the name of the server that serves it")
		}
	}
}

func overNoticeIn(t *testing.T, page *html.Node) *html.Node {
	t.Helper()

	for _, held := range elements(page, "p") {
		if hasClass(held, "over-notice") {
			return held
		}
	}
	t.Fatal("the page carries no notice that the exchange is over")
	return nil
}

func TestAnExchangeThatIsOverSaysSoUnderTheTitle(t *testing.T) {
	written, err := render.Page(oneDecision(document.One), render.Frame{Base: "/sessions/a1b2c3d4", Token: "a-token", Revision: 1, Over: true})
	if err != nil {
		t.Fatalf("rendering the page: %v", err)
	}
	over, err := html.Parse(strings.NewReader(written))
	if err != nil {
		t.Fatalf("reading back the rendered page: %v", err)
	}

	notice := overNoticeIn(t, over)
	if hasAttribute(notice, "hidden") {
		t.Error("a page rendered over hides the notice that says so")
	}
	if !strings.Contains(strings.ToLower(text(notice)), "over") {
		t.Errorf("the notice says %q, and it says the exchange is over", text(notice))
	}
	var heading *html.Node
	for _, held := range elements(over, "header") {
		if hasClass(held, "heading") {
			heading = held
		}
	}
	if heading == nil || !slices.Contains(elements(heading, "p"), notice) {
		t.Error("the notice that the exchange is over is not under the title")
	}

	if open := overNoticeIn(t, rendered(t, oneDecision(document.One))); !hasAttribute(open, "hidden") {
		t.Error("a page that is open says the exchange is over")
	}
}

func TestThePolicyStatesTheHashOfTheBytesEmbedded(t *testing.T) {
	// scenario: The embedded stylesheet paints
	page := rendered(t, painted(oneDecision(document.One)))

	embedded := embeddedStylesheet(t, page)
	for _, stated := range []string{
		".painted{--panel-face:#fbf9f4;}",
		".painted .drawn-panel{background:var(--panel-face);border:1px solid #d8d2c4;}",
		".painted .drawn-chart rect{fill:url(#drawn-band);}",
	} {
		if !strings.Contains(embedded, stated) {
			t.Errorf("the page embedded no %q:\n%s", stated, embedded)
		}
	}

	digest := sha256.Sum256([]byte(embedded))
	stated := "'sha256-" + base64.StdEncoding.EncodeToString(digest[:]) + "'"
	if policy := policyOf(t, page); !strings.Contains(policy, stated) {
		t.Errorf("the policy does not admit the bytes the page embedded as %s:\n%s", stated, policy)
	}
}

func TestInlineStyleIsNotAdmittedAsAClass(t *testing.T) {
	policy := policyOf(t, rendered(t, painted(oneDecision(document.One))))

	if !strings.Contains(policy, "style-src-attr 'none'") {
		t.Errorf("the policy admits style carried on an element, and style on an element cannot be scoped:\n%s", policy)
	}
	for _, admitted := range []string{"unsafe-inline", "unsafe-hashes", "unsafe-eval", "*"} {
		if strings.Contains(policy, admitted) {
			t.Errorf("the policy admits %q, and the reduction stops being the only door:\n%s", admitted, policy)
		}
	}
	if strings.Count(policy, "'sha256-") != 1 {
		t.Errorf("the policy states %d hashes, want the one stylesheet the page embedded:\n%s", strings.Count(policy, "'sha256-"), policy)
	}
}

func TestOneStylesheetPaintsEveryIllustrationInADocument(t *testing.T) {
	read := painted(oneDecision(document.One))
	second := read.Decisions[0]
	second.ID = "c-2"
	read.Decisions = append(read.Decisions, second)

	page := rendered(t, read)

	boxes := 0
	for _, box := range elements(page, "div") {
		if !hasClass(box, "illustration") {
			continue
		}
		boxes++
		inside := within(box, render.PaintedClass)
		if inside == nil {
			t.Errorf("an illustration renders outside the box a stylesheet paints in:\n%s", text(box))
			continue
		}
		if !strings.Contains(text(inside), "$1,240.00") {
			t.Errorf("the illustration's box holds none of what the caller supplied:\n%s", text(inside))
		}
	}

	if boxes != 2 {
		t.Fatalf("two decisions carrying an illustration rendered %d boxes", boxes)
	}
	if embedded := embeddedStylesheet(t, page); !strings.Contains(embedded, render.Illustration.Selector) {
		t.Errorf("the one stylesheet the page embedded names no illustration:\n%s", embedded)
	}
}

func TestADocumentThatCarriesNoStylesheetStillServes(t *testing.T) {
	read := painted(oneDecision(document.One))
	read.Stylesheet = ""

	page := rendered(t, read)

	if embedded := elements(page, "style"); len(embedded) != 0 {
		t.Errorf("a document carrying no stylesheet embedded %d of them", len(embedded))
	}
	if policy := policyOf(t, page); strings.Contains(policy, "sha256-") {
		t.Errorf("a document carrying no stylesheet states a hash for one:\n%s", policy)
	}

	unpainted := false
	for _, box := range elements(page, "div") {
		if hasClass(box, "illustration") && strings.Contains(text(box), "$1,240.00") {
			unpainted = true
		}
	}
	if !unpainted {
		t.Error("a document carrying no stylesheet renders no illustration either")
	}
}

func TestADocumentRendersWithoutTheCallerWritingMarkup(t *testing.T) {
	read := oneDecision(document.One)
	read.Decisions[0].Fragment = document.Fragment{}

	page := rendered(t, read)

	cards := cardsIn(page)
	if len(cards) != 1 {
		t.Fatalf("one decision rendered %d cards", len(cards))
	}

	card := cards[0]
	if named := attribute(card, "data-decision"); named != "c-1" {
		t.Errorf("the card names the decision %q, want %q", named, "c-1")
	}
	for _, expected := range []string{
		"should the cache be read before the lock is taken",
		"two callers arriving together read the same stale entry",
		"take the lock first",
		"leave it",
		"internal/cache/read.go:41",
	} {
		if !strings.Contains(text(card), expected) {
			t.Errorf("the card does not carry %q:\n%s", expected, text(card))
		}
	}
	for kind, value := range map[string]string{"severity": "error", "audience": "a person"} {
		mark := markIn(card, kind)
		if mark == nil {
			t.Errorf("the card carries no %s mark", kind)
			continue
		}
		if held := attribute(mark, "data-"+kind); held != value {
			t.Errorf("the card's %s mark carries %q, and the document stated %q", kind, held, value)
		}
		if spoken := attribute(mark, "aria-label"); !strings.Contains(strings.ToLower(spoken), value) {
			t.Errorf("the card's %s mark is announced as %q, and a person who cannot see the icon hears %q in it", kind, spoken, value)
		}
	}

	for _, choice := range elements(card, "input") {
		if attribute(choice, "type") != "radio" {
			t.Errorf("a single-choice decision renders a %q", attribute(choice, "type"))
		}
	}
}

func markIn(node *html.Node, kind string) *html.Node {
	for _, held := range elements(node, "span") {
		if hasClass(held, "mark") && hasClass(held, kind) {
			return held
		}
	}
	return nil
}

func threeOptions() document.Document {
	read := oneDecision(document.Many)
	read.Decisions[0].Options = append(read.Decisions[0].Options, document.Option{
		Label: "shard the map", Meaning: "contention matters", Command: "patch c-1 --shard",
	})
	return read
}

func optionsIn(page *html.Node) []*html.Node {
	var found []*html.Node
	for _, held := range elements(page, "label") {
		if hasClass(held, "option") {
			found = append(found, held)
		}
	}
	return found
}

func labelsIn(page *html.Node) []string {
	var read []string
	for _, option := range optionsIn(page) {
		for _, named := range elements(option, "span") {
			if hasClass(named, "option-label") {
				read = append(read, strings.TrimSpace(text(named)))
			}
		}
	}
	return read
}

func positionsIn(page *html.Node) map[string]string {
	held := map[string]string{}
	for _, option := range optionsIn(page) {
		label := ""
		for _, named := range elements(option, "span") {
			if hasClass(named, "option-label") {
				label = strings.TrimSpace(text(named))
			}
		}
		for _, choice := range elements(option, "input") {
			held[label] = attribute(choice, "value")
		}
	}
	return held
}

func TestTheOptionTheCallerWouldTakeIsReadFirst(t *testing.T) {
	read := threeOptions()
	read.Decisions[0].Options[2].Recommended = true

	page := rendered(t, read)

	want := []string{"shard the map", "take the lock first", "leave it"}
	if got := labelsIn(page); !slices.Equal(want, got) {
		t.Errorf("the options are read in the order %v, want %v", got, want)
	}

	marked := 0
	for _, option := range optionsIn(page) {
		if attribute(option, "data-recommended") == "yes" {
			marked++
		}
	}
	if marked != 1 {
		t.Errorf("%d options are marked as the one the caller would take, want the one the document marked", marked)
	}
}

func TestRearrangingTheOptionsDoesNotRearrangeTheAnswer(t *testing.T) {
	read := threeOptions()
	read.Decisions[0].Options[2].Recommended = true

	want := map[string]string{"take the lock first": "0", "leave it": "1", "shard the map": "2"}
	if diff := cmp.Diff(want, positionsIn(rendered(t, read))); diff != "" {
		t.Errorf("what each option posts back differs from the position the document gave it (-want +got):\n%s", diff)
	}
}

func TestTheRecommendedOptionArrivesChosen(t *testing.T) {
	read := threeOptions()
	read.Decisions[0].Options[2].Recommended = true

	page := rendered(t, read)

	var chosen []string
	for _, choice := range elements(page, "input") {
		if hasAttribute(choice, "checked") {
			chosen = append(chosen, attribute(choice, "value"))
		}
	}
	if diff := cmp.Diff([]string{"2"}, chosen); diff != "" {
		t.Errorf("what arrives chosen differs from the option the document marked (-want +got):\n%s", diff)
	}
}

func TestADocumentThatMarksNothingIsUnaffected(t *testing.T) {
	page := rendered(t, threeOptions())

	want := []string{"take the lock first", "leave it", "shard the map"}
	if got := labelsIn(page); !slices.Equal(want, got) {
		t.Errorf("a document that marks nothing renders its options as %v, want the order it listed them %v", got, want)
	}
	for _, option := range optionsIn(page) {
		if attribute(option, "data-recommended") != "" {
			t.Error("an option is marked, and the document marked nothing")
		}
		for _, named := range elements(option, "span") {
			if hasClass(named, "option-mark") {
				t.Error("the page carries a mark, and the document marked nothing")
			}
		}
	}
	for _, choice := range elements(page, "input") {
		if hasAttribute(choice, "checked") {
			t.Errorf("the option %q arrives chosen, and the document marked nothing", attribute(choice, "value"))
		}
	}
}

func noneControlIn(page *html.Node, decision string) *html.Node {
	for _, held := range elements(page, "label") {
		if !hasClass(held, "none") {
			continue
		}
		for _, input := range elements(held, "input") {
			if attribute(input, "data-decision") == decision {
				return held
			}
		}
	}
	return nil
}

func TestASingleChoiceDecisionOffersTakingNoneOfThem(t *testing.T) {
	page := rendered(t, oneDecision(document.One))

	none := noneControlIn(page, "c-1")
	if none == nil {
		t.Fatal("a single-choice decision renders no control for taking none of the options")
	}
	if hasClass(none, "option") {
		t.Error("the control for taking none of the options is counted as an option")
	}

	inputs := elements(none, "input")
	if len(inputs) != 1 {
		t.Fatalf("the control for taking none of the options carries %d inputs, want one", len(inputs))
	}

	control := inputs[0]
	if attribute(control, "type") != "radio" {
		t.Errorf("the control for taking none of the options renders as %q, want a radio", attribute(control, "type"))
	}
	if attribute(control, "data-part") != answering {
		t.Errorf("the control for taking none of the options carries data-part %q, want %q", attribute(control, "data-part"), answering)
	}
	if hasClass(control, "choice") {
		t.Error("the control for taking none of the options is counted as a choice")
	}

	var named string
	for _, choice := range elements(page, "input") {
		if hasClass(choice, "choice") && attribute(choice, "data-decision") == "c-1" {
			named = attribute(choice, "name")
			break
		}
	}
	if named == "" || attribute(control, "name") != named {
		t.Errorf("the control for taking none of the options is named %q, want the radio group %q", attribute(control, "name"), named)
	}

	if offered := len(optionsIn(page)); offered != 2 {
		t.Errorf("a single-choice decision offering two options and a way to take none of them counts %d options", offered)
	}
}

func TestAMultipleChoiceDecisionOffersNoWayToTakeNoneOfThem(t *testing.T) {
	page := rendered(t, oneDecision(document.Many))

	if none := noneControlIn(page, "c-1"); none != nil {
		t.Error("a multiple-choice decision renders a control for taking none of the options")
	}
}

func statingSomething() document.Document {
	read := oneDecision(document.One)
	read.Decisions[0].Options = nil
	read.Decisions[0].Mode = ""
	return read
}

func cardsIn(page *html.Node) []*html.Node {
	var found []*html.Node
	for _, held := range elements(page, "li") {
		if hasClass(held, "decision") {
			found = append(found, held)
		}
	}
	return found
}

func TestADecisionWithNoOptionsRendersAsACardThatAsksNothing(t *testing.T) {
	page := rendered(t, statingSomething())

	cards := cardsIn(page)
	if len(cards) != 1 {
		t.Fatalf("one decision that asks nothing rendered %d cards", len(cards))
	}

	card := cards[0]
	for _, stated := range []string{
		"should the cache be read before the lock is taken",
		"two callers arriving together read the same stale entry",
	} {
		if !strings.Contains(text(card), stated) {
			t.Errorf("the card does not carry %q:\n%s", stated, text(card))
		}
	}

	if offered := optionsIn(page); len(offered) != 0 {
		t.Errorf("a decision that asks nothing rendered %d options", len(offered))
	}
	for _, control := range elements(card, "input", "select", "fieldset", "legend") {
		t.Errorf("a decision that asks nothing renders a <%s class=%q>, and nothing on it is a choice",
			control.Data, attribute(control, "class"))
	}
	for _, said := range []string{"Choose one", "Choose as many as apply"} {
		if strings.Contains(text(card), said) {
			t.Errorf("a decision that asks nothing says %q, and there is nothing on it to choose among", said)
		}
	}
}

func TestADecisionWithNoOptionsStillCarriesTheField(t *testing.T) {
	read := statingSomething()
	fields := reflect.ValueOf(&read.Decisions[0]).Elem()

	for position := range fields.NumField() {
		named := fields.Type().Field(position).Name
		t.Run(named, func(t *testing.T) {
			emptied := statingSomething()
			held := reflect.ValueOf(&emptied.Decisions[0]).Elem().Field(position)
			held.Set(reflect.Zero(held.Type()))

			if !everyCardCarriesTheField(rendered(t, emptied)) {
				t.Errorf("emptying %s removed the field a person writes in", named)
			}
		})
	}

	page := rendered(t, read)
	if !everyCardCarriesTheField(page) {
		t.Fatal("a decision that asks nothing renders no field for a person to write in")
	}

	for _, written := range elements(page, "textarea") {
		if !hasClass(written, "words-field") {
			continue
		}
		if held := attribute(written, "rows"); held != "2" {
			t.Errorf("the field a person writes in is %q rows, and its size is the surface's", held)
		}
	}
	for _, labelled := range elements(page, "span") {
		if hasClass(labelled, "words-label") && strings.TrimSpace(text(labelled)) != "In your own words" {
			t.Errorf("the field a person writes in is labelled %q, and its label is the surface's", text(labelled))
		}
	}
}

func TestTheFreeTextFieldCannotBeSuppressed(t *testing.T) {
	for _, mode := range document.Modes() {
		t.Run(string(mode), func(t *testing.T) {
			read := oneDecision(mode)
			fields := reflect.ValueOf(&read.Decisions[0]).Elem()

			for position := range fields.NumField() {
				named := fields.Type().Field(position).Name
				t.Run(named, func(t *testing.T) {
					emptied := oneDecision(mode)
					held := reflect.ValueOf(&emptied.Decisions[0]).Elem().Field(position)
					held.Set(reflect.Zero(held.Type()))

					if !everyCardCarriesTheField(rendered(t, emptied)) {
						t.Errorf("emptying %s removed the field a person writes in", named)
					}
				})
			}

			if !everyCardCarriesTheField(rendered(t, read)) {
				t.Error("a whole decision renders no field for a person to write in")
			}
		})
	}
}

func everyCardCarriesTheField(page *html.Node) bool {
	cards := 0
	for _, card := range elements(page, "li") {
		if !hasClass(card, "decision") {
			continue
		}
		cards++
		fields := 0
		for _, written := range elements(card, "textarea") {
			if hasClass(written, "words-field") {
				fields++
			}
		}
		if fields != 1 {
			return false
		}
	}
	return cards > 0
}

func TestTheSurfaceOffersNoThirdInteraction(t *testing.T) {
	read := citingASource()
	read.Decisions[0].Fragment = document.Fragment{Markup: `<form><input name="smuggled"><button>publish</button></form>`}
	read.Decisions[0].Evidence = document.Evidence{Patch: "--- a/x\n+++ b/x\n@@ -1,2 +1,2 @@\n-old\n+new\n"}

	multiple := read.Decisions[0]
	multiple.ID = "c-4"
	multiple.Mode = document.Many
	multiple.Fragment = document.Fragment{}
	multiple.Evidence = document.Evidence{}
	read.Decisions = append(read.Decisions, multiple)

	page := rendered(t, read)

	for _, leading := range elements(page, "a") {
		if held := attribute(leading, "href"); !strings.HasPrefix(held, "#") {
			t.Errorf("the page offers a link to %q, and reading the source moves a person within the page they are on", held)
		}
	}

	var parts []string
	for _, control := range elements(page, "input", "textarea", "button", "select", "form", "output", "dialog") {
		part := attribute(control, "data-part")
		if part != answering && part != sending {
			t.Errorf("the page offers a <%s class=%q> that is part of neither answering nor sending", control.Data, attribute(control, "class"))
			continue
		}
		parts = append(parts, part+" "+attribute(control, "class"))
	}

	slices.Sort(parts)
	want := []string{
		"answering choice",
		"answering choice",
		"answering choice",
		"answering choice",
		"answering decline",
		"answering words-field",
		"answering words-field",
		"sending amend",
		"sending confirmation",
		"sending exchange shell",
		"sending finish",
		"sending go-back",
		"sending message-field",
		"sending receipt",
		"sending send",
		"sending send-anyway",
	}
	if !slices.Equal(want, parts) {
		t.Errorf("the controls the page offers are %v, want %v", parts, want)
	}
}

func TestTheConfirmationStandsOutsideTheFormAndSubmitsNothing(t *testing.T) {
	page := rendered(t, citingASource())

	asking := elements(page, "dialog")
	if len(asking) != 1 {
		t.Fatalf("the page carries %d dialogs, want the one a blank decision is confirmed in", len(asking))
	}

	for _, form := range elements(page, "form") {
		if len(elements(form, "dialog")) != 0 {
			t.Error("the confirmation stands inside the form, so its presses would submit it")
		}
	}

	for _, press := range elements(asking[0], "button") {
		if attribute(press, "type") != "button" {
			t.Errorf("the confirmation offers a <button type=%q>, and a press inside it submits nothing", attribute(press, "type"))
		}
	}

	submits := elements(page, "button")
	first := ""
	for _, press := range submits {
		if attribute(press, "type") == "submit" {
			first = attribute(press, "class")
			break
		}
	}
	if first != "send" {
		t.Errorf("the first submit in document order is %q, and enter reaches the press that waits for a reply", first)
	}
}

func citingASource() document.Document {
	read := oneDecision(document.One)
	read.Source = []document.Entry{
		{ID: "c-1", Name: "the read is unlocked", Severity: "error", Audience: "a person", Location: "internal/cache/read.go:41"},
		{ID: "c-2", Name: "the write counts twice", Severity: "warning", Audience: "the author patches", Location: "internal/cache/write.go:12"},
		{ID: "c-3", Severity: "info", Audience: "the author patches", Location: "internal/cache/read.go:80"},
	}
	return read
}

func nothingToDecide() document.Document {
	read := citingASource()
	read.Decisions = nil
	return read
}

func listedLeastSevereFirst() document.Document {
	read := document.Document{
		Title: "three things to settle",
		Source: []document.Entry{
			{ID: "c-3", Name: "the log line says nothing", Severity: "info", Audience: "the author patches", Location: "internal/cache/read.go:80"},
			{ID: "c-2", Name: "the write counts twice", Severity: "warning", Audience: "the author patches", Location: "internal/cache/write.go:12"},
			{ID: "c-1", Name: "the read is unlocked", Severity: "error", Audience: "a person", Location: "internal/cache/read.go:41"},
		},
	}
	for _, entry := range read.Source {
		read.Decisions = append(read.Decisions, document.Decision{
			ID:       entry.ID,
			Name:     entry.Name,
			Severity: entry.Severity,
			Audience: entry.Audience,
			Location: entry.Location,
			Question: "what should be done about " + entry.ID,
			Mode:     document.One,
			Options:  []document.Option{{Label: "leave it", Meaning: "acceptable here", Command: "agree " + entry.ID}},
		})
	}
	return read
}

func cardsRead(page *html.Node) []string {
	var named []string
	for _, card := range cardsIn(page) {
		named = append(named, attribute(card, "data-decision"))
	}
	return named
}

func entriesRead(page *html.Node) []string {
	var named []string
	for _, entry := range entriesIn(page) {
		named = append(named, attribute(elements(entry, "a", "span")[0], "href")+strings.TrimSpace(text(entry)))
	}
	return named
}

func severitiesRead(within []*html.Node) []string {
	var read []string
	for _, held := range within {
		stated := ""
		for _, badge := range elements(held, "span") {
			if hasClass(badge, "severity") {
				stated = attribute(badge, "data-severity")
			}
		}
		read = append(read, stated)
	}
	return read
}

func TestTheMostSevereCardIsTheFirstAPersonMeets(t *testing.T) {
	page := rendered(t, listedLeastSevereFirst())

	want := []string{"c-1", "c-2", "c-3"}
	if got := cardsRead(page); !slices.Equal(want, got) {
		t.Errorf("the cards arrive as %v, want %v, and a person triages by severity", got, want)
	}
	if got := severitiesRead(cardsIn(page)); !slices.Equal([]string{"error", "warning", "info"}, got) {
		t.Errorf("the cards arrive carrying the severities %v", got)
	}
}

func TestTheSourceListIsOrderedTheSameWay(t *testing.T) {
	read := listedLeastSevereFirst()
	page := rendered(t, read)

	want := []string{"error", "warning", "info"}
	if got := severitiesRead(entriesIn(page)); !slices.Equal(want, got) {
		t.Errorf("the entries are listed carrying the severities %v, want %v", got, want)
	}

	written, err := render.Page(read, render.Frame{Base: "/sessions/a1b2c3d4", Token: "a-token", Revision: 1})
	if err != nil {
		t.Fatalf("rendering the page: %v", err)
	}
	listed := strings.Index(written, `class="entries"`)
	decided := strings.Index(written, `class="decisions"`)
	if listed < 0 || decided < 0 || listed > decided {
		t.Errorf("the list is at %d and the decisions at %d, and ordering them changed which comes first", listed, decided)
	}
}

func TestEqualSeveritiesKeepTheOrderTheDocumentListedThem(t *testing.T) {
	read := listedLeastSevereFirst()
	for position := range read.Source {
		read.Source[position].Severity = "warning"
		read.Decisions[position].Severity = "warning"
	}

	page := rendered(t, read)

	want := []string{"c-3", "c-2", "c-1"}
	if got := cardsRead(page); !slices.Equal(want, got) {
		t.Errorf("three decisions of one severity arrive as %v, want the order the document listed them %v", got, want)
	}
	if got := entriesRead(page); len(got) != 3 || !strings.Contains(got[0], "the log line says nothing") {
		t.Errorf("three entries of one severity are listed as %v, want the order the document listed them", got)
	}
}

func TestTheCallerCannotOrderThePage(t *testing.T) {
	read := listedLeastSevereFirst()
	read.Decisions[0].Question = "the informational one the caller listed first"

	page := rendered(t, read)

	cards := cardsIn(page)
	if len(cards) != 3 {
		t.Fatalf("three decisions rendered %d cards", len(cards))
	}
	if named := attribute(cards[0], "data-decision"); named != "c-1" {
		t.Errorf("the first card a person meets is %q, and the error is %q", named, "c-1")
	}
	if strings.Contains(text(cards[0]), "the informational one the caller listed first") {
		t.Errorf("the card the caller listed first is the card a person meets first:\n%s", text(cards[0]))
	}
}

func TestAnEntryWhoseSourceStatesNoSeveritySortsLast(t *testing.T) {
	read := listedLeastSevereFirst()
	read.Source = append([]document.Entry{{ID: "c-0", Name: "nobody classified this one", Audience: "the author patches"}}, read.Source...)

	page := rendered(t, read)

	want := []string{"error", "warning", "info", ""}
	if got := severitiesRead(entriesIn(page)); !slices.Equal(want, got) {
		t.Errorf("the entries are listed carrying the severities %v, want %v, and an unclassified entry is not evidence of unimportance", got, want)
	}

	last := entriesIn(page)[3]
	if !strings.Contains(text(last), "nobody classified this one") {
		t.Errorf("the entry listed last is %q", strings.TrimSpace(text(last)))
	}
}

func entriesIn(page *html.Node) []*html.Node {
	var found []*html.Node
	for _, held := range elements(page, "li") {
		if hasClass(held, "entry") {
			found = append(found, held)
		}
	}
	return found
}

func TestTheSourceIsReadBeforeTheDecisions(t *testing.T) {
	page := rendered(t, citingASource())

	entries := entriesIn(page)
	if len(entries) != 3 {
		t.Fatalf("a document citing three source entries rendered %d of them", len(entries))
	}

	for position, expected := range [][]string{
		{"the read is unlocked", "internal/cache/read.go:41"},
		{"the write counts twice", "internal/cache/write.go:12"},
		{"c-3", "internal/cache/read.go:80"},
	} {
		for _, wanted := range expected {
			if !strings.Contains(text(entries[position]), wanted) {
				t.Errorf("entry %d does not carry %q:\n%s", position+1, wanted, text(entries[position]))
			}
		}
	}
	if got := severitiesRead(entries); !slices.Equal([]string{"error", "warning", "info"}, got) {
		t.Errorf("the entries carry the severities %v, and each carries its own source's", got)
	}

	written, err := render.Page(citingASource(), render.Frame{Base: "/sessions/a1b2c3d4", Token: "a-token", Revision: 1})
	if err != nil {
		t.Fatalf("rendering the page: %v", err)
	}
	listed := strings.Index(written, `class="entries"`)
	decided := strings.Index(written, `class="decisions"`)
	if listed < 0 || decided < 0 || listed > decided {
		t.Errorf("the list is at %d and the decisions at %d, and a person reads the list first", listed, decided)
	}
}

func TestAnEntryNoDecisionCitesIsStillShown(t *testing.T) {
	page := rendered(t, citingASource())

	var named []string
	for _, entry := range entriesIn(page) {
		named = append(named, strings.TrimSpace(text(entry)))
	}
	if len(named) != 3 {
		t.Fatalf("three entries and one decision rendered %d entries: %v", len(named), named)
	}

	cards := 0
	for _, held := range elements(page, "li") {
		if hasClass(held, "decision") {
			cards++
		}
	}
	if cards != 1 {
		t.Errorf("three entries and one decision rendered %d decisions, and an entry a decision covers is not a second decision", cards)
	}
}

func TestAnEntryLeadsToTheDecisionThatCoversIt(t *testing.T) {
	page := rendered(t, citingASource())

	leading := map[string]string{}
	for _, entry := range entriesIn(page) {
		for _, link := range elements(entry, "a") {
			for _, named := range elements(link, "span") {
				if hasClass(named, "entry-name") {
					leading[strings.TrimSpace(text(named))] = attribute(link, "href")
				}
			}
		}
	}

	if leading["the read is unlocked"] != "#decision-c-1" {
		t.Errorf("the entry a decision covers leads to %q, and the decision is rendered at %q", leading["the read is unlocked"], "#decision-c-1")
	}
	if len(leading) != 1 {
		t.Errorf("%d entries lead somewhere, and only the one a decision covers has anywhere to lead: %v", len(leading), leading)
	}

	held := false
	for _, card := range elements(page, "li") {
		if hasClass(card, "decision") && attribute(card, "id") == "decision-c-1" {
			held = true
		}
	}
	if !held {
		t.Error("the entry leads to an identifier the page does not carry")
	}
}

func TestTheListIsNotTheCallersProse(t *testing.T) {
	read := citingASource()
	read.Preamble = "there were three things and none of them matters"

	page := rendered(t, read)

	for _, entry := range entriesIn(page) {
		if strings.Contains(text(entry), "none of them matters") {
			t.Errorf("an entry carries the preamble the caller wrote:\n%s", text(entry))
		}
	}
	for _, wanted := range []string{"the read is unlocked", "the write counts twice", "c-3"} {
		held := false
		for _, entry := range entriesIn(page) {
			if strings.Contains(text(entry), wanted) {
				held = true
			}
		}
		if !held {
			t.Errorf("the list does not carry %q, and what it lists comes from the entries", wanted)
		}
	}
}

func TestADocumentCitingNoSourceIndexesItsDecisions(t *testing.T) {
	read := listedLeastSevereFirst()
	read.Source = nil

	page := rendered(t, read)

	entries := entriesIn(page)
	if got := severitiesRead(entries); !slices.Equal([]string{"error", "warning", "info"}, got) {
		t.Fatalf("a document citing no source indexes entries carrying the severities %v, want one per decision, most severe first", got)
	}
	for position, card := range cardsIn(page) {
		links := elements(entries[position], "a")
		if len(links) != 1 || attribute(links[0], "href") != "#"+attribute(card, "id") {
			t.Errorf("entry %d does not lead to the card at the same place, %q", position+1, attribute(card, "id"))
		}
	}
}

func TestAPageWithNothingToAnswerOffersNeitherControl(t *testing.T) {
	// scenario: A page with nothing to answer offers neither of the two
	page := rendered(t, nothingToDecide())

	for _, control := range elements(page, "input", "textarea", "button", "select", "form", "output", "dialog") {
		t.Errorf("a page holding no decision offers a <%s class=%q>, and there is nothing on it to answer and nothing to send",
			control.Data, attribute(control, "class"))
	}

	said := ""
	for _, held := range elements(page, "p") {
		if hasClass(held, "nothing-said") {
			said = strings.TrimSpace(text(held))
		}
	}
	for _, wanted := range []string{"nothing", "waiting"} {
		if !strings.Contains(strings.ToLower(said), wanted) {
			t.Errorf("a page holding no decision says %q, and it says there is nothing to answer and that nobody is waiting", said)
		}
	}
}

func TestThePageCarriesTheBaseItsRequestsAreMadeAgainst(t *testing.T) {
	const base = "/sessions/a1b2c3d4"

	written, err := render.Page(oneDecision(document.One), render.Frame{Base: base, Token: "a-token", Revision: 1})
	if err != nil {
		t.Fatalf("rendering the page: %v", err)
	}

	parsed, err := html.Parse(strings.NewReader(written))
	if err != nil {
		t.Fatalf("reading back the rendered page: %v", err)
	}

	bodies := elements(parsed, "body")
	if len(bodies) != 1 {
		t.Fatalf("the page holds %d bodies", len(bodies))
	}
	if carried := attribute(bodies[0], "data-base"); carried != base {
		t.Errorf("the page carries the base %q, and it was served under %q", carried, base)
	}
}
