package convention_test

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/internal/convention"
	"github.com/socrates-tools/socrates/internal/platform"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

func TestTheRequirementsComeBeforeTheExplanation(t *testing.T) {
	root := source.Root(t)

	required, err := convention.ToolchainFloorRequired(root)
	if err != nil {
		t.Fatalf("reading what the module requires: %v", err)
	}

	facts := append([]string{"claude", convention.ReviewerCLIFloor, "account", "Go " + required}, platform.Names()...)
	front, err := convention.RequirementsStatedFirst(root, facts)
	if err != nil {
		t.Fatalf("reading %s: %v", convention.FrontDocumentFile, err)
	}

	if len(front.Sections) == 0 {
		t.Fatalf("%s has no sections, so this check has nothing to read and would pass forever", convention.FrontDocumentFile)
	}
	if front.Sections[0] != convention.RequirementsHeading {
		t.Errorf("%s opens with %q rather than %q, so a reader learns how socrates works before whether they can use it at all", convention.FrontDocumentFile, front.Sections[0], convention.RequirementsHeading)
	}
	for _, missing := range front.Missing {
		t.Error(missing)
	}
}

func TestARequirementsSectionThatArrivesLateIsRead(t *testing.T) {
	root := t.TempDir()
	writeNested(t, filepath.Join(root, convention.FrontDocumentFile),
		"# socrates\n\n## How it works\n\nIt reviews.\n\n"+convention.RequirementsHeading+"\n\nThe claude CLI.\n")

	front, err := convention.RequirementsStatedFirst(root, []string{"claude", "account"})
	if err != nil {
		t.Fatalf("reading a fabricated front document: %v", err)
	}

	if front.Sections[0] == convention.RequirementsHeading {
		t.Errorf("a document explaining before it states its requirements was read as stating them first: %v", front.Sections)
	}
	want := []string{convention.FrontDocumentFile + " states what socrates review needs without naming account"}
	if diff := cmp.Diff(want, front.Missing); diff != "" {
		t.Errorf("what a requirements section missing a fact reports differs (-want +got):\n%s", diff)
	}
}

func TestAStatedToolchainFloorThatDisagreesWithTheModuleFails(t *testing.T) {
	root := source.Root(t)

	required, err := convention.ToolchainFloorRequired(root)
	if err != nil {
		t.Fatalf("reading what the module requires: %v", err)
	}

	stated, err := convention.ToolchainFloorsStated(root)
	if err != nil {
		t.Fatalf("reading what the documents state: %v", err)
	}
	if len(stated) == 0 {
		t.Fatalf("no document states a toolchain floor, so a reader is left to guess and this check has nothing to compare")
	}

	for _, disagreeing := range convention.FloorDisagreements(required, stated) {
		t.Errorf("%s: read the floor from the module rather than writing a second copy of it", disagreeing)
	}
}

func TestAFloorAheadOfTheModuleIsReportedNamingBoth(t *testing.T) {
	stated := []convention.StatedFloor{
		{Document: convention.FrontDocumentFile, Floor: "1.25"},
		{Document: convention.ContributingFile, Floor: "1.26"},
	}

	want := []string{convention.FrontDocumentFile + " states Go 1.25 and the module requires Go 1.26"}
	if diff := cmp.Diff(want, convention.FloorDisagreements("1.26", stated)); diff != "" {
		t.Errorf("what a document naming another floor reports differs (-want +got):\n%s", diff)
	}
}

func TestEveryInstalledBinaryIsExplained(t *testing.T) {
	root := source.Root(t)

	installed, err := convention.BinariesInstalled(root)
	if err != nil {
		t.Fatalf("reading what installing puts on a path: %v", err)
	}
	if len(installed) == 0 {
		t.Fatalf("installing puts nothing on a path, so this check has nothing to compare and would pass forever")
	}

	unexplained, err := convention.BinariesUnexplained(root)
	if err != nil {
		t.Fatalf("reading what %s explains: %v", convention.FrontDocumentFile, err)
	}
	for _, said := range unexplained {
		t.Errorf("%s, and a binary that appears with no explanation is indistinguishable from one nobody consented to", said)
	}
}

