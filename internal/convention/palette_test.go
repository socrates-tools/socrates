package convention_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/internal/convention"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

func TestEveryColourAndFontFamilyIsStatedOnceInTheLibrarysTokens(t *testing.T) {
	read, err := convention.PaintWrittenWhereItIsUsed(source.Root(t), convention.StyleDirectory, convention.TokensDirectory)
	if err != nil {
		t.Fatalf("reading what the page paints: %v", err)
	}

	for _, loose := range read.Colours {
		t.Errorf("%s, and a colour written where it is used survives neither the other preference nor the next reader", loose)
	}
	for _, loose := range read.Families {
		t.Errorf("%s, and a font family written where it is used survives neither the other preference nor the next reader", loose)
	}
}

func TestTheLibrarysTokensAreWhereThePaletteIsStated(t *testing.T) {
	tokens := filepath.Join(source.Root(t), convention.TokensDirectory)
	for _, stated := range []string{"tokens.css", "fonts.css"} {
		if _, err := os.Stat(filepath.Join(tokens, stated)); err != nil {
			t.Errorf("the library's tokens hold no %s: %v", stated, err)
		}
	}
}

func TestAColourWrittenWhereItIsUsedIsReported(t *testing.T) {
	root := t.TempDir()
	writeNested(t, filepath.Join(root, "styles", "tokens", "held.css"), `:root {
  --ink: #1b2529;
  --clay-950: #0e0e0d;
}

.stray {
  border-color: #aabbcc;
}
`)
	writeNested(t, filepath.Join(root, "styles", "packages", "question.css"), `:root {
  --own-ink: #1b2529;
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

	read, err := convention.PaintWrittenWhereItIsUsed(root, "styles", filepath.Join("styles", "tokens"))
	if err != nil {
		t.Fatalf("reading what a fabricated stylesheet paints: %v", err)
	}

	want := []string{
		`styles/packages/question.css:11 states "rgb" outside the library's tokens`,
		`styles/packages/question.css:2 states "#1b2529" outside the library's tokens`,
		`styles/packages/question.css:7 states "#b3bbb8" outside the library's tokens`,
		`styles/tokens/held.css:7 paints "#aabbcc" where it is used`,
	}
	if diff := cmp.Diff(want, read.Colours); diff != "" {
		t.Errorf("what a stylesheet painting its own colours reports differs (-want +got):\n%s", diff)
	}
}

func TestAFontFamilyNamedOutsideTheLibrarysTokensIsReported(t *testing.T) {
	root := t.TempDir()
	writeNested(t, filepath.Join(root, "styles", "tokens", "fonts.css"), `@font-face {
  font-family: 'Manrope';
  src: url('./fonts/Manrope-Variable.woff2') format('woff2');
}

:root {
  --font-body: 'Manrope', ui-sans-serif, sans-serif;
}
`)
	writeNested(t, filepath.Join(root, "styles", "packages", "card.css"), `@font-face {
  font-family: 'Own';
  src: url('./own.woff2') format('woff2');
}

.heading {
  font-family: var(--font-body);
}

.code {
  font-family: 'Courier New', monospace;
}

.reset {
  font-family: inherit;
  font: inherit;
}

.shorthand {
  font: 12px/1.4 Georgia, serif;
}

.shorthand-by-token {
  font: 600 12px/1.4 var(--font-body);
}
`)

	read, err := convention.PaintWrittenWhereItIsUsed(root, "styles", filepath.Join("styles", "tokens"))
	if err != nil {
		t.Fatalf("reading what a fabricated stylesheet names: %v", err)
	}

	want := []string{
		`styles/packages/card.css:11 names the font family 'Courier New', monospace outside the library's tokens`,
		`styles/packages/card.css:2 names the font family 'Own' outside the library's tokens`,
		`styles/packages/card.css:20 names the font family 12px/1.4 Georgia, serif outside the library's tokens`,
	}
	if diff := cmp.Diff(want, read.Families); diff != "" {
		t.Errorf("what a stylesheet naming its own font reports differs (-want +got):\n%s", diff)
	}
}
