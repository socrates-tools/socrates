package e2e_test

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/ktei/socrates/internal/convention"
	"github.com/ktei/socrates/internal/testing/source"
)

func TestHelpNamesEveryTargetAContributorIsToldToRun(t *testing.T) {
	root := requireMake(t)

	told, err := convention.TargetsDocumentsTell(root)
	if err != nil {
		t.Fatalf("reading what the documents tell a contributor to run: %v", err)
	}
	if len(told) == 0 {
		t.Fatal("no document names a target to run, so this check has nothing to compare and would pass forever")
	}

	printed := madeSay(t, root, "help")
	for _, target := range told {
		if !strings.Contains(printed, "make "+target) {
			t.Errorf("a document tells a contributor to run make %s and make help does not name it:\n%s", target, printed)
		}
	}
}

func madeSay(t *testing.T, root string, target string) string {
	t.Helper()

	command := exec.Command("make", target)
	command.Dir = root
	command.Env = append(os.Environ(), source.Toolchain()...)

	var out, failures bytes.Buffer
	command.Stdout = &out
	command.Stderr = &failures
	if err := command.Run(); err != nil {
		t.Fatalf("make %s: %v: %s", target, err, failures.String())
	}
	return out.String()
}
