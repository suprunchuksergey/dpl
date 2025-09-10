package parser

import (
	"fmt"
	"github.com/suprunchuksergey/dpl/internal/node"
	"github.com/suprunchuksergey/dpl/lexer"
	"github.com/suprunchuksergey/dpl/token"
	"strconv"
)

/*
parser реализует рекурсивный спуск для разбора выражений
с учетом приоритета операций.

слои (layer[p], где p — приоритет, меньший p = выше приоритет)
обрабатывают операторы (например, layer5: *, /, %; layer6: +, -) и операнды.

слои:
layer1 -> true, false, null, int, real, text, [], {}, ident
layer2 -> ()
layer3 -> [] (доступ к индексу)
layer4 -> - (унарный)
layer5 -> *, /, %
layer6 -> +, -
layer7 -> ||
layer8 -> ==, !=, <, <=, >, >=
layer9 -> not
layer10 -> and
layer11 -> or
layer12 -> =
layer13 -> if
layer14 -> for .. in
*/
type parser struct{ lex lexer.Lexer }

// true, false, null, int, real, text, [], {}, ident
func (p *parser) layer1() (node.Node, error) {
	var n node.Node

	switch tok := p.lex.Tok(); tok.ID() {
	case token.Fn:
		err := p.lex.Next()
		if err != nil {
			return nil, err
		}

		if !p.lex.Tok().Is(token.LParen) {
			return nil, unexpected(p.lex.Tok())
		}

		if err = p.lex.Next(); err != nil {
			return nil, err
		}

		var names []node.Node
		if p.lex.Tok().Is(token.RParen) {
			if err := p.lex.Next(); err != nil {
				return nil, err
			}
		} else {
			for {
				name, err := p.layer1()
				if err != nil {
					return nil, err
				}
				names = append(names, name)

				if p.lex.Tok().Is(token.Comma) {
					if err := p.lex.Next(); err != nil {
						return nil, err
					}

					if p.lex.Tok().Is(token.RParen) {
						break
					}

					continue
				}

				if !p.lex.Tok().Is(token.RParen) {
					return nil, unexpected(p.lex.Tok())
				}

				if err = p.lex.Next(); err != nil {
					return nil, err
				}
				break
			}
		}

		body, err := p.parseBody()
		if err != nil {
			return nil, err
		}
		if len(names) == 0 {
			return node.Function(body), nil
		}
		return node.Function(body, names...), nil

	case token.Ident:
		n = node.Ident(tok.(token.WithValue).Value())

	case token.LBrack:
		err := p.lex.Next()
		if err != nil {
			return nil, err
		}

		if p.lex.Tok().Is(token.RBrack) {
			n = node.Array()
			break
		}

		var items []node.Node
		for {
			n, err := p.layer12()
			if err != nil {
				return nil, err
			}
			items = append(items, n)

			if !p.lex.Tok().OneOf(token.Comma, token.RBrack) {
				return nil, unexpected(p.lex.Tok())
			}

			if p.lex.Tok().Is(token.Comma) {
				if err = p.lex.Next(); err != nil {
					return nil, err
				}
			}

			if p.lex.Tok().Is(token.RBrack) {
				break
			}
		}
		n = node.Array(items...)

	case token.LBrace:
		err := p.lex.Next()
		if err != nil {
			return nil, err
		}

		records := make([]node.KV, 0)

		if p.lex.Tok().Is(token.RBrace) {
			n = node.Object()
			break
		}

		for {
			k, err := p.layer12()
			if err != nil {
				return nil, err
			}

			if !p.lex.Tok().Is(token.Colon) {
				return nil, unexpected(p.lex.Tok())
			}
			if err = p.lex.Next(); err != nil {
				return nil, err
			}

			v, err := p.layer12()
			if err != nil {
				return nil, err
			}

			records = append(records, node.KV{Key: k, Value: v})

			if !p.lex.Tok().OneOf(token.Comma, token.RBrace) {
				return nil, unexpected(p.lex.Tok())
			}

			if p.lex.Tok().Is(token.Comma) {
				err = p.lex.Next()
				if err != nil {
					return nil, err
				}
			}

			if p.lex.Tok().Is(token.RBrace) {
				break
			}
		}
		n = node.Object(records...)

	case token.Null:
		n = node.Null()
	case token.True:
		n = node.Bool(true)
	case token.False:
		n = node.Bool(false)
	case token.Int:
		v := tok.(token.WithValue).Value()
		num, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		n = node.Int(num)
	case token.Real:
		v := tok.(token.WithValue).Value()
		num, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}
		n = node.Real(num)
	case token.Text:
		n = node.Text(
			tok.(token.WithValue).Value())
	default:
		return nil, unexpected(tok)
	}

	err := p.lex.Next()
	if err != nil {
		return nil, err
	}
	return n, nil
}

