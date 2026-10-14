package convention_test

import (
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/internal/convention"
	"github.com/socrates-tools/socrates/internal/platform"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

var aSemanticVersion = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)

func TestTheReleaseMatrixIsTheSupportedSet(t *testing.T) {
	root := source.Root(t)

	building, err := convention.AReleaseBuilds(root)
	if err != nil {
		t.Fatalf("reading what a release builds: %v", err)
	}
	if len(building) == 0 {
		t.Fatalf("%s builds nothing, so this check has nothing to compare and would pass forever", convention.ReleaseConfigurationFile)
	}

	for _, one := range building {
		where := convention.ReleaseConfigurationFile + " building " + one.Binary
		for _, disagreeing := range convention.PlatformDisagreements(platform.OperatingSystems(), one.OperatingSystems, where) {
			t.Errorf("%s: a release builds the supported set and no others", disagreeing)
		}
	}
}

func TestATargetOutsideTheSupportedSetIsReported(t *testing.T) {
	root := t.TempDir()
	writeNested(t, filepath.Join(root, convention.ReleaseConfigurationFile),
		"builds:\n  - id: one\n    binary: one\n    goos:\n      - linux\n      - darwin\n      - "+anUnsupportedPlatform+"\n")

	built, err := convention.OperatingSystemsAReleaseBuilds(root)
	if err != nil {
		t.Fatalf("reading a fabricated release: %v", err)
	}

	want := []string{convention.ReleaseConfigurationFile + " names " + anUnsupportedPlatform + ", and the supported set does not"}
	if diff := cmp.Diff(want, convention.PlatformDisagreements([]string{"linux", "darwin"}, built, convention.ReleaseConfigurationFile)); diff != "" {
		t.Errorf("what a release building an unsupported target reports differs (-want +got):\n%s", diff)
	}
}

func TestAReleaseBuildsOneBinaryAndTheArchiveCarriesItWithTheLicenceAndNotice(t *testing.T) {
	root := source.Root(t)

	building, err := convention.AReleaseBuilds(root)
	if err != nil {
		t.Fatalf("reading what a release builds: %v", err)
	}

	tool := convention.BuildsProducing(building, convention.Suite)
	if len(tool) != 1 {
		t.Fatalf("%s produces %d binaries named %s, and a release carries one", convention.ReleaseConfigurationFile, len(tool), convention.Suite)
	}
	if len(building) != 1 {
		t.Errorf("%s produces %d binaries, and the suite publishes one", convention.ReleaseConfigurationFile, len(building))
	}

	archives, err := convention.AReleaseArchives(root)
	if err != nil {
		t.Fatalf("reading what a release archives: %v", err)
	}
	if len(archives) == 0 {
		t.Fatalf("%s archives nothing, so nothing a person downloads is described there", convention.ReleaseConfigurationFile)
	}

	for _, archive := range archives {
		if !slices.Contains(archive.Builds, tool[0].ID) {
			t.Errorf("the archive %s carries %v, and not %s", archive.ID, archive.Builds, tool[0].ID)
		}
		for _, stated := range []string{convention.LicenceFile, convention.NoticeFile} {
			if !slices.Contains(archive.Files, stated) {
				t.Errorf("the archive %s carries %v, and not %s, which is a condition on the copy", archive.ID, archive.Files, stated)
			}
		}
	}
}

func TestAReleaseThatDidNotPassVerificationProducesNothing(t *testing.T) {
	gate, err := convention.HowTheReleaseIsGated(source.Root(t))
	if err != nil {
		t.Fatalf("reading how a release is gated: %v", err)
	}
	if gate.Verifying == "" {
		t.Fatalf("no job in %s runs the workflow that invokes make %s, so a tag publishes whatever it was pushed on", convention.ReleaseWorkflow, convention.TheEntryPoint)
	}
	if len(gate.Publishes) == 0 {
		t.Fatalf("%s verifies and publishes nothing, so this check has nothing to hold and would pass forever", convention.ReleaseWorkflow)
	}

	for _, ungated := range gate.Ungated {
		t.Errorf("the job %s in %s does not need %s, so it produces artifacts for a commit whose verification did not pass", ungated, convention.ReleaseWorkflow, gate.Verifying)
	}
}

