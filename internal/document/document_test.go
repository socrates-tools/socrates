package document_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates-review/internal/document"
)

func whole() document.Document {
	return document.Document{
		Title:      "two things to settle",
		Stylesheet: `.panel{border:1px solid #d8d2c4}`,
		Decisions: []document.Decision{{
			ID:       "d-1",
			Question: "does the cache belong here",
			Meaning:  "the cache is read on every call and nothing invalidates it",
			Mode:     document.One,
			Options: []document.Option{
				{Label: "keep it", Meaning: "the staleness is acceptable", Command: `agree d-1 --reason "the staleness is acceptable"`},
				{Label: "drop it", Meaning: "correctness first", Command: `patch d-1`},
			},
			Fragment: document.Fragment{Markup: `<div class="panel">one call, one read</div>`},
		}},
	}
}

func raw(t *testing.T, written document.Document) []byte {
	t.Helper()

	encoded, err := json.Marshal(written)
	if err != nil {
		t.Fatalf("writing a document out: %v", err)
	}
	return encoded
}

func briefly(read document.Document) string {
	return fmt.Sprintf("title %q, %d decisions, %+v, a stylesheet of %d bytes",
		read.Title, len(read.Decisions), read.Decisions, len(read.Stylesheet))
}

func TestADocumentThatDoesNotValidateIsRefusedBeforeServing(t *testing.T) {
	// scenario: A document that asks something is validated as before
	for _, refused := range []struct {
		named  string
		broken func(*document.Document)
		names  []string
	}{
		{
			named:  "a missing field",
			broken: func(d *document.Document) { d.Decisions[0].Question = "" },
			names:  []string{"question", "d-1"},
		},
		{
			named:  "a document with no title",
			broken: func(d *document.Document) { d.Title = "" },
			names:  []string{"title"},
		},
		{
			named:  "an unknown mode",
			broken: func(d *document.Document) { d.Decisions[0].Mode = "pick-a-few" },
			names:  []string{"pick-a-few", "d-1"},
		},
		{
			named:  "an option with no command",
			broken: func(d *document.Document) { d.Decisions[0].Options[1].Command = "" },
			names:  []string{"drop it", "d-1"},
		},
		{
			named:  "a decision with no identifier",
			broken: func(d *document.Document) { d.Decisions[0].ID = "" },
			names:  []string{"decision 1"},
		},
		{
			named:  "a document holding nothing at all",
			broken: func(d *document.Document) { d.Decisions = nil },
			names:  []string{"decision", "source", "preamble"},
		},
		{
			named: "a stylesheet larger than the stated size",
			broken: func(d *document.Document) {
				d.Stylesheet = strings.Repeat(".p{color:#222}", document.StylesheetCeiling)
			},
			names: []string{"stylesheet", strconv.Itoa(document.StylesheetCeiling)},
		},
	} {
		t.Run(refused.named, func(t *testing.T) {
			written := whole()
			refused.broken(&written)

			read, err := document.Read(raw(t, written))
			if err == nil {
				t.Fatalf("%s validated, and a document that does not validate is refused before serving: %s", refused.named, briefly(read))
			}
			for _, named := range refused.names {
				if !strings.Contains(err.Error(), named) {
					t.Errorf("refusing %s does not name %q: %v", refused.named, named, err)
				}
			}
		})
	}
}

func nothingToDecide() document.Document {
	return document.Document{
		Title:    "nothing to settle",
		Preamble: "what follows is what was seen",
		Source: []document.Entry{
			{ID: "c-1", Name: "the first thing", Severity: "warning", Audience: "the author patches", Location: "internal/cache/read.go:41"},
			{ID: "c-2", Name: "the second thing", Severity: "info", Audience: "the author patches", Location: "internal/cache/write.go:12"},
		},
	}
}

func TestADocumentHoldingNoDecisionIsRead(t *testing.T) {
	read, err := document.Read(raw(t, nothingToDecide()))
	if err != nil {
		t.Fatalf("reading a document holding no decision: %v", err)
	}
	if diff := cmp.Diff(nothingToDecide(), read); diff != "" {
		t.Errorf("what reading a document holding no decision gives back differs (-want +got):\n%s", diff)
	}
}

func TestADocumentHoldingOnlyAPreambleIsRead(t *testing.T) {
	written := nothingToDecide()
	written.Source = nil

	if _, err := document.Read(raw(t, written)); err != nil {
		t.Fatalf("reading a document holding only a title and a preamble: %v", err)
	}
}

func TestADocumentHoldingOnlyASourceIsRead(t *testing.T) {
	written := nothingToDecide()
	written.Preamble = ""

	if _, err := document.Read(raw(t, written)); err != nil {
		t.Fatalf("reading a document holding only a title and a source: %v", err)
	}
}

