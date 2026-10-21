package render

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"strings"
)

const primitivePrefix = "ui-"

//go:embed partials/*.gohtml
var partialSources embed.FS

type dialogView struct {
	Class       string
	Interaction string
	Body        string
	Actions     string
}

var drawnWith = map[string]any{
	"ui-dialog": dialogView{Class: "dialog"},
}

var page = parsedPage()

func parsedPage() *template.Template {
	var root *template.Template
	root = template.Must(template.New("page").Funcs(template.FuncMap{
		"fill": func(name string) (template.HTML, error) {
			var written bytes.Buffer
			if err := root.ExecuteTemplate(&written, name, nil); err != nil {
				return "", err
			}
			return template.HTML(written.String()), nil
		},
	}).Parse(pageSource))
	return template.Must(root.ParseFS(partialSources, "partials/*.gohtml"))
}

func Primitives() (map[string]string, error) {
	drawn := map[string]string{}
	for _, held := range page.Templates() {
		name := held.Name()
		if !strings.HasPrefix(name, primitivePrefix) || strings.Contains(name, ".") {
			continue
		}
		data, known := drawnWith[name]
		if !known {
			return nil, fmt.Errorf(partialDrawnWithNothing, name)
		}
		var written bytes.Buffer
		if err := held.Execute(&written, data); err != nil {
			return nil, fmt.Errorf(drawingAPartial, name, err)
		}
		drawn[name] = written.String()
	}
	return drawn, nil
}
