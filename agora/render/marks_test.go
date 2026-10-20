package render_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/agora/document"
	"golang.org/x/net/html"
)

func TestASeverityMarkSaysWhatTheSeverityMeans(t *testing.T) {
	for _, stated := range []struct {
		severity string
		spoken   []string
	}{
		{"error", []string{"Severity: Error", "should not land"}},
		{"warning", []string{"Severity: Warning", "may be fine"}},
		{"info", []string{"Severity: Info", "No action expected"}},
		{"cosmetic", []string{"Severity: cosmetic"}},
	} {
		t.Run(stated.severity, func(t *testing.T) {
			read := oneDecision(document.One)
			read.Decisions[0].Severity = stated.severity

			mark := markIn(cardsIn(rendered(t, read))[0], "severity")
			if mark == nil {
				t.Fatal("the card carries no severity mark")
			}
			for _, said := range stated.spoken {
				if !strings.Contains(attribute(mark, "aria-label"), said) {
					t.Errorf("the %s mark is announced as %q, and it says %q", stated.severity, attribute(mark, "aria-label"), said)
				}
			}
			if attribute(mark, "tabindex") != "0" {
				t.Errorf("the %s mark on a card cannot be reached by keyboard, and its meaning is shown only on hover or focus", stated.severity)
			}
		})
	}
}

func TestAMarkOnACardSaysItsNameAndAMarkInTheIndexOnlyDrawsIt(t *testing.T) {
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
		t.Run(stated.value, func(t *testing.T) {
			read := oneDecision(document.One)
			read.Decisions[0].Severity = stated.value
			read.Decisions[0].Audience = stated.value

			page := rendered(t, read)
			mark := markIn(cardsIn(page)[0], stated.kind)
			if mark == nil {
				t.Fatalf("the card carries no %s mark", stated.kind)
			}
			if said := strings.TrimSpace(text(mark)); said != stated.said {
				t.Errorf("the %s mark %q on a card says %q, want %q", stated.kind, stated.value, said, stated.said)
			}

			if stated.kind != "severity" {
				return
			}
			entries := entriesIn(page)
			if len(entries) == 0 {
				t.Fatal("the page indexes nothing")
			}
			drawn := markIn(entries[0], "severity")
			if drawn == nil {
				t.Fatal("the index entry carries no severity mark")
			}
			if said := strings.TrimSpace(text(drawn)); said != "" {
				t.Errorf("the severity mark in the index says %q, and the index has room only for the icon", said)
			}
		})
	}
}

func TestAnUnstatedSeverityOrAudienceIsNotMarked(t *testing.T) {
	read := oneDecision(document.One)
	read.Decisions[0].Severity = ""
	read.Decisions[0].Audience = ""

	card := cardsIn(rendered(t, read))[0]
	for _, kind := range []string{"severity", "audience"} {
		if mark := markIn(card, kind); mark != nil {
			t.Errorf("a decision stating no %s carries a %s mark announced as %q", kind, kind, attribute(mark, "aria-label"))
		}
	}
}

func progressSaid(page *html.Node) string {
	for _, held := range elements(page, "span") {
		if hasClass(held, "progress-said") {
			return strings.Join(strings.Fields(text(held)), " ")
		}
	}
	return ""
}

func TestTheIndexCountsWhatArrivesChosenAsAnswered(t *testing.T) {
	read := listedLeastSevereFirst()
	read.Decisions[0].Options[0].Recommended = true

	page := rendered(t, read)
	if said := progressSaid(page); said != "1 of 3 answered" {
		t.Errorf("the index says %q, and one of three decisions arrives with its recommendation chosen", said)
	}

	answered := map[string]string{}
	for _, entry := range entriesIn(page) {
		for _, link := range elements(entry, "a") {
			for _, mark := range elements(link, "span") {
				if hasClass(mark, "answered") {
					answered[attribute(link, "data-decision")] = attribute(mark, "data-answered")
				}
			}
		}
	}
	chosen := read.Decisions[0].ID
	for named, held := range answered {
		want := "no"
		if named == chosen {
			want = "yes"
		}
		if held != want {
			t.Errorf("the index marks %s as answered %q, want %q", named, held, want)
		}
	}
	if len(answered) != 3 {
		t.Errorf("the index marks %d decisions, want 3: %v", len(answered), answered)
	}
}

func TestADiffStatesWhatItAddsAndRemoves(t *testing.T) {
	read := oneDecision(document.One)
	read.Decisions[0].Evidence = document.Evidence{{Patch: "--- a/internal/cache/write.go\n+++ b/internal/cache/write.go\n@@ -10,3 +10,4 @@\n func Write(key string, entry Entry) {\n-\tstore[key] = entry\n+\tmu.Lock()\n+\tstore[key] = entry\n }\n"}}

	page := rendered(t, read)
	for _, held := range elements(page, "span") {
		if hasClass(held, "evidence-aside") {
			if said := strings.Join(strings.Fields(text(held)), " "); said != "+2 −1" {
				t.Errorf("the diff says %q of itself, and it adds two lines and removes one", said)
			}
			return
		}
	}
	t.Error("the diff states nothing of what it adds and removes")
}

func TestADiffIsWrittenInTheLanguageStatedOrElseItsFilesOwn(t *testing.T) {
	const patch = "--- a/internal/cache/write.go\n+++ b/internal/cache/write.go\n@@ -1,1 +1,1 @@\n-a\n+b\n"
	for _, stated := range []struct {
		language string
		patch    string
		written  string
	}{
		{"", patch, "go"},
		{"python", patch, "python"},
		{"", "--- a/Makefile\n+++ b/Makefile\n@@ -1,1 +1,1 @@\n-a\n+b\n", ""},
		{"", "--- a/web.v2/build.mjs\n+++ /dev/null\n@@ -1,1 +0,0 @@\n-a\n", "mjs"},
	} {
		read := oneDecision(document.One)
		read.Decisions[0].Evidence = document.Evidence{{Language: stated.language, Patch: stated.patch}}

		for _, held := range elements(rendered(t, read), "div") {
			if hasClass(held, "diff") {
				if written := attribute(held, "data-language"); written != stated.written {
					t.Errorf("a diff of %q stating %q is written in %q, want %q", stated.patch, stated.language, written, stated.written)
				}
			}
		}
	}
}

func TestEvidenceIsDrawnBlockAfterBlockInTheOrderWrittenEachUnderItsOwnHeading(t *testing.T) {
	read := oneDecision(document.One)
	read.Decisions[0].Evidence = document.Evidence{
		{Caption: "as a prop", Language: "go", Code: "first()\n"},
		{Language: "python", Code: "second()\n"},
		{Caption: "the lock", Patch: "--- a/internal/cache/write.go\n+++ b/internal/cache/write.go\n@@ -1,1 +1,1 @@\n-a\n+b\n"},
		{Caption: "nothing to show", Patch: "not a patch at all"},
	}

	var headings []string
	for _, held := range elements(rendered(t, read), "div") {
		if hasClass(held, "evidence") && !hasClass(held, "illustrated") {
			for _, heading := range elements(held, "span") {
				if hasClass(heading, "evidence-label") {
					headings = append(headings, strings.Join(strings.Fields(text(heading)), " "))
				}
			}
		}
	}

	want := []string{"as a prop go", "python", "the lock internal/cache/write.go"}
	if diff := cmp.Diff(want, headings); diff != "" {
		t.Errorf("the headings of the blocks differ from the order and captions written (-want +got):\n%s", diff)
	}
}
