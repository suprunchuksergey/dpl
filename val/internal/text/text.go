package text

import (
	"fmt"
	"strconv"
	"unicode"
)

type Text string

func (t Text) String() string { return fmt.Sprintf("text %s", string(t)) }

func (t Text) Int() int64 {
	var s int
	if t.skipSpaces(&s) {
		return 0
	}

	e := s
	if t[e] == '+' || t[e] == '-' {
		e++
	}

	t.skipNumbers(&e)

	n, _ := strconv.ParseInt(string(t[s:e]), 10, 64)
	return n
}

func (t Text) Real() float64 {
	var s int
	if t.skipSpaces(&s) {
		return 0
	}

	e := s
	if t[e] == '+' || t[e] == '-' {
		e++
	}

	if t.skipNumbers(&e) == false && t[e] == '.' {
		e++
		t.skipNumbers(&e)
	}

	n, _ := strconv.ParseFloat(string(t[s:e]), 64)
	return n
}

func (t Text) Text() string { return string(t) }

func (t Text) Bool() bool { return len(t) != 0 }

func (t Text) CanInt() bool {
	var i int
	if t.skipSpaces(&i) {
		return false
	}

	if t[i] == '+' || t[i] == '-' {
		i++
	}

	return i < len(t) && unicode.IsDigit(rune(t[i]))
}

func (t Text) CanReal() bool {
	var i int
	if t.skipSpaces(&i) {
		return false
	}

	if t[i] == '+' || t[i] == '-' {
		i++
		if i == len(t) {
			return false
		}
	}

	if unicode.IsDigit(rune(t[i])) {
		return t.skipNumbers(&i) == false && t[i] == '.'
	}

	if t[i] == '.' {
		i++
		return i < len(t) && unicode.IsDigit(rune(t[i]))
	}

	return false
}

func New(val string) Text { return Text(val) }

func (t Text) skip(i *int, c func(rune) bool) bool {
	for *i < len(t) && c(rune(t[*i])) {
		*i++
	}
	return *i >= len(t)
}

func (t Text) skipSpaces(i *int) bool { return t.skip(i, unicode.IsSpace) }

func (t Text) skipNumbers(i *int) bool { return t.skip(i, unicode.IsDigit) }
