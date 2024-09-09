package main

import (
	"flag"
	"fmt"
	"monkey/fileexec"
	"monkey/repl"
	"os"
	"os/user"
)

func main() {
	args := ParseArgs()

	if args.ShouldEnterRepl() {
		user, err := user.Current()
		if err != nil {
			panic(err)
		}

		fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
		fmt.Printf("Type in commands.\n")
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	fileexec.ExecFile(args.File)
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
