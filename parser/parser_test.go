package parser

import (
	"fmt"
	"github.com/suprunchuksergey/dpl/lexer"
	"github.com/suprunchuksergey/dpl/node"
	"github.com/suprunchuksergey/dpl/pos"
	"github.com/suprunchuksergey/dpl/token"
	"testing"
)

const (
	_ uint8 = iota
	parse
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
)

func call(p *parser, l uint8) (node.Node, error) {
	switch l {
	case parse:
		return p.Parse()
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

	if node.DeepEqual(r.expected, v) == false {
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
		r("true", node.True()),
		r("false", node.False()),
		r("null", node.Null()),
		r("19683", node.Int(19683)),
		r("19.683", node.Real(19.683)),
		r("'text'", node.Text("text")),
		r("fn", node.Ident("fn")),
		r("[]", node.Array([]node.Node{})),
		r("[true, 19683, 'text',]",
			node.Array([]node.Node{
				node.True(), node.Int(19683), node.Text("text"),
			})),
		r("[true, 19683, 'text', [true, 19683, 'text'],]",
			node.Array([]node.Node{
				node.True(), node.Int(19683), node.Text("text"),
				node.Array([]node.Node{node.True(), node.Int(19683), node.Text("text")}),
			})),
		r("{}", node.Map(node.Records{})),
		r(
			"{'text' : 19683}",
			node.Map(node.Records{
				node.NewRecord(node.Text("text"), node.Int(19683)),
			})),
		r(
			"{'text' : 19683,}",
			node.Map(node.Records{
				node.NewRecord(node.Text("text"), node.Int(19683)),
			})),
		r(
			"{'text'||19 : 19683, 3*9 : 19.683/683,}",
			node.Map(node.Records{
				node.NewRecord(
					node.Concat(node.Text("text"), node.Int(19)),
					node.Int(19683)),
				node.NewRecord(
					node.Mul(node.Int(3), node.Int(9)),
					node.Div(node.Real(19.683), node.Int(683))),
			})),
	).exec(t, layer1)
}

func Test_layer3(t *testing.T) {
	rs(
		r("true", node.True()),
		r("false", node.False()),
		r("null", node.Null()),
		r("19683", node.Int(19683)),
		r("19.683", node.Real(19.683)),
		r("'text'", node.Text("text")),
		r("[]", node.Array([]node.Node{})),
		r("[true, 19683, 'text',]",
			node.Array([]node.Node{
				node.True(), node.Int(19683), node.Text("text"),
			})),
		r("[true, 19683, 'text', [true, 19683, 'text'],]",
			node.Array([]node.Node{
				node.True(), node.Int(19683), node.Text("text"),
				node.Array([]node.Node{node.True(), node.Int(19683), node.Text("text")}),
			})),
		r("{}", node.Map(node.Records{})),
		r(
			"{'text' : 19683}",
			node.Map(node.Records{
				node.NewRecord(node.Text("text"), node.Int(19683)),
			})),
		r(
			"{'text' : 19683,}",
			node.Map(node.Records{
				node.NewRecord(node.Text("text"), node.Int(19683)),
			})),
		r(
			"{'text'||19 : 19683, 3*9 : 19.683/683,}",
			node.Map(node.Records{
				node.NewRecord(
					node.Concat(node.Text("text"), node.Int(19)),
					node.Int(19683)),
				node.NewRecord(
					node.Mul(node.Int(3), node.Int(9)),
					node.Div(node.Real(19.683), node.Int(683))),
			})),
		r("[683, 9][1]", node.IndexAccess(
			node.Array([]node.Node{node.Int(683), node.Int(9)}),
			node.Int(1),
		)),
		r("[683, 9][1][3+9]", node.IndexAccess(
			node.IndexAccess(
				node.Array([]node.Node{node.Int(683), node.Int(9)}),
				node.Int(1),
			), node.Add(node.Int(3), node.Int(9)),
		)),
	).exec(t, layer3)
}

func Test_layer4(t *testing.T) {
	rs(
		r("true", node.True()),
		r("false", node.False()),
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
		r("true", node.True()),
		r("false", node.False()),
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
		r("19%683", node.Rem(node.Int(19), node.Int(683))),
		r("19*683/-3",
			node.Div(
				node.Mul(node.Int(19), node.Int(683)),
				node.Neg(node.Int(3)))),
	).exec(t, layer5)
}

func Test_layer6(t *testing.T) {
	rs(
		r("true", node.True()),
		r("false", node.False()),
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
		r("19%683", node.Rem(node.Int(19), node.Int(683))),
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
		r("true", node.True()),
		r("false", node.False()),
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
		r("19%683", node.Rem(node.Int(19), node.Int(683))),
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
		r("true", node.True()),
		r("false", node.False()),
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
		r("19%683", node.Rem(node.Int(19), node.Int(683))),
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
		r("19=683", node.Eq(node.Int(19), node.Int(683))),
		r("19!=683", node.Neq(node.Int(19), node.Int(683))),
		r("19<683", node.Lt(node.Int(19), node.Int(683))),
		r("19<=683", node.Lte(node.Int(19), node.Int(683))),
		r("19>683", node.Gt(node.Int(19), node.Int(683))),
		r("19>=683", node.Gte(node.Int(19), node.Int(683))),
		r("19+83=68*-3",
			node.Eq(
				node.Add(node.Int(19), node.Int(83)),
				node.Mul(node.Int(68), node.Neg(node.Int(3))))),
		r("19=683!=true",
			node.Neq(
				node.Eq(node.Int(19), node.Int(683)),
				node.True())),
	).exec(t, layer8)
}

func Test_layer9(t *testing.T) {
	rs(
		r("true", node.True()),
		r("false", node.False()),
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
		r("19%683", node.Rem(node.Int(19), node.Int(683))),
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
		r("19=683", node.Eq(node.Int(19), node.Int(683))),
		r("19!=683", node.Neq(node.Int(19), node.Int(683))),
		r("19<683", node.Lt(node.Int(19), node.Int(683))),
		r("19<=683", node.Lte(node.Int(19), node.Int(683))),
		r("19>683", node.Gt(node.Int(19), node.Int(683))),
		r("19>=683", node.Gte(node.Int(19), node.Int(683))),
		r("19+83=68*-3",
			node.Eq(
				node.Add(node.Int(19), node.Int(83)),
				node.Mul(node.Int(68), node.Neg(node.Int(3))))),
		r("19=683!=true",
			node.Neq(
				node.Eq(node.Int(19), node.Int(683)),
				node.True())),
		r("not true", node.Not(node.True())),
		r("not not true",
			node.Not(
				node.Not(
					node.True()))),
		r("not not not true",
			node.Not(
				node.Not(
					node.Not(
						node.True())))),
		r("not 19+83=68*-3",
			node.Not(
				node.Eq(
					node.Add(node.Int(19), node.Int(83)),
					node.Mul(node.Int(68), node.Neg(node.Int(3)))))),
	).exec(t, layer9)
}

func Test_layer10(t *testing.T) {
	rs(
		r("true", node.True()),
		r("false", node.False()),
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
		r("19%683", node.Rem(node.Int(19), node.Int(683))),
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
		r("19=683", node.Eq(node.Int(19), node.Int(683))),
		r("19!=683", node.Neq(node.Int(19), node.Int(683))),
		r("19<683", node.Lt(node.Int(19), node.Int(683))),
		r("19<=683", node.Lte(node.Int(19), node.Int(683))),
		r("19>683", node.Gt(node.Int(19), node.Int(683))),
		r("19>=683", node.Gte(node.Int(19), node.Int(683))),
		r("19+83=68*-3",
			node.Eq(
				node.Add(node.Int(19), node.Int(83)),
				node.Mul(node.Int(68), node.Neg(node.Int(3))))),
		r("19=683!=true",
			node.Neq(
				node.Eq(node.Int(19), node.Int(683)),
				node.True())),
		r("not true", node.Not(node.True())),
		r("not not true",
			node.Not(
				node.Not(
					node.True()))),
		r("not not not true",
			node.Not(
				node.Not(
					node.Not(
						node.True())))),
		r("not 19+83=68*-3",
			node.Not(
				node.Eq(
					node.Add(node.Int(19), node.Int(83)),
					node.Mul(node.Int(68), node.Neg(node.Int(3)))))),
		r("true and false", node.And(node.True(), node.False())),
		r("not 19+83=68*-3 and true",
			node.And(
				node.Not(
					node.Eq(
						node.Add(node.Int(19), node.Int(83)),
						node.Mul(node.Int(68), node.Neg(node.Int(3))))),
				node.True())),
	).exec(t, layer10)
}

func Test_layer11(t *testing.T) {
	rs(
		r("true", node.True()),
		r("false", node.False()),
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
		r("19%683", node.Rem(node.Int(19), node.Int(683))),
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
		r("19=683", node.Eq(node.Int(19), node.Int(683))),
		r("19!=683", node.Neq(node.Int(19), node.Int(683))),
		r("19<683", node.Lt(node.Int(19), node.Int(683))),
		r("19<=683", node.Lte(node.Int(19), node.Int(683))),
		r("19>683", node.Gt(node.Int(19), node.Int(683))),
		r("19>=683", node.Gte(node.Int(19), node.Int(683))),
		r("19+83=68*-3",
			node.Eq(
				node.Add(node.Int(19), node.Int(83)),
				node.Mul(node.Int(68), node.Neg(node.Int(3))))),
		r("19=683!=true",
			node.Neq(
				node.Eq(node.Int(19), node.Int(683)),
				node.True())),
		r("not true", node.Not(node.True())),
		r("not not true",
			node.Not(
				node.Not(
					node.True()))),
		r("not not not true",
			node.Not(
				node.Not(
					node.Not(
						node.True())))),
		r("not 19+83=68*-3",
			node.Not(
				node.Eq(
					node.Add(node.Int(19), node.Int(83)),
					node.Mul(node.Int(68), node.Neg(node.Int(3)))))),
		r("true and false", node.And(node.True(), node.False())),
		r("not 19+83=68*-3 and true",
			node.And(
				node.Not(
					node.Eq(
						node.Add(node.Int(19), node.Int(83)),
						node.Mul(node.Int(68), node.Neg(node.Int(3))))),
				node.True())),
		r("true or false", node.Or(node.True(), node.False())),
		r("not 19+83=68*-3 or true",
			node.Or(
				node.Not(
					node.Eq(
						node.Add(node.Int(19), node.Int(83)),
						node.Mul(node.Int(68), node.Neg(node.Int(3))))),
				node.True())),
		r("true and false or true and true",
			node.Or(
				node.And(node.True(), node.False()),
				node.And(node.True(), node.True()),
			)),
	).exec(t, layer11)
}

func Test_layer2(t *testing.T) {
	rs(
		r("true", node.True()),
		r("false", node.False()),
		r("null", node.Null()),
		r("19683", node.Int(19683)),
		r("19.683", node.Real(19.683)),
		r("'text'", node.Text("text")),
		r("(19*683)", node.Mul(node.Int(19), node.Int(683))),
		r("(19/683)", node.Div(node.Int(19), node.Int(683))),
		r("(19%683)", node.Rem(node.Int(19), node.Int(683))),
	).exec(t, layer2)
}

func Test_Parse(t *testing.T) {
	rs(
		r("81+16*16",
			node.Add(
				node.Int(81),
				node.Mul(node.Int(16), node.Int(16)),
			),
		),
		r("(81+16)*16",
			node.Mul(
				node.Add(node.Int(81), node.Int(16)),
				node.Int(16),
			),
		),
		r("(81+16)*-(64/16)",
			node.Mul(
				node.Add(node.Int(81), node.Int(16)),
				node.Neg(node.Div(node.Int(64), node.Int(16))),
			),
		),
		r("(81+16)*-(64/16)<16%6-1",
			node.Lt(
				node.Mul(
					node.Add(node.Int(81), node.Int(16)),
					node.Neg(node.Div(node.Int(64), node.Int(16))),
				),
				node.Sub(
					node.Rem(node.Int(16), node.Int(6)),
					node.Int(1),
				),
			),
		),
		r("(81+16)*-(64/16)<16%6-1 and not true",
			node.And(
				node.Lt(
					node.Mul(
						node.Add(node.Int(81), node.Int(16)),
						node.Neg(node.Div(node.Int(64), node.Int(16))),
					),
					node.Sub(
						node.Rem(node.Int(16), node.Int(6)),
						node.Int(1),
					),
				),
				node.Not(node.True()),
			),
		),
		r("(81+16)*-(64/16)<16%6-1 and not true or not not true",
			node.Or(
				node.And(
					node.Lt(
						node.Mul(
							node.Add(node.Int(81), node.Int(16)),
							node.Neg(node.Div(node.Int(64), node.Int(16))),
						),
						node.Sub(
							node.Rem(node.Int(16), node.Int(6)),
							node.Int(1),
						),
					),
					node.Not(node.True()),
				),
				node.Not(node.Not(node.True())),
			),
		),
	).exec(t, parse)
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
