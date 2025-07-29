package op

import (
	"github.com/suprunchuksergey/dpl/val"
	"testing"
)

type row struct {
	a, b, expected val.Val
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
		{val.Int(8), val.Int(81), val.Int(89)},
		{val.Int(-8), val.Int(81), val.Int(73)},
		{val.Int(8), val.Text("81 text"), val.Int(89)},
		{val.Int(-8), val.Text("81 text"), val.Int(73)},
		{val.Int(-8), val.Text("-81 text"), val.Int(-89)},
		{val.Int(8), val.Real(81.881), val.Real(89.881)},
		{val.Int(8), val.Text("81.881"), val.Real(89.881)},
		{val.Int(8), val.True(), val.Int(9)},
		{val.Int(8), val.False(), val.Int(8)},
		{val.Int(8), val.Null(), val.Null()},
		{val.Int(81), val.Int(8), val.Int(89)},
		{val.Int(81), val.Int(-8), val.Int(73)},
		{val.Text("81 text"), val.Int(8), val.Int(89)},
		{val.Text("81 text"), val.Int(-8), val.Int(73)},
		{val.Text("-81 text"), val.Int(-8), val.Int(-89)},
		{val.Real(81.881), val.Int(8), val.Real(89.881)},
		{val.Text("81.881"), val.Int(8), val.Real(89.881)},
		{val.True(), val.Int(8), val.Int(9)},
		{val.False(), val.Int(8), val.Int(8)},
		{val.Null(), val.Int(8), val.Null()},
		{val.Text("81 text"), val.Text("216 text"), val.Int(297)},
		{val.Text("81 text"), val.Text("text"), val.Int(81)},
		{val.Text("text"), val.Text("216 text"), val.Int(216)},
		{val.Text("81 text"), val.Text("-216 text"), val.Int(-135)},
		{val.Text("-81 text"), val.Text("21.6 text"), val.Real(-59.4)},
		{val.Text("+.81 text"), val.Text("21.6 text"), val.Real(22.41)},
		{val.Text("text"), val.True(), val.Int(1)},
		{val.Text("text"), val.False(), val.Int(0)},
		{val.False(), val.False(), val.Int(0)},
		{val.True(), val.False(), val.Int(1)},
		{val.True(), val.True(), val.Int(2)},
	}.exec(t, Add)
}

func Test_Sub(t *testing.T) {
	rows{
		{val.Int(343), val.Int(27), val.Int(316)},
		{val.Int(343), val.Int(-27), val.Int(370)},
		{val.Int(343), val.Text("27 text"), val.Int(316)},
		{val.Int(343), val.Text("-27text"), val.Int(370)},
		{val.Int(343), val.Real(2.7), val.Real(340.3)},
		{val.Int(343), val.Real(-2.7), val.Real(345.7)},
		{val.Null(), val.Real(-2.7), val.Null()},
		{val.Int(343), val.Null(), val.Null()},
		{val.Int(343), val.False(), val.Int(343)},
		{val.Int(343), val.True(), val.Int(342)},
		{val.False(), val.False(), val.Int(0)},
		{val.True(), val.False(), val.Int(1)},
		{val.True(), val.True(), val.Int(0)},
		{val.Text("343"), val.Text("	 2.7text"), val.Real(340.3)},
		{val.Text("343"), val.Text("-2.7text"), val.Real(345.7)},
		{val.Text("3.43"), val.Text("-2.7text"), val.Real(6.130000000000001)},
		{val.Text("text"), val.Text("text"), val.Int(0)},
		{val.Text("text"), val.Text("343"), val.Int(-343)},
		{val.True(), val.Text("343"), val.Int(-342)},
		{val.Text(".9"), val.Text(".1"), val.Real(.8)},
		{val.Text(".9"), val.Text("-.1"), val.Real(1)},
	}.exec(t, Sub)
}

