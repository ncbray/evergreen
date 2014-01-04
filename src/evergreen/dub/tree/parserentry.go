package tree

import (
	"evergreen/dub/runtime"
	"evergreen/framework"
	"fmt"
	"strconv"
)

func FindLines(stream []rune) []int {
	lines := []int{0}
	for i, r := range stream {
		if r == '\n' {
			lines = append(lines, i+1)
		}
	}
	return lines
}

func GetLine(stream []rune, lines []int, line int) string {
	start := lines[line]
	end := len(stream)
	if line+1 < len(lines) {
		end = lines[line+1]
	}
	// HACK trim newline
	for end > start && (stream[end-1] == '\n' || stream[end-1] == '\r') {
		end -= 1
	}
	return string(stream[start:end])
}

func GetLocation(stream []rune, lines []int, pos int) (int, int, string) {
	// Stupid linear search
	var line int
	// If we don't find it, it must be on the last line.
	for line = 0; line < len(lines)-1; line++ {
		if pos >= lines[line] && pos < lines[line+1] {
			break
		}
	}
	col := pos - lines[line]
	return line, col, GetLine(stream, lines, line)
}

func GetRuneName(stream []rune, pos int) string {
	if pos < len(stream) {
		return strconv.QuoteRune(stream[pos])
	} else {
		return "EOF"
	}
}

func ParseDub(filename string, data []byte, status framework.Status) *File {
	stream := []rune(string(data))
	state := &runtime.State{Stream: stream}
	f := ParseFile(state)
	if state.Flow == 0 {
		return f
	} else {
		pos := state.Deepest
		lines := FindLines(stream)
		line, col, text := GetLocation(stream, lines, pos)
		loc := framework.Location{Filename: filename, Line: line, Col: col, Text: text}
		status.LocationError(loc, fmt.Sprintf("Unexpected %s", GetRuneName(stream, pos)))
		return nil
	}
}
