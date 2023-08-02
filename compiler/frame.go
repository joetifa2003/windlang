package compiler

import (
	"github.com/joetifa2003/windlang/opcode"
)

type VarScope int

const (
	VarScopeLocal VarScope = iota
	VarScopeGlobal
	VarScopeFree
)

type Var struct {
	Name  string
	Index int
	Scope VarScope
}

func (v *Var) Get() []opcode.OpCode {
	return opcode.Instructions{opcode.OP_GET, opcode.OpCode(v.Index)}
}

func (v *Var) Set() opcode.Instructions {
	return opcode.Instructions{opcode.OP_SET, opcode.OpCode(v.Index)}
}

type Frame struct {
	Parent       *Frame
	Instructions []opcode.OpCode
	Locals       []Var
	FreeVars     []Var
	blocks       [][]Var
}

func NewFrame(parent *Frame) Frame {
	return Frame{
		Parent:       parent,
		Instructions: []opcode.OpCode{},
		Locals:       []Var{},
		FreeVars:     []Var{},
		blocks: [][]Var{
			{},
		},
	}
}

func (f *Frame) currentBlock() *[]Var {
	return &f.blocks[len(f.blocks)-1]
}

func (f *Frame) beginBlock() {
	f.blocks = append(f.blocks, []Var{})
}

func (f *Frame) endBlock() {
	f.blocks = f.blocks[:len(f.blocks)-1]
}

func (f *Frame) define(name string) int {
	local := Var{Name: name, Index: len(f.Locals)}
	if f.Parent == nil {
		local.Scope = VarScopeGlobal
	} else {
		local.Scope = VarScopeLocal
	}

	f.Locals = append(f.Locals, local)
	block := f.currentBlock()
	*block = append(*block, local)

	return local.Index
}

func (f *Frame) findLocal(name string) (Var, bool) {
	for blockIdx := len(f.blocks) - 1; blockIdx >= 0; blockIdx-- {
		for _, v := range f.blocks[blockIdx] {
			if v.Name == name {
				return v, true
			}
		}
	}

	return Var{}, false
}

func (f *Frame) defineFree(v Var) Var {
	f.FreeVars = append(f.FreeVars, v)
	v.Scope = VarScopeFree
	v.Index = len(f.FreeVars) - 1

	return v
}

func (f *Frame) resolve(name string) Var {
	v, ok := f.findLocal(name)
	if !ok {
		if f.Parent == nil {
			panic("cannot resolve variable " + name)
		}

		parentV := f.Parent.resolve(name)

		if parentV.Scope == VarScopeGlobal {
			return parentV
		}

		return f.defineFree(parentV)
	}

	return v
}
