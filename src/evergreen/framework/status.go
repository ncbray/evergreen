package framework

import (
	"fmt"
)

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

func MakeStatus(loc LocationProvider) Status {
	return &simpleStatus{loc: loc}
}
