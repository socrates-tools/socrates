package gitrepo_test

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"

	"github.com/rui-ktei/socrates/internal/testing/gitrepo"
)

func TestGitStartsNoHousekeepingOfItsOwnInATestRepository(t *testing.T) {
	repository := gitrepo.New(t)
	repository.Write("reviewed.txt", "something to review\n")

	trace := traceOfCommitting(t, repository)

	for _, housekeeping := range []string{"maintenance", "gc --auto"} {
		if strings.Contains(trace, housekeeping) {
			t.Errorf("committing started %q, which writes in the repository while a test is watching what socrates writes there:\n%s", housekeeping, trace)
		}
	}
}

func traceOfCommitting(t *testing.T, repository *gitrepo.Repository) string {
	t.Helper()

	repository.Git("add", "--all")

	command := exec.Command("git", "-C", repository.Root, "commit", "--quiet", "--message", "a commit to watch")
	command.Env = append(command.Environ(), "GIT_TRACE=1")
	var traced bytes.Buffer
	command.Stderr = &traced
	if err := command.Run(); err != nil {
		t.Fatalf("committing under a trace: %v: %s", err, traced.String())
	}
	return traced.String()
}
