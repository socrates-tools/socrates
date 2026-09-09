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
	Binary  string
	Command string
	Flag    string
}

func Binaries() []string {
	return []string{Name, Surface}
}

func Named() []Invocation {
	var named []Invocation
	seen := map[Invocation]bool{}

	var written []string
	for _, one := range All() {
		written = append(written, codeIn(one.Document())...)
	}

	for _, code := range written {
		binary, command := "", ""
		awaitingCommand := false

		for word := range strings.FieldsSeq(code) {
			word = strings.Trim(word, "`.,:;\"'()[]")

			switch {
			case slices.Contains(Binaries(), word):
				binary, command, awaitingCommand = word, "", true
			case awaitingCommand:
				awaitingCommand = false
				if !strings.HasPrefix(word, "-") {
					command = word
					named = appended(named, seen, Invocation{Binary: binary, Command: command})
				}
			case strings.HasPrefix(word, "--") && command != "":
				named = appended(named, seen, Invocation{Binary: binary, Command: command, Flag: word})
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

func Unavailable(offered map[string]map[string][]string) []string {
	var absent []string

	for _, named := range Named() {
		commands, built := offered[named.Binary]
		if !built {
			absent = append(absent, fmt.Sprintf(noSuchBinary, named.Binary, named.Command, named.Binary))
			continue
		}
		flags, present := commands[named.Command]
		if !present {
			absent = append(absent, fmt.Sprintf(noSuchCommand, named.Binary, named.Command, named.Binary))
			continue
		}
		if named.Flag != "" && !slices.Contains(flags, named.Flag) {
			absent = append(absent, fmt.Sprintf(noSuchFlag, named.Binary, named.Command, named.Flag))
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
