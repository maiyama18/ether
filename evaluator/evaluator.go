package evaluator

import (
	"fmt"
	"github.com/muiscript/ether/ast"
	"github.com/muiscript/ether/object"
)

// TODO: add line to EvalError
func Eval(node ast.Node, env *object.Environment) (object.Object, error) {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.VarStatement:
		return evalVarStatement(node, env)
	case *ast.ExpressionStatement:
		return evalExpressionStatement(node, env)
	default:
		return nil, &EvalError{line: 1, msg: fmt.Sprintf("unable to eval node: %+v (%T)", node, node)}
	}
}

func evalProgram(program *ast.Program, env *object.Environment) (object.Object, error) {
	var evaluated object.Object
	for _, statement := range program.Statements {
		var err error
		evaluated, err = Eval(statement, env)
		if err != nil {
			return nil, err
		}
	}
	return evaluated, nil
}

func evalBlockStatement(blockStatement *ast.BlockStatement, env *object.Environment) (object.Object, error) {
	var evaluated object.Object
	for _, statement := range blockStatement.Statements {
		var err error
		evaluated, err = Eval(statement, env)
		if err != nil {
			return nil, err
		}
	}
	return evaluated, nil
}

func evalVarStatement(varStatement *ast.VarStatement, env *object.Environment) (object.Object, error) {
	value, err := evalExpression(varStatement.Expression, env)
	if err != nil {
		return nil, err
	}
	env.Set(varStatement.Identifier.Name, value)
	return nil, nil
}

func evalExpressionStatement(expressionStatement *ast.ExpressionStatement, env *object.Environment) (object.Object, error) {
	expression := expressionStatement.Expression
	return evalExpression(expression, env)
}

func evalExpression(expression ast.Expression, env *object.Environment) (object.Object, error) {
	switch expression := expression.(type) {
	case *ast.IntegerLiteral:
		return &object.Integer{Value: expression.Value}, nil
	case *ast.Identifier:
		value := env.Get(expression.Name)
		if value == nil {
			return nil, &EvalError{line: 1, msg: fmt.Sprintf("undefined identifier: %q", expression.Name)}
		}
		return value, nil
	case *ast.PrefixExpression:
		return evalPrefixExpression(expression, env)
	case *ast.InfixExpression:
		return evalInfixExpression(expression, env)
	case *ast.FunctionLiteral:
		return evalFunctionLiteral(expression, env)
	case *ast.FunctionCall:
		return evalFunctionCall(expression, env)
	default:
		return nil, &EvalError{line: 1, msg: fmt.Sprintf("unable to eval expression: %+v (%T)", expression, expression)}
	}
}

func evalPrefixExpression(prefixExpression *ast.PrefixExpression, env *object.Environment) (object.Object, error) {
	right, err := evalExpression(prefixExpression.Right, env)
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

func evalInfixExpression(infixExpression *ast.InfixExpression, env *object.Environment) (object.Object, error) {
	left, err := evalExpression(infixExpression.Left, env)
	if err != nil {
		return nil, err
	}
	right, err := evalExpression(infixExpression.Right, env)
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

func evalFunctionLiteral(functionLiteral *ast.FunctionLiteral, env *object.Environment) (object.Object, error) {
	return &object.Function{Parameters: functionLiteral.Parameters, Body: functionLiteral.Body, Env: env}, nil
}

func evalFunctionCall(functionCall *ast.FunctionCall, env *object.Environment) (object.Object, error) {
	functionExp, err := evalExpression(functionCall.Function, env)
	if err != nil {
		return nil, err
	}
	function, ok := functionExp.(*object.Function)
	if !ok {
		return nil, &EvalError{line: 1, msg: fmt.Sprintf("unable to convert to function: %+v (%T)", functionExp, functionExp)}
	}

	if len(functionCall.Arguments) != len(function.Parameters) {
		return nil, &EvalError{line: 1, msg: fmt.Sprintf("number of arguments for %+v wrong:\nwant=%d\ngot=%d\n", function, len(function.Parameters), len(functionCall.Arguments))}
	}

	var evaluatedArgs []object.Object
	for _, arg := range functionCall.Arguments {
		evaluatedArg, err := evalExpression(arg, env);
		if err != nil {
			return nil, err
		}

		evaluatedArgs = append(evaluatedArgs, evaluatedArg)
	}

	enclosedEnv := object.NewEnclosedEnvironment(function.Env)
	for i, evaluatedArg := range evaluatedArgs {
		ident := function.Parameters[i]
		enclosedEnv.Set(ident.Name, evaluatedArg)
	}

	return Eval(function.Body, enclosedEnv)
}
