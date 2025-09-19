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

	m := map[string]value.Value{
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
