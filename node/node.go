package node

import (
	"errors"
	"fmt"
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

type assign struct{ l, r Node }

func (a assign) Exec(ns namespace.Namespace) (val.Val, error) {
	indexes := make([]Node, 0)
	for {
		if access, ok := a.l.(indexAccess); ok {
			indexes = append(indexes, access.index)
			a.l = access.v
			continue
		}
		break
	}

	if _, ok := a.l.(ident); !ok {
		return nil, errors.New("невозможно присвоить значение неидентификатору")
	}

	v, err := a.r.Exec(ns)
	if err != nil {
		return nil, err
	}

	if len(indexes) == 0 {
		ns.Set(a.l.(ident).name, v)
		return v, nil
	}

	l, err := ns.Get(a.l.(ident).name)
	if err != nil {
		return nil, err
	}

	for i := len(indexes) - 1; i >= 0; i-- {
		ind := indexes[i]

		index, err := ind.Exec(ns)
		if err != nil {
			return nil, err
		}

		if i == 0 {
			if l.IsArray() {
				l.ToArray()[index.ToInt()] = v
			} else if l.IsMap() {
				l.ToMap()[index.ToText()] = v
			} else {
				return nil, fmt.Errorf("невозможно получить доступ по индексу к %s", l)
			}
			break
		}

		l, err = op.IndexAccess(l, index)
		if err != nil {
			return nil, err
		}
	}

	return v, nil
}

func newAssign(l, r Node) Node {
	return assign{
		l: l,
		r: r,
	}
}

type Branch struct {
	cond Node
	body Node
}

func NewBranch(cond, body Node) *Branch {
	return &Branch{
		cond: cond,
		body: body,
	}
}

type ifStmt struct {
	first  *Branch
	second []*Branch
	third  *Branch
}

func (stmt ifStmt) Exec(ns namespace.Namespace) (val.Val, error) {
	//условие if
	cond, err := stmt.first.cond.Exec(ns)
	if err != nil {
		return nil, err
	}
	//если true
	if cond.ToBool() {
		//вычислить тело и вернуть результат
		res, err := stmt.first.body.Exec(ns)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	//для каждого else if
	for _, b := range stmt.second {
		//условие else if
		cond, err := b.cond.Exec(ns)
		if err != nil {
			return nil, err
		}
		//если true
		if cond.ToBool() {
			//вычислить тело и вернуть результат
			res, err := b.body.Exec(ns)
			if err != nil {
				return nil, err
			}
			return res, nil
		}
		//иначе перейти к следующему else if
	}

	//если есть else
	if stmt.third != nil {
		//вычислить тело и вернуть результат
		res, err := stmt.third.body.Exec(ns)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	//иначе null
	return val.Null(), nil
}

func newIf(first *Branch, second []*Branch, third *Branch) Node {
	return ifStmt{
		first:  first,
		second: second,
		third:  third,
	}
}

type loop struct {
	rec1  Node
	rec2  Node //может быть nil
	value Node
	body  Node
}

func (l loop) exec1(ns namespace.Namespace) (val.Val, error) {
	id, ok := l.rec1.(ident)
	if !ok {
		return nil, errors.New("ожидался идентификатор")
	}
	name := id.name

	v, err := l.value.Exec(ns)
	if err != nil {
		return nil, err
	}

	if !v.CanIter() {
		return nil, errors.New("значение не поддерживает итерацию")
	}

	res := val.Null()
	for i := range v.Iter() {
		newns := namespace.WithParent(ns, map[string]val.Val{name: i})
		body, err := l.body.Exec(newns)
		if err != nil {
			return nil, err
		}
		res = body
	}

	return res, nil
}

func (l loop) exec2(ns namespace.Namespace) (val.Val, error) {
	id1, ok1 := l.rec1.(ident)
	if !ok1 {
		return nil, errors.New("ожидался идентификатор")
	}

	id2, ok2 := l.rec2.(ident)
	if !ok2 {
		return nil, errors.New("ожидался идентификатор")
	}

	a := id1.name
	b := id2.name

	v, err := l.value.Exec(ns)
	if err != nil {
		return nil, err
	}

	if !v.CanIter2() {
		return nil, errors.New("значение не поддерживает итерацию")
	}

	res := val.Null()
	for i, j := range v.Iter2() {
		newns := namespace.WithParent(ns, map[string]val.Val{a: i, b: j})
		body, err := l.body.Exec(newns)
		if err != nil {
			return nil, err
		}
		res = body
	}

	return res, nil
}

func (l loop) Exec(ns namespace.Namespace) (val.Val, error) {
	if l.rec1 == nil && l.rec2 == nil {
		return nil, errors.New("для выполнения требуется по крайней мере один получатель")
	}

	if l.rec1 != nil && l.rec2 != nil {
		return l.exec2(ns)
	}

	return l.exec1(ns)
}

type call struct {
	target Node
	args   []Node
}

func (c call) Exec(ns namespace.Namespace) (val.Val, error) {
	target, err := c.target.Exec(ns)
	if err != nil {
		return nil, err
	}
	if !target.IsFn() {
		return nil, fmt.Errorf("%s не функция", target)
	}
	args := make([]val.Val, 0, len(c.args))
	for _, arg := range c.args {
		v, err := arg.Exec(ns)
		if err != nil {
			return nil, err
		}
		args = append(args, v)
	}
	return target.Call(args)
}

type returnerr struct {
	v val.Val
}

func (r returnerr) Error() string {
	return "return может использоваться только в контексте функции"
}

type fnreturn struct{ val Node }

func (f fnreturn) Exec(ns namespace.Namespace) (val.Val, error) {
	v, err := f.val.Exec(ns)
	if err != nil {
		return nil, err
	}
	return nil, returnerr{v}
}

type fn struct {
	names []Node
	body  Node
}

func (f fn) Exec(ns namespace.Namespace) (val.Val, error) {
	names := make([]string, 0, len(f.names))
	for _, name := range f.names {
		if _, ok := name.(ident); !ok {
			return nil, errors.New("неверный аргумент функции")
		}
		names = append(names, name.(ident).name)
	}

	return val.Fn(func(args []val.Val) (val.Val, error) {
		m := make(map[string]val.Val)

		for i, name := range names {
			if i >= len(names) {
				m[name] = val.Null()
				continue
			}
			m[name] = args[i]
		}

		newns := namespace.WithParent(ns, m)
		res, err := f.body.Exec(newns)
		if err != nil {
			if err, ok := err.(returnerr); ok {
				return err.v, nil
			}
			return nil, err
		}
		return res, nil
	}, names), nil
}

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

func Assign(l, r Node) Node { return newAssign(l, r) }

func Call(target Node, args []Node) Node {
	return call{
		target: target,
		args:   args,
	}
}

func Return(val Node) Node {
	return fnreturn{val}
}

func Fn(body Node, names []Node) Node {
	return fn{names: names, body: body}
}

func For(rec1, rec2, value, body Node) Node {
	return loop{
		rec1:  rec1,
		rec2:  rec2,
		value: value,
		body:  body,
	}
}

func If(first *Branch, second []*Branch, third *Branch) Node { return newIf(first, second, third) }

func DeepEqual(a, b any) bool {
	if a == nil || b == nil {
		return a == b
	}
	switch aval := a.(type) {
	case *Branch:
		bval, ok := b.(*Branch)
		if !ok {
			return false
		}

		if aval == nil || bval == nil {
			return aval == bval
		}
		return DeepEqual(aval.cond, bval.cond) &&
			DeepEqual(aval.body, bval.body)

	case ifStmt:
		bval, ok := b.(ifStmt)
		if !ok || !DeepEqual(aval.first, bval.first) {
			return false
		}

		if len(aval.second) != len(bval.second) {
			return false
		}

		for i := range aval.second {
			if !DeepEqual(aval.second[i], bval.second[i]) {
				return false
			}
		}

		return DeepEqual(aval.third, bval.third)

	case value:
		bval, ok := b.(value)
		return ok && aval == bval

	case ident:
		bval, ok := b.(ident)
		return ok && aval.name == bval.name

	case assign:
		bval, ok := b.(assign)
		return ok &&
			DeepEqual(aval.l, bval.l) &&
			DeepEqual(aval.r, bval.r)

	case commands:
		bval, ok := b.(commands)
		if !ok || len(aval.commands) != len(bval.commands) {
			return false
		}
		for i := range aval.commands {
			if !DeepEqual(aval.commands[i], bval.commands[i]) {
				return false
			}
		}
		return true

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

	case loop:
		bval, ok := b.(loop)
		if !ok {
			return false
		}

		return DeepEqual(aval.rec1, bval.rec1) && DeepEqual(aval.rec2, bval.rec2) &&
			DeepEqual(aval.body, bval.body) && DeepEqual(aval.value, bval.value)

	case call:
		bval, ok := b.(call)
		if !ok || len(aval.args) != len(bval.args) {
			return false
		}

		for i := range aval.args {
			if !DeepEqual(aval.args[i], bval.args[i]) {
				return false
			}
		}

		return DeepEqual(aval.target, bval.target)

	case fn:
		bval, ok := b.(fn)
		if !ok || len(aval.names) != len(bval.names) {
			return false
		}
		for i := range aval.names {
			if !DeepEqual(aval.names[i], bval.names[i]) {
				return false
			}
		}
		return DeepEqual(aval.body, bval.body)

	default:
		return false
	}
}
