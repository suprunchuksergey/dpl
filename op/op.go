package op

import (
	. "github.com/suprunchuksergey/dpl/val"
	"math"
)

type Unary func(v Val) Val

type Binary func(a, b Val) Val

func Add(a, b Val) Val {
	if a.IsReal() || b.IsReal() {
		return Real(a.ToReal() + b.ToReal())
	}
	return Int(a.ToInt() + b.ToInt())
}

func Sub(a, b Val) Val {
	if a.IsReal() || b.IsReal() {
		return Real(a.ToReal() - b.ToReal())
	}
	return Int(a.ToInt() - b.ToInt())
}

func Mul(a, b Val) Val {
	if a.IsReal() || b.IsReal() {
		return Real(a.ToReal() * b.ToReal())
	}
	return Int(a.ToInt() * b.ToInt())
}

func Div(a, b Val) Val {
	if b.ToReal() == 0 {
		return Null()
	}
	if a.IsReal() || b.IsReal() {
		return Real(a.ToReal() / b.ToReal())
	}
	return Int(a.ToInt() / b.ToInt())
}

func Rem(a, b Val) Val {
	if b.ToReal() == 0 {
		return Null()
	}
	if a.IsReal() || b.IsReal() {
		return Real(math.Mod(a.ToReal(), b.ToReal()))
	}
	return Int(a.ToInt() % b.ToInt())
}

func Eq(a, b Val) Val {
	if a.IsText() && b.IsText() {
		return Bool(a.ToText() == b.ToText())
	}
	return Bool(a.ToReal() == b.ToReal())
}

func Neq(a, b Val) Val {
	if a.IsText() && b.IsText() {
		return Bool(a.ToText() != b.ToText())
	}
	return Bool(a.ToReal() != b.ToReal())
}

func Lt(a, b Val) Val {
	if a.IsText() && b.IsText() {
		return Bool(a.ToText() < b.ToText())
	}
	return Bool(a.ToReal() < b.ToReal())
}

func Lte(a, b Val) Val {
	if a.IsText() && b.IsText() {
		return Bool(a.ToText() <= b.ToText())
	}
	return Bool(a.ToReal() <= b.ToReal())
}

func Gt(a, b Val) Val {
	if a.IsText() && b.IsText() {
		return Bool(a.ToText() > b.ToText())
	}
	return Bool(a.ToReal() > b.ToReal())
}

func Gte(a, b Val) Val {
	if a.IsText() && b.IsText() {
		return Bool(a.ToText() >= b.ToText())
	}
	return Bool(a.ToReal() >= b.ToReal())
}

func Concat(a, b Val) Val {
	return Text(a.ToText() + b.ToText())
}

func And(a, b Val) Val {
	return Bool(a.ToBool() && b.ToBool())
}

func Or(a, b Val) Val {
	return Bool(a.ToBool() || b.ToBool())
}

func Not(v Val) Val {
	return Bool(!v.ToBool())
}

func Neg(v Val) Val {
	if v.IsReal() {
		return Real(-v.ToReal())
	}
	return Int(-v.ToInt())
}
