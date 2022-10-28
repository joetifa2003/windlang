package opcode

import (
	"fmt"

	"github.com/joetifa2003/windlang/value"
)

type OpType int

const (
	OP_CONST OpType = iota
	OP_ADD
	OP_SUBTRACT
	OP_MULTIPLY
	OP_DIVIDE
	OP_MODULO
	OP_LESS
	OP_LESSEQ
	OP_LET
	OP_EQ
	OP_JUMP_FALSE
	OP_JUMP
	OP_BLOCK
	OP_END_BLOCK
	OP_SET
	OP_GET
	OP_POP
	OP_ECHO
)

type OpCode interface {
	Type() OpType
	String() string
}

type ConstOpCode struct{ Value value.Value }

func (c ConstOpCode) String() string { return fmt.Sprintf("CONST %s", c.Value.String()) }
func (ConstOpCode) Type() OpType     { return OP_CONST }

type AddOpCode struct{}

func (AddOpCode) String() string { return "ADD" }
func (AddOpCode) Type() OpType   { return OP_ADD }

type SubtractOpCode struct{}

func (SubtractOpCode) String() string { return "SUBTRACT" }
func (SubtractOpCode) Type() OpType   { return OP_SUBTRACT }

type MultiplyOpCode struct{}

func (MultiplyOpCode) String() string { return "MULTIPLY" }
func (MultiplyOpCode) Type() OpType   { return OP_MULTIPLY }

type ModuloOpCode struct{}

func (ModuloOpCode) String() string { return "MODULO" }
func (ModuloOpCode) Type() OpType   { return OP_MODULO }

type DivideOpCode struct{}

func (DivideOpCode) String() string { return "DIVIDE" }
func (DivideOpCode) Type() OpType   { return OP_DIVIDE }

type EqualOpCode struct{}

func (EqualOpCode) String() string { return "EQ" }
func (EqualOpCode) Type() OpType   { return OP_EQ }

type LessThanOpCode struct{}

func (LessThanOpCode) String() string { return "LESS" }
func (LessThanOpCode) Type() OpType   { return OP_LESS }

type LessThanEqOpCode struct{}

func (LessThanEqOpCode) String() string { return "LESSEQ" }
func (LessThanEqOpCode) Type() OpType   { return OP_LESSEQ }

type JumpFalseOpCode struct{ Offset int }

func (j JumpFalseOpCode) String() string { return fmt.Sprintf("JMP_FALSE %d", j.Offset) }
func (JumpFalseOpCode) Type() OpType     { return OP_JUMP_FALSE }

type JumpOpCode struct{ Offset int }

func (j JumpOpCode) String() string { return fmt.Sprintf("JMP %d", j.Offset) }
func (JumpOpCode) Type() OpType     { return OP_JUMP }

type BlockOpCode struct{ VarCount int }

func (BlockOpCode) String() string { return "BLOCK" }
func (BlockOpCode) Type() OpType   { return OP_BLOCK }

type EndBlockOpCode struct{}

func (EndBlockOpCode) String() string { return "END_BLOCK" }
func (EndBlockOpCode) Type() OpType   { return OP_END_BLOCK }

type LetOpCode struct{ Index int }

func (s LetOpCode) String() string { return fmt.Sprintf("LET %d", s.Index) }
func (LetOpCode) Type() OpType     { return OP_LET }

type SetOpCode struct {
	Index      int
	ScopeIndex int
}

func (s SetOpCode) String() string { return fmt.Sprintf("SET %d %d", s.ScopeIndex, s.ScopeIndex) }
func (SetOpCode) Type() OpType     { return OP_SET }

type GetOpCode struct {
	Index      int
	ScopeIndex int
}

func (g GetOpCode) String() string { return fmt.Sprintf("GET %d %d", g.ScopeIndex, g.Index) }
func (GetOpCode) Type() OpType     { return OP_GET }

type PopOpCode struct{}

func (PopOpCode) String() string { return "POP" }
func (PopOpCode) Type() OpType   { return OP_POP }

type EchoOpCode struct{}

func (EchoOpCode) String() string { return "ECHO" }
func (EchoOpCode) Type() OpType   { return OP_ECHO }
