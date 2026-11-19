package render

import (
	"bytes"
	_ "embed"
	"html/template"
)

//go:embed catalogshell.gohtml
var catalogShellSource string

var catalogShellPage = template.Must(template.New("catalog-shell").Parse(catalogShellSource))

type catalogShellFrame struct {
	Title  string
	Policy string
}

func CatalogShell() (string, error) {
	var written bytes.Buffer
	if err := catalogShellPage.Execute(&written, catalogShellFrame{Title: catalogTitle, Policy: policy}); err != nil {
		return "", err
	}
	return written.String(), nil
}
