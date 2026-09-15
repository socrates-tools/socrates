package findings_test

import (
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/socrates-review/findings"
)

const aRealDocument = `{"findings":[{"id":"f1","severity":"warning","file":"internal/added.go","line":3,"description":"the error is dropped","action":"auto-fix"}],"risk_level":"low","risk_rationale":"one small change"}`

func TestAReviewerThatAnsweredWithoutADocument(t *testing.T) {
	for _, answered := range []struct {
		named   string
		message string
	}{
		{named: "nothing at all", message: ""},
		{named: "prose and no block", message: "I read the change and found nothing worth reporting."},
		{named: "a block that never closes", message: "here it is:\n```json\n" + aRealDocument},
	} {
		t.Run(answered.named, func(t *testing.T) {
			if _, err := findings.Extract(answered.message); err == nil {
				t.Errorf("a final message carrying %s was read as a findings document", answered.named)
			}
		})
	}
}

func TestALastBlockThatIsNotJSONIsNotADocument(t *testing.T) {
	message := "The findings:\n\n```json\n" + aRealDocument + "\n```\n\nAnd for reference:\n\n```go\nfunc broken() {}\n```\n"

	if _, err := findings.Extract(message); err == nil {
		t.Error("a final message ending with a block that is not JSON was read as a findings document")
	}
}

func TestReasoningBeforeTheAnswerIsNotMistakenForTheAnswer(t *testing.T) {
	message := strings.Join([]string{
		"The change drops an error here:",
		"",
		"```go",
		"func broken() { _ = doSomething() }",
		"```",
		"",
		"An earlier draft of the document looked like this, and it is not the answer:",
		"",
		"```json",
		`{"findings":[],"risk_level":"high","risk_rationale":"a draft nobody meant"}`,
		"```",
		"",
		"The findings:",
		"",
		"```json",
		aRealDocument,
		"```",
		"",
	}, "\n")

	extracted, err := findings.Extract(message)
	if err != nil {
		t.Fatalf("extracting the document from a final message that reasoned first: %v", err)
	}

	document, err := findings.Parse(extracted)
	if err != nil {
		t.Fatalf("parsing what was extracted: %v", err)
	}
	if len(document.Findings) != 1 || document.Findings[0].ID != "f1" {
		t.Errorf("what was extracted is %s, which is not the document at the end", extracted)
	}
	if document.RiskRationale == "a draft nobody meant" {
		t.Error("the draft earlier in the message was extracted rather than the answer at the end")
	}
}

func TestADocumentViolatingTheSchemaIsExtractedAndThenRefused(t *testing.T) {
	// scenario: An unparseable document fails on either tier
	message := "```json\n" + `{"findings":[{"id":"f1","severity":"catastrophic","file":"a.go","line":1,"description":"x","action":"no-op"}],"risk_level":"low","risk_rationale":"y"}` + "\n```\n"

	extracted, err := findings.Extract(message)
	if err != nil {
		t.Fatalf("extraction refused a well-formed block: %v", err)
	}
	if _, err := findings.Parse(extracted); err == nil {
		t.Error("a document naming a severity the schema does not have was accepted, so extraction reaches no validation")
	}
}

func TestADocumentIsExtractedFromAnUnlabelledBlock(t *testing.T) {
	message := "Here is what I found.\n\n```\n" + aRealDocument + "\n```\n"

	extracted, err := findings.Extract(message)
	if err != nil {
		t.Fatalf("extracting from a block carrying no language: %v", err)
	}
	if _, err := findings.Parse(extracted); err != nil {
		t.Errorf("what was extracted from a block carrying no language does not parse: %v", err)
	}
}
