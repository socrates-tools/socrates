package e2e_test

import (
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/review/report"

	designskill "github.com/socrates-tools/socrates/design/skill"
	implementskill "github.com/socrates-tools/socrates/implement/skill"
	skillkit "github.com/socrates-tools/socrates/internal/skill"
	"github.com/socrates-tools/socrates/internal/testing/harness"
	"github.com/socrates-tools/socrates/internal/testing/source"
	"github.com/socrates-tools/socrates/review/skill"
)

func allSkills() []skillkit.Skill {
	return slices.Concat(skill.All(), designskill.All(), implementskill.All())
}

func allNamespaces() map[string]string {
	combined := map[string]string{}
	maps.Copy(combined, skill.Namespaces)
	maps.Copy(combined, designskill.Namespaces)
	maps.Copy(combined, implementskill.Namespaces)
	return combined
}

func committedSkill(t *testing.T, session *harness.Session) string {
	t.Helper()
	return session.File(filepath.Join(skillkit.ClaudeCode.SkillsDirectoryUnder(source.Root(t)), skill.Name, skillkit.FileName))
}

func userLevelPath(t *testing.T, session *harness.Session) string {
	t.Helper()

	home := session.UserHome
	if home == "" {
		t.Fatal("the harness left no home directory, so this test would write into the person's own")
	}
	return filepath.Join(skillkit.ClaudeCode.SkillsDirectoryUnder(home), skill.Name, skillkit.FileName)
}

func TestInstallingNamesThePath(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	path := userLevelPath(t, session)

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
	t.Parallel()
	session := harness.Start(t)
	path := userLevelPath(t, session)

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
	t.Parallel()
	session := harness.Start(t)
	target := t.TempDir()

	printed, failures, code := session.Socrates("skill", "--host", skillkit.ClaudeCode.Flag, "--to", target)
	if code != 0 {
		t.Fatalf("socrates review skill --host claude --to exited %d: %s", code, failures)
	}

	elsewhere := filepath.Join(target, skill.Name, skillkit.FileName)
	if !strings.Contains(printed, elsewhere) {
		t.Errorf("socrates review skill --to does not say it wrote to %s:\n%s", elsewhere, printed)
	}
	if written := session.File(elsewhere); written != committedSkill(t, session) {
		t.Errorf("what was written to the named directory is not the skill:\n%s", written)
	}
	if _, err := os.Stat(userLevelPath(t, session)); !os.IsNotExist(err) {
		t.Errorf("socrates review skill --to also wrote %s", userLevelPath(t, session))
	}
}

