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
	Elements map[string]element  `json:"elements"`
	Parts    map[string][]string `json:"parts"`
}

type writing struct {
	attributes []string
	slots      []string
	parts      []string
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
		parts:      append(first.parts, second.parts...),
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
				note(within, writing{parts: []string{named}})
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
		if name == "" {
			for _, part := range held.parts {
				found = append(found, fmt.Sprintf("the markup writes the part %q outside any primitive, and a part belongs to one", part))
			}
			continue
		}
		moved, isMoved := declared.Elements[name]
		expected, hasParts := declared.Parts[name]
		if !isMoved && !hasParts {
			found = append(found, fmt.Sprintf("the markup writes <%s>, and the contract declares no such element", name))
			continue
		}
		if isMoved {
			found = append(found, undeclaredOfAnElement(name, moved, held)...)
			continue
		}
		for _, part := range held.parts {
			if !slices.Contains(expected, part) {
				found = append(found, fmt.Sprintf("the markup writes the part %q of <%s>, and the contract does not declare it", part, name))
			}
		}
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
	for _, part := range held.parts {
		found = append(found, fmt.Sprintf("the markup writes the part %q of <%s>, and the element draws its own parts", part, name))
	}
	return found
}

func unwritten(declared contract, written writings) []string {
	var found []string
	for name, expected := range declared.Parts {
		held, known := written[name]
		if !known {
			found = append(found, fmt.Sprintf("the contract declares <%s>, and no partial writes it", name))
			continue
		}
		for _, part := range expected {
			if !slices.Contains(held.parts, part) {
				found = append(found, fmt.Sprintf("the contract declares the part %q of <%s>, and no partial writes it", part, name))
			}
		}
	}
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

	want := element{Attributes: []string{"icon", "tone"}, Slots: []string{"label"}}
	if diff := cmp.Diff(want, declared.Elements["ui-badge"]); diff != "" {
		t.Errorf("the contract declares the badge differently (-want +got):\n%s", diff)
	}
	if _, stillParts := declared.Parts["ui-badge"]; stillParts {
		t.Errorf("the contract still declares parts for the badge: %v", declared.Parts["ui-badge"])
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
		if _, stillParts := declared.Parts[name]; stillParts {
			t.Errorf("the contract still declares parts for <%s>: %v", name, declared.Parts[name])
		}
	}
}

