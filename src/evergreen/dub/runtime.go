package dub

// TODO flow type?

const (
	// Real flows, used at runtime
	NORMAL = iota
	FAIL
	EXCEPTION
	// Virtual flows, only for graph construction
	RETURN
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

func (state *DubState) Peek() (r rune) {
	if state.Index < len(state.Stream) {
		return state.Stream[state.Index]
	} else {
		state.Fail()
		return 0
	}
}

func (state *DubState) Consume() {
	state.Index += 1
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
