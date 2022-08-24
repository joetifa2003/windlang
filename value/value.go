package value

import "fmt"

type ValueType int

const (
	VALUE_INT ValueType = iota
	VALUE_BOOL
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
