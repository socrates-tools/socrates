package render

import (
	_ "embed"
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/socrates-tools/socrates/agora/document"
)

const (
	IllustrationRoute = "/illustrations/"

	policy = "default-src 'none'; style-src 'self'; style-src-attr 'none'; script-src 'self'; font-src 'self'; img-src 'self' data:; connect-src 'self'; frame-src 'self'; form-action 'none'; base-uri 'none'; frame-ancestors 'none'"
)

var mostSevereFirst = []string{"error", "warning", "info"}

const chosenJoined = " · "

const viewing = "viewing"

//go:embed page.gohtml
var pageSource string

type Frame struct {
	Base     string
	Token    string
	Revision int
	Over     bool
}

type view struct {
	Title        string
	Preamble     string
	Base         string
	Token        string
	Revision     int
	Over         bool
	Policy       string
	Confirmation dialogView
	Tooltip      tooltipView
	Composer     cardView
	Shell        appShellView
	Nothing      cardView
	Index        []entryView
	Settled      []settledView
	Decisions    []decisionView
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
	Over     bool
}

func (d decisionView) Card() cardView {
	state := "open"
	if d.Over {
		state = "closed"
	}
	return cardView{
		ID:    "decision-" + d.ID,
		Class: "decision",
		Dataset: map[string]string{
			"decision": d.ID,
			"mode":     string(d.Mode),
			"state":    state,
		},
		Content: "decision-content",
		Data:    d,
	}
}

func (d decisionView) ChoiceGroup() choiceGroupView {
	legend := choiceOneLegend
	if d.Multiple {
		legend = choiceManyLegend
	}
	choices := make([]choiceView, 0, len(d.Options)+1)
	for _, option := range d.Options {
		choices = append(choices, choiceView{
			Class:       "option choice",
			Family:      "option",
			Variant:     "solid",
			Value:       strconv.Itoa(option.Position),
			Interaction: answering,
			Decision:    d.ID,
			Label:       option.Label,
			Meaning:     option.Meaning,
			Tag:         option.Tag,
			Recommended: option.Recommended,
			Checked:     option.Recommended,
			Disabled:    d.Over,
		})
	}
	if !d.Multiple {
		choices = append(choices, choiceView{
			Class:       "none decline",
			Family:      "option",
			Variant:     "dashed",
			Interaction: answering,
			Decision:    d.ID,
			Label:       noneOfThese,
			Meaning:     noneOfTheseMeaning,
			Disabled:    d.Over,
		})
	}
	return choiceGroupView{Class: "options-list", Label: legend, Multiple: d.Multiple, Choices: choices}
}

type blockView struct {
	Caption  string
	Language string
	Code     string
	Diff     *Diff
	Patch    string
	Written  string
	Copy     *copyView
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
	Tag         string
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
		Title:        read.Title,
		Preamble:     read.Preamble,
		Base:         frame.Base,
		Token:        frame.Token,
		Revision:     frame.Revision,
		Over:         frame.Over,
		Policy:       policy,
		Confirmation: dialogView{Class: "confirmation", Interaction: "sending", Body: "confirmation-body", Actions: "confirmation-actions"},
		Tooltip:      tooltipView{ID: "tooltip", Class: "tooltip"},
		Nothing:      cardView{Class: "nothing", Content: "nothing-content"},
	}
	shown.Composer = composerCard(shown)

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
	shown.Shell = shellOf(shown)
	return shown
}

func shellOf(v view) appShellView {
	shell := appShellView{
		Topbar:  shellPartView{Class: "topbar"},
		Content: shellPartView{Class: "main"},
		Slots:   shellSlots{Topbar: "shell-topbar", Content: "shell-content"},
		Data:    v,
	}
	if len(v.Index) > 0 {
		shell.Navigation = &shellPartView{Class: "index", Label: whatWasSeen}
		shell.Slots.Navigation = "shell-navigation"
	}
	if len(v.Decisions) > 0 {
		shell.Panel = &shellPartView{Class: "panel", Label: reply}
		shell.Slots.Panel = "shell-panel"
	}
	return shell
}

func composerCard(v view) cardView {
	phase := "open"
	if v.Over {
		phase = "over"
	}
	return cardView{
		Class:   "composer",
		Dataset: map[string]string{"phase": phase},
		Content: "composer-content",
		Data:    v,
	}
}

func (v view) Tracked() bool {
	return slices.ContainsFunc(v.Index, func(entry entryView) bool { return entry.Decided })
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
		Over:     frame.Over,
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
		if block.Code != "" {
			read.Copy = &copyView{Label: copyLabel, Hint: copyHint, Copied: copiedWord, Failed: notCopied, Interaction: viewing}
		}
		if block.Patch != "" {
			patched := ReadPatch(block.Patch)
			if patched.Empty() {
				continue
			}
			read.Diff = &patched
			read.Patch = block.Patch
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
			read.Tag = recommendedTag
			recommended = append(recommended, read)
			continue
		}
		rest = append(rest, read)
	}
	return append(recommended, rest...)
}
