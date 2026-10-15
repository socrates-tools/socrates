package e2e_test

import (
	"maps"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/convention"
	skillkit "github.com/socrates-tools/socrates/internal/skill"
	"github.com/socrates-tools/socrates/internal/testing/harness"
)

func TestADirectoryFlagWithoutAHostIsRefusedNamingHostBeforeAnythingIsWritten(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)

	for _, tool := range convention.Tools {
		if !installsATarget(t, session, tool.Namespace) {
			continue
		}

		flags := []string{skillkit.ToFlag}
		if tool.Namespace == "implement" {
			flags = append(flags, skillkit.AgentsToFlag)
		}
		for _, flag := range flags {
			directory := t.TempDir()
			printed, failures, code := session.Top(tool.Namespace, "skill", "--"+flag, directory)
			if code == 0 {
				t.Fatalf("socrates %s skill --%s exited 0 without --host:\n%s", tool.Namespace, flag, printed)
			}
			if !strings.Contains(failures, "--host") || !strings.Contains(failures, "--"+flag) {
				t.Errorf("socrates %s skill --%s without --host did not name both flags:\n%s", tool.Namespace, flag, failures)
			}
			if entries, _ := filepath.Glob(filepath.Join(directory, "*")); len(entries) != 0 {
				t.Errorf("socrates %s skill --%s without --host wrote %v", tool.Namespace, flag, entries)
			}
			if written := writtenUnder(t, session.UserHome); len(written) != 0 {
				t.Errorf("socrates %s skill --%s without --host wrote %v under the home", tool.Namespace, flag, slices.Sorted(maps.Keys(written)))
			}
		}
	}
}

func TestAHostWithADirectoryWritesThatHostsRenderingThereAndNothingElse(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	directory := t.TempDir()

	printed, failures, code := session.Top("design", "skill", "--host", skillkit.Codex.Flag, "--to", directory)
	if code != 0 {
		t.Fatalf("socrates design skill --host codex --to exited %d: %s\n%s", code, failures, printed)
	}

	written := writtenIn(t, directory)
	for _, one := range allSkills() {
		if one.Host != skillkit.Codex || !strings.HasPrefix(one.Name, "socrates-design") && one.Name != "socrates-tickets" {
			continue
		}
		if written[filepath.ToSlash(filepath.Join(one.Name, skillkit.FileName))] != one.Rendered() {
			t.Errorf("--host codex --to did not write Codex's rendering of %s", one.Name)
		}
		if written[filepath.ToSlash(filepath.Join(one.Name, skillkit.Codex.PolicyFile))] != skillkit.PolicyDocument {
			t.Errorf("--host codex --to did not write the policy beside %s", one.Name)
		}
	}
	if len(written) != 4 {
		t.Errorf("--host codex --to wrote %v, want two skills and a policy beside each", slices.Sorted(maps.Keys(written)))
	}
	if underHome := writtenUnder(t, session.UserHome); len(underHome) != 0 {
		t.Errorf("--host codex --to wrote %v under the home", slices.Sorted(maps.Keys(underHome)))
	}
}

func TestAHostAloneWritesOnlyThatHostsRenderingToItsOwnDirectory(t *testing.T) {
	t.Parallel()

	for _, host := range skillkit.Hosts() {
		t.Run(host.Flag, func(t *testing.T) {
			t.Parallel()
			session := harness.Start(t)

			for _, namespace := range []string{"design", "review", "implement"} {
				if printed, failures, code := session.Top(namespace, "skill", "--host", host.Flag); code != 0 {
					t.Fatalf("socrates %s skill --host %s exited %d: %s\n%s", namespace, host.Flag, code, failures, printed)
				}
			}

			written := writtenUnder(t, session.UserHome)
			for path := range written {
				if !strings.HasPrefix(path, filepath.ToSlash(host.SkillsDirectory)+"/") && !strings.HasPrefix(path, filepath.ToSlash(host.AgentsDirectory)+"/") {
					t.Errorf("--host %s wrote %s, which is not under that host's directories", host.Flag, path)
				}
			}
			for _, one := range allSkills() {
				if one.Host != host {
					continue
				}
				if written[installedKey(host, one.Name)] != one.Rendered() {
					t.Errorf("--host %s did not write %s as the package renders it", host.Flag, one.Name)
				}
			}
		})
	}
}
