package value

import (
	"fmt"
	"strconv"
)

type real struct{ v float64 }

func newReal(v float64) real { return real{v} }

func (r real) Int() int64 { return int64(r.v) }

func (r real) Real() float64 { return r.v }

func (r real) Text() string {
	return strconv.FormatFloat(r.v, 'g', -1, 64)
}

func (r real) Bool() bool { return r.v != 0 }

func (r real) IsInt() bool { return false }

func (r real) IsReal() bool { return true }

func (r real) IsText() bool { return false }

func (r real) IsBool() bool { return false }

func (r real) IsNull() bool { return false }

func (r real) String() string { return fmt.Sprintf("real %f", r.v) }

var _ Value = real{}
