package e2e_test

import (
	"bytes"
	"image/color"
	"image/png"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/agora/surface"
)

const somethingLiveOnThePage = `Array.from(document.querySelectorAll('button, .choice, textarea')).some((held) => !held.disabled)`

const everythingOnThePageIsLive = `Array.from(document.querySelectorAll('button, .choice, textarea')).every((held) => !held.disabled)`

const confirmationIsOpen = `document.querySelector('dialog.confirmation').open === true`

const confirmationIsClosed = `document.querySelector('dialog.confirmation').open === false`

const confirmationIsModal = `document.querySelector('dialog.confirmation').matches(':modal')`

const whatTheConfirmationSays = `document.querySelector('dialog.confirmation').textContent`

const whatTheConfirmationIsHeaded = `document.querySelector('.confirmation-heading').textContent`

const whatTheConfirmationNames = `Array.from(document.querySelectorAll('.confirmation-decision')).map((held) => held.dataset.decision)`

const whereTheCursorRests = `document.activeElement.className`

func pressPastTheConfirmation(page *reader, press string) {
	page.t.Helper()
	page.click(press)
	page.waitFor(confirmationIsOpen)
	page.click(".send-anyway")
}

func (s *surfaced) ended(t *testing.T, page *reader) {
	t.Helper()

	page.evaluate(chose("c-1", 0), nil)
	page.click(".finish")
	page.waitFor(`document.querySelector('.send').disabled === true`)

	if _, arrived := s.awaiting(5 * time.Second); !arrived {
		t.Fatal("finishing delivered nothing")
	}
}

func rendersClosed(t *testing.T, page *reader, how string) {
	t.Helper()

	var said string
	page.evaluate(`document.querySelector('.state').textContent`, &said)
	if !strings.Contains(strings.ToLower(said), "over") {
		t.Errorf("%s says %q, and the exchange it invites an answer to is over", how, said)
	}

	var accepting bool
	page.evaluate(somethingLiveOnThePage, &accepting)
	if accepting {
		t.Errorf("%s still offers a control that would send an answer", how)
	}

	var showing closedShowing
	page.evaluate(whatAClosedPageShows, &showing)
	if !showing.Notice || showing.Words || showing.Message {
		t.Errorf("%s shows %+v, and it says under the title that the exchange is over and draws no field nothing will read", how, showing)
	}
}

func TestAReloadedPageRendersClosed(t *testing.T) {
	held, page := read(t, twoDecisionDocument)
	held.ended(t, page)

	page.reload()

	rendersClosed(t, page, "a page reloaded onto an ended exchange")
}

func answersWhatTheCallerAskedAgain(t *testing.T, held *surfaced, page *reader, how string) {
	t.Helper()

	page.evaluate(chose("c-3", 0), nil)
	page.click(".send")

	delivered, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		var said string
		page.evaluate(`document.querySelector('.state').textContent`, &said)
		t.Fatalf("an answer sent from %s was not delivered, and the page says %q", how, said)
	}
	if len(delivered.Decisions) != 1 || delivered.Decisions[0].ID != "c-3" {
		t.Errorf("what %s delivered is %+v", how, delivered.Decisions)
	}

	var said string
	page.evaluate(`document.querySelector('.state').textContent`, &said)
	if strings.Contains(said, "not accepted") {
		t.Errorf("an answer sent from %s was refused: %q", how, said)
	}
}

func TestANewDocumentReopensAnEndedPage(t *testing.T) {
	held, page := read(t, twoDecisionDocument)
	held.ended(t, page)

	held.serveAgain(oneDecisionReply)
	page.waitFor(reachedThePage("c-3"))

	answersWhatTheCallerAskedAgain(t, held, page, "the page the caller asked again on, without it being reloaded")
}

func TestAReloadedPageStillReopens(t *testing.T) {
	held, page := read(t, twoDecisionDocument)
	held.ended(t, page)

	held.serveAgain(oneDecisionReply)
	page.reload()

	answersWhatTheCallerAskedAgain(t, held, page, "the page the caller asked again on, reloaded")
}

