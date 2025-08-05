package dpl

import (
	"github.com/suprunchuksergey/dpl/val"
	"testing"
)

func Test_Exec(t *testing.T) {
	v, err := Exec("(19.683+.6) * (512/64)")
	if err != nil {
		t.Error(err)
		return
	}

	expected := val.Real(162.264)

	if v != expected {
		t.Errorf("ожидалось %s, получено %s", expected, v)
	}
}
