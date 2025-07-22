package op

import (
	"github.com/suprunchuksergey/dpl/value"
	"testing"
)

func Test_Add(t *testing.T) {
	tests := []struct {
		a,
		b,
		expected value.Value
	}{
		//целое число
		{value.Int(8), value.Int(8), value.Int(16)},
		{value.Int(8), value.Text("8"), value.Int(16)},
		{value.Int(8), value.Real(8.8), value.Real(16.8)},
		{value.Int(8), value.Text("8.8"), value.Real(16.8)},
		{value.Int(8), value.Text("text"), value.Int(8)},
		{value.Int(8), value.True(), value.Int(9)},
		{value.Int(8), value.False(), value.Int(8)},

		//число с плавающей точкой
		{value.Real(8.8), value.Int(8), value.Real(16.8)},
		{value.Real(8.8), value.Text("8"), value.Real(16.8)},
		{value.Real(8.8), value.Real(8.8), value.Real(17.6)},
		{value.Real(8.8), value.Text("8.8"), value.Real(17.6)},
		{value.Real(8.8), value.Text("text"), value.Real(8.8)},
		{value.Real(8.8), value.True(), value.Real(9.8)},
		{value.Real(8.8), value.False(), value.Real(8.8)},

		//строка как строка
		{value.Text("text"), value.Int(8), value.Int(8)},
		{value.Text("text"), value.Text("8"), value.Int(8)},
		{value.Text("text"), value.Real(8.8), value.Real(8.8)},
		{value.Text("text"), value.Text("8.8"), value.Real(8.8)},
		{value.Text("text"), value.Text("text"), value.Int(0)},
		{value.Text("text"), value.True(), value.Int(1)},
		{value.Text("text"), value.False(), value.Int(0)},

		//строка как целое число
		{value.Text("8"), value.Int(8), value.Int(16)},
		{value.Text("8"), value.Text("8"), value.Int(16)},
		{value.Text("8"), value.Real(8.8), value.Real(16.8)},
		{value.Text("8"), value.Text("8.8"), value.Real(16.8)},
		{value.Text("8"), value.Text("text"), value.Int(8)},
		{value.Text("8"), value.True(), value.Int(9)},
		{value.Text("8"), value.False(), value.Int(8)},

		//строка как число с плавающей точкой
		{value.Text("8.8"), value.Int(8), value.Real(16.8)},
		{value.Text("8.8"), value.Text("8"), value.Real(16.8)},
		{value.Text("8.8"), value.Real(8.8), value.Real(17.6)},
		{value.Text("8.8"), value.Text("8.8"), value.Real(17.6)},
		{value.Text("8.8"), value.Text("text"), value.Real(8.8)},
		{value.Text("8.8"), value.True(), value.Real(9.8)},
		{value.Text("8.8"), value.False(), value.Real(8.8)},

		{value.Null(), value.Int(8), value.Null()},
		{value.Null(), value.Text("8"), value.Null()},
		{value.Null(), value.Real(8.8), value.Null()},
		{value.Null(), value.Text("8.8"), value.Null()},
		{value.Null(), value.Text("text"), value.Null()},
		{value.Null(), value.Null(), value.Null()},

		{value.Int(8), value.Null(), value.Null()},
		{value.Text("8"), value.Null(), value.Null()},
		{value.Real(8.8), value.Null(), value.Null()},
		{value.Text("8.8"), value.Null(), value.Null()},
		{value.Text("text"), value.Null(), value.Null()},
		{value.Null(), value.Null(), value.Null()},

		{value.Int(64), value.Text("-8"), value.Int(56)},

		{value.True(), value.True(), value.Int(2)},
		{value.True(), value.False(), value.Int(1)},
		{value.False(), value.False(), value.Int(0)},
	}

	for i, test := range tests {
		got := Add(test.a, test.b)
		if got != test.expected {
			t.Errorf("%d: ожидалось %s, получено %s",
				i, test.expected, got)
		}
	}
}

