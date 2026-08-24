package skill

import (
	"fmt"
	"slices"
	"strings"
)

const (
	codeIndent       = "    "
	codeMarker       = "`"
	inlineCodeMarker = "|"
)

type Invocation struct {
	Command string
	Flag    string
}

func Named() []Invocation {
	var named []Invocation
	seen := map[Invocation]bool{}

	for _, code := range codeIn(Document()) {
		command := ""
		awaitingCommand := false

		for word := range strings.FieldsSeq(code) {
			word = strings.Trim(word, "`.,:;\"'()[]")

			switch {
			case word == Name:
				command, awaitingCommand = "", true
			case awaitingCommand:
				awaitingCommand = false
				if !strings.HasPrefix(word, "-") {
					command = word
					named = appended(named, seen, Invocation{Command: command})
				}
			case strings.HasPrefix(word, "--") && command != "":
				named = appended(named, seen, Invocation{Command: command, Flag: word})
			}
		}
	}

	return named
}

func codeIn(document string) []string {
	var code []string

	for line := range strings.SplitSeq(document, "\n") {
		if strings.HasPrefix(line, codeIndent) {
			code = append(code, strings.TrimSpace(line))
			continue
		}
		spans := strings.Split(line, codeMarker)
		for position := 1; position < len(spans); position += 2 {
			code = append(code, spans[position])
		}
	}

	return code
}

func Unavailable(offered map[string][]string) []string {
	var absent []string

	for _, named := range Named() {
		flags, present := offered[named.Command]
		if !present {
			absent = append(absent, fmt.Sprintf(noSuchCommand, named.Command))
			continue
		}
		if named.Flag != "" && !slices.Contains(flags, named.Flag) {
			absent = append(absent, fmt.Sprintf(noSuchFlag, named.Command, named.Flag))
		}
	}

	return absent
}

func appended(named []Invocation, seen map[Invocation]bool, invocation Invocation) []Invocation {
	if seen[invocation] {
		return named
	}
	seen[invocation] = true
	return append(named, invocation)
}
