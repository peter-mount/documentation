package util

import (
	"fmt"
	"io"
)

// Writer assists in creating a LaTeX document
type Writer struct {
	w      io.WriteCloser // Underlying Writer
	parent *Writer        // Parent when nesting
	child  *Writer        // Child when nesting
	class  string         // document in begin{document} etc. Used with end
}

const (
	Command1 = "\\%s{%s}\n"
	Command2 = "\\%s{"
	Command3 = "}\n"
)

func NewWriter(w io.WriteCloser) *Writer {
	return &Writer{w: w}
}

// Write implements io.Writer
func (w *Writer) Write(p []byte) (n int, err error) {
	return w.w.Write(p)
}

// Close implements io.WriteCloser
func (w *Writer) Close() error {
	cw := w

	// Find youngest generation
	for cw.child != nil {
		cw = cw.child
	}

	// Now end until we get to the root
	for cw.parent != nil {
		cw = cw.End()
	}

	// Now close the writer
	return cw.w.Close()
}

func (w *Writer) WriteString(f string, a ...interface{}) *Writer {
	if len(a) == 0 {
		_, _ = w.Write([]byte(f))
	} else {
		_, _ = w.Write([]byte(fmt.Sprintf(f, a...)))
	}
	return w
}

func (w *Writer) DocumentClass() *Writer {
	return w.WriteString("\\documentclass[a4paper,10pt]{book}")
}

func (w *Writer) Begin(class string) *Writer {
	if w.child != nil {
		panic("begin when already with child")
	}
	w.WriteString(Command1, "begin", class)

	cw := &Writer{
		w:      w.w,
		parent: w,
		class:  class,
	}
	w.child = cw
	return cw
}

func (w *Writer) End() *Writer {
	if w.parent != nil {
		if w.parent.child == nil {
			panic("End from orphan")
		}
		_ = w.WriteString(Command1, "end", w.class)
		w.parent.child = nil
		return w.parent
	}
	return w
}

func (w *Writer) Comment(f string, a ...interface{}) *Writer {
	_ = w.WriteString("%% "+f+"\n", a...)
	return w
}

func (w *Writer) UsePackage(p ...string) *Writer {
	for _, s := range p {
		_ = w.WriteString(Command1, "usepackage", s)
	}
	return w
}
