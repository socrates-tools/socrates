package skill

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/socrates-tools/socrates/internal/dialogue"
	"github.com/socrates-tools/socrates/internal/refusal"
)

const (
	Name         = "socrates-review"
	HandoverName = Name + "-handover"
	Surface      = "agora"
	FileName     = "SKILL.md"

	Ceiling         = 15616
	HandoverCeiling = 16896

	vendorDirectory = ".claude"
	skillsDirectory = "skills"
)

func Document() string {
	return assembled(frontMatter, purpose, rules(), thePage(theAnswerAPageCarries), theLoop, prohibitions)
}

func thePage(answer string) string {
	return fmt.Sprintf(pageTemplate, answer)
}

func assembled(parts ...string) string {
	return strings.ReplaceAll(strings.Join(parts, "\n\n")+"\n", inlineCodeMarker, codeMarker)
}

func rules() string {
	return stating(dialogue.Rules())
}

func stating(rules []string) string {
	stated := []string{rulesHeading, ""}
	for _, rule := range rules {
		stated = append(stated, fmt.Sprintf(ruleLine, rule))
	}
	return strings.Join(stated, "\n")
}

func UserDirectory() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", refusal.Wrapping(err, noHomeDirectory, err)
	}
	return filepath.Join(home, vendorDirectory, skillsDirectory), nil
}

type Skill struct {
	Name     string
	Document func() string
	Ceiling  int
}

func All() []Skill {
	return []Skill{
		{Name: Name, Document: Document, Ceiling: Ceiling},
		{Name: HandoverName, Document: HandoverDocument, Ceiling: HandoverCeiling},
	}
}

type Installed struct {
	Name string
	Path string
}

func Install(directory string) ([]Installed, error) {
	var written []Installed

	for _, one := range All() {
		path := filepath.Join(directory, one.Name, FileName)

		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return nil, refusal.Wrapping(err, theTargetCannotBeWritten, path, err)
		}
		if err := os.WriteFile(path, []byte(one.Document()), 0o644); err != nil {
			return nil, refusal.Wrapping(err, theTargetCannotBeWritten, path, err)
		}
		written = append(written, Installed{Name: one.Name, Path: path})
	}

	return written, nil
}
