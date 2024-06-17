package evaluator_test

import (
	"monkey/evaluator/internal/evaluatortest"
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
