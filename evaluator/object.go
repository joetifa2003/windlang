package evaluator

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/joetifa2003/windlang/ast"
)

type ObjectType int

const (
	IntegerObj ObjectType = iota
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
	Clone() Object
}

type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Integer struct {
	Constant bool
	Value    int64
}

func (i *Integer) Type() ObjectType { return IntegerObj }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Clone() Object {
	c := *i
	return &c
}
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

type Float struct {
	Value float64
}

func (f *Float) Type() ObjectType { return FloatObj }
func (f *Float) Inspect() string  { return fmt.Sprintf("%f", f.Value) }
func (f *Float) Clone() Object {
	c := *f
	return &c
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BooleanObj }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Clone() Object {
	c := *b
	return &c
}
func (b *Boolean) HashKey() HashKey {
	var val uint64

	if b.Value {
		val = 1
	} else {
		val = 2
	}

	return HashKey{Type: b.Type(), Value: val}
}

type Nil struct{}

func (n *Nil) Type() ObjectType { return NilObj }
func (n *Nil) Inspect() string  { return "nil" }
func (n *Nil) Clone() Object    { return n }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return ReturnValueObj }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) Clone() Object {
	c := *rv
	return &c
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ErrorObj }
func (e *Error) Inspect() string  { return e.Message }
func (e *Error) Clone() Object {
	c := *e
	return &c
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
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
func (f *Function) Clone() Object {
	c := *f
	return &c
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return StringObj }
func (s *String) Inspect() string  { return s.Value }
func (s *String) Clone() Object {
	c := *s
	return &c
}
func (s *String) HashKey() HashKey {
	algo := fnv.New64a()
	algo.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: algo.Sum64()}
}

type BuiltinFunction func(evaluator *Evaluator, node *ast.CallExpression, args ...Object) (Object, *Error)

type BuiltinFn struct {
	ArgsCount int
	ArgsTypes []ObjectType
	Fn        BuiltinFunction
}

func (b *BuiltinFn) Type() ObjectType { return BuiltinObj }
func (b *BuiltinFn) Inspect() string  { return "builtin function" }
func (b *BuiltinFn) Clone() Object {
	c := *b
	return &c
}

type Array struct {
	Value []Object
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
func (a *Array) Clone() Object {
	c := *a
	return &c
}

type Hash struct {
	Pairs map[HashKey]Object
}

func (h *Hash) Type() ObjectType { return HashObj }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	out.WriteString("{")

	for _, value := range h.Pairs {
		out.WriteString(": ")
		out.WriteString(value.Inspect())
		out.WriteString(", ")
	}

	out.WriteString("}")

	return out.String()
}

func (h *Hash) Clone() Object {
	c := *h
	return &c
}

type IncludeObject struct {
	Value *Environment
}

func (i *IncludeObject) Type() ObjectType { return IncludeObj }
func (i *IncludeObject) Inspect() string {
	return "inclide_OBJ"
}
func (i *IncludeObject) Clone() Object {
	c := *i
	return &c
}
