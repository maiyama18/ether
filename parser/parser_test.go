package parser

import (
	"github.com/muiscript/ether/ast"
	"github.com/muiscript/ether/lexer"
	"testing"
)

func TestParser_ParseProgram_LetStatement(t *testing.T) {
	tests := []struct {
		desc         string
		input        string
		expectedName string
	}{
		{
			desc:         "simple",
			input:        "let a = 5;",
			expectedName: "a",
		},
		{
			desc:         "multiple-char identifier",
			input:        "let foo = 42;",
			expectedName: "foo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			program := parseProgram(tt.input)

			if len(program.Statements) != 1 {
				t.Errorf("statements length wrong.\nwant=%d\ngot=%d\n", 1, len(program.Statements))
			}
			letStatement, ok := program.Statements[0].(*ast.LetStatement)
			if !ok {
				t.Errorf("statement type wrong.\nwant=%T\ngot=%T (%v)\n", &ast.LetStatement{}, letStatement, letStatement)
			}
			if letStatement.Identifier.Name != tt.expectedName {
				t.Errorf("identifier name wrong.\nwant=%q\ngot=%q\n", tt.expectedName, letStatement.Identifier.Name)
			}
		})
	}
}

func parseProgram(input string) *ast.Program {
	lexer := lexer.New(input)
	parser := New(lexer)

	return parser.ParseProgram()
}
