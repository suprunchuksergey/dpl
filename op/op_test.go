package op

import (
	. "github.com/suprunchuksergey/dpl/val"
	"testing"
)

type row struct {
	a, b, expected Val
}

type rows []row

func (rows rows) exec(t *testing.T, f Binary) {
	for i, r := range rows {
		v := f(r.a, r.b)
		if v != r.expected {
			t.Errorf("%d: ожидалось %s, получено %s",
				i, r.expected, v)
		}
	}
}

func Test_Add(t *testing.T) {
	rows{
		{Int(8), Int(81), Int(89)},
		{Int(-8), Int(81), Int(73)},
		{Int(8), Text("81 text"), Int(89)},
		{Int(-8), Text("81 text"), Int(73)},
		{Int(-8), Text("-81 text"), Int(-89)},
		{Int(8), Real(81.881), Real(89.881)},
		{Int(8), Text("81.881"), Real(89.881)},
		{Int(8), True(), Int(9)},
		{Int(8), False(), Int(8)},
		{Int(81), Int(8), Int(89)},
		{Int(81), Int(-8), Int(73)},
		{Text("81 text"), Int(8), Int(89)},
		{Text("81 text"), Int(-8), Int(73)},
		{Text("-81 text"), Int(-8), Int(-89)},
		{Real(81.881), Int(8), Real(89.881)},
		{Text("81.881"), Int(8), Real(89.881)},
		{True(), Int(8), Int(9)},
		{False(), Int(8), Int(8)},
		{Text("81 text"), Text("216 text"), Int(297)},
		{Text("81 text"), Text("text"), Int(81)},
		{Text("text"), Text("216 text"), Int(216)},
		{Text("81 text"), Text("-216 text"), Int(-135)},
		{Text("-81 text"), Text("21.6 text"), Real(-59.4)},
		{Text("+.81 text"), Text("21.6 text"), Real(22.41)},
		{Text("text"), True(), Int(1)},
		{Text("text"), False(), Int(0)},
		{False(), False(), Int(0)},
		{True(), False(), Int(1)},
		{True(), True(), Int(2)},
		{Int(8), Null(), Int(8)},
	}.exec(t, Add)
}

func Test_Sub(t *testing.T) {
	rows{
		{Int(343), Int(27), Int(316)},
		{Int(343), Int(-27), Int(370)},
		{Int(343), Text("27 text"), Int(316)},
		{Int(343), Text("-27text"), Int(370)},
		{Int(343), Real(2.7), Real(340.3)},
		{Int(343), Real(-2.7), Real(345.7)},
		{Int(343), False(), Int(343)},
		{Int(343), True(), Int(342)},
		{False(), False(), Int(0)},
		{True(), False(), Int(1)},
		{True(), True(), Int(0)},
		{Text("343"), Text("	 2.7text"), Real(340.3)},
		{Text("343"), Text("-2.7text"), Real(345.7)},
		{Text("3.43"), Text("-2.7text"), Real(6.130000000000001)},
		{Text("text"), Text("text"), Int(0)},
		{Text("text"), Text("343"), Int(-343)},
		{True(), Text("343"), Int(-342)},
		{Text(".9"), Text(".1"), Real(.8)},
		{Text(".9"), Text("-.1"), Real(1)},
		{Int(8), Null(), Int(8)},
	}.exec(t, Sub)
}

func Test_Mul(t *testing.T) {
	rows{
		{Int(64), Int(16), Int(1024)},
		{Int(64), Int(-16), Int(-1024)},
		{Int(64), Text("16"), Int(1024)},
		{Int(64), Text("-16"), Int(-1024)},
		{Int(64), Real(1.6), Real(102.4)},
		{Text("64"), Text("16"), Int(1024)},
		{Text("64"), Text("-16"), Int(-1024)},
		{Text("64"), Text("1.6"), Real(102.4)},
		{Int(64), False(), Int(0)},
		{Int(64), True(), Int(64)},
		{True(), True(), Int(1)},
		{False(), True(), Int(0)},
		{Text(".5"), Text("-.3"), Real(-.15)},
		{Int(64), Int(0), Int(0)},
		{Int(0), Int(0), Int(0)},
		{Int(64), Text("text"), Int(0)},
		{Text("text"), Text("text"), Int(0)},
		{Int(16), Int(64), Int(1024)},
		{Int(-16), Int(64), Int(-1024)},
		{Text("16"), Int(64), Int(1024)},
		{Text("-16"), Int(64), Int(-1024)},
		{Real(1.6), Int(64), Real(102.4)},
		{Text("16"), Text("64"), Int(1024)},
		{Text("-16"), Text("64"), Int(-1024)},
		{Text("1.6"), Text("64"), Real(102.4)},
		{False(), Int(64), Int(0)},
		{True(), Int(64), Int(64)},
		{True(), True(), Int(1)},
		{True(), False(), Int(0)},
		{Text("-.3"), Text(".5"), Real(-.15)},
		{Int(0), Int(64), Int(0)},
		{Int(0), Int(0), Int(0)},
		{Text("text"), Int(64), Int(0)},
		{Text("text"), Text("text"), Int(0)},
		{Real(1.6), Real(6.4), Real(10.240000000000002)},
		{Int(8), Null(), Int(0)},
	}.exec(t, Mul)
}

