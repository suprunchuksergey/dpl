package text

import "testing"

func Test_Int(t *testing.T) {
	tests := []struct {
		data     string
		expected int64
	}{
		{"134217728", 134217728},
		{"+262144", 262144},
		{"-68719476736", -68719476736},
		{"	 16777216", 16777216},
		{"	 -65536", -65536},
		{"	 +65.536", 65},
		{"	 -65.data", -65},
		{"	 -65data", -65},
		{"	 ", 0},
		{"", 0},
		{"-", 0},
		{"+", 0},
		{".536", 0},
		{"--4294967296", 0},
		{"+-4294967296", 0},
		{"++4294967296", 0},
		{"data", 0},
		{"-data", 0},
	}

	for i, test := range tests {
		text := New(test.data)
		if text.Int() != test.expected {
			t.Errorf("%d: ожидалось %d, получено %d",
				i, test.expected, text.Int())
		}
	}
}

func Test_Real(t *testing.T) {
	tests := []struct {
		data     string
		expected float64
	}{
		{"134217728", 134217728},
		{"+262144", 262144},
		{"-68719476736", -68719476736},
		{"-68719.", -68719},
		{"	 16777216", 16777216},
		{"	 -65536", -65536},
		{"	 -65data", -65},
		{"	 ", 0},
		{"", 0},
		{"-", 0},
		{"+", 0},
		{"--4294967296", 0},
		{"+-4294967296", 0},
		{"++4294967296", 0},
		{"26.2144", 26.2144},
		{"-1342177.28", -1342177.28},
		{"+1677.7216", 1677.7216},
		{"-.7216", -.7216},
		{"	-.1677", -.1677},
		{"-.0", 0},
		{"672.", 672},
		{"..672", 0},
		{"672.2.672", 672.2},
		{"-2.62.267", -2.62},
		{"-2.data.267", -2},
		{"-22..67", -22},
		{"	+42949672.96.data", 42949672.96},
		{"	 -429.4967296data", -429.4967296},
		{"	 -. ", 0},
		{"	 +.data", 0},
		{"	 -. .429", 0},
		{"data", 0},
		{"-data", 0},
	}

	for i, test := range tests {
		text := New(test.data)
		if text.Real() != test.expected {
			t.Errorf("%d: ожидалось %f, получено %f",
				i, test.expected, text.Real())
		}
	}
}

func Test_CanInt(t *testing.T) {
	tests := []struct {
		data     string
		expected bool
	}{
		{"134217728", true},
		{"+262144", true},
		{"-68719476736", true},
		{"	 16777216", true},
		{"	 -65536", true},
		{"	 +65.536", true},
		{"	 -65.data", true},
		{"	 -65data", true},
		{"	 ", false},
		{"", false},
		{"-", false},
		{"+", false},
		{"	-", false},
		{"	+", false},
		{".536", false},
		{"--4294967296", false},
		{"+-4294967296", false},
		{"++4294967296", false},
		{"data", false},
		{"-data", false},
		{"0", true},
		{"-0", true},
	}

	for i, test := range tests {
		text := New(test.data)
		if text.CanInt() != test.expected {
			t.Errorf("%d: ожидалось %t, получено %t",
				i, test.expected, text.CanInt())
		}
	}
}

func Test_CanReal(t *testing.T) {
	tests := []struct {
		data     string
		expected bool
	}{
		{"134.217728", true},
		{"+2621.44", true},
		{"-68719476736.", true},
		{"-.68719", true},
		{"	 .16777216", true},
		{"	 -655.36", true},
		{"	 -65.data", true},
		{"	 ", false},
		{"", false},
		{"-", false},
		{"+", false},
		{"	 -", false},
		{"	 +", false},
		{"	 -.", false},
		{"	 +.", false},
		{".", false},
		{"--4294.967296", false},
		{"+-429496729.6", false},
		{"++42.94967296", false},
		{"26.2144", true},
		{"-1342177.28", true},
		{"+1677.7216", true},
		{"-.7216", true},
		{"	-.1677", true},
		{"-.0", true},
		{"672.", true},
		{"..672", false},
		{"672.2.672", true},
		{"-2.62.267", true},
		{"-2.data.267", true},
		{"-22..67", true},
		{"	+42949672.96.data", true},
		{"	 -429.4967296data", true},
		{"	 -. ", false},
		{"	 +.data", false},
		{"	 -. .429", false},
		{"data", false},
		{"-data", false},
		{"429", false},
	}

	for i, test := range tests {
		text := New(test.data)
		if text.CanReal() != test.expected {
			t.Errorf("%d: ожидалось %t, получено %t",
				i, test.expected, text.CanReal())
		}
	}
}

func Test_Bool(t *testing.T) {
	empty := New("")
	if empty.Bool() == true {
		t.Error("ожидалось false, получено true")
	}
	text := New("text")
	if text.Bool() == false {
		t.Error("ожидалось true, получено false")
	}
}
