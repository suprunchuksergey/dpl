package value

type Value interface {
	// какой бы тип ни был у значения,
	//он должен быть конвертируемым в любой другой тип
	//Text->Real, Int->Real, ...

	Int() int64
	Real() float64
	Text() string

	//узнать, является ли значение каким-то типом

	IsInt() bool
	IsReal() bool
	IsText() bool
}

func Int(v int64) Value { return newInteger(v) }

func Real(v float64) Value { return newReal(v) }

func Text(v string) Value { return newText(v) }
