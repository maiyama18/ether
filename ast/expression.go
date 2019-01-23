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

type PrefixExpression struct {
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) String() string {
	return "(" + pe.Operator + pe.Right.String() + ")"
}
func (pe *PrefixExpression) ExpressionNode() {}

type InfixExpression struct {
	Operator string
	Left     Expression
	Right    Expression
}

func (ie *InfixExpression) String() string {
	return "(" + ie.Left.String() + " " + ie.Operator + " " + ie.Right.String() + ")"
}
func (ie *InfixExpression) ExpressionNode() {}
