package snapshot

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const commitMessage = "socrates snapshot"

const worktreeDirectory = "worktree"

const NothingToReviewReason = "the working tree was identical to HEAD, so there was nothing to review"

func WorktreePath(runDirectory string) string {
	return filepath.Join(runDirectory, worktreeDirectory)
}

var commitIdentity = []string{
	"GIT_AUTHOR_NAME=socrates",
	"GIT_AUTHOR_EMAIL=socrates@localhost",
	"GIT_COMMITTER_NAME=socrates",
	"GIT_COMMITTER_EMAIL=socrates@localhost",
}

type Refusal struct {
	Directory string
	Reason    string
}

func (r *Refusal) Error() string {
	return r.Reason
}

type Snapshot struct {
	RepositoryRoot string
	Base           string
	Tree           string
	Commit         string
	Changed        []string
}

func (s Snapshot) NothingToReview() bool {
	return len(s.Changed) == 0
}

type Worktree struct {
	RepositoryRoot string
	Path           string
	Commit         string
}

func Take(ctx context.Context, directory string) (Snapshot, error) {
	root, base, err := resolve(ctx, directory)
	if err != nil {
		return Snapshot{}, err
	}

	baseTree, err := git(ctx, root, nil, "rev-parse", "--verify", base+"^{tree}")
	if err != nil {
		return Snapshot{}, err
	}

	indexDirectory, err := os.MkdirTemp("", "socrates-index-")
	if err != nil {
		return Snapshot{}, fmt.Errorf("creating a directory for the snapshot's own index: %w", err)
	}
	defer os.RemoveAll(indexDirectory)

	ownIndex := []string{"GIT_INDEX_FILE=" + filepath.Join(indexDirectory, "index")}

	if _, err := git(ctx, root, ownIndex, "read-tree", base); err != nil {
		return Snapshot{}, err
	}
	if _, err := git(ctx, root, ownIndex, "add", "--all"); err != nil {
		return Snapshot{}, err
	}
	tree, err := git(ctx, root, ownIndex, "write-tree")
	if err != nil {
		return Snapshot{}, err
	}

	taken := Snapshot{RepositoryRoot: root, Base: base, Tree: tree}
	if tree == baseTree {
		return taken, nil
	}

	commit, err := git(ctx, root, commitIdentity, "commit-tree", tree, "-p", base, "-m", commitMessage)
	if err != nil {
		return Snapshot{}, err
	}
	changed, err := git(ctx, root, nil, "diff-tree", "-r", "--name-only", "--no-commit-id", base, commit)
	if err != nil {
		return Snapshot{}, err
	}

	taken.Commit = commit
	taken.Changed = lines(changed)
	return taken, nil
}

func (s Snapshot) Checkout(ctx context.Context, path string) (Worktree, error) {
	if s.NothingToReview() {
		return Worktree{}, errors.New("there is nothing to review, so there is no snapshot to check out")
	}

	_, statErr := os.Stat(path)
	presentAlready := statErr == nil

	if _, err := git(ctx, s.RepositoryRoot, nil, "worktree", "add", "--detach", path, s.Commit); err != nil {
		if !presentAlready {
			if removeErr := os.RemoveAll(path); removeErr != nil {
				return Worktree{}, errors.Join(err, fmt.Errorf("removing %s: %w", path, removeErr))
			}
		}
		if _, pruneErr := git(ctx, s.RepositoryRoot, nil, "worktree", "prune"); pruneErr != nil {
			return Worktree{}, errors.Join(err, pruneErr)
		}
		return Worktree{}, err
	}

	return Worktree{RepositoryRoot: s.RepositoryRoot, Path: path, Commit: s.Commit}, nil
}

func (w Worktree) Remove(ctx context.Context) error {
	if err := os.RemoveAll(w.Path); err != nil {
		return fmt.Errorf("removing %s: %w", w.Path, err)
	}
	if _, err := git(ctx, w.RepositoryRoot, nil, "worktree", "prune"); err != nil {
		return err
	}
	return nil
}

func resolve(ctx context.Context, directory string) (string, string, error) {
	bare, err := git(ctx, directory, nil, "rev-parse", "--is-bare-repository")
	if err != nil {
		return "", "", &Refusal{
			Directory: directory,
			Reason:    fmt.Sprintf("%s is not inside a git working tree, so there is nothing to review: run scrat review from inside a repository", directory),
		}
	}
	if bare == "true" {
		return "", "", &Refusal{
			Directory: directory,
			Reason:    fmt.Sprintf("the repository at %s is bare, so it has no working tree to review", directory),
		}
	}

	root, err := git(ctx, directory, nil, "rev-parse", "--show-toplevel")
	if err != nil || root == "" {
		return "", "", &Refusal{
			Directory: directory,
			Reason:    fmt.Sprintf("%s is not inside a git working tree, so there is nothing to review: run scrat review from inside a repository", directory),
		}
	}

	base, err := git(ctx, root, nil, "rev-parse", "--verify", "HEAD^{commit}")
	if err != nil {
		return "", "", &Refusal{
			Directory: directory,
			Reason:    fmt.Sprintf("the repository at %s has no commit for HEAD to resolve to, so there is no base to review against: commit something first", root),
		}
	}

	return root, base, nil
}

func git(ctx context.Context, directory string, environment []string, args ...string) (string, error) {
	command := exec.CommandContext(ctx, "git", args...)
	command.Dir = directory
	if len(environment) > 0 {
		command.Env = append(os.Environ(), environment...)
	}

	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	if err := command.Run(); err != nil {
		reported := strings.TrimSpace(stderr.String())
		if reported == "" {
			reported = err.Error()
		}
		return "", fmt.Errorf("git %s in %s: %s", strings.Join(args, " "), directory, reported)
	}

	return strings.TrimSpace(stdout.String()), nil
}

func lines(output string) []string {
	if output == "" {
		return nil
	}
	return strings.Split(output, "\n")
}
