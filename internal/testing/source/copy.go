package source

import (
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func CopyOfTree(t *testing.T) string {
	t.Helper()

	root := Root(t)
	copied := t.TempDir()
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if slices.Contains([]string{".git", "bin", "node_modules"}, entry.Name()) {
				return filepath.SkipDir
			}
			return os.MkdirAll(filepath.Join(copied, relative), 0o755)
		}
		if !entry.Type().IsRegular() {
			return nil
		}
		held, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(copied, relative), held, 0o644)
	})
	if err != nil {
		t.Fatalf("copying the source tree: %v", err)
	}
	return copied
}

func BuildIn(t *testing.T, root, target, output string, tags ...string) {
	t.Helper()

	if err := build(root, target, output, tags, ""); err != nil {
		t.Fatalf("building %s in %s: %v", target, root, err)
	}
}
