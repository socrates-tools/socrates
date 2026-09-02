package convention

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"
)

const pageReloadMethod = "reload"

var PageDrivingTests = filepath.Join("internal", "e2e")

var TestsThatReloadThePage = map[string]string{
	"TestAReloadedPageRendersClosed": "the closed state a page renders when it is reloaded onto an exchange that ended",
	"TestAReloadedPageStillReopens":  "that reloading remains sufficient, so the tests that stopped reloading gave up no guarantee",
}

func TestsReloadingThePage(root string) (Naming, error) {
	var naming Naming

	err := filepath.WalkDir(filepath.Join(root, PageDrivingTests), func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || !strings.HasSuffix(path, "_test.go") {
			return nil
		}

		read, parseErr := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if parseErr != nil {
			return parseErr
		}

		for _, declared := range read.Decls {
			function, isFunction := declared.(*ast.FuncDecl)
			if !isFunction || function.Recv != nil || !reloadsThePage(function) {
				continue
			}
			if _, allowed := TestsThatReloadThePage[function.Name.Name]; allowed {
				naming.Named = append(naming.Named, function.Name.Name)
				continue
			}
			naming.Unnamed = append(naming.Unnamed, function.Name.Name)
		}
		return nil
	})
	if err != nil {
		return Naming{}, err
	}

	slices.Sort(naming.Unnamed)
	slices.Sort(naming.Named)

	for declared := range TestsThatReloadThePage {
		if !slices.Contains(naming.Named, declared) {
			naming.Stale = append(naming.Stale, declared)
		}
	}
	slices.Sort(naming.Stale)

	return naming, nil
}

func reloadsThePage(function *ast.FuncDecl) bool {
	reloads := false
	ast.Inspect(function, func(node ast.Node) bool {
		called, isCall := node.(*ast.CallExpr)
		if !isCall {
			return true
		}
		if selected, isSelector := called.Fun.(*ast.SelectorExpr); isSelector && selected.Sel.Name == pageReloadMethod {
			reloads = true
			return false
		}
		return true
	})
	return reloads
}
