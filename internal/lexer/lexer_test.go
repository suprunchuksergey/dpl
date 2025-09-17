package lexer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Tokenize(t *testing.T) {
	tests := []struct {
		data          string
		expectedValue []Token
		expectedError error
	}{
		{"+", []Token{newToken(Add), newToken(EOF)}, nil},

		{"-", []Token{newToken(Sub), newToken(EOF)}, nil},
		{"->", []Token{newToken(ArrowRight), newToken(EOF)}, nil},
		{"-->", []Token{
			newToken(Sub),
			newToken(ArrowRight),
			newToken(EOF)}, nil},

		{"*", []Token{newToken(Mul), newToken(EOF)}, nil},
		{"/", []Token{newToken(Div), newToken(EOF)}, nil},
		{"%", []Token{newToken(Mod), newToken(EOF)}, nil},

		{"+-*/%", []Token{
			newToken(Add),
			newToken(Sub),
			newToken(Mul),
			newToken(Div),
			newToken(Mod),
			newToken(EOF)}, nil},

		{"||", []Token{newToken(Concat), newToken(EOF)}, nil},
		{"||||", []Token{newToken(Concat), newToken(Concat), newToken(EOF)}, nil},

		{"=", []Token{newToken(Set), newToken(EOF)}, nil},
		{"==", []Token{newToken(Eq), newToken(EOF)}, nil},
		{"===", []Token{newToken(Eq), newToken(Set), newToken(EOF)}, nil},
		{"====", []Token{newToken(Eq), newToken(Eq), newToken(EOF)}, nil},

		{"!=", []Token{newToken(Neq), newToken(EOF)}, nil},

		{"<", []Token{newToken(Lt), newToken(EOF)}, nil},
		{"<=", []Token{newToken(Lte), newToken(EOF)}, nil},
		{"<	=", []Token{newToken(Lt), newToken(Set), newToken(EOF)}, nil},

		{">", []Token{newToken(Gt), newToken(EOF)}, nil},
		{">=", []Token{newToken(Gte), newToken(EOF)}, nil},
		{">	=", []Token{newToken(Gt), newToken(Set), newToken(EOF)}, nil},

		{"(", []Token{newToken(LParen), newToken(EOF)}, nil},
		{")", []Token{newToken(RParen), newToken(EOF)}, nil},
		{"()", []Token{newToken(LParen), newToken(RParen), newToken(EOF)}, nil},

		{"[", []Token{newToken(LBrack), newToken(EOF)}, nil},
		{"]", []Token{newToken(RBrack), newToken(EOF)}, nil},
		{"[]", []Token{newToken(LBrack), newToken(RBrack), newToken(EOF)}, nil},

		{"{", []Token{newToken(LBrace), newToken(EOF)}, nil},
		{"}", []Token{newToken(RBrace), newToken(EOF)}, nil},
		{"{}", []Token{newToken(LBrace), newToken(RBrace), newToken(EOF)}, nil},

		{";", []Token{newToken(Semicolon), newToken(EOF)}, nil},

		{":", []Token{newToken(Colon), newToken(EOF)}, nil},
		{":=", []Token{newToken(Create), newToken(EOF)}, nil},

		{",", []Token{newToken(Comma), newToken(EOF)}, nil},
		{".", []Token{newToken(Dot), newToken(EOF)}, nil},

		{"token", []Token{newTokenWithValue(Ident, "token"), newToken(EOF)}, nil},
		{"_token", []Token{newTokenWithValue(Ident, "_token"), newToken(EOF)}, nil},
		{"token2187", []Token{newTokenWithValue(Ident, "token2187"), newToken(EOF)}, nil},
		{"token_2187", []Token{newTokenWithValue(Ident, "token_2187"), newToken(EOF)}, nil},

		{"and", []Token{newToken(And), newToken(EOF)}, nil},
		{"or", []Token{newToken(Or), newToken(EOF)}, nil},
		{"not", []Token{newToken(Not), newToken(EOF)}, nil},
		{"andornot", []Token{newTokenWithValue(Ident, "andornot"), newToken(EOF)}, nil},
		{"and or not", []Token{
			newToken(And),
			newToken(Or),
			newToken(Not),
			newToken(EOF)}, nil},

		{"if", []Token{newToken(If), newToken(EOF)}, nil},
		{"elif", []Token{newToken(Elif), newToken(EOF)}, nil},
		{"else", []Token{newToken(Else), newToken(EOF)}, nil},

		{"for", []Token{newToken(For), newToken(EOF)}, nil},
		{"in", []Token{newToken(In), newToken(EOF)}, nil},

		{"return", []Token{newToken(Return), newToken(EOF)}, nil},

		{"true", []Token{newToken(True), newToken(EOF)}, nil},
		{"false", []Token{newToken(False), newToken(EOF)}, nil},
		{"null", []Token{newToken(Null), newToken(EOF)}, nil},

		{"", []Token{newToken(EOF)}, nil},

		{"2187", []Token{newTokenWithValue(Int, "2187"), newToken(EOF)}, nil},
		{"2.187", []Token{newTokenWithValue(Real, "2.187"), newToken(EOF)}, nil},
		{".2187", []Token{newTokenWithValue(Real, "0.2187"), newToken(EOF)}, nil},
		{"..2187", []Token{
			newToken(Dot),
			newTokenWithValue(Real, "0.2187"),
			newToken(EOF)}, nil},
		{"2187.", []Token{newTokenWithValue(Real, "2187.0"), newToken(EOF)}, nil},
		{"2187..", []Token{
			newTokenWithValue(Real, "2187.0"),
			newToken(Dot),
			newToken(EOF)}, nil},

		{`"token"`, []Token{newTokenWithValue(Text, "token"), newToken(EOF)}, nil},

		{`2187*19683 - 512%1 || "рублей"`, []Token{
			newTokenWithValue(Int, "2187"),
			newToken(Mul),
			newTokenWithValue(Int, "19683"),
			newToken(Sub),
			newTokenWithValue(Int, "512"),
			newToken(Mod),
			newTokenWithValue(Int, "1"),
			newToken(Concat),
			newTokenWithValue(Text, "рублей"),
			newToken(EOF)}, nil},

		{`.2187*19.683 - 512.%1.0001 || "рублей"`, []Token{
			newTokenWithValue(Real, "0.2187"),
			newToken(Mul),
			newTokenWithValue(Real, "19.683"),
			newToken(Sub),
			newTokenWithValue(Real, "512.0"),
			newToken(Mod),
			newTokenWithValue(Real, "1.0001"),
			newToken(Concat),
			newTokenWithValue(Text, "рублей"),
			newToken(EOF)}, nil},

		{`		 .2187 	*	19.683   -512.	 %1.0001||"рублей"	 	`, []Token{
			newTokenWithValue(Real, "0.2187"),
			newToken(Mul),
			newTokenWithValue(Real, "19.683"),
			newToken(Sub),
			newTokenWithValue(Real, "512.0"),
			newToken(Mod),
			newTokenWithValue(Real, "1.0001"),
			newToken(Concat),
			newTokenWithValue(Text, "рублей"),
			newToken(EOF)}, nil},

		{`
for i in 2187 {
	if i%2 == 0 {
		print(i);
	} else {
		print(i-1);
	};
};
`, []Token{
			newToken(For),
			newTokenWithValue(Ident, "i"),
			newToken(In),
			newTokenWithValue(Int, "2187"),
			newToken(LBrace),

			newToken(If),
			newTokenWithValue(Ident, "i"),
			newToken(Mod),
			newTokenWithValue(Int, "2"),
			newToken(Eq),
			newTokenWithValue(Int, "0"),
			newToken(LBrace),

			newTokenWithValue(Ident, "print"),
			newToken(LParen),
			newTokenWithValue(Ident, "i"),
			newToken(RParen),
			newToken(Semicolon),

			newToken(RBrace),
			newToken(Else),
			newToken(LBrace),

			newTokenWithValue(Ident, "print"),
			newToken(LParen),
			newTokenWithValue(Ident, "i"),
			newToken(Sub),
			newTokenWithValue(Int, "1"),
			newToken(RParen),
			newToken(Semicolon),

			newToken(RBrace),
			newToken(Semicolon),

			newToken(RBrace),
			newToken(Semicolon),

			newToken(EOF)}, nil},

		{`
factorial := (n) -> {
	if n <= 1 {
		return 1;
	};
	return n*factorial(n-1);
};
`, []Token{
			newTokenWithValue(Ident, "factorial"),
			newToken(Create),
			newToken(LParen),
			newTokenWithValue(Ident, "n"),
			newToken(RParen),
			newToken(ArrowRight),
			newToken(LBrace),

			newToken(If),
			newTokenWithValue(Ident, "n"),
			newToken(Lte),
			newTokenWithValue(Int, "1"),
			newToken(LBrace),

			newToken(Return),
			newTokenWithValue(Int, "1"),
			newToken(Semicolon),

			newToken(RBrace),
			newToken(Semicolon),

			newToken(Return),
			newTokenWithValue(Ident, "n"),
			newToken(Mul),
			newTokenWithValue(Ident, "factorial"),
			newToken(LParen),
			newTokenWithValue(Ident, "n"),
			newToken(Sub),
			newTokenWithValue(Int, "1"),
			newToken(RParen),
			newToken(Semicolon),

			newToken(RBrace),
			newToken(Semicolon),

			newToken(EOF)}, nil},

		{"#", nil, unexpected('#')},
		{"|", nil, expected('|')},
		{"!", nil, expected('=')},
	}

	for _, test := range tests {
		v, err := Tokenize(test.data)

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, v)
		}
	}
}
