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
	}

	for _, tt := range tests {
		evaluated := evaluatortest.DoEval(tt.input)
		evaluatortest.CheckBooleanObject(t, evaluated, tt.expected)
	}
}