func Test_Mul(t *testing.T) {
	rows{
		{val.Int(64), val.Int(16), val.Int(1024)},
		{val.Int(64), val.Int(-16), val.Int(-1024)},
		{val.Int(64), val.Text("16"), val.Int(1024)},
		{val.Int(64), val.Text("-16"), val.Int(-1024)},
		{val.Int(64), val.Real(1.6), val.Real(102.4)},
		{val.Int(64), val.Null(), val.Null()},
		{val.Text("64"), val.Text("16"), val.Int(1024)},
		{val.Text("64"), val.Text("-16"), val.Int(-1024)},
		{val.Text("64"), val.Text("1.6"), val.Real(102.4)},
		{val.Int(64), val.False(), val.Int(0)},
		{val.Int(64), val.True(), val.Int(64)},
		{val.True(), val.True(), val.Int(1)},
		{val.False(), val.True(), val.Int(0)},
		{val.Text(".5"), val.Text("-.3"), val.Real(-.15)},
		{val.Int(64), val.Int(0), val.Int(0)},
		{val.Int(0), val.Int(0), val.Int(0)},
		{val.Int(64), val.Text("text"), val.Int(0)},
		{val.Text("text"), val.Text("text"), val.Int(0)},
		{val.Int(16), val.Int(64), val.Int(1024)},
		{val.Int(-16), val.Int(64), val.Int(-1024)},
		{val.Text("16"), val.Int(64), val.Int(1024)},
		{val.Text("-16"), val.Int(64), val.Int(-1024)},
		{val.Real(1.6), val.Int(64), val.Real(102.4)},
		{val.Null(), val.Int(64), val.Null()},
		{val.Text("16"), val.Text("64"), val.Int(1024)},
		{val.Text("-16"), val.Text("64"), val.Int(-1024)},
		{val.Text("1.6"), val.Text("64"), val.Real(102.4)},
		{val.False(), val.Int(64), val.Int(0)},
		{val.True(), val.Int(64), val.Int(64)},
		{val.True(), val.True(), val.Int(1)},
		{val.True(), val.False(), val.Int(0)},
		{val.Text("-.3"), val.Text(".5"), val.Real(-.15)},
		{val.Int(0), val.Int(64), val.Int(0)},
		{val.Int(0), val.Int(0), val.Int(0)},
		{val.Text("text"), val.Int(64), val.Int(0)},
		{val.Text("text"), val.Text("text"), val.Int(0)},
		{val.Null(), val.Null(), val.Null()},
		{val.Real(1.6), val.Real(6.4), val.Real(10.240000000000002)},
	}.exec(t, Mul)
}

func Test_Div(t *testing.T) {
	rows{
		{val.Int(64), val.Int(16), val.Int(4)},
		{val.Int(64), val.Int(-16), val.Int(-4)},
		{val.Int(64), val.Text("	16 text"), val.Int(4)},
		{val.Int(64), val.Text("-16text"), val.Int(-4)},
		{val.Int(64), val.Real(1.6), val.Real(40)},
		{val.Int(64), val.Real(-.16), val.Real(-400)},
		{val.Text("	16 text"), val.Text("	4 text"), val.Int(4)},
		{val.Text("	64 text"), val.Text("1.6text"), val.Real(40)},
		{val.Int(0), val.Int(64), val.Int(0)},
		{val.Int(64), val.Int(0), val.Null()},
		{val.Int(64), val.Real(0), val.Null()},
		{val.Int(64), val.False(), val.Null()},
		{val.Int(64), val.Text("text"), val.Null()},
		{val.Int(64), val.Text("	-0"), val.Null()},
		{val.Int(64), val.Text("	+.0"), val.Null()},
		{val.Int(64), val.Null(), val.Null()},
		{val.Null(), val.Null(), val.Null()},
		{val.False(), val.True(), val.Int(0)},
		{val.True(), val.True(), val.Int(1)},
		{val.Real(6.4), val.Real(.6), val.Real(10.666666666666668)},
		{val.Real(6.4), val.Text("	.6.text"), val.Real(10.666666666666668)},
	}.exec(t, Div)
}

func Test_Rem(t *testing.T) {
	rows{
		{val.Int(625), val.Int(9), val.Int(4)},
		{val.Int(625), val.Int(-9), val.Int(4)},
		{val.Int(-625), val.Int(9), val.Int(-4)},
		{val.Int(625), val.Text("9"), val.Int(4)},
		{val.Int(625), val.Text("-9"), val.Int(4)},
		{val.Int(-625), val.Text("9"), val.Int(-4)},
		{val.Real(6.25), val.Real(.9), val.Real(.8499999999999999)},
		{val.Real(6.25), val.Text(".9"), val.Real(.8499999999999999)},
		{val.Real(62.5), val.Int(-9), val.Real(8.5)},
		{val.Real(-.625), val.Int(4), val.Real(-.625)},
		{val.Int(6), val.Int(2), val.Int(0)},
		{val.Int(7), val.Int(2), val.Int(1)},
		{val.Real(10.5), val.Int(6), val.Real(4.5)},
		{val.Text("6.25"), val.Text(".9"), val.Real(.8499999999999999)},
		{val.Text("	62.5"), val.Int(-9), val.Real(8.5)},
		{val.Text("	 -.625"), val.Int(4), val.Real(-.625)},
		{val.Int(625), val.Int(0), val.Null()},
		{val.Int(625), val.Real(0), val.Null()},
		{val.Int(625), val.Text("	0"), val.Null()},
		{val.Int(625), val.Text("text"), val.Null()},
		{val.Int(625), val.Text(".0"), val.Null()},
		{val.Int(625), val.Text("."), val.Null()},
		{val.Int(625), val.False(), val.Null()},
		{val.Int(0), val.Int(625), val.Int(0)},
		{val.Real(.9), val.Real(.5), val.Real(.4)},
		{val.Int(25), val.Int(2), val.Int(1)},
	}.exec(t, Rem)
}

