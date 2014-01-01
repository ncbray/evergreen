package runtime

// TODO flow type?

const (
	NORMAL = iota
	FAIL
	EXCEPTION
)

type State struct {
	Stream  []rune
	Index   int
	Deepest int
	Flow    int
}

func (state *State) Checkpoint() int {
	return state.Index
}

func (state *State) Recover(index int) {
	state.Index = index
	state.Flow = NORMAL
}

func (state *State) Read() (r rune) {
	if state.Index < len(state.Stream) {
		r = state.Stream[state.Index]
		state.Index += 1
	} else {
		state.Fail()
	}
	return
}

func (state *State) Peek() (r rune) {
	if state.Index < len(state.Stream) {
		return state.Stream[state.Index]
	} else {
		state.Fail()
		return 0
	}
}

func (state *State) Consume() {
	state.Index += 1
}

func (state *State) Slice(start int) string {
	return string(state.Stream[start:state.Index])
}

func (state *State) Fail() {
	if state.Index > state.Deepest {
		state.Deepest = state.Index
	}
	state.Flow = FAIL
}

// Utility function for generated tests
func MakeState(input string) *State {
	return &State{Stream: []rune(input)}
}
