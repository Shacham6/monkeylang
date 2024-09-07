package code_test

import (
	"fmt"
	"monkey/code"
	"testing"
)

func TestInstructionString(t *testing.T) {
	tests := []struct {
		instructions []code.Instructions
		expected     string
	}{
		{
			[]code.Instructions{
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 65535),
			},
			`0000 OpConstant 1
0003 OpConstant 2
0006 OpConstant 65535
`,
		},
		{
			[]code.Instructions{
				code.Make(code.OpAdd),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
			},
			`0000 OpAdd
0001 OpConstant 2
0004 OpConstant 3
`,
		},
	}

	for index, tt := range tests {
		t.Run(fmt.Sprintf("TestInstructionString[%d]", index), func(t *testing.T) {
			concatted := code.Instructions{}
			for _, ins := range tt.instructions {
				concatted = append(concatted, ins...)
			}

			if concatted.String() != tt.expected {
				t.Errorf("instructions wrongly formatted.\n want = %q\n  got = %q", tt.expected, concatted.String())
			}
		})
	}
}