func TestServingAgainOntoAnEndedExchangeYieldsAnAnswerablePage(t *testing.T) {
	held := standing(t)
	held.serve(oneDecisionDocument)

	held.answer(surface.Sent{Ends: true})
	if _, arrived := held.awaiting(5 * time.Second); !arrived {
		t.Fatal("ending without answering delivered nothing")
	}

	address := held.serveAgain(oneDecisionReply)
	if err := surface.CheckLoopback(originOf(t, address)[len("http://"):]); err != nil {
		t.Fatalf("serving a document again printed %q: %v", address, err)
	}

	held.answer(surface.Sent{Decisions: []surface.Choice{{ID: "c-3", Chosen: []int{1}}}})

	delivered, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("an answer sent to the address the caller was handed back was not delivered")
	}
	if diff := cmp.Diff([]string{"patch c-3"}, delivered.Decisions[0].Commands); diff != "" {
		t.Errorf("what the caller was handed back differs (-want +got):\n%s", diff)
	}
}

func TestServingAgainWhileTheExchangeIsOpenIsUnchanged(t *testing.T) {
	held := standing(t)
	first := held.serve(oneDecisionDocument)

	again := held.serveAgain(oneDecisionReply)
	if again != first {
		t.Errorf("serving a new document to an open surface printed %q, and the surface it was written to is at %q", again, first)
	}

	held.answer(surface.Sent{Decisions: []surface.Choice{{ID: "c-3", Chosen: []int{0}}}})

	delivered, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("an answer sent after a new document reached an open surface was not delivered")
	}
	if delivered.Ended {
		t.Error("serving a new document to an open surface reports the exchange as over")
	}
	if diff := cmp.Diff([]string{"agree c-3"}, delivered.Decisions[0].Commands); diff != "" {
		t.Errorf("what came back was not resolved against the document the surface took up (-want +got):\n%s", diff)
	}
}

func TestAPollOverAnEndedExchangeLeavesItEnded(t *testing.T) {
	held, page := read(t, twoDecisionDocument)
	held.ended(t, page)

	if delivered, arrived := held.awaiting(shortWindow); arrived {
		t.Errorf("a poll over an ended exchange delivered something further: %+v", delivered)
	}

	rendersClosed(t, page, "a page after a poll over an ended exchange")
}

func TestACallerWhoAsksAgainIsHandedAnAddressThatTakesTheAnswer(t *testing.T) {
	// scenario: Serving again onto an ended exchange yields an answerable page
	held, page := read(t, oneDecisionDocument)

	pressPastTheConfirmation(page, ".finish")
	page.waitFor(`document.querySelector('.send').disabled === true`)

	finished, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("finishing without answering delivered nothing")
	}
	if len(finished.Decisions) != 0 || !finished.Ended {
		t.Fatalf("finishing without answering delivered %+v", finished)
	}

	address := held.serveAgain(oneDecisionDocument)
	asked := opening(t, address)

	asked.evaluate(chose("c-1", 0), nil)
	asked.click(".finish")
	asked.waitFor(`document.querySelector('.send').disabled === true`)

	answered, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("the person answered at the address the caller was handed back, and nothing reached the caller")
	}
	if diff := cmp.Diff([]string{"patch c-1"}, answered.Decisions[0].Commands); diff != "" {
		t.Errorf("what reached the caller the second time differs (-want +got):\n%s", diff)
	}
}

func TestAPageOpenedForTheFirstTimeAfterTheEndRendersClosed(t *testing.T) {
	held := standing(t)
	address := held.serve(twoDecisionDocument)

	held.end(held.opened)

	page := opening(t, address)
	rendersClosed(t, page, "a page opened for the first time onto an ended exchange")

	var showing closedShowing
	page.evaluate(whatAClosedPageShows, &showing)
	if showing != (closedShowing{Notice: true}) {
		t.Errorf("a page opened for the first time onto an ended exchange shows %+v, and with nothing sent from it there is no composer to keep", showing)
	}
}

func TestAPageServedClosedFollowsWhatComesNext(t *testing.T) {
	// scenario: A new document reopens an ended page
	held := standing(t)
	address := held.serve(twoDecisionDocument)

	held.end(held.opened)

	page := opening(t, address)
	held.serveAgain(oneDecisionReply)
	page.waitFor(reachedThePage("c-3"))

	answersWhatTheCallerAskedAgain(t, held, page, "a page opened onto an ended exchange the caller then asked again on")

	streams := []string{}
	for _, reaching := range page.requests.all() {
		if strings.HasSuffix(reaching, "/live") {
			streams = append(streams, reaching)
		}
	}
	if len(streams) != 1 {
		t.Errorf("the page opened %d streams, and one held open is what reaches it without it reopening one every few seconds: %v", len(streams), streams)
	}
}

