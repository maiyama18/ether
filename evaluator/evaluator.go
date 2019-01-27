package evaluator

import (
	"fmt"
	"github.com/muiscript/ether/ast"
	"github.com/muiscript/ether/object"
)

var builtinFunctions map[string]*object.BuiltinFunction

func init() {
	builtinFunctions = map[string]*object.BuiltinFunction{
		"puts": {
			Fn: func(args ...object.Object) (object.Object, error) {
				for _, arg := range args {
					fmt.Println(arg)
				}
				return nil, nil
			},
		},
		"len": {
			Fn: func(args ...object.Object) (object.Object, error) {
				if len(args) != 1 {
					return nil, &EvalError{line: 1, msg: fmt.Sprintf("number of arguments for len wrong: want=%d got=%d\n", 1, len(args))}
				}
				array, ok := args[0].(*object.Array)
				if !ok {
					return nil, &EvalError{line: 1, msg: fmt.Sprintf("argument type for len wrong: want=%T\ngot=%T\n", &object.Array{}, array)}
				}

				return &object.Integer{Value: len(array.Elements)}, nil
			},
		},
		"map": {
			Fn: func(args ...object.Object) (object.Object, error) {
				if len(args) != 2 {
					return nil, &EvalError{line: 1, msg: fmt.Sprintf("number of arguments for map wrong: want=%d got=%d\n", 2, len(args))}
				}
				array, ok := args[0].(*object.Array)
				if !ok {
					return nil, &EvalError{line: 1, msg: fmt.Sprintf("first argument type for map wrong: want=%T\ngot=%T\n", &object.Array{}, array)}
				}
				function, ok := args[1].(*object.Function)
				if !ok {
					return nil, &EvalError{line: 1, msg: fmt.Sprintf("second argument type for map wrong: want=%T\ngot=%T\n", &object.Function{}, function)}
				}
				if len(function.Parameters) != 1 {
					return nil, &EvalError{line: 1, msg: fmt.Sprintf("number of parameters of map function wrong: want=%T\ngot=%T\n", 1, len(function.Parameters))}
				}

				var convertedElems []object.Object
				for _, elem := range array.Elements {
					enclosedEnv := object.NewEnclosedEnvironment(function.Env)
					enclosedEnv.Set(function.Parameters[0].Name, elem)

					evaluated, err := Eval(function.Body, enclosedEnv)
					if err != nil {
						return nil, err
					}
					convertedElems = append(convertedElems, evaluated)
				}

				return &object.Array{Elements: convertedElems}, nil
			},
		},
	}
}

