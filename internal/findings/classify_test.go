package findings_test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rui-ktei/socrates/internal/findings"
	"github.com/rui-ktei/socrates/internal/testing/source"
)

var everySeverity = []findings.Severity{findings.Error, findings.Warning, findings.Info, "", "critical", "INFO"}

var everyAction = []findings.Action{findings.NoOp, findings.AutoFix, findings.AskUser, "", "reformat", "AUTO-FIX", " auto-fix ", "auto_fix"}

func TestAMissingActionDefaultsToNeedingAPerson(t *testing.T) {
	finding := onlyFinding(t, `{"id":"f1","severity":"error","file":"a.go","line":1,"description":"d"}`)

	if finding.Action != "" {
		t.Errorf("the finding states the action %q, so this test asserts nothing", finding.Action)
	}
	if effective := finding.EffectiveAction(); effective != findings.AskUser {
		t.Errorf("a finding with no action is %q, want %q", effective, findings.AskUser)
	}
}

func TestAnEmptyActionDefaultsToNeedingAPerson(t *testing.T) {
	finding := onlyFinding(t, `{"id":"f1","severity":"error","file":"a.go","line":1,"description":"d","action":""}`)

	if effective := finding.EffectiveAction(); effective != findings.AskUser {
		t.Errorf("a finding with an empty action is %q, want %q", effective, findings.AskUser)
	}
}

func TestAnUnrecognisedActionDefaultsToNeedingAPerson(t *testing.T) {
	finding := onlyFinding(t, `{"id":"f1","severity":"error","file":"a.go","line":1,"description":"d","action":"reformat"}`)

	if effective := finding.EffectiveAction(); effective != findings.AskUser {
		t.Errorf("a finding actioned reformat is %q, want %q", effective, findings.AskUser)
	}
	if finding.Action != "reformat" {
		t.Errorf("the unrecognised action reads %q, so it is not visible as what the reviewer said", finding.Action)
	}
}

func TestEveryFindingIsInExactlyOneBucket(t *testing.T) {
	var document findings.Document
	document.RiskLevel = findings.LowRisk

	wanted := map[string]findings.Action{}
	for _, severity := range everySeverity {
		for _, action := range everyAction {
			identifier := fmt.Sprintf("f%d", len(document.Findings)+1)
			document.Findings = append(document.Findings, findings.Finding{
				ID:          identifier,
				Severity:    severity,
				File:        "a.go",
				Line:        1,
				Description: "d",
				Action:      action,
			})

			switch action {
			case findings.AutoFix, findings.NoOp:
				wanted[identifier] = action
			default:
				wanted[identifier] = findings.AskUser
			}
		}
	}

	classified := findings.Classify(document)
	if classified.Total() != len(document.Findings) {
		t.Errorf("%d findings were classified into %d buckets", len(document.Findings), classified.Total())
	}

	placed := map[string]findings.Action{}
	for bucket, held := range map[findings.Action][]findings.Finding{
		findings.AskUser: classified.NeedingAPerson,
		findings.AutoFix: classified.Patchable,
		findings.NoOp:    classified.NoOps,
	} {
		for _, finding := range held {
			if already, twice := placed[finding.ID]; twice {
				t.Errorf("%s is in both the %s and the %s buckets", finding.ID, already, bucket)
			}
			placed[finding.ID] = bucket
		}
	}

	if diff := cmp.Diff(wanted, placed); diff != "" {
		t.Errorf("the classification differs (-want +got):\n%s", diff)
	}

	if empty := findings.Classify(findings.Document{RiskLevel: findings.LowRisk}); empty.Total() != 0 {
		t.Errorf("an empty document classified into %d findings", empty.Total())
	}
}

