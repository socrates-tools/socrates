package config_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/implement/config"
	"github.com/socrates-tools/socrates/internal/testing/gitrepo"
)

func TestARepositoryIsKeyedByTheTopLevelOfItsMainWorktreeFromAnywhereInsideIt(t *testing.T) {
	repository := gitrepo.New(t)
	repository.Write("deep/inside/file.go", "package inside\n")

	for name, directory := range map[string]string{
		"the top level":     repository.Root,
		"a subdirectory":    filepath.Join(repository.Root, "deep", "inside"),
		"the git directory": filepath.Join(repository.Root, ".git"),
	} {
		got, err := config.RepositoryAt(directory)
		if err != nil || got != repository.Root {
			t.Errorf("from %s the repository is %q (%v), want %q", name, got, err, repository.Root)
		}
	}
}

func TestALinkedWorktreeIsKeyedByTheCloneItBelongsTo(t *testing.T) {
	repository := gitrepo.New(t)
	linked := filepath.Join(t.TempDir(), "linked")
	repository.Git("worktree", "add", "--quiet", "-b", "elsewhere", linked)

	got, err := config.RepositoryAt(linked)
	if err != nil || got != repository.Root {
		t.Errorf("from a worktree of the clone the repository is %q (%v), want the clone %q", got, err, repository.Root)
	}
}

func TestARepositoryWhosePathHoldsADotIsKeyedByThatPath(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "my.repo")
	if err := os.Mkdir(root, 0o755); err != nil {
		t.Fatal(err)
	}
	gitrepo.Hermetic(t)
	repository := gitrepo.At(t, root)
	repository.Git("-c", "init.defaultBranch=main", "init", "--quiet")
	repository.Write("README.md", "x\n")
	repository.Commit("first")

	got, err := config.RepositoryAt(root)
	if err != nil || got != repository.Root || !strings.Contains(got, "my.repo") {
		t.Fatalf("the repository is %q (%v), want %q", got, err, repository.Root)
	}
}

func TestOutsideAnyRepositoryThereIsNoRepository(t *testing.T) {
	gitrepo.Hermetic(t)

	_, err := config.RepositoryAt(t.TempDir())
	if !errors.Is(err, config.ErrNotARepository) {
		t.Errorf("outside a repository the error is %v, want ErrNotARepository", err)
	}
}

func TestABareRepositoryHasNoMainWorktreeToKeyAnEntryBy(t *testing.T) {
	bare := gitrepo.Bare(t)

	_, err := config.RepositoryAt(bare.Root)
	if !errors.Is(err, config.ErrBareRepository) {
		t.Errorf("in a bare repository the error is %v, want ErrBareRepository", err)
	}
}

func TestARepositoryHoldsALocalBranchAndARemoteTrackingBranchAndNothingElse(t *testing.T) {
	origin := gitrepo.New(t)
	origin.Git("branch", "dev")
	origin.Git("tag", "v1")
	clone := gitrepo.At(t, t.TempDir())
	clone.Git("clone", "--quiet", origin.Root, ".")
	clone.Git("branch", "local-only")

	for branch, held := range map[string]bool{
		"main":          true,
		"local-only":    true,
		"origin/dev":    true,
		"dev":           false,
		"v1":            false,
		"missing":       false,
		"":              false,
		"-d":            false,
		"HEAD":          false,
		"origin/nobody": false,
	} {
		if got := config.HoldsBranch(clone.Root, branch); got != held {
			t.Errorf("HoldsBranch(%q) is %v, want %v", branch, got, held)
		}
	}
}

