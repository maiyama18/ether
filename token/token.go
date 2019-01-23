package token

type Type string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT   = "IDENT"
	INTEGER = "INTEGER"

	ASSIGN = "="
	PLUS   = "+"
	MINUS  = "-"
	ASTER  = "*"
	SLASH  = "/"

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"

	VAR    = "VAR"
	RETURN = "RETURN"
)

type Token struct {
	Type    Type
	Literal string
	Line    int
}

func TypeByLiteral(literal string) Type {
	switch literal {
	case "var":
		return VAR
	case "return":
		return RETURN
	default:
		return IDENT
	}
}
