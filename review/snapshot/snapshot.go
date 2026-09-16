package snapshot

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/socrates-tools/socrates/review/hunks"
)

const commitMessage = "socrates snapshot"

const worktreeDirectory = "worktree"

func WorktreePath(runDirectory string) string {
	return filepath.Join(runDirectory, worktreeDirectory)
}

var commitIdentity = []string{
	"GIT_AUTHOR_NAME=socrates",
	"GIT_AUTHOR_EMAIL=socrates@localhost",
	"GIT_AUTHOR_DATE=@0 +0000",
	"GIT_COMMITTER_NAME=socrates",
	"GIT_COMMITTER_EMAIL=socrates@localhost",
	"GIT_COMMITTER_DATE=@0 +0000",
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

func Take(ctx context.Context, directory, ref string) (Snapshot, error) {
	root, head, err := resolve(ctx, directory)
	if err != nil {
		return Snapshot{}, err
	}
	if ref == "" {
		return take(ctx, root, head)
	}

	forked, err := ForkPoint(ctx, root, ref)
	if err != nil {
		return Snapshot{}, err
	}
	return take(ctx, root, forked)
}

func ForkPoint(ctx context.Context, repositoryRoot, ref string) (string, error) {
	if _, err := git(ctx, repositoryRoot, nil, "rev-parse", "--verify", ref+"^{commit}"); err != nil {
		return "", refusal.New(noSuchBase, ref, repositoryRoot)
	}

	forked, err := git(ctx, repositoryRoot, nil, "merge-base", ref, "HEAD")
	if err != nil || forked == "" {
		return "", refusal.New(noCommonAncestor, ref, repositoryRoot)
	}
	return forked, nil
}

func TakeAgainst(ctx context.Context, directory, base string) (Snapshot, error) {
	root, err := Root(ctx, directory)
	if err != nil {
		return Snapshot{}, err
	}
	if err := Reachable(ctx, root, base); err != nil {
		return Snapshot{}, err
	}
	return take(ctx, root, base)
}

func Reachable(ctx context.Context, repositoryRoot, commit string) error {
	if _, err := git(ctx, repositoryRoot, nil, "rev-parse", "--verify", commit+"^{commit}"); err != nil {
		return fmt.Errorf(baseHasGone, commit, repositoryRoot)
	}
	return nil
}

func take(ctx context.Context, root, base string) (Snapshot, error) {
	baseTree, err := git(ctx, root, nil, "rev-parse", "--verify", base+"^{tree}")
	if err != nil {
		return Snapshot{}, err
	}

	indexDirectory, err := os.MkdirTemp("", "socrates-index-")
	if err != nil {
		return Snapshot{}, fmt.Errorf(creatingAnIndexDirectory, err)
	}
	defer func() { _ = os.RemoveAll(indexDirectory) }()

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

func FileAt(ctx context.Context, repositoryRoot, commit, path string) ([]byte, bool, error) {
	listed, err := git(ctx, repositoryRoot, nil, "ls-tree", commit, "--", path)
	if err != nil {
		return nil, false, err
	}
	if !namesABlob(listed) {
		return nil, false, nil
	}

	content, err := gitOutput(ctx, repositoryRoot, nil, "cat-file", "blob", commit+":"+path)
	if err != nil {
		return nil, false, err
	}
	return []byte(content), true, nil
}

func BlobSizeAt(ctx context.Context, repositoryRoot, commit, path string) (int64, bool, error) {
	listed, err := git(ctx, repositoryRoot, nil, "ls-tree", "--long", commit, "--", path)
	if err != nil {
		return 0, false, err
	}
	if !namesABlob(listed) {
		return 0, false, nil
	}

	fields := strings.Fields(listed)
	if len(fields) < 4 {
		return 0, false, fmt.Errorf(unreadableTreeEntry, path, commit, listed)
	}
	size, err := strconv.ParseInt(fields[3], 10, 64)
	if err != nil {
		return 0, false, fmt.Errorf(unreadableTreeEntry, path, commit, listed)
	}
	return size, true, nil
}

func namesABlob(listed string) bool {
	if listed == "" {
		return false
	}
	fields := strings.Fields(listed)
	return len(fields) > 1 && fields[1] == "blob"
}

func Changed(ctx context.Context, repositoryRoot, base, commit string) ([]string, error) {
	changed, err := git(ctx, repositoryRoot, nil, "diff-tree", "-r", "--name-only", "--no-commit-id", base, commit)
	if err != nil {
		return nil, err
	}
	return lines(changed), nil
}

func Diff(ctx context.Context, repositoryRoot, base, commit string) (string, error) {
	return gitOutput(ctx, repositoryRoot, nil, "diff", "--no-color", "--no-ext-diff", "--find-renames",
		"--src-prefix="+hunks.OldSide, "--dst-prefix="+hunks.NewSide, base, commit)
}

func (s Snapshot) Checkout(ctx context.Context, path string) (Worktree, error) {
	if s.NothingToReview() {
		return Worktree{}, errors.New(nothingToCheckOut)
	}

	_, statErr := os.Stat(path)
	presentAlready := statErr == nil

	if _, err := git(ctx, s.RepositoryRoot, nil, "worktree", "add", "--detach", path, s.Commit); err != nil {
		if !presentAlready {
			if removeErr := os.RemoveAll(path); removeErr != nil {
				return Worktree{}, errors.Join(err, fmt.Errorf(removing, path, removeErr))
			}
		}
		if _, pruneErr := git(ctx, s.RepositoryRoot, nil, "worktree", "prune"); pruneErr != nil {
			return Worktree{}, errors.Join(err, pruneErr)
		}
		return Worktree{}, err
	}

	worktree := Worktree{RepositoryRoot: s.RepositoryRoot, Path: path, Commit: s.Commit}

	if err := worktree.Seal(); err != nil {
		return Worktree{}, worktree.discard(ctx, fmt.Errorf(sealingTheWorktree, path, err))
	}
	if err := worktree.ProveTheSeal(); err != nil {
		return Worktree{}, worktree.discard(ctx, err)
	}

	return worktree, nil
}

func (w Worktree) discard(ctx context.Context, reason error) error {
	if err := w.Remove(ctx); err != nil {
		return errors.Join(reason, err)
	}
	return reason
}

func (w Worktree) Differences(ctx context.Context) ([]string, error) {
	reported, err := git(ctx, w.Path, nil, "status", "--porcelain=v1", "--untracked-files=all")
	if err != nil {
		return nil, err
	}

	var differing []string
	for _, line := range lines(reported) {
		if len(line) < 4 {
			continue
		}
		differing = append(differing, strings.TrimSpace(line[3:]))
	}
	return differing, nil
}

func (w Worktree) Remove(ctx context.Context) error {
	if err := w.Unseal(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	if err := os.RemoveAll(w.Path); err != nil {
		return fmt.Errorf(removing, w.Path, err)
	}
	if _, err := git(ctx, w.RepositoryRoot, nil, "worktree", "prune"); err != nil {
		return err
	}
	return nil
}

func Root(ctx context.Context, directory string) (string, error) {
	root, err := git(ctx, directory, nil, "rev-parse", "--show-toplevel")
	if err != nil || root == "" {
		return "", refusal.New(notAWorkingTree, directory)
	}
	return root, nil
}

func resolve(ctx context.Context, directory string) (string, string, error) {
	bare, err := git(ctx, directory, nil, "rev-parse", "--is-bare-repository")
	if err != nil {
		return "", "", refusal.New(notAWorkingTree, directory)
	}
	if bare == "true" {
		return "", "", refusal.New(bareRepository, directory)
	}

	root, err := Root(ctx, directory)
	if err != nil {
		return "", "", err
	}

	base, err := git(ctx, root, nil, "rev-parse", "--verify", "HEAD^{commit}")
	if err != nil {
		return "", "", refusal.New(noBaseCommit, root)
	}

	return root, base, nil
}

func git(ctx context.Context, directory string, environment []string, args ...string) (string, error) {
	produced, err := gitOutput(ctx, directory, environment, args...)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(produced, " \t\r\n"), nil
}

func gitOutput(ctx context.Context, directory string, environment []string, args ...string) (string, error) {
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
		return "", fmt.Errorf(gitFailed, strings.Join(args, " "), directory, reported)
	}

	return stdout.String(), nil
}

func lines(output string) []string {
	if output == "" {
		return nil
	}
	return strings.Split(output, "\n")
}
