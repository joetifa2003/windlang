package evaluator

import (
	"fmt"
	"math"
)

var builtins = map[string]*Builtin{
	"len": {
		Fn: func(evaluator *Evaluator, args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}

			switch arg := args[0].(type) {
			case *String:
				return &Integer{Value: int64(len(arg.Value))}
			case *Array:
				return &Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported)")
			}
		},
	},
	"println": {
		Fn: func(evaluator *Evaluator, args ...Object) Object {
			argsString := []interface{}{}
			for _, arg := range args {
				argsString = append(argsString, arg.Inspect())
			}

			fmt.Println(argsString...)

			return NIL
		},
	},
	"print": {
		Fn: func(evaluator *Evaluator, args ...Object) Object {
			argsString := []interface{}{}
			for _, arg := range args {
				argsString = append(argsString, arg.Inspect())
			}

			fmt.Print(argsString...)

			return NIL
		},
	},
	"string": {
		Fn: func(evaluator *Evaluator, args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}

			switch arg := args[0].(type) {
			case *Integer:
				return &String{Value: fmt.Sprintf("%d", arg.Value)}
			}

			return newError("argument to `string` not supported")
		},
	},
	"sqrt": {
		Fn: func(evaluator *Evaluator, args ...Object) Object {
			switch arg := args[0].(type) {
			case *Integer:
				value := math.Sqrt(float64(arg.Value))

				return &Float{Value: value}

			case *Float:
				return &Float{Value: math.Sqrt(arg.Value)}
			}

			return newError("argument to `sqrt` not supported")
		},
	},
	"input": {
		Fn: func(evaluator *Evaluator, args ...Object) Object {
			if len(args) != 0 {
				prompt := args[0]
				fmt.Print(prompt.Inspect())
			}

			var input string
			fmt.Scanln(&input)

			return &String{Value: input}
		},
	},
}
