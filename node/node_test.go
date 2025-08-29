package node

import (
	"fmt"
	"github.com/suprunchuksergey/dpl/namespace"
	"github.com/suprunchuksergey/dpl/val"
	"reflect"
	"testing"
)

type row struct {
	n        Node
	expected val.Val
}

func (r row) exec(init map[string]val.Val) error {
	v, err := r.n.Exec(namespace.New(init))
	// предполагается, что узел должен выполниться без ошибок
	if err != nil {
		return fmt.Errorf("непредвиденная ошибка: %s", err.Error())
	}

	if !reflect.DeepEqual(v, r.expected) {
		return fmt.Errorf("ожидалось %s, получено %s",
			r.expected, v)
	}

	return nil
}

func (r row) exect(t *testing.T, init map[string]val.Val) {
	err := r.exec(init)
	if err != nil {
		t.Error(err)
	}
}

func newRow(n Node, expected val.Val) row {
	return row{
		n:        n,
		expected: expected,
	}
}

type rows []row

func (rows rows) exec(t *testing.T, init map[string]val.Val) {
	for i, r := range rows {
		err := r.exec(init)
		if err != nil {
			t.Errorf("%d: %s", i, err.Error())
		}
	}
}

func newRows(r ...row) rows { return r }

func Test_Add(t *testing.T) {
	newRow(
		Add(Val(val.Int(125)), Val(val.Int(343))),
		val.Int(468),
	).exect(t, nil)
}

func Test_Sub(t *testing.T) {
	newRow(
		Sub(Val(val.Int(468)), Val(val.Int(343))),
		val.Int(125),
	).exect(t, nil)
}

func Test_Mul(t *testing.T) {
	newRow(
		Mul(Val(val.Int(468)), Val(val.Int(343))),
		val.Int(160524),
	).exect(t, nil)
}

func Test_Div(t *testing.T) {
	newRows(
		newRow(
			Div(Val(val.Int(160524)), Val(val.Int(343))),
			val.Int(468),
		),
		newRow(
			Div(Val(val.Int(160524)), Val(val.Int(0))),
			val.Null(),
		),
	).exec(t, nil)
}

func Test_Rem(t *testing.T) {
	newRows(
		newRow(
			Rem(Val(val.Int(160524)), Val(val.Int(10))),
			val.Int(4),
		),
		newRow(
			Rem(Val(val.Int(160524)), Val(val.Int(0))),
			val.Null(),
		),
	).exec(t, nil)
}

func Test_Eq(t *testing.T) {
	newRow(
		Eq(Val(val.Int(160524)), Val(val.Int(160524))),
		val.True(),
	).exect(t, nil)
}

func Test_Neq(t *testing.T) {
	newRow(
		Neq(Val(val.Int(160524)), Val(val.Int(19683))),
		val.True(),
	).exect(t, nil)
}

func Test_Lt(t *testing.T) {
	newRows(
		newRow(
			Lt(Val(val.Int(19683)), Val(val.Int(160524))),
			val.True(),
		),
		newRow(
			Lt(Val(val.Int(160524)), Val(val.Int(160524))),
			val.False(),
		),
	).exec(t, nil)
}

func Test_Lte(t *testing.T) {
	newRows(
		newRow(
			Lte(Val(val.Int(19683)), Val(val.Int(160524))),
			val.True(),
		),
		newRow(
			Lte(Val(val.Int(160524)), Val(val.Int(160524))),
			val.True(),
		),
	).exec(t, nil)
}

func Test_Gt(t *testing.T) {
	newRows(
		newRow(
			Gt(Val(val.Int(160524)), Val(val.Int(19683))),
			val.True(),
		),
		newRow(
			Gt(Val(val.Int(160524)), Val(val.Int(160524))),
			val.False(),
		),
	).exec(t, nil)
}

