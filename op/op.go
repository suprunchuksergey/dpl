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

func Eq(a, b value.Value) value.Value {
	if a.IsNull() || b.IsNull() {
		return value.Null()
	}
	if a.IsText() && b.IsText() {
		return value.Bool(a.Text() == b.Text())
	}
	return value.Bool(a.Real() == b.Real())
}

func Neq(a, b value.Value) value.Value {
	if a.IsNull() || b.IsNull() {
		return value.Null()
	}
	if a.IsText() && b.IsText() {
		return value.Bool(a.Text() != b.Text())
	}
	return value.Bool(a.Real() != b.Real())
}

func Lt(a, b value.Value) value.Value {
	if a.IsNull() || b.IsNull() {
		return value.Null()
	}
	if a.IsText() && b.IsText() {
		return value.Bool(a.Text() < b.Text())
	}
	return value.Bool(a.Real() < b.Real())
}

func Lte(a, b value.Value) value.Value {
	if a.IsNull() || b.IsNull() {
		return value.Null()
	}
	if a.IsText() && b.IsText() {
		return value.Bool(a.Text() <= b.Text())
	}
	return value.Bool(a.Real() <= b.Real())
}

func Gt(a, b value.Value) value.Value {
	if a.IsNull() || b.IsNull() {
		return value.Null()
	}
	if a.IsText() && b.IsText() {
		return value.Bool(a.Text() > b.Text())
	}
	return value.Bool(a.Real() > b.Real())
}

func Gte(a, b value.Value) value.Value {
	if a.IsNull() || b.IsNull() {
		return value.Null()
	}
	if a.IsText() && b.IsText() {
		return value.Bool(a.Text() >= b.Text())
	}
	return value.Bool(a.Real() >= b.Real())
}

func Concat(a, b value.Value) value.Value {
	if a.IsNull() || b.IsNull() {
		return value.Null()
	}
	return value.Text(a.Text() + b.Text())
}
