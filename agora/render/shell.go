package render

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
)

//go:embed shell.gohtml
var shellSource string

var shellPage = template.Must(template.New("shell").Parse(shellSource))

type ShellFrame struct {
	Base  string
	Token string
}

type shellView struct {
	ShellFrame
	Title       string
	Policy      string
	NeedsScript string
}

func Shell(frame ShellFrame) (string, error) {
	var written bytes.Buffer
	shown := shellView{ShellFrame: frame, Title: shellTitle, Policy: policy, NeedsScript: pageNeedsJavaScript}
	if err := shellPage.Execute(&written, shown); err != nil {
		return "", fmt.Errorf(renderingTheShell, err)
	}
	return written.String(), nil
}