func statingSomething() document.Document {
	written := whole()
	written.Decisions[0].Options = nil
	written.Decisions[0].Mode = ""
	return written
}

func TestADecisionMayStateSomethingRatherThanAsk(t *testing.T) {
	read, err := document.Read(raw(t, statingSomething()))
	if err != nil {
		t.Fatalf("reading a document holding a decision that asks nothing: %v", err)
	}
	if len(read.Decisions) != 1 {
		t.Fatalf("a document holding one decision that asks nothing came back with %d", len(read.Decisions))
	}
	if len(read.Decisions[0].Options) != 0 {
		t.Errorf("a decision that asks nothing came back offering %d options", len(read.Decisions[0].Options))
	}
	if read.Decisions[0].Question != statingSomething().Decisions[0].Question {
		t.Errorf("what the decision states came back as %q", read.Decisions[0].Question)
	}
}

func TestAModeOnADecisionWithNoOptionsIsRefusedBeforeServing(t *testing.T) {
	for _, stated := range document.Modes() {
		t.Run(string(stated), func(t *testing.T) {
			written := statingSomething()
			written.Decisions[0].Mode = stated

			read, err := document.Read(raw(t, written))
			if err == nil {
				t.Fatalf("a decision offering nothing and stating the mode %q validated: %s", stated, briefly(read))
			}
			for _, named := range []string{"d-1", string(stated)} {
				if !strings.Contains(err.Error(), named) {
					t.Errorf("refusing a mode on a decision with no options does not name %q: %v", named, err)
				}
			}
		})
	}
}

func TestADecisionWithOptionsStillRequiresItsMode(t *testing.T) {
	written := whole()
	written.Decisions[0].Mode = ""

	read, err := document.Read(raw(t, written))
	if err == nil {
		t.Fatalf("a decision offering options and stating no mode validated: %s", briefly(read))
	}
	if !strings.Contains(err.Error(), "d-1") {
		t.Errorf("refusing a decision that states no mode does not name it: %v", err)
	}
}

func TestAnUnknownEvidenceLanguageIsRefusedBeforeServing(t *testing.T) {
	written := whole()
	written.Decisions[0].Evidence = document.Evidence{Language: "prolog", Code: "read(X) :- lock(X).\n"}

	read, err := document.Read(raw(t, written))
	if err == nil {
		t.Fatalf("a document naming a language the page cannot highlight validated: %s", briefly(read))
	}
	for _, named := range []string{"d-1", "prolog"} {
		if !strings.Contains(err.Error(), named) {
			t.Errorf("refusing an unknown evidence language does not name %q: %v", named, err)
		}
	}
	for _, known := range []string{"csharp", "diff", "go"} {
		if !strings.Contains(err.Error(), known) {
			t.Errorf("refusing an unknown evidence language does not name %q, and the way forward is in the message: %v", known, err)
		}
	}
}

func TestEvidenceWithNoLanguageIsServed(t *testing.T) {
	written := whole()
	written.Decisions[0].Evidence = document.Evidence{Code: "entry, held := cache.Read(key)\n"}

	read, err := document.Read(raw(t, written))
	if err != nil {
		t.Fatalf("reading a document carrying evidence that names no language: %v", err)
	}
	if read.Decisions[0].Evidence.Language != "" {
		t.Errorf("evidence naming no language came back naming %q", read.Decisions[0].Evidence.Language)
	}
}

func TestEveryLanguageTheDocumentKnowsIsServed(t *testing.T) {
	// scenario: The languages a caller actually passes are highlighted
	known := document.Languages()
	if len(known) == 0 {
		t.Fatal("the document contract knows no language at all, so this test proves nothing")
	}

	for _, named := range known {
		written := whole()
		written.Decisions[0].Evidence = document.Evidence{Language: named, Code: "one line\n"}

		if _, err := document.Read(raw(t, written)); err != nil {
			t.Errorf("a document naming the language %q is refused, and the page highlights it: %v", named, err)
		}
	}
}

func TestADecisionNamesTheOptionTheCallerWouldTake(t *testing.T) {
	written := whole()
	written.Decisions[0].Options[1].Recommended = true

	read, err := document.Read(raw(t, written))
	if err != nil {
		t.Fatalf("reading a document marking the option the caller would take: %v", err)
	}
	if diff := cmp.Diff(written.Decisions[0].Options, read.Decisions[0].Options); diff != "" {
		t.Errorf("what the marked options come back as differs (-want +got):\n%s", diff)
	}
	if !read.Decisions[0].Options[1].Recommended {
		t.Error("the option the document marked comes back unmarked")
	}
	if read.Decisions[0].Options[0].Recommended {
		t.Error("an option the document did not mark comes back marked")
	}
}

