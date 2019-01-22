package ast

import "github.com/muiscript/ether/token"

type Statement interface {
	Node
	StatementNode()
}

type LetStatement struct {
	token      token.Token
	Identifier *Identifier
	Expression Expression
}

func NewLetStatement(token token.Token, identifier *Identifier, expression Expression) *LetStatement {
	return &LetStatement{token: token, Identifier: identifier, Expression: expression}
}
func (ls *LetStatement) Token() token.Token { return ls.token }
func (ls *LetStatement) StatementNode()     {}

type ExpressionStatement struct {
	token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) Token() token.Token { return es.token }
func (es *ExpressionStatement) StatementNode()     {}

type ReturnStatement struct {
	token      token.Token
	Expression Expression
}

func NewReturnStatement(token token.Token, expression Expression) *ReturnStatement {
	return &ReturnStatement{token: token, Expression: expression}
}
func (rs *ReturnStatement) Token() token.Token { return rs.token }
func (rs *ReturnStatement) StatementNode()     {}
