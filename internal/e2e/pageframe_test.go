package e2e_test

import (
	"strings"
	"testing"

	"github.com/chromedp/chromedp"
)

func illustrated(decision sketch, path string) sketch {
	decision["illustration"] = sketch{"path": path}
	return aDocument("one thing to settle", nil, decision)
}

func withADrawing(t *testing.T, decision sketch) (*surfaced, *reader) {
	t.Helper()

	held := standing(t)
	aDrawing(t, held, "drawings/q3.html", "<!doctype html><p>a drawing</p>")
	return held, opening(t, held.serve(written(t, illustrated(decision, "drawings/q3.html"))))
}

type drawnFrame struct {
	Sandbox string `json:"sandbox"`
	Source  string `json:"source"`
	Title   string `json:"title"`
}

func framesOf(page *reader) []drawnFrame {
	page.t.Helper()

	var drawn []drawnFrame
	page.evaluate(`Array.from(document.querySelectorAll('iframe'), (held) => ({ sandbox: held.getAttribute('sandbox') || '', source: held.getAttribute('src') || '', title: held.getAttribute('title') || '' }))`, &drawn)
	return drawn
}

func TestThePageCarriesTheBrandRatherThanTheNameOfTheServer(t *testing.T) {
	t.Parallel()
	_, page := read(t, twoDecisionDocument)

	var drawn struct {
		Wordmark string `json:"wordmark"`
		Body     string `json:"body"`
	}
	page.evaluate(`({
  wordmark: (document.querySelector('.wordmark') || { textContent: '' }).textContent.trim(),
  body: document.body.textContent.toLowerCase(),
})`, &drawn)

	if drawn.Wordmark != "socrates" {
		t.Errorf("the wordmark says %q, want socrates", drawn.Wordmark)
	}
	if strings.Contains(drawn.Body, "agora") {
		t.Error("the page shows a person the name of the server that serves it")
	}
}

func TestAnExchangeThatIsOverSaysSoUnderTheTitle(t *testing.T) {
	t.Parallel()
	held, page := read(t, twoDecisionDocument)

	const whatTheNoticeIs = `(() => {
  const notice = document.querySelector('.over-notice');
  return {
    hidden: notice.hidden,
    said: notice.textContent.trim(),
    underTheTitle: !!notice.closest('header.heading'),
  };
})()`
	type notice struct {
		Hidden        bool   `json:"hidden"`
		Said          string `json:"said"`
		UnderTheTitle bool   `json:"underTheTitle"`
	}

	var open notice
	page.evaluate(whatTheNoticeIs, &open)
	if !open.Hidden {
		t.Error("a page that is open says the exchange is over")
	}

	held.end(held.opened)
	page.reload()

	var over notice
	page.evaluate(whatTheNoticeIs, &over)
	if over.Hidden {
		t.Error("a page drawn over hides the notice that says so")
	}
	if !strings.Contains(strings.ToLower(over.Said), "over") {
		t.Errorf("the notice says %q, and it says the exchange is over", over.Said)
	}
	if !over.UnderTheTitle {
		t.Error("the notice that the exchange is over is not under the title")
	}
}

func TestAnIllustrationIsAFrameSandboxedToScriptAlone(t *testing.T) {
	t.Parallel()
	held, page := withADrawing(t, aCacheDecision("single-choice"))

	frames := framesOf(page)
	if len(frames) != 1 {
		t.Fatalf("a decision carrying an illustration drew %d frames", len(frames))
	}
	if frames[0].Sandbox != "allow-scripts" {
		t.Errorf("the frame is sandboxed with %q, want allow-scripts and nothing else, so the drawing runs in an opaque origin", frames[0].Sandbox)
	}
	if want := "/sessions/" + held.opened + "/illustrations/c-1/q3.html"; frames[0].Source != want {
		t.Errorf("the frame is drawn from %q, want %q", frames[0].Source, want)
	}
	if frames[0].Title == "" {
		t.Error("the frame carries no title, and a screen reader has no name for it")
	}
}

