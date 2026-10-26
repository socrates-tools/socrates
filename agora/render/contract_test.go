package render_test

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/socrates-tools/socrates/agora/assets"
	"github.com/socrates-tools/socrates/agora/document"
	"github.com/socrates-tools/socrates/agora/render"
	"golang.org/x/net/html"
)

const primitivePrefix = "ui-"

const defaultSlot = ""

type element struct {
	Attributes []string `json:"attributes"`
	Slots      []string `json:"slots"`
}

type contract struct {
	Elements map[string]element `json:"elements"`
}

type writing struct {
	attributes []string
	slots      []string
	dataParts  []string
}

type writings map[string]writing

func declaredByTheContract(t *testing.T) contract {
	t.Helper()

	held, err := fs.ReadFile(assets.Files(), "contract.json")
	if err != nil {
		t.Fatalf("reading the contract the build writes: %v", err)
	}
	var declared contract
	if err := json.Unmarshal(held, &declared); err != nil {
		t.Fatalf("parsing the contract: %v", err)
	}
	return declared
}

func writtenByThePartials(t *testing.T) writings {
	t.Helper()

	drawn, err := render.Primitives()
	if err != nil {
		t.Fatalf("drawing the partials: %v", err)
	}
	written := writings{}
	for name, markup := range drawn {
		parsed, err := html.Parse(strings.NewReader(markup))
		if err != nil {
			t.Fatalf("reading back the partial %s: %v", name, err)
		}
		for key, held := range writingsOf(parsed) {
			written[key] = merged(written[key], held)
		}
	}
	return written
}

func merged(first, second writing) writing {
	return writing{
		attributes: append(first.attributes, second.attributes...),
		slots:      append(first.slots, second.slots...),
		dataParts:  append(first.dataParts, second.dataParts...),
	}
}

func writingsOf(node *html.Node) writings {
	found := writings{}
	note := func(element string, held writing) {
		found[element] = merged(found[element], held)
	}
	var visit func(node *html.Node, within string)
	visit = func(node *html.Node, within string) {
		if node.Type == html.ElementNode {
			if node.Data == "iframe" {
				if framed, err := html.Parse(strings.NewReader(attribute(node, "srcdoc"))); err == nil {
					for key, held := range writingsOf(framed) {
						note(key, held)
					}
				}
			}
			if strings.HasPrefix(node.Data, primitivePrefix) {
				within = node.Data
				note(within, writing{attributes: attributesOf(node), slots: slotsWrittenInto(node)})
			}
			if named := attribute(node, "data-part"); named != "" {
				note(within, writing{dataParts: []string{named}})
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			visit(child, within)
		}
	}
	visit(node, "")
	return found
}

func attributesOf(node *html.Node) []string {
	var named []string
	for _, held := range node.Attr {
		named = append(named, held.Key)
	}
	return named
}

func slotsWrittenInto(node *html.Node) []string {
	var named []string
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		switch {
		case child.Type == html.ElementNode:
			named = append(named, attribute(child, "slot"))
		case child.Type == html.TextNode && strings.TrimSpace(child.Data) != "":
			named = append(named, defaultSlot)
		}
	}
	return named
}

func isGlobal(attribute string) bool {
	switch attribute {
	case "class", "id", "hidden", "title", "role", "slot", "tabindex", "autofocus":
		return true
	}
	return strings.HasPrefix(attribute, "aria-") || strings.HasPrefix(attribute, "data-")
}

func undeclared(declared contract, written writings) []string {
	var found []string
	for name, held := range written {
		for _, part := range held.dataParts {
			if name == "" {
				found = append(found, fmt.Sprintf("the markup writes data-part=%q outside any primitive, and no part is named from outside a shadow root", part))
				continue
			}
			found = append(found, fmt.Sprintf("the markup writes data-part=%q inside <%s>, and the element draws its parts in its own shadow root", part, name))
		}
		if name == "" {
			continue
		}
		expected, known := declared.Elements[name]
		if !known {
			found = append(found, fmt.Sprintf("the markup writes <%s>, and the contract declares no such element", name))
			continue
		}
		found = append(found, undeclaredOfAnElement(name, expected, held)...)
	}
	slices.Sort(found)
	return found
}

