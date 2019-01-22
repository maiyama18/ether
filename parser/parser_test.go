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

func TestParser_ParseProgram_ReturnStatement(t *testing.T) {
	tests := []struct {
		desc  string
		input string
	}{
		{
			desc:  "simple",
			input: "return a;",
		},
		{
			desc:  "multiple-char identifier",
			input: "return foo;",
		},
		{
			desc:  "integer literal",
			input: "return 42;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			program := parseProgram(tt.input)

			if len(program.Statements) != 1 {
				t.Errorf("statements length wrong.\nwant=%d\ngot=%d\n", 1, len(program.Statements))
			}
			returnStatement, ok := program.Statements[0].(*ast.ReturnStatement)
			if !ok {
				t.Errorf("statement type wrong.\nwant=%T\ngot=%T (%v)\n", &ast.ReturnStatement{}, returnStatement, returnStatement)
			}
		})
	}
}

func TestParser_ParseProgram_IntegerLiteral(t *testing.T) {
	tests := []struct {
		desc     string
		input    string
		expected int
	}{
		{
			desc:     "5",
			input:    "5;",
			expected: 5,
		},
		{
			desc:     "42",
			input:    "42;",
			expected: 42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			program := parseProgram(tt.input)
			expression := convertProgramToSingleExpression(t, program)

			integerLiteral, ok := expression.(*ast.IntegerLiteral)
			if !ok {
				t.Errorf("statement type wrong.\nwant=%T\ngot=%T (%v)\n", &ast.IntegerLiteral{}, integerLiteral, integerLiteral)
			}
			if integerLiteral.Value != tt.expected {
				t.Errorf("integer value wrong.\nwant=%+v\ngot=%+v\n", tt.expected, integerLiteral.Value)
			}
		})
	}
}

func parseProgram(input string) *ast.Program {
	lexer := lexer.New(input)
	parser := New(lexer)

	return parser.ParseProgram()
}

func convertProgramToSingleExpression(t *testing.T, program *ast.Program) ast.Expression {
	t.Helper()

	if len(program.Statements) != 1 {
		t.Errorf("statements length wrong.\nwant=%d\ngot=%d\n", 1, len(program.Statements))
	}
	expressionStatement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("statement type wrong.\nwant=%T\ngot=%T (%v)\n", &ast.ExpressionStatement{}, expressionStatement, expressionStatement)
	}
	return expressionStatement.Expression
}
