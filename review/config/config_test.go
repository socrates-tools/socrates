package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/socrates-tools/socrates/internal/state"
	"github.com/socrates-tools/socrates/review/agent"
	"github.com/socrates-tools/socrates/review/config"
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

func TestASurfaceBlockIsRefusedByNamingTheFileItBelongsIn(t *testing.T) {
	stated := t.TempDir()
	t.Setenv(state.HomeVariable, stated)
	agoras := filepath.Join(stated, "agora", state.ConfigurationFile)

	for _, block := range []string{
		"surface:\n  bind: 127.0.0.1:7788\n",
		"diff_ceiling: 4096\nsurface:\n  reached_at: https://agora.example:8443/\n",
		"surface:\n",
	} {
		path := write(t, block)

		_, err := config.LoadFrom(path)
		if !refusal.Is(err) {
			t.Fatalf("a configuration file stating %q gave %v, want a refusal", block, err)
		}
		for _, named := range []string{path, "surface", agoras} {
			if !strings.Contains(err.Error(), named) {
				t.Errorf("the refusal of %q does not name %s: %s", block, named, err)
			}
		}
	}
}

func TestASurfaceBlockIsRefusedFromContentAloneToo(t *testing.T) {
	t.Setenv(state.HomeVariable, t.TempDir())

	if err := config.Parse("config.yaml", []byte("surface:\n  bind: 127.0.0.1:7788\n")); !refusal.Is(err) {
		t.Errorf("parsing a surface block gave %v, want the refusal the loader gives", err)
	}
}

func TestAFileThatIsNotAMappingIsRefusedWithoutMentioningASurfaceBlock(t *testing.T) {
	t.Setenv(state.HomeVariable, t.TempDir())

	_, err := config.LoadFrom(write(t, "- surface\n"))
	if !refusal.Is(err) {
		t.Fatalf("a configuration file holding a list gave %v, want a refusal", err)
	}
	if strings.Contains(err.Error(), "agora") {
		t.Errorf("a file holding no surface block was refused as though it did: %s", err)
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
	if want := filepath.Join(stated, "review", "config.yaml"); path != want {
		t.Errorf("socrates reads %s, and the top of its home under the brand's root is %s", path, want)
	}
}

func TestTheConfigurationIsResolvedByTheOneRuleThePackageDoesNotOwn(t *testing.T) {
	stated := t.TempDir()
	t.Setenv(state.HomeVariable, stated)

	path, err := config.Path()
	if err != nil {
		t.Fatalf("resolving the configuration path: %v", err)
	}
	home, err := state.Review.Home()
	if err != nil {
		t.Fatalf("resolving socrates' home: %v", err)
	}
	if directory := filepath.Dir(path); directory != home {
		t.Errorf("socrates reads its configuration from %s, and its home is %s, so the package resolves a directory of its own", directory, home)
	}
}

func TestAnAbsentFieldRunsTheCLISocratesShippedWith(t *testing.T) {
	absent, err := config.LoadFrom(filepath.Join(t.TempDir(), "socrates", "config.yaml"))
	if err != nil {
		t.Fatalf("loading an absent configuration file: %v", err)
	}
	if absent.Vendor() != agent.DefaultVendor {
		t.Errorf("an absent configuration file resolves to %q, want %q", absent.Vendor(), agent.DefaultVendor)
	}

	stated, err := config.LoadFrom(write(t, "diff_ceiling: 1000\n"))
	if err != nil {
		t.Fatalf("loading a configuration file naming no reviewer CLI: %v", err)
	}
	if stated.Vendor() != agent.DefaultVendor {
		t.Errorf("a configuration file naming no reviewer CLI resolves to %q, want %q", stated.Vendor(), agent.DefaultVendor)
	}
}

func TestAReviewerCLIIsCarriedInOneField(t *testing.T) {
	for _, one := range agent.Supported() {
		loaded, err := config.LoadFrom(write(t, "reviewer_cli: "+one.Name+"\n"))
		if err != nil {
			t.Fatalf("loading a configuration file naming %s: %v", one.Name, err)
		}
		if loaded.Vendor() != one.Name {
			t.Errorf("a configuration file naming %s resolves to %q", one.Name, loaded.Vendor())
		}
	}
}

func TestWhatTheLoaderRefusesIsRefusedFromContentAloneToo(t *testing.T) {
	for _, unusable := range []string{
		"diff_ceiling: -1\n",
		"reviewer_temperature: 0.5\n",
		"bind: 127.0.0.1:18018\n",
		"roles: [not a mapping\n",
	} {
		path := filepath.Join(t.TempDir(), "config.yaml")

		fromContent := config.Parse(path, []byte(unusable))
		if fromContent == nil {
			t.Errorf("%q parses from content alone, and the loader refuses it", unusable)
			continue
		}
		if !refusal.Is(fromContent) {
			t.Errorf("%q is rejected as %v rather than as a refusal", unusable, fromContent)
		}
		if !strings.Contains(fromContent.Error(), path) {
			t.Errorf("the refusal of %q reads %q, without naming the file", unusable, fromContent)
		}
	}
}

func TestContentTheLoaderAcceptsParsesFromContentAlone(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")

	if err := config.Parse(path, []byte("reviewer_cli: codex\nrun_retention: 72h\n")); err != nil {
		t.Errorf("content the loader accepts was refused: %v", err)
	}
	if err := config.Parse(path, nil); err != nil {
		t.Errorf("an empty file was refused: %v", err)
	}
}
