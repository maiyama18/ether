package ast

import (
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
func (i *Identifier) Line() int                       { return i.line }
func (i *Identifier) String() string                  { return i.Name }
func (i *Identifier) ExpressionNode()                 {}

type IntegerLiteral struct {
	Value int
	line  int
}

func NewIntegerLiteral(value, line int) *IntegerLiteral {
	return &IntegerLiteral{Value: value, line: line}
}
func (il *IntegerLiteral) Line() int       { return il.line }
func (il *IntegerLiteral) String() string  { return strconv.Itoa(il.Value) }
func (il *IntegerLiteral) ExpressionNode() {}

type BooleanLiteral struct {
	Value bool
	line  int
}

func NewBooleanLiteral(value bool, line int) *BooleanLiteral {
	return &BooleanLiteral{Value: value, line: line}
}
func (bl *BooleanLiteral) Line() int       { return bl.line }
func (bl *BooleanLiteral) String() string  { return strconv.FormatBool(bl.Value) }
func (bl *BooleanLiteral) ExpressionNode() {}

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

	return "|" + strings.Join(paramStrs, ", ") + "| " + fl.Body.String()
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

	return fc.Function.String() + "(" + strings.Join(argStrs, ", ") + ")"
}
func (fc *FunctionCall) ExpressionNode() {}

type ArrayLiteral struct {
	Elements []Expression
	line     int
}

func NewArrayLiteral(elements []Expression, line int) *ArrayLiteral {
	return &ArrayLiteral{Elements: elements, line: line}
}

func (al *ArrayLiteral) Line() int { return al.line }
func (al *ArrayLiteral) String() string {
	var elemStrs []string
	for _, elem := range al.Elements {
		elemStrs = append(elemStrs, elem.String())
	}

	return "[" + strings.Join(elemStrs, ", ") + "]"
}
func (al *ArrayLiteral) ExpressionNode() {}

type IndexExpression struct {
	Array Expression
	Index Expression
	line  int
}

func NewIndexExpression(array Expression, index Expression, line int) *IndexExpression {
	return &IndexExpression{Array: array, Index: index, line: line}
}

func (ie *IndexExpression) Line() int { return ie.line }
func (ie *IndexExpression) String() string {
	return ie.Array.String() + "[" + ie.Index.String() + "]"
}
func (ie *IndexExpression) ExpressionNode() {}
