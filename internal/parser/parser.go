package parser

import (
	"fmt"
	"github.com/suprunchuksergey/dpl/internal/lexer"
	"github.com/suprunchuksergey/dpl/internal/node"
	"strconv"
)

type parser struct {
	tokens []lexer.Token
	index  int
}

func newParser(tokens []lexer.Token) *parser {
	return &parser{tokens: tokens}
}

func (p *parser) next() { p.index++ }

func (p *parser) token() lexer.Token { return p.tokens[p.index] }

func (p *parser) id() uint8 { return p.token().ID() }

func unexpectedToken(token lexer.Token) error {
	return fmt.Errorf("неожиданный токен %s", token)
}

func (p *parser) commands(sep, stop uint8, handler func() (node.Node, error)) ([]node.Node, error) {
	if p.id() == stop {
		p.next()
		return nil, nil
	}

	var nodes []node.Node
	for {
		n, err := handler()
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, n)

		if p.id() == sep {
			p.next()
			if p.id() == stop {
				p.next()
				return nodes, nil
			}
			continue
		}
		break
	}

	if p.id() != stop {
		return nil, unexpectedToken(p.token())
	}
	p.next()

	return nodes, nil
}

func (p *parser) value() (node.Node, error) {
	switch p.id() {
	case lexer.Null:
		p.next()
		return node.Null(), nil
	case lexer.False:
		p.next()
		return node.Bool(false), nil
	case lexer.True:
		p.next()
		return node.Bool(true), nil
	case lexer.Ident:
		value := p.token().(lexer.TokenWithValue).Value()
		p.next()
		return node.Ident(value), nil
	case lexer.Int:
		value := p.token().(lexer.TokenWithValue).Value()
		p.next()
		n, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, err
		}
		return node.Int(n), nil
	case lexer.Real:
		value := p.token().(lexer.TokenWithValue).Value()
		p.next()
		n, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
		return node.Real(n), nil
	case lexer.Text:
		value := p.token().(lexer.TokenWithValue).Value()
		p.next()
		return node.Text(value), nil
	case lexer.LBrack:
		p.next()

		nodes, err := p.commands(lexer.Comma, lexer.RBrack, p.expression)
		if err != nil {
			return nil, err
		}

		return node.Array(nodes...), nil
	case lexer.LBrace:
		p.next()

		if p.id() == lexer.RBrace {
			p.next()
			return node.Object(), nil
		}

		pairs := make([]node.KV, 0)
		for {
			k, err := p.expression()
			if err != nil {
				return nil, err
			}

			if p.id() != lexer.Colon {
				return nil, unexpectedToken(p.token())
			}
			p.next()

			v, err := p.expression()
			if err != nil {
				return nil, err
			}

			pairs = append(pairs, node.KV{Key: k, Value: v})

			if p.id() == lexer.Comma {
				p.next()
				if p.id() == lexer.RBrace {
					p.next()
					return node.Object(pairs...), nil
				}
				continue
			}
			break
		}

		if p.id() != lexer.RBrace {
			return nil, unexpectedToken(p.token())
		}
		p.next()

		return node.Object(pairs...), nil

	default:
		return nil, unexpectedToken(p.token())
	}
}

func (p *parser) paren() (node.Node, error) {
	if p.id() != lexer.LParen {
		return p.value()
	}
	p.next()

	nodes, err := p.commands(lexer.Comma, lexer.RParen, p.expression)
	if err != nil {
		return nil, err
	}

	if p.id() == lexer.ArrowRight {
		p.next()

		if p.id() != lexer.LBrace {
			return nil, unexpectedToken(p.token())
		}
		p.next()

		cmds, err := p.commands(lexer.Semicolon, lexer.RBrace, p.construction)
		if err != nil {
			return nil, err
		}

		return node.Function(node.Block(cmds...), nodes...), nil
	}

	if len(nodes) == 1 {
		return nodes[0], nil
	}

	return nil, unexpectedToken(p.token())
}

func (p *parser) elByIndex() (node.Node, error) {
	n, err := p.paren()
	if err != nil {
		return nil, err
	}

	for {
		if p.id() == lexer.LBrack {
			p.next()

			i, err := p.expression()
			if err != nil {
				return nil, err
			}

			if p.id() != lexer.RBrack {
				return nil, unexpectedToken(p.token())
			}
			p.next()

			n = node.ElByIndex(n, i)
			continue
		}

		if p.id() == lexer.LParen {
			p.next()

			nodes, err := p.commands(lexer.Comma, lexer.RParen, p.expression)
			if err != nil {
				return nil, err
			}

			n = node.Call(n, nodes...)
			continue
		}

		break
	}

	return n, nil
}

func (p *parser) neg() (node.Node, error) {
	if p.id() != lexer.Sub {
		return p.elByIndex()
	}
	p.next()

	v, err := p.neg()
	if err != nil {
		return nil, err
	}
	return node.Neg(v), nil
}

func (p *parser) mul() (node.Node, error) {
	n, err := p.neg()
	if err != nil {
		return nil, err
	}

	for p.id() == lexer.Mul || p.id() == lexer.Div || p.id() == lexer.Mod {
		id := p.id()
		p.next()

		v, err := p.neg()
		if err != nil {
			return nil, err
		}

		switch id {
		case lexer.Mul:
			n = node.Mul(n, v)
		case lexer.Div:
			n = node.Div(n, v)
		case lexer.Mod:
			n = node.Mod(n, v)
		default:
			panic("недостижимый")
		}
	}

	return n, nil
}

