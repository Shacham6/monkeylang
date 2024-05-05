package ast

import (
	"fmt"
	"monkey/token"
	"strings"
)

type Node interface {
	String() string
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out strings.Builder
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type LetStatement struct {
	Token token.Token // The LET token.
	Name  *Identifier
	Value Expression
}

func NewLetStatement(token token.Token, name *Identifier, value Expression) *LetStatement {
	return &LetStatement{token, name, value}
}

func (ls *LetStatement) statementNode() {}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LetStatement) String() string {
	return fmt.Sprintf("%s %s = %s;", ls.TokenLiteral(), ls.Name.TokenLiteral(), ls.Value.String())
}

type ReturnStatement struct {
	Token token.Token // the RETURN token.
	Value Expression
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	return fmt.Sprintf("%s %s;", rs.TokenLiteral(), rs.Value.String())
}

type Identifier struct {
	Token token.Token // the IDENT token.
	Value string
}

func NewIdentifier(t token.Token, value string) *Identifier {
	return &Identifier{t, value}
}

func (*Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

type IntegerLiteral struct {
	Token token.Token // the INT token.
	Value int64
}

func NewIntegerLiteral(t token.Token, value int64) *IntegerLiteral {
	return &IntegerLiteral{t, value}
}

func (*IntegerLiteral) expressionNode() {}

func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

func (i IntegerLiteral) String() string {
	return i.Token.Literal
}

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (*ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	return fmt.Sprintf("%s;", es.Expression.String())
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (*PrefixExpression) expressionNode() {}

func (p *PrefixExpression) TokenLiteral() string {
	return p.Token.Literal
}

func (p *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", p.Operator, p.Right.String())
}
