package surface

import (
	"errors"
	"io/fs"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/socrates-tools/socrates/internal/refusal"
	"golang.org/x/net/html"
)

var (
	namesAScheme    = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9+.\-]*:`)
	cssComment      = regexp.MustCompile(`(?s)/\*.*?\*/`)
	cssURL          = regexp.MustCompile(`url\(\s*(?:"([^"]*)"|'([^']*)'|([^)\s"']*))\s*\)`)
	cssImportString = regexp.MustCompile(`@import\s+(?:"([^"]*)"|'([^']*)')`)
)

type attachment struct {
	written string
	from    string
}

func followed(written string) (string, bool) {
	reference := strings.TrimSpace(written)
	if reference == "" || strings.HasPrefix(reference, "/") || strings.HasPrefix(reference, "#") || namesAScheme.MatchString(reference) {
		return "", false
	}
	if cut := strings.IndexAny(reference, "?#"); cut >= 0 {
		reference = reference[:cut]
	}
	if reference == "" {
		return "", false
	}
	if decoded, err := url.PathUnescape(reference); err == nil {
		reference = decoded
	}
	return reference, true
}

func referencesInCSS(contents string) []string {
	bare := cssComment.ReplaceAllString(contents, "")

	var found []string
	for _, pattern := range []*regexp.Regexp{cssURL, cssImportString} {
		for _, groups := range pattern.FindAllStringSubmatch(bare, -1) {
			for _, group := range groups[1:] {
				if group != "" {
					found = append(found, group)
				}
			}
		}
	}
	return found
}

func referencesInSourceSet(sources string) []string {
	var found []string
	for rest := sources; ; {
		rest = strings.TrimLeft(rest, " \t\r\n\f,")
		if rest == "" {
			return found
		}
		end := strings.IndexAny(rest, " \t\r\n\f")
		if end < 0 {
			end = len(rest)
		}
		candidate := rest[:end]
		rest = rest[end:]

		if trimmed := strings.TrimRight(candidate, ","); trimmed != candidate {
			found = append(found, trimmed)
			continue
		}
		found = append(found, candidate)
		if comma := strings.Index(rest, ","); comma >= 0 {
			rest = rest[comma+1:]
		} else {
			rest = ""
		}
	}
}

func referencesInHTML(contents string) []string {
	var found []string
	tokens := html.NewTokenizer(strings.NewReader(contents))
	inStyle := false

	for {
		switch tokens.Next() {
		case html.ErrorToken:
			return found
		case html.TextToken:
			if inStyle {
				found = append(found, referencesInCSS(string(tokens.Text()))...)
			}
		case html.EndTagToken:
			inStyle = false
		case html.StartTagToken, html.SelfClosingTagToken:
			name, hasAttributes := tokens.TagName()
			inStyle = string(name) == "style"
			for hasAttributes {
				var key, value []byte
				key, value, hasAttributes = tokens.TagAttr()
				switch string(key) {
				case "src", "href":
					found = append(found, string(value))
				case "srcset":
					found = append(found, referencesInSourceSet(string(value))...)
				case "style":
					found = append(found, referencesInCSS(string(value))...)
				}
			}
		}
	}
}

func isStylesheet(file string) bool {
	return strings.EqualFold(path.Ext(file), ".css")
}

func gatherFiles(decision, written, source, entry string) ([]string, error) {
	folder, err := os.OpenRoot(filepath.Dir(source))
	if err != nil {
		return nil, refusal.New(illustrationCannotBeRead, decision, written, err)
	}
	defer func() { _ = folder.Close() }()

	contents, err := os.ReadFile(source)
	if err != nil {
		return nil, refusal.New(illustrationCannotBeRead, decision, written, err)
	}
	total := int64(len(contents))
	if total > IllustrationCeiling {
		return nil, refusal.New(illustrationIsTooLarge, decision, written, total, IllustrationCeiling)
	}

	files := []string{entry}
	held := map[string]bool{entry: true}
	pending := []attachment{}
	for _, reference := range referencesInHTML(string(contents)) {
		pending = append(pending, attachment{written: reference, from: entry})
	}

	for len(pending) > 0 {
		next := pending[0]
		pending = pending[1:]

		reference, follow := followed(next.written)
		if !follow {
			continue
		}
		target := path.Join(path.Dir(next.from), reference)
		if target == ".." || strings.HasPrefix(target, "../") || path.IsAbs(target) {
			return nil, refusal.New(illustrationReferenceLeavesItsFolder, decision, next.written, entry)
		}
		if held[target] {
			continue
		}
		if _, served := mediaTypes[strings.ToLower(path.Ext(target))]; !served {
			return nil, refusal.New(illustrationReferenceIsNotServed, decision, next.written)
		}

		found, err := folder.Stat(target)
		if errors.Is(err, fs.ErrNotExist) {
			return nil, refusal.New(illustrationReferenceDoesNotExist, decision, next.written)
		}
		if err != nil {
			return nil, refusal.New(illustrationReferenceCannotBeRead, decision, next.written, err)
		}
		if !found.Mode().IsRegular() {
			return nil, refusal.New(illustrationReferenceIsNotAFile, decision, next.written)
		}
		total += found.Size()
		if total > IllustrationCeiling {
			return nil, refusal.New(illustrationFilesAreTooLarge, decision, next.written, total, IllustrationCeiling)
		}

		held[target] = true
		files = append(files, target)
		if !isStylesheet(target) {
			continue
		}
		stylesheet, err := folder.ReadFile(target)
		if err != nil {
			return nil, refusal.New(illustrationReferenceCannotBeRead, decision, next.written, err)
		}
		for _, reference := range referencesInCSS(string(stylesheet)) {
			pending = append(pending, attachment{written: reference, from: target})
		}
	}
	return files, nil
}
