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

var page = parsedPage()

func parsedPage() *template.Template {
	var root *template.Template
	root = template.Must(template.New("page").Funcs(template.FuncMap{
		"dict": func(pairs ...any) (map[string]any, error) {
			if len(pairs)%2 != 0 {
				return nil, fmt.Errorf(dictionaryNeedsPairs, len(pairs))
			}
			held := make(map[string]any, len(pairs)/2)
			for at := 0; at < len(pairs); at += 2 {
				key, named := pairs[at].(string)
				if !named {
					return nil, fmt.Errorf(dictionaryKeyNotAWord, pairs[at])
				}
				held[key] = pairs[at+1]
			}
			return held, nil
		},
		"fill":            func(name string) (template.HTML, error) { return filled(root, name, nil) },
		"fillWith":        func(name string, data any) (template.HTML, error) { return filled(root, name, data) },
		"overline":        overlineOf,
		"overlineCounted": overlineCounting,
		"progressBar":     progressBarOf,
		"statusDot":       statusDotOf,
	}).Parse(pageSource))
	template.Must(root.ParseFS(partialSources, "partials/*.gohtml"))
	template.Must(root.New("catalog").Parse(catalogSource))
	return root
}

func partialDrawn(name string, data any) (template.HTML, error) {
	var written bytes.Buffer
	if err := page.ExecuteTemplate(&written, name, data); err != nil {
		return "", fmt.Errorf(drawingAPartial, name, err)
	}
	return template.HTML(written.String()), nil
}

func Primitives() (map[string]string, error) {
	written := map[string]string{}
	for _, held := range page.Templates() {
		name := held.Name()
		if !strings.HasPrefix(name, primitivePrefix) || strings.Contains(name, ".") {
			continue
		}
		states, known := catalogued[name]
		if !known || len(states) == 0 {
			return nil, fmt.Errorf(partialDrawnWithNothing, name)
		}
		for _, each := range states {
			markup, err := partialDrawn(name, each.View)
			if err != nil {
				return nil, err
			}
			written[name] += string(markup)
		}
	}
	return written, nil
}
