package evaluator_test

import (
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/testutils"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	p := parser.New(lexer.New(input))
	program := p.ParseProgram()

	return evaluator.Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result := testutils.CheckIsA[object.Integer](t, obj, "obj is not object.Integer")
	if result.Value != expected {
		t.Fatalf("object has wrong value. got = %d, expect = %d",
			result.Value, expected)
		return false
	}
	return false
}
