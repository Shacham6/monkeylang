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
	// var out bytes.Buffer

	params := []string{}
	for _, p := range f.parameters {
		params = append(params, p.String())
	}

	return fmt.Sprintf("(func [%s] %s)", strings.Join(params, ", "), f.body.String())

	// out.WriteString(f.TokenLiteral())
	// out.WriteString("(")
	// out.WriteString(strings.Join(params, ","))
	// out.WriteString(") ")
	// out.WriteString(f.body.String())
	//
	// return out.String()
}
