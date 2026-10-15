package skill

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/socrates-tools/socrates/implement/config"
	"github.com/socrates-tools/socrates/internal/refusal"
	skillkit "github.com/socrates-tools/socrates/internal/skill"
)

const (
	Name = "socrates-implement"

	AgentName       = "socrates-implementer"
	AgentCeiling    = 1700
	AgentFileSuffix = ".md"
)

var OverseerCeilings = skillkit.Ceilings{skillkit.ClaudeCode: 10500, skillkit.Codex: 11200}

var Namespaces = map[string]string{"implement": Name}

func Binaries() []string {
	return []string{Name}
}

func Document(host skillkit.Host) string {
	wording := wordingFor(host)
	parts := []string{skillkit.FrontMatter(host, Name, description), purpose, hostSection(host), theWritesNoCodeSection,
		fmt.Sprintf(startingSection, wording.noun, wording.verb), theDefaultBranchSection, theChangesBranchSection, theFrontierSection,
		fmt.Sprintf(worktreesSection, wording.startsTheImplementer), theReviewSection,
		fmt.Sprintf(landingSection, wording.runsTheVerification), theStatusSection, theResumingSection,
		fmt.Sprintf(personSection, wording.implementersRun), theLastReadSection, installingElsewhere(host)}
	return skillkit.Assembled(parts...)
}

type hostWording struct {
	noun, verb, startsTheImplementer, runsTheVerification, implementersRun string
}

func wordingFor(host skillkit.Host) hostWording {
	if host == skillkit.Codex {
		return hostWording{codexNoun, codexVerb, codexStartsTheImplementer, codexRunsTheVerification, codexImplementersRun}
	}
	return hostWording{claudeCodeNoun, claudeCodeVerb, claudeCodeStartsTheImplementer, claudeCodeRunsTheVerification, claudeCodeImplementersRun}
}

func hostSection(host skillkit.Host) string {
	if host != skillkit.Codex {
		return ""
	}
	return fmt.Sprintf(codexSection, host.Delegation)
}

func installingElsewhere(host skillkit.Host) string {
	if host == skillkit.ClaudeCode {
		return fmt.Sprintf(installingElsewhereWithAnAgent, filepath.ToSlash(host.SkillsDirectory), filepath.ToSlash(host.AgentsDirectory), host.Flag)
	}
	return fmt.Sprintf(installingElsewhereWithoutAnAgent, filepath.ToSlash(host.SkillsDirectory), host.Flag)
}

func For(host skillkit.Host) []skillkit.Skill {
	return []skillkit.Skill{{Name: Name, Host: host, Document: Document, Ceiling: OverseerCeilings.For(host)}}
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

func AgentsDirectory(host skillkit.Host) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", refusal.Wrapping(err, noHomeDirectory, err)
	}
	return AgentsDirectoryUnder(host, home), nil
}

func AgentsDirectoryUnder(host skillkit.Host, home string) string {
	return host.AgentsDirectoryUnder(home)
}

func defaultDefinitionFor(skillkit.Host) string {
	return DefaultDefinition()
}

func Ceilings() []skillkit.Skill {
	return append(All(), skillkit.Skill{Name: AgentName, Host: skillkit.ClaudeCode, Document: defaultDefinitionFor, Ceiling: AgentCeiling})
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
