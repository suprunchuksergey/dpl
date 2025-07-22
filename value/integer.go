package value

import (
	"fmt"
	"strconv"
)

type integer struct{ v int64 }

func newInteger(v int64) integer { return integer{v} }

func (i integer) Int() int64 { return i.v }

func (i integer) Real() float64 { return float64(i.v) }

func (i integer) Text() string { return strconv.FormatInt(i.v, 10) }

func (i integer) Bool() bool { return i.v != 0 }

func (i integer) IsInt() bool { return true }

func (i integer) IsReal() bool { return false }

func (i integer) IsText() bool { return false }

func (i integer) IsBool() bool { return false }

func (i integer) IsNull() bool { return false }

func (i integer) String() string { return fmt.Sprintf("int %d", i.v) }

var _ Value = integer{}
