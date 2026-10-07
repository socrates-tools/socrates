package state_test

import (
	"path/filepath"
	"testing"

	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/socrates-tools/socrates/internal/state"
)

func TestTheToolWritesUnderItsOwnDirectoryInTheBrandsRoot(t *testing.T) {
	own := t.TempDir()
	t.Setenv(state.HomeVariable, "")
	t.Setenv("HOME", own)
	t.Setenv("USERPROFILE", own)

	home, err := state.Review.Home()
	if err != nil {
		t.Fatalf("resolving the tool's home with no root stated: %v", err)
	}

	root := filepath.Join(own, ".socrates")
	if want := filepath.Join(root, "review"); home != want {
		t.Errorf("the tool's home is %s, and the brand's root holding a directory named for the tool is %s", home, want)
	}

	configuration, err := state.Review.Configuration()
	if err != nil {
		t.Fatalf("resolving the configuration file: %v", err)
	}
	if beside := filepath.Join(root, state.ConfigurationFile); configuration == beside {
		t.Errorf("the configuration file is %s, which is directly in the brand's root rather than in the tool's directory under it", configuration)
	}

	runs, err := state.Runs()
	if err != nil {
		t.Fatalf("resolving where the runs go: %v", err)
	}
	for _, under := range []string{configuration, runs} {
		if filepath.Dir(under) != home {
			t.Errorf("%s is not in the tool's directory %s, so removing the tool would not be removing one directory", under, home)
		}
	}
}

func TestEachToolReadsAConfigurationFileInItsOwnDirectory(t *testing.T) {
	stated := t.TempDir()
	t.Setenv(state.HomeVariable, stated)

	for tool, want := range map[state.Tool]string{
		state.Review: filepath.Join(stated, "review", state.ConfigurationFile),
		state.Agora:  filepath.Join(stated, "agora", state.ConfigurationFile),
	} {
		configuration, err := tool.Configuration()
		if err != nil {
			t.Fatalf("resolving the configuration file %s reads: %v", tool, err)
		}
		if configuration != want {
			t.Errorf("%s reads %s, and the file in its own directory under the brand's root is %s", tool, configuration, want)
		}
	}
}

func TestARelativeHomeIsRefusedForAgorasConfigurationToo(t *testing.T) {
	t.Setenv(state.HomeVariable, filepath.Join("somewhere", "below"))

	configuration, err := state.Agora.Configuration()
	if !refusal.Is(err) {
		t.Fatalf("a relative home resolved agora's configuration to %q, want a refusal: %v", configuration, err)
	}
}

func TestImplementWritesUnderItsOwnDirectoryInTheBrandsRoot(t *testing.T) {
	stated := t.TempDir()
	t.Setenv(state.HomeVariable, stated)

	configuration, err := state.Implement.Configuration()
	if err != nil {
		t.Fatalf("resolving implement's configuration file: %v", err)
	}
	if want := filepath.Join(stated, "implement", state.ConfigurationFile); configuration != want {
		t.Errorf("implement's configuration is %s, want %s", configuration, want)
	}
}
