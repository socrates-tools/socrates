package state_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/socrates-tools/socrates/internal/state"
)

func TestTheSurfacesDirectoryIsUnderTheStatedHome(t *testing.T) {
	stated := t.TempDir()
	t.Setenv(state.HomeVariable, stated)

	directory, err := state.Surface()
	if err != nil {
		t.Fatalf("resolving where the surface keeps what it writes: %v", err)
	}
	if want := filepath.Join(stated, "review", "surface"); directory != want {
		t.Errorf("the surface keeps what it writes at %s, and %s states the brand's root is %s", directory, state.HomeVariable, want)
	}
}

func TestTheSurfacesDirectorySitsBesideTheRuns(t *testing.T) {
	stated := t.TempDir()
	t.Setenv(state.HomeVariable, stated)

	directory, err := state.Surface()
	if err != nil {
		t.Fatalf("resolving where the surface keeps what it writes: %v", err)
	}
	runs, err := state.Runs()
	if err != nil {
		t.Fatalf("resolving where the runs are kept: %v", err)
	}
	if filepath.Dir(directory) != filepath.Dir(runs) {
		t.Errorf("the surface writes under %s and the runs under %s, and one home holds both beside each other", directory, runs)
	}
	if strings.HasPrefix(directory, runs+string(filepath.Separator)) {
		t.Errorf("the surface writes under %s, which is inside the runs at %s", directory, runs)
	}
}

func TestARelativeHomeIsRefusedForTheSurfacesDirectory(t *testing.T) {
	t.Setenv(state.HomeVariable, filepath.Join("somewhere", "below"))

	directory, err := state.Surface()
	if !refusal.Is(err) {
		t.Fatalf("a relative home resolved the surface's directory to %q, want a refusal: %v", directory, err)
	}
	if !strings.Contains(err.Error(), state.HomeVariable) || !strings.Contains(err.Error(), "somewhere") {
		t.Errorf("the refusal reads %q, without naming the variable and the path", err)
	}
}

func TestTheToolWritesUnderItsOwnDirectoryInTheBrandsRoot(t *testing.T) {
	own := t.TempDir()
	t.Setenv(state.HomeVariable, "")
	t.Setenv("HOME", own)
	t.Setenv("USERPROFILE", own)

	home, err := state.Home()
	if err != nil {
		t.Fatalf("resolving the tool's home with no root stated: %v", err)
	}

	root := filepath.Join(own, ".socrates")
	if want := filepath.Join(root, "review"); home != want {
		t.Errorf("the tool's home is %s, and the brand's root holding a directory named for the tool is %s", home, want)
	}

	configuration, err := state.Configuration()
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
