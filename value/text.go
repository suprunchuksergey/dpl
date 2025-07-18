package value

import (
	"strconv"
	"unicode"
)

type text struct{ v string }

func newText(v string) text { return text{v} }

func (t text) Int() int64 {
	var start int
	if t.skipSpaces(&start) {
		return 0
	}

	end := start
	if t.v[start] == '+' || t.v[start] == '-' {
		end++
	}

	t.skipNumbers(&end)

	n, _ := strconv.ParseInt(t.v[start:end], 10, 64)
	return n
}

func (t text) Real() float64 {
	var start int
	if t.skipSpaces(&start) {
		return 0
	}

	end := start
	if t.v[start] == '+' || t.v[start] == '-' {
		end++
	}

	if !t.skipNumbers(&end) && t.v[end] == '.' {
		end++
		t.skipNumbers(&end)
	}

	n, _ := strconv.ParseFloat(t.v[start:end], 64)
	return n
}

func (t text) Text() string { return t.v }

func (t text) IsInt() bool {
	var index int
	if t.skipSpaces(&index) {
		return false
	}

	if t.v[index] == '+' || t.v[index] == '-' {
		index++
	}

	return index < len(t.v) && unicode.IsDigit(rune(t.v[index]))
}

func (t text) IsReal() bool {
	var index int
	if t.skipSpaces(&index) {
		return false
	}

	if t.v[index] == '+' || t.v[index] == '-' {
		index++
	}

	if index == len(t.v) {
		return false
	}

	if unicode.IsDigit(rune(t.v[index])) {
		return !t.skipNumbers(&index) && t.v[index] == '.'
	} else if t.v[index] == '.' {
		index++
		return index < len(t.v) &&
			unicode.IsDigit(rune(t.v[index]))
	}

	return false
}

func (t text) IsText() bool { return true }

var _ Value = text{}

func (t text) skipSpaces(index *int) bool {
	for *index < len(t.v) && unicode.IsSpace(rune(t.v[*index])) {
		*index++
	}
	return *index == len(t.v)
}

func (t text) skipNumbers(index *int) bool {
	for *index < len(t.v) && unicode.IsDigit(rune(t.v[*index])) {
		*index++
	}
	return *index == len(t.v)
}