func TestAReleaseCarriesWhatIsNeededToVerifyIt(t *testing.T) {
	carried, err := convention.WhatAReleaseCarriesForVerification(source.Root(t))
	if err != nil {
		t.Fatalf("reading what a release carries: %v", err)
	}

	if carried.Checksums == "" {
		t.Errorf("%s writes no checksum file, so nobody can tell an intact download from a truncated one", convention.ReleaseConfigurationFile)
	}
	if len(carried.Archives) == 0 {
		t.Fatalf("%s archives nothing, so this check has nothing to hold and would pass forever", convention.ReleaseConfigurationFile)
	}
	if !slices.Contains(carried.BillsOver, "archive") {
		t.Errorf("%s draws up a bill of materials over %v and not over the archives, and what is inside a binary that embeds fonts and a client bundle is not obvious from the outside", convention.ReleaseConfigurationFile, carried.BillsOver)
	}

	signed := false
	for _, signature := range carried.Signatures {
		if signature.Over != "checksum" {
			continue
		}
		signed = true
		if !signature.Keyless {
			t.Errorf("the signature %s makes over the checksums names a key, and a release depends on no secret somebody can lose or be asked for", signature.Command)
		}
	}
	if !signed {
		t.Errorf("%s signs %v and not the checksums, so nothing says the checksums are the ones this project published", convention.ReleaseConfigurationFile, carried.Signatures)
	}

	if len(carried.Provenance) == 0 {
		t.Errorf("no job in %s runs %s, so nothing a compromised publisher cannot forge says which workflow from which commit produced the archives", convention.ReleaseWorkflow, convention.ProvenanceAction)
	}
	for _, subject := range carried.Provenance {
		if !strings.Contains(subject, "tar.gz") {
			t.Errorf("%s attests %q, which is not the archives a person downloads", convention.ReleaseWorkflow, subject)
		}
	}
}

func TestNothingInTheReleasePathReachesTheRealReviewerCLI(t *testing.T) {
	root := source.Root(t)

	reaching, err := convention.WhatTheReleaseReaches(root)
	if err != nil {
		t.Fatalf("reading what a release reaches: %v", err)
	}
	for _, reached := range reaching {
		t.Errorf("%s, and a release spends on the vendor no differently from a verification run", reached)
	}

	reach, err := convention.WorkflowsReachingTheRealCLI(root)
	if err != nil {
		t.Fatalf("reading what the workflows reach: %v", err)
	}
	if !slices.Contains(reach.Workflows, convention.ReleaseWorkflow) {
		t.Fatalf("%s was not among the workflows read, so this check says nothing about the release path", convention.ReleaseWorkflow)
	}
	for _, reached := range reach.Reaching {
		if strings.HasPrefix(reached, convention.ReleaseWorkflow) {
			t.Errorf("%s, and the boundary the verification run relies on is the one the release inherits", reached)
		}
	}
}

func TestAReleaseCompilingATagThatReachesTheRealCLIIsReported(t *testing.T) {
	root := t.TempDir()
	writeNested(t, filepath.Join(root, convention.ReleaseConfigurationFile),
		"builds:\n  - id: one\n    flags:\n      - -tags livesmoke\n")

	reaching, err := convention.WhatTheReleaseReaches(root)
	if err != nil {
		t.Fatalf("reading a fabricated release: %v", err)
	}

	want := []string{convention.ReleaseConfigurationFile + " compiles the livesmoke tag, which is behind the paths that reach the real CLI"}
	if diff := cmp.Diff(want, reaching); diff != "" {
		t.Errorf("what a release compiling the live smoke tag reports differs (-want +got):\n%s", diff)
	}
}

