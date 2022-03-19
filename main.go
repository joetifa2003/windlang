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
		let even = 0;
		let odd = 0;
		
		for (let i = 0; i < 100000; i++) {
			
			if (i % 2 == 0) {
				even = even + 1;
			} else {
				odd = odd + 1;
			}

		}

		println(even);
		println(odd);
	`

	lexer := lexer.New(input)
	parser := parser.New(lexer)
	program := parser.ParseProgram()

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)
	fmt.Println(evaluated.Inspect())
}
