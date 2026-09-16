package e2e_test

import (
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/review/report"

	"github.com/socrates-tools/socrates/internal/testing/harness"
	"github.com/socrates-tools/socrates/internal/testing/source"
	"github.com/socrates-tools/socrates/review/skill"
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

	printed, failures, code := session.Socrates("skill")
	if code != 0 {
		t.Fatalf("socrates review skill exited %d: %s\n%s", code, failures, printed)
	}
	if !strings.Contains(printed, path) {
		t.Errorf("socrates review skill does not say where it wrote, and %s is where it did:\n%s", path, printed)
	}
	if written := session.File(path); written != committedSkill(t, session) {
		t.Errorf("what was installed at user level is not the skill this repository ships:\n%s", written)
	}
}

func TestInstallingTwice(t *testing.T) {
	session := harness.Start(t)
	path := userLevelPath(t)

	first, failures, code := session.Socrates("skill")
	if code != 0 {
		t.Fatalf("socrates review skill exited %d: %s", code, failures)
	}

	second, refused, again := session.Socrates("skill")
	if again != 0 {
		t.Fatalf("socrates review skill run a second time exited %d: %s", again, refused)
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

	printed, failures, code := session.Socrates("skill", "--to", target)
	if code != 0 {
		t.Fatalf("socrates review skill --to exited %d: %s", code, failures)
	}

	elsewhere := filepath.Join(target, skill.Name, skill.FileName)
	if !strings.Contains(printed, elsewhere) {
		t.Errorf("socrates review skill --to does not say it wrote to %s:\n%s", elsewhere, printed)
	}
	if written := session.File(elsewhere); written != committedSkill(t, session) {
		t.Errorf("what was written to the named directory is not the skill:\n%s", written)
	}
	if _, err := os.Stat(userLevelPath(t)); !os.IsNotExist(err) {
		t.Errorf("socrates review skill --to also wrote %s", userLevelPath(t))
	}
}

func TestInstallingTheSkillWritesNothingInTheRepository(t *testing.T) {
	session := harness.Start(t)
	session.Repository.Write("internal/added.go", "package internal\n")

	before := session.Repository.Status()
	tracked := workingTree(t, session.Repository.Root)

	printed, failures, code := session.Socrates("skill")
	if code != 0 {
		t.Fatalf("socrates review skill exited %d: %s", code, failures)
	}

	if after := session.Repository.Status(); after != before {
		t.Errorf("git status changed:\nbefore:\n%s\nafter:\n%s", before, after)
	}
	if after := workingTree(t, session.Repository.Root); !maps.Equal(after, tracked) {
		t.Errorf("the working tree changed:\nbefore: %v\nafter: %v", tracked, after)
	}
	if strings.Contains(printed, session.Repository.Root) {
		t.Errorf("socrates review skill wrote inside the repository under review:\n%s", printed)
	}
}

func TestADirectoryThatCannotBeWritten(t *testing.T) {
	session := harness.Start(t)

	sealed := filepath.Join(t.TempDir(), "skills")
	if err := os.Mkdir(sealed, 0o500); err != nil {
		t.Fatalf("creating a directory nothing can write in: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(sealed, 0o700) })

	printed, failures, code := session.Socrates("skill", "--to", sealed)
	if code != 3 {
		t.Fatalf("socrates review skill --to a directory it cannot write exited %d, want 3: %s\n%s", code, failures, printed)
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

	printed, failures, code := session.Socrates("--help")
	if code != 0 {
		t.Fatalf("socrates --help exited %d: %s", code, failures)
	}

	offered := map[string][]string{}
	for _, command := range commandsIn(printed) {
		help, refused, exit := session.Socrates(command, "--help")
		if exit != 0 {
			t.Fatalf("socrates %s --help exited %d: %s", command, exit, refused)
		}
		offered[command] = flagsIn(help)
	}
	return offered
}

func surfaceOfTheSurface(t *testing.T) map[string][]string {
	t.Helper()

	held := servingSurface(t)
	printed, failures, code := held.run("--help")
	if code != 0 {
		t.Fatalf("agora --help exited %d: %s", code, failures)
	}

	offered := map[string][]string{}
	for _, command := range commandsIn(printed) {
		help, refused, exit := held.run(command, "--help")
		if exit != 0 {
			t.Fatalf("agora %s --help exited %d: %s", command, exit, refused)
		}
		offered[command] = flagsIn(help)
	}
	return offered
}

func everySurface(t *testing.T, session *harness.Session) map[string]map[string][]string {
	t.Helper()

	return map[string]map[string][]string{
		skill.Name:    surfaceOf(t, session),
		skill.Surface: surfaceOfTheSurface(t),
	}
}

func TestEveryCommandNamedIsReal(t *testing.T) {
	session := harness.Start(t)

	named := skill.Named()
	if len(named) < 6 {
		t.Fatalf("the skill names %d invocations, and a check that matched almost nothing would pass forever: %v", len(named), named)
	}

	attributed := map[string]int{}
	for _, invocation := range named {
		attributed[invocation.Binary]++
	}
	for _, binary := range skill.Binaries() {
		if attributed[binary] == 0 {
			t.Errorf("the skill names nothing of %s, and a binary nothing is attributed to is a binary this check never reads", binary)
		}
	}

	for _, absent := range skill.Unavailable(everySurface(t, session)) {
		t.Error(absent)
	}
}

func TestARemovedFlagFailsTheSuite(t *testing.T) {
	session := harness.Start(t)
	offered := everySurface(t, session)

	flagged := 0
	for _, named := range skill.Named() {
		if named.Flag == "" {
			continue
		}
		flagged++

		t.Run(named.Binary+" "+named.Command+" "+named.Flag, func(t *testing.T) {
			without := withoutTheFlag(offered, named)

			absent := skill.Unavailable(without)
			if len(absent) == 0 {
				t.Fatalf("the check passes with %s removed from %s %s", named.Flag, named.Binary, named.Command)
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

func TestAFlagOfTheSecondBinaryIsCheckedToo(t *testing.T) {
	session := harness.Start(t)
	offered := everySurface(t, session)

	surfaced := 0
	for _, named := range skill.Named() {
		if named.Binary != skill.Surface || named.Flag == "" {
			continue
		}
		surfaced++

		for _, reported := range skill.Unavailable(withoutTheFlag(offered, named)) {
			if !strings.Contains(reported, named.Flag) || !strings.Contains(reported, "skill") {
				t.Errorf("the check reports %q for a flag of %s, which names neither the flag nor the skill", reported, named.Binary)
			}
		}
		if len(skill.Unavailable(withoutTheFlag(offered, named))) == 0 {
			t.Errorf("the check passes with %s removed from %s %s", named.Flag, named.Binary, named.Command)
		}
	}

	if surfaced == 0 {
		t.Fatalf("the skill names no flag of %s, so this check reads nothing of the second binary", skill.Surface)
	}

	unrecognised := maps.Clone(offered)
	delete(unrecognised, skill.Surface)
	absent := skill.Unavailable(unrecognised)
	if len(absent) == 0 {
		t.Fatalf("the check passes with %s attributed to no binary at all, so an invocation it does not recognise is skipped rather than failed", skill.Surface)
	}
	for _, reported := range absent {
		if !strings.Contains(reported, skill.Surface) {
			t.Errorf("the check reports %q without naming the binary nothing is known about", reported)
		}
	}
}

func withoutTheFlag(offered map[string]map[string][]string, named skill.Invocation) map[string]map[string][]string {
	without := maps.Clone(offered)
	commands := maps.Clone(offered[named.Binary])
	commands[named.Command] = slices.DeleteFunc(slices.Clone(offered[named.Binary][named.Command]), func(flag string) bool { return flag == named.Flag })
	without[named.Binary] = commands
	return without
}

func TestOneSourceTwoSurfaces(t *testing.T) {
	session := harness.Start(t)
	name := parkedRun(t, session, dialogueScenario(patchableFindings, nothingFoundAgain))

	printed, failures, code := session.Socrates("findings", name)
	if code != 4 {
		t.Fatalf("socrates review findings of a parked run exited %d: %s", code, failures)
	}

	document := skill.Document()
	for _, rule := range report.ResponseRules() {
		if !strings.Contains(printed, rule) {
			t.Errorf("socrates review findings does not print the rule %q:\n%s", rule, printed)
		}
		if !strings.Contains(document, rule) {
			t.Errorf("the skill does not state the rule %q in the same words:\n%s", rule, document)
		}
	}
}

func TestRegeneratingIsOneCommand(t *testing.T) {
	session := harness.Start(t)
	target := t.TempDir()

	printed, failures, code := session.Socrates("skill", "--to", target)
	if code != 0 {
		t.Fatalf("socrates review skill --to exited %d: %s\n%s", code, failures, printed)
	}

	written := session.File(filepath.Join(target, skill.Name, skill.FileName))
	if written != committedSkill(t, session) {
		t.Error("the command the regeneration target runs does not reproduce the committed copy: run make skill")
	}
}

func TestBothSkillsAreInstalled(t *testing.T) {
	session := harness.Start(t)

	printed, failures, code := session.Socrates("skill")
	if code != 0 {
		t.Fatalf("socrates review skill exited %d: %s\n%s", code, failures, printed)
	}

	home := os.Getenv("HOME")
	if home == "" {
		t.Fatal("the harness left no home directory, so this test would write into the person's own")
	}

	for _, one := range skill.All() {
		path := filepath.Join(home, ".claude", "skills", one.Name, skill.FileName)

		if !strings.Contains(printed, path) {
			t.Errorf("socrates review skill does not say it wrote %s:\n%s", path, printed)
		}
		if written := session.File(path); written != one.Document() {
			t.Errorf("what was installed at %s is not what the skill's source generates", path)
		}
	}
	if len(skill.All()) < 2 {
		t.Fatal("this repository states one skill, and a session reviewing somebody else's change reads a second one")
	}
}
