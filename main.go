package main

import (
	"fmt"
	"wind-vm-go/evaluator"
	"wind-vm-go/lexer"
	"wind-vm-go/object"
	"wind-vm-go/parser"
)

func main() {

	input := `
		
	`

	lexer := lexer.New(input)
	parser := parser.New(lexer)
	program := parser.ParseProgram()

	if len(parser.Errors) != 0 {
		for _, err := range parser.Errors {
			fmt.Printf("[ERROR] %s\n", err)
		}

		return
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)

	if evaluated != nil {
		fmt.Printf("%s\n", evaluated.Inspect())
	}
}
