package val

import (
	"fmt"
	"github.com/suprunchuksergey/dpl/val/internal/boolean"
	"github.com/suprunchuksergey/dpl/val/internal/integer"
	"github.com/suprunchuksergey/dpl/val/internal/real"
	"github.com/suprunchuksergey/dpl/val/internal/text"
)

type Val interface {
	fmt.Stringer
	Int() int64
	Real() float64
	Text() string
	Bool() bool
}

func Int(val int64) Val { return integer.New(val) }

func Real(val float64) Val { return real.New(val) }

func False() Val { return boolean.NewFalse() }

func True() Val { return boolean.NewTrue() }

func Bool(val bool) Val {
	if val {
		return True()
	}
	return False()
}

func Text(val string) Val { return text.New(val) }

func Null() Val { return nil }

func IsInt(val Val) bool {
	switch v := val.(type) {
	case integer.Integer:
		return true
	case text.Text:
		return v.CanInt()
	default:
		return false
	}
}

func IsReal(val Val) bool {
	switch v := val.(type) {
	case real.Real:
		return true
	case text.Text:
		return v.CanReal()
	default:
		return false
	}
}

func IsText(val Val) bool {
	_, ok := val.(text.Text)
	return ok
}

func IsBool(val Val) bool {
	switch val.(type) {
	case boolean.False, boolean.True:
		return true
	default:
		return false
	}
}

func IsNull(val Val) bool { return val == nil }
