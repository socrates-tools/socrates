package e2e_test

import (
	"fmt"
	"testing"

	"github.com/chromedp/chromedp"
)

const (
	aWindowWideEnoughForEveryShellState = 1400
	aWindowTallEnoughForEveryShellState = 1000
	shellFrameNarrow                    = "iframe.catalog-frame-480"
	shellFrameAtNineSixty               = "iframe.catalog-frame-960"
	shellFrameAtTwelveEighty            = "iframe.catalog-frame-1280"
)

func shellCatalogued(t *testing.T) *reader {
	t.Helper()

	held := servingSurfaceBuiltWith(t, developmentBuildTag)
	address := held.serve(oneDecisionDocument)
	page := openingFor(t, originOf(t, address)+catalogRoute, "iframe.catalog-frame", ".catalog-state",
		chromedp.WindowSize(aWindowWideEnoughForEveryShellState, aWindowTallEnoughForEveryShellState))
	for _, frame := range []string{shellFrameNarrow, shellFrameAtNineSixty, shellFrameAtTwelveEighty} {
		page.waitFor(fmt.Sprintf(`(() => {
  const inside = document.querySelector(%q).contentWindow;
  return !!inside.customElements.get('ui-app-shell') && !!inside.document.querySelector('ui-app-shell [data-part="content"]');
})()`, frame))
	}
	return page
}

func inFrame(frame string) string {
	return fmt.Sprintf(`document.querySelector(%q).contentWindow`, frame)
}

func shellPart(frame, part string) string {
	return fmt.Sprintf(`%s.document.querySelector('ui-app-shell > [data-part="frame"] > [data-part="body"] > [data-part="%s"]')`, inFrame(frame), part)
}

func howTheShellStateScrolls(frame string) string {
	return fmt.Sprintf(`(() => {
  const inside = %s;
  const navigation = %s;
  const content = %s;
  const panel = %s;
  const root = inside.document.documentElement;
  const scrollsOnItsOwn = (held) => inside.getComputedStyle(held).overflowY === 'auto' && held.scrollHeight > held.clientHeight;
  return {
    windowScrolls: root.scrollHeight > inside.innerHeight + 1,
    navigationScrolls: scrollsOnItsOwn(navigation),
    contentScrolls: scrollsOnItsOwn(content),
    navigationBesideContent: navigation.getBoundingClientRect().right <= content.getBoundingClientRect().left + 1,
    navigationAboveContent: navigation.getBoundingClientRect().bottom <= content.getBoundingClientRect().top + 1,
    panelShown: panel.checkVisibility(),
    panelAtTheRight: Math.abs(panel.getBoundingClientRect().right - root.clientWidth) <= 1,
  };
})()`, inFrame(frame), shellPart(frame, "navigation"), shellPart(frame, "content"), shellPart(frame, "panel"))
}

type shellStateScrolling struct {
	WindowScrolls           bool `json:"windowScrolls"`
	NavigationScrolls       bool `json:"navigationScrolls"`
	ContentScrolls          bool `json:"contentScrolls"`
	NavigationBesideContent bool `json:"navigationBesideContent"`
	NavigationAboveContent  bool `json:"navigationAboveContent"`
	PanelShown              bool `json:"panelShown"`
	PanelAtTheRight         bool `json:"panelAtTheRight"`
}

func readingTheShellState(page *reader, frame string) shellStateScrolling {
	page.t.Helper()

	var read shellStateScrolling
	page.evaluate(howTheShellStateScrolls(frame), &read)
	return read
}