func TestRememberingABranchWritesThatRepositorysEntryAndLeavesTheRestOfTheFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "implement", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	held := "# mine\nconcurrency: 2\n\ndefault_branches:\n  /elsewhere/other.repo: main\n"
	if err := os.WriteFile(path, []byte(held), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := config.RememberDefaultBranch(path, "/home/me/my.repo", "dev"); err != nil {
		t.Fatalf("remembering: %v", err)
	}
	written, _ := os.ReadFile(path)
	if want := held + "  /home/me/my.repo: dev\n"; string(written) != want {
		t.Errorf("the file holds:\n%s\nwant:\n%s", written, want)
	}

	loaded, err := config.LoadFrom(path)
	if err != nil {
		t.Fatalf("loading: %v", err)
	}
	if got := loaded.DefaultBranchOf("/home/me/my.repo"); got != "dev" {
		t.Errorf("the entry reads %q, want dev", got)
	}
	if got := loaded.DefaultBranchOf("/home/me/my"); got != "" {
		t.Errorf("a prefix of the key reads %q, want nothing", got)
	}
}

func TestForgettingABranchRemovesOnlyThatRepositorysEntry(t *testing.T) {
	path := write(t, "default_branches:\n  /a/one.repo: main\n  /a/two.repo: dev\n")

	if err := config.ForgetDefaultBranch(path, "/a/one.repo"); err != nil {
		t.Fatalf("forgetting: %v", err)
	}
	loaded, err := config.LoadFrom(path)
	if err != nil {
		t.Fatalf("loading: %v", err)
	}
	if diff := cmp.Diff(map[string]string{"/a/two.repo": "dev"}, loaded.DefaultBranches); diff != "" {
		t.Errorf("the entries differ (-want +got):\n%s", diff)
	}

	if err := config.ForgetDefaultBranch(path, "/a/two.repo"); err != nil {
		t.Fatalf("forgetting the last: %v", err)
	}
	if written, _ := os.ReadFile(path); strings.TrimSpace(string(written)) != "" {
		t.Errorf("the emptied block left %q behind", written)
	}
}

func TestAnEntryTheLoaderWouldRefuseIsRefusedBeforeTheFileIsTouched(t *testing.T) {
	held := "concurrency: 2\n"
	for name, entry := range map[string][2]string{
		"a relative repository": {"relative/repo", "dev"},
		"a branch with a space": {"/a/repo", "my branch"},
	} {
		t.Run(name, func(t *testing.T) {
			path := write(t, held)
			if err := config.RememberDefaultBranch(path, entry[0], entry[1]); err == nil {
				t.Fatalf("%s was remembered", name)
			}
			if written, _ := os.ReadFile(path); string(written) != held {
				t.Errorf("a refused entry touched the file:\n%s", written)
			}
		})
	}
}

func TestAnEntryWithNoBranchInTheFileIsRefusedNamingTheFile(t *testing.T) {
	path := write(t, "default_branches:\n  /a/repo: \"\"\n")

	_, err := config.LoadFrom(path)
	if err == nil || !strings.Contains(err.Error(), path) {
		t.Errorf("an entry naming no branch loads with %v", err)
	}
}

func TestInsideARepositoryTheSettingsEndWithItsDefaultBranchAndEmptyWhereNothingIsRemembered(t *testing.T) {
	loaded, err := config.LoadFrom(write(t, "default_branches:\n  /a/repo: dev\n"))
	if err != nil {
		t.Fatalf("loading: %v", err)
	}

	remembered := loaded.SettingsIn("/a/repo")
	last := remembered[len(remembered)-1]
	if last.Key != "default_branch" || last.Value != "dev" || last.Source != config.Set {
		t.Errorf("the remembered branch reports %+v", last)
	}

	unremembered := loaded.SettingsIn("/b/repo")
	last = unremembered[len(unremembered)-1]
	if last.Key != "default_branch" || last.Value != "" || last.Source != config.Default {
		t.Errorf("an unremembered branch reports %+v, want empty and default", last)
	}

	if outside := loaded.SettingsIn(""); len(outside) != len(loaded.Settings()) {
		t.Errorf("outside a repository the settings hold %d rows, want %d", len(outside), len(loaded.Settings()))
	}
}
