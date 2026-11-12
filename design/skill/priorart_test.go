package skill_test

import (
	skillkit "github.com/socrates-tools/socrates/internal/skill"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/design/skill"
)

const seamsHeading = "## Seams and alternatives"

func designSaysEach(t *testing.T, heading, why string, said ...string) {
	t.Helper()

	section := sectionOf(t, skill.Document(skillkit.ClaudeCode), heading)
	for _, one := range said {
		if !strings.Contains(section, one) {
			t.Errorf("the %s section of the design skill does not say %q, so %s:\n%s", heading, one, why, section)
		}
	}
}

func TestALookupStillRunningHoldsBackOnlyTheQuestionsDownstreamOfIt(t *testing.T) {
	designSaysEach(t, factsHeading, "a slow lookup could hold back a whole round",
		"Look facts up in subagents, while the round goes out",
		"A lookup still running is an unsettled prerequisite",
		"only the questions downstream of it wait for it",
	)
}

func TestALookupStartsAFreshSubagentRatherThanAForkOfTheSession(t *testing.T) {
	designSaysEach(t, factsHeading, "a lookup could inherit the whole interview and answer from it rather than from the repository",
		"Start each one fresh rather than as a fork of this session",
	)
}

func TestTheGlossaryIsChallengedRatherThanOnlyChecked(t *testing.T) {
	designSaysEach(t, glossaryHeading, "a fuzzy word or a contradicted claim could reach the spec unchallenged",
		"Challenge the words rather than only checking them, each on a card of its own",
		"a fuzzy word, offered as the things it could mean",
		"an invented scenario at the boundary between two terms",
		"a claim the code contradicts",
	)
}

func TestAGlossaryEntrySaysWhatATermIsAndNamesTheSynonymsItGivesUp(t *testing.T) {
	designSaysEach(t, glossaryHeading, "an entry could drift into implementation detail or leave synonyms in use",
		"says what a term is rather than what it does",
		"holds no implementation detail",
		"names the synonyms it gives up",
	)
}

func TestACardAsksWhereTheChangeIsTestedPreferringTheHighestExistingSeam(t *testing.T) {
	designSaysEach(t, seamsHeading, "the spec's testing decisions could be the session's guess",
		"a card asks where it is tested",
		"existing seams before new ones",
		"the highest seam that can see the behaviour, ideally one",
		"Testing Decisions",
	)
}

func TestEachSeamOfferedSaysWhatATestThereCatchesAndWhatItMisses(t *testing.T) {
	designSaysEach(t, seamsHeading, "a person could choose a seam by its name without seeing what it trades",
		"what a test there catches and what it misses",
	)
}

func TestAnInterfaceDraftComesFromAFreshSubagentRatherThanAFork(t *testing.T) {
	designSaysEach(t, seamsHeading, "every draft could carry the session's own leaning and stop really differing",
		"have fresh subagents, not forks of this session",
	)
}

func TestAnInterfaceQuestionOffersDraftsMadeUnderDifferentConstraints(t *testing.T) {
	designSaysEach(t, seamsHeading, "a card could offer only the first design the session thought of",
		"each draft one under a different constraint",
		"offer the drafts that really differ as the card's options",
		"The Rejected section is built from those options",
	)
}

func TestAPrototypeAnswersAQuestionWordsCannotSettleOutsideTheRepository(t *testing.T) {
	designSaysEach(t, seamsHeading, "a prototype could be left in the repository or never built at all",
		"Where a card cannot settle a question in words",
		"build a prototype that answers that one question",
		"Write it outside the repository",
		"remove it once the question is settled",
		"marked as coming from the prototype",
	)
}

func TestTheSeamsSectionSitsBetweenTheGlossaryAndTheSpec(t *testing.T) {
	document := skill.Document(skillkit.ClaudeCode)

	glossary := strings.Index(document, glossaryHeading)
	seams := strings.Index(document, seamsHeading)
	spec := strings.Index(document, specHeading)
	if glossary < 0 || seams < 0 || spec < 0 {
		t.Fatalf("the skill is missing one of %q, %q or %q", glossaryHeading, seamsHeading, specHeading)
	}
	if glossary > seams || seams > spec {
		t.Error("the seams section is not stated after the glossary section and before the spec section")
	}
}

func TestTheDecisionsSectionNamesWhereAllThreeOftenHoldWithoutLooseningTheTest(t *testing.T) {
	designSaysEach(t, decisionsHeading, "a scope decision or a deliberate deviation could go unrecorded",
		"All three often hold for a scope decision",
		"a deliberate deviation the next engineer would otherwise fix",
		"a constraint the code does not show",
		"Each still has to pass all three",
	)
}
