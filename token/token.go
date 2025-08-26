package token

import (
	"fmt"
	"github.com/suprunchuksergey/dpl/pos"
)

const (
	_      uint8 = iota
	Add          // +
	Sub          // -
	Mul          // *
	Div          // /
	Rem          // %
	Eq           // ==
	Neq          // !=, <>
	Lt           // <
	Gt           // >
	Lte          // <=
	Gte          // >=
	Concat       // ||
	Int          // 10, 0, ...
	Real         // 10.0, .1, 10., ...
	Text         // '...'

	LParen // (
	RParen // )
	LBrack // [
	RBrack // ]
	LBrace // {
	RBrace // }

	Comma // ,
	Colon // :

	//начинается с буквы или подчеркивания и содержит
	//буквы, цифры и подчеркивания
	Ident // ident_10, ...

	//зарезервированные идентификаторы

	And   // and
	Or    // or
	Not   // not
	True  // true
	False // false
	Null  // null

	Semicolon // ;

	Assign // =

	EOF // конец файла
)

type token struct {
	id    uint8
	start pos.Pos
}

func newToken(id uint8, start pos.Pos) token {
	return token{
		id:    id,
		start: start,
	}
}

func (t token) ID() uint8 { return t.id }

func (t token) Start() pos.Pos { return t.start }

func (t token) Is(id uint8) bool { return t.id == id }

func (t token) OneOf(ids ...uint8) bool {
	for _, id := range ids {
		if t.Is(id) {
			return true
		}
	}
	return false
}

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
	case Rem:
		return "%"

	case Assign:
		return "="

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

	case Concat:
		return "||"
	case Int:
		return "int"
	case Real:
		return "real"
	case Text:
		return "text"

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
	case Comma:
		return ","
	case Colon:
		return ":"

	case Semicolon:
		return ";"

	case Ident:
		return "ident"

	case And:
		return "and"
	case Or:
		return "or"
	case Not:
		return "not"
	case True:
		return "true"
	case False:
		return "false"
	case Null:
		return "null"

	case EOF:
		return "eof"

	default:
		return "unknown"
	}
}

type withValue struct {
	token
	value string
}

func newWithValue(id uint8, start pos.Pos, value string) withValue {
	return withValue{
		token: newToken(id, start),
		value: value,
	}
}

func (t withValue) Value() string { return t.value }

type Token interface {
	ID() uint8
	Start() pos.Pos
	Is(id uint8) bool
	OneOf(ids ...uint8) bool
	fmt.Stringer
}

func New(id uint8, start pos.Pos) Token {
	return newToken(id, start)
}

type WithValue interface {
	Token
	Value() string
}

func NewWithValue(id uint8, start pos.Pos, value string) WithValue {
	return newWithValue(id, start, value)
}
