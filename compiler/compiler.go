package compiler

import (
	"github.com/joetifa2003/windlang/ast"
	"github.com/joetifa2003/windlang/opcode"
	"github.com/joetifa2003/windlang/value"
)

type Compiler struct {
	Scopes    [][]string
	Constants []value.Value
}

func NewCompiler() Compiler {
	return Compiler{
		Scopes:    [][]string{},
		Constants: []value.Value{},
	}
}

func (c *Compiler) addConstant(v value.Value) int {
	c.Constants = append(c.Constants, v)

	return len(c.Constants) - 1
}

func (c *Compiler) beginScope() {
	c.Scopes = append(c.Scopes, []string{})
}

func (c *Compiler) endScope() []string {
	lastScope := c.Scopes[len(c.Scopes)-1]
	c.Scopes = c.Scopes[:len(c.Scopes)-1]

	return lastScope
}

func (c *Compiler) addToScope(name string) int {
	scope := &c.Scopes[len(c.Scopes)-1]
	*scope = append(*scope, name)

	return len(*scope) - 1
}

// returns scope index and the value index inside it
func (c *Compiler) findInScope(name string) (int, int) {
	for scopesIndex := len(c.Scopes) - 1; scopesIndex >= 0; scopesIndex-- {
		for valueIndex, v := range c.Scopes[scopesIndex] {
			if v == name {
				return scopesIndex, valueIndex
			}
		}
	}

	panic("")
}

func (c *Compiler) Compile(node ast.Node) []opcode.OpCode {
	switch node := node.(type) {
	case *ast.Program:
		var instructions []opcode.OpCode

		c.beginScope()
		programInstructions := append(instructions, c.CompileProgram(node.Statements)...)
		scope := c.endScope()

		instructions = append(instructions, opcode.OP_BLOCK)
		instructions = append(instructions, opcode.OpCode(len(scope)))
		instructions = append(instructions, programInstructions...)
		instructions = append(instructions, opcode.OP_END_BLOCK)

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
			scope := c.endScope()

			instructions = append(instructions, opcode.OP_BLOCK)
			instructions = append(instructions, opcode.OpCode(len(scope)))
			instructions = append(instructions, bodyInstructions...)
			instructions = append(instructions, opcode.OP_END_BLOCK)

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
		scope := c.endScope()

		bodyInstructions = append(bodyInstructions, initializer...)
		bodyInstructions = append(bodyInstructions, condition...)
		bodyInstructions = append(bodyInstructions, opcode.OP_JUMP_FALSE)
		bodyInstructions = append(bodyInstructions, opcode.OpCode(len(body)+len(increment)+3))
		bodyInstructions = append(bodyInstructions, body...)
		bodyInstructions = append(bodyInstructions, increment...)
		bodyInstructions = append(bodyInstructions, opcode.OP_JUMP)
		bodyInstructions = append(bodyInstructions, opcode.OpCode(-len(body)-len(increment)-len(condition)-3))

		instructions = append(instructions, opcode.OP_BLOCK)
		instructions = append(instructions, opcode.OpCode(len(scope)))
		instructions = append(instructions, bodyInstructions...)
		instructions = append(instructions, opcode.OP_END_BLOCK)

		return instructions

	case *ast.LetStatement:
		var instructions []opcode.OpCode

		value := c.Compile(node.Value)
		index := c.addToScope(node.Name.Value)
		instructions = append(instructions, value...)
		instructions = append(instructions, opcode.OP_LET)
		instructions = append(instructions, opcode.OpCode(index))

		return instructions

	case *ast.Identifier:
		scopeIndex, valueIndex := c.findInScope(node.Value)
		return []opcode.OpCode{
			opcode.OP_GET,
			opcode.OpCode(valueIndex),
			opcode.OpCode(scopeIndex),
		}

	case *ast.AssignExpression:
		var instructions []opcode.OpCode
		scopeIndex, valueIndex := c.findInScope(node.Name.TokenLiteral())
		value := c.Compile(node.Value)
		instructions = append(instructions, value...)
		instructions = append(instructions,
			opcode.OP_SET,
			opcode.OpCode(valueIndex),
			opcode.OpCode(scopeIndex),
		)

		return instructions

	case *ast.PostfixExpression:
		var instructions []opcode.OpCode
		scopeIndex, valueIndex := c.findInScope(node.Left.TokenLiteral())
		instructions = append(instructions,
			opcode.OP_INC,
			opcode.OpCode(valueIndex),
			opcode.OpCode(scopeIndex),
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
