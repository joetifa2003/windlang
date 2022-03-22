package main

import (
	"fmt"
	"io/ioutil"
	"wind-vm-go/evaluator"
	"wind-vm-go/lexer"
	"wind-vm-go/object"
	"wind-vm-go/parser"
)

func main() {
	fileName := "main.wind"

	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Could not read file:", err)
		return
	}

	input := string(file)
	lexer := lexer.New(input)
	parser := parser.New(lexer)
	program := parser.ParseProgram()

	if len(parser.Errors) != 0 {
		for _, err := range parser.Errors {
			fmt.Printf("[ERROR] %s\n", err)
		}

		return
	}

	envManager := object.NewEnvironmentManager()
	env, _ := envManager.Get(fileName)
	evaluator := evaluator.New(envManager)
	evaluated := evaluator.Eval(program, env)

	switch evaluated.Type() {
	case object.ERROR_OBJ:
		fmt.Println(evaluated.Inspect())
	}
}