func TestAReleaseIsRecordedInTheChangelog(t *testing.T) {
	recorded, err := convention.ReleasesRecorded(source.Root(t))
	if err != nil {
		t.Fatalf("reading %s: %v", convention.ChangelogFile, err)
	}
	if len(recorded) == 0 {
		t.Fatalf("%s records no release, so somebody deciding whether to upgrade has nothing to read", convention.ChangelogFile)
	}

	for _, one := range recorded {
		if !aSemanticVersion.MatchString(one.Version) {
			t.Errorf("%s records %q, which is no semantic version", convention.ChangelogFile, one.Version)
		}
		if one.Dated == "" {
			t.Errorf("%s records %s and says when it was published nowhere", convention.ChangelogFile, one.Version)
		}
		if one.Says == "" {
			t.Errorf("%s names %s and says nothing about what changed in it", convention.ChangelogFile, one.Version)
		}
	}
}

func TestAReleaseRecordedWithNothingUnderItIsRead(t *testing.T) {
	root := t.TempDir()
	writeNested(t, filepath.Join(root, convention.ChangelogFile),
		"# Changelog\n\n## ["+convention.Unreleased+"]\n\n## [0.2.0] - 2026-09-01\n\nIt does the thing now.\n\n## [0.1.0] - 2026-08-31\n")

	recorded, err := convention.ReleasesRecorded(root)
	if err != nil {
		t.Fatalf("reading a fabricated changelog: %v", err)
	}

	want := []convention.RecordedRelease{
		{Version: "0.2.0", Dated: "2026-09-01", Says: "It does the thing now."},
		{Version: "0.1.0", Dated: "2026-08-31"},
	}
	if diff := cmp.Diff(want, recorded); diff != "" {
		t.Errorf("what a changelog records differs (-want +got):\n%s", diff)
	}
}

func TestAPublishingJobThatSkipsVerificationIsReported(t *testing.T) {
	root := t.TempDir()
	writeWorkflow(t, root, "verification.yml", "jobs:\n  suite:\n    steps:\n      - run: make verify\n")
	writeWorkflow(t, root, convention.ReleaseWorkflow,
		"jobs:\n  verification:\n    uses: ./.github/workflows/verification.yml\n  publication:\n    steps:\n      - run: make release\n")

	gate, err := convention.HowTheReleaseIsGated(root)
	if err != nil {
		t.Fatalf("reading a fabricated release: %v", err)
	}
	if gate.Verifying != "verification" {
		t.Fatalf("the job running the suite was read as %q", gate.Verifying)
	}

	want := []string{"publication"}
	if diff := cmp.Diff(want, gate.Ungated); diff != "" {
		t.Errorf("what a publishing job needing nothing reports differs (-want +got):\n%s", diff)
	}
}

func TestAPublishingJobBehindVerificationIsNotReported(t *testing.T) {
	root := t.TempDir()
	writeWorkflow(t, root, "verification.yml", "jobs:\n  suite:\n    steps:\n      - run: make verify\n")
	writeWorkflow(t, root, convention.ReleaseWorkflow,
		"jobs:\n  verification:\n    uses: ./.github/workflows/verification.yml\n  publication:\n    needs: verification\n    steps:\n      - run: make release\n")

	gate, err := convention.HowTheReleaseIsGated(root)
	if err != nil {
		t.Fatalf("reading a fabricated release: %v", err)
	}
	if len(gate.Ungated) != 0 {
		t.Errorf("a publishing job that needs the verification was reported as ungated: %v", gate.Ungated)
	}
}

func TestAReleaseStampsTheModulePathTheModuleDeclares(t *testing.T) {
	root := source.Root(t)

	declared, err := convention.ModulePath(root)
	if err != nil {
		t.Fatalf("reading the module path: %v", err)
	}

	building, err := convention.AReleaseBuilds(root)
	if err != nil {
		t.Fatalf("reading what a release builds: %v", err)
	}

	stamped := 0
	for _, one := range building {
		for _, named := range convention.SymbolsAReleaseStamps(one.Stamps) {
			if !strings.Contains(named, "/") {
				continue
			}
			stamped++
			if !strings.HasPrefix(named, declared+"/") {
				t.Errorf("%s stamps %s building %s, and go.mod declares %s, so the linker writes to a symbol the compiler resolves under another path", convention.ReleaseConfigurationFile, named, one.Binary, declared)
			}
		}
	}
	if stamped == 0 {
		t.Fatalf("%s stamps no symbol under a module path, so this check has nothing to compare and would pass forever", convention.ReleaseConfigurationFile)
	}
}

