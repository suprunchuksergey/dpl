package value

import "fmt"

type Value interface {
	// какой бы тип ни был у значения,
	//он должен быть конвертируемым в любой другой тип
	//Text->Real, Int->Real, ...

	Int() int64
	Real() float64
	Text() string
	Bool() bool

	//узнать, является ли значение каким-то типом

	IsInt() bool
	IsReal() bool
	IsText() bool
	IsBool() bool
	IsNull() bool

	fmt.Stringer
}

func Int(v int64) Value { return newInteger(v) }

func Real(v float64) Value { return newReal(v) }

func Text(v string) Value { return newText(v) }

func Null() Value { return newNull() }

func True() Value { return newTrue() }

func False() Value { return newFalse() }

func Bool(v bool) Value {
	if v {
		return newTrue()
	}
	return newFalse()
}
