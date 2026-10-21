package convention_test

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/internal/convention"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

const scratchTheme = `@import 'tailwindcss' source(none);

@source '../**/*.ts';
@source not '../**/node_modules';
@source '../../../render/**/*.gohtml';
@source '../../../../site/index.html';
`

func scratchClientRoot(t *testing.T) string {
	t.Helper()

	root := t.TempDir()
	writeNested(t, filepath.Join(root, convention.ClassSourcesStatement), scratchTheme)
	writeNested(t, filepath.Join(root, "agora", "web", "packages", "ui", "index.ts"), "export {}\n")
	writeNested(t, filepath.Join(root, "agora", "render", "page.gohtml"), "<html></html>\n")
	writeNested(t, filepath.Join(root, "site", "index.html"), "<html></html>\n")
	return root
}

func TestNoClassStringOfTheClientCarriesAnArbitraryValue(t *testing.T) {
	found, err := convention.ArbitraryValuesWrittenInClasses(source.Root(t), convention.ClassSourcesStatement)
	if err != nil {
		t.Fatalf("reading what the client writes as classes: %v", err)
	}

	for _, written := range found {
		t.Error(written)
	}
}

func TestTheClassesAreReadFromTheTypeScriptThePartialsAndTheLandingPage(t *testing.T) {
	sources, err := convention.ClassSourcesStated(source.Root(t), convention.ClassSourcesStatement)
	if err != nil {
		t.Fatalf("reading where classes are written: %v", err)
	}

	hasSuffix := func(suffix string) bool {
		return slices.ContainsFunc(sources, func(read string) bool { return strings.HasSuffix(read, suffix) })
	}
	for _, expected := range []string{".ts", "partials/ui-button.gohtml", "site/index.html"} {
		if !hasSuffix(expected) {
			t.Errorf("the @source lines of %s reach no %s", filepath.ToSlash(convention.ClassSourcesStatement), expected)
		}
	}
	for _, read := range sources {
		if strings.Contains(read, "node_modules") {
			t.Errorf("%s is read though the @source lines exclude node_modules", read)
		}
	}
}

func TestAnArbitraryValueInATypeScriptFileAPartialOrTheLandingPageIsReportedNamingTheFileAndTheLine(t *testing.T) {
	root := scratchClientRoot(t)
	writeNested(t, filepath.Join(root, "agora", "web", "packages", "ui", "badge", "badge.ts"), "const clean = 'flex gap-2'\n\nconst loose = 'flex p-[13px]'\n")
	writeNested(t, filepath.Join(root, "agora", "render", "partials", "ui-card.gohtml"), "<ui-card class=\"flex\">\n  <div class=\"bg-[#fff] p-2\"></div>\n</ui-card>\n")
	writeNested(t, filepath.Join(root, "site", "index.html"), "<html>\n<body>\n<p class=\"text-[14px]\">hello</p>\n</body>\n</html>\n")

	found, err := convention.ArbitraryValuesWrittenInClasses(root, convention.ClassSourcesStatement)
	if err != nil {
		t.Fatalf("reading what the client writes as classes: %v", err)
	}

	var located []string
	for _, written := range found {
		located = append(located, strings.SplitN(written, " ", 2)[0])
	}
	expected := []string{
		"agora/render/partials/ui-card.gohtml:2",
		"agora/web/packages/ui/badge/badge.ts:3",
		"site/index.html:3",
	}
	if diff := cmp.Diff(expected, located); diff != "" {
		t.Errorf("located arbitrary values (-want +got):\n%s", diff)
	}
}

