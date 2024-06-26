package evaluator_test

import (
	"monkey/evaluator/internal/evaluatortest"
	"monkey/object"
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
		{"-5", -5},
		{"-10", -10},
		{"5 + 5", 10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := evaluatortest.DoEval(tt.input)
		evaluatortest.CheckIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"1 == 1", true},
		{"2 > 1", true},
		{"2 < 1", false},
		{"1 != 1", false},
		{"1 >= 1", true},
		{"1 <= 1", true},
		{"1 != 2", true},
		{"(1 + 2) == 3", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := evaluatortest.DoEval(tt.input)
			evaluatortest.CheckBooleanObject(t, evaluated, tt.expected)
		})
	}
}

func TestEvalNullExpression(t *testing.T) {
	evaluated := evaluatortest.DoEval("null")
	evaluatortest.CheckNullObject(t, evaluated)
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := evaluatortest.DoEval(tt.input)
		evaluatortest.CheckBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := evaluatortest.DoEval(tt.input)
			integer, ok := tt.expected.(int)
			if ok {
				evaluatortest.CheckIntegerObject(t, evaluated, int64(integer))
			} else {
				evaluatortest.CheckNullObject(t, evaluated)
			}
		})
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 1;", 1},
		{"return 1; 2", 1},
		{"return 2 * 3; 1", 6},
		{"1; return 2; 3", 2},
		{`
			if (1 > 0) {
				if (2 > 1) {
					return 10;
				}
				return 1;
			}
		`, 10},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := evaluatortest.DoEval(tt.input)
			evaluatortest.CheckIntegerObject(t, evaluated, tt.expected)
		})
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
			if (10 > 1) {
				if (10 > 2) {
					return true + false;
				}
			}
			`, "unknown operator: BOOLEAN + BOOLEAN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := evaluatortest.DoEval(tt.input)

			errObj := testutils.CheckIsA[object.Error](t, evaluated, "evaluated is not an error object.")

			if errObj.Message != tt.expectedMessage {
				t.Errorf("wrong error message. expected = %q, got = %q", tt.expectedMessage, errObj.Message)
			}
		})
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = 5; b + a;", 10},
		{"let a = 5; let b = a; b;", 5},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluatortest.CheckIntegerObject(
				t,
				evaluatortest.DoEval(tt.input),
				tt.expected,
			)
		})
	}
}

func TestFunctionObject(t *testing.T) {
	input := `fn(x) {x + 2;};`
	evaluated := evaluatortest.DoEval(input)
	fn := testutils.CheckIsA[object.Function](t, evaluated, "evaluated is not a 'object.Function'")

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters = %+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x', got = %q", fn.Parameters[0])
	}

	expectedBody := "(x + 2);"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got = %q", expectedBody, fn.Body.String())
	}
}
