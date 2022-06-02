package evaluator

import (
	"math"

	"github.com/joetifa2003/windlang/ast"
)

func stdLibMath() *Environment {
	return &Environment{
		Store: map[string]Object{
			"abs": &GoFunction{
				ArgsCount: 1,
				ArgsTypes: []ObjectType{FloatObj},
				Fn: func(evaluator *Evaluator, node *ast.CallExpression, args ...Object) (Object, *Error) {
					num := args[0].(*Float).Value

					return &Float{Value: math.Abs(num)}, nil
				},
			},
		},
	}
}
