package convention_test

import (
	"path/filepath"
	"testing"

	"github.com/socrates-tools/socrates/internal/convention"
)

func nestedCheckout(t *testing.T, root, relative string) string {
	t.Helper()

	checkout := filepath.Join(root, relative)
	mkdir(t, checkout)
	write(t, filepath.Join(checkout, ".git"), "gitdir: /elsewhere/.git/worktrees/one\n")
	return checkout
}

func TestACheckoutNestedUnderTheRootIsNotTheRepositorysOwnSource(t *testing.T) {
	root := t.TempDir()

	mkdir(t, filepath.Join(root, "review", "command"))
	write(t, filepath.Join(root, "review", "command", "root.go"), "package command\n\nfunc NewRootCommand() {}\n")

	checkout := nestedCheckout(t, root, filepath.Join(".claude", "worktrees", "one"))
	mkdir(t, filepath.Join(checkout, "review", "command"))
	write(t, filepath.Join(checkout, "review", "command", "root.go"), "package command\n\n// a comment the conventions forbid\nfunc NewRootCommand() {}\n")

	found, err := convention.Comments(root)
	if err != nil {
		t.Fatalf("reading a tree holding a nested checkout: %v", err)
	}
	for _, commentary := range found {
		t.Errorf("%s is inside a checkout nested under the root, and the conventions are checked over this repository's own source rather than over a copy of it: %s", commentary.File, commentary.Text)
	}
}

func TestAStringInsideANestedCheckoutIsNotReadForTheRegistry(t *testing.T) {
	root := t.TempDir()

	checkout := nestedCheckout(t, root, filepath.Join(".claude", "worktrees", "one"))
	mkdir(t, filepath.Join(checkout, "review", "command"))
	write(t, filepath.Join(checkout, "review", "command", "root.go"),
		"package command\n\nimport \"fmt\"\n\nfunc Refuse() error {\n\treturn fmt.Errorf(\"a string no registry holds\")\n}\n")

	wording, err := convention.UserFacingStrings(root)
	if err != nil {
		t.Fatalf("reading a tree holding a nested checkout: %v", err)
	}
	for _, inline := range wording.Inline {
		t.Errorf("%s is inside a checkout nested under the root, and a copy of the repository is not the repository's own source", inline)
	}
}

func TestTheRootsOwnSourceIsStillReadWhenACheckoutIsNestedUnderIt(t *testing.T) {
	root := t.TempDir()

	mkdir(t, filepath.Join(root, "review", "command"))
	write(t, filepath.Join(root, "review", "command", "root.go"), "package command\n\n// a comment the conventions forbid\nfunc NewRootCommand() {}\n")

	nestedCheckout(t, root, filepath.Join(".claude", "worktrees", "one"))

	found, err := convention.Comments(root)
	if err != nil {
		t.Fatalf("reading a tree holding a nested checkout: %v", err)
	}
	if len(found) != 1 {
		t.Fatalf("the root's own source reports %d comments beside a nested checkout, and skipping the checkout must not skip the repository: %v", len(found), found)
	}
}
