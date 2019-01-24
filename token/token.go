package token

type Type string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT   = "IDENT"
	INTEGER = "INTEGER"

	ASSIGN = "ASSIGN"
	PLUS   = "PLUS"
	MINUS  = "MINUS"
	ASTER  = "ASTER"
	SLASH  = "SLASH"

	COMMA     = "COMMA"
	SEMICOLON = "SEMICOLON"

	LPAREN = "LPAREN"
	RPAREN = "RPAREN"
	LBRACE = "LBRACE"
	RBRACE = "RBRACE"
	BAR    = "BAR"

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
