package dpl

import (
	"errors"
	"fmt"
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

	return n.Exec(initNamespace(init))
}

func builtinLen(args ...value.Value) (value.Value, error) {
	if len(args) == 0 {
		return nil, errors.New("len: требуется один аргумент")
	}
	v := args[0]
	l, err := v.Len()
	if err != nil {
		return nil, err
	}
	return value.Int(l), nil
}

func builtinAppend(args ...value.Value) (value.Value, error) {
	if len(args) == 0 {
		return nil, errors.New("append: требуется один аргумент")
	}
	v := args[0]
	return v.Append(args[1:]...)
}

func builtinPrint(args ...value.Value) (value.Value, error) {
	a := make([]any, 0, len(args))
	for _, v := range args {
		a = append(a, v.String())
	}

	fmt.Print(a...)

	return value.Null(), nil
}

func builtinPrintln(args ...value.Value) (value.Value, error) {
	a := make([]any, 0, len(args))
	for _, v := range args {
		a = append(a, v.String())
	}

	fmt.Println(a...)

	return value.Null(), nil
}

func initNamespace(init map[string]value.Value) namespace.Namespace {
	m := map[string]value.Value{
		"len":     value.Function(builtinLen),
		"append":  value.Function(builtinAppend),
		"print":   value.Function(builtinPrint),
		"println": value.Function(builtinPrintln),
	}

	for k, v := range init {
		m[k] = v
	}

	return namespace.New(m)
}
