package parser

import "monkey/token"

type Precedence int

const (
	_                      = iota
	LOWEST      Precedence = iota
	EQUALS      Precedence = iota // == or !=
	LESSGREATER Precedence = iota // > or <
	SUM         Precedence = iota // - or +
	PRODUCT     Precedence = iota // / (slash) or *
	PREFIX      Precedence = iota // -X or !X
	CALL        Precedence = iota // myFunction(X)
)

var precedences = map[token.TokenType]Precedence{
	token.EQ:      EQUALS,
	token.NOT_EQ:  EQUALS,
	token.LT:      LESSGREATER,
	token.LT_EQ:   LESSGREATER,
	token.GT:      LESSGREATER,
	token.GT_EQ:   LESSGREATER,
	token.PLUS:    SUM,
	token.MINUS:   SUM,
	token.SLASH:   PRODUCT,
	token.ASTERIX: PRODUCT,
	token.LPAREN:  CALL,
}

func (p *Parser) peekPrecedence() Precedence {
	pr, ok := precedences[p.peekToken.Type]
	if !ok {
		return LOWEST
	}
	return pr
}

func (p *Parser) curPrecedence() Precedence {
	pr, ok := precedences[p.curToken.Type]
	if !ok {
		return LOWEST
	}
	return pr
}
