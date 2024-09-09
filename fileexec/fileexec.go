// filexec contains the entrypoint for executing a monkey script from the filesystem.
package fileexec

import (
	"fmt"
	"io"
	"monkey/compiler"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/vm"
	"os"
)

func ExecFile(filepath string) {
	buff, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed reading file at the given path with an error:\n%v\n", err)
		os.Exit(1)
	}

	parser := parser.New(lexer.New(string(buff)))
	program := parser.ParseProgram()
	if len(parser.Errors()) > 0 {
		printParserErrors(os.Stderr, parser.Errors())
		// If we have errors we cannot reliably continue to evaluate anything.
		os.Exit(1)
	}

	// env := object.NewEnvironment()
	macroEnv := object.NewEnvironment()

	evaluator.DefineMacros(program, macroEnv)
	expandedProgram := evaluator.ExpandMacros(program, macroEnv)

	comp := compiler.New()
	if err := comp.Compile(expandedProgram); err != nil {
		fmt.Fprintf(os.Stderr, "Ouch! Failed compiling program:\n%s\n", err)
		os.Exit(1)
	}

	machine := vm.New(comp.Bytecode())
	if err := machine.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoop! Failed execution with an error:\n%s\n", err)
		os.Exit(1)
	}

	// evaluationResult := evaluator.Eval(expandedProgram, env)
	//
	// if evaluationResult.Type() == object.ERROR_OBJ {
	// 	fmt.Fprintf(os.Stderr, "Encountered a runtime error: %s\n", evaluationResult.Inspect())
	// }
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Oops! We ran into some monkey business here!\n")
	io.WriteString(out, "parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
