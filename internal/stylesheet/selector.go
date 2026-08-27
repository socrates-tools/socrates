package stylesheet

import "strings"

var combinators = named(">", "+", "~")

const classMark = "."

func scoped(pieces []piece, within Scope) []string {
	var selectors []string
	for _, group := range commaSeparated(pieces) {
		if written, admitted := scopedSelector(group, within); admitted {
			selectors = append(selectors, written)
		}
	}
	return selectors
}

func commaSeparated(pieces []piece) [][]piece {
	var groups [][]piece
	var held []piece

	depth := 0
	for _, one := range pieces {
		switch one.sort {
		case opening, openingSquare:
			depth++
		case closing, closingSquare:
			depth--
		case comma:
			if depth == 0 {
				groups = append(groups, held)
				held = nil
				continue
			}
		}
		held = append(held, one)
	}
	return append(groups, held)
}

func scopedSelector(pieces []piece, within Scope) (string, bool) {
	trimmed := withoutEdgeSpace(pieces)
	if len(trimmed) == 0 {
		return "", false
	}

	for _, held := range trimmed {
		if held.sort == broken || held.sort == resource || strings.Contains(held.text, markupOpener) {
			return "", false
		}
	}
	if namesARefusedPseudo(trimmed) {
		return "", false
	}

	compound := leadingCompound(trimmed)
	if namesTheCallersRoot(trimmed[:compound]) {
		return within.Selector + renamed(trimmed[compound:], within.Prefix), true
	}
	return within.Selector + " " + renamed(trimmed, within.Prefix), true
}

func withoutEdgeSpace(pieces []piece) []piece {
	for len(pieces) > 0 && pieces[0].sort == space {
		pieces = pieces[1:]
	}
	for len(pieces) > 0 && pieces[len(pieces)-1].sort == space {
		pieces = pieces[:len(pieces)-1]
	}
	return pieces
}

func namesARefusedPseudo(pieces []piece) bool {
	following := false
	for _, held := range pieces {
		switch held.sort {
		case colon:
			following = true
		case word, call:
			if following && refusedPseudo[strings.ToLower(unescaped(strings.TrimSuffix(held.text, "(")))] {
				return true
			}
			following = false
		default:
			following = false
		}
	}
	return false
}

func leadingCompound(pieces []piece) int {
	depth := 0
	for position, held := range pieces {
		switch held.sort {
		case opening, openingSquare:
			depth++
		case closing, closingSquare:
			depth--
		case space:
			if depth == 0 {
				return position
			}
		case mark:
			if depth == 0 && combinators[held.text] {
				return position
			}
		}
	}
	return len(pieces)
}

func namesTheCallersRoot(compound []piece) bool {
	switch len(compound) {
	case 1:
		return compound[0].sort == word && theCallersRoot[strings.ToLower(unescaped(compound[0].text))]
	case 2:
		return compound[0].sort == colon && compound[1].sort == word &&
			strings.ToLower(unescaped(compound[1].text)) == rootPseudo
	default:
		return false
	}
}

func renamed(pieces []piece, prefix string) string {
	var written strings.Builder

	for position := 0; position < len(pieces); position++ {
		held := pieces[position]
		switch {
		case held.sort == name:
			written.WriteString(fragmentMark)
			written.WriteString(prefix)
			written.WriteString(strings.TrimPrefix(held.text, fragmentMark))
		case held.sort == mark && held.text == classMark && position+1 < len(pieces) && pieces[position+1].sort == word:
			written.WriteString(classMark)
			written.WriteString(prefix)
			written.WriteString(pieces[position+1].text)
			position++
		default:
			written.WriteString(held.text)
		}
	}
	return written.String()
}
