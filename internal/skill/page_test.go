package skill_test

import (
	"strings"
	"testing"

	"github.com/ktei/socrates/internal/skill"
)

const (
	pageHeading      = "## The page a person reads"
	cleanHeading     = "### clean"
	failedHeading    = "### failed"
	escalatedHeading = "### escalated"
	parkedHeading    = "### awaiting-response"
	cancelledHeading = "### cancelled"
	waitingStep      = "poll " + inlineCode + "socrates status"
)

var branchesThatHaveAResult = []string{cleanHeading, escalatedHeading, parkedHeading}

var branchesThatHaveNone = []string{failedHeading, cancelledHeading}

func TestTheBarIsStatedOnceAndRenderedFromThere(t *testing.T) {
	document := skill.Document()

	if !strings.Contains(document, skill.PageBar) {
		t.Fatalf("the skill does not render the bar its source states:\n%s", skill.PageBar)
	}
	for line := range strings.SplitSeq(skill.PageBar, "\n") {
		stated := strings.TrimSpace(line)
		if stated == "" {
			continue
		}
		if counted := strings.Count(document, stated); counted != 1 {
			t.Errorf("the skill states %q %d times, and two copies of a bar drift", stated, counted)
		}
	}
}

func TestTheBarNamesEachWayThePageFailsAPerson(t *testing.T) {
	bar := strings.ToLower(skill.PageBar)

	for _, required := range []struct {
		named  string
		stated []string
	}{
		{named: "getting to the point", stated: []string{"get to the point", "length is not thoroughness"}},
		{named: "plain english about impact", stated: []string{"plain english", "what the finding means for the change"}},
		{named: "an option said in the words of the thing being built", stated: []string{"what each option gets the person", "words of the thing being built", "without opening the code"}},
		{named: "saying which option you would take, and marking one and no other", stated: []string{"which option you would take", "marking that option and no other"}},
		{named: "an example where a finding is purely technical", stated: []string{"purely technical", "example"}},
		{named: "showing what a finding costs where that is what a person sees", stated: []string{"what a person sees", "show them it rather than describing it"}},
		{named: "a diagram only where a workflow needs one", stated: []string{"diagram only where", "not otherwise"}},
		{named: "restraint about tokens", stated: []string{"spend tokens"}},
	} {
		t.Run(required.named, func(t *testing.T) {
			for _, stated := range required.stated {
				if !strings.Contains(bar, stated) {
					t.Errorf("the bar does not say %q, and that is how a page fails a person by %s", stated, required.named)
				}
			}
		})
	}
}

func TestTheBarIsStatedBeforeTheBranchesThatReferToIt(t *testing.T) {
	document := skill.Document()

	bar := strings.Index(document, skill.PageBar)
	if bar < 0 {
		t.Fatal("the skill does not state the bar at all")
	}
	for _, branch := range branchesThatHaveAResult {
		at := strings.Index(document, branch)
		if at < 0 {
			t.Fatalf("the skill has no %s branch", branch)
		}
		if bar > at {
			t.Errorf("the bar is stated after %s, so the branch sends a session to make a page it has not been told how to make", branch)
		}
	}

	for _, branch := range branchesThatHaveAResult {
		if strings.Contains(sectionOf(t, document, branch), "Length is not thoroughness") {
			t.Errorf("%s restates the bar, and the copy in front of the session is the one that gets softened", branch)
		}
	}
}

