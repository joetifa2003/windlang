package value

import "fmt"

type ValueType byte

const (
	VALUE_INT ValueType = iota
	VALUE_BOOL
	VALUE_NIL
)

type Value interface {
	ValueType() ValueType
	String() string
}

type IntegerValue struct{ Value int64 }

func (IntegerValue) ValueType() ValueType { return VALUE_INT }
func (i IntegerValue) String() string     { return fmt.Sprint(i.Value) }

type BoolValue struct{ Value bool }

func (BoolValue) ValueType() ValueType { return VALUE_BOOL }
func (b BoolValue) String() string     { return fmt.Sprint(b.Value) }

type NilValue struct{}

func (NilValue) ValueType() ValueType { return VALUE_NIL }
func (NilValue) String() string       { return "nil" }