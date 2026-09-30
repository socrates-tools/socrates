package e2e_test

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/agora/surface"
	"github.com/socrates-tools/socrates/review/findings"
	"github.com/socrates-tools/socrates/review/report"
	"github.com/socrates-tools/socrates/review/runstate"
)

type reportedFinding struct {
	ID          string `json:"id"`
	Severity    string `json:"severity"`
	File        string `json:"file"`
	Line        int    `json:"line"`
	Description string `json:"description"`
	Action      string `json:"action"`
	Resolution  string `json:"resolution"`
	Inside      *bool  `json:"inside_the_change"`
}

type reportedRound struct {
	Run           string            `json:"run"`
	State         string            `json:"state"`
	AwaitsAPerson bool              `json:"awaits_a_person"`
	Findings      []reportedFinding `json:"findings"`
}

func (f reportedFinding) audience() string {
	return report.Audience(findings.Action(f.Action))
}

func (f reportedFinding) name() string {
	return "Finding " + strings.TrimPrefix(f.ID, "f")
}

func (f reportedFinding) location() string {
	return fmt.Sprintf("%s:%d", f.File, f.Line)
}

func (j *journey) reported(name string) reportedRound {
	j.t.Helper()

	printed, failures, _ := j.run("findings", name, "--json")
	if printed == "" {
		j.t.Fatalf("socrates review findings --json printed nothing for %s: %s", name, failures)
	}

	var read reportedRound
	if err := json.Unmarshal([]byte(printed), &read); err != nil {
		j.t.Fatalf("reading what socrates review findings --json reported: %v\n%s", err, printed)
	}
	return read
}

func (j *journey) serving() *surfaced {
	j.t.Helper()

	if j.surface == nil {
		j.surface = servingSurface(j.t)
		j.surface.where = j.session.Repository.Root
	}
	return j.surface
}

func (j *journey) serve(written string) string {
	j.t.Helper()

	held := j.serving()
	j.noted("agora", "serve", []string{"--document"})

	return held.serve(written)
}

func (j *journey) serveAgain(written string) string {
	j.t.Helper()

	held := j.serving()
	j.noted("agora", "serve", []string{"--document", "--session"})

	return held.serveAgain(written)
}

func (j *journey) polled() surface.Answered {
	j.t.Helper()

	held := j.serving()
	j.noted("agora", "poll", []string{"--session"})

	printed, failures, code := held.poll()
	if code != 0 {
		j.t.Fatalf("agora poll exited %d: %s", code, failures)
	}

	var read surface.Answered
	if err := json.Unmarshal([]byte(printed), &read); err != nil {
		j.t.Fatalf("reading what agora poll printed: %v\n%s", err, printed)
	}
	return read
}

func sourceOf(read reportedRound) []map[string]any {
	entries := make([]map[string]any, 0, len(read.Findings))
	for _, finding := range read.Findings {
		entries = append(entries, map[string]any{
			"id":       finding.ID,
			"name":     finding.name(),
			"severity": finding.Severity,
			"audience": finding.audience(),
			"location": finding.location(),
		})
	}
	return entries
}

func documentFor(t *testing.T, title string, read reportedRound, decisions []map[string]any) string {
	t.Helper()

	written, err := json.Marshal(map[string]any{
		"title":     title,
		"preamble":  "what an independent reviewer said about the change",
		"source":    sourceOf(read),
		"decisions": decisions,
	})
	if err != nil {
		t.Fatalf("writing the document the page is served from: %v", err)
	}
	return string(written)
}

func everyFindingAsADecision(read reportedRound) []map[string]any {
	decisions := make([]map[string]any, 0, len(read.Findings))
	for _, finding := range read.Findings {
		decisions = append(decisions, map[string]any{
			"id":       finding.ID,
			"name":     "whatever the session felt like",
			"severity": "whatever the session felt like",
			"audience": "whoever the session felt like",
			"location": "wherever the session felt like",
			"question": finding.Description,
			"meaning":  "this is what it costs the change",
			"mode":     "single-choice",
			"options": []map[string]any{
				{"label": "look again once I have changed this", "meaning": "a fresh reviewer reads the change as it then stands", "command": `socrates review begin --intent "` + finding.ID + `"`},
				{"label": "read what it said again", "meaning": "the findings as socrates printed them", "command": "socrates review findings " + read.Run},
			},
		})
	}
	return decisions
}

