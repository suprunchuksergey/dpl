package value

// true
type t struct{}

func newTrue() t { return t{} }

func (t t) Int() int64 { return 1 }

func (t t) Real() float64 { return 1 }

func (t t) Text() string { return "true" }

func (t t) Bool() bool { return true }

func (t t) IsInt() bool { return false }

func (t t) IsReal() bool { return false }

func (t t) IsText() bool { return false }

func (t t) IsBool() bool { return true }

func (t t) IsNull() bool { return false }

func (t t) String() string { return "bool true" }

var _ Value = t{}

// false
type f struct{}

func newFalse() f { return f{} }

func (f f) Int() int64 { return 0 }

func (f f) Real() float64 { return 0 }

func (f f) Text() string { return "false" }

func (f f) Bool() bool { return false }

func (f f) IsInt() bool { return false }

func (f f) IsReal() bool { return false }

func (f f) IsText() bool { return false }

func (f f) IsBool() bool { return true }

func (f f) IsNull() bool { return false }

func (f f) String() string { return "bool false" }

var _ Value = f{}
