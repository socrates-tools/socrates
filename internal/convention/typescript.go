package convention

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	stringToken      = "string"
	wordToken        = "word"
	punctuationToken = "punctuation"

	substitutionMark = "${}"
)

type typeScriptToken struct {
	Kind  string
	Text  string
	Start int
	End   int
	Line  int
}

type typeScriptLexer struct {
	source string
	at     int
	line   int
	tokens []typeScriptToken
}

func lexTypeScript(source string) ([]typeScriptToken, error) {
	lexer := &typeScriptLexer{source: source, line: 1}
	if err := lexer.lex(false); err != nil {
		return nil, err
	}
	return lexer.tokens, nil
}

func (lexer *typeScriptLexer) lex(untilBrace bool) error {
	depth := 0
	for lexer.at < len(lexer.source) {
		character := lexer.source[lexer.at]
		switch {
		case character == '\n':
			lexer.line++
			lexer.at++
		case character == ' ' || character == '\t' || character == '\r':
			lexer.at++
		case strings.HasPrefix(lexer.source[lexer.at:], "//"):
			lexer.skipLineComment()
		case strings.HasPrefix(lexer.source[lexer.at:], "/*"):
			if err := lexer.skipBlockComment(); err != nil {
				return err
			}
		case character == '\'' || character == '"':
			if err := lexer.quoted(character); err != nil {
				return err
			}
		case character == '`':
			if err := lexer.template(); err != nil {
				return err
			}
		case isWordCharacter(character):
			lexer.word()
		default:
			if untilBrace && character == '}' {
				if depth == 0 {
					lexer.at++
					return nil
				}
				depth--
			}
			if untilBrace && character == '{' {
				depth++
			}
			lexer.punctuation()
		}
	}
	if untilBrace {
		return fmt.Errorf("line %d: a template substitution is never closed", lexer.line)
	}
	return nil
}

func (lexer *typeScriptLexer) skipLineComment() {
	for lexer.at < len(lexer.source) && lexer.source[lexer.at] != '\n' {
		lexer.at++
	}
}

func (lexer *typeScriptLexer) skipBlockComment() error {
	opened := lexer.line
	end := strings.Index(lexer.source[lexer.at+2:], "*/")
	if end < 0 {
		return fmt.Errorf("line %d: a comment is never closed", opened)
	}
	lexer.line += strings.Count(lexer.source[lexer.at:lexer.at+2+end], "\n")
	lexer.at += end + 4
	return nil
}

func (lexer *typeScriptLexer) word() {
	start := lexer.at
	for lexer.at < len(lexer.source) && isWordCharacter(lexer.source[lexer.at]) {
		lexer.at++
	}
	lexer.tokens = append(lexer.tokens, typeScriptToken{Kind: wordToken, Text: lexer.source[start:lexer.at], Start: start, End: lexer.at, Line: lexer.line})
}

func (lexer *typeScriptLexer) punctuation() {
	_, width := utf8.DecodeRuneInString(lexer.source[lexer.at:])
	lexer.tokens = append(lexer.tokens, typeScriptToken{Kind: punctuationToken, Text: lexer.source[lexer.at : lexer.at+width], Start: lexer.at, End: lexer.at + width, Line: lexer.line})
	lexer.at += width
}

func (lexer *typeScriptLexer) quoted(quote byte) error {
	start, opened := lexer.at, lexer.line
	lexer.at++
	var text strings.Builder
	for lexer.at < len(lexer.source) {
		character := lexer.source[lexer.at]
		switch {
		case character == quote:
			lexer.at++
			lexer.tokens = append(lexer.tokens, typeScriptToken{Kind: stringToken, Text: text.String(), Start: start, End: lexer.at, Line: opened})
			return nil
		case character == '\n':
			return fmt.Errorf("line %d: a string is never closed", opened)
		case character == '\\':
			if err := lexer.escape(&text); err != nil {
				return err
			}
		default:
			text.WriteByte(character)
			lexer.at++
		}
	}
	return fmt.Errorf("line %d: a string is never closed", opened)
}

