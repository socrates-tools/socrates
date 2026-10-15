package convention_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/socrates-tools/socrates/internal/convention"
	"github.com/socrates-tools/socrates/internal/skill"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

func TestEverySocratesSkillAnyNamespaceShipsIsInvokedByAPersonOnly(t *testing.T) {
	invocable, err := convention.SkillsAModelCouldInvoke(source.Root(t))
	if err != nil {
		t.Fatalf("reading the committed skills: %v", err)
	}

	for _, one := range invocable {
		t.Errorf("%s for %s does not hold that host's person-only switch, so a model could start it", one.Name, one.Host.Name)
	}
}

func writeCommittedSkill(t *testing.T, root string, host skill.Host, name, document string) {
	t.Helper()
	directory := filepath.Join(host.SkillsDirectoryUnder(root), name)
	writePackaged(t, filepath.Join(directory, skill.FileName), document)
	if host.PolicyFile != "" && document != "" {
		writePackaged(t, filepath.Join(directory, host.PolicyFile), skill.PolicyDocument)
	}
}

func everyHostHoldsNothing(t *testing.T, root string) {
	t.Helper()
	for _, host := range skill.Hosts() {
		if err := os.MkdirAll(host.SkillsDirectoryUnder(root), 0o755); err != nil {
			t.Fatal(err)
		}
	}
}

func TestASkillOfAnyNamespaceWithoutTheFieldIsReported(t *testing.T) {
	root := t.TempDir()
	everyHostHoldsNothing(t, root)
	writeCommittedSkill(t, root, skill.ClaudeCode, "socrates-new", "---\nname: socrates-new\ndescription: x\n---\n\n# body\n")
	writeCommittedSkill(t, root, skill.ClaudeCode, "socrates-held", skill.FrontMatter(skill.ClaudeCode, "socrates-held", "x")+"\n")

	invocable, err := convention.SkillsAModelCouldInvoke(root)
	if err != nil {
		t.Fatalf("reading the committed skills: %v", err)
	}
	if len(invocable) != 1 || invocable[0].Name != "socrates-new" || invocable[0].Host != skill.ClaudeCode {
		t.Fatalf("the check reported %v, want only socrates-new for Claude Code", invocable)
	}
}

func TestACodexSkillWithoutItsPolicyIsReportedAndOneWithItIsNot(t *testing.T) {
	root := t.TempDir()
	everyHostHoldsNothing(t, root)
	writeCommittedSkill(t, root, skill.Codex, "socrates-held", skill.FrontMatter(skill.Codex, "socrates-held", "x")+"\n")
	writeCommittedSkill(t, root, skill.Codex, "socrates-open", skill.FrontMatter(skill.Codex, "socrates-open", "x")+"\n")
	if err := os.Remove(filepath.Join(skill.Codex.SkillsDirectoryUnder(root), "socrates-open", skill.Codex.PolicyFile)); err != nil {
		t.Fatal(err)
	}

	invocable, err := convention.SkillsAModelCouldInvoke(root)
	if err != nil {
		t.Fatalf("reading the committed skills: %v", err)
	}
	if len(invocable) != 1 || invocable[0].Name != "socrates-open" || invocable[0].Host != skill.Codex {
		t.Fatalf("the check reported %v, want only socrates-open for Codex", invocable)
	}
}

func TestTheFieldInTheBodyOfASkillDoesNotCountAsFrontMatter(t *testing.T) {
	root := t.TempDir()
	everyHostHoldsNothing(t, root)
	writeCommittedSkill(t, root, skill.ClaudeCode, "socrates-body", "---\nname: socrates-body\n---\n\n"+skill.ClaudeCode.PersonOnlyField+"\n")

	invocable, err := convention.SkillsAModelCouldInvoke(root)
	if err != nil {
		t.Fatalf("reading the committed skills: %v", err)
	}
	if len(invocable) != 1 {
		t.Fatalf("the check reported %v, want socrates-body", invocable)
	}
}

func TestASkillOutsideTheSocratesNamesIsLeftAlone(t *testing.T) {
	root := t.TempDir()
	everyHostHoldsNothing(t, root)
	for _, host := range skill.Hosts() {
		writeCommittedSkill(t, root, host, "release", "---\nname: release\n---\n")
	}

	invocable, err := convention.SkillsAModelCouldInvoke(root)
	if err != nil {
		t.Fatalf("reading the committed skills: %v", err)
	}
	if len(invocable) != 0 {
		t.Fatalf("the check reported %v", invocable)
	}
}

func TestAHostWhoseSkillsCannotBeReadFailsTheCheckRatherThanPassingIt(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(skill.ClaudeCode.SkillsDirectoryUnder(root), 0o755); err != nil {
		t.Fatal(err)
	}

	if _, err := convention.SkillsAModelCouldInvoke(root); err == nil {
		t.Fatal("a root with no skills directory for Codex was read as holding nothing a model could invoke")
	}
}
