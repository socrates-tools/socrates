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

func Document(host skill.Host) string {
	return assembled(skill.FrontMatter(host, Name, description), purpose, theRoundsSection, pageTemplate, theSettledSection, theFactsSection,
		theGlossarySection, theSeamsSection, theSpecSection, theDecisionsSection, theFrontierEmptiesSection, theEndingSection, installingElsewhere)
}

func assembled(parts ...string) string {
	return strings.ReplaceAll(strings.Join(parts, "\n\n")+"\n", skill.InlineCodeMarker, skill.CodeMarker)
}

func For(host skill.Host) []skill.Skill {
	return []skill.Skill{
		{Name: Name, Host: host, Document: Document, Ceiling: Ceiling},
		{Name: TicketsName, Host: host, Document: TicketsDocument, Ceiling: TicketsCeiling},
	}
}

func All() []skill.Skill {
	var every []skill.Skill
	for _, host := range skill.Hosts() {
		every = append(every, For(host)...)
	}
	return every
}
