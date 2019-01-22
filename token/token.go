package token

type Type string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT = "IDENT"
	INT   = "INT"

	ASSIGN = "="
	PLUS   = "+"
	MINUS  = "-"
	ASTER  = "*"
	SLASH  = "/"

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN   = "("
	RPAREN   = ")"

	LET    = "LET"
	RETURN = "RETURN"
)

type Token struct {
	Type    Type
	Literal string
	Line    int
}
