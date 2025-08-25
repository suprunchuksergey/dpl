package namespace

import (
	"github.com/suprunchuksergey/dpl/val"
	"testing"
)

func Test_Get(t *testing.T) {
	tests := []struct {
		ns       Namespace
		key      string
		expected val.Val
	}{
		{
			New(map[string]val.Val{"переменная": val.Int(4096)}),
			"переменная", val.Int(4096)},
		{WithParent(New(map[string]val.Val{"переменная": val.Int(4096)}), nil),
			"переменная", val.Int(4096)},
	}
	for i, test := range tests {
		v, err := test.ns.Get(test.key)
		if err != nil {
			t.Error(err)
			continue
		}

		if v != test.expected {
			t.Errorf("#%d: expected %s, got %s", i, test.expected, v)
		}
	}
}

func Test_Set(t *testing.T) {
	ns := New(nil)
	k, v := "переменная", val.Real(4.096)
	ns.Set(k, v)
	got, err := ns.Get(k)
	if err != nil {
		t.Error(err)
		return
	}
	if got != v {
		t.Errorf("expected %s, got %s", v, got)
	}
}
