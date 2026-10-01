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
	session := harness.Start(t)

	installing := 0
	for _, tool := range convention.Tools {
		help, failures, code := session.Top(tool.Namespace, "--help")
		if code != 0 {
			t.Fatalf("socrates %s --help exited %d: %s", tool.Namespace, code, failures)
		}
		if !strings.Contains(help, "\n  skill ") {
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