// ()
func (p *parser) layer2() (node.Node, error) {
	if p.lex.Tok().ID() != token.LParen {
		return p.layer1()
	}

	err := p.lex.Next()
	if err != nil {
		return nil, err
	}

	n, err := p.layer12()
	if err != nil {
		return nil, err
	}

	if p.lex.Tok().ID() != token.RParen {
		return nil, unexpected(p.lex.Tok())
	}

	err = p.lex.Next()
	if err != nil {
		return nil, err
	}

	return n, nil
}

// [] (доступ к индексу)
func (p *parser) layer3() (node.Node, error) {
	n, err := p.layer2()
	if err != nil {
		return nil, err
	}

	for {
		if p.lex.Tok().Is(token.LBrack) {
			if err = p.lex.Next(); err != nil {
				return nil, err
			}

			i, err := p.layer12()
			if err != nil {
				return nil, err
			}

			if !p.lex.Tok().Is(token.RBrack) {
				return nil, unexpected(p.lex.Tok())
			}
			if err = p.lex.Next(); err != nil {
				return nil, err
			}

			n = node.ElByIndex(n, i)
		} else if p.lex.Tok().Is(token.LParen) {
			if err = p.lex.Next(); err != nil {
				return nil, err
			}

			if p.lex.Tok().Is(token.RParen) {
				if err = p.lex.Next(); err != nil {
					return nil, err
				}
				n = node.Call(n)
				continue
			}

			var args []node.Node
			for {
				i, err := p.layer12()
				if err != nil {
					return nil, err
				}
				args = append(args, i)

				if p.lex.Tok().Is(token.Comma) {
					if err = p.lex.Next(); err != nil {
						return nil, err
					}

					if p.lex.Tok().Is(token.RParen) {
						if err = p.lex.Next(); err != nil {
							return nil, err
						}
						n = node.Call(n, args...)
						break
					}

					continue
				}

				if !p.lex.Tok().Is(token.RParen) {
					return nil, unexpected(p.lex.Tok())
				}

				if err = p.lex.Next(); err != nil {
					return nil, err
				}

				n = node.Call(n, args...)

				break
			}
		} else {
			break
		}
	}

	return n, nil
}

// - (унарный)
func (p *parser) layer4() (node.Node, error) {
	if p.lex.Tok().ID() != token.Sub {
		return p.layer3()
	}

	err := p.lex.Next()
	if err != nil {
		return nil, err
	}

	n, err := p.layer4()
	if err != nil {
		return nil, err
	}

	return node.Neg(n), nil
}

// *, /, %
func (p *parser) layer5() (node.Node, error) {
	n, err := p.layer4()
	if err != nil {
		return nil, err
	}

	for p.lex.Tok().OneOf(token.Mul, token.Div, token.Rem) {
		id := p.lex.Tok().ID()

		err := p.lex.Next()
		if err != nil {
			return nil, err
		}

		r, err := p.layer4()
		if err != nil {
			return nil, err
		}

		switch id {
		case token.Mul:
			n = node.Mul(n, r)
		case token.Div:
			n = node.Div(n, r)
		default:
			n = node.Mod(n, r)
		}
	}

	return n, nil
}

