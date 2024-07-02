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

	return evaluator.Eval(program, object.NewEnvironment())
}

func CheckIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result := testutils.CheckIsA[object.Integer](t, obj, "obj is not object.Integer")
	if result.Value != expected {
		t.Errorf("object has wrong value. got = %d, expect = %d",
			result.Value, expected)
		return false
	}
	return true
}

func CheckErrorObject(t *testing.T, obj object.Object, expectedMessage string) bool {
	result := testutils.CheckIsA[object.Error](t, obj, "obj is not object.Error")
	if result.Message != expectedMessage {
		t.Errorf("error has wrong message. got = %s, expect = %s",
			result.Message, expectedMessage)
		return false
	}
	return true
}

func CheckBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result := testutils.CheckIsA[object.Boolean](t, obj, "obj is not object.Boolean")
	if result.Value != expected {
		t.Errorf("object has wrong value. got = %v, expect = %v",
			result.Value, expected)
		return false
	}
	return true
}

func CheckStringObject(t *testing.T, obj object.Object, expected string) bool {
	result := testutils.CheckIsA[object.String](t, obj, "obj is not object.String")
	if result.Value != expected {
		t.Errorf("result has wrong value. got = %v, expect = %v", result.Value, expected)
		return false
	}
	return true
}

func CheckArrayValue(t *testing.T, obj object.Object, expected []CheckEvaluated) bool {
	arr := testutils.CheckIsA[object.Array](t, obj, "obj is not object.Array")
	hasFailed := false
	if len(expected) != len(arr.Elements) {
		t.Errorf("Lenth of arr is not the same expected. got = %d, want = %d", len(arr.Elements), len(expected))
		hasFailed = true
	}

	for i, elObj := range arr.Elements {
		if i >= len(expected) {
			t.Errorf(
				"Got in index [%d] an item whilst the expected finished previously already",
				i,
			)
			hasFailed = true
			continue
		}

		if checkResult := expected[i].CheckEvaluated(t, elObj); !checkResult {
			hasFailed = hasFailed && checkResult
		}
	}
	return hasFailed
}

func CheckNullObject(t *testing.T, obj object.Object) bool {
	testutils.CheckIsA[object.Null](t, obj, "obj is not object.Null")
	return true
}

type ResultInInt struct {
	n int64
}

func NewResultInInt(n int64) *ResultInInt {
	return &ResultInInt{n}
}

func (r *ResultInInt) CheckEvaluated(t *testing.T, obj object.Object) bool {
	return CheckIntegerObject(t, obj, r.n)
}

type ResultInError struct {
	message string
}

func NewResultInError(message string) *ResultInError {
	return &ResultInError{message}
}

func (r *ResultInError) CheckEvaluated(t *testing.T, obj object.Object) bool {
	return CheckErrorObject(t, obj, r.message)
}

type ResultInNil struct{}

func NewResultInNil() *ResultInNil {
	return &ResultInNil{}
}

func (r *ResultInNil) CheckEvaluated(t *testing.T, obj object.Object) bool {
	return CheckNullObject(t, obj)
}

type ResultInArray struct {
	elements []CheckEvaluated
}

func NewResultInArray(elements ...CheckEvaluated) *ResultInArray {
	return &ResultInArray{elements}
}

func (r *ResultInArray) CheckEvaluated(t *testing.T, obj object.Object) bool {
	return CheckArrayValue(t, obj, r.elements)
}

type ResultInString struct {
	expected string
}

func NewResultInString(expected string) *ResultInString {
	return &ResultInString{expected}
}

func (r *ResultInString) CheckEvaluated(t *testing.T, obj object.Object) bool {
	return CheckStringObject(t, obj, r.expected)
}

type CheckEvaluated interface {
	CheckEvaluated(t *testing.T, obj object.Object) bool
}
