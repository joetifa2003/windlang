package evaluator

import (
	"bufio"
	"fmt"
	"os"

	"github.com/joetifa2003/windlang/ast"
)

var builtins = map[string]*GoFunction{
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
			case *Float:
				return &String{Value: fmt.Sprintf("%f", arg.Value)}, nil
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
}
