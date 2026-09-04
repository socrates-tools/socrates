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

const RegistryFile = "messages.go"

var EmDash = string(rune(0x2014))

var packagesWithoutAMessageRegistry = []string{
	filepath.Join("cmd", "claudesim"),
	filepath.Join("cmd", "codexsim"),
	filepath.Join("cmd", "opencodesim"),
	filepath.Join("internal", "claudesim"),
	filepath.Join("internal", "codexsim"),
	filepath.Join("internal", "opencodesim"),
	filepath.Join("internal", "simulation"),
	filepath.Join("internal", "convention"),
	filepath.Join("internal", "prompt"),
	filepath.Join("internal", "testing"),
	filepath.Join("internal", "testquality"),
}

var messageCalls = map[string]map[string]int{
	"fmt":     {"Errorf": 0, "Sprintf": 0, "Printf": 0, "Println": 0, "Print": 0, "Fprintf": 1, "Fprintln": 1, "Fprint": 1},
	"errors":  {"New": 0},
	"refusal": {"New": 0, "Wrapping": 1},
}

var refusalCalls = map[string]int{"New": 0, "Wrapping": 1}

var nextActionVerbs = []string{
	"ask", "cancel", "check", "choose", "commit", "delete", "edit", "fix", "give",
	"move", "name", "pass", "point", "read", "remove", "rename", "run", "set",
	"start", "stop", "use", "wait",
}

type Wording struct {
	Inline     []string
	EmDashes   []string
	Actionless []string
}

func UserFacingStrings(root string) (Wording, error) {
	var wording Wording

	directories, err := packageDirectories(root)
	if err != nil {
		return Wording{}, err
	}

	for _, directory := range directories {
		registered, files, err := packageMessages(root, directory)
		if err != nil {
			return Wording{}, err
		}

		enforced := !exemptFromTheRegistry(directory)
		for _, file := range files {
			read, err := readWording(root, file, registered, enforced)
			if err != nil {
				return Wording{}, err
			}
			wording.Inline = append(wording.Inline, read.Inline...)
			wording.EmDashes = append(wording.EmDashes, read.EmDashes...)
			wording.Actionless = append(wording.Actionless, read.Actionless...)
		}
	}

	slices.Sort(wording.Inline)
	slices.Sort(wording.EmDashes)
	slices.Sort(wording.Actionless)
	return wording, nil
}

func exemptFromTheRegistry(directory string) bool {
	for _, exempt := range packagesWithoutAMessageRegistry {
		if directory == exempt || strings.HasPrefix(directory, exempt+string(filepath.Separator)) {
			return true
		}
	}
	return false
}

func packageDirectories(root string) ([]string, error) {
	seen := map[string]bool{}

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
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		relative, relErr := filepath.Rel(root, filepath.Dir(path))
		if relErr != nil {
			return relErr
		}
		seen[relative] = true
		return nil
	})
	if err != nil {
		return nil, err
	}

	directories := make([]string, 0, len(seen))
	for directory := range seen {
		directories = append(directories, directory)
	}
	slices.Sort(directories)
	return directories, nil
}

func packageMessages(root, directory string) (map[string]string, []string, error) {
	entries, err := readSourceFiles(filepath.Join(root, directory))
	if err != nil {
		return nil, nil, err
	}

	registered := map[string]string{}
	for _, path := range entries {
		set := token.NewFileSet()
		parsed, err := parser.ParseFile(set, path, nil, parser.SkipObjectResolution)
		if err != nil {
			return nil, nil, fmt.Errorf("parsing %s: %w", path, err)
		}

		for _, declaration := range parsed.Decls {
			constants, isConstant := declaration.(*ast.GenDecl)
			if !isConstant || constants.Tok != token.CONST {
				continue
			}
			for _, specification := range constants.Specs {
				valued, isValued := specification.(*ast.ValueSpec)
				if !isValued {
					continue
				}
				for position, name := range valued.Names {
					if position >= len(valued.Values) {
						continue
					}
					if text, quoted := stringLiteral(valued.Values[position]); quoted {
						registered[name.Name] = text
					}
				}
			}
		}
	}

	return registered, entries, nil
}

