package evaluatortest

import (
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/testutils"
	"testing"
)

func DoEval(input string) object.Object {
	p := parser.New(lexer.New(input))
	program := p.ParseProgram()

	return evaluator.Eval(program)
}

func CheckIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result := testutils.CheckIsA[object.Integer](t, obj, "obj is not object.Integer")
	if result.Value != expected {
		t.Fatalf("object has wrong value. got = %d, expect = %d",
			result.Value, expected)
		return false
	}
	return true
}

func CheckBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result := testutils.CheckIsA[object.Boolean](t, obj, "obj is not object.Boolean")
	if result.Value != expected {
		t.Fatalf("object has wrong value. got = %v, expect = %v",
			result.Value, expected)
		return false
	}
	return true
}

func CheckNullObject(t *testing.T, obj object.Object) bool {
	testutils.CheckIsA[object.Null](t, obj, "obj is not object.Null")
	return true
}
