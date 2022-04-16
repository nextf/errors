package stack

import (
	"fmt"
	"path"
	"runtime"
)

type CallStack []uintptr
type Frame runtime.Frame

func NewCallStack(skip, maxDepth int) CallStack {
	rpc := make([]uintptr, maxDepth)
	baseSkip := 2
	n := runtime.Callers(skip+baseSkip, rpc)
	if n < 1 {
		return nil
	}
	return rpc[:n]
}

func (s CallStack) String() string {
	var buff []byte
	for _, frame := range s.CallersFrames() {
		buff = append(buff, []byte(fmt.Sprintf("%s\n", frame))...)
	}
	return string(buff)
}

func (s CallStack) CallersFrames() []Frame {
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

func (s *CallStack) Format(st fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case st.Flag('+'):
			fmt.Fprintf(st, "%s", s.String())
		}
	}
}

func (f Frame) String() string {
	_, file := path.Split(f.File)
	return fmt.Sprintf("%s(%s:%d)", f.Function, file, f.Line)
}