func Test_Sub(t *testing.T) {
	tests := []struct {
		a,
		b,
		expected value.Value
	}{
		//целое число
		{value.Int(64), value.Int(8), value.Int(56)},
		{value.Int(64), value.Text("8"), value.Int(56)},
		{value.Int(64), value.Real(8.8), value.Real(55.2)},
		{value.Int(64), value.Text("8.8"), value.Real(55.2)},
		{value.Int(64), value.Text("text"), value.Int(64)},
		{value.Int(64), value.True(), value.Int(63)},
		{value.Int(64), value.False(), value.Int(64)},

		//число с плавающей точкой
		{value.Real(64.8), value.Int(8), value.Real(56.8)},
		{value.Real(64.8), value.Text("8"), value.Real(56.8)},
		{value.Real(64.8), value.Real(8.8), value.Real(56)},
		{value.Real(64.8), value.Text("8.8"), value.Real(56)},
		{value.Real(64.8), value.Text("text"), value.Real(64.8)},
		{value.Real(64.8), value.True(), value.Real(63.8)},
		{value.Real(64.8), value.False(), value.Real(64.8)},

		//строка как строка
		{value.Text("text"), value.Int(8), value.Int(-8)},
		{value.Text("text"), value.Text("8"), value.Int(-8)},
		{value.Text("text"), value.Real(8.8), value.Real(-8.8)},
		{value.Text("text"), value.Text("8.8"), value.Real(-8.8)},
		{value.Text("text"), value.Text("text"), value.Int(0)},
		{value.Text("text"), value.True(), value.Int(-1)},
		{value.Text("text"), value.False(), value.Int(0)},

		//строка как целое число
		{value.Text("64"), value.Int(8), value.Int(56)},
		{value.Text("64"), value.Text("8"), value.Int(56)},
		{value.Text("64"), value.Real(8.8), value.Real(55.2)},
		{value.Text("64"), value.Text("8.8"), value.Real(55.2)},
		{value.Text("64"), value.Text("text"), value.Int(64)},
		{value.Text("64"), value.True(), value.Int(63)},
		{value.Text("64"), value.False(), value.Int(64)},

		//строка как число с плавающей точкой
		{value.Text("64.8"), value.Int(8), value.Real(56.8)},
		{value.Text("64.8"), value.Text("8"), value.Real(56.8)},
		{value.Text("64.8"), value.Real(8.8), value.Real(56)},
		{value.Text("64.8"), value.Text("8.8"), value.Real(56)},
		{value.Text("64.8"), value.Text("text"), value.Real(64.8)},
		{value.Text("64.8"), value.True(), value.Real(63.8)},
		{value.Text("64.8"), value.False(), value.Real(64.8)},

		{value.Null(), value.Int(8), value.Null()},
		{value.Null(), value.Text("8"), value.Null()},
		{value.Null(), value.Real(8.8), value.Null()},
		{value.Null(), value.Text("8.8"), value.Null()},
		{value.Null(), value.Text("text"), value.Null()},
		{value.Null(), value.Null(), value.Null()},

		{value.Int(8), value.Null(), value.Null()},
		{value.Text("8"), value.Null(), value.Null()},
		{value.Real(8.8), value.Null(), value.Null()},
		{value.Text("8.8"), value.Null(), value.Null()},
		{value.Text("text"), value.Null(), value.Null()},
		{value.Null(), value.Null(), value.Null()},

		{value.Int(64), value.Text("-8"), value.Int(72)},

		{value.True(), value.True(), value.Int(0)},
		{value.True(), value.False(), value.Int(1)},
		{value.False(), value.False(), value.Int(0)},
	}

	for i, test := range tests {
		got := Sub(test.a, test.b)
		if got != test.expected {
			t.Errorf("%d: ожидалось %s, получено %s",
				i, test.expected, got)
		}
	}
}

