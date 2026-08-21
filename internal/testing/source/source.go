package source

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/rui-ktei/socrates/internal/corpus"
)

const ModulePath = corpus.ModulePath

var skippedDirectories = []string{".git", "bin", "node_modules", "testdata"}

func Root(t *testing.T) string {
	t.Helper()

	directory, err := os.Getwd()
	if err != nil {
		t.Fatalf("locating the working directory: %v", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(directory, "go.mod")); err == nil {
			return directory
		}
		parent := filepath.Dir(directory)
		if parent == directory {
			t.Fatalf("no go.mod above %s", directory)
		}
		directory = parent
	}
}

func Files(t *testing.T, root string) []string {
	t.Helper()

	var found []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if path != root && slices.Contains(skippedDirectories, entry.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			found = append(found, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("reading %s: %v", root, err)
	}

	return found
}

func DirectImports(t *testing.T, root, packageDirectory string) []string {
	t.Helper()

	seen := map[string]bool{}
	for _, path := range packageFiles(t, filepath.Join(root, packageDirectory)) {
		for _, imported := range importsOf(t, path) {
			seen[imported] = true
		}
	}
	return sorted(seen)
}

func PackagesCalled(t *testing.T, root, packageDirectory string) []string {
	t.Helper()

	seen := map[string]bool{}
	for _, path := range packageFiles(t, filepath.Join(root, packageDirectory)) {
		set := token.NewFileSet()
		parsed, err := parser.ParseFile(set, path, nil, parser.SkipObjectResolution)
		if err != nil {
			t.Fatalf("parsing %s: %v", path, err)
		}

		ast.Inspect(parsed, func(node ast.Node) bool {
			call, isCall := node.(*ast.CallExpr)
			if !isCall {
				return true
			}
			selected, isSelector := call.Fun.(*ast.SelectorExpr)
			if !isSelector {
				return true
			}
			if qualifier, isIdentifier := selected.X.(*ast.Ident); isIdentifier {
				seen[qualifier.Name] = true
			}
			return true
		})
	}
	return sorted(seen)
}

func packageFiles(t *testing.T, directory string) []string {
	t.Helper()

	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatalf("reading %s: %v", directory, err)
	}

	var found []string
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		found = append(found, filepath.Join(directory, name))
	}
	slices.Sort(found)
	return found
}

func sorted(seen map[string]bool) []string {
	found := make([]string, 0, len(seen))
	for name := range seen {
		found = append(found, name)
	}
	slices.Sort(found)
	return found
}

func Imports(t *testing.T, root, packageDirectory string) []string {
	t.Helper()

	seen := map[string]bool{}
	visited := map[string]bool{}

	var collect func(directory string)
	collect = func(directory string) {
		if visited[directory] {
			return
		}
		visited[directory] = true

		entries, err := os.ReadDir(directory)
		if err != nil {
			t.Fatalf("reading %s: %v", directory, err)
		}

		for _, entry := range entries {
			name := entry.Name()
			if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
				continue
			}
			for _, imported := range importsOf(t, filepath.Join(directory, name)) {
				seen[imported] = true
				if within, found := strings.CutPrefix(imported, ModulePath+"/"); found {
					collect(filepath.Join(root, within))
				}
			}
		}
	}
	collect(filepath.Join(root, packageDirectory))

	imported := make([]string, 0, len(seen))
	for path := range seen {
		imported = append(imported, path)
	}
	slices.Sort(imported)
	return imported
}

func importsOf(t *testing.T, path string) []string {
	t.Helper()

	set := token.NewFileSet()
	parsed, err := parser.ParseFile(set, path, nil, parser.ImportsOnly|parser.SkipObjectResolution)
	if err != nil {
		t.Fatalf("parsing %s: %v", path, err)
	}

	var imported []string
	for _, specification := range parsed.Imports {
		unquoted, err := strconv.Unquote(specification.Path.Value)
		if err != nil {
			t.Fatalf("reading an import path in %s: %v", path, err)
		}
		imported = append(imported, unquoted)
	}
	return imported
}