func Test_Gte(t *testing.T) {
	newRows(
		newRow(
			Gte(Val(val.Int(160524)), Val(val.Int(19683))),
			val.True(),
		),
		newRow(
			Gte(Val(val.Int(160524)), Val(val.Int(160524))),
			val.True(),
		),
	).exec(t, nil)
}

func Test_Concat(t *testing.T) {
	newRow(
		Concat(Val(val.Int(160524)), Val(val.Int(19683))),
		val.Text("16052419683"),
	).exect(t, nil)
}

func Test_And(t *testing.T) {
	newRows(
		newRow(
			And(Val(val.True()), Val(val.True())),
			val.True(),
		),
		newRow(
			And(Val(val.True()), Val(val.False())),
			val.False(),
		),
	).exec(t, nil)
}

func Test_Or(t *testing.T) {
	newRows(
		newRow(
			Or(Val(val.True()), Val(val.True())),
			val.True(),
		),
		newRow(
			Or(Val(val.True()), Val(val.False())),
			val.True(),
		),
	).exec(t, nil)
}

func Test_Not(t *testing.T) {
	newRows(
		newRow(Not(Val(val.True())), val.False()),
		newRow(Not(Val(val.False())), val.True()),
		newRow(Not(Val(val.Int(160524))), val.False()),
	).exec(t, nil)
}

func Test_Neg(t *testing.T) {
	newRow(
		Neg(Val(val.Int(160524))),
		val.Int(-160524),
	).exect(t, nil)
}

func Test_Val(t *testing.T) {
	newRows(
		newRow(Val(val.Int(160524)), val.Int(160524)),
		newRow(Val(val.Real(16.0524)), val.Real(16.0524)),
		newRow(Val(val.Text("16052419683")), val.Text("16052419683")),
		newRow(Val(val.True()), val.True()),
		newRow(Val(val.Null()), val.Null()),
	).exec(t, nil)
}

func Test_Int(t *testing.T) {
	newRow(
		Int(160524),
		val.Int(160524)).exect(t, nil)
}

func Test_Real(t *testing.T) {
	newRow(
		Real(16.0524),
		val.Real(16.0524)).exect(t, nil)
}

func Test_Text(t *testing.T) {
	newRow(
		Text("text"),
		val.Text("text")).exect(t, nil)
}

func Test_True(t *testing.T) {
	newRow(
		True(),
		val.True()).exect(t, nil)
}

func Test_False(t *testing.T) {
	newRow(
		False(),
		val.False()).exect(t, nil)
}

func Test_Null(t *testing.T) {
	newRow(
		Null(),
		val.Null()).exect(t, nil)
}

func Test_Array(t *testing.T) {
	newRows(
		newRow(Array([]Node{}), val.Array([]val.Val{})),
		newRow(Array([]Node{Int(81)}), val.Array([]val.Val{val.Int(81)})),
		newRow(
			Array([]Node{Int(81), Text("text")}),
			val.Array([]val.Val{val.Int(81), val.Text("text")})),
		newRow(
			Array([]Node{Add(Int(81), Int(27)), Text("text")}),
			val.Array([]val.Val{val.Int(108), val.Text("text")})),
	).exec(t, nil)
}

func Test_Map(t *testing.T) {
	newRows(
		newRow(
			Map(Records{}),
			val.Map(map[string]val.Val{})),
		newRow(
			Map(Records{
				Record{Text("text"), Int(81)},
				Record{Int(8), Real(8.1)},
				Record{True(), Null()},
			}),
			val.Map(map[string]val.Val{
				"text": val.Int(81),
				"8":    val.Real(8.1),
				"true": val.Null(),
			})),
		newRow(
			Map(Records{
				Record{Concat(Text("text"), Int(1)), Int(81)},
				Record{Add(Int(8), Int(1)), Sub(Real(8.1), Real(.1))},
				Record{And(True(), False()), Null()},
			}),
			val.Map(map[string]val.Val{
				"text1": val.Int(81),
				"9":     val.Real(8),
				"false": val.Null(),
			})),
	).exec(t, nil)
}

