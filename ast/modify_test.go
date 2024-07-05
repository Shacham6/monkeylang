package ast_test

import (
	"monkey/ast"
	"monkey/token"
	"reflect"
	"testing"
)

func TestModify(t *testing.T) {
	one := func() ast.Expression { return ast.NewIntegerLiteral(token.New(token.INT, "1"), 1) }
	two := func() ast.Expression { return ast.NewIntegerLiteral(token.New(token.INT, "2"), 2) }

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
		name     string
		input    ast.Node
		expected ast.Node
	}{
		{
			"flat",
			one(),
			two(),
		},
		{
			"nested",
			&ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{Expression: one()},
				},
			},
			&ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{Expression: two()},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modified := ast.Modify(tt.input, turnOneIntoTwo)
			equal := reflect.DeepEqual(modified, tt.expected)
			if !equal {
				t.Errorf("not equal. got = %#v, want = %#v", modified, tt.expected)
			}
		})
	}
}
