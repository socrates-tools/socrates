package e2e_test

import (
	"maps"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/ktei/socrates/internal/runstate"
	"github.com/ktei/socrates/internal/skill"
	"github.com/ktei/socrates/internal/surface"
	"github.com/ktei/socrates/internal/testing/harness"
)

const loopPatience = 30 * time.Second

type journey struct {
	t        *testing.T
	session  *harness.Session
	surface  *surfaced
	ran      []skill.Invocation
	status   string
	contents map[string]string
}

func driving(t *testing.T, scenario string, written string) *journey {
	t.Helper()

	session := harness.Start(t)
	session.Simulator.Scenario(scenario)

	driven := &journey{t: t, session: session}
	driven.author("internal/added.go", written)
	return driven
}

func (j *journey) author(path, content string) {
	j.t.Helper()

	j.session.Repository.Write(path, content)
	j.status = j.session.Repository.Status()
	j.contents = workingTree(j.t, j.session.Repository.Root)
}

func (j *journey) run(command string, arguments ...string) (string, string, int) {
	j.t.Helper()

	j.noted(skill.Name, command, arguments)
	return j.session.Socrates(append([]string{command}, arguments...)...)
}

func (j *journey) noted(binary, command string, arguments []string) {
	j.ran = append(j.ran, skill.Invocation{Binary: binary, Command: command})
	for _, argument := range arguments {
		if strings.HasPrefix(argument, "--") {
			j.ran = append(j.ran, skill.Invocation{Binary: binary, Command: command, Flag: argument})
		}
	}
}

func (j *journey) review(intent string) string {
	j.t.Helper()

	printed, failures, code := j.run("review", "--intent", intent)
	if code != 0 {
		j.t.Fatalf("socrates review --intent exited %d: %s\n%s", code, failures, printed)
	}
	return firstLine(printed)
}

func (j *journey) waitedFor(name string) (runstate.State, int) {
	j.t.Helper()

	stops := []runstate.State{runstate.Clean, runstate.Failed, runstate.Escalated, runstate.AwaitingResponse}
	deadline := time.Now().Add(loopPatience)

	for time.Now().Before(deadline) {
		printed, failures, code := j.run("status", name)
		state := stateIn(printed)
		if state == "" {
			j.t.Fatalf("socrates status printed no state for %s: %s%s", name, printed, failures)
		}
		if slices.Contains(stops, state) {
			return state, code
		}
		time.Sleep(20 * time.Millisecond)
	}

	j.t.Fatalf("polling socrates status for %s reported no state to act on within %s", name, loopPatience)
	return "", 0
}

func stateIn(printed string) runstate.State {
	for line := range strings.SplitSeq(printed, "\n") {
		if reported, found := strings.CutPrefix(strings.TrimSpace(line), "state: "); found {
			return runstate.State(reported)
		}
	}
	return ""
}

func (j *journey) findings(name string) string {
	j.t.Helper()

	printed, failures, _ := j.run("findings", name)
	if printed == "" {
		j.t.Fatalf("socrates findings printed nothing for %s: %s", name, failures)
	}
	return printed
}

func (j *journey) wroteNothingItself() {
	j.t.Helper()

	if after := j.session.Repository.Status(); after != j.status {
		j.t.Errorf("git status changed while the loop ran:\nbefore:\n%s\nafter:\n%s", j.status, after)
	}
	if after := workingTree(j.t, j.session.Repository.Root); !maps.Equal(after, j.contents) {
		j.t.Errorf("the working tree changed while the loop ran:\nbefore: %v\nafter: %v", j.contents, after)
	}
}

func (j *journey) namedNothingTheSkillDoesNot() {
	j.t.Helper()

	named := skill.Named()
	for _, ran := range j.ran {
		if !slices.Contains(named, ran) {
			j.t.Errorf("the journey ran %s %s %s, which the skill does not name, so it is not the loop the skill describes", ran.Binary, ran.Command, ran.Flag)
		}
	}
}

