package findings_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-devtool/socrates/internal/findings"
)

func TestTheReviewerIsGivenTheSchema(t *testing.T) {
	var schema map[string]any
	if err := json.Unmarshal([]byte(findings.Schema), &schema); err != nil {
		t.Fatalf("the schema is not JSON: %v", err)
	}

	properties, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Fatal("the schema states no properties")
	}
	for _, required := range []string{"findings", "risk_level", "risk_rationale"} {
		if _, present := properties[required]; !present {
			t.Errorf("the schema states no %q", required)
		}
	}

	finding, ok := properties["findings"].(map[string]any)["items"].(map[string]any)["properties"].(map[string]any)
	if !ok {
		t.Fatal("the schema does not describe a finding")
	}
	for _, required := range []string{"id", "severity", "file", "line", "description", "action"} {
		if _, present := finding[required]; !present {
			t.Errorf("the schema states no %q on a finding", required)
		}
	}
	for _, absent := range []string{"review_scope", "scope", "phase"} {
		if _, present := finding[absent]; present {
			t.Errorf("the schema classifies a finding by %q, and socrates has no delivery phase", absent)
		}
	}

	if diff := cmp.Diff([]any{"error", "warning", "info"}, finding["severity"].(map[string]any)["enum"]); diff != "" {
		t.Errorf("the severities differ (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]any{"no-op", "auto-fix", "ask-user"}, finding["action"].(map[string]any)["enum"]); diff != "" {
		t.Errorf("the actions differ (-want +got):\n%s", diff)
	}

	t.Run("a response that does not satisfy it fails the round", func(t *testing.T) {
		refused := map[string]string{
			"an unknown severity":      `{"findings":[{"id":"f1","severity":"critical","file":"a.go","line":1,"description":"d","action":"ask-user"}],"risk_level":"low","risk_rationale":"r"}`,
			"an unknown risk level":    `{"findings":[],"risk_level":"catastrophic","risk_rationale":"r"}`,
			"a field nobody asked for": `{"findings":[{"id":"f1","severity":"error","file":"a.go","line":1,"description":"d","action":"ask-user","review_scope":"pushed"}],"risk_level":"low","risk_rationale":"r"}`,
			"no findings array":        `{"risk_level":"low","risk_rationale":"r"}`,
			"no risk level":            `{"findings":[],"risk_rationale":"r"}`,
			"prose":                    `the change looks fine to me`,
			"two documents":            `{"findings":[],"risk_level":"low","risk_rationale":"r"} {"findings":[],"risk_level":"low","risk_rationale":"r"}`,
		}

		for name, document := range refused {
			t.Run(name, func(t *testing.T) {
				if parsed, err := findings.Parse([]byte(document)); err == nil {
					t.Errorf("the document was accepted as %+v", parsed)
				}
			})
		}
	})
}

func TestAFindingNamesAFileAndALine(t *testing.T) {
	parsed, err := findings.Parse([]byte(`{
		"findings": [{"id":"f1","severity":"error","file":"internal/snapshot/carve.go","line":42,"description":"d","action":"ask-user"}],
		"risk_level": "low",
		"risk_rationale": "r"
	}`))
	if err != nil {
		t.Fatalf("reading the document: %v", err)
	}

	finding := parsed.Findings[0]
	if finding.File != "internal/snapshot/carve.go" {
		t.Errorf("the finding names the file %q", finding.File)
	}
	if finding.Line != 42 {
		t.Errorf("the finding names line %d", finding.Line)
	}

	refused := map[string]string{
		"an absolute path":      `"/etc/passwd"`,
		"a path above the root": `"../elsewhere/a.go"`,
		"an unclean path":       `"./a.go"`,
		"no file at all":        `""`,
	}
	for name, file := range refused {
		t.Run(name, func(t *testing.T) {
			document := `{"findings":[{"id":"f1","severity":"error","file":` + file + `,"line":1,"description":"d","action":"ask-user"}],"risk_level":"low","risk_rationale":"r"}`
			if _, err := findings.Parse([]byte(document)); err == nil {
				t.Errorf("a finding naming the file %s was accepted", file)
			}
		})
	}

	for _, line := range []string{"0", "-1"} {
		t.Run("line "+line, func(t *testing.T) {
			document := `{"findings":[{"id":"f1","severity":"error","file":"a.go","line":` + line + `,"description":"d","action":"ask-user"}],"risk_level":"low","risk_rationale":"r"}`
			if _, err := findings.Parse([]byte(document)); err == nil {
				t.Errorf("a finding naming line %s was accepted, and lines are counted from one", line)
			}
		})
	}
}

func TestIdentifiersAreUniqueWithinARun(t *testing.T) {
	document := `{"findings":[
		{"id":"f1","severity":"error","file":"a.go","line":1,"description":"first","action":"ask-user"},
		{"id":"f1","severity":"warning","file":"b.go","line":2,"description":"second","action":"ask-user"}
	],"risk_level":"low","risk_rationale":"r"}`

	_, err := findings.Parse([]byte(document))
	if err == nil {
		t.Fatal("a document carrying two findings identified f1 was accepted")
	}
	if !strings.Contains(err.Error(), "f1") {
		t.Errorf("the failure reads %q, without naming the duplicate", err)
	}

	unnamed := `{"findings":[{"id":"","severity":"error","file":"a.go","line":1,"description":"d","action":"ask-user"}],"risk_level":"low","risk_rationale":"r"}`
	if _, err := findings.Parse([]byte(unnamed)); err == nil {
		t.Error("a finding with no identifier was accepted, so nothing could be attached to it")
	}
}
