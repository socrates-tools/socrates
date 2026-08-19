package report

import (
	"encoding/json"
	"fmt"

	"github.com/rui-ktei/socrates/internal/exit"
	"github.com/rui-ktei/socrates/internal/findings"
	"github.com/rui-ktei/socrates/internal/runstate"
	"github.com/rui-ktei/socrates/internal/runstore"
)

type Counts struct {
	Total          int `json:"total"`
	NeedingAPerson int `json:"needing_a_person"`
	Patchable      int `json:"patchable"`
	Informational  int `json:"informational"`
}

type Status struct {
	Run      string         `json:"run"`
	State    runstate.State `json:"state"`
	Round    int            `json:"round"`
	Reason   string         `json:"reason,omitempty"`
	Findings *Counts        `json:"findings,omitempty"`
}

type Finding struct {
	ID          string            `json:"id"`
	Severity    findings.Severity `json:"severity"`
	File        string            `json:"file"`
	Line        int               `json:"line"`
	Description string            `json:"description"`
	Resolved    findings.Action   `json:"action"`
}

type Findings struct {
	Run           string         `json:"run"`
	State         runstate.State `json:"state"`
	Reviewed      bool           `json:"reviewed"`
	Findings      []Finding      `json:"findings"`
	RiskLevel     findings.Risk  `json:"risk_level,omitempty"`
	RiskRationale string         `json:"risk_rationale,omitempty"`
	Reason        string         `json:"reason,omitempty"`
}

func StatusOf(store runstore.Store, name string) (Status, error) {
	run, document, reviewed, round, err := read(store, name)
	if err != nil {
		return Status{}, err
	}

	status := Status{Run: run.Name, State: run.State, Round: round, Reason: reasonOf(run)}
	if reviewed {
		classified := findings.Classify(document)
		status.Findings = &Counts{
			Total:          classified.Total(),
			NeedingAPerson: len(classified.NeedingAPerson),
			Patchable:      len(classified.Patchable),
			Informational:  len(classified.NoOps),
		}
	}
	return status, nil
}

func FindingsOf(store runstore.Store, name string) (Findings, error) {
	run, document, reviewed, _, err := read(store, name)
	if err != nil {
		return Findings{}, err
	}

	reported := Findings{Run: run.Name, State: run.State, Reviewed: reviewed, Findings: []Finding{}, Reason: reasonOf(run)}
	if !reviewed {
		return reported, nil
	}

	for _, finding := range document.Findings {
		reported.Findings = append(reported.Findings, Finding{
			ID:          finding.ID,
			Severity:    finding.Severity,
			File:        finding.File,
			Line:        finding.Line,
			Description: finding.Description,
			Resolved:    finding.EffectiveAction(),
		})
	}
	reported.RiskLevel = document.RiskLevel
	reported.RiskRationale = document.RiskRationale
	return reported, nil
}

func read(store runstore.Store, name string) (*runstore.Run, findings.Document, bool, int, error) {
	run, err := store.Load(name)
	if err != nil {
		return nil, findings.Document{}, false, 0, err
	}

	round, hasRound, err := store.LatestRound(name)
	if err != nil {
		return nil, findings.Document{}, false, 0, err
	}
	if !hasRound {
		return run, findings.Document{}, false, 0, nil
	}
	if !reviewCompleted(run.State) {
		return run, findings.Document{}, false, round.Number, nil
	}

	document, present, err := round.Findings()
	if err != nil {
		return nil, findings.Document{}, false, 0, err
	}
	return run, document, present, round.Number, nil
}

func reviewCompleted(state runstate.State) bool {
	return state == runstate.Clean || state == runstate.Escalated
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
	}
	if s.Reason != "" {
		lines = append(lines, fmt.Sprintf(verdictLine, s.Reason))
	}
	if s.Findings != nil {
		lines = append(lines, fmt.Sprintf(countsLine,
			s.Findings.Total, s.Findings.NeedingAPerson, s.Findings.Patchable, s.Findings.Informational))
	}
	return lines
}

func (f Findings) Code() int {
	return exit.For(f.State)
}

func (f Findings) Lines() []string {
	if !f.Reviewed {
		return []string{fmt.Sprintf(noReviewLine, f.stateReason())}
	}
	if len(f.Findings) == 0 {
		return []string{nothingToReportLine, fmt.Sprintf(riskLine, f.RiskLevel), fmt.Sprintf(rationaleLine, f.RiskRationale)}
	}

	var lines []string
	for _, finding := range f.Findings {
		lines = append(lines,
			fmt.Sprintf(findingLine, finding.ID, finding.Severity, finding.File, finding.Line, Audience(finding.Resolved)),
			fmt.Sprintf(descriptionLine, finding.Description))
	}
	return append(lines, fmt.Sprintf(riskLine, f.RiskLevel), fmt.Sprintf(rationaleLine, f.RiskRationale))
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
