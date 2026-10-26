package e2e_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/agora/surface"
)

const wordsOfADecision = `.decision[data-decision="c-1"] .words-field`

type clipboardFile struct {
	name string
	kind string
}

func filesOf(files ...clipboardFile) string {
	var built []string
	for _, file := range files {
		built = append(built, fmt.Sprintf(`new File(['held'], %q, { type: %q })`, file.name, file.kind))
	}
	return `[` + strings.Join(built, ", ") + `]`
}

func (r *reader) dropOn(selector string, files ...clipboardFile) {
	r.t.Helper()

	r.evaluate(`(() => {
  const transfer = new DataTransfer();
  for (const file of `+filesOf(files...)+`) { transfer.items.add(file); }
  const target = document.querySelector(`+jsonOf(r.t, selector)+`);
  target.dispatchEvent(new DragEvent('dragover', { dataTransfer: transfer, bubbles: true, composed: true, cancelable: true }));
  target.dispatchEvent(new DragEvent('drop', { dataTransfer: transfer, bubbles: true, composed: true, cancelable: true }));
  return true;
})()`, nil)
}

func (r *reader) pasteInto(selector, text string, files ...clipboardFile) bool {
	r.t.Helper()

	var left bool
	r.evaluate(`(() => {
  const transfer = new DataTransfer();
  for (const file of `+filesOf(files...)+`) { transfer.items.add(file); }
  if (`+jsonOf(r.t, text)+` !== '') { transfer.setData('text/plain', `+jsonOf(r.t, text)+`); }
  const target = document.querySelector(`+jsonOf(r.t, selector)+`);
  const pasted = new ClipboardEvent('paste', { clipboardData: transfer, bubbles: true, composed: true, cancelable: true });
  target.dispatchEvent(pasted);
  return !pasted.defaultPrevented;
})()`, &left)
	return left
}

func (r *reader) heldBy(selector string) []string {
	r.t.Helper()

	var names []string
	r.evaluate(`document.querySelector(`+jsonOf(r.t, selector)+`).files.map((file) => file.name)`, &names)
	return names
}

func (r *reader) refusalOf(selector string) string {
	r.t.Helper()

	var said string
	r.evaluate(`document.querySelector(`+jsonOf(r.t, selector)+`).drawnRefusal`, &said)
	return said
}

func attachedNamesOf(answered surface.Answered, decision string) []string {
	var names []string
	for _, decided := range answered.Decisions {
		if decided.ID == decision {
			for _, attachment := range decided.Attachments {
				names = append(names, attachment.Name)
			}
		}
	}
	return names
}

func TestDroppingTwoFilesOnADecisionAttachesBothToThatDecisionsTrayAndToNoOther(t *testing.T) {
	// scenario: Dropping two files on a decision attaches both to that decision's tray and to no other
	t.Parallel()
	_, browser := read(t, twoDecisionDocument)

	browser.dropOn(`.decision[data-decision="c-1"] .words-field`, clipboardFile{"before.png", "image/png"}, clipboardFile{"trace.txt", "text/plain"})

	if diff := cmp.Diff([]string{"before.png", "trace.txt"}, browser.heldBy(trayOfDecision("c-1"))); diff != "" {
		t.Errorf("the tray of the decision dropped on (-want +got):\n%s", diff)
	}
	if held := browser.heldBy(trayOfDecision("c-2")); len(held) != 0 {
		t.Errorf("the other decision's tray holds %v", held)
	}
	if held := browser.heldBy(composerTray); len(held) != 0 {
		t.Errorf("the composer's tray holds %v", held)
	}
}

func TestDroppingFilesOnTheComposerAttachesThemToTheComposersTray(t *testing.T) {
	// scenario: Dropping files on the composer attaches them to the composer's tray
	t.Parallel()
	_, browser := read(t, twoDecisionDocument)

	browser.dropOn(`.composer .message-field`, clipboardFile{"trace.txt", "text/plain"})

	if diff := cmp.Diff([]string{"trace.txt"}, browser.heldBy(composerTray)); diff != "" {
		t.Errorf("the composer's tray (-want +got):\n%s", diff)
	}
	if held := browser.heldBy(trayOfDecision("c-1")); len(held) != 0 {
		t.Errorf("a decision's tray holds %v", held)
	}
}

func TestDroppingFilesOnTheTrayItselfAttachesThemOnce(t *testing.T) {
	// scenario: Dropping files on the tray itself attaches each once
	t.Parallel()
	_, browser := read(t, twoDecisionDocument)

	browser.dropOn(trayOfDecision("c-1"), clipboardFile{"trace.txt", "text/plain"})

	if diff := cmp.Diff([]string{"trace.txt"}, browser.heldBy(trayOfDecision("c-1"))); diff != "" {
		t.Errorf("the tray dropped on (-want +got):\n%s", diff)
	}
}

func TestPastingAnImageWithTheCursorOutsideTheWordsFieldOfACardAttachesItToThatCardsTray(t *testing.T) {
	// scenario: Pasting an image with the cursor outside the words field of a card attaches it to that card's tray
	t.Parallel()
	_, browser := read(t, twoDecisionDocument)

	browser.pasteInto(`.decision[data-decision="c-2"] .choice`, "", clipboardFile{"shot.png", "image/png"})

	if diff := cmp.Diff([]string{"shot.png"}, browser.heldBy(trayOfDecision("c-2"))); diff != "" {
		t.Errorf("the tray of the card pasted into (-want +got):\n%s", diff)
	}
	if held := browser.heldBy(trayOfDecision("c-1")); len(held) != 0 {
		t.Errorf("the other card's tray holds %v", held)
	}
}

