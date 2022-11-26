package opcode

type OpCode int

const (
	// args: [const index]
	OP_CONST OpCode = iota
	OP_ADD
	OP_SUBTRACT
	OP_MULTIPLY
	OP_DIVIDE
	OP_MODULO
	OP_LESS
	OP_LESSEQ
	OP_LET // args: [index]
	OP_EQ
	OP_JUMP_FALSE // args: [offset]
	OP_JUMP       // args: [offset]
	OP_BLOCK      // args: [N of variables]
	OP_END_BLOCK
	OP_SET // args: [index, scope index]
	OP_GET // args: [index, scope index]
	OP_POP
	OP_ECHO
	OP_ARRAY // args: [n of elements]
)
