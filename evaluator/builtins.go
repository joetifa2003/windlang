package evaluator

import (
	"fmt"
	"wind-vm-go/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s",
					args[0].Type())
			}
		},
	},
	"println": {
		Fn: func(args ...object.Object) object.Object {
			argsString := []interface{}{}
			for _, arg := range args {
				argsString = append(argsString, arg.Inspect())
			}

			fmt.Println(argsString...)

			return NULL
		},
	},
	"print": {
		Fn: func(args ...object.Object) object.Object {
			argsString := []interface{}{}
			for _, arg := range args {
				argsString = append(argsString, arg.Inspect())
			}

			fmt.Print(argsString...)

			return NULL
		},
	},
	"string": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}

			switch arg := args[0].(type) {
			case *object.Integer:
				return &object.String{Value: fmt.Sprintf("%d", arg.Value)}
			}

			return newError("argument to `string` not supported, got %s",
				args[0].Type())
		},
	},
}
