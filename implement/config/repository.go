package config

import (
	"bytes"
	"errors"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/socrates-tools/socrates/internal/refusal"
)

const (
	worktreePrefix = "worktree "
	bareMark       = "bare"
	localBranches  = "refs/heads/"
	remoteBranches = "refs/remotes/"
)

var (
	ErrNotARepository = errors.New(notARepository)
	ErrBareRepository = errors.New(bareRepository)
)

func RepositoryAt(directory string) (string, error) {
	listed, err := git(directory, "worktree", "list", "--porcelain")
	if err != nil {
		var exited *exec.ExitError
		if errors.As(err, &exited) {
			return "", ErrNotARepository
		}
		return "", refusal.New(gitCannotRun, err)
	}

	main, _, _ := strings.Cut(listed, "\n\n")
	lines := strings.Split(main, "\n")
	if !strings.HasPrefix(lines[0], worktreePrefix) {
		return "", ErrNotARepository
	}
	for _, line := range lines[1:] {
		if strings.TrimSpace(line) == bareMark {
			return "", ErrBareRepository
		}
	}

	top := strings.TrimPrefix(lines[0], worktreePrefix)
	resolved, err := filepath.EvalSymlinks(top)
	if err != nil {
		return "", refusal.New(gitCannotRun, err)
	}
	return resolved, nil
}

func HoldsBranch(repository, branch string) bool {
	if branch == "" {
		return false
	}
	for _, namespace := range []string{localBranches, remoteBranches} {
		if _, err := git(repository, "show-ref", "--verify", "--quiet", namespace+branch); err == nil {
			return true
		}
	}
	return false
}

func RequireBranch(repository, branch string) error {
	if !HoldsBranch(repository, branch) {
		return refusal.New(repositoryHoldsNoSuchBranch, repository, branch)
	}
	return nil
}

func git(directory string, arguments ...string) (string, error) {
	command := exec.Command("git", append([]string{"-C", directory}, arguments...)...)
	var out bytes.Buffer
	command.Stdout = &out
	if err := command.Run(); err != nil {
		return "", err
	}
	return out.String(), nil
}
