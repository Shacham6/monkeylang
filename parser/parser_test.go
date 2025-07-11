package parser_test

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/parser"
	"monkey/testutils"
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
	switch v := interface{}(expected).(type) {
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
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := parser.New(lexer.New(tt.input))
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program.Statements does not contain 1 statements. got = %d",
					len(program.Statements))
			}

			stmt := program.Statements[0]
			if !testLetStatement(t, stmt, tt.expectedIdentifier) {
				return
			}

			val := stmt.(*ast.LetStatement).Value
			if !testLiteralExpression(t, val, tt.expectedValue) {
				return
			}
		})
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

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`
	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := testutils.CheckIsA[ast.ExpressionStatement](t, program.Statements[0], "progam.Statements[0] is not ast.ExpressionStatement")

	literal := testutils.CheckIsA[ast.StringLiteral](t, stmt.Expression, "stmt.Expression is not a ast.StringLiteral")

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got = %q", "hello world", literal.Value)
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

	if !testIdentifier(t, stmt.Expression, "foobar") {
		return
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
		input    string
		operator string
		value    interface{}
	}{
		{
			input:    "!5;",
			operator: "!",
			value:    5,
		},
		{
			input:    "-15;",
			operator: "-",
			value:    15,
		},
		{
			input:    "!true;",
			operator: "!",
			value:    true,
		},
		{
			input:    "!false;",
			operator: "!",
			value:    false,
		},
	}

	for _, tt := range prefixTests {
		t.Run(
			tt.input,
			func(t *testing.T) {
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

				if !testLiteralExpression(t, exp.Right, tt.value) {
					return
				}
			},
		)
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
		{"1 >= 2", 1, ">=", 2},
	}

	for _, tt := range infixTests {
		t.Run(
			tt.input,
			func(t *testing.T) {
				p := parser.New(lexer.New(tt.input))
				program := p.ParseProgram()
				checkParserErrors(t, p)

				if len(program.Statements) != 1 {
					t.Fatalf("len(program.Statements). got = %d. expect = %d", len(program.Statements), 1)
				}

				stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
				if !ok {
					t.Fatalf("program.Statements[0] is not ExpressionStatement, got = %T", program.Statements[0])
				}

				if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
					return
				}
			},
		)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "-a * b",
			expected: "(program (expr (infix (prefix - a) * b)))",
		},
		{
			input:    "!-a",
			expected: "(program (expr (prefix ! (prefix - a))))",
		},
		{
			input:    "a + b + c",
			expected: "(program (expr (infix (infix a + b) + c)))",
		},
		{
			input:    "a + b - c",
			expected: "(program (expr (infix (infix a + b) - c)))",
		},
		{
			input:    "a * b * c",
			expected: "(program (expr (infix (infix a * b) * c)))",
		},
		{
			input:    "a * b / c",
			expected: "(program (expr (infix (infix a * b) / c)))",
		},
		{
			input:    "a + b / c",
			expected: "(program (expr (infix a + (infix b / c))))",
		},
		{
			input:    "a + b * c + d / e - f",
			expected: "(program (expr (infix (infix (infix a + (infix b * c)) + (infix d / e)) - f)))",
		},
		{
			input:    "3 + 4; -5 * 5",
			expected: "(program (expr (infix 3 + 4)) (expr (infix (prefix - 5) * 5)))",
		},
		{
			input:    "5 > 4 == 3 < 4",
			expected: "(program (expr (infix (infix 5 > 4) == (infix 3 < 4))))",
		},
		{
			input:    "3 + 4 * 5 == 3 * 1 + 4 * 5",
			expected: "(program (expr (infix (infix 3 + (infix 4 * 5)) == (infix (infix 3 * 1) + (infix 4 * 5)))))",
		},
		{
			input:    "true",
			expected: "(program (expr true))",
		},
		{
			input:    "false",
			expected: "(program (expr false))",
		},
		{
			input:    "3 > 5 == false",
			expected: "(program (expr (infix (infix 3 > 5) == false)))",
		},
		{
			input:    "3 < 5 == true",
			expected: "(program (expr (infix (infix 3 < 5) == true)))",
		},
		{
			input:    "1 + (2 + 3) + 4",
			expected: "(program (expr (infix (infix 1 + (infix 2 + 3)) + 4)))",
		},
		{
			input:    "(5 + 5) * 2",
			expected: "(program (expr (infix (infix 5 + 5) * 2)))",
		},
		{
			input:    "2 / (5 + 5)",
			expected: "(program (expr (infix 2 / (infix 5 + 5))))",
		},
		{
			input:    "-(5 + 5)",
			expected: "(program (expr (prefix - (infix 5 + 5))))",
		},
		{
			input:    "!(true == true)",
			expected: "(program (expr (prefix ! (infix true == true))))",
		},
		{
			input:    "a + add(b * c) + d",
			expected: "(program (expr (infix (infix a + (call add (infix b * c))) + d)))",
		},
		{
			input:    "add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			expected: "(program (expr (call add a b 1 (infix 2 * 3) (infix 4 + 5) (call add 6 (infix 7 * 8)))))",
		},
		{
			input:    "add(a + b + c * d / f + g)",
			expected: "(program (expr (call add (infix (infix (infix a + b) + (infix (infix c * d) / f)) + g))))",
		},
		{
			input:    "a * [1, 2, 3, 4][b * c] * d",
			expected: "(program (expr (infix (infix a * (index [1 2 3 4] (infix b * c))) * d)))",
		},
		{
			input:    "add(a * b[2], b[1], 2 * [1, 2][1])",
			expected: "(program (expr (call add (infix a * (index b 2)) (index b 1) (infix 2 * (index [1 2] 1)))))",
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
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

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`
	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got = %d\n",
			1, program.Statements)
	}

	stmt := testutils.CheckIsA[ast.ExpressionStatement](t, program.Statements[0],
		"program.Statements[0] is not ast.ExpressionStatement")

	exp := testutils.CheckIsA[ast.IfExpression](t, stmt.Expression,
		"stmt.Expression is no ast.IfExpression")

	if !testInfixExpression(t, exp.Condition(), "x", "<", "y") {
		return
	}

	if lenConsequence := len(exp.Consequence().Statements()); lenConsequence != 1 {
		t.Errorf("consequence is not 1 statements. got = %d\n", lenConsequence)
	}

	consequence := testutils.CheckIsA[ast.ExpressionStatement](t, exp.Consequence().Statements()[0],
		"Statements[0] is not ExpressionStatement",
	)

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	alt, ok := exp.Alternative()

	if ok {
		t.Errorf("Got an unexpected alternative. got = %+v", alt)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`
	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got = %d\n",
			1, program.Statements)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got = %T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got = %T",
			stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition(), "x", "<", "y") {
		return
	}

	if lenConsequence := len(exp.Consequence().Statements()); lenConsequence != 1 {
		t.Errorf("consequence is not 1 statements. got = %d\n", lenConsequence)
	}

	consequence, ok := exp.Consequence().Statements()[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ExpressionStatement. got = %T",
			exp.Consequence().Statements()[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	alt, ok := exp.Alternative()

	if !ok {
		t.Errorf("Got an unexpected alternative. got = %+v", alt)
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`
	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got = %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("prgoram.Statements[0] is not a *ast.ExpressionStatement. got = %T", program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got = %T", stmt.Expression)
	}

	if len(function.Parameters()) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got = %d", len(function.Parameters()))
	}

	testLiteralExpression(t, function.Parameters()[0], "x")
	testLiteralExpression(t, function.Parameters()[1], "y")

	if len(function.Body().Statements()) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got = %d", len(function.Body().Statements()))
	}

	bodyStmt, ok := function.Body().Statements()[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement, got = %T", function.Body().Statements()[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{
			input:          "fn() {};",
			expectedParams: []string{},
		},
		{
			input:          "fn(x) {};",
			expectedParams: []string{"x"},
		},
		{
			input:          "fn(x, y, z) {};",
			expectedParams: []string{"x", "y", "z"},
		},
	}

	for _, tt := range tests {
		p := parser.New(lexer.New(tt.input))
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := testutils.CheckIsA[ast.ExpressionStatement](t, program.Statements[0], "program.Statements[0] is not an ast.ExpressionStatment")
		function := testutils.CheckIsA[ast.FunctionLiteral](t, stmt.Expression, "stmt.Expression is not a ast.FunctionLiteral")

		if len(function.Parameters()) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want = %d, got = %d", len(tt.expectedParams), len(function.Parameters()))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters()[i], ident)
		}
	}
}

