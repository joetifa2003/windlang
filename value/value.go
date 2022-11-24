package value

import "fmt"

type ValueType int

const (
	VALUE_INT ValueType = iota
	VALUE_BOOL
	VALUE_NIL
)

type Value struct {
	VType ValueType
	IntV  int
	BoolV bool
}

func (v Value) String() string {
	switch v.VType {
	case VALUE_INT:
		return fmt.Sprint(v.IntV)
	case VALUE_BOOL:
		return fmt.Sprint(v.BoolV)
	case VALUE_NIL:
		return "nil"
	}

	panic("Unimplemented String() for value type")
}

func NewNilValue() Value {
	return Value{
		VType: VALUE_NIL,
	}
}

func NewIntValue(v int) Value {
	return Value{
		VType: VALUE_INT,
		IntV:  v,
	}
}

func NewBoolValue(v bool) Value {
	return Value{
		VType: VALUE_BOOL,
		BoolV: v,
	}
}
