package render

import (
	"bytes"
	_ "embed"
	"html/template"
	"slices"
)

//go:embed catalog.gohtml
var catalogSource string

type dialogView struct {
	ID          string
	Class       string
	Interaction string
	Heading     string
	Body        string
	Actions     string
}

type state struct {
	Name  string
	Opens string
	View  any
}

func dialogState(name string, view dialogView) state {
	return state{Name: name, Opens: view.ID, View: view}
}

var catalogued = map[string][]state{
	"ui-dialog": {
		dialogState(dialogWithEverything, dialogView{
			ID:      "catalog-dialog",
			Class:   "catalog-dialog",
			Heading: catalogDialogHeading,
			Body:    "catalog-dialog-body",
			Actions: "catalog-dialog-actions",
		}),
		dialogState(dialogWithAHeadingAlone, dialogView{
			ID:      "catalog-dialog-bare",
			Class:   "catalog-dialog-bare",
			Heading: catalogDialogHeading,
		}),
	},
}

type catalogState struct {
	Name  string
	Opens string
	Drawn template.HTML
}

type catalogPrimitive struct {
	Name   string
	States []catalogState
}

type catalogView struct {
	Title      string
	Policy     string
	Primitives []catalogPrimitive
}

func Catalog() (string, error) {
	names := make([]string, 0, len(catalogued))
	for name := range catalogued {
		names = append(names, name)
	}
	slices.Sort(names)

	view := catalogView{Title: catalogTitle, Policy: policy}
	for _, name := range names {
		primitive := catalogPrimitive{Name: name}
		for _, held := range catalogued[name] {
			markup, err := partialDrawn(name, held.View)
			if err != nil {
				return "", err
			}
			primitive.States = append(primitive.States, catalogState{Name: held.Name, Opens: held.Opens, Drawn: markup})
		}
		view.Primitives = append(view.Primitives, primitive)
	}

	var written bytes.Buffer
	if err := page.ExecuteTemplate(&written, "catalog", view); err != nil {
		return "", err
	}
	return written.String(), nil
}
