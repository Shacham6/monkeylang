package object

import (
	"fmt"
	"monkey/ast"
)

type Quote struct {
	Node ast.Node
}

func (q *Quote) Type() ObjectType {
	return QUOTE_OBJ
}

func (q *Quote) Inspect() string {
	return fmt.Sprintf("QUOTE(%s)", q.Node.String())
}
