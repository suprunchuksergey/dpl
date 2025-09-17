package value

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"unicode"
)

func Test_Type(t *testing.T) {
	tests := []struct {
		value    Value
		expected string
	}{
		{Int(81), IntType},
		{Real(8.1), RealType},
		{Text("text"), TextType},
		{Text("8.1"), TextType},
		{Bool(false), BoolType},
		{Bool(true), BoolType},
		{Null(), NullType},
		{Array(), ArrayType},
		{Array(Int(81)), ArrayType},
		{Object(), ObjectType},
		{Object(KV{Text("text"), Int(81)}), ObjectType},
		{Function(nil), FunctionType},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, test.value.Type())
	}
}

func Test_Int(t *testing.T) {
	tests := []struct {
		value         Value
		expectedValue int64
		expectedError error
	}{
		{Int(243), 243, nil},
		{Int(-243), -243, nil},

		{Real(243), 243, nil},
		{Real(-243), -243, nil},
		{Real(24.3), 24, nil},
		{Real(-24.3), -24, nil},

		{Bool(false), 0, nil},
		{Bool(true), 1, nil},

		{Text("text"), 0, nil},
		{Text("243"), 243, nil},
		{Text("-243"), -243, nil},
		{Text("24.3"), 24, nil},
		{Text("-24.3"), -24, nil},

		{Null(), 0, nil},

		{Array(), 0, conversionError(ArrayType, IntType)},
		{Array(Int(243)), 0, conversionError(ArrayType, IntType)},

		{Object(), 0, conversionError(ObjectType, IntType)},
		{Object(KV{Text("text"), Int(243)}), 0, conversionError(ObjectType, IntType)},

		{Function(nil), 0, conversionError(FunctionType, IntType)},
	}

	for _, test := range tests {
		v, err := test.value.Int()

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_Real(t *testing.T) {
	tests := []struct {
		value         Value
		expectedValue float64
		expectedError error
	}{
		{Int(243), 243, nil},
		{Int(-243), -243, nil},

		{Real(243), 243, nil},
		{Real(-243), -243, nil},
		{Real(24.3), 24.3, nil},
		{Real(-24.3), -24.3, nil},

		{Bool(false), 0, nil},
		{Bool(true), 1, nil},

		{Text("text"), 0, nil},
		{Text("243"), 243, nil},
		{Text("-243"), -243, nil},
		{Text("24.3"), 24.3, nil},
		{Text("-24.3"), -24.3, nil},

		{Null(), 0, nil},

		{Array(), 0, conversionError(ArrayType, RealType)},
		{Array(Int(243)), 0, conversionError(ArrayType, RealType)},

		{Object(), 0, conversionError(ObjectType, RealType)},
		{Object(KV{Text("text"), Int(243)}), 0, conversionError(ObjectType, RealType)},

		{Function(nil), 0, conversionError(FunctionType, RealType)},
	}

	for _, test := range tests {
		v, err := test.value.Real()

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_Text(t *testing.T) {
	tests := []struct {
		value    Value
		expected string
	}{
		{Int(81), "81"},
		{Int(-81), "-81"},

		{Real(81), "81"},
		{Real(-81), "-81"},
		{Real(8.1), "8.1"},
		{Real(-8.1), "-8.1"},

		{Text("text"), "text"},
		{Text("8.1"), "8.1"},

		{Bool(false), "false"},
		{Bool(true), "true"},

		{Null(), "null"},

		{Array(), "[]"},
		{Array(Int(81)), "[81]"},
		{Array(Int(81), Array(Real(8.1))), "[81,[8.1]]"},

		{Object(), "{}"},
		{Object(KV{Text("text"), Int(81)}), "{text:81}"},

		{Function(nil), FunctionType},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, test.value.Text())
	}
}

func Test_Bool(t *testing.T) {
	tests := []struct {
		value         Value
		expectedValue bool
		expectedError error
	}{
		{Int(0), false, nil},
		{Int(243), true, nil},
		{Int(-243), true, nil},

		{Real(0), false, nil},
		{Real(2.43), true, nil},
		{Real(-2.43), true, nil},

		{Bool(false), false, nil},
		{Bool(true), true, nil},

		{Text(""), false, nil},
		{Text("0"), true, nil},
		{Text("text"), true, nil},

		{Null(), false, nil},

		{Array(), false, nil},
		{Array(Int(243)), true, nil},

		{Object(), false, nil},
		{Object(KV{Text("text"), Int(243)}), true, nil},

		{Function(nil), false, conversionError(FunctionType, BoolType)},
	}

	for _, test := range tests {
		v, err := test.value.Bool()

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_IsReal(t *testing.T) {
	tests := []struct {
		value    Value
		expected bool
	}{
		{Int(81), false},

		{Real(81), true},
		{Real(8.1), true},

		{Text("text"), false},
		{Text("81"), false},
		{Text("8.1"), true},
		{Text("-8.1"), true},
		{Text("+8.1"), true},

		{Bool(false), false},
		{Bool(true), false},

		{Null(), false},

		{Array(), false},
		{Array(Int(81)), false},

		{Object(), false},
		{Object(KV{Text("text"), Int(81)}), false},

		{Function(nil), false},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, test.value.IsReal())
	}
}

func Test_IsText(t *testing.T) {
	tests := []struct {
		value    Value
		expected bool
	}{
		{Int(81), false},

		{Real(81), false},
		{Real(8.1), false},

		{Text("text"), true},
		{Text("81"), true},
		{Text("8.1"), true},

		{Bool(false), false},
		{Bool(true), false},

		{Null(), false},

		{Array(), false},
		{Array(Int(81)), false},

		{Object(), false},
		{Object(KV{Text("text"), Int(81)}), false},

		{Function(nil), false},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, test.value.IsText())
	}
}

func Test_ElByIndex(t *testing.T) {
	tests := []struct {
		value,
		index,
		expectedValue Value
		expectedError error
	}{
		{Text("text"), Int(-1), nil, indexOutOfRange()},
		{Text("text"), Int(0), Text("t"), nil},
		{Text("text"), Int(1), Text("e"), nil},
		{Text("text"), Int(2), Text("x"), nil},
		{Text("text"), Int(3), Text("t"), nil},
		{Text("text"), Int(4), nil, indexOutOfRange()},

		{Text("текст"), Int(-1), nil, indexOutOfRange()},
		{Text("текст"), Int(0), Text("т"), nil},
		{Text("текст"), Int(1), Text("е"), nil},
		{Text("текст"), Int(2), Text("к"), nil},
		{Text("текст"), Int(3), Text("с"), nil},
		{Text("текст"), Int(4), Text("т"), nil},
		{Text("текст"), Int(5), nil, indexOutOfRange()},

		{Array(Int(81), Text("текст"), Bool(true)), Int(-1), nil, indexOutOfRange()},
		{Array(Int(81), Text("текст"), Bool(true)), Int(0), Int(81), nil},
		{Array(Int(81), Text("текст"), Bool(true)), Int(1), Text("текст"), nil},
		{Array(Int(81), Text("текст"), Bool(true)), Int(2), Bool(true), nil},
		{Array(Int(81), Text("текст"), Bool(true)), Int(3), nil, indexOutOfRange()},

		{Object(), Text("текст"), Null(), nil},
		{Object(
			KV{Text("текст"), Int(81)},
			KV{Text("текст2"), Bool(false)},
		), Text("текст"), Int(81), nil},
		{Object(
			KV{Text("текст"), Int(81)},
			KV{Text("текст2"), Bool(false)},
		), Text("текст2"), Bool(false), nil},

		{Int(81), Int(0), nil, noIndexSupport(IntType)},
		{Real(8.1), Int(0), nil, noIndexSupport(RealType)},
		{Bool(true), Int(0), nil, noIndexSupport(BoolType)},
		{Null(), Int(0), nil, noIndexSupport(NullType)},
		{Function(nil), Int(0), nil, noIndexSupport(FunctionType)},
	}

	for _, test := range tests {
		v, err := test.value.ElByIndex(test.index)

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_SetElByIndex(t *testing.T) {
	tests := []struct {
		value,
		index,
		newElement Value
		expected error
	}{
		{Array(Text("text"), Int(81)), Int(-1), nil, indexOutOfRange()},
		{Array(Text("text"), Int(81)), Int(0), Real(8.1), nil},
		{Array(Text("text"), Int(81)), Int(2), nil, indexOutOfRange()},

		{Object(), Text("text"), Int(81), nil},
		{Object(KV{Text("text"), Bool(false)}), Text("text"), Int(81), nil},
		{Object(KV{Text("text"), Bool(false)}), Int(1), Int(81), nil},

		{Text("text"), Int(0), nil, noSetIndexSupport(TextType)},
		{Int(81), Int(0), nil, noSetIndexSupport(IntType)},
		{Real(8.1), Int(0), nil, noSetIndexSupport(RealType)},
		{Bool(false), Int(0), nil, noSetIndexSupport(BoolType)},
		{Null(), Int(0), nil, noSetIndexSupport(NullType)},
		{Function(nil), Int(0), nil, noSetIndexSupport(FunctionType)},
	}

	for _, test := range tests {
		err := test.value.SetElByIndex(test.index, test.newElement)

		if test.expected != nil {
			assert.EqualError(t, err, test.expected.Error())
		} else {
			assert.NoError(t, err)

			newElement, err := test.value.ElByIndex(test.index)
			assert.NoError(t, err)

			assert.Equal(t, test.newElement, newElement)
		}
	}
}

func Test_Iter(t *testing.T) {
	tests := []struct {
		value          Value
		expectedValues []Value
		expectedError  error
	}{
		{Int(3), []Value{Int(0), Int(1), Int(2)}, nil},
		{Real(3.3), []Value{Int(0), Int(1), Int(2)}, nil},
		{Text("text"), []Value{Int(0), Int(1), Int(2), Int(3)}, nil},
		{Text("мяч"), []Value{Int(0), Int(1), Int(2)}, nil},
		{Array(Text("мяч"), Int(27), Real(27.7)), []Value{Int(0), Int(1), Int(2)}, nil},
		{Object(KV{Text("мяч"), Int(27)}), []Value{Text("мяч")}, nil},

		{Bool(false), nil, noIterSupport(BoolType)},
		{Null(), nil, noIterSupport(NullType)},
		{Function(nil), nil, noIterSupport(FunctionType)},
	}

	for _, test := range tests {
		iter, err := test.value.Iter()

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)

			values := make([]Value, 0, len(test.expectedValues))
			for i := range iter {
				values = append(values, i)
			}

			assert.Equal(t, test.expectedValues, values)
		}
	}
}

func Test_Iter2(t *testing.T) {
	tests := []struct {
		value          Value
		expectedValues [][2]Value
		expectedError  error
	}{
		{Text("мяч"), [][2]Value{
			{Int(0), Text("м")},
			{Int(1), Text("я")},
			{Int(2), Text("ч")},
		}, nil},
		{Text("text"), [][2]Value{
			{Int(0), Text("t")},
			{Int(1), Text("e")},
			{Int(2), Text("x")},
			{Int(3), Text("t")},
		}, nil},
		{Array(Text("мяч"), Real(27.7), Bool(true)), [][2]Value{
			{Int(0), Text("мяч")},
			{Int(1), Real(27.7)},
			{Int(2), Bool(true)},
		}, nil},
		{Object(KV{Text("мяч"), Int(7)}), [][2]Value{
			{Text("мяч"), Int(7)},
		}, nil},

		{Int(27), nil, noIterSupport2(IntType)},
		{Real(2.7), nil, noIterSupport2(RealType)},
		{Bool(false), nil, noIterSupport2(BoolType)},
		{Null(), nil, noIterSupport2(NullType)},
		{Function(nil), nil, noIterSupport2(FunctionType)},
	}

	for _, test := range tests {
		iter, err := test.value.Iter2()

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)

			values := make([][2]Value, 0, len(test.expectedValues))
			for i, j := range iter {
				values = append(values, [2]Value{i, j})
			}

			assert.Equal(t, test.expectedValues, values)
		}
	}
}

func Test_Call(t *testing.T) {
	tests := []struct {
		value         Value
		args          []Value
		expectedValue Value
		expectedError error
	}{
		{
			Function(func(args ...Value) (Value, error) {
				a, _ := args[0].Int()
				b, _ := args[1].Int()
				return Int(a + b), nil
			}),
			[]Value{Int(9), Int(27)},
			Int(36),
			nil,
		},

		{Int(27), nil, nil, noCallSupport(IntType)},
		{Real(2.7), nil, nil, noCallSupport(RealType)},
		{Bool(true), nil, nil, noCallSupport(BoolType)},
		{Text("text"), nil, nil, noCallSupport(TextType)},
		{Array(), nil, nil, noCallSupport(ArrayType)},
		{Object(), nil, nil, noCallSupport(ObjectType)},
		{Null(), nil, nil, noCallSupport(NullType)},
	}

	for _, test := range tests {
		value, err := test.value.Call(test.args...)

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, value)
		}
	}
}

