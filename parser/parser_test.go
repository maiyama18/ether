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

func TestParser_ParseProgram_Boolean(t *testing.T) {
	tests := []struct {
		desc     string
		input    string
		expected bool
	}{
		{
			desc:     "true",
			input:    "true;",
			expected: true,
		},
		{
			desc:     "false",
			input:    "false;",
			expected: false,
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
		expectedRight    interface{}
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
		{
			desc:             "!true",
			input:            "!true;",
			expectedOperator: "!",
			expectedRight:    true,
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

func TestParser_ParseProgram_IfExpression(t *testing.T) {
	tests := []struct {
		desc                 string
		input                string
		expectedConditionStr string
		expectedConsequence  interface{}
		expectedAlternative  interface{}
	}{
		{
			desc:  "if(true){10;}",
			input: "if (true) { 10; }",
			expectedConditionStr: "true",
			expectedConsequence: 10,
			expectedAlternative: nil,
		},
		{
			desc:  "if(!false){10;}",
			input: "if (!false) { 10; } else { 9; }",
			expectedConditionStr: "(!false)",
			expectedConsequence: 10,
			expectedAlternative: 9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			program := parseProgram(t, tt.input)
			expression := convertStatementsToSingleExpression(t, program.Statements)

			ifExpression, ok := expression.(*ast.IfExpression)
			if !ok {
				t.Errorf("statement type wrong.\nwant=%T\ngot=%T (%v)\n", &ast.IfExpression{}, ifExpression, ifExpression)
			}
			if ifExpression.Condition.String() != tt.expectedConditionStr {
				t.Errorf("condition string wrong\nwant=%s\ngot=%s\n", tt.expectedConditionStr, ifExpression.Condition.String())
			}
			consequence := convertStatementsToSingleExpression(t, ifExpression.Consequence.Statements)
			testLiteral(t, tt.expectedConsequence, consequence)
			if tt.expectedAlternative == nil {
				if ifExpression.Alternative != nil {
					t.Errorf("nil expected for alternative but got: %+v (%T)\n", ifExpression.Alternative, ifExpression.Alternative)
				}
			} else {
				alternative := convertStatementsToSingleExpression(t, ifExpression.Alternative.Statements)
				testLiteral(t, tt.expectedAlternative, alternative)
			}
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

func TestParser_ParseProgram_FunctionCall(t *testing.T) {
	tests := []struct {
		desc         string
		input        string
		expectedName string
		expectedArgs []interface{}
	}{
		{
			desc:         "g();",
			input:        "g();",
			expectedName: "g",
			expectedArgs: []interface{}{},
		},
		{
			desc:         "add(x,1);",
			input:        "add(x, 1);",
			expectedName: "add",
			expectedArgs: []interface{}{"x", 1},
		},
		{
			desc:         "|a,b|{a+b;}(2,3);",
			input:        "|a, b| { a + b; }(2, 3);",
			expectedArgs: []interface{}{2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			program := parseProgram(t, tt.input)
			expression := convertStatementsToSingleExpression(t, program.Statements)

			functionCall, ok := expression.(*ast.FunctionCall)
			if !ok {
				t.Errorf("statement type wrong.\nwant=%T\ngot=%T (%v)\n", &ast.FunctionCall{}, functionCall, functionCall)
			}

			if tt.expectedName != "" {
				if actualName := functionCall.Function.(*ast.Identifier).Name; actualName != tt.expectedName {
					t.Errorf("function name wrong.\nwant=%s\ngot=%s\n", tt.expectedName, actualName)
				}
			}
			for i, expectedArg := range tt.expectedArgs {
				actualArg := functionCall.Arguments[i]
				testLiteral(t, expectedArg, actualArg)
			}
		})
	}
}

func TestParser_ParseProgram_ArrayLiteral(t *testing.T) {
	tests := []struct {
		desc             string
		input            string
		expectedElements []interface{}
	}{
		{
			desc:             "[]",
			input:            "[];",
			expectedElements: []interface{}{},
		},
		{
			desc:             "[1]",
			input:            "[1];",
			expectedElements: []interface{}{1},
		},
		{
			desc:             "[1,2,3]",
			input:            "[1, 2, 3];",
			expectedElements: []interface{}{1, 2, 3},
		},
		{
			desc:             "[x,y]",
			input:            "[x, y];",
			expectedElements: []interface{}{"x", "y"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			program := parseProgram(t, tt.input)
			expression := convertStatementsToSingleExpression(t, program.Statements)

			arrayLiteral, ok := expression.(*ast.ArrayLiteral)
			if !ok {
				t.Errorf("statement type wrong.\nwant=%T\ngot=%T (%v)\n", &ast.ArrayLiteral{}, arrayLiteral, arrayLiteral)
			}

			for i, expectedElem := range tt.expectedElements {
				actualElem := arrayLiteral.Elements[i]
				testLiteral(t, expectedElem, actualElem)
			}
		})
	}
}

func TestParser_ParseProgram_IndexExpression(t *testing.T) {
	tests := []struct {
		desc          string
		input         string
		expectedArray string
		expectedIndex interface{}
	}{
		{
			desc:          "a[1]",
			input:         "a[1];",
			expectedArray: "a",
			expectedIndex: 1,
		},
		{
			desc:          "a[b]",
			input:         "a[b];",
			expectedArray: "a",
			expectedIndex: "b",
		},
		{
			desc:          "[0,1,2][1]",
			input:         "[0, 1, 2][1];",
			expectedIndex: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			program := parseProgram(t, tt.input)
			expression := convertStatementsToSingleExpression(t, program.Statements)

			indexExpression, ok := expression.(*ast.IndexExpression)
			if !ok {
				t.Errorf("statement type wrong.\nwant=%T\ngot=%T (%v)\n", &ast.IndexExpression{}, indexExpression, indexExpression)
			}

			if tt.expectedArray != "" {
				testLiteral(t, tt.expectedArray, indexExpression.Array)
			}
			testLiteral(t, tt.expectedIndex, indexExpression.Index)
		})
	}
}

func TestParser_ParseProgram_ArrowExpression(t *testing.T) {
	tests := []struct {
		desc         string
		input        string
		expectedName string
		expectedArgs []interface{}
	}{
		{
			desc:         "5->g();",
			input:        "5 -> g();",
			expectedName: "g",
			expectedArgs: []interface{}{5},
		},
		{
			desc:         "x->add(1);",
			input:        "x -> add(1);",
			expectedName: "add",
			expectedArgs: []interface{}{"x", 1},
		},
		{
			desc:         "2->|a,b|{a+b;}(3);",
			input:        "2 -> |a, b| { a + b; }(3);",
			expectedArgs: []interface{}{2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			program := parseProgram(t, tt.input)
			expression := convertStatementsToSingleExpression(t, program.Statements)

			functionCall, ok := expression.(*ast.FunctionCall)
			if !ok {
				t.Errorf("statement type wrong.\nwant=%T\ngot=%T (%v)\n", &ast.FunctionCall{}, functionCall, functionCall)
			}

			if tt.expectedName != "" {
				if actualName := functionCall.Function.(*ast.Identifier).Name; actualName != tt.expectedName {
					t.Errorf("function name wrong.\nwant=%s\ngot=%s\n", tt.expectedName, actualName)
				}
			}
			for i, expectedArg := range tt.expectedArgs {
				actualArg := functionCall.Arguments[i]
				testLiteral(t, expectedArg, actualArg)
			}
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
	case bool:
		booleanLiteral, ok := expression.(*ast.BooleanLiteral)
		if !ok {
			t.Errorf("expression type wrong.\nwant=%T\ngot=%T (%v)\n", &ast.BooleanLiteral{}, booleanLiteral, booleanLiteral)
		}
		if expected != booleanLiteral.Value {
			t.Errorf("boolean value wrong.\nwant=%+v\ngot=%+v\n", expected, booleanLiteral.Value)
		}
	case string:
		identifier, ok := expression.(*ast.Identifier)
		if !ok {
			t.Errorf("expression type wrong.\nwant=%T\ngot=%T (%v)\n", &ast.Identifier{}, identifier, identifier)
		}
		if expected != identifier.Name {
			t.Errorf("identifier name wrong.\nwant=%+v\ngot=%+v\n", expected, identifier.Name)
		}
	case nil:
		if expression != nil {
			t.Errorf("expression type wrong.\nwant=%T\ngot=%T (%v)\n", nil, expression, expression)
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
