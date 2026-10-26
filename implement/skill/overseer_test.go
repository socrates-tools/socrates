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

func TestTheOverseerSkillCarriesEachHostsPersonOnlySwitch(t *testing.T) {
	for _, one := range skill.All() {
		directory := t.TempDir()
		if _, err := skillkit.Install(directory, []skillkit.Skill{one}, skillkit.Dispatcher+" implement skill"); err != nil {
			t.Fatalf("installing the skill rendered for %s: %v", one.Host.Name, err)
		}
		if !one.Host.PersonOnly(filepath.Join(directory, one.Name)) {
			t.Errorf("the skill rendered for %s does not carry its person-only switch once installed, so a model could start an overseer", one.Host.Name)
		}
	}
}

func TestTheOverseerSkillDeclaresTheNameItIsInstalledUnder(t *testing.T) {
	for _, one := range skill.All() {
		if !strings.Contains(one.Rendered(), "\nname: "+skill.Name+"\n") {
			t.Errorf("the skill rendered for %s does not declare itself %q", one.Host.Name, skill.Name)
		}
	}
}

func TestTheCommittedOverseerSkillMatchesWhatIsRendered(t *testing.T) {
	for _, one := range skill.All() {
		path := filepath.Join(one.Host.SkillsDirectoryUnder(source.Root(t)), skill.Name, skillkit.FileName)

		committed, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("%s is not readable, and it is the copy this repository ships: %v. Run make skill to write it", path, err)
		}
		if string(committed) != one.Rendered() {
			t.Errorf("%s is not what the skill's source generates: run make skill to rewrite it", path)
		}
	}
}

