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