// +, -
func (p *parser) layer6() (node.Node, error) {
	n, err := p.layer5()
	if err != nil {
		return nil, err
	}

	for p.lex.Tok().OneOf(token.Add, token.Sub) {
		id := p.lex.Tok().ID()

		err := p.lex.Next()
		if err != nil {
			return nil, err
		}

		r, err := p.layer5()
		if err != nil {
			return nil, err
		}

		switch id {
		case token.Add:
			n = node.Add(n, r)
		default:
			n = node.Sub(n, r)
		}
	}

	return n, nil
}

// ||
func (p *parser) layer7() (node.Node, error) {
	n, err := p.layer6()
	if err != nil {
		return nil, err
	}

	for p.lex.Tok().Is(token.Concat) {
		err := p.lex.Next()
		if err != nil {
			return nil, err
		}

		r, err := p.layer6()
		if err != nil {
			return nil, err
		}

		n = node.Concat(n, r)
	}

	return n, nil
}

// ==, !=, <, <=, >, >=
func (p *parser) layer8() (node.Node, error) {
	n, err := p.layer7()
	if err != nil {
		return nil, err
	}

	for p.lex.Tok().OneOf(token.Eq, token.Neq, token.Lt,
		token.Gt, token.Lte, token.Gte) {
		id := p.lex.Tok().ID()

		err := p.lex.Next()
		if err != nil {
			return nil, err
		}

		r, err := p.layer7()
		if err != nil {
			return nil, err
		}

		switch id {
		case token.Eq:
			n = node.Eq(n, r)
		case token.Neq:
			n = node.Neq(n, r)
		case token.Lt:
			n = node.Lt(n, r)
		case token.Gt:
			n = node.Gt(n, r)
		case token.Lte:
			n = node.Lte(n, r)
		default:
			n = node.Gte(n, r)
		}
	}

	return n, nil
}

// not
func (p *parser) layer9() (node.Node, error) {
	if p.lex.Tok().ID() != token.Not {
		return p.layer8()
	}

	err := p.lex.Next()
	if err != nil {
		return nil, err
	}

	n, err := p.layer9()
	if err != nil {
		return nil, err
	}

	return node.Not(n), nil
}

// and
func (p *parser) layer10() (node.Node, error) {
	n, err := p.layer9()
	if err != nil {
		return nil, err
	}

	for p.lex.Tok().Is(token.And) {
		err := p.lex.Next()
		if err != nil {
			return nil, err
		}

		r, err := p.layer9()
		if err != nil {
			return nil, err
		}

		n = node.And(n, r)
	}

	return n, nil
}

// or
func (p *parser) layer11() (node.Node, error) {
	n, err := p.layer10()
	if err != nil {
		return nil, err
	}

	for p.lex.Tok().Is(token.Or) {
		err := p.lex.Next()
		if err != nil {
			return nil, err
		}

		r, err := p.layer10()
		if err != nil {
			return nil, err
		}

		n = node.Or(n, r)
	}

	return n, nil
}

// =
func (p *parser) layer12() (node.Node, error) {
	n, err := p.layer11()
	if err != nil {
		return nil, err
	}

	if p.lex.Tok().Is(token.Assign) {
		err := p.lex.Next()
		if err != nil {
			return nil, err
		}

		r, err := p.layer11()
		if err != nil {
			return nil, err
		}
		n = node.Set(n, r)
	}

	return n, nil
}

