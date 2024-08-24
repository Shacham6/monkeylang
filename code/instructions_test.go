package code_test

import (
	"monkey/code"
	"testing"
)

func TestInstructionString(t *testing.T) {
	instructions := []code.Instructions{
		code.Make(code.OpConstant, 1),
		code.Make(code.OpConstant, 2),
		code.Make(code.OpConstant, 65535),
	}

	expected := `0000 OpConstant 1
0003 OpConstant 2
0006 OpConstant 65535
`

	concatted := code.Instructions{}
	for _, ins := range instructions {
		concatted = append(concatted, ins...)
	}

	if concatted.String() != expected {
		t.Errorf("instructions wrongly formatted.\n want = %q\n  got = %q", expected, concatted.String())
	}
}
