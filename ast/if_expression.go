package ast

import (
	"bytes"
	"monkey/token"
)

type IfExpression struct {
	token       token.Token
	condition   Expression
	consequence *BlockStatement
	alternative *IfExpressionAlternative
}

func NewIfExpression(t token.Token, condition Expression, consequence *BlockStatement, alternative *IfExpressionAlternative) *IfExpression {
	return &IfExpression{t, condition, consequence, alternative}
}

func (i *IfExpression) Token() token.Token {
	return i.token
}

func (i *IfExpression) Condition() Expression {
	return i.condition
}

func (i *IfExpression) Consequence() *BlockStatement {
	return i.consequence
}

func (i *IfExpression) Alternative() *IfExpressionAlternative {
	return i.alternative
}

func (i *IfExpression) expressionNode() {}

func (i *IfExpression) TokenLiteral() string {
	return i.token.Literal
}

func (i *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(if")
	out.WriteString(i.condition.String())
	out.WriteString(" ")
	out.WriteString(i.consequence.String())

	if i.alternative.Ok() {
		out.WriteString("else ")
		out.WriteString(i.alternative.String())
	}
	out.WriteString(")")

	return out.String()
}