func Test_Ident(t *testing.T) {
	newRows(
		newRow(Ident("age"), val.Int(23)),
		newRow(Ident("num"), val.Real(2.3)),
		newRow(Add(Ident("num"), Ident("age")), val.Real(25.3)),
	).exec(t, map[string]val.Val{
		"age": val.Int(23),
		"num": val.Real(2.3),
	})
}

func Test_Commands(t *testing.T) {
	newRows(
		newRow(Commands([]Node{
			Add(Int(25), Real(2)),
			Add(Int(25), Real(.3)),
		}), val.Real(25.3)),
		newRow(Commands([]Node{
			Add(Int(25), Real(.3)),
			Add(Int(25), Int(2)),
		}), val.Int(27)),
	).exec(t, nil)
}

func Test_Call(t *testing.T) {
	newRows(
		newRow(Call(Ident("factorial"), []Node{Int(5)}), val.Int(120)),
	).exec(t, map[string]val.Val{
		"factorial": val.Fn(func(args []val.Val) (val.Val, error) {
			n := args[0].ToInt()

			factorial := func(n int64) int64 {
				if n <= 1 {
					return 1
				}
				res := int64(1)
				for i := int64(2); i <= n; i++ {
					res *= i
				}
				return res
			}

			return val.Int(factorial(n)), nil
		}, []string{"n"}),
	})
}

func Test_Assign(t *testing.T) {
	newRows(
		newRow(
			Commands([]Node{
				Assign(Ident("age"), Int(23)),
				Ident("age"),
			}), val.Int(23)),
		newRow(
			Commands([]Node{
				Assign(Ident("age"), Array([]Node{Int(23)})),
				Assign(IndexAccess(Ident("age"), Int(0)), Int(17)),
				IndexAccess(Ident("age"), Int(0)),
			}), val.Int(17)),
		newRow(
			Commands([]Node{
				Assign(Ident("user"), Map([]Record{{
					k: Text("name"),
					v: Text("sergey"),
				}})),
				Assign(IndexAccess(Ident("user"), Text("name")), Text("polina")),
				IndexAccess(Ident("user"), Text("name")),
			}), val.Text("polina")),
		newRow(
			Commands([]Node{
				Assign(Ident("user"), Array([]Node{
					Map([]Record{{k: Text("name"), v: Text("sergey")}}),
				})),
			}), val.Array([]val.Val{val.Map(map[string]val.Val{"name": val.Text("sergey")})}),
		),
		newRow(
			Commands([]Node{
				Assign(Ident("users"), Array([]Node{
					Map([]Record{{k: Text("name"), v: Text("sergey")}}),
				})),
				Assign(
					IndexAccess(
						IndexAccess(Ident("users"), Int(0)),
						Text("name"),
					), Text("polina")),
				IndexAccess(
					IndexAccess(Ident("users"), Int(0)),
					Text("name"),
				),
			}), val.Text("polina"),
		),
	).exec(t, nil)
}

func Test_If(t *testing.T) {
	newRows(
		newRow(
			If(NewBranch(
				Lt(Int(8), Int(27)), Text("8 меньше 27")),
				nil,
				nil), val.Text("8 меньше 27"),
		),
		newRow(
			If(
				NewBranch(Lt(Int(64), Int(27)), Text("64 меньше 27")),
				nil,
				NewBranch(nil, Text("64 не меньше 27"))), val.Text("64 не меньше 27"),
		),
		newRow(
			If(
				NewBranch(Lt(Int(64), Int(27)), Text("64 меньше 27")),
				[]*Branch{},
				NewBranch(nil, Text("64 не меньше 27"))), val.Text("64 не меньше 27"),
		),
		newRow(
			If(
				NewBranch(Lt(Int(64), Int(27)), Text("64 меньше 27")),
				[]*Branch{
					NewBranch(Lt(Int(8), Int(27)), Text("8 меньше 27")),
				},
				NewBranch(nil, Text("64 не меньше 27"))), val.Text("8 меньше 27"),
		),
		newRow(
			If(
				NewBranch(Lt(Int(64), Int(27)), Text("64 меньше 27")),
				[]*Branch{
					NewBranch(Eq(Int(8), Int(27)), Text("8 равен 27")),
					NewBranch(Lt(Int(8), Int(27)), Text("8 меньше 27")),
				},
				NewBranch(nil, Text("64 не меньше 27"))), val.Text("8 меньше 27"),
		),
	).exec(t, nil)
}

