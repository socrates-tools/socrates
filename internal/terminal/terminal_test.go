package terminal_test

import (
	"bytes"
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates-review/internal/terminal"
)

var lines = []string{"09:14:02  the run started", "09:15:33  the run ended as clean"}

func TestPipingWorks(t *testing.T) {
	var out, failures bytes.Buffer
	piped := terminal.Output{To: &out, Failures: &failures, Terminal: false, Pager: []string{"sed", "s/^/PAGED /"}}

	if err := piped.Write(context.Background(), lines); err != nil {
		t.Fatalf("writing to a pipe: %v", err)
	}
	if strings.Contains(out.String(), "PAGED") {
		t.Errorf("a pager was started for output that is not a terminal:\n%s", out.String())
	}
	for _, line := range lines {
		if !strings.Contains(out.String(), line) {
			t.Errorf("the plain output drops %q:\n%s", line, out.String())
		}
	}
}

func TestAPagerThatCannotBeStarted(t *testing.T) {
	var out, failures bytes.Buffer
	missing := filepath.Join(t.TempDir(), "no-such-pager")
	attached := terminal.Output{To: &out, Failures: &failures, Terminal: true, Pager: []string{missing}}

	if err := attached.Write(context.Background(), lines); err != nil {
		t.Fatalf("a pager that cannot be started was not fallen back from: %v", err)
	}
	for _, line := range lines {
		if !strings.Contains(out.String(), line) {
			t.Errorf("the output drops %q when the pager could not be started:\n%s", line, out.String())
		}
	}
}

func TestAConfiguredPagerIsUsedWhereTheOutputIsATerminal(t *testing.T) {
	var out, failures bytes.Buffer
	attached := terminal.Output{To: &out, Failures: &failures, Terminal: true, Pager: []string{"sed", "s/^/PAGED /"}}

	if err := attached.Write(context.Background(), lines); err != nil {
		t.Fatalf("paging: %v", err)
	}
	if !strings.Contains(out.String(), "PAGED") {
		t.Errorf("the configured pager was not used:\n%s", out.String())
	}
}

func TestWhatTheEnvironmentSaysAboutPaging(t *testing.T) {
	t.Setenv("PAGER", "less -R")

	var out bytes.Buffer
	asked := terminal.For(&out, &out)
	if len(asked.Pager) != 2 || asked.Pager[0] != "less" {
		t.Errorf("the pager reads as %v, want what PAGER says", asked.Pager)
	}
	if asked.Terminal {
		t.Error("a buffer reads as a terminal")
	}

	t.Setenv("PAGER", "")
	if pager := terminal.For(&out, &out).Pager; len(pager) != 0 {
		t.Errorf("no configured pager reads as %v", pager)
	}
}
