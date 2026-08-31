package convention

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

const (
	Organisation = "ktei"
	Repository   = "socrates"
)

var ArchivedChanges = filepath.Join("openspec", "changes", "archive")

var codeHostingHosts = []string{"github.com", "gitlab.com", "bitbucket.org", "codeberg.org"}

var accountHosting = regexp.MustCompile(
	`(?:` + strings.Join(quoted(codeHostingHosts), "|") + `)/([A-Za-z0-9._-]+)/` + regexp.QuoteMeta(Repository) + `\b`,
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
			if path != root && (skipped(entry.Name()) || relative == ArchivedChanges) {
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
