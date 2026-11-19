package e2e_test

import (
	"slices"
	"testing"
)

const whatEveryPrimitiveOfThePageIsWrittenWith = `(async () => {
  const globals = ['class', 'id', 'hidden', 'title', 'role', 'slot', 'tabindex', 'autofocus'];
  const isGlobal = (name) => globals.includes(name) || name.startsWith('aria-') || name.startsWith('data-');
  const found = [];
  const renderedOnlyWhileSaying = { 'ui-file-tray': ['refusal'] };
  const primitives = Array.from(document.querySelectorAll('*')).filter((held) => held.localName.startsWith('ui-'));
  for (const held of primitives) {
    const tag = held.localName;
    const defined = customElements.get(tag);
    if (!defined) {
      found.push('the page writes <' + tag + '>, and no element of that name is defined');
      continue;
    }
    await held.updateComplete;
    const declaredAttributes = defined.observedAttributes || [];
    for (const name of held.getAttributeNames()) {
      if (!isGlobal(name) && !declaredAttributes.includes(name)) {
        found.push('the page writes the attribute "' + name + '" on <' + tag + '>, and the element does not declare it');
      }
    }
    const declaredSlots = [...Array.from(held.shadowRoot.querySelectorAll('slot'), (slot) => slot.name), ...(renderedOnlyWhileSaying[tag] || [])];
    for (const child of held.childNodes) {
      if (child.nodeType === Node.ELEMENT_NODE) {
        const slot = child.getAttribute('slot') || '';
        if (!declaredSlots.includes(slot)) {
          found.push('the page writes the slot "' + slot + '" of <' + tag + '>, and the element does not render it');
        }
      } else if (child.nodeType === Node.TEXT_NODE && child.textContent.trim() !== '' && !declaredSlots.includes('')) {
        found.push('the page writes content into <' + tag + '> outside any named slot, and the element renders no default slot');
      }
    }
  }
  for (const control of document.querySelectorAll('button, textarea, input')) {
    found.push('the page writes a <' + control.localName + ' class="' + control.className + '">, and the controls draw theirs inside their own shadow roots');
  }
  for (const part of document.querySelectorAll('[data-part]')) {
    found.push('the page writes data-part="' + part.dataset.part + '" outside a shadow root');
  }
  return { found, tags: Array.from(new Set(primitives.map((held) => held.localName))) };
})()`

func aRichPage(t *testing.T) *reader {
	t.Helper()

	held := standing(t)
	aDrawing(t, held, "drawings/q3.html", "<!doctype html><p>a drawing</p>")

	first := aCacheDecision("single-choice")
	first["illustration"] = sketch{"path": "drawings/q3.html"}
	first["evidence"] = []sketch{
		{"language": "go", "code": "func Read() {}\n"},
		{"patch": patchedWrite},
	}
	second := aCacheDecision("multiple-choice")
	second["id"] = "c-2"
	second["options"].([]sketch)[0]["recommended"] = true
	document := aDocument("one thing to settle", []sketch{
		anEntry("c-1", "the read is unlocked", "error", "needs a person", "internal/cache/read.go:41"),
		anEntry("c-2", "the write counts twice", "warning", "can be patched", "internal/cache/write.go:12"),
	}, first, second)
	document["settled"] = whatWasSettled()
	return opening(t, held.serve(written(t, document)))
}

func TestARenderedPageWritesNoElementNoAttributeNoSlotAndNoPartThePrimitivesDoNotDeclare(t *testing.T) {
	t.Parallel()
	page := aRichPage(t)

	page.click(".send")
	page.waitFor(confirmationIsOpen)

	var drawn struct {
		Found []string `json:"found"`
		Tags  []string `json:"tags"`
	}
	page.evaluate(whatEveryPrimitiveOfThePageIsWrittenWith, &drawn, awaitingThePromise)

	for _, found := range drawn.Found {
		t.Error(found)
	}
	for _, reached := range []string{"ui-app-shell", "ui-badge", "ui-button", "ui-card", "ui-choice", "ui-choice-group", "ui-code-block", "ui-copy", "ui-dialog", "ui-diff-block", "ui-expand", "ui-field"} {
		if !slices.Contains(drawn.Tags, reached) {
			t.Errorf("the page drew no <%s>, so this proves nothing of it: %v", reached, drawn.Tags)
		}
	}
}

var piecesTheScriptFindsByTheirClass = []string{
	".surface", "form.exchange", "ui-app-shell", ".topbar", ".status", ".status-said", ".index", ".entries", ".entry", ".entry-link",
	".progress-bar", ".answered-count", ".main", ".main-inner", ".panel", ".summary", ".over-notice", ".over-notice-said", ".settled",
	".decision", ".decision .attribution .answered", ".decision .name", ".options-list", ".option", ".option-label", ".choice", ".decline",
	".words-field", ".tray", ".evidence", ".evidence.code", "ui-code-block", ".diff", "ui-diff-block", ".drawing-stage", "iframe.drawing",
	".drawing-close", "ui-expand", ".composer", ".message-field", ".send", ".finish", ".amend", ".receipt", ".receipt-decisions",
	".receipt-message", ".receipt-message-files", ".state", "ui-dialog.confirmation", ".confirmation-decisions", ".go-back",
	".send-anyway", "ui-tooltip", ".tooltip",
}

func TestEveryPieceTheScriptFindsItByIsDrawnUnderTheClassThatNamesIt(t *testing.T) {
	t.Parallel()
	page := aRichPage(t)

	var counts map[string]int
	page.evaluate(`Object.fromEntries(`+jsonOfList(t, piecesTheScriptFindsByTheirClass)+`.map((selector) => [selector, document.querySelectorAll(selector).length]))`, &counts)

	for _, selector := range piecesTheScriptFindsByTheirClass {
		if counts[selector] == 0 {
			t.Errorf("the page draws nothing the script can find by %q, and each piece keeps the one class that names it", selector)
		}
	}
}
