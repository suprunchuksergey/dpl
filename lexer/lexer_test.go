package lexer

import (
	"github.com/suprunchuksergey/dpl/pos"
	"github.com/suprunchuksergey/dpl/token"
	"testing"
)

func Test_next(t *testing.T) {
	const baseData = "lorem ipsum dolor sit amet,\nconsectetur adipiscing elit\nlorem"

	tests := []struct {
		n    int //количество вызовов методов next
		data string
		//ожидаемое состояние после всех вызовов, без src и index (index должен быть равен n)
		expected lexer
	}{
		{expected: lexer{char: 0, pos: pos.New()}},
		{
			n:        28,
			expected: lexer{char: 0, pos: pos.New()},
		},
		{
			data: baseData,
			expected: lexer{
				char: 'l',
				pos:  pos.New(),
			},
		},
		{
			n:    2,
			data: baseData,
			expected: lexer{
				char: 'r',
				pos:  pos.NewWithStart(1, 3),
			},
		},
		{
			n:    6,
			data: baseData,
			expected: lexer{
				char: 'i',
				pos:  pos.NewWithStart(1, 7),
			},
		},
		{
			n:    12,
			data: baseData,
			expected: lexer{
				char: 'd',
				pos:  pos.NewWithStart(1, 13),
			},
		},
		{
			n:    28,
			data: baseData,
			expected: lexer{
				char: 'c',
				pos:  pos.NewWithStart(2, 1),
			},
		},
		{
			n:    40,
			data: baseData,
			expected: lexer{
				char: 'a',
				pos:  pos.NewWithStart(2, 13),
			},
		},
		{
			n:    56,
			data: baseData,
			expected: lexer{
				char: 'l',
				pos:  pos.NewWithStart(3, 1),
			},
		},
		{
			n:    58,
			data: baseData,
			expected: lexer{
				char: 'r',
				pos:  pos.NewWithStart(3, 3),
			},
		},
		{
			n:    61,
			data: baseData,
			expected: lexer{
				char: 0,
				pos:  pos.NewWithStart(3, 6),
			},
		},
		{
			n:    70,
			data: baseData,
			expected: lexer{
				char: 0,
				pos:  pos.NewWithStart(3, 6),
			},
		},
		{
			n:    20,
			data: baseData,
			expected: lexer{
				char: 't',
				pos:  pos.NewWithStart(1, 21),
			},
		},
		{
			n:    26,
			data: baseData,
			expected: lexer{
				char: ',',
				pos:  pos.NewWithStart(1, 27),
			},
		},
		{
			n:    3,
			data: "ожидалось",
			expected: lexer{
				char: 'д',
				pos:  pos.NewWithStart(1, 4),
			},
		},
		{
			n:    4,
			data: "\nожидалось",
			expected: lexer{
				char: 'д',
				pos:  pos.NewWithStart(2, 4),
			},
		},
	}

	for i, test := range tests {
		lex := newLexer(test.data)
		for range test.n {
			lex.next()
		}

		index := test.n
		if test.n > len([]rune(test.data)) {
			index = len([]rune(test.data))
		}

		if lex.index != index {
			t.Errorf("%d: index: ожидалось %d, получено %d",
				i, index, lex.index)
		}

		if lex.char != test.expected.char {
			t.Errorf("%d: char: ожидалось %c, получено %c",
				i, test.expected.char, lex.char)
		}

		if lex.pos.String() != test.expected.pos.String() {
			t.Errorf("%d: pos: ожидалось %s, получено %s",
				i, test.expected.pos, lex.pos)
		}
	}
}

func Test_skipSpaces(t *testing.T) {
	tests := []struct {
		data string
		//ожидаемое состояние, без src
		expected lexer
	}{
		{
			data: "f",
			expected: lexer{
				index: 0,
				char:  'f',
				pos:   pos.NewWithStart(1, 1),
			},
		},
		{
			data: "	 f",
			expected: lexer{
				index: 2,
				char:  'f',
				pos:   pos.NewWithStart(1, 3),
			},
		},
		{
			data: "			 f",
			expected: lexer{
				index: 4,
				char:  'f',
				pos:   pos.NewWithStart(1, 5),
			},
		},
		{
			data: "\n\nf",
			expected: lexer{
				index: 2,
				char:  'f',
				pos:   pos.NewWithStart(3, 1),
			},
		},
		{
			data: "\n\n	 f",
			expected: lexer{
				index: 4,
				char:  'f',
				pos:   pos.NewWithStart(3, 3),
			},
		},
	}

	for i, test := range tests {
		lex := newLexer(test.data)
		lex.skipSpaces()

		if lex.index != test.expected.index {
			t.Errorf("%d: index: ожидалось %d, получено %d",
				i, test.expected.index, lex.index)
		}

		if lex.char != test.expected.char {
			t.Errorf("%d: char: ожидалось %c, получено %c",
				i, test.expected.char, lex.char)
		}

		if lex.pos.String() != test.expected.pos.String() {
			t.Errorf("%d: pos: ожидалось %s, получено %s",
				i, test.expected.pos, lex.pos)
		}
	}
}