func Test_Eq(t *testing.T) {
	rows{
		{val.Int(25), val.Int(25), val.True()},
		{val.Int(25), val.Real(25), val.True()},
		{val.Int(25), val.Text("25"), val.True()},
		{val.Int(25), val.Text("25.text"), val.True()},
		{val.Int(1), val.True(), val.True()},
		{val.Int(0), val.False(), val.True()},
		{val.Text("text"), val.Text("text"), val.True()},
		{val.True(), val.True(), val.True()},
		{val.False(), val.False(), val.True()},
		{val.Real(2.5), val.Real(2.5), val.True()},
		{val.Real(2.5), val.Text("2.5"), val.True()},
		{val.Int(25), val.Null(), val.Null()},
		{val.Int(25), val.Int(5), val.False()},
		{val.Int(25), val.Real(5), val.False()},
		{val.Int(25), val.Text("5"), val.False()},
		{val.Int(25), val.Text("2.text"), val.False()},
		{val.Int(1), val.False(), val.False()},
		{val.Int(0), val.True(), val.False()},
		{val.Text("0text"), val.Text("text"), val.False()},
		{val.True(), val.False(), val.False()},
		{val.False(), val.True(), val.False()},
		{val.Real(2.5), val.Real(25), val.False()},
		{val.Real(2.5), val.Text(".25"), val.False()},
		{val.Null(), val.Null(), val.Null()},
		{val.Real(25), val.Text("25.1"), val.False()},
		{val.Int(2187), val.Text("2187 рублей"), val.True()},
	}.exec(t, Eq)
}

func Test_Neq(t *testing.T) {
	rows{
		{val.Int(25), val.Int(25), val.False()},
		{val.Int(25), val.Real(25), val.False()},
		{val.Int(25), val.Text("25"), val.False()},
		{val.Int(25), val.Text("25.text"), val.False()},
		{val.Int(1), val.True(), val.False()},
		{val.Int(0), val.False(), val.False()},
		{val.Text("text"), val.Text("text"), val.False()},
		{val.True(), val.True(), val.False()},
		{val.False(), val.False(), val.False()},
		{val.Real(2.5), val.Real(2.5), val.False()},
		{val.Real(2.5), val.Text("2.5"), val.False()},
		{val.Int(25), val.Null(), val.Null()},
		{val.Int(25), val.Int(5), val.True()},
		{val.Int(25), val.Real(5), val.True()},
		{val.Int(25), val.Text("5"), val.True()},
		{val.Int(25), val.Text("2.text"), val.True()},
		{val.Int(1), val.False(), val.True()},
		{val.Int(0), val.True(), val.True()},
		{val.Text("0text"), val.Text("text"), val.True()},
		{val.True(), val.False(), val.True()},
		{val.False(), val.True(), val.True()},
		{val.Real(2.5), val.Real(25), val.True()},
		{val.Real(2.5), val.Text(".25"), val.True()},
		{val.Null(), val.Null(), val.Null()},
		{val.Real(25), val.Text("25.1"), val.True()},
		{val.Int(512), val.Text("2187 рублей"), val.True()},
	}.exec(t, Neq)
}

func Test_Lt(t *testing.T) {
	rows{
		{val.Int(512), val.Int(2187), val.True()},
		{val.Int(-512), val.Int(2187), val.True()},
		{val.Int(512), val.Text("2187"), val.True()},
		{val.Int(-512), val.Text("2187"), val.True()},
		{val.Int(512), val.Real(512.1), val.True()},
		{val.False(), val.True(), val.True()},
		{val.Text("hello"), val.Text("world"), val.True()},
		{val.Text("512hello"), val.Text("512world"), val.True()},
		{val.Int(2187), val.Int(512), val.False()},
		{val.Int(2187), val.Int(-512), val.False()},
		{val.Text("2187"), val.Int(512), val.False()},
		{val.Text("2187"), val.Int(-512), val.False()},
		{val.Real(512.1), val.Int(512), val.False()},
		{val.True(), val.False(), val.False()},
		{val.Text("world"), val.Text("hello"), val.False()},
		{val.Text("512world"), val.Text("512hello"), val.False()},
		{val.Int(512), val.Null(), val.Null()},
		{val.Null(), val.Int(512), val.Null()},
		{val.Int(512), val.Int(512), val.False()},
		{val.Int(512), val.Text("2187 рублей"), val.True()},
	}.exec(t, Lt)
}

