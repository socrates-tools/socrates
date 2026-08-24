package convention

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"
)

const buildConstraint = "//go:build"

var networkImports = []string{"net", "crypto/tls", "golang.org/x/net"}

var BuildConstrainedFiles = map[string]string{
	filepath.Join("internal", "agent", "deadline_test.go"):     "a deadline exercised with testing/synctest, which needs the experiment enabled",
	filepath.Join("internal", "fidelity", "drift_test.go"):     "comparing the simulator against recordings of the real CLI",
	filepath.Join("internal", "fidelity", "livesmoke_test.go"): "one real review against the real CLI",
	filepath.Join("internal", "fidelity", "record_test.go"):    "recording what the real CLI answers",
	filepath.Join("internal", "process", "process_other.go"):   "the process group on a platform that has no process groups",
	filepath.Join("internal", "process", "process_unix.go"):    "the process group on a platform that has them",
}

type Reach struct {
	Reaching   []string
	Tagged     []string
	Undeclared []string
	Stale      []string
}

func NetworkReach(root string, declared map[string]string) (Reach, error) {
	var reach Reach

	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if path != root && skipped(entry.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		relative, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return relErr
		}

		constrained, imported, readErr := readReach(path, relative)
		if readErr != nil {
			return readErr
		}

		if constrained {
			reach.Tagged = append(reach.Tagged, relative)
			if _, stated := declared[relative]; !stated {
				reach.Undeclared = append(reach.Undeclared, relative)
			}
		}
		for _, one := range imported {
			if reachesTheNetwork(one) && !constrained {
				reach.Reaching = append(reach.Reaching, relative+" imports "+one)
			}
		}
		return nil
	})
	if err != nil {
		return Reach{}, err
	}

	for stated := range declared {
		if !slices.Contains(reach.Tagged, stated) {
			reach.Stale = append(reach.Stale, stated)
		}
	}

	slices.Sort(reach.Reaching)
	slices.Sort(reach.Tagged)
	slices.Sort(reach.Undeclared)
	slices.Sort(reach.Stale)
	return reach, nil
}

func readReach(path, relative string) (bool, []string, error) {
	set := token.NewFileSet()
	parsed, err := parser.ParseFile(set, path, nil, parser.ImportsOnly|parser.ParseComments|parser.SkipObjectResolution)
	if err != nil {
		return false, nil, fmt.Errorf("parsing %s: %w", relative, err)
	}

	constrained := false
	for _, group := range parsed.Comments {
		for _, line := range group.List {
			if strings.HasPrefix(line.Text, buildConstraint) {
				constrained = true
			}
		}
	}

	var imported []string
	for _, one := range parsed.Imports {
		text, quoted := stringLiteral(one.Path)
		if quoted {
			imported = append(imported, text)
		}
	}
	return constrained, imported, nil
}

func reachesTheNetwork(imported string) bool {
	for _, reaching := range networkImports {
		if imported == reaching || strings.HasPrefix(imported, reaching+"/") {
			return true
		}
	}
	return false
}
