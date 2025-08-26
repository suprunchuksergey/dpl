package dpl

import (
	"github.com/suprunchuksergey/dpl/val"
	"testing"
)

func Test_Exec(t *testing.T) {
	ns := map[string]val.Val{
		"age": val.Int(23),
		"num": val.Real(2.3),
	}

	tests := []struct {
		data     string
		expected val.Val
	}{
		{"(19.683+.6) * (512/64)", val.Real(162.264)},
		{"(not true)+512", val.Int(512)},
		{"(not true)+512 == 512", val.True()},
		{"(not true)+512 == 512 and true", val.True()},
		{"(not true)+512 == 512 and true-false", val.True()},
		{"(not true)+512 == 512 and true-true", val.False()},
		{"(not true)+512 == 512 and true-true", val.False()},
		{"'512 рублей'+512", val.Int(1024)},
		{"'512 рублей'<=1024", val.True()},
		{"'512 рублей'+512<=1024", val.True()},
		{"'512 рублей'+512<=1024==1", val.True()},
		{"512||' '||'рублей'=='512 рублей'", val.True()},
		{"'a'<'b'", val.True()},
		{"true and true", val.True()},
		{"true and false", val.False()},
		{"true or false", val.True()},
		{"true and not false", val.True()},
		{"not (true and false)", val.True()},
		{"not null", val.True()},
		{"[5,1,2][0]", val.Int(5)},
		{"[[1,5],1,2][0][1]", val.Int(5)},
		{"'tests'[2]", val.Text("s")},
		{"[[1,'tests'],1,2][0][1]", val.Text("tests")},
		{"[[1,'tests'],1,2][0][1][1+1]", val.Text("s")},
		{"[{'tests':50},1,2][0]['tests']", val.Int(50)},
		{"[{'tests':50},1,2][0]['tests']*5||'рублей'", val.Text("250рублей")},
		{"age+num", val.Real(25.3)},
		{"age+num-[num][0];[5,1,2][0]", val.Int(5)},
	}

	for _, test := range tests {
		v, err := Exec(test.data, ns)
		if err != nil {
			t.Error(err)
			continue
		}

		if v != test.expected {
			t.Errorf("%q: ожидалось: %s, получено: %s", test.data, test.expected, v)
		}
	}
}
