package version_test

import (
	"strings"
	"testing"

	"github.com/socrates-tools/socrates-review/internal/version"
)

const (
	recordedCommit  = "7f3a9c1"
	recordedVersion = "v1.4.2"
)

func TestEveryRecordedBuildHasItsOwnLine(t *testing.T) {
	for _, described := range []struct {
		named           string
		build           version.Build
		namesTheCommit  bool
		namesTheVersion bool
		saysUncommitted bool
	}{
		{
			named:          "a build that recorded a commit",
			build:          version.Build{Commit: recordedCommit},
			namesTheCommit: true,
		},
		{
			named:           "a build that recorded a version names it with the commit",
			build:           version.Build{Version: recordedVersion, Commit: recordedCommit},
			namesTheCommit:  true,
			namesTheVersion: true,
		},
		{
			named:           "a build from a tree with uncommitted changes says so",
			build:           version.Build{Commit: recordedCommit, Uncommitted: true},
			namesTheCommit:  true,
			saysUncommitted: true,
		},
		{
			named:           "a released build from a tree with uncommitted changes says both",
			build:           version.Build{Version: recordedVersion, Commit: recordedCommit, Uncommitted: true},
			namesTheCommit:  true,
			namesTheVersion: true,
			saysUncommitted: true,
		},
		{
			named: "a build that recorded no commit says so",
			build: version.Build{},
		},
		{
			named: "a build that recorded no commit but was dirty says the same",
			build: version.Build{Uncommitted: true},
		},
		{
			named: "a build that recorded a version but no commit says what to build it with",
			build: version.Build{Version: recordedVersion},
		},
	} {
		t.Run(described.named, func(t *testing.T) {
			line := version.Describe(described.build)

			if strings.TrimSpace(line) == "" {
				t.Fatalf("%s is described by nothing at all", described.named)
			}
			if !strings.Contains(line, "socrates") {
				t.Errorf("%s is described without naming the binary: %q", described.named, line)
			}
			if strings.Count(line, "\n") != 0 {
				t.Errorf("%s is described in more than one line: %q", described.named, line)
			}

			if names := strings.Contains(line, recordedCommit); names != described.namesTheCommit {
				t.Errorf("%s names the commit: %t, want %t: %q", described.named, names, described.namesTheCommit, line)
			}
			if names := strings.Contains(line, recordedVersion); names != described.namesTheVersion {
				t.Errorf("%s names the version: %t, want %t: %q", described.named, names, described.namesTheVersion, line)
			}
			if says := strings.Contains(line, "uncommitted"); says != described.saysUncommitted {
				t.Errorf("%s says the tree held uncommitted changes: %t, want %t: %q", described.named, says, described.saysUncommitted, line)
			}
		})
	}
}

func TestAnUnstampedBuildSaysWhatToBuildItWith(t *testing.T) {
	line := version.Describe(version.Build{})

	if !strings.Contains(line, "make install") {
		t.Errorf("a build that recorded no commit does not say what to build it with: %q", line)
	}
	for _, inferred := range []string{"unknown", "devel", "0.0.0", "HEAD"} {
		if strings.Contains(line, inferred) {
			t.Errorf("a build that recorded no commit reports %q as its version: %q", inferred, line)
		}
	}
}

func TestABuildBetweenReleasesBorrowsNoTag(t *testing.T) {
	line := version.Describe(version.Build{Commit: recordedCommit})

	if !strings.Contains(line, "development build") {
		t.Errorf("a build that recorded no version does not say it is a development build: %q", line)
	}
	for _, borrowed := range []string{"v0.", "v1.", recordedVersion} {
		if strings.Contains(line, borrowed) {
			t.Errorf("a build that recorded no version reports %q as one: %q", borrowed, line)
		}
	}
}

func TestOneCommitIsDescribedOneWay(t *testing.T) {
	first := version.Describe(version.Build{Commit: recordedCommit})
	if again := version.Describe(version.Build{Commit: recordedCommit}); again != first {
		t.Errorf("the same build is described as %q and as %q", first, again)
	}
	if dirty := version.Describe(version.Build{Commit: recordedCommit, Uncommitted: true}); dirty == first {
		t.Errorf("a dirty build and a clean one are described identically: %q", first)
	}
}
