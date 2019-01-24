package parser

import (
	"github.com/muiscript/ether/ast"
	"github.com/muiscript/ether/lexer"
	"testing"
)

func TestParser_ParseProgram_VarStatement(t *testing.T) {
	tests := []struct {
		desc               string
		input              string
		expectedName       string
		expectedExpression interface{}
	}{
		{
			desc:               "var a = 5",
			input:              "var a = 5;",
			expectedName:       "a",
			expectedExpression: 5,
		},
		{
			desc:               "var foo = 42",
			input:              "var foo = 42;",
			expectedName:       "foo",
			expectedExpression: 42,
		},
		{
			desc:               "var foo = bar",
			input:              "var foo = bar;",
			expectedName:       "foo",
			expectedExpression: "bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			program := parseProgram(t, tt.input)

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
			testLiteral(t, tt.expectedExpression, varStatement.Expression)
		})
	}
}

func TestParser_ParseProgram_ReturnStatement(t *testing.T) {
	tests := []struct {
		desc               string
		input              string
		expectedExpression interface{}
	}{
		{
			desc:               "return a",
			input:              "return a;",
			expectedExpression: "a",
		},
		{
			desc:               "return foo",
			input:              "return foo;",
			expectedExpression: "foo",
		},
		{
			desc:               "return 42",
			input:              "return 42;",
			expectedExpression: 42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			program := parseProgram(t, tt.input)

			if len(program.Statements) != 1 {
				t.Errorf("statements length wrong.\nwant=%d\ngot=%d\n", 1, len(program.Statements))
			}
			returnStatement, ok := program.Statements[0].(*ast.ReturnStatement)
			if !ok {
				t.Errorf("statement type wrong.\nwant=%T\ngot=%T (%v)\n", &ast.ReturnStatement{}, returnStatement, returnStatement)
			}
			testLiteral(t, tt.expectedExpression, returnStatement.Expression)
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
			program := parseProgram(t, tt.input)
			expression := convertStatementsToSingleExpression(t, program.Statements)

			testLiteral(t, tt.expected, expression)
		})
	}
}

func TestParser_ParseProgram_Identifier(t *testing.T) {
	tests := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:     "a",
			input:    "a;",
			expected: "a",
		},
		{
			desc:     "foo",
			input:    "foo;",
			expected: "foo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			program := parseProgram(t, tt.input)
			expression := convertStatementsToSingleExpression(t, program.Statements)

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
			program := parseProgram(t, tt.input)
			expression := convertStatementsToSingleExpression(t, program.Statements)

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

func TestParser_ParseProgram_InfixExpression(t *testing.T) {
	tests := []struct {
		desc             string
		input            string
		expectedOperator string
		expectedLeft     int
		expectedRight    int
	}{
		{
			desc:             "addition",
			input:            "2 + 3;",
			expectedOperator: "+",
			expectedLeft:     2,
			expectedRight:    3,
		},
		{
			desc:             "subtraction",
			input:            "2 - 3;",
			expectedOperator: "-",
			expectedLeft:     2,
			expectedRight:    3,
		},
		{
			desc:             "multiplication",
			input:            "2 * 3;",
			expectedOperator: "*",
			expectedLeft:     2,
			expectedRight:    3,
		},
		{
			desc:             "division",
			input:            "2 / 3;",
			expectedOperator: "/",
			expectedLeft:     2,
			expectedRight:    3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			program := parseProgram(t, tt.input)
			expression := convertStatementsToSingleExpression(t, program.Statements)

			infixExpression, ok := expression.(*ast.InfixExpression)
			if !ok {
				t.Errorf("statement type wrong.\nwant=%T\ngot=%T (%v)\n", &ast.InfixExpression{}, infixExpression, infixExpression)
			}
			testInfixExpression(t, tt.expectedOperator, tt.expectedLeft, tt.expectedRight, infixExpression)
		})
	}
}

