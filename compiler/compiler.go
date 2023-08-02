package compiler

import (
	"fmt"

	"github.com/joetifa2003/windlang/ast"
	"github.com/joetifa2003/windlang/opcode"
	"github.com/joetifa2003/windlang/value"
)

type Compiler struct {
	Constants []value.Value
	Frames    []Frame
}

func NewCompiler() Compiler {
	return Compiler{
		Frames:    []Frame{NewFrame(nil)},
		Constants: []value.Value{},
	}
}

func (c *Compiler) addConstant(v value.Value) int {
	c.Constants = append(c.Constants, v)

	return len(c.Constants) - 1
}

func (c *Compiler) pushFrame() {
	c.Frames = append(c.Frames, NewFrame(c.curFrame()))
}

func (c *Compiler) popFrame() {
	c.Frames = c.Frames[:len(c.Frames)-1]
}

func (c *Compiler) curFrame() *Frame {
	return &c.Frames[len(c.Frames)-1]
}

func (c *Compiler) beginBlock() {
	c.curFrame().beginBlock()
}

func (c *Compiler) endBlock() {
	c.curFrame().endBlock()
}

func (c *Compiler) define(name string) int {
	return c.curFrame().define(name)
}

// resolve returns (frameOffset)
func (c *Compiler) resolve(name string) Var {
	return c.curFrame().resolve(name)
}

func (c *Compiler) Compile(node ast.Node) opcode.Instructions {
	switch node := node.(type) {
	case *ast.Program:
		return c.CompileProgram(node.Statements)

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

	case *ast.IfStatement:
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

		c.beginBlock()
		for _, stmt := range node.Statements {
			instructions = append(instructions, c.Compile(stmt)...)
		}
		c.endBlock()

		return instructions

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

		c.beginBlock()
		initializer := c.Compile(node.Initializer)
		condition := c.Compile(node.Condition)
		body := c.Compile(node.Body)
		increment := c.Compile(node.Increment)
		increment = append(increment, opcode.OP_POP)
		c.endBlock()

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

		frameOffset := c.define(node.Name.Value)
		value := c.Compile(node.Value)
		instructions = append(instructions, value...)
		instructions = append(instructions, opcode.OP_LET)
		instructions = append(instructions, opcode.OpCode(frameOffset))

		return instructions

	case *ast.Identifier:
		variable := c.resolve(node.Value)
		return variable.Get()

	case *ast.AssignExpression:
		var instructions []opcode.OpCode
		variable := c.resolve(node.Name.TokenLiteral())

		value := c.Compile(node.Value)
		instructions = append(instructions, value...)
		instructions = append(instructions, variable.Set()...)

		return instructions

	// case *ast.PostfixExpression:
	// 	var instructions []opcode.OpCode
	// 	frameOffset := c.resolve(node.Left.TokenLiteral())
	// 	instructions = append(instructions,
	// 		opcode.OP_INC,
	// 		opcode.OpCode(frameOffset),
	// 	)
	//
	// 	return instructions

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

	case *ast.FunctionLiteral:
		c.pushFrame()

		for _, param := range node.Parameters {
			c.define(param.Value)
		}

		bodyInstructions := c.Compile(node.Body)
		instructions := []opcode.OpCode{
			opcode.OP_CONST,
			opcode.OpCode(
				c.addConstant(
					value.NewFnValue(bodyInstructions, len(c.curFrame().Locals)),
				),
			),
		}

		c.popFrame()

		return instructions

	case *ast.CallExpression:
		instructions := c.Compile(node.Function)
		instructions = append(instructions, opcode.OP_CALL)

		return instructions

	default:
		panic(fmt.Sprintf("Unimplemented Ast %T", node))
	}
}

func (c *Compiler) CompileProgram(statements []ast.Statement) []opcode.OpCode {
	var instructions []opcode.OpCode
	for _, stmt := range statements {
		instructions = append(instructions, c.Compile(stmt)...)
	}

	return instructions
}
