package skill_test

import (
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/design/skill"
)

const (
	startingHeading   = "## The spec it starts from"
	slicingHeading    = "## Slicing"
	wavesHeading      = "## Blockers and waves"
	breakdownHeading  = "## Putting the breakdown on a page"
	lastCardHeading   = "## The last card"
	mainTicketHeading = "## The main ticket"
	aTicketHeading    = "## A ticket"
	committingHeading = "## Leaving them to be committed"
)

func ticketsSaysEach(t *testing.T, heading, why string, said ...string) {
	t.Helper()

	section := sectionOf(t, skill.TicketsDocument(), heading)
	for _, one := range said {
		if !strings.Contains(section, one) {
			t.Errorf("the %s section of the tickets skill does not say %q, so %s:\n%s", heading, one, why, section)
		}
	}
}

func TestTheTicketsSkillStartsFromASpecInAChangeFolderAndNothingElse(t *testing.T) {
	ticketsSaysEach(t, startingHeading, "a session could cut tickets from an argument nobody settled",
		"docs/changes/<slug>/spec.md",
		"and nothing else",
		"run a design first",
		"the ADRs its Decisions line names",
		"CONTEXT.md",
	)
}

func TestAnExistingTicketsFolderIsRefusedByName(t *testing.T) {
	ticketsSaysEach(t, startingHeading, "a second cut could overwrite what an implementation already recorded",
		"docs/changes/<slug>/tickets/",
		"already exists, write nothing",
		"name the folder",
	)
}

func TestTheTicketsSkillSlicesTracerBullets(t *testing.T) {
	ticketsSaysEach(t, slicingHeading, "a session could slice a change by layer rather than by behaviour",
		"tracer-bullet",
		"vertical slice",
		"one fresh context",
		"Prefactoring comes first",
		"wide refactor",
		"expand",
		"contract",
		"no file path and no code snippet",
	)
}

func TestATicketNamesItsBlockersAndItsWaveIsDerivedFromThem(t *testing.T) {
	ticketsSaysEach(t, wavesHeading, "a ticket's name and its blockers could disagree",
		"names the tickets that block it",
		"one more than the highest wave among its blockers",
		"1 for a ticket with none",
		"N-<task>.md",
		"N.1-<task>.md",
		"gates nothing",
	)
}

func TestANameIsFixedWhenTheTicketsAreWritten(t *testing.T) {
	ticketsSaysEach(t, wavesHeading, "abandoning a ticket could rename others and break every link to them",
		"A name is fixed once the tickets are written",
		"renames nothing",
	)
}

func TestTheBreakdownIsApprovedOnAPageWithOneCardPerTicketAndKeepItChosen(t *testing.T) {
	ticketsSaysEach(t, breakdownHeading, "agreeing with a breakdown could cost more than one press",
		"one card per proposed ticket",
		`"label": "keep it"`,
		`"recommended": true`,
		"words field",
		"merge, a split or a different blocker",
		"socrates agora serve --session <session>",
		"socrates agora poll --session <session>",
		"Never put a breakdown to a person in the terminal",
	)
}

func TestABreakdownAfterADesignIsServedToThePageThePersonAlreadyHas(t *testing.T) {
	ticketsSaysEach(t, breakdownHeading, "a design that handed on could cost the person a second address",
		"Where a design handed on to this skill, serve the breakdown to the design's session",
		"mint a new one only where there was none, or where that exchange ended or the session was collected",
	)
}

func TestWritingIsApprovedOnOneLastCardNamingTheFolderAndEveryFile(t *testing.T) {
	ticketsSaysEach(t, lastCardHeading, "a person would discover where the tickets went rather than approve it",
		"every card comes back kept",
		"one card, and it is the only card on that page",
		"every file the folder will hold",
		"docs/changes/<slug>/tickets/",
		"Nothing is written until that card is approved",
	)
}

func TestTheMainTicketIsTheReadmeAndPointsAtTheSpecRatherThanRepeatingIt(t *testing.T) {
	ticketsSaysEach(t, mainTicketHeading, "an implementer could be handed a second copy of the design",
		"README.md",
		"outside the numbering",
		"points at the spec and each ADR it names",
		"repeats no decision and argues for none",
	)
}

func TestStatusLivesOnlyInTheMainTicketInFourStates(t *testing.T) {
	ticketsSaysEach(t, mainTicketHeading, "an implementing session could lose track of what it handed out",
		"the one place a ticket's status is kept",
		"pending",
		"in progress",
		"done",
		"abandoned",
		"with a reason",
		"Mark a ticket in progress before handing it out",
		"once its acceptance criteria hold and its work is committed",
	)
	ticketsSaysEach(t, aTicketHeading, "a ticket's status could be kept in two places",
		"carries no status",
	)
}

func TestAnAbandonedTicketsBlockersPassToTheTicketsItBlocked(t *testing.T) {
	ticketsSaysEach(t, mainTicketHeading, "a ticket that waited on an abandoned one could start before what it waited on through it",
		"every ticket it blocked takes its blockers in its place",
	)
}

func TestTheMainTicketSaysToCommitBeforeAnythingIsHandedOut(t *testing.T) {
	ticketsSaysEach(t, mainTicketHeading, "an implementer in a worktree of its own could be handed a ticket it cannot read",
		"Commit this folder before handing out any ticket",
	)
	ticketsSaysEach(t, committingHeading, "the skill could commit on the person's behalf",
		"Leave the tickets uncommitted",
		"committing is the person's",
	)
}

func TestATicketsAcceptanceCriteriaNameNoTest(t *testing.T) {
	ticketsSaysEach(t, aTicketHeading, "a ticket could carry something code-shaped that goes stale",
		"acceptance criteria",
		"names no test and no scenario",
	)
}

func TestTheTicketsSkillNamesNoReview(t *testing.T) {
	if document := skill.TicketsDocument(); strings.Contains(document, "socrates review") {
		t.Error("the tickets skill names socrates review, and how a ticket is implemented is the implementing session's own")
	}
}

func TestTheTicketsSkillStatesItsSectionsInTheOrderTheWorkHappens(t *testing.T) {
	document := skill.TicketsDocument()

	previous := -1
	for _, heading := range []string{startingHeading, slicingHeading, wavesHeading, breakdownHeading, lastCardHeading, mainTicketHeading, aTicketHeading, committingHeading} {
		at := strings.Index(document, "\n"+heading+"\n")
		if at < 0 {
			t.Fatalf("the tickets skill has no %s section", heading)
		}
		if at < previous {
			t.Errorf("the %s section is stated before the section the work reaches ahead of it", heading)
		}
		previous = at
	}
}
