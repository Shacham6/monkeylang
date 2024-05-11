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

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got = %T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got = %s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() not %s. got = %s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, b bool) bool {
	astBool, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("expr not *ast.Boolean. got = %T", exp)
		return false
	}

	if astBool.Value() != b {
		t.Errorf("astBool.Value(). got = %v. expect = %v", astBool.Value(), b)
		return false
	}

	literalExpect := fmt.Sprintf("%v", b)
	if astBool.TokenLiteral() != literalExpect {
		t.Errorf("astBool.TokenLiteral(). got = %v. expect = %v", astBool.TokenLiteral(), literalExpect)
		return false
	}

	return true
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case bool:
		return testBooleanLiteral(t, exp, v)
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	}

	t.Errorf("type of exp not handled. got = %T", exp)

	return false
}

func testInfixExpression(
	t *testing.T,
	exp ast.Expression,
	left interface{},
	operator string,
	right interface{},
) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not *ast.InfixExpression. got = %T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("opExp.Operator, got = %s, expect = %s", opExp.Operator, operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
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

func checkAmountOfStatements(t *testing.T, program *ast.Program, expected int) {
	if got := len(program.Statements); got != expected {
		t.Fatalf("program has wrong amount of statements. expected = %d, got = %d", expected, got)
	}
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

	if !testLiteralExpression(t, stmt.Expression, 5) {
		return
	}
}

func TestBooleanExpression(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		expect bool
	}{
		{
			name:   "parsing true",
			input:  "true;",
			expect: true,
		},
		{
			name:   "parsing false",
			input:  "false",
			expect: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.New(lexer.New(tt.input))
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

			if !testLiteralExpression(t, stmt.Expression, tt.expect) {
				return
			}
		})
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

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, infixTest := range infixTests {
		p := parser.New(lexer.New(infixTest.input))
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("len(program.Statements). got = %d. expect = %d", len(program.Statements), 1)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ExpressionStatement, got = %T", program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, infixTest.leftValue, infixTest.operator, infixTest.rightValue) {
			return
		}
		//
		// exp, ok := stmt.Expression.(*ast.InfixExpression)
		// if !ok {
		// 	t.Fatalf("exp is not InfixExpression. got = %T", stmt.Expression)
		// }
		//
		// if !testIntegerLiteral(t, exp.Left, infixTest.leftValue) {
		// 	return
		// }
		//
		// if exp.Operator != infixTest.operator {
		// 	t.Fatalf("exp.Operator. got = %s. expect = %s", exp.Operator, infixTest.operator)
		// }
		//
		// if !testIntegerLiteral(t, exp.Right, infixTest.rightValue) {
		// 	return
		// }
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b);",
		},
		{
			"!-a",
			"(!(-a));",
		},
		{
			"a + b + c",
			"((a + b) + c);",
		},
		{
			"a + b - c",
			"((a + b) - c);",
		},
		{
			"a * b * c",
			"((a * b) * c);",
		},
		{
			"a * b / c",
			"((a * b) / c);",
		},
		{
			"a + b / c",
			"(a + (b / c));",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f);",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4);((-5) * 5);",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4));",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)));",
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("TestOperatorPrecedenceParsing[%d]", i), func(t *testing.T) {
			p := parser.New(lexer.New(test.input))
			program := p.ParseProgram()
			checkParserErrors(t, p)

			got := program.String()
			if got != test.expected {
				t.Fatalf("program.String() got = %s, expect = %s", got, test.expected)
			}
		})
	}
}
