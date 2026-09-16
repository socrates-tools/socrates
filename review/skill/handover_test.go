package skill_test

import (
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/review/dialogue"
	"github.com/socrates-tools/socrates/review/runstate"
	"github.com/socrates-tools/socrates/review/skill"
)

func TestTheHandoverSkillIsASkillOfItsOwn(t *testing.T) {
	handover := skill.HandoverDocument()

	if handover == skill.Document() {
		t.Fatal("the two skills render the same document, so a session reviewing its own work reads the forwarding half as well")
	}
	if !strings.Contains(handover, "name: "+skill.HandoverName) {
		t.Errorf("the handover skill does not declare itself %s:\n%s", skill.HandoverName, handover)
	}
	if strings.Contains(skill.Document(), skill.HandoverName) {
		t.Error("the skill a session reads about its own change names the handover skill, and what it costs that session is the prose it does not need")
	}
}

func TestTheHandoverSkillAnswersFindingsByItsOwnRules(t *testing.T) {
	handover := skill.HandoverDocument()

	for _, rule := range dialogue.HandoverRules() {
		if !strings.Contains(handover, rule) {
			t.Errorf("the handover skill does not state the rule %q in the same words", rule)
		}
	}
	for _, absent := range []dialogue.Resolution{dialogue.Fixed, dialogue.Decided} {
		if strings.Contains(handover, "--"+string(absent)) {
			t.Errorf("the handover skill names --%s somewhere, and every finding of such a round is forwarded, dropped or rejected:\n%s", absent, handover)
		}
	}
	for _, named := range []dialogue.Resolution{dialogue.Forwarded, dialogue.Dropped, dialogue.Rejected} {
		if !strings.Contains(handover, "--"+string(named)) {
			t.Errorf("the handover skill does not name --%s", named)
		}
	}
}

func TestTheHandoverSkillNamesTheBaseAndSaysWhatItResolvesTo(t *testing.T) {
	handover := skill.HandoverDocument()

	for _, stated := range []string{"--base", "--not-my-change", "fork point"} {
		if !strings.Contains(handover, stated) {
			t.Errorf("the handover skill does not name %q, so the review would be of the wrong change:\n%s", stated, handover)
		}
	}
	for _, said := range []string{"neither fetches nor rebases", "socrates-review status"} {
		if !strings.Contains(handover, said) {
			t.Errorf("the handover skill does not say %q", said)
		}
	}
}

func TestTheChangeIsReadInAWorktreeOfItsOwn(t *testing.T) {
	handover := skill.HandoverDocument()

	if !strings.Contains(handover, "git worktree add") {
		t.Errorf("the handover skill does not say to check the change out in a worktree of its own, so the review would be of their branch plus this session's edits:\n%s", handover)
	}
	if !strings.Contains(handover, "uncommitted edits") {
		t.Error("the handover skill does not say why a worktree of its own is what gets reviewed")
	}
}

func TestTheHeadOfTheChangeIsCapturedAndCheckedBeforeAnythingIsPosted(t *testing.T) {
	handover := skill.HandoverDocument()

	if !strings.Contains(handover, "Record the commit") {
		t.Errorf("the handover skill does not say to record the commit the change is at:\n%s", handover)
	}
	checking := sentenceNaming(handover, "Where it moved")
	if !strings.Contains(checking, "post nothing") {
		t.Errorf("the handover skill does not refuse to post where the head of the change moved: %q", checking)
	}
}

func TestARelayedAccountIsPassedRatherThanFollowed(t *testing.T) {
	handover := skill.HandoverDocument()

	for _, stated := range []string{"verbatim", "one argument", "never run"} {
		if !strings.Contains(handover, stated) {
			t.Errorf("the handover skill does not say the author's account is %q, and it is the one thing a stranger wrote that reaches this session:\n%s", stated, handover)
		}
	}
	for _, refused := range []string{"Do not repair it", "do not summarise it"} {
		if !strings.Contains(handover, refused) {
			t.Errorf("the handover skill does not say %q about the account it relays", refused)
		}
	}
}

func TestWhereACommentGoesComesFromSocrates(t *testing.T) {
	handover := skill.HandoverDocument()

	if !strings.Contains(handover, "inside the change") {
		t.Fatalf("the handover skill says nothing about whether a finding's location is inside the change:\n%s", handover)
	}
	if !strings.Contains(handover, "findings --json") {
		t.Error("the handover skill does not say where the placement is read from")
	}
	if !strings.Contains(handover, "joins one general comment, one rather than one each") {
		t.Errorf("the handover skill does not say the findings outside the change join one comment rather than one each:\n%s", handover)
	}
	if !strings.Contains(handover, "Say on each card which of the two acts forwarding it would be") {
		t.Errorf("the handover skill does not say the card tells the person which of the two forwarding a finding would be:\n%s", handover)
	}
	if !strings.Contains(handover, "treat that finding as outside the change") {
		t.Error("the handover skill does not say what to do where socrates-review states no location")
	}
}

