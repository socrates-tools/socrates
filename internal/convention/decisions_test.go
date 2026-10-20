package convention_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/internal/convention"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

var twoTools = []convention.Tool{
	{Namespace: "review", Directory: "review"},
	{Namespace: "agora", Directory: "agora"},
}

func writeDecision(t *testing.T, root, directory string, number int, header, body string) {
	t.Helper()

	name := filepath.Join(root, directory, "docs", "adr", numbered(number)+"-a-decision.md")
	writeNested(t, name, "# "+header+"\n\n- Status: accepted\n\n"+body+"\n")
}

func numbered(number int) string {
	return map[int]string{1: "0001", 2: "0002", 3: "0003", 7: "0007"}[number]
}

func TestEveryDecisionRecordHoldsANumberOnlyOneRecordHolds(t *testing.T) {
	records, err := convention.DecisionRecords(source.Root(t), convention.Tools)
	if err != nil {
		t.Fatalf("reading the decision records: %v", err)
	}
	if len(records) == 0 {
		t.Fatal("no decision record was found, so this check has nothing to read and would pass forever")
	}

	for _, number := range convention.NumbersHeldByMoreThanOneRecord(records) {
		t.Errorf("%s is the number of more than one record: %v", number, records[number])
	}
}

func TestEveryDecisionRecordIsNamedForTheNumberItsTitleCarries(t *testing.T) {
	disagreeing, err := convention.RecordsTitledForAnotherNumber(source.Root(t), convention.Tools)
	if err != nil {
		t.Fatalf("reading the decision records: %v", err)
	}
	for _, record := range disagreeing {
		t.Errorf("%s", record)
	}
}

func TestEveryDecisionReferenceNamesARecordThatExists(t *testing.T) {
	unresolved, err := convention.ReferencesToNoDecision(source.Root(t), convention.Tools)
	if err != nil {
		t.Fatalf("reading the documents: %v", err)
	}
	for _, reference := range unresolved {
		t.Errorf("%s", reference)
	}
}

func TestEveryRelativeLinkInADocumentReachesAFile(t *testing.T) {
	broken, err := convention.LinksReachingNothing(source.Root(t))
	if err != nil {
		t.Fatalf("reading the documents: %v", err)
	}
	for _, link := range broken {
		t.Errorf("%s", link)
	}
}

func TestEveryToolHasAGlossaryADecisionDirectoryAndInstructionsTheMapNames(t *testing.T) {
	gaps, err := convention.ContextGaps(source.Root(t), convention.Tools)
	if err != nil {
		t.Fatalf("reading the documents: %v", err)
	}
	for _, gap := range gaps {
		t.Error(gap)
	}
}

func TestAgorasGlossaryHoldsNoWordOfTheOtherTool(t *testing.T) {
	naming, err := convention.GlossaryNamingTheOtherTool(source.Root(t), convention.ToolDirectory("agora"))
	if err != nil {
		t.Fatalf("reading the glossary: %v", err)
	}
	for _, named := range naming {
		t.Errorf("%s, and agora's defining invariant is that it does not know what a review is", named)
	}
}

func TestTwoRecordsHoldingOneNumberInDifferentDirectoriesAreReported(t *testing.T) {
	root := t.TempDir()
	writeDecision(t, root, "", 7, "7. First", "Decided.")
	writeDecision(t, root, "agora", 7, "7. Second", "Decided.")

	records, err := convention.DecisionRecords(root, twoTools)
	if err != nil {
		t.Fatalf("reading a fabricated tree: %v", err)
	}

	if diff := cmp.Diff([]string{"0007"}, convention.NumbersHeldByMoreThanOneRecord(records)); diff != "" {
		t.Errorf("the numbers held twice differ (-want +got):\n%s", diff)
	}
}

func TestARecordWhoseTitleCarriesAnotherNumberThanItsNameIsReported(t *testing.T) {
	root := t.TempDir()
	writeDecision(t, root, "review", 2, "3. Named for another", "Decided.")

	got, err := convention.RecordsTitledForAnotherNumber(root, twoTools)
	if err != nil {
		t.Fatalf("reading a fabricated tree: %v", err)
	}

	want := []string{filepath.Join("review", "docs", "adr", "0002-a-decision.md") + " is titled for 3"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("the records titled for another number differ (-want +got):\n%s", diff)
	}
}

func TestAReferenceToADecisionNobodyHoldsIsReportedAndOneToAMovedRecordIsNot(t *testing.T) {
	root := t.TempDir()
	writeDecision(t, root, "agora", 1, "1. Moved", "Decided.")
	writeNested(t, filepath.Join(root, "README.md"), "It follows ADR 1 and ADR 0001, and contradicts ADR 9.\n")

	got, err := convention.ReferencesToNoDecision(root, twoTools)
	if err != nil {
		t.Fatalf("reading a fabricated tree: %v", err)
	}

	want := []string{"README.md:1 names ADR 9, which no directory holds"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("the references reaching nothing differ (-want +got):\n%s", diff)
	}
}

