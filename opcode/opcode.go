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
	// args: [index]
	OP_LET
	OP_EQ
	// args: [offset]
	OP_JUMP_FALSE
	// args: [offset]
	OP_JUMP
	// args: [N of variables]
	OP_BLOCK
	OP_END_BLOCK
	// args: [index, scope index]
	OP_SET
	// args: [index, scope index]
	OP_GET
	OP_POP
	OP_ECHO
)
