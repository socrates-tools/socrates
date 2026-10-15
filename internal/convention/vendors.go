package convention

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const (
	VendorStatementFile = "review/agent/vendor.go"
	HostStatementFile   = "internal/skill/host.go"
)

func VendorsNamedOutsideTheStatement(root string, named []string) ([]string, error) {
	var copies []string

	statements := []string{filepath.FromSlash(VendorStatementFile), filepath.FromSlash(HostStatementFile)}

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
		if slices.Contains(statements, relative) {
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

const (
	SeamFile     = "review/agent/seam.go"
	MessagesFile = "messages.go"
)

func SeamRulesRestatedOutside(root string, rules []string) ([]string, error) {
	var restated []string

	seam := filepath.FromSlash(SeamFile)
	within := filepath.Dir(seam)

	err := filepath.WalkDir(filepath.Join(root, within), func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		relative, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return relErr
		}
		if relative == seam || filepath.Base(relative) == MessagesFile {
			return nil
		}

		set := token.NewFileSet()
		parsed, parseErr := parser.ParseFile(set, path, nil, parser.SkipObjectResolution)
		if parseErr != nil {
			return fmt.Errorf("parsing %s: %w", relative, parseErr)
		}

		ast.Inspect(parsed, func(node ast.Node) bool {
			named, isName := node.(*ast.Ident)
			if !isName || !slices.Contains(rules, named.Name) {
				return true
			}
			restated = append(restated, fmt.Sprintf("%s:%d names %s", relative, set.Position(named.Pos()).Line, named.Name))
			return true
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	slices.Sort(restated)
	return restated, nil
}

const VendorsHeading = "### Which CLI reviews"

func VendorsDocumented(root string) ([]string, error) {
	held, err := os.ReadFile(filepath.Join(root, documentationFile))
	if err != nil {
		return nil, err
	}

	section, found := sectionUnder(string(held), VendorsHeading)
	if !found {
		return nil, nil
	}
	return everyNameTabulatedIn(section), nil
}

func everyNameTabulatedIn(section string) []string {
	var named []string
	afterASeparator := false

	for _, line := range strings.Split(section, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, tableRow) {
			afterASeparator = false
			continue
		}

		cell := strings.TrimSpace(firstCellOf(trimmed))
		if strings.HasPrefix(cell, tableSeparator) {
			afterASeparator = true
			continue
		}
		if afterASeparator {
			named = append(named, cell)
		}
	}
	return named
}

func VendorSection(root string) (string, error) {
	held, err := os.ReadFile(filepath.Join(root, documentationFile))
	if err != nil {
		return "", err
	}

	section, found := sectionUnder(string(held), VendorsHeading)
	if !found {
		return "", nil
	}
	return section, nil
}

func VendorDisagreements(stated, elsewhere []string, where string) []string {
	var disagreeing []string

	for _, one := range elsewhere {
		if !slices.Contains(stated, one) {
			disagreeing = append(disagreeing, fmt.Sprintf("%s names %s, and the statement does not", where, one))
		}
	}
	for _, one := range stated {
		if !slices.Contains(elsewhere, one) {
			disagreeing = append(disagreeing, fmt.Sprintf("the statement names %s, and %s does not", one, where))
		}
	}

	slices.Sort(disagreeing)
	return disagreeing
}
