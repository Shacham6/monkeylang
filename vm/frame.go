package vm

import (
	"monkey/code"
	"monkey/object"
)

type Frame struct {
	fn *object.CompiledFunction
	ip int
}

func NewFrame(fn *object.CompiledFunction) *Frame {
	ip := -1
	return &Frame{fn, ip}
}

func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}
