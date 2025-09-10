package parser

import (
	"fmt"
	"github.com/suprunchuksergey/dpl/internal/node"
	"github.com/suprunchuksergey/dpl/lexer"
	"github.com/suprunchuksergey/dpl/pos"
	"github.com/suprunchuksergey/dpl/token"
	"reflect"
	"testing"
)

const (
	_ uint8 = iota
	parse
	parseBody
	layer1
	layer2
	layer3
	layer4
	layer5
	layer6
	layer7
	layer8
	layer9
	layer10
	layer11
	layer12
	layer13
	layer14
)

func call(p *parser, l uint8) (node.Node, error) {
	switch l {
	case parse:
		return p.Parse()
	case parseBody:
		return p.parseBody()
	case layer1:
		return p.layer1()
	case layer2:
		return p.layer2()
	case layer3:
		return p.layer3()
	case layer4:
		return p.layer4()
	case layer5:
		return p.layer5()
	case layer6:
		return p.layer6()
	case layer7:
		return p.layer7()
	case layer8:
		return p.layer8()
	case layer9:
		return p.layer9()
	case layer10:
		return p.layer10()
	case layer11:
		return p.layer11()
	case layer12:
		return p.layer12()
	case layer13:
		return p.layer13()
	case layer14:
		return p.layer14()
	default:
		panic("неизвестный слой")
	}
}

func prepare(data string) (*parser, error) {
	lex := lexer.New(data)
	err := lex.Next()
	if err != nil {
		return nil, fmt.Errorf("%s: %s", data, err.Error())
	}
	return newParser(lex), nil
}

type row struct {
	data     string
	expected node.Node
}

func (r row) exec(t *testing.T, l uint8) {
	p, err := prepare(r.data)
	if err != nil {
		t.Error(err.Error())
		return
	}

	v, err := call(p, l)
	if err != nil {
		t.Errorf("%s: %s", r.data, err.Error())
		return
	}

	if reflect.DeepEqual(r.expected, v) == false {
		t.Errorf("%s: ожидалось %v, получено %v",
			r.data, r.expected, v)
	}
}

func r(data string, expected node.Node) row {
	return row{
		data:     data,
		expected: expected,
	}
}

type rows []row

func (rows rows) exec(t *testing.T, l uint8) {
	for _, row := range rows {
		row.exec(t, l)
	}
}

func rs(rs ...row) rows {
	return rs
}

func Test_layer1(t *testing.T) {
	rs(
		r("true", node.Bool(true)),
		r("false", node.Bool(false)),
		r("null", node.Null()),
		r("19683", node.Int(19683)),
		r("19.683", node.Real(19.683)),
		r("'text'", node.Text("text")),
		r("fun", node.Ident("fun")),
		r("[]", node.Array()),
		r("[true, 19683, 'text',]",
			node.Array(
				node.Bool(true), node.Int(19683), node.Text("text"),
			)),
		r("[true, 19683, 'text', [true, 19683, 'text'],]",
			node.Array(
				node.Bool(true), node.Int(19683), node.Text("text"),
				node.Array(node.Bool(true), node.Int(19683), node.Text("text")),
			)),
		r("{}", node.Object()),
		r(
			"{'text' : 19683}",
			node.Object(
				node.KV{node.Text("text"), node.Int(19683)},
			)),
		r(
			"{'text' : 19683,}",
			node.Object(
				node.KV{node.Text("text"), node.Int(19683)},
			)),
	).exec(t, layer1)
}

func Test_layer3(t *testing.T) {
	rs(
		r("true", node.Bool(true)),
		r("false", node.Bool(false)),
		r("null", node.Null()),
		r("19683", node.Int(19683)),
		r("19.683", node.Real(19.683)),
		r("'text'", node.Text("text")),
		r("[]", node.Array()),
		r("[true, 19683, 'text',]",
			node.Array(
				node.Bool(true), node.Int(19683), node.Text("text"),
			)),
		r("[true, 19683, 'text', [true, 19683, 'text'],]",
			node.Array(
				node.Bool(true), node.Int(19683), node.Text("text"),
				node.Array(node.Bool(true), node.Int(19683), node.Text("text")),
			)),
		r("{}", node.Object()),
		r(
			"{'text' : 19683}",
			node.Object(node.KV{node.Text("text"), node.Int(19683)})),
		r(
			"{'text' : 19683,}",
			node.Object(node.KV{node.Text("text"), node.Int(19683)})),
		r("[683, 9][1]", node.ElByIndex(
			node.Array(node.Int(683), node.Int(9)),
			node.Int(1),
		)),
		r("[683, 9][1][3+9]", node.ElByIndex(
			node.ElByIndex(
				node.Array(node.Int(683), node.Int(9)),
				node.Int(1),
			), node.Add(node.Int(3), node.Int(9)),
		)),
		r("r(3,9)", node.Call(
			node.Ident("r"),
			node.Int(3), node.Int(9),
		)),
		r("r(3,9,)", node.Call(
			node.Ident("r"),
			node.Int(3), node.Int(9),
		)),
		r("r(3)", node.Call(
			node.Ident("r"),
			node.Int(3),
		)),
		r("r()", node.Call(
			node.Ident("r"),
		)),
		r("[683, 9][1][3+9](9)", node.Call(
			node.ElByIndex(
				node.ElByIndex(
					node.Array(node.Int(683), node.Int(9)),
					node.Int(1),
				), node.Add(node.Int(3), node.Int(9)),
			),
			node.Int(9),
		)),
	).exec(t, layer3)
}

