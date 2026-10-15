package skill

import (
	"fmt"

	"github.com/socrates-tools/socrates/internal/skill"
)

const TicketsName = "socrates-tickets"

func TicketsDocument(host skill.Host) string {
	return assembled(skill.FrontMatter(host, TicketsName, ticketsDescription), ticketsPurpose, host.Guidance(ticketsPage, false), theStartingSection, theSlicingSection, theWavesSection,
		theBreakdownSection, theWritingThemSection, theMainTicketSection, theTicketSection, theCommittingSection, ticketsInstallingElsewhere(host))
}

func ticketsInstallingElsewhere(host skill.Host) string {
	return fmt.Sprintf(ticketsInstallingElsewhereTemplate, host.SkillsLocation(), host.Flag)
}
