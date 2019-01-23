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
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			identifier := &Identifier{Name: tt.name}
			testString(t, tt.expected, identifier)
		})
	}
}

func testString(t *testing.T, expected string, node Node) {
	if node.String() != expected {
		t.Errorf("string expression wrong: \nwant=%q\ngot=%q\n", expected, node.String())
	}
}

// func TestNode_String(t *testing.T) {
// 	tests := []struct {
// 		desc     string
// 		input    string
// 		expected string
// 	}{
// 		{
// 			desc:     "identifier",
// 			input:    "foo",
// 			expected: "foo",
// 		},
// 		{
// 			desc:     "integer literal",
// 			input:    "42",
// 			expected: "42",
// 		},
// 		{
// 			desc:     "prefix expression",
// 			input:    "-42",
// 			expected: "(-42)",
// 		},
// 		{
// 			desc:     "infix expression",
// 			input:    "2 + 3",
// 			expected: "(2 + 3)",
// 		},
// 		{
// 			desc:     "let statement",
// 			input:    "2 + 3",
// 			expected: "(2 + 3)",
// 		},
// 	}
// }
