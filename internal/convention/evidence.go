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
	filepath.Join("review", "agent", "boundary_test.go"):             "the rule that one package alone names the vendor's binary, which is a structural invariant of this module rather than a claim about behaviour",
	filepath.Join("review", "claudesim", "failure_test.go"):          "the file the simulator wrote into a temporary directory, which is the observable side effect that scenario asks for",
	filepath.Join("internal", "convention", "account_test.go"):       "the account this project answers to, where the source is the subject of the check and not a proxy for anything else",
	filepath.Join("internal", "convention", "comments_test.go"):      "this repository's no-commentary convention, where the source is the subject of the check and not a proxy for anything else",
	filepath.Join("internal", "convention", "crosstool_test.go"):     "the rule that no tool imports another tool, which is a structural invariant of this module",
	filepath.Join("internal", "convention", "decisions_test.go"):     "the rule that every decision record, glossary and link resolves wherever a tool keeps it, which is a property of this repository's documents and reads no Go file at all",
	filepath.Join("internal", "convention", "documents_test.go"):     "what a person reads before they trust the tool, install it or propose a change to it, which is a property of this repository's documents and reads no Go file at all",
	filepath.Join("internal", "convention", "evidence_test.go"):      "this convention itself, which is a property of the test suite",
	filepath.Join("internal", "convention", "ignorance_test.go"):     "the rule that the surface names nothing of the other tool, which is a property of this repository's source rather than a proxy for behaviour",
	filepath.Join("internal", "convention", "languages_test.go"):     "the rule that the languages the page can highlight are stated once, where the bundle that ships is the subject of the check and not a proxy for behaviour",
	filepath.Join("internal", "convention", "licence_test.go"):       "the rule that the terms are stated once at the root and never as a header, where the source is the subject of the check and not a proxy for anything else",
	filepath.Join("internal", "convention", "nested_test.go"):        "the rule that a checkout nested under the root is not this repository's own source, where which files a check reads is the subject of the check and not a proxy for behaviour",
	filepath.Join("internal", "convention", "notice_test.go"):        "the rule that everything embedded is accounted for in the notice, where what a binary distributes is the subject of the check and not a proxy for anything else",
	filepath.Join("internal", "convention", "palette_test.go"):       "the rule that the page paints from one stated palette, which is a property of this repository's source rather than a proxy for behaviour",
	filepath.Join("internal", "convention", "platforms_test.go"):     "the rule that the supported platforms are stated once, where the source, what a release builds and what a person reads are the subjects of the check and not proxies for behaviour",
	filepath.Join("internal", "convention", "reach_test.go"):         "what socrates itself reaches and what is compiled out of the ordinary suite, which is a structural invariant of this module",
	filepath.Join("internal", "convention", "reloading_test.go"):     "the rule that only a test about reloading reloads the page, which is a property of the test suite",
	filepath.Join("internal", "convention", "release_test.go"):       "what a release builds and what its archives carry, which is a property of this repository's release configuration and reads no Go file at all",
	filepath.Join("internal", "convention", "rendering_test.go"):     "the rule that one place alone turns an event into a line, which is a structural invariant of this module",
	filepath.Join("internal", "convention", "vendors_test.go"):       "the rule that what a person reads names every reviewer CLI socrates drives and why the others are not driven, where the document is the subject of the check and not a proxy for behaviour",
	filepath.Join("internal", "convention", "variables_test.go"):     "the rule that a message names only a variable the source reads, where the source is the subject of the check and not a proxy for anything else",
	filepath.Join("internal", "convention", "skills_test.go"):        "the rule that no shipped skill can be started by a model, where the committed skills are the subject of the check and no Go file is read",
	filepath.Join("internal", "convention", "site_test.go"):          "the rule that the landing page installs with the command the front document states, which is a property of this repository's documents and reads no Go file at all",
	filepath.Join("internal", "convention", "supplychain_test.go"):   "what this repository's workflows may do and what they pin, which is a property of its continuous integration and reads no Go file at all",
	filepath.Join("internal", "convention", "wording_test.go"):       "the wording convention itself, where the source is the subject of the check and not a proxy for anything else",
	filepath.Join("internal", "convention", "toolagreement_test.go"): "the rule that the tool list agrees with the tree, which is a structural invariant of this module: which directories hold a command package is a fact about the layout",
	filepath.Join("internal", "convention", "workflows_test.go"):     "what this repository's workflows reach, which is a property of its continuous integration and reads no Go file at all",
	filepath.Join("review", "dialogue", "dialogue_test.go"):          "the direction the dialogue package may import in, which is a structural invariant of this module",
	filepath.Join("agora", "document", "boundary_test.go"):           "the direction the document contract may import in, which is a structural invariant of this module",
	filepath.Join("review", "findings", "classify_test.go"):          "the rule that one function decides a finding's effective action, which is a structural invariant of this module",
	filepath.Join("review", "findings", "verdict_test.go"):           "the direction the findings package may import in, which is a structural invariant of this module",
	filepath.Join("review", "presentation", "boundary_test.go"):      "what the renderer may import and call, which is what makes it a function of the events and nothing else",
	filepath.Join("agora", "render", "boundary_test.go"):             "what the page renderer may import, which is what keeps it a function of the document and nothing else",
	filepath.Join("agora", "stylesheet", "boundary_test.go"):         "the rule that one file alone reaches the parser the reduction is built on, which is a structural invariant of this module",
	filepath.Join("agora", "surface", "boundary_test.go"):            "the rule that the surface starts no process, which is a structural invariant of this module",
	filepath.Join("internal", "version", "boundary_test.go"):         "the rule that what the binary says about itself is read from nothing at run time, which is a structural invariant of this module",
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
			if path != root && skipped(path, entry.Name()) {
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
