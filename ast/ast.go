package ast

import (
	"fmt"
	"monkey/token"
	"strings"
)

type Node interface {
	modifiable

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
	statementLines := []string{}
	for _, s := range p.Statements {
		statementLines = append(statementLines, s.String())
	}
	return fmt.Sprintf("(program %s)", strings.Join(statementLines, " "))
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
	return fmt.Sprintf(
		"(let %s %s)",
		ls.Name.String(),
		ls.Value.String(),
	)
}

type ReturnStatement struct {
	Token       token.Token // the RETURN token.
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	return fmt.Sprintf("(return %s)", rs.ReturnValue.String())
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
	return fmt.Sprintf("(expr %s)", es.Expression.String())
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func NewPrefixExpression(t token.Token, operator string, right Expression) *PrefixExpression {
	return &PrefixExpression{t, operator, right}
}

func (*PrefixExpression) expressionNode() {}

func (p *PrefixExpression) TokenLiteral() string {
	return p.Token.Literal
}

func (p *PrefixExpression) String() string {
	return fmt.Sprintf("(prefix %s %s)", p.Operator, p.Right.String())
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func NewInfixExpression(t token.Token, left Expression, operator string, right Expression) *InfixExpression {
	return &InfixExpression{t, left, operator, right}
}

func (*InfixExpression) expressionNode() {}

func (i *InfixExpression) TokenLiteral() string {
	return i.Token.Literal
}

func (i *InfixExpression) String() string {
	return fmt.Sprintf("(infix %s %s %s)", i.Left.String(), i.Operator, i.Right.String())
}

type Boolean struct {
	token token.Token
	value bool
}

func NewBoolean(t token.Token, value bool) *Boolean {
	return &Boolean{t, value}
}

func (b *Boolean) Value() bool {
	return b.value
}

func (*Boolean) expressionNode() {}

func (b *Boolean) TokenLiteral() string {
	return b.token.Literal
}

func (b *Boolean) String() string {
	return b.TokenLiteral()
}
