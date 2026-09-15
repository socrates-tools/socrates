package dialogue

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/socrates-tools/socrates/socrates-review/findings"
)

type Resolution string

const (
	Fixed     Resolution = "fixed"
	Rejected  Resolution = "rejected"
	Decided   Resolution = "decided"
	Forwarded Resolution = "forwarded"
	Dropped   Resolution = "dropped"
)

type Authorship bool

const (
	Mine    Authorship = false
	NotMine Authorship = true
)

var resolutions = []Resolution{Fixed, Rejected, Decided, Forwarded, Dropped}

var leaveAnArtifact = []Resolution{Fixed, Forwarded}

func Resolutions() []Resolution {
	return slices.Clone(resolutions)
}

func Flags(resolutions []Resolution) string {
	named := make([]string, 0, len(resolutions))
	for _, resolution := range resolutions {
		named = append(named, "--"+string(resolution))
	}
	if len(named) < 2 {
		return strings.Join(named, "")
	}
	return strings.Join(named[:len(named)-1], ", ") + " or " + named[len(named)-1]
}

func Rules() []string {
	return strings.Split(ResponseRules, "\n")
}

func HandoverRules() []string {
	return strings.Split(HandoverResponseRules, "\n")
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
	if slices.Contains(leaveAnArtifact, response.Resolution) {
		return nil
	}
	switch response.Resolution {
	case Rejected:
		return refusal.New(aRejectionIsAClaim)
	case Dropped:
		return refusal.New(aDropIsAClaim)
	}
	return refusal.New(aDecisionIsAClaim)
}

func Accepts(finding findings.Finding, authorship Authorship) []Resolution {
	if authorship == NotMine {
		return []Resolution{Forwarded, Dropped, Rejected}
	}
	if finding.EffectiveAction() == findings.AskUser {
		return []Resolution{Decided}
	}
	return []Resolution{Fixed, Rejected}
}

func Answering(finding findings.Finding, response Response, authorship Authorship) error {
	if slices.Contains(Accepts(finding, authorship), response.Resolution) {
		return nil
	}
	if authorship == NotMine {
		return refusal.New(theChangeIsNotMineToPatch, response.Resolution, finding.ID, Flags(Accepts(finding, NotMine)))
	}
	if response.Resolution == Forwarded || response.Resolution == Dropped {
		return refusal.New(theChangeIsMine, response.Resolution, finding.ID, Flags(Accepts(finding, Mine)))
	}
	if finding.EffectiveAction() == findings.AskUser {
		return refusal.New(theFindingAsksAPerson, finding.ID, response.Resolution, Decided)
	}
	return refusal.New(theFindingAsksForAPatch, finding.ID, response.Resolution, Fixed, Rejected)
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

func (r Responses) Unanswered(document findings.Document, authorship Authorship) []string {
	answerable := Answerable(Classified(document, authorship))

	var waiting []string
	for _, reported := range document.Findings {
		waitedOn := slices.ContainsFunc(answerable, func(one findings.Finding) bool { return one.ID == reported.ID })
		if _, answered := r.For(reported.ID); waitedOn && !answered {
			waiting = append(waiting, reported.ID)
		}
	}
	return waiting
}

func (r Responses) AwaitingAPerson(document findings.Document, authorship Authorship) []string {
	needing := AnsweredByAPerson(Classified(document, authorship))

	var waiting []string
	for _, unanswered := range r.Unanswered(document, authorship) {
		if slices.ContainsFunc(needing, func(one findings.Finding) bool { return one.ID == unanswered }) {
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
