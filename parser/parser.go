package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

type prefixParseFn func() ast.Expression

type infixParseFn func(ast.Expression) ast.Expression

func Parse(input string) (*ast.Program, []string) {
	p := New(lexer.New(input))
	program := p.ParseProgram()
	return program, p.errors
}

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = map[token.TokenType]prefixParseFn{
		token.IDENT:  p.parseIdentifier,
		token.INT:    p.parseIntegerLiteral,
		token.TRUE:   p.parseBoolean,
		token.FALSE:  p.parseBoolean,
		token.BANG:   p.parsePrefixExpression,
		token.MINUS:  p.parsePrefixExpression,
		token.LPAREN: p.parseGroupedExpression,
	}

	p.infixParseFns = map[token.TokenType]infixParseFn{
		token.PLUS:    p.parseInfixExpression,
		token.MINUS:   p.parseInfixExpression,
		token.SLASH:   p.parseInfixExpression,
		token.ASTERIX: p.parseInfixExpression,
		token.EQ:      p.parseInfixExpression,
		token.NOT_EQ:  p.parseInfixExpression,
		token.LT:      p.parseInfixExpression,
		token.GT:      p.parseInfixExpression,
	}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken() // Skip the start i.e "("

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil // Should we return an error? IS this even an error?
	}

	return exp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	token := p.curToken
	operator := p.curToken.Literal

	p.nextToken()

	right := p.parseExpression(PREFIX)

	return ast.NewPrefixExpression(token, operator, right)
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	token := p.curToken
	operator := p.curToken.Literal

	precedence := p.curPrecedence()
	p.nextToken()
	right := p.parseExpression(precedence)

	return ast.NewInfixExpression(token, left, operator, right)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	// TODO: We're skipping the expressions until we encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// We're skipping the expressions until we
	// encounter as semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) pushNoPrefixParseFnError(tokenType token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", tokenType)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(precedence Precedence) ast.Expression {
	prefix, ok := p.prefixParseFns[p.curToken.Type]

	if !ok {
		p.pushNoPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix, ok := p.infixParseFns[p.peekToken.Type]
		if !ok {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return ast.NewIdentifier(p.curToken, p.curToken.Literal)
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	strLiteral := p.curToken.Literal
	value, err := strconv.ParseInt(strLiteral, 0, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("Could not parse %q as integer", strLiteral))
		return nil
	}
	return ast.NewIntegerLiteral(p.curToken, value)
}

func (p *Parser) parseBoolean() ast.Expression {
	return ast.NewBoolean(p.curToken, p.curTokenIs(token.TRUE))
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	return false
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
