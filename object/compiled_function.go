package object

import (
	"fmt"
	"monkey/code"
)

type CompiledFunction struct {
	Instructions code.Instructions
	NumLocals    int
}

func (c *CompiledFunction) Type() ObjectType { return COMPILED_FUNCTION_OJ }

func (c *CompiledFunction) Inspect() string {
	return fmt.Sprintf("CompiledFunction[%p]", c)
}