func TestTheLoopStopsWhereTheSkillSaysToStop(t *testing.T) {
	t.Run(string(runstate.Escalated), func(t *testing.T) {
		driven := driving(t, dialogueScenario(unanswerableFinding, nothingFoundAgain), answeredWrongly)
		name := driven.review("count the characters of a document")

		state, code := driven.waitedFor(name)
		if state != runstate.Escalated || code != 1 {
			t.Fatalf("the run is %s and socrates status exited %d, want escalated", state, code)
		}

		reported := driven.findings(name)
		for _, handed := range []string{"f1", "f2", "f3"} {
			if !strings.Contains(reported, handed) {
				t.Errorf("the person is not handed %s:\n%s", handed, reported)
			}
		}
		if _, failures, refused := driven.run("respond", name, "--finding", "f2", "--fixed"); refused != 3 {
			t.Errorf("answering an escalated run exited %d, and the skill says to hand these to the person and stop: %s", refused, failures)
		}

		driven.wroteNothingItself()
		driven.namedNothingTheSkillDoesNot()
	})

	t.Run(string(runstate.Failed), func(t *testing.T) {
		driven := driving(t, "default:\n  behaviour: no-result\n", answeredWrongly)
		name := driven.review("count the characters of a document")

		state, code := driven.waitedFor(name)
		if state != runstate.Failed || code != 2 {
			t.Fatalf("the run is %s and socrates status exited %d, want failed", state, code)
		}
		if reported := driven.findings(name); !strings.Contains(reported, "no review completed") {
			t.Errorf("what the loop reads about a failed run does not say that no review completed:\n%s", reported)
		}

		driven.wroteNothingItself()
		driven.namedNothingTheSkillDoesNot()
	})

	t.Run(string(runstate.AwaitingResponse), func(t *testing.T) {
		driven := driving(t, dialogueScenario(patchableFindings, nothingFoundAgain), answeredWrongly)
		name := driven.review("count the characters of a document")

		state, code := driven.waitedFor(name)
		if state != runstate.AwaitingResponse || code != 4 {
			t.Fatalf("the run is %s and socrates status exited %d, want a parked run", state, code)
		}

		if _, failures, refused := driven.run("continue", name); refused != 3 {
			t.Errorf("continuing with nothing answered exited %d, and the skill says every finding that can be patched needs an answer first: %s", refused, failures)
		}
		if again, _ := driven.waitedFor(name); again != runstate.AwaitingResponse {
			t.Errorf("the run moved to %s on its own, and only an answer ends a parked round", again)
		}

		driven.wroteNothingItself()
		driven.namedNothingTheSkillDoesNot()
	})
}

