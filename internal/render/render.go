package render

import (
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"slices"
	"strings"

	"github.com/socrates-tools/socrates-review/internal/document"
	"github.com/socrates-tools/socrates-review/internal/sanitise"
	"github.com/socrates-tools/socrates-review/internal/stylesheet"
)

const (
	PaintedClass = "painted"
	NamePrefix   = "drawn-"

	policy         = "default-src 'none'; style-src 'self'%s; style-src-attr 'none'; script-src 'self'; font-src 'self'; img-src 'self' data:; connect-src 'self'; form-action 'none'; base-uri 'none'; frame-ancestors 'none'"
	admittedByHash = " 'sha256-%s'"
)

var Illustration = stylesheet.Scope{Selector: "." + PaintedClass, Prefix: NamePrefix}

var mostSevereFirst = []string{"error", "warning", "info"}

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
	Title      string
	Preamble   string
	Base       string
	Token      string
	Revision   int
	Over       bool
	Policy     string
	Stylesheet template.CSS
	Source     []entryView
	Decisions  []decisionView
}

type entryView struct {
	ID       string
	Name     string
	Severity string
	Audience string
	Location string
	Decided  bool
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
	painted := stylesheet.Reduce(read.Stylesheet, Illustration)

	shown := view{
		Title:      read.Title,
		Preamble:   read.Preamble,
		Base:       frame.Base,
		Token:      frame.Token,
		Revision:   frame.Revision,
		Over:       frame.Over,
		Policy:     policyAdmitting(painted),
		Stylesheet: template.CSS(painted),
	}

	decided := map[string]bool{}
	for _, decision := range read.Decisions {
		decided[decision.ID] = true
	}

	for _, entry := range read.Source {
		shown.Source = append(shown.Source, entryViewOf(entry, decided[entry.ID]))
	}
	for _, decision := range read.Decisions {
		shown.Decisions = append(shown.Decisions, decisionViewOf(decision))
	}

	slices.SortStableFunc(shown.Source, func(one, other entryView) int {
		return rankOf(one.Severity) - rankOf(other.Severity)
	})
	slices.SortStableFunc(shown.Decisions, func(one, other decisionView) int {
		return rankOf(one.Severity) - rankOf(other.Severity)
	})
	return shown
}

func rankOf(severity string) int {
	if at := slices.Index(mostSevereFirst, severity); at >= 0 {
		return at
	}
	return len(mostSevereFirst)
}

func entryViewOf(entry document.Entry, decided bool) entryView {
	name := entry.Name
	if strings.TrimSpace(name) == "" {
		name = entry.ID
	}
	return entryView{
		ID:       entry.ID,
		Name:     name,
		Severity: entry.Severity,
		Audience: entry.Audience,
		Location: entry.Location,
		Decided:  decided,
	}
}

func policyAdmitting(painted string) string {
	if painted == "" {
		return fmt.Sprintf(policy, "")
	}

	digest := sha256.Sum256([]byte(painted))
	return fmt.Sprintf(policy, fmt.Sprintf(admittedByHash, base64.StdEncoding.EncodeToString(digest[:])))
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
		Fragment: template.HTML(sanitise.Fragment(decision.Fragment.Markup, NamePrefix)),
	}

	if patched := ReadPatch(decision.Evidence.Patch); !patched.Empty() {
		shown.Diff = &patched
	}

	shown.Options = optionsRead(decision.Options)
	return shown
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
