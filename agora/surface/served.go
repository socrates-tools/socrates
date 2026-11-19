package surface

import (
	"net/url"
	"slices"

	"github.com/socrates-tools/socrates/agora/document"
)

const DocumentRoute = "/document"

type Served struct {
	Document ServedDocument `json:"document"`
	Standing Standing       `json:"standing"`
}

type Standing struct {
	Revision      int               `json:"revision"`
	Over          bool              `json:"over"`
	Illustrations map[string]string `json:"illustrations"`
	Answered      []string          `json:"answered"`
	FilesPerTray  int               `json:"files_per_tray"`
}

type ServedDocument struct {
	Title     string           `json:"title"`
	Preamble  string           `json:"preamble"`
	Decisions []ServedDecision `json:"decisions"`
	Source    []ServedEntry    `json:"source"`
	Settled   []ServedSettled  `json:"settled"`
}

type ServedDecision struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Severity string         `json:"severity"`
	Audience string         `json:"audience"`
	Location string         `json:"location"`
	Question string         `json:"question"`
	Meaning  string         `json:"meaning"`
	Mode     document.Mode  `json:"mode"`
	Options  []ServedOption `json:"options"`
	Evidence []ServedBlock  `json:"evidence"`
}

type ServedOption struct {
	Label       string `json:"label"`
	Meaning     string `json:"meaning"`
	Recommended bool   `json:"recommended"`
}

type ServedBlock struct {
	Caption  string `json:"caption"`
	Language string `json:"language"`
	Code     string `json:"code"`
	Patch    string `json:"patch"`
}

type ServedEntry struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Severity string `json:"severity"`
	Audience string `json:"audience"`
	Location string `json:"location"`
}

type ServedSettled struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Location string   `json:"location"`
	Question string   `json:"question"`
	Chosen   []string `json:"chosen"`
	Words    string   `json:"words"`
	Attached []string `json:"attached,omitempty"`
}

func servedOf(session string, read document.Document, revision int, over bool) Served {
	return Served{
		Document: servedDocumentOf(read),
		Standing: Standing{
			Revision:      revision,
			Over:          over,
			Illustrations: illustrationAddressesOf(session, read),
			Answered:      arrivingAnswered(read),
			FilesPerTray:  TrayCeiling,
		},
	}
}

func servedDocumentOf(read document.Document) ServedDocument {
	served := ServedDocument{
		Title:     read.Title,
		Preamble:  read.Preamble,
		Decisions: make([]ServedDecision, 0, len(read.Decisions)),
		Source:    make([]ServedEntry, 0, len(read.Source)),
		Settled:   make([]ServedSettled, 0, len(read.Settled)),
	}
	for _, decision := range read.Decisions {
		served.Decisions = append(served.Decisions, servedDecisionOf(decision))
	}
	for _, entry := range read.Source {
		served.Source = append(served.Source, ServedEntry(entry))
	}
	for _, answer := range read.Settled {
		served.Settled = append(served.Settled, ServedSettled{
			ID:       answer.ID,
			Name:     answer.Name,
			Location: answer.Location,
			Question: answer.Question,
			Chosen:   nonNil(answer.Chosen),
			Words:    answer.Words,
			Attached: answer.Attached,
		})
	}
	return served
}

func servedDecisionOf(decision document.Decision) ServedDecision {
	served := ServedDecision{
		ID:       decision.ID,
		Name:     decision.Name,
		Severity: decision.Severity,
		Audience: decision.Audience,
		Location: decision.Location,
		Question: decision.Question,
		Meaning:  decision.Meaning,
		Mode:     decision.Mode,
		Options:  make([]ServedOption, 0, len(decision.Options)),
		Evidence: make([]ServedBlock, 0, len(decision.Evidence)),
	}
	for _, option := range decision.Options {
		served.Options = append(served.Options, ServedOption{
			Label:       option.Label,
			Meaning:     option.Meaning,
			Recommended: option.Recommended,
		})
	}
	for _, block := range decision.Evidence {
		served.Evidence = append(served.Evidence, ServedBlock(block))
	}
	return served
}

func illustrationAddressesOf(session string, read document.Document) map[string]string {
	addresses := map[string]string{}
	for _, decision := range read.Decisions {
		if path := decision.Illustration.Path; path != "" {
			addresses[decision.ID] = SessionPrefix + session + IllustrationRoute + url.PathEscape(decision.ID) + "/" + url.PathEscape(baseName(path))
		}
	}
	return addresses
}

func arrivingAnswered(read document.Document) []string {
	answered := []string{}
	for _, decision := range read.Decisions {
		if slices.ContainsFunc(decision.Options, func(option document.Option) bool { return option.Recommended }) {
			answered = append(answered, decision.ID)
		}
	}
	return answered
}

func nonNil(held []string) []string {
	if held == nil {
		return []string{}
	}
	return held
}
