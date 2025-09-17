package lexer

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	_ uint8 = iota

	Add // +
	Sub // -
	Mul // *
	Div // /
	Mod // %

	Concat // ||

	Eq  // ==
	Neq // !=
	Lt  // <
	Gt  // >
	Lte // <=
	Gte // >=

	And // and
	Or  // or
	Not // not

	Ident  // ident
	Create // :=
	Set    // =

	If   // if
	Elif // elif
	Else // else

	For // for
	In  // in

	LParen // (
	RParen // )
	LBrack // [
	RBrack // ]
	LBrace // {
	RBrace // }

	Semicolon // ;
	Colon     // :
	Comma     // ,
	Dot       // .

	ArrowRight // ->
	Return     // return

	Int   // 2187
	Real  // 2.187, .2187, 2187.
	Text  // " ... "
	True  // true
	False // false
	Null  // null

	EOF // eof
)

type token struct{ id uint8 }

func (t token) String() string {
	switch t.id {
	case Add:
		return "+"
	case Sub:
		return "-"
	case Mul:
		return "*"
	case Div:
		return "/"
	case Mod:
		return "%"

	case Concat:
		return "||"

	case Eq:
		return "=="
	case Neq:
		return "!="
	case Lt:
		return "<"
	case Gt:
		return ">"
	case Lte:
		return "<="
	case Gte:
		return ">="

	case And:
		return "and"
	case Or:
		return "or"
	case Not:
		return "not"

	case Create:
		return ":="
	case Set:
		return "="

	case If:
		return "if"
	case Elif:
		return "elif"
	case Else:
		return "else"

	case For:
		return "for"
	case In:
		return "in"

	case LParen:
		return "("
	case RParen:
		return ")"
	case LBrack:
		return "["
	case RBrack:
		return "]"
	case LBrace:
		return "{"
	case RBrace:
		return "}"

	case Semicolon:
		return ";"
	case Colon:
		return ":"
	case Comma:
		return ","
	case Dot:
		return "."

	case ArrowRight:
		return "->"
	case Return:
		return "return"

	case True:
		return "true"
	case False:
		return "false"
	case Null:
		return "null"

	case EOF:
		return "eof"

	default:
		return "неизвестный"
	}
}

func (t token) ID() uint8 { return t.id }

func newToken(id uint8) token { return token{id: id} }

type tokenWithValue struct {
	token
	value string
}

func (t tokenWithValue) String() string {
	switch t.id {
	case Ident:
		return fmt.Sprintf("идентификатор %s", t.value)
	case Int:
		return fmt.Sprintf("целое число %s", t.value)
	case Real:
		return fmt.Sprintf("вещественное число %s", t.value)
	case Text:
		return fmt.Sprintf("строка %q", t.value)
	default:
		return "неизвестный"
	}
}

func (t tokenWithValue) Value() string { return t.value }

func newTokenWithValue(id uint8, value string) tokenWithValue {
	return tokenWithValue{
		token: newToken(id),
		value: value,
	}
}

type Token interface {
	fmt.Stringer
	ID() uint8
}

func NewToken(id uint8) Token { return newToken(id) }

type TokenWithValue interface {
	Token
	Value() string
}

func NewTokenWithValue(id uint8, value string) Token {
	return newTokenWithValue(id, value)
}

var keywords = map[string]uint8{
	"and":    And,
	"or":     Or,
	"not":    Not,
	"if":     If,
	"elif":   Elif,
	"else":   Else,
	"for":    For,
	"in":     In,
	"return": Return,
	"true":   True,
	"false":  False,
	"null":   Null,
}

func expected(char rune) error { return fmt.Errorf("ожидался символ %c", char) }

func unexpected(char rune) error { return fmt.Errorf("неожиданный символ %c", char) }

func helper(index int, id uint8) (int, Token) { return index + 1, newToken(id) }

func helper2(
	runes []rune,
	index int,
	char2 rune,
	id, id2 uint8,
) (int, Token) {
	if index+1 < len(runes) && runes[index+1] == char2 {
		return index + 2, newToken(id2)
	}
	return helper(index, id)
}

