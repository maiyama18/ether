package evaluator

import (
	"fmt"
	"github.com/muiscript/ether/ast"
	"github.com/muiscript/ether/object"
)

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
	default:
		return nil, &EvalError{line: 1, msg: fmt.Sprintf("unable to eval expression: %+v (%T)", expression, expression)}
	}
}