func undeclaredOfAnElement(name string, declared element, held writing) []string {
	var found []string
	for _, attribute := range held.attributes {
		if !isGlobal(attribute) && !slices.Contains(declared.Attributes, attribute) {
			found = append(found, fmt.Sprintf("the markup writes the attribute %q on <%s>, and the element does not declare it", attribute, name))
		}
	}
	for _, slot := range held.slots {
		switch {
		case slot == defaultSlot && !slices.Contains(declared.Slots, defaultSlot):
			found = append(found, fmt.Sprintf("the markup writes content into <%s> outside any named slot, and the element renders no default slot", name))
		case slot != defaultSlot && !slices.Contains(declared.Slots, slot):
			found = append(found, fmt.Sprintf("the markup writes the slot %q of <%s>, and the element does not render it", slot, name))
		}
	}
	return found
}

func unwritten(declared contract, written writings) []string {
	var found []string
	for name, expected := range declared.Elements {
		held, known := written[name]
		if !known {
			found = append(found, fmt.Sprintf("the contract declares <%s>, and no partial writes it", name))
			continue
		}
		for _, attribute := range expected.Attributes {
			if !slices.Contains(held.attributes, attribute) {
				found = append(found, fmt.Sprintf("the contract declares the attribute %q of <%s>, and no partial writes it", attribute, name))
			}
		}
		for _, slot := range expected.Slots {
			if !slices.Contains(held.slots, slot) {
				found = append(found, fmt.Sprintf("the contract declares the slot %q of <%s>, and no partial writes it", slot, name))
			}
		}
	}
	slices.Sort(found)
	return found
}

func disagreements(declared contract, written writings) []string {
	return append(undeclared(declared, written), unwritten(declared, written)...)
}

func TestEveryPartialWritesWhatTheContractDeclaresAndNothingElse(t *testing.T) {
	for _, found := range disagreements(declaredByTheContract(t), writtenByThePartials(t)) {
		t.Error(found)
	}
}

func TestARenderedPageWritesNoElementAndNoPartTheContractDoesNotDeclare(t *testing.T) {
	// scenario: A rendered page reaching the confirmation writes only what the contract declares
	page := rendered(t, oneDecision(document.One))

	written := writingsOf(page)
	if _, reached := written["ui-dialog"]; !reached {
		t.Fatalf("the rendered page writes %v, and it reaches the confirmation through the dialog's partial", written)
	}
	for _, found := range undeclared(declaredByTheContract(t), written) {
		t.Error(found)
	}
}

func TestARenderedPageWritesTheBadgeOnlyWithTheAttributesAndTheSlotItDeclares(t *testing.T) {
	// scenario: A rendered page carrying a severity and an audience writes each mark as a badge with its tone, its icon and its slotted label
	page := rendered(t, oneDecision(document.One))

	written := writingsOf(page)
	if _, reached := written["ui-badge"]; !reached {
		t.Fatalf("the rendered page writes %v, and it draws its marks as badges", written)
	}
	for _, found := range undeclared(declaredByTheContract(t), written) {
		t.Error(found)
	}
}

func TestTheContractDeclaresTheBadgeAsAnElementOfAttributesAndSlots(t *testing.T) {
	declared := declaredByTheContract(t)

	want := element{Attributes: []string{"icon", "tone", "variant"}, Slots: []string{"label"}}
	if diff := cmp.Diff(want, declared.Elements["ui-badge"]); diff != "" {
		t.Errorf("the contract declares the badge differently (-want +got):\n%s", diff)
	}
}

func TestTheContractDeclaresTheControlsAsElementsOfAttributesAndSlots(t *testing.T) {
	declared := declaredByTheContract(t)

	want := map[string]element{
		"ui-button":    {Attributes: []string{"disabled", "raised", "type", "variant"}, Slots: []string{defaultSlot}},
		"ui-field":     {Attributes: []string{"disabled", "embedded", "name", "placeholder", "rows"}, Slots: []string{"caption"}},
		"ui-copy":      {Attributes: []string{"copied-title", "failed-title", "hint", "label", "text"}},
		"ui-expand":    {Attributes: []string{"expanded", "hint", "label"}},
		"ui-file-tray": {Attributes: []string{"attach-label", "disabled", "max", "remove-label"}, Slots: []string{defaultSlot, "refusal"}},
	}
	for name, expected := range want {
		got := declared.Elements[name]
		slices.Sort(expected.Attributes)
		if diff := cmp.Diff(expected, got, cmpopts.EquateEmpty()); diff != "" {
			t.Errorf("the contract declares <%s> differently (-want +got):\n%s", name, diff)
		}
	}
}

