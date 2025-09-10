package dpl

import (
	"github.com/suprunchuksergey/dpl/internal/namespace"
	"github.com/suprunchuksergey/dpl/internal/value"
	"github.com/suprunchuksergey/dpl/lexer"
	"github.com/suprunchuksergey/dpl/parser"
)

func Exec(query string, ns map[string]value.Value) (value.Value, error) {
	lex := lexer.New(query)

	err := lex.Next()
	if err != nil {
		return nil, err
	}

	p := parser.New(lex)

	n, err := p.Parse()
	if err != nil {
		return nil, err
	}

	return n.Exec(namespace.New(ns))
}
