package agent_test

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/socrates-devtool/socrates/internal/agent"
	"github.com/socrates-devtool/socrates/internal/testing/source"
)

var simulatorTrees = []string{
	filepath.Join("cmd", "claudesim"),
	filepath.Join("internal", "claudesim"),
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
