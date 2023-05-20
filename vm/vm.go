package vm

import (
	"fmt"

	"github.com/joetifa2003/windlang/opcode"
	"github.com/joetifa2003/windlang/value"
)

type Frame struct {
	Instructions []opcode.OpCode
	ip           int
	NumOfLocals  int
}

type VM struct {
	Stack     Stack
	Constants []value.Value
	Frames    []Frame
}

func NewVM(constants []value.Value, mainFrame Frame) VM {
	vm := VM{
		Stack:     NewStack(),
		Constants: constants,
		Frames:    []Frame{mainFrame},
	}
	vm.initCurFrame()
	return vm
}

func (v *VM) curFrame() *Frame {
	return &v.Frames[len(v.Frames)-1]
}

// initCurFrame initializes the stack to hold local variables
func (v *VM) initCurFrame() {
	for i := 0; i < v.curFrame().NumOfLocals; i++ {
		v.Stack.push(value.NewNilValue())
	}
}

func (v *VM) Interpret() {
	for v.curFrame().ip < len(v.curFrame().Instructions) {
		instructions := v.curFrame().Instructions
		ip := &v.curFrame().ip

		switch instructions[*ip] {
		case opcode.OP_CONST:
			*ip++
			value := v.Constants[instructions[*ip]]

			v.Stack.push(value)

		case opcode.OP_ADD:
			right := v.Stack.pop()
			left := v.Stack.peek()

			switch {
			case left.VType == value.VALUE_INT && right.VType == value.VALUE_INT:
				leftNumber := left.GetInt()
				rightNumber := right.GetInt()

				v.Stack.update(value.NewIntValue(leftNumber + rightNumber))
			}

		case opcode.OP_SUBTRACT:
			right := v.Stack.pop()
			left := v.Stack.pop()

			switch {
			case left.VType == value.VALUE_INT && right.VType == value.VALUE_INT:
				leftNumber := left.GetInt()
				rightNumber := right.GetInt()

				v.Stack.push(value.NewIntValue(leftNumber - rightNumber))
			}

		case opcode.OP_MULTIPLY:
			right := v.Stack.pop()
			left := v.Stack.pop()

			switch {
			case left.VType == value.VALUE_INT && right.VType == value.VALUE_INT:
				leftNumber := left.GetInt()
				rightNumber := right.GetInt()

				v.Stack.push(value.NewIntValue(leftNumber * rightNumber))
			}

		case opcode.OP_MODULO:
			right := v.Stack.pop()
			left := v.Stack.pop()

			switch {
			case left.VType == value.VALUE_INT && right.VType == value.VALUE_INT:
				leftNumber := left.GetInt()
				rightNumber := right.GetInt()

				v.Stack.push(value.NewIntValue(leftNumber % rightNumber))
			}

		case opcode.OP_DIVIDE:
			right := v.Stack.pop()
			left := v.Stack.pop()

			switch {
			case left.VType == value.VALUE_INT && right.VType == value.VALUE_INT:
				leftNumber := left.GetInt()
				rightNumber := right.GetInt()

				v.Stack.push(value.NewIntValue(leftNumber / rightNumber))
			}

		case opcode.OP_EQ:
			right := v.Stack.pop()
			left := v.Stack.pop()

			switch {
			case left.VType == value.VALUE_INT && right.VType == value.VALUE_INT:
				leftNumber := left.GetInt()
				rightNumber := right.GetInt()

				v.Stack.push(value.NewBoolValue(leftNumber == rightNumber))
			}

		case opcode.OP_LESSEQ:
			right := v.Stack.pop()
			left := v.Stack.pop()

			switch {
			case left.VType == value.VALUE_INT && right.VType == value.VALUE_INT:
				leftNumber := left.GetInt()
				rightNumber := right.GetInt()

				v.Stack.push(value.NewBoolValue(leftNumber <= rightNumber))
			}

		case opcode.OP_JUMP_FALSE:
			operand := v.Stack.pop()
			*ip++
			offset := int(instructions[*ip])

			if !isTruthy(operand) {
				*ip += offset

				continue
			}

		case opcode.OP_JUMP:
			*ip++
			offset := int(instructions[*ip])

			*ip += offset

			continue

		case opcode.OP_LET:
			value := v.Stack.pop()
			*ip++
			offset := int(instructions[*ip])
			v.Stack.Value[offset] = value

		case opcode.OP_SET:
			panic("Unimplemented")
		case opcode.OP_GET:
			*ip++
			offset := int(instructions[*ip])
			v.Stack.push(v.Stack.Value[offset])
		case opcode.OP_GET_GLOBAL:
			panic("Unimplemented")

		case opcode.OP_SET_GLOBAL:
			panic("Unimplemented")

		case opcode.OP_POP:
			if len(v.Stack.Value) != 0 {
				v.Stack.pop()
			}

		case opcode.OP_ECHO:
			operand := v.Stack.pop()

			fmt.Println(operand.String())

		case opcode.OP_ARRAY:
			*ip++
			n := int(instructions[*ip])
			values := make([]value.Value, n)

			for i := 0; i < n; i++ {
				values[i] = v.Stack.pop()
			}

			v.Stack.push(value.NewArrayValue(values))

		case opcode.OP_INC:
			// ip++
			// index := int(instructions[ip])
			// ip++
			// scopeIndex := int(instructions[ip])

		case opcode.OP_CALL:

		case opcode.OP_RET:
			return

		default:
			panic("Unimplemented OpCode " + fmt.Sprint(instructions[*ip]))
		}

		*ip++
	}
}

func isTruthy(input value.Value) bool {
	switch input.VType {
	case value.VALUE_BOOL:
		return input.GetBool()
	default:
		return true
	}
}