func TestEveryBranchThatHasAResultSaysToServeThePage(t *testing.T) {
	document := skill.Document()

	for _, branch := range branchesThatHaveAResult {
		section := strings.ToLower(sectionOf(t, document, branch))
		if !strings.Contains(section, "page") {
			t.Errorf("%s says nothing about a page, and every round that reached a result is handed over as one:\n%s", branch, section)
		}
		for _, said := range []string{"agora serve", "address"} {
			if !strings.Contains(section, said) {
				t.Errorf("%s does not say %q, so the page it drafts is never put in front of anybody:\n%s", branch, said, section)
			}
		}
	}

	for _, branch := range branchesThatHaveNone {
		section := strings.ToLower(sectionOf(t, document, branch))
		if strings.Contains(section, "page") {
			t.Errorf("%s names a page, and no review completed there, so the page would hold the one line socrates already printed:\n%s", branch, section)
		}
	}

	t.Run("The clean branch drafts inline", func(t *testing.T) {
		clean := strings.ToLower(sectionOf(t, document, cleanHeading))
		for _, said := range []string{"draft it here, in this session", "nothing follows a run that is over"} {
			if !strings.Contains(clean, said) {
				t.Errorf("the clean branch does not say %q, so nothing says why the drafting is not delegated:\n%s", said, clean)
			}
		}
		if strings.Contains(clean, "delegate") {
			t.Errorf("the clean branch delegates, and a fresh subagent costs its own start-up where nothing follows:\n%s", clean)
		}
	})

	t.Run("The escalated branch drafts inline", func(t *testing.T) {
		escalated := strings.ToLower(sectionOf(t, document, escalatedHeading))
		for _, said := range []string{"draft it here, in this session", "for the same reason", "nothing follows an escalated run"} {
			if !strings.Contains(escalated, said) {
				t.Errorf("the escalated branch does not say %q, so nothing says why the drafting is not delegated:\n%s", said, escalated)
			}
		}
		if strings.Contains(escalated, "delegate") {
			t.Errorf("the escalated branch delegates, and a fresh subagent costs its own start-up where nothing follows:\n%s", escalated)
		}
	})

	t.Run("The parked branch delegates", func(t *testing.T) {
		parked := strings.ToLower(sectionOf(t, document, parkedHeading))
		for _, said := range []string{
			"delegate",
			"subagent",
			"whichever of the two it holds",
			"read again on every one of those turns",
			"what costs you is the finding prose rather than who answers it",
			"socrates findings --json",
			"the document and a sentence",
		} {
			if !strings.Contains(parked, said) {
				t.Errorf("the parked branch does not say %q, so the drafting lands in the context that is about to patch:\n%s", said, parked)
			}
		}
	})

	t.Run("The driving session owns the surface", func(t *testing.T) {
		parked := strings.ToLower(sectionOf(t, document, parkedHeading))
		for _, said := range []string{"serve that document yourself", "agora poll", "would exit owning the surface"} {
			if !strings.Contains(parked, said) {
				t.Errorf("the parked branch does not say %q, so the delegate is left owning a surface it is about to exit:\n%s", said, parked)
			}
		}
	})
}

func TestTheSkillSaysWhichBranchesWait(t *testing.T) {
	document := skill.Document()
	page := strings.ToLower(sectionOf(t, document, pageHeading))

	for _, said := range []string{
		`prints a line beginning "awaits a person"`,
		"where it says so, wait for them",
		"where it says nothing, hand over the address and get on with the work",
	} {
		if !strings.Contains(page, said) {
			t.Errorf("the page section does not say %q, so a session waits on every page or on none:\n%s", said, page)
		}
	}

	for _, waiting := range []string{escalatedHeading, parkedHeading} {
		if !strings.Contains(strings.ToLower(sectionOf(t, document, waiting)), "agora poll") {
			t.Errorf("%s does not say to wait on the page, and socrates reports rounds of it as ones that await a person", waiting)
		}
	}
	if clean := strings.ToLower(sectionOf(t, document, cleanHeading)); strings.Contains(clean, "agora poll") {
		t.Errorf("the clean branch waits on the page, and it stops a session for an answer nobody was asked for:\n%s", clean)
	}
}

func TestTheSessionDoesNotDecideTheConditionItself(t *testing.T) {
	page := sectionOf(t, skill.Document(), pageHeading)

	if !strings.Contains(page, "awaits a person") {
		t.Error("the skill does not name the line socrates prints, so a session has to decide for itself which rounds wait on a person")
	}
	for _, invented := range []string{"unanswered", "escalated", "awaiting-response", "count of"} {
		if strings.Contains(strings.ToLower(page), invented) {
			t.Errorf("the page section names %q, which is a condition of its own over what socrates already decided:\n%s", invented, page)
		}
	}
}

