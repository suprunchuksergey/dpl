package op

import (
	"github.com/suprunchuksergey/dpl/value"
	"math"
)

func Add(a, b value.Value) value.Value {
	if a.IsNull() || b.IsNull() {
		return value.Null()
	}
	if a.IsReal() || b.IsReal() {
		return value.Real(a.Real() + b.Real())
	}
	return value.Int(a.Int() + b.Int())
}

func Sub(a, b value.Value) value.Value {
	if a.IsNull() || b.IsNull() {
		return value.Null()
	}
	if a.IsReal() || b.IsReal() {
		return value.Real(a.Real() - b.Real())
	}
	return value.Int(a.Int() - b.Int())
}

func Mul(a, b value.Value) value.Value {
	if a.IsNull() || b.IsNull() {
		return value.Null()
	}
	if a.IsReal() || b.IsReal() {
		return value.Real(a.Real() * b.Real())
	}
	return value.Int(a.Int() * b.Int())
}

func Div(a, b value.Value) value.Value {
	if a.IsNull() || b.IsNull() || b.Real() == 0 {
		return value.Null()
	}
	if a.IsReal() || b.IsReal() {
		return value.Real(a.Real() / b.Real())
	}
	return value.Int(a.Int() / b.Int())
}

func Rem(a, b value.Value) value.Value {
	if a.IsNull() || b.IsNull() || b.Real() == 0 {
		return value.Null()
	}
	if a.IsReal() || b.IsReal() {
		return value.Real(math.Mod(a.Real(), b.Real()))
	}
	return value.Int(a.Int() % b.Int())
}
