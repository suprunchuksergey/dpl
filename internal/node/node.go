package node

import (
	"errors"
	"fmt"
	"github.com/suprunchuksergey/dpl/internal/namespace"
	"github.com/suprunchuksergey/dpl/internal/value"
	"math"
	"slices"
)

type Node interface {
	Exec(namespace namespace.Namespace) (value.Value, error)
}

type binary struct{ a, b Node }

// вычисляет оба узла
func (n binary) exec(
	namespace namespace.Namespace,
	validators ...func(value.Value, value.Value) error,
) (value.Value, value.Value, error) {

	a, err := n.a.Exec(namespace)
	if err != nil {
		return nil, nil, err
	}

	b, err := n.b.Exec(namespace)
	if err != nil {
		return nil, nil, err
	}

	for _, validator := range validators {
		if err := validator(a, b); err != nil {
			return nil, nil, err
		}
	}

	return a, b, nil
}

// шаблон для арифметических операторов
func (n binary) arithmetic(
	namespace namespace.Namespace,
	floatH func(float64, float64) float64,
	intH func(int64, int64) int64,
	validators ...func(value.Value, value.Value) error,
) (value.Value, error) {

	a, b, err := n.exec(namespace, validators...)
	if err != nil {
		return nil, err
	}

	if a.IsReal() || b.IsReal() {
		a, b, err := binaryToReal(a, b)
		if err != nil {
			return nil, err
		}
		return value.Real(floatH(a, b)), nil
	} else {
		a, b, err := binaryToInt(a, b)
		if err != nil {
			return nil, err
		}
		return value.Int(intH(a, b)), nil
	}
}

// шаблон для операторов сравнения
func (n binary) comparison(
	namespace namespace.Namespace,
	stringH func(string, string) bool,
	floatH func(float64, float64) bool,
	validators ...func(value.Value, value.Value) error,
) (value.Value, error) {

	a, b, err := n.exec(namespace, validators...)
	if err != nil {
		return nil, err
	}

	if a.IsText() && b.IsText() {
		return value.Bool(stringH(
			a.Text(), b.Text(),
		)), nil
	} else {
		a, b, err := binaryToReal(a, b)
		if err != nil {
			return nil, err
		}
		return value.Bool(floatH(a, b)), nil
	}
}

// шаблон для логических операторов
func (n binary) logic(
	namespace namespace.Namespace,
	boolH func(bool, bool) bool,
	validators ...func(value.Value, value.Value) error,
) (value.Value, error) {

	a, b, err := n.exec(namespace, validators...)
	if err != nil {
		return nil, err
	}

	aBool, bBool, err := binaryToBool(a, b)
	if err != nil {
		return nil, err
	}
	return value.Bool(boolH(aBool, bBool)), nil
}

func binaryToInt(a, b value.Value) (int64, int64, error) {
	aInt, err := a.Int()
	if err != nil {
		return 0, 0, err
	}
	bInt, err := b.Int()
	if err != nil {
		return 0, 0, err
	}
	return aInt, bInt, nil
}

func binaryToReal(a, b value.Value) (float64, float64, error) {
	aReal, err := a.Real()
	if err != nil {
		return 0, 0, err
	}
	bReal, err := b.Real()
	if err != nil {
		return 0, 0, err
	}
	return aReal, bReal, nil
}

func binaryToBool(a, b value.Value) (bool, bool, error) {
	aBool, err := a.Bool()
	if err != nil {
		return false, false, err
	}
	bBool, err := b.Bool()
	if err != nil {
		return false, false, err
	}
	return aBool, bBool, nil
}

func addOp[T int64 | float64 | string](a, b T) T { return a + b }
func subOp[T int64 | float64](a, b T) T          { return a - b }
func mulOp[T int64 | float64](a, b T) T          { return a * b }
func divOp[T int64 | float64](a, b T) T          { return a / b }

func eqOp[T float64 | string](a, b T) bool  { return a == b }
func neqOp[T float64 | string](a, b T) bool { return a != b }
func ltOp[T float64 | string](a, b T) bool  { return a < b }
func gtOp[T float64 | string](a, b T) bool  { return a > b }
func lteOp[T float64 | string](a, b T) bool { return a <= b }
func gteOp[T float64 | string](a, b T) bool { return a >= b }

