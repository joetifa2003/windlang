package evaluator

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/joetifa2003/windlang/ast"
)

func stdLibReq() *Environment {
	return &Environment{
		Store: map[string]Object{
			"get": &GoFunction{
				ArgsCount: 1,
				ArgsTypes: []ObjectType{StringObj},
				Fn: func(evaluator *Evaluator, node *ast.CallExpression, args ...Object) (Object, *Error) {
					url := args[0].(*String)

					resp, err := http.Get(url.Value)
					if err != nil {
						return NIL, evaluator.newError(node.Token, "get request failed")
					}

					respBytes, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						return NIL, evaluator.newError(node.Token, "get request failed")
					}

					result := make(map[string]interface{})
					json.Unmarshal(respBytes, &result)

					objectResults := make(map[HashKey]Object)
					for k, v := range result {
						key := &String{Value: k}

						objectResults[key.HashKey()] = GetObjectFromInterFace(v)
					}

					return &Hash{Pairs: objectResults}, nil

				},
			},
			// "post": &GoFunction{
			// 	ArgsCount: 1,
			// 	ArgsTypes: []ObjectType{StringObj},
			// 	Fn: func(evaluator *Evaluator, node *ast.CallExpression, args ...Object) (Object, *Error) {
			// 		http.Post()
			// 	},
			// },
		},
	}
}
