package convention_test

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/internal/convention"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

func TestTheSurfacesStylesheetHoldsNoRuleAUtilityCouldSay(t *testing.T) {
	// scenario: The surface's stylesheet holds no rule that names a class, an id or an element
	found, err := convention.RulesAUtilityCouldSay(source.Root(t), convention.SurfaceStyleDirectory)
	if err != nil {
		t.Fatalf("reading what the surface styles: %v", err)
	}

	for _, styled := range found {
		t.Error(styled)
	}
}

func TestARuleNamingAClassAnIdOrAnElementIsReportedNamingTheRule(t *testing.T) {
	root := t.TempDir()
	writeNested(t, filepath.Join(root, "surface", "page.css"), `::selection {
  color: inherit;
}

html {
  margin: 0;
}

.decision,
.decision-heading {
  gap: 8px;
}

*,
*::before {
  box-sizing: border-box;
}

#settled:target {
  order: 1;
}

:focus-visible {
  outline: none;
}

[hidden] {
  display: none;
}

.panel:has(.composer[data-phase='over']) {
  display: none;
}

@media (max-width: 44rem) {
  .composer,
  button:hover {
    padding: 16px;
  }

  *::after {
    content: '';
  }
}

@keyframes settle {
  from {
    opacity: 0;
  }
}

@media (prefers-reduced-motion: reduce) {
  * {
    transition-duration: 0ms !important;
  }
}
`)

	found, err := convention.RulesAUtilityCouldSay(root, "surface")
	if err != nil {
		t.Fatal(err)
	}

	want := []string{
		`surface/page.css:5 styles "html", and a utility on it says that`,
		`surface/page.css:9 styles ".decision", and a utility on it says that`,
		`surface/page.css:10 styles ".decision-heading", and a utility on it says that`,
		`surface/page.css:19 styles "#settled:target", and a utility on it says that`,
		`surface/page.css:31 styles ".panel:has(.composer[data-phase='over'])", and a utility on it says that`,
		`surface/page.css:36 styles ".composer", and a utility on it says that`,
		`surface/page.css:37 styles "button:hover", and a utility on it says that`,
	}
	if diff := cmp.Diff(want, found); diff != "" {
		t.Errorf("what the surface styles that a utility could say differs (-want +got):\n%s", diff)
	}
}