func TestEveryFormOfArbitraryValueIsReported(t *testing.T) {
	forms := []string{
		"p-[13px]",
		"bg-[#fff]",
		"grid-cols-[1fr_2fr]",
		"hover:bg-[var(--accent)]",
		"-mt-[2px]",
		"w-[calc(100%-2rem)]",
		"data-[state=open]:p-2",
		"has-[:checked]:bg-accent",
		"bg-accent/[0.5]",
		"[&_p]:text-fg-2",
		"[&>*]:p-2",
		"hover:[&_p]:flex",
		"[mask-type:luminance]",
		"[--gap:3px]",
		"@[300px]:flex",
		"bg-(--accent)",
		"p-(length:--gap)",
	}
	root := scratchClientRoot(t)
	var lines []string
	for _, form := range forms {
		lines = append(lines, "const c = '"+form+"'")
	}
	writeNested(t, filepath.Join(root, "agora", "web", "packages", "ui", "forms.ts"), strings.Join(lines, "\n")+"\n")

	found, err := convention.ArbitraryValuesWrittenInClasses(root, convention.ClassSourcesStatement)
	if err != nil {
		t.Fatalf("reading what the client writes as classes: %v", err)
	}

	for number, form := range forms {
		wanted := fmt.Sprintf("agora/web/packages/ui/forms.ts:%d ", number+1)
		if !slices.ContainsFunc(found, func(written string) bool {
			return strings.HasPrefix(written, wanted) && strings.Contains(written, "\""+form+"\"")
		}) {
			t.Errorf("%q was not reported on line %d, got %v", form, number+1, found)
		}
	}
}

func TestASquareBracketThatIsNotAClassStringFailsNothing(t *testing.T) {
	root := scratchClientRoot(t)
	writeNested(t, filepath.Join(root, "agora", "web", "packages", "surface", "page.ts"), `const buttons = document.querySelectorAll('[data-state="open"]')
const first = buttons[0]
const previous = rows[index - 1]
const names: string[] = []
const pair: Record<string, [number, number]> = {}
const link = anchor.querySelector('a[href]:not([hidden])')
const dialog = document.querySelector('ui-dialog[open]')
const selector = '[data-action=send]:hover'
const matrix = [[1, 2], [3, 4]]
const mixed = items[count-1]
const expression = /[a-z-]+/g
const classes = 'flex gap-2 p-4 md:p-6 hover:bg-accent'
`)
	writeNested(t, filepath.Join(root, "agora", "render", "partials", "ui-field.gohtml"), `<ui-field data-x="{{ index .Items 0 }}" class="flex">
  {{ range $i, $item := .Items }}{{ index $.Items $i }}{{ end }}
</ui-field>
`)

	found, err := convention.ArbitraryValuesWrittenInClasses(root, convention.ClassSourcesStatement)
	if err != nil {
		t.Fatalf("reading what the client writes as classes: %v", err)
	}

	for _, written := range found {
		t.Error(written)
	}
}

func TestAnArbitraryValueOutsideThePathsTheThemeStatesIsNotRead(t *testing.T) {
	root := t.TempDir()
	writeNested(t, filepath.Join(root, convention.ClassSourcesStatement), "@source 'client/**/*.ts';\n@source not 'client/**/vendor';\n")
	writeNested(t, filepath.Join(root, "agora", "web", "packages", "ui", "client", "inside.ts"), "const loose = 'p-[1px]'\n")
	writeNested(t, filepath.Join(root, "agora", "web", "packages", "ui", "client", "vendor", "library.ts"), "const loose = 'p-[2px]'\n")
	writeNested(t, filepath.Join(root, "agora", "web", "packages", "ui", "elsewhere", "outside.ts"), "const loose = 'p-[3px]'\n")

	found, err := convention.ArbitraryValuesWrittenInClasses(root, convention.ClassSourcesStatement)
	if err != nil {
		t.Fatalf("reading what the client writes as classes: %v", err)
	}

	if len(found) != 1 || !strings.HasPrefix(found[0], "agora/web/packages/ui/client/inside.ts:1 ") {
		t.Errorf("expected only client/inside.ts to be read, got %v", found)
	}
}

func TestAThemeStatingNoSourceFailsRatherThanReadingNothing(t *testing.T) {
	root := t.TempDir()
	writeNested(t, filepath.Join(root, convention.ClassSourcesStatement), "@import 'tailwindcss';\n")

	if _, err := convention.ArbitraryValuesWrittenInClasses(root, convention.ClassSourcesStatement); err == nil {
		t.Error("a theme stating no @source was read as holding nothing")
	}
}

func TestAStatedSourceThatIsNotThereFailsRatherThanReadingNothing(t *testing.T) {
	root := t.TempDir()
	writeNested(t, filepath.Join(root, convention.ClassSourcesStatement), "@source '../../../../site/index.html';\n")

	if _, err := convention.ArbitraryValuesWrittenInClasses(root, convention.ClassSourcesStatement); err == nil {
		t.Error("a stated source that does not exist was read as holding nothing")
	}
}
