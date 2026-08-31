package convention

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
)

var ContractsTestsRead = map[string]string{
	filepath.Join("internal", "agent", "boundary_test.go"):        "the rule that one package alone names the vendor's binary, which is a structural invariant of this module rather than a claim about behaviour",
	filepath.Join("internal", "claudesim", "failure_test.go"):     "the file the simulator wrote into a temporary directory, which is the observable side effect that scenario asks for",
	filepath.Join("internal", "convention", "account_test.go"):    "the account this project answers to, where the source is the subject of the check and not a proxy for anything else",
	filepath.Join("internal", "convention", "comments_test.go"):   "this repository's no-commentary convention, where the source is the subject of the check and not a proxy for anything else",
	filepath.Join("internal", "convention", "context_test.go"):    "the normative context itself, which is a document of this repository rather than a proxy for behaviour",
	filepath.Join("internal", "convention", "evidence_test.go"):   "this convention itself, which is a property of the test suite",
	filepath.Join("internal", "convention", "ignorance_test.go"):  "the rule that the surface names nothing of the other tool, which is a property of this repository's source rather than a proxy for behaviour",
	filepath.Join("internal", "convention", "palette_test.go"):    "the rule that the page paints from one stated palette, which is a property of this repository's source rather than a proxy for behaviour",
	filepath.Join("internal", "convention", "reach_test.go"):      "what socrates itself reaches and what is compiled out of the ordinary suite, which is a structural invariant of this module",
	filepath.Join("internal", "convention", "rendering_test.go"):  "the rule that one place alone turns an event into a line, which is a structural invariant of this module",
	filepath.Join("internal", "convention", "variables_test.go"):  "the rule that a message names only a variable the source reads, where the source is the subject of the check and not a proxy for anything else",
	filepath.Join("internal", "convention", "wording_test.go"):    "the wording convention itself, where the source is the subject of the check and not a proxy for anything else",
	filepath.Join("internal", "convention", "workflows_test.go"):  "what this repository's workflows reach, which is a property of its continuous integration and reads no Go file at all",
	filepath.Join("internal", "dialogue", "dialogue_test.go"):     "the direction the dialogue package may import in, which is a structural invariant of this module",
	filepath.Join("internal", "document", "boundary_test.go"):     "the direction the document contract may import in, which is a structural invariant of this module",
	filepath.Join("internal", "findings", "classify_test.go"):     "the rule that one function decides a finding's effective action, which is a structural invariant of this module",
	filepath.Join("internal", "findings", "verdict_test.go"):      "the direction the findings package may import in, which is a structural invariant of this module",
	filepath.Join("internal", "presentation", "boundary_test.go"): "what the renderer may import and call, which is what makes it a function of the events and nothing else",
	filepath.Join("internal", "render", "boundary_test.go"):       "what the page renderer may import, which is what keeps it a function of the document and nothing else",
	filepath.Join("internal", "stylesheet", "boundary_test.go"):   "the rule that one file alone reaches the parser the reduction is built on, which is a structural invariant of this module",
	filepath.Join("internal", "surface", "boundary_test.go"):      "the rule that the surface starts no process, which is a structural invariant of this module",
	filepath.Join("internal", "version", "boundary_test.go"):      "the rule that what the binary says about itself is read from nothing at run time, which is a structural invariant of this module",
}

var sourceListingCalls = map[string][]string{"source": {"Files", "Imports", "DirectImports", "PackagesCalled"}}

var fileReadingCalls = map[string][]string{
	"os":       {"ReadFile", "Open", "OpenFile", "ReadDir"},
	"filepath": {"Walk", "WalkDir", "Glob"},
	"fs":       {"ReadFile", "WalkDir", "Glob"},
}

var sourceReadingImports = []string{"go/ast", "go/parser", "go/token"}

var sourceReadingPackages = []string{filepath.Join("internal", "convention")}

type Reading struct {
	Unnamed []string
	Named   []string
	Stale   []string
}

func TestsReadingGoFiles(root string, named map[string]string) (Reading, error) {
	var reading Reading

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
		if !strings.HasSuffix(path, "_test.go") {
			return nil
		}

		relative, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return relErr
		}

		reads, readsErr := readsGoFiles(path, relative)
		if readsErr != nil {
			return readsErr
		}
		if !reads {
			return nil
		}

		if _, declared := named[relative]; declared {
			reading.Named = append(reading.Named, relative)
			return nil
		}
		reading.Unnamed = append(reading.Unnamed, relative)
		return nil
	})
	if err != nil {
		return Reading{}, err
	}

	slices.Sort(reading.Unnamed)
	slices.Sort(reading.Named)

	for declared := range named {
		if !slices.Contains(reading.Named, declared) {
			reading.Stale = append(reading.Stale, declared)
		}
	}
	slices.Sort(reading.Stale)

	return reading, nil
}

func readsGoFiles(path, relative string) (bool, error) {
	for _, owner := range sourceReadingPackages {
		if strings.HasPrefix(relative, owner+string(filepath.Separator)) {
			return true, nil
		}
	}

	set := token.NewFileSet()
	parsed, err := parser.ParseFile(set, path, nil, parser.SkipObjectResolution)
	if err != nil {
		return false, fmt.Errorf("parsing %s: %w", relative, err)
	}

	for _, specification := range parsed.Imports {
		imported, err := strconv.Unquote(specification.Path.Value)
		if err != nil {
			return false, fmt.Errorf("reading an import path in %s: %w", relative, err)
		}
		if slices.Contains(sourceReadingImports, imported) {
			return true, nil
		}
	}

	reads := false
	ast.Inspect(parsed, func(node ast.Node) bool {
		call, isCall := node.(*ast.CallExpr)
		if !isCall {
			return !reads
		}
		selected, isSelector := call.Fun.(*ast.SelectorExpr)
		if !isSelector {
			return !reads
		}
		qualifier, isIdentifier := selected.X.(*ast.Ident)
		if !isIdentifier {
			return !reads
		}
		if slices.Contains(sourceListingCalls[qualifier.Name], selected.Sel.Name) {
			reads = true
		}
		if slices.Contains(fileReadingCalls[qualifier.Name], selected.Sel.Name) && namesAGoFile(call) {
			reads = true
		}
		return !reads
	})

	return reads, nil
}

func namesAGoFile(call *ast.CallExpr) bool {
	named := false

	for _, argument := range call.Args {
		ast.Inspect(argument, func(node ast.Node) bool {
			literal, isLiteral := node.(*ast.BasicLit)
			if !isLiteral || literal.Kind != token.STRING {
				return !named
			}
			value, err := strconv.Unquote(literal.Value)
			if err != nil {
				return !named
			}
			if strings.HasSuffix(value, ".go") {
				named = true
			}
			return !named
		})
	}

	return named
}
