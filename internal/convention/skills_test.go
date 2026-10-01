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

	for _, name := range invocable {
		t.Errorf("%s does not carry %s in its front matter, so a model could start it", name, skill.PersonOnlyField)
	}
}

func writeCommittedSkill(t *testing.T, root, name, document string) {
	t.Helper()
	writePackaged(t, filepath.Join(root, ".claude", "skills", name, skill.FileName), document)
}

func TestASkillOfAnyNamespaceWithoutTheFieldIsReported(t *testing.T) {
	root := t.TempDir()
	writeCommittedSkill(t, root, "socrates-new", "---\nname: socrates-new\ndescription: x\n---\n\n# body\n")
	writeCommittedSkill(t, root, "socrates-held", skill.FrontMatter("socrates-held", "x")+"\n")

	invocable, err := convention.SkillsAModelCouldInvoke(root)
	if err != nil {
		t.Fatalf("reading the committed skills: %v", err)
	}
	if len(invocable) != 1 || invocable[0] != "socrates-new" {
		t.Fatalf("the check reported %v, want only socrates-new", invocable)
	}
}

func TestTheFieldInTheBodyOfASkillDoesNotCountAsFrontMatter(t *testing.T) {
	root := t.TempDir()
	writeCommittedSkill(t, root, "socrates-body", "---\nname: socrates-body\n---\n\n"+skill.PersonOnlyField+"\n")

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
	if err := os.MkdirAll(filepath.Join(root, ".claude", "skills"), 0o755); err != nil {
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