func TestParser_ParseProgram_FunctionLiteral(t *testing.T) {
	tests := []struct {
		desc                string
		input               string
		expectedParamNames  []string
		expectedReturnValue interface{}
	}{
		{
			desc:                "|| { 42; };",
			input:               "|| { 42; };",
			expectedParamNames:  []string{},
			expectedReturnValue: 42,
		},
		{
			desc:                "|x| { x; };",
			input:               "|x| { x; };",
			expectedParamNames:  []string{"x"},
			expectedReturnValue: "x",
		},
		{
			desc:                "|x, y| { x; };",
			input:               "|x, y| { x; };",
			expectedParamNames:  []string{"x", "y"},
			expectedReturnValue: "x",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			program := parseProgram(t, tt.input)
			expression := convertStatementsToSingleExpression(t, program.Statements)

			functionLiteral, ok := expression.(*ast.FunctionLiteral)
			if !ok {
				t.Errorf("statement type wrong.\nwant=%T\ngot=%T (%v)\n", &ast.FunctionLiteral{}, functionLiteral, functionLiteral)
			}

			if len(functionLiteral.Parameters) != len(tt.expectedParamNames) {
				t.Errorf("parameter length wrong.\nwant=%d\ngot=%d\n", len(tt.expectedParamNames), len(functionLiteral.Parameters))
			}
			for i, expectedParamName := range tt.expectedParamNames {
				if functionLiteral.Parameters[i].Name != expectedParamName {
					t.Errorf("%d-th parameter name wrong.\nwant=%+v\ngot=%+v\n", i, expectedParamName, functionLiteral.Parameters[i].Name)
				}
			}

			if len(functionLiteral.Body.Statements) != 1 {
				t.Errorf("number of statements wrong.\nwant=%d\ngot=%d\n", 1, len(functionLiteral.Body.Statements))
			}
			returnExpression := convertStatementsToSingleExpression(t, functionLiteral.Body.Statements)
			testLiteral(t, tt.expectedReturnValue, returnExpression)
		})
	}
}

func TestParser_ParseProgram_ComplexArithmetic(t *testing.T) {
	tests := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:     "1 + 2 + 3",
			input:    "1 + 2 + 3;",
			expected: "((1 + 2) + 3)",
		},
		{
			desc:     "1 + 2 - 3",
			input:    "1 + 2 - 3;",
			expected: "((1 + 2) - 3)",
		},
		{
			desc:     "1 * 2 + 3",
			input:    "1 * 2 + 3;",
			expected: "((1 * 2) + 3)",
		},
		{
			desc:     "1 + 2 * 3",
			input:    "1 + 2 * 3;",
			expected: "(1 + (2 * 3))",
		},
		{
			desc:     "1 + 2 * 3",
			input:    "1 + 2 * 3;",
			expected: "(1 + (2 * 3))",
		},
		{
			desc:     "- -42",
			input:    "- -42;",
			expected: "(-(-42))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			program := parseProgram(t, tt.input)
			expression := convertStatementsToSingleExpression(t, program.Statements)

			if actual := expression.String(); actual != tt.expected {
				t.Errorf("string expression wrong.\nwant=%q\ngot=%q\n", tt.expected, actual)
			}
		})
	}
}

func parseProgram(t *testing.T, input string) *ast.Program {
	lex := lexer.New(input)
	parser := New(lex)

	program, err := parser.ParseProgram()
	if err != nil {
		t.Fatalf("parse error: %s", err.Error())
	}
	return program
}

func convertStatementsToSingleExpression(t *testing.T, statements []ast.Statement) ast.Expression {
	t.Helper()

	if len(statements) != 1 {
		t.Errorf("statements length wrong.\nwant=%d\ngot=%d\n", 1, len(statements))
	}
	expressionStatement, ok := statements[0].(*ast.ExpressionStatement)
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
	case string:
		identifier, ok := expression.(*ast.Identifier)
		if !ok {
			t.Errorf("expression type wrong.\nwant=%T\ngot=%T (%v)\n", &ast.Identifier{}, identifier, identifier)
		}
		if expected != identifier.Name {
			t.Errorf("identifier name wrong.\nwant=%+v\ngot=%+v\n", expected, identifier.Name)
		}
	}
}

func testInfixExpression(t *testing.T, expectedOperator string, expectedLeft, expectedRight interface{}, infixExpression *ast.InfixExpression) {
	t.Helper()

	if infixExpression.Operator != expectedOperator {
		t.Errorf("operator wrong.\nwant=%+v\ngot=%+v\n", expectedOperator, infixExpression.Operator)
	}
	testLiteral(t, expectedLeft, infixExpression.Left)
	testLiteral(t, expectedRight, infixExpression.Right)
}
