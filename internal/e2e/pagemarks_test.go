package e2e_test

import (
	"strings"
	"testing"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
	"github.com/google/go-cmp/cmp"
)

type markRead struct {
	Label    string `json:"label"`
	Tabindex string `json:"tabindex"`
	Said     string `json:"said"`
	Variant  string `json:"variant"`
	Tone     string `json:"tone"`
	Icon     string `json:"icon"`
	Severity string `json:"severity"`
	Audience string `json:"audience"`
}

func markReading(within, kind string) string {
	return `(() => {
  const mark = document.querySelector(` + jsonText(within+" .mark."+kind) + `);
  if (!mark) return null;
  return {
    label: mark.getAttribute('aria-label') || '',
    tabindex: mark.getAttribute('tabindex') || '',
    said: mark.textContent.trim(),
    variant: mark.variant,
    tone: mark.tone,
    icon: mark.icon,
    severity: mark.dataset.severity || '',
    audience: mark.dataset.audience || '',
  };
})()`
}

func jsonText(value string) string {
	return "'" + strings.ReplaceAll(value, "'", `\'`) + "'"
}

func (r *revisions) markOf(within, kind string) *markRead {
	r.t.Helper()

	var read *markRead
	r.page.evaluate(markReading(within, kind), &read)
	return read
}

func (r *revisions) hovering(within, kind string) tipped {
	r.t.Helper()

	r.page.do(
		sending(emulation.SetFocusEmulationEnabled, emulation.SetFocusEmulationEnabledParams{Enabled: true}),
		chromedp.Focus(chromedp.CSS(within+" .mark."+kind)),
	)
	var said tipped
	r.page.evaluate(whatTheTipSays, &said)
	return said
}

func severityStated(severity string) sketch {
	return sketch{"severity": severity}
}

func audienceStated(audience string) sketch {
	return sketch{"audience": audience}
}

func oneCacheDocumentWith(stated sketch) sketch {
	decision := aCacheDecision("single-choice")
	delete(decision, "severity")
	delete(decision, "audience")
	for key, value := range stated {
		decision[key] = value
	}
	return aDocument("one thing to settle", nil, decision)
}

func TestASeverityMarkSaysWhatTheSeverityMeans(t *testing.T) {
	t.Parallel()
	shown := showing(t, oneCacheDocumentWith(severityStated("error")))

	for _, stated := range []struct {
		severity string
		spoken   []string
		hovered  []string
	}{
		{"error", []string{"Severity: Error", "should not land"}, []string{"Severity: Error", "should not land"}},
		{"warning", []string{"Severity: Warning", "may be fine"}, []string{"Severity: Warning", "may be fine"}},
		{"info", []string{"Severity: Info", "No action expected"}, []string{"Severity: Info", "No action expected"}},
		{"cosmetic", []string{"Severity: cosmetic"}, []string{"Severity: cosmetic"}},
	} {
		shown.show(oneCacheDocumentWith(severityStated(stated.severity)))

		mark := shown.markOf(".decision", "severity")
		if mark == nil {
			t.Fatalf("the card carries no severity mark for %s", stated.severity)
		}
		for _, said := range stated.spoken {
			if !strings.Contains(mark.Label, said) {
				t.Errorf("the %s mark is announced as %q, and it says %q", stated.severity, mark.Label, said)
			}
		}
		if mark.Tabindex != "0" {
			t.Errorf("the %s mark on a card cannot be reached by keyboard, and its meaning is shown only on hover or focus", stated.severity)
		}

		hovered := shown.hovering(".decision", "severity")
		if !hovered.Shown {
			t.Errorf("the %s mark shows nothing while it has the cursor", stated.severity)
		}
		for _, said := range stated.hovered {
			if !strings.Contains(hovered.Said, said) {
				t.Errorf("the %s mark shows %q while it has the cursor, and it says %q", stated.severity, hovered.Said, said)
			}
		}
	}
}