func TestARenderedPageWritesTheCardAndTheDialogOnlyWithTheSlotsTheyDeclare(t *testing.T) {
	// scenario: A rendered page writes each decision, the composer and the confirmation as a card or a dialog holding only the slots it renders
	read := oneDecision(document.One)
	written := writingsOf(rendered(t, read))

	for _, name := range []string{"ui-card", "ui-dialog"} {
		held, reached := written[name]
		if !reached {
			t.Fatalf("the rendered page writes %v, and it draws a decision and its confirmation on <%s>", written, name)
		}
		if len(held.parts) != 0 {
			t.Errorf("the rendered page writes the parts %v inside <%s>, and the element draws its own parts", held.parts, name)
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
	if _, stillParts := declared.Parts["ui-app-shell"]; stillParts {
		t.Errorf("the contract still declares parts for the app shell: %v", declared.Parts["ui-app-shell"])
	}
}

func TestARenderedPageWritesTheAppShellOnlyWithTheSlotsItDeclares(t *testing.T) {
	// scenario: A rendered page with an index and decisions fills the shell's topbar, navigation, content and panel slots and writes no part
	written := writingsOf(rendered(t, oneDecision(document.One)))

	held, reached := written["ui-app-shell"]
	if !reached {
		t.Fatalf("the rendered page writes %v, and it draws its frame on <ui-app-shell>", written)
	}
	if len(held.parts) != 0 {
		t.Errorf("the rendered page writes the parts %v inside <ui-app-shell>, and the element draws its own parts", held.parts)
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
	declared := contract{Parts: map[string][]string{"ui-dialog": {"panel"}}}
	written := writings{"ui-dialog": {parts: []string{"panel"}}, "ui-tooltip": {}}

	want := []string{`the markup writes <ui-tooltip>, and the contract declares no such element`}
	if diff := cmp.Diff(want, disagreements(declared, written)); diff != "" {
		t.Errorf("the disagreements differ (-want +got):\n%s", diff)
	}
}

func TestAPartNoContractDeclaresIsReportedNamingTheElementAndThePart(t *testing.T) {
	declared := contract{Parts: map[string][]string{"ui-dialog": {"panel"}}}
	written := writings{"ui-dialog": {parts: []string{"panel", "footer"}}}

	want := []string{`the markup writes the part "footer" of <ui-dialog>, and the contract does not declare it`}
	if diff := cmp.Diff(want, disagreements(declared, written)); diff != "" {
		t.Errorf("the disagreements differ (-want +got):\n%s", diff)
	}
}

func TestAPartNoPartialWritesIsReportedNamingTheElementAndThePart(t *testing.T) {
	declared := contract{Parts: map[string][]string{"ui-dialog": {"panel", "heading"}}}
	written := writings{"ui-dialog": {parts: []string{"panel"}}}

	want := []string{`the contract declares the part "heading" of <ui-dialog>, and no partial writes it`}
	if diff := cmp.Diff(want, disagreements(declared, written)); diff != "" {
		t.Errorf("the disagreements differ (-want +got):\n%s", diff)
	}
}

func TestAnElementNoPartialWritesIsReported(t *testing.T) {
	declared := contract{Parts: map[string][]string{"ui-dialog": {"panel"}, "ui-tooltip": {"bubble"}}}
	written := writings{"ui-dialog": {parts: []string{"panel"}}}

	want := []string{`the contract declares <ui-tooltip>, and no partial writes it`}
	if diff := cmp.Diff(want, disagreements(declared, written)); diff != "" {
		t.Errorf("the disagreements differ (-want +got):\n%s", diff)
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

func TestAPartWrittenInsideAnElementThatDrawsItsOwnPartsIsReported(t *testing.T) {
	written := writings{"ui-badge": {attributes: []string{"tone", "icon"}, slots: []string{"label"}, parts: []string{"label"}}}

	want := []string{`the markup writes the part "label" of <ui-badge>, and the element draws its own parts`}
	if diff := cmp.Diff(want, disagreements(aBadge(), written)); diff != "" {
		t.Errorf("the disagreements differ (-want +got):\n%s", diff)
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

func TestEveryControlOfARenderedPageIsHeldByAPrimitive(t *testing.T) {
	// scenario: No button, input or textarea of a rendered page stands outside a primitive's element
	read := oneDecision(document.One)
	read.Decisions[0].Evidence = document.Evidence{{Code: "x := 1", Language: "go"}}
	read.Decisions[0].Illustration = document.Illustration{Path: "drawings/c-1.html"}

	for _, control := range elements(rendered(t, read), "button", "input", "textarea") {
		if attribute(control, "data-part") != "control" {
			t.Errorf("the page offers a <%s class=%q> that is not the control of a primitive", control.Data, attribute(control, "class"))
		}
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

func TestAPartWrittenInsideADialogIsReportedBecauseTheDialogDrawsItsOwnParts(t *testing.T) {
	written := writings{
		"ui-card":   {slots: []string{"body", "footer", defaultSlot}},
		"ui-dialog": {slots: []string{"heading", "body", "actions"}, parts: []string{"panel"}},
	}

	want := []string{`the markup writes the part "panel" of <ui-dialog>, and the element draws its own parts`}
	if diff := cmp.Diff(want, disagreements(aCardAndADialog(), written)); diff != "" {
		t.Errorf("the disagreements differ (-want +got):\n%s", diff)
	}
}
