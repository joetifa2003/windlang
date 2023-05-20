package value

import (
	"fmt"
	"unsafe"

	"github.com/joetifa2003/windlang/opcode"
)

type ValueType int

const (
	VALUE_NIL ValueType = iota
	VALUE_INT
	VALUE_BOOL
	VALUE_ARRAY
	VALUE_OBJECT
	VALUE_FN
)

type Value struct {
	VType         ValueType
	primitiveData [8]byte // int64, float64, bool
	nonPrimitive  *NonPrimitiveData
}

type NonPrimitiveData struct {
	ArrayV []Value
	FnV    FnValue
}

type FnValue struct {
	Instructions []opcode.OpCode
	VarCount     int
}

func NewNilValue() Value {
	return Value{
		VType: VALUE_NIL,
	}
}

func NewIntValue(v int) Value {
	value := Value{
		VType:         VALUE_INT,
		primitiveData: [8]byte{},
	}

	*(*int)(unsafe.Pointer(&value.primitiveData[0])) = v

	return value
}

func NewBoolValue(v bool) Value {
	value := Value{
		VType:         VALUE_BOOL,
		primitiveData: [8]byte{},
	}

	*(*bool)(unsafe.Pointer(&value.primitiveData[0])) = v

	return value
}

func NewArrayValue(v []Value) Value {
	return Value{
		VType: VALUE_ARRAY,
		nonPrimitive: &NonPrimitiveData{
			ArrayV: v,
		},
	}
}

func NewFnValue(instructions []opcode.OpCode, varCount int) Value {
	return Value{
		VType: VALUE_FN,
		nonPrimitive: &NonPrimitiveData{
			FnV: FnValue{
				Instructions: instructions,
				VarCount:     varCount,
			},
		},
	}
}

func (v *Value) GetArray() []Value {
	return v.nonPrimitive.ArrayV
}

func (v *Value) GetFn() FnValue {
	return v.nonPrimitive.FnV
}

func (v *Value) GetInt() int {
	return *(*int)(unsafe.Pointer(&v.primitiveData[0]))
}

func (v *Value) GetIntPtr() *int {
	return (*int)(unsafe.Pointer(&v.primitiveData[0]))
}

func (v *Value) GetBool() bool {
	return *(*bool)(unsafe.Pointer(&v.primitiveData[0]))
}

func (v *Value) String() string {
	switch v.VType {
	case VALUE_INT:
		return fmt.Sprint(v.GetInt())
	case VALUE_BOOL:
		return fmt.Sprint(v.GetBool())
	case VALUE_NIL:
		return "nil"
	case VALUE_ARRAY:
		return fmt.Sprint(v.GetArray())
	}

	panic("Unimplemented String() for value type")
}
