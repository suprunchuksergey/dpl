package boolean

type False struct{}

func (f False) String() string { return "bool false" }

func (f False) Int() int64 { return 0 }

func (f False) Real() float64 { return 0 }

func (f False) Text() string { return "false" }

func (f False) Bool() bool { return false }

func NewFalse() False { return False{} }

type True struct{}

func (t True) String() string { return "bool true" }

func (t True) Int() int64 { return 1 }

func (t True) Real() float64 { return 1 }

func (t True) Text() string { return "true" }

func (t True) Bool() bool { return true }

func NewTrue() True { return True{} }
