package dub

// TODO flow type?

const (
	NORMAL = iota
	FAIL
)

type DubState struct {
	Stream  []rune
	Index   int
	Deepest int
	Flow    int
}

func (frame *DubState) Checkpoint() int {
	return frame.Index
}

func (frame *DubState) Recover(index int) {
	frame.Index = index
	frame.Flow = NORMAL
}

func (frame *DubState) Read() (r rune) {
	if frame.Index < len(frame.Stream) {
		r = frame.Stream[frame.Index]
		frame.Index += 1
	} else {
		frame.Fail()
	}
	return
}

func (frame *DubState) Fail() {
	if frame.Index > frame.Deepest {
		frame.Deepest = frame.Index
	}
	frame.Flow = FAIL
}
