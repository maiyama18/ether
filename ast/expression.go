package ast

import (
	"bytes"
	"strconv"
	"strings"
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

type FunctionLiteral struct {
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) String() string {
	var paramStrs []string
	for _, param := range fl.Parameters {
		paramStrs = append(paramStrs, param.String())
	}

	var out bytes.Buffer
	out.WriteString("|")
	out.WriteString(strings.Join(paramStrs, ", "))
	out.WriteString("| ")
	out.WriteString(fl.Body.String())

	return out.String()
}
func (fl *FunctionLiteral) ExpressionNode() {}

type FunctionCall struct {
	Function Expression // FunctionLiteral or Identifier
	Arguments []Expression
}

func (fc *FunctionCall) String() string {
	var argStrs []string
	for _, arg := range fc.Arguments {
		argStrs = append(argStrs, arg.String())
	}

	var out bytes.Buffer
	out.WriteString(fc.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(argStrs, ", "))
	out.WriteString(")")

	return out.String()
}
func (fc *FunctionCall) ExpressionNode() {}
