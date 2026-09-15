package render_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/agora/document"
	"github.com/socrates-tools/socrates/agora/render"
	"golang.org/x/net/html"
)

const patched = `--- a/internal/cache/write.go
+++ b/internal/cache/write.go
@@ -10,5 +10,7 @@ func Write(key string, entry Entry) {
 func Write(key string, entry Entry) {
-	store[key] = entry
+	mu.Lock()
+	defer mu.Unlock()
+	store[key] = entry
 }
`

func TestADiffRendersFromAPatch(t *testing.T) {
	read := oneDecision(document.One)
	read.Decisions[0].Evidence = document.Evidence{Patch: patched}

	page := rendered(t, read)

	shown := elements(page, "div")
	var diff *html.Node
	for _, held := range shown {
		if hasClass(held, "diff") {
			diff = held
		}
	}
	if diff == nil {
		t.Fatal("a decision carrying a patch renders no diff")
	}

	for _, named := range []string{"internal/cache/write.go", "func Write(key string, entry Entry) {"} {
		if !strings.Contains(text(diff), named) {
			t.Errorf("the diff does not carry %q:\n%s", named, text(diff))
		}
	}

	sides := map[string][]string{}
	for _, line := range elements(diff, "td") {
		if !hasClass(line, "line") {
			continue
		}
		sides[attribute(line, "data-side")+" "+attribute(line, "data-kind")] = append(
			sides[attribute(line, "data-side")+" "+attribute(line, "data-kind")], text(line))
	}

	if len(sides["left removed"]) != 1 {
		t.Errorf("the left side shows %v as removed, want one line", sides["left removed"])
	}
	if len(sides["right added"]) != 3 {
		t.Errorf("the right side shows %v as added, want three lines", sides["right added"])
	}
	if len(sides["left context"]) != 2 || len(sides["right context"]) != 2 {
		t.Errorf("the two sides show %v and %v as unchanged, want two lines each", sides["left context"], sides["right context"])
	}
	if len(sides["left absent"]) != 2 {
		t.Errorf("the left side shows %v where the right side gained lines, want two blanks", sides["left absent"])
	}
}

func TestAPatchIsReadIntoBothSidesWithTheirNumbers(t *testing.T) {
	// scenario: A diff renders from a patch
	read := render.ReadPatch(patched)

	if read.From != "internal/cache/write.go" || read.To != "internal/cache/write.go" {
		t.Errorf("the diff names %q and %q, and a person reads the path they know rather than the one the patch format prefixes", read.From, read.To)
	}
	if len(read.Hunks) != 1 {
		t.Fatalf("the patch read into %d hunks, want one", len(read.Hunks))
	}

	want := []render.Row{
		{
			Left:  render.Side{Kind: render.Context, Number: 10, Text: "func Write(key string, entry Entry) {"},
			Right: render.Side{Kind: render.Context, Number: 10, Text: "func Write(key string, entry Entry) {"},
		},
		{
			Left:  render.Side{Kind: render.Removed, Number: 11, Text: "\tstore[key] = entry"},
			Right: render.Side{Kind: render.Added, Number: 11, Text: "\tmu.Lock()"},
		},
		{
			Left:  render.Side{Kind: render.Absent},
			Right: render.Side{Kind: render.Added, Number: 12, Text: "\tdefer mu.Unlock()"},
		},
		{
			Left:  render.Side{Kind: render.Absent},
			Right: render.Side{Kind: render.Added, Number: 13, Text: "\tstore[key] = entry"},
		},
		{
			Left:  render.Side{Kind: render.Context, Number: 12, Text: "}"},
			Right: render.Side{Kind: render.Context, Number: 14, Text: "}"},
		},
	}
	if diff := cmp.Diff(want, read.Hunks[0].Rows); diff != "" {
		t.Errorf("what the patch read into differs (-want +got):\n%s", diff)
	}
}

func TestAnEmptyPatchRendersNoDiff(t *testing.T) {
	// scenario: A diff renders from a patch
	page := rendered(t, oneDecision(document.One))

	for _, held := range elements(page, "div") {
		if hasClass(held, "diff") {
			t.Error("a decision carrying no patch renders a diff")
		}
	}
}
