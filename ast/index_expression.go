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

func (i *IndexExpression) SetLeft(value Expression) {
	i.left = value
}

func (i *IndexExpression) Index() Expression {
	return i.index
}

func (i *IndexExpression) SetIndex(value Expression) {
	i.index = value
}

func (i *IndexExpression) TokenLiteral() string {
	return i.token.Literal
}

func (*IndexExpression) expressionNode() {}

func (i *IndexExpression) String() string {
	return fmt.Sprintf("(index %s %s)", i.left.String(), i.index.String())
}