func Test_readNums(t *testing.T) {
	tests := []struct {
		data string
		//ожидаемое состояние, без src
		state lexer
		//ожидаемый результат вызова
		nums string
	}{
		{
			data: "",
			state: lexer{
				index: 0,
				pos:   pos.NewWithStart(1, 1),
			},
			nums: "",
		},
		{
			data: "data",
			state: lexer{
				index: 0,
				char:  'd',
				pos:   pos.NewWithStart(1, 1),
			},
			nums: "",
		},
		{
			data: "235",
			state: lexer{
				index: 3,
				pos:   pos.NewWithStart(1, 4),
			},
			nums: "235",
		},
		{
			data: "235data",
			state: lexer{
				index: 3,
				char:  'd',
				pos:   pos.NewWithStart(1, 4),
			},
			nums: "235",
		},
	}

	for i, test := range tests {
		lex := newLexer(test.data)
		nums := lex.readNums()

		if nums != test.nums {
			t.Errorf("%d: nums: ожидалось %s, получено %s",
				i, test.nums, nums)
		}

		if lex.index != test.state.index {
			t.Errorf("%d: index: ожидалось %d, получено %d",
				i, test.state.index, lex.index)
		}

		if lex.char != test.state.char {
			t.Errorf("%d: char: ожидалось %c, получено %c",
				i, test.state.char, lex.char)
		}

		if lex.pos.String() != test.state.pos.String() {
			t.Errorf("%d: pos: ожидалось %s, получено %s",
				i, test.state.pos, lex.pos)
		}
	}
}

func Test_readTail(t *testing.T) {
	tests := []struct {
		data string
		//ожидаемое состояние, без src
		state lexer
		//ожидаемый результат вызова
		tail string
	}{
		{
			data: "",
			state: lexer{
				index: 0,
				pos:   pos.NewWithStart(1, 1),
			},
			tail: "",
		},
		{
			data: "data",
			state: lexer{
				index: 0,
				char:  'd',
				pos:   pos.NewWithStart(1, 1),
			},
			tail: "",
		},
		{
			data: "235",
			state: lexer{
				index: 0,
				char:  '2',
				pos:   pos.NewWithStart(1, 1),
			},
			tail: "",
		},
		{
			data: ".235",
			state: lexer{
				index: 4,
				pos:   pos.NewWithStart(1, 5),
			},
			tail: ".235",
		},
		{
			data: ".",
			state: lexer{
				index: 1,
				pos:   pos.NewWithStart(1, 2),
			},
			tail: ".0",
		},
	}

	for i, test := range tests {
		lex := newLexer(test.data)
		tail := lex.readTail()

		if tail != test.tail {
			t.Errorf("%d: tail: ожидалось %s, получено %s",
				i, test.tail, tail)
		}

		if lex.index != test.state.index {
			t.Errorf("%d: index: ожидалось %d, получено %d",
				i, test.state.index, lex.index)
		}

		if lex.char != test.state.char {
			t.Errorf("%d: char: ожидалось %c, получено %c",
				i, test.state.char, lex.char)
		}

		if lex.pos.String() != test.state.pos.String() {
			t.Errorf("%d: pos: ожидалось %s, получено %s",
				i, test.state.pos, lex.pos)
		}
	}
}

