package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/compiler"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/vm"
)

const PROMPT = ">> "

func StartTree(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	env := object.NewEnvironment()
	macroEnv := object.NewEnvironment()

	for {
		fmt.Fprintf(out, "%s", PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluator.DefineMacros(program, macroEnv)
		expanded := evaluator.ExpandMacros(program, macroEnv)

		evaluated := evaluator.Eval(expanded, env)
		if evaluated != nil {
			fmt.Fprintf(out, "%s\n", evaluated.Inspect())
		}
	}
}

func StartCompiled(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	macroEnv := object.NewEnvironment()

	constants := []object.Object{}
	globals := vm.InitGlobalsArray()

	symbolTable := compiler.NewSymbolTable()
	for idx, builtin := range object.Builtins {
		symbolTable.DefineBuiltin(idx, builtin.Name)
	}

	for {
		fmt.Fprintf(out, "%s", PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluator.DefineMacros(program, macroEnv)
		expanded := evaluator.ExpandMacros(program, macroEnv)

		comp := compiler.NewWithState(symbolTable, constants)

		if err := comp.Compile(expanded); err != nil {
			fmt.Fprintf(out, "Oops! Compilation failed:\n\t%s\n", err)
			continue
		}

		machine := vm.NewWithGlobalState(comp.Bytecode(), globals)
		if err := machine.Run(); err != nil {
			fmt.Fprintf(out, "Executing bytecode failed:\n\t%s\n", err)
			continue
		}

		stackTop := machine.LastPoppedStackElem()
		fmt.Fprintf(out, "%s\n", stackTop.Inspect())
	}
}

func printParserErrors(out io.Writer, errors []string) {
	fmt.Fprintf(out, "%s", "Oops! We ran into some monkey business here!\n")
	fmt.Fprintf(out, "%s", "parser errors:\n")
	for _, msg := range errors {
		fmt.Fprintf(out, "%s", "\t"+msg+"\n")
	}
}
