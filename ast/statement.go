package ast

import "bytes"

type Statement interface {
	Node
	StatementNode()
}

type VarStatement struct {
	Identifier *Identifier
	Expression Expression
	line       int
}

func NewVarStatement(identifier *Identifier, expression Expression, line int) *VarStatement {
	return &VarStatement{Identifier: identifier, Expression: expression, line: line}
}
func (vs *VarStatement) Line() int { return vs.line }
func (vs *VarStatement) String() string {
	return "var " + vs.Identifier.String() + " = " + vs.Expression.String() + ";"
}
func (vs *VarStatement) StatementNode() {}

type ReturnStatement struct {
	Expression Expression
	line       int
}

func NewReturnStatement(expression Expression, line int) *ReturnStatement {
	return &ReturnStatement{Expression: expression, line: line}
}
func (rs *ReturnStatement) Line() int { return rs.line }
func (rs *ReturnStatement) String() string { return "return " + rs.Expression.String() + ";" }
func (rs *ReturnStatement) StatementNode() {}

type ExpressionStatement struct {
	Expression Expression
	line       int
}

func NewExpressionStatement(expression Expression, line int) *ExpressionStatement {
	return &ExpressionStatement{Expression: expression, line: line}
}
func (es *ExpressionStatement) Line() int { return es.line }
func (es *ExpressionStatement) String() string { return es.Expression.String() + ";" }
func (es *ExpressionStatement) StatementNode() {}

type BlockStatement struct {
	Statements []Statement
	line       int
}

func NewBlockStatement(statements []Statement, line int) *BlockStatement {
	return &BlockStatement{Statements: statements, line: line}
}
func (bs *BlockStatement) Line() int { return bs.line }
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
