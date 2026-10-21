package render_test

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/socrates-tools/socrates/agora/document"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

func TestAnAudienceSocratesReviewWritesIsMarkedAsTheDesignDrewIt(t *testing.T) {
	for _, stated := range []struct {
		audience string
		tone     string
		spoken   string
	}{
		{"needs a person", "accent", "Your call. Only a person can decide this. The reviewer will not answer it."},
		{"can be patched", "neutral", "Author can patch. The session that wrote the change can fix this without you."},
		{"informational", "neutral", "For information. Nothing to decide."},
		{"the author patches", "neutral", "Decides: the author patches"},
	} {
		t.Run(stated.audience, func(t *testing.T) {
			read := oneDecision(document.One)
			read.Decisions[0].Audience = stated.audience

			mark := markIn(cardsIn(rendered(t, read))[0], "audience")
			if mark == nil {
				t.Fatal("the card carries no audience mark")
			}
			if spoken := attribute(mark, "aria-label"); spoken != stated.spoken {
				t.Errorf("the audience %q is announced as %q, want %q", stated.audience, spoken, stated.spoken)
			}
			if tone := attribute(mark, "tone"); tone != stated.tone {
				t.Errorf("the audience %q is drawn in the %q tone, want %q", stated.audience, tone, stated.tone)
			}
			if held := attribute(mark, "data-audience"); held != stated.audience {
				t.Errorf("the audience mark carries %q, and the document stated %q", held, stated.audience)
			}
		})
	}
}

func TestEachAudienceSocratesReviewWritesIsDrawnWithAnIconOfItsOwn(t *testing.T) {
	drawn := map[string]string{}
	for _, audience := range []string{"needs a person", "can be patched", "informational", "the author patches"} {
		read := oneDecision(document.One)
		read.Decisions[0].Audience = audience
		mark := markIn(cardsIn(rendered(t, read))[0], "audience")
		named := attribute(mark, "icon")
		if named == "" {
			t.Fatalf("the %q mark names no icon", audience)
		}
		for other, held := range drawn {
			if held == named {
				t.Errorf("%q and %q are drawn alike, and each is told apart at a glance", audience, other)
			}
		}
		drawn[audience] = named
	}
}

func TestEveryIconAMarkNamesIsOneTheBadgeDraws(t *testing.T) {
	// scenario: The icons the marks of every severity and audience name are all in the badge's own set
	held, err := os.ReadFile(filepath.Join(source.Root(t), "agora", "web", "packages", "ui", "badge", "icons.ts"))
	if err != nil {
		t.Fatalf("reading the badge's icons: %v", err)
	}
	drawn := map[string]bool{}
	for _, found := range regexp.MustCompile(`(?m)^  '?([a-z-]+)'?:`).FindAllStringSubmatch(string(held), -1) {
		drawn[found[1]] = true
	}
	if len(drawn) == 0 {
		t.Fatal("the badge draws no icon, so this proves nothing")
	}

	read := oneDecision(document.One)
	for _, severity := range []string{"error", "warning", "info", "cosmetic"} {
		for _, audience := range []string{"needs a person", "can be patched", "informational", "the author patches", "somebody else"} {
			read.Decisions[0].Severity = severity
			read.Decisions[0].Audience = audience
			for _, badge := range elements(rendered(t, read), "ui-badge") {
				if named := attribute(badge, "icon"); !drawn[named] {
					t.Errorf("the %s mark for %q names the icon %q, and the badge draws none of that name", severity, audience, named)
				}
			}
		}
	}
}
