package agent_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/convention"
	"github.com/socrates-tools/socrates/internal/testing/source"
	"github.com/socrates-tools/socrates/review/agent"
	"github.com/socrates-tools/socrates/review/findings"
)

const aFindingsDocument = `{"findings":[{"id":"f1","severity":"warning","file":"internal/added.go","line":3,"description":"the error is dropped","action":"auto-fix"}],"risk_level":"low","risk_rationale":"one small change"}`

var seamRules = []string{
	"failureTailBytes",
	"BoundedText",
	"beside",
	"theAgentExited",
	"theAgentExitedSayingNothing",
	"theResultReported",
	"theResultReportedSaying",
	"noStructuredOutput",
	"structuredOutputIsNotJSON",
	"schemaNotSatisfied",
	"noDocumentInTheAnswer",
	"noResultEvent",
}

func TestAVendorContributesOnlyWhereTheAccountIs(t *testing.T) {
	restated, err := convention.SeamRulesRestatedOutside(source.Root(t), seamRules)
	if err != nil {
		t.Fatalf("reading what states a rule about failures: %v", err)
	}
	for _, elsewhere := range restated {
		t.Errorf("%s, and the preference order, the bounding and the refusal to decide from the account are stated once in %s rather than once per vendor", elsewhere, convention.SeamFile)
	}
}

func TestEveryVendorReportsAFailureTheSameWay(t *testing.T) {
	said := agent.AccountOf("what the terminal event said", "what the error stream said")
	if !strings.HasPrefix(said, "what the terminal event said") {
		t.Errorf("the account reads %q, and the terminal event is preferred to the error stream", said)
	}
	if !strings.Contains(said, "what the error stream said") {
		t.Errorf("the account reads %q, and both channels are carried where both are there", said)
	}

	if only := agent.AccountOf("", "   "); only != "" {
		t.Errorf("an agent that said nothing on either channel reads as %q, want nothing quoted", only)
	}

	bounded := agent.Bounded(strings.Repeat("x", 64*1024))
	if len(bounded) >= 64*1024 {
		t.Errorf("an account of %d bytes was carried whole, and it is bounded on the ceiling the error stream already has", len(bounded))
	}
	if !strings.HasSuffix(bounded, "...") {
		t.Errorf("a bounded account reads %q at its end, without saying it was cut", bounded[len(bounded)-8:])
	}
}

func TestAFailureThatSaidNothingSaysSo(t *testing.T) {
	silent := agent.ExitedSaying("a-cli", errExited, "")
	if !strings.Contains(silent.Error(), "said nothing") {
		t.Errorf("an agent that exited saying nothing reads %q", silent)
	}

	spoke := agent.ExitedSaying("a-cli", errExited, "the model refused")
	if !strings.Contains(spoke.Error(), "the model refused") {
		t.Errorf("an agent that exited with an account reads %q, without carrying it", spoke)
	}
}

func TestTheEnforcedTierDeliversTheVendorsStructuredOutput(t *testing.T) {
	enforced := agent.Vendor{Name: "constrains", Binary: "constrains", Tier: agent.Enforced}
	request := agent.Request{Schema: json.RawMessage(findings.Schema)}

	document, err := agent.Document(enforced, request, enforced.Binary, json.RawMessage(aFindingsDocument), "prose nobody should read")
	if err != nil {
		t.Fatalf("the enforced tier refused its own structured output: %v", err)
	}
	if string(document) != aFindingsDocument {
		t.Errorf("the enforced tier delivered %s, want the structured output", document)
	}
	if _, err := findings.Parse(document); err != nil {
		t.Errorf("what the enforced tier delivered does not reach the validation: %v", err)
	}
}

