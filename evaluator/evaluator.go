package evaluator

import (
	"fmt"
	"github.com/muiscript/ether/ast"
	"github.com/muiscript/ether/object"
)

// TODO: add line to EvalError
func Eval(node ast.Node) (object.Object, error) {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node)
	case *ast.ExpressionStatement:
		return evalExpressionStatement(node)
	default:
		return nil, &EvalError{line: 1, msg: fmt.Sprintf("unable to eval node: %+v (%T)", node, node)}
	}
}

func evalProgram(program *ast.Program) (object.Object, error) {
	var evaluated object.Object
	for _, statement := range program.Statements {
		var err error
		evaluated, err = Eval(statement)
		if err != nil {
			return nil, err
		}
	}
	return evaluated, nil
}

func evalExpressionStatement(expressionStatement *ast.ExpressionStatement) (object.Object, error) {
	expression := expressionStatement.Expression
	return evalExpression(expression)
}

func evalExpression(expression ast.Expression) (object.Object, error) {
	switch expression := expression.(type) {
	case *ast.IntegerLiteral:
		return &object.Integer{Value: expression.Value}, nil
	case *ast.PrefixExpression:
		return evalPrefixExpression(expression)
	case *ast.InfixExpression:
		return evalInfixExpression(expression)
	default:
		return nil, &EvalError{line: 1, msg: fmt.Sprintf("unable to eval expression: %+v (%T)", expression, expression)}
	}
}

func evalPrefixExpression(prefixExpression *ast.PrefixExpression) (object.Object, error) {
	right, err := evalExpression(prefixExpression.Right)
	if err != nil {
		return nil, err
	}

	switch right := right.(type) {
	case *object.Integer:
		switch prefixExpression.Operator {
		case "-":
			return &object.Integer{Value: -right.Value}, nil
		default:
			return nil, &EvalError{line: 1, msg: fmt.Sprintf("unknown prefix operator for integer: %q", prefixExpression.Operator)}
		}
	default:
		return nil, &EvalError{line: 1, msg: fmt.Sprintf("invalid type for prefix expression: %+v (%T)", right, right)}
	}
}

func evalInfixExpression(infixExpression *ast.InfixExpression) (object.Object, error) {
	left, err := evalExpression(infixExpression.Left)
	if err != nil {
		return nil, err
	}
	right, err := evalExpression(infixExpression.Right)
	if err != nil {
		return nil, err
	}

	if left.Type() != right.Type() {
		return nil, &EvalError{line: 1, msg: fmt.Sprintf("type mismatch in infix expression: %+v %s %+v", left, infixExpression.Operator, right)}
	}
	switch left := left.(type) {
	case *object.Integer:
		right := right.(*object.Integer)
		switch infixExpression.Operator {
		case "+":
			return &object.Integer{Value: left.Value + right.Value}, nil
		case "-":
			return &object.Integer{Value: left.Value - right.Value}, nil
		case "*":
			return &object.Integer{Value: left.Value * right.Value}, nil
		case "/":
			return &object.Integer{Value: left.Value / right.Value}, nil
		default:
			return nil, &EvalError{line: 1, msg: fmt.Sprintf("unknown infix operator for integer: %q", infixExpression.Operator)}
		}
	default:
		return nil, &EvalError{line: 1, msg: fmt.Sprintf("invalid type for infix expression: %+v (%T)", left, left)}
	}
}
