package code

import "encoding/binary"

type Opcode byte

func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		// opcode not found
		// better error handling later?
		return []byte{}
	}

	instructionLen := 1
	for _, operandWidth := range def.OperandWidths {
		instructionLen += operandWidth
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, operand := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(operand))

		case 1:
			instruction[offset] = byte(operand)
		}
		offset += width
	}

	return instruction
}
