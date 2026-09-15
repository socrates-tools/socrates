package report_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/socrates-review/config"
	"github.com/socrates-tools/socrates/socrates-review/report"
)

func TestTheFileIsNamedBeforeWhatIsInForce(t *testing.T) {
	path := filepath.Join("somewhere", "config.yaml")

	lines := report.ConfigurationAt(path, config.Config{}).Lines()
	if lines[0] != path {
		t.Errorf("the first line reads %q, and it is the file socrates-review reads", lines[0])
	}
}

func TestEverySettingGetsOneLineOfItsOwn(t *testing.T) {
	loaded := config.Config{ReviewerCLI: "codex"}

	var settings []string
	for _, line := range report.ConfigurationAt("config.yaml", loaded).Lines()[1:] {
		if strings.TrimSpace(line) != "" {
			settings = append(settings, line)
		}
	}
	if len(settings) != len(loaded.Settings()) {
		t.Fatalf("%d settings are reported on %d lines:\n%s", len(loaded.Settings()), len(settings), strings.Join(settings, "\n"))
	}

	var keys []string
	for _, line := range settings {
		keys = append(keys, strings.Fields(line)[0])
	}
	var want []string
	for _, setting := range loaded.Settings() {
		want = append(want, setting.Key)
	}
	if diff := cmp.Diff(want, keys); diff != "" {
		t.Errorf("the keys reported differ from the settings (-want +got):\n%s", diff)
	}
}

func TestASettingSaysWhetherItIsSetOrADefault(t *testing.T) {
	loaded := config.Config{ReviewerCLI: "codex"}

	lines := report.ConfigurationAt("config.yaml", loaded).Lines()
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		last := fields[len(fields)-1]
		if last != string(config.Set) && last != string(config.Default) {
			t.Errorf("the line %q ends in %q, and every setting says whether it is set or a default", line, last)
		}
	}

	if !strings.Contains(lineFor(t, lines, "reviewer_cli"), "codex") {
		t.Errorf("the reviewer CLI line reads %q, and the file states codex", lineFor(t, lines, "reviewer_cli"))
	}
	if !strings.HasSuffix(strings.TrimRight(lineFor(t, lines, "reviewer_cli"), " "), string(config.Set)) {
		t.Errorf("a setting the file states does not read as set: %q", lineFor(t, lines, "reviewer_cli"))
	}
}

func TestASettingWithNoValueInForceLeavesItsValueEmpty(t *testing.T) {
	lines := report.ConfigurationAt("config.yaml", config.Config{}).Lines()

	reached := lineFor(t, lines, "surface.reached_at")
	if fields := strings.Fields(reached); len(fields) != 2 {
		t.Errorf("surface.reached_at reads %q, and nothing states one, so its value column is empty", reached)
	}
}

func TestTheColumnsLineUp(t *testing.T) {
	loaded := config.Config{ReviewerCLI: "codex"}
	lines := report.ConfigurationAt("config.yaml", loaded).Lines()

	source, value := -1, -1
	for at, setting := range loaded.Settings() {
		line := lines[at+2]

		if !strings.HasPrefix(line, setting.Key) {
			t.Errorf("the line %q does not open with %s", line, setting.Key)
		}
		startsSource := len(strings.TrimRight(line, " ")) - len(setting.Source)
		if source == -1 {
			source = startsSource
		}
		if startsSource != source {
			t.Errorf("the line %q says where its value came from at %d, and the line before said so at %d", line, startsSource, source)
		}
		if setting.Value == "" {
			continue
		}
		startsValue := strings.Index(line, setting.Value)
		if value == -1 {
			value = startsValue
		}
		if startsValue != value {
			t.Errorf("the line %q starts its value at %d, and the line before started at %d", line, startsValue, value)
		}
	}
}

func lineFor(t *testing.T, lines []string, key string) string {
	t.Helper()

	for _, line := range lines {
		if fields := strings.Fields(line); len(fields) > 0 && fields[0] == key {
			return line
		}
	}
	t.Fatalf("no line reports %s:\n%s", key, strings.Join(lines, "\n"))
	return ""
}
