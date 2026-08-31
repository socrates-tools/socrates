package surface_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-devtool/socrates/internal/refusal"
	"github.com/socrates-devtool/socrates/internal/state"
	"github.com/socrates-devtool/socrates/internal/surface"
)

func TestTheSurfacesStateSitsUnderTheDirectoryItIsGiven(t *testing.T) {
	home := t.TempDir()
	t.Setenv(state.HomeVariable, home)

	stated, err := state.Surface()
	if err != nil {
		t.Fatalf("resolving where the surface writes: %v", err)
	}

	record, err := surface.RecordPath()
	if err != nil {
		t.Fatalf("resolving where agora writes what it is reached at: %v", err)
	}
	if filepath.Dir(record) != stated {
		t.Errorf("the surface writes what it is reached at to %s, outside the directory it was given at %s", record, stated)
	}

	page, err := surface.PagePath("a1b2c3d4")
	if err != nil {
		t.Fatalf("resolving where the surface writes a session: %v", err)
	}
	if want := filepath.Join(stated, "sessions", "a1b2c3d4", "session.json"); page != want {
		t.Errorf("the surface writes a session at %s, and the directory it was given is %s", page, want)
	}
	if !strings.HasPrefix(stated, home+string(filepath.Separator)) {
		t.Errorf("the surface writes under %s, which is not under socrates' home at %s", stated, home)
	}
}

func TestARelativeHomeIsRefusedForWhatTheSurfaceWrites(t *testing.T) {
	t.Setenv(state.HomeVariable, filepath.Join("somewhere", "below"))

	for named, resolve := range map[string]func() (string, error){
		"what it is reached at": surface.RecordPath,
		"a session":             func() (string, error) { return surface.PagePath("a1b2c3d4") },
		"a served document":     func() (string, error) { return surface.DocumentPath("a1b2c3d4") },
	} {
		path, err := resolve()
		if !refusal.Is(err) {
			t.Errorf("agora resolved %s to %q with a relative home, want a refusal: %v", named, path, err)
			continue
		}
		if !strings.Contains(err.Error(), state.HomeVariable) || !strings.Contains(err.Error(), "somewhere") {
			t.Errorf("the refusal reads %q, without naming the variable and the path", err)
		}
	}
}
