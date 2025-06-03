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

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn

	curToken  token.Token
	peekToken token.Token

	errors []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{ //nolint:exhaustruct
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = map[token.TokenType]prefixParseFn{
		token.IDENT:    p.parseIdentifier,
		token.INT:      p.parseIntegerLiteral,
		token.TRUE:     p.parseBoolean,
		token.FALSE:    p.parseBoolean,
		token.BANG:     p.parsePrefixExpression,
		token.MINUS:    p.parsePrefixExpression,
		token.LPAREN:   p.parseGroupedExpression,
		token.IF:       p.parseIfExpression,
		token.FUNCTION: p.parseFunctionLiteral,
		token.MACRO:    p.parseMacroLiteral,
		token.STRING:   p.parseStringLiteral,
		token.LBRACKET: p.parseArrayLiteral,
		token.LBRACE:   p.parseHashLiteral,
	}

	p.infixParseFns = map[token.TokenType]infixParseFn{
		token.PLUS:     p.parseInfixExpression,
		token.MINUS:    p.parseInfixExpression,
		token.SLASH:    p.parseInfixExpression,
		token.ASTERIX:  p.parseInfixExpression,
		token.EQ:       p.parseInfixExpression,
		token.NOT_EQ:   p.parseInfixExpression,
		token.LT:       p.parseInfixExpression,
		token.LT_EQ:    p.parseInfixExpression,
		token.GT:       p.parseInfixExpression,
		token.GT_EQ:    p.parseInfixExpression,
		token.LPAREN:   p.parseCallExpression,
		token.LBRACKET: p.parseIndexExpression,
	}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) parseHashLiteral() ast.Expression {
	curToken := p.curToken

	pairs := map[ast.Expression]ast.Expression{}

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)
		if !p.expectPeek(token.COLON) {
			return nil
		}
		p.nextToken()

		val := p.parseExpression(LOWEST)
		pairs[key] = val

		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}
	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return ast.NewHashLiteral(curToken, pairs)
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	curToken := p.curToken

	p.nextToken()

	index := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return ast.NewIndexExpression(
		curToken, left, index,
	)
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken} //nolint:exhaustruct
	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}
	return list
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	arguments := p.parseExpressionList(token.RPAREN)
	exp := ast.NewCallExpression(p.curToken, function, arguments)
	return exp
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{} //nolint:exhaustruct
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		program.Statements = append(program.Statements, stmt)
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
	stmt := &ast.ReturnStatement{Token: p.curToken} //nolint:exhaustruct
	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken} //nolint:exhaustruct

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = ast.NewIdentifier(p.curToken, p.curToken.Literal)

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	// Tailoring an assignment of name to a function in case it is a function!
	fn, ok := stmt.Value.(*ast.FunctionLiteral)
	if ok {
		fn.SetName(stmt.Name.Value)
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken} //nolint:exhaustruct
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

func (p *Parser) parseIfExpression() ast.Expression {
	tok := p.curToken

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	condition := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	consequence := p.parseBlockStatement()

	var alternative *ast.BlockStatement
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		alternative = p.parseBlockStatement()
	}

	return ast.NewIfExpression(tok, condition, consequence, ast.NewIfExpressionAlternative(alternative))
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	blockToken := p.curToken
	blockStatements := []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		blockStatements = append(blockStatements, stmt)
		p.nextToken()
	}

	return ast.NewBlockStatement(blockToken, blockStatements)
}

func (p *Parser) parseIdentifier() ast.Expression {
	return ast.NewIdentifier(p.curToken, p.curToken.Literal)
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
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

func (p *Parser) parseFunctionLiteral() ast.Expression {
	curToken := p.curToken // the "fn" token
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	params := p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	body := p.parseBlockStatement()

	return ast.NewFunctionLiteral(curToken, params, body, "")
}

func (p *Parser) parseMacroLiteral() ast.Expression {
	curToken := p.curToken
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	params := p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	body := p.parseBlockStatement()

	return ast.NewMacroLiteral(curToken, params, body)
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()

		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		if p.curToken.Type != token.IDENT {
			p.errors = append(p.errors, "argument in function definition must be an identifier")
			continue
		}

		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
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
	p.peekError(t)
	return false
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
