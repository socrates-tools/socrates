package report

import (
	"encoding/json"
	"fmt"

	"github.com/socrates-tools/socrates/review/dialogue"
	"github.com/socrates-tools/socrates/review/exit"
	"github.com/socrates-tools/socrates/review/findings"
	"github.com/socrates-tools/socrates/review/runstate"
	"github.com/socrates-tools/socrates/review/runstore"
)

type Counts struct {
	Total          int `json:"total"`
	NeedingAPerson int `json:"needing_a_person"`
	Patchable      int `json:"patchable"`
	Informational  int `json:"informational"`
}

type Status struct {
	Run             string         `json:"run"`
	State           runstate.State `json:"state"`
	Round           int            `json:"round"`
	Base            string         `json:"base"`
	Reviewer        string         `json:"reviewer"`
	Tier            string         `json:"tier"`
	Reason          string         `json:"reason,omitempty"`
	Findings        *Counts        `json:"findings,omitempty"`
	Unanswered      *int           `json:"unanswered,omitempty"`
	AwaitingAPerson *int           `json:"awaiting_a_person,omitempty"`
}

type Finding struct {
	ID          string              `json:"id"`
	Severity    findings.Severity   `json:"severity"`
	File        string              `json:"file"`
	Line        int                 `json:"line"`
	Description string              `json:"description"`
	Resolved    findings.Action     `json:"action"`
	Resolution  dialogue.Resolution `json:"resolution,omitempty"`
	Inside      *bool               `json:"inside_the_change,omitempty"`
}

type Findings struct {
	Run           string         `json:"run"`
	State         runstate.State `json:"state"`
	Reviewed      bool           `json:"reviewed"`
	Findings      []Finding      `json:"findings"`
	RiskLevel     findings.Risk  `json:"risk_level,omitempty"`
	RiskRationale string         `json:"risk_rationale,omitempty"`
	AwaitsAPerson bool           `json:"awaits_a_person,omitempty"`
	Reason        string         `json:"reason,omitempty"`
	Guidance      []string       `json:"guidance,omitempty"`
}

func StatusOf(store runstore.Store, name string) (Status, error) {
	read, err := readRun(store, name)
	if err != nil {
		return Status{}, err
	}

	run := read.run
	status := Status{Run: run.Name, State: run.State, Round: read.number, Base: run.Base, Reviewer: run.Vendor(), Tier: run.Tier(), Reason: reasonOf(run)}
	if read.reviewed {
		classified := dialogue.Classified(read.document, run.Authorship())
		status.Findings = &Counts{
			Total:          classified.Total(),
			NeedingAPerson: len(classified.NeedingAPerson),
			Patchable:      len(classified.Patchable),
			Informational:  len(classified.NoOps),
		}
	}
	if run.State == runstate.AwaitingResponse {
		waiting := len(read.responses.Unanswered(read.document, run.Authorship()))
		asked := len(read.responses.AwaitingAPerson(read.document, run.Authorship()))
		status.Unanswered = &waiting
		status.AwaitingAPerson = &asked
	}
	return status, nil
}

func FindingsOf(store runstore.Store, name string) (Findings, error) {
	read, err := readRun(store, name)
	if err != nil {
		return Findings{}, err
	}

	run, document := read.run, read.document
	reported := Findings{Run: run.Name, State: run.State, Reviewed: read.reviewed, Findings: []Finding{}, Reason: reasonOf(run)}
	if run.State == runstate.AwaitingResponse {
		reported.Guidance = RulesFor(run.Authorship())
	}
	if !read.reviewed {
		return reported, nil
	}

	change, located, err := read.round.Change()
	if err != nil {
		return Findings{}, err
	}

	reported.AwaitsAPerson = dialogue.AwaitsAPerson(dialogue.Classified(document, run.Authorship()), run.State)
	for _, finding := range document.Findings {
		reporting := Finding{
			ID:          finding.ID,
			Severity:    finding.Severity,
			File:        finding.File,
			Line:        finding.Line,
			Description: finding.Description,
			Resolved:    finding.EffectiveAction(),
		}
		if answer, answered := read.responses.For(finding.ID); answered {
			reporting.Resolution = answer.Resolution
		}
		if located {
			inside := change.Holds(finding.File, finding.Line)
			reporting.Inside = &inside
		}
		reported.Findings = append(reported.Findings, reporting)
	}
	reported.RiskLevel = document.RiskLevel
	reported.RiskRationale = document.RiskRationale
	return reported, nil
}

type reading struct {
	run       *runstore.Run
	round     runstore.Round
	document  findings.Document
	responses dialogue.Responses
	reviewed  bool
	number    int
}