func TestTheContractDeclaresTheCardAndTheDialogAsElementsOfSlotsAlone(t *testing.T) {
	declared := declaredByTheContract(t)

	wants := map[string]element{
		"ui-card":   {Slots: []string{defaultSlot, "body", "footer"}},
		"ui-dialog": {Slots: []string{"actions", "body", "heading"}},
	}
	for name, want := range wants {
		if diff := cmp.Diff(want, declared.Elements[name], cmpopts.EquateEmpty()); diff != "" {
			t.Errorf("the contract declares <%s> differently (-want +got):\n%s", name, diff)
		}
	}
}

func TestARenderedPageWritesTheCardAndTheDialogOnlyWithTheSlotsTheyDeclare(t *testing.T) {
	// scenario: A rendered page writes each decision, the composer and the confirmation as a card or a dialog holding only the slots it renders
	read := oneDecision(document.One)
	written := writingsOf(rendered(t, read))

	for _, name := range []string{"ui-card", "ui-dialog"} {
		_, reached := written[name]
		if !reached {
			t.Fatalf("the rendered page writes %v, and it draws a decision and its confirmation on <%s>", written, name)
		}
	}
	for _, found := range undeclared(declaredByTheContract(t), written) {
		t.Error(found)
	}
	if !slices.Contains(written["ui-card"].slots, "body") || !slices.Contains(written["ui-card"].slots, "footer") {
		t.Errorf("a decision fills the card's slots %v, want its body and its footer", written["ui-card"].slots)
	}
	if want := []string{"heading", "body", "actions"}; !slices.Equal(written["ui-dialog"].slots, want) {
		t.Errorf("the confirmation fills the dialog's slots %v, want %v", written["ui-dialog"].slots, want)
	}
}

func TestTheContractDeclaresTheAppShellAsAnElementOfSlotsAlone(t *testing.T) {
	declared := declaredByTheContract(t)

	want := element{Slots: []string{"content", "navigation", "panel", "topbar"}}
	if diff := cmp.Diff(want, declared.Elements["ui-app-shell"], cmpopts.EquateEmpty()); diff != "" {
		t.Errorf("the contract declares the app shell differently (-want +got):\n%s", diff)
	}
}

func TestARenderedPageWritesTheAppShellOnlyWithTheSlotsItDeclares(t *testing.T) {
	// scenario: A rendered page with an index and decisions fills the shell's topbar, navigation, content and panel slots and writes no part
	written := writingsOf(rendered(t, oneDecision(document.One)))

	held, reached := written["ui-app-shell"]
	if !reached {
		t.Fatalf("the rendered page writes %v, and it draws its frame on <ui-app-shell>", written)
	}
	slices.Sort(held.slots)
	if want := []string{"content", "navigation", "panel", "topbar"}; !slices.Equal(held.slots, want) {
		t.Errorf("the page fills the shell's slots %v, want %v", held.slots, want)
	}
	for _, found := range undeclared(declaredByTheContract(t), written) {
		t.Error(found)
	}
}

func TestAnElementNoContractDeclaresIsReported(t *testing.T) {
	declared := contract{Elements: map[string]element{"ui-dialog": {Slots: []string{"body"}}}}
	written := writings{"ui-dialog": {slots: []string{"body"}}, "ui-tooltip": {}}

	want := []string{`the markup writes <ui-tooltip>, and the contract declares no such element`}
	if diff := cmp.Diff(want, disagreements(declared, written)); diff != "" {
		t.Errorf("the disagreements differ (-want +got):\n%s", diff)
	}
}

func TestAnElementNoPartialWritesIsReported(t *testing.T) {
	declared := contract{Elements: map[string]element{"ui-dialog": {Slots: []string{"body"}}, "ui-tooltip": {}}}
	written := writings{"ui-dialog": {slots: []string{"body"}}}

	want := []string{`the contract declares <ui-tooltip>, and no partial writes it`}
	if diff := cmp.Diff(want, disagreements(declared, written)); diff != "" {
		t.Errorf("the disagreements differ (-want +got):\n%s", diff)
	}
}

func TestAPartialWritingDataPartIsReportedNamingTheElementAndThePart(t *testing.T) {
	parsed, err := html.Parse(strings.NewReader(`<ui-dialog><div slot="body" data-part="panel">words</div></ui-dialog>`))
	if err != nil {
		t.Fatalf("reading the markup: %v", err)
	}
	declared := contract{Elements: map[string]element{"ui-dialog": {Slots: []string{"body"}}}}

	want := []string{`the markup writes data-part="panel" inside <ui-dialog>, and the element draws its parts in its own shadow root`}
	if diff := cmp.Diff(want, disagreements(declared, writingsOf(parsed))); diff != "" {
		t.Errorf("the disagreements differ (-want +got):\n%s", diff)
	}
}

