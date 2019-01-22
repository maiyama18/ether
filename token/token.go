package token

type Type string

const (
	ILLEGAL = "ILLEGAL"
	EOF = "EOF"

	IDENT = "IDENT"
	INT = "INT"

	ASSIGN = "ASSIGN"
	PLUS = "PLUS"

	COMMA = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"
	LBRACKET = "["
	RBRACKET = "]"

	FN = "FN"
	LET = "LET"
	RETURN = "RETURN"
)

type Token struct {
	Type Type
	Literal string
}