func Test_skip(t *testing.T) {
	tests := []struct {
		data     string
		index    int
		fun      func(rune) bool
		expected int
	}{
		{"", 0, unicode.IsSpace, 0},
		{"			", 0, unicode.IsSpace, 3},
		{"			", 1, unicode.IsSpace, 3},
		{"			243", 0, unicode.IsSpace, 3},
		{"243", 0, unicode.IsSpace, 0},
		{"			", 216, unicode.IsSpace, 216},

		{"", 0, unicode.IsDigit, 0},
		{"243", 0, unicode.IsDigit, 3},
		{"243", 1, unicode.IsDigit, 3},
		{"243data", 0, unicode.IsDigit, 3},
		{"243", 0, unicode.IsDigit, 3},
		{"243", 216, unicode.IsDigit, 216},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, skip([]rune(test.data), test.index, test.fun))
	}
}

func Test_skipDigits(t *testing.T) {
	tests := []struct {
		data     string
		index    int
		expected int
	}{
		{"", 0, 0},
		{"			", 0, 0},
		{"			", 1, 1},
		{"data", 3, 3},
		{"243", 0, 3},
		{"243", 216, 216},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, skipDigits([]rune(test.data), test.index))
	}
}

func Test_skipSpaces(t *testing.T) {
	tests := []struct {
		data     string
		index    int
		expected int
	}{
		{"", 0, 0},
		{"			", 0, 3},
		{"			", 1, 3},
		{"data", 3, 3},
		{"243", 0, 0},
		{"			", 216, 216},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, skipSpaces([]rune(test.data), test.index))
	}
}

