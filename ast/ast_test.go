package ast_test

import (
	"monkey/ast"
	"monkey/token"
	"testing"
)

func TestAstToString(t *testing.T) {
	ls := ast.LetStatement{
		Token: token.Token{
			Type:    token.LET,
			Literal: "let",
		},
		Name:  &ast.Identifier{token.Token{Type: token.IDENT, Literal: "x"}, "x"},
		Value: &ast.Identifier{token.Token{Type: token.IDENT, Literal: "y"}, "y"},
	}

	s := ls.String()
	expect := "(let x y)"
	if s != expect {
		t.Fatalf("ast.String(), got = %s, expect = %s", s, expect)
	}
}