func Test_Mul(t *testing.T) {
	tests := []struct {
		a,
		b,
		expected value.Value
	}{
		//целое число
		{value.Int(64), value.Int(8), value.Int(512)},
		{value.Int(64), value.Text("8"), value.Int(512)},
		{value.Int(64), value.Real(8.8), value.Real(563.2)},
		{value.Int(64), value.Text("8.8"), value.Real(563.2)},
		{value.Int(64), value.Text("text"), value.Int(0)},

		//число с плавающей точкой
		{value.Real(64.8), value.Int(8), value.Real(518.4)},
		{value.Real(64.8), value.Text("8"), value.Real(518.4)},
		{value.Real(64.8), value.Real(8.8), value.Real(570.24)},
		{value.Real(64.8), value.Text("8.8"), value.Real(570.24)},
		{value.Real(64.8), value.Text("text"), value.Real(0)},

		//строка как строка
		{value.Text("text"), value.Int(8), value.Int(0)},
		{value.Text("text"), value.Text("8"), value.Int(0)},
		{value.Text("text"), value.Real(8.8), value.Real(0)},
		{value.Text("text"), value.Text("8.8"), value.Real(0)},
		{value.Text("text"), value.Text("text"), value.Int(0)},

		//строка как целое число
		{value.Text("64"), value.Int(8), value.Int(512)},
		{value.Text("64"), value.Text("8"), value.Int(512)},
		{value.Text("64"), value.Real(8.8), value.Real(563.2)},
		{value.Text("64"), value.Text("8.8"), value.Real(563.2)},
		{value.Text("64"), value.Text("text"), value.Int(0)},

		//строка как число с плавающей точкой
		{value.Text("64.8"), value.Int(8), value.Real(518.4)},
		{value.Text("64.8"), value.Text("8"), value.Real(518.4)},
		{value.Text("64.8"), value.Real(8.8), value.Real(570.24)},
		{value.Text("64.8"), value.Text("8.8"), value.Real(570.24)},
		{value.Text("64.8"), value.Text("text"), value.Real(0)},

		{value.Null(), value.Int(8), value.Null()},
		{value.Null(), value.Text("8"), value.Null()},
		{value.Null(), value.Real(8.8), value.Null()},
		{value.Null(), value.Text("8.8"), value.Null()},
		{value.Null(), value.Text("text"), value.Null()},
		{value.Null(), value.Null(), value.Null()},

		{value.Int(8), value.Null(), value.Null()},
		{value.Text("8"), value.Null(), value.Null()},
		{value.Real(8.8), value.Null(), value.Null()},
		{value.Text("8.8"), value.Null(), value.Null()},
		{value.Text("text"), value.Null(), value.Null()},
		{value.Null(), value.Null(), value.Null()},

		{value.Int(64), value.Text("-8"), value.Int(-512)},
	}

	for i, test := range tests {
		got := Mul(test.a, test.b)
		if got != test.expected {
			t.Errorf("%d: ожидалось %s, получено %s",
				i, test.expected, got)
		}
	}
}

