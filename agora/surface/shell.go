package surface

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"strings"
)

const (
	IllustrationRoute = "/illustrations/"

	pagePolicy = "default-src 'none'; style-src 'self'; style-src-attr 'none'; script-src 'self'; font-src 'self'; img-src 'self' data: blob:; connect-src 'self'; frame-src 'self'; form-action 'none'; base-uri 'none'; frame-ancestors 'none'"
)

//go:embed shell.gohtml
var shellSource string

var shellPage = template.Must(template.New("shell").Parse(shellSource))

type shellFrame struct {
	Base  string
	Token string
}

type shellView struct {
	shellFrame
	Title       string
	Policy      string
	NeedsScript string
}

func writeShell(frame shellFrame) (string, error) {
	var written bytes.Buffer
	shown := shellView{shellFrame: frame, Title: shellTitle, Policy: pagePolicy, NeedsScript: pageNeedsJavaScript}
	if err := shellPage.Execute(&written, shown); err != nil {
		return "", fmt.Errorf(renderingTheShell, err)
	}
	return written.String(), nil
}

func baseName(path string) string {
	return path[strings.LastIndexAny(path, `/\`)+1:]
}
