package ast

import "github.com/muiscript/ether/token"

type Expression interface {
	Node
	ExpressionNode()
}

type Identifier struct {
	token token.Token
	Name  string
}

func (i *Identifier) Token() token.Token { return i.token }
func (i *Identifier) ExpressionNode()    {}

type IntegerLiteral struct {
	token token.Token
	Value int
}

func (il *IntegerLiteral) Token() token.Token { return il.token }
func (il *IntegerLiteral) ExpressionNode()    {}
