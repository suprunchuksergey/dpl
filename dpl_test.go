package dpl

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/suprunchuksergey/dpl/internal/value"
	"testing"
)

func Test_Exec(t *testing.T) {
	tests := []struct {
		program       string
		expectedValue value.Value
		expectedError error
	}{
		{`
factorial := (n) -> {
	if n <= 1 {
		return 1;
	};
	return n * factorial(n-1);
};

sum := 0;

for i in 8 {
	sum = sum + factorial(i);
};

sum;
`, value.Int(5914), nil},

		{`
arr := [];

for i in 4 {
	arr = append(arr,i);
};

arr = append(arr,arr[len(arr)-1]);

arr;
`, value.Array(value.Int(0), value.Int(1), value.Int(2), value.Int(3), value.Int(3)),
			nil},

		{`
arr := [];

for i in 4 {
	arr = append(arr,i);
};

len(arr);
`, value.Int(4),
			nil},
	}

	for _, test := range tests {
		v, err := Exec(test.program, nil)

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_builtinLen(t *testing.T) {
	tests := []struct {
		args          []value.Value
		expectedValue value.Value
		expectedError error
	}{
		{nil, nil, errors.New("len: требуется один аргумент")},

		{[]value.Value{value.Array()}, value.Int(0), nil},
		{[]value.Value{value.Array(value.Null())}, value.Int(1), nil},
		{[]value.Value{value.Array(value.Null(), value.Null())}, value.Int(2), nil},

		{[]value.Value{
			value.Array(value.Null(), value.Null()),
			value.Array(value.Null()),
		}, value.Int(2), nil},
	}

	for _, test := range tests {
		v, err := builtinLen(test.args...)

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_builtinAppend(t *testing.T) {
	tests := []struct {
		args          []value.Value
		expectedValue value.Value
		expectedError error
	}{
		{nil, nil, errors.New("append: требуется один аргумент")},

		{[]value.Value{value.Array()}, value.Array(), nil},
		{[]value.Value{value.Array(), value.Int(512)}, value.Array(value.Int(512)), nil},
		{
			[]value.Value{value.Array(value.Int(512)), value.Int(2187)},
			value.Array(value.Int(512), value.Int(2187)),
			nil,
		},
		{
			[]value.Value{value.Array(value.Int(512)), value.Int(2187), value.Int(1024)},
			value.Array(value.Int(512), value.Int(2187), value.Int(1024)),
			nil,
		},
	}

	for _, test := range tests {
		v, err := builtinAppend(test.args...)

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}
