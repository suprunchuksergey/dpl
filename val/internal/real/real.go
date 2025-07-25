package real

import (
	"fmt"
	"strconv"
)

type Real float64

func (r Real) String() string { return fmt.Sprintf("real %s", r.Text()) }

func (r Real) Int() int64 { return int64(r) }

func (r Real) Real() float64 { return float64(r) }

func (r Real) Text() string {
	return strconv.FormatFloat(
		float64(r), 'g', -1, 64)
}

func (r Real) Bool() bool { return r != 0 }

func New(val float64) Real { return Real(val) }