func Test_Next(t *testing.T) {
	//все тесты должны пройти без ошибок
	tests := []struct {
		data string
		//ожидаемое состояние, после i-го вызова, без src
		expected []lexer
	}{
		{
			data: "+",
			expected: []lexer{
				{
					index: 1,
					pos:   pos.NewWithStart(1, 2),
					tok:   token.New(token.Add, pos.NewWithStart(1, 1)),
				},
			},
		},
		{
			data: "=",
			expected: []lexer{
				{
					index: 1,
					pos:   pos.NewWithStart(1, 2),
					tok:   token.New(token.Assign, pos.NewWithStart(1, 1)),
				},
			},
		},
		{
			data: "for",
			expected: []lexer{
				{
					index: 3,
					pos:   pos.NewWithStart(1, 4),
					tok:   token.New(token.For, pos.NewWithStart(1, 1)),
				},
			},
		},
		{
			data: "in",
			expected: []lexer{
				{
					index: 2,
					pos:   pos.NewWithStart(1, 3),
					tok:   token.New(token.In, pos.NewWithStart(1, 1)),
				},
			},
		},
		{
			data: "if",
			expected: []lexer{
				{
					index: 2,
					pos:   pos.NewWithStart(1, 3),
					tok:   token.New(token.If, pos.NewWithStart(1, 1)),
				},
			},
		},
		{
			data: "elif",
			expected: []lexer{
				{
					index: 4,
					pos:   pos.NewWithStart(1, 5),
					tok:   token.New(token.Elif, pos.NewWithStart(1, 1)),
				},
			},
		},
		{
			data: "else",
			expected: []lexer{
				{
					index: 4,
					pos:   pos.NewWithStart(1, 5),
					tok:   token.New(token.Else, pos.NewWithStart(1, 1)),
				},
			},
		},
		{
			data: "-",
			expected: []lexer{
				{
					index: 1,
					pos:   pos.NewWithStart(1, 2),
					tok:   token.New(token.Sub, pos.NewWithStart(1, 1)),
				},
			},
		},
		{
			data: "*",
			expected: []lexer{
				{
					index: 1,
					pos:   pos.NewWithStart(1, 2),
					tok:   token.New(token.Mul, pos.NewWithStart(1, 1)),
				},
			},
		},
		{
			data: "/",
			expected: []lexer{
				{
					index: 1,
					pos:   pos.NewWithStart(1, 2),
					tok:   token.New(token.Div, pos.NewWithStart(1, 1)),
				},
			},
		},
		{
			data: "%",
			expected: []lexer{
				{
					index: 1,
					pos:   pos.NewWithStart(1, 2),
					tok:   token.New(token.Rem, pos.NewWithStart(1, 1)),
				},
			},
		},
		{
			data: "==",
			expected: []lexer{
				{
					index: 2,
					pos:   pos.NewWithStart(1, 3),
					tok:   token.New(token.Eq, pos.NewWithStart(1, 1)),
				},
			},
		},
		{
			data: "!=",
			expected: []lexer{
				{
					index: 2,
					pos:   pos.NewWithStart(1, 3),
					tok:   token.New(token.Neq, pos.NewWithStart(1, 1)),
				},
			},
		},
		{
			data: "<>",
			expected: []lexer{
				{
					index: 2,
					pos:   pos.NewWithStart(1, 3),
					tok:   token.New(token.Neq, pos.NewWithStart(1, 1)),
				},
			},
		},
		{
			data: "<",
			expected: []lexer{
				{
					index: 1,
					pos:   pos.NewWithStart(1, 2),
					tok:   token.New(token.Lt, pos.NewWithStart(1, 1)),
				},
			},
		},
		{
			data: "<=",
			expected: []lexer{
				{
					index: 2,
					pos:   pos.NewWithStart(1, 3),
					tok:   token.New(token.Lte, pos.NewWithStart(1, 1)),
				},
			},
		},
		{
			data: ">",
			expected: []lexer{
				{
					index: 1,
					pos:   pos.NewWithStart(1, 2),
					tok:   token.New(token.Gt, pos.NewWithStart(1, 1)),
				},
			},
		},
		{
			data: ">=",
			expected: []lexer{
				{
					index: 2,
					pos:   pos.NewWithStart(1, 3),
					tok:   token.New(token.Gte, pos.NewWithStart(1, 1)),
				},
			},
		},
		{
			data: "||",
			expected: []lexer{
				{
					index: 2,
					pos:   pos.NewWithStart(1, 3),
					tok:   token.New(token.Concat, pos.NewWithStart(1, 1)),
				},
			},
		},
		{
			data: "	 \n	+",
			expected: []lexer{
				{
					index: 5,
					pos:   pos.NewWithStart(2, 3),
					tok:   token.New(token.Add, pos.NewWithStart(2, 2)),
				},
			},
		},
		{
			data: "	 \n	+5",
			expected: []lexer{
				{
					index: 5,
					char:  '5',
					pos:   pos.NewWithStart(2, 3),
					tok:   token.New(token.Add, pos.NewWithStart(2, 2)),
				},
			},
		},
		{
			data: "	 \n	||",
			expected: []lexer{
				{
					index: 6,
					pos:   pos.NewWithStart(2, 4),
					tok:   token.New(token.Concat, pos.NewWithStart(2, 2)),
				},
			},
		},
		{
			data: "	 \n	|||",
			expected: []lexer{
				{
					index: 6,
					char:  '|',
					pos:   pos.NewWithStart(2, 4),
					tok:   token.New(token.Concat, pos.NewWithStart(2, 2)),
				},
			},
		},
		{
			data: "***",
			expected: []lexer{
				{
					index: 1,
					char:  '*',
					pos:   pos.NewWithStart(1, 2),
					tok:   token.New(token.Mul, pos.NewWithStart(1, 1)),
				},
				{
					index: 2,
					char:  '*',
					pos:   pos.NewWithStart(1, 3),
					tok:   token.New(token.Mul, pos.NewWithStart(1, 2)),
				},
				{
					index: 3,
					pos:   pos.NewWithStart(1, 4),
					tok:   token.New(token.Mul, pos.NewWithStart(1, 3)),
				},
			},
		},
		{
			data: "	*\n+	\n *",
			expected: []lexer{
				{
					index: 2,
					char:  '\n',
					pos:   pos.NewWithStart(1, 3),
					tok:   token.New(token.Mul, pos.NewWithStart(1, 2)),
				},
				{
					index: 4,
					char:  '\t',
					pos:   pos.NewWithStart(2, 2),
					tok:   token.New(token.Add, pos.NewWithStart(2, 1)),
				},
				{
					index: 8,
					pos:   pos.NewWithStart(3, 3),
					tok:   token.New(token.Mul, pos.NewWithStart(3, 2)),
				},
			},
		},
		{
			data: "====",
			expected: []lexer{
				{
					index: 2,
					char:  '=',
					pos:   pos.NewWithStart(1, 3),
					tok:   token.New(token.Eq, pos.NewWithStart(1, 1)),
				},
				{
					index: 4,
					pos:   pos.NewWithStart(1, 5),
					tok:   token.New(token.Eq, pos.NewWithStart(1, 3)),
				},
			},
		},
		{
			data: "0",
			expected: []lexer{
				{
					index: 1,
					pos:   pos.NewWithStart(1, 2),
					tok: token.NewWithValue(
						token.Int,
						pos.NewWithStart(1, 1),
						"0",
					),
				},
			},
		},
		{
			data: "01",
			expected: []lexer{
				{
					index: 1,
					char:  '1',
					pos:   pos.NewWithStart(1, 2),
					tok: token.NewWithValue(
						token.Int,
						pos.NewWithStart(1, 1),
						"0",
					),
				},
			},
		},
		{
			data: "0.",
			expected: []lexer{
				{
					index: 2,
					pos:   pos.NewWithStart(1, 3),
					tok: token.NewWithValue(
						token.Real,
						pos.NewWithStart(1, 1),
						"0.0",
					),
				},
			},
		},
		{
			data: "0.l",
			expected: []lexer{
				{
					index: 2,
					char:  'l',
					pos:   pos.NewWithStart(1, 3),
					tok: token.NewWithValue(
						token.Real,
						pos.NewWithStart(1, 1),
						"0.0",
					),
				},
			},
		},
		{
			data: "0.235",
			expected: []lexer{
				{
					index: 5,
					pos:   pos.NewWithStart(1, 6),
					tok: token.NewWithValue(
						token.Real,
						pos.NewWithStart(1, 1),
						"0.235",
					),
				},
			},
		},
		{
			data: "0.0",
			expected: []lexer{
				{
					index: 3,
					pos:   pos.NewWithStart(1, 4),
					tok: token.NewWithValue(
						token.Real,
						pos.NewWithStart(1, 1),
						"0.0",
					),
				},
			},
		},
		{
			data: ".0",
			expected: []lexer{
				{
					index: 2,
					pos:   pos.NewWithStart(1, 3),
					tok: token.NewWithValue(
						token.Real,
						pos.NewWithStart(1, 1),
						"0.0",
					),
				},
			},
		},
		{
			data: ".235",
			expected: []lexer{
				{
					index: 4,
					pos:   pos.NewWithStart(1, 5),
					tok: token.NewWithValue(
						token.Real,
						pos.NewWithStart(1, 1),
						"0.235",
					),
				},
			},
		},
		{
			data: "235",
			expected: []lexer{
				{
					index: 3,
					pos:   pos.NewWithStart(1, 4),
					tok: token.NewWithValue(
						token.Int,
						pos.NewWithStart(1, 1),
						"235",
					),
				},
			},
		},
		{
			data: "235.",
			expected: []lexer{
				{
					index: 4,
					pos:   pos.NewWithStart(1, 5),
					tok: token.NewWithValue(
						token.Real,
						pos.NewWithStart(1, 1),
						"235.0",
					),
				},
			},
		},
		{
			data: "35.5",
			expected: []lexer{
				{
					index: 4,
					pos:   pos.NewWithStart(1, 5),
					tok: token.NewWithValue(
						token.Real,
						pos.NewWithStart(1, 1),
						"35.5",
					),
				},
			},
		},
		{
			data: "5. + -.5 * 5",
			expected: []lexer{
				{
					index: 2,
					char:  ' ',
					pos:   pos.NewWithStart(1, 3),
					tok: token.NewWithValue(
						token.Real,
						pos.NewWithStart(1, 1),
						"5.0",
					),
				},
				{
					index: 4,
					char:  ' ',
					pos:   pos.NewWithStart(1, 5),
					tok:   token.New(token.Add, pos.NewWithStart(1, 4)),
				},
				{
					index: 6,
					char:  '.',
					pos:   pos.NewWithStart(1, 7),
					tok:   token.New(token.Sub, pos.NewWithStart(1, 6)),
				},
				{
					index: 8,
					char:  ' ',
					pos:   pos.NewWithStart(1, 9),
					tok: token.NewWithValue(
						token.Real,
						pos.NewWithStart(1, 7),
						"0.5",
					),
				},
				{
					index: 10,
					char:  ' ',
					pos:   pos.NewWithStart(1, 11),
					tok:   token.New(token.Mul, pos.NewWithStart(1, 10)),
				},
				{
					index: 12,
					pos:   pos.NewWithStart(1, 13),
					tok: token.NewWithValue(
						token.Int,
						pos.NewWithStart(1, 12),
						"5",
					),
				},
			},
		},
		{
			data: "5.5 == 5",
			expected: []lexer{
				{
					index: 3,
					char:  ' ',
					pos:   pos.NewWithStart(1, 4),
					tok: token.NewWithValue(
						token.Real,
						pos.NewWithStart(1, 1),
						"5.5",
					),
				},
				{
					index: 6,
					char:  ' ',
					pos:   pos.NewWithStart(1, 7),
					tok:   token.New(token.Eq, pos.NewWithStart(1, 5)),
				},
				{
					index: 8,
					pos:   pos.NewWithStart(1, 9),
					tok: token.NewWithValue(
						token.Int,
						pos.NewWithStart(1, 8),
						"5",
					),
				},
			},
		},
		{
			data: ".5 5 5. 5.5",
			expected: []lexer{
				{
					index: 2,
					char:  ' ',
					pos:   pos.NewWithStart(1, 3),
					tok: token.NewWithValue(
						token.Real,
						pos.NewWithStart(1, 1),
						"0.5",
					),
				},
				{
					index: 4,
					char:  ' ',
					pos:   pos.NewWithStart(1, 5),
					tok: token.NewWithValue(
						token.Int,
						pos.NewWithStart(1, 4),
						"5",
					),
				},
				{
					index: 7,
					char:  ' ',
					pos:   pos.NewWithStart(1, 8),
					tok: token.NewWithValue(
						token.Real,
						pos.NewWithStart(1, 6),
						"5.0",
					),
				},
				{
					index: 11,
					pos:   pos.NewWithStart(1, 12),
					tok: token.NewWithValue(
						token.Real,
						pos.NewWithStart(1, 9),
						"5.5",
					),
				},
			},
		},
		{
			//будет прочитано как два числа, так как начальные нули не допускаются,
			//если только число не 0
			data: "05.5",
			expected: []lexer{
				{
					index: 1,
					char:  '5',
					pos:   pos.NewWithStart(1, 2),
					tok: token.NewWithValue(
						token.Int,
						pos.NewWithStart(1, 1),
						"0",
					),
				},
				{
					index: 4,
					pos:   pos.NewWithStart(1, 5),
					tok: token.NewWithValue(
						token.Real,
						pos.NewWithStart(1, 2),
						"5.5",
					),
				},
			},
		},
		{
			data: "'lexer'",
			expected: []lexer{
				{
					index: 7,
					pos:   pos.NewWithStart(1, 8),
					tok: token.NewWithValue(
						token.Text,
						pos.NewWithStart(1, 1),
						"lexer",
					),
				},
			},
		},
		{
			data: "'l''exer'",
			expected: []lexer{
				{
					index: 9,
					pos:   pos.NewWithStart(1, 10),
					tok: token.NewWithValue(
						token.Text,
						pos.NewWithStart(1, 1),
						"l'exer",
					),
				},
			},
		},
		{
			data: "'l''''exer'",
			expected: []lexer{
				{
					index: 11,
					pos:   pos.NewWithStart(1, 12),
					tok: token.NewWithValue(
						token.Text,
						pos.NewWithStart(1, 1),
						"l''exer",
					),
				},
			},
		},
		{
			data: "''",
			expected: []lexer{
				{
					index: 2,
					pos:   pos.NewWithStart(1, 3),
					tok: token.NewWithValue(
						token.Text,
						pos.NewWithStart(1, 1),
						"",
					),
				},
			},
		},
		{
			data: "''''''",
			expected: []lexer{
				{
					index: 6,
					pos:   pos.NewWithStart(1, 7),
					tok: token.NewWithValue(
						token.Text,
						pos.NewWithStart(1, 1),
						"''",
					),
				},
			},
		},
		{
			data: `'l\nexer'`,
			expected: []lexer{
				{
					index: 9,
					pos:   pos.NewWithStart(1, 10),
					tok: token.NewWithValue(
						token.Text,
						pos.NewWithStart(1, 1),
						`l\nexer`,
					),
				},
			},
		},
		{
			data: "235+'l''exer'",
			expected: []lexer{
				{
					index: 3,
					char:  '+',
					pos:   pos.NewWithStart(1, 4),
					tok: token.NewWithValue(
						token.Int,
						pos.NewWithStart(1, 1),
						"235",
					),
				},
				{
					index: 4,
					char:  '\'',
					pos:   pos.NewWithStart(1, 5),
					tok:   token.New(token.Add, pos.NewWithStart(1, 4)),
				},
				{
					index: 13,
					pos:   pos.NewWithStart(1, 14),
					tok: token.NewWithValue(
						token.Text,
						pos.NewWithStart(1, 5),
						"l'exer",
					),
				},
			},
		},
		{
			data: `(5)`,
			expected: []lexer{
				{
					index: 1,
					char:  '5',
					pos:   pos.NewWithStart(1, 2),
					tok:   token.New(token.LParen, pos.NewWithStart(1, 1)),
				},
				{
					index: 2,
					char:  ')',
					pos:   pos.NewWithStart(1, 3),
					tok: token.NewWithValue(
						token.Int,
						pos.NewWithStart(1, 2),
						"5",
					),
				},
				{
					index: 3,
					pos:   pos.NewWithStart(1, 4),
					tok:   token.New(token.RParen, pos.NewWithStart(1, 3)),
				},
			},
		},
		{
			data: "PI + -.5 * i",
			expected: []lexer{
				{
					index: 2,
					char:  ' ',
					pos:   pos.NewWithStart(1, 3),
					tok: token.NewWithValue(
						token.Ident,
						pos.NewWithStart(1, 1),
						"pi",
					),
				},
				{
					index: 4,
					char:  ' ',
					pos:   pos.NewWithStart(1, 5),
					tok:   token.New(token.Add, pos.NewWithStart(1, 4)),
				},
				{
					index: 6,
					char:  '.',
					pos:   pos.NewWithStart(1, 7),
					tok:   token.New(token.Sub, pos.NewWithStart(1, 6)),
				},
				{
					index: 8,
					char:  ' ',
					pos:   pos.NewWithStart(1, 9),
					tok: token.NewWithValue(
						token.Real,
						pos.NewWithStart(1, 7),
						"0.5",
					),
				},
				{
					index: 10,
					char:  ' ',
					pos:   pos.NewWithStart(1, 11),
					tok:   token.New(token.Mul, pos.NewWithStart(1, 10)),
				},
				{
					index: 12,
					pos:   pos.NewWithStart(1, 13),
					tok: token.NewWithValue(
						token.Ident,
						pos.NewWithStart(1, 12),
						"i",
					),
				},
			},
		},
		{
			data: "and or not true false null lexer",
			expected: []lexer{
				{
					index: 3,
					char:  ' ',
					pos:   pos.NewWithStart(1, 4),
					tok:   token.New(token.And, pos.NewWithStart(1, 1)),
				},
				{
					index: 6,
					char:  ' ',
					pos:   pos.NewWithStart(1, 7),
					tok:   token.New(token.Or, pos.NewWithStart(1, 5)),
				},
				{
					index: 10,
					char:  ' ',
					pos:   pos.NewWithStart(1, 11),
					tok:   token.New(token.Not, pos.NewWithStart(1, 8)),
				},
				{
					index: 15,
					char:  ' ',
					pos:   pos.NewWithStart(1, 16),
					tok:   token.New(token.True, pos.NewWithStart(1, 12)),
				},
				{
					index: 21,
					char:  ' ',
					pos:   pos.NewWithStart(1, 22),
					tok:   token.New(token.False, pos.NewWithStart(1, 17)),
				},
				{
					index: 26,
					char:  ' ',
					pos:   pos.NewWithStart(1, 27),
					tok:   token.New(token.Null, pos.NewWithStart(1, 23)),
				},
				{
					index: 32,
					pos:   pos.NewWithStart(1, 33),
					tok: token.NewWithValue(
						token.Ident,
						pos.NewWithStart(1, 28),
						"lexer",
					),
				},
			},
		},
		{
			data: "",
			expected: []lexer{
				{
					index: 0,
					pos:   pos.NewWithStart(1, 1),
					tok:   token.New(token.EOF, pos.NewWithStart(1, 1)),
				},
				{
					index: 0,
					pos:   pos.NewWithStart(1, 1),
					tok:   token.New(token.EOF, pos.NewWithStart(1, 1)),
				},
				{
					index: 0,
					pos:   pos.NewWithStart(1, 1),
					tok:   token.New(token.EOF, pos.NewWithStart(1, 1)),
				},
			},
		},
		{
			data: "{}",
			expected: []lexer{
				{
					index: 1,
					char:  '}',
					pos:   pos.NewWithStart(1, 2),
					tok:   token.New(token.LBrace, pos.NewWithStart(1, 1)),
				},
				{
					index: 2,
					pos:   pos.NewWithStart(1, 3),
					tok:   token.New(token.RBrace, pos.NewWithStart(1, 2)),
				},
			},
		},
		{
			data: "[]",
			expected: []lexer{
				{
					index: 1,
					char:  ']',
					pos:   pos.NewWithStart(1, 2),
					tok:   token.New(token.LBrack, pos.NewWithStart(1, 1)),
				},
				{
					index: 2,
					pos:   pos.NewWithStart(1, 3),
					tok:   token.New(token.RBrack, pos.NewWithStart(1, 2)),
				},
			},
		},
		{
			data: ",:",
			expected: []lexer{
				{
					index: 1,
					char:  ':',
					pos:   pos.NewWithStart(1, 2),
					tok:   token.New(token.Comma, pos.NewWithStart(1, 1)),
				},
				{
					index: 2,
					pos:   pos.NewWithStart(1, 3),
					tok:   token.New(token.Colon, pos.NewWithStart(1, 2)),
				},
			},
		},
		{
			data: ";",
			expected: []lexer{
				{
					index: 1,
					pos:   pos.NewWithStart(1, 2),
					tok:   token.New(token.Semicolon, pos.NewWithStart(1, 1)),
				},
			}},
	}

	for i, test := range tests {
		lex := newLexer(test.data)

		for j, state := range test.expected {
			err := lex.Next()
			if err != nil {
				t.Errorf("%d %d: непредвиденная ошибка: %s", i, j, err)
				continue
			}

			if lex.index != state.index {
				t.Errorf("тест %d, токен %d: index: ожидалось %d, получено %d",
					i, j, state.index, lex.index)
			}

			if lex.char != state.char {
				t.Errorf("тест %d, токен %d: char: ожидалось %c, получено %c",
					i, j, state.char, lex.char)
			}

			if lex.pos.String() != state.pos.String() {
				t.Errorf("тест %d, токен %d: pos: ожидалось %s, получено %s",
					i, j, state.pos, lex.pos)
			}

			if lex.tok.ID() != state.tok.ID() {
				t.Errorf("тест %d, токен %d: tok.id: ожидалось %d, получено %d",
					i, j, state.tok.ID(), lex.tok.ID())
			}

			if lex.tok.Start().String() != state.tok.Start().String() {
				t.Errorf("тест %d, токен %d: tok.start: ожидалось %s, получено %s",
					i, j, state.tok.Start(), lex.tok.Start())
			}

			val, ok := state.tok.(token.WithValue)
			if !ok {
				continue
			}

			lexval, ok := lex.tok.(token.WithValue)
			if !ok {
				t.Errorf("тест %d, токен %d: tok.(WithValue): ожидалось что токен будет иметь значение",
					i, j)
				continue
			}

			if lexval.Value() != val.Value() {
				t.Errorf("тест %d, токен %d: tok.value: ожидалось %s, получено %s",
					i, j, val.Value(), lexval.Value())
			}
		}
	}
}

