package compiler

import (
	"github.com/joetifa2003/windlang/ast"
	"github.com/joetifa2003/windlang/opcode"
	"github.com/joetifa2003/windlang/value"
)

type Local struct {
	Name         string
	Scope        int
	IndexInFrame int
}

type Frame struct {
	Instructions []opcode.OpCode
	Locals       []Local
	scopes       [][]Local
}

func NewFrame() Frame {
	return Frame{
		Instructions: []opcode.OpCode{},
		Locals:       []Local{},
		scopes:       [][]Local{},
	}
}

type Compiler struct {
	Constants []value.Value
	Frames    []Frame
}

func NewCompiler() Compiler {
	return Compiler{
		Frames:    []Frame{NewFrame()},
		Constants: []value.Value{},
	}
}

func (c *Compiler) addConstant(v value.Value) int {
	c.Constants = append(c.Constants, v)

	return len(c.Constants) - 1
}

func (c *Compiler) curFrame() *Frame {
	return &c.Frames[len(c.Frames)-1]
}

func (c *Compiler) beginScope() {
	frame := c.curFrame()
	frame.scopes = append(frame.scopes, []Local{})
}

func (c *Compiler) endScope() {
	frame := c.curFrame()
	frame.scopes = frame.scopes[:len(frame.scopes)-1]
}

func (c *Compiler) addToScope(name string) int {
	frame := c.curFrame()
	local := Local{Name: name, Scope: len(frame.scopes) - 1, IndexInFrame: len(frame.Locals)}
	frame.Locals = append(frame.Locals, local)

	scope := &frame.scopes[len(frame.scopes)-1]
	*scope = append(*scope, local)

	return len(frame.Locals) - 1
}

// findInScope returns (frameOffset)
func (c *Compiler) findInScope(name string) int {
	frame := c.curFrame()

	for scopesIndex := len(frame.scopes) - 1; scopesIndex >= 0; scopesIndex-- {
		for _, v := range frame.scopes[scopesIndex] {
			if v.Name == name {
				return v.IndexInFrame
			}
		}
	}

	panic("")
}

