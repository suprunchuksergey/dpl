package value

type null struct{}

func newNull() null { return null{} }

func (n null) Int() int64 { panic("невозможно преобразовать null в int") }

func (n null) Real() float64 { panic("невозможно преобразовать null в real") }

func (n null) Text() string { panic("невозможно преобразовать null в text") }

func (n null) IsInt() bool { return false }

func (n null) IsReal() bool { return false }

func (n null) IsText() bool { return false }

func (n null) IsNull() bool { return true }

func (n null) String() string { return "null" }

var _ Value = null{}
