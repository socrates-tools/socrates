package state_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/rui-ktei/socrates/internal/refusal"
	"github.com/rui-ktei/socrates/internal/state"
)

func TestTheSurfacesDirectoryIsUnderTheStatedHome(t *testing.T) {
	stated := t.TempDir()
	t.Setenv(state.HomeVariable, stated)

	directory, err := state.Surface()
	if err != nil {
		t.Fatalf("resolving where the surface keeps what it writes: %v", err)
	}
	if want := filepath.Join(stated, "surface"); directory != want {
		t.Errorf("the surface keeps what it writes at %s, and %s states the home is %s", directory, state.HomeVariable, want)
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
