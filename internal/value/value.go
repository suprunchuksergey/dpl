package value

import (
	"errors"
	"fmt"
	"iter"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

type Value interface {
	Type() string

	Int() (int64, error)
	Real() (float64, error)
	Text() string
	Bool() (bool, error)

	IsReal() bool
	IsText() bool

	ElByIndex(index Value) (Value, error)
	SetElByIndex(index, value Value) error

	Iter() (iter.Seq[Value], error)
	Iter2() (iter.Seq2[Value, Value], error)

	Call(args ...Value) (Value, error)

	Len() (int64, error)

	Append(values ...Value) (Value, error)

	Value() any

	fmt.Stringer
}

const (
	IntType      = "int"
	RealType     = "real"
	TextType     = "text"
	BoolType     = "bool"
	ArrayType    = "array"
	ObjectType   = "object"
	FunctionType = "function"
	NullType     = "null"
)

type valueT interface {
	int64 |
		float64 |
		string |
		bool |
		[]Value |
		map[string]Value |
		func(args ...Value) (Value, error) |
		struct{} //nil
}

type value[T valueT] struct{ value T }

func (v value[T]) Value() any {
	switch v := any(v.value).(type) {
	case int64, float64, string, bool,
		func(...Value) (Value, error):
		return v

	case struct{}:
		return nil

	case []Value:
		sl := make([]any, 0, len(v))
		for _, i := range v {
			sl = append(sl, i.Value())
		}
		return sl

	case map[string]Value:
		m := make(map[string]any, len(v))
		for k, v := range v {
			m[k] = v.Value()
		}
		return m

	default:
		panic("неизвестный тип данных")
	}
}

func noAppendSupport(typ string) error {
	return fmt.Errorf("тип %s не поддерживает добавление новых элементов", typ)
}

func (v value[T]) Append(values ...Value) (Value, error) {
	arr, ok := any(v.value).([]Value)
	if !ok {
		return nil, noAppendSupport(v.Type())
	}

	newArr := make([]Value, len(arr))
	copy(newArr, arr)
	newArr = append(newArr, values...)

	return Array(newArr...), nil
}

func (v value[T]) String() string { return v.Text() }

//МЕТОДЫ КОНВЕРТАЦИИ:

func conversionError(from, to string) error {
	return fmt.Errorf("невозможно преобразовать %s в %s", from, to)
}

func (v value[T]) Int() (int64, error) {
	switch value := any(v.value).(type) {
	case int64:
		return value, nil
	case float64:
		return int64(value), nil
	case string:
		return textToInt([]rune(value)), nil
	case bool:
		if value {
			return 1, nil
		}
		return 0, nil
	case struct{}:
		return 0, nil
	default:
		return 0, conversionError(v.Type(), IntType)
	}
}

func (v value[T]) Real() (float64, error) {
	switch value := any(v.value).(type) {
	case int64:
		return float64(value), nil
	case float64:
		return value, nil
	case string:
		return textToReal([]rune(value)), nil
	case bool:
		if value {
			return 1, nil
		}
		return 0, nil
	case struct{}:
		return 0, nil
	default:
		return 0, conversionError(v.Type(), RealType)
	}
}

func (v value[T]) Text() string {
	switch value := any(v.value).(type) {
	case int64:
		return strconv.FormatInt(value, 10)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case string:
		return value
	case bool:
		return strconv.FormatBool(value)
	case struct{}:
		return "null"
	case []Value:
		strs := make([]string, 0, len(value))
		for _, v := range value {
			strs = append(strs, v.Text())
		}
		return fmt.Sprintf("[%s]", strings.Join(strs, ","))
	case map[string]Value:
		strs := make([]string, 0, len(value))
		for k, v := range value {
			strs = append(strs, fmt.Sprintf("%s:%s", k, v.Text()))
		}
		return fmt.Sprintf("{%s}", strings.Join(strs, ","))
	case func(...Value) (Value, error):
		//в дальнейшем это может быть изменено на что-то другое
		return FunctionType
	default:
		panic("неизвестный тип данных")
	}
}

func (v value[T]) Bool() (bool, error) {
	switch value := any(v.value).(type) {
	case int64:
		return value != 0, nil
	case float64:
		return value != 0, nil
	case string:
		return len(value) != 0, nil
	case bool:
		return value, nil
	case struct{}:
		return false, nil
	case []Value:
		return len(value) != 0, nil
	case map[string]Value:
		return len(value) != 0, nil
	default:
		return false, conversionError(v.Type(), BoolType)
	}
}

//МЕТОДЫ ПРОВЕРКИ ТИПА:

func (v value[T]) IsReal() bool {
	switch any(v.value).(type) {
	case float64:
		return true
	case string:
		return textIsReal([]rune(any(v.value).(string)))
	default:
		return false
	}
}

func (v value[T]) IsText() bool {
	_, ok := any(v.value).(string)
	return ok
}

//МЕТОДЫ РАБОТЫ С ЭЛЕМЕНТАМИ КОЛЛЕКЦИИ (ПОЛУЧЕНИЕ, МОДИФИКАЦИЯ):

func indexOutOfRange() error {
	return errors.New("индекс находится вне диапазона")
}

func noIndexSupport(typ string) error {
	return fmt.Errorf("тип %s не поддерживает доступ по индексу", typ)
}

func noSetIndexSupport(typ string) error {
	return fmt.Errorf("тип %s не поддерживает изменение элемента по индексу", typ)
}

func (v value[T]) ElByIndex(index Value) (Value, error) {
	switch value := any(v.value).(type) {
	case string:
		i, err := index.Int()
		if err != nil {
			return nil, err
		}
		runes := []rune(value)
		if i < 0 || int(i) >= len(runes) {
			return nil, indexOutOfRange()
		}
		return Text(string(runes[int(i)])), nil

	case []Value:
		i, err := index.Int()
		if err != nil {
			return nil, err
		}
		if i < 0 || int(i) >= len(value) {
			return nil, indexOutOfRange()
		}
		return value[int(i)], nil

	case map[string]Value:
		v, ok := value[index.Text()]
		if !ok {
			return Null(), nil
		}
		return v, nil

	default:
		return nil, noIndexSupport(v.Type())
	}
}

func (v value[T]) SetElByIndex(index, value Value) error {
	switch target := any(v.value).(type) {
	case []Value:
		i, err := index.Int()
		if err != nil {
			return err
		}
		if i < 0 || int(i) >= len(target) {
			return indexOutOfRange()
		}
		target[int(i)] = value
		return nil

	case map[string]Value:
		target[index.Text()] = value
		return nil

	default:
		return noSetIndexSupport(v.Type())
	}
}

//МЕТОДЫ ИТЕРАЦИИ:

func noIterSupport(typ string) error {
	return fmt.Errorf("тип %s не поддерживает итерацию", typ)
}

func noIterSupport2(typ string) error {
	return fmt.Errorf("тип %s не поддерживает итерацию на два элемента", typ)
}

func (v value[T]) Iter() (iter.Seq[Value], error) {
	switch target := any(v.value).(type) {
	case int64:
		return func(yield func(Value) bool) {
			for i := range target {
				if !yield(Int(i)) {
					return
				}
			}
		}, nil

	case float64:
		i, err := v.Int()
		if err != nil {
			return nil, err
		}
		return Int(i).Iter()

	case string, []Value:
		l, err := v.Len()
		if err != nil {
			return nil, err
		}
		return Int(l).Iter()

	case map[string]Value:
		return func(yield func(Value) bool) {
			for k := range target {
				if !yield(Text(k)) {
					return
				}
			}
		}, nil

	default:
		return nil, noIterSupport(v.Type())
	}
}

func (v value[T]) Iter2() (iter.Seq2[Value, Value], error) {
	switch any(v.value).(type) {
	case string, []Value, map[string]Value:
		return func(yield func(Value, Value) bool) {
			iter, err := v.Iter()
			if err != nil {
				panic("невозможная ошибка: " + err.Error())
			}

			for i := range iter {
				value, err := v.ElByIndex(i)
				if err != nil {
					panic("невозможная ошибка: " + err.Error())
				}

				if !yield(i, value) {
					return
				}
			}
		}, nil

	default:
		return nil, noIterSupport2(v.Type())
	}
}

//МЕТОДЫ РАБОТЫ С ФУНКЦИЯМИ:

func noCallSupport(typ string) error {
	return fmt.Errorf("тип %s не поддерживает вызовы", typ)
}

func (v value[T]) Call(args ...Value) (Value, error) {
	fun, ok := any(v.value).(func(...Value) (Value, error))
	if !ok {
		return nil, noCallSupport(v.Type())
	}
	return fun(args...)
}

//ПОЛУЧЕНИЕ ТИПА ЗНАЧЕНИЯ:

func (v value[T]) Type() string {
	switch any(v.value).(type) {
	case int64:
		return IntType
	case float64:
		return RealType
	case string:
		return TextType
	case bool:
		return BoolType
	case struct{}:
		return NullType
	case []Value:
		return ArrayType
	case map[string]Value:
		return ObjectType
	case func(...Value) (Value, error):
		return FunctionType
	default:
		panic("неизвестный тип данных")
	}
}

//ПОЛУЧЕНИЕ ДЛИНЫ ЗНАЧЕНИЯ:

func noLenSupport(typ string) error {
	return fmt.Errorf("тип %s не поддерживает получение длины", typ)
}

func (v value[T]) Len() (int64, error) {
	switch target := any(v.value).(type) {
	case string:
		return int64(len([]rune(target))), nil
	case []Value:
		return int64(len(target)), nil
	case map[string]Value:
		return int64(len(target)), nil
	default:
		return 0, noLenSupport(v.Type())
	}
}

//КОНСТРУКТОРЫ:

func Int(v int64) Value    { return value[int64]{v} }
func Real(v float64) Value { return value[float64]{v} }
func Text(v string) Value  { return value[string]{v} }
func Bool(v bool) Value    { return value[bool]{v} }

func Array(v ...Value) Value {
	if v == nil {
		v = make([]Value, 0)
	}
	return value[[]Value]{v}
}

type KV struct{ Key, Value Value }

func Object(v ...KV) Value {
	m := make(map[string]Value, len(v))

	for _, kv := range v {
		m[kv.Key.Text()] = kv.Value
	}

	return value[map[string]Value]{m}
}

func Function(v func(args ...Value) (Value, error)) Value {
	return value[func(...Value) (Value, error)]{v}
}

func Null() Value { return value[struct{}]{} }

//ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ ДЛЯ РАБОТЫ СО СТРОКАМИ:

func skip(sl []rune, index int, fun func(rune) bool) int {
	for index < len(sl) && fun(sl[index]) {
		index++
	}
	return index
}

func skipDigits(sl []rune, index int) int {
	return skip(sl, index, unicode.IsDigit)
}

func skipSpaces(sl []rune, index int) int {
	return skip(sl, index, unicode.IsSpace)
}

func textIsReal(sl []rune) bool {
	index := skipSpaces(sl, 0)
	if index == len(sl) {
		return false
	}

	if sl[index] == '-' || sl[index] == '+' {
		index++
		if index == len(sl) {
			return false
		}
	}

	if unicode.IsDigit(sl[index]) {
		index = skipDigits(sl, index)
		return index != len(sl) && sl[index] == '.'
	}

	return sl[index] == '.' && index+1 != len(sl) && unicode.IsDigit(sl[index+1])
}

func textToInt(sl []rune) int64 {
	start := skipSpaces(sl, 0)
	if start == len(sl) {
		return 0
	}

	end := start
	if sl[end] == '-' || sl[end] == '+' {
		end++
		if end == len(sl) {
			return 0
		}
	}

	end = skipDigits(sl, end)

	num, _ := strconv.ParseInt(string(sl[start:end]), 10, 64)
	return num
}

func textToReal(sl []rune) float64 {
	start := skipSpaces(sl, 0)
	if start == len(sl) {
		return 0
	}

	end := start
	if sl[end] == '-' || sl[end] == '+' {
		end++
		if end == len(sl) {
			return 0
		}
	}

	end = skipDigits(sl, end)

	if end != len(sl) && sl[end] == '.' {
		end = skipDigits(sl, end+1)
	}

	num, _ := strconv.ParseFloat(string(sl[start:end]), 64)
	return num
}

func Of(v any) (Value, error) {
	if v == nil {
		return Null(), nil
	}

	switch val := reflect.ValueOf(v); val.Kind() {
	case
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return Int(val.Int()), nil

	case reflect.Float64:
		return Real(val.Float()), nil

	case reflect.Bool:
		return Bool(val.Bool()), nil

	case reflect.String:
		return Text(val.String()), nil

	case reflect.Slice:
		values := make([]Value, 0, val.Len())

		for i := range val.Len() {
			v, err := Of(val.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			values = append(values, v)
		}

		return Array(values...), nil

	case reflect.Map:
		values := make([]KV, 0, val.Len())

		for iter := val.MapRange(); iter.Next(); {
			k, err := Of(iter.Key().Interface())
			if err != nil {
				return nil, err
			}

			v, err := Of(iter.Value().Interface())
			if err != nil {
				return nil, err
			}

			values = append(values, KV{Key: k, Value: v})
		}

		return Object(values...), nil

	default:
		return nil, conversionError(val.Kind().String(), "Value")
	}
}
