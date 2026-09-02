package e2e_test

import (
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/ktei/socrates/internal/surface"
)

const somethingLiveOnThePage = `Array.from(document.querySelectorAll('button, .choice, textarea')).some((held) => !held.disabled)`

const everythingOnThePageIsLive = `Array.from(document.querySelectorAll('button, .choice, textarea')).every((held) => !held.disabled)`

func (s *surfaced) ended(t *testing.T, page *reader) {
	t.Helper()

	page.evaluate(chose("c-1", 0), nil)
	page.click(within("c-1", ".commit"))
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
}

func TestAReloadedPageRendersClosed(t *testing.T) {
	held, page := read(t, twoDecisionDocument)
	held.ended(t, page)

	page.reload()

	rendersClosed(t, page, "a page reloaded onto an ended exchange")
}

func TestFinishingOverAnUncommittedAnswerAsksFirst(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.click(".finish")
	page.waitFor(`!document.querySelector('.discarding').hidden`)

	var said string
	page.evaluate(`document.querySelector('.discarding-said').textContent`, &said)
	if !strings.Contains(said, "c-1") {
		t.Errorf("what the page says before it discards an uncommitted answer does not name the decision holding it: %q", said)
	}
	if !strings.Contains(said, "will not reach") {
		t.Errorf("what the page says before it discards an uncommitted answer does not say it will not reach whoever asked: %q", said)
	}

	var live bool
	page.evaluate(`!document.querySelector('.send').disabled`, &live)
	if !live {
		t.Error("finishing over an uncommitted answer closed the page before the person confirmed")
	}
	if delivered, arrived := held.awaiting(shortWindow); arrived {
		t.Errorf("finishing over an uncommitted answer sent before the person confirmed: %+v", delivered)
	}
}

func TestConfirmingCarriesNothingThatWasNotCommitted(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.evaluate(chose("c-2", 0), nil)
	page.click(within("c-2", ".commit"))
	page.click(".finish")
	page.waitFor(`!document.querySelector('.discarding').hidden`)
	page.click(".discard")
	page.waitFor(`document.querySelector('.send').disabled === true`)

	delivered, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("confirming the finish delivered nothing")
	}
	if !delivered.Ended {
		t.Error("confirming the finish does not report the exchange as over")
	}
	for _, decided := range delivered.Decisions {
		if decided.ID == "c-1" {
			t.Errorf("confirming the finish carried an answer that was never committed: %+v", decided)
		}
	}
	if len(delivered.Decisions) != 1 || delivered.Decisions[0].ID != "c-2" {
		t.Errorf("confirming the finish carried %+v, want the one answer that was committed", delivered.Decisions)
	}
}

func TestDecliningLeavesTheExchangeOpen(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.click(".finish")
	page.waitFor(`!document.querySelector('.discarding').hidden`)
	page.click(".keep")
	page.waitFor(`document.querySelector('.discarding').hidden === true`)

	var live bool
	page.evaluate(`!document.querySelector('.send').disabled`, &live)
	if !live {
		t.Error("declining the finish closed the page")
	}
	if delivered, arrived := held.awaiting(shortWindow); arrived {
		t.Errorf("declining the finish delivered something: %+v", delivered)
	}

	var ticked bool
	page.evaluate(`document.querySelector('.decision[data-decision="c-1"] .choice[value="0"]').checked`, &ticked)
	if !ticked {
		t.Fatal("declining the finish took away the selection the decision held")
	}

	page.click(within("c-1", ".commit"))
	page.click(".send")
	page.waitFor(`document.querySelector('.state').textContent.length > 0`)

	delivered, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("what was still there to commit was committed and sent, and nothing was delivered")
	}
	if len(delivered.Decisions) != 1 || delivered.Decisions[0].ID != "c-1" {
		t.Errorf("what was committed after declining came back as %+v", delivered.Decisions)
	}
}