func TestTheFrameSourceEscapesTheDecisionIdentifierAndKeepsTheFileName(t *testing.T) {
	t.Parallel()
	decision := aCacheDecision("single-choice")
	decision["id"] = "a b/c?d"
	held, page := withADrawing(t, decision)

	frames := framesOf(page)
	if len(frames) != 1 {
		t.Fatalf("a decision carrying an illustration drew %d frames", len(frames))
	}
	if want := "/sessions/" + held.opened + "/illustrations/a%20b%2Fc%3Fd/q3.html"; frames[0].Source != want {
		t.Errorf("the frame is drawn from %q, want %q", frames[0].Source, want)
	}
}

func TestADecisionWithNoIllustrationDrawsNoFrame(t *testing.T) {
	t.Parallel()
	_, page := read(t, written(t, oneCacheDocument("single-choice")))

	if frames := framesOf(page); len(frames) != 0 {
		t.Errorf("a decision carrying no illustration drew %d frames", len(frames))
	}
}

func policyOfThePage(page *reader) string {
	page.t.Helper()

	var policy string
	page.evaluate(`document.querySelector('meta[http-equiv="Content-Security-Policy" i]').content`, &policy)
	return policy
}

func TestThePolicyAdmitsTheFrameAndEmbedsNothingACallerWrote(t *testing.T) {
	t.Parallel()
	_, page := withADrawing(t, aCacheDecision("single-choice"))
	policy := policyOfThePage(page)

	if !strings.Contains(policy, "frame-src 'self'") {
		t.Errorf("the policy does not admit the frame the illustration is drawn in:\n%s", policy)
	}
	for _, admitted := range []string{"unsafe-inline", "unsafe-hashes", "unsafe-eval", "sha256-", "*"} {
		if strings.Contains(policy, admitted) {
			t.Errorf("the policy admits %q, and the page embeds nothing a caller wrote:\n%s", admitted, policy)
		}
	}
	if !strings.Contains(policy, "style-src-attr 'none'") {
		t.Errorf("the policy admits style carried on an element:\n%s", policy)
	}

	var embedded int
	page.evaluate(`document.querySelectorAll('style').length`, &embedded)
	if embedded != 0 {
		t.Errorf("the page embeds %d stylesheets, and a caller's drawing is in a frame of its own", embedded)
	}
}

func TestThePolicyAdmitsBlobOnlyForImages(t *testing.T) {
	t.Parallel()
	_, page := read(t, written(t, oneCacheDocument("single-choice")))

	for directive := range strings.SplitSeq(policyOfThePage(page), ";") {
		directive = strings.TrimSpace(directive)
		admitsBlob := strings.Contains(directive, "blob:")
		isImageDirective := strings.HasPrefix(directive, "img-src ")
		if admitsBlob && !isImageDirective {
			t.Errorf("the policy admits blob: in %q, and only images are previewed from memory", directive)
		}
		if isImageDirective && !admitsBlob {
			t.Errorf("the policy does not admit blob: for images in %q", directive)
		}
	}
}

func TestThePageCarriesTheBaseItsRequestsAreMadeAgainst(t *testing.T) {
	t.Parallel()
	held, page := read(t, written(t, oneCacheDocument("single-choice")))

	var carried string
	page.evaluate(`document.body.dataset.base`, &carried)
	if want := "/sessions/" + held.opened; carried != want {
		t.Errorf("the page carries the base %q, and it was served under %q", carried, want)
	}
}

