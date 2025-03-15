package vm

import (
	"monkey/code"
	"monkey/object"
)

type Frame struct {
	fn          *object.CompiledFunction
	ip          int
	basePointer int
}

func NewFrame(fn *object.CompiledFunction, basePointer int) *Frame {
	ip := -1
	return &Frame{fn, ip, basePointer}
}

func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}
