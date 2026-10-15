package convention

import (
	"html"
	"regexp"
	"slices"
	"strings"
)

const LandingPage = "site/index.html"

var installCommand = regexp.MustCompile(`(?m)^curl -fsSL \S+/` + regexp.QuoteMeta(InstallerFile) + ` \| sh$`)

var markup = regexp.MustCompile(`<[^>]*>`)

func InstallCommandsTheLandingPageOmits(root string) ([]string, error) {
	front, err := Document(root, FrontDocumentFile)
	if err != nil {
		return nil, err
	}
	page, err := Document(root, LandingPage)
	if err != nil {
		return nil, err
	}

	shown := html.UnescapeString(markup.ReplaceAllString(page, ""))

	var omitted []string
	for _, command := range installCommand.FindAllString(front, -1) {
		if !strings.Contains(shown, command) && !slices.Contains(omitted, command) {
			omitted = append(omitted, command)
		}
	}
	return omitted, nil
}
