package convention

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

func TestSourceHoldsNoCommentary(t *testing.T) {
	found, err := Comments(source.Root(t))
	if err != nil {
		t.Fatalf("reading the source: %v", err)
	}

	for _, commentary := range found {
		t.Errorf("%s:%d holds a comment: %s", commentary.File, commentary.Line, commentary.Text)
	}
}

func TestCommentaryReporting(t *testing.T) {
	cases := []struct {
		name   string
		source string
		want   []Commentary
	}{
		{
			name:   "a line comment is reported",
			source: "package example\n\nfunc Example() {}\n\n// a note nobody asked for\n",
			want:   []Commentary{{File: "source.go", Line: 5, Text: "// a note nobody asked for"}},
		},
		{
			name:   "a build directive is not reported",
			source: "//go:build linux\n\npackage example\n\nfunc Example() {}\n",
		},
		{
			name:   "a generate directive is not reported",
			source: "package example\n\n//go:generate true\nfunc Example() {}\n",
		},
		{
			name:   "a scenario marker is not reported",
			source: "package example\n\nfunc Example() {\n\t// scenario: Runs sort by when they started\n}\n",
		},
		{
			name:   "a block comment is reported at its opening line",
			source: "package example\n\n/* first\nsecond */\nfunc Example() {}\n",
			want:   []Commentary{{File: "source.go", Line: 3, Text: "/* first"}},
		},
		{
			name:   "commentary inside a string literal is not reported",
			source: "package example\n\nvar Note = \"// not a comment\"\n",
		},
		{
			name:   "a trailing comment on a statement is reported",
			source: "package example\n\nvar Note = \"\" // why\n",
			want:   []Commentary{{File: "source.go", Line: 3, Text: "// why"}},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			root := t.TempDir()
			if err := os.WriteFile(filepath.Join(root, "source.go"), []byte(testCase.source), 0o644); err != nil {
				t.Fatalf("writing the source: %v", err)
			}

			found, err := Comments(root)
			if err != nil {
				t.Fatalf("reading the source: %v", err)
			}

			if diff := cmp.Diff(testCase.want, found); diff != "" {
				t.Errorf("commentary differs (-want +got):\n%s", diff)
			}
		})
	}
}

func TestSkippedDirectoriesAreNotRead(t *testing.T) {
	root := t.TempDir()
	for _, directory := range skippedDirectories {
		nested := filepath.Join(root, directory)
		if err := os.MkdirAll(nested, 0o755); err != nil {
			t.Fatalf("creating %s: %v", directory, err)
		}
		source := "package example\n\n// a note in " + directory + "\n"
		if err := os.WriteFile(filepath.Join(nested, "source.go"), []byte(source), 0o644); err != nil {
			t.Fatalf("writing into %s: %v", directory, err)
		}
	}

	found, err := Comments(root)
	if err != nil {
		t.Fatalf("reading the source: %v", err)
	}
	if len(found) != 0 {
		t.Errorf("read commentary from a skipped directory: %v", found)
	}
}
