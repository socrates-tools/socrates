package source

import (
	"strings"
	"testing"
)

func TestAModuleLinkedAtTwoVersionsForTwoBuildsFailsNamingBoth(t *testing.T) {
	_, err := mergeLinkedModulesOfEveryBuild(map[string][]LinkedModule{
		"linux/amd64":  {{Path: "example.test/dep", Version: "v1.0.0", Directory: "/one"}},
		"darwin/arm64": {{Path: "example.test/dep", Version: "v1.1.0", Directory: "/two"}},
	})

	if err == nil {
		t.Fatal("two versions of one module were merged into one without a word")
	}
	for _, want := range []string{"example.test/dep", "v1.0.0", "v1.1.0", "linux/amd64", "darwin/arm64"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("the failure %q does not name %q", err, want)
		}
	}
}

func TestAModuleLinkedForEveryBuildAtOneVersionIsListedOnce(t *testing.T) {
	merged, err := mergeLinkedModulesOfEveryBuild(map[string][]LinkedModule{
		"linux/amd64":  {{Path: "example.test/b", Version: "v1.0.0", Directory: "/b"}, {Path: "example.test/a", Version: "v1.0.0", Directory: "/a"}},
		"darwin/arm64": {{Path: "example.test/a", Version: "v1.0.0", Directory: "/a"}},
	})
	if err != nil {
		t.Fatalf("merging agreeing builds: %v", err)
	}

	if len(merged) != 2 || merged[0].Path != "example.test/a" || merged[1].Path != "example.test/b" {
		t.Errorf("merged = %v, want a then b once each", merged)
	}
}
