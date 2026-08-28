package convention

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var environmentVariable = regexp.MustCompile(`\b[A-Z][A-Z0-9]*(?:_[A-Z0-9]+)+\b`)

type NamedVariable struct {
	File string
	Name string
}

func VariablesNamedToAPerson(root string) ([]NamedVariable, error) {
	read, err := variablesTheSourceReads(root)
	if err != nil {
		return nil, err
	}

	var unread []NamedVariable
	err = filepath.WalkDir(root, func(path string, entry fs.DirEntry, walked error) error {
		if walked != nil {
			return walked
		}
		if entry.IsDir() {
			if path != root && skipped(entry.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.Name() != RegistryFile {
			return nil
		}

		named, err := variablesNamedIn(path)
		if err != nil {
			return err
		}
		relative, relErr := filepath.Rel(root, path)
		if relErr != nil {
			relative = path
		}
		for _, name := range named {
			if !read[name] {
				unread = append(unread, NamedVariable{File: relative, Name: name})
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	slices.SortFunc(unread, func(a, b NamedVariable) int {
		if a.File != b.File {
			return strings.Compare(a.File, b.File)
		}
		return strings.Compare(a.Name, b.Name)
	})
	return unread, nil
}

func variablesNamedIn(path string) ([]string, error) {
	set := token.NewFileSet()
	parsed, err := parser.ParseFile(set, path, nil, parser.SkipObjectResolution)
	if err != nil {
		return nil, err
	}

	var named []string
	ast.Inspect(parsed, func(node ast.Node) bool {
		literal, isLiteral := node.(*ast.BasicLit)
		if !isLiteral || literal.Kind != token.STRING {
			return true
		}
		text, unquoteErr := strconv.Unquote(literal.Value)
		if unquoteErr != nil {
			return true
		}
		named = append(named, environmentVariable.FindAllString(text, -1)...)
		return true
	})

	slices.Sort(named)
	return slices.Compact(named), nil
}

func variablesTheSourceReads(root string) (map[string]bool, error) {
	read := map[string]bool{}

	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walked error) error {
		if walked != nil {
			return walked
		}
		if entry.IsDir() {
			if path != root && skipped(entry.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") || filepath.Base(path) == RegistryFile {
			return nil
		}

		set := token.NewFileSet()
		parsed, err := parser.ParseFile(set, path, nil, parser.SkipObjectResolution)
		if err != nil {
			return err
		}

		ast.Inspect(parsed, func(node ast.Node) bool {
			literal, isLiteral := node.(*ast.BasicLit)
			if !isLiteral || literal.Kind != token.STRING {
				return true
			}
			text, unquoteErr := strconv.Unquote(literal.Value)
			if unquoteErr != nil {
				return true
			}
			if environmentVariable.FindString(text) == text {
				read[text] = true
			}
			return true
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return read, nil
}