func Test_textIsReal(t *testing.T) {
	tests := []struct {
		data     string
		expected bool
	}{
		{"", false},
		{"	", false},
		{"-", false},
		{"+", false},
		{"--", false},
		{"++", false},
		{"	-", false},
		{"	+", false},
		{"	--", false},
		{"	++", false},
		{"data", false},
		{"	-data", false},
		{"	+data", false},
		{"	--data", false},
		{"	++data", false},
		{"216", false},
		{"-216", false},
		{"+216", false},
		{"216data", false},
		{"-216data", false},
		{"+216data", false},
		{".", false},
		{"..", false},
		{"-.", false},
		{"+.", false},
		{"-..", false},
		{"+..", false},
		{"	-.", false},
		{"	+.", false},
		{"	-..", false},
		{"	+..", false},

		{"2.16", true},
		{"-2.16", true},
		{"+2.16", true},
		{"	2.16", true},
		{"	-2.16", true},
		{"	+2.16", true},
		{"	.16", true},
		{"	-.16", true},
		{"	+.16", true},
		{"	2.", true},
		{"	-2.", true},
		{"	+2.", true},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, textIsReal([]rune(test.data)))
	}
}

func Test_textToInt(t *testing.T) {
	tests := []struct {
		data     string
		expected int64
	}{
		{"", 0},
		{"	", 0},
		{"-", 0},
		{"+", 0},
		{"--", 0},
		{"++", 0},
		{"	-", 0},
		{"	+", 0},
		{"	--", 0},
		{"	++", 0},
		{"data", 0},
		{"	-data", 0},
		{"	+data", 0},
		{"	--data", 0},
		{"	++data", 0},
		{".", 0},
		{"..", 0},
		{"-.", 0},
		{"+.", 0},
		{"-..", 0},
		{"+..", 0},
		{"	-.", 0},
		{"	+.", 0},
		{"	-..", 0},
		{"	+..", 0},
		{"	.16", 0},
		{"	-.16", 0},
		{"	+.16", 0},

		{"216", 216},
		{"-216", -216},
		{"+216", 216},
		{"216data", 216},
		{"-216data", -216},
		{"+216data", 216},
		{"2.16", 2},
		{"-2.16", -2},
		{"+2.16", 2},
		{"	2.16", 2},
		{"	-2.16", -2},
		{"	+2.16", 2},
		{"	2.", 2},
		{"	-2.", -2},
		{"	+2.", 2},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, textToInt([]rune(test.data)))
	}
}