func Test_Next_ids_only(t *testing.T) {
	tests := []struct {
		data     string
		expected []uint8
	}{
		{"\n\n\n(9\n+27)	 *\n(\n\n(	 51.2-.2\n) /\n5)", []uint8{
			token.LParen, token.Int, token.Add, token.Int, token.RParen,
			token.Mul, token.LParen, token.LParen, token.Real, token.Sub,
			token.Real, token.RParen, token.Div, token.Int, token.RParen,
		}},
		{"27	 !=51.2\n", []uint8{
			token.Int, token.Neq, token.Real,
		}},
		{"	 'data'||\n'51''.2 data'\n||.2", []uint8{
			token.Text, token.Concat, token.Text, token.Concat, token.Real,
		}},
		{"9\n\n<= 27", []uint8{token.Int, token.Lte, token.Int}},
		{"9 >=\n27", []uint8{token.Int, token.Gte, token.Int}},
		{"9	 <\n	 27", []uint8{token.Int, token.Lt, token.Int}},
		{"\n\n9\n	 >	 27\n", []uint8{token.Int, token.Gt, token.Int}},
		{"and AND And", []uint8{token.And, token.And, token.And}},
		{"or OR Or", []uint8{token.Or, token.Or, token.Or}},
		{"not NOT Not", []uint8{token.Not, token.Not, token.Not}},
		{"true TRUE True", []uint8{token.True, token.True, token.True}},
		{"false FALSE False", []uint8{token.False, token.False, token.False}},
		{"null NULL Null", []uint8{token.Null, token.Null, token.Null}},
		{"res_8 RES_8 Res_8", []uint8{token.Ident, token.Ident, token.Ident}},
		{"_8 8res", []uint8{token.Ident, token.Int, token.Ident}},
		{"привет ПРИВЕТ Привет", []uint8{token.Ident, token.Ident, token.Ident}},
		{"(привет+27)	 *((	 51.2-привет) /привет)*NULL", []uint8{
			token.LParen, token.Ident, token.Add, token.Int, token.RParen,
			token.Mul, token.LParen, token.LParen, token.Real, token.Sub,
			token.Ident, token.RParen, token.Div, token.Ident, token.RParen,
			token.Mul, token.Null,
		}},
		{"oro anda notn", []uint8{token.Ident, token.Ident, token.Ident}},
	}

	for i, test := range tests {
		lex := newLexer(test.data)

		for j, id := range test.expected {
			err := lex.Next()
			if err != nil {
				t.Errorf("%d: %d: непредвиденная ошибка: %s", i, j, err.Error())
				return
			}

			if id != lex.Tok().ID() {
				t.Errorf("%d: %d: ожидалось %d, получено %d",
					i, j, id, lex.Tok().ID())
			}
		}
	}
}

