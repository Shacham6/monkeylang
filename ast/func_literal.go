package ast

import (
	"fmt"
	"monkey/token"
	"strings"
)

type FunctionLiteral struct {
	token      token.Token
	parameters []*Identifier
	body       *BlockStatement
	name       string
}

func NewFunctionLiteral(token token.Token, parameters []*Identifier, body *BlockStatement, name string) *FunctionLiteral {
	return &FunctionLiteral{token, parameters, body, name}
}

func (f *FunctionLiteral) Parameters() []*Identifier { return f.parameters }

func (f *FunctionLiteral) Body() *BlockStatement { return f.body }

func (f *FunctionLiteral) Name() (string, bool) { return f.name, f.name != "" }

func (f *FunctionLiteral) SetName(s string) { f.name = s }

func (*FunctionLiteral) expressionNode() {}

func (f *FunctionLiteral) TokenLiteral() string { return f.token.Literal }

func (f *FunctionLiteral) String() string {
	params := []string{}
	for _, p := range f.parameters {
		params = append(params, p.String())
	}

	if len(f.name) == 0 {
		return fmt.Sprintf("(func [%s] %s)", strings.Join(params, ", "), f.body.String())
	}

	return fmt.Sprintf("(func %s [%s] %s)", f.name, strings.Join(params, ", "), f.body.String())
}
