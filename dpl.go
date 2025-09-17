package dpl

import (
	"github.com/suprunchuksergey/dpl/internal/lexer"
	"github.com/suprunchuksergey/dpl/internal/namespace"
	"github.com/suprunchuksergey/dpl/internal/parser"
	"github.com/suprunchuksergey/dpl/internal/value"
)

func Exec(program string, init map[string]value.Value) (value.Value, error) {
	tokens, err := lexer.Tokenize(program)
	if err != nil {
		return nil, err
	}

	n, err := parser.Parse(tokens)
	if err != nil {
		return nil, err
	}

	return n.Exec(namespace.New(init))
}
