package ast_test

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/parser"
	"monkey/token"
	"testing"
)

func TestModify(t *testing.T) {
	turnOneIntoTwo := func(node ast.Node) ast.Node {
		integer, ok := node.(*ast.IntegerLiteral)
		if !ok {
			return node
		}

		if integer.Value != 1 {
			return integer
		}

		return ast.NewIntegerLiteral(token.New(token.INT, "2"), 2)
	}

	tests := []struct {
		input   string
		preMod  string
		postMod string
	}{
		{
			"1",
			"(program (expr 1))",
			"(program (expr 2))",
		},
		{
			"1 + 2",
			"(program (expr (infix 1 + 2)))",
			"(program (expr (infix 2 + 2)))",
		},
		{
			"2[1]",
			"(program (expr (index 2 1)))",
			"(program (expr (index 2 2)))",
		},
		{
			"1[2]",
			"(program (expr (index 1 2)))",
			"(program (expr (index 2 2)))",
		},
		{
			"[1, 2, 3]",
			"(program (expr [1 2 3]))",
			"(program (expr [2 2 3]))",
		},
		{
			"[1, 2][2]",
			"(program (expr (index [1 2] 2)))",
			"(program (expr (index [2 2] 2)))",
		},
		{
			"if (1 == 2) {}",
			"(program (expr (if (infix 1 == 2) (block ))))",
			"(program (expr (if (infix 2 == 2) (block ))))",
		},
		{
			"if (true) { 1 }",
			"(program (expr (if true (block (expr 1)))))",
			"(program (expr (if true (block (expr 2)))))",
		},
		{
			"if (true) {} else { 1 }",
			"(program (expr (if true (block ) (block (expr 1)))))",
			"(program (expr (if true (block ) (block (expr 2)))))",
		},
		{
			"return 1",
			"(program (return 1))",
			"(program (return 2))",
		},
		{
			"let a = 1",
			"(program (let a 1))",
			"(program (let a 2))",
		},
		{
			"let a = 1 + 2",
			"(program (let a (infix 1 + 2)))",
			"(program (let a (infix 2 + 2)))",
		},
		{
			"fn(){ 1 }",
			"(program (expr (func [] (block (expr 1)))))",
			"(program (expr (func [] (block (expr 2)))))",
		},
		{
			"fn() {2 + 1}",
			"(program (expr (func [] (block (expr (infix 2 + 1))))))",
			"(program (expr (func [] (block (expr (infix 2 + 2))))))",
		},
		{
			"return 1",
			"(program (return 1))",
			"(program (return 2))",
		},
		{
			"fn() { return 1 }",
			"(program (expr (func [] (block (return 1)))))",
			"(program (expr (func [] (block (return 2)))))",
		},
		{
			"return 1 + 2",
			"(program (return (infix 1 + 2)))",
			"(program (return (infix 2 + 2)))",
		},
		{
			"func(1)",
			"(program (expr (call func 1)))",
			"(program (expr (call func 2)))",
		},
		{
			"func(1 + 2)",
			"(program (expr (call func (infix 1 + 2))))",
			"(program (expr (call func (infix 2 + 2))))",
		},
		{
			`{"a": 1}`,
			`(program (expr (hash (pair a 1))))`,
			`(program (expr (hash (pair a 2))))`,
		},
		{
			`{1: "a"}`,
			`(program (expr (hash (pair 1 a))))`,
			`(program (expr (hash (pair 2 a))))`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := parser.New(lexer.New(tt.input))
			program := p.ParseProgram()

			checkParserErrors(t, p)

			if tt.preMod != program.String() {
				t.Errorf("parsing result unexpected. got = %s, want = %s", program.String(), tt.preMod)
			}

			modified := ast.Modify(program, turnOneIntoTwo)

			if tt.postMod != modified.String() {
				t.Errorf("modification result unexpected. got = %s, want = %s", modified.String(), tt.postMod)
			}
		})
	}
}

func checkParserErrors(t *testing.T, p *parser.Parser) {
	if len(p.Errors()) > 0 {
		t.Fatalf("parsing resulted in unexpected errors: %#v", p.Errors())
	}
}
