package base

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

type CodeWriter struct {
	currentMargin     string
	marginStack       []string
	pendingEmptyLines int
	Out               io.Writer
}

func (w *CodeWriter) flushPending() {
	for w.pendingEmptyLines > 0 {
		w.writeLine("")
		w.pendingEmptyLines -= 1
	}
}

func (w *CodeWriter) writeLine(text string) {
	line := strings.TrimRightFunc(w.currentMargin+text, unicode.IsSpace)
	io.WriteString(w.Out, line)
	io.WriteString(w.Out, "\n")
}

func (w *CodeWriter) Line(text string) {
	w.flushPending()
	w.writeLine(text)
}

// The type signature of this function is stricter than what is done in fmt.
// This helps catch refactoring bugs where types being passed to Sprintf change.
// In general, code generators should not be using format strings for turning
// non-string data types into strings. Be explicit, don't implicitly rely on how
// Go formats data.
func (w *CodeWriter) Linef(format string, args ...string) {
	// "arg..." passes the slice straight through, so the type needs to match.
	// This is a gotcha for Python programers - the argument list is not expanded.
	rewrapped := make([]interface{}, len(args))
	for i, arg := range args {
		rewrapped[i] = arg
	}
	w.Line(fmt.Sprintf(format, rewrapped...))
}

func (w *CodeWriter) EmptyLines(count int) {
	if count > w.pendingEmptyLines {
		w.pendingEmptyLines = count
	}
}

func (w *CodeWriter) PushMargin(margin string) {
	w.marginStack = append(w.marginStack, w.currentMargin)
	w.currentMargin += margin
}

func (w *CodeWriter) PopMargin() {
	top := len(w.marginStack) - 1
	w.currentMargin = w.marginStack[top]
	w.marginStack = w.marginStack[:top]
}
