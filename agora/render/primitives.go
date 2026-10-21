package render

import (
	"bytes"
	"html/template"
	"strconv"
)

type shellPartView struct {
	Class string
	Label string
}

type appShellView struct {
	Class       string
	BodyClass   string
	Form        bool
	Interaction string
	Topbar      shellPartView
	Navigation  *shellPartView
	Content     shellPartView
	Panel       *shellPartView
	Slots       shellSlots
	Data        any
}

type shellSlots struct {
	Topbar     string
	Navigation string
	Content    string
	Panel      string
	Overlay    string
}

type badgeView struct {
	Class     string
	Tone      string
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
	Hunks    []Hunk
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

func filled(root *template.Template, name string, data any) (template.HTML, error) {
	var written bytes.Buffer
	if err := root.ExecuteTemplate(&written, name, data); err != nil {
		return "", err
	}
	return template.HTML(written.String()), nil
}
