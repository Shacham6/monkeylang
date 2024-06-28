package ast

import (
	"fmt"
	"monkey/token"
)

type IndexExpression struct {
	left  Expression
	index Expression
	token token.Token
}

func NewIndexExpression(token token.Token, left Expression, index Expression) *IndexExpression {
	return &IndexExpression{left, index, token}
}

func (i *IndexExpression) Left() Expression {
	return i.left
}

func (i *IndexExpression) Index() Expression {
	return i.index
}

func (i *IndexExpression) TokenLiteral() string {
	return i.token.Literal
}

func (*IndexExpression) expressionNode() {}

func (i *IndexExpression) String() string {
	return fmt.Sprintf("(index %s %s)", i.left.String(), i.index.String())
	// var out bytes.Buffer
	//
	// out.WriteString("(")
	// out.WriteString(i.left.String())
	// out.WriteString("[")
	// out.WriteString(i.index.String())
	// out.WriteString("])")
	//
	// return out.String()
}
