package skill_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/implement/skill"
	"github.com/socrates-tools/socrates/internal/convention"
	skillkit "github.com/socrates-tools/socrates/internal/skill"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

func TestTheOverseerSkillIsHeldToACeilingOfItsOwn(t *testing.T) {
	for _, violation := range skillkit.CeilingViolations(skill.All()) {
		t.Error(violation)
	}
}

func TestTheOverseerSkillCarriesThePersonOnlyField(t *testing.T) {
	if !convention.CarriesThePersonOnlyField(skill.Document()) {
		t.Errorf("the skill does not carry %s, so a model could start an overseer", skillkit.PersonOnlyField)
	}
}

func TestTheOverseerSkillDeclaresTheNameItIsInstalledUnder(t *testing.T) {
	if !strings.Contains(skill.Document(), "\nname: "+skill.Name+"\n") {
		t.Errorf("the skill does not declare itself %q", skill.Name)
	}
}

func TestTheCommittedOverseerSkillMatchesWhatIsRendered(t *testing.T) {
	path := filepath.Join(source.Root(t), ".claude", "skills", skill.Name, skillkit.FileName)

	committed, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("%s is not readable, and it is the copy this repository ships: %v. Run make skill to write it", path, err)
	}
	if string(committed) != skill.Document() {
		t.Errorf("%s is not what the skill's source generates: run make skill to rewrite it", path)
	}
}

func TestInstallingTheOverseerSkillWritesItUnderItsOwnName(t *testing.T) {
	directory := t.TempDir()

	written, err := skillkit.Install(directory, skill.All(), skillkit.Dispatcher+" implement skill")
	if err != nil {
		t.Fatalf("installing: %v", err)
	}

	if len(written) != 1 || written[0].Name != skill.Name {
		t.Fatalf("installing wrote %v, and implement ships the one skill %s", written, skill.Name)
	}
	if want := filepath.Join(directory, skill.Name, skillkit.FileName); written[0].Path != want {
		t.Errorf("installing wrote %s, want %s", written[0].Path, want)
	}
}

func TestTheOverseerSkillNamesInvocationsOfImplementOnly(t *testing.T) {
	named := skillkit.Named(skill.All(), skill.Namespaces)
	if len(named) == 0 {
		t.Fatal("the skill names no invocation at all, so this check reads nothing")
	}

	commands := map[string]bool{}
	for _, invocation := range named {
		if invocation.Binary != skill.Name {
			t.Errorf("the skill names an invocation attributed to %s", invocation.Binary)
		}
		commands[invocation.Command] = true
	}
	for _, command := range []string{"config", "default-branch", "skill"} {
		if !commands[command] {
			t.Errorf("the skill never names socrates implement %s", command)
		}
	}
}

func sectionOf(t *testing.T, heading string) string {
	t.Helper()

	_, after, found := strings.Cut(skill.Document(), heading)
	if !found {
		t.Fatalf("the skill has no %s section", heading)
	}
	if next := strings.Index(after, "\n## "); next >= 0 {
		return after[:next]
	}
	return after
}

func saysEach(t *testing.T, heading, why string, said ...string) {
	t.Helper()

	section := sectionOf(t, heading)
	for _, one := range said {
		if !strings.Contains(section, one) {
			t.Errorf("the %s section does not say %q, so %s:\n%s", heading, one, why, section)
		}
	}
}

const (
	writesNoCodeHeading = "## You write no code"
	startingHeading     = "## Starting a run"
	defaultHeading      = "## The default branch"
	branchHeading       = "## The change's branch"
	frontierHeading     = "## The frontier"
	worktreesHeading    = "## Worktrees and briefs"
	reviewHeading       = "## Review"
	landingHeading      = "## Landing"
	statusHeading       = "## Status"
	resumingHeading     = "## Resuming"
	personHeading       = "## The person"
	lastReadHeading     = "## The last read"
)

func TestTheOverseerWritesNoCode(t *testing.T) {
	saysEach(t, writesNoCodeHeading, "an overseer could patch what an implementer got wrong",
		"You write no code, no test and no documentation",
		"The one file you edit is the main ticket, and only its status table",
		"goes back to that implementer as feedback",
		"A conflict is never yours to resolve",
	)
}

func TestTheOverseerReadsTheMainTicketTheSpecAndTheADRsBeforeHandingOutAnything(t *testing.T) {
	saysEach(t, startingHeading, "tickets would be reviewed against their own text alone",
		"reading the main ticket, the spec, every ADR the spec names and `CONTEXT.md`, before anything is handed out",
	)
}

