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
	filepath.Join("internal", "fidelity", "criteria_test.go"):  "establishing by recording which reviewer CLIs meet the criteria",
	filepath.Join("internal", "fidelity", "drift_test.go"):     "comparing the simulator against recordings of the real CLI",
	filepath.Join("internal", "fidelity", "livesmoke_test.go"): "one real review against the real CLI",
	filepath.Join("internal", "fidelity", "record_test.go"):    "recording what the real CLI answers",
	filepath.Join("internal", "process", "process_other.go"):   "the process group on a platform that has no process groups",
	filepath.Join("internal", "process", "process_unix.go"):    "the process group on a platform that has them",
}

type Scope struct {
	Module string
	From   []string
}

var Socrates = Scope{
	Module: "github.com/ktei/socrates",
	From:   []string{filepath.Join("cmd", "socrates")},
}

type Reach struct {
	Measured   []string
	Reaching   []string
	Tagged     []string
	Undeclared []string
	Stale      []string
}

type sourceFile struct {
	relative    string
	directory   string
	constrained bool
	imported    []string
	test        bool
}

func NetworkReach(root string, scope Scope, declared map[string]string) (Reach, error) {
	files, err := readTree(root)
	if err != nil {
		return Reach{}, err
	}

	measured := scope.closure(files)

	var reach Reach
	for directory := range measured {
		reach.Measured = append(reach.Measured, directory)
	}

	for _, file := range files {
		if file.constrained {
			reach.Tagged = append(reach.Tagged, file.relative)
			if _, stated := declared[file.relative]; !stated {
				reach.Undeclared = append(reach.Undeclared, file.relative)
			}
			continue
		}
		if !measured[file.directory] {
			continue
		}
		for _, one := range file.imported {
			if reachesTheNetwork(one) {
				reach.Reaching = append(reach.Reaching, file.relative+" imports "+one)
			}
		}
	}

	for stated := range declared {
		if !slices.Contains(reach.Tagged, stated) {
			reach.Stale = append(reach.Stale, stated)
		}
	}

	slices.Sort(reach.Measured)
	slices.Sort(reach.Reaching)
	slices.Sort(reach.Tagged)
	slices.Sort(reach.Undeclared)
	slices.Sort(reach.Stale)
	return reach, nil
}

func readTree(root string) ([]sourceFile, error) {
	var files []sourceFile

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

		files = append(files, sourceFile{
			relative:    relative,
			directory:   filepath.Dir(relative),
			constrained: constrained,
			imported:    imported,
			test:        strings.HasSuffix(relative, "_test.go"),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (scope Scope) closure(files []sourceFile) map[string]bool {
	byDirectory := map[string][]sourceFile{}
	for _, file := range files {
		byDirectory[file.directory] = append(byDirectory[file.directory], file)
	}

	reached := map[string]bool{}

	var follow func(directory string)
	follow = func(directory string) {
		holding, present := byDirectory[directory]
		if !present || reached[directory] {
			return
		}
		reached[directory] = true

		for _, file := range holding {
			if file.test {
				continue
			}
			for _, one := range file.imported {
				if next, inside := scope.within(one); inside {
					follow(next)
				}
			}
		}
	}

	for _, start := range scope.From {
		follow(start)
	}
	return reached
}

func (scope Scope) within(imported string) (string, bool) {
	if scope.Module == "" {
		return "", false
	}
	if imported == scope.Module {
		return ".", true
	}
	beneath, inside := strings.CutPrefix(imported, scope.Module+"/")
	if !inside {
		return "", false
	}
	return filepath.FromSlash(beneath), true
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
