package boolean

import (
	"testing"
)

func Test_False_Int(t *testing.T) {
	b := NewFalse()
	if b.Int() != 0 {
		t.Errorf("ожидалось 0, получено %d", b.Int())
	}
}

func Test_False_Real(t *testing.T) {
	b := NewFalse()
	if b.Real() != 0 {
		t.Errorf("ожидалось 0, получено %f", b.Real())
	}
}

func Test_False_Text(t *testing.T) {
	b := NewFalse()
	if b.Text() != "false" {
		t.Errorf("ожидалось 'false', получено %s", b.Text())
	}
}

func Test_False_Bool(t *testing.T) {
	b := NewFalse()
	if b.Bool() != false {
		t.Error("ожидалось false, получено true")
	}
}

func Test_True_Int(t *testing.T) {
	b := NewTrue()
	if b.Int() != 1 {
		t.Errorf("ожидалось 1, получено %d", b.Int())
	}
}

func Test_True_Real(t *testing.T) {
	b := NewTrue()
	if b.Real() != 1 {
		t.Errorf("ожидалось 1, получено %f", b.Real())
	}
}

func Test_True_Text(t *testing.T) {
	b := NewTrue()
	if b.Text() != "true" {
		t.Errorf("ожидалось 'true', получено %s", b.Text())
	}
}

func Test_True_Bool(t *testing.T) {
	b := NewTrue()
	if b.Bool() != true {
		t.Error("ожидалось true, получено false")
	}
}
