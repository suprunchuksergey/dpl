package val

import (
	"reflect"
	"testing"
	"unicode"
)

func Test_ToInt(t *testing.T) {
	tests := []struct {
		val      Val
		expected int64
	}{
		{Int(2187), 2187},
		{Int(-2187), -2187},
		{Real(2187), 2187},
		{Real(-2187), -2187},
		{Real(21.87), 21},
		{Real(-21.87), -21},
		{True(), 1},
		{False(), 0},
		{Text("text"), 0},
		{Text("-text"), 0},
		{Text("+text"), 0},
		{Text("	-text"), 0},
		{Text("	+text"), 0},
		{Text("0"), 0},
		{Text("-0"), 0},
		{Text("+0"), 0},
		{Text("2187"), 2187},
		{Text("-2187"), -2187},
		{Text("+2187"), 2187},
		{Text("2187text"), 2187},
		{Text("-2187text"), -2187},
		{Text("+2187text"), 2187},
		{Text("21.87"), 21},
		{Text("-21.87"), -21},
		{Text("+21.87"), 21},
		{Text("	21.87"), 21},
		{Text("	-21.87"), -21},
		{Text("	+21.87"), 21},
		{Text("	21.87text"), 21},
		{Text("	-21.87text"), -21},
		{Text("	+21.87text"), 21},
		{Null(), 0},
	}
	for i, test := range tests {
		n := test.val.ToInt()
		if n != test.expected {
			t.Errorf("%d: ожидалось %d, получено %d", i, test.expected, n)
		}
	}
}

func Test_ToReal(t *testing.T) {
	tests := []struct {
		val      Val
		expected float64
	}{
		{Int(2187), 2187},
		{Int(-2187), -2187},
		{Real(2187), 2187},
		{Real(-2187), -2187},
		{Real(21.87), 21.87},
		{Real(-21.87), -21.87},
		{True(), 1},
		{False(), 0},
		{Text("text"), 0},
		{Text("-text"), 0},
		{Text("+text"), 0},
		{Text("	-text"), 0},
		{Text("	+text"), 0},
		{Text("0"), 0},
		{Text("-0"), 0},
		{Text("+0"), 0},
		{Text("2187"), 2187},
		{Text("-2187"), -2187},
		{Text("+2187"), 2187},
		{Text("2187text"), 2187},
		{Text("-2187text"), -2187},
		{Text("+2187text"), 2187},
		{Text("21.87"), 21.87},
		{Text("-21.87"), -21.87},
		{Text("+21.87"), 21.87},
		{Text("	21.87"), 21.87},
		{Text("	-21.87"), -21.87},
		{Text("	+21.87"), 21.87},
		{Text("	21.87text"), 21.87},
		{Text("	-21.87text"), -21.87},
		{Text("	+21.87text"), 21.87},
		{Text(".87"), .87},
		{Text("-.87"), -.87},
		{Text("+.87"), .87},
		{Null(), 0},
	}
	for i, test := range tests {
		n := test.val.ToReal()
		if n != test.expected {
			t.Errorf("%d: ожидалось %f, получено %f", i, test.expected, n)
		}
	}
}

func Test_ToText(t *testing.T) {
	tests := []struct {
		val      Val
		expected string
	}{
		{Int(2187), "2187"},
		{Int(-2187), "-2187"},
		{Int(+2187), "2187"},
		{Real(21.87), "21.87"},
		{Real(.87), "0.87"},
		{Real(-.87), "-0.87"},
		{Real(+.87), "0.87"},
		{True(), "true"},
		{False(), "false"},
		{Text("text"), "text"},
		{Text("2187text"), "2187text"},
		{Text("21.87text"), "21.87text"},
		{Null(), ""},
	}
	for i, test := range tests {
		n := test.val.ToText()
		if n != test.expected {
			t.Errorf("%d: ожидалось %s, получено %s", i, test.expected, n)
		}
	}
}

