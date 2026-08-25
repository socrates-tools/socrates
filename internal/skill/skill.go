package skill

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rui-ktei/socrates/internal/dialogue"
	"github.com/rui-ktei/socrates/internal/refusal"
)

const (
	Name     = "socrates"
	Surface  = "agora"
	FileName = "SKILL.md"

	vendorDirectory = ".claude"
	skillsDirectory = "skills"
)

func Document() string {
	assembled := strings.Join([]string{frontMatter, purpose, rules(), thePage, theLoop, prohibitions}, "\n\n") + "\n"
	return strings.ReplaceAll(assembled, inlineCodeMarker, codeMarker)
}

func rules() string {
	stated := []string{rulesHeading, ""}
	for _, rule := range dialogue.Rules() {
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

func Install(directory string) (string, error) {
	path := filepath.Join(directory, Name, FileName)

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", refusal.Wrapping(err, theTargetCannotBeWritten, path, err)
	}
	if err := os.WriteFile(path, []byte(Document()), 0o644); err != nil {
		return "", refusal.Wrapping(err, theTargetCannotBeWritten, path, err)
	}
	return path, nil
}