func TestABinaryNamedOnlyInPassingIsReported(t *testing.T) {
	root := t.TempDir()
	writeNested(t, filepath.Join(root, "Makefile"), "install:\n\tgo install ./cmd/socrates\n\tgo install ./cmd/agora\n")
	writeNested(t, filepath.Join(root, convention.FrontDocumentFile),
		"# socrates\n\n## What socrates is\n\nsocrates reviews. socrates freezes. socrates spawns.\n\n## Where it writes\n\nWhat agora writes is beside the runs.\n")

	unexplained, err := convention.BinariesUnexplained(root)
	if err != nil {
		t.Fatalf("reading a fabricated front document: %v", err)
	}

	want := []string{"make install puts agora on a path and " + convention.FrontDocumentFile + " explains it nowhere"}
	if diff := cmp.Diff(want, unexplained); diff != "" {
		t.Errorf("what a binary named only in passing reports differs (-want +got):\n%s", diff)
	}
}

func TestTheReachIsStatedInBothDirections(t *testing.T) {
	root := source.Root(t)

	unstated, err := convention.ReachUnstated(root)
	if err != nil {
		t.Fatalf("reading where what socrates reads goes: %v", err)
	}
	for _, said := range unstated {
		t.Error(said)
	}

	named, err := convention.NonAffiliationStated(root)
	if err != nil {
		t.Fatalf("reading %s: %v", convention.FrontDocumentFile, err)
	}
	if !named {
		t.Errorf("%s drives a vendor's CLI and never says socrates is %s with it", convention.FrontDocumentFile, convention.NotAffiliated)
	}
}

func TestStatingOnlyTheFlatteringHalfOfTheReachIsReported(t *testing.T) {
	root := t.TempDir()
	writeNested(t, filepath.Join(root, convention.FrontDocumentFile), "# socrates\n\n"+convention.SendsNothing+" anywhere.\n")
	writeNested(t, filepath.Join(root, convention.SecurityPolicyFile), "# Security\n\nReport it.\n")

	unstated, err := convention.ReachUnstated(root)
	if err != nil {
		t.Fatalf("reading a fabricated pair of documents: %v", err)
	}
	if len(unstated) != 2 {
		t.Fatalf("a document stating one half and one stating neither reported %d complaints: %v", len(unstated), unstated)
	}
	if !strings.Contains(unstated[0], convention.FrontDocumentFile) || !strings.Contains(unstated[0], convention.SendsToTheVendor) {
		t.Errorf("what a document stating only the flattering half reports is %q", unstated[0])
	}
	if !strings.Contains(unstated[1], "neither half") {
		t.Errorf("what a document stating neither half reports is %q", unstated[1])
	}
}

func TestAReporterHasAPrivateChannel(t *testing.T) {
	gaps, err := convention.SecurityPolicyGaps(source.Root(t))
	if err != nil {
		t.Fatalf("reading %s: %v", convention.SecurityPolicyFile, err)
	}
	for _, gap := range gaps {
		t.Error(gap)
	}
}

// scenario: The policy states what is out of scope
func TestAPolicyWithNoScopeSectionIsReported(t *testing.T) {
	root := t.TempDir()
	writeNested(t, filepath.Join(root, convention.SecurityPolicyFile),
		"# Security\n\nReport it at https://example.invalid/"+convention.PrivateChannelPath+" and you will hear back within 7 days.\n")

	gaps, err := convention.SecurityPolicyGaps(root)
	if err != nil {
		t.Fatalf("reading a fabricated policy: %v", err)
	}
	if len(gaps) != 1 || !strings.Contains(gaps[0], convention.OutOfScopeHeading) {
		t.Errorf("what a policy stating no scope reports is %v", gaps)
	}
}

