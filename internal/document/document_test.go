package document_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rui-ktei/socrates/internal/document"
)

func whole() document.Document {
	return document.Document{
		Title: "two things to settle",
		Decisions: []document.Decision{{
			ID:       "d-1",
			Question: "does the cache belong here",
			Meaning:  "the cache is read on every call and nothing invalidates it",
			Mode:     document.One,
			Options: []document.Option{
				{Label: "keep it", Meaning: "the staleness is acceptable", Command: `agree d-1 --reason "the staleness is acceptable"`},
				{Label: "drop it", Meaning: "correctness first", Command: `patch d-1`},
			},
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

func TestADocumentThatDoesNotValidateIsRefusedBeforeServing(t *testing.T) {
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
			named:  "a decision with no options and no field of its own to answer in",
			broken: func(d *document.Document) { d.Decisions[0].Options = nil },
			names:  []string{"d-1"},
		},
		{
			named:  "a decision with no identifier",
			broken: func(d *document.Document) { d.Decisions[0].ID = "" },
			names:  []string{"decision 1"},
		},
		{
			named:  "a document holding no decision",
			broken: func(d *document.Document) { d.Decisions = nil },
			names:  []string{"decision"},
		},
	} {
		t.Run(refused.named, func(t *testing.T) {
			written := whole()
			refused.broken(&written)

			read, err := document.Read(raw(t, written))
			if err == nil {
				t.Fatalf("%s validated, and a document that does not validate is refused before serving: %+v", refused.named, read)
			}
			for _, named := range refused.names {
				if !strings.Contains(err.Error(), named) {
					t.Errorf("refusing %s does not name %q: %v", refused.named, named, err)
				}
			}
		})
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
	written.Source = []document.Entry{{
		ID:       "d-1",
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
		if decision.Severity != cited.Severity || decision.Audience != cited.Audience || decision.Location != cited.Location {
			t.Errorf("%s does not carry what the sample's source entry stated: %+v", decision.ID, decision)
		}
	}
}
