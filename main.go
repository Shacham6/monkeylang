package main

import (
	"flag"
	"fmt"
	"io"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/repl"
	"os"
	"os/user"
)

func main() {
	args := ParseArgs()

	if args.ShouldEnterRepl() {
		startRepl()
		return
	}

	startFile(args.File)
}

func startRepl() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
	fmt.Printf("Type in commands.\n")
	repl.Start(os.Stdin, os.Stdout)
}

func startFile(file string) {
	buff, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	parser := parser.New(lexer.New(string(buff)))
	program := parser.ParseProgram()
	if len(parser.Errors()) > 0 {
		printParserErrors(os.Stderr, parser.Errors())
		// If we have errors we cannot reliably continue to evaluate anything.
		os.Exit(1)
	}

	env := object.NewEnvironment()
	evaluator.Eval(program, env)
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Oops! We ran into some monkey business here!\n")
	io.WriteString(out, "parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

type MonkeyProgArgs struct {
	File string
}

func ParseArgs() *MonkeyProgArgs {
	fileFlag := flag.String("file", "", "Path to the file to be evaluated. If omitted, will enter REPL instead.")
	flag.Parse()
	return &MonkeyProgArgs{
		File: *fileFlag,
	}
}

func (m *MonkeyProgArgs) ShouldEnterRepl() bool {
	return m.File == ""
}
