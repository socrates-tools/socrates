package skill

import (
	"strings"

	"github.com/socrates-tools/socrates/internal/skill"
)

const (
	Name    = "socrates-design"
	Surface = "agora"

	Ceiling = 1095
)

var Namespaces = map[string]string{
	"design": Name,
	"agora":  Surface,
}

func Binaries() []string {
	return []string{Name, Surface}
}

func Document() string {
	return assembled(frontMatter, purpose)
}

func assembled(parts ...string) string {
	return strings.ReplaceAll(strings.Join(parts, "\n\n")+"\n", skill.InlineCodeMarker, skill.CodeMarker)
}

func All() []skill.Skill {
	return []skill.Skill{
		{Name: Name, Document: Document, Ceiling: Ceiling},
	}
}
