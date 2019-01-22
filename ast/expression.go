package ast

import (
	"strconv"
)

type Expression interface {
	Node
	ExpressionNode()
}

type Identifier struct {
	Name string
}

func (i *Identifier) String() string  { return i.Name }
func (i *Identifier) ExpressionNode() {}

type IntegerLiteral struct {
	Value int
}

func (il *IntegerLiteral) String() string {
	return strconv.Itoa(il.Value)
}
func (il *IntegerLiteral) ExpressionNode() {}
