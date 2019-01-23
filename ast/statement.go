package ast

type Statement interface {
	Node
	StatementNode()
}

type VarStatement struct {
	Identifier *Identifier
	Expression Expression
}

func (vs *VarStatement) String() string {
	return "let " + vs.Identifier.String() + " = " + vs.Expression.String()
}
func (vs *VarStatement) StatementNode() {}

type ReturnStatement struct {
	Expression Expression
}

func (rs *ReturnStatement) String() string { return "return " + rs.Expression.String() }
func (rs *ReturnStatement) StatementNode() {}

type ExpressionStatement struct {
	Expression Expression
}

func (es *ExpressionStatement) String() string { return es.Expression.String() }
func (es *ExpressionStatement) StatementNode() {}
