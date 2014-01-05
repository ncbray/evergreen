package framework

import (
	"fmt"
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

type LocationProvider interface {
	AddFile(filename string, stream []rune)
	GetLocationInfo(pos int) (filename string, line int, col int, text string)
}

type simpleProvider struct {
	filename string
	stream   []rune
	lines    []int
}

func (p *simpleProvider) AddFile(filename string, stream []rune) {
	p.filename = filename
	p.stream = stream
	p.lines = FindLines(stream)
}

func (p *simpleProvider) GetLocationInfo(pos int) (string, int, int, string) {
	line, col, text := GetLocation(p.stream, p.lines, pos)
	return p.filename, line, col, text
}

type Status interface {
	CreateChild() Status
	Error(format string, a ...interface{})
	LocationError(loc int, message string)
	ShouldHalt() bool
	ErrorCount() int
}

type simpleStatus struct {
	parent     *simpleStatus
	loc        LocationProvider
	errorCount int
}

func (status *simpleStatus) CreateChild() Status {
	return &simpleStatus{parent: status, loc: status.loc}
}

func (status *simpleStatus) incErrorCount() {
	status.errorCount += 1
	if status.parent != nil {
		status.parent.incErrorCount()
	}
}

func (status *simpleStatus) Error(format string, a ...interface{}) {
	fmt.Printf("ERROR %s\n", fmt.Sprintf(format, a...))
	status.incErrorCount()
}
func (status *simpleStatus) LocationError(loc int, message string) {
	filename, line, col, text := status.loc.GetLocationInfo(loc)
	fmt.Printf("ERROR %s:%d:%d: %s\n", filename, line+1, col, message)
	fmt.Printf("    %s\n", text)
	// TODO tabs?
	arrow := []rune{}
	for i := 0; i < col; i++ {
		arrow = append(arrow, ' ')
	}
	arrow = append(arrow, '^')
	fmt.Printf("    %s\n", string(arrow))
	status.incErrorCount()
}

func (status *simpleStatus) ShouldHalt() bool {
	return status.errorCount > 0
}

func (status *simpleStatus) ErrorCount() int {
	return status.errorCount
}

func MakeProvider() LocationProvider {
	return &simpleProvider{}
}

func MakeStatus(loc LocationProvider) Status {
	return &simpleStatus{loc: loc}
}
