package convention_test

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/internal/convention"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

func TestEveryRuleOfTheSurfacesStylesheetStylesAPieceTheGlossaryNames(t *testing.T) {
	found, err := convention.PiecesStyledBySurface(source.Root(t), convention.SurfaceStyleDirectory, convention.PiecesOfTheSurface)
	if err != nil {
		t.Fatalf("reading what the surface styles: %v", err)
	}

	for _, styled := range found {
		t.Error(styled)
	}
}

func TestARuleStylingWhatNoPieceOfTheGlossaryOwnsIsReportedNamingTheRule(t *testing.T) {
	root := t.TempDir()
	writeNested(t, filepath.Join(root, "agora", "CONTEXT.md"), "- **decision**: one card.\n- **composer**: where a person sends.\n")
	writeNested(t, filepath.Join(root, "surface", "page.css"), `html {
  margin: 0;
}

.decision,
.decision-heading {
  gap: 8px;
}

.decisions {
  gap: 8px;
}

button {
  height: 36px;
}

.send,
.copy:hover {
  order: 1;
}

.panel .sends > input {
  flex: 1;
}

.composer:not(.mystery) {
  order: -1;
}

.index-group > ui-overline {
  display: block;
}

.question:has(.decision[data-state='closed']) {
  color: var(--fg-2);
}

@media (max-width: 44rem) {
  .composer,
  .tooltip {
    padding: 16px;
  }
}

@keyframes settle {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.composer {
  transition: color 0.15s, background url(x.png);
  font-family:
    var(--font-body),
    var(--font-mono);
}
`)
	pieces := map[string][]string{
		"decision": {"decision", "question"},
		"composer": {"composer", "send"},
		"receipt":  {"receipt"},
	}

	found, err := convention.PiecesStyledBySurface(root, "surface", pieces)
	if err != nil {
		t.Fatal(err)
	}

	want := []string{
		`surface/page.css:10 styles ".decisions", and ".decisions" belongs to no piece the glossary names`,
		`surface/page.css:14 styles "button", and "button" belongs to no piece the glossary names`,
		`surface/page.css:19 styles ".copy:hover", and ".copy" belongs to no piece the glossary names`,
		`surface/page.css:23 styles ".panel .sends > input", and ".panel" belongs to no piece the glossary names`,
		`surface/page.css:27 styles ".composer:not(.mystery)", and ".mystery" belongs to no piece the glossary names`,
		`surface/page.css:31 styles ".index-group > ui-overline", and ".index-group" belongs to no piece the glossary names`,
		`surface/page.css:41 styles ".tooltip", and ".tooltip" belongs to no piece the glossary names`,
		`agora/CONTEXT.md does not name the piece "receipt"`,
	}
	if diff := cmp.Diff(want, found); diff != "" {
		t.Errorf("what the surface styles beyond the glossary differs (-want +got):\n%s", diff)
	}
}

func TestAPrimitiveNamedByItsTagIsReportedWhateverClassItCarries(t *testing.T) {
	root := t.TempDir()
	writeNested(t, filepath.Join(root, "agora", "CONTEXT.md"), "- **composer**: where a person sends.\n")
	writeNested(t, filepath.Join(root, "surface", "page.css"), `.composer ui-overline {
  display: block;
}

.composer > textarea,
.composer dialog.composer {
  flex: 1;
}

.composer a {
  color: inherit;
}
`)

	found, err := convention.PiecesStyledBySurface(root, "surface", map[string][]string{"composer": {"composer"}})
	if err != nil {
		t.Fatal(err)
	}

	want := []string{
		`surface/page.css:1 styles ".composer ui-overline", and "ui-overline" belongs to no piece the glossary names`,
		`surface/page.css:5 styles ".composer > textarea", and "textarea" belongs to no piece the glossary names`,
		`surface/page.css:6 styles ".composer dialog.composer", and "dialog" belongs to no piece the glossary names`,
	}
	if diff := cmp.Diff(want, found); diff != "" {
		t.Errorf("what the surface styles by tag differs (-want +got):\n%s", diff)
	}
}
