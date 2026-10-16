package platform_test

import (
	"runtime"
	"slices"
	"testing"

	"github.com/socrates-tools/socrates/internal/convention"
	"github.com/socrates-tools/socrates/internal/platform"
	"github.com/socrates-tools/socrates/internal/process"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

func TestTheSupportedSetIsStatedOnce(t *testing.T) {
	if len(platform.Supported()) == 0 {
		t.Fatal("the supported set names no platform, so nothing that reads it is being checked against anything")
	}

	if platform.Supports(runtime.GOOS) != process.ReapsByGroup {
		t.Errorf("the supported set says %s is supported: %t, and reaping by process group is available here: %t", runtime.GOOS, platform.Supports(runtime.GOOS), process.ReapsByGroup)
	}

	named := append(platform.OperatingSystems(), platform.Names()...)
	copies, err := convention.PlatformsNamedOutsideTheStatement(source.Root(t), named)
	if err != nil {
		t.Fatalf("reading what names a platform: %v", err)
	}
	for _, copied := range copies {
		t.Errorf("%s, which is a second copy of a set stated in %s", copied, convention.StatementDirectory)
	}
}

func TestAPlatformQualifiesOnTheGuaranteesRatherThanOnCompiling(t *testing.T) {
	for _, one := range platform.Supported() {
		if !one.ReapsByProcessGroup || !one.FilesystemEnforcesTheSeal || !one.TheSuiteRunsThere {
			t.Errorf("%s is supported while %+v, and each of the three is what supported means", one.OperatingSystem, one)
		}
	}

	unproven := platform.Platform{OperatingSystem: "plan9", Name: "Plan 9", ReapsByProcessGroup: true, FilesystemEnforcesTheSeal: true}
	if unproven.Qualifies() {
		t.Error("a platform the suite has never run on qualifies, so compiling is the basis after all")
	}
}

func TestAReviewIsSupportedOnlyWhereTheSetAndTheBuildAgree(t *testing.T) {
	if platform.ReviewIsSupported() != (process.ReapsByGroup && platform.Supports(platform.TheOneRunning())) {
		t.Error("whether a review is supported is decided by something other than the stated set and what this build can reap")
	}
}

func TestEveryBuildPairsASupportedSystemWithAnArchitectureItIsStatedFor(t *testing.T) {
	builds := platform.Builds()

	for _, system := range platform.OperatingSystems() {
		var architectures []string
		for _, build := range builds {
			if build.OperatingSystem == system {
				architectures = append(architectures, build.Architecture)
			}
		}
		if len(architectures) == 0 {
			t.Errorf("%s is supported and no build is stated for it", system)
		}
	}

	for _, build := range builds {
		if !platform.Supports(build.OperatingSystem) {
			t.Errorf("a build is stated for %s, which is not supported", build.OperatingSystem)
		}
		if !slices.Contains(platform.Architectures(), build.Architecture) {
			t.Errorf("a build is stated for %s, which is not among the architectures", build.Architecture)
		}
	}
}
