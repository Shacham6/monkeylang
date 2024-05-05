package parser_test

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/parser"
	"testing"
)

func checkParserErrors(t *testing.T, p *parser.Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	statementCount := len(program.Statements)
	if statementCount != 3 {
		t.Fatalf(
			"program.Statements does not contain 3 statements, got = %d",
			statementCount,
		)
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let', got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s is not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf(
			"letStmt.Name.Value not %s. got=%s",
			name,
			letStmt.Name.Value,
		)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf(
			"s.Name not '%s'. got='%s'",
			name,
			letStmt.Name,
		)
		return false
	}
	return true
}

func TestReturnStatement(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`
	l := lexer.New(input)
	p := parser.New(l)
	checkParserErrors(t, p)

	program := p.ParseProgram()
	if d := len(program.Statements); d != 3 {
		t.Fatalf("program.Statements does not contain 3 statements, got = %v", d)
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.returnStatement. got=%T", stmt)
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}
	}
}

func checkAmountOfStatements(t *testing.T, program *ast.Program, expected int) {
	if got := len(program.Statements); got != expected {
		t.Fatalf("program has wrong amount of statements. expected = %d, got = %d", expected, got)
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`
	p := parser.New(lexer.New(input))
	program := p.ParseProgram()

	checkParserErrors(t, p)
	checkAmountOfStatements(t, program, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ExpressionStatement like expected. Got %T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expression is not Identifier. got = %T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("identifier value not %s. got = %v", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if gotLen := len(program.Statements); gotLen != 1 {
		t.Fatalf("program has not enough statements. got = %d, expected = 1", gotLen)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionStatement. got = %T",
			program.Statements[0],
		)
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf(
			"stmt.Expression is not ast.IntegerLiteral. got = %T",
			stmt.Expression,
		)
	}

	if literal.Value != 5 {
		t.Errorf(
			"literal.Value not %s. got = %d",
			"5",
			literal.Value,
		)
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, tt := range prefixTests {
		parser := parser.New(lexer.New(tt.input))
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d elements. got = %d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got = %T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got = %T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator. got = %s. expect = %s", exp.Operator, tt.operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, expected int64) bool {
	integ, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("exp not IntegerLiteral. got = %T", exp)
		return false
	}

	if integ.Value != expected {
		t.Errorf("integ.Value. got = %d. expect = %d.", integ.Value, expected)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", expected) {
		t.Errorf("integ.TokenLiteral(). got = %s. expect = %d", integ.TokenLiteral(), expected)
		return false
	}

	return true
}
