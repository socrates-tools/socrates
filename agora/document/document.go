package document

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
)

type Mode string

const (
	One  Mode = "single-choice"
	Many Mode = "multiple-choice"
)

var modes = []Mode{One, Many}

func Modes() []Mode {
	return slices.Clone(modes)
}

type Document struct {
	Title     string     `json:"title"`
	Preamble  string     `json:"preamble"`
	Decisions []Decision `json:"decisions"`
	Source    []Entry    `json:"source"`
	Settled   []Settled  `json:"settled"`
}

type Settled struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Location string   `json:"location"`
	Question string   `json:"question"`
	Chosen   []string `json:"chosen"`
	Words    string   `json:"words"`
	Attached []string `json:"attached,omitempty"`
}

type Entry struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Severity string `json:"severity"`
	Audience string `json:"audience"`
	Location string `json:"location"`
}

type Decision struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Severity     string       `json:"severity"`
	Audience     string       `json:"audience"`
	Location     string       `json:"location"`
	Question     string       `json:"question"`
	Meaning      string       `json:"meaning"`
	Mode         Mode         `json:"mode"`
	Options      []Option     `json:"options"`
	Evidence     Evidence     `json:"evidence"`
	Illustration Illustration `json:"illustration"`
}

type Illustration struct {
	Path string `json:"path"`
}

type Option struct {
	Label       string `json:"label"`
	Meaning     string `json:"meaning"`
	Command     string `json:"command"`
	Recommended bool   `json:"recommended"`
}

type Evidence []Block

type Block struct {
	Caption  string `json:"caption"`
	Language string `json:"language"`
	Code     string `json:"code"`
	Patch    string `json:"patch"`
}

func Read(raw []byte) (Document, error) {
	var read Document

	if err := refuseEvidenceWrittenAsAnObject(raw); err != nil {
		return Document{}, err
	}

	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&read); err != nil {
		return Document{}, fmt.Errorf(readingTheDocument, err)
	}
	if decoder.More() {
		return Document{}, errors.New(moreThanOneDocument)
	}

	if err := Validate(read); err != nil {
		return Document{}, err
	}
	return read.copiedFromSource(), nil
}

func refuseEvidenceWrittenAsAnObject(raw []byte) error {
	var written struct {
		Decisions []struct {
			ID       string          `json:"id"`
			Evidence json.RawMessage `json:"evidence"`
		} `json:"decisions"`
	}
	_ = json.Unmarshal(raw, &written)
	for position, decision := range written.Decisions {
		if bytes.HasPrefix(bytes.TrimSpace(decision.Evidence), []byte("{")) {
			return fmt.Errorf(evidenceWrittenAsAnObject, decisionName(position, decision.ID))
		}
	}
	return nil
}

func decisionName(position int, identifier string) string {
	if identifier == "" {
		return fmt.Sprintf(decisionAt, position+1)
	}
	return fmt.Sprintf(decisionNamed, identifier)
}

