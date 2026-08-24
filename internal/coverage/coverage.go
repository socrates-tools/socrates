package coverage

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

const (
	scenarioHeading  = "#### Scenario:"
	scenarioMarker   = "scenario:"
	specFile         = "spec.md"
	tasksFile        = "tasks.md"
	doneTask         = "- [x]"
	changesDirectory = "changes"
	archiveDirectory = "archive"

	renamedSection = "## RENAMED"
	sectionHeading = "## "
	renameFrom     = "- FROM:"
	renameTo       = "- TO:"
)

var specificationDirectories = []string{
	filepath.Join("openspec", "specs"),
	filepath.Join("openspec", "changes"),
}

var skippedDirectories = []string{".git", "bin", "node_modules", "testdata"}

type Scenario struct {
	Spec     string
	Line     int
	Text     string
	Underway bool
}

func Scenarios(root string) ([]Scenario, error) {
	var found []Scenario

	renames, err := renamedScenarios(root)
	if err != nil {
		return nil, err
	}

	for _, directory := range specificationDirectories {
		walked := filepath.Join(root, directory)
		if _, err := os.Stat(walked); err != nil {
			continue
		}

		err := filepath.WalkDir(walked, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() {
				if entry.Name() == archiveDirectory {
					return filepath.SkipDir
				}
				return nil
			}
			if filepath.Base(path) != specFile {
				return nil
			}

			content, readErr := os.ReadFile(path)
			if readErr != nil {
				return fmt.Errorf(readingASpecification, path, readErr)
			}

			relative, relErr := filepath.Rel(root, path)
			if relErr != nil {
				relative = path
			}

			underway, err := underway(root, relative)
			if err != nil {
				return err
			}

			for number, line := range strings.Split(string(content), "\n") {
				stated, headed := strings.CutPrefix(strings.TrimSpace(line), scenarioHeading)
				if !headed {
					continue
				}
				text := strings.TrimSpace(stated)
				if renamed, was := renames[Normalise(text)]; was && !strings.Contains(relative, changesDirectory) {
					text = renamed
				}
				found = append(found, Scenario{Spec: relative, Line: number + 1, Text: text, Underway: underway})
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	slices.SortFunc(found, func(one, other Scenario) int {
		if one.Spec != other.Spec {
			return strings.Compare(one.Spec, other.Spec)
		}
		return one.Line - other.Line
	})
	return found, nil
}

func renamedScenarios(root string) (map[string]string, error) {
	renames := map[string]string{}
	changes := filepath.Join(root, "openspec", changesDirectory)

	err := filepath.WalkDir(changes, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if entry.IsDir() {
			if entry.Name() == archiveDirectory {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Base(path) != specFile {
			return nil
		}

		relative, relErr := filepath.Rel(root, path)
		if relErr != nil {
			relative = path
		}
		started, err := underway(root, relative)
		if err != nil {
			return err
		}
		if !started {
			return nil
		}

		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return fmt.Errorf(readingASpecification, path, readErr)
		}

		section, from := "", ""
		for line := range strings.SplitSeq(string(content), "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, sectionHeading) && !strings.HasPrefix(trimmed, scenarioHeading) {
				section = trimmed
				continue
			}
			if !strings.HasPrefix(section, renamedSection) {
				continue
			}
			if named, marked := strings.CutPrefix(trimmed, renameFrom); marked {
				from = scenarioNamed(named)
			}
			if named, marked := strings.CutPrefix(trimmed, renameTo); marked && from != "" {
				if to := scenarioNamed(named); to != "" {
					renames[Normalise(from)] = to
				}
				from = ""
			}
		}
		return nil
	})
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return renames, nil
}

func scenarioNamed(line string) string {
	quoted := strings.Trim(strings.TrimSpace(line), "`")
	named, headed := strings.CutPrefix(strings.TrimSpace(quoted), scenarioHeading)
	if !headed {
		return ""
	}
	return strings.TrimSpace(named)
}

func underway(root, spec string) (bool, error) {
	parts := strings.Split(filepath.ToSlash(spec), "/")
	if len(parts) < 3 || parts[1] != changesDirectory {
		return true, nil
	}

	tasks := filepath.Join(root, parts[0], changesDirectory, parts[2], tasksFile)
	content, err := os.ReadFile(tasks)
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf(readingASpecification, tasks, err)
	}
	return strings.Contains(string(content), doneTask), nil
}

func Claims(root string) ([]string, error) {
	var claims []string

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
		if !strings.HasSuffix(path, "_test.go") {
			return nil
		}

		set := token.NewFileSet()
		parsed, parseErr := parser.ParseFile(set, path, nil, parser.ParseComments|parser.SkipObjectResolution)
		if parseErr != nil {
			return fmt.Errorf(readingATest, path, parseErr)
		}

		for _, declaration := range parsed.Decls {
			function, isFunction := declaration.(*ast.FuncDecl)
			if isFunction && strings.HasPrefix(function.Name.Name, "Test") {
				claims = append(claims, function.Name.Name)
			}
		}

		ast.Inspect(parsed, func(node ast.Node) bool {
			call, isCall := node.(*ast.CallExpr)
			if !isCall || len(call.Args) == 0 {
				return true
			}
			selected, isSelector := call.Fun.(*ast.SelectorExpr)
			if !isSelector || selected.Sel.Name != "Run" {
				return true
			}
			if receiver, isIdentifier := selected.X.(*ast.Ident); !isIdentifier || receiver.Name != "t" {
				return true
			}
			if literal, isLiteral := call.Args[0].(*ast.BasicLit); isLiteral && literal.Kind == token.STRING {
				if named, unquoteErr := strconv.Unquote(literal.Value); unquoteErr == nil {
					claims = append(claims, named)
				}
			}
			return true
		})

		for _, group := range parsed.Comments {
			for _, comment := range group.List {
				text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
				if stated, marked := strings.CutPrefix(text, scenarioMarker); marked {
					claims = append(claims, strings.TrimSpace(stated))
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	normalised := make([]string, 0, len(claims))
	for _, claim := range claims {
		normalised = append(normalised, Normalise(claim))
	}
	slices.Sort(normalised)
	return slices.Compact(normalised), nil
}

func Gaps(scenarios []Scenario, claims []string) []Scenario {
	var gaps []Scenario

	for _, scenario := range scenarios {
		wanted := Normalise(scenario.Text)
		if wanted == "" {
			gaps = append(gaps, scenario)
			continue
		}

		covered := false
		for _, claim := range claims {
			if strings.Contains(claim, wanted) {
				covered = true
				break
			}
		}
		if !covered {
			gaps = append(gaps, scenario)
		}
	}
	return gaps
}

func Normalise(text string) string {
	var built strings.Builder
	for _, character := range strings.ToLower(text) {
		if unicode.IsLetter(character) || unicode.IsDigit(character) {
			built.WriteRune(character)
		}
	}
	return built.String()
}