func Test_Next_errors(t *testing.T) {
	//все тесты должны пройти с ошибками
	tests := []struct {
		data     string
		expected error
	}{
		{
			data: "!",
			expected: errExpected(
				pos.NewWithStart(1, 2),
				'=',
			),
		},
		{
			data: "!!",
			expected: errExpected(
				pos.NewWithStart(1, 2),
				'=',
			),
		},
		{
			data: "|",
			expected: errExpected(
				pos.NewWithStart(1, 2),
				'|',
			),
		},
		{
			data: "'",
			expected: errExpected(
				pos.NewWithStart(1, 2),
				'\'',
			),
		},
		{
			data: "'''",
			expected: errExpected(
				pos.NewWithStart(1, 4),
				'\'',
			),
		},
		{
			data: "+#",
			expected: errUnexpected(
				pos.NewWithStart(1, 2),
				'#',
			),
		},
		{
			data: ".",
			expected: errUnexpected(
				pos.NewWithStart(1, 1),
				'.',
			),
		},
		{
			data: "+.",
			expected: errUnexpected(
				pos.NewWithStart(1, 2),
				'.',
			),
		},
	}

	for i, test := range tests {
		lex := newLexer(test.data)
		var err error
		for err == nil {
			err = lex.Next()
		}
		if err.Error() != test.expected.Error() {
			t.Errorf("%d: ожидалось %q, получено %q",
				i, test.expected.Error(), err.Error())
		}
	}
}
