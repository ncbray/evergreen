package framework

import (
	"fmt"
)

type Location struct {
	Filename string
	Line     int
	Col      int
	Text     string
}

type Status interface {
	CreateChild() Status
	Error(message string)
	LocationError(loc Location, message string)
	ShouldHalt() bool
	ErrorCount() int
}

type simpleStatus struct {
	parent     *simpleStatus
	errorCount int
}

func (status *simpleStatus) CreateChild() Status {
	return &simpleStatus{parent: status}
}

func (status *simpleStatus) incErrorCount() {
	status.errorCount += 1
	if status.parent != nil {
		status.parent.incErrorCount()
	}
}

func (status *simpleStatus) Error(message string) {
	fmt.Printf("ERROR %s\n", message)
	status.incErrorCount()
}
func (status *simpleStatus) LocationError(loc Location, message string) {
	fmt.Printf("ERROR %s:%d:%d: %s\n", loc.Filename, loc.Line+1, loc.Col, message)
	fmt.Printf("    %s\n", loc.Text)
	// TODO tabs?
	arrow := []rune{}
	for i := 0; i < loc.Col; i++ {
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

func MakeStatus() Status {
	return &simpleStatus{}
}
