package skill

import (
	"fmt"
	"slices"
	"strings"
)

type Invocation struct {
	Command string
	Flag    string
}

func Named() []Invocation {
	var named []Invocation
	seen := map[Invocation]bool{}

	for line := range strings.SplitSeq(Document(), "\n") {
		command := ""
		awaitingCommand := false

		for _, word := range strings.Fields(line) {
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
