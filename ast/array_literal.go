package ast

import (
	"fmt"
	"monkey/token"
	"strings"
)

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func NewArrayLiteral(token token.Token, elements []Expression) *ArrayLiteral {
	return &ArrayLiteral{token, elements}
}

func (*ArrayLiteral) expressionNode() {}

func (a *ArrayLiteral) TokenLiteral() string {
	return a.Token.Literal
}

func (a *ArrayLiteral) String() string {
	items := []string{}
	for _, el := range a.Elements {
		items = append(items, el.String())
	}

	return fmt.Sprintf("[%s]", strings.Join(items, " "))
}