func TestADataPartWrittenOutsideAnyPrimitiveIsReported(t *testing.T) {
	parsed, err := html.Parse(strings.NewReader(`<div data-part="control"></div>`))
	if err != nil {
		t.Fatalf("reading the markup: %v", err)
	}

	want := []string{`the markup writes data-part="control" outside any primitive, and no part is named from outside a shadow root`}
	if diff := cmp.Diff(want, undeclared(contract{}, writingsOf(parsed))); diff != "" {
		t.Errorf("the disagreements differ (-want +got):\n%s", diff)
	}
}

func TestNoPartialAndNoRenderedPageWritesDataPart(t *testing.T) {
	// scenario: Every partial and a rendered page carrying code, a patch and an illustration write no data-part attribute
	read := oneDecision(document.One)
	read.Decisions[0].Evidence = document.Evidence{{Code: "x := 1", Language: "go"}, {Patch: "--- a/x\n+++ b/x\n@@ -1,1 +1,1 @@\n-a\n+b\n"}}
	read.Decisions[0].Illustration = document.Illustration{Path: "drawings/c-1.html"}

	for place, written := range map[string]writings{"the partials": writtenByThePartials(t), "the rendered page": writingsOf(rendered(t, read)), "the catalog": writingsOf(theCatalog(t))} {
		for name, held := range written {
			for _, part := range held.dataParts {
				t.Errorf("%s write data-part=%q on or inside <%s>", place, part, name)
			}
		}
	}
}

func TestTheContractCarriesElementsAndNoPartsBesideThem(t *testing.T) {
	held, err := fs.ReadFile(assets.Files(), "contract.json")
	if err != nil {
		t.Fatalf("reading the contract the build writes: %v", err)
	}
	var sections map[string]json.RawMessage
	if err := json.Unmarshal(held, &sections); err != nil {
		t.Fatalf("parsing the contract: %v", err)
	}

	var names []string
	for name := range sections {
		names = append(names, name)
	}
	if diff := cmp.Diff([]string{"elements"}, names); diff != "" {
		t.Errorf("the contract's sections differ (-want +got):\n%s", diff)
	}
}

func aBadge() contract {
	return contract{Elements: map[string]element{"ui-badge": {Attributes: []string{"icon", "tone"}, Slots: []string{"label"}}}}
}

func TestAnAttributeABadgeDoesNotDeclareIsReportedNamingTheBadgeAndTheAttribute(t *testing.T) {
	written := writings{"ui-badge": {attributes: []string{"tone", "icon", "colour"}, slots: []string{"label"}}}

	want := []string{`the markup writes the attribute "colour" on <ui-badge>, and the element does not declare it`}
	if diff := cmp.Diff(want, disagreements(aBadge(), written)); diff != "" {
		t.Errorf("the disagreements differ (-want +got):\n%s", diff)
	}
}

func TestASlotABadgeDoesNotRenderIsReportedNamingTheBadgeAndTheSlot(t *testing.T) {
	written := writings{"ui-badge": {attributes: []string{"tone", "icon"}, slots: []string{"label", "trailing"}}}

	want := []string{`the markup writes the slot "trailing" of <ui-badge>, and the element does not render it`}
	if diff := cmp.Diff(want, disagreements(aBadge(), written)); diff != "" {
		t.Errorf("the disagreements differ (-want +got):\n%s", diff)
	}
}

func TestContentWrittenIntoABadgeOutsideAnySlotIsReportedNamingTheBadge(t *testing.T) {
	written := writings{"ui-badge": {attributes: []string{"tone", "icon"}, slots: []string{"label", defaultSlot}}}

	want := []string{`the markup writes content into <ui-badge> outside any named slot, and the element renders no default slot`}
	if diff := cmp.Diff(want, disagreements(aBadge(), written)); diff != "" {
		t.Errorf("the disagreements differ (-want +got):\n%s", diff)
	}
}