func TestAMarkOnACardSaysItsNameAndAMarkInTheIndexOnlyDrawsIt(t *testing.T) {
	t.Parallel()
	shown := showing(t, oneCacheDocumentWith(severityStated("error")))

	for _, stated := range []struct {
		kind  string
		value string
		said  string
	}{
		{"severity", "error", "Error"},
		{"severity", "warning", "Warning"},
		{"severity", "info", "Info"},
		{"severity", "cosmetic", "cosmetic"},
		{"audience", "needs a person", "Your call"},
		{"audience", "can be patched", "Author can patch"},
		{"audience", "informational", "For information"},
		{"audience", "the author patches", "Decides: the author patches"},
	} {
		shown.show(oneCacheDocumentWith(sketch{stated.kind: stated.value}))

		mark := shown.markOf(".decision", stated.kind)
		if mark == nil {
			t.Fatalf("the card carries no %s mark for %q", stated.kind, stated.value)
		}
		if mark.Said != stated.said {
			t.Errorf("the %s mark %q on a card says %q, want %q", stated.kind, stated.value, mark.Said, stated.said)
		}

		if stated.kind != "severity" {
			continue
		}
		drawn := shown.markOf(".entry", "severity")
		if drawn == nil {
			t.Fatalf("the index entry carries no severity mark for %q", stated.value)
		}
		if drawn.Said != "" {
			t.Errorf("the severity mark in the index says %q, and the index has room only for the icon", drawn.Said)
		}
	}
}

func TestAnUnstatedSeverityOrAudienceIsNotMarked(t *testing.T) {
	t.Parallel()
	_, page := read(t, written(t, oneCacheDocumentWith(nil)))

	for _, kind := range []string{"severity", "audience"} {
		var marks int
		page.evaluate(`document.querySelectorAll('.decision .mark.`+kind+`').length`, &marks)
		if marks != 0 {
			t.Errorf("a decision stating no %s carries %d %s marks", kind, marks, kind)
		}
	}
}

func TestAnAudienceSocratesReviewWritesIsMarkedAsTheDesignDrewIt(t *testing.T) {
	t.Parallel()
	shown := showing(t, oneCacheDocumentWith(audienceStated("needs a person")))

	for _, stated := range []struct {
		audience string
		tone     string
		spoken   string
		hovered  string
	}{
		{"needs a person", "accent", "Your call. Only a person can decide this. The reviewer will not answer it.", "Only a person can decide this. The reviewer will not answer it."},
		{"can be patched", "neutral", "Author can patch. The session that wrote the change can fix this without you.", "The session that wrote the change can fix this without you."},
		{"informational", "neutral", "For information. Nothing to decide.", "Nothing to decide."},
		{"the author patches", "neutral", "Decides: the author patches", ""},
	} {
		shown.show(oneCacheDocumentWith(audienceStated(stated.audience)))

		mark := shown.markOf(".decision", "audience")
		if mark == nil {
			t.Fatalf("the card carries no audience mark for %q", stated.audience)
		}
		if mark.Label != stated.spoken {
			t.Errorf("the audience %q is announced as %q, want %q", stated.audience, mark.Label, stated.spoken)
		}
		if mark.Tone != stated.tone {
			t.Errorf("the audience %q is drawn in the %q tone, want %q", stated.audience, mark.Tone, stated.tone)
		}
		if mark.Audience != stated.audience {
			t.Errorf("the audience mark carries %q, and the document stated %q", mark.Audience, stated.audience)
		}

		hovered := shown.hovering(".decision", "audience")
		if !hovered.Shown || !strings.Contains(hovered.Said, stated.hovered) {
			t.Errorf("the audience %q shows %+v while it has the cursor, and it says %q", stated.audience, hovered, stated.hovered)
		}
	}
}

func TestEachAudienceSocratesReviewWritesIsDrawnWithAnIconOfItsOwn(t *testing.T) {
	t.Parallel()
	shown := showing(t, oneCacheDocumentWith(audienceStated("needs a person")))

	drawn := map[string]string{}
	for _, audience := range []string{"needs a person", "can be patched", "informational", "the author patches"} {
		shown.show(oneCacheDocumentWith(audienceStated(audience)))

		mark := shown.markOf(".decision", "audience")
		if mark == nil || mark.Icon == "" {
			t.Fatalf("the %q mark names no icon", audience)
		}
		for other, held := range drawn {
			if held == mark.Icon {
				t.Errorf("%q and %q are drawn alike, and each is told apart at a glance", audience, other)
			}
		}
		drawn[audience] = mark.Icon
	}
}

