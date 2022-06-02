package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/joetifa2003/windlang/evaluator"
	"github.com/joetifa2003/windlang/lexer"
	"github.com/joetifa2003/windlang/parser"

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
		parser := parser.New(lexer, filePath)
		program := parser.ParseProgram()
		parserErrors := parser.ReportErrors()
		if len(parserErrors) > 0 {
			for _, err := range parserErrors {
				fmt.Println(err)
			}

			os.Exit(1)
		}

		envManager := evaluator.NewEnvironmentManager()
		env, _ := envManager.Get(filePath)
		ev := evaluator.New(envManager, filePath)
		evaluated, evErr := ev.Eval(program, env, nil)
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
