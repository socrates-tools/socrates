package skill

import (
	"os"
	"path/filepath"
	"strings"
)

const frontMatterFence = "---"

type Host struct {
	Name                   string
	ConfigurationDirectory string
	Starter                string
	SkillsDirectory        string
	AgentsDirectory        string
	PersonOnlyField        string
}

var ClaudeCode = Host{
	Name:                   "Claude Code",
	ConfigurationDirectory: ".claude",
	Starter:                "/",
	SkillsDirectory:        filepath.Join(".claude", "skills"),
	AgentsDirectory:        filepath.Join(".claude", "agents"),
	PersonOnlyField:        "disable-model-invocation: true",
}

func Hosts() []Host {
	return []Host{ClaudeCode}
}

func (h Host) Starts(name string) string {
	return h.Starter + name
}

func (h Host) SkillsDirectoryUnder(home string) string {
	return filepath.Join(home, h.SkillsDirectory)
}

func (h Host) AgentsDirectoryUnder(home string) string {
	return filepath.Join(home, h.AgentsDirectory)
}

func (h Host) CarriesPersonOnly(document string) bool {
	lines := strings.Split(document, "\n")
	if len(lines) == 0 || lines[0] != frontMatterFence {
		return false
	}
	for _, one := range lines[1:] {
		if one == frontMatterFence {
			return false
		}
		if one == h.PersonOnlyField {
			return true
		}
	}
	return false
}

func (h Host) PersonOnly(directory string) bool {
	document, err := os.ReadFile(filepath.Join(directory, FileName))
	return err == nil && h.CarriesPersonOnly(string(document))
}