func TestTheLoopAnswersAQuestionAndReachesAVerdict(t *testing.T) {
	driven := driving(t, dialogueScenario(questionOfIntent, nothingFoundAgain), answeredWrongly)
	name := driven.review("count the characters of a document")

	state, code := driven.waitedFor(name)
	if state != runstate.AwaitingResponse || code != 4 {
		t.Fatalf("a round holding a question left the run %s and socrates status exited %d, want a parked run", state, code)
	}
	driven.wroteNothingItself()

	reported := driven.findings(name)
	for _, needed := range []string{"f1", "needs a person", "can be patched"} {
		if !strings.Contains(reported, needed) {
			t.Fatalf("the findings the loop reads do not name %q:\n%s", needed, reported)
		}
	}

	driven.author("internal/added.go", answeredBetter)
	if _, failures, code := driven.run("respond", name, "--finding", "f2", "--fixed", "--note", "the guard is above the read now"); code != 4 {
		t.Fatalf("recording a fix exited %d: %s", code, failures)
	}
	if _, failures, code := driven.run("respond", name, "--finding", "f3", "--rejected", "--reason", "there is no retry loop in this function at all"); code != 4 {
		t.Fatalf("recording a rejection exited %d: %s", code, failures)
	}
	if _, failures, refused := driven.run("continue", name); refused != 3 {
		t.Fatalf("continuing with the question unanswered exited %d, and only a person closes it: %s", refused, failures)
	}

	const decided = "Vic decided the function stays, because the importer reads it"
	if _, failures, code := driven.run("respond", name, "--finding", "f1", "--decided", "--reason", decided); code != 4 {
		t.Fatalf("recording what a person decided exited %d: %s", code, failures)
	}
	if _, failures, code := driven.run("continue", name); code != 0 {
		t.Fatalf("socrates continue with everything answered exited %d: %s", code, failures)
	}

	state, code = driven.waitedFor(name)
	if state != runstate.Clean || code != 0 {
		t.Fatalf("the round that judged the decision left the run %s and socrates status exited %d, want clean", state, code)
	}

	invocations := driven.session.Simulator.Invocations()
	if len(invocations) != 2 {
		t.Fatalf("the reviewer was invoked %d times, want twice", len(invocations))
	}
	judging := invocations[1].Prompt
	if !strings.Contains(judging, "A person answered this finding, and the author records what they decided:\n"+decided) {
		t.Errorf("the round that judges the decision was not told a person made it:\n%s", judging)
	}

	driven.wroteNothingItself()
	driven.namedNothingTheSkillDoesNot()
}

func TestAParkedRoundWaitsOnAQuestionAndOnAPatch(t *testing.T) {
	driven := driving(t, dialogueScenario(questionOfIntent, nothingFoundAgain), answeredWrongly)
	name := driven.review("count the characters of a document")

	state, _ := driven.waitedFor(name)
	if state != runstate.AwaitingResponse {
		t.Fatalf("a round holding a question and two patchable findings left the run %s", state)
	}

	waiting, _, _ := driven.run("status", name)
	for _, said := range []string{"unanswered: 3", "needing a person rather than a patch: 1"} {
		if !strings.Contains(waiting, said) {
			t.Errorf("the status does not say %q, so a caller cannot tell what the round waits on:\n%s", said, waiting)
		}
	}

	if _, failures, refused := driven.run("continue", name); refused != 3 {
		t.Errorf("continuing with nothing answered exited %d: %s", refused, failures)
	}

	if _, failures, code := driven.run("respond", name, "--finding", "f1", "--decided", "--reason", "Vic decided the function stays"); code != 4 {
		t.Fatalf("recording what a person decided exited %d: %s", code, failures)
	}
	if _, failures, refused := driven.run("continue", name); refused != 3 {
		t.Errorf("continuing with the question answered and the patches unanswered exited %d: %s", refused, failures)
	}

	driven.author("internal/added.go", answeredBetter)
	if _, failures, code := driven.run("respond", name, "--finding", "f2", "--fixed", "--note", "the guard is above the read now"); code != 4 {
		t.Fatalf("recording a fix exited %d: %s", code, failures)
	}
	if _, failures, refused := driven.run("continue", name); refused != 3 {
		t.Errorf("continuing with one patchable finding still unanswered exited %d: %s", refused, failures)
	}
	if _, failures, code := driven.run("respond", name, "--finding", "f3", "--rejected", "--reason", "there is no retry loop in this function at all"); code != 4 {
		t.Fatalf("recording a rejection exited %d: %s", code, failures)
	}

	answered, _, _ := driven.run("status", name)
	if !strings.Contains(answered, "unanswered: 0") {
		t.Errorf("the status does not report zero unanswered once every finding has an answer:\n%s", answered)
	}
	if _, failures, code := driven.run("continue", name); code != 0 {
		t.Fatalf("socrates continue with everything answered exited %d: %s", code, failures)
	}

	driven.wroteNothingItself()
	driven.namedNothingTheSkillDoesNot()
}

