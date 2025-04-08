package vm

import (
	"fmt"
	"monkey/code"
	"monkey/object"
	"strings"
)

type VmRunError struct {
	Err          error
	Instructions code.Instructions

	Stack        []object.Object
	Globals      []object.Object
	StackPointer int // "stack pointer". Always points to the next value. Top of stack is stack[sp-1]
}

func (e *VmRunError) Error() string {
	lines := []string{}
	lines = append(lines, fmt.Sprintf("Got runtime error: %s", e.Err))

	globalsLines := []string{"Globals:"}
	for idx, g := range e.Globals {
		if g == nil {
			continue
		}
		globalsLines = append(globalsLines, fmt.Sprintf("[%04d] %s", idx, g.Inspect()))
	}
	globalsString := strings.Join(globalsLines, "\n")
	lines = append(lines, globalsString)

	stackLines := []string{"Stack:"}
	for idx, so := range e.Stack {
		if so == nil {
			continue
		}

		var pointer string
		if idx == e.StackPointer {
			pointer = "<-------"
		} else {
			pointer = ""
		}

		stackLines = append(stackLines, fmt.Sprintf("[%04d] %s %s", idx, so.Inspect(), pointer))
	}
	stackString := strings.Join(stackLines, "\n")
	lines = append(lines, stackString)

	runtimeDump := fmt.Sprintf(
		"Runtime Dump:\n%s",
		e.Instructions.String(),
	)

	lines = append(lines, runtimeDump)

	return strings.Join(lines, "\n=====================\n")
}
