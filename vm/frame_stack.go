package vm

type FramesStack struct {
	frames      []*Frame
	framesIndex int
}

func NewFramesStack(size int) FramesStack {
	return FramesStack{
		frames:      make([]*Frame, size),
		framesIndex: 0,
	}
}

func (fs *FramesStack) Current() *Frame {
	return fs.frames[fs.framesIndex]
}

func (fs *FramesStack) Push(f *Frame) {
	fs.frames = append(fs.frames, f)
	fs.framesIndex++
}

func (fs *FramesStack) Pop() *Frame {
	// for performance purposes I'm not going to check bounds, so just know that
	// popping too much/pushing too much can cause panic-y shit.
	f := fs.frames[fs.framesIndex]
	fs.framesIndex--
	return f
}