func (p *parser) add() (node.Node, error) {
	n, err := p.mul()
	if err != nil {
		return nil, err
	}

	for p.id() == lexer.Add || p.id() == lexer.Sub {
		id := p.id()
		p.next()

		v, err := p.mul()
		if err != nil {
			return nil, err
		}

		switch id {
		case lexer.Add:
			n = node.Add(n, v)
		case lexer.Sub:
			n = node.Sub(n, v)
		default:
			panic("недостижимый")
		}
	}

	return n, nil
}

func (p *parser) concat() (node.Node, error) {
	n, err := p.add()
	if err != nil {
		return nil, err
	}

	for p.id() == lexer.Concat {
		p.next()

		v, err := p.add()
		if err != nil {
			return nil, err
		}

		n = node.Concat(n, v)
	}

	return n, nil
}

func (p *parser) eq() (node.Node, error) {
	n, err := p.concat()
	if err != nil {
		return nil, err
	}

	for p.id() == lexer.Eq ||
		p.id() == lexer.Neq ||
		p.id() == lexer.Lt ||
		p.id() == lexer.Lte ||
		p.id() == lexer.Gt ||
		p.id() == lexer.Gte {

		id := p.id()
		p.next()

		v, err := p.concat()
		if err != nil {
			return nil, err
		}

		switch id {
		case lexer.Eq:
			n = node.Eq(n, v)
		case lexer.Neq:
			n = node.Neq(n, v)
		case lexer.Lt:
			n = node.Lt(n, v)
		case lexer.Lte:
			n = node.Lte(n, v)
		case lexer.Gt:
			n = node.Gt(n, v)
		case lexer.Gte:
			n = node.Gte(n, v)
		default:
			panic("недостижимый")
		}
	}

	return n, nil
}

func (p *parser) not() (node.Node, error) {
	if p.id() != lexer.Not {
		return p.eq()
	}
	p.next()

	v, err := p.not()
	if err != nil {
		return nil, err
	}
	return node.Not(v), nil
}

func (p *parser) and() (node.Node, error) {
	n, err := p.not()
	if err != nil {
		return nil, err
	}

	for p.id() == lexer.And {
		p.next()

		v, err := p.not()
		if err != nil {
			return nil, err
		}

		n = node.And(n, v)
	}

	return n, nil
}

func (p *parser) or() (node.Node, error) {
	n, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.id() == lexer.Or {
		p.next()

		v, err := p.and()
		if err != nil {
			return nil, err
		}

		n = node.Or(n, v)
	}

	return n, nil
}

func (p *parser) set() (node.Node, error) {
	n, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.id() != lexer.Set && p.id() != lexer.Create {
		return n, nil
	}

	id := p.id()
	p.next()

	v, err := p.set()
	if err != nil {
		return nil, err
	}

	if id == lexer.Set {
		return node.Set(n, v), nil
	}
	return node.Create(n, v), nil
}

func (p *parser) expression() (node.Node, error) { return p.set() }

func (p *parser) branch() (node.Node, error) {
	if p.id() != lexer.If {
		return p.expression()
	}

	branches := make([]node.Branch, 0)
	for {
		p.next()

		cond, err := p.expression()
		if err != nil {
			return nil, err
		}

		if p.id() != lexer.LBrace {
			return nil, unexpectedToken(p.token())
		}
		p.next()

		cmds, err := p.commands(lexer.Semicolon, lexer.RBrace, p.construction)
		if err != nil {
			return nil, err
		}

		branches = append(branches, node.Branch{Cond: cond, Body: node.Block(cmds...)})

		if p.id() != lexer.Elif {
			break
		}
	}

	if p.id() == lexer.Else {
		p.next()

		if p.id() != lexer.LBrace {
			return nil, unexpectedToken(p.token())
		}
		p.next()

		cmds, err := p.commands(lexer.Semicolon, lexer.RBrace, p.construction)
		if err != nil {
			return nil, err
		}

		branches = append(branches, node.Branch{Cond: node.Bool(true), Body: node.Block(cmds...)})
	}

	return node.If(branches...), nil
}

func (p *parser) loop() (node.Node, error) {
	if p.id() != lexer.For {
		return p.branch()
	}

	p.next()

	recipients, err := p.commands(lexer.Comma, lexer.In, p.value)
	if err != nil {
		return nil, err
	}

	from, err := p.expression()
	if err != nil {
		return nil, err
	}

	if p.id() != lexer.LBrace {
		return nil, unexpectedToken(p.token())
	}
	p.next()

	cmds, err := p.commands(lexer.Semicolon, lexer.RBrace, p.construction)
	if err != nil {
		return nil, err
	}

	return node.For(recipients, from, node.Block(cmds...)), nil
}

func (p *parser) ret() (node.Node, error) {
	if p.id() != lexer.Return {
		return p.loop()
	}

	p.next()

	v, err := p.expression()
	if err != nil {
		return nil, err
	}

	return node.Return(v), nil
}

func (p *parser) construction() (node.Node, error) { return p.ret() }

func (p *parser) parse() (node.Node, error) {
	cmds, err := p.commands(lexer.Semicolon, lexer.EOF, p.construction)
	if err != nil {
		return nil, err
	}
	return node.Block(cmds...), nil
}

func Parse(tokens []lexer.Token) (node.Node, error) { return newParser(tokens).parse() }
