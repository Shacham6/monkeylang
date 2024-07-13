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

			if quote.Inspect() != tt.expected {
				t.Errorf("quote.Inspect() is not as expected. got = %s, want = %s", quote.Inspect(), tt.expected)
			}
		})
	}
}

func TestQuoteUnquote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`quote(unquote(4))`,
			`QUOTE(4)`,
		},
		{
			`quote(unquote(1 + 2))`,
			`QUOTE(3)`,
		},
		{
			`quote(1 + unquote(2 + 3))`,
			`QUOTE((infix 1 + 5))`,
		},
		{
			`quote(unquote(1 + 2) + 3)`,
			`QUOTE((infix 3 + 3))`,
		},
		{
			`quote(unquote(true))`,
			`QUOTE(true)`,
		},
		{
			`quote(unquote(false))`,
			`QUOTE(false)`,
		},
		{
			`quote(unquote(!true))`,
			`QUOTE(false)`,
		},
		{
			`quote(unquote("hello " + "world"))`,
			`QUOTE("hello world")`,
		},
		{
			`quote(unquote(null))`,
			`QUOTE(null)`,
		},
		{
			`quote(unquote(quote(4 + 4)))`,
			`QUOTE((infix 4 + 4))`,
		},
		{
			`
			let quotedInfixExpression = quote(4 + 4);
			quote(unquote(4 + 4) + unquote(quotedInfixExpression))
			`,
			`QUOTE((infix 8 + (infix 4 + 4)))`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := DoEval(tt.input)
			quote := testutils.CheckIsA[object.Quote](t, evaluated, "evaluated is not object.Quote")

			if quote.Node == nil {
				t.Fatalf("quote.Node is nil")
			}

			if quote.Inspect() != tt.expected {
				t.Errorf("quote.Inspect() not equal to expected. got = %s, want = %s", quote.Inspect(), tt.expected)
			}
		})
	}
}
