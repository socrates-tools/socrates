package convention

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/socrates-tools/socrates/internal/skill"
)

const (
	committedSkillsDirectory = ".claude/skills"
	socratesSkillPrefix      = "socrates-"
	frontMatterFence         = "---"
)

func SkillsAModelCouldInvoke(root string) ([]string, error) {
	entries, err := os.ReadDir(filepath.Join(root, filepath.FromSlash(committedSkillsDirectory)))
	if err != nil {
		return nil, err
	}

	var invocable []string
	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), socratesSkillPrefix) {
			continue
		}
		document, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(committedSkillsDirectory), entry.Name(), skill.FileName))
		if err != nil {
			return nil, err
		}
		if !frontMatterCarries(string(document), skill.PersonOnlyField) {
			invocable = append(invocable, entry.Name())
		}
	}
	return invocable, nil
}

func frontMatterCarries(document, line string) bool {
	lines := strings.Split(document, "\n")
	if len(lines) == 0 || lines[0] != frontMatterFence {
		return false
	}
	for _, one := range lines[1:] {
		if one == frontMatterFence {
			return false
		}
		if one == line {
			return true
		}
	}
	return false
}
