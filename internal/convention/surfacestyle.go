package convention

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var ClientStyleDirectory = filepath.Join(ToolDirectory("agora"), "web", "packages")

var attributeSelector = regexp.MustCompile(`\[[^\]]*\]`)

var pseudoName = regexp.MustCompile(`::?[a-zA-Z-]+`)

var namedInASelector = regexp.MustCompile(`[#.]?[a-zA-Z_-][\w-]*`)

var atRulesHoldingRules = map[string]bool{"@media": true, "@supports": true, "@container": true, "@layer": true}

type selectorWritten struct {
	text string
	line int
}

func RulesAUtilityCouldSay(root, directory string) ([]string, error) {
	var found []string

	err := filepath.WalkDir(filepath.Join(root, directory), func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() && entry.Name() == "node_modules" {
			return filepath.SkipDir
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
		for _, written := range selectorsWrittenIn(string(held)) {
			if namesWhatAUtilityStyles(written.text) {
				found = append(found, fmt.Sprintf("%s:%d styles %q, and a utility on it says that", filepath.ToSlash(relative), written.line, written.text))
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return found, nil
}

func namesWhatAUtilityStyles(selector string) bool {
	bare := pseudoName.ReplaceAllString(attributeSelector.ReplaceAllString(selector, ""), " ")
	return namedInASelector.MatchString(bare)
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
