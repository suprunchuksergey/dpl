package value

import "testing"

func Test_integer_Int(t *testing.T) {
	tests := []struct{ data int64 }{{64}, {-64}}

	for i, test := range tests {
		v := newInteger(test.data)
		got := v.Int()
		if got != v.v {
			t.Errorf("%d: ожидалось %d, получено %d",
				i, v.v, got)
		}
	}
}

func Test_integer_Real(t *testing.T) {
	tests := []struct{ data int64 }{{64}, {-64}}

	for i, test := range tests {
		v := newInteger(test.data)
		got := v.Real()
		if got != float64(v.v) {
			t.Errorf("%d: ожидалось %f, получено %f",
				i, float64(v.v), got)
		}
	}
}

func Test_integer_Text(t *testing.T) {
	tests := []struct {
		data     int64
		expected string
	}{
		{data: 64, expected: "64"},
		{data: -128, expected: "-128"},
	}

	for i, test := range tests {
		got := newInteger(test.data).Text()
		if got != test.expected {
			t.Errorf("%d: ожидалось %s, получено %s",
				i, test.expected, got)
		}
	}
}