func TestInstallingTheSkillWritesNothingInTheRepository(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	session := harness.Start(t)

	sealed := filepath.Join(t.TempDir(), "skills")
	if err := os.Mkdir(sealed, 0o500); err != nil {
		t.Fatalf("creating a directory nothing can write in: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(sealed, 0o700) })

	printed, failures, code := session.Socrates("skill", "--host", skillkit.ClaudeCode.Flag, "--to", sealed)
	if code != 3 {
		t.Fatalf("socrates review skill --host claude --to a directory it cannot write exited %d, want 3: %s\n%s", code, failures, printed)
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

func committedDesignSkill(t *testing.T, session *harness.Session) string {
	t.Helper()
	return session.File(filepath.Join(skillkit.ClaudeCode.SkillsDirectoryUnder(source.Root(t)), designskill.Name, skillkit.FileName))
}

func designUserLevelPath(t *testing.T, session *harness.Session) string {
	t.Helper()

	home := session.UserHome
	if home == "" {
		t.Fatal("the harness left no home directory, so this test would write into the person's own")
	}
	return filepath.Join(skillkit.ClaudeCode.SkillsDirectoryUnder(home), designskill.Name, skillkit.FileName)
}

func TestInstallingDesignNamesThePath(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	path := designUserLevelPath(t, session)

	printed, failures, code := session.Top("design", "skill")
	if code != 0 {
		t.Fatalf("socrates design skill exited %d: %s\n%s", code, failures, printed)
	}
	if !strings.Contains(printed, path) {
		t.Errorf("socrates design skill does not say where it wrote, and %s is where it did:\n%s", path, printed)
	}
	if written := session.File(path); written != committedDesignSkill(t, session) {
		t.Errorf("what was installed at user level is not the skill this repository ships:\n%s", written)
	}
}

func TestInstallingDesignTwice(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	path := designUserLevelPath(t, session)

	first, failures, code := session.Top("design", "skill")
	if code != 0 {
		t.Fatalf("socrates design skill exited %d: %s", code, failures)
	}

	second, refused, again := session.Top("design", "skill")
	if again != 0 {
		t.Fatalf("socrates design skill run a second time exited %d: %s", again, refused)
	}
	if second != first {
		t.Errorf("the second installation says something else:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
	if written := session.File(path); written != committedDesignSkill(t, session) {
		t.Errorf("installing twice left something other than the skill at %s:\n%s", path, written)
	}
}

func TestInstallingDesignWritesNothingInTheRepository(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	session.Repository.Write("internal/added.go", "package internal\n")

	before := session.Repository.Status()
	tracked := workingTree(t, session.Repository.Root)

	printed, failures, code := session.Top("design", "skill")
	if code != 0 {
		t.Fatalf("socrates design skill exited %d: %s", code, failures)
	}

	if after := session.Repository.Status(); after != before {
		t.Errorf("git status changed:\nbefore:\n%s\nafter:\n%s", before, after)
	}
	if after := workingTree(t, session.Repository.Root); !maps.Equal(after, tracked) {
		t.Errorf("the working tree changed:\nbefore: %v\nafter: %v", tracked, after)
	}
	if strings.Contains(printed, session.Repository.Root) {
		t.Errorf("socrates design skill wrote inside the repository under review:\n%s", printed)
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

func surfaceOfNamespace(t *testing.T, session *harness.Session, namespace string) map[string][]string {
	t.Helper()

	printed, failures, code := session.Top(namespace, "--help")
	if code != 0 {
		t.Fatalf("socrates %s --help exited %d: %s", namespace, code, failures)
	}

	offered := map[string][]string{}
	for _, command := range commandsIn(printed) {
		help, refused, exit := session.Top(namespace, command, "--help")
		if exit != 0 {
			t.Fatalf("socrates %s %s --help exited %d: %s", namespace, command, exit, refused)
		}
		offered[command] = flagsIn(help)
	}
	return offered
}

func surfaceOfDesign(t *testing.T, session *harness.Session) map[string][]string {
	t.Helper()
	return surfaceOfNamespace(t, session, "design")
}

func surfaceOfImplement(t *testing.T, session *harness.Session) map[string][]string {
	t.Helper()
	return surfaceOfNamespace(t, session, "implement")
}

func everySurface(t *testing.T, session *harness.Session) map[string]map[string][]string {
	t.Helper()

	return map[string]map[string][]string{
		skill.Name:          surfaceOf(t, session),
		skill.Surface:       surfaceOfTheSurface(t),
		designskill.Name:    surfaceOfDesign(t, session),
		implementskill.Name: surfaceOfImplement(t, session),
	}
}

func TestEveryCommandNamedIsReal(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)

	named := skillkit.Named(allSkills(), allNamespaces())
	if len(named) < 6 {
		t.Fatalf("the skills name %d invocations, and a check that matched almost nothing would pass forever: %v", len(named), named)
	}

	attributed := map[string]int{}
	for _, invocation := range named {
		attributed[invocation.Binary]++
	}
	for _, binary := range slices.Concat(skill.Binaries(), designskill.Binaries(), implementskill.Binaries()) {
		if attributed[binary] == 0 {
			t.Errorf("the skills name nothing of %s, and a binary nothing is attributed to is a binary this check never reads", binary)
		}
	}

	for _, absent := range skillkit.Unavailable(allSkills(), allNamespaces(), everySurface(t, session)) {
		t.Error(absent)
	}
}

func TestARemovedFlagFailsTheSuite(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	offered := everySurface(t, session)

	flagged := 0
	for _, named := range skillkit.Named(allSkills(), allNamespaces()) {
		if named.Flag == "" {
			continue
		}
		flagged++

		t.Run(named.Binary+" "+named.Command+" "+named.Flag, func(t *testing.T) {
			without := withoutTheFlag(offered, named)

			absent := skillkit.Unavailable(allSkills(), allNamespaces(), without)
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
		t.Fatal("the skills name no flag at all, so this check has nothing to remove")
	}
}

func TestAFlagOfASecondBinaryIsCheckedToo(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	offered := everySurface(t, session)

	for _, binary := range []string{skill.Surface, designskill.Name, implementskill.Name} {
		t.Run(binary, func(t *testing.T) {
			surfaced := 0
			for _, named := range skillkit.Named(allSkills(), allNamespaces()) {
				if named.Binary != binary || named.Flag == "" {
					continue
				}
				surfaced++

				for _, reported := range skillkit.Unavailable(allSkills(), allNamespaces(), withoutTheFlag(offered, named)) {
					if !strings.Contains(reported, named.Flag) || !strings.Contains(reported, "skill") {
						t.Errorf("the check reports %q for a flag of %s, which names neither the flag nor the skill", reported, named.Binary)
					}
				}
				if len(skillkit.Unavailable(allSkills(), allNamespaces(), withoutTheFlag(offered, named))) == 0 {
					t.Errorf("the check passes with %s removed from %s %s", named.Flag, named.Binary, named.Command)
				}
			}

			if surfaced == 0 {
				t.Fatalf("the skills name no flag of %s, so this check reads nothing of that binary", binary)
			}

			unrecognised := maps.Clone(offered)
			delete(unrecognised, binary)
			absent := skillkit.Unavailable(allSkills(), allNamespaces(), unrecognised)
			if len(absent) == 0 {
				t.Fatalf("the check passes with %s attributed to no binary at all, so an invocation it does not recognise is skipped rather than failed", binary)
			}
			for _, reported := range absent {
				if !strings.Contains(reported, binary) {
					t.Errorf("the check reports %q without naming the binary nothing is known about", reported)
				}
			}
		})
	}
}

func withoutTheFlag(offered map[string]map[string][]string, named skillkit.Invocation) map[string]map[string][]string {
	without := maps.Clone(offered)
	commands := maps.Clone(offered[named.Binary])
	commands[named.Command] = slices.DeleteFunc(slices.Clone(offered[named.Binary][named.Command]), func(flag string) bool { return flag == named.Flag })
	without[named.Binary] = commands
	return without
}

func TestOneSourceTwoSurfaces(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	name := parkedRun(t, session, dialogueScenario(patchableFindings, nothingFoundAgain))

	printed, failures, code := session.Socrates("findings", name)
	if code != 4 {
		t.Fatalf("socrates review findings of a parked run exited %d: %s", code, failures)
	}

	document := skill.Document(skillkit.ClaudeCode)
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
	t.Parallel()
	session := harness.Start(t)
	target := t.TempDir()

	printed, failures, code := session.Socrates("skill", "--host", skillkit.ClaudeCode.Flag, "--to", target)
	if code != 0 {
		t.Fatalf("socrates review skill --host claude --to exited %d: %s\n%s", code, failures, printed)
	}

	written := session.File(filepath.Join(target, skill.Name, skillkit.FileName))
	if written != committedSkill(t, session) {
		t.Error("the command the regeneration target runs does not reproduce the committed copy: run make skill")
	}
}

func TestBothSkillsAreInstalled(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)

	printed, failures, code := session.Socrates("skill")
	if code != 0 {
		t.Fatalf("socrates review skill exited %d: %s\n%s", code, failures, printed)
	}

	home := session.UserHome
	if home == "" {
		t.Fatal("the harness left no home directory, so this test would write into the person's own")
	}

	for _, one := range skill.All() {
		path := filepath.Join(one.Host.SkillsDirectoryUnder(home), one.Name, skillkit.FileName)

		if !strings.Contains(printed, path) {
			t.Errorf("socrates review skill does not say it wrote %s:\n%s", path, printed)
		}
		if written := session.File(path); written != one.Rendered() {
			t.Errorf("what was installed at %s is not what the skill's source generates", path)
		}
	}
	if len(skill.All()) < 2 {
		t.Fatal("this repository states one skill, and a session reviewing somebody else's change reads a second one")
	}
}

func TestBothDesignSkillsAreInstalled(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)

	printed, failures, code := session.Top("design", "skill")
	if code != 0 {
		t.Fatalf("socrates design skill exited %d: %s\n%s", code, failures, printed)
	}

	home := session.UserHome
	if home == "" {
		t.Fatal("the harness left no home directory, so this test would write into the person's own")
	}

	for _, one := range designskill.All() {
		path := filepath.Join(one.Host.SkillsDirectoryUnder(home), one.Name, skillkit.FileName)

		if !strings.Contains(printed, path) {
			t.Errorf("socrates design skill does not say it wrote %s:\n%s", path, printed)
		}
		if written := session.File(path); written != one.Rendered() {
			t.Errorf("what was installed at %s is not what the skill's source generates", path)
		}
		if committed := session.File(filepath.Join(one.Host.SkillsDirectoryUnder(source.Root(t)), one.Name, skillkit.FileName)); committed != one.Rendered() {
			t.Errorf("the committed copy of %s is not what its source generates: run make skill", one.Name)
		}
	}
	if len(designskill.All()) < 2 {
		t.Fatal("design states one skill, and a spec is cut into tickets by a second one")
	}
}

func TestTheOverseerSkillIsInstalledAndTheCommittedCopyIsWhatItsSourceGenerates(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)

	printed, failures, code := session.Top("implement", "skill")
	if code != 0 {
		t.Fatalf("socrates implement skill exited %d: %s\n%s", code, failures, printed)
	}

	home := session.UserHome
	if home == "" {
		t.Fatal("the harness left no home directory, so this test would write into the person's own")
	}
	path := filepath.Join(skillkit.ClaudeCode.SkillsDirectoryUnder(home), implementskill.Name, skillkit.FileName)

	if !strings.Contains(printed, path) {
		t.Errorf("socrates implement skill does not say it wrote %s:\n%s", path, printed)
	}
	if written := session.File(path); written != implementskill.Document(skillkit.ClaudeCode) {
		t.Errorf("what was installed at %s is not what the skill's source generates", path)
	}
	if committed := session.File(filepath.Join(skillkit.ClaudeCode.SkillsDirectoryUnder(source.Root(t)), implementskill.Name, skillkit.FileName)); committed != implementskill.Document(skillkit.ClaudeCode) {
		t.Errorf("the committed copy of %s is not what its source generates: run make skill", implementskill.Name)
	}
}

func TestTheOverseerSkillIsInstalledBesideTheDefinitionUnderTheDirectoriesNamed(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	skills, agents := t.TempDir(), t.TempDir()

	printed, failures, code := session.Top("implement", "skill", "--host", skillkit.ClaudeCode.Flag, "--to", skills, "--agents-to", agents)
	if code != 0 {
		t.Fatalf("socrates implement skill exited %d: %s\n%s", code, failures, printed)
	}

	if written := session.File(filepath.Join(skills, implementskill.Name, skillkit.FileName)); written != implementskill.Document(skillkit.ClaudeCode) {
		t.Error("what was installed under --to is not what the skill's source generates")
	}
	if _, err := os.Stat(filepath.Join(agents, implementskill.AgentFileName())); err != nil {
		t.Errorf("the agent definition is not beside the skill under --agents-to: %v", err)
	}
}

func TestACodexHostWithAnAgentsDirectoryWritesTheRoleThereAndNothingElse(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	agents := filepath.Join(t.TempDir(), "agents")

	printed, failures, code := session.Top("implement", "skill", "--host", skillkit.Codex.Flag, "--agents-to", agents)
	if code != 0 {
		t.Fatalf("socrates implement skill --host codex --agents-to exited %d: %s\n%s", code, failures, printed)
	}

	role := filepath.Join(agents, implementskill.CodexRoleFileName())
	if written := session.File(role); written != implementskill.DefaultCodexRole() {
		t.Errorf("what was installed at %s is not what the role's source generates", role)
	}
	if !strings.Contains(printed, role) {
		t.Errorf("socrates implement skill does not say it wrote %s:\n%s", role, printed)
	}
	if written := writtenIn(t, agents); len(written) != 1 {
		t.Errorf("--host codex --agents-to wrote %v, want the role alone", slices.Sorted(maps.Keys(written)))
	}
	if _, claude := writtenUnder(t, session.UserHome)[installedAgentKey(skillkit.ClaudeCode)]; claude {
		t.Error("--host codex wrote Claude Code's agent definition")
	}
}

func TestEveryHostsFrontMatterCarriesItsOwnKeysAndNoOthers(t *testing.T) {
	t.Parallel()

	for _, one := range allSkills() {
		front, _, found := strings.Cut(strings.TrimPrefix(one.Rendered(), "---\n"), "\n---\n")
		if !found {
			t.Fatalf("%s for %s has no front matter", one.Name, one.Host.Name)
		}

		var keys []string
		for _, line := range strings.Split(front, "\n") {
			key, _, _ := strings.Cut(line, ":")
			keys = append(keys, key)
		}

		want := []string{"name", "description"}
		if one.Host.PersonOnlyField != "" {
			key, _, _ := strings.Cut(one.Host.PersonOnlyField, ":")
			want = append(want, key)
		}
		if !slices.Equal(keys, want) {
			t.Errorf("%s for %s has front matter keys %v, want %v", one.Name, one.Host.Name, keys, want)
		}
	}
}

func servingSkills() []skillkit.Skill {
	var serving []skillkit.Skill
	for _, one := range allSkills() {
		if strings.Contains(one.Rendered(), "socrates agora serve") {
			serving = append(serving, one)
		}
	}
	return serving
}

func TestEveryCodexRenderingThatServesAPageEscalatesServePollAndEnd(t *testing.T) {
	t.Parallel()

	var checked int
	for _, one := range servingSkills() {
		if one.Host != skillkit.Codex {
			continue
		}
		checked++
		for _, want := range []string{"sandbox_permissions", "require_escalated", "socrates agora serve", "socrates agora poll", "socrates agora end", "justification"} {
			if !strings.Contains(one.Rendered(), want) {
				t.Errorf("%s for Codex serves a page and does not say %q", one.Name, want)
			}
		}
	}
	if checked != 4 {
		t.Errorf("%d Codex renderings serve a page, want the four design, tickets, review and review handover", checked)
	}
}

func TestEveryCodexRenderingThatDelegatesStartsAgentsWithSpawnAgentAndCollectsThemWithWaitAgent(t *testing.T) {
	t.Parallel()

	for _, one := range servingSkills() {
		if one.Host != skillkit.Codex || !strings.Contains(one.Rendered(), "subagent") {
			continue
		}
		for _, want := range []string{"spawn_agent", "wait_agent", "answer the person between"} {
			if !strings.Contains(one.Rendered(), want) {
				t.Errorf("%s for Codex sends work to a subagent and does not say %q", one.Name, want)
			}
		}
	}
}

func TestNoClaudeCodeRenderingNamesCodexsEscalationOrAgentCalls(t *testing.T) {
	t.Parallel()

	for _, one := range allSkills() {
		if one.Host != skillkit.ClaudeCode {
			continue
		}
		for _, unsaid := range []string{"require_escalated", "spawn_agent", "wait_agent"} {
			if strings.Contains(one.Rendered(), unsaid) {
				t.Errorf("%s for Claude Code names %q, which is Codex's", one.Name, unsaid)
			}
		}
	}
}

func TestNoCodexRenderingThatServesAPageDescribesBackgroundWork(t *testing.T) {
	t.Parallel()

	for _, one := range servingSkills() {
		if one.Host == skillkit.Codex && strings.Contains(strings.ToLower(one.Rendered()), "background") {
			t.Errorf("%s for Codex describes work in the background, which is Claude Code's", one.Name)
		}
	}
}

func TestEveryCodexRenderingThatServesAPageTellsTheSessionToInvokeNoSkill(t *testing.T) {
	t.Parallel()

	for _, one := range servingSkills() {
		if one.Host == skillkit.Codex && !strings.Contains(one.Rendered(), "Invoke no skill") {
			t.Errorf("%s for Codex does not tell the session to invoke no skill", one.Name)
		}
	}
}

func TestImplementAttributesReviewsCommandsToReviewsOwnBinary(t *testing.T) {
	t.Parallel()

	if got := implementskill.Namespaces["review"]; got != skill.Name {
		t.Errorf("implement attributes socrates review commands to %q, and review's own binary is %q, which the merged namespaces would overwrite", got, skill.Name)
	}
}