func andOp(a, b bool) bool { return a && b }
func orOp(a, b bool) bool  { return a || b }

func opNotDefined(op, typ string) error {
	return fmt.Errorf("оператор %s не определен для типа %s", op, typ)
}

func getCheckOpNotDefined(
	op string,
	whitelist ...string,
) func(value.Value) error {
	return func(v value.Value) error {
		for _, i := range whitelist {
			if v.Type() == i {
				return nil
			}
		}
		return opNotDefined(op, v.Type())
	}
}

func getBinaryCheckOpNotDefined(
	op string,
	whitelist ...string,
) func(value.Value, value.Value) error {
	check := getCheckOpNotDefined(op, whitelist...)
	return func(a, b value.Value) error {
		if err := check(a); err != nil {
			return err
		}
		return check(b)
	}
}

var baseWhitelist = []string{
	value.IntType,
	value.RealType,
	value.TextType,
	value.BoolType,
	value.NullType,
}

type add struct{ binary }

func (n add) Exec(namespace namespace.Namespace) (value.Value, error) {
	return n.arithmetic(
		namespace,
		addOp[float64],
		addOp[int64],
		getBinaryCheckOpNotDefined("+", baseWhitelist...),
	)
}

func Add(a, b Node) Node { return add{binary{a: a, b: b}} }

type sub struct{ binary }

func (n sub) Exec(namespace namespace.Namespace) (value.Value, error) {
	return n.arithmetic(
		namespace,
		subOp[float64],
		subOp[int64],
		getBinaryCheckOpNotDefined("-", baseWhitelist...),
	)
}

func Sub(a, b Node) Node { return sub{binary{a: a, b: b}} }

type mul struct{ binary }

func (n mul) Exec(namespace namespace.Namespace) (value.Value, error) {
	return n.arithmetic(
		namespace,
		mulOp[float64],
		mulOp[int64],
		getBinaryCheckOpNotDefined("*", baseWhitelist...),
	)
}

func Mul(a, b Node) Node { return mul{binary{a: a, b: b}} }

func divByZero() error { return errors.New("деление на ноль") }

func checkDivByZero(_, b value.Value) error {
	bReal, err := b.Real()
	if err != nil {
		return err
	}
	if bReal == 0 {
		return divByZero()
	}
	return nil
}

type div struct{ binary }

func (n div) Exec(namespace namespace.Namespace) (value.Value, error) {
	return n.arithmetic(
		namespace,
		divOp[float64],
		divOp[int64],
		getBinaryCheckOpNotDefined("/", baseWhitelist...),
		checkDivByZero,
	)
}

func Div(a, b Node) Node { return div{binary{a: a, b: b}} }

type mod struct{ binary }

func (n mod) Exec(namespace namespace.Namespace) (value.Value, error) {
	return n.arithmetic(
		namespace,
		func(a, b float64) float64 { return math.Mod(a, b) },
		func(a, b int64) int64 { return a % b },
		getBinaryCheckOpNotDefined("%", baseWhitelist...),
		checkDivByZero,
	)
}

func Mod(a, b Node) Node { return mod{binary{a: a, b: b}} }

type concat struct{ binary }

func (n concat) Exec(namespace namespace.Namespace) (value.Value, error) {
	a, b, err := n.exec(namespace)
	if err != nil {
		return nil, err
	}

	return value.Text(a.Text() + b.Text()), nil
}

func Concat(a, b Node) Node { return concat{binary{a: a, b: b}} }

type eq struct{ binary }

func (n eq) Exec(namespace namespace.Namespace) (value.Value, error) {
	return n.comparison(
		namespace,
		eqOp[string],
		eqOp[float64],
		getBinaryCheckOpNotDefined("==", baseWhitelist...),
	)
}

func Eq(a, b Node) Node { return eq{binary{a: a, b: b}} }

type neq struct{ binary }

func (n neq) Exec(namespace namespace.Namespace) (value.Value, error) {
	return n.comparison(
		namespace,
		neqOp[string],
		neqOp[float64],
		getBinaryCheckOpNotDefined("!=", baseWhitelist...),
	)
}

func Neq(a, b Node) Node { return neq{binary{a: a, b: b}} }

