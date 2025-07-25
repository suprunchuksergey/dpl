package real

import (
	"strconv"
	"testing"
)

var val = 64.4

func Test_Int(t *testing.T) {
	r := New(val)
	if r.Int() != int64(val) {
		t.Errorf("ожидалось %d, получено %d", int64(val), r.Int())
	}
}

func Test_Real(t *testing.T) {
	r := New(val)
	if r.Real() != val {
		t.Errorf("ожидалось %f, получено %f", val, r.Real())
	}
}

func Test_Text(t *testing.T) {
	r := New(val)
	expected := strconv.FormatFloat(
		val, 'g', -1, 64)
	if r.Text() != expected {
		t.Errorf("ожидалось %s, получено %s", expected, r.Text())
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