func Test_layer4(t *testing.T) {
	rs(
		r("true", node.Bool(true)),
		r("false", node.Bool(false)),
		r("null", node.Null()),
		r("19683", node.Int(19683)),
		r("19.683", node.Real(19.683)),
		r("'text'", node.Text("text")),
		r("-19", node.Neg(node.Int(19))),
		r("--19",
			node.Neg(
				node.Neg(
					node.Int(19)))),
		r("---19",
			node.Neg(
				node.Neg(
					node.Neg(
						node.Int(19))))),
	).exec(t, layer4)
}

func Test_layer5(t *testing.T) {
	rs(
		r("true", node.Bool(true)),
		r("false", node.Bool(false)),
		r("null", node.Null()),
		r("19683", node.Int(19683)),
		r("19.683", node.Real(19.683)),
		r("'text'", node.Text("text")),
		r("-19", node.Neg(node.Int(19))),
		r("--19",
			node.Neg(
				node.Neg(
					node.Int(19)))),
		r("---19",
			node.Neg(
				node.Neg(
					node.Neg(
						node.Int(19))))),
		r("19*683", node.Mul(node.Int(19), node.Int(683))),
		r("19/683", node.Div(node.Int(19), node.Int(683))),
		r("19%683", node.Mod(node.Int(19), node.Int(683))),
		r("19*683/-3",
			node.Div(
				node.Mul(node.Int(19), node.Int(683)),
				node.Neg(node.Int(3)))),
	).exec(t, layer5)
}

func Test_layer6(t *testing.T) {
	rs(
		r("true", node.Bool(true)),
		r("false", node.Bool(false)),
		r("null", node.Null()),
		r("19683", node.Int(19683)),
		r("19.683", node.Real(19.683)),
		r("'text'", node.Text("text")),
		r("-19", node.Neg(node.Int(19))),
		r("--19",
			node.Neg(
				node.Neg(
					node.Int(19)))),
		r("---19",
			node.Neg(
				node.Neg(
					node.Neg(
						node.Int(19))))),
		r("19*683", node.Mul(node.Int(19), node.Int(683))),
		r("19/683", node.Div(node.Int(19), node.Int(683))),
		r("19%683", node.Mod(node.Int(19), node.Int(683))),
		r("19*683/-3",
			node.Div(
				node.Mul(node.Int(19), node.Int(683)),
				node.Neg(node.Int(3)))),
		r("19+683", node.Add(node.Int(19), node.Int(683))),
		r("19-683", node.Sub(node.Int(19), node.Int(683))),
		r("19--683", node.Sub(node.Int(19), node.Neg(node.Int(683)))),
		r("19*683/-3+83",
			node.Add(
				node.Div(
					node.Mul(node.Int(19), node.Int(683)),
					node.Neg(node.Int(3))),
				node.Int(83))),
	).exec(t, layer6)
}

func Test_layer7(t *testing.T) {
	rs(
		r("true", node.Bool(true)),
		r("false", node.Bool(false)),
		r("null", node.Null()),
		r("19683", node.Int(19683)),
		r("19.683", node.Real(19.683)),
		r("'text'", node.Text("text")),
		r("-19", node.Neg(node.Int(19))),
		r("--19",
			node.Neg(
				node.Neg(
					node.Int(19)))),
		r("---19",
			node.Neg(
				node.Neg(
					node.Neg(
						node.Int(19))))),
		r("19*683", node.Mul(node.Int(19), node.Int(683))),
		r("19/683", node.Div(node.Int(19), node.Int(683))),
		r("19%683", node.Mod(node.Int(19), node.Int(683))),
		r("19*683/-3",
			node.Div(
				node.Mul(node.Int(19), node.Int(683)),
				node.Neg(node.Int(3)))),
		r("19+683", node.Add(node.Int(19), node.Int(683))),
		r("19-683", node.Sub(node.Int(19), node.Int(683))),
		r("19--683", node.Sub(node.Int(19), node.Neg(node.Int(683)))),
		r("19*683/-3+83",
			node.Add(
				node.Div(
					node.Mul(node.Int(19), node.Int(683)),
					node.Neg(node.Int(3))),
				node.Int(83))),
		r("'hello'||'world'",
			node.Concat(
				node.Text("hello"),
				node.Text("world"))),
		r("19*683/-3+83||'рублей'",
			node.Concat(
				node.Add(
					node.Div(
						node.Mul(node.Int(19), node.Int(683)),
						node.Neg(node.Int(3))),
					node.Int(83)),
				node.Text("рублей"),
			)),
		r("'рублей'||19*683/-3+83",
			node.Concat(
				node.Text("рублей"),
				node.Add(
					node.Div(
						node.Mul(node.Int(19), node.Int(683)),
						node.Neg(node.Int(3))),
					node.Int(83)),
			)),
	).exec(t, layer7)
}

