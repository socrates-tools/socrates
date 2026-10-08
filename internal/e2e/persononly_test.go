package e2e_test

import (
	"io/fs"
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/convention"
	skillkit "github.com/socrates-tools/socrates/internal/skill"
	"github.com/socrates-tools/socrates/internal/testing/harness"
)

func TestEverySkillAnyNamespaceInstallsIsInvokedByAPersonOnly(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)

	installing := 0
	for _, tool := range convention.Tools {
		help, _, code := session.Top(tool.Namespace, "skill", "--help")
		if code != 0 || !strings.Contains(help, "\n      --to ") {
			continue
		}

		directory := t.TempDir()
		printed, failures, code := session.Top(tool.Namespace, "skill", "--to", directory)
		if code != 0 {
			t.Fatalf("socrates %s skill exited %d: %s\n%s", tool.Namespace, code, failures, printed)
		}

		installedHere := 0
		walked := filepath.WalkDir(directory, func(path string, entry fs.DirEntry, err error) error {
			if err != nil || entry.IsDir() || entry.Name() != skillkit.FileName {
				return err
			}
			installedHere++
			if !convention.CarriesThePersonOnlyField(session.File(path)) {
				t.Errorf("socrates %s skill installed %s without %s in its front matter, so a model could start it", tool.Namespace, path, skillkit.PersonOnlyField)
			}
			return nil
		})
		if walked != nil {
			t.Fatalf("reading what socrates %s skill installed: %v", tool.Namespace, walked)
		}
		if installedHere == 0 {
			t.Errorf("socrates %s skill offers an install and wrote no %s", tool.Namespace, skillkit.FileName)
		}
		installing += installedHere
	}

	if installing == 0 {
		t.Fatal("no namespace offers a skill command, so this check would pass forever")
	}
}

func TestEveryInstallSaysAPersonStartsTheSkillAndNoSessionDoes(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)

	announced := 0
	for _, tool := range convention.Tools {
		help, _, code := session.Top(tool.Namespace, "skill", "--help")
		if code != 0 || !strings.Contains(help, "\n      --to ") {
			continue
		}

		directory := t.TempDir()
		printed, failures, code := session.Top(tool.Namespace, "skill", "--to", directory)
		if code != 0 {
			t.Fatalf("socrates %s skill exited %d: %s\n%s", tool.Namespace, code, failures, printed)
		}

		installed, err := filepath.Glob(filepath.Join(directory, "*", skillkit.FileName))
		if err != nil {
			t.Fatalf("reading what socrates %s skill installed: %v", tool.Namespace, err)
		}
		for _, path := range installed {
			name := filepath.Base(filepath.Dir(path))
			want := "the " + name + " skill is at " + path + ", and only a person starts it, by typing /" + name
			if !strings.Contains(printed, want+"\n") {
				t.Errorf("socrates %s skill printed:\n%s\nand a person reading it is not told that only they start %s, with:\n%s", tool.Namespace, printed, name, want)
			}
			announced++
		}
		if strings.Contains(printed, "a session") {
			t.Errorf("socrates %s skill printed:\n%s\nwhich says a session starts a skill that only a person can start", tool.Namespace, printed)
		}
	}

	if announced == 0 {
		t.Fatal("no namespace installed a skill, so this check would pass forever")
	}
}