func TestThePageStopsSayingTheExchangeIsOver(t *testing.T) {
	held, page := read(t, twoDecisionDocument)
	held.ended(t, page)
	rendersClosed(t, page, "a page whose exchange the person ended")

	held.serveAgain(oneDecisionReply)
	page.waitFor(reachedThePage("c-3"))

	var said string
	page.evaluate(`document.querySelector('.state').textContent`, &said)
	if strings.Contains(strings.ToLower(said), "over") {
		t.Errorf("a page the caller asked again on still says %q", said)
	}

	var operable bool
	page.evaluate(everythingOnThePageIsLive, &operable)
	if !operable {
		t.Error("a page the caller asked again on carries a control that would not take an answer")
	}
}

func TestAnAnswerChosenBeforeTheEndIsNotCarriedIntoWhatComesNext(t *testing.T) {
	// scenario: A new document reopens an ended page
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.click(".finish")
	page.waitFor(`document.querySelector('.send').disabled === true`)
	if _, arrived := held.awaiting(5 * time.Second); !arrived {
		t.Fatal("finishing delivered nothing")
	}

	held.serveAgain(oneDecisionReply)
	page.waitFor(reachedThePage("c-3"))

	pressPastTheConfirmation(page, ".send")

	delivered, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("sending from the page the caller asked again on delivered nothing")
	}
	for _, decided := range delivered.Decisions {
		if decided.ID == "c-1" {
			t.Errorf("an answer chosen before the exchange ended was carried into content that does not hold it: %+v", decided)
		}
	}
}

func TestEitherPressOverABlankDecisionRaisesTheConfirmationAndSendsNothing(t *testing.T) {
	for _, press := range []string{".send", ".finish"} {
		t.Run(press, func(t *testing.T) {
			held, page := read(t, twoDecisionDocument)

			page.evaluate(chose("c-1", 0), nil)
			page.evaluate(chose("c-2", 2), nil)
			page.click(press)
			page.waitFor(confirmationIsOpen)

			if delivered, arrived := held.awaiting(shortWindow); arrived {
				t.Fatalf("pressing %s over a blank decision delivered an answer: %+v", press, delivered)
			}

			var modal bool
			page.evaluate(confirmationIsModal, &modal)
			if !modal {
				t.Errorf("the confirmation pressing %s raised is not modal, and the page behind it is still answerable", press)
			}

			var said string
			page.evaluate(whatTheConfirmationSays, &said)
			if !strings.Contains(said, "c-2") {
				t.Errorf("the confirmation pressing %s raised reads %q, and it does not name the decision left blank", press, said)
			}

			var cursor string
			page.evaluate(whereTheCursorRests, &cursor)
			if cursor != "send-anyway" {
				t.Errorf("the cursor rests on %q when the confirmation pressing %s raised opens, and it opens on Send anyway", cursor, press)
			}
		})
	}
}

func TestTheConfirmationNamesEveryBlankDecisionAndOnlyThose(t *testing.T) {
	_, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-2", 2), nil)
	page.click(".send")
	page.waitFor(confirmationIsOpen)

	var named []string
	page.evaluate(whatTheConfirmationNames, &named)
	if diff := cmp.Diff([]string{"c-1", "c-2"}, named); diff != "" {
		t.Errorf("the decisions the confirmation names over a page holding two blanks differ (-want +got):\n%s", diff)
	}

	var heading string
	page.evaluate(whatTheConfirmationIsHeaded, &heading)
	if !strings.Contains(heading, "2 of 2") {
		t.Errorf("the confirmation over two blank decisions is headed %q, and it says how many of how many are blank", heading)
	}

	page.click(".go-back")
	page.waitFor(confirmationIsClosed)

	page.evaluate(chose("c-1", 0), nil)
	page.click(".send")
	page.waitFor(confirmationIsOpen)

	page.evaluate(whatTheConfirmationNames, &named)
	if diff := cmp.Diff([]string{"c-2"}, named); diff != "" {
		t.Errorf("the decisions the confirmation names after one of the two was answered differ (-want +got):\n%s", diff)
	}

	page.evaluate(whatTheConfirmationIsHeaded, &heading)
	if !strings.Contains(heading, "1 of 2") {
		t.Errorf("the confirmation after one of the two was answered is headed %q, and it says how many of how many are blank", heading)
	}
}

