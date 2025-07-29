package op

import (
	"github.com/suprunchuksergey/dpl/val"
	"math"
)

type Unary func(v val.Val) val.Val

func unary(f Unary) Unary {
	return func(v val.Val) val.Val {
		if val.IsNull(v) {
			return v
		}
		return f(v)
	}
}

type Binary func(a, b val.Val) val.Val

func binary(f Binary) Binary {
	return func(a, b val.Val) val.Val {
		if val.IsNull(a) || val.IsNull(b) {
			return val.Null()
		}
		return f(a, b)
	}
}

var Add = binary(func(a, b val.Val) val.Val {
	if val.IsReal(a) || val.IsReal(b) {
		return val.Real(a.Real() + b.Real())
	}
	return val.Int(a.Int() + b.Int())
})

var Sub = binary(func(a, b val.Val) val.Val {
	if val.IsReal(a) || val.IsReal(b) {
		return val.Real(a.Real() - b.Real())
	}
	return val.Int(a.Int() - b.Int())
})

var Mul = binary(func(a, b val.Val) val.Val {
	if val.IsReal(a) || val.IsReal(b) {
		return val.Real(a.Real() * b.Real())
	}
	return val.Int(a.Int() * b.Int())
})

var Div = binary(func(a, b val.Val) val.Val {
	if b.Real() == 0 {
		return val.Null()
	}
	if val.IsReal(a) || val.IsReal(b) {
		return val.Real(a.Real() / b.Real())
	}
	return val.Int(a.Int() / b.Int())
})

var Rem = binary(func(a, b val.Val) val.Val {
	if b.Real() == 0 {
		return val.Null()
	}
	if val.IsReal(a) || val.IsReal(b) {
		return val.Real(math.Mod(a.Real(), b.Real()))
	}
	return val.Int(a.Int() % b.Int())
})

var Eq = binary(func(a, b val.Val) val.Val {
	if val.IsText(a) && val.IsText(b) {
		return val.Bool(a.Text() == b.Text())
	}
	return val.Bool(a.Real() == b.Real())
})

var Neq = binary(func(a, b val.Val) val.Val {
	if val.IsText(a) && val.IsText(b) {
		return val.Bool(a.Text() != b.Text())
	}
	return val.Bool(a.Real() != b.Real())
})

var Lt = binary(func(a, b val.Val) val.Val {
	if val.IsText(a) && val.IsText(b) {
		return val.Bool(a.Text() < b.Text())
	}
	return val.Bool(a.Real() < b.Real())
})

var Lte = binary(func(a, b val.Val) val.Val {
	if val.IsText(a) && val.IsText(b) {
		return val.Bool(a.Text() <= b.Text())
	}
	return val.Bool(a.Real() <= b.Real())
})

var Gt = binary(func(a, b val.Val) val.Val {
	if val.IsText(a) && val.IsText(b) {
		return val.Bool(a.Text() > b.Text())
	}
	return val.Bool(a.Real() > b.Real())
})

var Gte = binary(func(a, b val.Val) val.Val {
	if val.IsText(a) && val.IsText(b) {
		return val.Bool(a.Text() >= b.Text())
	}
	return val.Bool(a.Real() >= b.Real())
})

var Concat = binary(func(a, b val.Val) val.Val {
	return val.Text(a.Text() + b.Text())
})

var And = binary(func(a, b val.Val) val.Val {
	return val.Bool(a.Bool() && b.Bool())
})

var Or = binary(func(a, b val.Val) val.Val {
	return val.Bool(a.Bool() || b.Bool())
})

var Not = unary(func(v val.Val) val.Val {
	return val.Bool(!v.Bool())
})
