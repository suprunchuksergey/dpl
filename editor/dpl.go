//go:build js && wasm
// +build js,wasm

package main

import (
	"github.com/suprunchuksergey/dpl"
	"github.com/suprunchuksergey/dpl/internal/value"
	"strings"
	"syscall/js"
)

func exec(_ js.Value, args []js.Value) any {
	program := args[0].String()
	output := args[1]
	draw := args[2]

	m := map[string]value.Value{
		"draw": value.Function(func(args ...value.Value) (value.Value, error) {
			draw.Invoke(
				js.ValueOf(args[0].Value()),
				js.ValueOf(args[1].Value()),
				js.ValueOf(args[2].Value()),
			)
			return value.Null(), nil
		}),

		"print": value.Function(func(args ...value.Value) (value.Value, error) {
			var str strings.Builder
			for _, arg := range args {
				str.WriteString(arg.String())
			}
			output.Invoke(js.ValueOf(str.String()))
			return value.Null(), nil
		}),

		"println": value.Function(func(args ...value.Value) (value.Value, error) {
			var str strings.Builder
			for _, arg := range args {
				str.WriteString(arg.String())
			}
			str.WriteRune('\n')
			output.Invoke(js.ValueOf(str.String()))
			return value.Null(), nil
		}),
	}

	_, err := dpl.Exec(program, m)
	if err != nil {
		output.Invoke(js.ValueOf("ошибка: " + err.Error()))
	}

	return js.Undefined()
}

func main() {
	js.Global().Set("exec", js.FuncOf(exec))
	<-make(chan struct{})
}
