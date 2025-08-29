package lexer

import (
	"fmt"
	"github.com/suprunchuksergey/dpl/pos"
	"github.com/suprunchuksergey/dpl/token"
	"strings"
	"unicode"
)

type lexer struct {
	src   []rune      //источник
	index int         //индекс текущего символа
	char  rune        //текущий символ, 0 это EOF
	pos   pos.Pos     //текущая позиция
	tok   token.Token //последний прочитанный токен
}

func newLexer(data string) *lexer {
	src := []rune(data)
	var char rune //по умолчанию символ равен 0 (EOF)
	//если срез не пустой,
	//то текущий символ становится первым символом из среза
	if len(src) != 0 {
		char = src[0]
	}
	return &lexer{
		src:  src,
		char: char,
		pos:  pos.New(),
	}
}

// сдвигает index на 1 вперед, изменяет pos и char
func (l *lexer) next() {
	//индекс уже вышел за пределы файла,
	//предполагается, что char уже равен 0,
	//а pos там, где и должен быть
	if l.index >= len(l.src) {
		return
	}

	//сдвиг позиции
	if /*следующий символ на новой строке*/ l.char == '\n' {
		l.pos.MoveLine()
	} else /*следующий символ в соседнем столбце*/ {
		l.pos.MoveColumn()
	}

	//сдвиг индекса
	l.index++
	if l.index == len(l.src) /*индекс указывает на конец файла*/ {
		l.char = 0 //EOF
		return
	}
	l.char = l.src[l.index]
}

// пропустить все пробельные символы
func (l *lexer) skipSpaces() {
	for unicode.IsSpace(l.char) {
		l.next()
	}
}

// вспомогательная функция, меняет tok на новый токен и вызывает метод next.
// используется для токенов, состоящих из одного символа и не имеющих значения (+, *, ...).
func (l *lexer) h(id uint8, start pos.Pos) {
	l.tok = token.New(id, start)
	l.next()
}

// считывает символы до тех пор, пока они соответствуют условию
// и символ не 0
func (l *lexer) readWhile(cond func(rune) bool) string {
	var res strings.Builder
	for cond(l.char) && l.char != 0 {
		res.WriteRune(l.char)
		l.next()
	}
	return res.String()
}

// считывает символы от кавычки до кавычки, учитывает экранирование,
// возвращает ошибку, если строка не завершена
func (l *lexer) readText() (string, error) {
	//если символ не является кавычкой
	if l.char != '\'' {
		//выйти из функции
		return "", nil
	}
	//перейти к следующему символу
	l.next()

	//прочитать все символы до кавычки или EOF
	res := l.readWhile(func(r rune) bool { return r != '\'' })
	if l.char == 0 /*EOF*/ {
		return "", errExpected(l.pos, '\'')
	}
	//если это не EOF, то это кавычка
	//перейти к следующему символу
	l.next()

	//если следующий символ не является кавычкой
	if l.char != '\'' {
		//выходим из функции
		return res, nil
	}
	//в противном случае это экранирование кавычки
	res += "'"

	text, err := l.readText()
	if err != nil {
		return "", err
	}

	return res + text, nil
}

// вспомогательная функция,
// считывает все цифры до первой не цифры, возвращает все считанные цифры
func (l *lexer) readNums() string {
	return l.readWhile(unicode.IsDigit)
}

// прочитать хвост числа от точки,
// если после точки нет цифр, будет возвращен ".0"
func (l *lexer) readTail() string {
	if l.char != '.' {
		return ""
	}

	var res strings.Builder

	res.WriteRune('.')
	l.next()

	nums := l.readNums()
	if len(nums) == 0 {
		res.WriteRune('0')
	} else {
		res.WriteString(nums)
	}

	return res.String()
}

