package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/joetifa2003/windlang/compiler"
	"github.com/joetifa2003/windlang/lexer"
	"github.com/joetifa2003/windlang/parser"
	"github.com/joetifa2003/windlang/vm"
	"github.com/pkg/profile"
	"github.com/spf13/cobra"
)

var (
	debug = false
)

var vmCommand = &cobra.Command{
	Use:   "vm [file]",
	Short: "Run a Wind script using vm",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("requires 1 argument")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]

		file, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalln("Could not read file:", err)
			return
		}

		defer profile.Start().Stop()
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

		compiler := compiler.NewCompiler()
		instructions := compiler.Compile(program)

		virtualM := vm.NewVM(compiler.Constants)
		virtualM.Interpret(instructions)
	},
}

func init() {
	vmCommand.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Show debug info")
	rootCmd.AddCommand(vmCommand)
}