func TestSendAnywayDeliversWhatTheControlsHold(t *testing.T) {
	for _, raising := range []struct {
		press string
		ends  bool
	}{{".send", false}, {".finish", true}} {
		t.Run(raising.press, func(t *testing.T) {
			held, page := read(t, twoDecisionDocument)

			page.evaluate(chose("c-1", 0), nil)
			page.evaluate(chose("c-2", 2), nil)
			page.click(raising.press)
			page.waitFor(confirmationIsOpen)

			page.click(".send-anyway")

			delivered, arrived := held.awaiting(5 * time.Second)
			if !arrived {
				t.Fatalf("sending anyway after pressing %s delivered nothing", raising.press)
			}
			if len(delivered.Decisions) != 1 || delivered.Decisions[0].ID != "c-1" {
				t.Errorf("sending anyway after pressing %s carried %+v, want only the decision that was not blank", raising.press, delivered.Decisions)
			}
			if delivered.Ended != raising.ends {
				t.Errorf("sending anyway after pressing %s ended the exchange: %v, want %v", raising.press, delivered.Ended, raising.ends)
			}

			page.waitFor(confirmationIsClosed)
		})
	}
}

const whatTheControlsHold = `(() => ({
  chosen: Array.from(document.querySelectorAll('.choice:checked')).map((held) => held.dataset.decision + ':' + held.value),
  message: document.querySelector('.message-field').value,
}))()`

type holding struct {
	Chosen  []string `json:"chosen"`
	Message string   `json:"message"`
}

func TestEveryWayOutOfTheConfirmationSendsNothingAndLeavesThePageAsItWas(t *testing.T) {
	const wrote = "the second one is blank on purpose"

	for _, door := range []struct {
		name  string
		taken func(page *reader)
	}{
		{"go back and answer", func(page *reader) { page.click(".go-back") }},
		{"escape", func(page *reader) { page.press(kb.Escape) }},
		{"a press on the backdrop", func(page *reader) { page.clickAt(4, 4) }},
	} {
		t.Run(door.name, func(t *testing.T) {
			held, page := read(t, twoDecisionDocument)

			page.evaluate(chose("c-1", 0), nil)
			page.evaluate(chose("c-2", 2), nil)
			page.write(".message-field", wrote)

			page.click(".finish")
			page.waitFor(confirmationIsOpen)

			door.taken(page)
			page.waitFor(confirmationIsClosed)

			if delivered, arrived := held.awaiting(shortWindow); arrived {
				t.Fatalf("leaving the confirmation by %s delivered an answer anyway: %+v", door.name, delivered)
			}

			var read holding
			page.evaluate(whatTheControlsHold, &read)
			if diff := cmp.Diff([]string{"c-1:0"}, read.Chosen); diff != "" {
				t.Errorf("leaving the confirmation by %s moved what the controls hold (-want +got):\n%s", door.name, diff)
			}
			if read.Message != wrote {
				t.Errorf("leaving the confirmation by %s left the message as %q, want %q", door.name, read.Message, wrote)
			}

			var cursor string
			page.evaluate(whereTheCursorRests, &cursor)
			if cursor != "finish" {
				t.Errorf("leaving the confirmation by %s rested the cursor on %q, and it returns to the button that raised it", door.name, cursor)
			}

			var accepting bool
			page.evaluate(everythingOnThePageIsLive, &accepting)
			if !accepting {
				t.Errorf("leaving the confirmation by %s left the page unable to take a further answer", door.name)
			}
		})
	}
}

