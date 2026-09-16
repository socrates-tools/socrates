package convention_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/internal/convention"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

// scenario: A verification run reaches no reviewer CLI
func TestNoWorkflowReachesTheRealCLI(t *testing.T) {
	reach, err := convention.WorkflowsReachingTheRealCLI(source.Root(t))
	if err != nil {
		t.Fatalf("reading the workflows: %v", err)
	}

	for _, reaching := range reach.Reaching {
		t.Errorf("%s, and nothing in continuous integration spends on the vendor a person can spend for themselves", reaching)
	}
	if len(reach.Workflows) == 0 {
		t.Fatalf("there is no workflow under %s, so this check has nothing to read and would pass forever", convention.WorkflowsDirectory)
	}
}

// scenario: A verification run reaches no forge
func TestNoWorkflowReachesAForge(t *testing.T) {
	reach, err := convention.WorkflowsReachingAForge(source.Root(t))
	if err != nil {
		t.Fatalf("reading the workflows: %v", err)
	}

	for _, reaching := range reach.Reaching {
		t.Errorf("%s, and nothing in continuous integration writes in a repository socrates does not own", reaching)
	}
	if len(reach.Workflows) == 0 {
		t.Fatalf("there is no workflow under %s, so this check has nothing to read and would pass forever", convention.WorkflowsDirectory)
	}
}

func TestAWorkflowReachingAForgeIsReported(t *testing.T) {
	root := t.TempDir()
	writeWorkflow(t, root, "posting.yml", "steps:\n  - run: make forge-smoke\n  - run: go test -tags forgesmoke ./internal/fidelity/\n")

	reach, err := convention.WorkflowsReachingAForge(root)
	if err != nil {
		t.Fatalf("reading a fabricated workflow: %v", err)
	}

	want := []string{
		"posting.yml runs make forge-smoke, which reaches a forge",
		"posting.yml compiles the forgesmoke tag, which is behind the paths that reach a forge",
	}
	if diff := cmp.Diff(want, reach.Reaching); diff != "" {
		t.Errorf("what a workflow reaching a forge reports differs (-want +got):\n%s", diff)
	}
}

func TestTheRunInvokesTheEntryPointRatherThanTheGates(t *testing.T) {
	gates, err := convention.WorkflowsRunningTheGates(source.Root(t))
	if err != nil {
		t.Fatalf("reading the workflows: %v", err)
	}

	if len(gates.Invoking) == 0 {
		t.Errorf("no workflow runs make %s, so what a run proves and what a contributor runs are two lists that can disagree", convention.TheEntryPoint)
	}
	for _, listing := range gates.Listing {
		t.Errorf("%s: run make %s instead, so a gate added to it is added to both", listing, convention.TheEntryPoint)
	}
}

func TestAWorkflowListingAGateOfItsOwnIsReported(t *testing.T) {
	root := t.TempDir()
	writeWorkflow(t, root, "listing.yml", "steps:\n  - run: make verify\n  - run: make corpus-outcomes\n")

	gates, err := convention.WorkflowsRunningTheGates(root)
	if err != nil {
		t.Fatalf("reading a fabricated workflow: %v", err)
	}

	want := []string{"listing.yml runs make corpus-outcomes of its own"}
	if diff := cmp.Diff(want, gates.Listing); diff != "" {
		t.Errorf("what a workflow running a gate of its own reports differs (-want +got):\n%s", diff)
	}
}

func TestAWorkflowRecordingAgainstTheRealCLIIsReported(t *testing.T) {
	root := t.TempDir()
	writeWorkflow(t, root, "paid.yml", "steps:\n  - run: make fixtures\n")

	reach, err := convention.WorkflowsReachingTheRealCLI(root)
	if err != nil {
		t.Fatalf("reading a fabricated workflow: %v", err)
	}

	want := []string{"paid.yml runs make fixtures, which reaches the real CLI"}
	if diff := cmp.Diff(want, reach.Reaching); diff != "" {
		t.Errorf("what a workflow running make fixtures reports differs (-want +got):\n%s", diff)
	}
}

func TestAWorkflowNamingTheCredentialIsReported(t *testing.T) {
	root := t.TempDir()
	writeWorkflow(t, root, "keyed.yml", "env:\n  ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}\n")

	reach, err := convention.WorkflowsReachingTheRealCLI(root)
	if err != nil {
		t.Fatalf("reading a fabricated workflow: %v", err)
	}

	want := []string{"keyed.yml names ANTHROPIC_API_KEY, which is a credential the real CLI reads"}
	if diff := cmp.Diff(want, reach.Reaching); diff != "" {
		t.Errorf("what a workflow naming the credential reports differs (-want +got):\n%s", diff)
	}
}

func TestAWorkflowCompilingATagThatReachesTheRealCLIIsReported(t *testing.T) {
	root := t.TempDir()
	writeWorkflow(t, root, "tagged.yml", "steps:\n  - run: go test -tags livesmoke ./internal/fidelity/\n")

	reach, err := convention.WorkflowsReachingTheRealCLI(root)
	if err != nil {
		t.Fatalf("reading a fabricated workflow: %v", err)
	}

	want := []string{"tagged.yml compiles the livesmoke tag, which is behind the paths that reach the real CLI"}
	if diff := cmp.Diff(want, reach.Reaching); diff != "" {
		t.Errorf("what a workflow compiling the live smoke tag reports differs (-want +got):\n%s", diff)
	}
}

func TestComparingTheSimulatorAgainstTheRecordingsIsNotReaching(t *testing.T) {
	root := t.TempDir()
	writeWorkflow(t, root, "free.yml", "steps:\n  - run: make fixtures-diff\n")

	reach, err := convention.WorkflowsReachingTheRealCLI(root)
	if err != nil {
		t.Fatalf("reading a fabricated workflow: %v", err)
	}

	if len(reach.Reaching) != 0 {
		t.Errorf("make fixtures-diff was reported as reaching the real CLI: %v", reach.Reaching)
	}
}

func writeWorkflow(t *testing.T, root, name, content string) {
	t.Helper()

	directory := filepath.Join(root, convention.WorkflowsDirectory)
	if err := os.MkdirAll(directory, 0o755); err != nil {
		t.Fatalf("creating %s: %v", directory, err)
	}
	write(t, filepath.Join(directory, name), content)
}
