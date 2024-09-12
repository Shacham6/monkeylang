package main

import (
	"flag"
	"fmt"
	"log"
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

		switch args.Engine {
		case ENGINE_VM:
			repl.StartCompiled(os.Stdin, os.Stdout)
			return
		case ENGINE_TREE:
			repl.StartTree(os.Stdin, os.Stdout)
			return
		default:
			log.Fatalf("REPL with engine of type %s is not supported", args.Engine)
		}
		return
	}

	fileexec.ExecFile(args.File)
}

type EngineType string

const (
	ENGINE_VM   EngineType = "vm"
	ENGINE_TREE EngineType = "tree"
)

func isEngineTypeAllowed(engineType EngineType) bool {
	return engineType == ENGINE_VM || engineType == ENGINE_TREE
}

type MonkeyProgArgs struct {
	File   string
	Engine EngineType
}

func ParseArgs() *MonkeyProgArgs {
	fileFlag := flag.String("file", "", "Path to the file to be evaluated. If omitted, will enter REPL instead.")
	engineFlag := flag.String("engine", "vm", "The backend engine to evaluate the language. [vm, tree]")

	flag.Parse()

	engine := EngineType(*engineFlag)

	if !isEngineTypeAllowed(engine) {
		fmt.Fprintf(os.Stderr, "provided invalid '-engine' (%s), possible values are: [vm, tree]", *engineFlag)
		os.Exit(1)
	}

	return &MonkeyProgArgs{
		File:   *fileFlag,
		Engine: engine,
	}
}

func (m *MonkeyProgArgs) ShouldEnterRepl() bool {
	return m.File == ""
}
