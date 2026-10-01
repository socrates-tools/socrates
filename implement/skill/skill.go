package skill

import (
	"os"
	"path/filepath"

	"github.com/socrates-tools/socrates/implement/config"
	"github.com/socrates-tools/socrates/internal/refusal"
	skillkit "github.com/socrates-tools/socrates/internal/skill"
)

const (
	AgentName       = "socrates-implementer"
	AgentCeiling    = 1400
	AgentFileSuffix = ".md"

	vendorDirectory = ".claude"
	agentsDirectory = "agents"
)

func AgentFileName() string {
	return AgentName + AgentFileSuffix
}

func Definition(model, effort string) string {
	return "---\n" +
		"name: " + AgentName + "\n" +
		"description: " + agentDescription + "\n" +
		"model: " + model + "\n" +
		"effort: " + effort + "\n" +
		"background: true\n" +
		"---\n\n" +
		standingInstructions + "\n"
}

func DefaultDefinition() string {
	return Definition(config.DefaultModel, config.DefaultEffort)
}

func DefinitionFor(stated config.Config) string {
	role := stated.Implementer()
	return Definition(role.Model, role.Effort)
}

func AgentsDirectory() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", refusal.Wrapping(err, noHomeDirectory, err)
	}
	return AgentsDirectoryUnder(home), nil
}

func AgentsDirectoryUnder(home string) string {
	return filepath.Join(home, vendorDirectory, agentsDirectory)
}

func Ceilings() []skillkit.Skill {
	return []skillkit.Skill{{Name: AgentName, Document: DefaultDefinition, Ceiling: AgentCeiling}}
}

func WriteDefinition(directory string, stated config.Config) (string, error) {
	path := filepath.Join(directory, AgentFileName())

	if err := os.MkdirAll(directory, 0o755); err != nil {
		return "", refusal.Wrapping(err, theDefinitionCannotBeWritten, path, err)
	}
	if err := os.WriteFile(path, []byte(DefinitionFor(stated)), 0o644); err != nil {
		return "", refusal.Wrapping(err, theDefinitionCannotBeWritten, path, err)
	}
	return path, nil
}

func Install(directory string, stated config.Config) (string, error) {
	absolute, err := filepath.Abs(directory)
	if err != nil {
		return "", refusal.Wrapping(err, theTargetCannotBeMadeAbsolute, directory, err)
	}

	path, err := WriteDefinition(absolute, stated)
	if err != nil {
		return "", err
	}
	if err := Record(absolute); err != nil {
		return "", err
	}
	return path, nil
}

func Reinstall(stated config.Config) (string, bool, error) {
	directory, installed, err := Recorded()
	if err != nil || !installed {
		return "", false, err
	}
	path, err := WriteDefinition(directory, stated)
	return path, err == nil, err
}
