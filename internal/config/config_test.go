package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/rui-ktei/socrates/internal/agent"
	"github.com/rui-ktei/socrates/internal/config"
	"github.com/rui-ktei/socrates/internal/refusal"
	"github.com/rui-ktei/socrates/internal/state"
)

func TestNoConfigurationFile(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, "socrates", "config.yaml")

	loaded, err := config.LoadFrom(path)
	if err != nil {
		t.Fatalf("loading an absent configuration file: %v", err)
	}
	if diff := cmp.Diff(config.Config{}, loaded); diff != "" {
		t.Errorf("an absent configuration file is not the defaults (-want +got):\n%s", diff)
	}

	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatalf("reading %s: %v", directory, err)
	}
	if len(entries) != 0 {
		t.Errorf("loading an absent configuration file created %v", entries)
	}
}

func TestAnEmptyFileIsTheDefaults(t *testing.T) {
	loaded, err := config.LoadFrom(write(t, ""))
	if err != nil {
		t.Fatalf("loading an empty configuration file: %v", err)
	}
	if diff := cmp.Diff(config.Config{}, loaded); diff != "" {
		t.Errorf("an empty configuration file is not the defaults (-want +got):\n%s", diff)
	}
}

func TestEverySettingIsRead(t *testing.T) {
	path := write(t, `roles:
  reviewer:
    model: sonnet
    effort: max
diff_ceiling: 4096
agent_deadline: 45m
log_ceiling: 8192
run_retention: 72h
surface:
  bind: 127.0.0.1:7788
  reached_at: https://agora.example:8443/
`)

	loaded, err := config.LoadFrom(path)
	if err != nil {
		t.Fatalf("loading %s: %v", path, err)
	}

	want := config.Config{
		Roles:         agent.Configuration{agent.Reviewer: {Model: "sonnet", Effort: "max"}},
		DiffCeiling:   4096,
		AgentDeadline: 45 * time.Minute,
		LogCeiling:    8192,
		RunRetention:  72 * time.Hour,
		Surface:       config.Surface{Bind: "127.0.0.1:7788", ReachedAt: "https://agora.example:8443/"},
	}
	if diff := cmp.Diff(want, loaded); diff != "" {
		t.Errorf("the configuration differs (-want +got):\n%s", diff)
	}
}

func TestAnUnreadableConfigurationFileIsARefusal(t *testing.T) {
	path := write(t, "roles: [this is not a mapping\n")

	_, err := config.LoadFrom(path)
	if err == nil {
		t.Fatal("a configuration file that does not parse was accepted")
	}
	if !refusal.Is(err) {
		t.Errorf("a configuration file that does not parse gave %T, want a refusal", err)
	}
	if !strings.Contains(err.Error(), path) {
		t.Errorf("the refusal does not name the file: %s", err)
	}
}

func TestASettingSocratesDoesNotHaveIsARefusal(t *testing.T) {
	path := write(t, "reviewer_temperature: 0.7\n")

	_, err := config.LoadFrom(path)
	if !refusal.Is(err) {
		t.Fatalf("a configuration file naming a setting socrates does not have gave %v, want a refusal", err)
	}
	if !strings.Contains(err.Error(), "reviewer_temperature") {
		t.Errorf("the refusal does not name the setting: %s", err)
	}
}

func TestASurfaceBlockThatNamesAFieldSocratesDoesNotKnowIsRefused(t *testing.T) {
	path := write(t, "surface:\n  proxy: something.example\n")

	_, err := config.LoadFrom(path)
	if !refusal.Is(err) {
		t.Fatalf("a surface block naming a field socrates does not know gave %v, want a refusal", err)
	}
	if !strings.Contains(err.Error(), "proxy") {
		t.Errorf("the refusal does not name the field: %s", err)
	}
}

func TestWhatTheSurfaceOnceReadOnItsOwnIsRefusedOutsideTheBlock(t *testing.T) {
	for _, stated := range []string{"bind: 127.0.0.1:0\n", "reached_at: https://agora.example:8443/\n"} {
		named, _, _ := strings.Cut(stated, ":")

		_, err := config.LoadFrom(write(t, stated))
		if !refusal.Is(err) {
			t.Fatalf("a configuration file stating %s at the top gave %v, want a refusal", named, err)
		}
		if !strings.Contains(err.Error(), named) {
			t.Errorf("the refusal does not name the field: %s", err)
		}
	}
}

func TestARoleSocratesDoesNotHaveIsARefusal(t *testing.T) {
	path := write(t, "roles:\n  judge:\n    model: opus\n")

	_, err := config.LoadFrom(path)
	if !refusal.Is(err) {
		t.Fatalf("a configuration file naming a role socrates does not have gave %v, want a refusal", err)
	}
	if !strings.Contains(err.Error(), "judge") {
		t.Errorf("the refusal does not name the role: %s", err)
	}
}

func TestANegativeMeasureIsARefusal(t *testing.T) {
	path := write(t, "diff_ceiling: -1\n")

	_, err := config.LoadFrom(path)
	if !refusal.Is(err) {
		t.Fatalf("a negative ceiling gave %v, want a refusal", err)
	}
	if !strings.Contains(err.Error(), "diff_ceiling") {
		t.Errorf("the refusal does not name the setting: %s", err)
	}
}

func TestNoSettingChangesIt(t *testing.T) {
	for _, stated := range []string{"round_cap: 3\n", "rounds: 3\n", "cap: 3\n", "max_rounds: 3\n"} {
		named, _, _ := strings.Cut(stated, ":")

		_, err := config.LoadFrom(write(t, stated))
		if !refusal.Is(err) {
			t.Fatalf("a configuration file stating %s gave %v, want a refusal", named, err)
		}
		if !strings.Contains(err.Error(), named) {
			t.Errorf("the refusal does not name the setting: %s", err)
		}
	}
}

func write(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
	return path
}

func TestTheConfigurationSitsAtTheTopOfTheHomeOnEveryPlatform(t *testing.T) {
	stated := t.TempDir()
	t.Setenv(state.HomeVariable, stated)

	path, err := config.Path()
	if err != nil {
		t.Fatalf("resolving the configuration path: %v", err)
	}
	if want := filepath.Join(stated, "config.yaml"); path != want {
		t.Errorf("socrates reads %s, and its home is %s", path, want)
	}
}

func TestTheConfigurationIsResolvedByTheOneRuleThePackageDoesNotOwn(t *testing.T) {
	stated := t.TempDir()
	t.Setenv(state.HomeVariable, stated)

	path, err := config.Path()
	if err != nil {
		t.Fatalf("resolving the configuration path: %v", err)
	}
	home, err := state.Home()
	if err != nil {
		t.Fatalf("resolving socrates' home: %v", err)
	}
	if directory := filepath.Dir(path); directory != home {
		t.Errorf("socrates reads its configuration from %s, and its home is %s, so the package resolves a directory of its own", directory, home)
	}
}