func TestARoundTheAuthorAnswersAloneGetsAPageNobodyWaitsOn(t *testing.T) {
	parked := sectionOf(t, skill.Document(), parkedHeading)

	for _, said := range []string{
		"Hand the round over as a page whichever of the two it holds",
		"where it says nothing, get on with the patching and leave the page for them to read whenever they like",
	} {
		if !strings.Contains(parked, said) {
			t.Errorf("the parked branch does not say %q, so a round the author answers alone reaches nobody or stops the loop:\n%s", said, parked)
		}
	}
	for _, said := range []string{"socrates respond", "socrates continue"} {
		if !strings.Contains(parked, said) {
			t.Errorf("the parked branch no longer says %q", said)
		}
	}
}

func TestARoundWithNothingToAnswerIsAPageWithNothingToAnswer(t *testing.T) {
	page := sectionOf(t, skill.Document(), pageHeading)

	for _, said := range []string{
		"the document holds no decision at all",
		"nothing on it to answer",
		"nobody is waiting on it",
		"Do not invent a decision to fill such a page",
	} {
		if !strings.Contains(page, said) {
			t.Errorf("the page section does not say %q, so a round with nothing to answer is a page with a decision nobody asked for:\n%s", said, page)
		}
	}
}

func TestThePageIsNotInThePollingLoop(t *testing.T) {
	document := skill.Document()

	step := stepOf(t, document, waitingStep)
	if strings.Contains(strings.ToLower(step), "page") {
		t.Errorf("the step a session runs in a loop names a page, and a page is made once after the run has stopped moving:\n%s", step)
	}
	if counted := strings.Count(document, "socrates watch"); counted != 1 {
		t.Errorf("the skill names socrates watch %d times, want once", counted)
	}
}

func TestAControlCarriesItsCommandAndShowsItToNobody(t *testing.T) {
	page := strings.ToLower(sectionOf(t, skill.Document(), pageHeading))

	for _, said := range []string{
		"carries the command it becomes",
		"nothing is recorded until you run it here",
		"shows that command to nobody",
		"an option whose command you cannot name is one you do not offer",
	} {
		if !strings.Contains(page, said) {
			t.Errorf("the page section does not say %q, so a control can take a decision whose destination nobody can see:\n%s", said, page)
		}
	}
	if strings.Contains(page, "names the command it becomes") {
		t.Errorf("the page section still says a control names the command it becomes, and the person reading it is not at a prompt:\n%s", page)
	}
}

func TestTheSkillSaysHowAnAnswerReachesTheRun(t *testing.T) {
	page := strings.ToLower(sectionOf(t, skill.Document(), pageHeading))

	for _, said := range []string{"the person answers on the page", "agora poll", "socrates respond", "their own words"} {
		if !strings.Contains(page, said) {
			t.Errorf("the page section does not say %q, so an answer given on the page reaches no run:\n%s", said, page)
		}
	}
	if !strings.Contains(page, "do not ask them here for a decision they have already made there") {
		t.Errorf("the page section does not forbid asking the person to say again here what they answered there:\n%s", page)
	}
}

func TestAPageIsServedLocallyAndPublishedNowhere(t *testing.T) {
	page := strings.ToLower(sectionOf(t, skill.Document(), pageHeading))

	if !strings.Contains(page, "agora serve") {
		t.Errorf("the page section names no local surface to serve the page with:\n%s", page)
	}
	for _, elsewhere := range []string{"artifact", "upload", "publish", "host", "link", "claude.ai", "https://"} {
		if strings.Contains(page, elsewhere) {
			t.Errorf("the page section names %q, and a document quoting the code it was written about leaves this machine for nowhere:\n%s", elsewhere, page)
		}
	}
}

