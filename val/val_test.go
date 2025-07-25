package val

import "testing"

func Test_IsInt(t *testing.T) {
	tests := []struct {
		val      Val
		expected bool
	}{
		{Int(672), true},
		{Int(-42949672), true},
		{Text("	 672"), true},
		{Text("	 -42949672text"), true},
		{Text("	 +217728 text"), true},
		{Text(" +text67"), false},
		{Text(" 	-text7"), false},
		{Real(67.2), false},
		{Real(-4.2949672), false},
		{True(), false},
		{False(), false},
		{Null(), false},
	}

	for i, test := range tests {
		if IsInt(test.val) != test.expected {
			t.Errorf("%d: ожидалось %t, получено %t",
				i, test.expected, !test.expected)
		}
	}
}

func Test_IsReal(t *testing.T) {
	tests := []struct {
		val      Val
		expected bool
	}{
		{Int(672), false},
		{Int(-42949672), false},
		{Text("	 672"), false},
		{Text("	 -42949672text"), false},
		{Text("	 +217728 text"), false},
		{Text(" +text67"), false},
		{Text(" 	-text7"), false},
		{Real(67.2), true},
		{Real(-4.2949672), true},
		{Text("67.text"), true},
		{Text("	 -4.2949672"), true},
		{Text("21772.8"), true},
		{Text("21772."), true},
		{Text("	 -.21772"), true},
		{True(), false},
		{False(), false},
		{Null(), false},
	}

	for i, test := range tests {
		if IsReal(test.val) != test.expected {
			t.Errorf("%d: ожидалось %t, получено %t",
				i, test.expected, !test.expected)
		}
	}
}

func Test_IsText(t *testing.T) {
	tests := []struct {
		val      Val
		expected bool
	}{
		{Int(672), false},
		{Int(-42949672), false},
		{Text("	 672"), true},
		{Text("	 text67"), true},
		{Real(67.2), false},
		{Real(-4.2949672), false},
		{True(), false},
		{False(), false},
		{Null(), false},
	}

	for i, test := range tests {
		if IsText(test.val) != test.expected {
			t.Errorf("%d: ожидалось %t, получено %t",
				i, test.expected, !test.expected)
		}
	}
}

func Test_IsBool(t *testing.T) {
	tests := []struct {
		val      Val
		expected bool
	}{
		{Int(672), false},
		{Int(-42949672), false},
		{Text("	 672"), false},
		{Text("	 text67"), false},
		{Real(67.2), false},
		{Real(-4.2949672), false},
		{True(), true},
		{False(), true},
		{Null(), false},
	}

	for i, test := range tests {
		if IsBool(test.val) != test.expected {
			t.Errorf("%d: ожидалось %t, получено %t",
				i, test.expected, !test.expected)
		}
	}
}

func Test_IsNull(t *testing.T) {
	tests := []struct {
		val      Val
		expected bool
	}{
		{Int(672), false},
		{Int(-42949672), false},
		{Text("	 672"), false},
		{Text("	 text67"), false},
		{Real(67.2), false},
		{Real(-4.2949672), false},
		{True(), false},
		{False(), false},
		{Null(), true},
		{nil, true},
	}

	for i, test := range tests {
		if IsNull(test.val) != test.expected {
			t.Errorf("%d: ожидалось %t, получено %t",
				i, test.expected, !test.expected)
		}
	}
}
