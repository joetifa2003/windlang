package vm

import (
	"fmt"

	"github.com/joetifa2003/windlang/opcode"
	"github.com/joetifa2003/windlang/value"
)

type VM struct {
	Stack     Stack
	EnvStack  EnvironmentStack
	Constants []value.Value
}

func NewVM(constants []value.Value) VM {
	stack := NewStack()
	envStack := NewEnvironmentStack()

	return VM{
		Stack:     stack,
		EnvStack:  envStack,
		Constants: constants,
	}
}

func (v *VM) Interpret(instructions []opcode.OpCode) {
	ip := 0
	for ip < len(instructions) {
		switch instructions[ip] {
		case opcode.OP_CONST:
			ip++
			value := v.Constants[instructions[ip]]

			v.Stack.push(value)

		case opcode.OP_ADD:
			right := v.Stack.pop()
			left := v.Stack.pop()

			switch {
			case left.VType == value.VALUE_INT && right.VType == value.VALUE_INT:
				leftNumber := left.GetInt()
				rightNumber := right.GetInt()

				v.Stack.push(value.NewIntValue(leftNumber + rightNumber))
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
			ip++
			offset := int(instructions[ip])

			if !isTruthy(operand) {
				ip += offset

				continue
			}

		case opcode.OP_JUMP:
			ip++
			offset := int(instructions[ip])

			ip += offset

			continue

		case opcode.OP_BLOCK:
			ip++
			varCount := int(instructions[ip])

			v.EnvStack.push(NewEnvironment(varCount))

		case opcode.OP_END_BLOCK:
			v.EnvStack.pop()

		case opcode.OP_LET:
			value := v.Stack.pop()

			ip++
			index := int(instructions[ip])

			v.EnvStack.let(index, value)

		case opcode.OP_SET:
			value := v.Stack.pop()

			ip++
			index := int(instructions[ip])
			ip++
			scopeIndex := int(instructions[ip])

			newVal := v.EnvStack.set(scopeIndex, index, value)
			v.Stack.push(newVal)

		case opcode.OP_GET:
			ip++
			index := int(instructions[ip])
			ip++
			scopeIndex := int(instructions[ip])

			value := v.EnvStack.get(scopeIndex, index)
			v.Stack.push(value)

		case opcode.OP_POP:
			if len(v.Stack.Value) != 0 {
				v.Stack.pop()
			}

		case opcode.OP_ECHO:
			operand := v.Stack.pop()

			fmt.Println(operand.String())

		case opcode.OP_ARRAY:
			ip++
			n := int(instructions[ip])
			values := make([]value.Value, n)

			for i := 0; i < n; i++ {
				values[i] = v.Stack.pop()
			}

			v.Stack.push(value.NewArrayValue(values))

		case opcode.OP_INC:
			ip++
			index := int(instructions[ip])
			ip++
			scopeIndex := int(instructions[ip])

			ok := v.EnvStack.increment(scopeIndex, index)
			if !ok {
				panic("")
			}

		default:
			panic("Unimplemented OpCode " + fmt.Sprint(instructions[ip]))
		}

		ip++
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