func TestARelativeLinkToAFileThatMovedIsReportedAndOneThatResolvesIsNot(t *testing.T) {
	root := t.TempDir()
	writeNested(t, filepath.Join(root, "docs", "kept.md"), "kept\n")
	writeNested(t, filepath.Join(root, "docs", "reader.md"), "[kept](kept.md) [moved](gone.md#top) [site](https://example.com/x.md) [here](#section)\n")

	got, err := convention.LinksReachingNothing(root)
	if err != nil {
		t.Fatalf("reading a fabricated tree: %v", err)
	}

	want := []string{filepath.Join("docs", "reader.md") + ":1 links to gone.md, which does not exist"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("the links reaching nothing differ (-want +got):\n%s", diff)
	}
}

func TestAToolWithNoGlossaryNoDecisionDirectoryAndNoInstructionsIsReportedWhereverTheMapIsSilent(t *testing.T) {
	root := t.TempDir()
	writeNested(t, filepath.Join(root, "CONTEXT.md"), "shared\n")
	writeNested(t, filepath.Join(root, "CONTEXT-MAP.md"), "See CONTEXT.md and review/CONTEXT.md.\n")
	writeNested(t, filepath.Join(root, "AGENTS.md"), "map\n")
	writeNested(t, filepath.Join(root, "review", "CONTEXT.md"), "review\n")
	writeNested(t, filepath.Join(root, "review", "AGENTS.md"), "review\n")
	writeDecision(t, root, "review", 1, "1. Held", "Decided.")
	writeNested(t, filepath.Join(root, "CLAUDE.md"), "@AGENTS.md\n")
	writeNested(t, filepath.Join(root, "review", "CLAUDE.md"), "@AGENTS.md\n")

	got, err := convention.ContextGaps(root, twoTools)
	if err != nil {
		t.Fatalf("reading a fabricated tree: %v", err)
	}

	want := []string{
		"CONTEXT-MAP.md does not name agora/CONTEXT.md",
		"agora/CONTEXT.md does not exist",
		"agora/AGENTS.md does not exist",
		"agora/CLAUDE.md does not hold @AGENTS.md alone",
		"agora/docs/adr holds no decision record",
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("the gaps differ (-want +got):\n%s", diff)
	}
}

func TestAClaudeInstructionsFileThatLinksToAgentsRatherThanImportingItIsReported(t *testing.T) {
	root := t.TempDir()
	writeNested(t, filepath.Join(root, "CONTEXT.md"), "shared\n")
	writeNested(t, filepath.Join(root, "CONTEXT-MAP.md"), "See CONTEXT.md and review/CONTEXT.md.\n")
	writeNested(t, filepath.Join(root, "AGENTS.md"), "map\n")
	writeNested(t, filepath.Join(root, "review", "CONTEXT.md"), "review\n")
	writeNested(t, filepath.Join(root, "review", "AGENTS.md"), "review\n")
	writeDecision(t, root, "review", 1, "1. Held", "Decided.")
	if err := os.Symlink("AGENTS.md", filepath.Join(root, "CLAUDE.md")); err != nil {
		t.Fatalf("linking the instructions: %v", err)
	}
	writeNested(t, filepath.Join(root, "review", "CLAUDE.md"), "@AGENTS.md\nAnd one more rule.\n")

	got, err := convention.ContextGaps(root, []convention.Tool{{Namespace: "review", Directory: "review"}})
	if err != nil {
		t.Fatalf("reading a fabricated tree: %v", err)
	}

	want := []string{
		"CLAUDE.md does not hold @AGENTS.md alone",
		"review/CLAUDE.md does not hold @AGENTS.md alone",
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("the gaps differ (-want +got):\n%s", diff)
	}
}

func TestAWordOfTheOtherToolInAGlossaryIsReported(t *testing.T) {
	root := t.TempDir()
	writeNested(t, filepath.Join(root, "agora", "CONTEXT.md"), "- **option**: one answer.\n- **verdict**: what was concluded.\n")

	got, err := convention.GlossaryNamingTheOtherTool(root, "agora")
	if err != nil {
		t.Fatalf("reading a fabricated tree: %v", err)
	}

	want := []string{filepath.Join("agora", "CONTEXT.md") + `:2 names "verdict"`}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("what a glossary naming the other tool reports differs (-want +got):\n%s", diff)
	}
}
