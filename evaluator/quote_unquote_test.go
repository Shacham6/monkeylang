package evaluator_test

import (
	. "monkey/evaluator/internal/evaluatortest"
	"monkey/object"
	"monkey/testutils"
	"testing"
)

func TestQuote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`quote(5)`, `QUOTE(5)`},
		{`quote(foobar)`, `QUOTE(foobar)`},
		{`quote(1 + 2)`, `QUOTE((infix 1 + 2))`},
		{`quote(foobar + barfoo)`, `QUOTE((infix foobar + barfoo))`},
		{`quote(-10)`, `QUOTE((prefix - 10))`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := DoEval(tt.input)
			quote := testutils.CheckIsA[object.Quote](t, evaluated, "evaluated is not object.Quote")

			if quote.Node == nil {
				t.Fatalf("quote.Node is nil")
			}

			asString := quote.Inspect()
			if asString != tt.expected {
				t.Errorf("quote.Inspect() is not as expected. got = %s, want = %s", asString, tt.expected)
			}
		})
	}
}
