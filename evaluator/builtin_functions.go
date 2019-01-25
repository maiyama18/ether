package evaluator

import (
	"fmt"
	"github.com/muiscript/ether/object"
)

var builtinFunctions = map[string]*object.BuiltinFunction{
	"puts": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg)
			}
		},
	},
}