func Test_textToReal(t *testing.T) {
	tests := []struct {
		data     string
		expected float64
	}{
		{"", 0},
		{"	", 0},
		{"-", 0},
		{"+", 0},
		{"--", 0},
		{"++", 0},
		{"	-", 0},
		{"	+", 0},
		{"	--", 0},
		{"	++", 0},
		{"data", 0},
		{"	-data", 0},
		{"	+data", 0},
		{"	--data", 0},
		{"	++data", 0},
		{".", 0},
		{"..", 0},
		{"-.", 0},
		{"+.", 0},
		{"-..", 0},
		{"+..", 0},
		{"	-.", 0},
		{"	+.", 0},
		{"	-..", 0},
		{"	+..", 0},

		{"216", 216},
		{"-216", -216},
		{"+216", 216},
		{"216data", 216},
		{"-216data", -216},
		{"+216data", 216},

		{"	.16", .16},
		{"	-.16", -.16},
		{"	+.16", .16},
		{"	2.", 2},
		{"	-2.", -2},
		{"	+2.", 2},
		{"2.16", 2.16},
		{"-2.16", -2.16},
		{"+2.16", 2.16},
		{"	2.16", 2.16},
		{"	-2.16", -2.16},
		{"	+2.16", 2.16},

		{"2.16.6", 2.16},
		{"-2.16.6", -2.16},
		{"+2.16.6", 2.16},
		{"2.data.6", 2},
		{"-2.data.6", -2},
		{"+2.data.6", 2},
		{"2..6", 2},
		{"-2..6", -2},
		{"+2..6", 2},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, textToReal([]rune(test.data)))
	}
}

