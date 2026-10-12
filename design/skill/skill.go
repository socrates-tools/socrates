package skill

import (
	"strings"

	"github.com/socrates-tools/socrates/internal/skill"
)

const (
	Name    = "socrates-design"
	Surface = "agora"

	Ceiling        = 13300
	TicketsCeiling = 9800
)

var Namespaces = map[string]string{
	"design": Name,
	"agora":  Surface,
}

func Binaries() []string {
	return []string{Name, Surface}
}

func Document() string {
	return assembled(frontMatter, purpose, theRoundsSection, pageTemplate, theSettledSection, theFactsSection,
		theGlossarySection, theSeamsSection, theSpecSection, theDecisionsSection, theFrontierEmptiesSection, theEndingSection, installingElsewhere)
}

func assembled(parts ...string) string {
	return strings.ReplaceAll(strings.Join(parts, "\n\n")+"\n", skill.InlineCodeMarker, skill.CodeMarker)
}

func All() []skill.Skill {
	return []skill.Skill{
		{Name: Name, Document: Document, Ceiling: Ceiling},
		{Name: TicketsName, Document: TicketsDocument, Ceiling: TicketsCeiling},
	}
}
