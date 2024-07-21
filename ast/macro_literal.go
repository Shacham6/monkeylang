package ast

import (
	"fmt"
	"monkey/token"
	"strings"
)

type MacroLiteral struct {
	token      token.Token // the MACRO token
	parameters []*Identifier
	body       *BlockStatement
}

func NewMacroLiteral(token token.Token, parameters []*Identifier, body *BlockStatement) *MacroLiteral {
	return &MacroLiteral{token, parameters, body}
}

func (m *MacroLiteral) TokenLiteral() string {
	return m.token.Literal
}

func (m *MacroLiteral) expressionNode() {}

func (m *MacroLiteral) String() string {
	parametersStrings := []string{}
	for _, p := range m.parameters {
		parametersStrings = append(parametersStrings, p.String())
	}

	return fmt.Sprintf("(macro [%s] %s)", strings.Join(parametersStrings, " "), m.body.String())
}

func (m *MacroLiteral) Parameters() []*Identifier {
	return m.parameters
}

func (m *MacroLiteral) Body() *BlockStatement {
	return m.body
}
