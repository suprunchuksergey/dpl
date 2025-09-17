package dpl

import (
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
