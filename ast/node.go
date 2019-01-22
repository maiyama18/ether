package ast

import "github.com/muiscript/ether/token"

type Node interface {
	Token() token.Token
}

type Program struct {
	Statements []Statement
}

func (p *Program) Token() token.Token {
	if len(p.Statements) == 0 {
		return token.Token{Type: token.EOF, Literal: "", Line: 1}
	}
	return p.Statements[0].Token()
}
