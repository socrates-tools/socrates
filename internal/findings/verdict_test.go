package findings_test

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/internal/findings"
	"github.com/socrates-tools/socrates/internal/runstate"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

func TestNoFindingsAtAll(t *testing.T) {
	document, err := findings.Parse([]byte(`{"findings":[],"risk_level":"low","risk_rationale":"nothing much changed"}`))
	if err != nil {
		t.Fatalf("reading the document: %v", err)
	}

	if verdict := document.Verdict(); verdict != runstate.Clean {
		t.Errorf("an empty findings array is %q, want %q", verdict, runstate.Clean)
	}
	if document.NeedsAPerson() {
		t.Error("an empty findings array needs a person")
	}
}

func TestOnlyNoOpsAndInformationalFindings(t *testing.T) {
	document, err := findings.Parse([]byte(`{"findings":[
		{"id":"f1","severity":"info","file":"a.go","line":1,"description":"a name could be better","action":"no-op"},
		{"id":"f2","severity":"info","file":"b.go","line":2,"description":"a doc could say more","action":"ask-user"}
	],"risk_level":"low","risk_rationale":"r"}`))
	if err != nil {
		t.Fatalf("reading the document: %v", err)
	}

	if verdict := document.Verdict(); verdict != runstate.Clean {
		t.Errorf("only informational findings is %q, want %q", verdict, runstate.Clean)
	}
	if len(document.Findings) != 2 {
		t.Errorf("a clean verdict kept %d findings, want both of them", len(document.Findings))
	}
	if diff := cmp.Diff([]string{"f2"}, identifiers(findings.Classify(document).NeedingAPerson)); diff != "" {
		t.Errorf("the classification of a clean run differs (-want +got):\n%s", diff)
	}
}

func TestAnythingBlocking(t *testing.T) {
	for _, first := range everySeverity {
		for _, action := range everyAction {
			single := document(finding("f1", first, action))
			if got, want := single.Verdict(), verdictFor(first); got != want {
				t.Fatalf("severity %q actioned %q is %q, want %q", first, action, got, want)
			}

			for _, second := range everySeverity {
				for _, otherAction := range everyAction {
					pair := document(finding("f1", first, action), finding("f2", second, otherAction))
					want := runstate.Clean
					if first.Blocks() || second.Blocks() {
						want = runstate.Escalated
					}
					if got := pair.Verdict(); got != want {
						t.Fatalf("severities %q and %q are %q, want %q", first, second, got, want)
					}
				}
			}
		}
	}
}

func TestTheRuleReachesNothingOutward(t *testing.T) {
	forbidden := []string{
		"embed",
		"io/fs",
		"io/ioutil",
		"log",
		"math/rand",
		"os",
		"path/filepath",
		"syscall",
		"time",
	}

	for _, imported := range source.Imports(t, source.Root(t), filepath.Join("internal", "findings")) {
		reaching := slices.Contains(forbidden, imported) ||
			strings.HasPrefix(imported, "os/") ||
			strings.HasPrefix(imported, "net") ||
			strings.HasPrefix(imported, "golang.org/x/sys")
		if reaching {
			t.Errorf("the package holding the verdict rule imports %s", imported)
		}
	}
}

func TestAnInformationalFindingDoesNotBlock(t *testing.T) {
	document := document(finding("f1", findings.Info, findings.AskUser), finding("f2", findings.Info, findings.NoOp))

	if verdict := document.Verdict(); verdict != runstate.Clean {
		t.Errorf("informational findings are %q, want %q", verdict, runstate.Clean)
	}
}

func TestAWarningBlocks(t *testing.T) {
	document := document(finding("f1", findings.Warning, findings.NoOp))

	if verdict := document.Verdict(); verdict == runstate.Clean {
		t.Errorf("a warning is %q", verdict)
	}
	if verdict := document.Verdict(); verdict != runstate.Escalated {
		t.Errorf("a warning is %q, want %q", verdict, runstate.Escalated)
	}
}

func TestActionDoesNotChangeWhetherARunBlocks(t *testing.T) {
	patchable := document(finding("f1", findings.Error, findings.AutoFix))
	asked := document(finding("f1", findings.Error, findings.AskUser))

	if patchable.Verdict() != runstate.Escalated || asked.Verdict() != runstate.Escalated {
		t.Errorf("the same error is %q when patchable and %q when asked", patchable.Verdict(), asked.Verdict())
	}

	if diff := cmp.Diff([]string{"f1"}, identifiers(patchable.Patchable())); diff != "" {
		t.Errorf("the patchable finding differs (-want +got):\n%s", diff)
	}
	if !asked.NeedsAPerson() || asked.Patchable() != nil {
		t.Error("the finding actioned ask-user was offered to a patching agent")
	}
	if patchable.NeedsAPerson() {
		t.Error("the finding actioned auto-fix was sent to the person as well")
	}
}

func TestAnInformationalFindingCanStillBeActionable(t *testing.T) {
	document := document(finding("f1", findings.Info, findings.AutoFix))

	if verdict := document.Verdict(); verdict != runstate.Clean {
		t.Errorf("an informational finding actioned auto-fix is %q, want %q", verdict, runstate.Clean)
	}
	if diff := cmp.Diff([]string{"f1"}, identifiers(document.Patchable())); diff != "" {
		t.Errorf("the patchable findings differ (-want +got):\n%s", diff)
	}
}

func TestAHighRiskLevelWithNoBlockingFindings(t *testing.T) {
	document, err := findings.Parse([]byte(`{"findings":[
		{"id":"f1","severity":"info","file":"a.go","line":1,"description":"d","action":"no-op"}
	],"risk_level":"high","risk_rationale":"the change touches the snapshot path, which everything else stands on"}`))
	if err != nil {
		t.Fatalf("reading the document: %v", err)
	}

	if verdict := document.Verdict(); verdict != runstate.Clean {
		t.Errorf("a high risk level with no blocking findings is %q, want %q", verdict, runstate.Clean)
	}
	if document.RiskLevel != findings.HighRisk {
		t.Errorf("the risk level reads %q, so it was not recorded", document.RiskLevel)
	}
	if document.RiskRationale == "" {
		t.Error("the risk rationale was not recorded")
	}
}

func TestTheRationaleIsOneSentenceAndIsKept(t *testing.T) {
	const rationale = "  The change rewrites the snapshot path, and every other capability stands on it.  "

	document, err := findings.Parse([]byte(`{"findings":[],"risk_level":"medium","risk_rationale":"` + rationale + `"}`))
	if err != nil {
		t.Fatalf("reading the document: %v", err)
	}

	if document.RiskRationale != rationale {
		t.Errorf("the rationale reads %q, want %q verbatim", document.RiskRationale, rationale)
	}
}

func finding(identifier string, severity findings.Severity, action findings.Action) findings.Finding {
	return findings.Finding{
		ID:          identifier,
		Severity:    severity,
		File:        "internal/snapshot/snapshot.go",
		Line:        1,
		Description: "what is wrong",
		Action:      action,
	}
}

func document(held ...findings.Finding) findings.Document {
	return findings.Document{
		Findings:      held,
		RiskLevel:     findings.LowRisk,
		RiskRationale: "nothing much",
	}
}

func verdictFor(severity findings.Severity) runstate.State {
	if severity.Blocks() {
		return runstate.Escalated
	}
	return runstate.Clean
}
