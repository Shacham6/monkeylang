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
}

func NewFunctionLiteral(token token.Token, parameters []*Identifier, body *BlockStatement) *FunctionLiteral {
	return &FunctionLiteral{token, parameters, body}
}

func (f *FunctionLiteral) Parameters() []*Identifier { return f.parameters }

func (f *FunctionLiteral) Body() *BlockStatement { return f.body }

func (*FunctionLiteral) expressionNode() {}

func (f *FunctionLiteral) TokenLiteral() string { return f.token.Literal }

func (f *FunctionLiteral) String() string {
	params := []string{}
	for _, p := range f.parameters {
		params = append(params, p.String())
	}

	return fmt.Sprintf("(func [%s] %s)", strings.Join(params, ", "), f.body.String())
}
