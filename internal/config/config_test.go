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

func write(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
	return path
}
