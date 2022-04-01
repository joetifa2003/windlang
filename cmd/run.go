package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"wind-vm-go/evaluator"
	"wind-vm-go/lexer"
	"wind-vm-go/parser"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [file]",
	Short: "Run a Wind script",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("requires 1 argument")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		if filePath == "" {
			log.Fatal("File path is required")
		}

		file, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatalln("Could not read file:", err)
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

		envManager := evaluator.NewEnvironmentManager()
		env, _ := envManager.Get(filePath)
		ev := evaluator.New(envManager)
		evaluated, evErr := ev.Eval(program, env)
		if evErr != nil {
			fmt.Println(evErr.Inspect())
		}

		if evaluated == nil {
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
