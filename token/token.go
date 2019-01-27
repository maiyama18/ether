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

	ARROW = "ARROW"

	COMMA     = "COMMA"
	SEMICOLON = "SEMICOLON"

	LPAREN   = "LPAREN"
	RPAREN   = "RPAREN"
	LBRACE   = "LBRACE"
	RBRACE   = "RBRACE"
	LBRACKET = "LBRACKET"
	RBRACKET = "RBRACKET"
	BAR      = "BAR"

	VAR    = "VAR"
	RETURN = "RETURN"
	TRUE   = "TRUE"
	FALSE  = "FALSE"
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
	case "true":
		return TRUE
	case "false":
		return FALSE
	default:
		return IDENT
	}
}
