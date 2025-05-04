package vm

import (
	"monkey/code"
	"monkey/object"
)

type Frame struct {
	cl          *object.Closure
	ip          int
	basePointer int
}

func NewFrame(cl *object.Closure, basePointer int) *Frame {
	ip := -1
	return &Frame{cl, ip, basePointer}
}

func (f *Frame) Instructions() code.Instructions {
	return f.cl.Fn.Instructions
}