func Test_layer8(t *testing.T) {
	rs(
		r("true", node.Bool(true)),
		r("false", node.Bool(false)),
		r("null", node.Null()),
		r("19683", node.Int(19683)),
		r("19.683", node.Real(19.683)),
		r("'text'", node.Text("text")),
		r("-19", node.Neg(node.Int(19))),
		r("--19",
			node.Neg(
				node.Neg(
					node.Int(19)))),
		r("---19",
			node.Neg(
				node.Neg(
					node.Neg(
						node.Int(19))))),
		r("19*683", node.Mul(node.Int(19), node.Int(683))),
		r("19/683", node.Div(node.Int(19), node.Int(683))),
		r("19%683", node.Mod(node.Int(19), node.Int(683))),
		r("19*683/-3",
			node.Div(
				node.Mul(node.Int(19), node.Int(683)),
				node.Neg(node.Int(3)))),
		r("19+683", node.Add(node.Int(19), node.Int(683))),
		r("19-683", node.Sub(node.Int(19), node.Int(683))),
		r("19--683", node.Sub(node.Int(19), node.Neg(node.Int(683)))),
		r("19*683/-3+83",
			node.Add(
				node.Div(
					node.Mul(node.Int(19), node.Int(683)),
					node.Neg(node.Int(3))),
				node.Int(83))),
		r("'hello'||'world'",
			node.Concat(
				node.Text("hello"),
				node.Text("world"))),
		r("19*683/-3+83||'рублей'",
			node.Concat(
				node.Add(
					node.Div(
						node.Mul(node.Int(19), node.Int(683)),
						node.Neg(node.Int(3))),
					node.Int(83)),
				node.Text("рублей"),
			)),
		r("'рублей'||19*683/-3+83",
			node.Concat(
				node.Text("рублей"),
				node.Add(
					node.Div(
						node.Mul(node.Int(19), node.Int(683)),
						node.Neg(node.Int(3))),
					node.Int(83)),
			)),
		r("19==683", node.Eq(node.Int(19), node.Int(683))),
		r("19!=683", node.Neq(node.Int(19), node.Int(683))),
		r("19<683", node.Lt(node.Int(19), node.Int(683))),
		r("19<=683", node.Lte(node.Int(19), node.Int(683))),
		r("19>683", node.Gt(node.Int(19), node.Int(683))),
		r("19>=683", node.Gte(node.Int(19), node.Int(683))),
		r("19+83==68*-3",
			node.Eq(
				node.Add(node.Int(19), node.Int(83)),
				node.Mul(node.Int(68), node.Neg(node.Int(3))))),
		r("19==683!=true",
			node.Neq(
				node.Eq(node.Int(19), node.Int(683)),
				node.Bool(true))),
	).exec(t, layer8)
}

func Test_layer9(t *testing.T) {
	rs(
		r("true", node.Bool(true)),
		r("false", node.Bool(false)),
		r("null", node.Null()),
		r("19683", node.Int(19683)),
		r("19.683", node.Real(19.683)),
		r("'text'", node.Text("text")),
		r("-19", node.Neg(node.Int(19))),
		r("--19",
			node.Neg(
				node.Neg(
					node.Int(19)))),
		r("---19",
			node.Neg(
				node.Neg(
					node.Neg(
						node.Int(19))))),
		r("19*683", node.Mul(node.Int(19), node.Int(683))),
		r("19/683", node.Div(node.Int(19), node.Int(683))),
		r("19%683", node.Mod(node.Int(19), node.Int(683))),
		r("19*683/-3",
			node.Div(
				node.Mul(node.Int(19), node.Int(683)),
				node.Neg(node.Int(3)))),
		r("19+683", node.Add(node.Int(19), node.Int(683))),
		r("19-683", node.Sub(node.Int(19), node.Int(683))),
		r("19--683", node.Sub(node.Int(19), node.Neg(node.Int(683)))),
		r("19*683/-3+83",
			node.Add(
				node.Div(
					node.Mul(node.Int(19), node.Int(683)),
					node.Neg(node.Int(3))),
				node.Int(83))),
		r("'hello'||'world'",
			node.Concat(
				node.Text("hello"),
				node.Text("world"))),
		r("19*683/-3+83||'рублей'",
			node.Concat(
				node.Add(
					node.Div(
						node.Mul(node.Int(19), node.Int(683)),
						node.Neg(node.Int(3))),
					node.Int(83)),
				node.Text("рублей"),
			)),
		r("'рублей'||19*683/-3+83",
			node.Concat(
				node.Text("рублей"),
				node.Add(
					node.Div(
						node.Mul(node.Int(19), node.Int(683)),
						node.Neg(node.Int(3))),
					node.Int(83)),
			)),
		r("19==683", node.Eq(node.Int(19), node.Int(683))),
		r("19!=683", node.Neq(node.Int(19), node.Int(683))),
		r("19<683", node.Lt(node.Int(19), node.Int(683))),
		r("19<=683", node.Lte(node.Int(19), node.Int(683))),
		r("19>683", node.Gt(node.Int(19), node.Int(683))),
		r("19>=683", node.Gte(node.Int(19), node.Int(683))),
		r("19+83==68*-3",
			node.Eq(
				node.Add(node.Int(19), node.Int(83)),
				node.Mul(node.Int(68), node.Neg(node.Int(3))))),
		r("19==683!=true",
			node.Neq(
				node.Eq(node.Int(19), node.Int(683)),
				node.Bool(true))),
		r("not true", node.Not(node.Bool(true))),
		r("not not true",
			node.Not(
				node.Not(
					node.Bool(true)))),
		r("not not not true",
			node.Not(
				node.Not(
					node.Not(
						node.Bool(true))))),
		r("not 19+83==68*-3",
			node.Not(
				node.Eq(
					node.Add(node.Int(19), node.Int(83)),
					node.Mul(node.Int(68), node.Neg(node.Int(3)))))),
	).exec(t, layer9)
}

