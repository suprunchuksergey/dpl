package integer

import (
	"strconv"
	"testing"
)

var val int64 = 64

func Test_Int(t *testing.T) {
	i := New(val)
	if i.Int() != val {
		t.Errorf("ожидалось %d, получено %d", val, i.Int())
	}
}

func Test_Real(t *testing.T) {
	i := New(val)
	if i.Real() != float64(val) {
		t.Errorf("ожидалось %f, получено %f", float64(val), i.Real())
	}
}

func Test_Text(t *testing.T) {
	i := New(val)
	expected := strconv.FormatInt(val, 10)
	if i.Text() != expected {
		t.Errorf("ожидалось %s, получено %s", expected, i.Text())
	}
}

func Test_Bool(t *testing.T) {
	msg := "ожидалось %t, получено %t"

	neg := New(-1)
	if neg.Bool() == false {
		t.Errorf(msg, true, false)
	}

	pos := New(1)
	if pos.Bool() == false {
		t.Errorf(msg, true, false)
	}

	zero := New(0)
	if zero.Bool() == true {
		t.Errorf(msg, false, true)
	}
}
