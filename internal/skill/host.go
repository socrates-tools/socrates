package skill

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.yaml.in/yaml/v3"
)

const (
	frontMatterFence = "---"

	HostFlag = "host"

	PolicyDocument = "policy:\n  allow_implicit_invocation: false\n"
)

type Host struct {
	Name            string
	Flag            string
	Starter         string
	SkillsDirectory string
	AgentsDirectory string

	ConfigurationVariable    string
	SkillsUnderConfiguration bool
	PersonOnlyField          string
	PolicyFile               string
	Escalation               string
	Delegation               string
	NoSkill                  string
}

var ClaudeCode = Host{
	Name:            "Claude Code",
	Flag:            "claude",
	Starter:         "/",
	SkillsDirectory: filepath.Join(".claude", "skills"),
	AgentsDirectory: filepath.Join(".claude", "agents"),
	PersonOnlyField: "disable-model-invocation: true",

	ConfigurationVariable:    "CLAUDE_CONFIG_DIR",
	SkillsUnderConfiguration: true,
}

var Codex = Host{
	Name:            "Codex",
	Flag:            "codex",
	Starter:         "$",
	SkillsDirectory: filepath.Join(".agents", "skills"),
	AgentsDirectory: filepath.Join(".codex", "agents"),
	PolicyFile:      filepath.Join("agents", "openai.yaml"),
	Escalation:      escalatesThePage,
	Delegation:      spawnsAndWaitsForAgents,
	NoSkill:         invokesNoSkill,

	ConfigurationVariable: "CODEX_HOME",
}

func Hosts() []Host {
	return []Host{ClaudeCode, Codex}
}

func HostFlagged(value string) (Host, bool) {
	for _, host := range Hosts() {
		if host.Flag == value {
			return host, true
		}
	}
	return Host{}, false
}

func (h Host) Guidance(page string, delegates bool) string {
	var lines []string
	if h.Escalation != "" {
		lines = append(lines, fmt.Sprintf(h.Escalation, page))
	}
	if delegates && h.Delegation != "" {
		lines = append(lines, h.Delegation)
	}
	if h.NoSkill != "" {
		lines = append(lines, h.NoSkill)
	}
	if len(lines) == 0 {
		return ""
	}
	return strings.Join(append([]string{fmt.Sprintf(guidanceHeading, h.Name), ""}, lines...), "\n")
}

func (h Host) Starts(name string) string {
	return h.Starter + name
}

func (h Host) ConfigurationStated() string {
	if h.ConfigurationVariable == "" {
		return ""
	}
	return os.Getenv(h.ConfigurationVariable)
}

func (h Host) SkillsLocation() string {
	if h.SkillsUnderConfiguration {
		return h.location(h.SkillsDirectory)
	}
	return fmt.Sprintf(locationUnderTheHome, filepath.ToSlash(h.SkillsDirectory))
}

func (h Host) AgentsLocation() string {
	if h.ConfigurationVariable == "" {
		return fmt.Sprintf(locationUnderTheHome, filepath.ToSlash(h.AgentsDirectory))
	}
	return h.location(h.AgentsDirectory)
}

func (h Host) location(underTheHome string) string {
	return fmt.Sprintf(locationUnderAVariable, h.ConfigurationVariable, filepath.Base(underTheHome), filepath.ToSlash(underTheHome))
}

func (h Host) SkillsDirectoryUnder(home string) string {
	return filepath.Join(home, h.SkillsDirectory)
}

func (h Host) AgentsDirectoryUnder(home string) string {
	return filepath.Join(home, h.AgentsDirectory)
}

func (h Host) frontMatterHoldsSwitch(document string) bool {
	if h.PersonOnlyField == "" {
		return true
	}
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
	if err != nil || !h.frontMatterHoldsSwitch(string(document)) {
		return false
	}
	return h.PolicyFile == "" || forbidsImplicitInvocation(filepath.Join(directory, h.PolicyFile))
}

func forbidsImplicitInvocation(path string) bool {
	held, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	var policy struct {
		Policy struct {
			AllowImplicitInvocation *bool `yaml:"allow_implicit_invocation"`
		} `yaml:"policy"`
	}
	if err := yaml.Unmarshal(held, &policy); err != nil {
		return false
	}
	allowed := policy.Policy.AllowImplicitInvocation
	return allowed != nil && !*allowed
}
