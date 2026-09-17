package convention

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

const (
	Organisation = "socrates-tools"
	Suite        = "socrates"
	Repository   = "socrates-review"
)

const moduleDirective = "module"

func ModulePath(root string) (string, error) {
	held, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(string(held), "\n") {
		stated := strings.Fields(line)
		if len(stated) == 2 && stated[0] == moduleDirective {
			return stated[1], nil
		}
	}
	return "", fmt.Errorf("go.mod states no module path")
}

var codeHostingHosts = []string{"github.com", "gitlab.com", "bitbucket.org", "codeberg.org"}

var namesTheProjectAnswersTo = []string{Suite, Repository}

var accountHosting = regexp.MustCompile(
	`(?:` + strings.Join(quoted(codeHostingHosts), "|") + `)/([A-Za-z0-9._-]+)/(?:` + strings.Join(quoted(namesTheProjectAnswersTo), "|") + `)\b`,
)

type NamedAccount struct {
	File    string
	Line    int
	Account string
}

func AccountsOtherThanTheOrganisation(root string) ([]NamedAccount, error) {
	var named []NamedAccount

	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walked error) error {
		if walked != nil {
			return walked
		}

		relative, relErr := filepath.Rel(root, path)
		if relErr != nil {
			relative = path
		}

		if entry.IsDir() {
			if path != root && skipped(path, entry.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if !entry.Type().IsRegular() {
			return nil
		}

		held, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		if !utf8.Valid(held) {
			return nil
		}

		named = append(named, accountsNamedIn(relative, string(held))...)
		return nil
	})
	if err != nil {
		return nil, err
	}

	slices.SortFunc(named, func(one, other NamedAccount) int {
		if one.File != other.File {
			return strings.Compare(one.File, other.File)
		}
		return one.Line - other.Line
	})
	return named, nil
}

func accountsNamedIn(relative, held string) []NamedAccount {
	var named []NamedAccount

	for position, line := range strings.Split(held, "\n") {
		for _, hosted := range accountHosting.FindAllStringSubmatch(line, -1) {
			if hosted[1] == Organisation {
				continue
			}
			named = append(named, NamedAccount{File: relative, Line: position + 1, Account: hosted[1]})
		}
	}
	return named
}

func quoted(names []string) []string {
	escaped := make([]string, 0, len(names))
	for _, name := range names {
		escaped = append(escaped, regexp.QuoteMeta(name))
	}
	return escaped
}
