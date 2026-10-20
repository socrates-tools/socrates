package render

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/url"
	"slices"
	"strings"

	"github.com/socrates-tools/socrates/agora/document"
)

const (
	IllustrationRoute = "/illustrations/"

	policy = "default-src 'none'; style-src 'self'; style-src-attr 'none'; script-src 'self'; font-src 'self'; img-src 'self' data:; connect-src 'self'; frame-src 'self'; form-action 'none'; base-uri 'none'; frame-ancestors 'none'"
)

var mostSevereFirst = []string{"error", "warning", "info"}

const chosenJoined = " · "

//go:embed page.gohtml
var pageSource string

var page = template.Must(template.New("page").Parse(pageSource))

type Frame struct {
	Base     string
	Token    string
	Revision int
	Over     bool
}

type view struct {
	Title     string
	Preamble  string
	Base      string
	Token     string
	Revision  int
	Over      bool
	Policy    string
	Index     []entryView
	Settled   []settledView
	Decisions []decisionView
}

type entryView struct {
	ID       string
	Name     string
	Severity *markView
	Location string
	Decided  bool
	Answered bool
}

type settledView struct {
	ID       string
	Name     string
	Location string
	Question string
	Chosen   string
	Words    string
}

type decisionView struct {
	ID       string
	Name     string
	Severity *markView
	Audience *markView
	Location string
	Question string
	Meaning  string
	Mode     document.Mode
	Multiple bool
	Options  []optionView
	Evidence []blockView
	Drawing  *drawingView
	Answered bool
}

type blockView struct {
	Caption  string
	Language string
	Code     string
	Diff     *Diff
	Written  string
}

type drawingView struct {
	Source string
	Title  string
}

type optionView struct {
	Position    int
	Label       string
	Meaning     string
	Recommended bool
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
		Base:     frame.Base,
		Token:    frame.Token,
		Revision: frame.Revision,
		Over:     frame.Over,
		Policy:   policy,
	}

	answered := map[string]bool{}
	for _, decision := range read.Decisions {
		shown.Decisions = append(shown.Decisions, decisionViewOf(decision, frame))
		answered[decision.ID] = arrivesAnswered(decision)
	}

	for _, entry := range read.Source {
		held, decided := answered[entry.ID]
		shown.Index = append(shown.Index, entryView{
			ID:       entry.ID,
			Name:     nameOr(entry.Name, entry.ID),
			Severity: quietly(severityMark(entry.Severity)),
			Location: entry.Location,
			Decided:  decided,
			Answered: held,
		})
	}
	if len(read.Source) == 0 {
		for _, decision := range shown.Decisions {
			shown.Index = append(shown.Index, entryView{
				ID:       decision.ID,
				Name:     decision.Name,
				Severity: quietly(decision.Severity),
				Location: decision.Location,
				Decided:  true,
				Answered: decision.Answered,
			})
		}
	}

	for _, answer := range read.Settled {
		shown.Settled = append(shown.Settled, settledView{
			ID:       answer.ID,
			Name:     nameOr(answer.Name, answer.ID),
			Location: answer.Location,
			Question: answer.Question,
			Chosen:   strings.Join(answer.Chosen, chosenJoined),
			Words:    strings.TrimSpace(answer.Words),
		})
	}

	slices.SortStableFunc(shown.Index, func(one, other entryView) int {
		return rankOf(one.Severity) - rankOf(other.Severity)
	})
	slices.SortStableFunc(shown.Decisions, func(one, other decisionView) int {
		return rankOf(one.Severity) - rankOf(other.Severity)
	})
	return shown
}

func (v view) AnsweredCount() int {
	count := 0
	for _, decision := range v.Decisions {
		if decision.Answered {
			count++
		}
	}
	return count
}

func arrivesAnswered(decision document.Decision) bool {
	return slices.ContainsFunc(decision.Options, func(option document.Option) bool { return option.Recommended })
}

func rankOf(severity *markView) int {
	if severity == nil {
		return len(mostSevereFirst)
	}
	if at := slices.Index(mostSevereFirst, severity.Value); at >= 0 {
		return at
	}
	return len(mostSevereFirst)
}

func nameOr(name, identifier string) string {
	if strings.TrimSpace(name) == "" {
		return identifier
	}
	return name
}

func decisionViewOf(decision document.Decision, frame Frame) decisionView {
	shown := decisionView{
		ID:       decision.ID,
		Name:     nameOr(decision.Name, decision.ID),
		Severity: severityMark(decision.Severity),
		Audience: audienceMark(decision.Audience),
		Location: decision.Location,
		Question: decision.Question,
		Meaning:  decision.Meaning,
		Mode:     decision.Mode,
		Multiple: decision.Mode == document.Many,
		Evidence: blocksRead(decision.Evidence),
		Answered: arrivesAnswered(decision),
	}

	if path := decision.Illustration.Path; path != "" {
		shown.Drawing = &drawingView{
			Source: frame.Base + IllustrationRoute + url.PathEscape(decision.ID) + "/" + url.PathEscape(BaseName(path)),
			Title:  fmt.Sprintf(illustrationTitle, shown.Name),
		}
	}

	shown.Options = optionsRead(decision.Options)
	return shown
}

func blocksRead(written document.Evidence) []blockView {
	var shown []blockView

	for _, block := range written {
		read := blockView{Caption: block.Caption, Language: block.Language, Code: block.Code}
		if block.Patch != "" {
			patched := ReadPatch(block.Patch)
			if patched.Empty() {
				continue
			}
			read.Diff = &patched
			read.Written = writtenIn(block.Language, patched)
		}
		if read.Code == "" && read.Diff == nil {
			continue
		}
		shown = append(shown, read)
	}
	return shown
}

func BaseName(path string) string {
	return path[strings.LastIndexAny(path, `/\`)+1:]
}

func writtenIn(stated string, patched Diff) string {
	if stated != "" {
		return stated
	}
	named := patched.To
	if named == "" || named == "/dev/null" {
		named = patched.From
	}
	at := strings.LastIndex(named, ".")
	if at < 0 || at < strings.LastIndex(named, "/") {
		return ""
	}
	return strings.ToLower(named[at+1:])
}

func optionsRead(offered []document.Option) []optionView {
	var recommended, rest []optionView

	for position, option := range offered {
		read := optionView{
			Position:    position,
			Label:       option.Label,
			Meaning:     option.Meaning,
			Recommended: option.Recommended,
		}
		if option.Recommended {
			recommended = append(recommended, read)
			continue
		}
		rest = append(rest, read)
	}
	return append(recommended, rest...)
}