func TestTheSymbolsAReleaseStampsAreReadFromItsLinkerFlags(t *testing.T) {
	named := convention.SymbolsAReleaseStamps([]string{
		"-X " + hostedBy("github.com", "an-impostor", convention.Suite) + "/internal/version.version={{ .Tag }}",
		"-X main.commit={{ .ShortCommit }}",
		"-trimpath",
	})

	want := []string{hostedBy("github.com", "an-impostor", convention.Suite) + "/internal/version", "main"}
	if diff := cmp.Diff(want, named); diff != "" {
		t.Errorf("the symbols read from a release's linker flags differ (-want +got):\n%s", diff)
	}
}

func TestTheCodeIsScannedOnRequestAndBeforeARelease(t *testing.T) {
	root := source.Root(t)

	scanning, err := convention.WorkflowsChecking(root, convention.TheCodeScan)
	if err != nil {
		t.Fatalf("reading what the workflows check: %v", err)
	}
	if len(scanning) == 0 {
		t.Fatalf("no workflow runs %s, so nothing analyses the code and the security tab stays empty", convention.TheCodeScan)
	}
	for _, named := range scanning {
		start, err := convention.WhatStartsAWorkflow(root, named)
		if err != nil {
			t.Fatalf("reading what starts %s: %v", named, err)
		}
		if !start.OnRequest {
			t.Errorf("%s runs %s and nobody can ask for one: state %s", named, convention.TheCodeScan, convention.OnRequest)
		}
		for _, unasked := range start.Unasked {
			t.Errorf("%s starts on %s, and a run happens only when somebody asks for one or a release gates on one", named, unasked)
		}
	}

	scan, err := convention.HowTheReleaseIsScanned(root)
	if err != nil {
		t.Fatalf("reading how the release is scanned: %v", err)
	}
	if scan.Scanning == "" {
		t.Fatalf("no job in %s runs a workflow that runs %s, so a release publishes code nobody analysed", convention.ReleaseWorkflow, convention.TheCodeScan)
	}
	if len(scan.Unscanned) != 0 {
		t.Errorf("%v in %s publish without needing %s", scan.Unscanned, convention.ReleaseWorkflow, scan.Scanning)
	}
}

func TestAPublishingJobThatSkipsTheScanIsReported(t *testing.T) {
	root := t.TempDir()
	writeWorkflow(t, root, "verification.yml", "jobs:\n  suite:\n    steps:\n      - run: make verify\n")
	writeWorkflow(t, root, "scan.yml", "jobs:\n  analysis:\n    steps:\n      - uses: github/codeql-action/analyze@0000000000000000000000000000000000000000\n")
	writeWorkflow(t, root, convention.ReleaseWorkflow,
		"jobs:\n  verification:\n    uses: ./.github/workflows/verification.yml\n  scanning:\n    needs: verification\n    uses: ./.github/workflows/scan.yml\n  publication:\n    needs: verification\n    steps:\n      - run: make release\n")

	scan, err := convention.HowTheReleaseIsScanned(root)
	if err != nil {
		t.Fatalf("reading a fabricated release: %v", err)
	}
	if scan.Scanning != "scanning" {
		t.Fatalf("the job running the scan was read as %q", scan.Scanning)
	}

	want := []string{"publication"}
	if diff := cmp.Diff(want, scan.Unscanned); diff != "" {
		t.Errorf("what a publishing job skipping the scan reports differs (-want +got):\n%s", diff)
	}
}

func TestAPublishingJobBehindTheScanIsNotReported(t *testing.T) {
	root := t.TempDir()
	writeWorkflow(t, root, "verification.yml", "jobs:\n  suite:\n    steps:\n      - run: make verify\n")
	writeWorkflow(t, root, "scan.yml", "jobs:\n  analysis:\n    steps:\n      - uses: github/codeql-action/analyze@0000000000000000000000000000000000000000\n")
	writeWorkflow(t, root, convention.ReleaseWorkflow,
		"jobs:\n  verification:\n    uses: ./.github/workflows/verification.yml\n  scanning:\n    needs: verification\n    uses: ./.github/workflows/scan.yml\n  publication:\n    needs: [verification, scanning]\n    steps:\n      - run: make release\n")

	scan, err := convention.HowTheReleaseIsScanned(root)
	if err != nil {
		t.Fatalf("reading a fabricated release: %v", err)
	}
	if len(scan.Unscanned) != 0 {
		t.Errorf("a publishing job that needs the scan was reported as unscanned: %v", scan.Unscanned)
	}
}