func Test_Div(t *testing.T) {
	rows{
		{Int(64), Int(16), Int(4)},
		{Int(64), Int(-16), Int(-4)},
		{Int(64), Text("	16 text"), Int(4)},
		{Int(64), Text("-16text"), Int(-4)},
		{Int(64), Real(1.6), Real(40)},
		{Int(64), Real(-.16), Real(-400)},
		{Text("	16 text"), Text("	4 text"), Int(4)},
		{Text("	64 text"), Text("1.6text"), Real(40)},
		{Int(0), Int(64), Int(0)},
		{Int(64), Int(0), Null()},
		{Int(64), Real(0), Null()},
		{Int(64), False(), Null()},
		{Int(64), Text("text"), Null()},
		{Int(64), Text("	-0"), Null()},
		{Int(64), Text("	+.0"), Null()},
		{Int(64), Null(), Null()},
		{Null(), Null(), Null()},
		{False(), True(), Int(0)},
		{True(), True(), Int(1)},
		{Real(6.4), Real(.6), Real(10.666666666666668)},
		{Real(6.4), Text("	.6.text"), Real(10.666666666666668)},
	}.exec(t, Div)
}

func Test_Rem(t *testing.T) {
	rows{
		{Int(625), Int(9), Int(4)},
		{Int(625), Int(-9), Int(4)},
		{Int(-625), Int(9), Int(-4)},
		{Int(625), Text("9"), Int(4)},
		{Int(625), Text("-9"), Int(4)},
		{Int(-625), Text("9"), Int(-4)},
		{Real(6.25), Real(.9), Real(.8499999999999999)},
		{Real(6.25), Text(".9"), Real(.8499999999999999)},
		{Real(62.5), Int(-9), Real(8.5)},
		{Real(-.625), Int(4), Real(-.625)},
		{Int(6), Int(2), Int(0)},
		{Int(7), Int(2), Int(1)},
		{Real(10.5), Int(6), Real(4.5)},
		{Text("6.25"), Text(".9"), Real(.8499999999999999)},
		{Text("	62.5"), Int(-9), Real(8.5)},
		{Text("	 -.625"), Int(4), Real(-.625)},
		{Int(625), Int(0), Null()},
		{Int(625), Real(0), Null()},
		{Int(625), Text("	0"), Null()},
		{Int(625), Text("text"), Null()},
		{Int(625), Text(".0"), Null()},
		{Int(625), Text("."), Null()},
		{Int(625), False(), Null()},
		{Int(0), Int(625), Int(0)},
		{Real(.9), Real(.5), Real(.4)},
		{Int(25), Int(2), Int(1)},
	}.exec(t, Rem)
}

func Test_Eq(t *testing.T) {
	rows{
		{Int(25), Int(25), True()},
		{Int(25), Real(25), True()},
		{Int(25), Text("25"), True()},
		{Int(25), Text("25.text"), True()},
		{Int(1), True(), True()},
		{Int(0), False(), True()},
		{Text("text"), Text("text"), True()},
		{True(), True(), True()},
		{False(), False(), True()},
		{Real(2.5), Real(2.5), True()},
		{Real(2.5), Text("2.5"), True()},
		{Int(25), Int(5), False()},
		{Int(25), Real(5), False()},
		{Int(25), Text("5"), False()},
		{Int(25), Text("2.text"), False()},
		{Int(1), False(), False()},
		{Int(0), True(), False()},
		{Text("0text"), Text("text"), False()},
		{True(), False(), False()},
		{False(), True(), False()},
		{Real(2.5), Real(25), False()},
		{Real(2.5), Text(".25"), False()},
		{Real(25), Text("25.1"), False()},
		{Int(2187), Text("2187 рублей"), True()},
		{Int(0), Null(), True()},
	}.exec(t, Eq)
}