func Test_layer10(t *testing.T) {
	rs(
		r("true", node.Bool(true)),
		r("false", node.Bool(false)),
		r("null", node.Null()),
		r("19683", node.Int(19683)),
		r("19.683", node.Real(19.683)),
		r("'text'", node.Text("text")),
		r("-19", node.Neg(node.Int(19))),
		r("--19",
			node.Neg(
				node.Neg(
					node.Int(19)))),
		r("---19",
			node.Neg(
				node.Neg(
					node.Neg(
						node.Int(19))))),
		r("19*683", node.Mul(node.Int(19), node.Int(683))),
		r("19/683", node.Div(node.Int(19), node.Int(683))),
		r("19%683", node.Mod(node.Int(19), node.Int(683))),
		r("19*683/-3",
			node.Div(
				node.Mul(node.Int(19), node.Int(683)),
				node.Neg(node.Int(3)))),
		r("19+683", node.Add(node.Int(19), node.Int(683))),
		r("19-683", node.Sub(node.Int(19), node.Int(683))),
		r("19--683", node.Sub(node.Int(19), node.Neg(node.Int(683)))),
		r("19*683/-3+83",
			node.Add(
				node.Div(
					node.Mul(node.Int(19), node.Int(683)),
					node.Neg(node.Int(3))),
				node.Int(83))),
		r("'hello'||'world'",
			node.Concat(
				node.Text("hello"),
				node.Text("world"))),
		r("19*683/-3+83||'рублей'",
			node.Concat(
				node.Add(
					node.Div(
						node.Mul(node.Int(19), node.Int(683)),
						node.Neg(node.Int(3))),
					node.Int(83)),
				node.Text("рублей"),
			)),
		r("'рублей'||19*683/-3+83",
			node.Concat(
				node.Text("рублей"),
				node.Add(
					node.Div(
						node.Mul(node.Int(19), node.Int(683)),
						node.Neg(node.Int(3))),
					node.Int(83)),
			)),
		r("19==683", node.Eq(node.Int(19), node.Int(683))),
		r("19!=683", node.Neq(node.Int(19), node.Int(683))),
		r("19<683", node.Lt(node.Int(19), node.Int(683))),
		r("19<=683", node.Lte(node.Int(19), node.Int(683))),
		r("19>683", node.Gt(node.Int(19), node.Int(683))),
		r("19>=683", node.Gte(node.Int(19), node.Int(683))),
		r("19+83==68*-3",
			node.Eq(
				node.Add(node.Int(19), node.Int(83)),
				node.Mul(node.Int(68), node.Neg(node.Int(3))))),
		r("19==683!=true",
			node.Neq(
				node.Eq(node.Int(19), node.Int(683)),
				node.Bool(true))),
		r("not true", node.Not(node.Bool(true))),
		r("not not true",
			node.Not(
				node.Not(
					node.Bool(true)))),
		r("not not not true",
			node.Not(
				node.Not(
					node.Not(
						node.Bool(true))))),
		r("not 19+83==68*-3",
			node.Not(
				node.Eq(
					node.Add(node.Int(19), node.Int(83)),
					node.Mul(node.Int(68), node.Neg(node.Int(3)))))),
		r("true and false", node.And(node.Bool(true), node.Bool(false))),
		r("not 19+83==68*-3 and true",
			node.And(
				node.Not(
					node.Eq(
						node.Add(node.Int(19), node.Int(83)),
						node.Mul(node.Int(68), node.Neg(node.Int(3))))),
				node.Bool(true))),
	).exec(t, layer10)
}

func Test_layer11(t *testing.T) {
	rs(
		r("true", node.Bool(true)),
		r("false", node.Bool(false)),
		r("null", node.Null()),
		r("19683", node.Int(19683)),
		r("19.683", node.Real(19.683)),
		r("'text'", node.Text("text")),
		r("-19", node.Neg(node.Int(19))),
		r("--19",
			node.Neg(
				node.Neg(
					node.Int(19)))),
		r("---19",
			node.Neg(
				node.Neg(
					node.Neg(
						node.Int(19))))),
		r("19*683", node.Mul(node.Int(19), node.Int(683))),
		r("19/683", node.Div(node.Int(19), node.Int(683))),
		r("19%683", node.Mod(node.Int(19), node.Int(683))),
		r("19*683/-3",
			node.Div(
				node.Mul(node.Int(19), node.Int(683)),
				node.Neg(node.Int(3)))),
		r("19+683", node.Add(node.Int(19), node.Int(683))),
		r("19-683", node.Sub(node.Int(19), node.Int(683))),
		r("19--683", node.Sub(node.Int(19), node.Neg(node.Int(683)))),
		r("19*683/-3+83",
			node.Add(
				node.Div(
					node.Mul(node.Int(19), node.Int(683)),
					node.Neg(node.Int(3))),
				node.Int(83))),
		r("'hello'||'world'",
			node.Concat(
				node.Text("hello"),
				node.Text("world"))),
		r("19*683/-3+83||'рублей'",
			node.Concat(
				node.Add(
					node.Div(
						node.Mul(node.Int(19), node.Int(683)),
						node.Neg(node.Int(3))),
					node.Int(83)),
				node.Text("рублей"),
			)),
		r("'рублей'||19*683/-3+83",
			node.Concat(
				node.Text("рублей"),
				node.Add(
					node.Div(
						node.Mul(node.Int(19), node.Int(683)),
						node.Neg(node.Int(3))),
					node.Int(83)),
			)),
		r("19==683", node.Eq(node.Int(19), node.Int(683))),
		r("19!=683", node.Neq(node.Int(19), node.Int(683))),
		r("19<683", node.Lt(node.Int(19), node.Int(683))),
		r("19<=683", node.Lte(node.Int(19), node.Int(683))),
		r("19>683", node.Gt(node.Int(19), node.Int(683))),
		r("19>=683", node.Gte(node.Int(19), node.Int(683))),
		r("19+83==68*-3",
			node.Eq(
				node.Add(node.Int(19), node.Int(83)),
				node.Mul(node.Int(68), node.Neg(node.Int(3))))),
		r("19==683!=true",
			node.Neq(
				node.Eq(node.Int(19), node.Int(683)),
				node.Bool(true))),
		r("not true", node.Not(node.Bool(true))),
		r("not not true",
			node.Not(
				node.Not(
					node.Bool(true)))),
		r("not not not true",
			node.Not(
				node.Not(
					node.Not(
						node.Bool(true))))),
		r("not 19+83==68*-3",
			node.Not(
				node.Eq(
					node.Add(node.Int(19), node.Int(83)),
					node.Mul(node.Int(68), node.Neg(node.Int(3)))))),
		r("true and false", node.And(node.Bool(true), node.Bool(false))),
		r("not 19+83==68*-3 and true",
			node.And(
				node.Not(
					node.Eq(
						node.Add(node.Int(19), node.Int(83)),
						node.Mul(node.Int(68), node.Neg(node.Int(3))))),
				node.Bool(true))),
		r("true or false", node.Or(node.Bool(true), node.Bool(false))),
		r("not 19+83==68*-3 or true",
			node.Or(
				node.Not(
					node.Eq(
						node.Add(node.Int(19), node.Int(83)),
						node.Mul(node.Int(68), node.Neg(node.Int(3))))),
				node.Bool(true))),
		r("true and false or true and true",
			node.Or(
				node.And(node.Bool(true), node.Bool(false)),
				node.And(node.Bool(true), node.Bool(true)),
			)),
	).exec(t, layer11)
}