func TestTheOverseerReadsTheSettingsFromConfigAtTheStartOfEveryRun(t *testing.T) {
	saysEach(t, startingHeading, "a changed cap would hold only after something was reinstalled",
		"Read the settings with `socrates implement config` at the start of every run",
		"`concurrency` is the most implementers running at once",
	)
}

func TestWhereNoDefaultBranchIsRememberedTheOverseerProposesTheNamedCandidates(t *testing.T) {
	saysEach(t, defaultHeading, "a default branch could be guessed or taken from a feature prefix",
		"reports no `default_branch`",
		"you do not guess one",
		"named `main`, `master`, `dev` or `develop`",
		"local or remote-tracking",
		"Never propose a branch under a prefix such as `feature/` or `hotfix/`",
		"Ask in the terminal",
		"`socrates implement default-branch <branch>`",
	)
}

func TestTheChangesBranchIsCutFromTheCurrentHeadAndStartsWithWhatTheDesignLeft(t *testing.T) {
	saysEach(t, branchHeading, "the branch could start from the wrong commit or sweep up the person's own work",
		"`socrates-implement/<slug>`",
		"cut from the current `HEAD`",
		"the change's folder, the ADRs whose Part of names its slug, and `CONTEXT.md`",
		"Make those the first commit of the branch",
		"Anything else uncommitted is the person's",
		"start nothing",
	)
}

func TestTheFrontierIsEveryTicketWhoseBlockersAreDoneBoundedByConcurrency(t *testing.T) {
	saysEach(t, frontierHeading, "a ticket could start before its blockers are done or more could run than the person allows",
		"The frontier is every ticket whose blockers are done",
		"Run at most `concurrency` implementers at once",
	)
}

func TestTheLongestChainGoesFirstThenTheMostUnblockedThenNameOrder(t *testing.T) {
	section := sectionOf(t, frontierHeading)

	chain := strings.Index(section, "the ticket with the longest chain of tickets waiting behind it")
	unblocks := strings.Index(section, "the ticket that unblocks the most tickets")
	name := strings.Index(section, "name order")
	if chain < 0 || unblocks < 0 || name < 0 || chain > unblocks || unblocks > name {
		t.Errorf("the frontier section does not order longest chain, then most unblocked, then name order, so two runs over one graph could hand out in different orders:\n%s", section)
	}
}

func TestAReadyTicketSharingASurfaceIsHeldBackAndNoBlockerIsWritten(t *testing.T) {
	saysEach(t, frontierHeading, "two implementers could add one thing under two names, or the tickets as cut could change",
		"Hold back a ready ticket that shares a surface with a running one",
		"a message catalogue, a type or a registry",
		"Write no blocker into the tickets",
	)
}

func TestEachImplementerWorksInAWorktreeTheOverseerMadeOnABranchSheNamed(t *testing.T) {
	saysEach(t, worktreesHeading, "an implementer could work anywhere, on a branch the overseer cannot find",
		"Make each implementer's worktree and branch yourself, with git, off the change's branch",
		"git worktree add -b socrates-implement/<slug>-<ticket>",
		"~/.socrates/implement/worktrees/",
		"a sibling of the change's branch rather than a child",
	)
}

func TestAnImplementerIsBriefedWithPointersAndRunsInTheBackground(t *testing.T) {
	saysEach(t, worktreesHeading, "the overseer could block on a ticket or hand over a paraphrase that drifts",
		"Run the implementer in the background",
		"so you stay free to talk to the person",
		"Brief it with pointers rather than a paraphrase",
		"its worktree, its branch, the change's branch, its ticket, and the spec and ADRs the ticket rests on",
	)
}

func TestTheOverseerReviewsEachDiffAndSendsFeedbackToTheSameImplementerWithNoCap(t *testing.T) {
	saysEach(t, reviewHeading, "a diff could land unread or a fix could come from a context that does not know the code",
		"review its diff against the ticket, the spec and the ADRs",
		"send feedback to the same implementer",
		"Set no cap on the exchanges",
		"hand it to a fresh implementer briefed with what the first one got wrong",
	)
}

func TestLandingIsAFastForwardFollowedByTheRepositorysOwnVerificationAfterEveryLanding(t *testing.T) {
	saysEach(t, landingHeading, "the overseer could resolve a conflict herself or leave two tickets that fail together unnoticed",
		"git merge --ff-only socrates-implement/<slug>-<ticket>",
		"run the repository's own verification",
		"Run it after every landing",
		"A failure goes back to the implementer whose ticket landed last",
	)
}

