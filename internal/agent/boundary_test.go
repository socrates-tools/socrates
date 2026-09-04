package agent_test

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/ktei/socrates/internal/agent"
	"github.com/ktei/socrates/internal/testing/source"
)

var simulatorTrees = []string{
	filepath.Join("cmd", "claudesim"),
	filepath.Join("cmd", "codexsim"),
	filepath.Join("cmd", "opencodesim"),
	filepath.Join("internal", "claudesim"),
	filepath.Join("internal", "codexsim"),
	filepath.Join("internal", "opencodesim"),
	filepath.Join("internal", "simulation"),
	filepath.Join("internal", "testing"),
}

func TestOnePlaceStartsTheProcess(t *testing.T) {
	root := source.Root(t)
	owner := filepath.Join("internal", "agent")

	naming := map[string]bool{}
	for _, path := range source.Files(t, root) {
		relative, err := filepath.Rel(root, path)
		if err != nil {
			t.Fatalf("placing %s: %v", path, err)
		}
		if withinSimulator(relative) {
			continue
		}

		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("reading %s: %v", path, err)
		}
		if strings.Contains(string(content), `"claude"`) {
			naming[filepath.Dir(relative)] = true
		}
	}

	if len(naming) != 1 || !naming[owner] {
		t.Errorf("the agent binary is named in %v, want only %s", keys(naming), owner)
	}

	constructor := reflect.TypeOf(agent.NewClaude)
	if answered := constructor.Out(0); answered.Kind() != reflect.Interface {
		t.Errorf("the constructor answers %s, want the interface every caller reaches it through", answered)
	}
}

func withinSimulator(relative string) bool {
	for _, tree := range simulatorTrees {
		if strings.HasPrefix(relative, tree+string(filepath.Separator)) {
			return true
		}
	}
	return false
}

func keys(named map[string]bool) []string {
	var listed []string
	for name := range named {
		listed = append(listed, name)
	}
	return listed
}

var vendorFlags = []string{
	"--output-schema",
	"--json-schema",
	"--ignore-user-config",
	"--setting-sources",
	"--safe-mode",
	"--sandbox",
	"--variant",
	"project_doc_max_bytes",
	"OPENCODE_DISABLE_PROJECT_CONFIG",
}

var vendorEventShapes = []string{
	"turn.completed",
	"turn.failed",
	"thread.started",
	"agent_message",
	"step_finish",
	"structured_output",
	"stream-json",
	"mcp_servers",
}

func TestNothingAboveTheSeamKnowsWhichVendorAnswered(t *testing.T) {
	root := source.Root(t)
	seam := filepath.Join("internal", "agent")

	naming := append(append(append([]string{}, agent.Names()...), vendorFlags...), vendorEventShapes...)

	for _, path := range source.Files(t, root) {
		relative, err := filepath.Rel(root, path)
		if err != nil {
			t.Fatalf("placing %s: %v", path, err)
		}
		if withinSimulator(relative) || strings.HasPrefix(relative, seam+string(filepath.Separator)) {
			continue
		}
		if strings.HasSuffix(relative, "_test.go") {
			continue
		}

		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("reading %s: %v", path, err)
		}
		for _, named := range naming {
			if strings.Contains(string(content), `"`+named+`"`) {
				t.Errorf("%s names %q, and what varies between vendors is confined to %s", relative, named, seam)
			}
		}
	}
}