func (c *Compiler) Compile(node ast.Node) opcode.Instructions {
	switch node := node.(type) {
	case *ast.Program:
		var instructions []opcode.OpCode

		c.beginScope()
		programInstructions := append(instructions, c.CompileProgram(node.Statements)...)
		c.endScope()

		instructions = append(instructions, programInstructions...)

		return instructions

	case *ast.ExpressionStatement:
		expr := c.Compile(node.Expression)
		expr = append(expr, opcode.OP_POP)

		return expr

	case *ast.InfixExpression:
		var instructions []opcode.OpCode

		left := c.Compile(node.Left)
		right := c.Compile(node.Right)

		instructions = append(instructions, left...)
		instructions = append(instructions, right...)

		switch node.Operator {
		case "+":
			instructions = append(instructions, opcode.OP_ADD)
		case "-":
			instructions = append(instructions, opcode.OP_SUBTRACT)
		case "*":
			instructions = append(instructions, opcode.OP_MULTIPLY)
		case "/":
			instructions = append(instructions, opcode.OP_DIVIDE)
		case "<=":
			instructions = append(instructions, opcode.OP_LESSEQ)
		case "%":
			instructions = append(instructions, opcode.OP_MODULO)
		case "==":
			instructions = append(instructions, opcode.OP_EQ)

		default:
			panic("Unimplemented operator " + node.Operator)
		}

		return instructions

	case *ast.IfExpression:
		var instructions []opcode.OpCode

		condition := c.Compile(node.Condition)
		thenBranch := c.Compile(node.ThenBranch)
		elseBranch := []opcode.OpCode{}
		if node.ElseBranch != nil {
			elseBranch = c.Compile(node.ElseBranch)
		}

		instructions = append(instructions, condition...)
		instructions = append(instructions, opcode.OP_JUMP_FALSE)
		instructions = append(instructions, opcode.OpCode(len(thenBranch)+3))
		instructions = append(instructions, thenBranch...)
		instructions = append(instructions, opcode.OP_JUMP)
		instructions = append(instructions, opcode.OpCode(len(elseBranch)+1))
		instructions = append(instructions, elseBranch...)

		return instructions

	case *ast.IntegerLiteral:
		return []opcode.OpCode{
			opcode.OP_CONST,
			opcode.OpCode(
				c.addConstant(value.NewIntValue(node.Value)),
			),
		}

	case *ast.Boolean:
		return []opcode.OpCode{
			opcode.OP_CONST,
			opcode.OpCode(
				c.addConstant(value.NewBoolValue(node.Value)),
			),
		}

	case *ast.BlockStatement:
		var instructions []opcode.OpCode
		var bodyInstructions []opcode.OpCode

		if node.VarCount != 0 {
			c.beginScope()
			for _, stmt := range node.Statements {
				bodyInstructions = append(bodyInstructions, c.Compile(stmt)...)
			}
			c.endScope()

			instructions = append(instructions, bodyInstructions...)

			return instructions
		} else {
			for _, stmt := range node.Statements {
				bodyInstructions = append(bodyInstructions, c.Compile(stmt)...)
			}

			return bodyInstructions
		}

	case *ast.WhileStatement:
		var instructions []opcode.OpCode

		condition := c.Compile(node.Condition)
		body := c.Compile(node.Body)

		instructions = append(instructions, condition...)
		instructions = append(instructions, opcode.OP_JUMP_FALSE)
		instructions = append(instructions, opcode.OpCode(len(body)+3))
		instructions = append(instructions, body...)
		instructions = append(instructions, opcode.OP_JUMP)
		instructions = append(instructions, opcode.OpCode(-len(body)-len(condition)-3))

		return instructions

	case *ast.ForStatement:
		var bodyInstructions []opcode.OpCode
		var instructions []opcode.OpCode

		c.beginScope()
		initializer := c.Compile(node.Initializer)
		condition := c.Compile(node.Condition)
		body := c.Compile(node.Body)
		increment := c.Compile(node.Increment)
		increment = append(increment, opcode.OP_POP)
		c.endScope()

		bodyInstructions = append(bodyInstructions, initializer...)
		bodyInstructions = append(bodyInstructions, condition...)
		bodyInstructions = append(bodyInstructions, opcode.OP_JUMP_FALSE)
		bodyInstructions = append(bodyInstructions, opcode.OpCode(len(body)+len(increment)+3))
		bodyInstructions = append(bodyInstructions, body...)
		bodyInstructions = append(bodyInstructions, increment...)
		bodyInstructions = append(bodyInstructions, opcode.OP_JUMP)
		bodyInstructions = append(bodyInstructions, opcode.OpCode(-len(body)-len(increment)-len(condition)-3))

		instructions = append(instructions, bodyInstructions...)

		return instructions

	case *ast.LetStatement:
		var instructions []opcode.OpCode

		frameOffset := c.addToScope(node.Name.Value)
		value := c.Compile(node.Value)
		instructions = append(instructions, value...)
		instructions = append(instructions, opcode.OP_LET)
		instructions = append(instructions, opcode.OpCode(frameOffset))

		return instructions

	case *ast.Identifier:
		frameOffset := c.findInScope(node.Value)
		return []opcode.OpCode{
			opcode.OP_GET,
			opcode.OpCode(frameOffset),
		}

	case *ast.AssignExpression:
		var instructions []opcode.OpCode
		frameOffset := c.findInScope(node.Name.TokenLiteral())

		value := c.Compile(node.Value)
		instructions = append(instructions, value...)
		instructions = append(instructions,
			opcode.OP_SET,
			opcode.OpCode(frameOffset),
		)

		return instructions

	case *ast.PostfixExpression:
		var instructions []opcode.OpCode
		frameOffset := c.findInScope(node.Left.TokenLiteral())
		instructions = append(instructions,
			opcode.OP_INC,
			opcode.OpCode(frameOffset),
		)

		return instructions

	case *ast.EchoStatement:
		var instructions []opcode.OpCode

		instructions = append(instructions, c.Compile(node.Value)...)
		instructions = append(instructions, opcode.OP_ECHO)

		return instructions

	case *ast.NilLiteral:
		return []opcode.OpCode{opcode.OP_CONST, opcode.OpCode(c.addConstant(value.NewNilValue()))}

	case *ast.ArrayLiteral:
		instructions := []opcode.OpCode{}

		for i := len(node.Value) - 1; i >= 0; i-- {
			instructions = append(instructions, c.Compile(node.Value[i])...)
		}

		instructions = append(instructions, opcode.OP_ARRAY, opcode.OpCode(len(node.Value)))

		return instructions

	default:
		panic("Unimplemented Ast %d")
	}
}

func (c *Compiler) CompileProgram(statements []ast.Statement) []opcode.OpCode {
	var instructions []opcode.OpCode
	for _, stmt := range statements {
		instructions = append(instructions, c.Compile(stmt)...)
	}

	return instructions
}
