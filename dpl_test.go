package dpl

import (
	"github.com/suprunchuksergey/dpl/internal/value"
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
	ns := map[string]value.Value{
		"age": value.Int(23),
		"num": value.Real(2.3),
		"factorial": value.Function(func(args ...value.Value) (value.Value, error) {
			a := args[0]
			i, _ := a.Int()
			return value.Int(factorial(i)), nil
		}),
	}

	tests := []struct {
		data     string
		expected value.Value
	}{
		{"(19.683+.6) * (512/64)", value.Real(162.264)},
		{"(not true)+512", value.Int(512)},
		{"(not true)+512 == 512", value.Bool(true)},
		{"(not true)+512 == 512 and true", value.Bool(true)},
		{"(not true)+512 == 512 and true-false", value.Bool(true)},
		{"(not true)+512 == 512 and true-true", value.Bool(false)},
		{"(not true)+512 == 512 and true-true", value.Bool(false)},
		{"'512 рублей'+512", value.Int(1024)},
		{"'512 рублей'<=1024", value.Bool(true)},
		{"'512 рублей'+512<=1024", value.Bool(true)},
		{"'512 рублей'+512<=1024==1", value.Bool(true)},
		{"512||' '||'рублей'=='512 рублей'", value.Bool(true)},
		{"'a'<'b'", value.Bool(true)},
		{"true and true", value.Bool(true)},
		{"true and false", value.Bool(false)},
		{"true or false", value.Bool(true)},
		{"true and not false", value.Bool(true)},
		{"not (true and false)", value.Bool(true)},
		{"not null", value.Bool(true)},
		{"[5,1,2][0]", value.Int(5)},
		{"[[1,5],1,2][0][1]", value.Int(5)},
		{"'tests'[2]", value.Text("s")},
		{"[[1,'tests'],1,2][0][1]", value.Text("tests")},
		{"[[1,'tests'],1,2][0][1][1+1]", value.Text("s")},
		{"[{'tests':50},1,2][0]['tests']", value.Int(50)},
		{"[{'tests':50},1,2][0]['tests']*5||'рублей'", value.Text("250рублей")},
		{"age+num", value.Real(25.3)},
		{"age+num-[num][0];[5,1,2][0]", value.Int(5)},
		{`
users = [8,27,64];
users[0] = 125-16;
users[0]+users[1];
`, value.Int(136)},
		{`
users = [8,27,64];
users[0] = 125-16;
users = [32,1024];
users[0]+users[1];
`, value.Int(1056)},
		{`
users = [[8],105];
users[0][0] = 1056;
users[0][0];
`, value.Int(1056)},
		{`
users = [{ 'name': 'sergey', 'age': 23 }];
users[0]['name'] = 'polina';
users[0]['age'] = 26;
users[0]['name']||' '||users[0]['age'];
`, value.Text("polina 26")},
		{`
users = [{ 'name': 'sergey', 'age': 23 }];
users[0] = { 'name': 'polina', 'age': 26 };
users[0]['name']||' '||users[0]['age'];
`, value.Text("polina 26")},
		{`
a = 23;
b = a;
a==b;`, value.Bool(true)},
		{`
a = 23;
b = a;
a = 26;
a!=b;`, value.Bool(true)},
		{`
a = [23];
b = a;
a[0] = 26;
a[0]==b[0];`, value.Bool(true)},
		{`
a = [23];
b = a;
a[0] = 26;
a[0];`, value.Int(26)},
		{`
a = [23];
b = a;
a[0] = 26;
b[0];`, value.Int(26)},
		{`
a = [23];
b = a;
b[0] = 26;
a[0];`, value.Int(26)},

		{`if 8<81 {'8<81'}`, value.Text("8<81")},
		{`if 8>81 {'8>81'}`, value.Null()},
		{`if 8>81 {'8>81'} else {'8<=81'}`, value.Text("8<=81")},
		{`if 8>81 {'8>81'}
elif 81==81 {'81==81'}
else {'8<=81'}`, value.Text("81==81")},
		{`if 8>81 {'8>81'}
elif 81>81 {'81>81'}
elif 'polina'=='polina' {'polina'}
else {'8<=81'}`, value.Text("polina")},
		{`if 8>81 {'8>81'}
elif 81>81 {'81>81'}
elif 'polina'=='polina' {'polina';'sergey'}
else {'8<=81'}`, value.Text("sergey")},
		{`
a = [8,16,32];
for i in a {
	a[i]=a[i]/2;
};
a;
`, value.Array(value.Int(4), value.Int(8), value.Int(16)),
		},
		{`
a = [8,16,32];
for i,k in a {
	a[i]=k/2;
};
a;
`, value.Array(value.Int(4), value.Int(8), value.Int(16)),
		},
		{"factorial(8)", value.Int(40320)},
		{"factorial(7)", value.Int(5040)},
		{"factorial(8)+factorial(7)", value.Int(45360)},
		{"i = factorial(8)+factorial(7);i", value.Int(45360)},

		{"i = fn (j,k) {j+k};i(20,30)", value.Int(50)},
		{"i = fn () {50};i()", value.Int(50)},
		{"i = fn () {return 50};i()", value.Int(50)},

		{`i = fn (n) {
		if n > 10 {
		return 'n > 10'
	} else {return 'n<=10'}};i(20);
`, value.Text("n > 10")},
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
