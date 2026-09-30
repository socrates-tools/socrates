package render

import (
	"fmt"
	"html/template"
)

type markView struct {
	Kind      string
	Focusable bool
	Value     string
	Tone      string
	Icon      string
	Title     string
	Hint      string
	Spoken    string
}

type severityStated struct {
	tone  string
	icon  string
	title string
	hint  string
}

var severitiesKnown = map[string]severityStated{
	"error":   {tone: "error", icon: "circle-alert", title: severityError, hint: severityErrorHint},
	"warning": {tone: "warning", icon: "triangle-alert", title: severityWarning, hint: severityWarningHint},
	"info":    {tone: "neutral", icon: "info", title: severityInfo, hint: severityInfoHint},
}

func severityMark(value string) *markView {
	if value == "" {
		return nil
	}
	stated, known := severitiesKnown[value]
	if !known {
		stated = severityStated{tone: "neutral", icon: "info", title: fmt.Sprintf(severityNamed, value)}
	}
	return marked(markView{Kind: "severity", Value: value, Tone: stated.tone, Icon: stated.icon, Title: stated.title, Hint: stated.hint})
}

func audienceMark(value string) *markView {
	if value == "" {
		return nil
	}
	stated, known := audiencesSocratesReviewWrites[value]
	if !known {
		stated = audienceStated{tone: "neutral", icon: "info", title: fmt.Sprintf(audienceNamed, value)}
	}
	return marked(markView{Kind: "audience", Value: value, Tone: stated.tone, Icon: stated.icon, Title: stated.title, Hint: stated.hint})
}

func marked(shown markView) *markView {
	shown.Focusable = true
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
	return &quiet
}

func (m markView) Drawn() template.HTML {
	return icon(m.Icon)
}