func Test_IndexAccess(t *testing.T) {
	newRows(
		newRow(
			IndexAccess(
				Array([]Node{
					Int(512),
					Int(6561),
					Int(16384),
				}), Int(1)),
			val.Int(6561)),
		newRow(
			IndexAccess(
				Array([]Node{
					Int(512),
					Int(6561),
					Int(16384),
				}), Int(2)),
			val.Int(16384)),
		newRow(IndexAccess(Text("text"), Int(2)), val.Text("x")),
		newRow(IndexAccess(Map(Records{
			Record{Concat(Text("text"), Int(1)), Int(81)},
			Record{Add(Int(8), Int(1)), Sub(Real(8.1), Real(.1))},
			Record{And(True(), False()), Null()},
		}), Text("text1")), val.Int(81)),
	).exec(t, nil)
}

func Test_Expr(t *testing.T) {
	newRows(
		newRow(
			Add(
				Mul(Val(val.Int(27)), Val(val.Int(81))),
				Div(Val(val.Int(25)), Val(val.Int(5))),
			),
			val.Int(2192),
		),
		newRow(
			Sub(
				Mul(Val(val.Int(27)), Val(val.Int(81))),
				Mul(Val(val.Int(256)), Val(val.Int(6))),
			),
			val.Int(651),
		),
		newRow(
			Eq(
				Add(Val(val.Int(2)), Val(val.Int(2))),
				Mul(Val(val.Int(2)), Val(val.Int(2))),
			),
			val.True(),
		),
		newRow(
			Add(
				Sub(Val(val.Int(81)), Val(val.Int(27))),
				Add(
					Mul(Val(val.Int(256)), Val(val.Int(2))),
					Div(Val(val.Int(256)), Val(val.Int(8))),
				),
			),
			val.Int(598),
		),
		newRow(
			Eq(
				Add(
					Sub(Val(val.Int(81)), Val(val.Int(27))),
					Add(
						Mul(Val(val.Int(256)), Val(val.Int(2))),
						Div(Val(val.Int(256)), Val(val.Int(8))),
					),
				),
				Val(val.Int(598)),
			),
			val.True(),
		),
		newRow(
			Lt(
				Add(
					Sub(Val(val.Int(81)), Val(val.Int(27))),
					Add(
						Mul(Val(val.Int(256)), Val(val.Int(2))),
						Div(Val(val.Int(256)), Val(val.Int(8))),
					),
				),
				Val(val.Int(600)),
			),
			val.True(),
		),
	).exec(t, nil)
}

func Test_For(t *testing.T) {
	newRows(
		newRow(For(Ident("i"), nil, Int(6), Ident("i")), val.Int(5)),
		newRow(For(Ident("i"), Ident("j"), Array([]Node{
			Int(256),
			Real(2.56),
		}), Concat(Ident("i"), Ident("j"))), val.Text("12.56")),
	).exec(t, nil)
}

