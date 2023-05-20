package opcode

import (
	"fmt"
	"strings"
)

type OpCode int

const (
	OP_CONST OpCode = iota // args: [const index]
	OP_ADD
	OP_SUBTRACT
	OP_MULTIPLY
	OP_DIVIDE
	OP_MODULO
	OP_LESS
	OP_LESSEQ
	OP_EQ
	OP_POP
	OP_ECHO
	OP_RET

	OP_LET        // args: [frameOffset]
	OP_JUMP_FALSE // args: [offset]
	OP_JUMP       // args: [offset]
	OP_SET        // args: [frameOffset]
	OP_SET_GLOBAL // args: [index]
	OP_GET        // args: [frameOffset]
	OP_GET_GLOBAL // args [index]
	OP_INC        // args: [frameOffset]
	OP_INC_GLOBAL // args: [index]
	OP_ARRAY      // args: [n of elements]
	OP_CALL       // args: [n of args]
)

type Instructions []OpCode

func (instructions Instructions) String() string {
	var out strings.Builder
	ip := 0
	for ip < len(instructions) {
		op := instructions[ip]
		switch op {
		case OP_CONST:
			ip++
			idx := int(instructions[ip])
			out.WriteString(fmt.Sprintf("const %d", idx))
		case OP_ADD:
			out.WriteString("add")
		case OP_SUBTRACT:
			out.WriteString("sub")
		case OP_MULTIPLY:
			out.WriteString("mul")
		case OP_DIVIDE:
			out.WriteString("div")
		case OP_MODULO:
			out.WriteString("mod")
		case OP_LESS:
			out.WriteString("less")
		case OP_LESSEQ:
			out.WriteString("lessq")
		case OP_EQ:
			out.WriteString("eq")
		case OP_POP:
			out.WriteString("pop")
		case OP_ECHO:
			out.WriteString("echo")
		case OP_RET:
			out.WriteString("ret")
		case OP_LET:
			ip++
			frameOffset := int(instructions[ip])
			out.WriteString(fmt.Sprintf("let %d", frameOffset))
		case OP_GET:
			ip++
			frameOffset := int(instructions[ip])
			out.WriteString(fmt.Sprintf("get %d", frameOffset))

		default:
			panic(fmt.Sprintf("Unimplemented opcode %d", op))
		}

		ip++

		if ip < len(instructions) {
			out.WriteString("\n")
		}
	}

	return out.String()
}
