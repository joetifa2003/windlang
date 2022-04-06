package evaluator

import (
	"bufio"
	"fmt"
	"os"

	"github.com/joetifa2003/windlang/ast"
)

var builtins = map[string]*BuiltinFn{
	"len": {
		ArgsCount: 1,
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
		ArgsCount: -1,
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
		ArgsCount: -1,
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
		ArgsCount: 1,
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, args ...Object) (Object, *Error) {
			switch arg := args[0].(type) {
			case *Integer:
				return &String{Value: fmt.Sprintf("%d", arg.Value)}, nil
			}

			return nil, evaluator.newError(node.Token, "argument to `string` not supported")
		},
	},
	"input": {
		ArgsCount: -1,
		ArgsTypes: []ObjectType{StringObj},
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, args ...Object) (Object, *Error) {
			if len(args) != 0 {
				prompt := args[0]
				fmt.Print(prompt.Inspect())
			}

			var input string
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			input = scanner.Text()

			return &String{Value: input}, nil
		},
	},
	"append": {
		ArgsCount: 2,
		ArgsTypes: []ObjectType{ArrayObj, Any},
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, args ...Object) (Object, *Error) {
			array := args[0].(*Array)
			value := args[1]

			array.Value = append(array.Value, value)

			return array, nil
		},
	},
}
