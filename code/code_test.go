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
		{
			code.OpAdd,
			[]int{}, // no operands because operates on values that are already pushed in the stack.
			[]byte{byte(code.OpAdd)},
		},
		{
			code.OpGetLocal,
			[]int{255},
			[]byte{byte(code.OpGetLocal), 255},
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