func Eval(node ast.Node, env *object.Environment) (object.Object, error) {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.VarStatement:
		return evalVarStatement(node, env)
	case *ast.ReturnStatement:
		return evalReturnStatement(node, env)
	case *ast.ExpressionStatement:
		return evalExpressionStatement(node, env)
	default:
		return nil, &EvalError{line: node.Line(), msg: fmt.Sprintf("unable to eval node: %+v (%T)", node, node)}
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
		if returnValue, ok := evaluated.(*object.ReturnValue); ok {
			return returnValue.Value, nil
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
		if returnValue, ok := evaluated.(*object.ReturnValue); ok {
			return returnValue, nil
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

func evalReturnStatement(returnStatement *ast.ReturnStatement, env *object.Environment) (object.Object, error) {
	value, err := evalExpression(returnStatement.Expression, env)
	if err != nil {
		return nil, err
	}
	return &object.ReturnValue{Value: value}, nil
}

func evalExpressionStatement(expressionStatement *ast.ExpressionStatement, env *object.Environment) (object.Object, error) {
	return evalExpression(expressionStatement.Expression, env)
}

func evalExpression(expression ast.Expression, env *object.Environment) (object.Object, error) {
	switch expression := expression.(type) {
	case *ast.IntegerLiteral:
		return &object.Integer{Value: expression.Value}, nil
	case *ast.Identifier:
		value := env.Get(expression.Name)
		if value == nil {
			if builtin, ok := builtinFunctions[expression.Name]; ok {
				return builtin, nil
			} else {
				return nil, &EvalError{line: expression.Line(), msg: fmt.Sprintf("undefined identifier: %q", expression.Name)}
			}
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
	case *ast.ArrayLiteral:
		return evalArrayLiteral(expression, env)
	case *ast.IndexExpression:
		return evalIndexExpression(expression, env)
	default:
		return nil, &EvalError{line: expression.Line(), msg: fmt.Sprintf("unable to eval expression: %+v (%T)", expression, expression)}
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
			return nil, &EvalError{line: prefixExpression.Line(), msg: fmt.Sprintf("unknown prefix operator for integer: %q", prefixExpression.Operator)}
		}
	default:
		return nil, &EvalError{line: prefixExpression.Right.Line(), msg: fmt.Sprintf("invalid type for prefix expression: %+v (%T)", right, right)}
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
		return nil, &EvalError{line: infixExpression.Line(), msg: fmt.Sprintf("type mismatch in infix expression: %+v %s %+v", left, infixExpression.Operator, right)}
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
			return nil, &EvalError{line: infixExpression.Line(), msg: fmt.Sprintf("unknown infix operator for integer: %q", infixExpression.Operator)}
		}
	default:
		return nil, &EvalError{line: infixExpression.Line(), msg: fmt.Sprintf("invalid type for infix expression: %+v (%T)", left, left)}
	}
}

func evalFunctionLiteral(functionLiteral *ast.FunctionLiteral, env *object.Environment) (object.Object, error) {
	return &object.Function{Parameters: functionLiteral.Parameters, Body: functionLiteral.Body, Env: env}, nil
}

func evalFunctionCall(functionCall *ast.FunctionCall, env *object.Environment) (object.Object, error) {
	var evaluatedArgs []object.Object
	for _, arg := range functionCall.Arguments {
		evaluatedArg, err := evalExpression(arg, env)
		if err != nil {
			return nil, err
		}

		evaluatedArgs = append(evaluatedArgs, evaluatedArg)
	}

	function, err := evalExpression(functionCall.Function, env)
	if err != nil {
		return nil, err
	}

	switch function := function.(type) {
	case *object.Function:
		if len(functionCall.Arguments) != len(function.Parameters) {
			return nil, &EvalError{line: functionCall.Arguments[0].Line(), msg: fmt.Sprintf("number of arguments for %+v wrong:\nwant=%d\ngot=%d\n", function, len(function.Parameters), len(functionCall.Arguments))}
		}

		enclosedEnv := object.NewEnclosedEnvironment(function.Env)
		for i, evaluatedArg := range evaluatedArgs {
			ident := function.Parameters[i]
			enclosedEnv.Set(ident.Name, evaluatedArg)
		}

		evaluated, err := Eval(function.Body, enclosedEnv)
		if err != nil {
			return nil, err
		}
		return unwrapReturnValue(evaluated), nil
	case *object.BuiltinFunction:
		return function.Fn(evaluatedArgs...)
	default:
		return nil, &EvalError{line: functionCall.Line(), msg: fmt.Sprintf("unable to convert to function: %+v (%T)", function, function)}
	}
}

func evalArrayLiteral(arrayLiteral *ast.ArrayLiteral, env *object.Environment) (object.Object, error) {
	var evaluatedElements []object.Object
	for _, elem := range arrayLiteral.Elements {
		evaluatedElem, err := evalExpression(elem, env)
		if err != nil {
			return nil, err
		}
		evaluatedElements = append(evaluatedElements, evaluatedElem)
	}
	return &object.Array{Elements: evaluatedElements}, nil
}

func evalIndexExpression(indexExpression *ast.IndexExpression, env *object.Environment) (object.Object, error) {
	evaluatedArray, err := evalExpression(indexExpression.Array, env)
	if err != nil {
		return nil, err
	}
	array, ok := evaluatedArray.(*object.Array)
	if !ok {
		return nil, &EvalError{line: indexExpression.Line(), msg: fmt.Sprintf("unable to convert to array: %+v (%T)", evaluatedArray, evaluatedArray)}
	}

	evaluatedIndex, err := evalExpression(indexExpression.Index, env)
	if err != nil {
		return nil, err
	}
	index, ok := evaluatedIndex.(*object.Integer)
	if !ok {
		return nil, &EvalError{line: indexExpression.Line(), msg: fmt.Sprintf("unable to convert to integer: %+v (%T)", evaluatedIndex, evaluatedIndex)}
	}

	if index.Value < 0 || len(array.Elements) <= index.Value {
		return nil, &EvalError{line: indexExpression.Line(), msg: fmt.Sprintf("index out of range: %v[%d]\n", array, index.Value)}
	}

	return array.Elements[index.Value], nil
}

func unwrapReturnValue(obj object.Object) object.Object {
	switch obj := obj.(type) {
	case *object.ReturnValue:
		return obj.Value
	default:
		return obj
	}
}