func Test_ToBool(t *testing.T) {
	tests := []struct {
		val      Val
		expected bool
	}{
		{Int(2187), true},
		{Int(-2187), true},
		{Int(0), false},
		{Int(-0), false},
		{Real(2187), true},
		{Real(-2187), true},
		{Real(0), false},
		{Real(-0), false},
		{Text("text"), true},
		{Text(""), false},
		{Text("2187text"), true},
		{Text("21.87text"), true},
		{True(), true},
		{False(), false},
		{Null(), false},
	}
	for i, test := range tests {
		n := test.val.ToBool()
		if n != test.expected {
			t.Errorf("%d: ожидалось %t, получено %t", i, test.expected, n)
		}
	}
}

func Test_ToArray(t *testing.T) {
	tests := []struct {
		val      Val
		expected []Val
	}{
		{Array([]Val{Int(2187)}), []Val{Int(2187)}},
		{Text("text"), []Val{
			Text("t"),
			Text("e"),
			Text("x"),
			Text("t"),
		}},
		{Text(""), []Val{}},
	}
	for i, test := range tests {
		n := test.val.ToArray()
		if !reflect.DeepEqual(n, test.expected) {
			t.Errorf("%d: ожидалось %v, получено %v", i, test.expected, n)
		}
	}
}

func Test_ToMap(t *testing.T) {
	tests := []struct {
		val      Val
		expected map[string]Val
	}{
		{Map(map[string]Val{"a": Int(2187)}), map[string]Val{"a": Int(2187)}},
	}
	for i, test := range tests {
		n := test.val.ToMap()
		if !reflect.DeepEqual(n, test.expected) {
			t.Errorf("%d: ожидалось %v, получено %v", i, test.expected, n)
		}
	}
}

func Test_IsInt(t *testing.T) {
	tests := []struct {
		val      Val
		expected bool
	}{
		{Int(2187), true},
		{Real(2187), false},
		{Text("text"), false},
		{Text(""), false},
		{Text("2187text"), true},
		{Text("21.87text"), true},
		{Text(".87text"), false},
		{True(), false},
		{False(), false},
		{Null(), false},
	}
	for i, test := range tests {
		n := test.val.IsInt()
		if n != test.expected {
			t.Errorf("%d: ожидалось %t, получено %t", i, test.expected, n)
		}
	}
}

func Test_IsReal(t *testing.T) {
	tests := []struct {
		val      Val
		expected bool
	}{
		{Real(2187), true},
		{Int(2187), false},
		{Text("text"), false},
		{Text(""), false},
		{Text("2187text"), false},
		{Text("21.87text"), true},
		{Text("-21.87text"), true},
		{Text("+21.87text"), true},
		{Text(".87text"), true},
		{Text("-.87text"), true},
		{Text("+.87text"), true},
		{True(), false},
		{False(), false},
		{Null(), false},
	}
	for i, test := range tests {
		n := test.val.IsReal()
		if n != test.expected {
			t.Errorf("%d: ожидалось %t, получено %t", i, test.expected, n)
		}
	}
}

func Test_IsText(t *testing.T) {
	tests := []struct {
		val      Val
		expected bool
	}{
		{Text("text"), true},
		{Text(""), true},
		{Text("2187"), true},
		{Text("-2187"), true},
		{Text("+2187"), true},
		{Text("21.87"), true},
		{Text("-21.87"), true},
		{Text("+21.87"), true},
		{Int(2187), false},
		{Real(2187), false},
		{True(), false},
		{False(), false},
		{Null(), false},
	}
	for i, test := range tests {
		n := test.val.IsText()
		if n != test.expected {
			t.Errorf("%d: ожидалось %t, получено %t", i, test.expected, n)
		}
	}
}

func Test_IsBool(t *testing.T) {
	tests := []struct {
		val      Val
		expected bool
	}{
		{True(), true},
		{False(), true},
		{Bool(true), true},
		{Bool(false), true},
		{Text("text"), false},
		{Int(2187), false},
		{Real(2187), false},
		{Null(), false},
	}
	for i, test := range tests {
		n := test.val.IsBool()
		if n != test.expected {
			t.Errorf("%d: ожидалось %t, получено %t", i, test.expected, n)
		}
	}
}

func Test_IsNull(t *testing.T) {
	tests := []struct {
		val      Val
		expected bool
	}{
		{Int(2187), false},
		{Real(21.87), false},
		{Text("text"), false},
		{True(), false},
		{False(), false},
		{Null(), true},
	}
	for i, test := range tests {
		n := test.val.IsNull()
		if n != test.expected {
			t.Errorf("%d: ожидалось %t, получено %t", i, test.expected, n)
		}
	}
}