func TestAVendorThatWasToConstrainTheSchemaAndDidNotFails(t *testing.T) {
	enforced := agent.Vendor{Name: "constrains", Binary: "constrains", Tier: agent.Enforced}
	request := agent.Request{Schema: json.RawMessage(findings.Schema)}

	answered := "Here is what I found.\n\n```json\n" + aFindingsDocument + "\n```\n"

	if _, err := agent.Document(enforced, request, enforced.Binary, nil, answered); err == nil {
		t.Error("a vendor stated to constrain the schema returned no structured output and the text was read instead, so the tier is not load-bearing")
	}

	if _, err := agent.Document(enforced, request, enforced.Binary, json.RawMessage("not json at all"), answered); err == nil {
		t.Error("a vendor stated to constrain the schema returned structured output that is not JSON and the call did not fail")
	}
}

func TestTheParsedTierDeliversTheDocumentFromTheFinalMessage(t *testing.T) {
	extracting := agent.Vendor{Name: "extracts", Binary: "extracts", Tier: agent.Parsed}
	request := agent.Request{Schema: json.RawMessage(findings.Schema)}

	answered := "I read the change.\n\n```go\nfunc broken() {}\n```\n\nThe findings:\n\n```json\n" + aFindingsDocument + "\n```\n"

	document, err := agent.Document(extracting, request, extracting.Binary, nil, answered)
	if err != nil {
		t.Fatalf("the parsed tier did not find the document in the final message: %v", err)
	}
	parsed, err := findings.Parse(document)
	if err != nil {
		t.Fatalf("what the parsed tier delivered does not reach the same validation: %v", err)
	}
	if len(parsed.Findings) != 1 || parsed.Findings[0].ID != "f1" {
		t.Errorf("the parsed tier delivered %s", document)
	}
}

func TestAnUnparseableDocumentFailsOnEitherTier(t *testing.T) {
	request := agent.Request{
		Schema:  json.RawMessage(findings.Schema),
		Honours: func(document json.RawMessage) error { _, err := findings.Parse(document); return err },
	}

	violating := `{"findings":[{"id":"f1","severity":"catastrophic","file":"a.go","line":1,"description":"x","action":"no-op"}],"risk_level":"low","risk_rationale":"y"}`

	for _, tier := range []agent.Tier{agent.Enforced, agent.Parsed} {
		vendor := agent.Vendor{Name: "either", Binary: "either", Tier: tier}

		document, err := agent.Document(vendor, request, vendor.Binary, json.RawMessage(violating), "```json\n"+violating+"\n```\n")
		if err != nil {
			t.Fatalf("the %s tier refused a well-formed document before validating it: %v", tier, err)
		}
		if err := agent.Honoured(request, document); err == nil {
			t.Errorf("the %s tier delivered a document violating the schema and the call did not fail", tier)
		}

		if _, err := agent.Document(vendor, request, vendor.Binary, nil, "no document here at all"); err == nil {
			t.Errorf("the %s tier delivered a document where the reviewer returned none", tier)
		}
	}
}

func TestAReviewerThatAnsweredWithoutADocumentFailsTheCall(t *testing.T) {
	extracting := agent.Vendor{Name: "extracts", Binary: "extracts", Tier: agent.Parsed}
	request := agent.Request{Schema: json.RawMessage(findings.Schema)}

	_, err := agent.Document(extracting, request, extracting.Binary, nil, "I read the change and it looks fine to me.")
	if err == nil {
		t.Fatal("a reviewer that answered in prose with no document was read as having delivered one")
	}
	if !strings.Contains(err.Error(), extracting.Binary) {
		t.Errorf("the failure reads %q, without naming which reviewer answered", err)
	}
}

func TestASchemaIsHonouredOnceOrTheCallFails(t *testing.T) {
	// scenario: A schema is honoured or the call fails
	asked := 0
	request := agent.Request{
		Schema: json.RawMessage(findings.Schema),
		Honours: func(json.RawMessage) error {
			asked++
			return errNotHonoured
		},
	}

	if err := agent.Honoured(request, json.RawMessage(aFindingsDocument)); err == nil {
		t.Error("a document the schema check refused did not fail the call")
	}
	if asked != 1 {
		t.Errorf("the schema check was asked %d times, want once", asked)
	}
}

type sentinel string

func (s sentinel) Error() string { return string(s) }

const (
	errExited      = sentinel("exit status 1")
	errNotHonoured = sentinel("the document does not satisfy the schema")
)
