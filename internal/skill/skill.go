package skill

import (
	"os"
	"path/filepath"

	"github.com/socrates-tools/socrates/internal/refusal"
)

const (
	FileName = "SKILL.md"

	vendorDirectory = ".claude"
	skillsDirectory = "skills"
)

type Skill struct {
	Name     string
	Document func() string
	Ceiling  int
}

type Installed struct {
	Name string
	Path string
}

func UserDirectory(invocation string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", refusal.Wrapping(err, noHomeDirectory, err, invocation)
	}
	return UserDirectoryUnder(home), nil
}

func UserDirectoryUnder(home string) string {
	return filepath.Join(home, vendorDirectory, skillsDirectory)
}

func Install(directory string, skills []Skill, invocation string) ([]Installed, error) {
	var written []Installed

	for _, one := range skills {
		path := filepath.Join(directory, one.Name, FileName)

		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return nil, refusal.Wrapping(err, theTargetCannotBeWritten, path, err, invocation)
		}
		if err := os.WriteFile(path, []byte(one.Document()), 0o644); err != nil {
			return nil, refusal.Wrapping(err, theTargetCannotBeWritten, path, err, invocation)
		}
		written = append(written, Installed{Name: one.Name, Path: path})
	}

	return written, nil
}
