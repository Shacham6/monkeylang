package code_test

import (
	"monkey/code"
	"testing"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       code.Opcode
		operands []int
		expected []byte
	}{
		{
			code.OpConstant,
			[]int{0xfffffe}, // damn this is hard
			[]byte{byte(code.OpConstant), 0xff, 0xfe},
		},
	}

	for _, tt := range tests {
		instruction := code.Make(tt.op, tt.operands...)

		if len(instruction) != len(tt.expected) {
			t.Errorf("instruction has wrong length. want = %d, got = %d", len(tt.expected), len(instruction))
		}

		for i, expectedByte := range tt.expected {
			if instruction[i] != tt.expected[i] {
				t.Errorf("wrong byte at pos %d. want = %d, got = %d", i, expectedByte, instruction[i])
			}
		}
	}
}

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