func TestTheCatalogDrawsTheShellBelowNineSixtyStackedAndScrollingWithTheWindow(t *testing.T) {
	// scenario: Below 960px the shell stacks the navigation over the content and nothing in it scrolls on its own
	t.Parallel()
	page := shellCatalogued(t)

	read := readingTheShellState(page, shellFrameNarrow)

	if !read.NavigationAboveContent || read.NavigationBesideContent {
		t.Errorf("in a window narrower than 960px the navigation sits %+v, and it stacks above the content", read)
	}
	if read.NavigationScrolls || read.ContentScrolls {
		t.Errorf("in a window narrower than 960px a part scrolls on its own: %+v", read)
	}
	if !read.WindowScrolls {
		t.Errorf("in a window narrower than 960px the window does not scroll: %+v, and the stacked parts scroll with it", read)
	}
	if read.PanelShown {
		t.Error("in a window narrower than 960px the panel is shown, and it comes from 1280px")
	}
}

func TestTheCatalogDrawsTheShellFromNineSixtyWithOnlyTheNavigationAndTheContentScrolling(t *testing.T) {
	// scenario: From 960px the shell fills the window and only the navigation and the content scroll
	t.Parallel()
	page := shellCatalogued(t)

	for _, frame := range []string{shellFrameAtNineSixty, shellFrameAtTwelveEighty} {
		t.Run(frame, func(t *testing.T) {
			read := readingTheShellState(page, frame)

			if read.WindowScrolls {
				t.Errorf("the window scrolls: %+v, and only the navigation and the content do", read)
			}
			if !read.NavigationScrolls {
				t.Errorf("the navigation does not scroll on its own: %+v", read)
			}
			if !read.ContentScrolls {
				t.Errorf("the content does not scroll on its own: %+v", read)
			}
			if !read.NavigationBesideContent {
				t.Errorf("the navigation does not sit beside the content: %+v", read)
			}
		})
	}
}

func TestTheCatalogDrawsThePanelOnlyFromTwelveEighty(t *testing.T) {
	// scenario: The reply panel's place appears at the right of the shell from 1280px and not at 960px
	t.Parallel()
	page := shellCatalogued(t)

	atNineSixty := readingTheShellState(page, shellFrameAtNineSixty)
	atTwelveEighty := readingTheShellState(page, shellFrameAtTwelveEighty)

	if atNineSixty.PanelShown {
		t.Error("the panel is shown in a window 960px wide, and it comes from 1280px")
	}
	if !atTwelveEighty.PanelShown || !atTwelveEighty.PanelAtTheRight {
		t.Errorf("in a window 1280px wide the panel is %+v, and it sits at the right", atTwelveEighty)
	}
}

func markedSectionsOf(frame string) string {
	return fmt.Sprintf(`JSON.stringify(Array.from(%s.document.querySelectorAll('ui-app-shell > [data-part="frame"] > [data-part="body"] > [data-part="navigation"] a[data-section]'))
  .filter((link) => link.getAttribute('aria-current') === 'location')
  .map((link) => link.dataset.section))`, inFrame(frame))
}

func scrollTheContentTo(frame, section string) string {
	return fmt.Sprintf(`(() => {
  const content = %s;
  const target = %s.document.getElementById(%q);
  content.scrollTop += target.getBoundingClientRect().top - content.getBoundingClientRect().top;
})()`, shellPart(frame, "content"), inFrame(frame), section)
}

func TestTheShellMarksInTheNavigationTheSectionInView(t *testing.T) {
	// scenario: Scrolling the shell's content moves the navigation's mark to the section in view
	t.Parallel()
	page := shellCatalogued(t)

	for _, frame := range []string{shellFrameAtNineSixty, shellFrameAtTwelveEighty} {
		t.Run(frame, func(t *testing.T) {
			markedOnly := func(section string) string {
				return fmt.Sprintf(`%s === JSON.stringify([%q])`, markedSectionsOf(frame), section)
			}

			page.waitFor(markedOnly("section-1"))

			page.evaluate(scrollTheContentTo(frame, "section-5"), nil)
			page.waitFor(markedOnly("section-5"))

			page.evaluate(scrollTheContentTo(frame, "section-1"), nil)
			page.waitFor(markedOnly("section-1"))
		})
	}
}
