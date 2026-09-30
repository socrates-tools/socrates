package skill

const TicketsName = "socrates-tickets"

func TicketsDocument() string {
	return assembled(ticketsFrontMatter, ticketsPurpose, theStartingSection, theSlicingSection, theWavesSection,
		theBreakdownSection, theLastCardSection, theMainTicketSection, theTicketSection, theCommittingSection, ticketsInstallingElsewhere)
}
