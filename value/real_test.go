package value

import "testing"

func Test_real_Int(t *testing.T) {
	tests := []struct {
		data     float64
		expected int64
	}{
		{64, 64},
		{64.64, 64},
		{-64.64, -64},
		{.64, 0},
	}

	for i, test := range tests {
		got := newReal(test.data).Int()
		if got != test.expected {
			t.Errorf("%d: ожидалось %d, получено %d",
				i, test.expected, got)
		}
	}
}

func Test_real_Real(t *testing.T) {
	tests := []struct{ data float64 }{{64.64}, {64}}

	for i, test := range tests {
		v := newReal(test.data)
		got := v.Real()
		if got != v.v {
			t.Errorf("%d: ожидалось %f, получено %f",
				i, v.v, got)
		}
	}
}

func Test_real_Text(t *testing.T) {
	tests := []struct {
		data     float64
		expected string
	}{
		{data: 64, expected: "64"},
		{data: 64.64, expected: "64.64"},
		{data: -64.64, expected: "-64.64"},
	}

	for i, test := range tests {
		got := newReal(test.data).Text()
		if got != test.expected {
			t.Errorf("%d: ожидалось %s, получено %s",
				i, test.expected, got)
		}
	}
}
