package surface_test

import (
	"net/url"
	"path"
	"testing"

	"github.com/socrates-tools/socrates-review/internal/surface"
)

func TestTwoMintedSessionsDifferAndEachIsAPathSegment(t *testing.T) {
	seen := map[string]bool{}

	for range 256 {
		minted, err := surface.MintSession()
		if err != nil {
			t.Fatalf("minting a session: %v", err)
		}
		if seen[minted] {
			t.Fatalf("minting produced %q twice, and two callers must not land on the same page", minted)
		}
		seen[minted] = true

		if minted == "" {
			t.Fatal("minting produced an empty session, and a session is a path segment")
		}
		if escaped := url.PathEscape(minted); escaped != minted {
			t.Errorf("the session %q is written into a URL as %q, and a session is a safe path segment by construction", minted, escaped)
		}
		if cleaned := path.Clean("/" + minted); cleaned != "/"+minted {
			t.Errorf("the session %q reads as the path %q, and a session names one segment", minted, cleaned)
		}
	}
}
