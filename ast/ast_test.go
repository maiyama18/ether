package ast

import "testing"

func TestIdentifier_String(t *testing.T) {
	tests := []struct {
		desc     string
		name     string
		expected string
	}{
		{
			desc:     "one-char",
			name:     "a",
			expected: "a",
		},
		{
			desc:     "multi-char",
			name:     "foo",
			expected: "foo",
		},
		{
			desc:     "with-number",
			name:     "foo2bar3",
			expected: "foo2bar3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			identifier := &Identifier{Name: tt.name}
			testString(t, tt.expected, identifier)
		})
	}
}

func TestIntegerLiteral_String(t *testing.T) {
	tests := []struct {
		desc     string
		value    int
		expected string
	}{
		{
			desc:     "5",
			value:    5,
			expected: "5",
		},
		{
			desc:     "42",
			value:    42,
			expected: "42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			integerLiteral := &IntegerLiteral{Value: tt.value}
			testString(t, tt.expected, integerLiteral)
		})
	}
}

func TestPrefixExpression_String(t *testing.T) {
	tests := []struct {
		desc     string
		operator string
		right    Expression
		expected string
	}{
		{
			desc:     "-5",
			operator: "-",
			right:    &IntegerLiteral{Value: 5},
			expected: "(-5)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			prefixExpression := &PrefixExpression{Operator: tt.operator, Right: tt.right}
			testString(t, tt.expected, prefixExpression)
		})
	}
}

func TestInfixExpression_String(t *testing.T) {
	tests := []struct {
		desc     string
		operator string
		left     Expression
		right    Expression
		expected string
	}{
		{
			desc:     "3 + 4",
			operator: "+",
			left:     &IntegerLiteral{Value: 3},
			right:    &IntegerLiteral{Value: 4},
			expected: "(3 + 4)",
		},
		{
			desc:     "42-6",
			operator: "-",
			left:     &IntegerLiteral{Value: 42},
			right:    &IntegerLiteral{Value: 6},
			expected: "(42 - 6)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			infixExpression := &InfixExpression{Operator: tt.operator, Left: tt.left, Right: tt.right}
			testString(t, tt.expected, infixExpression)
		})
	}
}

func testString(t *testing.T, expected string, node Node) {
	if node.String() != expected {
		t.Errorf("string expression wrong: \nwant=%q\ngot=%q\n", expected, node.String())
	}
}
