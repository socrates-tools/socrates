package stylesheet_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/testing/source"
)

const theSeam = "parse.go"

var handedOut = []string{"Scope", "string", "bool"}

func reduction(t *testing.T) map[string]*ast.File {
	t.Helper()

	paths, err := filepath.Glob(filepath.Join(source.Root(t), "agora", "stylesheet", "*.go"))
	if err != nil {
		t.Fatalf("reading the reduction: %v", err)
	}

	read := map[string]*ast.File{}
	for _, path := range paths {
		named := filepath.Base(path)
		if strings.HasSuffix(named, "_test.go") {
			continue
		}
		set := token.NewFileSet()
		parsed, parseErr := parser.ParseFile(set, path, nil, parser.SkipObjectResolution)
		if parseErr != nil {
			t.Fatalf("reading %s: %v", named, parseErr)
		}
		read[named] = parsed
	}

	if len(read) == 0 {
		t.Fatal("the reduction holds no source at all, and a check that saw nothing would pass forever")
	}
	return read
}

func TestTheParserIsReachedFromOneFileOnly(t *testing.T) {
	// scenario: A stylesheet that does not parse embeds nothing
	reaching := 0

	for named, parsed := range reduction(t) {
		for _, imported := range parsed.Imports {
			if !strings.Contains(imported.Path.Value, "tdewolff") {
				continue
			}
			reaching++
			if named != theSeam {
				t.Errorf("%s reaches %s directly, and the reduction owns one seam over a parser in %s", named, imported.Path.Value, theSeam)
			}
		}
	}

	if reaching == 0 {
		t.Fatal("no file reaches a parser at all, so the reduction matches patterns rather than reading structure")
	}
}

func TestTheParsersTypesStayBehindTheSeam(t *testing.T) {
	// scenario: A stylesheet outside the allowlist is reduced
	exported := 0

	for _, parsed := range reduction(t) {
		for _, declared := range parsed.Decls {
			function, isFunction := declared.(*ast.FuncDecl)
			if !isFunction || !function.Name.IsExported() {
				continue
			}
			exported++
			for _, mentioned := range typesNamedBy(function.Type) {
				if !slices.Contains(handedOut, mentioned) {
					t.Errorf("%s names %s, and what the reduction hands out is its own vocabulary rather than a parser's", function.Name.Name, mentioned)
				}
			}
		}
	}

	if exported == 0 {
		t.Fatal("the reduction exports nothing, and a check that saw nothing would pass forever")
	}
}

func typesNamedBy(signature *ast.FuncType) []string {
	var named []string

	for _, group := range []*ast.FieldList{signature.Params, signature.Results} {
		if group == nil {
			continue
		}
		for _, field := range group.List {
			ast.Inspect(field.Type, func(node ast.Node) bool {
				switch held := node.(type) {
				case *ast.Ident:
					named = append(named, held.Name)
				case *ast.SelectorExpr:
					named = append(named, held.Sel.Name)
					return false
				}
				return true
			})
		}
	}
	return named
}
