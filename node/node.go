package node

import (
	"github.com/suprunchuksergey/dpl/namespace"
	"github.com/suprunchuksergey/dpl/op"
	"github.com/suprunchuksergey/dpl/val"
	"reflect"
)

type binary struct {
	left, right Node
	op          op.Binary
}

func (b binary) Exec(ns namespace.Namespace) (val.Val, error) {
	l, err := b.left.Exec(ns)
	if err != nil {
		return nil, err
	}

	r, err := b.right.Exec(ns)
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

func (v value) Exec(_ namespace.Namespace) (val.Val, error) { return v.val, nil }

func newValue(val val.Val) value { return value{val} }

type unary struct {
	n  Node
	op op.Unary
}

func (u unary) Exec(ns namespace.Namespace) (val.Val, error) {
	v, err := u.n.Exec(ns)
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

type array struct{ items []Node }

func (a array) Exec(ns namespace.Namespace) (val.Val, error) {
	items := make([]val.Val, 0, len(a.items))
	for _, item := range a.items {
		v, err := item.Exec(ns)
		if err != nil {
			return nil, err
		}
		items = append(items, v)
	}
	return val.Array(items), nil
}

func newArray(items []Node) array { return array{items} }

type Record struct{ k, v Node }

func NewRecord(k, v Node) Record {
	return Record{k: k, v: v}
}

type Records []Record

type dict struct{ records Records }

func (d dict) Exec(ns namespace.Namespace) (val.Val, error) {
	m := make(map[string]val.Val, len(d.records))

	for _, record := range d.records {
		k, err := record.k.Exec(ns)
		if err != nil {
			return nil, err
		}
		v, err := record.v.Exec(ns)
		if err != nil {
			return nil, err
		}
		m[k.ToText()] = v
	}

	return val.Map(m), nil
}

func newDict(records Records) dict { return dict{records} }

type indexAccess struct{ v, index Node }

func (i indexAccess) Exec(ns namespace.Namespace) (val.Val, error) {
	v, err := i.v.Exec(ns)
	if err != nil {
		return nil, err
	}
	index, err := i.index.Exec(ns)
	if err != nil {
		return nil, err
	}
	return op.IndexAccess(v, index)
}

func newIndexAccess(v, index Node) indexAccess {
	return indexAccess{v: v, index: index}
}

type ident struct{ name string }

func (i ident) Exec(ns namespace.Namespace) (val.Val, error) {
	return ns.Get(i.name)
}

func newIdent(name string) ident { return ident{name: name} }

type commands struct{ commands []Node }

func (cmds commands) Exec(ns namespace.Namespace) (val.Val, error) {
	res := val.Null()
	for _, cmd := range cmds.commands {
		v, err := cmd.Exec(ns)
		if err != nil {
			return nil, err
		}
		res = v
	}
	return res, nil
}

func newCommands(cmds []Node) Node { return commands{commands: cmds} }

type Node interface {
	Exec(ns namespace.Namespace) (val.Val, error)
}

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

func Ident(name string) Node { return newIdent(name) }

func Array(v []Node) Node { return newArray(v) }
func Map(v Records) Node  { return newDict(v) }

func IndexAccess(v, index Node) Node { return newIndexAccess(v, index) }

func Commands(cmds []Node) Node { return newCommands(cmds) }

func DeepEqual(a, b Node) bool {
	if a == nil || b == nil {
		return a == b
	}
	switch aval := a.(type) {
	case value:
		bval, ok := b.(value)
		return ok && aval == bval

	case ident:
		bval, ok := b.(ident)
		return ok && aval.name == bval.name

	case array:
		bval, ok := b.(array)
		if !ok || len(aval.items) != len(bval.items) {
			return false
		}
		for i := range aval.items {
			if !DeepEqual(aval.items[i], bval.items[i]) {
				return false
			}
		}
		return true

	case dict:
		bval, ok := b.(dict)
		if !ok || len(aval.records) != len(bval.records) {
			return false
		}

		for i := range aval.records {
			if !DeepEqual(aval.records[i].k, bval.records[i].k) ||
				!DeepEqual(aval.records[i].v, bval.records[i].v) {
				return false
			}
		}

		return true

	case indexAccess:
		bval, ok := b.(indexAccess)
		return ok && DeepEqual(aval.v, bval.v) && DeepEqual(aval.index, bval.index)

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