func (l *lexer) Next() error {
	l.skipSpaces()         // пропустить все пробельные символы
	start := l.pos.Clone() //начало токена
	switch l.char {
	case 0 /*EOF*/ :
		l.tok = token.New(token.EOF, start)
	case '(':
		l.h(token.LParen, start)
	case ')':
		l.h(token.RParen, start)
	case '[':
		l.h(token.LBrack, start)
	case ']':
		l.h(token.RBrack, start)
	case '{':
		l.h(token.LBrace, start)
	case '}':
		l.h(token.RBrace, start)
	case ',':
		l.h(token.Comma, start)
	case ':':
		l.h(token.Colon, start)

	case ';':
		l.h(token.Semicolon, start)

	case '+':
		l.h(token.Add, start)
	case '-':
		l.h(token.Sub, start)
	case '*':
		l.h(token.Mul, start)
	case '/':
		l.h(token.Div, start)
	case '%':
		l.h(token.Rem, start)
	case '=':
		l.next()
		if l.char == '=' {
			l.next()
			l.tok = token.New(token.Eq, start)
			break
		}
		l.tok = token.New(token.Assign, start)

	case '!':
		l.next()
		if l.char != '=' {
			return errExpected(l.pos, '=')
		}
		l.next()
		l.tok = token.New(token.Neq, start)
	case '<':
		l.next()
		if l.char == '>' {
			l.next()
			l.tok = token.New(token.Neq, start)
		} else if l.char == '=' {
			l.next()
			l.tok = token.New(token.Lte, start)
		} else {
			l.tok = token.New(token.Lt, start)
		}
	case '>':
		l.next()
		if l.char == '=' {
			l.next()
			l.tok = token.New(token.Gte, start)
		} else {
			l.tok = token.New(token.Gt, start)
		}
	case '|':
		l.next()
		if l.char != '|' {
			return errExpected(l.pos, '|')
		}
		l.next()
		l.tok = token.New(token.Concat, start)

	//строка
	case '\'':
		value, err := l.readText()
		if err != nil {
			return err
		}
		l.tok = token.NewWithValue(token.Text, start, value)

	default:
		//случай, когда простого сравнения с символом недостаточно,
		//а необходимы полноценные проверки (является ли символ цифрой, буквой и т. д.)
		switch {
		//число, начинающееся с нуля
		case l.char == '0':
			var res strings.Builder
			id := token.Int //id по умолчанию int
			//записал ноль и перешел к следующему символу
			res.WriteRune('0')
			l.next()

			tail := l.readTail()
			if len(tail) != 0 {
				id = token.Real
				res.WriteString(tail)
			}
			l.tok = token.NewWithValue(id, start, res.String())

		//число, начинающееся с точки
		case l.char == '.':
			var res strings.Builder

			//записать точку, дополненную нулем
			res.WriteString("0.")
			l.next()

			nums := l.readNums()
			//точка без начального числа, за которым не следует никаких цифр
			if len(nums) == 0 {
				return errUnexpected(start, '.')
			}
			res.WriteString(nums)
			l.tok = token.NewWithValue(token.Real, start, res.String())

		//число
		case unicode.IsDigit(l.char):
			var res strings.Builder
			id := token.Int //id по умолчанию int

			res.WriteString(l.readNums())

			tail := l.readTail()
			if len(tail) != 0 {
				id = token.Real
				res.WriteString(tail)
			}
			l.tok = token.NewWithValue(id, start, res.String())

		//идентификатор
		case unicode.IsLetter(l.char) || l.char == '_':
			var res strings.Builder

			for unicode.IsDigit(l.char) || l.char == '_' || unicode.IsLetter(l.char) {
				res.WriteRune(unicode.ToLower(l.char))
				l.next()
			}

			switch res.String() {
			case "if":
				l.tok = token.New(token.If, start)
			case "elif":
				l.tok = token.New(token.Elif, start)
			case "else":
				l.tok = token.New(token.Else, start)

			case "return":
				l.tok = token.New(token.Return, start)
			case "fn":
				l.tok = token.New(token.Fn, start)

			case "for":
				l.tok = token.New(token.For, start)
			case "in":
				l.tok = token.New(token.In, start)

			case "and":
				l.tok = token.New(token.And, start)
			case "or":
				l.tok = token.New(token.Or, start)
			case "not":
				l.tok = token.New(token.Not, start)
			case "true":
				l.tok = token.New(token.True, start)
			case "false":
				l.tok = token.New(token.False, start)
			case "null":
				l.tok = token.New(token.Null, start)
			default:
				l.tok = token.NewWithValue(token.Ident, start, res.String())
			}

		default:
			return errUnexpected(l.pos, l.char)
		}
	}
	return nil
}

func (l *lexer) Tok() token.Token {
	return l.tok
}

func errExpected(pos pos.Pos, char rune) error {
	return fmt.Errorf("%s ожидался символ %c", pos, char)
}

func errUnexpected(pos pos.Pos, char rune) error {
	return fmt.Errorf("%s неожиданный символ %c", pos, char)
}

type Lexer interface {
	// Next считывает новый токен из последовательности,
	//вернет ошибку, если токен не может быть прочитан (встретил неизвестный символ и т. д.)
	Next() error
	// Tok возвращает последний прочитанный токен
	Tok() token.Token
}

func New(data string) Lexer {
	return newLexer(data)
}
