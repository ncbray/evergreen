package tree

import (
	"evergreen/dub/runtime"
	"evergreen/framework"
	"fmt"
	"strconv"
)

func GetRuneName(stream []rune, pos int) string {
	if pos < len(stream) {
		return strconv.QuoteRune(stream[pos])
	} else {
		return "EOF"
	}
}

func ParseDub(data []byte, status framework.Status) *File {
	stream := []rune(string(data))
	state := &runtime.State{Stream: stream}
	f := ParseFile(state)
	if state.Flow == 0 {
		return f
	} else {
		pos := state.Deepest
		status.LocationError(pos, fmt.Sprintf("Unexpected %s", GetRuneName(stream, pos)))
		return nil
	}
}
