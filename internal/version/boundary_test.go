package version_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/testing/source"
)

var recordedWhenLinked = []string{"version", "commit", "uncommitted"}

func TestTheReportedBuildIsReadFromNothing(t *testing.T) {
	forbidden := []string{
		"io/fs",
		"log",
		"os",
		"os/exec",
		"path/filepath",
		"syscall",
		"time",
	}

	for _, imported := range source.Imports(t, source.Root(t), filepath.Join("internal", "version")) {
		reaching := slices.Contains(forbidden, imported) ||
			strings.HasPrefix(imported, "os/") ||
			strings.HasPrefix(imported, "net") ||
			strings.HasPrefix(imported, "golang.org/x/sys")
		if reaching {
			t.Errorf("the package that says which build this is imports %s, so what it reports could come from the repository the command is run in", imported)
		}
	}
}

// scenario: The version is read where the binary was linked, not where it runs
func TestEveryFactAboutTheBuildIsSetByTheLinkerAlone(t *testing.T) {
	initialised := map[string]bool{}

	for _, path := range source.Files(t, filepath.Join(source.Root(t), "internal", "version")) {
		set := token.NewFileSet()
		parsed, err := parser.ParseFile(set, path, nil, parser.SkipObjectResolution)
		if err != nil {
			t.Fatalf("reading %s: %v", path, err)
		}

		for _, declaration := range parsed.Decls {
			general, isGeneral := declaration.(*ast.GenDecl)
			if !isGeneral || general.Tok != token.VAR {
				continue
			}
			for _, specification := range general.Specs {
				valued, isValued := specification.(*ast.ValueSpec)
				if !isValued {
					continue
				}
				for _, named := range valued.Names {
					initialised[named.Name] = len(valued.Values) > 0
				}
			}
		}
	}

	for _, fact := range recordedWhenLinked {
		declared, is := initialised[fact]
		if !is {
			t.Errorf("%s is not a package variable of the one that says which build this is, so the linker has nothing to set", fact)
			continue
		}
		if declared {
			t.Errorf("%s is given a value in the source, so what a binary says about itself is written where it was compiled rather than recorded where it was linked", fact)
		}
	}
}