func Test_Len(t *testing.T) {
	tests := []struct {
		value         Value
		expectedValue int64
		expectedError error
	}{
		{Text(""), 0, nil},
		{Text("text"), 4, nil},
		{Text("текст"), 5, nil},

		{Array(), 0, nil},
		{Array(Int(81)), 1, nil},
		{Array(Int(81), Text("text")), 2, nil},

		{Object(), 0, nil},
		{Object(KV{Text("text"), Null()}, KV{Text("текст"), Null()}), 2, nil},

		{Int(81), 0, noLenSupport(IntType)},
		{Real(8.1), 0, noLenSupport(RealType)},
		{Bool(true), 0, noLenSupport(BoolType)},
		{Null(), 0, noLenSupport(NullType)},
		{Function(nil), 0, noLenSupport(FunctionType)},
	}

	for _, test := range tests {
		v, err := test.value.Len()

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_Append(t *testing.T) {
	tests := []struct {
		value         Value
		values        []Value
		expectedValue Value
		expectedError error
	}{
		{Array(), []Value{Int(243)}, Array(Int(243)), nil},

		{Array(), nil, Array(), nil},
		{Array(), []Value{}, Array(), nil},

		{Array(Int(243)), nil, Array(Int(243)), nil},
		{Array(Int(243)), []Value{}, Array(Int(243)), nil},

		{Array(Int(243)), []Value{Int(32)}, Array(Int(243), Int(32)), nil},
		{
			Array(Int(243), Int(32)),
			[]Value{Bool(true), Text("text")},
			Array(Int(243), Int(32), Bool(true), Text("text")),
			nil,
		},

		{Int(243), nil, nil, noAppendSupport(IntType)},
		{Real(2.43), nil, nil, noAppendSupport(RealType)},
		{Bool(false), nil, nil, noAppendSupport(BoolType)},
		{Bool(true), nil, nil, noAppendSupport(BoolType)},
		{Text("text"), nil, nil, noAppendSupport(TextType)},
		{Null(), nil, nil, noAppendSupport(NullType)},
		{Object(), nil, nil, noAppendSupport(ObjectType)},
		{Function(nil), nil, nil, noAppendSupport(FunctionType)},
	}

	for _, test := range tests {
		v, err := test.value.Append(test.values...)

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}
