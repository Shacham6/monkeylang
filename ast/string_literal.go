package ast

import (
	"fmt"
	"monkey/token"
)

type StringLiteral struct {
	Token token.Token
	Value string
}

func NewStringLiteral(token token.Token, value string) *StringLiteral {
	return &StringLiteral{token, value}
}

func (*StringLiteral) expressionNode() {}

func (s *StringLiteral) TokenLiteral() string {
	return s.Token.Literal
}

func (s *StringLiteral) String() string {
	return fmt.Sprintf(`"%s"`, s.Value)
}