func Test_layer2(t *testing.T) {
	rs(
		r("true", node.Bool(true)),
		r("false", node.Bool(false)),
		r("null", node.Null()),
		r("19683", node.Int(19683)),
		r("19.683", node.Real(19.683)),
		r("'text'", node.Text("text")),
		r("(19*683)", node.Mul(node.Int(19), node.Int(683))),
		r("(19/683)", node.Div(node.Int(19), node.Int(683))),
		r("(19%683)", node.Mod(node.Int(19), node.Int(683))),
	).exec(t, layer2)
}

func Test_layer12(t *testing.T) {
	rs(
		r("true", node.Bool(true)),
		r("false", node.Bool(false)),
		r("null", node.Null()),
		r("19683", node.Int(19683)),
		r("19.683", node.Real(19.683)),
		r("'text'", node.Text("text")),
		r("-19", node.Neg(node.Int(19))),
		r("--19",
			node.Neg(
				node.Neg(
					node.Int(19)))),
		r("---19",
			node.Neg(
				node.Neg(
					node.Neg(
						node.Int(19))))),
		r("19*683", node.Mul(node.Int(19), node.Int(683))),
		r("19/683", node.Div(node.Int(19), node.Int(683))),
		r("19%683", node.Mod(node.Int(19), node.Int(683))),
		r("19*683/-3",
			node.Div(
				node.Mul(node.Int(19), node.Int(683)),
				node.Neg(node.Int(3)))),
		r("19+683", node.Add(node.Int(19), node.Int(683))),
		r("19-683", node.Sub(node.Int(19), node.Int(683))),
		r("19--683", node.Sub(node.Int(19), node.Neg(node.Int(683)))),
		r("19*683/-3+83",
			node.Add(
				node.Div(
					node.Mul(node.Int(19), node.Int(683)),
					node.Neg(node.Int(3))),
				node.Int(83))),
		r("'hello'||'world'",
			node.Concat(
				node.Text("hello"),
				node.Text("world"))),
		r("19*683/-3+83||'рублей'",
			node.Concat(
				node.Add(
					node.Div(
						node.Mul(node.Int(19), node.Int(683)),
						node.Neg(node.Int(3))),
					node.Int(83)),
				node.Text("рублей"),
			)),
		r("'рублей'||19*683/-3+83",
			node.Concat(
				node.Text("рублей"),
				node.Add(
					node.Div(
						node.Mul(node.Int(19), node.Int(683)),
						node.Neg(node.Int(3))),
					node.Int(83)),
			)),
		r("19==683", node.Eq(node.Int(19), node.Int(683))),
		r("19!=683", node.Neq(node.Int(19), node.Int(683))),
		r("19<683", node.Lt(node.Int(19), node.Int(683))),
		r("19<=683", node.Lte(node.Int(19), node.Int(683))),
		r("19>683", node.Gt(node.Int(19), node.Int(683))),
		r("19>=683", node.Gte(node.Int(19), node.Int(683))),
		r("19+83==68*-3",
			node.Eq(
				node.Add(node.Int(19), node.Int(83)),
				node.Mul(node.Int(68), node.Neg(node.Int(3))))),
		r("19==683!=true",
			node.Neq(
				node.Eq(node.Int(19), node.Int(683)),
				node.Bool(true))),
		r("not true", node.Not(node.Bool(true))),
		r("not not true",
			node.Not(
				node.Not(
					node.Bool(true)))),
		r("not not not true",
			node.Not(
				node.Not(
					node.Not(
						node.Bool(true))))),
		r("not 19+83==68*-3",
			node.Not(
				node.Eq(
					node.Add(node.Int(19), node.Int(83)),
					node.Mul(node.Int(68), node.Neg(node.Int(3)))))),
		r("true and false", node.And(node.Bool(true), node.Bool(false))),
		r("not 19+83==68*-3 and true",
			node.And(
				node.Not(
					node.Eq(
						node.Add(node.Int(19), node.Int(83)),
						node.Mul(node.Int(68), node.Neg(node.Int(3))))),
				node.Bool(true))),
		r("true or false", node.Or(node.Bool(true), node.Bool(false))),
		r("not 19+83==68*-3 or true",
			node.Or(
				node.Not(
					node.Eq(
						node.Add(node.Int(19), node.Int(83)),
						node.Mul(node.Int(68), node.Neg(node.Int(3))))),
				node.Bool(true))),
		r("true and false or true and true",
			node.Or(
				node.And(node.Bool(true), node.Bool(false)),
				node.And(node.Bool(true), node.Bool(true)),
			)),
		r("name='sergey'", node.Set(node.Ident("name"), node.Text("sergey"))),
		r("users[0]='sergey'",
			node.Set(node.ElByIndex(node.Ident("users"), node.Int(0)), node.Text("sergey"))),
	).exec(t, layer12)
}

