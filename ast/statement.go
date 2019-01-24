package ast

import "bytes"

type Statement interface {
	Node
	StatementNode()
}

type VarStatement struct {
	Identifier *Identifier
	Expression Expression
	Line       int
}

func (vs *VarStatement) String() string {
	return "var " + vs.Identifier.String() + " = " + vs.Expression.String() + ";"
}
func (vs *VarStatement) StatementNode() {}

type ReturnStatement struct {
	Expression Expression
	Line       int
}

func (rs *ReturnStatement) String() string { return "return " + rs.Expression.String() + ";" }
func (rs *ReturnStatement) StatementNode() {}

type ExpressionStatement struct {
	Expression Expression
	Line       int
}

func (es *ExpressionStatement) String() string { return es.Expression.String() + ";" }
func (es *ExpressionStatement) StatementNode() {}

type BlockStatement struct {
	Statements []Statement
	Line       int
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{")
	for _, statement := range bs.Statements {
		out.WriteString(statement.String())
	}
	out.WriteString("}")
	return out.String()
}
func (bs *BlockStatement) StatementNode() {}
