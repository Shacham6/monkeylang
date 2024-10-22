package compiler

import "monkey/code"

type CompilationScope struct {
	Instructions    code.Instructions
	LastInstruction EmittedInstruction
	PrevInstruction EmittedInstruction
}

func NewCompilationScope() CompilationScope {
	return CompilationScope{
		Instructions:    code.Instructions{},
		LastInstruction: ZeroEmittedInstruction(),
		PrevInstruction: ZeroEmittedInstruction(),
	}
}

func (c *CompilationScope) Emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.AddInstruction(ins)
	c.UpdateLastInstruction(op, pos)
	return pos
}

func (c *CompilationScope) UpdateLastInstruction(op code.Opcode, pos int) {
	previous := c.LastInstruction
	last := EmittedInstruction{
		Opcode:   op,
		Position: pos,
	}
	c.PrevInstruction = previous
	c.LastInstruction = last
}

func (c *CompilationScope) RemoveLastInstruction() {
	c.Instructions = c.Instructions[:c.LastInstruction.Position]
	c.LastInstruction = c.PrevInstruction
}

// ReplaceInstruction replaces an instruction, defined at pos, with the newInstruction.
//
// Please be careful using this. Unless the instruction at pos and the new instruction
// are of the same size, you are risking an overflow. Tread carefully.
func (c *CompilationScope) ReplaceInstruction(pos int, newInstruction []byte) {
	for i := 0; i < len(newInstruction); i++ {
		c.Instructions[pos+i] = newInstruction[i]
	}
}

// ChangeOperand changes the operand of the instruction at opPos.
func (c *CompilationScope) ChangeOperand(opPos int, operand int) {
	op := code.Opcode(c.Instructions[opPos])
	newInstruction := code.Make(op, operand)

	c.ReplaceInstruction(opPos, newInstruction)
}

func (c *CompilationScope) AddInstruction(ins []byte) int {
	posNewInstruction := len(c.Instructions)
	c.Instructions = append(c.Instructions, ins...)
	return posNewInstruction
}

func (c *CompilationScope) LastInstructionIs(op code.Opcode) bool {
	if len(c.Instructions) == 0 {
		return false
	}
	return c.LastInstruction.Opcode == op
}
