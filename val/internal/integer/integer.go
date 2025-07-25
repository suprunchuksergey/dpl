package integer

import (
	"fmt"
	"strconv"
)

type Integer int64

func (i Integer) String() string { return fmt.Sprintf("int %d", i) }

func (i Integer) Int() int64 { return int64(i) }

func (i Integer) Real() float64 { return float64(i) }

func (i Integer) Text() string { return strconv.FormatInt(int64(i), 10) }

func (i Integer) Bool() bool { return i != 0 }

func New(val int64) Integer { return Integer(val) }
