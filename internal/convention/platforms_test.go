package convention_test

import (
	"path/filepath"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-devtool/socrates/internal/convention"
	"github.com/socrates-devtool/socrates/internal/platform"
	"github.com/socrates-devtool/socrates/internal/testing/source"
)

const anUnsupportedPlatform = "windows"

func TestAReleaseBuildsTheSupportedPlatformsAndNoOthers(t *testing.T) {
	root := source.Root(t)

	built, err := convention.PlatformsAReleaseBuilds(root, platform.ReleaseTargetsFile)
	if err != nil {
		t.Fatalf("reading what a release builds: %v", err)
	}
	if len(built) == 0 {
		t.Fatalf("%s names no platform, so this check has nothing to compare and would pass forever", platform.ReleaseTargetsFile)
	}

	var operating, running []string
	for _, one := range built {
		operating = append(operating, one.OS)
		running = append(running, one.Runner)
	}

	for _, disagreeing := range convention.PlatformDisagreements(platform.OperatingSystems(), operating, platform.ReleaseTargetsFile) {
		t.Errorf("%s: regenerate it with make platforms", disagreeing)
	}
	for _, disagreeing := range convention.PlatformDisagreements(platform.Runners(), running, platform.ReleaseTargetsFile) {
		t.Errorf("%s: regenerate it with make platforms", disagreeing)
	}
}

func TestAPlatformInTheMatrixButNotInTheSupportedSetFails(t *testing.T) {
	named, err := convention.RunnersNamedInWorkflows(source.Root(t))
	if err != nil {
		t.Fatalf("reading what the workflows run on: %v", err)
	}
	if len(named) == 0 {
		t.Fatal("no workflow names a runner at all, so this check has nothing to read and would pass forever")
	}

	for _, running := range named {
		if !slices.Contains(platform.Runners(), running.Runner) {
			t.Errorf("%s runs on %s, and the supported set names %v: read the matrix from %s rather than naming a platform in a workflow", running.Workflow, running.Runner, platform.Runners(), platform.ReleaseTargetsFile)
		}
	}
}

func TestAPlatformAddedToAReleaseWithoutBeingStatedFailsTheSuite(t *testing.T) {
	built := append(platform.OperatingSystems(), anUnsupportedPlatform)

	want := []string{platform.ReleaseTargetsFile + " names " + anUnsupportedPlatform + ", and the supported set does not"}
	if diff := cmp.Diff(want, convention.PlatformDisagreements(platform.OperatingSystems(), built, platform.ReleaseTargetsFile)); diff != "" {
		t.Errorf("what a release building an unstated platform reports differs (-want +got):\n%s", diff)
	}
}

func TestAPlatformStatedButNeverBuiltFailsTheSuite(t *testing.T) {
	stated := platform.OperatingSystems()
	if len(stated) < 2 {
		t.Fatalf("the supported set names %v, and this check needs one to leave out", stated)
	}
	left := stated[len(stated)-1]

	want := []string{"the supported set names " + left + ", and " + platform.ReleaseTargetsFile + " does not"}
	if diff := cmp.Diff(want, convention.PlatformDisagreements(stated, stated[:len(stated)-1], platform.ReleaseTargetsFile)); diff != "" {
		t.Errorf("what a stated platform no release builds reports differs (-want +got):\n%s", diff)
	}
}

func TestASecondCopyOfTheSupportedSetIsReported(t *testing.T) {
	root := t.TempDir()
	writeNested(t, filepath.Join(root, "internal", "thing", "thing.go"), "package thing\n\nvar built = []string{\"linux\", \"darwin\"}\n")

	copies, err := convention.PlatformsNamedOutsideTheStatement(root, []string{"linux", "darwin"})
	if err != nil {
		t.Fatalf("reading the fabricated package: %v", err)
	}

	want := []string{
		filepath.Join("internal", "thing", "thing.go") + `:3 states "darwin"`,
		filepath.Join("internal", "thing", "thing.go") + `:3 states "linux"`,
	}
	if diff := cmp.Diff(want, copies); diff != "" {
		t.Errorf("the reported second copies differ (-want +got):\n%s", diff)
	}
}

func TestTheStatementItselfIsNotASecondCopy(t *testing.T) {
	root := t.TempDir()
	writeNested(t, filepath.Join(root, filepath.FromSlash(convention.StatementDirectory), "platform.go"), "package platform\n\nvar stated = []string{\"linux\"}\n")

	copies, err := convention.PlatformsNamedOutsideTheStatement(root, []string{"linux"})
	if err != nil {
		t.Fatalf("reading the fabricated statement: %v", err)
	}
	if len(copies) != 0 {
		t.Errorf("the one place the set is stated was reported as a second copy: %v", copies)
	}
}

func TestWhatAPersonReadsMatchesWhatTheSourceStates(t *testing.T) {
	root := source.Root(t)

	documented, err := convention.PlatformsDocumented(root)
	if err != nil {
		t.Fatalf("reading what a person is told: %v", err)
	}
	if len(documented) == 0 {
		t.Fatalf("%s has no %q section naming the platforms, so a stranger is left to assume socrates runs wherever Go runs", "README.md", convention.DocumentedHeading)
	}

	for _, disagreeing := range convention.PlatformDisagreements(platform.Names(), documented, "README.md") {
		t.Error(disagreeing)
	}

	named, err := convention.TheWayForwardIsDocumented(root)
	if err != nil {
		t.Fatalf("reading what a person is told to do instead: %v", err)
	}
	if !named {
		t.Errorf("the section does not name %s and %s, and the refusal points a person at the second of them", convention.WhereTheWayForwardApplies, convention.TheWayForward)
	}
}

func TestAPlatformDocumentedInNoTableIsRead(t *testing.T) {
	root := t.TempDir()
	writeNested(t, filepath.Join(root, "README.md"), "# thing\n\n"+convention.DocumentedHeading+"\n\nIt runs on Linux.\n\n## something else\n")

	documented, err := convention.PlatformsDocumented(root)
	if err != nil {
		t.Fatalf("reading the fabricated documentation: %v", err)
	}
	if len(documented) != 0 {
		t.Errorf("prose outside a table was read as a documented platform: %v", documented)
	}
}

func TestTheDocumentedPlatformsAreTheRowsOfTheTable(t *testing.T) {
	root := t.TempDir()
	writeNested(t, filepath.Join(root, "README.md"), "# thing\n\n"+convention.DocumentedHeading+
		"\n\n| platform | what holds |\n| --- | --- |\n| Linux | both |\n| macOS | both |\n\nOn Windows, use WSL.\n\n## something else\n")

	documented, err := convention.PlatformsDocumented(root)
	if err != nil {
		t.Fatalf("reading the fabricated documentation: %v", err)
	}
	if diff := cmp.Diff([]string{"Linux", "macOS"}, documented); diff != "" {
		t.Errorf("what the table documents differs (-want +got):\n%s", diff)
	}

	named, err := convention.TheWayForwardIsDocumented(root)
	if err != nil {
		t.Fatalf("reading the fabricated documentation: %v", err)
	}
	if !named {
		t.Error("a section naming Windows and WSL was read as naming neither")
	}
}
