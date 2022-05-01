package stack

import (
	"fmt"
	"io"
	"path"
	"runtime"
)

type CallStack []uintptr
type Frame runtime.Frame

func RecordCallStack(skip, maxDepth int) CallStack {
	rpc := make([]uintptr, maxDepth)
	baseSkip := 2
	n := runtime.Callers(skip+baseSkip, rpc)
	if n < 1 {
		return nil
	}
	return rpc[:n]
}

func (s CallStack) string() string {
	var buff []byte
	for _, frame := range s.StackTrace() {
		buff = append(buff, []byte(fmt.Sprintf("%s\n", frame.Describe()))...)
	}
	if len(buff) > 0 {
		buff = buff[:len(buff)-1]
	}
	return string(buff)
}

const tab string = "\x20\x20\x20\x20"

func (c CallStack) Format(s fmt.State, verb rune) {
	var indent string
	if s.Flag('+') {
		indent = tab
	}
	switch verb {
	case 'v':
		frames := c.StackTrace()
		framesSize := len(frames)
		maxDepth := framesSize
		if wid, ok := s.Width(); ok && wid < maxDepth {
			maxDepth = wid
		}
		i := 0
		for {
			fmt.Fprintf(s, "%s%s", indent, frames[i].Describe())
			i++
			if i < maxDepth {
				// Has more lines
				io.WriteString(s, "\n")
			} else {
				break
			}
		}
		if maxDepth < framesSize {
			// Collapse
			fmt.Fprintf(s, "\n%s...(more:%d)", indent, framesSize-maxDepth)
		}
	case 's':
		io.WriteString(s, c.string())
	case 'q':
		fmt.Fprintf(s, "%q", c.string())
	}
}

func (s CallStack) StackTrace() []Frame {
	frames := runtime.CallersFrames(s)
	var buff []Frame
	for {
		frame, more := frames.Next()
		buff = append(buff, Frame(frame))
		if !more {
			break
		}
	}
	return buff
}

func (f Frame) Describe() string {
	_, file := path.Split(f.File)
	return fmt.Sprintf("%s(%s:%d)", f.Function, file, f.Line)
}
