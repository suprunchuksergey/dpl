package node

import (
	"github.com/suprunchuksergey/dpl/op"
	"github.com/suprunchuksergey/dpl/val"
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

type Node interface{ Exec() (val.Val, error) }

func Add(l, r Node) Node    { return newBinary(l, r, op.Add) }
func Sub(l, r Node) Node    { return newBinary(l, r, op.Sub) }
func Mul(l, r Node) Node    { return newBinary(l, r, op.Mul) }
func Div(l, r Node) Node    { return newBinary(l, r, op.Div) }
func Rem(l, r Node) Node    { return newBinary(l, r, op.Rem) }
func Eq(l, r Node) Node     { return newBinary(l, r, op.Eq) }
func Neq(l, r Node) Node    { return newBinary(l, r, op.Neq) }
func Lt(l, r Node) Node     { return newBinary(l, r, op.Lt) }
func Lte(l, r Node) Node    { return newBinary(l, r, op.Lte) }
func Gt(l, r Node) Node     { return newBinary(l, r, op.Gt) }
func Gte(l, r Node) Node    { return newBinary(l, r, op.Gte) }
func Concat(l, r Node) Node { return newBinary(l, r, op.Concat) }

func Val(val val.Val) Node { return newValue(val) }
