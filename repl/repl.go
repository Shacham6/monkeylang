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

		comp := compiler.New()

		if err := comp.Compile(expanded); err != nil {
			fmt.Fprintf(out, "Oops! Compilation failed:\n%s\n", err)
		}

		machine := vm.New(comp.Bytecode())
		if err := machine.Run(); err != nil {
			fmt.Fprintf(out, "Executing bytecode failed:\n%s\n", err)
		}

		stackTop := machine.StackTop()
		fmt.Fprintf(out, "%s\n", stackTop.Inspect())
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Oops! We ran into some monkey business here!\n")
	io.WriteString(out, "parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
