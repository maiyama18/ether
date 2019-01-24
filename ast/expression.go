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
	line int
}

func NewIdentifier(name string, line int) *Identifier { return &Identifier{Name: name, line: line} }
func (i *Identifier) Line() int { return i.line }
func (i *Identifier) String() string  { return i.Name }
func (i *Identifier) ExpressionNode() {}

type IntegerLiteral struct {
	Value int
	line  int
}

func NewIntegerLiteral(value, line int) *IntegerLiteral { return &IntegerLiteral{Value: value, line: line} }
func (il *IntegerLiteral) Line() int { return il.line }
func (il *IntegerLiteral) String() string { return strconv.Itoa(il.Value) }
func (il *IntegerLiteral) ExpressionNode() {}

type PrefixExpression struct {
	Operator string
	Right    Expression
	line     int
}

func NewPrefixExpression(operator string, right Expression, line int) *PrefixExpression {
	return &PrefixExpression{Operator: operator, Right: right, line: line}
}
func (pe *PrefixExpression) Line() int { return pe.line }
func (pe *PrefixExpression) String() string {
	return "(" + pe.Operator + pe.Right.String() + ")"
}
func (pe *PrefixExpression) ExpressionNode() {}

type InfixExpression struct {
	Operator string
	Left     Expression
	Right    Expression
	line     int
}

func NewInfixExpression(operator string, left, right Expression, line int) *InfixExpression {
	return &InfixExpression{Operator: operator, Left: left, Right: right, line: line}
}
func (ie *InfixExpression) Line() int { return ie.line }
func (ie *InfixExpression) String() string {
	return "(" + ie.Left.String() + " " + ie.Operator + " " + ie.Right.String() + ")"
}
func (ie *InfixExpression) ExpressionNode() {}

type FunctionLiteral struct {
	Parameters []*Identifier
	Body       *BlockStatement
	line       int
}

func NewFunctionLiteral(parameters []*Identifier, body *BlockStatement, line int) *FunctionLiteral {
	return &FunctionLiteral{Parameters: parameters, Body: body, line: line}
}
func (fl *FunctionLiteral) Line() int { return fl.line }
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
	Function  Expression // FunctionLiteral or Identifier
	Arguments []Expression
	line      int
}

func NewFunctionCall(function Expression, arguments []Expression, line int) *FunctionCall {
	return &FunctionCall{Function: function, Arguments: arguments, line: line}
}
func (fc *FunctionCall) Line() int { return fc.line }
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
