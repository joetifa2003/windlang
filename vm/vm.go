package vm

import (
	"github.com/joetifa2003/windlang/opcode"
	"github.com/joetifa2003/windlang/value"
)

type VM struct {
	Stack    Stack
	EnvStack EnvironmentStack
}

func NewVM() VM {
	envStack := NewEnvironmentStack()

	return VM{
		Stack: Stack{
			Value: make([]value.Value, 2048),
		},
		EnvStack: envStack,
	}
}

func (v *VM) Interpret(instructions []opcode.OpCode) {
	ip := 0
	for ip < len(instructions) {
		curInstruction := instructions[ip]
		switch instruction := curInstruction.(type) {
		case opcode.ConstOpCode:
			v.Stack.push(instruction.Value)

		case opcode.AddOpCode:
			right := v.Stack.pop()
			left := v.Stack.pop()

			switch {
			case left.ValueType() == value.VALUE_INT && right.ValueType() == value.VALUE_INT:
				leftInt := left.(value.IntegerValue)
				rightInt := right.(value.IntegerValue)

				v.Stack.push(value.IntegerValue{Value: leftInt.Value + rightInt.Value})
			}

		case opcode.SubtractOpCode:
			right := v.Stack.pop()
			left := v.Stack.pop()

			switch {
			case left.ValueType() == value.VALUE_INT && right.ValueType() == value.VALUE_INT:
				leftInt := left.(value.IntegerValue)
				rightInt := right.(value.IntegerValue)

				v.Stack.push(value.IntegerValue{Value: leftInt.Value - rightInt.Value})
			}

		case opcode.MultiplyOpCode:
			right := v.Stack.pop()
			left := v.Stack.pop()

			switch {
			case left.ValueType() == value.VALUE_INT && right.ValueType() == value.VALUE_INT:
				leftInt := left.(value.IntegerValue)
				rightInt := right.(value.IntegerValue)

				v.Stack.push(value.IntegerValue{Value: leftInt.Value * rightInt.Value})
			}

		case opcode.ModuloOpCode:
			right := v.Stack.pop()
			left := v.Stack.pop()

			switch {
			case left.ValueType() == value.VALUE_INT && right.ValueType() == value.VALUE_INT:
				leftInt := left.(value.IntegerValue)
				rightInt := right.(value.IntegerValue)

				v.Stack.push(value.IntegerValue{Value: leftInt.Value % rightInt.Value})
			}

		case opcode.DivideOpCode:
			right := v.Stack.pop()
			left := v.Stack.pop()

			switch {
			case left.ValueType() == value.VALUE_INT && right.ValueType() == value.VALUE_INT:
				leftInt := left.(value.IntegerValue)
				rightInt := right.(value.IntegerValue)

				v.Stack.push(value.IntegerValue{Value: leftInt.Value / rightInt.Value})
			}

		case opcode.EqualOpCode:
			right := v.Stack.pop()
			left := v.Stack.pop()

			switch {
			case left.ValueType() == value.VALUE_INT && right.ValueType() == value.VALUE_INT:
				leftInt := left.(value.IntegerValue)
				rightInt := right.(value.IntegerValue)

				v.Stack.push(value.BoolValue{Value: leftInt.Value == rightInt.Value})
			}

		case opcode.LessThanEqOpCode:
			right := v.Stack.pop()
			left := v.Stack.pop()

			switch {
			case left.ValueType() == value.VALUE_INT && right.ValueType() == value.VALUE_INT:
				leftInt := left.(value.IntegerValue)
				rightInt := right.(value.IntegerValue)

				v.Stack.push(value.BoolValue{Value: leftInt.Value <= rightInt.Value})
			}

		case opcode.JumpFalseOpCode:
			operand := v.Stack.pop()

			if !isTruthy(operand) {
				ip += instruction.Offset

				continue
			}

		case opcode.JumpOpCode:
			ip += instruction.Offset

			continue

		case opcode.BlockOpCode:
			v.EnvStack.push(NewEnvironment(instruction.VarCount))

		case opcode.EndBlockOpCode:
			v.EnvStack.pop()

		case opcode.LetOpCode:
			value := v.Stack.pop()
			v.EnvStack.let(instruction.Index, value)

		case opcode.SetOpCode:
			value := v.Stack.pop()
			newVal := v.EnvStack.set(instruction.ScopeIndex, instruction.Index, value)
			v.Stack.push(newVal)

		case opcode.GetOpCode:
			value := v.EnvStack.get(instruction.ScopeIndex, instruction.Index)
			v.Stack.push(value)

		case opcode.PopOpCode:
			if len(v.Stack.Value) != 0 {
				v.Stack.pop()
			}

		default:
			panic("Unimplemented OpCode")
		}

		ip++
	}
}

func isTruthy(input value.Value) bool {
	switch input := input.(type) {
	case value.BoolValue:
		return input.Value
	default:
		return true
	}
}
