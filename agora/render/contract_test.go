package render_test

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/agora/assets"
	"github.com/socrates-tools/socrates/agora/document"
	"github.com/socrates-tools/socrates/agora/render"
	"golang.org/x/net/html"
)

const primitivePrefix = "ui-"

type parts map[string][]string

func declaredByTheContract(t *testing.T) parts {
	t.Helper()

	held, err := fs.ReadFile(assets.Files(), "contract.json")
	if err != nil {
		t.Fatalf("reading the contract the build writes: %v", err)
	}
	var declared parts
	if err := json.Unmarshal(held, &declared); err != nil {
		t.Fatalf("parsing the contract: %v", err)
	}
	return declared
}

func writtenByThePartials(t *testing.T) parts {
	t.Helper()

	drawn, err := render.Primitives()
	if err != nil {
		t.Fatalf("drawing the partials: %v", err)
	}
	written := parts{}
	for name, markup := range drawn {
		parsed, err := html.Parse(strings.NewReader(markup))
		if err != nil {
			t.Fatalf("reading back the partial %s: %v", name, err)
		}
		for key, held := range partsOf(parsed) {
			written[key] = append(written[key], held...)
		}
	}
	return written
}

func partsOf(node *html.Node) parts {
	found := parts{}
	var visit func(node *html.Node, within string)
	visit = func(node *html.Node, within string) {
		if node.Type == html.ElementNode {
			if node.Data == "iframe" {
				if framed, err := html.Parse(strings.NewReader(attribute(node, "srcdoc"))); err == nil {
					for key, held := range partsOf(framed) {
						found[key] = append(found[key], held...)
					}
				}
			}
			if strings.HasPrefix(node.Data, primitivePrefix) {
				within = node.Data
				if _, seen := found[within]; !seen {
					found[within] = nil
				}
			}
			if named := attribute(node, "data-part"); named != "" {
				found[within] = append(found[within], named)
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			visit(child, within)
		}
	}
	visit(node, "")
	return found
}

func undeclared(declared, written parts) []string {
	var found []string
	for element, held := range written {
		if element == "" {
			for _, part := range held {
				found = append(found, fmt.Sprintf("the markup writes the part %q outside any primitive, and a part belongs to one", part))
			}
			continue
		}
		expected, known := declared[element]
		if !known {
			found = append(found, fmt.Sprintf("the markup writes <%s>, and the contract declares no such element", element))
			continue
		}
		for _, part := range held {
			if !slices.Contains(expected, part) {
				found = append(found, fmt.Sprintf("the markup writes the part %q of <%s>, and the contract does not declare it", part, element))
			}
		}
	}
	slices.Sort(found)
	return found
}

func unwritten(declared, written parts) []string {
	var found []string
	for element, expected := range declared {
		held, known := written[element]
		if !known {
			found = append(found, fmt.Sprintf("the contract declares <%s>, and no partial writes it", element))
			continue
		}
		for _, part := range expected {
			if !slices.Contains(held, part) {
				found = append(found, fmt.Sprintf("the contract declares the part %q of <%s>, and no partial writes it", part, element))
			}
		}
	}
	slices.Sort(found)
	return found
}

func disagreements(declared, written parts) []string {
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

	written := partsOf(page)
	if _, reached := written["ui-dialog"]; !reached {
		t.Fatalf("the rendered page writes %v, and it reaches the confirmation through the dialog's partial", written)
	}
	for _, found := range undeclared(declaredByTheContract(t), written) {
		t.Error(found)
	}
}

func TestTheContractDeclaresTheDialog(t *testing.T) {
	declared := declaredByTheContract(t)
	if _, known := declared["ui-dialog"]; !known {
		t.Errorf("the contract declares %v, and the dialog is the first primitive it declares", declared)
	}
}

func TestAnElementNoContractDeclaresIsReported(t *testing.T) {
	declared := parts{"ui-dialog": {"panel"}}
	written := parts{"ui-dialog": {"panel"}, "ui-tooltip": nil}

	want := []string{`the markup writes <ui-tooltip>, and the contract declares no such element`}
	if diff := cmp.Diff(want, disagreements(declared, written)); diff != "" {
		t.Errorf("the disagreements differ (-want +got):\n%s", diff)
	}
}

func TestAPartNoContractDeclaresIsReportedNamingTheElementAndThePart(t *testing.T) {
	declared := parts{"ui-dialog": {"panel"}}
	written := parts{"ui-dialog": {"panel", "footer"}}

	want := []string{`the markup writes the part "footer" of <ui-dialog>, and the contract does not declare it`}
	if diff := cmp.Diff(want, disagreements(declared, written)); diff != "" {
		t.Errorf("the disagreements differ (-want +got):\n%s", diff)
	}
}

func TestAPartNoPartialWritesIsReportedNamingTheElementAndThePart(t *testing.T) {
	declared := parts{"ui-dialog": {"panel", "heading"}}
	written := parts{"ui-dialog": {"panel"}}

	want := []string{`the contract declares the part "heading" of <ui-dialog>, and no partial writes it`}
	if diff := cmp.Diff(want, disagreements(declared, written)); diff != "" {
		t.Errorf("the disagreements differ (-want +got):\n%s", diff)
	}
}

func TestAnElementNoPartialWritesIsReported(t *testing.T) {
	declared := parts{"ui-dialog": {"panel"}, "ui-tooltip": {"bubble"}}
	written := parts{"ui-dialog": {"panel"}}

	want := []string{`the contract declares <ui-tooltip>, and no partial writes it`}
	if diff := cmp.Diff(want, disagreements(declared, written)); diff != "" {
		t.Errorf("the disagreements differ (-want +got):\n%s", diff)
	}
}

func TestARenderedPageReachesEveryControlThroughThePartialThatDrawsIt(t *testing.T) {
	// scenario: A page with options, a code block and an illustration draws the button, field, choice, copy and expand through their partials
	read := oneDecision(document.One)
	read.Decisions[0].Evidence = document.Evidence{{Code: "x := 1", Language: "go"}}
	read.Decisions[0].Illustration = document.Illustration{Path: "drawings/c-1.html"}

	written := partsOf(rendered(t, read))

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
