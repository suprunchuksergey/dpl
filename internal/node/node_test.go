package node

import (
	"github.com/stretchr/testify/assert"
	"github.com/suprunchuksergey/dpl/internal/namespace"
	"github.com/suprunchuksergey/dpl/internal/value"
	"testing"
)

func Test_Add(t *testing.T) {
	tests := []struct {
		node          Node
		expectedValue value.Value
		expectedError error
	}{
		{Add(Int(8), Int(9)), value.Int(17), nil},
		{Add(Int(8), Real(1.7)), value.Real(9.7), nil},
		{Add(Int(8), Bool(true)), value.Int(9), nil},
		{Add(Int(8), Null()), value.Int(8), nil},
		{Add(Int(8), Text("text")), value.Int(8), nil},
		{Add(Int(8), Text("9")), value.Int(17), nil},
		{Add(Int(8), Text("1.7")), value.Real(9.7), nil},

		{Add(Int(8), Array()), nil, opNotDefined("+", value.ArrayType)},
		{Add(Int(8), Object()), nil, opNotDefined("+", value.ObjectType)},
		{Add(Int(8), Function(nil)), nil, opNotDefined("+", value.FunctionType)},
	}

	for _, test := range tests {
		v, err := test.node.Exec(nil)

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_Sub(t *testing.T) {
	tests := []struct {
		node          Node
		expectedValue value.Value
		expectedError error
	}{
		{Sub(Int(17), Int(9)), value.Int(8), nil},
		{Sub(Int(17), Real(1.7)), value.Real(15.3), nil},
		{Sub(Int(17), Bool(true)), value.Int(16), nil},
		{Sub(Int(17), Null()), value.Int(17), nil},
		{Sub(Int(17), Text("text")), value.Int(17), nil},
		{Sub(Int(17), Text("9")), value.Int(8), nil},
		{Sub(Int(17), Text("1.7")), value.Real(15.3), nil},

		{Sub(Int(17), Array()), nil, opNotDefined("-", value.ArrayType)},
		{Sub(Int(17), Object()), nil, opNotDefined("-", value.ObjectType)},
		{Sub(Int(17), Function(nil)), nil, opNotDefined("-", value.FunctionType)},
	}

	for _, test := range tests {
		v, err := test.node.Exec(nil)

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_Mul(t *testing.T) {
	tests := []struct {
		node          Node
		expectedValue value.Value
		expectedError error
	}{
		{Mul(Int(8), Int(9)), value.Int(72), nil},
		{Mul(Int(8), Real(1.7)), value.Real(13.6), nil},
		{Mul(Int(8), Bool(true)), value.Int(8), nil},
		{Mul(Int(8), Null()), value.Int(0), nil},
		{Mul(Int(8), Text("text")), value.Int(0), nil},
		{Mul(Int(8), Text("9")), value.Int(72), nil},
		{Mul(Int(8), Text("1.7")), value.Real(13.6), nil},

		{Mul(Int(8), Array()), nil, opNotDefined("*", value.ArrayType)},
		{Mul(Int(8), Object()), nil, opNotDefined("*", value.ObjectType)},
		{Mul(Int(8), Function(nil)), nil, opNotDefined("*", value.FunctionType)},
	}

	for _, test := range tests {
		v, err := test.node.Exec(nil)

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_Div(t *testing.T) {
	tests := []struct {
		node          Node
		expectedValue value.Value
		expectedError error
	}{
		{Div(Int(243), Int(3)), value.Int(81), nil},
		{Div(Int(8), Real(1.25)), value.Real(6.4), nil},
		{Div(Int(8), Bool(true)), value.Int(8), nil},
		{Div(Int(8), Null()), nil, divByZero()},
		{Div(Int(8), Text("text")), nil, divByZero()},
		{Div(Int(243), Text("3")), value.Int(81), nil},
		{Div(Int(8), Text("1.25")), value.Real(6.4), nil},

		{Div(Int(243), Int(0)), nil, divByZero()},

		{Div(Int(8), Array()), nil, opNotDefined("/", value.ArrayType)},
		{Div(Int(8), Object()), nil, opNotDefined("/", value.ObjectType)},
		{Div(Int(8), Function(nil)), nil, opNotDefined("/", value.FunctionType)},
	}

	for _, test := range tests {
		v, err := test.node.Exec(nil)

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_Mod(t *testing.T) {
	tests := []struct {
		node          Node
		expectedValue value.Value
		expectedError error
	}{
		{Mod(Int(243), Int(16)), value.Int(3), nil},
		{Mod(Int(243), Real(1.75)), value.Real(1.5), nil},
		{Mod(Int(243), Bool(true)), value.Int(0), nil},

		{Mod(Int(243), Bool(false)), nil, divByZero()},
		{Mod(Int(243), Null()), nil, divByZero()},
		{Mod(Int(243), Int(0)), nil, divByZero()},
		{Mod(Int(243), Text("text")), nil, divByZero()},

		{Mod(Int(243), Text("16")), value.Int(3), nil},
		{Mod(Int(243), Text("1.75")), value.Real(1.5), nil},

		{Mod(Int(8), Array()), nil, opNotDefined("%", value.ArrayType)},
		{Mod(Int(8), Object()), nil, opNotDefined("%", value.ObjectType)},
		{Mod(Int(8), Function(nil)), nil, opNotDefined("%", value.FunctionType)},
	}

	for _, test := range tests {
		v, err := test.node.Exec(nil)

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_Concat(t *testing.T) {
	tests := []struct {
		node          Node
		expectedValue value.Value
		expectedError error
	}{
		{Concat(Text("имя: "), Text("сергей")), value.Text("имя: сергей"), nil},
		{Concat(Text("фамилия: "), Null()), value.Text("фамилия: null"), nil},
		{Concat(Text("возраст: "), Int(23)), value.Text("возраст: 23"), nil},
		{Concat(Text("является администратором: "), Bool(true)), value.Text("является администратором: true"), nil},
		{Concat(Real(2.3), Int(19683)), value.Text("2.319683"), nil},
		{Concat(
			Array(Int(23)),
			Object(KV{Text("имя"), Text("сергей")}),
		), value.Text("[23]{имя:сергей}"), nil},
		{Concat(Function(nil), Int(23)), value.Text(value.FunctionType + "23"), nil},
	}

	for _, test := range tests {
		v, err := test.node.Exec(nil)

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_Eq(t *testing.T) {
	tests := []struct {
		node          Node
		expectedValue value.Value
		expectedError error
	}{
		{Eq(Int(8), Int(8)), value.Bool(true), nil},
		{Eq(Int(8), Int(9)), value.Bool(false), nil},

		{Eq(Int(8), Real(8)), value.Bool(true), nil},
		{Eq(Int(8), Real(8.9)), value.Bool(false), nil},
		{Eq(Int(8), Real(9)), value.Bool(false), nil},

		{Eq(Int(0), Bool(false)), value.Bool(true), nil},
		{Eq(Int(1), Bool(true)), value.Bool(true), nil},

		{Eq(Bool(false), Bool(false)), value.Bool(true), nil},
		{Eq(Bool(true), Bool(true)), value.Bool(true), nil},
		{Eq(Bool(true), Bool(false)), value.Bool(false), nil},

		{Eq(Text("фамилия"), Text("фамилия")), value.Bool(true), nil},
		{Eq(Text("имя"), Text("фамилия")), value.Bool(false), nil},

		{Eq(Text("8рублей"), Int(8)), value.Bool(true), nil},
		{Eq(Text("8.8рублей"), Int(8)), value.Bool(false), nil},

		{Eq(Text("8долларов"), Text("8рублей")), value.Bool(false), nil},

		{Eq(Int(8), Array()), nil, opNotDefined("==", value.ArrayType)},
		{Eq(Int(8), Object()), nil, opNotDefined("==", value.ObjectType)},
		{Eq(Int(8), Function(nil)), nil, opNotDefined("==", value.FunctionType)},

		//спорные моменты
		{Eq(Int(0), Null()), value.Bool(true), nil},
		{Eq(Bool(false), Null()), value.Bool(true), nil},
	}

	for _, test := range tests {
		v, err := test.node.Exec(nil)

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_Neq(t *testing.T) {
	tests := []struct {
		node          Node
		expectedValue value.Value
		expectedError error
	}{
		{Neq(Int(8), Int(8)), value.Bool(false), nil},
		{Neq(Int(8), Int(9)), value.Bool(true), nil},

		{Neq(Int(8), Real(8)), value.Bool(false), nil},
		{Neq(Int(8), Real(8.9)), value.Bool(true), nil},
		{Neq(Int(8), Real(9)), value.Bool(true), nil},

		{Neq(Int(0), Bool(false)), value.Bool(false), nil},
		{Neq(Int(1), Bool(true)), value.Bool(false), nil},

		{Neq(Bool(false), Bool(false)), value.Bool(false), nil},
		{Neq(Bool(true), Bool(true)), value.Bool(false), nil},
		{Neq(Bool(true), Bool(false)), value.Bool(true), nil},

		{Neq(Text("фамилия"), Text("фамилия")), value.Bool(false), nil},
		{Neq(Text("имя"), Text("фамилия")), value.Bool(true), nil},

		{Neq(Text("8рублей"), Int(8)), value.Bool(false), nil},
		{Neq(Text("8.8рублей"), Int(8)), value.Bool(true), nil},

		{Neq(Text("8долларов"), Text("8рублей")), value.Bool(true), nil},

		{Neq(Int(8), Array()), nil, opNotDefined("!=", value.ArrayType)},
		{Neq(Int(8), Object()), nil, opNotDefined("!=", value.ObjectType)},
		{Neq(Int(8), Function(nil)), nil, opNotDefined("!=", value.FunctionType)},

		//спорные моменты
		{Neq(Int(0), Null()), value.Bool(false), nil},
		{Neq(Bool(false), Null()), value.Bool(false), nil},
	}

	for _, test := range tests {
		v, err := test.node.Exec(nil)

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_Lt(t *testing.T) {
	tests := []struct {
		node          Node
		expectedValue value.Value
		expectedError error
	}{
		{Lt(Int(8), Int(8)), value.Bool(false), nil},
		{Lt(Int(8), Int(9)), value.Bool(true), nil},

		{Lt(Int(8), Real(8)), value.Bool(false), nil},
		{Lt(Int(8), Real(8.9)), value.Bool(true), nil},
		{Lt(Int(8), Real(9)), value.Bool(true), nil},

		{Lt(Int(0), Bool(false)), value.Bool(false), nil},
		{Lt(Int(1), Bool(true)), value.Bool(false), nil},
		{Lt(Int(0), Bool(true)), value.Bool(true), nil},

		{Lt(Bool(false), Bool(true)), value.Bool(true), nil},
		{Lt(Bool(true), Bool(false)), value.Bool(false), nil},

		{Lt(Text("фамилия"), Text("фамилия")), value.Bool(false), nil},
		{Lt(Text("имя"), Text("фамилия")), value.Bool(true), nil},

		{Lt(Text("8рублей"), Int(8)), value.Bool(false), nil},
		{Lt(Int(8), Text("8.8рублей")), value.Bool(true), nil},

		{Lt(Text("8долларов"), Text("8рублей")), value.Bool(true), nil},

		{Lt(Int(8), Array()), nil, opNotDefined("<", value.ArrayType)},
		{Lt(Int(8), Object()), nil, opNotDefined("<", value.ObjectType)},
		{Lt(Int(8), Function(nil)), nil, opNotDefined("<", value.FunctionType)},

		//спорные моменты
		{Lt(Int(0), Null()), value.Bool(false), nil},
		{Lt(Bool(false), Null()), value.Bool(false), nil},
	}

	for _, test := range tests {
		v, err := test.node.Exec(nil)

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_Gt(t *testing.T) {
	tests := []struct {
		node          Node
		expectedValue value.Value
		expectedError error
	}{
		{Gt(Int(8), Int(8)), value.Bool(false), nil},
		{Gt(Int(9), Int(8)), value.Bool(true), nil},

		{Gt(Int(8), Real(8)), value.Bool(false), nil},
		{Gt(Real(8.9), Int(8)), value.Bool(true), nil},
		{Gt(Real(9), Int(8)), value.Bool(true), nil},

		{Gt(Int(0), Bool(false)), value.Bool(false), nil},
		{Gt(Int(1), Bool(true)), value.Bool(false), nil},
		{Gt(Bool(true), Int(0)), value.Bool(true), nil},

		{Gt(Bool(true), Bool(false)), value.Bool(true), nil},
		{Gt(Bool(false), Bool(true)), value.Bool(false), nil},

		{Gt(Text("фамилия"), Text("фамилия")), value.Bool(false), nil},
		{Gt(Text("фамилия"), Text("имя")), value.Bool(true), nil},

		{Gt(Text("8рублей"), Int(8)), value.Bool(false), nil},
		{Gt(Text("8.8рублей"), Int(8)), value.Bool(true), nil},

		{Gt(Text("8рублей"), Text("8долларов")), value.Bool(true), nil},

		{Gt(Int(8), Array()), nil, opNotDefined(">", value.ArrayType)},
		{Gt(Int(8), Object()), nil, opNotDefined(">", value.ObjectType)},
		{Gt(Int(8), Function(nil)), nil, opNotDefined(">", value.FunctionType)},

		//спорные моменты
		{Gt(Int(0), Null()), value.Bool(false), nil},
		{Gt(Bool(false), Null()), value.Bool(false), nil},
	}

	for _, test := range tests {
		v, err := test.node.Exec(nil)

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_Lte(t *testing.T) {
	tests := []struct {
		node          Node
		expectedValue value.Value
		expectedError error
	}{
		{Lte(Int(8), Int(8)), value.Bool(true), nil},
		{Lte(Int(8), Int(9)), value.Bool(true), nil},
		{Lte(Int(9), Int(8)), value.Bool(false), nil},

		{Lte(Int(8), Real(8)), value.Bool(true), nil},
		{Lte(Int(8), Real(8.9)), value.Bool(true), nil},
		{Lte(Int(8), Real(9)), value.Bool(true), nil},
		{Lte(Int(9), Real(8)), value.Bool(false), nil},

		{Lte(Int(0), Bool(false)), value.Bool(true), nil},
		{Lte(Int(1), Bool(true)), value.Bool(true), nil},
		{Lte(Int(0), Bool(true)), value.Bool(true), nil},

		{Lte(Bool(false), Bool(true)), value.Bool(true), nil},
		{Lte(Bool(true), Bool(false)), value.Bool(false), nil},

		{Lte(Text("фамилия"), Text("фамилия")), value.Bool(true), nil},
		{Lte(Text("имя"), Text("фамилия")), value.Bool(true), nil},

		{Lte(Text("8рублей"), Int(8)), value.Bool(true), nil},
		{Lte(Int(8), Text("8.8рублей")), value.Bool(true), nil},

		{Lte(Text("8долларов"), Text("8рублей")), value.Bool(true), nil},

		{Lte(Int(8), Array()), nil, opNotDefined("<=", value.ArrayType)},
		{Lte(Int(8), Object()), nil, opNotDefined("<=", value.ObjectType)},
		{Lte(Int(8), Function(nil)), nil, opNotDefined("<=", value.FunctionType)},

		//спорные моменты
		{Lte(Int(0), Null()), value.Bool(true), nil},
		{Lte(Bool(false), Null()), value.Bool(true), nil},
	}

	for _, test := range tests {
		v, err := test.node.Exec(nil)

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_Gte(t *testing.T) {
	tests := []struct {
		node          Node
		expectedValue value.Value
		expectedError error
	}{
		{Gte(Int(8), Int(8)), value.Bool(true), nil},
		{Gte(Int(9), Int(8)), value.Bool(true), nil},
		{Gte(Int(8), Int(9)), value.Bool(false), nil},

		{Gte(Int(8), Real(8)), value.Bool(true), nil},
		{Gte(Real(8.9), Int(8)), value.Bool(true), nil},
		{Gte(Real(9), Int(8)), value.Bool(true), nil},
		{Gte(Int(8), Real(9)), value.Bool(false), nil},

		{Gte(Int(0), Bool(false)), value.Bool(true), nil},
		{Gte(Int(1), Bool(true)), value.Bool(true), nil},
		{Gte(Bool(true), Int(0)), value.Bool(true), nil},

		{Gte(Bool(true), Bool(false)), value.Bool(true), nil},
		{Gte(Bool(false), Bool(true)), value.Bool(false), nil},

		{Gte(Text("фамилия"), Text("фамилия")), value.Bool(true), nil},
		{Gte(Text("фамилия"), Text("имя")), value.Bool(true), nil},

		{Gte(Text("8рублей"), Int(8)), value.Bool(true), nil},
		{Gte(Text("8.8рублей"), Int(8)), value.Bool(true), nil},

		{Gte(Text("8рублей"), Text("8долларов")), value.Bool(true), nil},

		{Gte(Int(8), Array()), nil, opNotDefined(">=", value.ArrayType)},
		{Gte(Int(8), Object()), nil, opNotDefined(">=", value.ObjectType)},
		{Gte(Int(8), Function(nil)), nil, opNotDefined(">=", value.FunctionType)},

		//спорные моменты
		{Gte(Int(0), Null()), value.Bool(true), nil},
		{Gte(Bool(false), Null()), value.Bool(true), nil},
	}

	for _, test := range tests {
		v, err := test.node.Exec(nil)

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_And(t *testing.T) {
	tests := []struct {
		node          Node
		expectedValue value.Value
		expectedError error
	}{
		{And(Bool(true), Bool(true)), value.Bool(true), nil},
		{And(Bool(true), Bool(false)), value.Bool(false), nil},
		{And(Bool(false), Bool(true)), value.Bool(false), nil},
		{And(Bool(false), Bool(false)), value.Bool(false), nil},

		{And(Int(1), Int(1)), value.Bool(true), nil},
		{And(Int(1), Int(0)), value.Bool(false), nil},
		{And(Int(0), Int(1)), value.Bool(false), nil},
		{And(Int(0), Int(0)), value.Bool(false), nil},

		{And(Real(1), Real(1)), value.Bool(true), nil},
		{And(Real(1), Real(0)), value.Bool(false), nil},
		{And(Real(0), Real(1)), value.Bool(false), nil},
		{And(Real(0), Real(0)), value.Bool(false), nil},

		{And(Array(Int(0)), Object(KV{Text("text"), Int(0)})), value.Bool(true), nil},
		{And(Array(Int(0)), Object()), value.Bool(false), nil},
		{And(Array(), Object(KV{Text("text"), Int(0)})), value.Bool(false), nil},
		{And(Array(), Object()), value.Bool(false), nil},

		{And(Text("text"), Text("text")), value.Bool(true), nil},
		{And(Text("text"), Text("")), value.Bool(false), nil},
		{And(Text(""), Text("text")), value.Bool(false), nil},
		{And(Text(""), Text("")), value.Bool(false), nil},

		{And(Bool(true), Null()), value.Bool(false), nil},
		{And(Null(), Bool(true)), value.Bool(false), nil},
		{And(Null(), Null()), value.Bool(false), nil},

		{And(Function(nil), Bool(true)), nil, opNotDefined("and", value.FunctionType)},
	}

	for _, test := range tests {
		v, err := test.node.Exec(nil)

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_Or(t *testing.T) {
	tests := []struct {
		node          Node
		expectedValue value.Value
		expectedError error
	}{
		{Or(Bool(true), Bool(true)), value.Bool(true), nil},
		{Or(Bool(true), Bool(false)), value.Bool(true), nil},
		{Or(Bool(false), Bool(true)), value.Bool(true), nil},
		{Or(Bool(false), Bool(false)), value.Bool(false), nil},

		{Or(Int(1), Int(1)), value.Bool(true), nil},
		{Or(Int(1), Int(0)), value.Bool(true), nil},
		{Or(Int(0), Int(1)), value.Bool(true), nil},
		{Or(Int(0), Int(0)), value.Bool(false), nil},

		{Or(Real(1), Real(1)), value.Bool(true), nil},
		{Or(Real(1), Real(0)), value.Bool(true), nil},
		{Or(Real(0), Real(1)), value.Bool(true), nil},
		{Or(Real(0), Real(0)), value.Bool(false), nil},

		{Or(Array(Int(0)), Object(KV{Text("text"), Int(0)})), value.Bool(true), nil},
		{Or(Array(Int(0)), Object()), value.Bool(true), nil},
		{Or(Array(), Object(KV{Text("text"), Int(0)})), value.Bool(true), nil},
		{Or(Array(), Object()), value.Bool(false), nil},

		{Or(Text("text"), Text("text")), value.Bool(true), nil},
		{Or(Text("text"), Text("")), value.Bool(true), nil},
		{Or(Text(""), Text("text")), value.Bool(true), nil},
		{Or(Text(""), Text("")), value.Bool(false), nil},

		{Or(Bool(true), Null()), value.Bool(true), nil},
		{Or(Null(), Bool(true)), value.Bool(true), nil},
		{Or(Null(), Null()), value.Bool(false), nil},

		{Or(Function(nil), Bool(true)), nil, opNotDefined("or", value.FunctionType)},
	}

	for _, test := range tests {
		v, err := test.node.Exec(nil)

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_Not(t *testing.T) {
	tests := []struct {
		node          Node
		expectedValue value.Value
		expectedError error
	}{
		{Not(Bool(false)), value.Bool(true), nil},
		{Not(Bool(true)), value.Bool(false), nil},

		{Not(Int(0)), value.Bool(true), nil},
		{Not(Int(1)), value.Bool(false), nil},

		{Not(Real(0)), value.Bool(true), nil},
		{Not(Real(1)), value.Bool(false), nil},

		{Not(Text("")), value.Bool(true), nil},
		{Not(Text("text")), value.Bool(false), nil},

		{Not(Array()), value.Bool(true), nil},
		{Not(Array(Int(0))), value.Bool(false), nil},

		{Not(Object()), value.Bool(true), nil},
		{Not(Object(KV{Text("text"), Int(0)})), value.Bool(false), nil},

		{Not(Null()), value.Bool(true), nil},

		{Not(Function(nil)), nil, opNotDefined("not", value.FunctionType)},
	}

	for _, test := range tests {
		v, err := test.node.Exec(nil)

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_Neg(t *testing.T) {
	tests := []struct {
		node          Node
		expectedValue value.Value
		expectedError error
	}{

		{Neg(Int(27)), value.Int(-27), nil},
		{Neg(Real(2.7)), value.Real(-2.7), nil},

		{Neg(Bool(true)), value.Int(-1), nil},

		{Neg(Text("")), value.Int(0), nil},
		{Neg(Text("27")), value.Int(-27), nil},
		{Neg(Text("2.7")), value.Real(-2.7), nil},

		{Neg(Null()), value.Int(0), nil},

		{Neg(Array()), nil, opNotDefined("унарный -", value.ArrayType)},
		{Neg(Object()), nil, opNotDefined("унарный -", value.ObjectType)},
		{Neg(Function(nil)), nil, opNotDefined("унарный -", value.FunctionType)},
	}

	for _, test := range tests {
		v, err := test.node.Exec(nil)

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_ElByIndex(t *testing.T) {
	tests := []struct {
		node          Node
		expectedValue value.Value
		expectedError error
	}{
		{ElByIndex(Text("value"), Int(0)), value.Text("v"), nil},
		{ElByIndex(Text("value"), Int(1)), value.Text("a"), nil},

		{ElByIndex(Text("value"), Bool(false)), value.Text("v"), nil},
		{ElByIndex(Text("value"), Bool(true)), value.Text("a"), nil},

		{ElByIndex(Text("value"), Null()), value.Text("v"), nil},

		{ElByIndex(Text("value"), Text("0")), value.Text("v"), nil},
		{ElByIndex(Text("value"), Text("1")), value.Text("a"), nil},

		{ElByIndex(Text("value"), Text("0.27")), value.Text("v"), nil},
		{ElByIndex(Text("value"), Text("1.27")), value.Text("a"), nil},

		{ElByIndex(Text("value"), Array()), nil, wrongIndex(value.ArrayType)},
		{ElByIndex(Text("value"), Object()), nil, wrongIndex(value.ObjectType)},
		{ElByIndex(Text("value"), Function(nil)), nil, wrongIndex(value.FunctionType)},

		{ElByIndex(Array(Int(27), Text("value")), Int(0)), value.Int(27), nil},
		{ElByIndex(Array(Int(27), Text("value")), Int(1)), value.Text("value"), nil},

		{ElByIndex(Array(Int(27)), Array()), nil, wrongIndex(value.ArrayType)},
		{ElByIndex(Array(Int(27)), Object()), nil, wrongIndex(value.ObjectType)},
		{ElByIndex(Array(Int(27)), Function(nil)), nil, wrongIndex(value.FunctionType)},

		{ElByIndex(Object(KV{Text("value"), Int(27)}), Text("value")), value.Int(27), nil},
		{ElByIndex(Object(KV{Text("[27]"), Text("value")}), Array(Int(27))), value.Text("value"), nil},

		{ElByIndex(Int(27), Int(0)), nil, opNotDefined("[<index>]", value.IntType)},
		{ElByIndex(Real(2.7), Int(0)), nil, opNotDefined("[<index>]", value.RealType)},
		{ElByIndex(Bool(true), Int(0)), nil, opNotDefined("[<index>]", value.BoolType)},
		{ElByIndex(Null(), Int(0)), nil, opNotDefined("[<index>]", value.NullType)},
		{ElByIndex(Function(nil), Int(0)), nil, opNotDefined("[<index>]", value.FunctionType)},
	}

	for _, test := range tests {
		v, err := test.node.Exec(nil)

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_Ident(t *testing.T) {
	n := namespace.New(map[string]value.Value{
		"name": value.Text("сергей"),
		"age":  value.Int(23),
	})

	tests := []struct {
		node          Node
		expectedValue value.Value
		expectedError error
	}{
		{Ident("name"), value.Text("сергей"), nil},
		{Ident("age"), value.Int(23), nil},

		{Ident("surname"), nil, namespace.VarDoesNotExist("surname")},
	}

	for _, test := range tests {
		v, err := test.node.Exec(n)

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_Create(t *testing.T) {
	n := namespace.New(map[string]value.Value{
		"surname": value.Text("вишенка"),
	})

	tests := []struct {
		node          Node
		expectedValue value.Value
		expectedError error
	}{
		{Create(Ident("name"), Text("сергей")), value.Text("сергей"), nil},
		{Create(Ident("age"), Int(23)), value.Int(23), nil},

		{Create(Ident("surname"), Text("вишенка")), nil, namespace.VarAlreadyExists("surname")},
	}

	for _, test := range tests {
		v, err := test.node.Exec(n)

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)

			name := test.node.(create).name.(ident).v
			v, err := n.Get(name)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_Set(t *testing.T) {
	n := namespace.New(map[string]value.Value{
		"surname": value.Text("не вишенка"),
		"array": value.Array(
			value.Int(625),
			value.Int(3125),
			value.Array(value.Bool(true)),
		),
		"object": value.Object(
			value.KV{value.Text("value"), value.Int(23)},
		),
	})

	tests := []struct {
		node          Node
		expectedValue value.Value
		expectedError error
	}{
		{Set(Ident("name"), Text("сергей")), value.Text("сергей"), nil},
		{Set(Ident("age"), Int(23)), value.Int(23), nil},
		{Set(Ident("surname"), Text("вишенка")), value.Text("вишенка"), nil},
		{Set(
			ElByIndex(Ident("array"), Int(0)),
			Int(512),
		), value.Int(512), nil},
		{Set(
			ElByIndex(ElByIndex(Ident("array"), Int(2)), Int(0)),
			Bool(false),
		), value.Bool(false), nil},
		{Set(
			ElByIndex(Ident("object"), Text("value")),
			Int(512),
		), value.Int(512), nil},
		{Set(
			ElByIndex(Ident("object"), Text("missing")),
			Bool(true),
		), value.Bool(true), nil},

		{Set(Int(512), Int(23)), nil, idExpected()},
		{Set(ElByIndex(Int(512), Int(0)), Int(23)), nil, idExpected()},
	}

	for _, test := range tests {
		v, err := test.node.Exec(n)

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)

			if id, ok := test.node.(set).name.(ident); ok {
				name := id.v
				v, err := n.Get(name)
				assert.NoError(t, err)
				assert.Equal(t, test.expectedValue, v)
			}
		}
	}
}

func Test_Block(t *testing.T) {
	tests := []struct {
		node          Node
		expectedValue value.Value
		expectedError error
	}{
		{
			Block(
				Text("text"),
			), value.Text("text"), nil,
		},
		{
			Block(
				Text("text"),
				Bool(true),
				Int(19683),
			), value.Int(19683), nil,
		},
	}

	for _, test := range tests {
		v, err := test.node.Exec(nil)

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_If(t *testing.T) {
	tests := []struct {
		node          Node
		expectedValue value.Value
		expectedError error
	}{
		{
			If(
				Branch{Bool(false), Text("первый")},
			), value.Null(), nil,
		},
		{
			If(
				Branch{Bool(true), Text("первый")},
			), value.Text("первый"), nil,
		},
		{
			If(
				Branch{Bool(false), Text("первый")},
				Branch{Bool(true), Text("второй")},
			), value.Text("второй"), nil,
		},
		{
			If(
				Branch{Bool(false), Text("первый")},
				Branch{Bool(true), Text("второй")},
				Branch{Bool(true), Text("третий")},
			), value.Text("второй"), nil,
		},
		{
			If(
				Branch{Bool(false), Text("первый")},
				Branch{Bool(false), Text("второй")},
				Branch{Bool(true), Text("третий")},
			), value.Text("третий"), nil,
		},
		{
			If(
				Branch{Bool(false), Text("первый")},
				Branch{Bool(false), Text("второй")},
				Branch{Bool(false), Text("третий")},
			), value.Null(), nil,
		},
	}

	for _, test := range tests {
		v, err := test.node.Exec(namespace.New(nil))

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_For(t *testing.T) {
	tests := []struct {
		node          Node
		expectedValue value.Value
		expectedError error
	}{
		{
			For(
				[]Node{Ident("i")},
				Int(10),
				Ident("i"),
			), value.Int(9), nil,
		},
		{
			For(
				[]Node{Ident("i"), Ident("j")},
				Text("text"),
				Concat(Ident("i"), Ident("j")),
			), value.Text("3t"), nil,
		},
		{
			For(
				[]Node{Ident("i"), Ident("j")},
				Array(Int(9), Int(27), Text("рубля")),
				Concat(Ident("i"), Ident("j")),
			), value.Text("2рубля"), nil,
		},

		{
			For(
				[]Node{},
				Int(10),
				Ident("i"),
			), nil, tooFewRecipients(),
		},
		{
			For(
				[]Node{Ident("i"), Ident("j"), Ident("k")},
				Int(10),
				Ident("i"),
			), nil, tooManyRecipients(),
		},
		{
			For(
				[]Node{Int(0)},
				Int(10),
				Ident("i"),
			), nil, idExpected(),
		},
	}

	for _, test := range tests {
		v, err := test.node.Exec(namespace.New(nil))

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_Call(t *testing.T) {
	tests := []struct {
		node          Node
		expectedValue value.Value
		expectedError error
	}{
		{Call(Ident("тестовая функция")), value.Text("привет имя по умолчанию"), nil},
		{Call(Ident("тестовая функция"), Text("сергей")), value.Text("привет сергей"), nil},
		{Call(Ident("тестовая функция"), Text("полина")), value.Text("привет полина"), nil},
		{Call(Ident("тестовая функция"), Text("полина"), Int(26)), value.Text("привет полина"), nil},

		{Call(Int(26)), nil, opNotDefined("вызов функции", value.IntType)},
		{Call(Real(2.6)), nil, opNotDefined("вызов функции", value.RealType)},
		{Call(Bool(true)), nil, opNotDefined("вызов функции", value.BoolType)},
		{Call(Null()), nil, opNotDefined("вызов функции", value.NullType)},
		{Call(Text("text")), nil, opNotDefined("вызов функции", value.TextType)},
		{Call(Array()), nil, opNotDefined("вызов функции", value.ArrayType)},
		{Call(Object()), nil, opNotDefined("вызов функции", value.ObjectType)},
	}

	for _, test := range tests {
		v, err := test.node.Exec(namespace.New(map[string]value.Value{
			"тестовая функция": value.Function(
				func(args ...value.Value) (value.Value, error) {
					n := value.Text("имя по умолчанию")
					if len(args) != 0 {
						n = args[0]
					}
					return value.Text("привет " + n.Text()), nil
				},
			),
		}))

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}

func Test_Function(t *testing.T) {
	tests := []struct {
		node          Node
		expectedValue value.Value
		expectedError error
	}{
		{Call(
			Function(
				Concat(Text("привет "), Ident("name")),
				Ident("name"),
			),
			Text("сергей"),
		),
			value.Text("привет сергей"), nil},
		{Call(
			Function(
				Div(Ident("a"), Ident("b")),
				Ident("a"), Ident("b"),
			),
			Int(729), Int(9),
		),
			value.Int(81), nil},
		{Call(
			Function(
				Div(Ident("a"), Ident("b")),
				Ident("a"), Ident("b"), Ident("c"),
			),
			Int(729), Int(9),
		),
			value.Int(81), nil},
		{Call(
			Function(
				Div(Ident("a"), Ident("b")),
				Ident("a"), Ident("b"),
			),
			Int(729), Int(9), Int(2187),
		),
			value.Int(81), nil},
		{Call(Function(Text("привет мир"))), value.Text("привет мир"), nil},
		{Call(Function(Text("привет мир")), Int(729)), value.Text("привет мир"), nil},
		{Call(
			Function(
				Block(
					Mul(Ident("a"), Ident("b")),
					Div(Ident("a"), Ident("b")),
				),
				Ident("a"), Ident("b"),
			),
			Int(729), Int(9),
		),
			value.Int(81), nil},
		{Call(
			Function(
				Block(
					Return(Mul(Ident("a"), Ident("b"))),
					Div(Ident("a"), Ident("b")),
				),
				Ident("a"), Ident("b"),
			),
			Int(729), Int(9),
		),
			value.Int(6561), nil},
	}

	for _, test := range tests {
		v, err := test.node.Exec(namespace.New(nil))

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}
