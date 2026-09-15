package stylesheet

import (
	"strings"
)

const (
	markupOpener     = "<"
	fragmentMark     = "#"
	customPrefix     = "--"
	resourceFunction = "url"
)

type Scope struct {
	Selector string
	Prefix   string
}

func Reduce(supplied string, within Scope) string {
	if strings.TrimSpace(supplied) == "" {
		return ""
	}
	return reduced(reader(supplied), within)
}

func reduced(read reading, within Scope) string {
	var written strings.Builder
	dropping, open := 0, 0

	for {
		held, stated, pieces := read.next()
		if held == ended {
			break
		}

		if dropping > 0 {
			switch held {
			case openingAtRule, openingRule:
				dropping++
			case closingAtRule, closingRule:
				dropping--
			}
			continue
		}

		switch held {
		case bodilessAtRule:
			if prelude, admitted := atRuleAdmitted(stated, pieces); admitted {
				written.WriteString(prelude)
				written.WriteString(";")
			}
		case openingAtRule:
			prelude, admitted := atRuleAdmitted(stated, pieces)
			if !admitted {
				dropping = 1
				continue
			}
			written.WriteString(prelude)
			written.WriteString("{")
			open++
		case closingAtRule, closingRule:
			if open == 0 {
				continue
			}
			written.WriteString("}")
			open--
		case openingRule:
			selectors := scoped(pieces, within)
			if len(selectors) == 0 {
				dropping = 1
				continue
			}
			written.WriteString(strings.Join(selectors, ","))
			written.WriteString("{")
			open++
		case declared:
			if open == 0 {
				continue
			}
			if emitted, admitted := declarationAdmitted(stated, pieces, within); admitted {
				written.WriteString(emitted)
			}
		case customDeclared:
			if open == 0 {
				continue
			}
			if emitted, admitted := customAdmitted(read, stated, pieces, within); admitted {
				written.WriteString(emitted)
			}
		}
	}

	if !read.whole() {
		return ""
	}
	return written.String() + strings.Repeat("}", open)
}

func atRuleAdmitted(stated string, pieces []piece) (string, bool) {
	if !allowedAtRules[strings.ToLower(unescaped(stated))] {
		return "", false
	}
	for _, held := range pieces {
		if !pieceAdmitted(held) {
			return "", false
		}
	}
	return stated + " " + joined(pieces), true
}

func declarationAdmitted(stated string, pieces []piece, within Scope) (string, bool) {
	property := strings.ToLower(unescaped(stated))
	if strings.HasPrefix(property, customPrefix) || !allowedProperties[property] {
		return "", false
	}

	value, admitted := valueAdmitted(pieces, within)
	if !admitted {
		return "", false
	}
	return property + ":" + value + ";", true
}

func customAdmitted(read reading, stated string, pieces []piece, within Scope) (string, bool) {
	property := unescaped(stated)
	if !strings.HasPrefix(property, customPrefix) || strings.Contains(property, markupOpener) {
		return "", false
	}

	var value string
	for _, held := range pieces {
		if held.sort != customValue {
			return "", false
		}
		value = held.text
	}

	lexed, whole := read.lex(value)
	if !whole {
		return "", false
	}
	rewritten, admitted := valueAdmitted(lexed, within)
	if !admitted {
		return "", false
	}
	return property + ":" + rewritten + ";", true
}

func valueAdmitted(pieces []piece, within Scope) (string, bool) {
	var written []piece

	for _, held := range pieces {
		if !pieceAdmitted(held) {
			return "", false
		}
		if held.sort == resource {
			reference, admitted := fragmentReference(held.text, within)
			if !admitted {
				return "", false
			}
			written = append(written, piece{sort: resource, text: reference})
			continue
		}
		written = append(written, held)
	}

	value := strings.TrimSpace(joined(written))
	if value == "" {
		return "", false
	}
	return value, true
}

func pieceAdmitted(held piece) bool {
	if held.sort == broken || strings.Contains(held.text, markupOpener) {
		return false
	}
	if held.sort != call {
		return true
	}
	return allowedFunctions[strings.ToLower(unescaped(strings.TrimSuffix(held.text, "(")))]
}

func fragmentReference(stated string, within Scope) (string, bool) {
	read := unescaped(stated)

	opened := strings.Index(read, "(")
	if opened < 0 || !strings.EqualFold(read[:opened], resourceFunction) || !strings.HasSuffix(read, ")") {
		return "", false
	}

	inside := strings.TrimSpace(read[opened+1 : len(read)-1])
	for _, quote := range []string{`"`, "'"} {
		if len(inside) >= 2 && strings.HasPrefix(inside, quote) && strings.HasSuffix(inside, quote) {
			inside = inside[1 : len(inside)-1]
			break
		}
	}

	referenced, fragmentary := strings.CutPrefix(inside, fragmentMark)
	if !fragmentary || referenced == "" || strings.ContainsAny(referenced, `"'() `) {
		return "", false
	}
	return resourceFunction + "(" + fragmentMark + within.Prefix + referenced + ")", true
}

func joined(pieces []piece) string {
	var written strings.Builder
	for _, held := range pieces {
		written.WriteString(held.text)
	}
	return written.String()
}

func unescaped(stated string) string {
	if !strings.Contains(stated, `\`) {
		return stated
	}

	var written strings.Builder
	for position := 0; position < len(stated); position++ {
		if stated[position] != '\\' || position+1 == len(stated) {
			written.WriteByte(stated[position])
			continue
		}

		position++
		digits := 0
		value := rune(0)
		for digits < 6 && position+digits < len(stated) && hexadecimal(stated[position+digits]) {
			value = value*16 + rune(figure(stated[position+digits]))
			digits++
		}
		if digits == 0 {
			written.WriteByte(stated[position])
			continue
		}

		written.WriteRune(value)
		position += digits
		if position < len(stated) && (stated[position] == ' ' || stated[position] == '\t' || stated[position] == '\n') {
			continue
		}
		position--
	}
	return written.String()
}

func hexadecimal(character byte) bool {
	return character >= '0' && character <= '9' ||
		character >= 'a' && character <= 'f' ||
		character >= 'A' && character <= 'F'
}

func figure(character byte) int {
	switch {
	case character <= '9':
		return int(character - '0')
	case character >= 'a':
		return int(character-'a') + 10
	default:
		return int(character-'A') + 10
	}
}
