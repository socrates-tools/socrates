package convention_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-devtool/socrates/internal/convention"
	"github.com/socrates-devtool/socrates/internal/testing/source"
)

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
