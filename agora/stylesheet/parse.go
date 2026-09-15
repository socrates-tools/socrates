package stylesheet

import (
	"errors"
	"io"

	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/css"
)

type shape int

const (
	ended shape = iota
	passedOver
	bodilessAtRule
	openingAtRule
	closingAtRule
	openingRule
	closingRule
	declared
	customDeclared
)

type sort int

const (
	loose sort = iota
	word
	call
	resource
	name
	quoted
	space
	mark
	colon
	comma
	opening
	closing
	openingSquare
	closingSquare
	customName
	customValue
	broken
)

type piece struct {
	sort sort
	text string
}

type reading interface {
	next() (shape, string, []piece)
	whole() bool
	lex(value string) ([]piece, bool)
}

type walk struct {
	parser *css.Parser
	err    error
}

func reader(supplied string) *walk {
	return &walk{parser: css.NewParser(parse.NewInputString(supplied), false)}
}

func (w *walk) next() (shape, string, []piece) {
	grammar, _, data := w.parser.Next()
	if grammar == css.ErrorGrammar {
		w.err = w.parser.Err()
		return ended, "", nil
	}
	return shapeOf(grammar), string(data), piecesOf(w.parser.Values())
}

func (w *walk) whole() bool {
	return errors.Is(w.err, io.EOF)
}

func (w *walk) lex(value string) ([]piece, bool) {
	lexer := css.NewLexer(parse.NewInputString(value))

	var held []piece
	for {
		token, data := lexer.Next()
		if token == css.ErrorToken {
			return held, errors.Is(lexer.Err(), io.EOF)
		}
		held = append(held, pieceOf(token, data))
	}
}

func shapeOf(grammar css.GrammarType) shape {
	switch grammar {
	case css.AtRuleGrammar:
		return bodilessAtRule
	case css.BeginAtRuleGrammar:
		return openingAtRule
	case css.EndAtRuleGrammar:
		return closingAtRule
	case css.BeginRulesetGrammar:
		return openingRule
	case css.EndRulesetGrammar:
		return closingRule
	case css.DeclarationGrammar:
		return declared
	case css.CustomPropertyGrammar:
		return customDeclared
	default:
		return passedOver
	}
}

func piecesOf(tokens []css.Token) []piece {
	held := make([]piece, 0, len(tokens))
	for _, token := range tokens {
		held = append(held, pieceOf(token.TokenType, token.Data))
	}
	return held
}

func pieceOf(token css.TokenType, data []byte) piece {
	return piece{sort: sortOf(token), text: string(data)}
}

func sortOf(token css.TokenType) sort {
	switch token {
	case css.IdentToken:
		return word
	case css.FunctionToken:
		return call
	case css.URLToken:
		return resource
	case css.HashToken:
		return name
	case css.StringToken:
		return quoted
	case css.WhitespaceToken:
		return space
	case css.DelimToken:
		return mark
	case css.ColonToken:
		return colon
	case css.CommaToken:
		return comma
	case css.LeftParenthesisToken:
		return opening
	case css.RightParenthesisToken:
		return closing
	case css.LeftBracketToken:
		return openingSquare
	case css.RightBracketToken:
		return closingSquare
	case css.CustomPropertyNameToken:
		return customName
	case css.CustomPropertyValueToken:
		return customValue
	case css.BadStringToken, css.BadURLToken, css.ErrorToken:
		return broken
	default:
		return loose
	}
}
