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

	evaluated, err := Eval(program)
	if err != nil {
		t.Errorf("eval error: %s\n", err.Error())
	}

	return evaluated
}