func Test_layer13(t *testing.T) {
	rs(
		r("true", node.Bool(true)),
		r("false", node.Bool(false)),
		r("null", node.Null()),
		r("19683", node.Int(19683)),
		r("19.683", node.Real(19.683)),
		r("'text'", node.Text("text")),
		r("-19", node.Neg(node.Int(19))),
		r("--19",
			node.Neg(
				node.Neg(
					node.Int(19)))),
		r("---19",
			node.Neg(
				node.Neg(
					node.Neg(
						node.Int(19))))),
		r("19*683", node.Mul(node.Int(19), node.Int(683))),
		r("19/683", node.Div(node.Int(19), node.Int(683))),
		r("19%683", node.Mod(node.Int(19), node.Int(683))),
		r("19*683/-3",
			node.Div(
				node.Mul(node.Int(19), node.Int(683)),
				node.Neg(node.Int(3)))),
		r("19+683", node.Add(node.Int(19), node.Int(683))),
		r("19-683", node.Sub(node.Int(19), node.Int(683))),
		r("19--683", node.Sub(node.Int(19), node.Neg(node.Int(683)))),
		r("19*683/-3+83",
			node.Add(
				node.Div(
					node.Mul(node.Int(19), node.Int(683)),
					node.Neg(node.Int(3))),
				node.Int(83))),
		r("'hello'||'world'",
			node.Concat(
				node.Text("hello"),
				node.Text("world"))),
		r("19*683/-3+83||'рублей'",
			node.Concat(
				node.Add(
					node.Div(
						node.Mul(node.Int(19), node.Int(683)),
						node.Neg(node.Int(3))),
					node.Int(83)),
				node.Text("рублей"),
			)),
		r("'рублей'||19*683/-3+83",
			node.Concat(
				node.Text("рублей"),
				node.Add(
					node.Div(
						node.Mul(node.Int(19), node.Int(683)),
						node.Neg(node.Int(3))),
					node.Int(83)),
			)),
		r("19==683", node.Eq(node.Int(19), node.Int(683))),
		r("19!=683", node.Neq(node.Int(19), node.Int(683))),
		r("19<683", node.Lt(node.Int(19), node.Int(683))),
		r("19<=683", node.Lte(node.Int(19), node.Int(683))),
		r("19>683", node.Gt(node.Int(19), node.Int(683))),
		r("19>=683", node.Gte(node.Int(19), node.Int(683))),
		r("19+83==68*-3",
			node.Eq(
				node.Add(node.Int(19), node.Int(83)),
				node.Mul(node.Int(68), node.Neg(node.Int(3))))),
		r("19==683!=true",
			node.Neq(
				node.Eq(node.Int(19), node.Int(683)),
				node.Bool(true))),
		r("not true", node.Not(node.Bool(true))),
		r("not not true",
			node.Not(
				node.Not(
					node.Bool(true)))),
		r("not not not true",
			node.Not(
				node.Not(
					node.Not(
						node.Bool(true))))),
		r("not 19+83==68*-3",
			node.Not(
				node.Eq(
					node.Add(node.Int(19), node.Int(83)),
					node.Mul(node.Int(68), node.Neg(node.Int(3)))))),
		r("true and false", node.And(node.Bool(true), node.Bool(false))),
		r("not 19+83==68*-3 and true",
			node.And(
				node.Not(
					node.Eq(
						node.Add(node.Int(19), node.Int(83)),
						node.Mul(node.Int(68), node.Neg(node.Int(3))))),
				node.Bool(true))),
		r("true or false", node.Or(node.Bool(true), node.Bool(false))),
		r("not 19+83==68*-3 or true",
			node.Or(
				node.Not(
					node.Eq(
						node.Add(node.Int(19), node.Int(83)),
						node.Mul(node.Int(68), node.Neg(node.Int(3))))),
				node.Bool(true))),
		r("true and false or true and true",
			node.Or(
				node.And(node.Bool(true), node.Bool(false)),
				node.And(node.Bool(true), node.Bool(true)),
			)),
		r("name='sergey'", node.Set(node.Ident("name"), node.Text("sergey"))),
		r("users[0]='sergey'",
			node.Set(node.ElByIndex(node.Ident("users"), node.Int(0)), node.Text("sergey"))),
		r("if 19<68 {68;19}",
			node.If(node.Branch{node.Lt(node.Int(19), node.Int(68)), node.Block(
				node.Int(68),
				node.Int(19),
			)}),
		),
		r("if 19<68 {68;19} else {'sergey'}",
			node.If(node.Branch{node.Lt(node.Int(19), node.Int(68)), node.Block(
				node.Int(68),
				node.Int(19),
			)}, node.Branch{node.Bool(true), node.Block(node.Text("sergey"))}),
		),
	).exec(t, layer13)
}

