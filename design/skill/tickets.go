package skill

import "github.com/socrates-tools/socrates/internal/skill"

const TicketsName = "socrates-tickets"

func TicketsDocument(host skill.Host) string {
	return assembled(skill.FrontMatter(host, TicketsName, ticketsDescription), ticketsPurpose, theStartingSection, theSlicingSection, theWavesSection,
		theBreakdownSection, theWritingThemSection, theMainTicketSection, theTicketSection, theCommittingSection, ticketsInstallingElsewhere)
}