func TestEveryIconAMarkNamesIsOneTheBadgeDraws(t *testing.T) {
	t.Parallel()
	shown := showing(t, oneCacheDocumentWith(audienceStated("needs a person")))

	const drawnIcons = `(async () => {
  const badges = Array.from(document.querySelectorAll('ui-badge.mark'));
  await Promise.all(badges.map((held) => held.updateComplete));
  return badges.map((held) => ({ icon: held.icon, drawn: !!held.shadowRoot.querySelector('svg') }));
})()`

	for _, severity := range []string{"error", "warning", "info", "cosmetic"} {
		for _, audience := range []string{"needs a person", "can be patched", "informational", "the author patches", "somebody else"} {
			shown.show(oneCacheDocumentWith(sketch{"severity": severity, "audience": audience}))

			var badges []struct {
				Icon  string `json:"icon"`
				Drawn bool   `json:"drawn"`
			}
			shown.page.evaluate(drawnIcons, &badges, awaitingThePromise)
			if len(badges) == 0 {
				t.Fatalf("the %s mark for %q draws no badge, so this proves nothing", severity, audience)
			}
			for _, badge := range badges {
				if !badge.Drawn {
					t.Errorf("the %s mark for %q names the icon %q, and the badge draws none of that name", severity, audience, badge.Icon)
				}
			}
		}
	}
}

func TestASeverityIsDrawnFilledAndAnAudienceOutlinedSoHowSeriousAndWhoDecidesReadApart(t *testing.T) {
	t.Parallel()
	shown := showing(t, oneCacheDocumentWith(nil))

	for _, stated := range []struct{ severity, audience string }{
		{"error", "needs a person"},
		{"warning", "can be patched"},
		{"info", "informational"},
		{"cosmetic", "the author patches"},
	} {
		shown.show(oneCacheDocumentWith(sketch{"severity": stated.severity, "audience": stated.audience}))

		if mark := shown.markOf(".decision", "severity"); mark == nil || mark.Variant != "filled" {
			t.Errorf("the %s mark is drawn %+v, want filled", stated.severity, mark)
		}
		if mark := shown.markOf(".decision", "audience"); mark == nil || mark.Variant != "outlined" {
			t.Errorf("the %q mark is drawn %+v, want outlined", stated.audience, mark)
		}
	}
}

func TestAnIndexEntryWhoseDecisionIsYourCallCarriesThatMarkAndNoOtherAudience(t *testing.T) {
	t.Parallel()
	shown := showing(t, oneCacheDocumentWith(audienceStated("needs a person")))

	for _, stated := range []struct {
		audience string
		marked   bool
	}{
		{"needs a person", true},
		{"can be patched", false},
		{"informational", false},
		{"the author patches", false},
	} {
		shown.show(oneCacheDocumentWith(audienceStated(stated.audience)))

		drawn := shown.markOf(".entry", "audience")
		if !stated.marked {
			if drawn != nil {
				t.Errorf("the index marks %q, and only a decision that is the person's own call is marked there", stated.audience)
			}
			continue
		}
		if drawn == nil {
			t.Fatalf("the index does not mark %q, and a person scans the index for what is theirs to decide", stated.audience)
		}
		if drawn.Said != "" {
			t.Errorf("the audience mark in the index says %q, and the index has room only for the icon", drawn.Said)
		}
		if drawn.Variant != "outlined" {
			t.Errorf("the audience mark in the index is drawn %q, want outlined as on the card", drawn.Variant)
		}
	}
}

func TestAnIndexEntryReadFromTheSourceMarksYourCallToo(t *testing.T) {
	t.Parallel()
	decision := aDecision("c-1", nil)
	shown := showing(t, aDocument("one thing to settle", []sketch{anEntry("c-1", "", "error", "needs a person", "")}, decision))

	if shown.markOf(".entry", "audience") == nil {
		t.Error("an entry the source lists as needing a person carries no Your call mark in the index")
	}
}

func TestTheIndexCountsWhatArrivesChosenAsAnswered(t *testing.T) {
	t.Parallel()
	recommended := aDecision("c-3", sketch{"options": []sketch{{"label": "leave it", "meaning": "acceptable here", "command": "agree c-3", "recommended": true}}})
	_, page := read(t, written(t, aDocument("three things to settle", leastSevereFirstSource(), recommended, aDecision("c-2", nil), aDecision("c-1", nil))))

	var said string
	page.evaluate(`document.querySelector('.progress-said').textContent.replace(/\s+/g, ' ').trim()`, &said)
	if said != "1 of 3 answered" {
		t.Errorf("the index says %q, and one of three decisions arrives with its recommendation chosen", said)
	}

	var answered map[string]string
	page.evaluate(`Object.fromEntries(Array.from(document.querySelectorAll('a.entry-link')).map((held) => [held.dataset.decision, held.querySelector('.answered').dataset.answered]))`, &answered)
	if diff := cmp.Diff(map[string]string{"c-1": "no", "c-2": "no", "c-3": "yes"}, answered); diff != "" {
		t.Errorf("what the index marks as answered differs (-want +got):\n%s", diff)
	}
}