func Test_layer14(t *testing.T) {
	rs(
		r("true", node.Bool(true)),
		r("false", node.Bool(false)),
		r("null", node.Null()),
		r("19683", node.Int(19683)),
		r("19.683", node.Real(19.683)),
		r("'text'", node.Text("text")),
		r("-19", node.Neg(node.Int(19))),
		r("--19",
			node.Neg(
				node.Neg(
					node.Int(19)))),
		r("---19",
			node.Neg(
				node.Neg(
					node.Neg(
						node.Int(19))))),
		r("19*683", node.Mul(node.Int(19), node.Int(683))),
		r("19/683", node.Div(node.Int(19), node.Int(683))),
		r("19%683", node.Mod(node.Int(19), node.Int(683))),
		r("19*683/-3",
			node.Div(
				node.Mul(node.Int(19), node.Int(683)),
				node.Neg(node.Int(3)))),
		r("19+683", node.Add(node.Int(19), node.Int(683))),
		r("19-683", node.Sub(node.Int(19), node.Int(683))),
		r("19--683", node.Sub(node.Int(19), node.Neg(node.Int(683)))),
		r("19*683/-3+83",
			node.Add(
				node.Div(
					node.Mul(node.Int(19), node.Int(683)),
					node.Neg(node.Int(3))),
				node.Int(83))),
		r("'hello'||'world'",
			node.Concat(
				node.Text("hello"),
				node.Text("world"))),
		r("19*683/-3+83||'рублей'",
			node.Concat(
				node.Add(
					node.Div(
						node.Mul(node.Int(19), node.Int(683)),
						node.Neg(node.Int(3))),
					node.Int(83)),
				node.Text("рублей"),
			)),
		r("'рублей'||19*683/-3+83",
			node.Concat(
				node.Text("рублей"),
				node.Add(
					node.Div(
						node.Mul(node.Int(19), node.Int(683)),
						node.Neg(node.Int(3))),
					node.Int(83)),
			)),
		r("19==683", node.Eq(node.Int(19), node.Int(683))),
		r("19!=683", node.Neq(node.Int(19), node.Int(683))),
		r("19<683", node.Lt(node.Int(19), node.Int(683))),
		r("19<=683", node.Lte(node.Int(19), node.Int(683))),
		r("19>683", node.Gt(node.Int(19), node.Int(683))),
		r("19>=683", node.Gte(node.Int(19), node.Int(683))),
		r("19+83==68*-3",
			node.Eq(
				node.Add(node.Int(19), node.Int(83)),
				node.Mul(node.Int(68), node.Neg(node.Int(3))))),
		r("19==683!=true",
			node.Neq(
				node.Eq(node.Int(19), node.Int(683)),
				node.Bool(true))),
		r("not true", node.Not(node.Bool(true))),
		r("not not true",
			node.Not(
				node.Not(
					node.Bool(true)))),
		r("not not not true",
			node.Not(
				node.Not(
					node.Not(
						node.Bool(true))))),
		r("not 19+83==68*-3",
			node.Not(
				node.Eq(
					node.Add(node.Int(19), node.Int(83)),
					node.Mul(node.Int(68), node.Neg(node.Int(3)))))),
		r("true and false", node.And(node.Bool(true), node.Bool(false))),
		r("not 19+83==68*-3 and true",
			node.And(
				node.Not(
					node.Eq(
						node.Add(node.Int(19), node.Int(83)),
						node.Mul(node.Int(68), node.Neg(node.Int(3))))),
				node.Bool(true))),
		r("true or false", node.Or(node.Bool(true), node.Bool(false))),
		r("not 19+83==68*-3 or true",
			node.Or(
				node.Not(
					node.Eq(
						node.Add(node.Int(19), node.Int(83)),
						node.Mul(node.Int(68), node.Neg(node.Int(3))))),
				node.Bool(true))),
		r("true and false or true and true",
			node.Or(
				node.And(node.Bool(true), node.Bool(false)),
				node.And(node.Bool(true), node.Bool(true)),
			)),
		r("name='sergey'", node.Set(node.Ident("name"), node.Text("sergey"))),
		r("users[0]='sergey'",
			node.Set(node.ElByIndex(node.Ident("users"), node.Int(0)), node.Text("sergey"))),

		r("for i in a {i}",
			node.For([]node.Node{node.Ident("i")}, node.Ident("a"), node.Block(
				node.Ident("i"),
			)),
		),
		r("for i,j in a {i;j}",
			node.For([]node.Node{node.Ident("i"), node.Ident("j")}, node.Ident("a"), node.Block(
				node.Ident("i"),
				node.Ident("j"),
			)),
		),
	).exec(t, layer14)
}

