package tree

import (
	"evergreen/dub/runtime"
	"fmt"
	"io/ioutil"
)

func FindLines(stream []rune) []int {
	lines := []int{}
	for i, r := range stream {
		if r == '\n' {
			lines = append(lines, i+1)
		}
	}
	lines = append(lines, len(stream))
	return lines
}

func PrintError(filename string, deepest int, stream []rune, lines []int) {
	// Stupid linear search
	var line int
	var col int
	var start int
	var end int
	for i, s := range lines {
		start = end
		end = s
		line = i + 1
		if deepest < end {
			break
		}
	}
	col = deepest - start

	// HACK trim newline
	for end > start && (end > len(stream) || stream[end-1] == '\n' || stream[end-1] == '\t') {
		end -= 1
	}
	text := string(stream[start:end])
	fmt.Printf("%s: Unexpected @ %d:%d\n%s\n", filename, line, col, text)
	// TODO tabs?
	arrow := []rune{}
	for i := 0; i < col; i++ {
		arrow = append(arrow, ' ')
	}
	arrow = append(arrow, '^')
	fmt.Println(string(arrow))
}

func ParseDub(filename string) *File {
	data, _ := ioutil.ReadFile(filename)
	stream := []rune(string(data))
	state := &runtime.State{Stream: stream}
	f := ParseFile(state)
	if state.Flow != 0 {
		lines := FindLines(stream)
		PrintError(filename, state.Deepest, stream, lines)
		panic(state.Deepest)
	}
	return f
}
