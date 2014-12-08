// Package compiler implements generic compiler functionality.
package compiler

import (
	"fmt"
	"strings"
	"time"
)

var Verbosity int = 0

type StatusReporter interface {
	GlobalError(message string)
	LocationError(loc int, message string)
	ShouldHalt() bool
}

type ParentStatus interface {
	StatusReporter
	ChildEnded()
}

type CompileStatus interface {
	ParentStatus
	Pass(name string) PassStatus
	ErrorCount() int
}

type PassStatus interface {
	ParentStatus
	Pass(name string) PassStatus
	Task(name string) TaskStatus
	Begin()
	End()
}

type TaskStatus interface {
	StatusReporter
	Begin()
	End()
}

type compileStatus struct {
	loc        LocationProvider
	errorCount int
	liveChild  bool
}

func (status *compileStatus) Pass(name string) PassStatus {
	if status.liveChild {
		panic(name)
	}
	status.liveChild = true
	return &passStatus{
		parent: status,
		path:   []string{name},
	}
}

func (status *compileStatus) incErrorCount() {
	status.errorCount += 1
}

func (status *compileStatus) GlobalError(message string) {
	fmt.Printf("ERROR: %s\n", message)
	status.incErrorCount()
}
func (status *compileStatus) LocationError(loc int, message string) {
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

func (status *compileStatus) ShouldHalt() bool {
	return status.errorCount > 0
}

func (status *compileStatus) ChildEnded() {
	if !status.liveChild {
		panic("???")
	}
	status.liveChild = false
}

func (status *compileStatus) ErrorCount() int {
	return status.errorCount
}

func MakeStatus(loc LocationProvider) CompileStatus {
	return &compileStatus{loc: loc}
}

type passStatus struct {
	parent    ParentStatus
	path      []string
	errored   bool
	live      bool
	start     time.Time
	liveChild bool
}

func (status *passStatus) Pass(name string) PassStatus {
	if !status.live || status.liveChild {
		panic(status.name())
	}
	status.liveChild = true
	return &passStatus{
		parent: status,
		path:   append(status.path, name),
	}
}

func (status *passStatus) Task(name string) TaskStatus {
	if !status.live || status.liveChild {
		panic(status.name())
	}
	status.liveChild = true
	return &taskStatus{
		parent: status,
		path:   append(status.path, name),
	}
}

func (status *passStatus) GlobalError(message string) {
	if !status.live {
		panic(status.name())
	}
	status.parent.GlobalError(message)
	status.errored = true
}
func (status *passStatus) LocationError(loc int, message string) {
	if !status.live {
		panic(status.name())
	}
	status.parent.LocationError(loc, message)
	status.errored = true
}

func (status *passStatus) ShouldHalt() bool {
	if !status.live {
		panic(status.name())
	}
	return status.errored
}

func (status *passStatus) ChildEnded() {
	if !status.live || !status.liveChild {
		panic(status.name())
	}
	status.liveChild = false
}

func (status *passStatus) name() string {
	return strings.Join(status.path, "|")
}

func (status *passStatus) Begin() {
	if status.live || status.liveChild {
		panic(status.name())
	}
	if Verbosity > 0 {
		fmt.Printf(">>> [%s]\n", status.name())
		status.start = time.Now()
	}
	status.live = true
}

func (status *passStatus) End() {
	status.parent.ChildEnded()
	if !status.live || status.liveChild {
		panic(status.name())
	}
	status.live = false
	delta := time.Since(status.start)
	if Verbosity > 0 {
		fmt.Printf("<<< [%s] %d us\n", status.name(), delta/time.Microsecond)
	}
}

type taskStatus struct {
	parent  ParentStatus
	path    []string
	errored bool
	live    bool
}

func (status *taskStatus) GlobalError(message string) {
	if !status.live {
		panic(status.name())
	}
	status.parent.GlobalError(message)
	status.errored = true
}

func (status *taskStatus) LocationError(loc int, message string) {
	if !status.live {
		panic(status.name())
	}
	status.parent.LocationError(loc, message)
	status.errored = true
}

func (status *taskStatus) ShouldHalt() bool {
	if !status.live {
		panic(status.name())
	}
	return status.errored
}

func (status *taskStatus) name() string {
	return strings.Join(status.path, "|")
}

func (status *taskStatus) Begin() {
	if status.live {
		panic(status.name())
	}
	status.live = true
}

func (status *taskStatus) End() {
	if !status.live {
		panic(status.name())
	}
	status.live = false
	status.parent.ChildEnded()
}