func Test_Parse(t *testing.T) {
	rs(
		r("81+16*16",
			node.Block(
				node.Add(
					node.Int(81),
					node.Mul(node.Int(16), node.Int(16)),
				),
			),
		),
		r("(81+16)*16",
			node.Block(
				node.Mul(
					node.Add(node.Int(81), node.Int(16)),
					node.Int(16),
				),
			),
		),
		r("(81+16)*-(64/16)",
			node.Block(
				node.Mul(
					node.Add(node.Int(81), node.Int(16)),
					node.Neg(node.Div(node.Int(64), node.Int(16))),
				),
			),
		),
		r("(81+16)*-(64/16)<16%6-1",
			node.Block(
				node.Lt(
					node.Mul(
						node.Add(node.Int(81), node.Int(16)),
						node.Neg(node.Div(node.Int(64), node.Int(16))),
					),
					node.Sub(
						node.Mod(node.Int(16), node.Int(6)),
						node.Int(1),
					),
				),
			),
		),
		r("(81+16)*-(64/16)<16%6-1 and not true",
			node.Block(
				node.And(
					node.Lt(
						node.Mul(
							node.Add(node.Int(81), node.Int(16)),
							node.Neg(node.Div(node.Int(64), node.Int(16))),
						),
						node.Sub(
							node.Mod(node.Int(16), node.Int(6)),
							node.Int(1),
						),
					),
					node.Not(node.Bool(true)),
				),
			),
		),
		r("(81+16)*-(64/16)<16%6-1 and not true or not not true",
			node.Block(
				node.Or(
					node.And(
						node.Lt(
							node.Mul(
								node.Add(node.Int(81), node.Int(16)),
								node.Neg(node.Div(node.Int(64), node.Int(16))),
							),
							node.Sub(
								node.Mod(node.Int(16), node.Int(6)),
								node.Int(1),
							),
						),
						node.Not(node.Bool(true)),
					),
					node.Not(node.Not(node.Bool(true))),
				),
			),
		),
		r("81+16*16;16*81",
			node.Block(
				node.Add(
					node.Int(81),
					node.Mul(node.Int(16), node.Int(16)),
				),
				node.Mul(
					node.Int(16),
					node.Int(81),
				),
			),
		),
		r("81+16*16;16*81;",
			node.Block(
				node.Add(
					node.Int(81),
					node.Mul(node.Int(16), node.Int(16)),
				),
				node.Mul(
					node.Int(16),
					node.Int(81),
				),
			),
		),
	).exec(t, parse)
}

func Test_parseBody(t *testing.T) {
	rs(
		r("{81+16*16}",
			node.Block(
				node.Add(
					node.Int(81),
					node.Mul(node.Int(16), node.Int(16)),
				),
			),
		),
		r("{(81+16)*16}",
			node.Block(
				node.Mul(
					node.Add(node.Int(81), node.Int(16)),
					node.Int(16),
				),
			),
		),
		r("{(81+16)*-(64/16)}",
			node.Block(
				node.Mul(
					node.Add(node.Int(81), node.Int(16)),
					node.Neg(node.Div(node.Int(64), node.Int(16))),
				),
			),
		),
		r("{(81+16)*-(64/16)<16%6-1}",
			node.Block(
				node.Lt(
					node.Mul(
						node.Add(node.Int(81), node.Int(16)),
						node.Neg(node.Div(node.Int(64), node.Int(16))),
					),
					node.Sub(
						node.Mod(node.Int(16), node.Int(6)),
						node.Int(1),
					),
				),
			),
		),
		r("{(81+16)*-(64/16)<16%6-1 and not true}",
			node.Block(
				node.And(
					node.Lt(
						node.Mul(
							node.Add(node.Int(81), node.Int(16)),
							node.Neg(node.Div(node.Int(64), node.Int(16))),
						),
						node.Sub(
							node.Mod(node.Int(16), node.Int(6)),
							node.Int(1),
						),
					),
					node.Not(node.Bool(true)),
				),
			),
		),
		r("{(81+16)*-(64/16)<16%6-1 and not true or not not true}",
			node.Block(
				node.Or(
					node.And(
						node.Lt(
							node.Mul(
								node.Add(node.Int(81), node.Int(16)),
								node.Neg(node.Div(node.Int(64), node.Int(16))),
							),
							node.Sub(
								node.Mod(node.Int(16), node.Int(6)),
								node.Int(1),
							),
						),
						node.Not(node.Bool(true)),
					),
					node.Not(node.Not(node.Bool(true))),
				),
			),
		),
		r("{81+16*16;16*81}",
			node.Block(
				node.Add(
					node.Int(81),
					node.Mul(node.Int(16), node.Int(16)),
				),
				node.Mul(
					node.Int(16),
					node.Int(81),
				),
			),
		),
		r("{81+16*16;16*81;}",
			node.Block(
				node.Add(
					node.Int(81),
					node.Mul(node.Int(16), node.Int(16)),
				),
				node.Mul(
					node.Int(16),
					node.Int(81),
				),
			),
		),
	).exec(t, parseBody)
}

func Test_Parse_errs(t *testing.T) {
	//все тесты должны пройти с ошибками
	tests := []struct {
		data     string
		expected error
	}{
		{"(81+16", unexpected(token.New(token.EOF, pos.NewWithStart(1, 7)))},
		{"81+16*", unexpected(token.New(token.EOF, pos.NewWithStart(1, 7)))},
		{"-", unexpected(token.New(token.EOF, pos.NewWithStart(1, 2)))},
		{"16 true", unexpected(token.New(token.True, pos.NewWithStart(1, 4)))},
		{")", unexpected(token.New(token.RParen, pos.NewWithStart(1, 1)))},
	}

	for i, test := range tests {
		p, err := prepare(test.data)
		if err != nil {
			t.Error(err.Error())
			continue
		}

		_, err = p.Parse()
		if err.Error() != test.expected.Error() {
			t.Errorf("%d: ожидалось %q, получено %q",
				i, test.expected.Error(), err.Error())
		}
	}
}
