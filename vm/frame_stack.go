package vm

type FramesStack struct {
	frames []*Frame

	// topIndex is the position of the top of the stack.
	//
	// This position can be used to signify the safest reading position.
	// This position+1 is the writing position.
	topIndex int
}

func NewFramesStack(size int) FramesStack {
	return FramesStack{
		frames:   make([]*Frame, size),
		topIndex: -1,
	}
}

func (fs *FramesStack) Current() *Frame {
	return fs.frames[fs.topIndex]
}

func (fs *FramesStack) Push(f *Frame) {
	fs.frames[fs.topIndex+1] = f
	fs.topIndex++
}

func (fs *FramesStack) Pop() *Frame {
	// for performance purposes I'm not going to check bounds, so just know that
	// popping too much/pushing too much can cause panic-y shit.
	f := fs.frames[fs.topIndex]
	fs.topIndex--
	return f
}