func Test_Neq(t *testing.T) {
	rows{
		{Int(25), Int(25), False()},
		{Int(25), Real(25), False()},
		{Int(25), Text("25"), False()},
		{Int(25), Text("25.text"), False()},
		{Int(1), True(), False()},
		{Int(0), False(), False()},
		{Text("text"), Text("text"), False()},
		{True(), True(), False()},
		{False(), False(), False()},
		{Real(2.5), Real(2.5), False()},
		{Real(2.5), Text("2.5"), False()},
		{Int(25), Int(5), True()},
		{Int(25), Real(5), True()},
		{Int(25), Text("5"), True()},
		{Int(25), Text("2.text"), True()},
		{Int(1), False(), True()},
		{Int(0), True(), True()},
		{Text("0text"), Text("text"), True()},
		{True(), False(), True()},
		{False(), True(), True()},
		{Real(2.5), Real(25), True()},
		{Real(2.5), Text(".25"), True()},
		{Real(25), Text("25.1"), True()},
		{Int(512), Text("2187 рублей"), True()},
		{Int(0), Null(), False()},
	}.exec(t, Neq)
}

func Test_Lt(t *testing.T) {
	rows{
		{Int(512), Int(2187), True()},
		{Int(-512), Int(2187), True()},
		{Int(512), Text("2187"), True()},
		{Int(-512), Text("2187"), True()},
		{Int(512), Real(512.1), True()},
		{False(), True(), True()},
		{Text("hello"), Text("world"), True()},
		{Text("512hello"), Text("512world"), True()},
		{Int(2187), Int(512), False()},
		{Int(2187), Int(-512), False()},
		{Text("2187"), Int(512), False()},
		{Text("2187"), Int(-512), False()},
		{Real(512.1), Int(512), False()},
		{True(), False(), False()},
		{Text("world"), Text("hello"), False()},
		{Text("512world"), Text("512hello"), False()},
		{Int(512), Int(512), False()},
		{Int(512), Text("2187 рублей"), True()},
		{Null(), Int(512), True()},
	}.exec(t, Lt)
}

func Test_Lte(t *testing.T) {
	rows{
		{Int(512), Int(2187), True()},
		{Int(-512), Int(2187), True()},
		{Int(512), Text("2187"), True()},
		{Int(-512), Text("2187"), True()},
		{Int(512), Real(512.1), True()},
		{False(), True(), True()},
		{Text("hello"), Text("world"), True()},
		{Text("512hello"), Text("512world"), True()},
		{Int(2187), Int(512), False()},
		{Int(2187), Int(-512), False()},
		{Text("2187"), Int(512), False()},
		{Text("2187"), Int(-512), False()},
		{Real(512.1), Int(512), False()},
		{True(), False(), False()},
		{Text("world"), Text("hello"), False()},
		{Text("512world"), Text("512hello"), False()},
		{Int(512), Int(512), True()},
		{Int(512), Text("2187 рублей"), True()},
		{Null(), Int(512), True()},
	}.exec(t, Lte)
}

func Test_Gt(t *testing.T) {
	rows{
		{Int(512), Int(2187), False()},
		{Int(-512), Int(2187), False()},
		{Int(512), Text("2187"), False()},
		{Int(-512), Text("2187"), False()},
		{Int(512), Real(512.1), False()},
		{False(), True(), False()},
		{Text("hello"), Text("world"), False()},
		{Text("512hello"), Text("512world"), False()},
		{Int(2187), Int(512), True()},
		{Int(2187), Int(-512), True()},
		{Text("2187"), Int(512), True()},
		{Text("2187"), Int(-512), True()},
		{Real(512.1), Int(512), True()},
		{True(), False(), True()},
		{Text("world"), Text("hello"), True()},
		{Text("512world"), Text("512hello"), True()},
		{Int(512), Int(512), False()},
		{Int(2187), Text("512 рублей"), True()},
		{Int(2187), Null(), True()},
	}.exec(t, Gt)
}

