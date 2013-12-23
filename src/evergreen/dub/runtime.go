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

func (state *DubState) Checkpoint() int {
	return state.Index
}

func (state *DubState) Recover(index int) {
	state.Index = index
	state.Flow = NORMAL
}

func (state *DubState) Read() (r rune) {
	if state.Index < len(state.Stream) {
		r = state.Stream[state.Index]
		state.Index += 1
	} else {
		state.Fail()
	}
	return
}

func (state *DubState) Slice(start int) string {
	return string(state.Stream[start:state.Index])
}

func (state *DubState) Fail() {
	if state.Index > state.Deepest {
		state.Deepest = state.Index
	}
	state.Flow = FAIL
}
