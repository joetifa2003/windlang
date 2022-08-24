package compiler

import (
	"github.com/joetifa2003/windlang/ast"
	"github.com/joetifa2003/windlang/opcode"
	"github.com/joetifa2003/windlang/value"
)

type Compiler struct {
	Scopes [][]string
}

func NewCompiler() Compiler {
	return Compiler{
		Scopes: [][]string{},
	}
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
		for scopeIndex, v := range c.Scopes[scopesIndex] {
			if v == name {
				return scopesIndex, scopeIndex
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

		instructions = append(instructions, opcode.BlockOpCode{VarCount: len(scope)})
		instructions = append(instructions, programInstructions...)
		instructions = append(instructions, opcode.EndBlockOpCode{})

		return instructions

	case *ast.ExpressionStatement:
		expr := c.Compile(node.Expression)
		expr = append(expr, opcode.PopOpCode{})

		return expr

	case *ast.InfixExpression:
		var instructions []opcode.OpCode

		left := c.Compile(node.Left)
		right := c.Compile(node.Right)

		instructions = append(instructions, left...)
		instructions = append(instructions, right...)

		switch node.Operator {
		case "+":
			instructions = append(instructions, opcode.AddOpCode{})
		case "-":
			instructions = append(instructions, opcode.SubtractOpCode{})
		case "*":
			instructions = append(instructions, opcode.MultiplyOpCode{})
		case "/":
			instructions = append(instructions, opcode.DivideOpCode{})
		case "<=":
			instructions = append(instructions, opcode.LessThanEqOpCode{})
		case "%":
			instructions = append(instructions, opcode.ModuloOpCode{})
		case "==":
			instructions = append(instructions, opcode.EqualOpCode{})

		default:
			panic("Unimplemented operator")
		}

		return instructions

	case *ast.IfExpression:
		var instructions []opcode.OpCode

		condition := c.Compile(node.Condition)
		thenBranch := c.Compile(node.ThenBranch)
		removeLastPop(&thenBranch)
		elseBranch := []opcode.OpCode{}
		if node.ElseBranch != nil {
			elseBranch = c.Compile(node.ElseBranch)
			removeLastPop(&elseBranch)
		}

		instructions = append(instructions, condition...)
		instructions = append(instructions, opcode.JumpFalseOpCode{Offset: len(thenBranch) + 2})
		instructions = append(instructions, thenBranch...)
		instructions = append(instructions, opcode.JumpOpCode{Offset: len(elseBranch) + 1})
		instructions = append(instructions, elseBranch...)

		return instructions

	case *ast.IntegerLiteral:
		return []opcode.OpCode{opcode.ConstOpCode{Value: value.IntegerValue{Value: node.Value}}}

	case *ast.Boolean:
		return []opcode.OpCode{opcode.ConstOpCode{Value: value.BoolValue{Value: node.Value}}}

	case *ast.BlockStatement:
		var instructions []opcode.OpCode
		var bodyInstructions []opcode.OpCode

		c.beginScope()
		for _, stmt := range node.Statements {
			bodyInstructions = append(bodyInstructions, c.Compile(stmt)...)
		}
		scope := c.endScope()

		instructions = append(instructions, opcode.BlockOpCode{VarCount: len(scope)})
		instructions = append(instructions, bodyInstructions...)
		instructions = append(instructions, opcode.EndBlockOpCode{})

		return instructions

	case *ast.WhileStatement:
		var instructions []opcode.OpCode

		condition := c.Compile(node.Condition)
		body := c.Compile(node.Body)

		instructions = append(instructions, condition...)
		instructions = append(instructions, opcode.JumpFalseOpCode{Offset: len(body) + 2})
		instructions = append(instructions, body...)
		instructions = append(instructions, opcode.JumpOpCode{Offset: -len(body) - len(condition) - 1})

		return instructions

	case *ast.ForStatement:
		var bodyInstructions []opcode.OpCode
		var instructions []opcode.OpCode

		c.beginScope()
		initializer := c.Compile(node.Initializer)
		condition := c.Compile(node.Condition)
		body := c.Compile(node.Body)
		increment := c.Compile(node.Increment)
		increment = append(increment, opcode.PopOpCode{})
		scope := c.endScope()

		bodyInstructions = append(bodyInstructions, initializer...)
		bodyInstructions = append(bodyInstructions, condition...)
		bodyInstructions = append(bodyInstructions, opcode.JumpFalseOpCode{Offset: len(body) + len(increment) + 2})
		bodyInstructions = append(bodyInstructions, body...)
		bodyInstructions = append(bodyInstructions, increment...)
		bodyInstructions = append(bodyInstructions, opcode.JumpOpCode{Offset: -len(body) - len(increment) - len(condition) - 1})

		instructions = append(instructions, opcode.BlockOpCode{VarCount: len(scope)})
		instructions = append(instructions, bodyInstructions...)
		instructions = append(instructions, opcode.EndBlockOpCode{})

		return instructions

	case *ast.LetStatement:
		var instructions []opcode.OpCode

		value := c.Compile(node.Value)
		index := c.addToScope(node.Name.Value)
		instructions = append(instructions, value...)
		instructions = append(instructions, opcode.LetOpCode{Index: index})

		return instructions

	case *ast.Identifier:
		offset, index := c.findInScope(node.Value)
		return []opcode.OpCode{opcode.GetOpCode{Index: index, ScopeIndex: offset}}

	case *ast.AssignExpression:
		var instructions []opcode.OpCode
		offset, index := c.findInScope(node.Name.TokenLiteral())
		value := c.Compile(node.Value)
		instructions = append(instructions, value...)
		instructions = append(instructions, opcode.SetOpCode{Index: index, ScopeIndex: offset})

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

func removeLastPop(instructions *[]opcode.OpCode) {
	instructionsLen := len((*instructions))
	lastInstruction := (*instructions)[instructionsLen-1]

	switch lastInstruction.(type) {
	case opcode.PopOpCode:
		*instructions = (*instructions)[:instructionsLen-1]

	case opcode.EndBlockOpCode:
		_, ok := (*instructions)[instructionsLen-2].(opcode.PopOpCode)
		if ok {
			(*instructions) = removeIndex(*instructions, instructionsLen-2)
		}
	}
}

func removeIndex[T any](s []T, index int) []T {
	return append(s[:index], s[index+1:]...)
}
