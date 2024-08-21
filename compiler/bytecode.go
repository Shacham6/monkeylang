package compiler

import (
	"monkey/code"
	"monkey/object"
)

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
