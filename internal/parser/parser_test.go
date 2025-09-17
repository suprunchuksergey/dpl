package parser

import (
	"github.com/stretchr/testify/assert"
	"github.com/suprunchuksergey/dpl/internal/lexer"
	"github.com/suprunchuksergey/dpl/internal/node"
	"testing"
)

func Test_value(t *testing.T) {
	tests := []struct {
		data          string
		expectedValue node.Node
		expectedError error
	}{
		{"null", node.Null(), nil},
		{"false", node.Bool(false), nil},
		{"true", node.Bool(true), nil},
		{"2187", node.Int(2187), nil},
		{"2.187", node.Real(2.187), nil},
		{`"text"`, node.Text("text"), nil},

		{"[]", node.Array(), nil},
		{"[2187]", node.Array(node.Int(2187)), nil},
		{"[2187,]", node.Array(node.Int(2187)), nil},
		{`[2187,"text"]`, node.Array(node.Int(2187), node.Text("text")), nil},
		{`[2187,"text",[]]`,
			node.Array(
				node.Int(2187),
				node.Text("text"),
				node.Array(),
			), nil},
		{`[2187,"text",[true]]`,
			node.Array(
				node.Int(2187),
				node.Text("text"),
				node.Array(node.Bool(true)),
			), nil},

		{"{}", node.Object(), nil},
		{`{"text": 2187}`, node.Object(
			node.KV{Key: node.Text("text"), Value: node.Int(2187)},
		), nil},
		{`{"text": 2187,}`, node.Object(
			node.KV{Key: node.Text("text"), Value: node.Int(2187)},
		), nil},
		{`{"text": 2187,"name": "сергей"}`, node.Object(
			node.KV{Key: node.Text("text"), Value: node.Int(2187)},
			node.KV{Key: node.Text("name"), Value: node.Text("сергей")},
		), nil},

		{"name", node.Ident("name"), nil},

		{"+", nil, unexpectedToken(lexer.NewToken(lexer.Add))},
		{"[2187", nil, unexpectedToken(lexer.NewToken(lexer.EOF))},
		{`{"text": 2187`, nil, unexpectedToken(lexer.NewToken(lexer.EOF))},
		{`{"text" 2187}`, nil, unexpectedToken(lexer.NewTokenWithValue(lexer.Int, "2187"))},
	}

	for _, test := range tests {
		tokens, err := lexer.Tokenize(test.data)
		assert.NoError(t, err)

		p := newParser(tokens)

		v, err := p.value()

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_paren(t *testing.T) {
	tests := []struct {
		data          string
		expectedValue node.Node
		expectedError error
	}{
		{"2187", node.Int(2187), nil},

		{"(2187)", node.Int(2187), nil},

		{"() -> {}", node.Function(node.Block()), nil},
		{"(name) -> {}", node.Function(node.Block(), node.Ident("name")), nil},
		{"(name, age) -> {}", node.Function(
			node.Block(),
			node.Ident("name"), node.Ident("age"),
		), nil},
		{"(name, age) -> {name}", node.Function(
			node.Block(node.Ident("name")),
			node.Ident("name"), node.Ident("age"),
		), nil},
		{"(name, age) -> {name;}", node.Function(
			node.Block(node.Ident("name")),
			node.Ident("name"), node.Ident("age"),
		), nil},
		{"(name, age) -> {name;age}", node.Function(
			node.Block(node.Ident("name"), node.Ident("age")),
			node.Ident("name"), node.Ident("age"),
		), nil},

		{"()", nil, unexpectedToken(lexer.NewToken(lexer.EOF))},
		{"(name, age)", nil, unexpectedToken(lexer.NewToken(lexer.EOF))},
		{"(name, age) ->", nil, unexpectedToken(lexer.NewToken(lexer.EOF))},
		{"(name, age) -> {name;", nil, unexpectedToken(lexer.NewToken(lexer.EOF))},
		{"(name, age) -> []", nil, unexpectedToken(lexer.NewToken(lexer.LBrack))},
	}

	for _, test := range tests {
		tokens, err := lexer.Tokenize(test.data)
		assert.NoError(t, err)

		p := newParser(tokens)

		v, err := p.paren()

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_elByIndex(t *testing.T) {
	tests := []struct {
		data          string
		expectedValue node.Node
		expectedError error
	}{
		{"2187", node.Int(2187), nil},
		{"(2187)", node.Int(2187), nil},
		{"() -> {}", node.Function(node.Block()), nil},

		{"array[1]", node.ElByIndex(node.Ident("array"), node.Int(1)), nil},
		{"array[1][3]",
			node.ElByIndex(
				node.ElByIndex(node.Ident("array"), node.Int(1)),
				node.Int(3),
			), nil},

		{"factorial()", node.Call(node.Ident("factorial")), nil},
		{"factorial(3)", node.Call(node.Ident("factorial"), node.Int(3)), nil},
		{"factorial(3)(1)",
			node.Call(
				node.Call(node.Ident("factorial"), node.Int(3)),
				node.Int(1),
			), nil},
		{"factorial(3)(1,8)",
			node.Call(
				node.Call(node.Ident("factorial"), node.Int(3)),
				node.Int(1), node.Int(8),
			), nil},

		{"array[1", nil, unexpectedToken(lexer.NewToken(lexer.EOF))},
		{"factorial(1", nil, unexpectedToken(lexer.NewToken(lexer.EOF))},
	}

	for _, test := range tests {
		tokens, err := lexer.Tokenize(test.data)
		assert.NoError(t, err)

		p := newParser(tokens)

		v, err := p.elByIndex()

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_neg(t *testing.T) {
	tests := []struct {
		data          string
		expectedValue node.Node
		expectedError error
	}{
		{"2187", node.Int(2187), nil},
		{"(2187)", node.Int(2187), nil},
		{"() -> {}", node.Function(node.Block()), nil},
		{"array[1]", node.ElByIndex(node.Ident("array"), node.Int(1)), nil},
		{"factorial()", node.Call(node.Ident("factorial")), nil},

		{"-2187", node.Neg(node.Int(2187)), nil},
		{"--2187", node.Neg(node.Neg(node.Int(2187))), nil},
	}

	for _, test := range tests {
		tokens, err := lexer.Tokenize(test.data)
		assert.NoError(t, err)

		p := newParser(tokens)

		v, err := p.neg()

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_mul(t *testing.T) {
	tests := []struct {
		data          string
		expectedValue node.Node
		expectedError error
	}{
		{"2187", node.Int(2187), nil},
		{"(2187)", node.Int(2187), nil},
		{"() -> {}", node.Function(node.Block()), nil},
		{"array[1]", node.ElByIndex(node.Ident("array"), node.Int(1)), nil},
		{"factorial()", node.Call(node.Ident("factorial")), nil},
		{"-2187", node.Neg(node.Int(2187)), nil},

		{"27*8", node.Mul(node.Int(27), node.Int(8)), nil},
		{"27/8", node.Div(node.Int(27), node.Int(8)), nil},
		{"27%8", node.Mod(node.Int(27), node.Int(8)), nil},

		{"27*8*16", node.Mul(node.Mul(node.Int(27), node.Int(8)), node.Int(16)), nil},
		{"27*8/16", node.Div(node.Mul(node.Int(27), node.Int(8)), node.Int(16)), nil},
		{"27*8/16%4", node.Mod(node.Div(node.Mul(node.Int(27), node.Int(8)), node.Int(16)), node.Int(4)), nil},
		{"27*8/16%-4", node.Mod(node.Div(node.Mul(node.Int(27), node.Int(8)), node.Int(16)), node.Neg(node.Int(4))), nil},

		{"27*(8*16)", node.Mul(node.Int(27), node.Mul(node.Int(8), node.Int(16))), nil},
	}

	for _, test := range tests {
		tokens, err := lexer.Tokenize(test.data)
		assert.NoError(t, err)

		p := newParser(tokens)

		v, err := p.mul()

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_add(t *testing.T) {
	tests := []struct {
		data          string
		expectedValue node.Node
		expectedError error
	}{
		{"2187", node.Int(2187), nil},
		{"(2187)", node.Int(2187), nil},
		{"() -> {}", node.Function(node.Block()), nil},
		{"array[1]", node.ElByIndex(node.Ident("array"), node.Int(1)), nil},
		{"factorial()", node.Call(node.Ident("factorial")), nil},
		{"-2187", node.Neg(node.Int(2187)), nil},
		{"27*8", node.Mul(node.Int(27), node.Int(8)), nil},
		{"27/8", node.Div(node.Int(27), node.Int(8)), nil},
		{"27%8", node.Mod(node.Int(27), node.Int(8)), nil},

		{"27+8", node.Add(node.Int(27), node.Int(8)), nil},
		{"27-8", node.Sub(node.Int(27), node.Int(8)), nil},

		{"27+8+16", node.Add(node.Add(node.Int(27), node.Int(8)), node.Int(16)), nil},
		{"27+8-16", node.Sub(node.Add(node.Int(27), node.Int(8)), node.Int(16)), nil},

		{"27-8*16", node.Sub(node.Int(27), node.Mul(node.Int(8), node.Int(16))), nil},
	}

	for _, test := range tests {
		tokens, err := lexer.Tokenize(test.data)
		assert.NoError(t, err)

		p := newParser(tokens)

		v, err := p.add()

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_concat(t *testing.T) {
	tests := []struct {
		data          string
		expectedValue node.Node
		expectedError error
	}{
		{"2187", node.Int(2187), nil},
		{"(2187)", node.Int(2187), nil},
		{"() -> {}", node.Function(node.Block()), nil},
		{"array[1]", node.ElByIndex(node.Ident("array"), node.Int(1)), nil},
		{"factorial()", node.Call(node.Ident("factorial")), nil},
		{"-2187", node.Neg(node.Int(2187)), nil},
		{"27*8", node.Mul(node.Int(27), node.Int(8)), nil},
		{"27/8", node.Div(node.Int(27), node.Int(8)), nil},
		{"27%8", node.Mod(node.Int(27), node.Int(8)), nil},
		{"27+8", node.Add(node.Int(27), node.Int(8)), nil},
		{"27-8", node.Sub(node.Int(27), node.Int(8)), nil},

		{`"привет "||"мир"`, node.Concat(node.Text("привет "), node.Text("мир")), nil},
		{`27+8||"рублей"`, node.Concat(node.Add(node.Int(27), node.Int(8)), node.Text("рублей")), nil},
	}

	for _, test := range tests {
		tokens, err := lexer.Tokenize(test.data)
		assert.NoError(t, err)

		p := newParser(tokens)

		v, err := p.concat()

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_eq(t *testing.T) {
	tests := []struct {
		data          string
		expectedValue node.Node
		expectedError error
	}{
		{"2187", node.Int(2187), nil},
		{"(2187)", node.Int(2187), nil},
		{"() -> {}", node.Function(node.Block()), nil},
		{"array[1]", node.ElByIndex(node.Ident("array"), node.Int(1)), nil},
		{"factorial()", node.Call(node.Ident("factorial")), nil},
		{"-2187", node.Neg(node.Int(2187)), nil},
		{"27*8", node.Mul(node.Int(27), node.Int(8)), nil},
		{"27/8", node.Div(node.Int(27), node.Int(8)), nil},
		{"27%8", node.Mod(node.Int(27), node.Int(8)), nil},
		{"27+8", node.Add(node.Int(27), node.Int(8)), nil},
		{"27-8", node.Sub(node.Int(27), node.Int(8)), nil},
		{`"привет "||"мир"`, node.Concat(node.Text("привет "), node.Text("мир")), nil},

		{"27==8", node.Eq(node.Int(27), node.Int(8)), nil},
		{"27!=8", node.Neq(node.Int(27), node.Int(8)), nil},
		{"27<8", node.Lt(node.Int(27), node.Int(8)), nil},
		{"27<=8", node.Lte(node.Int(27), node.Int(8)), nil},
		{"27>8", node.Gt(node.Int(27), node.Int(8)), nil},
		{"27>=8", node.Gte(node.Int(27), node.Int(8)), nil},
	}

	for _, test := range tests {
		tokens, err := lexer.Tokenize(test.data)
		assert.NoError(t, err)

		p := newParser(tokens)

		v, err := p.eq()

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_not(t *testing.T) {
	tests := []struct {
		data          string
		expectedValue node.Node
		expectedError error
	}{
		{"2187", node.Int(2187), nil},
		{"(2187)", node.Int(2187), nil},
		{"() -> {}", node.Function(node.Block()), nil},
		{"array[1]", node.ElByIndex(node.Ident("array"), node.Int(1)), nil},
		{"factorial()", node.Call(node.Ident("factorial")), nil},
		{"-2187", node.Neg(node.Int(2187)), nil},
		{"27*8", node.Mul(node.Int(27), node.Int(8)), nil},
		{"27/8", node.Div(node.Int(27), node.Int(8)), nil},
		{"27%8", node.Mod(node.Int(27), node.Int(8)), nil},
		{"27+8", node.Add(node.Int(27), node.Int(8)), nil},
		{"27-8", node.Sub(node.Int(27), node.Int(8)), nil},
		{`"привет "||"мир"`, node.Concat(node.Text("привет "), node.Text("мир")), nil},
		{"27==8", node.Eq(node.Int(27), node.Int(8)), nil},
		{"27!=8", node.Neq(node.Int(27), node.Int(8)), nil},
		{"27<8", node.Lt(node.Int(27), node.Int(8)), nil},
		{"27<=8", node.Lte(node.Int(27), node.Int(8)), nil},
		{"27>8", node.Gt(node.Int(27), node.Int(8)), nil},
		{"27>=8", node.Gte(node.Int(27), node.Int(8)), nil},

		{"not true", node.Not(node.Bool(true)), nil},
		{"not not true", node.Not(node.Not(node.Bool(true))), nil},
	}

	for _, test := range tests {
		tokens, err := lexer.Tokenize(test.data)
		assert.NoError(t, err)

		p := newParser(tokens)

		v, err := p.not()

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_and(t *testing.T) {
	tests := []struct {
		data          string
		expectedValue node.Node
		expectedError error
	}{
		{"2187", node.Int(2187), nil},
		{"(2187)", node.Int(2187), nil},
		{"() -> {}", node.Function(node.Block()), nil},
		{"array[1]", node.ElByIndex(node.Ident("array"), node.Int(1)), nil},
		{"factorial()", node.Call(node.Ident("factorial")), nil},
		{"-2187", node.Neg(node.Int(2187)), nil},
		{"27*8", node.Mul(node.Int(27), node.Int(8)), nil},
		{"27/8", node.Div(node.Int(27), node.Int(8)), nil},
		{"27%8", node.Mod(node.Int(27), node.Int(8)), nil},
		{"27+8", node.Add(node.Int(27), node.Int(8)), nil},
		{"27-8", node.Sub(node.Int(27), node.Int(8)), nil},
		{`"привет "||"мир"`, node.Concat(node.Text("привет "), node.Text("мир")), nil},
		{"27==8", node.Eq(node.Int(27), node.Int(8)), nil},
		{"27!=8", node.Neq(node.Int(27), node.Int(8)), nil},
		{"27<8", node.Lt(node.Int(27), node.Int(8)), nil},
		{"27<=8", node.Lte(node.Int(27), node.Int(8)), nil},
		{"27>8", node.Gt(node.Int(27), node.Int(8)), nil},
		{"27>=8", node.Gte(node.Int(27), node.Int(8)), nil},
		{"not true", node.Not(node.Bool(true)), nil},

		{"false and true", node.And(node.Bool(false), node.Bool(true)), nil},
		{"false and true and false", node.And(
			node.And(node.Bool(false), node.Bool(true)),
			node.Bool(false),
		), nil},
	}

	for _, test := range tests {
		tokens, err := lexer.Tokenize(test.data)
		assert.NoError(t, err)

		p := newParser(tokens)

		v, err := p.and()

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_or(t *testing.T) {
	tests := []struct {
		data          string
		expectedValue node.Node
		expectedError error
	}{
		{"2187", node.Int(2187), nil},
		{"(2187)", node.Int(2187), nil},
		{"() -> {}", node.Function(node.Block()), nil},
		{"array[1]", node.ElByIndex(node.Ident("array"), node.Int(1)), nil},
		{"factorial()", node.Call(node.Ident("factorial")), nil},
		{"-2187", node.Neg(node.Int(2187)), nil},
		{"27*8", node.Mul(node.Int(27), node.Int(8)), nil},
		{"27/8", node.Div(node.Int(27), node.Int(8)), nil},
		{"27%8", node.Mod(node.Int(27), node.Int(8)), nil},
		{"27+8", node.Add(node.Int(27), node.Int(8)), nil},
		{"27-8", node.Sub(node.Int(27), node.Int(8)), nil},
		{`"привет "||"мир"`, node.Concat(node.Text("привет "), node.Text("мир")), nil},
		{"27==8", node.Eq(node.Int(27), node.Int(8)), nil},
		{"27!=8", node.Neq(node.Int(27), node.Int(8)), nil},
		{"27<8", node.Lt(node.Int(27), node.Int(8)), nil},
		{"27<=8", node.Lte(node.Int(27), node.Int(8)), nil},
		{"27>8", node.Gt(node.Int(27), node.Int(8)), nil},
		{"27>=8", node.Gte(node.Int(27), node.Int(8)), nil},
		{"not true", node.Not(node.Bool(true)), nil},
		{"false and true", node.And(node.Bool(false), node.Bool(true)), nil},

		{"false or true", node.Or(node.Bool(false), node.Bool(true)), nil},
		{"false or true or false", node.Or(
			node.Or(node.Bool(false), node.Bool(true)),
			node.Bool(false),
		), nil},
	}

	for _, test := range tests {
		tokens, err := lexer.Tokenize(test.data)
		assert.NoError(t, err)

		p := newParser(tokens)

		v, err := p.or()

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_set(t *testing.T) {
	tests := []struct {
		data          string
		expectedValue node.Node
		expectedError error
	}{
		{"2187", node.Int(2187), nil},
		{"(2187)", node.Int(2187), nil},
		{"() -> {}", node.Function(node.Block()), nil},
		{"array[1]", node.ElByIndex(node.Ident("array"), node.Int(1)), nil},
		{"factorial()", node.Call(node.Ident("factorial")), nil},
		{"-2187", node.Neg(node.Int(2187)), nil},
		{"27*8", node.Mul(node.Int(27), node.Int(8)), nil},
		{"27/8", node.Div(node.Int(27), node.Int(8)), nil},
		{"27%8", node.Mod(node.Int(27), node.Int(8)), nil},
		{"27+8", node.Add(node.Int(27), node.Int(8)), nil},
		{"27-8", node.Sub(node.Int(27), node.Int(8)), nil},
		{`"привет "||"мир"`, node.Concat(node.Text("привет "), node.Text("мир")), nil},
		{"27==8", node.Eq(node.Int(27), node.Int(8)), nil},
		{"27!=8", node.Neq(node.Int(27), node.Int(8)), nil},
		{"27<8", node.Lt(node.Int(27), node.Int(8)), nil},
		{"27<=8", node.Lte(node.Int(27), node.Int(8)), nil},
		{"27>8", node.Gt(node.Int(27), node.Int(8)), nil},
		{"27>=8", node.Gte(node.Int(27), node.Int(8)), nil},
		{"not true", node.Not(node.Bool(true)), nil},
		{"false and true", node.And(node.Bool(false), node.Bool(true)), nil},
		{"false or true", node.Or(node.Bool(false), node.Bool(true)), nil},

		{"age:=27", node.Create(node.Ident("age"), node.Int(27)), nil},
		{"age=27", node.Set(node.Ident("age"), node.Int(27)), nil},
		{"age = number = 27", node.Set(
			node.Ident("age"),
			node.Set(node.Ident("number"), node.Int(27)),
		), nil},
	}

	for _, test := range tests {
		tokens, err := lexer.Tokenize(test.data)
		assert.NoError(t, err)

		p := newParser(tokens)

		v, err := p.set()

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_expression(t *testing.T) {
	tests := []struct {
		data          string
		expectedValue node.Node
		expectedError error
	}{
		{"res:=(4+9)*27/16",
			node.Create(
				node.Ident("res"),
				node.Div(
					node.Mul(
						node.Add(node.Int(4), node.Int(9)),
						node.Int(27)),
					node.Int(16)),
			), nil},

		{"res=(4+9)*27/16 <= 19683%2187",
			node.Set(
				node.Ident("res"),
				node.Lte(
					node.Div(
						node.Mul(
							node.Add(node.Int(4), node.Int(9)),
							node.Int(27)),
						node.Int(16)),
					node.Mod(node.Int(19683), node.Int(2187)),
				),
			), nil},

		{"-(res[4+9*27]) > -(27*9)",
			node.Gt(
				node.Neg(
					node.ElByIndex(
						node.Ident("res"),
						node.Add(
							node.Int(4),
							node.Mul(node.Int(9), node.Int(27)),
						)),
				),
				node.Neg(node.Mul(node.Int(27), node.Int(9))),
			), nil},

		{"(res) -> {-(res[4+9*27]) > -(27*9)}",
			node.Function(
				node.Block(
					node.Gt(
						node.Neg(
							node.ElByIndex(
								node.Ident("res"),
								node.Add(
									node.Int(4),
									node.Mul(node.Int(9), node.Int(27)),
								)),
						),
						node.Neg(node.Mul(node.Int(27), node.Int(9))),
					),
				), node.Ident("res"),
			), nil},
	}

	for _, test := range tests {
		tokens, err := lexer.Tokenize(test.data)
		assert.NoError(t, err)

		p := newParser(tokens)

		v, err := p.expression()

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_branch(t *testing.T) {
	tests := []struct {
		data          string
		expectedValue node.Node
		expectedError error
	}{
		{"if 21<87 {81}", node.If(
			node.Branch{
				Cond: node.Lt(node.Int(21), node.Int(87)),
				Body: node.Block(node.Int(81)),
			},
		), nil},

		{"if 21<87 {81} elif 16==27 {125;625}", node.If(
			node.Branch{
				Cond: node.Lt(node.Int(21), node.Int(87)),
				Body: node.Block(node.Int(81)),
			},
			node.Branch{
				Cond: node.Eq(node.Int(16), node.Int(27)),
				Body: node.Block(node.Int(125), node.Int(625)),
			},
		), nil},

		{"if 21<87 {81} elif 16==27 {125;625} elif 9>3 {1}", node.If(
			node.Branch{
				Cond: node.Lt(node.Int(21), node.Int(87)),
				Body: node.Block(node.Int(81)),
			},
			node.Branch{
				Cond: node.Eq(node.Int(16), node.Int(27)),
				Body: node.Block(node.Int(125), node.Int(625)),
			},
			node.Branch{
				Cond: node.Gt(node.Int(9), node.Int(3)),
				Body: node.Block(node.Int(1)),
			},
		), nil},

		{"if 21<87 {81} elif 16==27 {125;625} else {1}", node.If(
			node.Branch{
				Cond: node.Lt(node.Int(21), node.Int(87)),
				Body: node.Block(node.Int(81)),
			},
			node.Branch{
				Cond: node.Eq(node.Int(16), node.Int(27)),
				Body: node.Block(node.Int(125), node.Int(625)),
			},
			node.Branch{
				Cond: node.Bool(true),
				Body: node.Block(node.Int(1)),
			},
		), nil},

		{"if 21<87 {81} elif 16==27 {125;625} elif 9>3 {1} else {64}", node.If(
			node.Branch{
				Cond: node.Lt(node.Int(21), node.Int(87)),
				Body: node.Block(node.Int(81)),
			},
			node.Branch{
				Cond: node.Eq(node.Int(16), node.Int(27)),
				Body: node.Block(node.Int(125), node.Int(625)),
			},
			node.Branch{
				Cond: node.Gt(node.Int(9), node.Int(3)),
				Body: node.Block(node.Int(1)),
			},
			node.Branch{
				Cond: node.Bool(true),
				Body: node.Block(node.Int(64)),
			},
		), nil},
	}

	for _, test := range tests {
		tokens, err := lexer.Tokenize(test.data)
		assert.NoError(t, err)

		p := newParser(tokens)

		v, err := p.branch()

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_loop(t *testing.T) {
	tests := []struct {
		data          string
		expectedValue node.Node
		expectedError error
	}{
		{"for i in 81 {i}", node.For(
			[]node.Node{node.Ident("i")},
			node.Int(81),
			node.Block(node.Ident("i")),
		), nil},
		{"for i, in 81 {i}", node.For(
			[]node.Node{node.Ident("i")},
			node.Int(81),
			node.Block(node.Ident("i")),
		), nil},
		{"for i,j in [81] {i;j}", node.For(
			[]node.Node{node.Ident("i"), node.Ident("j")},
			node.Array(node.Int(81)),
			node.Block(node.Ident("i"), node.Ident("j")),
		), nil},
	}

	for _, test := range tests {
		tokens, err := lexer.Tokenize(test.data)
		assert.NoError(t, err)

		p := newParser(tokens)

		v, err := p.loop()

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_ret(t *testing.T) {
	tests := []struct {
		data          string
		expectedValue node.Node
		expectedError error
	}{
		{"return 81", node.Return(node.Int(81)), nil},
	}

	for _, test := range tests {
		tokens, err := lexer.Tokenize(test.data)
		assert.NoError(t, err)

		p := newParser(tokens)

		v, err := p.ret()

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_construction(t *testing.T) {
	tests := []struct {
		data          string
		expectedValue node.Node
		expectedError error
	}{
		{
			`if 21<87 {for i in 81 {i}};`,
			node.If(
				node.Branch{
					Cond: node.Lt(node.Int(21), node.Int(87)),
					Body: node.Block(
						node.For([]node.Node{node.Ident("i")}, node.Int(81), node.Block(node.Ident("i"))),
					),
				},
			), nil},

		{
			`
if 21<87 {for i in 81 {i}}
elif 16==27 {for j in 21 {j}};
`,
			node.If(
				node.Branch{
					Cond: node.Lt(node.Int(21), node.Int(87)),
					Body: node.Block(
						node.For([]node.Node{node.Ident("i")}, node.Int(81), node.Block(node.Ident("i"))),
					),
				},
				node.Branch{
					Cond: node.Eq(node.Int(16), node.Int(27)),
					Body: node.Block(
						node.For([]node.Node{node.Ident("j")}, node.Int(21), node.Block(node.Ident("j"))),
					),
				},
			), nil},

		{
			`
if 21<87 {for i in 81 {i}}
elif 16==27 {for j in 21 {j}}
else {for k in 1 {k}};
`,
			node.If(
				node.Branch{
					Cond: node.Lt(node.Int(21), node.Int(87)),
					Body: node.Block(
						node.For([]node.Node{node.Ident("i")}, node.Int(81), node.Block(node.Ident("i"))),
					),
				},
				node.Branch{
					Cond: node.Eq(node.Int(16), node.Int(27)),
					Body: node.Block(
						node.For([]node.Node{node.Ident("j")}, node.Int(21), node.Block(node.Ident("j"))),
					),
				},
				node.Branch{
					Cond: node.Bool(true),
					Body: node.Block(
						node.For([]node.Node{node.Ident("k")}, node.Int(1), node.Block(node.Ident("k"))),
					),
				},
			), nil},

		{`
for n in 81 {
	if n<87 {for i in 81 {i}}
	elif n==27 {for j in 21 {j}}
	else {for k in 1 {k}};
};
`,
			node.For(
				[]node.Node{node.Ident("n")},
				node.Int(81),
				node.Block(
					node.If(
						node.Branch{
							Cond: node.Lt(node.Ident("n"), node.Int(87)),
							Body: node.Block(
								node.For([]node.Node{node.Ident("i")}, node.Int(81), node.Block(node.Ident("i"))),
							),
						},
						node.Branch{
							Cond: node.Eq(node.Ident("n"), node.Int(27)),
							Body: node.Block(
								node.For([]node.Node{node.Ident("j")}, node.Int(21), node.Block(node.Ident("j"))),
							),
						},
						node.Branch{
							Cond: node.Bool(true),
							Body: node.Block(
								node.For([]node.Node{node.Ident("k")}, node.Int(1), node.Block(node.Ident("k"))),
							),
						},
					),
				),
			),
			nil},

		{`
(l) -> {
	for n in l {
		if n<87 {for i in 81 {i}}
		elif n==27 {for j in 21 {j}}
		else {for k in 1 {k}};
	};
	return l;
};
`,
			node.Function(
				node.Block(
					node.For(
						[]node.Node{node.Ident("n")},
						node.Ident("l"),
						node.Block(
							node.If(
								node.Branch{
									Cond: node.Lt(node.Ident("n"), node.Int(87)),
									Body: node.Block(
										node.For([]node.Node{node.Ident("i")}, node.Int(81), node.Block(node.Ident("i"))),
									),
								},
								node.Branch{
									Cond: node.Eq(node.Ident("n"), node.Int(27)),
									Body: node.Block(
										node.For([]node.Node{node.Ident("j")}, node.Int(21), node.Block(node.Ident("j"))),
									),
								},
								node.Branch{
									Cond: node.Bool(true),
									Body: node.Block(
										node.For([]node.Node{node.Ident("k")}, node.Int(1), node.Block(node.Ident("k"))),
									),
								},
							),
						),
					),
					node.Return(node.Ident("l")),
				),
				node.Ident("l"),
			),
			nil},
	}

	for _, test := range tests {
		tokens, err := lexer.Tokenize(test.data)
		assert.NoError(t, err)

		p := newParser(tokens)

		v, err := p.construction()

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_parse(t *testing.T) {
	tests := []struct {
		data          string
		expectedValue node.Node
		expectedError error
	}{
		{`
n := 2187;

if n<87 {
	n = 7;
};

print(n);
`,
			node.Block(
				node.Create(node.Ident("n"), node.Int(2187)),
				node.If(
					node.Branch{
						Cond: node.Lt(node.Ident("n"), node.Int(87)),
						Body: node.Block(node.Set(node.Ident("n"), node.Int(7))),
					},
				),
				node.Call(node.Ident("print"), node.Ident("n")),
			),
			nil},
	}

	for _, test := range tests {
		tokens, err := lexer.Tokenize(test.data)
		assert.NoError(t, err)

		p := newParser(tokens)

		v, err := p.parse()

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}
