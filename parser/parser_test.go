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
			input:        "var a = 5;",
			expectedName: "a",
		},
		{
			desc:         "multiple-char identifier",
			input:        "var foo = 42;",
			expectedName: "foo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			program := parseProgram(tt.input)

			if len(program.Statements) != 1 {
				t.Errorf("statements length wrong.\nwant=%d\ngot=%d\n", 1, len(program.Statements))
			}
			varStatement, ok := program.Statements[0].(*ast.VarStatement)
			if !ok {
				t.Errorf("statement type wrong.\nwant=%T\ngot=%T (%v)\n", &ast.VarStatement{}, varStatement, varStatement)
			}
			if varStatement.Identifier.Name != tt.expectedName {
				t.Errorf("identifier name wrong.\nwant=%q\ngot=%q\n", tt.expectedName, varStatement.Identifier.Name)
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

			testLiteral(t, tt.expected, expression)
		})
	}
}

func TestParser_ParseProgram_PrefixExpression(t *testing.T) {
	tests := []struct {
		desc             string
		input            string
		expectedOperator string
		expectedRight    int
	}{
		{
			desc:             "-5",
			input:            "-5;",
			expectedOperator: "-",
			expectedRight:    5,
		},
		{
			desc:             "-42",
			input:            "-42;",
			expectedOperator: "-",
			expectedRight:    42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			program := parseProgram(tt.input)
			expression := convertProgramToSingleExpression(t, program)

			prefixExpression, ok := expression.(*ast.PrefixExpression)
			if !ok {
				t.Errorf("statement type wrong.\nwant=%T\ngot=%T (%v)\n", &ast.PrefixExpression{}, prefixExpression, prefixExpression)
			}
			if prefixExpression.Operator != tt.expectedOperator {
				t.Errorf("operator wrong.\nwant=%+v\ngot=%+v\n", tt.expectedOperator, prefixExpression.Operator)
			}
			testLiteral(t, tt.expectedRight, prefixExpression.Right)
		})
	}
}

// func TestParser_ParseProgram_InfixExpression(t *testing.T) {
// 	tests := []struct {
// 		desc             string
// 		input            string
// 		expectedOperator string
// 		expectedLeft     int
// 		expectedRight    int
// 	}{
// 		{
// 			desc:             "addition",
// 			input:            "2 + 3;",
// 			expectedOperator: "+",
// 			expectedLeft:     2,
// 			expectedRight:    3,
// 		},
// 		{
// 			desc:             "subtraction",
// 			input:            "2 - 3;",
// 			expectedOperator: "-",
// 			expectedLeft:     2,
// 			expectedRight:    3,
// 		},
// 		{
// 			desc:             "multiplication",
// 			input:            "2 * 3;",
// 			expectedOperator: "*",
// 			expectedLeft:     2,
// 			expectedRight:    3,
// 		},
// 		{
// 			desc:             "division",
// 			input:            "2 / 3;",
// 			expectedOperator: "*",
// 			expectedLeft:     2,
// 			expectedRight:    3,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.desc, func(t *testing.T) {
// 			program := parseProgram(tt.input)
// 			expression := convertProgramToSingleExpression(t, program)
//
// 			infixExpression, ok := expression.(*ast.InfixExpression)
// 			if !ok {
// 				t.Errorf("statement type wrong.\nwant=%T\ngot=%T (%v)\n", &ast.InfixExpression{}, infixExpression, infixExpression)
// 			}
// 			if infixExpression.Operator != tt.expectedOperator {
// 				t.Errorf("operator wrong.\nwant=%+v\ngot=%+v\n", tt.expectedOperator, infixExpression.Operator)
// 			}
// 			testLiteral(t, tt.expectedLeft, infixExpression.Left)
// 			testLiteral(t, tt.expectedRight, infixExpression.Right)
// 		})
// 	}
// }

func parseProgram(input string) *ast.Program {
	lex := lexer.New(input)
	parser := New(lex)

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

func testLiteral(t *testing.T, expected interface{}, expression ast.Expression) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		integerLiteral, ok := expression.(*ast.IntegerLiteral)
		if !ok {
			t.Errorf("expression type wrong.\nwant=%T\ngot=%T (%v)\n", &ast.IntegerLiteral{}, integerLiteral, integerLiteral)
		}
		if expected != integerLiteral.Value {
			t.Errorf("integer value wrong.\nwant=%+v\ngot=%+v\n", expected, integerLiteral.Value)
		}
	}
}
