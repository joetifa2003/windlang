package evaluator

import (
	"fmt"
	"math"
	"wind-vm-go/ast"
)

var builtins = map[string]*Builtin{
	"len": {
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, args ...Object) (Object, *Error) {
			if len(args) != 1 {
				return nil, evaluator.newError(node.Token, "wrong number of arguments. got=%d, want=1",
					len(args))
			}

			switch arg := args[0].(type) {
			case *String:
				return &Integer{Value: int64(len(arg.Value))}, nil
			case *Array:
				return &Integer{Value: int64(len(arg.Value))}, nil
			default:
				return nil, evaluator.newError(node.Token, "argument to `len` not supported)")
			}
		},
	},
	"println": {
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, args ...Object) (Object, *Error) {
			argsString := []interface{}{}
			for _, arg := range args {
				argsString = append(argsString, arg.Inspect())
			}

			fmt.Println(argsString...)

			return NIL, nil
		},
	},
	"print": {
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, args ...Object) (Object, *Error) {
			argsString := []interface{}{}
			for _, arg := range args {
				argsString = append(argsString, arg.Inspect())
			}

			fmt.Print(argsString...)

			return NIL, nil
		},
	},
	"string": {
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, args ...Object) (Object, *Error) {
			if len(args) != 1 {
				return nil, evaluator.newError(node.Token, "wrong number of arguments. got=%d, want=1",
					len(args))
			}

			switch arg := args[0].(type) {
			case *Integer:
				return &String{Value: fmt.Sprintf("%d", arg.Value)}, nil
			}

			return nil, evaluator.newError(node.Token, "argument to `string` not supported")
		},
	},
	"sqrt": {
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, args ...Object) (Object, *Error) {
			switch arg := args[0].(type) {
			case *Integer:
				value := math.Sqrt(float64(arg.Value))

				return &Float{Value: value}, nil

			case *Float:
				return &Float{Value: math.Sqrt(arg.Value)}, nil
			}

			return nil, evaluator.newError(node.Token, "argument to `sqrt` not supported")
		},
	},
	"input": {
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, args ...Object) (Object, *Error) {
			if len(args) != 0 {
				prompt := args[0]
				fmt.Print(prompt.Inspect())
			}

			var input string
			fmt.Scanln(&input)

			return &String{Value: input}, nil
		},
	},
}
