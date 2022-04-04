package evaluator

import (
	"strings"
	"wind-vm-go/ast"
)

func StdlibStrings() *Environment {
	return &Environment{
		Store: map[string]Object{
			"contains": &BuiltinFn{
				ArgsCount: 2,
				ArgsTypes: []ObjectType{StringObj, StringObj},
				Fn: func(evaluator *Evaluator, node *ast.CallExpression, args ...Object) (Object, *Error) {
					str1 := args[0].(*String)
					str2 := args[1].(*String)

					if strings.Contains(str1.Value, str2.Value) {
						return TRUE, nil
					} else {
						return FALSE, nil
					}
				},
			},
			"containsAny": &BuiltinFn{
				ArgsCount: 2,
				ArgsTypes: []ObjectType{StringObj, StringObj},
				Fn: func(evaluator *Evaluator, node *ast.CallExpression, args ...Object) (Object, *Error) {
					str1 := args[0].(*String)
					str2 := args[1].(*String)

					if strings.ContainsAny(str1.Value, str2.Value) {
						return TRUE, nil
					} else {
						return FALSE, nil
					}
				},
			},
			"count": &BuiltinFn{
				ArgsCount: 2,
				ArgsTypes: []ObjectType{StringObj, StringObj},
				Fn: func(evaluator *Evaluator, node *ast.CallExpression, args ...Object) (Object, *Error) {
					str1 := args[0].(*String)
					str2 := args[1].(*String)

					return &Integer{
						Value: int64(strings.Count(str1.Value, str2.Value)),
					}, nil
				},
			},
			"join": &BuiltinFn{
				ArgsCount: 2,
				ArgsTypes: []ObjectType{ArrayObj, StringObj},
				Fn: func(evaluator *Evaluator, node *ast.CallExpression, args ...Object) (Object, *Error) {
					arg1 := args[0].(*Array)
					arg2 := args[1].(*String)

					arrStrings := []string{}
					for _, v := range arg1.Value {
						arrStrings = append(arrStrings, v.Inspect())
					}

					return &String{
						Value: strings.Join(arrStrings, arg2.Value),
					}, nil
				},
			},
			"replaceN": &BuiltinFn{
				ArgsCount: 4,
				ArgsTypes: []ObjectType{StringObj, StringObj, StringObj, IntegerObj},
				Fn: func(evaluator *Evaluator, node *ast.CallExpression, args ...Object) (Object, *Error) {
					arg1 := args[0].(*String)
					arg2 := args[1].(*String)
					arg3 := args[2].(*String)
					arg4 := args[3].(*Integer)

					return &String{
						Value: strings.Replace(arg1.Value, arg2.Value, arg3.Value, int(arg4.Value)),
					}, nil
				},
			},
			"replace": &BuiltinFn{
				ArgsCount: 3,
				ArgsTypes: []ObjectType{StringObj, StringObj, StringObj},
				Fn: func(evaluator *Evaluator, node *ast.CallExpression, args ...Object) (Object, *Error) {
					arg1 := args[0].(*String)
					arg2 := args[1].(*String)
					arg3 := args[2].(*String)

					return &String{
						Value: strings.Replace(arg1.Value, arg2.Value, arg3.Value, 1),
					}, nil
				},
			},
			"replaceAll": &BuiltinFn{
				ArgsCount: 3,
				ArgsTypes: []ObjectType{StringObj, StringObj, StringObj},
				Fn: func(evaluator *Evaluator, node *ast.CallExpression, args ...Object) (Object, *Error) {
					arg1 := args[0].(*String)
					arg2 := args[1].(*String)
					arg3 := args[2].(*String)

					return &String{
						Value: strings.ReplaceAll(arg1.Value, arg2.Value, arg3.Value),
					}, nil
				},
			},
			"toLowerCase": &BuiltinFn{
				ArgsCount: 1,
				ArgsTypes: []ObjectType{StringObj},
				Fn: func(evaluator *Evaluator, node *ast.CallExpression, args ...Object) (Object, *Error) {
					arg1 := args[0].(*String)

					return &String{
						Value: strings.ToLower(arg1.Value),
					}, nil
				},
			},
			"toUpperCase": &BuiltinFn{
				ArgsCount: 1,
				ArgsTypes: []ObjectType{StringObj},
				Fn: func(evaluator *Evaluator, node *ast.CallExpression, args ...Object) (Object, *Error) {
					arg1 := args[0].(*String)

					return &String{
						Value: strings.ToUpper(arg1.Value),
					}, nil
				},
			},
			"index": &BuiltinFn{
				ArgsCount: 1,
				ArgsTypes: []ObjectType{StringObj},
				Fn: func(evaluator *Evaluator, node *ast.CallExpression, args ...Object) (Object, *Error) {
					arg1 := args[0].(*String)
					arg2 := args[1].(*String)

					return &Integer{
						Value: int64(strings.Index(arg1.Value, arg2.Value)),
					}, nil
				},
			},
			"lastIndex": &BuiltinFn{
				ArgsCount: 1,
				ArgsTypes: []ObjectType{StringObj},
				Fn: func(evaluator *Evaluator, node *ast.CallExpression, args ...Object) (Object, *Error) {
					arg1 := args[0].(*String)
					arg2 := args[1].(*String)

					return &Integer{
						Value: int64(strings.LastIndex(arg1.Value, arg2.Value)),
					}, nil
				},
			},
		},
	}
}
