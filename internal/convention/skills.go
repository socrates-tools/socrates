package convention

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/socrates-tools/socrates/internal/skill"
)

const socratesSkillPrefix = "socrates-"

type InvocableSkill struct {
	Host skill.Host
	Name string
}

func SkillsAModelCouldInvoke(root string) ([]InvocableSkill, error) {
	var invocable []InvocableSkill

	for _, host := range skill.Hosts() {
		directory := filepath.Join(root, host.SkillsDirectory)
		entries, err := os.ReadDir(directory)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			if !entry.IsDir() || !strings.HasPrefix(entry.Name(), socratesSkillPrefix) {
				continue
			}
			if !host.PersonOnly(filepath.Join(directory, entry.Name())) {
				invocable = append(invocable, InvocableSkill{Host: host, Name: entry.Name()})
			}
		}
	}
	return invocable, nil
}
