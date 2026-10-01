package config_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/implement/config"
)

func TestTheFileIsNamedBeforeWhatIsInForce(t *testing.T) {
	path := filepath.Join("somewhere", "config.yaml")

	lines := config.ReportAt(path, config.Config{}).Lines()
	if lines[0] != path {
		t.Errorf("the first line reads %q, and it is the file socrates implement reads", lines[0])
	}
}

func TestEverySettingGetsOneLineOfItsOwnEndingInWhereItCameFrom(t *testing.T) {
	loaded, err := config.LoadFrom(write(t, "concurrency: 2\n"))
	if err != nil {
		t.Fatalf("loading: %v", err)
	}

	var keys, sources []string
	for _, line := range config.ReportAt("config.yaml", loaded).Lines()[1:] {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		keys = append(keys, fields[0])
		sources = append(sources, fields[len(fields)-1])
	}
	if diff := cmp.Diff([]string{"concurrency", "roles.implementer.model", "roles.implementer.effort"}, keys); diff != "" {
		t.Errorf("the keys reported differ from the settings (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]string{"set", "default", "default"}, sources); diff != "" {
		t.Errorf("where each setting came from differs (-want +got):\n%s", diff)
	}
}

func TestTheColumnsLineUp(t *testing.T) {
	lines := config.ReportAt("config.yaml", config.Config{}).Lines()[2:]

	value := strings.Index(lines[0], "4")
	for _, line := range lines {
		fields := strings.Fields(line)
		if at := strings.Index(line, fields[1]); at != value {
			t.Errorf("the line %q starts its value at %d, and the first line started at %d", line, at, value)
		}
	}
}