type lt struct{ binary }

func (n lt) Exec(namespace namespace.Namespace) (value.Value, error) {
	return n.comparison(
		namespace,
		ltOp[string],
		ltOp[float64],
		getBinaryCheckOpNotDefined("<", baseWhitelist...),
	)
}

func Lt(a, b Node) Node { return lt{binary{a: a, b: b}} }

type gt struct{ binary }

func (n gt) Exec(namespace namespace.Namespace) (value.Value, error) {
	return n.comparison(
		namespace,
		gtOp[string],
		gtOp[float64],
		getBinaryCheckOpNotDefined(">", baseWhitelist...),
	)
}

func Gt(a, b Node) Node { return gt{binary{a: a, b: b}} }

type lte struct{ binary }

func (n lte) Exec(namespace namespace.Namespace) (value.Value, error) {
	return n.comparison(
		namespace,
		lteOp[string],
		lteOp[float64],
		getBinaryCheckOpNotDefined("<=", baseWhitelist...),
	)
}

func Lte(a, b Node) Node { return lte{binary{a: a, b: b}} }

type gte struct{ binary }

func (n gte) Exec(namespace namespace.Namespace) (value.Value, error) {
	return n.comparison(
		namespace,
		gteOp[string],
		gteOp[float64],
		getBinaryCheckOpNotDefined(">=", baseWhitelist...),
	)
}

func Gte(a, b Node) Node { return gte{binary{a: a, b: b}} }

var logicWhitelist = []string{
	value.IntType,
	value.RealType,
	value.TextType,
	value.BoolType,
	value.NullType,
	value.ArrayType,
	value.ObjectType,
}

type and struct{ binary }

func (n and) Exec(namespace namespace.Namespace) (value.Value, error) {
	return n.logic(
		namespace,
		andOp,
		getBinaryCheckOpNotDefined("and", logicWhitelist...),
	)
}

func And(a, b Node) Node { return and{binary{a: a, b: b}} }

type or struct{ binary }

func (n or) Exec(namespace namespace.Namespace) (value.Value, error) {
	return n.logic(
		namespace,
		orOp,
		getBinaryCheckOpNotDefined("or", logicWhitelist...),
	)
}

func Or(a, b Node) Node { return or{binary{a: a, b: b}} }

type unary struct{ v Node }

func (n unary) exec(
	namespace namespace.Namespace,
	validators ...func(value.Value) error,
) (value.Value, error) {

	v, err := n.v.Exec(namespace)
	if err != nil {
		return nil, err
	}

	for _, validator := range validators {
		if err := validator(v); err != nil {
			return nil, err
		}
	}

	return v, nil
}

type neg struct{ unary }

func (n neg) Exec(namespace namespace.Namespace) (value.Value, error) {
	v, err := n.exec(namespace, getCheckOpNotDefined("унарный -", baseWhitelist...))
	if err != nil {
		return nil, err
	}

	if v.IsReal() {
		v, err := v.Real()
		if err != nil {
			return nil, err
		}
		return value.Real(-v), nil
	} else {
		v, err := v.Int()
		if err != nil {
			return nil, err
		}
		return value.Int(-v), nil
	}
}

func Neg(v Node) Node { return neg{unary{v: v}} }

type not struct{ unary }

func (n not) Exec(namespace namespace.Namespace) (value.Value, error) {
	v, err := n.exec(namespace, getCheckOpNotDefined("not", logicWhitelist...))
	if err != nil {
		return nil, err
	}

	vBool, err := v.Bool()
	if err != nil {
		return nil, err
	}

	return value.Bool(!vBool), nil
}

func Not(v Node) Node { return not{unary{v: v}} }

type valueNode struct{ v value.Value }

func (n valueNode) Exec(_ namespace.Namespace) (value.Value, error) {
	return n.v, nil
}

func Int(v int64) Node    { return valueNode{v: value.Int(v)} }
func Real(v float64) Node { return valueNode{v: value.Real(v)} }
func Text(v string) Node  { return valueNode{v: value.Text(v)} }
func Bool(v bool) Node    { return valueNode{v: value.Bool(v)} }

type array struct{ nodes []Node }