func TestAReleaseCarriesTheInstallerAndListsItInTheChecksums(t *testing.T) {
	root := source.Root(t)

	attaching, err := convention.WhatAReleaseAttaches(root)
	if err != nil {
		t.Fatalf("reading what a release attaches: %v", err)
	}
	if !slices.Contains(attaching.Attached, convention.InstallerFile) {
		t.Errorf("%s attaches %v, and not %s, so the one command a person installs with downloads nothing", convention.ReleaseConfigurationFile, attaching.Attached, convention.InstallerFile)
	}
	if !slices.Contains(attaching.Checksummed, convention.InstallerFile) {
		t.Errorf("%s lists %v beside the archives in the checksums, and not %s, so the signature does not cover it", convention.ReleaseConfigurationFile, attaching.Checksummed, convention.InstallerFile)
	}
	if !convention.Carries(root, convention.InstallerFile) {
		t.Errorf("the repository holds no %s for a release to attach", convention.InstallerFile)
	}
}

func TestAReleaseThatAttachesNoInstallerIsRead(t *testing.T) {
	root := t.TempDir()
	writePackaged(t, filepath.Join(root, convention.ReleaseConfigurationFile), "checksum:\n  name_template: checksums.txt\nrelease:\n  extra_files:\n    - glob: ./other.sh\n")

	attaching, err := convention.WhatAReleaseAttaches(root)
	if err != nil {
		t.Fatalf("reading a fabricated release: %v", err)
	}
	if diff := cmp.Diff([]string{"other.sh"}, attaching.Attached); diff != "" {
		t.Errorf("what a fabricated release attaches was read otherwise (-want +got):\n%s", diff)
	}
	if len(attaching.Checksummed) != 0 {
		t.Errorf("a release listing nothing extra in its checksums was read as listing %v", attaching.Checksummed)
	}
}

func TestTheOneLineTheReadmeInstallsWithDownloadsWhatAReleaseAttaches(t *testing.T) {
	root := source.Root(t)

	downloaded, err := convention.AssetsTheFrontDocumentDownloads(root)
	if err != nil {
		t.Fatalf("reading what %s downloads: %v", convention.FrontDocumentFile, err)
	}
	if !slices.Contains(downloaded, convention.InstallerFile) {
		t.Errorf("%s downloads %v from the latest release, and not %s, so a person installing reads no one command", convention.FrontDocumentFile, downloaded, convention.InstallerFile)
	}

	attaching, err := convention.WhatAReleaseAttaches(root)
	if err != nil {
		t.Fatalf("reading what a release attaches: %v", err)
	}
	for _, asset := range downloaded {
		if !slices.Contains(attaching.Attached, asset) {
			t.Errorf("%s downloads %s from the latest release, and %s attaches only %v", convention.FrontDocumentFile, asset, convention.ReleaseConfigurationFile, attaching.Attached)
		}
	}
}

func TestAnAssetTheReadmeDownloadsIsReadFromTheLatestRelease(t *testing.T) {
	root := t.TempDir()
	writePackaged(t, filepath.Join(root, convention.FrontDocumentFile), "# x\n\n```sh\ncurl -fsSL https://github.com/socrates-tools/socrates/releases/latest/download/setup.sh | sh\n```\n")

	downloaded, err := convention.AssetsTheFrontDocumentDownloads(root)
	if err != nil {
		t.Fatalf("reading a fabricated front document: %v", err)
	}
	if diff := cmp.Diff([]string{"setup.sh"}, downloaded); diff != "" {
		t.Errorf("what a fabricated front document downloads was read otherwise (-want +got):\n%s", diff)
	}
}
