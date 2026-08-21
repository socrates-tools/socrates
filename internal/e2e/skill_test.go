package e2e_test

import (
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/rui-ktei/socrates/internal/report"

	"github.com/rui-ktei/socrates/internal/skill"
	"github.com/rui-ktei/socrates/internal/testing/harness"
	"github.com/rui-ktei/socrates/internal/testing/source"
)

func committedSkill(t *testing.T, session *harness.Session) string {
	t.Helper()
	return session.File(filepath.Join(source.Root(t), ".claude", "skills", skill.Name, skill.FileName))
}

func userLevelPath(t *testing.T) string {
	t.Helper()

	home := os.Getenv("HOME")
	if home == "" {
		t.Fatal("the harness left no home directory, so this test would write into the person's own")
	}
	return filepath.Join(home, ".claude", "skills", skill.Name, skill.FileName)
}

func TestInstallingNamesThePath(t *testing.T) {
	session := harness.Start(t)
	path := userLevelPath(t)

	printed, failures, code := session.Scrat("skill")
	if code != 0 {
		t.Fatalf("scrat skill exited %d: %s\n%s", code, failures, printed)
	}
	if !strings.Contains(printed, path) {
		t.Errorf("scrat skill does not say where it wrote, and %s is where it did:\n%s", path, printed)
	}
	if written := session.File(path); written != committedSkill(t, session) {
		t.Errorf("what was installed at user level is not the skill this repository ships:\n%s", written)
	}
}