func TestAQuestionThatReachesTheCapEscalatesWithEverythingKept(t *testing.T) {
	driven := driving(t, dialogueScenario(patchableFindings, questionAtTheCap), answeredWrongly)
	name := driven.review("count the characters of a document")

	state, _ := driven.waitedFor(name)
	if state != runstate.AwaitingResponse {
		t.Fatalf("the first round left the run %s, want a parked run", state)
	}

	driven.author("internal/added.go", answeredBetter)
	if _, failures, code := driven.run("respond", name, "--finding", "f1", "--fixed", "--note", "an empty document now answers minus one"); code != 4 {
		t.Fatalf("recording a fix exited %d: %s", code, failures)
	}
	if _, failures, code := driven.run("respond", name, "--finding", "f2", "--rejected", "--reason", "there is no retry loop in this function at all"); code != 4 {
		t.Fatalf("recording a rejection exited %d: %s", code, failures)
	}
	if _, failures, code := driven.run("continue", name); code != 0 {
		t.Fatalf("socrates continue exited %d: %s", code, failures)
	}

	state, code := driven.waitedFor(name)
	if state != runstate.Escalated || code != 1 {
		t.Fatalf("a question at the final permitted round left the run %s and socrates status exited %d, want escalated", state, code)
	}

	run := driven.session.Run(name)
	if !strings.Contains(reason(run), "cap") {
		t.Errorf("the verdict does not say the round cap was reached: %q", reason(run))
	}

	reported := driven.findings(name)
	if !strings.Contains(reported, "g1") || !strings.Contains(reported, "needs a person") {
		t.Errorf("the question is not handed to the person:\n%s", reported)
	}
	if _, failures, refused := driven.run("respond", name, "--finding", "g1", "--decided", "--reason", "Vic decided the sentinel stays"); refused != 3 {
		t.Errorf("answering a run that ran out of rounds exited %d: %s", refused, failures)
	}

	for _, path := range []string{
		driven.session.Round(name, 1).PromptPath(),
		driven.session.Round(name, 1).FindingsPath(),
		driven.session.Round(name, 1).ResponsesPath(),
		driven.session.Round(name, 2).PromptPath(),
		driven.session.Round(name, 2).FindingsPath(),
	} {
		written, err := os.Stat(path)
		if err != nil {
			t.Errorf("%s is not readable after the run escalated: %v", path, err)
			continue
		}
		if written.Size() == 0 {
			t.Errorf("%s is empty", path)
		}
	}

	driven.wroteNothingItself()
	driven.namedNothingTheSkillDoesNot()
}

func TestTheLoopTheSkillDescribesReachesAVerdict(t *testing.T) {
	driven := driving(t, dialogueScenario(patchableFindings, nothingFoundAgain), answeredWrongly)
	name := driven.review("count the characters of a document")

	state, code := driven.waitedFor(name)
	if state != runstate.AwaitingResponse || code != 4 {
		t.Fatalf("the first round left the run %s and socrates status exited %d, want a parked run", state, code)
	}
	driven.wroteNothingItself()

	reported := driven.findings(name)
	for _, needed := range []string{"f1", "f2", "can be patched"} {
		if !strings.Contains(reported, needed) {
			t.Fatalf("the findings the loop reads do not name %q:\n%s", needed, reported)
		}
	}

	driven.author("internal/added.go", answeredBetter)
	if _, failures, code := driven.run("respond", name, "--finding", "f1", "--fixed", "--note", "an empty document now answers minus one"); code != 4 {
		t.Fatalf("recording a fix exited %d: %s", code, failures)
	}
	if _, failures, code := driven.run("respond", name, "--finding", "f2", "--rejected", "--reason", "there is no retry loop in this function at all"); code != 4 {
		t.Fatalf("recording a rejection exited %d: %s", code, failures)
	}
	if _, failures, code := driven.run("continue", name); code != 0 {
		t.Fatalf("socrates continue exited %d: %s", code, failures)
	}

	state, code = driven.waitedFor(name)
	if state != runstate.Clean || code != 0 {
		t.Fatalf("the second round left the run %s and socrates status exited %d, want clean", state, code)
	}

	driven.wroteNothingItself()
	driven.namedNothingTheSkillDoesNot()
}

