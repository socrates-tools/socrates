package dialogue

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/socrates-devtool/socrates/internal/findings"
	"github.com/socrates-devtool/socrates/internal/refusal"
)

type Resolution string

const (
	Fixed    Resolution = "fixed"
	Rejected Resolution = "rejected"
	Decided  Resolution = "decided"
)

var resolutions = []Resolution{Fixed, Rejected, Decided}

func Resolutions() []Resolution {
	return slices.Clone(resolutions)
}

func Rules() []string {
	return strings.Split(ResponseRules, "\n")
}

type Response struct {
	Finding    string     `json:"finding"`
	Resolution Resolution `json:"resolution"`
	Reason     string     `json:"reason,omitempty"`
}

type Responses struct {
	Recorded []Response `json:"responses"`
}

func Validate(response Response) error {
	if strings.TrimSpace(response.Finding) == "" {
		return refusal.New(aResponseNamesItsFinding)
	}
	if !slices.Contains(resolutions, response.Resolution) {
		return refusal.New(notAResolution, response.Resolution, resolutions)
	}
	if strings.TrimSpace(response.Reason) != "" {
		return nil
	}
	if response.Resolution == Rejected {
		return refusal.New(aRejectionIsAClaim)
	}
	if response.Resolution == Decided {
		return refusal.New(aDecisionIsAClaim)
	}
	return nil
}

func Accepts(finding findings.Finding) []Resolution {
	if finding.EffectiveAction() == findings.AskUser {
		return []Resolution{Decided}
	}
	return []Resolution{Fixed, Rejected}
}

func Answering(finding findings.Finding, response Response) error {
	if slices.Contains(Accepts(finding), response.Resolution) {
		return nil
	}
	if finding.EffectiveAction() == findings.AskUser {
		return refusal.New(theFindingAsksAPerson, finding.ID, response.Resolution, Decided)
	}
	return refusal.New(theFindingAsksForAPatch, finding.ID, Fixed, Rejected)
}

func (r Responses) For(finding string) (Response, bool) {
	for _, recorded := range slices.Backward(r.Recorded) {
		if recorded.Finding == finding {
			return recorded, true
		}
	}
	return Response{}, false
}

func (r Responses) Latest() []Response {
	var latest []Response
	for _, recorded := range r.Recorded {
		if slices.ContainsFunc(latest, func(seen Response) bool { return seen.Finding == recorded.Finding }) {
			continue
		}
		answer, _ := r.For(recorded.Finding)
		latest = append(latest, answer)
	}
	return latest
}

func (r Responses) Answers(finding string) int {
	answers := 0
	for _, recorded := range r.Recorded {
		if recorded.Finding == finding {
			answers++
		}
	}
	return answers
}

func (r Responses) Unanswered(document findings.Document) []string {
	answerable := Answerable(findings.Classify(document))

	var waiting []string
	for _, reported := range document.Findings {
		waitedOn := slices.ContainsFunc(answerable, func(one findings.Finding) bool { return one.ID == reported.ID })
		if _, answered := r.For(reported.ID); waitedOn && !answered {
			waiting = append(waiting, reported.ID)
		}
	}
	return waiting
}

func (r Responses) AwaitingAPerson(document findings.Document) []string {
	var waiting []string
	for _, unanswered := range r.Unanswered(document) {
		reported, _ := Reported(document, unanswered)
		if reported.EffectiveAction() == findings.AskUser {
			waiting = append(waiting, unanswered)
		}
	}
	return waiting
}

func Reported(document findings.Document, finding string) (findings.Finding, bool) {
	for _, reported := range document.Findings {
		if reported.ID == finding {
			return reported, true
		}
	}
	return findings.Finding{}, false
}

func Parse(document []byte) (Responses, error) {
	var parsed Responses

	for line := range strings.SplitSeq(string(document), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}

		var response Response
		decoder := json.NewDecoder(strings.NewReader(line))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&response); err != nil {
			return Responses{}, fmt.Errorf(readingTheDocument, err)
		}
		if err := Validate(response); err != nil {
			return Responses{}, err
		}
		parsed.Recorded = append(parsed.Recorded, response)
	}
	return parsed, nil
}

func Line(response Response) ([]byte, error) {
	encoded, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf(writingTheDocument, err)
	}
	return append(encoded, '\n'), nil
}
