package value

import (
	"testing"
)

func Test_text_Int(t *testing.T) {
	tests := []struct {
		data     string
		expected int64
	}{
		{data: "", expected: 0},
		{data: "			", expected: 0},
		{data: "-", expected: 0},
		{data: "+", expected: 0},
		{data: "	-", expected: 0},
		{data: "	+", expected: 0},
		{data: "рубля", expected: 0},
		{data: "0", expected: 0},
		{data: "64 рубля", expected: 64},
		{data: "+64 рубля", expected: 64},
		{data: "-64 рубля", expected: -64},
		{data: "	64рубля", expected: 64},
		{data: "	+64рубля", expected: 64},
		{data: "	-64рубля", expected: -64},
		{data: "64.64 рубля", expected: 64},
		{data: "+64.64 рубля", expected: 64},
		{data: "-64.64 рубля", expected: -64},
		{data: "--64.64 рубля", expected: 0},
		{data: "- 64.64 рубля", expected: 0},
		{data: "++64.64 рубля", expected: 0},
		{data: "+ 64.64 рубля", expected: 0},
		{data: "64 666", expected: 64},
		{data: "			64		рубля	", expected: 64},
	}

	for i, test := range tests {
		got := newText(test.data).Int()
		if got != test.expected {
			t.Errorf("%d: ожидалось %d, получено %d",
				i, test.expected, got)
		}
	}
}

func Test_text_Real(t *testing.T) {
	tests := []struct {
		data     string
		expected float64
	}{
		{data: "", expected: 0},
		{data: "			", expected: 0},
		{data: "-", expected: 0},
		{data: "+", expected: 0},
		{data: "	-", expected: 0},
		{data: "	+", expected: 0},
		{data: "рубля", expected: 0},
		{data: "0", expected: 0},
		{data: "0.0", expected: 0},
		{data: "64 рубля", expected: 64},
		{data: "+64 рубля", expected: 64},
		{data: "-64 рубля", expected: -64},
		{data: "	64рубля", expected: 64},
		{data: "	+64рубля", expected: 64},
		{data: "	-64рубля", expected: -64},
		{data: "64.64 рубля", expected: 64.64},
		{data: "+64.64 рубля", expected: 64.64},
		{data: "-64.64 рубля", expected: -64.64},
		{data: "--64.64 рубля", expected: 0},
		{data: "- 64.64 рубля", expected: 0},
		{data: "++64.64 рубля", expected: 0},
		{data: "+ 64.64 рубля", expected: 0},
		{data: "64 666", expected: 64},
		{data: "			64		рубля	", expected: 64},
		{data: ".64", expected: .64},
		{data: "		.64", expected: .64},
		{data: "..64", expected: 0},
		{data: "64.64.64", expected: 64.64},
		{data: ".64.64", expected: .64},
		{data: "64.64.64рубля", expected: 64.64},
		{data: ".64.64рубля", expected: .64},
		{data: "-.64.64", expected: -.64},
		{data: "-64.64.64", expected: -64.64},
		{data: "-64..64", expected: -64},
		{data: "-64-64", expected: -64},
		{data: "+64.64", expected: 64.64},
		{data: "+0.64", expected: .64},
		{data: "-0.64", expected: -.64},
	}

	for i, test := range tests {
		got := newText(test.data).Real()
		if got != test.expected {
			t.Errorf("%d: ожидалось %f, получено %f",
				i, test.expected, got)
		}
	}
}

func Test_text_Text(t *testing.T) {
	tests := []struct{ data string }{
		{"data"},
		{"			data"},
		{"				-64..64"},
	}

	for i, test := range tests {
		txt := newText(test.data)
		got := txt.Text()
		if got != txt.v {
			t.Errorf("%d: ожидалось %s, получено %s",
				i, txt.v, got)
		}
	}
}

func Test_text_IsInt(t *testing.T) {
	tests := []struct {
		data     string
		expected bool
	}{
		{"64", true},
		{"	 64", true},
		{"-64", true},
		{"+64", true},
		{"	 -64", true},
		{"	 +64", true},
		{"64 рубля", true},
		{"	 64 рубля", true},
		{"", false},
		{"		 ", false},
		{"руб", false},
		{"-", false},
		{"	 -", false},
		{"	 +", false},
		{"+руб", false},
		{"-руб", false},
		{" +руб", false},
		{" -руб", false},
		{".64", false},
	}

	for i, test := range tests {
		if newText(test.data).IsInt() != test.expected {
			t.Errorf("%d: %q: ожидалось %t, получено %t",
				i, test.data, test.expected, !test.expected)
		}
	}
}

func Test_text_IsReal(t *testing.T) {
	tests := []struct {
		data     string
		expected bool
	}{
		{"64", false},
		{"	 64", false},
		{"-64", false},
		{"+64", false},
		{"	 -64", false},
		{"	 +64", false},
		{"64 рубля", false},
		{"	 64 рубля", false},
		{"", false},
		{"		 ", false},
		{"руб", false},
		{"-", false},
		{"	 -", false},
		{"	 +", false},
		{"+руб", false},
		{"-руб", false},
		{" +руб", false},
		{" -руб", false},
		{"64.64", true},
		{"64.", true},
		{".64", true},
		{"	 64.64", true},
		{"	 64.", true},
		{"	 .64", true},
		{"	 -64.64", true},
		{"	 -64.", true},
		{"	 -.64", true},
		{"	 +64.64", true},
		{"	 +64.", true},
		{"	 +.64", true},
		{"	 +64.64руб", true},
		{"	 +64.руб", true},
		{"	 +.64руб", true},
		{".", false},
		{"+.", false},
		{"-.", false},
		{"-.6", true},
		{"+.6", true},
		{"+.0", true},
		{"0.0", true},
		{"0.", true},
	}

	for i, test := range tests {
		if newText(test.data).IsReal() != test.expected {
			t.Errorf("%d: %q: ожидалось %t, получено %t",
				i, test.data, test.expected, !test.expected)
		}
	}
}

func Test_text_Bool(t *testing.T) {
	tests := []struct {
		data     string
		expected bool
	}{
		{"64", true},
		{"data", true},
		{"", false},
	}

	for i, test := range tests {
		if newText(test.data).Bool() != test.expected {
			t.Errorf("%d: %q: ожидалось %t, получено %t",
				i, test.data, test.expected, !test.expected)
		}
	}
}
