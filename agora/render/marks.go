package render

import "fmt"

type markView struct {
	Kind      string
	Focusable bool
	Named     bool
	Value     string
	Tone      string
	Icon      string
	Title     string
	Name      string
	Hint      string
	Spoken    string
}

type severityStated struct {
	tone  string
	icon  string
	name  string
	title string
	hint  string
}

var severitiesKnown = map[string]severityStated{
	"error":   {tone: "error", icon: "circle-alert", name: severityErrorName, title: severityError, hint: severityErrorHint},
	"warning": {tone: "warning", icon: "triangle-alert", name: severityWarningName, title: severityWarning, hint: severityWarningHint},
	"info":    {tone: "neutral", icon: "info", name: severityInfoName, title: severityInfo, hint: severityInfoHint},
}

func severityMark(value string) *markView {
	if value == "" {
		return nil
	}
	stated, known := severitiesKnown[value]
	if !known {
		stated = severityStated{tone: "neutral", icon: "info", name: value, title: fmt.Sprintf(severityNamed, value)}
	}
	return marked(markView{Kind: "severity", Value: value, Tone: stated.tone, Icon: stated.icon, Name: stated.name, Title: stated.title, Hint: stated.hint})
}

func audienceMark(value string) *markView {
	if value == "" {
		return nil
	}
	stated, known := audiencesSocratesReviewWrites[value]
	if !known {
		stated = audienceStated{tone: "neutral", icon: "info", title: fmt.Sprintf(audienceNamed, value)}
	}
	return marked(markView{Kind: "audience", Value: value, Tone: stated.tone, Icon: stated.icon, Name: stated.title, Title: stated.title, Hint: stated.hint})
}

func indexedAudienceMark(value string) *markView {
	if !audiencesSocratesReviewWrites[value].askedOfMe {
		return nil
	}
	return quietly(audienceMark(value))
}

var variantOfKind = map[string]string{
	"severity": "filled",
	"audience": "outlined",
}

func marked(shown markView) *markView {
	shown.Focusable = true
	shown.Named = true
	shown.Spoken = shown.Title
	if shown.Hint != "" {
		shown.Spoken = fmt.Sprintf(spokenWithItsHint, shown.Title, shown.Hint)
	}
	return &shown
}

func quietly(mark *markView) *markView {
	if mark == nil {
		return nil
	}
	quiet := *mark
	quiet.Focusable = false
	quiet.Named = false
	return &quiet
}

func (m markView) Badge() badgeView {
	shown := badgeView{
		Class:     "mark " + m.Kind,
		Tone:      m.Tone,
		Variant:   variantOfKind[m.Kind],
		Dataset:   map[string]string{m.Kind: m.Value},
		Title:     m.Title,
		Hint:      m.Hint,
		Spoken:    m.Spoken,
		Focusable: m.Focusable,
		Icon:      m.Icon,
	}
	if m.Named {
		shown.Label = m.Name
	}
	return shown
}
