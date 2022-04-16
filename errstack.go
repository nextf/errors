package errors

import (
	"fmt"

	"github.com/nextf/errors/stack"
)

const (
	maxDepth int = 32
)

type errorStack struct {
	stack stack.CallStack
	cause error
}

func (c *errorStack) CallersFrames() []stack.Frame {
	return c.stack.CallersFrames()
}

func (c *errorStack) Error() string {
	return c.cause.Error()
}

func (c *errorStack) Unwrap() error {
	return c.cause
}

func (c *errorStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		fmt.Fprintf(s, "@callstack\n%v", c.stack)
		if s.Flag('+') && c.cause != nil {
			fmt.Fprintf(s, "\nCaused by: %+v", c.cause)
		}
	}
}

func newErrorStack(skip int) error {
	return &errorStack{stack.NewCallStack(skip+1, maxDepth), nil}
}

func withErrorStack(err error, skip int) error {
	return &errorStack{stack.NewCallStack(skip+1, maxDepth), err}
}