func Test_Div(t *testing.T) {
	tests := []struct {
		a,
		b,
		expected value.Value
	}{
		//целое число
		{value.Int(64), value.Int(8), value.Int(8)},
		{value.Int(64), value.Text("8"), value.Int(8)},
		{value.Int(64), value.Real(8.8), value.Real(7.2727272727272725)},
		{value.Int(64), value.Text("8.8"), value.Real(7.2727272727272725)},
		{value.Int(64), value.Text("text"), value.Null()},

		//число с плавающей точкой
		{value.Real(64.8), value.Int(8), value.Real(8.1)},
		{value.Real(64.8), value.Text("8"), value.Real(8.1)},
		{value.Real(64.8), value.Real(8.8), value.Real(7.363636363636362)},
		{value.Real(64.8), value.Text("8.8"), value.Real(7.363636363636362)},
		{value.Real(64.8), value.Text("text"), value.Null()},

		//строка как строка
		{value.Text("text"), value.Int(8), value.Int(0)},
		{value.Text("text"), value.Text("8"), value.Int(0)},
		{value.Text("text"), value.Real(8.8), value.Real(0)},
		{value.Text("text"), value.Text("8.8"), value.Real(0)},
		{value.Text("text"), value.Text("text"), value.Null()},

		//строка как целое число
		{value.Text("64"), value.Int(8), value.Int(8)},
		{value.Text("64"), value.Text("8"), value.Int(8)},
		{value.Text("64"), value.Real(8.8), value.Real(7.2727272727272725)},
		{value.Text("64"), value.Text("8.8"), value.Real(7.2727272727272725)},
		{value.Text("64"), value.Text("text"), value.Null()},

		//строка как число с плавающей точкой
		{value.Text("64.8"), value.Int(8), value.Real(8.1)},
		{value.Text("64.8"), value.Text("8"), value.Real(8.1)},
		{value.Text("64.8"), value.Real(8.8), value.Real(7.363636363636362)},
		{value.Text("64.8"), value.Text("8.8"), value.Real(7.363636363636362)},
		{value.Text("64.8"), value.Text("text"), value.Null()},

		{value.Int(0), value.Int(8), value.Int(0)},
		{value.Text("0"), value.Int(8), value.Int(0)},
		{value.Real(0), value.Int(8), value.Real(0)},
		{value.Text("0.0"), value.Int(8), value.Real(0)},
		{value.Text("text"), value.Int(8), value.Int(0)},

		{value.Null(), value.Int(8), value.Null()},
		{value.Null(), value.Text("8"), value.Null()},
		{value.Null(), value.Real(8.8), value.Null()},
		{value.Null(), value.Text("8.8"), value.Null()},
		{value.Null(), value.Text("text"), value.Null()},
		{value.Null(), value.Null(), value.Null()},

		{value.Int(8), value.Null(), value.Null()},
		{value.Text("8"), value.Null(), value.Null()},
		{value.Real(8.8), value.Null(), value.Null()},
		{value.Text("8.8"), value.Null(), value.Null()},
		{value.Text("text"), value.Null(), value.Null()},
		{value.Null(), value.Null(), value.Null()},

		{value.Int(8), value.Int(0), value.Null()},
		{value.Int(8), value.Text("0"), value.Null()},
		{value.Int(8), value.Real(0), value.Null()},
		{value.Int(8), value.Text("0"), value.Null()},
		{value.Int(8), value.Text("text"), value.Null()},

		{value.Int(8), value.Text(".0"), value.Null()},
		{value.Int(8), value.Text("0."), value.Null()},
		{value.Int(8), value.Text("-.0"), value.Null()},
		{value.Int(8), value.Text("-0."), value.Null()},
		{value.Int(8), value.Text("+.0"), value.Null()},
		{value.Int(8), value.Text("+0."), value.Null()},

		{value.Int(64), value.Text("-8"), value.Int(-8)},

		{value.Int(10), value.Int(8), value.Int(1)},
		{value.Int(10), value.Int(4), value.Int(2)},
		{value.Int(10), value.Real(8), value.Real(1.25)},
		{value.Int(10), value.Real(4), value.Real(2.5)},
		{value.Real(10), value.Int(8), value.Real(1.25)},
		{value.Real(10), value.Int(4), value.Real(2.5)},

		{value.Int(10), value.False(), value.Null()},
		{value.Int(10), value.True(), value.Int(10)},
	}

	for i, test := range tests {
		got := Div(test.a, test.b)
		if got != test.expected {
			t.Errorf("%d: ожидалось %s, получено %s",
				i, test.expected, got)
		}
	}
}