func readSourceFiles(directory string) ([]string, error) {
	entries, err := filepath.Glob(filepath.Join(directory, "*.go"))
	if err != nil {
		return nil, err
	}

	var found []string
	for _, path := range entries {
		if !strings.HasSuffix(path, "_test.go") {
			found = append(found, path)
		}
	}
	slices.Sort(found)
	return found, nil
}

func readWording(root, path string, registered map[string]string, enforced bool) (Wording, error) {
	relative, err := filepath.Rel(root, path)
	if err != nil {
		relative = path
	}

	set := token.NewFileSet()
	parsed, err := parser.ParseFile(set, path, nil, parser.SkipObjectResolution)
	if err != nil {
		return Wording{}, fmt.Errorf("parsing %s: %w", relative, err)
	}

	var wording Wording
	at := func(position token.Pos) string {
		return fmt.Sprintf("%s:%d", relative, set.Position(position).Line)
	}

	ast.Inspect(parsed, func(node ast.Node) bool {
		if literal, isLiteral := node.(*ast.BasicLit); isLiteral {
			text, quoted := stringLiteral(literal)
			if quoted && strings.Contains(text, EmDash) {
				wording.EmDashes = append(wording.EmDashes, at(literal.Pos()))
			}
			return true
		}

		call, isCall := node.(*ast.CallExpr)
		if !isCall {
			return true
		}
		qualifier, name, found := calledFunction(call)
		if !found {
			return true
		}

		if qualifier == "refusal" {
			if position, carries := refusalCalls[name]; carries && position < len(call.Args) {
				if stated, quoted := messageOf(call.Args[position], registered); quoted && !namesANextAction(stated) {
					wording.Actionless = append(wording.Actionless, at(call.Pos()))
				}
			}
		}

		if enforced {
			position, carries := messageCalls[qualifier][name]
			if carries && position < len(call.Args) {
				if text, quoted := stringLiteral(call.Args[position]); quoted && isProse(text) {
					wording.Inline = append(wording.Inline, at(call.Args[position].Pos()))
				}
			}
		}
		return true
	})

	return wording, nil
}

func calledFunction(call *ast.CallExpr) (string, string, bool) {
	selected, isSelector := call.Fun.(*ast.SelectorExpr)
	if !isSelector {
		return "", "", false
	}
	qualifier, isIdentifier := selected.X.(*ast.Ident)
	if !isIdentifier {
		return "", "", false
	}
	return qualifier.Name, selected.Sel.Name, true
}

func messageOf(argument ast.Expr, registered map[string]string) (string, bool) {
	if text, quoted := stringLiteral(argument); quoted {
		return text, true
	}
	if named, isIdentifier := argument.(*ast.Ident); isIdentifier {
		text, found := registered[named.Name]
		return text, found
	}
	return "", false
}

func stringLiteral(expression ast.Expr) (string, bool) {
	literal, isLiteral := expression.(*ast.BasicLit)
	if !isLiteral || literal.Kind != token.STRING {
		return "", false
	}
	text, err := strconv.Unquote(literal.Value)
	if err != nil {
		return "", false
	}
	return text, true
}

func isProse(text string) bool {
	return strings.Contains(strings.TrimSpace(text), " ")
}

func namesANextAction(message string) bool {
	for _, fragment := range fragmentsOf(message) {
		words := strings.Fields(fragment)
		if len(words) == 0 {
			continue
		}
		if slices.Contains(nextActionVerbs, strings.ToLower(strings.Trim(words[0], ".,:;\"'"))) {
			return true
		}
	}
	return false
}

func fragmentsOf(message string) []string {
	separated := message
	for _, separator := range []string{": ", ". ", ", so ", ", and ", "\n"} {
		separated = strings.ReplaceAll(separated, separator, "\x00")
	}

	var fragments []string
	for fragment := range strings.SplitSeq(separated, "\x00") {
		fragments = append(fragments, strings.TrimSpace(fragment))
	}
	return fragments
}
