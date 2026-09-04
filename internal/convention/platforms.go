package convention

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/ktei/socrates/internal/platform"
)

const (
	StatementDirectory        = "internal/platform"
	DocumentedHeading         = "## Supported platforms"
	TheWayForward             = "WSL"
	WhereTheWayForwardApplies = "Windows"

	documentationFile = "README.md"
	tableRow          = "|"
	tableSeparator    = "-"
)

func PlatformsNamedOutsideTheStatement(root string, named []string) ([]string, error) {
	var copies []string

	statement := filepath.FromSlash(StatementDirectory)

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

		relative, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return relErr
		}
		if relative == statement || strings.HasPrefix(relative, statement+string(filepath.Separator)) {
			return nil
		}

		set := token.NewFileSet()
		parsed, parseErr := parser.ParseFile(set, path, nil, parser.SkipObjectResolution)
		if parseErr != nil {
			return fmt.Errorf("parsing %s: %w", relative, parseErr)
		}

		ast.Inspect(parsed, func(node ast.Node) bool {
			literal, isLiteral := node.(*ast.BasicLit)
			if !isLiteral {
				return true
			}
			if text, quoted := stringLiteral(literal); quoted && slices.Contains(named, text) {
				copies = append(copies, fmt.Sprintf("%s:%d states %q", relative, set.Position(literal.Pos()).Line, text))
			}
			return true
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	slices.Sort(copies)
	return copies, nil
}

func PlatformsAReleaseBuilds(root, targets string) ([]platform.Target, error) {
	held, err := os.ReadFile(filepath.Join(root, targets))
	if err != nil {
		return nil, err
	}

	var built []platform.Target
	if err := json.Unmarshal(held, &built); err != nil {
		return nil, fmt.Errorf("reading %s: %w", targets, err)
	}
	return built, nil
}

func PlatformsDocumented(root string) ([]string, error) {
	held, err := os.ReadFile(filepath.Join(root, documentationFile))
	if err != nil {
		return nil, err
	}

	section, found := sectionUnder(string(held), DocumentedHeading)
	if !found {
		return nil, nil
	}
	return namesTabulatedIn(section), nil
}

func TheWayForwardIsDocumented(root string) (bool, error) {
	held, err := os.ReadFile(filepath.Join(root, documentationFile))
	if err != nil {
		return false, err
	}

	section, found := sectionUnder(string(held), DocumentedHeading)
	if !found {
		return false, nil
	}
	return strings.Contains(section, WhereTheWayForwardApplies) && strings.Contains(section, TheWayForward), nil
}

func PlatformDisagreements(stated, elsewhere []string, where string) []string {
	var disagreeing []string

	for _, one := range elsewhere {
		if !slices.Contains(stated, one) {
			disagreeing = append(disagreeing, fmt.Sprintf("%s names %s, and the supported set does not", where, one))
		}
	}
	for _, one := range stated {
		if !slices.Contains(elsewhere, one) {
			disagreeing = append(disagreeing, fmt.Sprintf("the supported set names %s, and %s does not", one, where))
		}
	}

	slices.Sort(disagreeing)
	return disagreeing
}

func sectionUnder(document, heading string) (string, bool) {
	_, under, found := strings.Cut(document, heading)
	if !found {
		return "", false
	}
	if before, _, ended := strings.Cut(under, "\n## "); ended {
		return before, true
	}
	return under, true
}

func namesTabulatedIn(section string) []string {
	var rows []string
	separated := -1

	for _, line := range strings.Split(section, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, tableRow) {
			continue
		}
		named := strings.TrimSpace(firstCellOf(trimmed))
		if strings.HasPrefix(named, tableSeparator) {
			separated = len(rows)
		}
		rows = append(rows, named)
	}

	if separated < 0 {
		return nil
	}
	return rows[separated+1:]
}

func firstCellOf(row string) string {
	cells := strings.Split(strings.Trim(row, tableRow), tableRow)
	if len(cells) == 0 {
		return ""
	}
	return cells[0]
}