func Test_DeepEqual(t *testing.T) {
	tests := []struct {
		a, b     Node
		expected bool
	}{
		{True(), True(), true},
		{False(), False(), true},
		{Null(), Null(), true},
		{Int(256), Int(256), true},
		{Int(256), Int(600), false},
		{Int(256), Real(2.56), false},
		{
			Add(Int(256), Real(2.56)),
			Add(Int(256), Real(2.56)),
			true,
		},
		{
			Add(Int(256), Real(2.56)),
			Mul(Int(256), Real(2.56)),
			false,
		},
		{
			Add(Int(256), Real(2.56)),
			Add(Real(2.56), Int(256)),
			false,
		},
		{
			Add(Int(256), Real(2.56)),
			Mul(Int(256), Real(2.56)),
			false,
		},
		{
			Add(Int(256), Real(2.56)),
			Add(Int(256), Int(256)),
			false,
		},
		{
			Mul(Add(Int(256), Real(2.56)), Real(2.56)),
			Mul(Add(Int(256), Real(2.56)), Real(2.56)),
			true,
		},
		{
			Mul(Add(Int(256), Real(2.56)), Real(2.56)),
			Mul(Sub(Int(256), Real(2.56)), Real(2.56)),
			false,
		},
		{True(), False(), false},
		{Null(), False(), false},
		{Null(), Int(256), false},
		{
			Add(Int(256), Real(2.56)),
			Int(256),
			false,
		},
		{
			Add(Int(256), Real(2.56)),
			Not(True()),
			false,
		},
		{Int(256), Not(True()), false},
		{Int(256), nil, false},
		{nil, nil, true},
		{Not(True()), Not(True()), true},
		{Array([]Node{Int(108), Text("text")}),
			Array([]Node{Int(108), Text("text")}), true},
		{Array([]Node{Int(108), Text("text")}),
			Array([]Node{Int(8), Text("text")}), false},
		{Array([]Node{Int(108), Text("text")}),
			Array([]Node{Int(108)}), false},
		{Array([]Node{Add(Int(108), Text("108"))}),
			Array([]Node{Add(Int(108), Text("108"))}),
			true},
		{Array([]Node{Add(Int(108), Text("108"))}),
			Array([]Node{Add(Int(108), Text("108")), Add(Int(108), Text("108"))}),
			false},
		{Map(Records{Record{Text("text"), Int(108)}}),
			Map(Records{Record{Text("text"), Int(108)}}),
			true,
		},
		{Map(Records{Record{Text("text"), Int(108)}}),
			Map(Records{Record{Text("txt"), Int(108)}}),
			false,
		},
		{Map(Records{Record{Text("text"), Int(108)}}),
			Map(Records{Record{Text("text"), Int(18)}}),
			false,
		},
		{
			Map(Records{
				Record{
					Concat(Text("text"), Text("false")),
					Add(Int(108), Int(8))}}),
			Map(Records{
				Record{
					Concat(Text("text"), Text("false")),
					Add(Int(108), Int(8))}}),
			true,
		},
		{
			Map(Records{
				Record{
					Concat(Text("text"), Text("false")),
					Add(Int(108), Int(8))},
				Record{
					Concat(Text("text"), Text("false")),
					Add(Int(108), Int(8))},
			}),
			Map(Records{
				Record{
					Concat(Text("text"), Text("false")),
					Add(Int(108), Int(8))}}),
			false,
		},
		{
			Commands([]Node{
				Add(Int(25), Real(2)),
				Add(Int(25), Real(.3)),
			}),
			Commands([]Node{
				Add(Int(25), Real(2)),
				Add(Int(25), Real(.3)),
			}),
			true,
		},
		{
			Commands([]Node{
				Add(Int(25), Real(2)),
				Add(Int(25), Real(.3)),
			}),
			Commands([]Node{
				Add(Int(25), Real(.3)),
				Add(Int(25), Real(2)),
			}),
			false,
		},
		{
			Commands([]Node{
				Add(Int(25), Real(2)),
				Add(Int(25), Real(.3)),
			}),
			Commands([]Node{
				Add(Int(25), Real(2)),
			}),
			false,
		},
		{
			If(
				NewBranch(Lt(Int(64), Int(27)), Text("64 меньше 27")),
				[]*Branch{
					NewBranch(Eq(Int(8), Int(27)), Text("8 равен 27")),
					NewBranch(Lt(Int(8), Int(27)), Text("8 меньше 27")),
				},
				NewBranch(nil, Text("64 не меньше 27"))),
			If(
				NewBranch(Lt(Int(64), Int(27)), Text("64 меньше 27")),
				[]*Branch{
					NewBranch(Eq(Int(8), Int(27)), Text("8 равен 27")),
					NewBranch(Lt(Int(8), Int(27)), Text("8 меньше 27")),
				},
				NewBranch(nil, Text("64 не меньше 27"))),
			true,
		},
		{
			If(
				NewBranch(Gt(Int(64), Int(27)), Text("64 меньше 27")),
				[]*Branch{
					NewBranch(Eq(Int(8), Int(27)), Text("8 равен 27")),
					NewBranch(Lt(Int(8), Int(27)), Text("8 меньше 27")),
				},
				NewBranch(nil, Text("64 не меньше 27"))),
			If(
				NewBranch(Lt(Int(64), Int(27)), Text("64 меньше 27")),
				[]*Branch{
					NewBranch(Eq(Int(8), Int(27)), Text("8 равен 27")),
					NewBranch(Lt(Int(8), Int(27)), Text("8 меньше 27")),
				},
				NewBranch(nil, Text("64 не меньше 27"))),
			false,
		},
		{
			If(
				NewBranch(Lt(Int(64), Int(27)), Text("64 меньше 27")),
				[]*Branch{
					NewBranch(Eq(Int(8), Int(27)), Text("8 равен 27")),
				},
				NewBranch(nil, Text("64 не меньше 27"))),
			If(
				NewBranch(Lt(Int(64), Int(27)), Text("64 меньше 27")),
				[]*Branch{
					NewBranch(Eq(Int(8), Int(27)), Text("8 равен 27")),
					NewBranch(Lt(Int(8), Int(27)), Text("8 меньше 27")),
				},
				NewBranch(nil, Text("64 не меньше 27"))),
			false,
		},
		{
			If(
				NewBranch(Lt(Int(64), Int(27)), Text("64 меньше 27")),
				[]*Branch{
					NewBranch(Eq(Int(8), Int(27)), Text("8 равен 27")),
					NewBranch(Lt(Int(8), Int(27)), Text("8 меньше 27")),
				},
				nil),
			If(
				NewBranch(Lt(Int(64), Int(27)), Text("64 меньше 27")),
				[]*Branch{
					NewBranch(Eq(Int(8), Int(27)), Text("8 равен 27")),
					NewBranch(Lt(Int(8), Int(27)), Text("8 меньше 27")),
				},
				NewBranch(nil, Text("64 не меньше 27"))),
			false,
		},
		{
			For(Ident("i"), Ident("j"), Array([]Node{
				Int(256),
				Real(2.56),
			}), Concat(Ident("i"), Ident("j"))),
			For(Ident("i"), Ident("j"), Array([]Node{
				Int(256),
				Real(2.56),
			}), Concat(Ident("i"), Ident("j"))),
			true,
		},
		{
			For(Ident("i"), Ident("j"), Array([]Node{
				Int(256),
				Real(2.56),
			}), Concat(Ident("i"), Ident("j"))),
			For(Ident("i"), nil, Array([]Node{
				Int(256),
				Real(2.56),
			}), Concat(Ident("i"), Ident("j"))),
			false,
		},
		{Call(Ident("factorial"), []Node{Int(5)}),
			Call(Ident("factorial"), []Node{Int(5)}),
			true,
		},
	}

	for i, test := range tests {
		if DeepEqual(test.a, test.b) != test.expected {
			t.Errorf("%d: ожидалось %t, получено %t",
				i, test.expected, !test.expected)
		}
	}
}