func Validate(read Document) error {
	if strings.TrimSpace(read.Title) == "" {
		return fmt.Errorf(theDocumentStatesNothingFor, "title")
	}
	if len(read.Decisions) == 0 && len(read.Source) == 0 && len(read.Settled) == 0 && strings.TrimSpace(read.Preamble) == "" {
		return errors.New(theDocumentHoldsNothing)
	}

	cited, err := read.citedByIdentifier()
	if err != nil {
		return err
	}

	seen := map[string]bool{}
	for position, decision := range read.Decisions {
		named := decisionName(position, decision.ID)

		if decision.ID == "" {
			return fmt.Errorf(noIdentifier, named)
		}
		if seen[decision.ID] {
			return fmt.Errorf(repeatedIdentifier, decision.ID)
		}
		seen[decision.ID] = true

		if strings.TrimSpace(decision.Question) == "" {
			return fmt.Errorf(statesNothingFor, named, "question")
		}
		for shown, block := range decision.Evidence {
			if err := validateBlock(named, shown+1, block); err != nil {
				return err
			}
		}
		if len(decision.Options) == 0 && decision.Mode != "" {
			return fmt.Errorf(aModeWithoutOptions, named, decision.Mode)
		}
		if len(decision.Options) > 0 && !slices.Contains(modes, decision.Mode) {
			return fmt.Errorf(unknownMode, named, decision.Mode, modes)
		}
		recommended, commanded, bare := 0, 0, 0
		for offered, option := range decision.Options {
			if strings.TrimSpace(option.Label) == "" {
				return fmt.Errorf(optionStatesNothingFor, offered+1, named, "label")
			}
			if strings.TrimSpace(option.Command) == "" {
				bare++
			} else {
				commanded++
			}
			if option.Recommended {
				recommended++
			}
		}
		if commanded > 0 && bare > 0 {
			return fmt.Errorf(decisionMixesCommandedOptions, named)
		}
		if recommended > 1 {
			return fmt.Errorf(moreThanOneRecommended, named, recommended)
		}

		if len(cited) > 0 {
			if _, citing := cited[decision.ID]; !citing {
				return fmt.Errorf(namesNothingInTheSource, named)
			}
		}
	}

	return validateWhatWasSettled(read.Settled, seen)
}

func validateBlock(named string, position int, block Block) error {
	holdsCode := block.Code != ""
	holdsPatch := block.Patch != ""
	if holdsCode && holdsPatch {
		return fmt.Errorf(blockHoldsCodeAndAPatch, named, position)
	}
	if !holdsCode && !holdsPatch {
		return fmt.Errorf(blockHoldsNeitherCodeNorAPatch, named, position)
	}
	if stated := block.Language; stated != "" && !slices.Contains(languages, stated) {
		return fmt.Errorf(unknownLanguage, named, stated, languages)
	}
	return nil
}

func validateWhatWasSettled(settled []Settled, open map[string]bool) error {
	seen := map[string]bool{}
	for position, answered := range settled {
		if answered.ID == "" {
			return fmt.Errorf(settledAnswerCarriesNoIdentifier, position+1)
		}
		if seen[answered.ID] {
			return fmt.Errorf(repeatedSettledIdentifier, answered.ID)
		}
		seen[answered.ID] = true

		if open[answered.ID] {
			return fmt.Errorf(settledAndStillOpen, answered.ID)
		}
		if strings.TrimSpace(answered.Question) == "" {
			return fmt.Errorf(settledStatesNothingFor, answered.ID, "question")
		}
		for chosen, label := range answered.Chosen {
			if strings.TrimSpace(label) == "" {
				return fmt.Errorf(settledChoiceStatesNothing, answered.ID, chosen+1)
			}
		}
		for attached, name := range answered.Attached {
			if strings.TrimSpace(name) == "" {
				return fmt.Errorf(settledAttachmentStatesNothing, answered.ID, attached+1)
			}
		}
	}
	return nil
}

func (d Document) citedByIdentifier() (map[string]Entry, error) {
	if len(d.Source) == 0 {
		return nil, nil
	}

	cited := map[string]Entry{}
	for position, entry := range d.Source {
		if entry.ID == "" {
			return nil, fmt.Errorf(sourceEntryCarriesNoIdentifier, position+1)
		}
		if _, twice := cited[entry.ID]; twice {
			return nil, fmt.Errorf(repeatedSourceIdentifier, entry.ID)
		}
		cited[entry.ID] = entry
	}
	return cited, nil
}

func (d Document) copiedFromSource() Document {
	cited, err := d.citedByIdentifier()
	if err != nil || len(cited) == 0 {
		return d
	}

	copied := d
	copied.Decisions = slices.Clone(d.Decisions)
	for position, decision := range copied.Decisions {
		entry := cited[decision.ID]
		decision.Name = entry.Name
		decision.Severity = entry.Severity
		decision.Audience = entry.Audience
		decision.Location = entry.Location
		copied.Decisions[position] = decision
	}
	return copied
}