func (n array) Exec(namespace namespace.Namespace) (value.Value, error) {
	values := make([]value.Value, 0, len(n.nodes))

	for _, node := range n.nodes {
		v, err := node.Exec(namespace)
		if err != nil {
			return nil, err
		}
		values = append(values, v)
	}

	return value.Array(values...), nil
}

func Array(nodes ...Node) Node { return array{nodes: nodes} }

type KV struct{ Key, Value Node }

type object struct{ pairs []KV }

func (n object) Exec(namespace namespace.Namespace) (value.Value, error) {
	pairs := make([]value.KV, 0, len(n.pairs))

	for _, pair := range n.pairs {
		key, err := pair.Key.Exec(namespace)
		if err != nil {
			return nil, err
		}

		v, err := pair.Value.Exec(namespace)
		if err != nil {
			return nil, err
		}

		pairs = append(pairs, value.KV{Key: key, Value: v})
	}

	return value.Object(pairs...), nil
}

func Object(pairs ...KV) Node {
	return object{pairs: pairs}
}

func Null() Node { return valueNode{v: value.Null()} }

func wrongIndex(typ string) error {
	return fmt.Errorf(`тип %s не подходит в качестве индекса`, typ)
}

type elByIndex struct{ v, index Node }

func (n elByIndex) Exec(namespace namespace.Namespace) (value.Value, error) {
	v, err := n.v.Exec(namespace)
	if err != nil {
		return nil, err
	}

	check := getCheckOpNotDefined(
		"[<index>]",
		value.TextType, value.ArrayType, value.ObjectType)
	if err := check(v); err != nil {
		return nil, err
	}

	index, err := n.index.Exec(namespace)
	if err != nil {
		return nil, err
	}

	if v.Type() == value.ObjectType {
		return v.ElByIndex(index)
	}

	if !slices.Contains(baseWhitelist, index.Type()) {
		return nil, wrongIndex(index.Type())
	}

	return v.ElByIndex(index)
}

func ElByIndex(v, index Node) Node { return elByIndex{v: v, index: index} }

type ident struct{ v string }

func (n ident) Exec(namespace namespace.Namespace) (value.Value, error) {
	return namespace.Get(n.v)
}

func Ident(v string) Node { return ident{v: v} }

func idExpected() error {
	return errors.New("ожидался идентификатор")
}

type create struct{ name, v Node }

func (n create) Exec(namespace namespace.Namespace) (value.Value, error) {
	id, ok := n.name.(ident)
	if !ok {
		return nil, idExpected()
	}

	v, err := n.v.Exec(namespace)
	if err != nil {
		return nil, err
	}

	if err := namespace.Create(id.v, v); err != nil {
		return nil, err
	}

	return v, nil
}

func Create(name, v Node) Node { return create{name: name, v: v} }

type set struct{ name, v Node }

func (n set) Exec(namespace namespace.Namespace) (value.Value, error) {
	v, err := n.v.Exec(namespace)
	if err != nil {
		return nil, err
	}

	if id, ok := n.name.(ident); ok {
		namespace.Set(id.v, v)
		return v, nil
	}

	id := n.name
	indexes := make([]value.Value, 0, 1)
	for {
		index, ok := id.(elByIndex)
		if !ok {
			break
		}
		i, err := index.index.Exec(namespace)
		if err != nil {
			return nil, err
		}
		indexes = append(indexes, i)
		id = index.v
	}

	if _, ok := id.(ident); !ok {
		return nil, idExpected()
	}

	name := id.(ident).v
	target, err := namespace.Get(name)
	if err != nil {
		return nil, err
	}

	for i := len(indexes) - 1; i >= 0; i-- {
		if i == 0 {
			if err := target.SetElByIndex(indexes[i], v); err != nil {
				return nil, err
			}
			break
		}

		target, err = target.ElByIndex(indexes[i])
		if err != nil {
			return nil, err
		}
	}

	return v, nil
}

func Set(name, v Node) Node { return set{name: name, v: v} }

type block struct{ cmds []Node }

func (n block) Exec(namespace namespace.Namespace) (value.Value, error) {
	v := value.Null()

	for _, cmd := range n.cmds {
		val, err := cmd.Exec(namespace)
		if err != nil {
			return nil, err
		}
		v = val
	}

	return v, nil
}

