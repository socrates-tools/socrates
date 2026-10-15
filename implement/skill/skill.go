package skill

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/socrates-tools/socrates/implement/config"
	"github.com/socrates-tools/socrates/internal/refusal"
	skillkit "github.com/socrates-tools/socrates/internal/skill"
)

const (
	Name    = "socrates-implement"
	Ceiling = 10500

	AgentName       = "socrates-implementer"
	AgentCeiling    = 1700
	AgentFileSuffix = ".md"

	CodexRoleCeiling    = 1750
	CodexRoleFileSuffix = ".toml"
)

var Namespaces = map[string]string{"implement": Name}

func Binaries() []string {
	return []string{Name}
}

func Document(host skillkit.Host) string {
	parts := []string{skillkit.FrontMatter(host, Name, description), purpose, theWritesNoCodeSection, theStartingSection, theDefaultBranchSection,
		theChangesBranchSection, theFrontierSection, theWorktreesSection, theReviewSection, theLandingSection,
		theStatusSection, theResumingSection, thePersonSection, theLastReadSection, installingElsewhere(host)}
	return strings.ReplaceAll(strings.Join(parts, "\n\n")+"\n", skillkit.InlineCodeMarker, skillkit.CodeMarker)
}

func installingElsewhere(host skillkit.Host) string {
	return fmt.Sprintf(installingElsewhereTemplate, filepath.ToSlash(host.SkillsDirectory), AgentKind(host), whereAgentsInstall(host), host.Flag)
}

func AgentKind(host skillkit.Host) string {
	return agentKinds[host]
}

func whereAgentsInstall(host skillkit.Host) string {
	under := "~/" + filepath.ToSlash(host.AgentsDirectory)
	if host.ConfigurationVariable == "" {
		return fmt.Sprintf(codeSpan, under)
	}
	return fmt.Sprintf(agentsUnderAVariable, host.ConfigurationVariable, filepath.Base(host.AgentsDirectory), under)
}

func For(host skillkit.Host) []skillkit.Skill {
	return []skillkit.Skill{{Name: Name, Host: host, Document: Document, Ceiling: Ceiling}}
}

func All() []skillkit.Skill {
	var every []skillkit.Skill
	for _, host := range skillkit.Hosts() {
		every = append(every, For(host)...)
	}
	return every
}

func AgentFileName() string {
	return AgentName + AgentFileSuffix
}

func CodexRoleFileName() string {
	return AgentName + CodexRoleFileSuffix
}

func AgentFileNameFor(host skillkit.Host) string {
	if host == skillkit.Codex {
		return CodexRoleFileName()
	}
	return AgentFileName()
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

func CodexRole(model, effort string) string {
	lines := []string{
		"name = " + strconv.Quote(AgentName),
		"description = " + strconv.Quote(agentDescription),
	}
	if model != "" {
		lines = append(lines, "model = "+strconv.Quote(model))
	}
	lines = append(lines,
		"model_reasoning_effort = "+strconv.Quote(effort),
		"developer_instructions = '''",
		standingInstructions,
		"'''",
	)
	return strings.Join(lines, "\n") + "\n"
}

func DefaultDefinition() string {
	return Definition(config.DefaultModel, config.DefaultEffort)
}

func DefaultCodexRole() string {
	return CodexRole("", config.DefaultCodexEffort)
}

func DefinitionFor(stated config.Config) string {
	role := stated.Implementer()
	return Definition(role.Model, role.Effort)
}

func AgentsDirectory(host skillkit.Host) (string, error) {
	if stated := os.Getenv(host.ConfigurationVariable); host.ConfigurationVariable != "" && stated != "" {
		return filepath.Join(stated, filepath.Base(host.AgentsDirectory)), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", refusal.Wrapping(err, noHomeDirectory, err)
	}
	return AgentsDirectoryUnder(host, home), nil
}

func AgentsDirectoryUnder(host skillkit.Host, home string) string {
	return host.AgentsDirectoryUnder(home)
}

func defaultAgentFor(host skillkit.Host) string {
	if host == skillkit.Codex {
		return DefaultCodexRole()
	}
	return DefaultDefinition()
}

func Ceilings() []skillkit.Skill {
	return append(All(),
		skillkit.Skill{Name: AgentName, Host: skillkit.ClaudeCode, Document: defaultAgentFor, Ceiling: AgentCeiling},
		skillkit.Skill{Name: AgentName, Host: skillkit.Codex, Document: defaultAgentFor, Ceiling: CodexRoleCeiling},
	)
}

func agentFor(host skillkit.Host, stated config.Config) string {
	if host == skillkit.Codex {
		return DefaultCodexRole()
	}
	return DefinitionFor(stated)
}

func writeAgent(host skillkit.Host, directory string, stated config.Config) (string, error) {
	path := filepath.Join(directory, AgentFileNameFor(host))

	if err := os.MkdirAll(directory, 0o755); err != nil {
		return "", refusal.Wrapping(err, theDefinitionCannotBeWritten, path, err)
	}
	if err := os.WriteFile(path, []byte(agentFor(host, stated)), 0o644); err != nil {
		return "", refusal.Wrapping(err, theDefinitionCannotBeWritten, path, err)
	}
	return path, nil
}

func Install(host skillkit.Host, directory string, stated config.Config) (string, error) {
	absolute, err := filepath.Abs(directory)
	if err != nil {
		return "", refusal.Wrapping(err, theTargetCannotBeMadeAbsolute, directory, err)
	}

	path, err := writeAgent(host, absolute, stated)
	if err != nil {
		return "", err
	}
	if err := Record(host, absolute); err != nil {
		return "", err
	}
	return path, nil
}

func Reinstall(host skillkit.Host, stated config.Config) (string, bool, error) {
	directory, installed, err := Recorded(host)
	if err != nil || !installed {
		return "", false, err
	}
	path, err := writeAgent(host, directory, stated)
	return path, err == nil, err
}
