package ast

import (
	"bytes"
	"monkey/token"
	"strings"
)

type CallExpression struct {
	token     token.Token // the '(' token
	function  Expression  // identifier or FunctionLiteral
	arguments []Expression
}

func NewCallExpression(token token.Token, function Expression, arguments []Expression) *CallExpression {
	return &CallExpression{token, function, arguments}
}

func (c *CallExpression) Function() Expression { return c.function }

func (c *CallExpression) Arguments() []Expression { return c.arguments }

func (*CallExpression) expressionNode() {}

func (c *CallExpression) TokenLiteral() string { return c.token.Literal }

func (c *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range c.arguments {
		args = append(args, a.String())
	}

	out.WriteString(c.function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}
