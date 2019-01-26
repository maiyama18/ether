package evaluator

import (
	"github.com/muiscript/ether/lexer"
	"github.com/muiscript/ether/object"
	"github.com/muiscript/ether/parser"
	"testing"
)

func TestEval_Integer(t *testing.T) {
	tests := []struct {
		desc     string
		input    string
		expected int
	}{
		{
			desc:     "42",
			input:    "42;",
			expected: 42,
		},
		{
			desc:     "-42",
			input:    "-42;",
			expected: -42,
		},
		{
			desc:     "-(-42)",
			input:    "-(-42);",
			expected: 42,
		},
		{
			desc:     "15 + 3",
			input:    "15 + 3;",
			expected: 18,
		},
		{
			desc:     "15 - 3",
			input:    "15 - 3;",
			expected: 12,
		},
		{
			desc:     "15 * 3",
			input:    "15 * 3;",
			expected: 45,
		},
		{
			desc:     "15 / 3",
			input:    "15 / 3;",
			expected: 5,
		},
		{
			desc:     "1 + 2 * 3",
			input:    "1 + 2 * 3;",
			expected: 7,
		},
		{
			desc:     "(1 + 2) * 3",
			input:    "(1 + 2) * 3;",
			expected: 9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			evaluated := eval(t, tt.input)
			testObject(t, tt.expected, evaluated)
		})
	}
}

// since the parse of function literal is tested in parser package,
// here we only test whether...
// - the function literal is evaluated as function object
// - the environment is properly included in function.
func TestEval_Function(t *testing.T) {
	tests := []struct {
		desc                string
		input               string
		expectedEnvVarName  []string
		expectedEnvVarValue []interface{}
	}{
		{
			desc:                "|| { 42; };",
			input:               "|| { 42; };",
			expectedEnvVarName:  []string{},
			expectedEnvVarValue: []interface{}{},
		},
		{
			desc:                "var c = 1; || { 42; };",
			input:               "var c = 1; || { 42; };",
			expectedEnvVarName:  []string{"c"},
			expectedEnvVarValue: []interface{}{1},
		},
		{
			desc:                "var a = 2; var b = 3; || { 42; };",
			input:               "var a = 2; var b = 3; || { 42; };",
			expectedEnvVarName:  []string{"a", "b"},
			expectedEnvVarValue: []interface{}{2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			evaluated := eval(t, tt.input)
			function, ok := evaluated.(*object.Function)
			if !ok {
				t.Errorf("unable to convert to function: %+v\n", evaluated)
			}
			for i, expectedName := range tt.expectedEnvVarName {
				expectedValue := tt.expectedEnvVarValue[i]
				actual := function.Env.Get(expectedName)
				if actual == nil {
					t.Errorf("undefined identifier: %s\n", expectedName)
				}
				testObject(t, expectedValue, actual)
			}
		})
	}
}

func TestEval_FunctionCall(t *testing.T) {
	tests := []struct {
		desc     string
		input    string
		expected interface{}
	}{
		{
			desc:     "var five=||{5;};five();",
			input:    "var five = || { 5; }; five();",
			expected: 5,
		},
		{
			desc:     "var two=2;var three=3;var five=||{two+three;};five();",
			input:    "var two = 2;var three = 3; var five = || { two + three; }; five();",
			expected: 5,
		},
		{
			desc:     "var identity=|x|{x;};identity(42);",
			input:    "var identity = |x| { x; }; identity(42);",
			expected: 42,
		},
		{
			desc:     "|x|{x;}(42);",
			input:    "|x| { x; }(42);",
			expected: 42,
		},
		{
			desc:     "var add=|x,y|{x+y;};add(7,8);",
			input:    "var add = |x, y| { x + y; }; add(7, 8);",
			expected: 15,
		},
		{
			desc:     "var five=||{return 5;};five();",
			input:    "var five = || { return 5; }; five();",
			expected: 5,
		},
		{
			desc:     "var five=||{4;return 5;};five();",
			input:    "var five = || { 4; return 5; }; five();",
			expected: 5,
		},
		{
			desc:     "var five=||{4;return 5;4;};five();",
			input:    "var five = || { 4; return 5; 4; }; five();",
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			evaluated := eval(t, tt.input)
			testObject(t, tt.expected, evaluated)
		})
	}
}