func TestWhereOnlyTheOverseersOwnStatusCommitsAreMissingSheMergesThemHerself(t *testing.T) {
	saysEach(t, landingHeading, "every status commit would send a merged implementer back for a merge that cannot conflict",
		"Every status commit you make on the change's branch puts an implementer that already merged one commit behind",
		"Implementers never touch the main ticket",
		"the only commits its branch lacks are your own status commits",
		"git -C <worktree> merge socrates-implement/<slug>",
		"That merge cannot conflict",
	)
}

func TestAnyOtherNewCommitOnTheChangesBranchGoesBackToTheImplementerToMerge(t *testing.T) {
	saysEach(t, landingHeading, "the overseer could merge somebody's code and own the conflict",
		"Any other new commit on the change's branch goes back to the implementer to merge",
	)
}

func TestInProgressIsCommittedBeforeAHandOutAndDoneOnceATicketHasLandedAndPassed(t *testing.T) {
	saysEach(t, statusHeading, "a session resuming after a lost context could hand a ticket out twice",
		"You are the one writer of the status table",
		"Commit in progress before a ticket is handed out",
		"Commit done once a ticket has landed and passed the verification",
		"abandoned, with a reason",
	)
}

func TestAResumedOverseerContinuesEachInProgressTicketFromItsOwnBranchWithAFreshImplementer(t *testing.T) {
	saysEach(t, resumingHeading, "committed work could be thrown away",
		"Hand each in-progress ticket to a fresh implementer",
		"continue from that ticket's own branch",
	)
}

func TestThePersonIsReachedInTheTerminalAndOnlyAboutWhatCodingCannotFix(t *testing.T) {
	saysEach(t, personHeading, "the person could be asked about what a capable engineer would decide",
		"Reach the person in the terminal",
		"only about what coding cannot fix",
		"a security block, a credential, an account or a permission you cannot grant yourself",
		"every implementer the question does not concern keeps working",
	)
}

func TestWhatThePersonSaysMidRunIsRelayedAtOnceToEveryImplementerItConcerns(t *testing.T) {
	saysEach(t, personHeading, "an implementer could keep building to a brief the person has since corrected",
		"relay it at once to every implementer it concerns",
		"the one channel through which the person talks about the change",
	)
}

func TestTheWholeChangeIsReadAgainstTheSpecFromTheDefaultBranchsForkPointAndNothingIsPushed(t *testing.T) {
	saysEach(t, lastReadHeading, "tickets that each passed could add up to something the spec did not ask for, or the change could be published unasked",
		"read the whole change against the spec",
		"git merge-base <default branch> socrates-implement/<slug>",
		"Hand what falls short to one implementer",
		"remove every worktree",
		"the branch `socrates-implement/<slug>` is ready",
		"Push nothing",
	)
}

func TestTheSectionsComeInTheOrderOfTheLoop(t *testing.T) {
	document := skill.Document()

	previous := -1
	for _, heading := range []string{
		writesNoCodeHeading, startingHeading, defaultHeading, branchHeading, frontierHeading, worktreesHeading,
		reviewHeading, landingHeading, statusHeading, resumingHeading, personHeading, lastReadHeading,
	} {
		at := strings.Index(document, heading)
		if at < 0 || at < previous {
			t.Fatalf("%q is missing or comes before the section that must precede it", heading)
		}
		previous = at
	}
}

func TestTheOverseerSkillNamesNoSocratesReviewInvocation(t *testing.T) {
	if found := regexp.MustCompile(`(?i)socrates\s+review`).FindString(skill.Document()); found != "" {
		t.Errorf("the skill names %q, and how a ticket is reviewed is the overseer's own", found)
	}
}

func TestTheOverseerSkillNamesNoNamespaceButImplement(t *testing.T) {
	for _, tool := range convention.Tools {
		if tool.Namespace == "implement" {
			continue
		}
		if found := regexp.MustCompile(`(?i)socrates\s+` + regexp.QuoteMeta(tool.Namespace) + `\b`).FindString(skill.Document()); found != "" {
			t.Errorf("the skill names %q, and it drives no other namespace", found)
		}
	}
}

func TestTheOverseerSkillNamesNoOtherSkill(t *testing.T) {
	allowed := map[string]bool{skill.Name: true, skill.AgentName: true}

	for _, named := range regexp.MustCompile(`socrates-[a-z]+(-[a-z]+)*`).FindAllString(skill.Document(), -1) {
		if !allowed[named] {
			t.Errorf("the skill names %q, and an overseer invokes no skill", named)
		}
	}
}

func TestTheOverseerIsToldToInvokeNoSkill(t *testing.T) {
	saysEach(t, writesNoCodeHeading, "an overseer could start another workflow underneath the person's",
		"Invoke no skill",
	)
}