func TestTheLoopCarriesTheSessionFromServeToPollAndReachesAVerdict(t *testing.T) {
	const inTheirOwnWords = "the function stays, because the importer reads it every morning"

	driven := driving(t, dialogueScenario(questionOfIntent, nothingFoundAgain), answeredWrongly)
	name := driven.review("count the characters of a document")

	state, code := driven.waitedFor(name)
	if state != runstate.AwaitingResponse || code != 4 {
		t.Fatalf("a round holding a question left the run %s and socrates status exited %d, want a parked run", state, code)
	}

	read := driven.reported(name)
	var asked reportedFinding
	for _, finding := range read.Findings {
		if finding.Action == "ask-user" {
			asked = finding
		}
	}
	if asked.ID == "" {
		t.Fatalf("no finding of the parked round needs a person: %+v", read.Findings)
	}

	address := driven.serve(documentFor(t, "one thing only you can settle", read, []map[string]any{{
		"id":       asked.ID,
		"question": asked.Description,
		"meaning":  "nothing socrates can do closes this one",
		"mode":     "single-choice",
		"options": []map[string]any{
			{"label": "the function stays", "meaning": "the importer reads it", "command": `socrates respond --finding ` + asked.ID + ` --decided --reason "the function stays"`},
			{"label": "take it out", "meaning": "nothing reads it", "command": `socrates respond --finding ` + asked.ID + ` --decided --reason "take the function out"`},
		},
	}}))

	held := driven.serving()
	if carried := theSessionIn(t, address); carried != held.opened {
		t.Fatalf("the address agora serve printed holds %q, and the session it opened is %q", carried, held.opened)
	}

	held.answer(surface.Sent{
		Decisions: []surface.Choice{{ID: asked.ID, Chosen: []int{0}, Words: inTheirOwnWords}},
		Ends:      true,
	})

	answered := driven.polled()
	if len(answered.Decisions) != 1 || answered.Decisions[0].ID != asked.ID {
		t.Fatalf("agora poll on the session agora serve printed returned %+v", answered.Decisions)
	}

	driven.author("internal/added.go", answeredBetter)
	if _, failures, code := driven.run("respond", name, "--finding", asked.ID, "--decided", "--reason", answered.Decisions[0].Words); code != 4 {
		t.Fatalf("recording what the person answered on the page exited %d: %s", code, failures)
	}
	if _, failures, code := driven.run("respond", name, "--finding", "f2", "--fixed", "--note", "the guard is above the read now"); code != 4 {
		t.Fatalf("recording a fix exited %d: %s", code, failures)
	}
	if _, failures, code := driven.run("respond", name, "--finding", "f3", "--rejected", "--reason", "there is no retry loop in this function at all"); code != 4 {
		t.Fatalf("recording a rejection exited %d: %s", code, failures)
	}
	if _, failures, code := driven.run("continue", name); code != 0 {
		t.Fatalf("socrates continue with everything answered exited %d: %s", code, failures)
	}

	state, code = driven.waitedFor(name)
	if state != runstate.Clean || code != 0 {
		t.Fatalf("the round that judged what the person answered on the page left the run %s and socrates status exited %d, want clean", state, code)
	}

	judging := driven.session.Simulator.Invocations()[1].Prompt
	if !strings.Contains(judging, inTheirOwnWords) {
		t.Errorf("the next round was not told what the person answered on the page:\n%s", judging)
	}

	driven.wroteNothingItself()
	driven.namedNothingTheSkillDoesNot()
}