func Test_skip(t *testing.T) {
	tests := []struct {
		val      string
		c        func(rune) bool
		expected int
	}{
		{"2187rune", unicode.IsDigit, 4},
		{"			 rune", unicode.IsSpace, 4},
		{"2187", unicode.IsDigit, 4},
		{"			 ", unicode.IsSpace, 4},
	}
	for i, test := range tests {
		n := skip([]rune(test.val), 0, test.c)
		if n != test.expected {
			t.Errorf("%d: ожидалось %d, получено %d", i, test.expected, n)
		}
	}
}

func Test_textIsInt(t *testing.T) {
	tests := []struct {
		val      string
		expected bool
	}{
		{"", false},
		{"	 ", false},
		{"	 -", false},
		{"	 +", false},
		{"string", false},
		{"	 string", false},
		{"	 -string", false},
		{"	 +string", false},
		{"2187", true},
		{"-2187", true},
		{"+2187", true},
		{"21.87", true},
		{"-21.87", true},
		{"+21.87", true},
		{"	21.87string", true},
		{"	-21.87string", true},
		{"	+21.87string", true},
	}
	for i, test := range tests {
		n := textIsInt([]rune(test.val))
		if n != test.expected {
			t.Errorf("%d: ожидалось %t, получено %t", i, test.expected, n)
		}
	}
}

func Test_textIsReal(t *testing.T) {
	tests := []struct {
		val      string
		expected bool
	}{
		{"", false},
		{"	 ", false},
		{"	 -", false},
		{"	 +", false},
		{"string", false},
		{"	 string", false},
		{"	 -string", false},
		{"	 +string", false},
		{"2187", false},
		{"-2187", false},
		{"+2187", false},
		{"21.87", true},
		{"-21.87", true},
		{"+21.87", true},
		{"	21.87string", true},
		{"	-21.87string", true},
		{"	+21.87string", true},
		{"	 .87string", true},
		{"	 -.87string", true},
		{"	 +.87string", true},
	}
	for i, test := range tests {
		n := textIsReal([]rune(test.val))
		if n != test.expected {
			t.Errorf("%d: ожидалось %t, получено %t", i, test.expected, n)
		}
	}
}

func Test_textToInt(t *testing.T) {
	tests := []struct {
		val      string
		expected int64
	}{
		{"", 0},
		{"	 ", 0},
		{"	 -", 0},
		{"	 +", 0},
		{"string", 0},
		{"	 string", 0},
		{"	 -string", 0},
		{"	 +string", 0},
		{"2187", 2187},
		{"-2187", -2187},
		{"+2187", 2187},
		{"21.87", 21},
		{"-21.87", -21},
		{"+21.87", 21},
		{"	21.87string", 21},
		{"	-21.87string", -21},
		{"	+21.87string", 21},
	}
	for i, test := range tests {
		n := textToInt([]rune(test.val))
		if n != test.expected {
			t.Errorf("%d: ожидалось %d, получено %d", i, test.expected, n)
		}
	}
}

func Test_textToReal(t *testing.T) {
	tests := []struct {
		val      string
		expected float64
	}{
		{"", 0},
		{"	 ", 0},
		{"	 -", 0},
		{"	 +", 0},
		{"string", 0},
		{"	 string", 0},
		{"	 -string", 0},
		{"	 +string", 0},
		{"2187", 2187},
		{"-2187", -2187},
		{"+2187", 2187},
		{"21.87", 21.87},
		{"-21.87", -21.87},
		{"+21.87", 21.87},
		{"	21.87string", 21.87},
		{"	-21.87string", -21.87},
		{"	+21.87string", 21.87},
		{".87", .87},
		{"-.87", -.87},
		{"+.87", .87},
		{"0.87", .87},
		{"-0.87", -.87},
		{"+0.87", .87},
	}
	for i, test := range tests {
		n := textToReal([]rune(test.val))
		if n != test.expected {
			t.Errorf("%d: ожидалось %f, получено %f", i, test.expected, n)
		}
	}
}
