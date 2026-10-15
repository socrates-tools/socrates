package skill

import (
	"fmt"

	"github.com/socrates-tools/socrates/internal/skill"
)

const (
	Name    = "socrates-design"
	Surface = "agora"
)

var (
	Ceilings        = skill.Ceilings{skill.ClaudeCode: 13300, skill.Codex: 13662}
	TicketsCeilings = skill.Ceilings{skill.ClaudeCode: 9800, skill.Codex: 10017}
)

var Namespaces = map[string]string{
	"design": Name,
	"agora":  Surface,
}

func Binaries() []string {
	return []string{Name, Surface}
}

func Document(host skill.Host) string {
	return assembled(skill.FrontMatter(host, Name, description), purpose, host.Guidance(designPage, true), theRoundsSection, pageTemplate, theSettledSection, theFactsSection,
		theGlossarySection, theSeamsSection, theSpecSection, theDecisionsSection, theFrontierEmptiesSection, theEndingSection, installingElsewhere(host))
}

func assembled(parts ...string) string {
	return skill.Assembled(parts...)
}

func installingElsewhere(host skill.Host) string {
	return fmt.Sprintf(installingElsewhereTemplate, host.SkillsLocation(), host.Flag)
}

func For(host skill.Host) []skill.Skill {
	return []skill.Skill{
		{Name: Name, Host: host, Document: Document, Ceiling: Ceilings.For(host)},
		{Name: TicketsName, Host: host, Document: TicketsDocument, Ceiling: TicketsCeilings.For(host)},
	}
}

func All() []skill.Skill {
	var every []skill.Skill
	for _, host := range skill.Hosts() {
		every = append(every, For(host)...)
	}
	return every
}
