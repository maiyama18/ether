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
			desc:  "one-char operators",
			input: "=+-*/(),;",
			expectedTokens: []token.Token{
				{Type: token.ASSIGN, Literal: "=", Line: 1},
				{Type: token.PLUS, Literal: "+", Line: 1},
				{Type: token.MINUS, Literal: "-", Line: 1},
				{Type: token.ASTER, Literal: "*", Line: 1},
				{Type: token.SLASH, Literal: "/", Line: 1},
				{Type: token.LPAREN, Literal: "(", Line: 1},
				{Type: token.RPAREN, Literal: ")", Line: 1},
				{Type: token.COMMA, Literal: ",", Line: 1},
				{Type: token.SEMICOLON, Literal: ";", Line: 1},
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