func Block(cmds ...Node) Node { return block{cmds: cmds} }

type Branch struct{ Cond, Body Node }

type branch struct {
	branches []Branch
}

func (n branch) Exec(namespace namespace.Namespace) (value.Value, error) {
	for _, b := range n.branches {
		cond, err := b.Cond.Exec(namespace)
		if err != nil {
			return nil, err
		}

		condBool, err := cond.Bool()
		if err != nil {
			return nil, err
		}

		if condBool {
			return b.Body.Exec(namespace.New(nil))
		}
	}

	return value.Null(), nil
}

func If(branches ...Branch) Node { return branch{branches: branches} }

func tooFewRecipients() error { return errors.New("слишком мало получателей") }

func tooManyRecipients() error { return errors.New("слишком много получателей") }

type loop struct {
	recipients []Node
	from       Node
	body       Node
}

func (n loop) Exec(namespace namespace.Namespace) (value.Value, error) {
	names := make([]string, 0, len(n.recipients))
	for _, recipient := range n.recipients {
		r, ok := recipient.(ident)
		if !ok {
			return nil, idExpected()
		}
		names = append(names, r.v)
	}

	from, err := n.from.Exec(namespace)
	if err != nil {
		return nil, err
	}

	switch len(names) {
	case 0:
		return nil, tooFewRecipients()

	case 1:
		iter, err := from.Iter()
		if err != nil {
			return nil, err
		}

		res := value.Null()
		for i := range iter {
			val, err := n.body.Exec(namespace.New(map[string]value.Value{
				names[0]: i,
			}))
			if err != nil {
				return nil, err
			}

			res = val
		}

		return res, nil

	case 2:
		iter, err := from.Iter2()
		if err != nil {
			return nil, err
		}

		res := value.Null()
		for i, j := range iter {
			val, err := n.body.Exec(namespace.New(map[string]value.Value{
				names[0]: i,
				names[1]: j,
			}))
			if err != nil {
				return nil, err
			}

			res = val
		}

		return res, nil

	default:
		return nil, tooManyRecipients()
	}
}

func For(recipients []Node, from, body Node) Node {
	return loop{
		recipients: recipients,
		from:       from,
		body:       body,
	}
}

type call struct {
	target Node
	args   []Node
}

func (n call) Exec(namespace namespace.Namespace) (value.Value, error) {
	target, err := n.target.Exec(namespace)
	if err != nil {
		return nil, err
	}

	if target.Type() != value.FunctionType {
		return nil, opNotDefined("вызов функции", target.Type())
	}

	args := make([]value.Value, 0, len(n.args))
	for _, arg := range n.args {
		val, err := arg.Exec(namespace)
		if err != nil {
			return nil, err
		}
		args = append(args, val)
	}

	return target.Call(args...)
}

func Call(target Node, args ...Node) Node {
	return call{
		target: target,
		args:   args,
	}
}

type returnErr struct{ v value.Value }

func (r returnErr) Error() string {
	return "return может использоваться только в контексте функции"
}

type returnNode struct{ v Node }

func (n returnNode) Exec(namespace namespace.Namespace) (value.Value, error) {
	v, err := n.v.Exec(namespace)
	if err != nil {
		return nil, err
	}
	return nil, returnErr{v: v}
}

func Return(v Node) Node { return returnNode{v: v} }

type function struct {
	params []Node
	body   Node
}

func (n function) Exec(namespace namespace.Namespace) (value.Value, error) {
	names := make([]string, 0, len(n.params))
	for _, param := range n.params {
		id, ok := param.(ident)
		if !ok {
			return nil, idExpected()
		}
		names = append(names, id.v)
	}

	return value.Function(
		func(args ...value.Value) (value.Value, error) {
			init := make(map[string]value.Value, len(names))

			for i, name := range names {
				if i >= len(args) {
					init[name] = value.Null()
					continue
				}

				init[name] = args[i]
			}

			res, err := n.body.Exec(namespace.New(init))
			if err != nil {
				if err, ok := err.(returnErr); ok {
					return err.v, nil
				}

				return nil, err
			}

			return res, nil
		},
	), nil
}

func Function(body Node, params ...Node) Node {
	return function{
		body:   body,
		params: params,
	}
}