func TestNamedFunctionParsing(t *testing.T) {
	input := `let something = fn() {}`
	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	fmt.Printf("%+v", program)
	if len(program.Statements) != 1 {
		t.Fatalf("wrong num of statements, want = %d, got = %d", 4, len(program.Statements))
	}

	letStmt := testutils.CheckIsA[ast.LetStatement](
		t, program.Statements[0],
		"program.Statements[0] is not ast.LetStatement as expected")

	fn := testutils.CheckIsA[ast.FunctionLiteral](
		t, letStmt.Value,
		"letStmt is not ast.FunctionLiteral as expected",
	)

	fnName, ok := fn.Name()
	if !ok {
		t.Fatalf("expected fn to have name but it does not")
	}
	if fnName != "something" {
		t.Fatalf("fn has wrong Name(), want = %s, got = %s", "something", fnName)
	}

	// checkParserErrors(t, p)

	// testutils.CheckIsA[]()
}

func TestCallExpressionParsing(t *testing.T) {
	input := `add(1, 2 * 3,  4 + 5);`

	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements, got = %d", len(program.Statements))
	}

	stmt := testutils.CheckIsA[ast.ExpressionStatement](t, program.Statements[0], "stmt is not a ast.ExpressionStatement")

	exp := testutils.CheckIsA[ast.CallExpression](t, stmt.Expression, "stmt.Expression is not a ast.CallExpression")

	if !testIdentifier(t, exp.Function(), "add") {
		return
	}

	arguments := exp.Arguments()

	if len(arguments) != 3 {
		t.Fatalf("wrong length of arguments. got = %d", len(exp.Arguments()))
	}

	testLiteralExpression(t, arguments[0], 1)
	testInfixExpression(t, arguments[1], 2, "*", 3)
	testInfixExpression(t, arguments[2], 4, "+", 5)
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 3, 3 + 3]"
	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := testutils.CheckIsA[ast.ExpressionStatement](t, program.Statements[0],
		"program.Statements[0] is not ast.ExpressionStatement")
	array := testutils.CheckIsA[ast.ArrayLiteral](t, stmt.Expression,
		"stmt.Expression is not ast.ArrayLiteral")
	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) is not 3. got = %d", len(array.Elements))
	}
	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 3)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpression(t *testing.T) {
	input := "array[1 + 1]"
	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := testutils.CheckIsA[ast.ExpressionStatement](t, program.Statements[0],
		"program.Statements[0] is not ast.ExpressionStatement")

	indexExp := testutils.CheckIsA[ast.IndexExpression](t, stmt.Expression,
		"stmt.Expression is not ast.IndexExpression")

	if !testIdentifier(t, indexExp.Left(), "array") {
		return
	}

	if !testInfixExpression(t, indexExp.Index(), 1, "+", 1) {
		return
	}
}

func TestParsingHashLiteralStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := testutils.CheckIsA[ast.ExpressionStatement](t, program.Statements[0], "program.Statements[0] is not ast.ExpressionStatement")
	hashLit := testutils.CheckIsA[ast.HashLiteral](t, stmt.Expression, "stmt.Expression is not ast.HashLiteral")

	if len(hashLit.Pairs()) != 3 {
		t.Errorf("hashLit.Pairs() has wrong length, got = %d, want = %d", len(hashLit.Pairs()), 3)
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, value := range hashLit.Pairs() {
		literal := testutils.CheckIsA[ast.StringLiteral](t, key, "key is not ast.StringLiteral")
		expectedValue := expected[literal.Value]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	input := "{}"
	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := testutils.CheckIsA[ast.ExpressionStatement](t, program.Statements[0], "program.Statements[0] is not ast.ExpressionStatement")
	hashLit := testutils.CheckIsA[ast.HashLiteral](t, stmt.Expression, "stmt.Expression is not ast.HashLiteral")

	if len(hashLit.Pairs()) != 0 {
		t.Errorf("hashLit.Pairs() has wrong length, got = %d, expect = %d", len(hashLit.Pairs()), 0)
	}
}

func TestParsingHashLiteralWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8}`
	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := testutils.CheckIsA[ast.ExpressionStatement](t, program.Statements[0], "program.Statements[0] is not ast.ExpressionStatement")
	hashLit := testutils.CheckIsA[ast.HashLiteral](t, stmt.Expression, "stmt.Expression is not ast.HashLiteral")

	expected := map[string]func(expr ast.Expression){
		"one": func(expr ast.Expression) {
			testInfixExpression(t, expr, 0, "+", 1)
		},
		"two": func(expr ast.Expression) {
			testInfixExpression(t, expr, 10, "-", 8)
		},
	}

	for key, value := range hashLit.Pairs() {
		keyString := testutils.CheckIsA[ast.StringLiteral](t, key, "key is not ast.StringLiteral")
		expectedFunc := expected[keyString.Value]
		expectedFunc(value)
	}
}

func TestMacroLiteralParsing(t *testing.T) {
	input := `macro(x, y) { x + y; };`

	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d elements. got = %d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got = %T", program.Statements[0])
	}

	macro, ok := stmt.Expression.(*ast.MacroLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.MacroLiteral, got = %T", stmt.Expression)
	}

	testLiteralExpression(t, macro.Parameters()[0], "x")
	testLiteralExpression(t, macro.Parameters()[1], "y")

	if amountOfStatements := len(macro.Body().Statements()); amountOfStatements != 1 {
		t.Fatalf("len(macro.Body().Statements()) is not %d, got = %d", 1, amountOfStatements)
	}

	exprStatement, ok := macro.Body().Statements()[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("macro.Body.Statements[0] is not *ast.ExpressionStatement, got = %T", macro.Body().Statements()[0])
	}

	testInfixExpression(t, exprStatement.Expression, "x", "+", "y")
}