func Test_Lte(t *testing.T) {
	rows{
		{val.Int(512), val.Int(2187), val.True()},
		{val.Int(-512), val.Int(2187), val.True()},
		{val.Int(512), val.Text("2187"), val.True()},
		{val.Int(-512), val.Text("2187"), val.True()},
		{val.Int(512), val.Real(512.1), val.True()},
		{val.False(), val.True(), val.True()},
		{val.Text("hello"), val.Text("world"), val.True()},
		{val.Text("512hello"), val.Text("512world"), val.True()},
		{val.Int(2187), val.Int(512), val.False()},
		{val.Int(2187), val.Int(-512), val.False()},
		{val.Text("2187"), val.Int(512), val.False()},
		{val.Text("2187"), val.Int(-512), val.False()},
		{val.Real(512.1), val.Int(512), val.False()},
		{val.True(), val.False(), val.False()},
		{val.Text("world"), val.Text("hello"), val.False()},
		{val.Text("512world"), val.Text("512hello"), val.False()},
		{val.Int(512), val.Null(), val.Null()},
		{val.Null(), val.Int(512), val.Null()},
		{val.Int(512), val.Int(512), val.True()},
		{val.Int(512), val.Text("2187 рублей"), val.True()},
	}.exec(t, Lte)
}

func Test_Gt(t *testing.T) {
	rows{
		{val.Int(512), val.Int(2187), val.False()},
		{val.Int(-512), val.Int(2187), val.False()},
		{val.Int(512), val.Text("2187"), val.False()},
		{val.Int(-512), val.Text("2187"), val.False()},
		{val.Int(512), val.Real(512.1), val.False()},
		{val.False(), val.True(), val.False()},
		{val.Text("hello"), val.Text("world"), val.False()},
		{val.Text("512hello"), val.Text("512world"), val.False()},
		{val.Int(2187), val.Int(512), val.True()},
		{val.Int(2187), val.Int(-512), val.True()},
		{val.Text("2187"), val.Int(512), val.True()},
		{val.Text("2187"), val.Int(-512), val.True()},
		{val.Real(512.1), val.Int(512), val.True()},
		{val.True(), val.False(), val.True()},
		{val.Text("world"), val.Text("hello"), val.True()},
		{val.Text("512world"), val.Text("512hello"), val.True()},
		{val.Int(512), val.Null(), val.Null()},
		{val.Null(), val.Int(512), val.Null()},
		{val.Int(512), val.Int(512), val.False()},
		{val.Int(2187), val.Text("512 рублей"), val.True()},
	}.exec(t, Gt)
}

func Test_Gte(t *testing.T) {
	rows{
		{val.Int(512), val.Int(2187), val.False()},
		{val.Int(-512), val.Int(2187), val.False()},
		{val.Int(512), val.Text("2187"), val.False()},
		{val.Int(-512), val.Text("2187"), val.False()},
		{val.Int(512), val.Real(512.1), val.False()},
		{val.False(), val.True(), val.False()},
		{val.Text("hello"), val.Text("world"), val.False()},
		{val.Text("512hello"), val.Text("512world"), val.False()},
		{val.Int(2187), val.Int(512), val.True()},
		{val.Int(2187), val.Int(-512), val.True()},
		{val.Text("2187"), val.Int(512), val.True()},
		{val.Text("2187"), val.Int(-512), val.True()},
		{val.Real(512.1), val.Int(512), val.True()},
		{val.True(), val.False(), val.True()},
		{val.Text("world"), val.Text("hello"), val.True()},
		{val.Text("512world"), val.Text("512hello"), val.True()},
		{val.Int(512), val.Null(), val.Null()},
		{val.Null(), val.Int(512), val.Null()},
		{val.Int(512), val.Int(512), val.True()},
		{val.Int(2187), val.Text("512 рублей"), val.True()},
	}.exec(t, Gte)
}

func Test_Concat(t *testing.T) {
	rows{
		{val.Int(512), val.Int(2187), val.Text("5122187")},
		{val.Real(5.12), val.Int(2187), val.Text("5.122187")},
		{val.Real(5.12), val.Real(2.187), val.Text("5.122.187")},
		{val.Real(5.12), val.Text("рублей"), val.Text("5.12рублей")},
		{val.True(), val.False(), val.Text("truefalse")},
		{val.Null(), val.Text("рублей"), val.Null()},
		{val.Real(.12), val.Text("рублей"), val.Text("0.12рублей")},
		{val.Text("hello "), val.Text("world"), val.Text("hello world")},
	}.exec(t, Concat)
}
