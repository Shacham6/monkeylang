package ast

import (
	"fmt"
	"monkey/token"
	"strings"
)

type HashLiteral struct {
	pairs map[Expression]Expression
	token token.Token // The '{' token
}

func NewHashLiteral(token token.Token, pairs map[Expression]Expression) *HashLiteral {
	return &HashLiteral{pairs, token}
}

func (h *HashLiteral) Pairs() map[Expression]Expression {
	return h.pairs
}

func (h *HashLiteral) TokenLiteral() string {
	return h.token.Literal
}

func (h *HashLiteral) expressionNode() {}

func (h *HashLiteral) String() string {
	pairs := []string{}
	for key, value := range h.pairs {
		pairs = append(pairs, fmt.Sprintf("(pair %s %s)", key.String(), value.String()))
	}

	return fmt.Sprintf("(hash %s)", strings.Join(pairs, " "))
}
