package ast

import "testing"

func TestIdentifier_String(t *testing.T) {
	tests := []struct {
		desc     string
		name     string
		expected string
	}{
		{
			desc:     "a",
			name:     "a",
			expected: "a",
		},
		{
			desc:     "foo",
			name:     "foo",
			expected: "foo",
		},
		{
			desc:     "foo2bar3",
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

func TestVarStatement_String(t *testing.T) {
	tests := []struct {
		desc       string
		name       string
		expression Expression
		expected   string
	}{
		{
			desc:       "var foo = 42;",
			name:       "foo",
			expression: &IntegerLiteral{Value: 42},
			expected:   "var foo = 42;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			varStatement := &VarStatement{Identifier: &Identifier{Name: tt.name}, Expression: tt.expression}
			testString(t, tt.expected, varStatement)
		})
	}
}

func TestReturnStatement_String(t *testing.T) {
	tests := []struct {
		desc       string
		expression Expression
		expected   string
	}{
		{
			desc:       "return 42;",
			expression: &IntegerLiteral{Value: 42},
			expected:   "return 42;",
		},
		{
			desc:       "return foo;",
			expression: &Identifier{Name: "foo"},
			expected:   "return foo;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			returnStatement := &ReturnStatement{Expression: tt.expression}
			testString(t, tt.expected, returnStatement)
		})
	}
}

func TestExpressionStatement_String(t *testing.T) {
	tests := []struct {
		desc       string
		expression Expression
		expected   string
	}{
		{
			desc:       "42;",
			expression: &IntegerLiteral{Value: 42},
			expected:   "42;",
		},
		{
			desc:       "foo;",
			expression: &Identifier{Name: "foo"},
			expected:   "foo;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			expressionStatement := &ExpressionStatement{Expression: tt.expression}
			testString(t, tt.expected, expressionStatement)
		})
	}
}

func testString(t *testing.T, expected string, node Node) {
	if node.String() != expected {
		t.Errorf("string expression wrong: \nwant=%q\ngot=%q\n", expected, node.String())
	}
}
