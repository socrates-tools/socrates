package render

import (
	"bytes"
	"fmt"
	"html/template"
	"strconv"
)

type shellPartView struct {
	Class string
	Label string
}

type appShellView struct {
	Class      string
	Topbar     shellPartView
	Navigation *shellPartView
	Content    shellPartView
	Panel      *shellPartView
	Slots      shellSlots
	Data       any
}

type shellSlots struct {
	Topbar     string
	Navigation string
	Content    string
	Panel      string
}

type badgeView struct {
	Class     string
	Tone      string
	Variant   string
	Dataset   map[string]string
	Title     string
	Hint      string
	Spoken    string
	Focusable bool
	Icon      string
	Label     string
}

type codeBlockView struct {
	Language string
	Code     string
}

type diffBlockView struct {
	Language string
	Patch    string
}

type cardView struct {
	ID      string
	Class   string
	Dataset map[string]string
	Content string
	Data    any
}

type brandView struct {
	Class         string
	MarkClass     string
	WordmarkClass string
	Mark          string
	Name          string
}

type overlineView struct {
	Class string
	Text  string
	Count string
}

type progressBarView struct {
	Class string
	Max   int
	Value int
	Label string
}

type statusDotView struct {
	Class string
	Tone  string
}

type fileTrayView struct {
	Class       string
	Interaction string
	Sample      string
	Max         int
	AttachLabel string
	Refusal     string
	RemoveLabel string
	Content     string
	Data        any
	Disabled    bool
}

type tooltipView struct {
	ID    string
	Class string
}

func overlineOf(class, text string) overlineView {
	return overlineView{Class: class, Text: text}
}

func overlineCounting(class, text string, count int) overlineView {
	return overlineView{Class: class, Text: text, Count: strconv.Itoa(count)}
}

func progressBarOf(class string, max, value int) progressBarView {
	return progressBarView{Class: class, Max: max, Value: value}
}

func statusDotOf(class, tone string) statusDotView {
	return statusDotView{Class: class, Tone: tone}
}

const FilesPerTray = 10

func fileTrayOf(class string, max int) fileTrayView {
	return fileTrayView{
		Class:       class,
		Max:         max,
		AttachLabel: fileTrayAttachLabel,
		Refusal:     fmt.Sprintf(fileTrayRefusal, max),
		RemoveLabel: fileTrayRemoveLabel,
	}
}

func fileTrayFor(class, interaction string, disabled bool, content string, data any) fileTrayView {
	view := fileTrayOf(class, FilesPerTray)
	view.Interaction = interaction
	view.Disabled = disabled
	view.Content = content
	view.Data = data
	return view
}

func filled(root *template.Template, name string, data any) (template.HTML, error) {
	var written bytes.Buffer
	if err := root.ExecuteTemplate(&written, name, data); err != nil {
		return "", err
	}
	return template.HTML(written.String()), nil
}