func TestAPieceKeepsTheClassThatNamesItBesideTheUtilitiesThatLayItOut(t *testing.T) {
	t.Parallel()
	_, page := read(t, written(t, oneCacheDocument("single-choice")))

	var laid struct {
		CardOverflow     string `json:"cardOverflow"`
		OptionsDisplay   string `json:"optionsDisplay"`
		ComposerDisplay  string `json:"composerDisplay"`
		ComposerDirected string `json:"composerDirected"`
	}
	page.evaluate(`(() => {
  const card = getComputedStyle(document.querySelector('ui-card.decision'));
  const options = getComputedStyle(document.querySelector('ui-choice-group.options-list'));
  const composer = getComputedStyle(document.querySelector('ui-card.composer'));
  return { cardOverflow: card.overflow, optionsDisplay: options.display, composerDisplay: composer.display, composerDirected: composer.flexDirection };
})()`, &laid)

	if laid.CardOverflow != "hidden" {
		t.Errorf("the decision lays itself out with overflow %q, and its layout is written beside its name", laid.CardOverflow)
	}
	if laid.OptionsDisplay != "flex" && laid.OptionsDisplay != "grid" {
		t.Errorf("the options list lays itself out as %q, and its layout is written beside its name", laid.OptionsDisplay)
	}
	if laid.ComposerDisplay != "flex" || laid.ComposerDirected != "column" {
		t.Errorf("the composer lays itself out as %q %q, and its layout is written beside its name", laid.ComposerDisplay, laid.ComposerDirected)
	}
}

func TestThePanelTakesItsLayoutFromThePageAndGivesWayWhenNothingWasSent(t *testing.T) {
	t.Parallel()
	held := standing(t)
	page := opening(t, held.serve(twoDecisionDocument), chromedp.WindowSize(aWindowWideEnoughForThePanel, 900))

	var open bool
	page.evaluate(`document.querySelector('aside.panel').checkVisibility()`, &open)
	if !open {
		t.Fatal("an open exchange shows no panel beside the decisions")
	}

	held.end(held.opened)
	page.reload()

	var shown bool
	page.evaluate(`document.querySelector('aside.panel').checkVisibility()`, &shown)
	if shown {
		t.Error("an exchange that is over with nothing sent shows a panel, and the shell's part takes its layout from the page")
	}
}

func TestAnIndexOfEntriesNoneOfWhichIsDecidedLaysOutOneColumnAndHidesTheMarks(t *testing.T) {
	t.Parallel()
	held := standing(t)
	page := openingNothingToAnswer(t, held.serve(written(t, aDocument("what was seen", []sketch{anEntry("s-1", "Cache", "info", "", "a.go:1")}))), chromedp.WindowSize(aWindowWideEnoughForThePanel, 900))

	var laid struct {
		Columns int    `json:"columns"`
		Mark    string `json:"mark"`
	}
	page.evaluate(`(() => {
  const link = document.querySelector('span.entry-link');
  return {
    columns: getComputedStyle(link).gridTemplateColumns.split(' ').length,
    mark: getComputedStyle(document.querySelector('.entry .answered')).display,
  };
})()`, &laid)

	if laid.Columns != 1 {
		t.Errorf("an undecided entry lays out %d columns, and its index has nothing to track", laid.Columns)
	}
	if laid.Mark != "none" {
		t.Errorf("an undecided entry's mark displays as %q, and it is hidden where nothing is tracked", laid.Mark)
	}
}

func TestAnIndexWithADecidedEntryLaysOutTheMarkColumnAndShowsTheMarks(t *testing.T) {
	t.Parallel()
	held := standing(t)
	page := opening(t, held.serve(written(t, aDocument("one thing to settle", []sketch{anEntry("c-1", "Cache", "error", "", "")}, aCacheDecision("single-choice")))), chromedp.WindowSize(aWindowWideEnoughForThePanel, 900))

	var laid struct {
		Columns int    `json:"columns"`
		Mark    string `json:"mark"`
	}
	page.evaluate(`(() => {
  const link = document.querySelector('a.entry-link');
  return {
    columns: getComputedStyle(link).gridTemplateColumns.split(' ').length,
    mark: getComputedStyle(link.querySelector('.answered')).display,
  };
})()`, &laid)

	if laid.Columns != 2 {
		t.Errorf("a decided entry lays out %d columns, and its index tracks answers in a column of its own", laid.Columns)
	}
	if laid.Mark == "none" {
		t.Error("a tracked entry's mark is hidden, and it is shown")
	}
}
