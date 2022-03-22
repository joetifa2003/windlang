package object

import (
	"bytes"
	"fmt"
	"strings"
	"wind-vm-go/ast"
)

type ObjectType string

const (
	INTEGER_OBJ      ObjectType = "INTEGER"
	BOOLEAN_OBJ      ObjectType = "BOOLEAN"
	NULL_OBJ         ObjectType = "NULL"
	RETURN_VALUE_OBJ ObjectType = "RETURN_VALUE"
	ERROR_OBJ        ObjectType = "ERROR"
	FUNCTION_OBJ     ObjectType = "FUNCTION"
	STRING_OBJ       ObjectType = "STRING"
	BUILTIN_OBJ      ObjectType = "BUILTIN"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i Integer) Type() ObjectType { return INTEGER_OBJ }
func (i Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Float struct {
	Value float64
}

func (f Float) Type() ObjectType { return INTEGER_OBJ }
func (f Float) Inspect() string  { return fmt.Sprintf("%f", f.Value) }

type Boolean struct {
	Value bool
}

func (b Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type Null struct{}

func (n Null) Type() ObjectType { return NULL_OBJ }
func (n Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (rv ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Error struct {
	Message string
}

func (e Error) Type() ObjectType { return ERROR_OBJ }
func (e Error) Inspect() string  { return "ERROR: " + e.Message }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f Function) Type() ObjectType { return FUNCTION_OBJ }
func (f Function) Inspect() string {
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

type String struct {
	Value string
}

func (s String) Type() ObjectType { return STRING_OBJ }
func (s String) Inspect() string  { return s.Value }

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b Builtin) Inspect() string  { return "builtin function" }
