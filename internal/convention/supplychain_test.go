package convention_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ktei/socrates/internal/convention"
	"github.com/ktei/socrates/internal/platform"
	"github.com/ktei/socrates/internal/testing/source"
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
		t.Fatalf("no workflow runs %s, so nothing reads the module graph for a published advisory at all", convention.TheVulnerabilityCheck)
	}

	for _, named := range checking {
		start, err := convention.WhatStartsAWorkflow(root, named)
		if err != nil {
			t.Fatalf("reading what starts %s: %v", named, err)
		}
		if !start.OnRequest {
			t.Errorf("%s runs %s and nobody can ask for one: state %s, because an advisory now waits for somebody to ask", named, convention.TheVulnerabilityCheck, convention.OnRequest)
		}
	}
}

func TestARunIsStartedByAPersonOrByARelease(t *testing.T) {
	root := source.Root(t)

	gates, err := convention.WorkflowsRunningTheGates(root)
	if err != nil {
		t.Fatalf("reading which workflows run the suite: %v", err)
	}
	checking, err := convention.WorkflowsChecking(root, convention.TheVulnerabilityCheck)
	if err != nil {
		t.Fatalf("reading what the workflows check: %v", err)
	}

	asked := append(slices.Clone(gates.Invoking), checking...)
	if len(asked) == 0 {
		t.Fatal("no workflow runs the suite or reads the module graph, so this check has nothing to read and would pass forever")
	}

	for _, named := range asked {
		start, err := convention.WhatStartsAWorkflow(root, named)
		if err != nil {
			t.Fatalf("reading what starts %s: %v", named, err)
		}
		for _, unasked := range start.Unasked {
			t.Errorf("%s starts on %s, and a run nobody asked for is one this account cannot pay for", named, unasked)
		}
		if !start.OnRequest {
			t.Errorf("%s states no %s, so nobody can start it at all", named, convention.OnRequest)
		}
	}
}

func TestAWorkflowStartingOnAPushIsReported(t *testing.T) {
	root := t.TempDir()
	writeWorkflow(t, root, "eager.yml", "on:\n  push:\n    branches: [main]\n  pull_request:\n  workflow_dispatch:\n\njobs:\n  one:\n    steps:\n      - run: make verify\n")

	start, err := convention.WhatStartsAWorkflow(root, "eager.yml")
	if err != nil {
		t.Fatalf("reading a fabricated workflow: %v", err)
	}

	if diff := cmp.Diff([]string{"push", "pull_request"}, start.Unasked); diff != "" {
		t.Errorf("what a workflow starting on a push reports differs (-want +got):\n%s", diff)
	}
	if !start.OnRequest {
		t.Error("a workflow stating workflow_dispatch was read as unstartable")
	}
}

func TestAReleaseStillGatesOnTheFullMatrix(t *testing.T) {
	root := source.Root(t)

	gate, err := convention.HowTheReleaseIsGated(root)
	if err != nil {
		t.Fatalf("reading how a release is gated: %v", err)
	}
	if gate.Verifying == "" {
		t.Fatalf("no job in %s runs the workflow that invokes make %s, so a release gates on nothing", convention.ReleaseWorkflow, convention.TheEntryPoint)
	}
	for _, ungated := range gate.Ungated {
		t.Errorf("the job %s in %s does not need %s", ungated, convention.ReleaseWorkflow, gate.Verifying)
	}

	reading, err := convention.WorkflowsReadingTheStatement(root, platform.ReleaseTargetsFile)
	if err != nil {
		t.Fatalf("reading what the suite runs on: %v", err)
	}
	if len(reading.Invoking) == 0 {
		t.Fatalf("the workflow a release gates on runs the suite without reading %s, so the matrix it runs is a second copy", platform.ReleaseTargetsFile)
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
