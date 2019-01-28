package lexer

import (
	"github.com/muiscript/ether/token"
	"testing"
)

func TestLexer_NextToken(t *testing.T) {
	tests := []struct {
		desc           string
		input          string
		expectedTokens []token.Token
	}{
		{
			desc:  "empty",
			input: "",
			expectedTokens: []token.Token{
				{Type: token.EOF, Literal: "", Line: 1},
			},
		},
		{
			desc:  "single-char operators",
			input: "=+-*/!<>(){}[]|,;",
			expectedTokens: []token.Token{
				{Type: token.ASSIGN, Literal: "=", Line: 1},
				{Type: token.PLUS, Literal: "+", Line: 1},
				{Type: token.MINUS, Literal: "-", Line: 1},
				{Type: token.ASTER, Literal: "*", Line: 1},
				{Type: token.SLASH, Literal: "/", Line: 1},
				{Type: token.BANG, Literal: "!", Line: 1},
				{Type: token.LT, Literal: "<", Line: 1},
				{Type: token.GT, Literal: ">", Line: 1},
				{Type: token.LPAREN, Literal: "(", Line: 1},
				{Type: token.RPAREN, Literal: ")", Line: 1},
				{Type: token.LBRACE, Literal: "{", Line: 1},
				{Type: token.RBRACE, Literal: "}", Line: 1},
				{Type: token.LBRACKET, Literal: "[", Line: 1},
				{Type: token.RBRACKET, Literal: "]", Line: 1},
				{Type: token.BAR, Literal: "|", Line: 1},
				{Type: token.COMMA, Literal: ",", Line: 1},
				{Type: token.SEMICOLON, Literal: ";", Line: 1},
				{Type: token.EOF, Literal: "", Line: 1},
			},
		},
		{
			desc:  "multi-char operators",
			input: "== !=",
			expectedTokens: []token.Token{
				{Type: token.EQ, Literal: "==", Line: 1},
				{Type: token.NEQ, Literal: "!=", Line: 1},
				{Type: token.EOF, Literal: "", Line: 1},
			},
		},
		{
			desc:  "arrow operator",
			input: "3 - 2 -> double()",
			expectedTokens: []token.Token{
				{Type: token.INTEGER, Literal: "3", Line: 1},
				{Type: token.MINUS, Literal: "-", Line: 1},
				{Type: token.INTEGER, Literal: "2", Line: 1},
				{Type: token.ARROW, Literal: "->", Line: 1},
				{Type: token.IDENT, Literal: "double", Line: 1},
				{Type: token.LPAREN, Literal: "(", Line: 1},
				{Type: token.RPAREN, Literal: ")", Line: 1},
				{Type: token.EOF, Literal: "", Line: 1},
			},
		},
		{
			desc: "multiple-line input",
			input: `=
+-
*/`,
			expectedTokens: []token.Token{
				{Type: token.ASSIGN, Literal: "=", Line: 1},
				{Type: token.PLUS, Literal: "+", Line: 2},
				{Type: token.MINUS, Literal: "-", Line: 2},
				{Type: token.ASTER, Literal: "*", Line: 3},
				{Type: token.SLASH, Literal: "/", Line: 3},
				{Type: token.EOF, Literal: "", Line: 3},
			},
		},
		{
			desc: "single-char keywords and literals",
			input: `var a = 5;
true;
false;
return a;`,
			expectedTokens: []token.Token{
				{Type: token.VAR, Literal: "var", Line: 1},
				{Type: token.IDENT, Literal: "a", Line: 1},
				{Type: token.ASSIGN, Literal: "=", Line: 1},
				{Type: token.INTEGER, Literal: "5", Line: 1},
				{Type: token.SEMICOLON, Literal: ";", Line: 1},

				{Type: token.TRUE, Literal: "true", Line: 2},
				{Type: token.SEMICOLON, Literal: ";", Line: 2},

				{Type: token.FALSE, Literal: "false", Line: 3},
				{Type: token.SEMICOLON, Literal: ";", Line: 3},

				{Type: token.RETURN, Literal: "return", Line: 4},
				{Type: token.IDENT, Literal: "a", Line: 4},
				{Type: token.SEMICOLON, Literal: ";", Line: 4},
				{Type: token.EOF, Literal: "", Line: 4},
			},
		},
		{
			desc: "var and return statements",
			input: `var foo = 42;
return foo;`,
			expectedTokens: []token.Token{
				{Type: token.VAR, Literal: "var", Line: 1},
				{Type: token.IDENT, Literal: "foo", Line: 1},
				{Type: token.ASSIGN, Literal: "=", Line: 1},
				{Type: token.INTEGER, Literal: "42", Line: 1},
				{Type: token.SEMICOLON, Literal: ";", Line: 1},

				{Type: token.RETURN, Literal: "return", Line: 2},
				{Type: token.IDENT, Literal: "foo", Line: 2},
				{Type: token.SEMICOLON, Literal: ";", Line: 2},
				{Type: token.EOF, Literal: "", Line: 2},
			},
		},
		{
			desc:  "if-else expression",
			input: "if (true) { 10 } else { 9 };",
			expectedTokens: []token.Token{
				{Type: token.IF, Literal: "if", Line: 1},
				{Type: token.LPAREN, Literal: "(", Line: 1},
				{Type: token.TRUE, Literal: "true", Line: 1},
				{Type: token.RPAREN, Literal: ")", Line: 1},
				{Type: token.LBRACE, Literal: "{", Line: 1},
				{Type: token.INTEGER, Literal: "10", Line: 1},
				{Type: token.RBRACE, Literal: "}", Line: 1},
				{Type: token.ELSE, Literal: "else", Line: 1},
				{Type: token.LBRACE, Literal: "{", Line: 1},
				{Type: token.INTEGER, Literal: "9", Line: 1},
				{Type: token.RBRACE, Literal: "}", Line: 1},
				{Type: token.SEMICOLON, Literal: ";", Line: 1},
				{Type: token.EOF, Literal: "", Line: 1},
			},
		},
		{
			desc: "comment",
			input: `var foo = 42;
# ignore me
return foo; # this is comment`,
			expectedTokens: []token.Token{
				{Type: token.VAR, Literal: "var", Line: 1},
				{Type: token.IDENT, Literal: "foo", Line: 1},
				{Type: token.ASSIGN, Literal: "=", Line: 1},
				{Type: token.INTEGER, Literal: "42", Line: 1},
				{Type: token.SEMICOLON, Literal: ";", Line: 1},

				{Type: token.RETURN, Literal: "return", Line: 3},
				{Type: token.IDENT, Literal: "foo", Line: 3},
				{Type: token.SEMICOLON, Literal: ";", Line: 3},
				{Type: token.EOF, Literal: "", Line: 3},
				{Type: token.EOF, Literal: "", Line: 3},
			},
		},
	}

	for _, tt := range tests {
		lexer := New(tt.input)

		t.Run(tt.desc, func(t *testing.T) {
			for _, expected := range tt.expectedTokens {
				actual := lexer.NextToken()
				if actual != expected {
					t.Errorf("wrong token. \nwant:%+v\ngot:%+v\n", expected, actual)
				}
			}
		})
	}
}