func TestADocumentMixingEveryCase(t *testing.T) {
	document, err := findings.Parse([]byte(`{"findings":[
		{"id":"patchable","severity":"info","file":"a.go","line":1,"description":"d","action":"auto-fix"},
		{"id":"asked","severity":"error","file":"a.go","line":2,"description":"d","action":"ask-user"},
		{"id":"nothing","severity":"info","file":"a.go","line":3,"description":"d","action":"no-op"},
		{"id":"unknown","severity":"warning","file":"a.go","line":4,"description":"d","action":"reformat"},
		{"id":"empty","severity":"warning","file":"a.go","line":5,"description":"d","action":""},
		{"id":"absent","severity":"error","file":"a.go","line":6,"description":"d"}
	],"risk_level":"medium","risk_rationale":"r"}`))
	if err != nil {
		t.Fatalf("reading the document: %v", err)
	}

	classified := findings.Classify(document)

	if diff := cmp.Diff([]string{"asked", "unknown", "empty", "absent"}, identifiers(classified.NeedingAPerson)); diff != "" {
		t.Errorf("the findings needing a person differ (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]string{"patchable"}, identifiers(classified.Patchable)); diff != "" {
		t.Errorf("the patchable findings differ (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]string{"nothing"}, identifiers(classified.NoOps)); diff != "" {
		t.Errorf("the no-ops differ (-want +got):\n%s", diff)
	}
	if classified.Total() != len(document.Findings) {
		t.Errorf("%d findings became %d, so something was lost", len(document.Findings), classified.Total())
	}

	if !document.NeedsAPerson() {
		t.Error("a document holding four findings for the person says nobody is needed")
	}
	if diff := cmp.Diff([]string{"patchable"}, identifiers(document.Patchable())); diff != "" {
		t.Errorf("the patchable findings the consumer sees differ (-want +got):\n%s", diff)
	}
}

func TestTheTwoConsumersAgreeByConstruction(t *testing.T) {
	root := source.Root(t)
	owner := filepath.Join(root, "internal", "findings", "findings.go")

	const decider = "EffectiveAction"

	for _, path := range source.Files(t, root) {
		set := token.NewFileSet()
		parsed, err := parser.ParseFile(set, path, nil, parser.SkipObjectResolution)
		if err != nil {
			t.Fatalf("parsing %s: %v", path, err)
		}

		holder := ""
		ast.Inspect(parsed, func(node ast.Node) bool {
			if function, isFunction := node.(*ast.FuncDecl); isFunction {
				holder = function.Name.Name
				return true
			}
			selector, isSelector := node.(*ast.SelectorExpr)
			if !isSelector || selector.Sel.Name != "Action" {
				return true
			}
			if qualifier, isIdentifier := selector.X.(*ast.Ident); isIdentifier && qualifier.Name == "findings" {
				return true
			}
			if path == owner && holder == decider {
				return true
			}
			t.Errorf("%s:%d reads a finding's raw action inside %s, and only %s may decide it",
				path, set.Position(selector.Pos()).Line, holder, decider)
			return true
		})
	}

	set := token.NewFileSet()
	parsed, err := parser.ParseFile(set, owner, nil, parser.SkipObjectResolution)
	if err != nil {
		t.Fatalf("parsing %s: %v", owner, err)
	}
	deciders := 0
	for _, declaration := range parsed.Decls {
		if function, isFunction := declaration.(*ast.FuncDecl); isFunction && function.Name.Name == decider {
			deciders++
		}
	}
	if deciders != 1 {
		t.Errorf("%d functions are named %s, want exactly one", deciders, decider)
	}
}

func onlyFinding(t *testing.T, finding string) findings.Finding {
	t.Helper()

	document, err := findings.Parse([]byte(`{"findings":[` + finding + `],"risk_level":"low","risk_rationale":"r"}`))
	if err != nil {
		t.Fatalf("reading the document: %v", err)
	}
	if len(document.Findings) != 1 {
		t.Fatalf("the document holds %d findings, want one", len(document.Findings))
	}
	return document.Findings[0]
}

func identifiers(held []findings.Finding) []string {
	var named []string
	for _, finding := range held {
		named = append(named, finding.ID)
	}
	return named
}
