package dpl

import (
	"github.com/suprunchuksergey/dpl/lexer"
	"github.com/suprunchuksergey/dpl/namespace"
	"github.com/suprunchuksergey/dpl/parser"
	"github.com/suprunchuksergey/dpl/val"
)

func Exec(query string, ns map[string]val.Val) (val.Val, error) {
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
