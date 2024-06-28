package ast

import (
	"bytes"
	"fmt"
	"monkey/token"
)

type BlockStatement struct {
	token      token.Token
	statements []Statement
}

func NewBlockStatement(token token.Token, statements []Statement) *BlockStatement {
	return &BlockStatement{token, statements}
}

func (b *BlockStatement) Statements() []Statement {
	return b.statements
}

func (*BlockStatement) statementNode() {}

func (b *BlockStatement) TokenLiteral() string {
	return b.token.Literal
}

func (b *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range b.statements {
		out.WriteString(s.String())
	}

	return fmt.Sprintf("(block %s)", out.String())
}
