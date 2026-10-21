package convention

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

var SurfaceStyleDirectory = filepath.Join(ToolDirectory("agora"), "web", "packages", "surface")

var AgoraGlossary = filepath.Join(ToolDirectory("agora"), GlossaryFile)

var PiecesOfTheSurface = map[string][]string{
	"decision":       {"decision", "decisions", "decisions-item", "attribution", "name", "question", "meaning", "answering", "answered", "words", "check", "still-open"},
	"source":         {"location"},
	"index":          {"index", "entries", "entry", "progress"},
	"settled answer": {"settled", "when-open", "when-closed"},
	"option":         {"options"},
	"evidence":       {"evidence", "evidences", "diff", "illustrated"},
	"illustration":   {"drawing"},
	"confirmation":   {"confirmation"},
	"composer":       {"composer", "bar", "send", "sends", "finish", "amend", "state", "message"},
	"receipt":        {"receipt", "notice"},
	"summary":        {"summary", "panel"},
	"content":        {"main-inner", "heading", "meta", "preamble", "over-notice"},
	"topbar":         {"brand", "crumb", "wordmark", "status"},
	"empty page":     {"nothing"},
}

var primitiveElements = map[string]bool{
	"button":   true,
	"input":    true,
	"textarea": true,
	"select":   true,
	"dialog":   true,
}

var glossaryTerm = regexp.MustCompile(`(?m)^- \*\*([^*]+)\*\*:`)

var attributeSelector = regexp.MustCompile(`\[[^\]]*\]`)

var pseudoName = regexp.MustCompile(`::?[a-zA-Z-]+`)

var classSelector = regexp.MustCompile(`\.(-?[_a-zA-Z][\w-]*)`)

var idSelector = regexp.MustCompile(`#[\w-]+`)

var elementSelector = regexp.MustCompile(`[a-zA-Z][\w-]*`)

var atRulesHoldingRules = map[string]bool{"@media": true, "@supports": true, "@container": true, "@layer": true}

type selectorWritten struct {
	text string
	line int
}

func PiecesStyledBySurface(root, directory string, pieces map[string][]string) ([]string, error) {
	var found []string

	err := filepath.WalkDir(filepath.Join(root, directory), func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filepath.Ext(path) != ".css" {
			return nil
		}
		held, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		relative, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return relErr
		}
		found = append(found, stylesOfWhatNoPieceOwnsIn(relative, string(held), pieces)...)
		return nil
	})
	if err != nil {
		return nil, err
	}

	named, err := glossaryTermsIn(filepath.Join(root, AgoraGlossary))
	if err != nil {
		return nil, err
	}
	for _, piece := range sortedKeys(pieces) {
		if !slices.Contains(named, piece) {
			found = append(found, fmt.Sprintf("%s does not name the piece %q", filepath.ToSlash(AgoraGlossary), piece))
		}
	}

	return found, nil
}

func sortedKeys(held map[string][]string) []string {
	keys := make([]string, 0, len(held))
	for key := range held {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}

func glossaryTermsIn(path string) ([]string, error) {
	held, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var terms []string
	for _, match := range glossaryTerm.FindAllStringSubmatch(string(held), -1) {
		terms = append(terms, match[1])
	}
	return terms, nil
}

func stylesOfWhatNoPieceOwnsIn(named, held string, pieces map[string][]string) []string {
	var found []string

	for _, written := range selectorsWrittenIn(held) {
		if offender := whatNoPieceOwns(written.text, pieces); offender != "" {
			found = append(found, fmt.Sprintf("%s:%d styles %q, and %q belongs to no piece the glossary names", named, written.line, written.text, offender))
		}
	}

	return found
}

func whatNoPieceOwns(selector string, pieces map[string][]string) string {
	plain := idSelector.ReplaceAllString(pseudoName.ReplaceAllString(attributeSelector.ReplaceAllString(selector, ""), " "), " ")

	type written struct {
		at   int
		what string
	}
	var offences []written

	for _, class := range classSelector.FindAllStringSubmatchIndex(plain, -1) {
		name := plain[class[2]:class[3]]
		if !ownedByAPiece(name, pieces) {
			offences = append(offences, written{at: class[0], what: "." + name})
		}
	}
	withoutClasses := classSelector.ReplaceAllStringFunc(plain, func(class string) string { return strings.Repeat(" ", len(class)) })
	for _, element := range elementSelector.FindAllStringIndex(withoutClasses, -1) {
		name := withoutClasses[element[0]:element[1]]
		if primitiveElements[name] || strings.HasPrefix(name, "ui-") {
			offences = append(offences, written{at: element[0], what: name})
		}
	}

	if len(offences) == 0 {
		return ""
	}
	first := offences[0]
	for _, offence := range offences[1:] {
		if offence.at < first.at {
			first = offence
		}
	}
	return first.what
}

func ownedByAPiece(class string, pieces map[string][]string) bool {
	for _, classes := range pieces {
		for _, owned := range classes {
			if class == owned || strings.HasPrefix(class, owned+"-") {
				return true
			}
		}
	}
	return false
}

func selectorsWrittenIn(held string) []selectorWritten {
	unannotated := commentWritten.ReplaceAllStringFunc(held, func(comment string) string {
		return notALineBreak.ReplaceAllString(comment, " ")
	})

	var written []selectorWritten
	start := 0
	line := 1
	startLine := 1
	for position := 0; position < len(unannotated); position++ {
		switch unannotated[position] {
		case '\n':
			line++
		case ';', '}':
			start = position + 1
			startLine = line
		case '{':
			prelude := unannotated[start:position]
			trimmed := strings.TrimSpace(prelude)
			next := position
			if strings.HasPrefix(trimmed, "@") {
				name, _, _ := strings.Cut(trimmed, " ")
				if !atRulesHoldingRules[name] {
					next = endOfBlock(unannotated, position)
					line += strings.Count(unannotated[position:next], "\n")
				}
			} else {
				written = append(written, selectorsOf(prelude, startLine)...)
				next = endOfBlock(unannotated, position)
				line += strings.Count(unannotated[position:next], "\n")
			}
			position = next
			start = position + 1
			startLine = line
		}
	}
	return written
}

func endOfBlock(held string, open int) int {
	depth := 0
	for position := open; position < len(held); position++ {
		switch held[position] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return position
			}
		}
	}
	return len(held) - 1
}

func selectorsOf(prelude string, firstLine int) []selectorWritten {
	var selectors []selectorWritten

	depth := 0
	partStart := 0
	flush := func(end int) {
		part := prelude[partStart:end]
		text := strings.TrimSpace(part)
		if text == "" {
			return
		}
		leading := len(part) - len(strings.TrimLeft(part, " \t\r\n"))
		line := firstLine + strings.Count(prelude[:partStart+leading], "\n")
		selectors = append(selectors, selectorWritten{text: strings.Join(strings.Fields(text), " "), line: line})
	}
	for position := 0; position < len(prelude); position++ {
		switch prelude[position] {
		case '(', '[':
			depth++
		case ')', ']':
			depth--
		case ',':
			if depth == 0 {
				flush(position)
				partStart = position + 1
			}
		}
	}
	flush(len(prelude))
	return selectors
}
