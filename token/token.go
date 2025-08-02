package token

import (
	"github.com/suprunchuksergey/dpl/pos"
)

const (
	_      uint8 = iota
	Add          // +
	Sub          // -
	Mul          // *
	Div          // /
	Rem          // %
	Eq           // =, ==
	Neq          // !=, <>
	Lt           // <
	Gt           // >
	Lte          // <=
	Gte          // >=
	Concat       // ||
	Int          // 10, 0, ...
	Real         // 10.0, .1, 10., ...
	Text         // '...'
	LParen       // (
	RParen       // )

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

func (t withValue) Value() string {
	return t.value
}

type Token interface {
	ID() uint8
	Start() pos.Pos
	Is(id uint8) bool
	OneOf(ids ...uint8) bool
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
