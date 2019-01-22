package ast

import (
	"github.com/muiscript/ether/token"
	"strconv"
)

type Expression interface {
	Node
	ExpressionNode()
}

type Identifier struct {
	token token.Token
	Name  string
}

func NewIdentifier(token token.Token) *Identifier {
	return &Identifier{token: token, Name: token.Literal}
}
func (i *Identifier) Token() token.Token { return i.token }
func (i *Identifier) ExpressionNode()    {}

type IntegerLiteral struct {
	token token.Token
	Value int
}

func NewIntegerLiteral(token token.Token) *IntegerLiteral {
	value, _ := strconv.Atoi(token.Literal)
	return &IntegerLiteral{token: token, Value: value}
}
func (il *IntegerLiteral) Token() token.Token { return il.token }
func (il *IntegerLiteral) ExpressionNode()    {}
