package render

import (
	_ "embed"
	"fmt"
	"html/template"
	"strings"

	"github.com/rui-ktei/socrates/internal/document"
	"github.com/rui-ktei/socrates/internal/sanitise"
)

//go:embed page.gohtml
var pageSource string

var page = template.Must(template.New("page").Parse(pageSource))

type Frame struct {
	Token    string
	Revision int
	Over     bool
}

type view struct {
	Title     string
	Preamble  string
	Token     string
	Revision  int
	Over      bool
	Decisions []decisionView
}

type decisionView struct {
	ID       string
	Name     string
	Severity string
	Audience string
	Location string
	Question string
	Meaning  string
	Mode     document.Mode
	Multiple bool
	Options  []optionView
	Language string
	Code     string
	Diff     *Diff
	Fragment template.HTML
}

type optionView struct {
	Position int
	Label    string
	Meaning  string
}

func Page(read document.Document, frame Frame) (string, error) {
	var written strings.Builder
	if err := page.Execute(&written, viewOf(read, frame)); err != nil {
		return "", fmt.Errorf(renderingThePage, err)
	}
	return written.String(), nil
}

func viewOf(read document.Document, frame Frame) view {
	shown := view{
		Title:    read.Title,
		Preamble: read.Preamble,
		Token:    frame.Token,
		Revision: frame.Revision,
		Over:     frame.Over,
	}

	for _, decision := range read.Decisions {
		shown.Decisions = append(shown.Decisions, decisionViewOf(decision))
	}
	return shown
}

func named(decision document.Decision) string {
	if strings.TrimSpace(decision.Name) == "" {
		return decision.ID
	}
	return decision.Name
}

func decisionViewOf(decision document.Decision) decisionView {
	shown := decisionView{
		ID:       decision.ID,
		Name:     named(decision),
		Severity: decision.Severity,
		Audience: decision.Audience,
		Location: decision.Location,
		Question: decision.Question,
		Meaning:  decision.Meaning,
		Mode:     decision.Mode,
		Multiple: decision.Mode == document.Many,
		Language: decision.Evidence.Language,
		Code:     decision.Evidence.Code,
		Fragment: template.HTML(sanitise.Fragment(decision.Fragment)),
	}

	if patched := ReadPatch(decision.Evidence.Patch); !patched.Empty() {
		shown.Diff = &patched
	}

	for position, option := range decision.Options {
		shown.Options = append(shown.Options, optionView{
			Position: position,
			Label:    option.Label,
			Meaning:  option.Meaning,
		})
	}
	return shown
}