func TestADecisionIsStillThePersonsToMake(t *testing.T) {
	document := skill.Document()
	parked := sectionOf(t, document, parkedHeading)

	if !strings.Contains(parked, "You do not answer such a finding on your own behalf.") {
		t.Errorf("the branch holding a finding a person answers does not say the session answers none of them itself:\n%s", parked)
	}

	page := strings.ToLower(sectionOf(t, document, pageHeading))
	for _, permitted := range []string{"answer it yourself", "decide it yourself", "on their behalf"} {
		if strings.Contains(page, permitted) {
			t.Errorf("the page section says %q, and a decision arriving through a surface is no more the session's than one arriving here:\n%s", permitted, page)
		}
	}
}

func TestTheSkillSaysThePageAndTheTerminalUseTheSameIdentifiers(t *testing.T) {
	page := strings.ToLower(sectionOf(t, skill.Document(), pageHeading))

	for _, said := range []string{"identifier", "severity", "audience", "f3", "socrates findings --json"} {
		if !strings.Contains(page, said) {
			t.Errorf("the page section does not say %q, so a person naming a finding leaves the session guessing which one", said)
		}
	}
}

func TestThePageSaysWhatChoosingAnOptionMeansForTheChange(t *testing.T) {
	page := strings.ToLower(sectionOf(t, skill.Document(), pageHeading))

	for _, said := range []string{
		"name each decision from its identifier",
		"never title it yourself",
		"what choosing each option means for the change",
		"nothing about what you will type",
	} {
		if !strings.Contains(page, said) {
			t.Errorf("the page section does not say %q, so the person is handed the reviewer's filing system and the session's typing:\n%s", said, page)
		}
	}
}

const recommendingLine = "Say which option you would take, where you would take one, by marking that option and no other."

func TestTheSkillNamesTheMarkWhereItNamesTheDocument(t *testing.T) {
	page := sectionOf(t, skill.Document(), pageHeading)

	shape, _, found := strings.Cut(page, "The person answers on the page")
	if !found {
		t.Fatalf("the page section no longer shows the shape of a document before saying how it is answered:\n%s", page)
	}
	if !strings.Contains(shape, `"recommended": true`) {
		t.Errorf("the shape a session copies carries no marked option, so nothing it copies recommends anything:\n%s", shape)
	}
}

func TestTheSkillSaysOneOptionCarriesIt(t *testing.T) {
	page := strings.ToLower(sectionOf(t, skill.Document(), pageHeading))

	for _, said := range []string{"at most one", "refused before anything is served"} {
		if !strings.Contains(page, said) {
			t.Errorf("the page section does not say %q, so a session may mark two options and learn that only from a failure:\n%s", said, page)
		}
	}
}

func TestTheSkillSaysMarkingIsNotAnswering(t *testing.T) {
	page := strings.ToLower(sectionOf(t, skill.Document(), pageHeading))

	for _, said := range []string{"reads that option first", "in words", "selects nothing"} {
		if !strings.Contains(page, said) {
			t.Errorf("the page section does not say %q, so what marking an option does is left to be guessed:\n%s", said, page)
		}
	}
	if !strings.Contains(page, "still theirs") {
		t.Errorf("the page section does not say the decision remains the person's once an option is marked:\n%s", page)
	}
}

func TestWhenToRecommendIsStatedOnce(t *testing.T) {
	document := skill.Document()

	if counted := strings.Count(document, recommendingLine); counted != 1 {
		t.Fatalf("the skill states %q %d times, want once, and the bar is where it belongs", recommendingLine, counted)
	}

	beside := strings.Replace(sectionOf(t, document, pageHeading), skill.PageBar, "", 1)
	for _, judged := range []string{
		"which option you would take",
		"looks equally likely",
		"where you would take one",
	} {
		if strings.Contains(beside, judged) {
			t.Errorf("the page section says %q beside the bar, and the judgement about when to recommend is the bar's:\n%s", judged, beside)
		}
	}
}