func TestAnAttributeABadgeDeclaresAndNoPartialWritesIsReportedNamingTheBadgeAndTheAttribute(t *testing.T) {
	written := writings{"ui-badge": {attributes: []string{"tone"}, slots: []string{"label"}}}

	want := []string{`the contract declares the attribute "icon" of <ui-badge>, and no partial writes it`}
	if diff := cmp.Diff(want, disagreements(aBadge(), written)); diff != "" {
		t.Errorf("the disagreements differ (-want +got):\n%s", diff)
	}
}

func TestASlotABadgeRendersAndNoPartialWritesIsReportedNamingTheBadgeAndTheSlot(t *testing.T) {
	written := writings{"ui-badge": {attributes: []string{"tone", "icon"}}}

	want := []string{`the contract declares the slot "label" of <ui-badge>, and no partial writes it`}
	if diff := cmp.Diff(want, disagreements(aBadge(), written)); diff != "" {
		t.Errorf("the disagreements differ (-want +got):\n%s", diff)
	}
}

func TestAnElementTheContractDeclaresWithAttributesAndSlotsThatNoPartialWritesIsReported(t *testing.T) {
	want := []string{`the contract declares <ui-badge>, and no partial writes it`}
	if diff := cmp.Diff(want, disagreements(aBadge(), writings{})); diff != "" {
		t.Errorf("the disagreements differ (-want +got):\n%s", diff)
	}
}

func TestTheGlobalAttributesAreNotPartOfTheContract(t *testing.T) {
	written := writings{"ui-badge": {
		attributes: []string{"tone", "icon", "class", "id", "hidden", "title", "role", "slot", "tabindex", "autofocus", "aria-label", "data-title", "data-hint"},
		slots:      []string{"label"},
	}}

	if diff := cmp.Diff([]string(nil), disagreements(aBadge(), written)); diff != "" {
		t.Errorf("a global attribute was held to the contract (-want +got):\n%s", diff)
	}
}