func TestPressingTheSameButtonAgainAfterGoingBackAsksAgainRatherThanSending(t *testing.T) {
	// scenario: no armed state survives the confirmation
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.evaluate(chose("c-2", 2), nil)

	page.click(".send")
	page.waitFor(confirmationIsOpen)
	page.click(".go-back")
	page.waitFor(confirmationIsClosed)

	page.click(".send")
	page.waitFor(confirmationIsOpen)

	if delivered, arrived := held.awaiting(shortWindow); arrived {
		t.Fatalf("pressing send a second time after going back sent rather than asking again: %+v", delivered)
	}
}

func TestAPressWithNothingBlankIsNotAskedAbout(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.evaluate(chose("c-2", 0), nil)
	page.click(".send")

	delivered, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("sending with nothing blank delivered nothing")
	}
	if len(delivered.Decisions) != 2 {
		t.Errorf("sending with nothing blank carried %+v, want both decisions", delivered.Decisions)
	}

	var closed bool
	page.evaluate(confirmationIsClosed, &closed)
	if !closed {
		t.Error("a press with nothing blank raised the confirmation anyway")
	}
}

func TestARevisionClosesTheConfirmation(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.evaluate(chose("c-2", 2), nil)
	page.click(".send")
	page.waitFor(confirmationIsOpen)

	held.serveAgain(oneDecisionReply)
	page.waitFor(reachedThePage("c-3"))
	page.waitFor(confirmationIsClosed)

	var open int
	page.evaluate(`document.querySelectorAll('dialog[open]').length`, &open)
	if open != 0 {
		t.Errorf("a revision left %d dialogs open over cards that are no longer on the page", open)
	}
}

const whereTheCursorLanded = `(() => {
  const held = document.activeElement;
  const card = held.closest('.decision');
  const box = card ? card.getBoundingClientRect() : null;
  return {
    decision: card ? card.dataset.decision : '',
    control: held.className,
    inView: !!box && box.top < window.innerHeight && box.bottom > 0,
    hash: location.hash,
  };
})()`

type landed struct {
	Decision string `json:"decision"`
	Control  string `json:"control"`
	InView   bool   `json:"inView"`
	Hash     string `json:"hash"`
}

const offeringNoneBesideOfferingSome = `{
  "title": "one question left before the spec is written",
  "decisions": [
    {"id": "q1", "question": "what should the loopback server be called in the glossary", "meaning": "the surface, matching what agora's own source already calls it"},
    {"id": "q2", "question": "should a design round carry a cap the way a review round does", "meaning": "no cap, because nothing here spawns a process for a cap to protect"},
    {
      "id": "q3",
      "question": "should CONTEXT.md gain an entry for the surface today or wait until the spec is written",
      "meaning": "the glossary is maintained as terms settle",
      "mode": "single-choice",
      "options": [
        {"label": "write it now", "meaning": "the term is already settled", "recommended": true},
        {"label": "wait for the spec", "meaning": "batch every glossary change with the write"}
      ]
    }
  ]
}`

func TestABlankDecisionsNameTakesThePersonToItsCard(t *testing.T) {
	for _, jump := range []struct {
		named   string
		to      string
		control string
	}{
		{"a decision offering options", "q3", "choice"},
		{"a decision offering none", "q1", "words-field"},
	} {
		t.Run(jump.named, func(t *testing.T) {
			held := standing(t)
			page := opening(t, held.serve(offeringNoneBesideOfferingSome))

			page.evaluate(declined("q3"), nil)
			page.write(within("q2", ".words-field"), inTheirOwnWordsAboutTheCap)

			page.click(".send")
			page.waitFor(confirmationIsOpen)

			var named []string
			page.evaluate(whatTheConfirmationNames, &named)
			if diff := cmp.Diff([]string{"q1", "q3"}, named); diff != "" {
				t.Fatalf("the decisions the confirmation names differ (-want +got):\n%s", diff)
			}

			page.click(`.confirmation-decision[data-decision="` + jump.to + `"]`)
			page.waitFor(confirmationIsClosed)

			var read landed
			page.evaluate(whereTheCursorLanded, &read)
			if read.Decision != jump.to {
				t.Errorf("pressing the name of %s left the cursor on the card %q, want %q", jump.named, read.Decision, jump.to)
			}
			if read.Control != jump.control {
				t.Errorf("pressing the name of %s left the cursor on a %q, want the card's %q", jump.named, read.Control, jump.control)
			}
			if !read.InView {
				t.Errorf("pressing the name of %s left its card out of view", jump.named)
			}
			if read.Hash != "" {
				t.Errorf("pressing the name of %s left %q in the address bar", jump.named, read.Hash)
			}

			if delivered, arrived := held.awaiting(shortWindow); arrived {
				t.Fatalf("pressing the name of %s delivered an answer: %+v", jump.named, delivered)
			}

			var still holding
			page.evaluate(whatTheControlsHold, &still)
			if len(still.Chosen) != 0 {
				t.Errorf("pressing the name of %s moved what the options hold: %v", jump.named, still.Chosen)
			}

			var words string
			page.evaluate(`document.querySelector('`+within("q2", ".words-field")+`').value`, &words)
			if words != inTheirOwnWordsAboutTheCap {
				t.Errorf("pressing the name of %s left the words typed on another card as %q", jump.named, words)
			}
		})
	}
}