func TestNothingIsPostedThatThePersonDidNotChoose(t *testing.T) {
	handover := skill.HandoverDocument()

	if !strings.Contains(handover, "Post nothing the person did not choose on the page") {
		t.Errorf("the handover skill does not say that nothing is posted the person did not choose:\n%s", handover)
	}

	for _, stated := range []string{"Read every comment body back from", "findings --json", "what socrates-review recorded"} {
		if !strings.Contains(handover, stated) {
			t.Errorf("the handover skill does not say %q, so a comment body could be anything the page handed back:\n%s", stated, handover)
		}
	}
	if !strings.Contains(handover, "which findings go and never a word of what goes") {
		t.Error("the handover skill does not say what the page decides and what it does not")
	}
}

func TestADropsReasonStaysOnThisMachine(t *testing.T) {
	handover := skill.HandoverDocument()

	kept := sentenceNaming(handover, "leaves this machine")
	if !strings.Contains(kept, "dropped") && !strings.Contains(kept, "reason") {
		t.Errorf("the handover skill does not say a drop's reason stays here: %q", kept)
	}
	if !strings.Contains(handover, "posted nowhere") {
		t.Error("the handover skill does not say a dropped or rejected finding reaches nobody")
	}
}

func TestNoSecondRoundFollowsAHandover(t *testing.T) {
	said := sectionAbout(t, skill.HandoverDocument(), string(runstate.AwaitingResponse))

	if !strings.Contains(said, "socrates-review continue") {
		t.Errorf("the handover skill does not say to run socrates-review continue once every finding is answered:\n%s", said)
	}
	if !strings.Contains(said, "No second round follows") {
		t.Errorf("the handover skill does not say that no second round follows:\n%s", said)
	}
	if !strings.Contains(said, string(runstate.Escalated)) {
		t.Errorf("the handover skill does not say the run ends %s:\n%s", runstate.Escalated, said)
	}
}

func TestARepeatReviewReadsWhatIsAlreadyThereAndJudgesNoDuplicate(t *testing.T) {
	said := sectionAbout(t, skill.HandoverDocument(), string(runstate.AwaitingResponse))

	for _, stated := range []string{"comments already on the change", "file and line only", "evidence"} {
		if !strings.Contains(said, stated) {
			t.Errorf("the handover skill does not say to %q before drafting the page:\n%s", stated, said)
		}
	}
	if !strings.Contains(said, "withhold none") {
		t.Errorf("the handover skill does not refuse to withhold a finding it takes for a duplicate:\n%s", said)
	}
}

func TestEveryStateAHandoverRunRestsInIsNamed(t *testing.T) {
	handover := skill.HandoverDocument()

	for _, state := range runstate.States() {
		if !state.Terminal() && state != runstate.AwaitingResponse {
			continue
		}
		if !strings.Contains(handover, sectionHeading(string(state))) {
			t.Errorf("a run rests in %s and the handover skill says nothing about it", state)
		}
	}
}

func TestEveryBranchOfTheHandoverThatHasAResultServesAPage(t *testing.T) {
	handover := skill.HandoverDocument()

	for _, state := range []runstate.State{runstate.Clean, runstate.Escalated, runstate.AwaitingResponse} {
		said := sectionAbout(t, handover, string(state))
		if !strings.Contains(said, "agora serve") {
			t.Errorf("the %s branch of the handover skill hands nothing over as a page:\n%s", state, said)
		}
	}
	for _, state := range []runstate.State{runstate.Failed, runstate.Cancelled} {
		said := sectionAbout(t, handover, string(state))
		if strings.Contains(said, "agora serve") {
			t.Errorf("the %s branch of the handover skill serves a page, and nothing was reviewed to put on one:\n%s", state, said)
		}
		if !strings.Contains(said, "post nothing") {
			t.Errorf("the %s branch of the handover skill does not say to post nothing:\n%s", state, said)
		}
	}
}

func TestThePageIsHeldToTheSameBarInBothSkills(t *testing.T) {
	handover := skill.HandoverDocument()

	if !strings.Contains(handover, skill.PageBar) {
		t.Errorf("the handover skill does not hold the page to the bar the other skill holds it to:\n%s", handover)
	}
	if strings.Count(handover, skill.PageBar) != 1 {
		t.Errorf("the handover skill states the bar %d times", strings.Count(handover, skill.PageBar))
	}
	if !strings.Contains(handover, "Give every finding in the round a card of its own") {
		t.Error("the handover skill does not say every finding of the round gets a card")
	}
}

func TestNoInlineMarkerSurvivesIntoEitherDocument(t *testing.T) {
	for _, one := range skill.All() {
		document := one.Document()

		if counted := strings.Count(document, "|"); counted != 0 {
			t.Errorf("%s carries %d inline markers that were never turned into code", one.Name, counted)
		}
		if counted := strings.Count(document, inlineCode); counted%2 != 0 {
			t.Errorf("%s holds %d code markers, so one code span is unclosed", one.Name, counted)
		}
	}
}
