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

func (w *CodeWriter) Linef(format string, args ...interface{}) {
	w.Line(fmt.Sprintf(format, args...))
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