func readDigits(runes []rune, index int) (int, string) {
	var digits strings.Builder

	for index < len(runes) && unicode.IsDigit(runes[index]) {
		digits.WriteRune(runes[index])
		index++
	}

	str := digits.String()
	if len(str) == 0 {
		str = "0"
	}

	return index, str
}

func Tokenize(text string) ([]Token, error) {
	runes := []rune(text)
	var index int

	tokens := make([]Token, 0)
	for index < len(runes) {
		if unicode.IsSpace(runes[index]) {
			index++
			continue
		}

		var tok Token
		switch runes[index] {
		case '+':
			index, tok = helper(index, Add)
		case '-':
			index, tok = helper2(runes, index, '>', Sub, ArrowRight)
		case '*':
			index, tok = helper(index, Mul)
		case '/':
			index, tok = helper(index, Div)
		case '%':
			index, tok = helper(index, Mod)

		case '|':
			if index+1 < len(runes) && runes[index+1] == '|' {
				index, tok = index+2, newToken(Concat)
				break
			}
			return nil, expected('|')

		case '=':
			index, tok = helper2(runes, index, '=', Set, Eq)

		case '!':
			if index+1 < len(runes) && runes[index+1] == '=' {
				index, tok = index+2, newToken(Neq)
				break
			}
			return nil, expected('=')

		case '<':
			index, tok = helper2(runes, index, '=', Lt, Lte)
		case '>':
			index, tok = helper2(runes, index, '=', Gt, Gte)

		case '(':
			index, tok = helper(index, LParen)
		case ')':
			index, tok = helper(index, RParen)
		case '[':
			index, tok = helper(index, LBrack)
		case ']':
			index, tok = helper(index, RBrack)
		case '{':
			index, tok = helper(index, LBrace)
		case '}':
			index, tok = helper(index, RBrace)

		case ':':
			index, tok = helper2(runes, index, '=', Colon, Create)

		case ';':
			index, tok = helper(index, Semicolon)
		case ',':
			index, tok = helper(index, Comma)

		default:
			switch {
			case unicode.IsLetter(runes[index]) || runes[index] == '_':
				var ident strings.Builder

				for index < len(runes) &&
					(unicode.IsLetter(runes[index]) ||
						unicode.IsDigit(runes[index]) ||
						runes[index] == '_') {
					ident.WriteRune(runes[index])
					index++
				}

				kwd, ok := keywords[ident.String()]
				if ok {
					tok = newToken(kwd)
					break
				}
				tok = newTokenWithValue(Ident, ident.String())

			case runes[index] == '"':
				index++

				var str strings.Builder
				for index < len(runes) {
					if runes[index] == '"' {
						break
					}
					str.WriteRune(runes[index])
					index++
				}

				if index == len(runes) {
					return nil, expected('"')
				}

				index++
				tok = newTokenWithValue(Text, str.String())

			case runes[index] == '0':
				var value strings.Builder

				value.WriteRune('0')
				index++

				if index < len(runes) && runes[index] == '.' {
					value.WriteRune('.')

					var digits string
					index, digits = readDigits(runes, index+1)

					value.WriteString(digits)
					tok = newTokenWithValue(Real, value.String())
					break
				}

				tok = newTokenWithValue(Int, value.String())

			case unicode.IsDigit(runes[index]):
				var value strings.Builder

				var digits string
				index, digits = readDigits(runes, index)

				value.WriteString(digits)

				if index < len(runes) && runes[index] == '.' {
					value.WriteRune('.')

					var digits string
					index, digits = readDigits(runes, index+1)

					value.WriteString(digits)
					tok = newTokenWithValue(Real, value.String())
					break
				}

				tok = newTokenWithValue(Int, value.String())

			case runes[index] == '.':
				if index+1 < len(runes) && unicode.IsDigit(runes[index+1]) {
					var value strings.Builder
					value.WriteString("0.")

					var digits string
					index, digits = readDigits(runes, index+1)

					value.WriteString(digits)

					tok = newTokenWithValue(Real, value.String())
					break
				}
				index, tok = helper(index, Dot)

			default:
				return nil, unexpected(runes[index])
			}
		}
		tokens = append(tokens, tok)
	}
	tokens = append(tokens, newToken(EOF))

	return tokens, nil
}
