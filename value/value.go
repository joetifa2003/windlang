package value

import (
	"fmt"
	"unsafe"
)

type ValueType int

const (
	VALUE_INT ValueType = iota
	VALUE_BOOL
	VALUE_NIL
	VALUE_ARRAY
	VALUE_OBJECT
)

type Value struct {
	VType        ValueType
	IntV         int
	BoolV        bool
	FloatV       float32
	nonPrimitive unsafe.Pointer
}

type nonPrimitiveValue struct {
	ArrayV  []Value
	ObjectV map[string]Value
	StringV string
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

func NewArrayValue(v []Value) Value {
	return Value{
		VType:        VALUE_ARRAY,
		nonPrimitive: unsafe.Pointer(&nonPrimitiveValue{ArrayV: v}),
	}
}

func (v Value) GetArray() []Value {
	return (*nonPrimitiveValue)(v.nonPrimitive).ArrayV
}

func (v Value) String() string {
	switch v.VType {
	case VALUE_INT:
		return fmt.Sprint(v.IntV)
	case VALUE_BOOL:
		return fmt.Sprint(v.BoolV)
	case VALUE_NIL:
		return "nil"
	case VALUE_ARRAY:
		return fmt.Sprint(v.GetArray())
	}

	panic("Unimplemented String() for value type")
}
