package op

import (
	"github.com/suprunchuksergey/dpl/val"
	"math"
)

type Op func(a, b val.Val) val.Val

func check(f Op) Op {
	return func(a, b val.Val) val.Val {
		if val.IsNull(a) || val.IsNull(b) {
			return val.Null()
		}
		return f(a, b)
	}
}

var Add = check(func(a, b val.Val) val.Val {
	if val.IsReal(a) || val.IsReal(b) {
		return val.Real(a.Real() + b.Real())
	}
	return val.Int(a.Int() + b.Int())
})

var Sub = check(func(a, b val.Val) val.Val {
	if val.IsReal(a) || val.IsReal(b) {
		return val.Real(a.Real() - b.Real())
	}
	return val.Int(a.Int() - b.Int())
})

var Mul = check(func(a, b val.Val) val.Val {
	if val.IsReal(a) || val.IsReal(b) {
		return val.Real(a.Real() * b.Real())
	}
	return val.Int(a.Int() * b.Int())
})

var Div = check(func(a, b val.Val) val.Val {
	if b.Real() == 0 {
		return val.Null()
	}
	if val.IsReal(a) || val.IsReal(b) {
		return val.Real(a.Real() / b.Real())
	}
	return val.Int(a.Int() / b.Int())
})

var Rem = check(func(a, b val.Val) val.Val {
	if b.Real() == 0 {
		return val.Null()
	}
	if val.IsReal(a) || val.IsReal(b) {
		return val.Real(math.Mod(a.Real(), b.Real()))
	}
	return val.Int(a.Int() % b.Int())
})

var Eq = check(func(a, b val.Val) val.Val {
	if val.IsText(a) && val.IsText(b) {
		return val.Bool(a.Text() == b.Text())
	}
	return val.Bool(a.Real() == b.Real())
})

var Neq = check(func(a, b val.Val) val.Val {
	if val.IsText(a) && val.IsText(b) {
		return val.Bool(a.Text() != b.Text())
	}
	return val.Bool(a.Real() != b.Real())
})

var Lt = check(func(a, b val.Val) val.Val {
	if val.IsText(a) && val.IsText(b) {
		return val.Bool(a.Text() < b.Text())
	}
	return val.Bool(a.Real() < b.Real())
})

var Lte = check(func(a, b val.Val) val.Val {
	if val.IsText(a) && val.IsText(b) {
		return val.Bool(a.Text() <= b.Text())
	}
	return val.Bool(a.Real() <= b.Real())
})

var Gt = check(func(a, b val.Val) val.Val {
	if val.IsText(a) && val.IsText(b) {
		return val.Bool(a.Text() > b.Text())
	}
	return val.Bool(a.Real() > b.Real())
})

var Gte = check(func(a, b val.Val) val.Val {
	if val.IsText(a) && val.IsText(b) {
		return val.Bool(a.Text() >= b.Text())
	}
	return val.Bool(a.Real() >= b.Real())
})

var Concat = check(func(a, b val.Val) val.Val {
	return val.Text(a.Text() + b.Text())
})
