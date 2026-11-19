package surface

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
)

//go:embed catalogshell.gohtml
var catalogShellSource string

var catalogShellPage = template.Must(template.New("catalog-shell").Parse(catalogShellSource))

type catalogShellView struct {
	Title  string
	Policy string
}

func CatalogShell() (string, error) {
	var written bytes.Buffer
	if err := catalogShellPage.Execute(&written, catalogShellView{Title: catalogTitle, Policy: pagePolicy}); err != nil {
		return "", fmt.Errorf(renderingTheCatalogShell, err)
	}
	return written.String(), nil
}
