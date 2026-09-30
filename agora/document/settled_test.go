package document_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/agora/document"
)

func withWhatWasSettled() document.Document {
	written := whole()
	written.Settled = []document.Settled{
		{ID: "q1", Name: "the surface", Location: "agora/surface", Question: "what is the page called", Chosen: []string{"the surface"}},
		{ID: "q2", Question: "should a design carry a cap", Chosen: []string{"no cap", "no ceiling either"}, Words: "nothing here spawns a process for a cap to protect"},
		{ID: "q3", Question: "who writes the glossary entry"},
	}
	return written
}

func TestWhatWasSettledIsReadAsTheCallerWroteIt(t *testing.T) {
	read, err := document.Read(raw(t, withWhatWasSettled()))
	if err != nil {
		t.Fatalf("reading a document carrying what was settled: %v", err)
	}
	if diff := cmp.Diff(withWhatWasSettled().Settled, read.Settled); diff != "" {
		t.Errorf("what was settled is read differently from how the caller wrote it (-want +got):\n%s", diff)
	}
}

func TestADocumentHoldingOnlyWhatWasSettledIsRead(t *testing.T) {
	written := withWhatWasSettled()
	written.Decisions = nil

	if _, err := document.Read(raw(t, written)); err != nil {
		t.Fatalf("reading a document holding a title and what was settled: %v", err)
	}
}

func TestASettledAnswerThatCannotBeShownIsRefusedBeforeServing(t *testing.T) {
	for _, refused := range []struct {
		named  string
		broken func(*document.Document)
		names  []string
	}{
		{
			named:  "a settled answer with no identifier",
			broken: func(d *document.Document) { d.Settled[1].ID = "" },
			names:  []string{"settled answer 2"},
		},
		{
			named:  "two settled answers under one identifier",
			broken: func(d *document.Document) { d.Settled[1].ID = "q1" },
			names:  []string{"q1"},
		},
		{
			named:  "a question both settled and still open",
			broken: func(d *document.Document) { d.Settled[0].ID = "d-1" },
			names:  []string{"d-1", "settled"},
		},
		{
			named:  "a settled answer stating no question",
			broken: func(d *document.Document) { d.Settled[0].Question = " " },
			names:  []string{"q1", "question"},
		},
		{
			named:  "a settled answer choosing a blank label",
			broken: func(d *document.Document) { d.Settled[1].Chosen[1] = "" },
			names:  []string{"q2", "label"},
		},
	} {
		t.Run(refused.named, func(t *testing.T) {
			written := withWhatWasSettled()
			refused.broken(&written)

			read, err := document.Read(raw(t, written))
			if err == nil {
				t.Fatalf("%s validated, and a settled answer the page cannot show is refused before serving: %+v", refused.named, read.Settled)
			}
			for _, named := range refused.names {
				if !strings.Contains(err.Error(), named) {
					t.Errorf("refusing %s does not name %q: %v", refused.named, named, err)
				}
			}
		})
	}
}
