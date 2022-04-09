package evaluator

import (
	"bytes"
	"strings"

	"github.com/joetifa2003/windlang/ast"
)

type Array struct {
	Value []Object
}

func (a *Array) GetFunction(name string) (*GoFunction, bool) {
	return GetFunctionFromObject(name, a, arrayFunctions)
}
func (a *Array) Type() ObjectType { return ArrayObj }
func (a *Array) Inspect() string {
	var out bytes.Buffer

	out.WriteString("[")
	for _, obj := range a.Value {
		out.WriteString(obj.Inspect())
		out.WriteString(",")
	}
	out.WriteString("]")

	return out.String()
}

var arrayFunctions = map[string]OwnedFunction[*Array]{
	"len": {
		ArgsCount: 0,
		ArgsTypes: []ObjectType{},
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, this *Array, args ...Object) (Object, *Error) {
			return &Integer{
				Value: int64(len(this.Value)),
			}, nil
		},
	},
	"join": {
		ArgsCount: 1,
		ArgsTypes: []ObjectType{StringObj},
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, this *Array, args ...Object) (Object, *Error) {
			strArr := []string{}
			for _, obj := range this.Value {
				strArr = append(strArr, obj.Inspect())
			}

			return &String{
				Value: strings.Join(strArr, args[0].(*String).Value),
			}, nil
		},
	},
	"filter": {
		ArgsCount: 1,
		ArgsTypes: []ObjectType{FunctionObj},
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, this *Array, args ...Object) (Object, *Error) {
			fn := args[0].(*Function)

			filtered := []Object{}
			for _, obj := range this.Value {
				result, err := evaluator.applyFunction(node, fn, []Object{obj})
				if err != nil {
					return nil, err
				}

				if result == TRUE {
					filtered = append(filtered, obj)
				}
			}

			return &Array{
				Value: filtered,
			}, nil
		},
	},
	"map": {
		ArgsCount: 1,
		ArgsTypes: []ObjectType{FunctionObj},
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, this *Array, args ...Object) (Object, *Error) {
			fn := args[0].(*Function)

			mapped := []Object{}
			for _, obj := range this.Value {
				result, err := evaluator.applyFunction(node, fn, []Object{obj})
				if err != nil {
					return nil, err
				}

				mapped = append(mapped, result)
			}

			return &Array{
				Value: mapped,
			}, nil
		},
	},
	"reduce": {
		ArgsCount: 2,
		ArgsTypes: []ObjectType{FunctionObj, Any},
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, this *Array, args ...Object) (Object, *Error) {
			fn := args[0].(*Function)
			initial := args[1]

			accumulator := initial
			for _, obj := range this.Value {
				result, err := evaluator.applyFunction(node, fn, []Object{accumulator, obj})
				if err != nil {
					return nil, err
				}

				accumulator = result
			}

			return accumulator, nil
		},
	},
	"push": {
		ArgsCount: 1,
		ArgsTypes: []ObjectType{Any},
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, this *Array, args ...Object) (Object, *Error) {
			this.Value = append(this.Value, args[0])
			return this, nil
		},
	},
	"pop": {
		ArgsCount: 0,
		ArgsTypes: []ObjectType{},
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, this *Array, args ...Object) (Object, *Error) {
			last := this.Value[len(this.Value)-1]
			this.Value = this.Value[:len(this.Value)-1]
			return last, nil
		},
	},
	"contains": {
		ArgsCount: 1,
		ArgsTypes: []ObjectType{FunctionObj},
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, this *Array, args ...Object) (Object, *Error) {
			fn := args[0].(*Function)

			for _, obj := range this.Value {
				result, err := evaluator.applyFunction(node, fn, []Object{obj})
				if err != nil {
					return nil, err
				}

				if result == TRUE {
					return TRUE, nil
				}
			}

			return FALSE, nil
		},
	},
	"count": {
		ArgsCount: 1,
		ArgsTypes: []ObjectType{FunctionObj},
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, this *Array, args ...Object) (Object, *Error) {
			fn := args[0].(*Function)

			count := 0
			for _, obj := range this.Value {
				result, err := evaluator.applyFunction(node, fn, []Object{obj})
				if err != nil {
					return nil, err
				}

				if result == TRUE {
					count++
				}
			}

			return &Integer{
				Value: int64(count),
			}, nil
		},
	},
	"clone": {
		ArgsCount: 0,
		ArgsTypes: []ObjectType{},
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, this *Array, args ...Object) (Object, *Error) {
			newValue := make([]Object, len(this.Value))

			for i, obj := range this.Value {
				newValue[i] = obj
			}

			return &Array{
				Value: newValue,
			}, nil
		},
	},
	"removeAt": {
		ArgsCount: 1,
		ArgsTypes: []ObjectType{IntegerObj},
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, this *Array, args ...Object) (Object, *Error) {
			index := args[0].(*Integer).Value

			if index < 0 || index >= int64(len(this.Value)) {
				return nil, evaluator.newError(node.Token, "index %d out of bounds", index)
			}

			newValue := []Object{}
			for i, v := range this.Value {
				if i != int(index) {
					newValue = append(newValue, v)
				}
			}
			removedValue := this.Value[index]
			this.Value = newValue

			return removedValue, nil
		},
	},
}
