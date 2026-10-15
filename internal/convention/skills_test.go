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
		t.Errorf("%s does not carry %s in its front matter, so a model could start it", one.Name, one.Host.PersonOnlyField)
	}
}

func writeCommittedSkill(t *testing.T, root, name, document string) {
	t.Helper()
	writePackaged(t, filepath.Join(skill.ClaudeCode.SkillsDirectoryUnder(root), name, skill.FileName), document)
}

func TestASkillOfAnyNamespaceWithoutTheFieldIsReported(t *testing.T) {
	root := t.TempDir()
	writeCommittedSkill(t, root, "socrates-new", "---\nname: socrates-new\ndescription: x\n---\n\n# body\n")
	writeCommittedSkill(t, root, "socrates-held", skill.FrontMatter(skill.ClaudeCode, "socrates-held", "x")+"\n")

	invocable, err := convention.SkillsAModelCouldInvoke(root)
	if err != nil {
		t.Fatalf("reading the committed skills: %v", err)
	}
	if len(invocable) != 1 || invocable[0].Name != "socrates-new" {
		t.Fatalf("the check reported %v, want only socrates-new", invocable)
	}
}

func TestTheFieldInTheBodyOfASkillDoesNotCountAsFrontMatter(t *testing.T) {
	root := t.TempDir()
	writeCommittedSkill(t, root, "socrates-body", "---\nname: socrates-body\n---\n\n"+skill.ClaudeCode.PersonOnlyField+"\n")

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
	writeCommittedSkill(t, root, "release", "---\nname: release\n---\n")
	if err := os.MkdirAll(skill.ClaudeCode.SkillsDirectoryUnder(root), 0o755); err != nil {
		t.Fatal(err)
	}

	invocable, err := convention.SkillsAModelCouldInvoke(root)
	if err != nil {
		t.Fatalf("reading the committed skills: %v", err)
	}
	if len(invocable) != 0 {
		t.Fatalf("the check reported %v", invocable)
	}
}
