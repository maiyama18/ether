package evaluator

import (
	"fmt"
	"github.com/muiscript/ether/object"
)

var builtinFunctions = map[string]*object.BuiltinFunction{
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
				return nil, &EvalError{line: 1, msg: fmt.Sprintf("number of arguments for len wrong: want=%d\n", 1)}
			}
			array, ok := args[0].(*object.Array)
			if !ok {
				return nil, &EvalError{line: 1, msg: fmt.Sprintf("argument type for len wrong: want=%T\ngot=%T\n", &object.Array{}, array)}
			}

			return &object.Integer{Value: len(array.Elements)}, nil
		},
	},
	// "map": {
	// 	Fn: func(args ...object.Object) (object.Object, error) {
	// 		if len(args) != 2 {
	// 			return nil, &EvalError{line: 1, msg: fmt.Sprintf("number of arguments for map wrong: want=%d\n")}
	// 		}
	// 	}
	// }
}
