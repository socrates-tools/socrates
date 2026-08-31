package convention_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-devtool/socrates/internal/convention"
	"github.com/socrates-devtool/socrates/internal/platform"
	"github.com/socrates-devtool/socrates/internal/testing/source"
)

func TestNoActionIsNamedByATag(t *testing.T) {
	tagged, err := convention.ActionsNamedByATag(source.Root(t))
	if err != nil {
		t.Fatalf("reading what the workflows use: %v", err)
	}

	for _, named := range tagged {
		t.Errorf("%s names %s by %q, which somebody outside this repository can move: pin it to a commit with its version beside it", named.Workflow, named.Uses, named.Named)
	}
}

func TestAnActionNamedByATagFailsTheCheck(t *testing.T) {
	root := t.TempDir()
	writeWorkflow(t, root, "moving.yml", "jobs:\n  one:\n    steps:\n      - uses: actions/checkout@v7\n      - uses: actions/setup-go@0123456789abcdef0123456789abcdef01234567 # v7.0.0\n")

	tagged, err := convention.ActionsNamedByATag(root)
	if err != nil {
		t.Fatalf("reading a fabricated workflow: %v", err)
	}

	want := []convention.NamedAction{{Workflow: "moving.yml", Uses: "actions/checkout", Named: "v7"}}
	if diff := cmp.Diff(want, tagged); diff != "" {
		t.Errorf("what a workflow naming an action by a tag reports differs (-want +got):\n%s", diff)
	}
}

func TestVerificationRunsReadOnly(t *testing.T) {
	asked, err := convention.TokensWorkflowsAskFor(source.Root(t))
	if err != nil {
		t.Fatalf("reading what the workflows ask for: %v", err)
	}
	if len(asked) == 0 {
		t.Fatal("there is no workflow to read, so this check would pass forever")
	}

	for _, claim := range asked {
		if claim.Stated == "" {
			t.Errorf("%s states no permissions, so what its token may do is a repository setting nobody reviews: state %q at the top", claim.Workflow, convention.ReadOnlyToken)
			continue
		}
		if strings.TrimSpace(claim.Stated) != convention.ReadOnlyToken {
			t.Errorf("%s asks for %q at the top rather than %q: escalate in the job that needs more", claim.Workflow, claim.Stated, convention.ReadOnlyToken)
		}
	}
}

func TestASupersededRunDoesNotReport(t *testing.T) {
	superseding, err := convention.WhatWorkflowsSupersede(source.Root(t))
	if err != nil {
		t.Fatalf("reading what the workflows supersede: %v", err)
	}
	if len(superseding) == 0 {
		t.Fatal("there is no workflow to read, so this check would pass forever")
	}

	for _, claim := range superseding {
		if !strings.Contains(claim.Stated, convention.CancelSuperseded) {
			t.Errorf("%s states %q, so a run left behind by a newer commit keeps going and reports about a state the branch has left: state %q", claim.Workflow, claim.Stated, convention.CancelSuperseded)
		}
	}
}

func TestThePageHalfRunsAgainstAKnownBrowser(t *testing.T) {
	driving, err := convention.BrowsersWorkflowsDrive(source.Root(t))
	if err != nil {
		t.Fatalf("reading how the workflows find a browser: %v", err)
	}
	if len(driving) == 0 {
		t.Fatalf("no workflow runs make %s, so nothing drives the page and this check would pass forever", convention.TheEntryPoint)
	}

	for _, found := range driving {
		if found.Version == "" {
			t.Errorf("%s drives the page against no stated browser version, so a runner image change decides what the page half runs against", found.Workflow)
		}
		if found.FromTheRunner {
			t.Errorf("%s resolves a browser from the runner image, and its removal would take the half of the suite that drives a real page with it", found.Workflow)
		}
		if !found.NamesTheSeam {
			t.Errorf("%s installs a browser and does not name %s, so the suite drives whatever it finds instead", found.Workflow, convention.TheBrowserSeam)
		}
	}
}

func TestAKnownVulnerabilityInADependencyIsReported(t *testing.T) {
	root := source.Root(t)

	checking, err := convention.WorkflowsChecking(root, convention.TheVulnerabilityCheck)
	if err != nil {
		t.Fatalf("reading what the workflows check: %v", err)
	}
	if len(checking) == 0 {
		t.Fatalf("no workflow runs %s, so an advisory published against the module graph waits for somebody to look", convention.TheVulnerabilityCheck)
	}

	when, err := convention.WhenWorkflowsRun(root)
	if err != nil {
		t.Fatalf("reading when the workflows run: %v", err)
	}

	for _, claim := range when {
		if !slices.Contains(checking, claim.Workflow) {
			continue
		}
		for _, trigger := range []string{"schedule", "push"} {
			if !strings.Contains(claim.Stated, trigger) {
				t.Errorf("%s runs %s and not on %s, and an advisory published between runs has to arrive without anybody looking", claim.Workflow, convention.TheVulnerabilityCheck, trigger)
			}
		}
	}
}

func TestTheSuiteRunsOnEachSupportedPlatform(t *testing.T) {
	reading, err := convention.WorkflowsReadingTheStatement(source.Root(t), platform.ReleaseTargetsFile)
	if err != nil {
		t.Fatalf("reading what the workflows run the suite on: %v", err)
	}
	if len(reading.Invoking) == 0 {
		t.Fatalf("no workflow runs the suite from %s, so what it proves is about a platform list somebody wrote twice", platform.ReleaseTargetsFile)
	}
	for _, listing := range reading.Listing {
		t.Errorf("%s: derive the matrix from it, so a platform is run because the set names it", listing)
	}
}

func TestADependencyUpdateArrivesAsAChangeSomebodyReviews(t *testing.T) {
	unwatched, err := convention.EcosystemsUnwatched(source.Root(t))
	if err != nil {
		t.Fatalf("reading %s: %v", convention.DependabotFile, err)
	}

	for _, said := range unwatched {
		t.Errorf("%s, so an update there arrives silently or not at all", said)
	}
}
