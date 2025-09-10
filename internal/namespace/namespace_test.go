package namespace

import (
	"github.com/stretchr/testify/assert"
	"github.com/suprunchuksergey/dpl/internal/value"
	"testing"
)

func Test_Get(t *testing.T) {
	parent := New(map[string]value.Value{
		"name": value.Text("сергей"),
	})

	n := parent.New(map[string]value.Value{
		"age": value.Int(23),
	})

	tests := []struct {
		name          string
		expectedValue value.Value
		expectedError error
	}{
		{"age", value.Int(23), nil},
		{"name", value.Text("сергей"), nil},
		{"surname", nil, VarDoesNotExist("surname")},
	}

	for _, test := range tests {
		v, err := n.Get(test.name)

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_Create(t *testing.T) {
	n := New(map[string]value.Value{
		"age": value.Int(23),
	})

	tests := []struct {
		name          string
		value         value.Value
		expectedError error
	}{
		{"name", value.Text("сергей"), nil},
		{"surname", value.Text("без фамилии"), nil},
		{"age", value.Int(23), VarAlreadyExists("age")},
	}

	for _, test := range tests {
		err := n.Create(test.name, test.value)

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)

			v, err := n.Get(test.name)
			assert.NoError(t, err)

			assert.Equal(t, test.value, v)
		}
	}
}

func Test_Set(t *testing.T) {
	parent := New(map[string]value.Value{
		"name": value.Text("сергей"),
	})

	n := parent.New(map[string]value.Value{
		"age": value.Int(23),
	})

	tests := []struct {
		name  string
		value value.Value
	}{
		{"name", value.Text("полина")},
		{"surname", value.Text("вишенка")},
		{"age", value.Int(26)},
	}

	for _, test := range tests {
		n.Set(test.name, test.value)

		v, err := n.Get(test.name)
		assert.NoError(t, err)

		assert.Equal(t, test.value, v)
	}
}
