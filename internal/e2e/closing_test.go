package e2e_test

import (
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/agora/surface"
)

const somethingLiveOnThePage = `Array.from(document.querySelectorAll('button, .choice, textarea')).some((held) => !held.disabled)`

const everythingOnThePageIsLive = `Array.from(document.querySelectorAll('button, .choice, textarea')).every((held) => !held.disabled)`

const noticeIsShowing = `!document.querySelector('.notice').hidden`

const noticeIsHidden = `document.querySelector('.notice').hidden`

func pressPastTheNotice(page *reader, press string) {
	page.t.Helper()
	page.click(press)
	page.waitFor(noticeIsShowing)
	page.click(press)
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

	pressPastTheNotice(page, ".finish")
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

	pressPastTheNotice(page, ".send")

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

func TestEitherPressOverABlankDecisionNamesItAndSendsNothing(t *testing.T) {
	for _, press := range []string{".send", ".finish"} {
		t.Run(press, func(t *testing.T) {
			held, page := read(t, twoDecisionDocument)

			page.evaluate(chose("c-1", 0), nil)
			page.evaluate(chose("c-2", 2), nil)
			page.click(press)
			page.waitFor(noticeIsShowing)

			if delivered, arrived := held.awaiting(shortWindow); arrived {
				t.Fatalf("pressing %s over a blank decision delivered an answer: %+v", press, delivered)
			}

			var said string
			page.evaluate(`document.querySelector('.notice').textContent`, &said)
			if !strings.Contains(said, "c-2") {
				t.Errorf("the notice pressing %s raised reads %q, and it does not name the decision left blank", press, said)
			}
			if !strings.Contains(strings.ToLower(said), "again") {
				t.Errorf("the notice pressing %s raised reads %q, and it does not say in plain words that pressing again sends", press, said)
			}
		})
	}
}

func TestTheNoticeNamesEveryBlankDecisionAndEachNameLinksToItsCard(t *testing.T) {
	// scenario: An entry leads to the decision that covers it
	_, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-2", 2), nil)
	page.click(".send")
	page.waitFor(noticeIsShowing)

	var hrefs []string
	page.evaluate(`Array.from(document.querySelectorAll('.notice-decisions a')).map((held) => held.getAttribute('href'))`, &hrefs)

	want := []string{"#decision-c-1", "#decision-c-2"}
	if diff := cmp.Diff(want, hrefs); diff != "" {
		t.Errorf("the notice's links to the decisions left blank differ (-want +got):\n%s", diff)
	}

	page.click(`.notice-decisions a`)
	page.waitFor(`location.hash === '#decision-c-1'`)

	var reached string
	page.evaluate(`document.querySelector(location.hash).dataset.decision`, &reached)
	if reached != "c-1" {
		t.Errorf("following the notice's first link arrived at %q, want %q", reached, "c-1")
	}
}

func TestPressingAgainSendsCarryingEveryDecisionThatWasNotBlank(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.evaluate(chose("c-2", 2), nil)
	page.click(".send")
	page.waitFor(noticeIsShowing)

	if delivered, arrived := held.awaiting(shortWindow); arrived {
		t.Fatalf("the first press over a blank decision delivered an answer: %+v", delivered)
	}

	page.click(".send")

	delivered, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("pressing send again over the same blank decision delivered nothing")
	}
	if len(delivered.Decisions) != 1 || delivered.Decisions[0].ID != "c-1" {
		t.Errorf("pressing send again carried %+v, want only the decision that was not blank", delivered.Decisions)
	}

	var hidden bool
	page.evaluate(noticeIsHidden, &hidden)
	if !hidden {
		t.Error("sending after the notice was confirmed leaves the notice showing")
	}
}

func TestDecliningTheNoticeLeavesTheExchangeOpenAndSendsNothing(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.evaluate(chose("c-2", 2), nil)
	page.click(".send")
	page.waitFor(noticeIsShowing)

	page.click(".keep-answering")
	page.waitFor(noticeIsHidden)

	if delivered, arrived := held.awaiting(shortWindow); arrived {
		t.Fatalf("declining the notice delivered an answer anyway: %+v", delivered)
	}

	var accepting bool
	page.evaluate(everythingOnThePageIsLive, &accepting)
	if !accepting {
		t.Error("declining the notice left the page unable to take a further answer")
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

	var hidden bool
	page.evaluate(noticeIsHidden, &hidden)
	if !hidden {
		t.Error("a press with nothing blank raised the notice anyway")
	}
}

func TestPressingTheOtherPressWhileTheNoticeIsShowingRaisesItForThatPressInstead(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.evaluate(chose("c-2", 2), nil)
	page.click(".send")
	page.waitFor(noticeIsShowing)

	page.click(".finish")

	if delivered, arrived := held.awaiting(shortWindow); arrived {
		t.Fatalf("pressing the other press while the notice is showing delivered an answer: %+v", delivered)
	}

	var hidden bool
	page.evaluate(noticeIsHidden, &hidden)
	if hidden {
		t.Fatal("pressing the other press while the notice is showing closed the notice instead of raising it again")
	}

	page.click(".finish")
	delivered, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("pressing finish again after it raised the notice delivered nothing")
	}
	if !delivered.Ended {
		t.Error("pressing finish again does not end the exchange")
	}
}
