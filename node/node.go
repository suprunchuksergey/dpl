package node

import (
	"github.com/suprunchuksergey/dpl/op"
	"github.com/suprunchuksergey/dpl/val"
	"reflect"
)

type binary struct {
	left, right Node
	op          op.Binary
}

func (b binary) Exec() (val.Val, error) {
	l, err := b.left.Exec()
	if err != nil {
		return nil, err
	}

	r, err := b.right.Exec()
	if err != nil {
		return nil, err
	}

	return b.op(l, r), nil
}

func newBinary(l, r Node, op op.Binary) binary {
	return binary{
		left:  l,
		right: r,
		op:    op,
	}
}

type value struct{ val val.Val }

func (v value) Exec() (val.Val, error) { return v.val, nil }

func newValue(val val.Val) value { return value{val} }

type unary struct {
	n  Node
	op op.Unary
}

func (u unary) Exec() (val.Val, error) {
	v, err := u.n.Exec()
	if err != nil {
		return nil, err
	}
	return u.op(v), nil
}

func newUnary(n Node, op op.Unary) unary {
	return unary{
		n:  n,
		op: op,
	}
}

type Node interface{ Exec() (val.Val, error) }

func Add(l, r Node) Node { return newBinary(l, r, op.Add) }
func Sub(l, r Node) Node { return newBinary(l, r, op.Sub) }
func Mul(l, r Node) Node { return newBinary(l, r, op.Mul) }
func Div(l, r Node) Node { return newBinary(l, r, op.Div) }
func Rem(l, r Node) Node { return newBinary(l, r, op.Rem) }

func Eq(l, r Node) Node  { return newBinary(l, r, op.Eq) }
func Neq(l, r Node) Node { return newBinary(l, r, op.Neq) }
func Lt(l, r Node) Node  { return newBinary(l, r, op.Lt) }
func Lte(l, r Node) Node { return newBinary(l, r, op.Lte) }
func Gt(l, r Node) Node  { return newBinary(l, r, op.Gt) }
func Gte(l, r Node) Node { return newBinary(l, r, op.Gte) }

func Concat(l, r Node) Node { return newBinary(l, r, op.Concat) }

func And(l, r Node) Node { return newBinary(l, r, op.And) }
func Or(l, r Node) Node  { return newBinary(l, r, op.Or) }
func Not(n Node) Node    { return newUnary(n, op.Not) }

func Neg(n Node) Node { return newUnary(n, op.Neg) }

func Val(val val.Val) Node { return newValue(val) }

func Int(v int64) Node    { return Val(val.Int(v)) }
func Real(v float64) Node { return Val(val.Real(v)) }
func Text(v string) Node  { return Val(val.Text(v)) }
func True() Node          { return Val(val.True()) }
func False() Node         { return Val(val.False()) }
func Null() Node          { return Val(val.Null()) }

func DeepEqual(a, b Node) bool {
	if a == nil || b == nil {
		return a == b
	}
	switch aval := a.(type) {
	case value:
		bval, ok := b.(value)
		return ok && aval == bval

	case unary:
		bval, ok := b.(unary)
		return ok &&
			reflect.ValueOf(aval.op).Pointer() ==
				reflect.ValueOf(bval.op).Pointer() &&
			DeepEqual(aval.n, bval.n)

	case binary:
		bval, ok := b.(binary)
		return ok &&
			reflect.ValueOf(aval.op).Pointer() ==
				reflect.ValueOf(bval.op).Pointer() &&
			DeepEqual(aval.left, bval.left) &&
			DeepEqual(aval.right, bval.right)

	default:
		return false
	}
}
