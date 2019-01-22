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
			desc:  "empty input",
			input: "",
			expectedTokens: []token.Token{
				{Type: token.EOF, Literal: "", Line: 1},
			},
		},
	}

	for _, tt := range tests {
		lexer := New(tt.input)

		t.Run(tt.desc, func(t *testing.T) {
			for _, expected := range tt.expectedTokens {
				actual := lexer.NextToken()
				if actual != expected {
					t.Errorf("wrong token. \nwant:%+v\ngot:%+v\n", actual, expected)
				}
			}
		})
	}
}