func TestInstallingTwice(t *testing.T) {
	session := harness.Start(t)
	path := userLevelPath(t)

	first, failures, code := session.Scrat("skill")
	if code != 0 {
		t.Fatalf("scrat skill exited %d: %s", code, failures)
	}

	second, refused, again := session.Scrat("skill")
	if again != 0 {
		t.Fatalf("scrat skill run a second time exited %d: %s", again, refused)
	}
	if second != first {
		t.Errorf("the second installation says something else:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
	if written := session.File(path); written != committedSkill(t, session) {
		t.Errorf("installing twice left something other than the skill at %s:\n%s", path, written)
	}
}

func TestInstallingElsewhere(t *testing.T) {
	session := harness.Start(t)
	target := t.TempDir()

	printed, failures, code := session.Scrat("skill", "--to", target)
	if code != 0 {
		t.Fatalf("scrat skill --to exited %d: %s", code, failures)
	}

	elsewhere := filepath.Join(target, skill.Name, skill.FileName)
	if !strings.Contains(printed, elsewhere) {
		t.Errorf("scrat skill --to does not say it wrote to %s:\n%s", elsewhere, printed)
	}
	if written := session.File(elsewhere); written != committedSkill(t, session) {
		t.Errorf("what was written to the named directory is not the skill:\n%s", written)
	}
	if _, err := os.Stat(userLevelPath(t)); !os.IsNotExist(err) {
		t.Errorf("scrat skill --to also wrote %s", userLevelPath(t))
	}
}

func TestInstallingTheSkillWritesNothingInTheRepository(t *testing.T) {
	session := harness.Start(t)
	session.Repository.Write("internal/added.go", "package internal\n")

	before := session.Repository.Status()
	tracked := workingTree(t, session.Repository.Root)

	printed, failures, code := session.Scrat("skill")
	if code != 0 {
		t.Fatalf("scrat skill exited %d: %s", code, failures)
	}

	if after := session.Repository.Status(); after != before {
		t.Errorf("git status changed:\nbefore:\n%s\nafter:\n%s", before, after)
	}
	if after := workingTree(t, session.Repository.Root); !maps.Equal(after, tracked) {
		t.Errorf("the working tree changed:\nbefore: %v\nafter: %v", tracked, after)
	}
	if strings.Contains(printed, session.Repository.Root) {
		t.Errorf("scrat skill wrote inside the repository under review:\n%s", printed)
	}
}

func TestADirectoryThatCannotBeWritten(t *testing.T) {
	session := harness.Start(t)

	sealed := filepath.Join(t.TempDir(), "skills")
	if err := os.Mkdir(sealed, 0o500); err != nil {
		t.Fatalf("creating a directory nothing can write in: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(sealed, 0o700) })

	printed, failures, code := session.Scrat("skill", "--to", sealed)
	if code != 3 {
		t.Fatalf("scrat skill --to a directory it cannot write exited %d, want 3: %s\n%s", code, failures, printed)
	}
	if !strings.Contains(failures, sealed) {
		t.Errorf("the refusal does not name the directory: %s", failures)
	}
	if !strings.Contains(failures, "permission") {
		t.Errorf("the refusal does not say why the write did not happen: %s", failures)
	}
	if entries, err := os.ReadDir(sealed); err != nil || len(entries) != 0 {
		t.Errorf("the refusal left %d entries in %s: %v", len(entries), sealed, err)
	}
}

func surfaceOf(t *testing.T, session *harness.Session) map[string][]string {
	t.Helper()

	printed, failures, code := session.Scrat("--help")
	if code != 0 {
		t.Fatalf("scrat --help exited %d: %s", code, failures)
	}

	offered := map[string][]string{}
	for _, command := range commandsIn(printed) {
		help, refused, exit := session.Scrat(command, "--help")
		if exit != 0 {
			t.Fatalf("scrat %s --help exited %d: %s", command, exit, refused)
		}
		offered[command] = flagsIn(help)
	}
	return offered
}

func TestEveryCommandNamedIsReal(t *testing.T) {
	session := harness.Start(t)

	named := skill.Named()
	if len(named) < 6 {
		t.Fatalf("the skill names %d invocations, and a check that matched almost nothing would pass forever: %v", len(named), named)
	}

	for _, absent := range skill.Unavailable(surfaceOf(t, session)) {
		t.Error(absent)
	}
}

func TestARemovedFlagFailsTheSuite(t *testing.T) {
	session := harness.Start(t)
	offered := surfaceOf(t, session)

	flagged := 0
	for _, named := range skill.Named() {
		if named.Flag == "" {
			continue
		}
		flagged++

		t.Run(named.Command+" "+named.Flag, func(t *testing.T) {
			without := maps.Clone(offered)
			without[named.Command] = slices.DeleteFunc(slices.Clone(offered[named.Command]), func(flag string) bool { return flag == named.Flag })

			absent := skill.Unavailable(without)
			if len(absent) == 0 {
				t.Fatalf("the check passes with %s removed from scrat %s", named.Flag, named.Command)
			}
			for _, reported := range absent {
				if !strings.Contains(reported, named.Flag) || !strings.Contains(reported, "skill") {
					t.Errorf("the check reports %q, which names neither the flag nor the skill", reported)
				}
			}
		})
	}

	if flagged == 0 {
		t.Fatal("the skill names no flag at all, so this check has nothing to remove")
	}
}

func TestOneSourceTwoSurfaces(t *testing.T) {
	session := harness.Start(t)
	name := parkedRun(t, session, dialogueScenario(patchableFindings, nothingFoundAgain))

	printed, failures, code := session.Scrat("findings", name)
	if code != 4 {
		t.Fatalf("scrat findings of a parked run exited %d: %s", code, failures)
	}

	document := skill.Document()
	for _, rule := range report.ResponseRules() {
		if !strings.Contains(printed, rule) {
			t.Errorf("scrat findings does not print the rule %q:\n%s", rule, printed)
		}
		if !strings.Contains(document, rule) {
			t.Errorf("the skill does not state the rule %q in the same words:\n%s", rule, document)
		}
	}
}

func TestRegeneratingIsOneCommand(t *testing.T) {
	session := harness.Start(t)
	target := t.TempDir()

	printed, failures, code := session.Scrat("skill", "--to", target)
	if code != 0 {
		t.Fatalf("scrat skill --to exited %d: %s\n%s", code, failures, printed)
	}

	written := session.File(filepath.Join(target, skill.Name, skill.FileName))
	if written != committedSkill(t, session) {
		t.Error("the command the regeneration target runs does not reproduce the committed copy: run make skill")
	}
}