func (lexer *typeScriptLexer) template() error {
	start, opened := lexer.at, lexer.line
	lexer.at++
	var text strings.Builder
	for lexer.at < len(lexer.source) {
		character := lexer.source[lexer.at]
		switch {
		case character == '`':
			lexer.at++
			lexer.tokens = append(lexer.tokens, typeScriptToken{Kind: stringToken, Text: text.String(), Start: start, End: lexer.at, Line: opened})
			return nil
		case character == '\\':
			if err := lexer.escape(&text); err != nil {
				return err
			}
		case character == '$' && strings.HasPrefix(lexer.source[lexer.at:], "${"):
			lexer.at += 2
			text.WriteString(substitutionMark)
			if err := lexer.lex(true); err != nil {
				return err
			}
		default:
			if character == '\n' {
				lexer.line++
			}
			text.WriteByte(character)
			lexer.at++
		}
	}
	return fmt.Errorf("line %d: a template is never closed", opened)
}

func (lexer *typeScriptLexer) escape(text *strings.Builder) error {
	lexer.at++
	if lexer.at >= len(lexer.source) {
		return fmt.Errorf("line %d: a string ends in a backslash", lexer.line)
	}
	character := lexer.source[lexer.at]
	lexer.at++
	switch character {
	case 'n':
		text.WriteByte('\n')
	case 't':
		text.WriteByte('\t')
	case 'r':
		text.WriteByte('\r')
	case 'b':
		text.WriteByte('\b')
	case 'f':
		text.WriteByte('\f')
	case 'v':
		text.WriteByte('\v')
	case '0':
		text.WriteByte(0)
	case '\n':
		lexer.line++
	case 'x':
		return lexer.codePoint(text, 2)
	case 'u':
		if strings.HasPrefix(lexer.source[lexer.at:], "{") {
			end := strings.IndexByte(lexer.source[lexer.at:], '}')
			if end < 0 {
				return fmt.Errorf("line %d: an escape is never closed", lexer.line)
			}
			digits := lexer.source[lexer.at+1 : lexer.at+end]
			lexer.at += end + 1
			return writeCodePoint(text, digits, lexer.line)
		}
		return lexer.codePoint(text, 4)
	default:
		text.WriteByte(character)
	}
	return nil
}

func (lexer *typeScriptLexer) codePoint(text *strings.Builder, width int) error {
	if lexer.at+width > len(lexer.source) {
		return fmt.Errorf("line %d: an escape is cut short", lexer.line)
	}
	digits := lexer.source[lexer.at : lexer.at+width]
	lexer.at += width
	return writeCodePoint(text, digits, lexer.line)
}

func writeCodePoint(text *strings.Builder, digits string, line int) error {
	value, err := strconv.ParseUint(digits, 16, 32)
	if err != nil {
		return fmt.Errorf("line %d: %q is not a code point", line, digits)
	}
	text.WriteRune(rune(value))
	return nil
}

func isWordCharacter(character byte) bool {
	return character == '_' || character == '$' ||
		(character >= 'a' && character <= 'z') ||
		(character >= 'A' && character <= 'Z') ||
		(character >= '0' && character <= '9')
}

func stringsOf(tokens []typeScriptToken) []typeScriptToken {
	var found []typeScriptToken
	for _, token := range tokens {
		if token.Kind == stringToken {
			found = append(found, token)
		}
	}
	return found
}

func isPropertyKey(tokens []typeScriptToken, position int) bool {
	next := position + 1
	return next < len(tokens) && tokens[next].Kind == punctuationToken && tokens[next].Text == ":"
}

func declaredGroup(tokens []typeScriptToken, name string) ([]typeScriptToken, bool) {
	for position := 2; position < len(tokens); position++ {
		if tokens[position].Kind != wordToken || tokens[position].Text != name {
			continue
		}
		if tokens[position-1].Text != "const" || tokens[position-2].Text != "export" {
			continue
		}

		opened := -1
		for ahead := position + 1; ahead < len(tokens); ahead++ {
			if tokens[ahead].Kind == punctuationToken && tokens[ahead].Text == "=" {
				opened = ahead + 1
				break
			}
		}
		if opened < 0 || opened >= len(tokens) || tokens[opened].Text != "{" {
			return nil, false
		}

		depth := 0
		for ahead := opened; ahead < len(tokens); ahead++ {
			if tokens[ahead].Kind == punctuationToken {
				switch tokens[ahead].Text {
				case "{":
					depth++
				case "}":
					depth--
				}
			}
			if depth == 0 {
				return tokens[opened : ahead+1], true
			}
		}
		return nil, false
	}
	return nil, false
}
