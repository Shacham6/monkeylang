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

func (*ArrayLiteral) expressionNode() {}

func (a *ArrayLiteral) TokenLiteral() string {
	return a.Token.Literal
}

func (a *ArrayLiteral) String() string {
	// var out bytes.Buffer
	//
	// elements := []string{}
	// for _, el := range a.Elements {
	// 	elements = append(elements, el.String())
	// }
	// out.WriteString("[")
	// out.WriteString(strings.Join(elements, ", "))
	// out.WriteString("]")
	//
	// return out.String()
	items := []string{}
	for _, el := range a.Elements {
		items = append(items, el.String())
	}

	return fmt.Sprintf("[%s]", strings.Join(items, " "))
}
