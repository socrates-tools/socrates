package convention

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"
)

const eventConstantPrefix = "Event"

var EventKindsNamed = map[string]string{
	filepath.Join("review", "presentation", "messages.go"): "the one place an event kind becomes a line a person reads",
	filepath.Join("review", "follow", "follow.go"):         "the rotation record, which tells a follower what the log kept when it was truncated",
	filepath.Join("review", "corpus", "runs", "runs.go"):   "which model answered a round, taken out of the log as data rather than rendered",
}

type Naming struct {
	Unnamed []string
	Named   []string
	Stale   []string
}

func FilesNamingEventKinds(root string, kinds, defining []string, named map[string]string) (Naming, error) {
	var naming Naming

	qualifiers := make([]string, 0, len(defining))
	for _, directory := range defining {
		qualifiers = append(qualifiers, filepath.Base(directory))
	}

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
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		relative, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return relErr
		}
		for _, directory := range defining {
			if strings.HasPrefix(relative, directory+string(filepath.Separator)) {
				return nil
			}
		}

		names, namesErr := namesAnEventKind(path, relative, kinds, qualifiers)
		if namesErr != nil {
			return namesErr
		}
		if !names {
			return nil
		}

		if _, declared := named[relative]; declared {
			naming.Named = append(naming.Named, relative)
			return nil
		}
		naming.Unnamed = append(naming.Unnamed, relative)
		return nil
	})
	if err != nil {
		return Naming{}, err
	}

	slices.Sort(naming.Unnamed)
	slices.Sort(naming.Named)

	for declared := range named {
		if !slices.Contains(naming.Named, declared) {
			naming.Stale = append(naming.Stale, declared)
		}
	}
	slices.Sort(naming.Stale)

	return naming, nil
}

func namesAnEventKind(path, relative string, kinds, qualifiers []string) (bool, error) {
	set := token.NewFileSet()
	parsed, err := parser.ParseFile(set, path, nil, parser.SkipObjectResolution)
	if err != nil {
		return false, fmt.Errorf("parsing %s: %w", relative, err)
	}

	names := false
	ast.Inspect(parsed, func(node ast.Node) bool {
		if literal, isLiteral := node.(*ast.BasicLit); isLiteral {
			if text, quoted := stringLiteral(literal); quoted && slices.Contains(kinds, text) {
				names = true
			}
			return !names
		}

		selected, isSelector := node.(*ast.SelectorExpr)
		if !isSelector {
			return !names
		}
		qualifier, isIdentifier := selected.X.(*ast.Ident)
		if !isIdentifier {
			return !names
		}
		kind := selected.Sel.Name
		if slices.Contains(qualifiers, qualifier.Name) && kind != eventConstantPrefix && strings.HasPrefix(kind, eventConstantPrefix) {
			names = true
		}
		return !names
	})

	return names, nil
}
