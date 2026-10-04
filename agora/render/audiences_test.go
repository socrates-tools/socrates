package render_test

import (
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/agora/document"
	"golang.org/x/net/html"
)

func TestAnAudienceSocratesReviewWritesIsMarkedAsTheDesignDrewIt(t *testing.T) {
	for _, stated := range []struct {
		audience string
		tone     string
		spoken   string
	}{
		{"needs a person", "clay", "Your call. Only a person can decide this. The reviewer will not answer it."},
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
			if tone := attribute(mark, "data-tone"); tone != stated.tone {
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
		icons := elements(mark, "svg")
		if len(icons) != 1 {
			t.Fatalf("the %q mark draws %d icons", audience, len(icons))
		}
		var inner strings.Builder
		if err := html.Render(&inner, icons[0]); err != nil {
			t.Fatalf("reading back the mark: %v", err)
		}
		for other, held := range drawn {
			if held == inner.String() {
				t.Errorf("%q and %q are drawn alike, and each is told apart at a glance", audience, other)
			}
		}
		drawn[audience] = inner.String()
	}
}
