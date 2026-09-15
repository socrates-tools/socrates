package findings_test

import (
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/findings"
	"github.com/socrates-tools/socrates/internal/prompt"
	"github.com/socrates-tools/socrates/internal/skill"
)

var shapingTheDocument = []string{
	"risk_rationale",
	"risk_level",
	"auto-fix",
	"ask-user",
}

func TestTheAuthorNeverShapesTheDocument(t *testing.T) {
	read := skill.Document()
	if strings.TrimSpace(read) == "" {
		t.Fatal("the skill the authoring session reads is empty, so this check reads nothing")
	}

	for _, shaping := range shapingTheDocument {
		if strings.Contains(read, shaping) {
			t.Errorf("what the authoring session reads names %q, which is a field of the findings document: the author has context and does the patching, and a document the author shaped is one whose severities and actions the author chose", shaping)
		}
	}
	if strings.Contains(read, findings.Schema) {
		t.Error("what the authoring session reads carries the findings schema, so it is being asked to produce a document")
	}
	if strings.Contains(read, findings.DeliveryInstruction()) {
		t.Error("what the authoring session reads carries the instruction for delivering a findings document, and that instruction is the reviewer's")
	}
}

func TestOnlyTheReviewerIsAskedForADocument(t *testing.T) {
	// scenario: The author never shapes the document
	asked := prompt.Request{
		Intent:                 "add a parser",
		Repository:             "/repo",
		Base:                   "1111111111111111111111111111111111111111",
		Snapshot:               "2222222222222222222222222222222222222222",
		Changed:                []string{"internal/added.go"},
		Diff:                   "diff --git a/internal/added.go b/internal/added.go\n",
		Ceiling:                1 << 20,
		DocumentEndsTheMessage: true,
	}

	delivered := prompt.Build(asked).Text
	if !strings.Contains(delivered, findings.DeliveryInstruction()) {
		t.Fatal("the reviewer's own prompt does not ask it for the document, so nothing here establishes who was asked")
	}

	for _, named := range []string{"the author", "the coding session", "the session that wrote"} {
		for _, line := range strings.Split(delivered, "\n") {
			if strings.Contains(strings.ToLower(line), named) && strings.Contains(line, "findings document") {
				t.Errorf("the delivered prompt reads %q, which asks somebody other than the reviewer for the document", line)
			}
		}
	}
}
