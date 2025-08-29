package dpl

import (
	"github.com/suprunchuksergey/dpl/val"
	"reflect"
	"testing"
)

func factorial(n int64) int64 {
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}

func Test_Exec(t *testing.T) {
	ns := map[string]val.Val{
		"age": val.Int(23),
		"num": val.Real(2.3),
		"factorial": val.Fn(func(args []val.Val) (val.Val, error) {
			return val.Int(factorial(args[0].ToInt())), nil
		}, nil),
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
		{`
users = [8,27,64];
users[0] = 125-16;
users[0]+users[1];
`, val.Int(136)},
		{`
users = [8,27,64];
users[0] = 125-16;
users = [32,1024];
users[0]+users[1];
`, val.Int(1056)},
		{`
users = [[8],105];
users[0][0] = 1056;
users[0][0];
`, val.Int(1056)},
		{`
users = [{ 'name': 'sergey', 'age': 23 }];
users[0]['name'] = 'polina';
users[0]['age'] = 26;
users[0]['name']||' '||users[0]['age'];
`, val.Text("polina 26")},
		{`
users = [{ 'name': 'sergey', 'age': 23 }];
users[0] = { 'name': 'polina', 'age': 26 };
users[0]['name']||' '||users[0]['age'];
`, val.Text("polina 26")},
		{`
a = 23;
b = a;
a==b;`, val.True()},
		{`
a = 23;
b = a;
a = 26;
a!=b;`, val.True()},
		{`
a = [23];
b = a;
a[0] = 26;
a[0]==b[0];`, val.True()},
		{`
a = [23];
b = a;
a[0] = 26;
a[0];`, val.Int(26)},
		{`
a = [23];
b = a;
a[0] = 26;
b[0];`, val.Int(26)},
		{`
a = [23];
b = a;
b[0] = 26;
a[0];`, val.Int(26)},

		{`if 8<81 {'8<81'}`, val.Text("8<81")},
		{`if 8>81 {'8>81'}`, val.Null()},
		{`if 8>81 {'8>81'} else {'8<=81'}`, val.Text("8<=81")},
		{`if 8>81 {'8>81'}
elif 81==81 {'81==81'}
else {'8<=81'}`, val.Text("81==81")},
		{`if 8>81 {'8>81'}
elif 81>81 {'81>81'}
elif 'polina'=='polina' {'polina'}
else {'8<=81'}`, val.Text("polina")},
		{`if 8>81 {'8>81'}
elif 81>81 {'81>81'}
elif 'polina'=='polina' {'polina';'sergey'}
else {'8<=81'}`, val.Text("sergey")},
		{`
a = [8,16,32];
for i in a {
	a[i]=a[i]/2;
};
a;
`, val.Array([]val.Val{val.Int(4), val.Int(8), val.Int(16)}),
		},
		{`
a = [8,16,32];
for i,k in a {
	a[i]=k/2;
};
a;
`, val.Array([]val.Val{val.Int(4), val.Int(8), val.Int(16)}),
		},
		{"factorial(8)", val.Int(40320)},
		{"factorial(7)", val.Int(5040)},
		{"factorial(8)+factorial(7)", val.Int(45360)},
		{"i = factorial(8)+factorial(7);i", val.Int(45360)},

		{"i = fn (j,k) {j+k};i(20,30)", val.Int(50)},
		{"i = fn () {50};i()", val.Int(50)},
		{"i = fn () {return 50};i()", val.Int(50)},

		{`i = fn (n) {
		if n > 10 {
		return 'n > 10'
	} else {return 'n<=10'}};i(20);
`, val.Text("n > 10")},
	}

	for _, test := range tests {
		v, err := Exec(test.data, ns)
		if err != nil {
			t.Error(err)
			continue
		}

		if !reflect.DeepEqual(v, test.expected) {
			t.Errorf("%q: ожидалось: %s, получено: %s", test.data, test.expected, v)
		}
	}
}
