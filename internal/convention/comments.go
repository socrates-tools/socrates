package convention

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type Commentary struct {
	File string
	Line int
	Text string
}

var exemptPrefixes = []string{
	"//go:",
	"//nolint",
	"// scenario:",
	"//scenario:",
}

var skippedDirectories = []string{".git", "bin", "node_modules", "testdata"}

func Comments(root string) ([]Commentary, error) {
	var found []Commentary

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
		if filepath.Ext(path) != ".go" {
			return nil
		}

		relative, relErr := filepath.Rel(root, path)
		if relErr != nil {
			relative = path
		}

		set := token.NewFileSet()
		parsed, parseErr := parser.ParseFile(set, path, nil, parser.ParseComments|parser.SkipObjectResolution)
		if parseErr != nil {
			return fmt.Errorf("parsing %s: %w", relative, parseErr)
		}

		for _, group := range parsed.Comments {
			for _, comment := range group.List {
				if strings.HasPrefix(comment.Text, "//") && exempt(comment.Text) {
					continue
				}
				found = append(found, Commentary{
					File: relative,
					Line: set.Position(comment.Slash).Line,
					Text: firstLine(comment.Text),
				})
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return found, nil
}

func exempt(text string) bool {
	for _, prefix := range exemptPrefixes {
		if strings.HasPrefix(text, prefix) {
			return true
		}
	}
	return false
}

func skipped(path, name string) bool {
	return slices.Contains(skippedDirectories, name) || nestedCheckout(path)
}

func nestedCheckout(path string) bool {
	_, err := os.Stat(filepath.Join(path, ".git"))
	return err == nil
}

func firstLine(text string) string {
	first, _, _ := strings.Cut(text, "\n")
	return strings.TrimSpace(first)
}
