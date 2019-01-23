package evaluator

import (
	"github.com/muiscript/ether/lexer"
	"github.com/muiscript/ether/object"
	"github.com/muiscript/ether/parser"
	"testing"
)

func TestEval_IntegerExpression(t *testing.T) {
	tests := []struct {
		desc     string
		input    string
		expected int
	}{
		{
			desc: "42",
			input: "42;",
			expected: 42,
		},
		{
			desc: "-42",
			input: "-42;",
			expected: -42,
		},
		{
			desc: "-(-42)",
			input: "-(-42);",
			expected: 42,
		},
		{
			desc: "15 + 3",
			input: "15 + 3;",
			expected: 18,
		},
		{
			desc: "15 - 3",
			input: "15 - 3;",
			expected: 12,
		},
		{
			desc: "15 * 3",
			input: "15 * 3;",
			expected: 45,
		},
		{
			desc: "15 / 3",
			input: "15 / 3;",
			expected: 5,
		},
		{
			desc: "1 + 2 * 3",
			input: "1 + 2 * 3;",
			expected: 7,
		},
		{
			desc: "(1 + 2) * 3",
			input: "(1 + 2) * 3;",
			expected: 9,
		},
	}

	for _, tt := range tests {
		evaluated := eval(t, tt.input)
		integer, ok := evaluated.(*object.Integer)
		if !ok {
			t.Errorf("unable to convert to integer: %+v\n", evaluated)
		}
		if integer.Value != tt.expected {
			t.Errorf("integer value wrong.\nwant=%d\ngot=%d\n", tt.expected, integer.Value)
		}
	}
}

func TestEval_VarStatement(t *testing.T) {
	tests := []struct {
		desc     string
		input    string
		expected int
	}{
		{
			desc: "assignment",
			input: "var a = 42; a;",
			expected: 42,
		},
		{
			desc: "operation using identifier",
			input: "var a = 42; a / 2;",
			expected: 21,
		},
		{
			desc: "re-assignment",
			input: "var a = 42; var b = a; b;",
			expected: 42,
		},
	}

	for _, tt := range tests {
		evaluated := eval(t, tt.input)
		integer, ok := evaluated.(*object.Integer)
		if !ok {
			t.Errorf("unable to convert to integer: %+v\n", evaluated)
		}
		if integer.Value != tt.expected {
			t.Errorf("integer value wrong.\nwant=%d\ngot=%d\n", tt.expected, integer.Value)
		}
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
