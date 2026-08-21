package e2e_test

import (
	"maps"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/rui-ktei/socrates/internal/runstate"
	"github.com/rui-ktei/socrates/internal/skill"
	"github.com/rui-ktei/socrates/internal/testing/harness"
)

const loopPatience = 30 * time.Second

type journey struct {
	t        *testing.T
	session  *harness.Session
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

	j.ran = append(j.ran, skill.Invocation{Command: command})
	for _, argument := range arguments {
		if strings.HasPrefix(argument, "--") {
			j.ran = append(j.ran, skill.Invocation{Command: command, Flag: argument})
		}
	}
	return j.session.Scrat(append([]string{command}, arguments...)...)
}

func (j *journey) review(intent string) string {
	j.t.Helper()

	printed, failures, code := j.run("review", "--intent", intent)
	if code != 0 {
		j.t.Fatalf("scrat review --intent exited %d: %s\n%s", code, failures, printed)
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
			j.t.Fatalf("scrat status printed no state for %s: %s%s", name, printed, failures)
		}
		if slices.Contains(stops, state) {
			return state, code
		}
		time.Sleep(20 * time.Millisecond)
	}

	j.t.Fatalf("polling scrat status for %s reported no state to act on within %s", name, loopPatience)
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
		j.t.Fatalf("scrat findings printed nothing for %s: %s", name, failures)
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
			j.t.Errorf("the journey ran scrat %s %s, which the skill does not name, so it is not the loop the skill describes", ran.Command, ran.Flag)
		}
	}
}

func TestTheLoopStopsWhereTheSkillSaysToStop(t *testing.T) {
	t.Run(string(runstate.Escalated), func(t *testing.T) {
		driven := driving(t, dialogueScenario(questionOfIntent, nothingFoundAgain), answeredWrongly)
		name := driven.review("count the characters of a document")

		state, code := driven.waitedFor(name)
		if state != runstate.Escalated || code != 1 {
			t.Fatalf("the run is %s and scrat status exited %d, want escalated", state, code)
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
			t.Fatalf("the run is %s and scrat status exited %d, want failed", state, code)
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
			t.Fatalf("the run is %s and scrat status exited %d, want a parked run", state, code)
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

func TestTheLoopTheSkillDescribesReachesAVerdict(t *testing.T) {
	driven := driving(t, dialogueScenario(patchableFindings, nothingFoundAgain), answeredWrongly)
	name := driven.review("count the characters of a document")

	state, code := driven.waitedFor(name)
	if state != runstate.AwaitingResponse || code != 4 {
		t.Fatalf("the first round left the run %s and scrat status exited %d, want a parked run", state, code)
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
		t.Fatalf("scrat continue exited %d: %s", code, failures)
	}

	state, code = driven.waitedFor(name)
	if state != runstate.Clean || code != 0 {
		t.Fatalf("the second round left the run %s and scrat status exited %d, want clean", state, code)
	}

	driven.wroteNothingItself()
	driven.namedNothingTheSkillDoesNot()
}
