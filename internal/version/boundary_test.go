package version_test

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/socrates-devtool/socrates/internal/testing/source"
)

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
