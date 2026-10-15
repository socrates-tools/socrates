package e2e_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/convention"
	skillkit "github.com/socrates-tools/socrates/internal/skill"
	"github.com/socrates-tools/socrates/internal/testing/harness"
)

func installsATarget(t *testing.T, session *harness.Session, namespace string) bool {
	t.Helper()
	help, _, code := session.Top(namespace, "skill", "--help")
	return code == 0 && strings.Contains(help, "\n      --to ")
}

func TestEverySkillAnyNamespaceInstallsIsInvokedByAPersonOnly(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)

	installing := 0
	for _, tool := range convention.Tools {
		if !installsATarget(t, session, tool.Namespace) {
			continue
		}

		for _, host := range skillkit.Hosts() {
			directory := t.TempDir()
			printed, failures, code := session.Top(tool.Namespace, "skill", "--host", host.Flag, "--to", directory)
			if code != 0 {
				t.Fatalf("socrates %s skill --host %s exited %d: %s\n%s", tool.Namespace, host.Flag, code, failures, printed)
			}

			installedHere := 0
			walked := filepath.WalkDir(directory, func(path string, entry fs.DirEntry, err error) error {
				if err != nil || entry.IsDir() || entry.Name() != skillkit.FileName {
					return err
				}
				installedHere++
				if !host.PersonOnly(filepath.Dir(path)) {
					t.Errorf("socrates %s skill --host %s installed %s without %s's person-only switch, so a model could start it", tool.Namespace, host.Flag, path, host.Name)
				}
				return nil
			})
			if walked != nil {
				t.Fatalf("reading what socrates %s skill installed: %v", tool.Namespace, walked)
			}
			if installedHere == 0 {
				t.Errorf("socrates %s skill --host %s offers an install and wrote no %s", tool.Namespace, host.Flag, skillkit.FileName)
			}
			installing += installedHere
		}
	}

	if installing == 0 {
		t.Fatal("no namespace offers a skill command, so this check would pass forever")
	}
}

func TestACodexSkillThatLosesItsPolicyIsNoLongerPersonOnlyAsAClaudeCodeOneLosingItsFieldIsNot(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)

	for _, host := range skillkit.Hosts() {
		directory := t.TempDir()
		if _, failures, code := session.Top("design", "skill", "--host", host.Flag, "--to", directory); code != 0 {
			t.Fatalf("socrates design skill --host %s exited %d: %s", host.Flag, code, failures)
		}
		installed := filepath.Join(directory, "socrates-design")
		if !host.PersonOnly(installed) {
			t.Fatalf("the skill installed for %s is not person-only before anything is removed", host.Name)
		}

		if host.PolicyFile != "" {
			if err := os.Remove(filepath.Join(installed, host.PolicyFile)); err != nil {
				t.Fatal(err)
			}
		} else {
			document := filepath.Join(installed, skillkit.FileName)
			held, err := os.ReadFile(document)
			if err != nil {
				t.Fatal(err)
			}
			stripped := strings.Replace(string(held), host.PersonOnlyField+"\n", "", 1)
			if err := os.WriteFile(document, []byte(stripped), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		if host.PersonOnly(installed) {
			t.Errorf("the skill installed for %s is still read as person-only after its switch was removed", host.Name)
		}
	}
}

func TestEveryInstallSaysAPersonStartsTheSkillAndNoSessionDoes(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)

	announced := 0
	for _, tool := range convention.Tools {
		if !installsATarget(t, session, tool.Namespace) {
			continue
		}

		for _, host := range skillkit.Hosts() {
			directory := t.TempDir()
			printed, failures, code := session.Top(tool.Namespace, "skill", "--host", host.Flag, "--to", directory)
			if code != 0 {
				t.Fatalf("socrates %s skill --host %s exited %d: %s\n%s", tool.Namespace, host.Flag, code, failures, printed)
			}

			installed, err := filepath.Glob(filepath.Join(directory, "*", skillkit.FileName))
			if err != nil {
				t.Fatalf("reading what socrates %s skill installed: %v", tool.Namespace, err)
			}
			for _, path := range installed {
				name := filepath.Base(filepath.Dir(path))
				want := "the " + name + " skill is at " + path + ", and only a person starts it, by typing " + host.Starts(name)
				if !strings.Contains(printed, want+"\n") {
					t.Errorf("socrates %s skill --host %s printed:\n%s\nand a person reading it is not told that only they start %s, with:\n%s", tool.Namespace, host.Flag, printed, name, want)
				}
				announced++
			}
			if strings.Contains(printed, "a session") {
				t.Errorf("socrates %s skill printed:\n%s\nwhich says a session starts a skill that only a person can start", tool.Namespace, printed)
			}
		}
	}

	if announced == 0 {
		t.Fatal("no namespace installed a skill, so this check would pass forever")
	}
}