func TestPastingAnImageIntoTheMessageAttachesItToTheComposersTray(t *testing.T) {
	// scenario: Pasting an image into the composer attaches it to the composer's tray
	t.Parallel()
	_, browser := read(t, twoDecisionDocument)

	browser.pasteInto(`.composer .message-field`, "", clipboardFile{"shot.png", "image/png"})

	if diff := cmp.Diff([]string{"shot.png"}, browser.heldBy(composerTray)); diff != "" {
		t.Errorf("the composer's tray (-want +got):\n%s", diff)
	}
}

func TestTwoPastedUnnamedImagesComeBackFromPollAsPasted1AndPasted2(t *testing.T) {
	// scenario: Two pasted unnamed images arrive in the poll as pasted-1.png and pasted-2.png
	t.Parallel()
	held, browser := read(t, twoDecisionDocument)

	browser.evaluate(chose("c-1", 0), nil)
	browser.pasteInto(wordsOfADecision, "", clipboardFile{"", "image/png"})
	browser.pasteInto(wordsOfADecision, "", clipboardFile{"", "image/png"})
	browser.click(".send")
	browser.waitFor(`!document.querySelector('.receipt').hidden`)

	if diff := cmp.Diff([]string{"pasted-1.png", "pasted-2.png"}, attachedNamesOf(held.polledAnswer(), "c-1")); diff != "" {
		t.Errorf("the names poll printed (-want +got):\n%s", diff)
	}
}

func TestAPastedImageTheBrowserCallsImageIsNamedLikeAnUnnamedOne(t *testing.T) {
	// scenario: A pasted screenshot the browser names image.png is named pasted-1.png
	t.Parallel()
	_, browser := read(t, twoDecisionDocument)

	browser.pasteInto(wordsOfADecision, "", clipboardFile{"image.png", "image/png"}, clipboardFile{"trace.txt", "text/plain"})

	if diff := cmp.Diff([]string{"pasted-1.png", "trace.txt"}, browser.heldBy(trayOfDecision("c-1"))); diff != "" {
		t.Errorf("the tray pasted into (-want +got):\n%s", diff)
	}
}

func TestPastingTextIntoTheWordsLeavesTheTextToPasteAndAttachesNothing(t *testing.T) {
	// scenario: Pasting text into the words leaves the text in the words and attaches nothing
	t.Parallel()
	_, browser := read(t, twoDecisionDocument)

	leftToPaste := browser.pasteInto(wordsOfADecision, "only words")

	if !leftToPaste {
		t.Error("the page took a paste of text away from the words")
	}
	for _, tray := range []string{trayOfDecision("c-1"), trayOfDecision("c-2"), composerTray} {
		if held := browser.heldBy(tray); len(held) != 0 {
			t.Errorf("%s holds %v after a paste of text", tray, held)
		}
	}
}

func TestPastingAClipboardOfFilesAndTextAttachesTheFilesAndLeavesTheTextToPaste(t *testing.T) {
	// scenario: Pasting a clipboard of files and text attaches the files and pastes the text
	t.Parallel()
	_, browser := read(t, twoDecisionDocument)

	leftToPaste := browser.pasteInto(wordsOfADecision, "some words", clipboardFile{"trace.txt", "text/plain"})

	if !leftToPaste {
		t.Error("the page took the text of a mixed paste away from the words")
	}
	if diff := cmp.Diff([]string{"trace.txt"}, browser.heldBy(trayOfDecision("c-1"))); diff != "" {
		t.Errorf("the tray pasted into (-want +got):\n%s", diff)
	}
}

func TestPastingIntoAFullTrayRefusesInPlaceNamingTheCeiling(t *testing.T) {
	// scenario: Pasting into a full tray refuses in place, naming the ceiling
	t.Parallel()
	_, browser := read(t, twoDecisionDocument)
	var ceiling int
	browser.evaluate(`document.querySelector(`+jsonOf(t, trayOfDecision("c-1"))+`).max`, &ceiling)
	var names []string
	for at := range ceiling {
		names = append(names, fmt.Sprintf("held-%d.txt", at))
	}
	browser.addToTray(trayOfDecision("c-1"), names...)

	browser.pasteInto(wordsOfADecision, "", clipboardFile{"shot.png", "image/png"})

	if held := browser.heldBy(trayOfDecision("c-1")); len(held) != ceiling {
		t.Errorf("the full tray holds %d files, want %d", len(held), ceiling)
	}
	if said := browser.refusalOf(trayOfDecision("c-1")); !strings.Contains(said, fmt.Sprint(ceiling)) {
		t.Errorf("the refusal says %q, want it to name the ceiling %d", said, ceiling)
	}
}

func TestAPasteReadsOnlyTheFilesThePasteEventCarriesAndNeverReadsTheClipboard(t *testing.T) {
	// scenario: A paste reads only the files the paste event carries and never reads the clipboard itself
	t.Parallel()
	_, browser := read(t, twoDecisionDocument)
	browser.evaluate(`(() => {
  window.clipboardReads = 0;
  const counting = () => { window.clipboardReads += 1; return Promise.resolve([]); };
  Object.defineProperty(navigator, 'clipboard', { configurable: true, value: { read: counting, readText: counting } });
  return true;
})()`, nil)

	browser.pasteInto(wordsOfADecision, "some words", clipboardFile{"shot.png", "image/png"})

	var reads int
	browser.evaluate(`window.clipboardReads`, &reads)
	if reads != 0 {
		t.Errorf("the page read the clipboard %d times during a paste", reads)
	}
	if diff := cmp.Diff([]string{"shot.png"}, browser.heldBy(trayOfDecision("c-1"))); diff != "" {
		t.Errorf("the tray pasted into (-want +got):\n%s", diff)
	}
}
