package evaluator

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"
	"wind-vm-go/ast"
)

type ObjectType int

const (
	INTEGER_OBJ ObjectType = iota
	FLOAT_OBJ
	BOOLEAN_OBJ
	NIL_OBJ
	RETURN_VALUE_OBJ
	ERROR_OBJ
	FUNCTION_OBJ
	STRING_OBJ
	BUILTIN_OBJ
	ARRAY_OBJ
	HASH_OBJ
)

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
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
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

func (f *Float) Type() ObjectType { return FLOAT_OBJ }
func (f *Float) Inspect() string  { return fmt.Sprintf("%f", f.Value) }
func (f *Float) Clone() Object {
	c := *f
	return &c
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
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

func (n *Nil) Type() ObjectType { return NIL_OBJ }
func (n *Nil) Inspect() string  { return "nil" }
func (n *Nil) Clone() Object    { return n }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) Clone() Object {
	c := *rv
	return &c
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
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

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
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

func (s *String) Type() ObjectType { return STRING_OBJ }
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

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }
func (b *Builtin) Clone() Object {
	c := *b
	return &c
}

type Array struct {
	Value []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
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

func (h *Hash) Type() ObjectType { return HASH_OBJ }
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
func HashMapFromEnv(env *Environment) *Hash {
	hash := make(map[HashKey]Object)

	for key, value := range env.Store {
		stringKey := String{Value: key}
		hash[stringKey.HashKey()] = value
	}

	return &Hash{Pairs: hash}
}