func TestAFinishWithNothingUncommittedIsNotAskedAbout(t *testing.T) {
	for _, held := range []struct {
		how    string
		answer func(page *reader)
	}{
		{how: "every decision committed", answer: func(page *reader) {
			page.evaluate(chose("c-1", 0), nil)
			page.click(within("c-1", ".commit"))
			page.evaluate(chose("c-2", 1), nil)
			page.click(within("c-2", ".commit"))
		}},
		{how: "no decision touched at all", answer: func(page *reader) {}},
	} {
		t.Run(held.how, func(t *testing.T) {
			surfaced, page := read(t, twoDecisionDocument)
			held.answer(page)

			page.click(".finish")
			page.waitFor(`document.querySelector('.send').disabled === true`)

			var asking bool
			page.evaluate(`!document.querySelector('.discarding').hidden`, &asking)
			if asking {
				t.Errorf("finishing with %s asked before it sent", held.how)
			}

			delivered, arrived := surfaced.awaiting(5 * time.Second)
			if !arrived {
				t.Fatalf("finishing with %s delivered nothing", held.how)
			}
			if !delivered.Ended {
				t.Errorf("finishing with %s does not report the exchange as over", held.how)
			}
		})
	}
}

func answersWhatTheCallerAskedAgain(t *testing.T, held *surfaced, page *reader, how string) {
	t.Helper()

	page.evaluate(chose("c-3", 0), nil)
	page.click(within("c-3", ".commit"))
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

func TestAPersonWhoFinishesEarlyIsToldBeforeTheirAnswerGoes(t *testing.T) {
	// scenario: Finishing over an uncommitted answer asks first
	held, page := read(t, oneDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.click(".finish")
	page.waitFor(`!document.querySelector('.discarding').hidden`)

	var said string
	page.evaluate(`document.querySelector('.discarding-said').textContent`, &said)
	if !strings.Contains(said, "c-1") || !strings.Contains(said, "will not reach") {
		t.Errorf("the page did not say what finishing was about to discard: %q", said)
	}

	page.click(".keep")
	page.waitFor(`document.querySelector('.discarding').hidden === true`)

	page.click(within("c-1", ".commit"))
	page.write(".message-field", "this is the one")
	page.click(".finish")
	page.waitFor(`document.querySelector('.send').disabled === true`)

	delivered, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("the answer the person went back and committed was never delivered")
	}
	if !delivered.Ended {
		t.Error("finishing after the person confirmed does not report the exchange as over")
	}
	if delivered.Message != "this is the one" {
		t.Errorf("the message the person wrote came back as %q", delivered.Message)
	}
	if diff := cmp.Diff([]string{"patch c-1"}, delivered.Decisions[0].Commands); diff != "" {
		t.Errorf("what reached the caller differs (-want +got):\n%s", diff)
	}
}

func TestACallerWhoAsksAgainIsHandedAnAddressThatTakesTheAnswer(t *testing.T) {
	// scenario: Serving again onto an ended exchange yields an answerable page
	held, page := read(t, oneDecisionDocument)

	page.click(".finish")
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
	asked.click(within("c-1", ".commit"))
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

	rendersClosed(t, opening(t, address), "a page opened for the first time onto an ended exchange")
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

func TestAnAnswerCommittedBeforeTheEndIsNotCarriedIntoWhatComesNext(t *testing.T) {
	// scenario: A new document reopens an ended page
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.click(within("c-1", ".commit"))
	page.click(".finish")
	page.waitFor(`document.querySelector('.send').disabled === true`)
	if _, arrived := held.awaiting(5 * time.Second); !arrived {
		t.Fatal("finishing delivered nothing")
	}

	held.serveAgain(oneDecisionReply)
	page.waitFor(reachedThePage("c-3"))

	page.click(".send")

	delivered, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("sending from the page the caller asked again on delivered nothing")
	}
	for _, decided := range delivered.Decisions {
		if decided.ID == "c-1" {
			t.Errorf("an answer committed before the exchange ended was carried into content that does not hold it: %+v", decided)
		}
	}
}
