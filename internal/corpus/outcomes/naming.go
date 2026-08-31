package outcomes

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/socrates-devtool/socrates/internal/runstate"
)

const (
	statePackage = "runstate"
	testPrefix   = "Test"
	testSuffix   = "_test.go"
)

var skippedDirectories = []string{".git", "bin", "node_modules", "testdata"}

func Naming(root string) (map[runstate.State][]string, error) {
	asserting := map[runstate.State]map[string]bool{}

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
		if !strings.HasSuffix(path, testSuffix) {
			return nil
		}

		set := token.NewFileSet()
		parsed, err := parser.ParseFile(set, path, nil, parser.SkipObjectResolution)
		if err != nil {
			return fmt.Errorf(readingATest, path, err)
		}

		for _, declaration := range parsed.Decls {
			function, isFunction := declaration.(*ast.FuncDecl)
			if !isFunction || !strings.HasPrefix(function.Name.Name, testPrefix) {
				continue
			}
			for state := range statesNamedIn(function) {
				if asserting[state] == nil {
					asserting[state] = map[string]bool{}
				}
				asserting[state][function.Name.Name] = true
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	named := map[runstate.State][]string{}
	for state, functions := range asserting {
		for function := range functions {
			named[state] = append(named[state], function)
		}
		slices.Sort(named[state])
	}
	return named, nil
}

func statesNamedIn(function *ast.FuncDecl) map[runstate.State]bool {
	named := map[runstate.State]bool{}

	ast.Inspect(function, func(node ast.Node) bool {
		switch stated := node.(type) {
		case *ast.SelectorExpr:
			qualifier, isIdentifier := stated.X.(*ast.Ident)
			if isIdentifier && qualifier.Name == statePackage {
				if state, known := stateCalled(stated.Sel.Name); known {
					named[state] = true
				}
			}
		case *ast.BasicLit:
			if stated.Kind != token.STRING {
				return true
			}
			text, err := strconv.Unquote(stated.Value)
			if err != nil {
				return true
			}
			if state := runstate.State(text); state.Valid() {
				named[state] = true
			}
		}
		return true
	})

	return named
}

func stateCalled(name string) (runstate.State, bool) {
	for _, state := range runstate.States() {
		if strings.EqualFold(strings.ReplaceAll(string(state), "-", ""), name) {
			return state, true
		}
	}
	return "", false
}