func TestWhatAPartialWritesIntoAnElementIsReadAsItsAttributesAndTheSlotsOfItsChildren(t *testing.T) {
	parsed, err := html.Parse(strings.NewReader(`<ui-badge tone="error" data-title="t"><span slot="label">word</span>loose<svg></svg></ui-badge>`))
	if err != nil {
		t.Fatalf("reading the markup: %v", err)
	}

	got := writingsOf(parsed)["ui-badge"]

	if diff := cmp.Diff([]string{"tone", "data-title"}, got.attributes); diff != "" {
		t.Errorf("the attributes differ (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]string{"label", defaultSlot, defaultSlot}, got.slots); diff != "" {
		t.Errorf("the slots differ (-want +got):\n%s", diff)
	}
}

func TestARenderedPageReachesEveryControlThroughThePartialThatDrawsIt(t *testing.T) {
	// scenario: A page with options, a code block and an illustration draws the button, field, choice, copy and expand through their partials
	read := oneDecision(document.One)
	read.Decisions[0].Evidence = document.Evidence{{Code: "x := 1", Language: "go"}}
	read.Decisions[0].Illustration = document.Illustration{Path: "drawings/c-1.html"}

	written := writingsOf(rendered(t, read))

	for _, element := range []string{"ui-button", "ui-field", "ui-choice", "ui-copy", "ui-expand"} {
		if _, reached := written[element]; !reached {
			t.Errorf("the rendered page writes %v, and it reaches the controls through <%s>", written, element)
		}
	}
	for _, found := range undeclared(declaredByTheContract(t), written) {
		t.Error(found)
	}
}

func TestNoButtonAndNoTextareaOfARenderedPageStandsOutsideAPrimitivesShadowRoot(t *testing.T) {
	// scenario: A rendered page writes no button and no textarea of its own, since the button and the field draw theirs inside their elements
	read := oneDecision(document.One)
	read.Decisions[0].Evidence = document.Evidence{{Code: "x := 1", Language: "go"}}
	read.Decisions[0].Illustration = document.Illustration{Path: "drawings/c-1.html"}

	page := rendered(t, read)
	for _, control := range elements(page, "button", "textarea") {
		t.Errorf("the page writes a <%s class=%q>, and the button, the field, the copy control and the expand control draw theirs inside their own shadow roots", control.Data, attribute(control, "class"))
	}
	for _, control := range elements(page, "input") {
		t.Errorf("the page writes an <input class=%q> of its own, and the controls are drawn inside their primitives' shadow roots", attribute(control, "class"))
	}
}

var indicatorsTheBrandAndTheTooltip = map[string]element{
	"ui-status-dot":   {Attributes: []string{"tone"}},
	"ui-overline":     {Attributes: []string{"count"}, Slots: []string{defaultSlot}},
	"ui-progress-bar": {Attributes: []string{"label", "max", "value"}},
	"ui-brand":        {Slots: []string{defaultSlot, "mark"}},
	"ui-tooltip":      {},
}

func TestTheContractDeclaresEachIndicatorTheBrandAndTheTooltipAsAnElementOfAttributesAndSlots(t *testing.T) {
	declared := declaredByTheContract(t)

	for name, want := range indicatorsTheBrandAndTheTooltip {
		t.Run(name, func(t *testing.T) {
			if diff := cmp.Diff(want, declared.Elements[name], cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("the contract declares <%s> differently (-want +got):\n%s", name, diff)
			}
		})
	}
}

func TestARenderedPageWritesTheIndicatorsTheBrandAndTheTooltipOnlyWithTheAttributesAndSlotsTheyDeclare(t *testing.T) {
	// scenario: A rendered page writes its status dots, overlines, progress bar, brand and tooltip with only declared attributes and slotted content
	read := oneDecision(document.One)
	written := writingsOf(rendered(t, read))

	for name := range indicatorsTheBrandAndTheTooltip {
		if _, reached := written[name]; !reached {
			t.Errorf("the rendered page writes %v, and it draws <%s>", written, name)
		}
	}
	for _, found := range undeclared(declaredByTheContract(t), written) {
		t.Error(found)
	}
}

func TestEveryIndicatorTheBrandAndTheTooltipWrittenWithSomethingItDoesNotDeclareIsReported(t *testing.T) {
	declared := contract{Elements: indicatorsTheBrandAndTheTooltip}
	for _, mutation := range []struct {
		scenario string
		written  writing
		element  string
		want     string
	}{
		{"a status dot given a colour", writing{attributes: []string{"tone", "colour"}}, "ui-status-dot", `the markup writes the attribute "colour" on <ui-status-dot>, and the element does not declare it`},
		{"an overline given a total", writing{attributes: []string{"count", "total"}, slots: []string{defaultSlot}}, "ui-overline", `the markup writes the attribute "total" on <ui-overline>, and the element does not declare it`},
		{"an overline given a slot it does not render", writing{attributes: []string{"count"}, slots: []string{defaultSlot, "suffix"}}, "ui-overline", `the markup writes the slot "suffix" of <ui-overline>, and the element does not render it`},
		{"a progress bar given a percentage", writing{attributes: []string{"label", "max", "value", "percent"}}, "ui-progress-bar", `the markup writes the attribute "percent" on <ui-progress-bar>, and the element does not declare it`},
		{"a progress bar given content", writing{attributes: []string{"label", "max", "value"}, slots: []string{defaultSlot}}, "ui-progress-bar", `the markup writes content into <ui-progress-bar> outside any named slot, and the element renders no default slot`},
		{"a brand given a slot it does not render", writing{slots: []string{defaultSlot, "mark", "tagline"}}, "ui-brand", `the markup writes the slot "tagline" of <ui-brand>, and the element does not render it`},
		{"a brand given an attribute", writing{attributes: []string{"mark"}, slots: []string{defaultSlot, "mark"}}, "ui-brand", `the markup writes the attribute "mark" on <ui-brand>, and the element does not declare it`},
		{"a tooltip given a title", writing{attributes: []string{"heading"}}, "ui-tooltip", `the markup writes the attribute "heading" on <ui-tooltip>, and the element does not declare it`},
		{"a tooltip given its parts", writing{dataParts: []string{"title"}}, "ui-tooltip", `the markup writes data-part="title" inside <ui-tooltip>, and the element draws its parts in its own shadow root`},
	} {
		t.Run(mutation.scenario, func(t *testing.T) {
			written := writings{}
			for name, held := range indicatorsTheBrandAndTheTooltip {
				written[name] = writing{attributes: held.Attributes, slots: held.Slots}
			}
			written[mutation.element] = mutation.written

			if diff := cmp.Diff([]string{mutation.want}, disagreements(declared, written)); diff != "" {
				t.Errorf("the disagreements differ (-want +got):\n%s", diff)
			}
		})
	}
}

func TestEveryAttributeAndSlotAnIndicatorTheBrandOrTheTooltipDeclaresAndNoPartialWritesIsReported(t *testing.T) {
	declared := contract{Elements: indicatorsTheBrandAndTheTooltip}
	for _, mutation := range []struct {
		scenario string
		element  string
		written  writing
		want     string
	}{
		{"a status dot whose tone nobody writes", "ui-status-dot", writing{}, `the contract declares the attribute "tone" of <ui-status-dot>, and no partial writes it`},
		{"an overline whose count nobody writes", "ui-overline", writing{slots: []string{defaultSlot}}, `the contract declares the attribute "count" of <ui-overline>, and no partial writes it`},
		{"an overline whose words nobody slots", "ui-overline", writing{attributes: []string{"count"}}, `the contract declares the slot "" of <ui-overline>, and no partial writes it`},
		{"a progress bar whose label nobody writes", "ui-progress-bar", writing{attributes: []string{"max", "value"}}, `the contract declares the attribute "label" of <ui-progress-bar>, and no partial writes it`},
		{"a brand whose mark nobody slots", "ui-brand", writing{slots: []string{defaultSlot}}, `the contract declares the slot "mark" of <ui-brand>, and no partial writes it`},
	} {
		t.Run(mutation.scenario, func(t *testing.T) {
			written := writings{}
			for name, held := range indicatorsTheBrandAndTheTooltip {
				written[name] = writing{attributes: held.Attributes, slots: held.Slots}
			}
			written[mutation.element] = mutation.written

			if diff := cmp.Diff([]string{mutation.want}, disagreements(declared, written)); diff != "" {
				t.Errorf("the disagreements differ (-want +got):\n%s", diff)
			}
		})
	}
}

func aCardAndADialog() contract {
	return contract{Elements: map[string]element{
		"ui-card":   {Slots: []string{defaultSlot, "body", "footer"}},
		"ui-dialog": {Slots: []string{"heading", "body", "actions"}},
	}}
}

func TestASlotACardDoesNotRenderIsReportedNamingTheCardAndTheSlot(t *testing.T) {
	written := writings{
		"ui-card":   {slots: []string{"body", "footer", defaultSlot, "aside"}},
		"ui-dialog": {slots: []string{"heading", "body", "actions"}},
	}

	want := []string{`the markup writes the slot "aside" of <ui-card>, and the element does not render it`}
	if diff := cmp.Diff(want, disagreements(aCardAndADialog(), written)); diff != "" {
		t.Errorf("the disagreements differ (-want +got):\n%s", diff)
	}
}

func TestAnAttributeACardDoesNotDeclareIsReportedNamingTheCardAndTheAttribute(t *testing.T) {
	written := writings{
		"ui-card":   {attributes: []string{"class", "padding"}, slots: []string{"body", "footer", defaultSlot}},
		"ui-dialog": {slots: []string{"heading", "body", "actions"}},
	}

	want := []string{`the markup writes the attribute "padding" on <ui-card>, and the element does not declare it`}
	if diff := cmp.Diff(want, disagreements(aCardAndADialog(), written)); diff != "" {
		t.Errorf("the disagreements differ (-want +got):\n%s", diff)
	}
}

func TestASlotADialogRendersAndNoPartialWritesIsReportedNamingTheDialogAndTheSlot(t *testing.T) {
	written := writings{
		"ui-card":   {slots: []string{"body", "footer", defaultSlot}},
		"ui-dialog": {slots: []string{"heading", "body"}},
	}

	want := []string{`the contract declares the slot "actions" of <ui-dialog>, and no partial writes it`}
	if diff := cmp.Diff(want, disagreements(aCardAndADialog(), written)); diff != "" {
		t.Errorf("the disagreements differ (-want +got):\n%s", diff)
	}
}

func TestAnAttributeADialogDoesNotDeclareIsReportedNamingTheDialogAndTheAttribute(t *testing.T) {
	written := writings{
		"ui-card":   {slots: []string{"body", "footer", defaultSlot}},
		"ui-dialog": {attributes: []string{"open"}, slots: []string{"heading", "body", "actions"}},
	}

	want := []string{`the markup writes the attribute "open" on <ui-dialog>, and the element does not declare it`}
	if diff := cmp.Diff(want, disagreements(aCardAndADialog(), written)); diff != "" {
		t.Errorf("the disagreements differ (-want +got):\n%s", diff)
	}
}

func TestContentWrittenIntoADialogOutsideAnySlotIsReportedNamingTheDialog(t *testing.T) {
	written := writings{
		"ui-card":   {slots: []string{"body", "footer", defaultSlot}},
		"ui-dialog": {slots: []string{"heading", "body", "actions", defaultSlot}},
	}

	want := []string{`the markup writes content into <ui-dialog> outside any named slot, and the element renders no default slot`}
	if diff := cmp.Diff(want, disagreements(aCardAndADialog(), written)); diff != "" {
		t.Errorf("the disagreements differ (-want +got):\n%s", diff)
	}
}

var blocksOfEvidence = map[string]element{
	"ui-code-block": {Attributes: []string{"language"}, Slots: []string{defaultSlot}},
	"ui-diff-block": {Attributes: []string{"language"}, Slots: []string{defaultSlot}},
}

func TestTheContractDeclaresTheCodeBlockAndTheDiffBlockAsElementsOfALanguageAndTheContentTheyHold(t *testing.T) {
	declared := declaredByTheContract(t)

	for name, want := range blocksOfEvidence {
		t.Run(name, func(t *testing.T) {
			if diff := cmp.Diff(want, declared.Elements[name], cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("the contract declares <%s> differently (-want +got):\n%s", name, diff)
			}
		})
	}
}

func TestARenderedPageWritesTheCodeBlockAndTheDiffBlockOnlyWithTheLanguageAndTheContentTheyHold(t *testing.T) {
	// scenario: A decision carrying code and a patch writes a code block holding its code and a diff block holding its tables, each naming a language and nothing inside
	read := oneDecision(document.One)
	read.Decisions[0].Evidence = document.Evidence{
		{Code: "x := 1\n", Language: "go"},
		{Patch: "--- a/internal/cache/write.go\n+++ b/internal/cache/write.go\n@@ -1,1 +1,1 @@\n-a\n+b\n"},
	}

	written := writingsOf(rendered(t, read))

	for name := range blocksOfEvidence {
		held, reached := written[name]
		if !reached {
			t.Fatalf("the rendered page writes %v, and it draws evidence on <%s>", written, name)
		}
		if !slices.Contains(held.attributes, "language") {
			t.Errorf("the rendered page writes <%s> with the attributes %v, want it to name its language", name, held.attributes)
		}
	}
	for _, found := range undeclared(declaredByTheContract(t), written) {
		t.Error(found)
	}
}

func TestEveryCodeBlockAndDiffBlockWrittenWithSomethingItDoesNotDeclareIsReported(t *testing.T) {
	declared := contract{Elements: blocksOfEvidence}
	for _, mutation := range []struct {
		scenario string
		element  string
		written  writing
		want     string
	}{
		{"a code block given a highlight switch", "ui-code-block", writing{attributes: []string{"language", "highlight"}, slots: []string{defaultSlot}}, `the markup writes the attribute "highlight" on <ui-code-block>, and the element does not declare it`},
		{"a code block given a slot it does not render", "ui-code-block", writing{attributes: []string{"language"}, slots: []string{defaultSlot, "gutter"}}, `the markup writes the slot "gutter" of <ui-code-block>, and the element does not render it`},
		{"a code block given its code part", "ui-code-block", writing{attributes: []string{"language"}, slots: []string{defaultSlot}, dataParts: []string{"code"}}, `the markup writes data-part="code" inside <ui-code-block>, and the element draws its parts in its own shadow root`},
		{"a diff block given a layout", "ui-diff-block", writing{attributes: []string{"language", "layout"}, slots: []string{defaultSlot}}, `the markup writes the attribute "layout" on <ui-diff-block>, and the element does not declare it`},
		{"a diff block given a hunk part", "ui-diff-block", writing{attributes: []string{"language"}, slots: []string{defaultSlot}, dataParts: []string{"hunk"}}, `the markup writes data-part="hunk" inside <ui-diff-block>, and the element draws its parts in its own shadow root`},
		{"a code block whose language nobody writes", "ui-code-block", writing{slots: []string{defaultSlot}}, `the contract declares the attribute "language" of <ui-code-block>, and no partial writes it`},
		{"a diff block whose tables nobody writes", "ui-diff-block", writing{attributes: []string{"language"}}, `the contract declares the slot "" of <ui-diff-block>, and no partial writes it`},
	} {
		t.Run(mutation.scenario, func(t *testing.T) {
			written := writings{}
			for name, held := range blocksOfEvidence {
				written[name] = writing{attributes: held.Attributes, slots: held.Slots}
			}
			written[mutation.element] = mutation.written

			if diff := cmp.Diff([]string{mutation.want}, disagreements(declared, written)); diff != "" {
				t.Errorf("the disagreements differ (-want +got):\n%s", diff)
			}
		})
	}
}