func TestEval_ArrowExpression(t *testing.T) {
	tests := []struct {
		desc     string
		input    string
		expected interface{}
	}{
		{
			desc:     "var identity=|x|{x;};42->identity();",
			input:    "var identity = |x| { x; }; 42 -> identity();",
			expected: 42,
		},
		{
			desc:     "42->|x|{x;}();",
			input:    "42 -> |x| { x; }();",
			expected: 42,
		},
		{
			desc:     "var add=|x,y|{x+y;};7->add(8);",
			input:    "var add = |x, y| { x + y; }; 7 -> add(8);",
			expected: 15,
		},
		{
			desc:     "var add=|x,y|{x+y;};var double=|x|{2*x;};7->double()->add(1);",
			input:    "var add = |x, y| { x + y; }; var double = |x| { 2 * x; }; 7 -> double() -> add(1);",
			expected: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			evaluated := eval(t, tt.input)
			testObject(t, tt.expected, evaluated)
		})
	}
}

func TestEval_ArrayLiteral(t *testing.T) {
	tests := []struct {
		desc        string
		input       string
		expectedLen int
		expected    []interface{}
	}{
		{
			desc:     "[]",
			input:    "[]",
			expected: []interface{}{},
		},
		{
			desc:     "[1]",
			input:    "[1]",
			expected: []interface{}{1},
		},
		{
			desc:     "[1,2,3]",
			input:    "[1, 2, 3]",
			expected: []interface{}{1, 2, 3},
		},
		{
			desc:     "var arr=[1,2,3];arr",
			input:    "var arr = [1, 2, 3]; arr",
			expected: []interface{}{1, 2, 3},
		},

	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			evaluated := eval(t, tt.input)
			actual, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("unable to convert to Array: %+v (%T)", evaluated, evaluated)
			}
			if len(actual.Elements) != len(tt.expected) {
				t.Errorf("array length wrong.\nwant=%d\ngot=%d\n", len(tt.expected), len(actual.Elements))
			}
			for i, expected := range tt.expected {
				testObject(t, expected, actual.Elements[i])
			}
		})
	}
}

func TestEval_VarStatement(t *testing.T) {
	tests := []struct {
		desc     string
		input    string
		expected int
	}{
		{
			desc:     "assignment",
			input:    "var a = 42; a;",
			expected: 42,
		},
		{
			desc:     "operation using identifier",
			input:    "var a = 42; a / 2;",
			expected: 21,
		},
		{
			desc:     "re-assignment",
			input:    "var a = 42; var b = a; b;",
			expected: 42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			evaluated := eval(t, tt.input)
			testObject(t, tt.expected, evaluated)
		})
	}
}

func TestEval_ReturnStatement(t *testing.T) {
	tests := []struct {
		desc     string
		input    string
		expected int
	}{
		{
			desc:     "return 42",
			input:    "return 42;",
			expected: 42,
		},
		{
			desc:     "1; return 42",
			input:    "42;",
			expected: 42,
		},
		{
			desc:     "1; return 42; 1",
			input:    "42;",
			expected: 42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			evaluated := eval(t, tt.input)
			testObject(t, tt.expected, evaluated)
		})
	}
}

func eval(t *testing.T, input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)

	program, err := p.ParseProgram()
	if err != nil {
		t.Errorf("parse error: %s\n", err.Error())
	}

	env := object.NewEnvironment()
	evaluated, err := Eval(program, env)
	if err != nil {
		t.Errorf("eval error: %s\n", err.Error())
	}

	return evaluated
}

func testObject(t *testing.T, expectedValue interface{}, actual object.Object) {
	switch expectedValue := expectedValue.(type) {
	case int:
		integer, ok := actual.(*object.Integer)
		if !ok {
			t.Fatalf("unable to convert to integer: %+v\n", actual)
		}
		if integer.Value != expectedValue {
			t.Errorf("integer value wrong:\nwant=%d\ngot=%d\n", expectedValue, integer.Value)
		}
	default:
		t.Errorf("unexpected type: %T", expectedValue)
	}
}