func Test_Rem(t *testing.T) {
	tests := []struct {
		a,
		b,
		expected value.Value
	}{
		//целое число
		{value.Int(69), value.Int(8), value.Int(5)},
		{value.Int(69), value.Text("8"), value.Int(5)},
		{value.Int(69), value.Real(8.8), value.Real(7.399999999999995)},
		{value.Int(69), value.Text("8.8"), value.Real(7.399999999999995)},
		{value.Int(69), value.Text("text"), value.Null()},

		//число с плавающей точкой
		{value.Real(69.8), value.Int(8), value.Real(5.799999999999997)},
		{value.Real(69.8), value.Text("8"), value.Real(5.799999999999997)},
		{value.Real(69.8), value.Real(8.8), value.Real(8.199999999999992)},
		{value.Real(69.8), value.Text("8.8"), value.Real(8.199999999999992)},
		{value.Real(69.8), value.Text("text"), value.Null()},

		//строка как строка
		{value.Text("text"), value.Int(8), value.Int(0)},
		{value.Text("text"), value.Text("8"), value.Int(0)},
		{value.Text("text"), value.Real(8.8), value.Real(0)},
		{value.Text("text"), value.Text("8.8"), value.Real(0)},
		{value.Text("text"), value.Text("text"), value.Null()},

		//строка как целое число
		{value.Text("69"), value.Int(8), value.Int(5)},
		{value.Text("69"), value.Text("8"), value.Int(5)},
		{value.Text("69"), value.Real(8.8), value.Real(7.399999999999995)},
		{value.Text("69"), value.Text("8.8"), value.Real(7.399999999999995)},
		{value.Text("69"), value.Text("text"), value.Null()},

		//строка как число с плавающей точкой
		{value.Text("69.8"), value.Int(8), value.Real(5.799999999999997)},
		{value.Text("69.8"), value.Text("8"), value.Real(5.799999999999997)},
		{value.Text("69.8"), value.Real(8.8), value.Real(8.199999999999992)},
		{value.Text("69.8"), value.Text("8.8"), value.Real(8.199999999999992)},
		{value.Text("69.8"), value.Text("text"), value.Null()},

		{value.Int(0), value.Int(8), value.Int(0)},
		{value.Text("0"), value.Int(69), value.Int(0)},
		{value.Real(0), value.Int(6), value.Real(0)},
		{value.Text("0.0"), value.Int(9), value.Real(0)},
		{value.Text("text"), value.Int(1), value.Int(0)},

		{value.Null(), value.Int(8), value.Null()},
		{value.Null(), value.Text("8"), value.Null()},
		{value.Null(), value.Real(8.8), value.Null()},
		{value.Null(), value.Text("8.8"), value.Null()},
		{value.Null(), value.Text("text"), value.Null()},
		{value.Null(), value.Null(), value.Null()},

		{value.Int(8), value.Null(), value.Null()},
		{value.Text("8"), value.Null(), value.Null()},
		{value.Real(8.8), value.Null(), value.Null()},
		{value.Text("8.8"), value.Null(), value.Null()},
		{value.Text("text"), value.Null(), value.Null()},
		{value.Null(), value.Null(), value.Null()},

		{value.Int(8), value.Int(0), value.Null()},
		{value.Int(69), value.Text("0"), value.Null()},
		{value.Int(6), value.Real(0), value.Null()},
		{value.Int(9), value.Text("0"), value.Null()},
		{value.Int(1), value.Text("text"), value.Null()},

		{value.Int(8), value.Text(".0"), value.Null()},
		{value.Int(69), value.Text("0."), value.Null()},
		{value.Int(6), value.Text("-.0"), value.Null()},
		{value.Int(9), value.Text("-0."), value.Null()},
		{value.Int(1), value.Text("+.0"), value.Null()},
		{value.Int(19), value.Text("+0."), value.Null()},

		{value.Int(10), value.Int(8), value.Int(2)},
		{value.Int(9), value.Int(6), value.Int(3)},
		{value.Int(69), value.Int(6), value.Int(3)},
		{value.Real(69.9), value.Int(8), value.Real(5.900000000000006)},
		{value.Int(69), value.Real(8.8), value.Real(7.399999999999995)},
		{value.Int(8), value.Int(8), value.Int(0)},

		{value.Text("69.9"), value.Text("8"), value.Real(5.900000000000006)},
		{value.Text("69.9"), value.Int(8), value.Real(5.900000000000006)},
		{value.Text("69.9"), value.Real(8), value.Real(5.900000000000006)},
		{value.Text("9"), value.Text("6"), value.Int(3)},
		{value.Int(512), value.Int(81), value.Int(26)},
		{value.Real(512.512), value.Int(81), value.Real(26.511999999999944)},
		{value.Int(512), value.Real(81.81), value.Real(21.139999999999986)},

		{value.Int(10), value.False(), value.Null()},
		{value.Int(10), value.True(), value.Int(0)},
	}

	for i, test := range tests {
		got := Rem(test.a, test.b)
		if got != test.expected {
			t.Errorf("%d: ожидалось %s, получено %s",
				i, test.expected, got)
		}
	}
}