func TestTheDecisionsAreInFullAndTheRestIsAList(t *testing.T) {
	page := strings.ToLower(sectionOf(t, skill.Document(), pageHeading))

	for _, said := range []string{"in full", "short list"} {
		if !strings.Contains(page, said) {
			t.Errorf("the page section does not say %q, so the only thing the person has to act on is buried among the rest", said)
		}
	}
	for _, said := range []string{
		"rather than prose you write",
		"every entry of it, whether or not a decision cites that entry",
		"retype none of them",
	} {
		if !strings.Contains(page, said) {
			t.Errorf("the page section does not say %q, so the list is one the session writes and one that can disagree with its own terminal:\n%s", said, page)
		}
	}
}

func sectionOf(t *testing.T, document, heading string) string {
	t.Helper()

	_, after, found := strings.Cut(document, heading)
	if !found {
		t.Fatalf("the skill has no %s section", heading)
	}

	depth := strings.Count(strings.Fields(heading)[0], "#")
	for line := range strings.SplitSeq(after, "\n") {
		if !strings.HasPrefix(line, "#") {
			continue
		}
		if strings.Count(strings.Fields(line)[0], "#") > depth {
			continue
		}
		if next := strings.Index(after, "\n"+line); next >= 0 {
			return after[:next]
		}
	}
	return after
}

func stepOf(t *testing.T, document, named string) string {
	t.Helper()

	at := strings.Index(document, named)
	if at < 0 {
		t.Fatalf("the skill has no step naming %q", named)
	}

	start := strings.LastIndex(document[:at], "\n3.")
	if start < 0 {
		t.Fatalf("the step naming %q is not numbered", named)
	}
	end := strings.Index(document[start+1:], "\n4.")
	if end < 0 {
		t.Fatalf("the step naming %q does not end", named)
	}
	return document[start : start+1+end]
}

func TestTheSkillSaysToCarryTheSessionFromServeToPoll(t *testing.T) {
	page := sectionOf(t, skill.Document(), pageHeading)

	for _, said := range []string{
		"agora serve --document ./round.json" + inlineCode + " prints the loopback address of its page",
		"agora poll --session",
		"the session at the end of the address",
	} {
		if !strings.Contains(page, said) {
			t.Errorf("the page section does not say %q, so a session polls a surface holding more than one page and is told about none of them:\n%s", said, page)
		}
	}
}

func TestALaterRoundReachesThePageThePersonAlreadyHas(t *testing.T) {
	document := skill.Document()
	page := sectionOf(t, document, pageHeading)

	for _, said := range []string{
		"Keep that session for the whole review",
		"the page they already have",
		"agora serve --session",
		"you hand over no second address",
	} {
		if !strings.Contains(page, said) {
			t.Errorf("the page section does not say %q, so a second round a person reads is served as a page of its own and the first goes stale:\n%s", said, page)
		}
	}

	if counted := strings.Count(document, "Keep that session for the whole review"); counted != 1 {
		t.Errorf("the skill says to keep the session %d times, want once, where the session is already stated", counted)
	}

	for _, branch := range branchesThatHaveAResult {
		section := sectionOf(t, document, branch)
		if !strings.Contains(section, "unless an earlier round of this run already gave them one") {
			t.Errorf("%s says to hand over the address with no condition on it, and a branch reached as the second round hands over a second address:\n%s", branch, section)
		}
	}
}

func TestTheSessionIsStatedOnceRatherThanOnEveryBranch(t *testing.T) {
	document := skill.Document()

	if counted := strings.Count(document, "agora poll --session"); counted != 1 {
		t.Errorf("the skill says how the session reaches agora poll %d times, want once, and the copy a session reads first is the one that gets softened", counted)
	}
	for _, branch := range []string{escalatedHeading, parkedHeading} {
		if strings.Contains(sectionOf(t, document, branch), "--session") {
			t.Errorf("%s restates how the session reaches agora poll, and it is stated once where the page is described", branch)
		}
	}
}