func TestARoundThatWaitsOnAPersonIsWaitedOn(t *testing.T) {
	driven := driving(t, dialogueScenario(unanswerableFinding, nothingFoundAgain), answeredWrongly)
	name := driven.begin("count the characters of a document")

	state, _ := driven.waitedFor(name)
	if state != runstate.Escalated {
		t.Fatalf("a round escalating for a finding no answer addresses left the run %s", state)
	}

	if reported := driven.findings(name); !strings.Contains(reported, awaitsAPerson) {
		t.Fatalf("socrates does not report the round as one that waits on a person:\n%s", reported)
	}

	read := driven.reported(name)
	if !read.AwaitsAPerson {
		t.Fatalf("socrates review findings --json does not report the round as one that waits on a person: %+v", read)
	}
	if len(read.Findings) == 0 {
		t.Fatal("socrates review findings --json reported no finding, so there is nothing to hand over")
	}

	address := driven.serve(documentFor(t, "three things this change was told", read, everyFindingAsADecision(read)))
	if err := surface.CheckLoopback(originOf(t, address)[len("http://"):]); err != nil {
		t.Fatalf("agora serve handed back %q: %v", address, err)
	}

	page := opening(t, address)
	for _, finding := range read.Findings {
		var shown bool
		page.evaluate(fmt.Sprintf(`!!document.querySelector('.decision[data-decision="%s"]')`, finding.ID), &shown)
		if !shown {
			t.Errorf("the page the person was handed does not hold %s, which socrates reported for them to answer", finding.ID)
		}
	}

	t.Run("Every finding keeps the identifier socrates gave it", func(t *testing.T) {
		for _, finding := range read.Findings {
			var name, severity, audience, location string
			page.evaluate(shownIn(finding.ID, ".name"), &name)
			page.evaluate(markedIn(finding.ID, "severity"), &severity)
			page.evaluate(markedIn(finding.ID, "audience"), &audience)
			page.evaluate(shownIn(finding.ID, ".location"), &location)

			if strings.TrimSpace(name) != finding.name() {
				t.Errorf("%s is shown under the name %q, and the source the decision cites named it %q", finding.ID, name, finding.name())
			}
			if severity != finding.Severity {
				t.Errorf("%s is shown with the severity %q, and socrates review findings --json reported %q", finding.ID, severity, finding.Severity)
			}
			if audience != finding.audience() {
				t.Errorf("%s is shown to %q, and socrates review findings --json reported %q", finding.ID, audience, finding.audience())
			}
			if location != finding.location() {
				t.Errorf("%s is shown at %q, and socrates review findings --json reported %q", finding.ID, location, finding.location())
			}
		}
	})

	nothingLeftTheMachine(t, address, page.requests.all())
	driven.wroteNothingItself()
	driven.namedNothingTheSkillDoesNot()
}

func shownIn(decision, selector string) string {
	return fmt.Sprintf(`(document.querySelector('.decision[data-decision="%s"] %s') || {textContent: ""}).textContent`, decision, selector)
}

func markedIn(decision, kind string) string {
	return fmt.Sprintf(`(document.querySelector('.decision[data-decision="%s"] .mark.%s') || {dataset: {}}).dataset.%s || ""`, decision, kind, kind)
}

func TestADecisionMadeOnThePageIsRecordedInTheTerminal(t *testing.T) {
	const inTheirOwnWords = "the function stays, because the importer reads it every morning"

	driven := driving(t, dialogueScenario(questionOfIntent, nothingFoundAgain), answeredWrongly)
	name := driven.begin("count the characters of a document")

	state, _ := driven.waitedFor(name)
	if state != runstate.AwaitingResponse {
		t.Fatalf("a round parked on a question about intent left the run %s", state)
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
			{"label": "the function stays", "meaning": "the importer reads it", "command": `socrates review respond --finding ` + asked.ID + ` --decided --reason "the function stays"`},
			{"label": "take it out", "meaning": "nothing reads it", "command": `socrates review respond --finding ` + asked.ID + ` --decided --reason "take the function out"`},
		},
	}}))

	page := opening(t, address)
	page.evaluate(chose(asked.ID, 0), nil)
	page.write(within(asked.ID, ".words-field"), inTheirOwnWords)
	page.click(".finish")
	page.waitFor(`document.querySelector('.send').disabled === true`)

	answered := driven.polled()
	if len(answered.Decisions) != 1 {
		t.Fatalf("what the person sent carried %d decisions, want the one they committed: %+v", len(answered.Decisions), answered)
	}
	decided := answered.Decisions[0]
	if decided.ID != asked.ID {
		t.Errorf("the answer came back against %q, and the person answered %q", decided.ID, asked.ID)
	}
	if decided.Words != inTheirOwnWords {
		t.Errorf("what the person wrote came back as %q", decided.Words)
	}

	if _, err := os.Stat(driven.session.Round(name, 1).ResponsesPath()); !os.IsNotExist(err) {
		t.Errorf("the run holds a response before socrates review respond was run: %v", err)
	}
	if waiting, _, _ := driven.run("status", name); !strings.Contains(waiting, "needing a person rather than a patch: 1") {
		t.Errorf("the run does not still wait on the person after they answered on the page:\n%s", waiting)
	}

	if _, failures, code := driven.run("respond", name, "--finding", asked.ID, "--decided", "--reason", decided.Words); code != 4 {
		t.Fatalf("recording what the person decided exited %d: %s", code, failures)
	}

	recorded := readResponses(t, driven.session, name)
	held, found := recorded.For(asked.ID)
	if !found {
		t.Fatalf("the run holds no response for %s once socrates review respond has run", asked.ID)
	}

	t.Run("The person's own words reach the record", func(t *testing.T) {
		if held.Reason != inTheirOwnWords {
			t.Errorf("the reason the next reviewer receives reads %q, and the person wrote %q", held.Reason, inTheirOwnWords)
		}
	})

	nothingLeftTheMachine(t, address, page.requests.all())
	driven.wroteNothingItself()
	driven.namedNothingTheSkillDoesNot()
}
