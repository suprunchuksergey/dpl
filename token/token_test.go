package token

import (
	"github.com/suprunchuksergey/dpl/pos"
	"testing"
)

func Test_Is(t *testing.T) {
	tok := New(Eq, pos.New())
	if tok.Is(Eq) == false {
		t.Error("ожидалось true, получено false")
	}

	if tok.Is(Lte) == true {
		t.Error("ожидалось false, получено true")
	}
}

func Test_OneOf(t *testing.T) {
	tok := New(Eq, pos.New())

	if tok.OneOf(Lte, Gte, Sub, Add, Mul, Eq, Neq) == false {
		t.Error("ожидалось true, получено false")
	}

	if tok.OneOf(Concat, Int, Real, Text) == true {
		t.Error("ожидалось false, получено true")
	}
}
