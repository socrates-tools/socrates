package skill

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/socrates-tools/socrates/internal/refusal"
)

const FileName = "SKILL.md"

func FrontMatter(host Host, name, description string) string {
	lines := []string{frontMatterFence, "name: " + name, "description: " + description}
	if host.PersonOnlyField != "" {
		lines = append(lines, host.PersonOnlyField)
	}
	return strings.Join(append(lines, frontMatterFence), "\n")
}

type Skill struct {
	Name     string
	Host     Host
	Document func(Host) string
	Ceiling  int
}

func (s Skill) Rendered() string {
	return s.Document(s.Host)
}

type Installed struct {
	Name string
	Host Host
	Path string
}

func (i Installed) Announcement() string {
	return fmt.Sprintf(installedForAPerson, i.Name, i.Path, i.Host.Starts(i.Name))
}

func UserDirectory(host Host, invocation string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", refusal.Wrapping(err, noHomeDirectory, err, invocation)
	}
	return UserDirectoryUnder(host, home), nil
}

func UserDirectoryUnder(host Host, home string) string {
	return host.SkillsDirectoryUnder(home)
}

func Install(directory string, skills []Skill, invocation string) ([]Installed, error) {
	var written []Installed

	for _, one := range skills {
		path := filepath.Join(directory, one.Name, FileName)

		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return nil, refusal.Wrapping(err, theTargetCannotBeWritten, path, err, invocation)
		}
		if err := os.WriteFile(path, []byte(one.Rendered()), 0o644); err != nil {
			return nil, refusal.Wrapping(err, theTargetCannotBeWritten, path, err, invocation)
		}
		if one.Host.PolicyFile != "" {
			policy := filepath.Join(filepath.Dir(path), one.Host.PolicyFile)
			if err := os.MkdirAll(filepath.Dir(policy), 0o755); err != nil {
				return nil, refusal.Wrapping(err, theTargetCannotBeWritten, policy, err, invocation)
			}
			if err := os.WriteFile(policy, []byte(PolicyDocument), 0o644); err != nil {
				return nil, refusal.Wrapping(err, theTargetCannotBeWritten, policy, err, invocation)
			}
		}
		written = append(written, Installed{Name: one.Name, Host: one.Host, Path: path})
	}

	return written, nil
}