func Test_Gte(t *testing.T) {
	rows{
		{Int(512), Int(2187), False()},
		{Int(-512), Int(2187), False()},
		{Int(512), Text("2187"), False()},
		{Int(-512), Text("2187"), False()},
		{Int(512), Real(512.1), False()},
		{False(), True(), False()},
		{Text("hello"), Text("world"), False()},
		{Text("512hello"), Text("512world"), False()},
		{Int(2187), Int(512), True()},
		{Int(2187), Int(-512), True()},
		{Text("2187"), Int(512), True()},
		{Text("2187"), Int(-512), True()},
		{Real(512.1), Int(512), True()},
		{True(), False(), True()},
		{Text("world"), Text("hello"), True()},
		{Text("512world"), Text("512hello"), True()},
		{Int(512), Int(512), True()},
		{Int(2187), Text("512 рублей"), True()},
		{Int(2187), Null(), True()},
	}.exec(t, Gte)
}

func Test_Concat(t *testing.T) {
	rows{
		{Int(512), Int(2187), Text("5122187")},
		{Real(5.12), Int(2187), Text("5.122187")},
		{Real(5.12), Real(2.187), Text("5.122.187")},
		{Real(5.12), Text("рублей"), Text("5.12рублей")},
		{True(), False(), Text("truefalse")},
		{Real(.12), Text("рублей"), Text("0.12рублей")},
		{Text("hello "), Text("world"), Text("hello world")},
		{Text("hello"), Null(), Text("hello")},
	}.exec(t, Concat)
}

func Test_And(t *testing.T) {
	rows{
		{True(), True(), True()},
		{True(), False(), False()},
		{False(), True(), False()},
		{False(), False(), False()},
		{Int(512), Text("text"), True()},
		{Int(0), Text("text"), False()},
		{Int(512), Text(""), False()},
		{Int(0), Text(""), False()},
		{Null(), Text("text"), False()},
	}.exec(t, And)
}

func Test_Or(t *testing.T) {
	rows{
		{True(), True(), True()},
		{True(), False(), True()},
		{False(), True(), True()},
		{False(), False(), False()},
		{Int(512), Text("text"), True()},
		{Int(0), Text("text"), True()},
		{Int(512), Text(""), True()},
		{Int(0), Text(""), False()},
		{Null(), Text("text"), True()},
	}.exec(t, Or)
}

func Test_Not(t *testing.T) {
	tests := []struct{ v, expected Val }{
		{True(), False()},
		{False(), True()},
		{Int(512), False()},
		{Int(0), True()},
		{Text("text"), False()},
		{Text(""), True()},
		{Null(), True()},
	}

	for i, test := range tests {
		v := Not(test.v)
		if v != test.expected {
			t.Errorf("%d: ожидалось %s, получено %s",
				i, test.expected, v)
		}
	}
}

func Test_Neg(t *testing.T) {
	tests := []struct{ v, expected Val }{
		{Int(0), Int(0)},
		{Int(512), Int(-512)},
		{Real(5.12), Real(-5.12)},
		{True(), Int(-1)},
		{False(), Int(0)},
		{Text("text"), Int(0)},
		{Text("512text"), Int(-512)},
		{Text("5.12text"), Real(-5.12)},
	}

	for i, test := range tests {
		v := Neg(test.v)
		if v != test.expected {
			t.Errorf("%d: ожидалось %s, получено %s",
				i, test.expected, v)
		}
	}
}

func Test_IndexAccess(t *testing.T) {
	tests := []struct {
		v,
		index,
		expected Val
	}{
		{Array([]Val{Text("hello"), Int(512), True()}),
			Int(0), Text("hello")},
		{Array([]Val{Text("hello"), Int(512), True()}),
			Int(1), Int(512)},
		{Array([]Val{Text("hello"), Int(512), True()}),
			Int(2), True()},
		{Text("hello"), Int(0), Text("h")},
		{Text("hello"), Int(1), Text("e")},
		{Text("hello"), Int(2), Text("l")},
		{Map(map[string]Val{
			"hello": Int(512),
			"5":     True(),
		}), Text("hello"), Int(512)},
		{Map(map[string]Val{
			"hello": Int(512),
			"5":     True(),
		}), Int(5), True()},
	}
	for i, test := range tests {
		v, err := IndexAccess(test.v, test.index)
		if err != nil {
			return
		}
		if test.expected != v {
			t.Errorf("%d: ожидалось %s, получено %s",
				i, test.expected, v)
		}
	}
}
