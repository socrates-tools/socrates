package findings

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"slices"
	"strings"

	"github.com/socrates-tools/socrates-review/internal/runstate"
)

type Severity string

const (
	Error   Severity = "error"
	Warning Severity = "warning"
	Info    Severity = "info"
)

type Action string

const (
	NoOp    Action = "no-op"
	AutoFix Action = "auto-fix"
	AskUser Action = "ask-user"
)

type Risk string

const (
	LowRisk    Risk = "low"
	MediumRisk Risk = "medium"
	HighRisk   Risk = "high"
)

var severities = []Severity{Error, Warning, Info}

var actions = []Action{NoOp, AutoFix, AskUser}

var risks = []Risk{LowRisk, MediumRisk, HighRisk}

func Severities() []Severity {
	return slices.Clone(severities)
}

func Actions() []Action {
	return slices.Clone(actions)
}

func Risks() []Risk {
	return slices.Clone(risks)
}

type Finding struct {
	ID          string   `json:"id"`
	Severity    Severity `json:"severity"`
	File        string   `json:"file"`
	Line        int      `json:"line"`
	Description string   `json:"description"`
	Action      Action   `json:"action"`
}

type Document struct {
	Findings      []Finding `json:"findings"`
	RiskLevel     Risk      `json:"risk_level"`
	RiskRationale string    `json:"risk_rationale"`
}

func (s Severity) Blocks() bool {
	return s != Info
}

func (f Finding) EffectiveAction() Action {
	if slices.Contains(actions, f.Action) {
		return f.Action
	}
	return AskUser
}

type Classification struct {
	NeedingAPerson []Finding
	Patchable      []Finding
	NoOps          []Finding
}

func Classify(document Document) Classification {
	var classified Classification

	for _, finding := range document.Findings {
		switch finding.EffectiveAction() {
		case AutoFix:
			classified.Patchable = append(classified.Patchable, finding)
		case NoOp:
			classified.NoOps = append(classified.NoOps, finding)
		default:
			classified.NeedingAPerson = append(classified.NeedingAPerson, finding)
		}
	}

	return classified
}

func (c Classification) Total() int {
	return len(c.NeedingAPerson) + len(c.Patchable) + len(c.NoOps)
}

func (d Document) NeedsAPerson() bool {
	return len(Classify(d).NeedingAPerson) > 0
}

func (d Document) Patchable() []Finding {
	return Classify(d).Patchable
}

func (c Classification) Verdict() runstate.State {
	for _, group := range [][]Finding{c.NeedingAPerson, c.Patchable, c.NoOps} {
		for _, finding := range group {
			if finding.Severity.Blocks() {
				return runstate.Escalated
			}
		}
	}
	return runstate.Clean
}

func (d Document) Verdict() runstate.State {
	return Classify(d).Verdict()
}

func Parse(document []byte) (Document, error) {
	var stated map[string]json.RawMessage
	if err := json.Unmarshal(document, &stated); err != nil {
		return Document{}, fmt.Errorf(readingTheDocument, err)
	}
	for _, required := range []string{"findings", "risk_level", "risk_rationale"} {
		if _, present := stated[required]; !present {
			return Document{}, fmt.Errorf(statesNothingFor, required)
		}
	}

	var parsed Document
	decoder := json.NewDecoder(bytes.NewReader(document))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&parsed); err != nil {
		return Document{}, fmt.Errorf(readingTheDocument, err)
	}
	if decoder.More() {
		return Document{}, errors.New(moreThanOneDocument)
	}

	if err := Validate(parsed); err != nil {
		return Document{}, err
	}
	return parsed, nil
}

func Validate(document Document) error {
	if !slices.Contains(risks, document.RiskLevel) {
		return fmt.Errorf(unknownRiskLevel, document.RiskLevel, risks)
	}

	seen := map[string]bool{}
	for position, finding := range document.Findings {
		named := fmt.Sprintf(findingAt, position+1)
		if finding.ID != "" {
			named = fmt.Sprintf(findingNamed, finding.ID)
		}

		if finding.ID == "" {
			return fmt.Errorf(noIdentifier, named)
		}
		if seen[finding.ID] {
			return fmt.Errorf(repeatedIdentifier, finding.ID)
		}
		seen[finding.ID] = true

		if !slices.Contains(severities, finding.Severity) {
			return fmt.Errorf(unknownSeverity, named, finding.Severity, severities)
		}
		if err := validatePath(finding.File); err != nil {
			return fmt.Errorf(namesTheFile, named, finding.File, err)
		}
		if finding.Line < 1 {
			return fmt.Errorf(linesAreCountedFromOne, named, finding.Line)
		}
		if strings.TrimSpace(finding.Description) == "" {
			return fmt.Errorf(describesNothing, named)
		}
	}

	return nil
}

func validatePath(file string) error {
	if file == "" {
		return errors.New(aFindingNamesItsFile)
	}
	if path.IsAbs(file) || strings.HasPrefix(file, "\\") || strings.Contains(file, ":\\") {
		return errors.New(relativeToTheRoot)
	}
	if file != path.Clean(file) {
		return errors.New(asTheRepositorySpellsIt)
	}
	if file == ".." || strings.HasPrefix(file, "../") {
		return errors.New(insideTheRepository)
	}
	return nil
}