func cornerOfTheViewport(t *testing.T, page *reader) color.Color {
	t.Helper()

	var shot []byte
	page.do(chromedp.CaptureScreenshot(&shot))

	painted, err := png.Decode(bytes.NewReader(shot))
	if err != nil {
		t.Fatalf("reading what the page painted: %v", err)
	}
	return painted.At(4, 4)
}

func brightnessOf(painted color.Color) float64 {
	red, green, blue, _ := painted.RGBA()
	return float64(red+green+blue) / 3
}

func TestTheConfirmationDimsThePageBehindIt(t *testing.T) {
	for _, preference := range []string{"light", "dark"} {
		t.Run(preference, func(t *testing.T) {
			_, page := read(t, twoDecisionDocument)
			page.emulate(preference)
			page.evaluate(chose("c-2", 2), nil)

			before := cornerOfTheViewport(t, page)

			page.click(".send")
			page.waitFor(confirmationIsOpen)

			behind := cornerOfTheViewport(t, page)

			if brightnessOf(behind) >= brightnessOf(before) {
				t.Errorf("under a %s preference the page behind the confirmation is painted %v against %v before it opened, and what is behind a modal is dimmed", preference, behind, before)
			}
		})
	}
}

func TestAPressInsideTheConfirmationIsNotAPressOnTheBackdrop(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.evaluate(chose("c-2", 2), nil)
	page.click(".send")
	page.waitFor(confirmationIsOpen)

	var padding struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}
	page.evaluate(`(() => {
  const box = document.querySelector('dialog.confirmation').getBoundingClientRect();
  return { x: box.left + 4, y: box.top + 4 };
})()`, &padding)

	page.clickAt(padding.X, padding.Y)

	var open bool
	page.evaluate(confirmationIsOpen, &open)
	if !open {
		t.Error("a press on the confirmation's own padding closed it, and only a press outside it is a press on the backdrop")
	}

	if delivered, arrived := held.awaiting(shortWindow); arrived {
		t.Fatalf("a press on the confirmation's own padding delivered an answer: %+v", delivered)
	}
}

const whatAClosedPageShows = `(() => ({
  notice: document.querySelector('.over-notice').checkVisibility(),
  words: Array.from(document.querySelectorAll('.words-field')).some((held) => held.checkVisibility()),
  message: document.querySelector('.message-field').checkVisibility(),
  composer: document.querySelector('.composer').checkVisibility(),
  receipt: document.querySelector('.receipt').checkVisibility(),
}))()`

type closedShowing struct {
	Notice   bool `json:"notice"`
	Words    bool `json:"words"`
	Message  bool `json:"message"`
	Composer bool `json:"composer"`
	Receipt  bool `json:"receipt"`
}

func TestAClosedExchangeHidesWhatCanNoLongerBeAnsweredAndSaysSoUnderTheTitle(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	var open closedShowing
	page.evaluate(whatAClosedPageShows, &open)
	if open != (closedShowing{Words: true, Message: true, Composer: true}) {
		t.Errorf("an open page shows %+v, and it offers every field and says nothing of the exchange being over", open)
	}

	held.ended(t, page)

	var finished closedShowing
	page.evaluate(whatAClosedPageShows, &finished)
	if finished != (closedShowing{Notice: true, Composer: true, Receipt: true}) {
		t.Errorf("a page that finished shows %+v, and it keeps what went and drops every field", finished)
	}
}
