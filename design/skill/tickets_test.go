package skill_test

import (
	skillkit "github.com/socrates-tools/socrates/internal/skill"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/design/skill"
)

const (
	startingHeading    = "## The spec it starts from"
	slicingHeading     = "## Slicing"
	wavesHeading       = "## Blockers and waves"
	breakdownHeading   = "## Putting the breakdown on a page"
	writingThemHeading = "## Writing them"
	mainTicketHeading  = "## The main ticket"
	aTicketHeading     = "## A ticket"
	committingHeading  = "## Leaving them to be committed"
)

func ticketsSaysEach(t *testing.T, heading, why string, said ...string) {
	t.Helper()

	section := sectionOf(t, skill.TicketsDocument(skillkit.ClaudeCode), heading)
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

func TestTheBreakdownMintsItsOwnSessionAndNamesNoDesignHandingOn(t *testing.T) {
	ticketsSaysEach(t, breakdownHeading, "the tickets skill could depend on a session a design left behind",
		"Mint the session once, with the first `socrates agora serve`",
	)
	for _, unsaid := range []string{"handed on", "hand on", "design's session", "a design that"} {
		if strings.Contains(skill.TicketsDocument(skillkit.ClaudeCode), unsaid) {
			t.Errorf("the tickets skill says %q, so it could reach for the session of a skill it cannot invoke", unsaid)
		}
	}
}

func TestTheTicketsAreWrittenOnceEveryCardIsKeptWithoutAFurtherPage(t *testing.T) {
	ticketsSaysEach(t, writingThemHeading, "a person would be asked a second time to approve a breakdown they already kept",
		"Nothing is written before a round comes back with every card kept",
		"serve no further page",
		"socrates agora end --session <session>",
		"docs/changes/<slug>/tickets/",
		"say in the terminal every path you wrote",
	)

	writing := sectionOf(t, skill.TicketsDocument(skillkit.ClaudeCode), writingThemHeading)
	if strings.Contains(writing, `"label"`) {
		t.Errorf("the writing section still composes a card for the person to approve:\n%s", writing)
	}
}

func TestEachCardCarriesTheWholeTicketItWouldBeWrittenAs(t *testing.T) {
	ticketsSaysEach(t, breakdownHeading, "a person could approve a summary while the acceptance criteria and blockers were written unseen",
		`"evidence": [{"language": "markdown", "code":`,
		"one block of the whole file it would be written as",
		"acceptance criteria, its stories and its blockers",
	)
	ticketsSaysEach(t, writingThemHeading, "a ticket could be rewritten from memory after the person approved its card",
		"exactly the evidence its card carried in that round",
	)

	breakdown := sectionOf(t, skill.TicketsDocument(skillkit.ClaudeCode), breakdownHeading)
	if strings.Contains(breakdown, "blocked by: 1-prefactor") {
		t.Errorf("the card example still states blockers as free text beside the file that states them:\n%s", breakdown)
	}
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
	if document := skill.TicketsDocument(skillkit.ClaudeCode); strings.Contains(document, "socrates review") {
		t.Error("the tickets skill names socrates review, and how a ticket is implemented is the implementing session's own")
	}
}

func TestTheTicketsSkillStatesItsSectionsInTheOrderTheWorkHappens(t *testing.T) {
	document := skill.TicketsDocument(skillkit.ClaudeCode)

	previous := -1
	for _, heading := range []string{startingHeading, slicingHeading, wavesHeading, breakdownHeading, writingThemHeading, mainTicketHeading, aTicketHeading, committingHeading} {
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

func TestAChangeThatFitsOneFreshContextIsOneTicket(t *testing.T) {
	ticketsSaysEach(t, slicingHeading, "a small change could be cut into a dozen tickets",
		"A change that fits one fresh context whole is one ticket",
		"one ticket is a breakdown",
	)
}

func TestABlockerIsNamedOnlyWhereTheTicketCannotStartWithoutIt(t *testing.T) {
	ticketsSaysEach(t, wavesHeading, "a blocker that does not gate could serialise tickets that could run side by side",
		"Name a blocker only where the ticket cannot start without its work",
	)
}

func TestTheBreakdownAsksAboutGranularityAndBlockersAndNamesAStoryNoCardDelivers(t *testing.T) {
	ticketsSaysEach(t, breakdownHeading, "a cut could be too fine, gated by blockers that do not gate, or leave a story out",
		"whether it is cut too coarse or too fine",
		"whether every blocker gates the ticket it blocks",
		"the person reads its acceptance criteria, its stories and its blockers",
		"names every user story no card delivers",
	)
}

func TestATicketNamesItsUserStoriesAndWhatItNeedsAPersonFor(t *testing.T) {
	ticketsSaysEach(t, aTicketHeading, "a ticket could deliver no story anybody can trace, or stall a run on a credential",
		"## User stories",
		"## Needs a person",
		"a credential, an account, a permission or a manual step",
		"leaves out Needs a person",
	)
}

func TestEveryAcceptanceCriterionIsFalseWhereTheTicketStarts(t *testing.T) {
	ticketsSaysEach(t, aTicketHeading, "a criterion could be ticked without anything being built",
		"names what somebody would observe were it false",
		"is false at the commit the ticket starts from",
		"already true there",
		"only another ticket's work can meet",
		"restates the title",
	)
}