func readRun(store runstore.Store, name string) (reading, error) {
	run, err := store.Load(name)
	if err != nil {
		return reading{}, err
	}

	round, hasRound, err := store.LatestRound(name)
	if err != nil {
		return reading{}, err
	}
	if !hasRound {
		return reading{run: run}, nil
	}
	if !reviewCompleted(run.State) {
		return reading{run: run, round: round, number: round.Number}, nil
	}

	document, present, err := round.Findings()
	if err != nil {
		return reading{}, err
	}
	responses, err := round.Responses()
	if err != nil {
		return reading{}, err
	}
	return reading{run: run, round: round, document: document, responses: responses, reviewed: present, number: round.Number}, nil
}

func reviewCompleted(state runstate.State) bool {
	return state == runstate.Clean || state == runstate.Escalated || state == runstate.AwaitingResponse
}

func reasonOf(run *runstore.Run) string {
	if run.Verdict == nil {
		return ""
	}
	return run.Verdict.Reason
}

func (s Status) Code() int {
	return exit.For(s.State)
}

func (s Status) Lines() []string {
	lines := []string{
		fmt.Sprintf(runLine, s.Run),
		fmt.Sprintf(stateLine, s.State),
		fmt.Sprintf(roundLine, s.Round),
		fmt.Sprintf(baseLine, s.Base),
		fmt.Sprintf(reviewerLine, s.Reviewer, s.Tier),
	}
	if s.Reason != "" {
		lines = append(lines, fmt.Sprintf(verdictLine, s.Reason))
	}
	if s.Findings != nil {
		lines = append(lines, fmt.Sprintf(countsLine,
			s.Findings.Total, s.Findings.NeedingAPerson, s.Findings.Patchable, s.Findings.Informational))
	}
	if s.Unanswered == nil {
		return lines
	}

	lines = append(lines, awaitingLine, fmt.Sprintf(unansweredLine, *s.Unanswered))
	if s.AwaitingAPerson != nil {
		lines = append(lines, fmt.Sprintf(awaitingAPersonLine, *s.AwaitingAPerson))
	}
	if *s.Unanswered == 0 {
		lines = append(lines, readyLine)
	}
	return lines
}

func (f Findings) Code() int {
	return exit.For(f.State)
}

func (f Findings) Lines() []string {
	if !f.Reviewed {
		return append([]string{fmt.Sprintf(noReviewLine, f.stateReason())}, f.rules()...)
	}
	if len(f.Findings) == 0 {
		return append([]string{nothingToReportLine, fmt.Sprintf(riskLine, f.RiskLevel), fmt.Sprintf(rationaleLine, f.RiskRationale)}, f.rules()...)
	}

	var lines []string
	for _, finding := range f.Findings {
		lines = append(lines,
			fmt.Sprintf(findingLine, finding.ID, finding.Severity, finding.File, finding.Line, Audience(finding.Resolved)),
			fmt.Sprintf(descriptionLine, finding.Description))
	}
	lines = append(lines, fmt.Sprintf(riskLine, f.RiskLevel), fmt.Sprintf(rationaleLine, f.RiskRationale))
	lines = append(lines, f.forAPerson()...)
	return append(lines, f.rules()...)
}

func (f Findings) forAPerson() []string {
	if !f.AwaitsAPerson {
		return nil
	}
	return []string{awaitsAPersonLine}
}

func (f Findings) rules() []string {
	if len(f.Guidance) == 0 {
		return nil
	}

	lines := []string{"", rulesHeading}
	for _, rule := range f.Guidance {
		lines = append(lines, fmt.Sprintf(ruleLine, rule))
	}
	return lines
}

func ResponseRules() []string {
	return dialogue.Rules()
}

func HandoverRules() []string {
	return dialogue.HandoverRules()
}

func RulesFor(authorship dialogue.Authorship) []string {
	if authorship == dialogue.NotMine {
		return HandoverRules()
	}
	return ResponseRules()
}

func (f Findings) stateReason() string {
	if f.Reason != "" {
		return f.Reason
	}
	return fmt.Sprintf(stillWorkingReason, f.State)
}

func Audience(action findings.Action) string {
	switch action {
	case findings.AutoFix:
		return patchableAudience
	case findings.NoOp:
		return informationalAudience
	default:
		return personAudience
	}
}

func AsJSON(reported any) (string, error) {
	encoded, err := json.Marshal(reported)
	if err != nil {
		return "", fmt.Errorf(encodingFailed, err)
	}
	return string(encoded) + "\n", nil
}
