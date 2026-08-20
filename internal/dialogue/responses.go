package dialogue

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/rui-ktei/socrates/internal/findings"
	"github.com/rui-ktei/socrates/internal/refusal"
)

type Resolution string

const (
	Fixed    Resolution = "fixed"
	Rejected Resolution = "rejected"
)

var resolutions = []Resolution{Fixed, Rejected}

func Resolutions() []Resolution {
	return slices.Clone(resolutions)
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
	if response.Resolution == Rejected && strings.TrimSpace(response.Reason) == "" {
		return refusal.New(aRejectionIsAClaim)
	}
	return nil
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
	var waiting []string
	for _, finding := range Answerable(findings.Classify(document)) {
		if _, answered := r.For(finding.ID); !answered {
			waiting = append(waiting, finding.ID)
		}
	}
	return waiting
}

func Names(document findings.Document, finding string) bool {
	for _, reported := range document.Findings {
		if reported.ID == finding {
			return true
		}
	}
	return false
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

func (r Responses) Encode() ([]byte, error) {
	var encoded []byte
	for _, response := range r.Recorded {
		line, err := Line(response)
		if err != nil {
			return nil, err
		}
		encoded = append(encoded, line...)
	}
	return encoded, nil
}
