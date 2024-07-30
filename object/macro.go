package object

import (
	"bytes"
	"monkey/ast"
	"strings"
)

type Macro struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (m *Macro) Type() ObjectType {
	return MACRO_OBJ
}

func (m *Macro) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, param := range m.Parameters {
		params = append(params, param.String())
	}

	out.WriteString("macro")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(m.Body.String())
	out.WriteString("\n}")

	return out.String()
}

var x Object = &Macro{}
