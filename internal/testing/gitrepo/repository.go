package gitrepo

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type Repository struct {
	t           *testing.T
	Root        string
	Environment Environment
}

type Environment struct {
	Home      string
	Variables []string
}

func NewEnvironment(t *testing.T) Environment {
	t.Helper()

	home := t.TempDir()
	return Environment{
		Home: home,
		Variables: []string{
			"TMPDIR=" + t.TempDir(),
			"HOME=" + home,
			"CLAUDE_CONFIG_DIR=",
			"CODEX_HOME=",
			"GIT_CONFIG_GLOBAL=" + os.DevNull,
			"GIT_CONFIG_SYSTEM=" + os.DevNull,
			"GIT_CONFIG_NOSYSTEM=1",
			"GIT_TERMINAL_PROMPT=0",
			"GIT_CONFIG_COUNT=2",
			"GIT_CONFIG_KEY_0=maintenance.auto",
			"GIT_CONFIG_VALUE_0=false",
			"GIT_CONFIG_KEY_1=gc.auto",
			"GIT_CONFIG_VALUE_1=0",
			"GIT_AUTHOR_NAME=A Person",
			"GIT_AUTHOR_EMAIL=person@example.com",
			"GIT_COMMITTER_NAME=A Person",
			"GIT_COMMITTER_EMAIL=person@example.com",
		},
	}
}

func (e Environment) Of(inherited []string, extra ...string) []string {
	combined := append([]string(nil), inherited...)
	combined = append(combined, e.Variables...)
	return append(combined, extra...)
}

func (e Environment) ApplyToThisProcess(t *testing.T) {
	t.Helper()

	for _, variable := range e.Variables {
		name, value, _ := strings.Cut(variable, "=")
		t.Setenv(name, value)
	}
}

func Hermetic(t *testing.T) {
	t.Helper()
	NewEnvironment(t).ApplyToThisProcess(t)
}

func New(t *testing.T) *Repository {
	t.Helper()

	environment := NewEnvironment(t)
	environment.ApplyToThisProcess(t)
	return initialised(t, environment)
}

func Isolated(t *testing.T) *Repository {
	t.Helper()
	return initialised(t, NewEnvironment(t))
}

func initialised(t *testing.T, environment Environment) *Repository {
	t.Helper()

	repository := At(t, t.TempDir())
	repository.Environment = environment
	repository.Git("-c", "init.defaultBranch=main", "init", "--quiet")
	repository.Write("README.md", "socrates\n")
	repository.Commit("the first commit")
	return repository
}

func Bare(t *testing.T) *Repository {
	t.Helper()

	environment := NewEnvironment(t)
	environment.ApplyToThisProcess(t)
	repository := At(t, t.TempDir())
	repository.Environment = environment
	repository.Git("-c", "init.defaultBranch=main", "init", "--quiet", "--bare")
	return repository
}

func Empty(t *testing.T) *Repository {
	t.Helper()

	environment := NewEnvironment(t)
	environment.ApplyToThisProcess(t)
	repository := At(t, t.TempDir())
	repository.Environment = environment
	repository.Git("-c", "init.defaultBranch=main", "init", "--quiet")
	return repository
}

func At(t *testing.T, root string) *Repository {
	t.Helper()

	resolved, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatalf("resolving %s: %v", root, err)
	}
	return &Repository{t: t, Root: resolved}
}

func (r *Repository) Git(args ...string) string {
	r.t.Helper()

	command := exec.Command("git", append([]string{"-C", r.Root}, args...)...)
	if r.Environment.Variables != nil {
		command.Env = r.Environment.Of(os.Environ())
	}
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	if err := command.Run(); err != nil {
		r.t.Fatalf("git %s: %v: %s", strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return strings.TrimSpace(stdout.String())
}

func (r *Repository) Write(path, content string) {
	r.t.Helper()

	full := filepath.Join(r.Root, path)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		r.t.Fatalf("creating the directory for %s: %v", path, err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		r.t.Fatalf("writing %s: %v", path, err)
	}
}

func (r *Repository) Read(path string) string {
	r.t.Helper()

	content, err := os.ReadFile(filepath.Join(r.Root, path))
	if err != nil {
		r.t.Fatalf("reading %s: %v", path, err)
	}
	return string(content)
}

func (r *Repository) Delete(path string) {
	r.t.Helper()

	if err := os.RemoveAll(filepath.Join(r.Root, path)); err != nil {
		r.t.Fatalf("removing %s: %v", path, err)
	}
}

func (r *Repository) Commit(message string) string {
	r.t.Helper()

	r.Git("add", "--all")
	r.Git("commit", "--quiet", "--message", message)
	return r.Head()
}

func (r *Repository) Head() string {
	r.t.Helper()
	return r.Git("rev-parse", "HEAD")
}

func (r *Repository) Status() string {
	r.t.Helper()
	return r.Git("status", "--porcelain=v2", "--branch", "--untracked-files=all")
}

func (r *Repository) Bytes(path string) []byte {
	r.t.Helper()

	content, err := os.ReadFile(filepath.Join(r.Root, path))
	if err != nil {
		r.t.Fatalf("reading %s: %v", path, err)
	}
	return content
}

func (r *Repository) Tracked(commit string) []string {
	r.t.Helper()

	listed := r.Git("ls-tree", "-r", "--name-only", commit)
	if listed == "" {
		return nil
	}
	return strings.Split(listed, "\n")
}

func (r *Repository) Contents(commit, path string) string {
	r.t.Helper()
	return r.Git("show", commit+":"+path)
}

func (r *Repository) Worktrees() string {
	r.t.Helper()
	return r.Git("worktree", "list", "--porcelain")
}

func (r *Repository) AdministrativeFiles() []string {
	r.t.Helper()

	administrative := filepath.Join(r.Root, ".git")
	var found []string

	err := filepath.WalkDir(administrative, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		relative, relErr := filepath.Rel(administrative, path)
		if relErr != nil {
			return relErr
		}
		found = append(found, filepath.ToSlash(relative))
		return nil
	})
	if err != nil {
		r.t.Fatalf("reading %s: %v", administrative, err)
	}

	return found
}
