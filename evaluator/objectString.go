package evaluator

import (
	"hash/fnv"
	"strings"

	"github.com/joetifa2003/windlang/ast"
)

type String struct {
	Value string
}

func (s *String) GetFunction(name string) (*GoFunction, bool) {
	return GetFunctionFromObject(name, s, stringFunctions)
}
func (s *String) Type() ObjectType { return StringObj }
func (s *String) Inspect() string  { return s.Value }
func (s *String) HashKey() HashKey {
	algo := fnv.New64a()
	algo.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: algo.Sum64(), InspectValue: s.Inspect()}
}

var stringFunctions = map[string]OwnedFunction[*String]{
	"len": {
		ArgsCount: 0,
		ArgsTypes: []ObjectType{},
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, this *String, args ...Object) (Object, *Error) {
			return &Integer{
				Value: int64(len(this.Value)),
			}, nil
		},
	},
	"charAt": {
		ArgsCount: 1,
		ArgsTypes: []ObjectType{IntegerObj},
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, this *String, args ...Object) (Object, *Error) {
			index := args[0].(*Integer)

			if index.Value >= int64(len(this.Value)) {
				return NIL, nil
			}

			return &String{
				Value: string([]rune(this.Value)[index.Value]),
			}, nil
		},
	},
	"contains": {
		ArgsCount: 1,
		ArgsTypes: []ObjectType{StringObj},
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, this *String, args ...Object) (Object, *Error) {
			substr := args[0].(*String)

			if strings.Contains(this.Value, substr.Value) {
				return TRUE, nil
			} else {
				return FALSE, nil
			}
		},
	},
	"containsAny": {
		ArgsCount: 1,
		ArgsTypes: []ObjectType{StringObj},
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, this *String, args ...Object) (Object, *Error) {
			substr := args[0].(*String)

			if strings.ContainsAny(this.Value, substr.Value) {
				return TRUE, nil
			} else {
				return FALSE, nil
			}
		},
	},
	"count": {
		ArgsCount: 1,
		ArgsTypes: []ObjectType{StringObj},
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, this *String, args ...Object) (Object, *Error) {
			substr := args[0].(*String)

			return &Integer{
				Value: int64(strings.Count(this.Value, substr.Value)),
			}, nil
		},
	},
	"replace": {
		ArgsCount: 2,
		ArgsTypes: []ObjectType{StringObj, StringObj},
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, this *String, args ...Object) (Object, *Error) {
			old := args[0].(*String)
			new := args[1].(*String)

			return &String{
				Value: strings.Replace(this.Value, old.Value, new.Value, 1),
			}, nil
		},
	},
	"replaceN": {
		ArgsCount: 3,
		ArgsTypes: []ObjectType{StringObj, StringObj, IntegerObj},
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, this *String, args ...Object) (Object, *Error) {
			old := args[0].(*String)
			new := args[1].(*String)
			n := args[2].(*Integer)

			return &String{
				Value: strings.Replace(this.Value, old.Value, new.Value, int(n.Value)),
			}, nil
		},
	},
	"replaceAll": {
		ArgsCount: 2,
		ArgsTypes: []ObjectType{StringObj, StringObj},
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, this *String, args ...Object) (Object, *Error) {
			old := args[0].(*String)
			new := args[1].(*String)

			return &String{
				Value: strings.ReplaceAll(this.Value, old.Value, new.Value),
			}, nil
		},
	},
	"toLowerCase": {
		ArgsCount: 0,
		ArgsTypes: []ObjectType{},
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, this *String, args ...Object) (Object, *Error) {
			return &String{
				Value: strings.ToLower(this.Value),
			}, nil
		},
	},
	"toUpperCase": {
		ArgsCount: 0,
		ArgsTypes: []ObjectType{},
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, this *String, args ...Object) (Object, *Error) {
			return &String{
				Value: strings.ToUpper(this.Value),
			}, nil
		},
	},
	"indexOf": {
		ArgsCount: 1,
		ArgsTypes: []ObjectType{StringObj},
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, this *String, args ...Object) (Object, *Error) {
			substr := args[0].(*String)

			return &Integer{
				Value: int64(strings.Index(this.Value, substr.Value)),
			}, nil
		},
	},
	"lastIndexOf": {
		ArgsCount: 1,
		ArgsTypes: []ObjectType{StringObj},
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, this *String, args ...Object) (Object, *Error) {
			substr := args[0].(*String)

			return &Integer{
				Value: int64(strings.LastIndex(this.Value, substr.Value)),
			}, nil
		},
	},
	"changeAt": {
		ArgsCount: 2,
		ArgsTypes: []ObjectType{IntegerObj, StringObj},
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, this *String, args ...Object) (Object, *Error) {
			index := args[0].(*Integer)
			newValue := args[1].(*String)

			if index.Value >= int64(len(this.Value)) {
				return nil, evaluator.newError(node.Token, "index out of range: got %d max %d", index.Value, len(this.Value)-1)
			}

			if len(newValue.Value) > 1 {
				return nil, evaluator.newError(node.Token, "new value can be at most one character")
			}

			return &String{
				Value: this.Value[:index.Value] + newValue.Value + this.Value[index.Value+1:],
			}, nil
		},
	},
	"trim": {
		ArgsCount: 0,
		ArgsTypes: []ObjectType{},
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, this *String, args ...Object) (Object, *Error) {
			return &String{
				Value: strings.TrimSpace(this.Value),
			}, nil
		},
	},
	"split": {
		ArgsCount: 1,
		ArgsTypes: []ObjectType{StringObj},
		Fn: func(evaluator *Evaluator, node *ast.CallExpression, this *String, args ...Object) (Object, *Error) {
			seperator := args[0].(*String)

			strArr := strings.Split(this.Value, seperator.Value)
			objArr := []Object{}

			for _, str := range strArr {
				objArr = append(objArr, &String{Value: str})
			}

			return &Array{
				Value: objArr,
			}, nil
		},
	},
}