func TestInstallingTheOverseerSkillWritesItUnderItsOwnName(t *testing.T) {
	directory := t.TempDir()

	written, err := skillkit.Install(directory, skill.For(skillkit.ClaudeCode), skillkit.Dispatcher+" implement skill")
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

func TestTheOverseerSkillNamesInvocationsOfImplementAndAtTheLastReadOfReview(t *testing.T) {
	named := skillkit.Named(skill.All(), skill.Namespaces)
	if len(named) == 0 {
		t.Fatal("the skill names no invocation at all, so this check reads nothing")
	}

	commands := map[string]map[string]bool{}
	for _, invocation := range named {
		if invocation.Binary != skill.Name && invocation.Binary != reviewBinary {
			t.Errorf("the skill names an invocation attributed to %s", invocation.Binary)
		}
		if commands[invocation.Binary] == nil {
			commands[invocation.Binary] = map[string]bool{}
		}
		commands[invocation.Binary][invocation.Command] = true
	}
	for _, command := range []string{"config", "default-branch", "skill"} {
		if !commands[skill.Name][command] {
			t.Errorf("the skill never names socrates implement %s", command)
		}
	}
	for _, command := range []string{"begin", "findings", "respond", "continue"} {
		if !commands[reviewBinary][command] {
			t.Errorf("the skill never names socrates review %s", command)
		}
	}
}

func TestTheOverseerSkillAttributesReviewsCommandsToReviewsOwnBinary(t *testing.T) {
	if got := skill.Namespaces["review"]; got != reviewBinary {
		t.Errorf("implement attributes socrates review commands to %q, and a name other than %q overwrites review's own entry when the suite merges both", got, reviewBinary)
	}
}

func sectionOf(t *testing.T, host skillkit.Host, heading string) string {
	t.Helper()

	_, after, found := strings.Cut(skill.Document(host), heading)
	if !found {
		t.Fatalf("the skill for %s has no %s section", host.Name, heading)
	}
	if next := strings.Index(after, "\n## "); next >= 0 {
		return after[:next]
	}
	return after
}

func saysEach(t *testing.T, heading, why string, said ...string) {
	t.Helper()

	for _, host := range skillkit.Hosts() {
		saysEachOn(t, host, heading, why, said...)
	}
}

func saysEachOn(t *testing.T, host skillkit.Host, heading, why string, said ...string) {
	t.Helper()

	section := sectionOf(t, host, heading)
	for _, one := range said {
		if !strings.Contains(section, one) {
			t.Errorf("the %s section for %s does not say %q, so %s:\n%s", heading, host.Name, one, why, section)
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

	reviewBinary = "socrates-review"
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
		"reading the main ticket, the spec, every ADR the spec names and the glossaries, before anything is handed out",
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
		"Cut it from the current `HEAD`",
		"the change's folder, the ADRs whose Part of names its slug, and the glossaries",
		"Make those the first commit of the branch",
		"Anything else uncommitted is the person's",
		"start nothing",
	)
}

func TestTheChangesBranchIsNamedForTheKindOfChangeTheSpecAsksFor(t *testing.T) {
	saysEach(t, branchHeading, "the change could land on a branch whose name says nothing about what kind of change it is",
		"The change's branch is `<type>/<slug>`",
		"`<type>` is `feature`, `fix`, `refactor`, `docs` or `chore`, whichever the spec says the change is",
		"or for a change that is several, what it gives somebody using the project",
	)
}

func TestAResumedOverseerKeepsTheTypeTheChangesBranchAlreadyHas(t *testing.T) {
	saysEach(t, branchHeading, "a resumed run could judge the type afresh and start a second branch beside the first",
		"Where `<type>/<slug>` exists already under any type, this run is a resumed one and keeps that branch",
	)
}

func TestTheFrontierIsEveryTicketWhoseBlockersAreDoneBoundedByConcurrency(t *testing.T) {
	saysEach(t, frontierHeading, "a ticket could start before its blockers are done or more could run than the person allows",
		"The frontier is every ticket whose blockers are done",
		"Run at most `concurrency` implementers at once",
	)
}

func TestTheLongestChainGoesFirstThenTheMostUnblockedThenNameOrder(t *testing.T) {
	section := sectionOf(t, skillkit.ClaudeCode, frontierHeading)

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
		"git worktree add -b <type>/<slug>-<ticket>",
		"the `worktrees/` directory beside the file `socrates implement config` prints on its first line",
		"`~/.socrates/implement/worktrees/` by default",
		"a sibling of the change's branch rather than a child",
	)
}

func TestAnImplementerIsBriefedWithPointersAndRunsInTheBackgroundOnClaudeCode(t *testing.T) {
	saysEachOn(t, skillkit.ClaudeCode, worktreesHeading, "the overseer could block on a ticket",
		"Run the implementer in the background",
		"so you stay free to talk to the person",
	)
}

func TestAnImplementerIsBriefedWithPointers(t *testing.T) {
	saysEach(t, worktreesHeading, "the overseer could hand over a paraphrase that drifts",
		"Brief it with pointers rather than a paraphrase",
		"its worktree, its branch, the change's branch, its ticket, the notes, and the spec and ADRs the ticket rests on",
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
		"git merge --ff-only <type>/<slug>-<ticket>",
		"run the repository's own verification",
		"Run it after every landing",
		"A failure goes back to the implementer whose ticket landed last",
	)
}

func TestTheVerificationAfterALandingRunsInTheBackgroundSoTheOverseerStaysFree(t *testing.T) {
	saysEachOn(t, skillkit.ClaudeCode, landingHeading, "a verification that takes minutes could leave the person unanswered",
		"Run it in the background, so you stay free to talk to the person while it runs, as the implementers do",
	)
}

func TestTheCommandsTheOverseerRunsAreNotCountedAsAllThatTheNamespaceHolds(t *testing.T) {
	for _, one := range skill.All() {
		if strings.Contains(one.Rendered(), "offers three commands") {
			t.Errorf("the skill for %s counts the binary's commands, which stops being true when the namespace gains one", one.Host.Name)
		}
	}
	saysEach(t, "# socrates implement", "the commands the overseer runs are not stated", "The commands you run are these")
}

func TestWhereOnlyTheOverseersOwnStatusCommitsAreMissingSheMergesThemHerself(t *testing.T) {
	saysEach(t, landingHeading, "every status commit would send a merged implementer back for a merge that cannot conflict",
		"Every status commit you make on the change's branch puts an implementer that already merged one commit behind",
		"Implementers never touch the main ticket",
		"the only commits its branch lacks are your own status commits",
		"git -C <worktree> merge <type>/<slug>",
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
		"git merge-base <default branch> <type>/<slug>",
		"to one fresh implementer",
		"remove every worktree",
		"the branch `<type>/<slug>` is ready",
		"Push nothing",
	)
}

func TestTheSectionsComeInTheOrderOfTheLoop(t *testing.T) {
	for _, host := range skillkit.Hosts() {
		document := skill.Document(host)

		previous := -1
		for _, heading := range []string{
			writesNoCodeHeading, startingHeading, defaultHeading, branchHeading, frontierHeading, worktreesHeading,
			reviewHeading, landingHeading, statusHeading, resumingHeading, personHeading, lastReadHeading,
		} {
			at := strings.Index(document, heading)
			if at < 0 || at < previous {
				t.Fatalf("%q is missing or comes before the section that must precede it in the skill for %s", heading, host.Name)
			}
			previous = at
		}
	}
}

func TestTheOverseerSkillNamesReviewsCommandsAtTheLastReadAndNowhereElse(t *testing.T) {
	for _, host := range skillkit.Hosts() {
		document := skill.Document(host)
		lastRead := sectionOf(t, host, lastReadHeading)

		if !regexp.MustCompile(`(?i)socrates\s+review`).MatchString(lastRead) {
			t.Errorf("the last read for %s names no socrates review command", host.Name)
		}
		outside := strings.Replace(document, lastRead, "", 1)
		if found := regexp.MustCompile(`(?i)socrates\s+review`).FindString(outside); found != "" {
			t.Errorf("the skill for %s names %q outside the last read, and how a ticket is reviewed is the overseer's own", host.Name, found)
		}
	}
}

func TestTheOverseerSkillNamesNoReviewSkillToInvoke(t *testing.T) {
	for _, one := range skill.All() {
		if found := regexp.MustCompile(`socrates-review`).FindString(one.Rendered()); found != "" {
			t.Errorf("the skill for %s names %q, and the overseer drives review's commands rather than its skill", one.Host.Name, found)
		}
	}
}

func TestTheOverseerSkillNamesNoNamespaceButImplementAndReview(t *testing.T) {
	for _, tool := range convention.Tools {
		if tool.Namespace == "implement" || tool.Namespace == "review" {
			continue
		}
		for _, one := range skill.All() {
			if found := regexp.MustCompile(`(?i)socrates\s+` + regexp.QuoteMeta(tool.Namespace) + `\b`).FindString(one.Rendered()); found != "" {
				t.Errorf("the skill for %s names %q, and it drives no other namespace", one.Host.Name, found)
			}
		}
	}
}

func TestTheOverseerSkillNamesNoOtherSkill(t *testing.T) {
	allowed := map[string]bool{skill.Name: true, skill.AgentName: true}

	for _, one := range skill.All() {
		for _, named := range regexp.MustCompile(`socrates-[a-z]+(-[a-z]+)*`).FindAllString(one.Rendered(), -1) {
			if !allowed[named] {
				t.Errorf("the skill for %s names %q, and an overseer invokes no skill", one.Host.Name, named)
			}
		}
	}
}

func TestTheOverseerIsToldToInvokeNoSkill(t *testing.T) {
	saysEach(t, writesNoCodeHeading, "an overseer could start another workflow underneath the person's",
		"Invoke no skill",
	)
}

func TestOneSubagentWritesNotesOutsideTheRepositoryBeforeTheFirstHandOut(t *testing.T) {
	saysEach(t, startingHeading, "every implementer could explore the same code from nothing and name one thing twice",
		"Before the first hand-out, have one subagent read the code the change touches and write notes to `<notes>/<slug>.md`",
		"the names the shared surfaces will use",
		"`~/.socrates/implement/notes/` by default",
		"outside the repository",
		"You write no notes yourself",
		"a resumed run reads the notes already there",
	)
}

func TestTheOverseerReviewsOnTwoCountsQuotingTheLineOrRuleMissed(t *testing.T) {
	saysEach(t, reviewHeading, "feedback could say something falls short without saying of what",
		"two counts kept apart",
		"quoting the line it falls short of",
		"citing the rule it breaks",
		"An acceptance criterion holds only where something would fail without the work",
	)
}

func TestATicketThatNeedsAPersonIsAskedAboutBeforeItIsHandedOut(t *testing.T) {
	saysEach(t, personHeading, "an implementer could stall halfway through a ticket waiting on a credential",
		"Where a ticket names what it needs a person for, ask for it in the terminal before handing the ticket out",
	)
}

func TestTheLastReadGoesThroughTheSpecsUserStoriesOneByOne(t *testing.T) {
	saysEach(t, lastReadHeading, "a story no ticket delivered could reach a branch called ready",
		"Go through the spec's user stories one by one",
		"every story the change does not deliver",
	)
}

func TestTheOverseerNamesBehaviourTheTicketDidNotAskForAboveAllASiblingTicketsWork(t *testing.T) {
	saysEach(t, reviewHeading, "an implementer could build ahead into a sibling ticket under way in another worktree",
		"Name too what the diff does that the ticket does not ask for",
		"above all what another ticket delivers",
	)
}

func TestTheHandBackMapsEachStoryToItsTicketAndNamesEveryOneWayDoor(t *testing.T) {
	saysEach(t, lastReadHeading, "the person could be asked to publish a change without knowing what it delivers or what cannot be walked back",
		"each user story with the ticket that delivered it",
		"every one-way door the change walks through",
		"say so where there is none",
	)
}

func TestTheHandBackNamesFeedbackThatRecurredAsACheckWorthWriting(t *testing.T) {
	saysEach(t, lastReadHeading, "a mistake fed back to implementer after implementer could be forgotten when the run ends",
		"feedback you sent for the same mistake to more than one implementer",
		"a check the repository could run instead",
	)
}

func TestTheLastReadBeginsSocratesReviewOverTheWholeChangeWithNobodyToAsk(t *testing.T) {
	for _, host := range skillkit.Hosts() {
		saysEachOn(t, host, lastReadHeading, "the whole change could reach the person unreviewed, or the review could stop to ask somebody who is not there",
			"socrates review begin --wait --base <default branch> --nobody-to-ask --intent \"This change builds docs/changes/<slug>/spec.md\" --host "+host.Flag,
			"on the change's branch in the person's own checkout",
		)
	}
}

func TestTheReviewRunsInTheBackgroundOnEachHostSoTheOverseerStaysFree(t *testing.T) {
	saysEachOn(t, skillkit.ClaudeCode, lastReadHeading, "a review that takes minutes could leave the person unanswered",
		"Run it in the background",
	)
	saysEachOn(t, skillkit.Codex, lastReadHeading, "a review that takes minutes could leave the person unanswered",
		"Have a subagent run it",
	)
}

func TestTheStoriesAreWalkedWhileTheReviewerReadsAndFindingsAreReadAsJSON(t *testing.T) {
	saysEach(t, lastReadHeading, "the story read could wait on the review, or the findings could be read from prose",
		"While the reviewer reads",
		"socrates review findings --json",
	)
}

func TestAFindingThatAsksAPersonIsSettledByTheOverseersOwnJudgementAndAFalseOneIsRejected(t *testing.T) {
	saysEach(t, lastReadHeading, "the overseer could stop to ask a person who is not there, or record a decision no person made",
		"Settle each finding that asks a person by your own judgement",
		"reject each you judge false",
		"`--settled`, never `--decided`",
	)
}

func TestEveryReviewFixAndEveryStoryShortfallGoesToOneFreshImplementerLandedAndVerifiedBeforeTheSecondRound(t *testing.T) {
	saysEach(t, lastReadHeading, "fixes could be built twice, or reviewed again before they landed",
		"every fix, settled answers included",
		"to one fresh implementer together with the stories the change falls short of",
		"in a worktree of its own made the same way as a ticket's",
		"land and verify its work the same way",
		"before the second round",
	)
}

func TestAResponseIsRecordedToEveryFindingBeforeTheSecondRoundRuns(t *testing.T) {
	saysEach(t, lastReadHeading, "the second round could be asked to judge answers that were never given",
		"Record a response to every finding with `socrates review respond --finding <id>`",
		"`--fixed`",
		"`--rejected`",
		"`socrates review continue --wait`",
	)
}

func TestACleanFirstRoundHasNoSecondRoundAndItsStoryShortfallsStillGoToTheImplementer(t *testing.T) {
	saysEach(t, lastReadHeading, "a clean review could be run again for nothing, or the stories it never read could be dropped",
		"A first round that ended clean has no second round",
		"its story shortfalls still go to the one implementer",
	)
}

func TestAReviewThatCannotBeginOrEndsEscalatedOrFailedStartsNoFurtherRunAndIsHandedBackAnyway(t *testing.T) {
	saysEach(t, lastReadHeading, "the overseer could loop on a review that cannot pass, try to repair a reviewer CLI that is not installed, or stop before the person is told",
		"Where the review cannot begin or ends escalated or failed, start no further run and hand back anyway, with the refusal or the failure it gave",
	)
}

func TestTheHandBackSaysHowTheReviewEndedWhatIsStillOpenAndEachSettledFindingWithItsReason(t *testing.T) {
	saysEach(t, lastReadHeading, "the person could publish a change without knowing what the review left open or what the overseer decided for them",
		"say how the review ended",
		"every finding still open or the failure the reviewer gave",
		"each settled finding with its reason",
	)
}