// if
func (p *parser) layer13() (node.Node, error) {
	if !p.lex.Tok().Is(token.If) {
		return p.layer12()
	}

	if err := p.lex.Next(); err != nil {
		return nil, err
	}

	branches := make([]node.Branch, 0)

	cond, err := p.layer12()
	if err != nil {
		return nil, err
	}

	body, err := p.parseBody()
	if err != nil {
		return nil, err
	}

	branches = append(branches, node.Branch{Cond: cond, Body: body})

	for p.lex.Tok().Is(token.Elif) {
		if err := p.lex.Next(); err != nil {
			return nil, err
		}

		cond, err := p.layer12()
		if err != nil {
			return nil, err
		}

		body, err := p.parseBody()
		if err != nil {
			return nil, err
		}

		branches = append(branches, node.Branch{Cond: cond, Body: body})
	}

	if p.lex.Tok().Is(token.Else) {
		if err = p.lex.Next(); err != nil {
			return nil, err
		}

		body, err := p.parseBody()
		if err != nil {
			return nil, err
		}

		branches = append(branches, node.Branch{Cond: node.Bool(true), Body: body})
	}

	return node.If(branches...), nil
}

func (p *parser) layer14() (node.Node, error) {
	if !p.lex.Tok().Is(token.For) {
		return p.layer13()
	}

	err := p.lex.Next()
	if err != nil {
		return nil, err
	}

	id1, err := p.layer1()
	if err != nil {
		return nil, err
	}

	var id2 node.Node
	if p.lex.Tok().Is(token.Comma) {
		if err = p.lex.Next(); err != nil {
			return nil, err
		}

		if id2, err = p.layer1(); err != nil {
			return nil, err
		}
	}

	if !p.lex.Tok().Is(token.In) {
		return nil, unexpected(p.lex.Tok())
	}

	if err = p.lex.Next(); err != nil {
		return nil, err
	}

	v, err := p.layer12()
	if err != nil {
		return nil, err
	}

	body, err := p.parseBody()
	if err != nil {
		return nil, err
	}

	recs := []node.Node{}
	if id1 != nil {
		recs = append(recs, id1)
	}

	if id2 != nil {
		recs = append(recs, id2)
	}

	return node.For(recs, v, body), nil
}

func (p *parser) layer15() (node.Node, error) {
	if !p.lex.Tok().Is(token.Return) {
		return p.layer14()
	}

	if err := p.lex.Next(); err != nil {
		return nil, err
	}

	v, err := p.layer12()
	if err != nil {
		return nil, err
	}

	return node.Return(v), nil
}

func (p *parser) parseBody() (node.Node, error) {
	if !p.lex.Tok().Is(token.LBrace) {
		return nil, unexpected(p.lex.Tok())
	}

	if err := p.lex.Next(); err != nil {
		return nil, err
	}

	var nodes []node.Node
	for {
		n, err := p.layer15()
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, n)

		if p.lex.Tok().Is(token.Semicolon) {
			if err = p.lex.Next(); err != nil {
				return nil, err
			}

			if p.lex.Tok().Is(token.RBrace) {
				break
			}

			continue
		}

		if !p.lex.Tok().Is(token.RBrace) {
			return nil, unexpected(p.lex.Tok())
		}

		break
	}

	if err := p.lex.Next(); err != nil {
		return nil, err
	}

	return node.Block(nodes...), nil
}

func (p *parser) Parse() (node.Node, error) {
	//если это конец
	if p.lex.Tok().Is(token.EOF) {
		return nil, nil
	}

	var cmds []node.Node
	for {
		n, err := p.layer14()
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, n)

		if p.lex.Tok().Is(token.Semicolon) {
			err = p.lex.Next()
			if err != nil {
				return nil, err
			}

			if p.lex.Tok().Is(token.EOF) {
				break
			}

			continue
		}

		if !p.lex.Tok().Is(token.EOF) {
			return nil, unexpected(p.lex.Tok())
		}
		break
	}

	return node.Block(cmds...), nil
}

func newParser(lex lexer.Lexer) *parser { return &parser{lex} }

type Parser interface {
	//возвращает узел или ошибку,
	//узел будет nil, если первый токен EOF
	Parse() (node.Node, error)
}

// у лексера должен быть прочитан первый токен
func New(lex lexer.Lexer) Parser { return newParser(lex) }

func unexpected(tok token.Token) error {
	return fmt.Errorf("%s неожиданный токен %s", tok.Start(), tok)
}