func TestADocumentMarkingTwoOptionsIsRefusedBeforeServing(t *testing.T) {
	written := whole()
	written.Decisions[0].Options[0].Recommended = true
	written.Decisions[0].Options[1].Recommended = true

	read, err := document.Read(raw(t, written))
	if err == nil {
		t.Fatalf("a document marking two options validated: %+v", read)
	}
	if !strings.Contains(err.Error(), "d-1") {
		t.Errorf("refusing a document marking two options does not name the decision: %v", err)
	}
}

func TestADocumentThatMarksNothingIsUnaffected(t *testing.T) {
	written := whole()

	read, err := document.Read(raw(t, written))
	if err != nil {
		t.Fatalf("reading a document that marks nothing: %v", err)
	}
	if diff := cmp.Diff(written, read); diff != "" {
		t.Errorf("a document that marks nothing reads back differently (-want +got):\n%s", diff)
	}
	for _, option := range read.Decisions[0].Options {
		if option.Recommended {
			t.Errorf("the option %q comes back marked, and the document marked nothing", option.Label)
		}
	}
}

func TestAnIllustrationWrittenAsABareStringIsRefusedByName(t *testing.T) {
	// scenario: An illustration the stylesheet paints survives neutralising
	written := map[string]any{
		"title": "one thing to settle",
		"decisions": []any{map[string]any{
			"id":       "d-1",
			"question": "does the cache belong here",
			"mode":     string(document.One),
			"options":  []any{map[string]any{"label": "keep it", "command": "agree d-1"}},
			"fragment": `<div class="panel">one call, one read</div>`,
		}},
	}

	raw, err := json.Marshal(written)
	if err != nil {
		t.Fatalf("writing a document out: %v", err)
	}

	read, err := document.Read(raw)
	if err == nil {
		t.Fatalf("an illustration written as a bare string was read, and the field says what it is made of: %+v", read)
	}
	if !strings.Contains(err.Error(), "fragment") {
		t.Errorf("refusing an illustration written as a bare string does not name the field: %v", err)
	}
}

func TestAWholeDocumentIsRead(t *testing.T) {
	read, err := document.Read(raw(t, whole()))
	if err != nil {
		t.Fatalf("reading a whole document: %v", err)
	}
	if diff := cmp.Diff(whole(), read); diff != "" {
		t.Errorf("what reading a whole document gives back differs (-want +got):\n%s", diff)
	}
}

func TestACallerCannotRenameAnIdentifier(t *testing.T) {
	written := whole()
	written.Decisions[0].Severity = "info"
	written.Decisions[0].Audience = "nobody in particular"
	written.Decisions[0].Location = "somewhere friendlier.go"
	written.Decisions[0].Name = "a heading of its own"
	written.Source = []document.Entry{{
		ID:       "d-1",
		Name:     "the first thing",
		Severity: "error",
		Audience: "a person",
		Location: "internal/cache/read.go:41",
	}}

	read, err := document.Read(raw(t, written))
	if err != nil {
		t.Fatalf("reading a document citing a source: %v", err)
	}

	said := read.Decisions[0]
	if diff := cmp.Diff(written.Source[0], document.Entry{
		ID:       said.ID,
		Name:     said.Name,
		Severity: said.Severity,
		Audience: said.Audience,
		Location: said.Location,
	}); diff != "" {
		t.Errorf("what the page would show is not what the source said (-source +shown):\n%s", diff)
	}
}

func TestADecisionNamingNothingInTheSourceIsRefused(t *testing.T) {
	written := whole()
	written.Source = []document.Entry{{ID: "d-2", Severity: "error", Audience: "a person", Location: "internal/cache/read.go:41"}}

	if _, err := document.Read(raw(t, written)); err == nil {
		t.Fatal("a decision naming nothing in the source was read, and it is refused")
	} else if !strings.Contains(err.Error(), "d-1") {
		t.Errorf("refusing a decision naming nothing in the source does not name the identifier: %v", err)
	}
}

func TestTheSourceShapeACallerWritesIsWhatIsRead(t *testing.T) {
	// scenario: A caller cannot rename an identifier
	sample, err := os.ReadFile(filepath.Join("testdata", "source.json"))
	if err != nil {
		t.Fatalf("reading the checked-in sample of the source shape: %v", err)
	}

	read, err := document.Read(sample)
	if err != nil {
		t.Fatalf("reading the checked-in sample of the source shape: %v", err)
	}

	if len(read.Source) == 0 {
		t.Fatal("the checked-in sample cites no source, so it holds nothing for a document to be copied from")
	}
	for _, decision := range read.Decisions {
		var cited document.Entry
		for _, entry := range read.Source {
			if entry.ID == decision.ID {
				cited = entry
			}
		}
		if decision.Name != cited.Name || decision.Severity != cited.Severity || decision.Audience != cited.Audience || decision.Location != cited.Location {
			t.Errorf("%s does not carry what the sample's source entry stated: %+v", decision.ID, decision)
		}
	}
}
