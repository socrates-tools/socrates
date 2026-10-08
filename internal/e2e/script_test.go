package e2e_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
	"github.com/socrates-tools/socrates/internal/state"
	"github.com/socrates-tools/socrates/internal/testing/source"
	"github.com/socrates-tools/socrates/review/simulation"
)

func TestTheCommandLineBehavesAsItsScriptsSay(t *testing.T) {
	binaries := t.TempDir()
	if err := os.Symlink(source.Shared(t, "./cmd/socrates"), filepath.Join(binaries, "socrates")); err != nil {
		t.Fatalf("linking socrates: %v", err)
	}
	if err := os.Symlink(source.Shared(t, "./review/cmd/claudesim"), filepath.Join(binaries, "claude")); err != nil {
		t.Fatalf("linking the simulator as claude: %v", err)
	}

	testscript.Run(t, testscript.Params{
		Dir: filepath.Join("testdata", "script"),
		Setup: func(env *testscript.Env) error {
			env.Setenv("PATH", binaries+string(os.PathListSeparator)+os.Getenv("PATH"))
			env.Setenv("HOME", env.WorkDir)
			env.Setenv("TMPDIR", env.WorkDir)
			root := filepath.Join(env.WorkDir, "home")
			env.Setenv(state.HomeVariable, root)
			env.Setenv("TOOL_HOME", state.Review.HomeUnder(root))
			env.Setenv(simulation.ScenarioVariable, filepath.Join(env.WorkDir, "clean.yaml"))
			env.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
			env.Setenv("GIT_CONFIG_SYSTEM", os.DevNull)
			env.Setenv("GIT_CONFIG_NOSYSTEM", "1")
			env.Setenv("GIT_TERMINAL_PROMPT", "0")
			env.Setenv("GIT_AUTHOR_NAME", "A Person")
			env.Setenv("GIT_AUTHOR_EMAIL", "person@example.com")
			env.Setenv("GIT_COMMITTER_NAME", "A Person")
			env.Setenv("GIT_COMMITTER_EMAIL", "person@example.com")
			return nil
		},
	})
}