func TestALaterRoundIsWrittenToThePageThePersonAlreadyHas(t *testing.T) {
	driven := driving(t, dialogueScenario(questionOfIntent, questionAtTheCap), answeredWrongly)
	name := driven.review("count the characters of a document")

	state, code := driven.waitedFor(name)
	if state != runstate.AwaitingResponse || code != 4 {
		t.Fatalf("a round holding a question left the run %s and socrates status exited %d, want a parked run", state, code)
	}

	first := driven.reported(name)
	firstAddress := driven.serve(documentFor(t, "one thing only you can settle", first, oneQuestionAsADecision(t, first)))

	held := driven.serving()
	opened := held.opened
	held.answer(surface.Sent{
		Decisions: []surface.Choice{{ID: whatOnlyAPersonAnswers(t, first).ID, Chosen: []int{0}, Words: "the function stays, because the importer reads it"}},
		Ends:      true,
	})
	answered := driven.polled()

	driven.author("internal/added.go", answeredBetter)
	if _, failures, refused := driven.run("respond", name, "--finding", "f1", "--decided", "--reason", answered.Decisions[0].Words); refused != 4 {
		t.Fatalf("recording what the person answered on the page exited %d: %s", refused, failures)
	}
	if _, failures, refused := driven.run("respond", name, "--finding", "f2", "--fixed", "--note", "the guard is above the read now"); refused != 4 {
		t.Fatalf("recording a fix exited %d: %s", refused, failures)
	}
	if _, failures, refused := driven.run("respond", name, "--finding", "f3", "--rejected", "--reason", "there is no retry loop in this function at all"); refused != 4 {
		t.Fatalf("recording a rejection exited %d: %s", refused, failures)
	}
	if _, failures, refused := driven.run("continue", name); refused != 0 {
		t.Fatalf("socrates continue with everything answered exited %d: %s", refused, failures)
	}

	state, code = driven.waitedFor(name)
	if state != runstate.Escalated || code != 1 {
		t.Fatalf("the second round left the run %s and socrates status exited %d, want escalated", state, code)
	}

	second := driven.reported(name)
	secondAddress := driven.serveAgain(documentFor(t, "what the second round could not settle", second, oneQuestionAsADecision(t, second)))

	if secondAddress != firstAddress {
		t.Fatalf("the person was handed %q for the first round and %q for the second, and one review is one page", firstAddress, secondAddress)
	}
	if held.opened != opened {
		t.Fatalf("the second round was written to the session %q, and the first was served under %q", held.opened, opened)
	}

	held.answer(surface.Sent{
		Decisions: []surface.Choice{{ID: whatOnlyAPersonAnswers(t, second).ID, Chosen: []int{0}, Words: "the sentinel stays"}},
		Ends:      true,
	})
	if read := driven.polled(); len(read.Decisions) != 1 || read.Decisions[0].Words != "the sentinel stays" {
		t.Fatalf("agora poll on the one session returned %+v, and the second round was answered on the page the first was served under", read.Decisions)
	}

	driven.wroteNothingItself()
	driven.namedNothingTheSkillDoesNot()
}

func whatOnlyAPersonAnswers(t *testing.T, read reportedRound) reportedFinding {
	t.Helper()

	for _, finding := range read.Findings {
		if finding.Action == "ask-user" {
			return finding
		}
	}
	t.Fatalf("no finding of the round needs a person: %+v", read.Findings)
	return reportedFinding{}
}

func oneQuestionAsADecision(t *testing.T, read reportedRound) []map[string]any {
	t.Helper()

	asked := whatOnlyAPersonAnswers(t, read)
	return []map[string]any{{
		"id":       asked.ID,
		"question": asked.Description,
		"meaning":  "nothing socrates can do closes this one",
		"mode":     "single-choice",
		"options": []map[string]any{
			{"label": "it stays", "meaning": "what is there is what the product wants", "command": `socrates respond --finding ` + asked.ID + ` --decided --reason "it stays"`},
			{"label": "take it out", "meaning": "nothing reads it", "command": `socrates respond --finding ` + asked.ID + ` --decided --reason "take it out"`},
		},
	}}
}
