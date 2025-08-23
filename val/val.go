package val

import (
	"fmt"
	"strconv"
	"unicode"
)

type val struct {
	//int64(INT) | float64(REAL) | string(TEXT) | bool(BOOL) |
	//[]Val(ARRAY) | map[string]Val(MAP) | nil(NULL)
	val any
}

func (val val) ToInt() int64 {
	switch v := val.val.(type) {
	case int64:
		return v
	case float64:
		return int64(v)
	case bool:
		if v {
			return 1
		}
		return 0
	case string:
		return textToInt([]rune(v))
	case nil:
		return 0
	default:
		panic(fmt.Sprintf("невозможно преобразовать %s в INT", val))
	}
}

func (val val) ToReal() float64 {
	switch v := val.val.(type) {
	case int64:
		return float64(v)
	case float64:
		return v
	case bool:
		if v {
			return 1
		}
		return 0
	case string:
		return textToReal([]rune(v))
	case nil:
		return 0
	default:
		panic(fmt.Sprintf("невозможно преобразовать %s в REAL", val))
	}
}

func (val val) ToText() string {
	switch v := val.val.(type) {
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case string:
		return v
	case nil:
		return ""
	default:
		panic(fmt.Sprintf("невозможно преобразовать %s в TEXT", val))
	}
}

func (val val) ToBool() bool {
	switch v := val.val.(type) {
	case int64:
		return v != 0
	case float64:
		return v != 0
	case bool:
		return v
	case string:
		return len(v) > 0
	case nil:
		return false
	default:
		panic(fmt.Sprintf("невозможно преобразовать %s в BOOL", val))
	}
}

func (val val) IsInt() bool {
	switch val.val.(type) {
	case int64:
		return true
	case string:
		return textIsInt([]rune(val.val.(string)))
	}
	return false
}

func (val val) IsReal() bool {
	switch val.val.(type) {
	case float64:
		return true
	case string:
		return textIsReal([]rune(val.val.(string)))
	}
	return false
}

func (val val) IsText() bool {
	_, ok := val.val.(string)
	return ok
}

func (val val) IsBool() bool {
	_, ok := val.val.(bool)
	return ok
}

func (val val) IsArray() bool {
	_, ok := val.val.([]Val)
	return ok
}

func (val val) IsMap() bool {
	_, ok := val.val.(map[string]Val)
	return ok
}

func (val val) IsNull() bool {
	return val.val == nil
}

func (val val) String() string {
	return fmt.Sprint(val.val)
}

type Val interface {
	//типы могут быть легко преобразованы друг в друга.

	ToInt() int64
	ToReal() float64
	ToText() string
	ToBool() bool

	//каждый метод проверки возвращает true, только если
	//значение точно соответствует указанному типу,
	//исключением являются строки, которые могут интерпретироваться как числа,
	//если начинаются с числа.

	IsInt() bool
	IsReal() bool

	IsText() bool
	IsBool() bool
	IsArray() bool
	IsMap() bool
	IsNull() bool

	fmt.Stringer
}

func Int(v int64) Val    { return val{v} }
func Real(v float64) Val { return val{v} }
func Text(v string) Val  { return val{v} }

func Bool(v bool) Val { return val{v} }
func True() Val       { return val{true} }
func False() Val      { return val{false} }

func Array(v []Val) Val        { return val{v} }
func Map(v map[string]Val) Val { return val{v} }
func Null() Val                { return val{nil} }

func skip(v []rune, i int, c func(rune) bool) int {
	for i < len(v) && c(v[i]) {
		i++
	}
	return i
}

func skipDigits(v []rune, i int) int {
	return skip(v, i, unicode.IsDigit)
}

func skipSpaces(v []rune, i int) int {
	return skip(v, i, unicode.IsSpace)
}

func textIsInt(v []rune) bool {
	i := skipSpaces(v, 0)
	if i == len(v) {
		return false
	}
	if v[i] == '-' || v[i] == '+' {
		i++
		if i == len(v) {
			return false
		}
	}
	return unicode.IsDigit(v[i])
}

func textIsReal(v []rune) bool {
	i := skipSpaces(v, 0)
	if i == len(v) {
		return false
	}
	if v[i] == '-' || v[i] == '+' {
		i++
		if i == len(v) {
			return false
		}
	}
	if unicode.IsDigit(v[i]) {
		i := skipDigits(v, i)
		return i != len(v) && v[i] == '.'
	}
	return v[i] == '.' &&
		i+1 != len(v) &&
		unicode.IsDigit(v[i+1])
}

func textToInt(v []rune) int64 {
	s := skipSpaces(v, 0)
	if s == len(v) {
		return 0
	}

	e := s
	if v[e] == '-' || v[e] == '+' {
		e++
		if e == len(v) {
			return 0
		}
	}

	e = skipDigits(v, e)

	n, _ := strconv.ParseInt(string(v[s:e]), 10, 64)
	return n
}

func textToReal(v []rune) float64 {
	s := skipSpaces(v, 0)
	if s == len(v) {
		return 0
	}

	e := s
	if v[e] == '-' || v[e] == '+' {
		e++
		if e == len(v) {
			return 0
		}
	}

	e = skipDigits(v, e)
	if e != len(v) && v[e] == '.' {
		e++
		e = skipDigits(v, e)
	}
	n, _ := strconv.ParseFloat(string(v[s:e]), 64)
	return n
}
