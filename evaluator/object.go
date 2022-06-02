package evaluator

import (
	"bytes"
	"fmt"
	"math"
	"strings"

	"github.com/joetifa2003/windlang/ast"
)

type ObjectType int

const (
	Any ObjectType = iota
	IntegerObj
	FloatObj
	BooleanObj
	NilObj
	ReturnValueObj
	ErrorObj
	FunctionObj
	StringObj
	BuiltinObj
	ArrayObj
	HashObj
	IncludeObj
)

func (ot ObjectType) String() string {
	switch ot {
	case IntegerObj:
		return "INTEGER"
	case FloatObj:
		return "FLOAT"
	case BooleanObj:
		return "BOOLEAN"
	case NilObj:
		return "NIL"
	case ReturnValueObj:
		return "RETURN_VALUE"
	case ErrorObj:
		return "ERROR"
	case FunctionObj:
		return "FUNCTION"
	case StringObj:
		return "STRING"
	case BuiltinObj:
		return "BUILTIN"
	case ArrayObj:
		return "ARRAY"
	case HashObj:
		return "HASH"
	case IncludeObj:
		return "INCLUDE"
	default:
		return "UNKNOWN"
	}
}

type Object interface {
	Type() ObjectType
	Inspect() string
}

type OwnedFunction[T Object] struct {
	ArgsCount int
	ArgsTypes []ObjectType
	Fn        func(evaluator *Evaluator, node *ast.CallExpression, this T, args ...Object) (Object, *Error)
}

type ObjectWithFunctions interface {
	GetFunction(name string) (*GoFunction, bool)
}

func GetFunctionFromObject[T Object](name string, object T, functions map[string]OwnedFunction[T]) (*GoFunction, bool) {
	if fn, ok := functions[name]; ok {
		return &GoFunction{
			ArgsCount: fn.ArgsCount,
			ArgsTypes: fn.ArgsTypes,
			Fn: func(evaluator *Evaluator, node *ast.CallExpression, args ...Object) (Object, *Error) {
				return fn.Fn(evaluator, node, object, args...)
			},
		}, true
	}

	return nil, false
}

type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type         ObjectType
	Value        uint64
	InspectValue string
}

func (hk *HashKey) Inspect() string {
	return hk.InspectValue
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return IntegerObj }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

type Float struct {
	Value float64
}

func (f *Float) Type() ObjectType { return FloatObj }
func (f *Float) Inspect() string  { return fmt.Sprintf("%f", f.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BooleanObj }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) HashKey() HashKey {
	var val uint64

	if b.Value {
		val = 1
	} else {
		val = 2
	}

	return HashKey{Type: b.Type(), Value: val, InspectValue: b.Inspect()}
}

type Nil struct{}

func (n *Nil) Type() ObjectType { return NilObj }
func (n *Nil) Inspect() string  { return "nil" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return ReturnValueObj }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ErrorObj }
func (e *Error) Inspect() string  { return e.Message }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
	This       Object
}

func (f *Function) Type() ObjectType { return FunctionObj }
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}

type BuiltinFunction func(evaluator *Evaluator, node *ast.CallExpression, args ...Object) (Object, *Error)

type GoFunction struct {
	ArgsCount int
	ArgsTypes []ObjectType
	Fn        BuiltinFunction
}

func (b *GoFunction) Type() ObjectType { return BuiltinObj }
func (b *GoFunction) Inspect() string  { return "builtin function" }

type Hash struct {
	Pairs map[HashKey]Object
}

func (h *Hash) Type() ObjectType { return HashObj }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	out.WriteString("{")

	for key, value := range h.Pairs {
		out.WriteString(key.Inspect())
		out.WriteString(": ")
		out.WriteString(value.Inspect())
		out.WriteString(", ")
	}

	out.WriteString("}")

	return out.String()
}

type IncludeObject struct {
	Value *Environment
}

func (i *IncludeObject) Type() ObjectType { return IncludeObj }
func (i *IncludeObject) Inspect() string {
	return "include_OBJ"
}

func GetObjectFromInterFace(v interface{}) Object {
	switch v := v.(type) {
	case float64:
		if v == math.Trunc(v) {
			return &Integer{Value: int64(v)}
		} else {
			return &Float{Value: v}
		}

	case string:
		return &String{Value: v}

	case bool:
		if v {
			return TRUE
		} else {
			return FALSE
		}

	case []interface{}:
		res := make([]Object, len(v))
		for i, val := range v {
			res[i] = GetObjectFromInterFace(val)
		}

		return &Array{Value: res}
	}

	return NIL
}
