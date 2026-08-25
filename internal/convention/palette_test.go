package convention_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rui-ktei/socrates/internal/convention"
	"github.com/rui-ktei/socrates/internal/testing/source"
)

func TestEveryColourIsStatedOnce(t *testing.T) {
	read, err := convention.ColoursWrittenWhereTheyAreUsed(source.Root(t), convention.StyleDirectory)
	if err != nil {
		t.Fatalf("reading what the page paints: %v", err)
	}

	for _, loose := range read.Loose {
		t.Errorf("%s, and a colour written where it is used survives neither the other preference nor the next reader", loose)
	}
}

func TestAColourWrittenWhereItIsUsedIsReported(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "styles"), 0o755); err != nil {
		t.Fatalf("making a fabricated tree: %v", err)
	}
	write(t, filepath.Join(root, "styles", "held.css"), `:root {
  --ink: #1b2529;
}

.question {
  color: var(--ink);
  border-color: #b3bbb8;
}

.severity {
  --shouted: rgb(255, 0, 0);
  background: var(--shouted);
}
`)

	read, err := convention.ColoursWrittenWhereTheyAreUsed(root, "styles")
	if err != nil {
		t.Fatalf("reading what a fabricated stylesheet paints: %v", err)
	}

	want := []string{
		`styles/held.css:11 states "rgb" outside the block that states the palette`,
		`styles/held.css:7 paints "#b3bbb8" where it is used`,
	}
	if diff := cmp.Diff(want, read.Loose); diff != "" {
		t.Errorf("what a stylesheet painting its own colours reports differs (-want +got):\n%s", diff)
	}
}
