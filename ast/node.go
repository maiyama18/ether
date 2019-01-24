package ast

import (
	"bytes"
)

type Node interface {
	Line() int
	String() string
}

type Program struct {
	Statements []Statement
}

func (p *Program) Line() int { return 1 }
func (p *Program) String() string {
	var out bytes.Buffer
	for _, statement := range p.Statements {
		out.WriteString(statement.String())
	}
	return out.String()
}
