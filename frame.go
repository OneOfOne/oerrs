package oerrs

import (
	"fmt"
	"runtime"
	"strings"
)

type Framer interface {
	Frame() *Frame
}

// A Frame contains part of a call stack.
// Copied from xerrors
type Frame struct {
	runtime.Frame
	trim bool
}

// Caller returns a Frame that describes a frame on the caller's stack.
// The argument skip is the number of frames to skip over.
// Caller(0) returns the frame for the caller of Caller.
func Caller(skip int) *Frame {
	var callers [6]uintptr
	n := runtime.Callers(skip+1, callers[:])
	frs := runtime.CallersFrames(callers[:n])

	for fr, ok := frs.Next(); ok; fr, ok = frs.Next() {
		if !strings.HasPrefix(fr.Function, "runtime.") {
			return &Frame{fr, true}
		}
	}
	return &Frame{}
}

// location reports the file, line, and function of a frame.
//
// The returned function may be "" even if file and line are not.
func (fr *Frame) Location() (function, file string, line int) {
	if fr == nil {
		return
	}
	fn := fr.Function
	if fr.trim {
		if idx := strings.LastIndexByte(fn, '/'); idx > 0 {
			fn = fn[idx+1:]
		}
	}
	return fn, fr.File, fr.Line
}

func (f *Frame) String() string {
	function, file, line := f.Location()
	return fmt.Sprintf("%s:%d [%s]", file, line, function)
}

// Format prints the stack as error detail.
// It should be called from an error's Format implementation
// after printing any other error detail.
func (f *Frame) Format(p Printer) {
	if f == nil || !p.Detail() {
		return
	}

	function, file, line := f.Location()
	if function != "" {
		p.Printf("%s @ ", function)
	}
	if file != "" {
		p.Printf("%s:%d\n", file, line)
	} else {
		p.Print("n/a")
	}
}