func TestTheGuideLeadsWithTheProposal(t *testing.T) {
	reading, err := convention.HowTheGuideReads(source.Root(t))
	if err != nil {
		t.Fatalf("reading %s: %v", convention.ContributingFile, err)
	}

	if reading.ProposalAt < 0 {
		t.Fatalf("%s never names %s, so a contributor learns nowhere that a change starts as a proposal", convention.ContributingFile, convention.ProposalHome)
	}
	if reading.ToolchainAt < 0 {
		t.Fatalf("%s has no %q section, so this check has nothing to order it against", convention.ContributingFile, convention.ToolchainHeading)
	}
	if reading.ProposalAt > reading.ToolchainAt {
		t.Errorf("%s states the toolchain before it states that a change starts as a proposal in %s, and a contributor who learns the toolchain first writes code first", convention.ContributingFile, convention.ProposalHome)
	}
}

func TestTheGuideNamesOneVerificationCommand(t *testing.T) {
	reading, err := convention.HowTheGuideReads(source.Root(t))
	if err != nil {
		t.Fatalf("reading %s: %v", convention.ContributingFile, err)
	}

	if !reading.NamesEntry {
		t.Errorf("%s never names make %s, so how a change is checked is a list a contributor has to keep in step by hand", convention.ContributingFile, convention.TheEntryPoint)
	}
	for _, listing := range reading.Listing {
		t.Errorf("%s: name the one entry point instead, so a gate added to it is added to what a contributor runs", listing)
	}
}

func TestAGuideListingTheGatesItselfIsReported(t *testing.T) {
	root := t.TempDir()
	writeNested(t, filepath.Join(root, convention.ContributingFile),
		"# Contributing\n\nA change starts as a "+convention.ProposalHome+".\n\n"+convention.ToolchainHeading+
			"\n\nGo 1.26.\n\nRun make verify, then make corpus-outcomes.\n")

	reading, err := convention.HowTheGuideReads(root)
	if err != nil {
		t.Fatalf("reading a fabricated guide: %v", err)
	}

	want := []string{convention.ContributingFile + " tells a contributor to run make corpus-outcomes of its own"}
	if diff := cmp.Diff(want, reading.Listing); diff != "" {
		t.Errorf("what a guide listing a gate of its own reports differs (-want +got):\n%s", diff)
	}
}

func TestAFirstTimeContributorFindsTheDocumentsTheyNeed(t *testing.T) {
	missing, err := convention.DocumentsMissing(source.Root(t))
	if err != nil {
		t.Fatalf("reading what the repository carries: %v", err)
	}
	for _, said := range missing {
		t.Error(said)
	}
}

func TestADocumentAContributorNeedsAndCannotFindIsReported(t *testing.T) {
	missing, err := convention.DocumentsMissing(t.TempDir())
	if err != nil {
		t.Fatalf("reading an empty repository: %v", err)
	}
	if len(missing) == 0 {
		t.Fatal("an empty repository was read as carrying every document a first change needs")
	}
	for _, named := range []string{convention.CodeOfConductFile, convention.OwnershipFile, convention.IssueTemplateDirectory} {
		if !slices.ContainsFunc(missing, func(said string) bool { return strings.Contains(said, named) }) {
			t.Errorf("an empty repository was read as carrying %s: %v", named, missing)
		}
	}
}

func TestAContributorsEditorDoesNotFightTheConventions(t *testing.T) {
	disagreeing, err := convention.EditorSettingsDisagreeing(source.Root(t))
	if err != nil {
		t.Fatalf("reading the editor and attribute settings: %v", err)
	}
	for _, said := range disagreeing {
		t.Error(said)
	}
}

func TestSettingsThatFightTheConventionsAreReported(t *testing.T) {
	root := t.TempDir()
	writeNested(t, filepath.Join(root, convention.EditorSettingsFile), "root = true\ncharset = utf-8\n")
	writeNested(t, filepath.Join(root, convention.AttributeSettingsFile), "* text=auto\n")

	disagreeing, err := convention.EditorSettingsDisagreeing(root)
	if err != nil {
		t.Fatalf("reading fabricated settings: %v", err)
	}

	for _, wanted := range []string{"end_of_line = lf", "insert_final_newline = true", "indent_style = tab", "eol=lf"} {
		if !slices.ContainsFunc(disagreeing, func(said string) bool { return strings.Contains(said, wanted) }) {
			t.Errorf("settings stating no %q were read as agreeing with the conventions: %v", wanted, disagreeing)
		}
	}
}
