package ast

type Statement interface {
	Node
	StatementNode()
}

type LetStatement struct {
	Identifier *Identifier
	Expression Expression
}

func (ls *LetStatement) String() string {
	return "let " + ls.Identifier.String() + " = " + ls.Expression.String()
}
func (ls *LetStatement) StatementNode() {}

type ReturnStatement struct {
	Expression Expression
}

func (rs *ReturnStatement) String() string { return "return " + rs.Expression.String() }
func (rs *ReturnStatement) StatementNode() {}

type ExpressionStatement struct {
	Expression Expression
}

func (es *ExpressionStatement) String() string { return es.String() }
func (es *ExpressionStatement) StatementNode() {}
